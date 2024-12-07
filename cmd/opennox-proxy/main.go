package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/opennox/libs/log"
	"github.com/opennox/libs/noxnet"
)

//go:generate d2 diagram.d2 diagram.svg
//go:generate d2 diagram.d2 diagram.png

var (
	fServer = flag.String("server", "127.0.0.1:18590", "server address to proxy requests to")
	fHost   = flag.String("host", "0.0.0.0:18600", "address to host proxy on")
	fFile   = flag.String("file", "", "file name to dump messages to")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	srv, err := netip.ParseAddrPort(*fServer)
	if err != nil {
		return err
	}
	p := NewProxy(srv)
	defer p.Close()
	log.Printf("serving proxy %v -> %v", *fHost, srv)
	return p.ListenAndServe(*fHost)
}

func NewProxy(srv netip.AddrPort) *Proxy {
	p := &Proxy{
		realSrv: srv,
		clients: make(map[netip.AddrPort]*clientPort),
	}
	return p
}

type Proxy struct {
	realSrv  netip.AddrPort
	clientID uint32 // atomic

	emu   sync.Mutex
	efile *os.File
	enc   *json.Encoder

	wmu sync.Mutex
	lis *net.UDPConn

	cmu     sync.RWMutex
	clients map[netip.AddrPort]*clientPort
}

func (p *Proxy) Close() error {
	p.emu.Lock()
	if p.efile != nil {
		p.efile.Close()
	}
	p.emu.Unlock()
	p.wmu.Lock()
	if p.lis != nil {
		p.lis.Close()
	}
	p.wmu.Unlock()
	return nil
}

func (p *Proxy) ListenAndServe(addr string) error {
	ip, err := netip.ParseAddrPort(addr)
	if err != nil {
		return err
	}
	lis, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   ip.Addr().AsSlice(),
		Port: int(ip.Port()),
	})
	if err != nil {
		return err
	}
	defer lis.Close()
	return p.Serve(lis)
}

func (p *Proxy) Serve(lis *net.UDPConn) error {
	p.lis = lis
	var buf [4096]byte
	for {
		n, addr, err := lis.ReadFromUDPAddrPort(buf[:])
		if err != nil {
			return err
		}
		data := buf[:n]
		p.sendAsClient(addr, data)
	}
}

func (p *Proxy) getClient(addr netip.AddrPort) (*clientPort, error) {
	p.cmu.RLock()
	c := p.clients[addr]
	p.cmu.RUnlock()
	if c != nil {
		return c, nil
	}

	p.cmu.Lock()
	c = p.clients[addr]
	if c != nil {
		p.cmu.Unlock()
		return c, nil
	}
	c = p.newClient(addr)
	p.clients[addr] = c
	p.cmu.Unlock()
	if err := c.listen(addr.Addr()); err != nil {
		p.cmu.Lock()
		delete(p.clients, addr)
		p.cmu.Unlock()
		return nil, err
	}
	log.Printf("NEW %d: %v (real) <=> %v (proxy)", c.id, addr, c.lis.LocalAddr())
	go c.serve()
	return c, nil
}

// sendAsClient sends data from the client using unique client server port to the real server.
func (p *Proxy) sendAsClient(realCli netip.AddrPort, data []byte) {
	c, err := p.getClient(realCli)
	if err != nil {
		log.Printf("cannot host client %v: %v", realCli, err)
		return
	}
	p.recordPacket(c.id, 0, data)
	log.Printf("CLI%d(%v) -> SP(%v): [%d]: %x", c.id, realCli, p.lis.LocalAddr(), len(data), data)
	err = c.SendToServer(data)
	if err != nil {
		log.Printf("cannot send client %v packet: %v", realCli, err)
		return
	}
}

// sendToClient sends data from the proxy server port to the client.
func (p *Proxy) sendToClient(id uint32, addr netip.AddrPort, data []byte) error {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	p.recordPacket(0, id, data)
	log.Printf("SP(%v) -> CLI%d(%v): [%d]: %x", p.lis.LocalAddr(), id, addr, len(data), data)
	_, err := p.lis.WriteTo(data, net.UDPAddrFromAddrPort(addr))
	return err
}

func (p *Proxy) newClient(addr netip.AddrPort) *clientPort {
	id := atomic.AddUint32(&p.clientID, 1)
	return &clientPort{
		id:      id,
		p:       p,
		realCli: addr,
	}
}

type clientPort struct {
	id      uint32 // our own id for debugging
	p       *Proxy
	realCli netip.AddrPort
	xor     uint32 // atomic

	wmu sync.Mutex
	lis *net.UDPConn
}

func (c *clientPort) listen(addr netip.Addr) error {
	lis, err := net.ListenUDP("udp4", &net.UDPAddr{IP: addr.AsSlice()})
	if err != nil {
		return err
	}
	c.lis = lis
	return nil
}

// serve accepts packets from the real server and redirects it to the proxied client.
func (c *clientPort) serve() {
	var buf [4096]byte
	for {
		n, addr, err := c.lis.ReadFromUDPAddrPort(buf[:])
		if err != nil {
			log.Printf("client %v listener: %v", c.realCli, err)
			return
		}
		data := buf[:n]
		if addr != c.p.realSrv {
			log.Printf("???(%v) -> CP%d(%v): [%d]: %x", addr, c.id, c.lis.LocalAddr(), len(data), data)
			continue
		}
		log.Printf("SRV(%v) -> CP%d(%v): [%d]: %x", c.p.realSrv, c.id, c.lis.LocalAddr(), len(data), data)
		if xor := byte(atomic.LoadUint32(&c.xor)); xor != 0 {
			xorData(xor, data)
		}
		data = c.interceptServer(data)
		if len(data) == 0 {
			continue
		}
		err = c.p.sendToClient(c.id, c.realCli, data)
		if err != nil {
			log.Printf("client %v send: %v", c.realCli, err)
		}
	}
}

func modifyMessage[T noxnet.Message](data []byte, fnc func(p T)) []byte {
	var zero T
	msg := reflect.New(reflect.TypeOf(zero).Elem()).Interface().(T)
	_, err := msg.Decode(data[3:])
	if err != nil {
		log.Printf("cannot decode %v: %v", msg.NetOp(), err)
		return data
	}
	fnc(msg)
	buf := make([]byte, 3+msg.EncodeSize())
	copy(buf, data[:3])
	_, err = msg.Encode(buf[3:])
	if err != nil {
		log.Printf("cannot encode %v: %v", msg.NetOp(), err)
		return data
	}
	return buf
}

func (c *clientPort) interceptServer(data []byte) []byte {
	if len(data) < 3 {
		return data
	}
	if data[0] == 0 && data[1] == 0 {
		switch noxnet.Op(data[2]) {
		case noxnet.MSG_SERVER_INFO:
			return modifyMessage(data, func(p *noxnet.MsgServerInfo) {
				p.ServerName = "Proxy: " + p.ServerName
			})
		}
	} else if data[0] == 0x80 && data[1] == 0 {
		switch noxnet.Op(data[2]) {
		case noxnet.MSG_SERVER_ACCEPT:
			return modifyMessage(data, func(p *noxnet.MsgServerAccept) {
				atomic.StoreUint32(&c.xor, uint32(p.XorKey))
				p.XorKey = 0
			})
		case noxnet.MSG_ACCEPTED:
			var accept noxnet.MsgAccept
			left := data[2:]
			n, err := accept.Decode(left[1:])
			if err != nil {
				return data
			}
			left = left[1+n:]

			if len(left) == 0 || noxnet.Op(left[0]) != noxnet.MSG_SERVER_ACCEPT {
				return data
			}
			var saccept noxnet.MsgServerAccept
			n, err = saccept.Decode(left[1:])
			if err != nil {
				return data
			}
			atomic.StoreUint32(&c.xor, uint32(saccept.XorKey))
			saccept.XorKey = 0

			out := append([]byte{}, data[0], data[1])
			out, err = noxnet.AppendPacket(out, &accept)
			if err != nil {
				return data
			}
			out, err = noxnet.AppendPacket(out, &saccept)
			if err != nil {
				return data
			}
			return out
		}
	}
	return data
}

// SendToServer data to the real server using client's unique proxy port.
func (c *clientPort) SendToServer(data []byte) error {
	if xor := byte(atomic.LoadUint32(&c.xor)); xor != 0 {
		xorData(xor, data)
	}
	log.Printf("CP%d(%v) -> SRV(%v): [%d]: %x", c.id, c.lis.LocalAddr(), c.p.realSrv, len(data), data)
	c.wmu.Lock()
	defer c.wmu.Unlock()
	_, err := c.lis.WriteTo(data, net.UDPAddrFromAddrPort(c.p.realSrv))
	return err
}

func xorData(key byte, p []byte) {
	for i := range p {
		p[i] ^= key
	}
}
