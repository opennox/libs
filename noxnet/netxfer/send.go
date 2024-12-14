package netxfer

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/opennox/libs/binenc"
)

const (
	sendClosed = sendState(iota)
	sendStarted
	sendAccepted
)

type DoneFunc func()
type AbortFunc func(reason Error)

type sendState byte

type sendChunk struct {
	ind      Chunk
	data     []byte
	lastSent time.Duration
	retries  uint16
	next     *sendChunk
	prev     *sendChunk
}

type sendStream[C Conn] struct {
	x         *Sender[C]
	conn      C
	id        SendID
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

func (p *sendStream[C]) Close() {
	if p == nil || p.x == nil {
		return
	}
	p.x.arr[p.id] = nil
	p.x.active--
	*p = sendStream[C]{}
}

func (p *sendStream[C]) Done() {
	if p == nil || p.x == nil {
		return
	}
	p.callDone()
	p.Close()
}

func (p *sendStream[C]) Abort() {
	if p == nil || p.x == nil {
		return
	}
	p.callAborted(ErrClosed)
	p.Close()
}

func (p *sendStream[C]) Cancel(reason Error) {
	if p.state != sendAccepted {
		return
	}
	_ = p.conn.SendReliable(&MsgXfer{&MsgCancel{
		RecvID: p.recvID,
		Reason: reason,
	}})
	p.callAborted(reason)
	p.Close()
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
			*it = sendChunk{} // GC
			return
		}
	}
}

func (p *sendStream[C]) Accept(rid RecvID) {
	if p == nil || p.x == nil {
		return
	}
	p.recvID = rid
	p.state = sendAccepted
}

func (p *sendStream[C]) Start(d Data) {
	p.state = sendStarted
	p.action = d.Action
	p.typ = d.Type
	p.size = len(d.Data)
	p.chunkCnt = (len(d.Data)-1)/blockSize + 1
	p.nextChunk = 1
	left := d.Data
	for i := range p.chunkCnt {
		b := &sendChunk{
			ind: Chunk(i + 1),
		}

		n := blockSize
		if len(left) <= n {
			n = len(left)
		}
		b.data = slices.Clone(left[:n])
		left = left[n:]

		b.prev, b.next = p.last, nil
		if prev := p.last; prev != nil {
			prev.next = b
		} else {
			p.first = b
		}
		p.last = b
	}
	_ = p.conn.SendReliable(&MsgXfer{&MsgStart{
		Act:    d.Action,
		Size:   uint32(len(d.Data)),
		Type:   binenc.String{Value: d.Type},
		SendID: p.id,
	}})
}

func (p *sendStream[C]) Update(ts time.Duration) {
	if p == nil || p.x == nil {
		return
	}
	if p.state != sendAccepted {
		return
	}
	for j, b := 0, p.first; j < 2 && b != nil; j, b = j+1, b.next {
		if t := b.lastSent; t == 0 {
			_ = p.conn.SendUnreliable(&MsgXfer{&MsgData{
				RecvID: p.recvID,
				Token:  0,
				Chunk:  b.ind,
				Data:   b.data,
			}})
			p.nextChunk++
			b.lastSent = ts
		} else if ts > t+retryInterval {
			if b.retries >= maxRetries {
				p.Cancel(ErrSendTimeout)
				return
			}
			_ = p.conn.SendUnreliable(&MsgXfer{&MsgData{
				RecvID: p.recvID,
				Token:  0,
				Chunk:  b.ind,
				Data:   b.data,
			}})
			b.lastSent = ts
			b.retries++
		}
	}
}

func (p *sendStream[C]) callAborted(reason Error) {
	if p.onAbort != nil {
		p.onAbort(reason)
	}
}

func (p *sendStream[C]) callDone() {
	if p.onDone != nil {
		p.onDone()
	}
}

type Sender[C Conn] struct {
	arr    []*sendStream[C]
	active int
}

func (x *Sender[C]) Reset(n int) {
	if n < 0 {
		n = minStreams
	} else if n > maxStreams {
		n = maxStreams
	}
	x.arr = make([]*sendStream[C], n)
	x.active = 0
}

func (x *Sender[C]) getS(id SendID) *sendStream[C] {
	if int(id) >= len(x.arr) {
		return nil
	}
	return x.arr[id]
}

func (x *Sender[C]) getR(conn C, rid RecvID) *sendStream[C] {
	for _, it := range x.arr {
		if it == nil || it.state != sendAccepted {
			continue
		}
		if it.conn == conn && it.recvID == rid {
			return it
		}
	}
	return nil
}

func (x *Sender[C]) HandleAccept(conn C, m *MsgAccept) {
	x.getS(m.SendID).Accept(m.RecvID)
}

func (x *Sender[C]) HandleAck(conn C, m *MsgAck) {
	x.getR(conn, m.RecvID).Ack(m.Chunk)
}

func (x *Sender[C]) HandleDone(conn C, m *MsgDone) {
	x.getR(conn, m.RecvID).Done()
}

func (x *Sender[C]) HandleAbort(conn C, m *MsgAbort) {
	x.getR(conn, m.RecvID).Abort()
}

func (x *Sender[C]) add(s *sendStream[C]) bool {
	for i, it := range x.arr {
		if it == nil {
			x.arr[i] = s
			x.active++
			s.id = SendID(i)
			return true
		}
	}
	return false
}

func (x *Sender[C]) StartSend(conn C, p Data, onDone DoneFunc, onAbort AbortFunc) (SendID, bool) {
	s := &sendStream[C]{
		x:       x,
		conn:    conn,
		recvID:  0,
		size:    len(p.Data),
		onDone:  onDone,
		onAbort: onAbort,
	}
	if !x.add(s) {
		return 0, false
	}
	s.Start(p)
	return s.id, true
}

func (x *Sender[C]) Send(ctx context.Context, conn C, p Data) error {
	done := make(chan struct{})
	abort := make(chan Error, 1)
	id, ok := x.StartSend(conn, p, func() {
		close(done)
	}, func(reason Error) {
		select {
		case abort <- reason:
		default:
		}
	})
	if !ok {
		return errors.New("send failed, too many streams")
	}
	select {
	case <-ctx.Done():
		x.Cancel(conn, id)
		return ctx.Err()
	case err := <-abort:
		return err
	case <-done:
		return nil
	}
}

func (x *Sender[C]) CancelAll(conn C) {
	for _, it := range x.arr {
		if it == nil {
			continue
		}
		if it.conn == conn {
			it.Cancel(ErrClosed)
		}
	}
}

func (x *Sender[C]) Cancel(conn C, id SendID) {
	x.getS(id).Cancel(ErrClosed)
}

func (x *Sender[C]) Update(ts time.Duration) {
	if x.active == 0 {
		return
	}
	for _, it := range x.arr {
		if it == nil {
			continue
		}
		it.Update(ts)
	}
}

func (x *Sender[C]) Handle(conn C, ts time.Duration, m Msg) {
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
