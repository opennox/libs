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

	"github.com/opennox/libs/noxnet/discover"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/noxnet/udpconn"
)

var (
	ErrPasswordRequired = errors.New("password required")
	ErrJoinFailed       = errors.New("join failed")
)

func NewClient(log *slog.Logger, conn udpconn.PacketConn) *Client {
	p := udpconn.NewPort(log, conn, netmsg.Options{IsClient: true})
	return NewClientWithPort(log, p)
}

func NewClientWithPort(log *slog.Logger, port *udpconn.Port) *Client {
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
	Port *udpconn.Port

	discover struct {
		sync.RWMutex
		byToken map[uint32]chan<- ServerInfoResp
	}

	join struct {
		sync.RWMutex
		res chan<- netmsg.Message
	}

	smu  sync.RWMutex
	port *udpconn.Conn
	srv  udpconn.Stream
	pid  uint32
	own  udpconn.Stream
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
	c.srv = udpconn.Stream{}
	c.pid = 0
	c.own = udpconn.Stream{}
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
		c.srv = c.port.Stream(udpconn.ServerStreamID)
		c.smu.Unlock()
	}
}

func (c *Client) handleMsg(s udpconn.Stream, m netmsg.Message, flags udpconn.PacketFlags) bool {
	c.smu.RLock()
	port := c.port
	c.smu.RUnlock()
	if port != nil && port == s.Conn() {
		switch s.SID() {
		case 0: // from server
			return c.handleServerMsg(m)
		}
		return false
	}
	if s.SID() != udpconn.ServerStreamID {
		return false
	}
	switch m := m.(type) {
	default:
		return false
	case *discover.MsgServerInfo:
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
			Addr: s.Conn().RemoteAddr(),
			Info: v,
		}:
		default:
		}
		return true
	}
}

func (c *Client) handleServerMsg(m netmsg.Message) bool {
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
	Info discover.MsgServerInfo
}

func (c *Client) Discover(ctx context.Context, port int, out chan<- ServerInfoResp) error {
	if port <= 0 {
		port = udpconn.DefaultPort
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
	if err := c.Port.BroadcastMsg(port, &discover.MsgDiscover{Token: token}); err != nil {
		return err
	}
	<-ctx.Done()
	err := ctx.Err()
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		err = nil
	}
	return err
}

func (c *Client) joinSrv(ctx context.Context, req netmsg.Message, out chan<- netmsg.Message, reliable bool) (func(), error) {
	c.smu.RLock()
	srv := c.srv
	c.smu.RUnlock()
	if !srv.Valid() {
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
		err = srv.SendReliable(ctx, req)
	} else {
		err = srv.SendUnreliable(req)
	}
	if err != nil {
		cancel()
		return nil, err
	}
	return cancel, nil
}

func (c *Client) joinOwn(ctx context.Context, req netmsg.Message, out chan<- netmsg.Message, reliable bool) (func(), error) {
	c.smu.RLock()
	own := c.own
	c.smu.RUnlock()
	if !own.Valid() {
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
		err = own.SendReliable(ctx, req)
	} else {
		err = own.SendUnreliable(req)
	}
	if err != nil {
		cancel()
		return nil, err
	}
	return cancel, nil
}

func (c *Client) TryJoin(ctx context.Context, addr netip.AddrPort, req MsgServerTryJoin) error {
	out := make(chan netmsg.Message, 1)
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
	out := make(chan netmsg.Message, 1)
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
	out := make(chan netmsg.Message, 1)
	cancel, err := c.joinSrv(ctx, &netmsg.Unknown{Op: netmsg.MSG_SERVER_CONNECT}, out, true)
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
			c.own = c.port.Stream(udpconn.SID(resp.ID))
			return nil
		}
	}
}

func (c *Client) clientAccept(ctx context.Context, req *MsgClientAccept) error {
	out := make(chan netmsg.Message, 1)
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
