package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgForgetDrawables{}, false)
}

type MsgForgetDrawables struct {
	Unk0 uint32
}

func (*MsgForgetDrawables) NetOp() netmsg.Op {
	return netmsg.MSG_FORGET_DRAWABLES
}

func (*MsgForgetDrawables) EncodeSize() int {
	return 4
}

func (p *MsgForgetDrawables) Encode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], p.Unk0)
	return 4, nil
}

func (p *MsgForgetDrawables) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = binary.BigEndian.Uint32(data[0:4])
	return 4, nil
}
