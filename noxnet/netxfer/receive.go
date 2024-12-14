package netxfer

import (
	"slices"
	"time"
)

type ReceiveFunc[C Conn] func(conn C, p Data)

type recvChunk struct {
	ind  Chunk
	data []byte
	next *recvChunk
	prev *recvChunk
}

type recvStream[C Conn] struct {
	x          *Receiver[C]
	conn       C
	id         RecvID
	full       Data
	received   int // bytes
	lastUpdate time.Duration

	nextChunk Chunk // for large payloads
	first     *recvChunk
	last      *recvChunk
}

func (p *recvStream[C]) AddChunk(ts time.Duration, m *MsgData) {
	if p == nil {
		return
	}
	p.lastUpdate = ts
	_ = p.conn.SendReliable(&MsgXfer{&MsgAck{
		RecvID: p.id,
		Token:  m.Token,
		Chunk:  m.Chunk,
	}})
	if len(m.Data) == 0 {
		return
	}
	defer func() {
		if p.received >= len(p.full.Data) {
			p.Done()
		}
	}()
	chunk := m.Chunk
	data := m.Data
	if len(p.full.Data) < maxDataSize {
		if chunk == 0 {
			return // chunks must start from 1
		}
		// Simplified code path for small payloads.
		// Just put the block where it's supposed to be.
		off := (int(chunk) - 1) * blockSize
		if off > len(p.full.Data) {
			return
		}
		p.received += copy(p.full.Data[off:], data)
		return
	}
	// For large payloads, we have to handle overflows and interpret
	// chunk number as an ID and not an index.
	if chunk == p.nextChunk {
		cur := p.received
		p.received += copy(p.full.Data[cur:cur+len(data)], data)
		p.nextChunk++
		return
	}
	b := &recvChunk{
		ind:  chunk,
		data: slices.Clone(data),
	}
	b.prev, b.next = p.last, nil
	if last := p.last; last != nil {
		last.next = b
	}
	p.last = b
	if first := p.first; first == nil {
		p.first = b
	}

	for it := p.first; it != nil; it = it.next {
		if p.nextChunk == it.ind {
			cur := p.received
			p.received += copy(p.full.Data[cur:], it.data)
			p.nextChunk++
			if prev := it.prev; prev != nil {
				prev.next = it.next
			} else {
				p.first = it.next
			}
			if next := it.next; next != nil {
				next.prev = it.prev
			} else {
				p.last = it.prev
			}
			*it = recvChunk{} // GC
		}
	}
}

func (p *recvStream[C]) Close() {
	if p == nil || p.x == nil {
		return
	}
	p.x.arr[p.id] = nil
	p.x.active--
	*p = recvStream[C]{}
}

func (p *recvStream[C]) Abort(reason Error) {
	if p == nil || p.x == nil {
		return
	}
	_ = p.conn.SendReliable(&MsgXfer{&MsgAbort{
		RecvID: p.id,
		Reason: reason,
	}})
	p.Close()
}

func (p *recvStream[C]) Done() {
	if p == nil || p.x == nil {
		return
	}
	_ = p.conn.SendReliable(&MsgXfer{&MsgDone{
		RecvID: p.id,
	}})
	if p.x.onReceive != nil {
		p.x.onReceive(p.conn, p.full)
	}
	p.Close()
}

func (p *recvStream[C]) Accept(sid SendID) {
	if p == nil || p.x == nil {
		return
	}
	_ = p.conn.SendReliable(&MsgXfer{&MsgAccept{
		RecvID: p.id,
		SendID: sid,
	}})
	if len(p.full.Data) == 0 {
		p.Done()
	}
}

func (p *recvStream[C]) Update(ts time.Duration) {
	if p == nil || p.x == nil {
		return
	}
	if ts > p.lastUpdate+receiveTimeout {
		p.Abort(ErrRecvTimeout)
	}
}

type Receiver[C Conn] struct {
	arr    []*recvStream[C]
	active int

	onReceive ReceiveFunc[C]
}

func (x *Receiver[C]) OnReceive(fnc ReceiveFunc[C]) {
	x.onReceive = fnc
}

func (x *Receiver[C]) Reset(n int) {
	if n < 0 {
		n = minStreams
	} else if n > maxStreams {
		n = maxStreams
	}
	x.active = 0
	x.arr = make([]*recvStream[C], n)
}

func (x *Receiver[C]) add(s *recvStream[C]) bool {
	for i, p := range x.arr {
		if p == nil {
			s.id = RecvID(i)
			x.arr[i] = s
			x.active++
			return true
		}
	}
	return false
}

func (x *Receiver[C]) get(conn C, rid RecvID) *recvStream[C] {
	if int(rid) >= len(x.arr) {
		return nil
	}
	return x.arr[rid]
}

func (x *Receiver[C]) Update(ts time.Duration) {
	if x.active == 0 {
		return
	}
	for _, it := range x.arr {
		it.Update(ts)
	}
}

func (x *Receiver[C]) HandleStart(conn C, ts time.Duration, m *MsgStart) {
	s := &recvStream[C]{
		x:          x,
		conn:       conn,
		lastUpdate: ts,
		nextChunk:  1,
		full: Data{
			Action: m.Act,
			Type:   m.Type.Value,
			Data:   make([]byte, m.Size),
		},
	}
	if x.add(s) {
		s.Accept(m.SendID)
	}
}

func (x *Receiver[C]) HandleData(conn C, ts time.Duration, m *MsgData) {
	x.get(conn, m.RecvID).AddChunk(ts, m)
}

func (x *Receiver[C]) HandleCancel(conn C, m *MsgCancel) {
	x.get(conn, m.RecvID).Close()
}

func (x *Receiver[C]) Handle(conn C, ts time.Duration, m Msg) {
	switch m := m.(type) {
	case *MsgStart:
		x.HandleStart(conn, ts, m)
	case *MsgCancel:
		x.HandleCancel(conn, m)
	case *MsgData:
		x.HandleData(conn, ts, m)
	}
}
