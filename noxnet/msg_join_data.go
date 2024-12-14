package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgJoinData{}, false)
}

type MsgJoinData struct {
	NetCode NetCode
	Unk2    uint32
}

func (*MsgJoinData) NetOp() netmsg.Op {
	return netmsg.MSG_JOIN_DATA
}

func (*MsgJoinData) EncodeSize() int {
	return 6
}

func (m *MsgJoinData) Encode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], uint16(m.NetCode))
	binary.LittleEndian.PutUint32(data[2:6], m.Unk2)
	return 6, nil
}

func (m *MsgJoinData) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, io.ErrUnexpectedEOF
	}
	m.NetCode = NetCode(binary.LittleEndian.Uint16(data[0:2]))
	m.Unk2 = binary.LittleEndian.Uint32(data[2:6])
	return 6, nil
}
