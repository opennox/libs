package netxfer

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/binenc"
)

type TestConn struct {
	TS   *Timestamp
	Recv interface {
		Handle(conn *TestConn, ts Timestamp, m Msg)
	}
	Peer *TestConn
	Buf  []Msg
}

func (c *TestConn) SendReliableMsg(m Msg) error {
	c.Buf = append(c.Buf, m)
	c.Recv.Handle(c.Peer, *c.TS, m)
	return nil
}

func (c *TestConn) SendUnreliableMsg(m Msg) error {
	c.Buf = append(c.Buf, m)
	c.Recv.Handle(c.Peer, *c.TS, m)
	return nil
}

func TestXfer(t *testing.T) {
	var (
		ts    Timestamp
		send  Sender[*TestConn]
		recv  Receiver[*TestConn]
		check func(p Data)
	)
	cs := &TestConn{TS: &ts, Recv: &send}
	cr := &TestConn{TS: &ts, Recv: &recv}
	cs.Peer, cr.Peer = cr, cs

	send.Init(minStreams)
	recv.Init(minStreams, func(_ *TestConn, p Data) {
		check(p)
	})

	tick := func(expSend, expRecv []Msg) {
		send.Update(ts)
		must.Eq(t, expSend, cr.Buf, must.Sprint("tick", ts))
		cr.Buf = nil

		recv.Update(ts)
		must.Eq(t, expRecv, cs.Buf, must.Sprint("tick", ts))
		cs.Buf = nil

		ts++
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
	gotTS := Timestamp(0)
	sentTS := Timestamp(0)
	check = func(p2 Data) {
		must.EqOp(t, p.Action, p2.Action)
		must.EqOp(t, p.Type, p2.Type)
		must.Eq(t, data, p2.Data)
		gotCnt++
		gotTS = ts
	}
	send.Send(cr, p, func() {
		sentTS = ts
	}, func() {
		t.Fatal("aborted")
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
				&MsgData{
					Chunk: 1,
					Data:  data[blockSize*0 : blockSize*1],
				},
			},
			FromRecv: []Msg{
				&MsgAccept{SendID: 0, RecvID: 0},
				&MsgAck{RecvID: 0, Chunk: 1},
			},
		},
		1: {
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
		2: {
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
		3: {},
		4: {},
	}
	for _, exp := range ticks {
		tick(exp.FromSend, exp.FromRecv)
	}
	must.EqOp(t, 1, gotCnt)
	must.EqOp(t, 1+2, gotTS)
	must.EqOp(t, 1+2, sentTS)
}
