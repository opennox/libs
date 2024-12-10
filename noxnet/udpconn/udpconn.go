package udpconn

import (
	"context"
	"encoding/hex"
	"log/slog"
	"net"
	"net/netip"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opennox/libs/noxnet/netmsg"
)

const (
	DefaultPort = 18590
)

const (
	MaxStreams     = 128
	ServerStreamID = StreamID(0)
	MaxStreamID    = StreamID(MaxStreams - 1)
	maskID         = byte(MaxStreams - 1) // 0x7F
	reliableFlag   = byte(MaxStreams)     // 0x80
)

const (
	resendTick     = 20 * time.Millisecond
	resendInterval = time.Second
	resendRetries  = 5
	defaultTimeout = resendInterval*resendRetries + resendTick
	maxAckDelay    = 100 * time.Millisecond
	maxAckMsgs     = 50
)

var (
	broadcastIP4 = netip.AddrFrom4([4]byte{255, 255, 255, 255})
)

type PacketConn interface {
	LocalAddr() net.Addr
	WriteToUDPAddrPort(b []byte, addr netip.AddrPort) (int, error)
	ReadFromUDPAddrPort(b []byte) (int, netip.AddrPort, error)
	Close() error
}

type onMessageFuncs struct {
	mu    sync.RWMutex
	funcs []OnMessageFunc
}

func (f *onMessageFuncs) Add(fnc OnMessageFunc) {
	if fnc == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.funcs = append(f.funcs, fnc)
}

func (f *onMessageFuncs) Call(conn *Conn, sid StreamID, m netmsg.Message, flags PacketFlags) bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	list := f.funcs
	f.mu.RUnlock()
	for _, fnc := range list {
		if fnc(conn, sid, m, flags) {
			return true
		}
	}
	return false
}

type StreamID byte
type PacketFlags byte

func (f PacketFlags) Has(f2 PacketFlags) bool {
	return f&f2 != 0
}

const (
	PacketReliable = PacketFlags(1 << iota)
)

type OnMessageFunc func(conn *Conn, sid StreamID, m netmsg.Message, flags PacketFlags) bool

func NewPort(log *slog.Logger, conn PacketConn, isServer bool) *Port {
	if log == nil {
		log = slog.Default()
	}
	return &Port{
		log:      log,
		conn:     conn,
		isServer: isServer,
		debug:    log.Enabled(context.Background(), slog.LevelDebug),
		byAddr:   make(map[netip.AddrPort]*Conn),
		closed:   make(chan struct{}),
	}
}

type Port struct {
	log *slog.Logger

	wmu  sync.Mutex
	wbuf []byte
	conn PacketConn

	OnConn    func(c *Conn) bool
	onMessage onMessageFuncs

	hmu    sync.RWMutex
	byAddr map[netip.AddrPort]*Conn

	closed   chan struct{}
	isServer bool
	debug    bool
}

func (p *Port) Close() {
	select {
	case <-p.closed:
	default:
		close(p.closed)
		_ = p.conn.Close()
	}
}

func (p *Port) LocalAddr() netip.AddrPort {
	addr := p.conn.LocalAddr().(*net.UDPAddr)
	ip, _ := netip.AddrFromSlice(addr.IP)
	return netip.AddrPortFrom(ip, uint16(addr.Port))
}

func (p *Port) OnMessage(fnc OnMessageFunc) {
	p.onMessage.Add(fnc)
}

func (p *Port) Reset() {
	p.hmu.Lock()
	defer p.hmu.Unlock()
	for _, h := range p.byAddr {
		h.Reset()
	}
	p.byAddr = make(map[netip.AddrPort]*Conn)
}

func (p *Port) getConn(addr netip.AddrPort) *Conn {
	p.hmu.RLock()
	defer p.hmu.RUnlock()
	return p.byAddr[addr]
}

func (p *Port) Conn(addr netip.AddrPort) *Conn {
	p.hmu.RLock()
	h := p.byAddr[addr]
	p.hmu.RUnlock()
	if h != nil {
		return h
	}
	p.hmu.Lock()
	defer p.hmu.Unlock()
	h = p.byAddr[addr]
	if h != nil {
		return h
	}
	h = &Conn{
		p:    p,
		addr: addr,
		log:  p.log.With("remote", addr),
	}
	if p.OnConn != nil && !p.OnConn(h) {
		return nil
	}
	p.byAddr[addr] = h
	return h
}

func (p *Port) writeRaw(addr netip.AddrPort, b1, b2 byte, data []byte, xor byte) error {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	p.wbuf = p.wbuf[:0]
	p.wbuf = append(p.wbuf, b1, b2)
	p.wbuf = append(p.wbuf, data...)
	if xor != 0 {
		xorBuf(xor, p.wbuf)
	}
	_, err := p.conn.WriteToUDPAddrPort(p.wbuf, addr)
	return err
}

func (p *Port) WriteMsg(addr netip.AddrPort, b1, b2 byte, m netmsg.Message, xor byte) error {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	var err error
	p.wbuf = p.wbuf[:0]
	p.wbuf = append(p.wbuf, b1, b2)
	p.wbuf, err = netmsg.Append(p.wbuf, m)
	if err != nil {
		return err
	}
	if xor != 0 {
		xorBuf(xor, p.wbuf)
	}
	_, err = p.conn.WriteToUDPAddrPort(p.wbuf, addr)
	return err
}

func (p *Port) BroadcastMsg(port int, m netmsg.Message) error {
	if port <= 0 {
		port = DefaultPort
	}
	addr := netip.AddrPortFrom(broadcastIP4, uint16(port))
	return p.WriteMsg(addr, 0, 0, m, 0)
}

func (p *Port) Start() {
	go p.readLoop()
	go p.resendLoop()
}

func (p *Port) readLoop() {
	var buf [4096]byte
	for {
		n, addr, err := p.conn.ReadFromUDPAddrPort(buf[:])
		if err != nil {
			select {
			default:
				p.log.Error("cannot read packet", "err", err)
			case <-p.closed:
			}
			return
		}
		data := buf[:n]
		if len(data) < 2 {
			continue
		}
		h := p.Conn(addr)
		if h == nil {
			continue // ignore
		}
		h.handlePacket(data)
	}
}

func (p *Port) resendLoop() {
	ticker := time.NewTicker(resendTick)
	defer ticker.Stop()
	for {
		select {
		case <-p.closed:
			return
		case <-ticker.C:
		}
		p.resendAll()
	}
}

func (p *Port) resendAll() {
	p.hmu.Lock()
	defer p.hmu.Unlock()
	for _, h := range p.byAddr {
		_ = h.SendQueue()
	}
}

func xorBuf(key byte, p []byte) {
	for i := range p {
		p[i] ^= key
	}
}

func seqBefore(v, cur byte) bool {
	return v <= cur || (v >= 0xff-maxAckMsgs && cur-v < maxAckMsgs)
}

type PacketID uintptr

type packet struct {
	pid      PacketID
	sid      StreamID
	seq      byte
	xor      byte
	lastSend time.Time
	deadline time.Time
	data     []byte
	done     func()
	timeout  func()
}

type Conn struct {
	p    *Port
	log  *slog.Logger
	addr netip.AddrPort

	packetID atomic.Uintptr
	mu       sync.RWMutex
	xor      byte
	syn      byte
	ack      byte
	needAck  int
	nextPing time.Time
	queue    []*packet

	onMessage onMessageFuncs

	smu  sync.RWMutex
	byID [MaxStreams]*Stream
}

func (p *Conn) Port() *Port {
	return p.p
}

func (p *Conn) RemoteAddr() netip.AddrPort {
	return p.addr
}

func (p *Conn) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.xor = 0
	p.syn = 0
	p.ack = 0
	p.nextPing = time.Time{}
	p.needAck = 0
	p.queue = nil
}

func (p *Conn) Encrypt(key byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.xor = key
}

func (p *Conn) OnMessage(fnc OnMessageFunc) {
	p.onMessage.Add(fnc)
}

func (p *Conn) handlePacket(data []byte) {
	var doneFuncs []func()
	p.mu.Lock()
	if xor := p.xor; xor != 0 {
		xorBuf(xor, data)
	}
	b1, b2 := data[0], data[1]
	data = data[2:]
	reliable, sid, seq := b1&reliableFlag != 0, StreamID(b1&maskID), b2
	if p.p.debug {
		sdata := hex.EncodeToString(data)
		if reliable {
			p.log.Debug("RECV", "syn", seq, "sid", sid, "data", sdata)
		} else {
			p.log.Debug("RECV", "ack", seq, "sid", sid, "data", sdata)
		}
	}
	if reliable {
		// New reliable message that we should ACK in the future.
		exp := p.ack
		if seq != exp {
			p.mu.Unlock()
			return // Ignore out of order packets.
		}
		p.ack = seq + 1
		p.needAck++
		if p.needAck-1 >= maxAckMsgs {
			_ = p.sendAckPing()
		} else {
			p.nextPing = time.Now().Add(maxAckDelay)
		}
	} else {
		// Unreliable message with ACK for our reliable messages.
		p.queue = slices.DeleteFunc(p.queue, func(m *packet) bool {
			del := seqBefore(m.seq, seq)
			if del && m.done != nil {
				doneFuncs = append(doneFuncs, m.done)
			}
			return del
		})
	}
	p.mu.Unlock()
	var onMsgID *onMessageFuncs
	if s := p.getWithID(sid); s != nil {
		onMsgID = &s.onMessage
	}
	onMsg := &p.onMessage
	onMsgGlobal := &p.p.onMessage
	for _, done := range doneFuncs {
		done()
	}

	var flags PacketFlags
	if reliable {
		flags |= PacketReliable
	}
	for len(data) > 0 {
		m, n, err := netmsg.DecodeAny(data)
		if err != nil {
			op := data[0]
			p.log.Error("Failed to decode packet", "op", op, "err", err)
			break
		}
		data = data[n:]
		if p.p.debug {
			p.log.Debug("RECV", "type", reflect.TypeOf(m).String(), "msg", m)
		}
		if onMsgID.Call(p, sid, m, flags) {
			continue
		}
		if onMsg.Call(p, sid, m, flags) {
			continue
		}
		if onMsgGlobal.Call(p, sid, m, flags) {
			continue
		}
	}
}

func (p *Conn) ViewQueue(fnc func(sid StreamID, data []byte)) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, m := range p.queue {
		fnc(m.sid, m.data)
	}
}

func (p *Conn) QueuedFor(sid StreamID) int {
	n := 0
	p.ViewQueue(func(sid2 StreamID, data []byte) {
		if sid == sid2 {
			n++
		}
	})
	return n
}

func (p *Conn) UpdateQueue(fnc func(pid PacketID, sid StreamID, data []byte) bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.queue = slices.DeleteFunc(p.queue, func(m *packet) bool {
		return !fnc(m.pid, m.sid, m.data)
	})
}

func (p *Conn) ResetFor(sid StreamID) {
	p.UpdateQueue(func(_ PacketID, sid2 StreamID, _ []byte) bool {
		return sid != sid2 // remove sid == sid2
	})
}

func (p *Conn) sendUnreliable(sid StreamID, data []byte) error {
	seq := p.ack
	p.nextPing = time.Time{}
	if p.p.debug {
		p.log.Debug("SEND", "ack", seq, "sid", sid, "data", hex.EncodeToString(data))
	}
	return p.p.writeRaw(p.addr, byte(sid)&maskID, seq, data, p.xor)
}

func (p *Conn) sendAckPing() error {
	p.needAck = 0
	p.nextPing = time.Time{}
	return p.sendUnreliable(0, nil)
}

func (p *Conn) Ack() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sendAckPing()
}

func (p *Conn) SendUnreliable(sid StreamID, data []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sendUnreliable(sid, data)
}

func (p *Conn) SendUnreliableMsg(sid StreamID, m netmsg.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	seq := p.ack
	p.nextPing = time.Time{}
	if p.p.debug {
		p.log.Debug("SEND", "ack", seq, "sid", sid, "type", reflect.TypeOf(m).String(), "msg", m)
	}
	return p.p.WriteMsg(p.addr, byte(sid)&maskID, seq, m, p.xor)
}

func (p *Conn) QueueReliable(ctx context.Context, sid StreamID, data []byte, done, timeout func()) PacketID {
	data = slices.Clone(data)
	now := time.Now()
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = now.Add(defaultTimeout)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	seq := p.syn
	p.syn++
	pid := PacketID(p.packetID.Add(1))
	p.queue = append(p.queue, &packet{
		pid:      pid,
		sid:      sid,
		seq:      seq,
		xor:      p.xor,
		lastSend: time.Time{}, // send in the next tick
		deadline: deadline,
		data:     data,
		done:     done,
		timeout:  timeout,
	})
	return pid
}

func (p *Conn) CancelReliable(pid PacketID) {
	p.UpdateQueue(func(pid2 PacketID, _ StreamID, _ []byte) bool {
		return pid != pid2 // remove pid == pid2
	})
}

func (p *Conn) QueueReliableMsg(ctx context.Context, sid StreamID, arr []netmsg.Message, done, timeout func()) (PacketID, error) {
	var (
		data []byte
		err  error
	)
	for _, m := range arr {
		data, err = netmsg.Append(data, m)
		if err != nil {
			return 0, err
		}
	}
	pid := p.QueueReliable(ctx, sid, data, done, timeout)
	return pid, nil
}

func (p *Conn) SendReliable(ctx context.Context, sid StreamID, data []byte) error {
	var cancel func()
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	acked := make(chan struct{})
	pid := p.QueueReliable(ctx, sid, data, func() {
		close(acked)
	}, cancel)
	if err := p.sendQueue(func(p *packet) bool {
		return p.pid == pid
	}); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-acked:
		return nil
	}
}

func (p *Conn) SendReliableMsg(ctx context.Context, sid StreamID, m netmsg.Message) error {
	data, err := netmsg.Append(nil, m)
	if err != nil {
		return err
	}
	return p.SendReliable(ctx, sid, data)
}

func (p *Conn) sendQueue(filter func(p *packet) bool) error {
	var (
		doneFuncs []func()
		lastErr   error
	)
	now := time.Now()
	p.mu.Lock()
	p.queue = slices.DeleteFunc(p.queue, func(m *packet) bool {
		if filter != nil && !filter(m) {
			return false // keep
		}
		del := m.deadline.Before(now)
		if del && m.timeout != nil {
			doneFuncs = append(doneFuncs, m.timeout)
		}
		return del
	})
	for i := range p.queue {
		m := p.queue[i]
		if filter != nil && !filter(m) {
			continue
		}
		if m.lastSend.Add(resendInterval).Before(now) {
			m.lastSend = now
			if p.p.debug {
				p.log.Debug("SEND", "syn", m.seq, "sid", m.sid, "data", hex.EncodeToString(m.data))
			}
			if err := p.p.writeRaw(p.addr, byte(m.sid)|reliableFlag, m.seq, m.data, m.xor); err != nil {
				lastErr = err
			}
		}
	}
	if !p.nextPing.IsZero() && p.nextPing.Before(now) {
		if err := p.sendAckPing(); err != nil {
			lastErr = err
		}
	}
	p.mu.Unlock()
	for _, done := range doneFuncs {
		done()
	}
	return lastErr
}

func (p *Conn) SendQueue() error {
	return p.sendQueue(nil)
}

func (p *Conn) getWithID(sid StreamID) *Stream {
	p.smu.RLock()
	defer p.smu.RUnlock()
	return p.byID[sid]
}

func (p *Conn) WithID(sid StreamID) *Stream {
	p.smu.RLock()
	s := p.byID[sid]
	p.smu.RUnlock()
	if s != nil {
		return s
	}
	p.smu.Lock()
	defer p.smu.Unlock()
	s = p.byID[sid]
	if s != nil {
		return s
	}
	s = &Stream{p: p, sid: sid}
	p.byID[sid] = s
	return s
}

type Stream struct {
	p   *Conn
	sid StreamID

	onMessage onMessageFuncs
}

func (p *Stream) Conn() *Conn {
	return p.p
}

func (p *Stream) SID() StreamID {
	return p.sid
}

func (p *Stream) Addr() netip.AddrPort {
	return p.p.RemoteAddr()
}

func (p *Stream) Reset() {
	p.p.ResetFor(p.sid)
}

func (p *Stream) SendQueue() error {
	return p.p.sendQueue(func(m *packet) bool {
		return m.sid == p.sid
	})
}

func (p *Stream) OnMessage(fnc OnMessageFunc) {
	p.onMessage.Add(fnc)
}

func (p *Stream) SendUnreliable(data []byte) error {
	return p.p.SendUnreliable(p.sid, data)
}

func (p *Stream) SendUnreliableMsg(m netmsg.Message) error {
	return p.p.SendUnreliableMsg(p.sid, m)
}

func (p *Stream) QueueReliable(ctx context.Context, data []byte, done, timeout func()) PacketID {
	return p.p.QueueReliable(ctx, p.sid, data, done, timeout)
}

func (p *Stream) CancelReliable(id PacketID) {
	p.p.CancelReliable(id)
}

func (p *Stream) QueueReliableMsg(ctx context.Context, arr []netmsg.Message, done, timeout func()) (PacketID, error) {
	return p.p.QueueReliableMsg(ctx, p.sid, arr, done, timeout)
}

func (p *Stream) SendReliable(ctx context.Context, data []byte) error {
	return p.p.SendReliable(ctx, p.sid, data)
}

func (p *Stream) SendReliableMsg(ctx context.Context, m netmsg.Message) error {
	return p.p.SendReliableMsg(ctx, p.sid, m)
}
