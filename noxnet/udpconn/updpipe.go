package udpconn

import (
	"encoding/hex"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"slices"

	"github.com/opennox/libs/noxnet/netmsg"
)

type PipeConn struct {
	recv   chan []byte
	peer   *PipeConn
	closed chan struct{}

	Log   *slog.Logger
	Addr  netip.AddrPort
	Port  *Port
	Drop  func(data []byte) bool
	Debug bool
}

func (c *PipeConn) LocalAddr() net.Addr {
	return &net.UDPAddr{
		IP:   c.Addr.Addr().AsSlice(),
		Port: int(c.Addr.Port()),
	}
}

func (c *PipeConn) WriteToUDPAddrPort(data []byte, addr netip.AddrPort) (int, error) {
	if addr != c.peer.Addr && (addr.Addr() != broadcastIP4 || addr.Port() != c.peer.Addr.Port()) {
		return len(data), nil // ignore
	}
	drop := c.Drop
	if drop != nil && drop(data) {
		if c.Log != nil && c.Debug {
			c.Log.Info("DROP", "srv", c.Addr, "dst", addr, "data", hex.EncodeToString(data))
		}
		return len(data), nil
	}
	select {
	case c.peer.recv <- slices.Clone(data):
		if c.Log != nil && c.Debug {
			c.Log.Info("SEND", "srv", c.Addr, "dst", addr, "data", hex.EncodeToString(data))
		}
	case <-c.closed:
		return 0, io.ErrClosedPipe
	case <-c.peer.closed:
		return 0, io.ErrClosedPipe
	}
	return len(data), nil
}

func (c *PipeConn) ReadFromUDPAddrPort(data []byte) (int, netip.AddrPort, error) {
	select {
	case <-c.closed:
		return 0, netip.AddrPort{}, io.ErrClosedPipe
	case <-c.peer.closed:
		return 0, netip.AddrPort{}, io.ErrClosedPipe
	case m, ok := <-c.recv:
		if !ok {
			return 0, netip.AddrPort{}, io.ErrClosedPipe
		}
		if len(m) > len(data) {
			return 0, netip.AddrPort{}, io.ErrShortBuffer
		}
		n := copy(data, m)
		if c.Log != nil && c.Debug {
			c.Log.Info("RECV", "src", c.peer.Addr, "dst", c.Addr, "data", hex.EncodeToString(data[:n]))
		}
		return n, c.peer.Addr, nil
	}
}

func (c *PipeConn) Close() error {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	return nil
}

func NewPipe(log *slog.Logger, buf int) (srv *PipeConn, cli *PipeConn) {
	s2c := make(chan []byte, buf)
	c2s := make(chan []byte, buf)
	srv = &PipeConn{
		Log: log,
		Addr: netip.AddrPortFrom(
			netip.AddrFrom4([4]byte{1, 1, 1, 1}),
			10000,
		),
		recv:   c2s,
		closed: make(chan struct{}),
	}
	cli = &PipeConn{
		Log: log,
		Addr: netip.AddrPortFrom(
			netip.AddrFrom4([4]byte{2, 2, 2, 2}),
			20000,
		),
		recv:   s2c,
		closed: make(chan struct{}),
	}
	srv.peer, cli.peer = cli, srv
	srv.Port = NewPort(log, srv, netmsg.Options{IsClient: false})
	cli.Port = NewPort(log, cli, netmsg.Options{IsClient: true})
	return
}
