package netxfer

const (
	minStreams    = 16
	maxStreams    = 256
	blockSize     = 512
	retryInterval = Timestamp(90)
	maxRetries    = 20
)

type Data struct {
	Action Action
	Type   string
	Data   []byte
}

type Conn interface {
	comparable
	SendReliableMsg(m Msg) error
	SendUnreliableMsg(m Msg) error
}

type Timestamp uint32

type State[C Conn] struct {
	Sender[C]
	Receiver[C]
}

func (x *State[C]) Init(n int, onRecv ReceiveFunc[C]) {
	x.Sender.Init(n)
	x.Receiver.Init(n, onRecv)
}

func (x *State[C]) Free() {
	x.Sender.Free()
	x.Receiver.Free()
}

func (x *State[C]) Update(ts Timestamp) {
	x.Sender.Update(ts)
	x.Receiver.Update(ts)
}

func (x *State[C]) Handle(conn C, ts Timestamp, m Msg) {
	switch m := m.(type) {
	case *MsgStart:
		x.Receiver.HandleStart(conn, ts, m)
	case *MsgAccept:
		x.Sender.HandleAccept(conn, m)
	case *MsgCancel:
		x.Receiver.HandleCancel(m)
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
