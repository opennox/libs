package noxnet

import (
	"encoding/binary"
	"image"
	"io"

	"github.com/opennox/opennox-lib/binenc"
)

func init() {
	RegisterMessage(&MsgAccept{}, false)
	RegisterMessage(&MsgServerAccept{}, false)
	RegisterMessage(&MsgClientAccept{}, false)
}

type MsgAccept struct {
	ID byte
}

func (*MsgAccept) NetOp() Op {
	return MSG_ACCEPTED
}

func (*MsgAccept) EncodeSize() int {
	return 1
}

func (p *MsgAccept) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.ID
	return 1, nil
}

func (p *MsgAccept) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	p.ID = data[0]
	return 1, nil
}

type MsgServerAccept struct {
	ID     uint32
	XorKey byte
}

func (*MsgServerAccept) NetOp() Op {
	return MSG_SERVER_ACCEPT
}

func (*MsgServerAccept) EncodeSize() int {
	return 5
}

func (p *MsgServerAccept) Encode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint32(data[0:4], p.ID)
	data[4] = p.XorKey
	return 5, nil
}

func (p *MsgServerAccept) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	p.ID = binary.LittleEndian.Uint32(data[0:4])
	p.XorKey = data[4]
	return 5, nil
}

type MsgClientAccept struct {
	PlayerName  string // 0-65
	PlayerClass byte   // 66
	IsFemale    byte   // 67
	Unk70       [29]byte
	Screen      image.Point // 97-104
	Serial      string      // 105-126
	Unk129      [26]byte
}

func (*MsgClientAccept) NetOp() Op {
	return MSG_CLIENT_ACCEPT
}

func (*MsgClientAccept) EncodeSize() int {
	return 153
}

func (p *MsgClientAccept) Encode(data []byte) (int, error) {
	if len(data) < 153 {
		return 0, io.ErrShortBuffer
	}
	binenc.CStringSet16(data[0:66], p.PlayerName)
	data[66] = p.PlayerClass
	data[67] = p.IsFemale
	copy(data[68:97], p.Unk70[:])
	binary.LittleEndian.PutUint32(data[97:101], uint32(p.Screen.X))
	binary.LittleEndian.PutUint32(data[101:105], uint32(p.Screen.Y))
	binenc.CStringSet(data[105:127], p.Serial)
	copy(data[127:153], p.Unk129[:])
	return 153, nil
}

func (p *MsgClientAccept) Decode(data []byte) (int, error) {
	if len(data) < 153 {
		return 0, io.ErrUnexpectedEOF
	}
	p.PlayerName = binenc.CString16(data[0:66])
	p.PlayerClass = data[66]
	p.IsFemale = data[67]
	copy(p.Unk70[:], data[68:99])
	p.Screen.X = int(binary.LittleEndian.Uint32(data[97:101]))
	p.Screen.Y = int(binary.LittleEndian.Uint32(data[101:105]))
	p.Serial = binenc.CString(data[105:127])
	copy(p.Unk129[:], data[127:153])
	return 153, nil
}
