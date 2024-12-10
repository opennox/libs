package netxfer

import (
	"math"
	"time"

	"github.com/opennox/libs/noxnet/netmsg"
)

const (
	minStreams     = 16
	maxStreams     = 256
	blockSize      = 512
	maxBlocks      = math.MaxUint16 - 1
	maxDataSize    = maxBlocks * blockSize
	retryInterval  = 3 * time.Second
	receiveTimeout = 30 * time.Second
	maxRetries     = 20
)

type Data struct {
	Action Action
	Type   string
	Data   []byte
}

type Conn interface {
	comparable
	SendReliableMsg(m netmsg.Message) error
	SendUnreliableMsg(m netmsg.Message) error
}

type State[C Conn] struct {
	Sender[C]
	Receiver[C]
}

func (x *State[C]) Reset(n int) {
	x.Sender.Reset(n)
	x.Receiver.Reset(n)
}

func (x *State[C]) Update(ts time.Duration) {
	x.Sender.Update(ts)
	x.Receiver.Update(ts)
}

func (x *State[C]) Handle(conn C, ts time.Duration, m Msg) {
	switch m := m.(type) {
	case *MsgStart:
		x.Receiver.HandleStart(conn, ts, m)
	case *MsgAccept:
		x.Sender.HandleAccept(conn, m)
	case *MsgCancel:
		x.Receiver.HandleCancel(conn, m)
	case *MsgAbort:
		x.Sender.HandleAbort(conn, m)
	case *MsgData:
		x.Receiver.HandleData(conn, ts, m)
	case *MsgAck:
		x.Sender.HandleAck(conn, m)
	case *MsgDone:
		x.Sender.HandleDone(conn, m)
	}
}
