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
	ServerStreamID = SID(0)
	MaxStreamID    = SID(MaxStreams - 1)
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

func (f *onMessageFuncs) Call(s Stream, m netmsg.Message, flags PacketFlags) bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	list := f.funcs
	f.mu.RUnlock()
	for _, fnc := range list {
		if fnc(s, m, flags) {
			return true
		}
	}
	return false
}

type Header struct {
	SID   SID
	Seq   Seq
	Flags PacketFlags
}

func (h *Header) Decode(b [2]byte) {
	h.SID = SID(b[0] & maskID)
	h.Flags = Unreliable
	if b[0]&reliableFlag != 0 {
		h.Flags |= Reliable
	}
	h.Seq = Seq(b[1])
}

func (h Header) Encode() [2]byte {
	b1 := byte(h.SID) & maskID
	if h.Flags.Has(Reliable) {
		b1 |= reliableFlag
	}
	b2 := byte(h.Seq)
	return [2]byte{b1, b2}
}

// SID is a stream ID.
type SID byte
type PacketFlags byte

func (f PacketFlags) Has(f2 PacketFlags) bool {
	return f&f2 != 0
}

const (
	Unreliable = PacketFlags(0)
	Reliable   = PacketFlags(1 << iota)
)

type OnMessageFunc func(s Stream, m netmsg.Message, flags PacketFlags) bool

func NewPort(log *slog.Logger, conn PacketConn, opts netmsg.Options) *Port {
	if log == nil {
		log = slog.Default()
	}
	return &Port{
		log:    log,
		opts:   opts,
		conn:   conn,
		debug:  log.Enabled(context.Background(), slog.LevelDebug),
		byAddr: make(map[netip.AddrPort]*Conn),
		closed: make(chan struct{}),
	}
}

type Port struct {
	log  *slog.Logger
	opts netmsg.Options

	wmu  sync.Mutex
	wbuf []byte
	conn PacketConn

	OnConn    func(c *Conn) bool
	onMessage onMessageFuncs

	hmu    sync.RWMutex
	byAddr map[netip.AddrPort]*Conn

	closed chan struct{}
	debug  bool
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
	h.enc.Options = p.opts
	if p.OnConn != nil && !p.OnConn(h) {
		return nil
	}
	p.byAddr[addr] = h
	return h
}

func (p *Port) WriteMsg(addr netip.AddrPort, xor byte, hdr Header, enc *netmsg.State, msgs ...netmsg.Message) error {
	h := hdr.Encode()
	p.wmu.Lock()
	defer p.wmu.Unlock()
	var err error
	p.wbuf = p.wbuf[:0]
	p.wbuf = append(p.wbuf, h[0], h[1])
	for _, m := range msgs {
		p.wbuf, err = enc.Append(p.wbuf, m)
		if err != nil {
			return err
		}
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
	return p.WriteMsg(addr, 0, Header{Flags: Unreliable}, nil, m)
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

type Seq byte

func (v Seq) Before(v2 Seq) bool {
	return v <= v2 || (v >= 0xff-maxAckMsgs && v2-v < maxAckMsgs)
}

// PID is a packet ID.
type PID uintptr

type packet struct {
	pid      PID
	hdr      Header
	xor      byte
	lastSend time.Time
	deadline time.Time
	msgs     []netmsg.Message
	done     func()
	timeout  func()
}

func (p *packet) QueueID() QueueID {
	return QueueID{
		Stream: p.hdr.SID,
		Packet: p.pid,
	}
}

type Conn struct {
	p    *Port
	log  *slog.Logger
	addr netip.AddrPort
	enc  netmsg.State

	packetID atomic.Uintptr
	mu       sync.RWMutex
	xor      byte
	syn      Seq
	ack      Seq
	needAck  int
	nextPing time.Time
	queue    []*packet

	onMessage onMessageFuncs
}

func (p *Conn) EncodeState() *netmsg.State {
	return &p.enc
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
	var h Header
	h.Decode([2]byte{data[0], data[1]})
	reliable := h.Flags.Has(Reliable)
	data = data[2:]
	if p.p.debug {
		sdata := hex.EncodeToString(data)
		if reliable {
			p.log.Debug("RECV", "syn", h.Seq, "sid", h.SID, "data", sdata)
		} else {
			p.log.Debug("RECV", "ack", h.Seq, "sid", h.SID, "data", sdata)
		}
	}
	if reliable {
		// New reliable message that we should ACK in the future.
		exp := p.ack
		if h.Seq != exp {
			p.mu.Unlock()
			return // Ignore out of order packets.
		}
		p.ack = h.Seq + 1
		p.needAck++
		if p.needAck-1 >= maxAckMsgs {
			_ = p.sendAckPing()
		} else {
			p.nextPing = time.Now().Add(maxAckDelay)
		}
	} else {
		// Unreliable message with ACK for our reliable messages.
		p.queue = slices.DeleteFunc(p.queue, func(m *packet) bool {
			del := m.hdr.Seq.Before(h.Seq)
			if del && m.done != nil {
				doneFuncs = append(doneFuncs, m.done)
			}
			return del
		})
	}
	p.mu.Unlock()
	s := p.Stream(h.SID)
	onMsg := &p.onMessage
	onMsgGlobal := &p.p.onMessage
	for _, done := range doneFuncs {
		done()
	}

	var flags PacketFlags
	if reliable {
		flags |= Reliable
	}
	for len(data) > 0 {
		m, n, err := p.enc.DecodeNext(data)
		if err != nil {
			op := data[0]
			p.log.Error("Failed to decode packet", "op", op, "err", err)
			break
		}
		data = data[n:]
		if p.p.debug {
			p.log.Debug("RECV", "type", reflect.TypeOf(m).String(), "msg", m)
		}
		if onMsg.Call(s, m, flags) {
			continue
		}
		if onMsgGlobal.Call(s, m, flags) {
			continue
		}
	}
}

type QueueID struct {
	Stream SID
	Packet PID
}

func (p *Conn) ViewQueue(fnc func(id QueueID, m netmsg.Message)) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, m := range p.queue {
		id := m.QueueID()
		for _, msg := range m.msgs {
			fnc(id, msg)
		}
	}
}

func (p *Conn) QueuedFor(sid SID, ops ...netmsg.Op) int {
	n := 0
	p.ViewQueue(func(id QueueID, m netmsg.Message) {
		if sid != id.Stream {
			return
		}
		if len(ops) == 0 || slices.Contains(ops, m.NetOp()) {
			n++
		}
	})
	return n
}

func (p *Conn) DeleteQueue(fnc func(id QueueID, msgs []netmsg.Message) bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.queue = slices.DeleteFunc(p.queue, func(m *packet) bool {
		return fnc(m.QueueID(), m.msgs)
	})
}

func (p *Conn) ResetFor(sid SID) {
	p.DeleteQueue(func(id QueueID, _ []netmsg.Message) bool {
		return sid == id.Stream
	})
}

func (p *Conn) SendUnreliable(sid SID, msgs ...netmsg.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sendUnreliable(sid, msgs...)
}

func (p *Conn) sendUnreliable(sid SID, msgs ...netmsg.Message) error {
	seq := p.ack
	p.nextPing = time.Time{}
	if p.p.debug {
		typ := ""
		if len(msgs) > 0 {
			typ = reflect.TypeOf(msgs[0]).String()
		}
		p.log.Debug("SEND", "ack", seq, "sid", sid, "type", typ, "msgs", msgs)
	}
	h := Header{SID: sid, Seq: seq, Flags: Unreliable}
	return p.p.WriteMsg(p.addr, p.xor, h, &p.enc, msgs...)
}

func (p *Conn) sendAckPing() error {
	p.needAck = 0
	p.nextPing = time.Time{}
	return p.sendUnreliable(0)
}

func (p *Conn) Ack() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.sendAckPing()
}

type Options struct {
	Context   context.Context
	Deadline  time.Time
	OnDone    func()
	OnTimeout func()
}

func (p *Conn) QueueReliable(sid SID, opts Options, msgs ...netmsg.Message) PID {
	msgs = slices.Clone(msgs)
	now := time.Now()
	var deadline time.Time
	if !opts.Deadline.IsZero() {
		deadline = opts.Deadline
	}
	if opts.Context != nil {
		deadline, _ = opts.Context.Deadline()
	}
	if deadline.IsZero() {
		deadline = now.Add(defaultTimeout)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	seq := p.syn
	p.syn++
	pid := PID(p.packetID.Add(1))
	p.queue = append(p.queue, &packet{
		pid: pid,
		hdr: Header{
			SID:   sid,
			Seq:   seq,
			Flags: Reliable,
		},
		xor:      p.xor,
		lastSend: time.Time{}, // send in the next tick
		deadline: deadline,
		msgs:     msgs,
		done:     opts.OnDone,
		timeout:  opts.OnTimeout,
	})
	return pid
}

func (p *Conn) CancelReliable(pid PID) {
	p.DeleteQueue(func(id QueueID, _ []netmsg.Message) bool {
		return pid != id.Packet
	})
}

func (p *Conn) SendReliable(ctx context.Context, sid SID, msgs ...netmsg.Message) error {
	var cancel func()
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	acked := make(chan struct{})
	pid := p.QueueReliable(sid, Options{
		Context: ctx,
		OnDone: func() {
			close(acked)
		},
		OnTimeout: cancel,
	}, msgs...)
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
				p.log.Debug("SEND", "syn", m.hdr.Seq, "sid", m.hdr.SID, "msgs", m.msgs)
			}
			if err := p.p.WriteMsg(p.addr, m.xor, m.hdr, &p.enc, m.msgs...); err != nil {
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

func (p *Conn) Stream(sid SID) Stream {
	return Stream{p: p, sid: sid}
}

type Stream struct {
	p   *Conn
	sid SID
}

func (p Stream) Valid() bool {
	return p.p != nil
}

func (p Stream) Conn() *Conn {
	return p.p
}

func (p Stream) SID() SID {
	return p.sid
}

func (p Stream) Addr() netip.AddrPort {
	return p.p.RemoteAddr()
}

func (p Stream) Reset() {
	p.p.ResetFor(p.sid)
}

func (p Stream) SendQueue() error {
	return p.p.sendQueue(func(m *packet) bool {
		return m.hdr.SID == p.sid
	})
}

func (p Stream) OnMessage(fnc OnMessageFunc) {
	p.p.onMessage.Add(func(s Stream, m netmsg.Message, flags PacketFlags) bool {
		if p != s {
			return false
		}
		return fnc(s, m, flags)
	})
}

func (p Stream) SendUnreliable(msgs ...netmsg.Message) error {
	return p.p.SendUnreliable(p.sid, msgs...)
}

func (p Stream) QueueReliable(opts Options, msgs ...netmsg.Message) PID {
	return p.p.QueueReliable(p.sid, opts, msgs...)
}

func (p Stream) CancelReliable(id PID) {
	p.p.CancelReliable(id)
}

func (p Stream) SendReliable(ctx context.Context, msgs ...netmsg.Message) error {
	return p.p.SendReliable(ctx, p.sid, msgs...)
}
