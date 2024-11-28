package noxnet

import (
	"github.com/opennox/libs/noxnet/netxfer"
)

func init() {
	RegisterMessage(&MsgXfer{}, true)
}

type MsgXfer struct {
	Msg netxfer.Msg
}

func (*MsgXfer) NetOp() Op {
	return MSG_XFER_MSG
}

func (m *MsgXfer) EncodeSize() int {
	return netxfer.EncodeSize(m.Msg)
}

func (m *MsgXfer) Encode(data []byte) (int, error) {
	return netxfer.Encode(data, m.Msg)
}

func (m *MsgXfer) Decode(data []byte) (int, error) {
	msg, n, err := netxfer.Decode(data)
	m.Msg = msg
	return n, err
}
