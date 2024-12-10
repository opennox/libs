package netxfer

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/binenc"
)

const tickStep = 10 * time.Millisecond

type TestState struct {
	TS       time.Duration
	Sender   Sender[*TestConn]
	Receiver Receiver[*TestConn]
	Send     *TestConn
	Recv     *TestConn
}

type TestConn struct {
	s    *TestState
	Recv interface {
		Handle(conn *TestConn, ts time.Duration, m Msg)
	}
	Peer *TestConn
	Buf  []Msg
}

func (s *TestState) Tick() {
	s.TS += 10 * time.Millisecond
	s.Send.Buf = nil
	s.Recv.Buf = nil
	s.Sender.Update(s.TS)
	s.Receiver.Update(s.TS)
}

func (s *TestState) StartSend(p Data, onDone DoneFunc, onAbort AbortFunc) bool {
	_, ok := s.Sender.StartSend(s.Recv, p, onDone, onAbort)
	return ok
}

func (c *TestConn) SendReliableMsg(m Msg) error {
	c.Buf = append(c.Buf, m)
	c.Recv.Handle(c.Peer, c.s.TS, m)
	return nil
}

func (c *TestConn) SendUnreliableMsg(m Msg) error {
	c.Buf = append(c.Buf, m)
	c.Recv.Handle(c.Peer, c.s.TS, m)
	return nil
}

func newState(check func(s *TestState, p Data)) *TestState {
	s := &TestState{}
	cs := &TestConn{s: s, Recv: &s.Sender}
	cr := &TestConn{s: s, Recv: &s.Receiver}
	cs.Peer, cr.Peer = cr, cs
	s.Send = cs
	s.Recv = cr

	s.Sender.Reset(minStreams)
	s.Receiver.Reset(minStreams)
	s.Receiver.OnReceive(func(_ *TestConn, p Data) {
		check(s, p)
	})
	return s
}

func TestXferFlow(t *testing.T) {
	var check func(p Data)
	s := newState(func(_ *TestState, p Data) {
		check(p)
	})

	tick := func(expSend, expRecv []Msg) {
		opts := must.Sprint("tick", s.TS)
		must.Eq(t, expSend, s.Recv.Buf, opts)
		must.Eq(t, expRecv, s.Send.Buf, opts)
		s.Tick()
	}

	tick(nil, nil)

	data := make([]byte, blockSize*5/2)
	r := rand.NewChaCha8([32]byte{1, 2, 3})
	r.Read(data)

	p := Data{
		Action: 123,
		Type:   "Test",
		Data:   slices.Clone(data),
	}

	gotCnt := 0
	gotTS := time.Duration(0)
	sentTS := time.Duration(0)
	check = func(p2 Data) {
		must.EqOp(t, p.Action, p2.Action)
		must.EqOp(t, p.Type, p2.Type)
		must.Eq(t, data, p2.Data)
		gotCnt++
		gotTS = s.TS
	}
	s.StartSend(p, func() {
		sentTS = s.TS
	}, func(reason Error) {
		t.Fatal("aborted", reason)
	})
	ticks := []struct {
		FromSend []Msg
		FromRecv []Msg
	}{
		0: {
			FromSend: []Msg{
				&MsgStart{
					SendID: 0,
					Act:    p.Action,
					Type:   binenc.String{Value: p.Type},
					Size:   uint32(len(data)),
				},
			},
			FromRecv: []Msg{
				&MsgAccept{SendID: 0, RecvID: 0},
			},
		},
		1: {
			FromSend: []Msg{
				&MsgData{
					Chunk: 1,
					Data:  data[blockSize*0 : blockSize*1],
				},
			},
			FromRecv: []Msg{
				&MsgAck{RecvID: 0, Chunk: 1},
			},
		},
		2: {
			FromSend: []Msg{
				&MsgData{
					Chunk: 2,
					Data:  data[blockSize*1 : blockSize*2],
				},
			},
			FromRecv: []Msg{
				&MsgAck{RecvID: 0, Chunk: 2},
			},
		},
		3: {
			FromSend: []Msg{
				&MsgData{
					Chunk: 3,
					Data:  data[blockSize*2:],
				},
			},
			FromRecv: []Msg{
				&MsgAck{RecvID: 0, Chunk: 3},
				&MsgDone{RecvID: 0},
			},
		},
		4: {},
		5: {},
	}
	for _, exp := range ticks {
		tick(exp.FromSend, exp.FromRecv)
	}
	must.EqOp(t, 1, gotCnt)
	must.EqOp(t, (1+3)*tickStep, gotTS)
	must.EqOp(t, (1+3)*tickStep, sentTS)
}

func TestXferSize(t *testing.T) {
	cases := []struct {
		name string
		size int
	}{
		{"empty", 0},
		{"tiny", 10},
		{"block", blockSize},
		{"multiple", blockSize * 3},
		{"with half", blockSize * 5 / 2},
		{"large", maxDataSize + blockSize + 10},
	}
	r := rand.NewChaCha8([32]byte{1, 2, 3})
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := make([]byte, c.size)
			r.Read(data)

			done := false
			s := newState(func(_ *TestState, p Data) {
				done = true
				must.True(t, bytes.Equal(data, p.Data))
			})
			ok := s.StartSend(Data{
				Action: 1,
				Type:   c.name,
				Data:   slices.Clone(data),
			}, nil, func(reason Error) {
				t.Fatal("aborted", reason)
			})
			must.True(t, ok, must.Sprint("send failed"))
			for i := 0; !done; i++ {
				if i >= maxBlocks+10 {
					t.Fatal("too many ticks")
				}
				s.Tick()
			}
			must.True(t, done)
		})
	}
}
