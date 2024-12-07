package noxnet

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/netip"
	"reflect"
	"sync"
)

var (
	ErrPasswordRequired = errors.New("password required")
	ErrJoinFailed       = errors.New("join failed")
)

func NewClient(log *slog.Logger, conn PacketConn) *Client {
	return NewClientWithPort(log, NewPort(log, conn, false))
}

func NewClientWithPort(log *slog.Logger, port *Port) *Client {
	c := &Client{
		log:  log,
		Port: port,
	}
	c.Port.OnMessage(c.handleMsg)
	c.Port.Start()
	return c
}

type Client struct {
	log  *slog.Logger
	Port *Port

	discover struct {
		sync.RWMutex
		byToken map[uint32]chan<- ServerInfoResp
	}

	join struct {
		sync.RWMutex
		res chan<- Message
	}

	smu  sync.RWMutex
	port *Conn
	srv  *Stream
	pid  uint32
	own  *Stream
}

func (c *Client) LocalAddr() netip.AddrPort {
	return c.Port.LocalAddr()
}

func (c *Client) Close() {
	c.Reset()
	c.Port.Close()
}

func (c *Client) Reset() {
	c.smu.Lock()
	defer c.smu.Unlock()
	c.port = nil
	c.srv = nil
	c.pid = 0
	c.own = nil
	c.Port.Reset()
}

func (c *Client) SetServerAddr(addr netip.AddrPort) {
	var cur netip.AddrPort
	if c.port != nil {
		cur = c.port.RemoteAddr()
	}
	if addr == cur {
		return
	}
	c.Reset()
	if addr.IsValid() {
		c.smu.Lock()
		c.port = c.Port.Conn(addr)
		c.srv = c.port.WithID(ServerStreamID)
		c.smu.Unlock()
	}
}

func (c *Client) handleMsg(conn *Conn, sid StreamID, m Message, flags PacketFlags) bool {
	c.smu.RLock()
	port := c.port
	c.smu.RUnlock()
	if port != nil && port == conn {
		switch sid {
		case 0: // from server
			return c.handleServerMsg(m)
		}
		return false
	}
	if sid != ServerStreamID {
		return false
	}
	switch m := m.(type) {
	default:
		return false
	case *MsgServerInfo:
		c.discover.RLock()
		ch := c.discover.byToken[m.Token]
		c.discover.RUnlock()
		if ch == nil {
			return true // ignore
		}
		v := *m
		v.Token = 0
		select {
		case ch <- ServerInfoResp{
			Addr: conn.RemoteAddr(),
			Info: v,
		}:
		default:
		}
		return true
	}
}

func (c *Client) handleServerMsg(m Message) bool {
	switch m := m.(type) {
	case *MsgJoinOK, ErrorMsg:
		c.join.RLock()
		res := c.join.res
		c.join.RUnlock()
		if res != nil {
			select {
			case res <- m:
			default:
			}
			return true
		}
		return false
	default:
		c.log.Warn("unhandled server message", "type", reflect.TypeOf(m).String(), "msg", m)
		return false
	}
}

type ServerInfoResp struct {
	Addr netip.AddrPort
	Info MsgServerInfo
}

func (c *Client) Discover(ctx context.Context, port int, out chan<- ServerInfoResp) error {
	if port <= 0 {
		port = DefaultPort
	}
	token := rand.Uint32()
	c.discover.Lock()
	if c.discover.byToken == nil {
		c.discover.byToken = make(map[uint32]chan<- ServerInfoResp)
	}
	c.discover.byToken[token] = out
	c.discover.Unlock()
	defer func() {
		c.discover.Lock()
		delete(c.discover.byToken, token)
		c.discover.Unlock()
	}()
	if err := c.Port.BroadcastMsg(port, &MsgDiscover{Token: token}); err != nil {
		return err
	}
	<-ctx.Done()
	err := ctx.Err()
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		err = nil
	}
	return err
}

func (c *Client) joinSrv(ctx context.Context, req Message, out chan<- Message, reliable bool) (func(), error) {
	c.smu.RLock()
	srv := c.srv
	c.smu.RUnlock()
	if srv == nil {
		return nil, errors.New("server address must be set")
	}
	c.join.Lock()
	if c.join.res != nil {
		c.join.Unlock()
		return nil, errors.New("already joining")
	}
	c.join.res = out
	c.join.Unlock()
	cancel := func() {
		c.join.Lock()
		c.join.res = nil
		c.join.Unlock()
	}
	var err error
	if reliable {
		err = srv.SendReliableMsg(ctx, req)
	} else {
		err = srv.SendUnreliableMsg(req)
	}
	if err != nil {
		cancel()
		return nil, err
	}
	return cancel, nil
}

func (c *Client) joinOwn(ctx context.Context, req Message, out chan<- Message, reliable bool) (func(), error) {
	c.smu.RLock()
	own := c.own
	c.smu.RUnlock()
	if own == nil {
		return nil, errors.New("not connected")
	}
	c.join.Lock()
	if c.join.res != nil {
		c.join.Unlock()
		return nil, errors.New("already joining")
	}
	c.join.res = out
	c.join.Unlock()
	cancel := func() {
		c.join.Lock()
		c.join.res = nil
		c.join.Unlock()
	}
	var err error
	if reliable {
		err = own.SendReliableMsg(ctx, req)
	} else {
		err = own.SendUnreliableMsg(req)
	}
	if err != nil {
		cancel()
		return nil, err
	}
	return cancel, nil
}

func (c *Client) TryJoin(ctx context.Context, addr netip.AddrPort, req MsgServerTryJoin) error {
	out := make(chan Message, 1)
	c.SetServerAddr(addr)
	cancel, err := c.joinSrv(ctx, &req, out, false)
	if err != nil {
		return err
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-out:
		switch resp := resp.(type) {
		case *MsgJoinOK:
			return nil
		case ErrorMsg:
			return resp.Error()
		default:
			return fmt.Errorf("unexpected response: %v", resp.NetOp())
		}
	}
}

func (c *Client) TryPassword(ctx context.Context, pass string) error {
	out := make(chan Message, 1)
	cancel, err := c.joinSrv(ctx, &MsgServerPass{Pass: pass}, out, false)
	if err != nil {
		return err
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-out:
		switch resp := resp.(type) {
		case *MsgJoinOK:
			return nil
		case ErrorMsg:
			return resp.Error()
		default:
			return fmt.Errorf("unexpected response: %v", resp.NetOp())
		}
	}
}

func (c *Client) connect(ctx context.Context, addr netip.AddrPort) error {
	c.SetServerAddr(addr)
	out := make(chan Message, 1)
	cancel, err := c.joinSrv(ctx, &MsgUnknown{Op: MSG_SERVER_CONNECT}, out, true)
	if err != nil {
		return err
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-out:
		if err := c.port.Ack(); err != nil {
			return err
		}
		switch resp := resp.(type) {
		default:
			return fmt.Errorf("unexpected response: %v", resp.NetOp())
		case *MsgServerAccept:
			c.smu.Lock()
			defer c.smu.Unlock()
			c.port.Encrypt(resp.XorKey)
			c.pid = resp.ID
			c.own = c.port.WithID(StreamID(resp.ID))
			return nil
		}
	}
}

func (c *Client) clientAccept(ctx context.Context, req *MsgClientAccept) error {
	out := make(chan Message, 1)
	cancel, err := c.joinOwn(ctx, req, out, true)
	if err != nil {
		return err
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-out:
		if err := c.port.Ack(); err != nil {
			return err
		}
		switch resp := resp.(type) {
		default:
			return fmt.Errorf("unexpected response: %v", resp.NetOp())
		}
	}
}

func (c *Client) Connect(ctx context.Context, addr netip.AddrPort, req *MsgClientAccept) error {
	if err := c.connect(ctx, addr); err != nil {
		return err
	}
	if err := c.clientAccept(ctx, req); err != nil {
		return err
	}
	return nil
}
