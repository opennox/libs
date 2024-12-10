package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgServerTryJoin{}, false)
	netmsg.Register(&MsgJoinOK{}, false)
	netmsg.Register(&MsgConnect{}, false)
}

type MsgServerTryJoin struct {
	Unk0        byte     // 0
	PlayerName  string   // 1-50
	PlayerClass byte     // 51
	PlayerLevel byte     // 52
	Serial      string   // 53-76
	Version     uint32   // 77-80
	Team        uint32   // 81-84
	Unk85       [10]byte // 85-88
	Unk95       byte     // 95
	Unk96       byte     // 96
}

func (*MsgServerTryJoin) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_TRY_JOIN
}

func (*MsgServerTryJoin) EncodeSize() int {
	return 97
}

func (p *MsgServerTryJoin) Encode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Unk0
	binenc.CStringSet16(data[1:51], p.PlayerName)
	binenc.CStringSet(data[53:77], p.Serial)
	binary.LittleEndian.PutUint32(data[77:81], p.Version)
	binary.LittleEndian.PutUint32(data[81:85], p.Team)
	copy(data[85:95], p.Unk85[:])
	data[95] = p.Unk95
	data[96] = p.Unk96
	return 97, nil
}

func (p *MsgServerTryJoin) Decode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Unk0 = data[0]
	p.PlayerName = binenc.CString16(data[1:53])
	p.Serial = binenc.CString(data[53:77])
	p.Version = binary.LittleEndian.Uint32(data[77:81])
	p.Team = binary.LittleEndian.Uint32(data[81:85])
	copy(p.Unk85[:], data[85:95])
	p.Unk95 = data[95]
	p.Unk96 = data[96]
	return 97, nil
}

type MsgJoinOK struct {
}

func (*MsgJoinOK) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_JOIN_OK
}

func (*MsgJoinOK) EncodeSize() int {
	return 0
}

func (p *MsgJoinOK) Encode(data []byte) (int, error) {
	return 0, nil
}

func (p *MsgJoinOK) Decode(data []byte) (int, error) {
	return 0, nil
}

type MsgConnect struct {
}

func (*MsgConnect) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_CONNECT
}

func (*MsgConnect) EncodeSize() int {
	return 0
}

func (p *MsgConnect) Encode(data []byte) (int, error) {
	return 0, nil
}

func (p *MsgConnect) Decode(data []byte) (int, error) {
	return 0, nil
}
