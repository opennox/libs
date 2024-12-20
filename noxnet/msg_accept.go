package noxnet

import (
	"encoding/binary"
	"image"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgAccept{}, false)
	netmsg.Register(&MsgServerAccept{}, false)
	netmsg.Register(&MsgClientAccept{}, false)
}

type MsgAccept struct {
	ID byte
}

func (*MsgAccept) NetOp() netmsg.Op {
	return netmsg.MSG_ACCEPTED
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

func (*MsgServerAccept) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_ACCEPT
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
	PlayerInfo
	Screen image.Point // 97-104
	Serial string      // 105-126
	Unk129 [26]byte
}

func (*MsgClientAccept) NetOp() netmsg.Op {
	return netmsg.MSG_CLIENT_ACCEPT
}

func (*MsgClientAccept) EncodeSize() int {
	return 153
}

func (p *MsgClientAccept) Encode(data []byte) (int, error) {
	if len(data) < 153 {
		return 0, io.ErrShortBuffer
	}
	_, err := p.PlayerInfo.Encode(data[0:97])
	if err != nil {
		return 0, err
	}
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
	_, err := p.PlayerInfo.Decode(data[0:97])
	if err != nil {
		return 0, err
	}
	p.Screen.X = int(binary.LittleEndian.Uint32(data[97:101]))
	p.Screen.Y = int(binary.LittleEndian.Uint32(data[101:105]))
	p.Serial = binenc.CString(data[105:127])
	copy(p.Unk129[:], data[127:153])
	return 153, nil
}
