package netxfer

import (
	"slices"

	"github.com/opennox/libs/binenc"
)

const (
	sendFree     = sendState(0)
	sendStarted  = sendState(1)
	sendAccepted = sendState(2)
)

type DoneFunc func()
type AbortFunc func()

type sendState byte

type sendChunk struct {
	ind      Chunk
	data     []byte
	lastSent Timestamp
	retries  uint16
	next     *sendChunk
	prev     *sendChunk
}

type sendStream[C Conn] struct {
	conn      C
	recvID    RecvID
	state     sendState
	size      int
	nextChunk int
	chunkCnt  int
	action    Action
	typ       string
	first     *sendChunk
	last      *sendChunk
	onDone    DoneFunc
	onAbort   AbortFunc
}

func (p *sendStream[C]) Reset() {
	*p = sendStream[C]{
		state:     sendFree,
		nextChunk: 1,
	}
}

func (p *sendStream[C]) Free() {
	for it := p.first; it != nil; it = it.next {
		*it = sendChunk{}
	}
}

func (p *sendStream[C]) Ack(chunk Chunk) {
	if p == nil {
		return
	}
	for it := p.first; it != nil; it = it.next {
		if it.ind == chunk {
			if next := it.next; next != nil {
				next.prev = it.prev
			} else {
				p.last = it.prev
			}
			if prev := it.prev; prev != nil {
				prev.next = it.next
			} else {
				p.first = it.next
			}
			*it = sendChunk{}
			return
		}
	}
}

func (p *sendStream[C]) callAborted() {
	if p.onAbort != nil {
		p.onAbort()
	}
}

func (p *sendStream[C]) callDone() {
	if p.onDone != nil {
		p.onDone()
	}
}

type Sender[C Conn] struct {
	arr    []sendStream[C]
	cnt    int
	active int
}

func (x *Sender[C]) Init(n int) {
	if n < 0 {
		n = minStreams
	} else if n > maxStreams {
		n = maxStreams
	}
	x.cnt = n
	x.arr = make([]sendStream[C], n)
	for i := 0; i < n; i++ {
		x.arr[i].Reset()
	}
	x.active = 0
}

func (x *Sender[C]) find(conn C, rid RecvID) *sendStream[C] {
	for i := 0; i < x.cnt; i++ {
		it := &x.arr[i]
		if it.conn == conn && it.recvID == rid {
			return it
		}
	}
	return nil
}

func (x *Sender[C]) HandleAccept(conn C, m *MsgAccept) {
	id := m.SendID
	if int(id) >= x.cnt {
		return
	}
	p := &x.arr[id]
	p.recvID = m.RecvID
	p.state = sendAccepted
}

func (x *Sender[C]) HandleAck(conn C, m *MsgAck) {
	x.find(conn, m.RecvID).Ack(m.Chunk)
}

func (x *Sender[C]) HandleDone(conn C, m *MsgDone) {
	p := x.find(conn, m.RecvID)
	if p == nil {
		return
	}
	p.callDone()
	x.active--
	p.Reset()
}

func (x *Sender[C]) HandleAbort(conn C, m *MsgAbort) {
	p := x.find(conn, m.RecvID)
	if p == nil {
		return
	}
	p.callAborted()
	for it := p.first; it != nil; it = it.next {
		*it = sendChunk{}
	}
	p.Reset()
}

func (x *Sender[C]) Send(conn C, p Data, onDone DoneFunc, onAbort AbortFunc) bool {
	if len(p.Data) == 0 {
		return false
	}
	s, id := x.newStream()
	if s == nil {
		return false
	}
	x.active++
	left := p.Data
	blocks := (len(p.Data)-1)/blockSize + 1
	for i := range blocks {
		b := &sendChunk{
			ind: Chunk(i + 1),
		}

		n := blockSize
		if len(left) <= n {
			n = len(left)
		}
		b.data = slices.Clone(left[:n])
		left = left[n:]

		b.prev, b.next = s.last, nil
		if prev := s.last; prev != nil {
			prev.next = b
		} else {
			s.first = b
		}
		s.last = b
	}
	s.conn = conn
	s.recvID = 0
	s.state = sendStarted
	s.size = len(p.Data)
	s.nextChunk = 1
	s.chunkCnt = blocks
	if p.Type != "" {
		s.typ = p.Type
	}
	s.action = p.Action
	s.onDone = onDone
	s.onAbort = onAbort
	_ = conn.SendReliableMsg(&MsgStart{
		Act:    p.Action,
		Size:   uint32(len(p.Data)),
		Type:   binenc.String{Value: p.Type},
		SendID: id,
	})
	return true
}

func (x *Sender[C]) Cancel(conn C) {
	for i := 0; i < x.cnt; i++ {
		it := &x.arr[i]
		if it.state == sendAccepted && it.conn == conn {
			x.cancel(it, ErrClosed)
		}
	}
}

func (x *Sender[C]) newStream() (*sendStream[C], SendID) {
	for i := 0; i < x.cnt; i++ {
		it := &x.arr[i]
		if it.state == sendFree && it.size == 0 {
			return it, SendID(i)
		}
	}
	return nil, 0
}

func (x *Sender[C]) Free() {
	x.arr = nil
}

func (x *Sender[C]) Update(ts Timestamp) {
	if x.active == 0 {
		return
	}
	if x.cnt <= 0 {
		return
	}
	for i := 0; i < x.cnt; i++ {
		s := &x.arr[i]
		if s.state != sendAccepted {
			continue
		}
		for j, b := 0, s.first; j < 2 && b != nil; j, b = j+1, b.next {
			if t := b.lastSent; t == 0 {
				_ = s.conn.SendUnreliableMsg(&MsgData{
					RecvID: s.recvID,
					Token:  0,
					Chunk:  b.ind,
					Data:   b.data,
				})
				s.nextChunk++
				b.lastSent = ts
			} else if ts > t+retryInterval {
				if b.retries < maxRetries {
					_ = s.conn.SendUnreliableMsg(&MsgData{
						RecvID: s.recvID,
						Token:  0,
						Chunk:  b.ind,
						Data:   b.data,
					})
					b.lastSent = ts
					b.retries++
				} else if s.state == sendAccepted {
					x.cancel(s, ErrSendTimeout)
					break
				}
			}
		}
	}
}

func (x *Sender[C]) cancel(s *sendStream[C], reason Error) {
	_ = s.conn.SendReliableMsg(&MsgCancel{
		RecvID: s.recvID,
		Reason: reason,
	})
	s.callAborted()
	s.Free()
	s.Reset()
	if x.active != 0 {
		x.active--
	}
}

func (x *Sender[C]) Handle(conn C, ts Timestamp, m Msg) {
	switch m := m.(type) {
	case *MsgAccept:
		x.HandleAccept(conn, m)
	case *MsgAbort:
		x.HandleAbort(conn, m)
	case *MsgAck:
		x.HandleAck(conn, m)
	case *MsgDone:
		x.HandleDone(conn, m)
	}
}
