package noxnet

import (
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgFxJiggle{}, false)
}

type MsgFxJiggle struct {
	Val byte
}

func (*MsgFxJiggle) NetOp() netmsg.Op {
	return netmsg.MSG_FX_JIGGLE
}

func (*MsgFxJiggle) EncodeSize() int {
	return 1
}

func (p *MsgFxJiggle) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Val
	return 1, nil
}

func (p *MsgFxJiggle) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Val = data[0]
	return 1, nil
}
