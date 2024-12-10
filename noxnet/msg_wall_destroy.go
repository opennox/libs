package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgWallDestroy{}, false)
}

type MsgWallDestroy struct {
	ID uint16
}

func (*MsgWallDestroy) NetOp() netmsg.Op {
	return netmsg.MSG_DESTROY_WALL
}

func (*MsgWallDestroy) EncodeSize() int {
	return 2
}

func (p *MsgWallDestroy) Encode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:], p.ID)
	return 2, nil
}

func (p *MsgWallDestroy) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	p.ID = binary.LittleEndian.Uint16(data[0:])
	return 2, nil
}
