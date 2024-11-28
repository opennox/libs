package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterMessage(&MsgUseMap{}, false)
}

type MsgUseMap struct {
	MapName binenc.String
	CRC     uint32
	T       Timestamp
}

func (*MsgUseMap) NetOp() Op {
	return MSG_USE_MAP
}

func (*MsgUseMap) EncodeSize() int {
	return 40
}

func (m *MsgUseMap) Encode(data []byte) (int, error) {
	if len(data) < 40 {
		return 0, io.ErrShortBuffer
	}
	m.MapName.Encode(data[0:32])
	binary.LittleEndian.PutUint32(data[32:36], m.CRC)
	binary.LittleEndian.PutUint32(data[36:40], uint32(m.T))
	return 40, nil
}

func (m *MsgUseMap) Decode(data []byte) (int, error) {
	if len(data) < 40 {
		return 0, io.ErrUnexpectedEOF
	}
	m.MapName.Decode(data[0:32])
	m.CRC = binary.LittleEndian.Uint32(data[32:36])
	m.T = Timestamp(binary.LittleEndian.Uint32(data[36:40]))
	return 40, nil
}
