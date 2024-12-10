package discover

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgDiscover{}, false)
	netmsg.Register(&MsgServerInfo{}, true)
}

type MsgDiscover struct {
	Unk0  [5]byte
	Token uint32
}

func (*MsgDiscover) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_DISCOVER
}

func (*MsgDiscover) EncodeSize() int {
	return 9
}

func (p *MsgDiscover) Encode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrShortBuffer
	}
	copy(data[0:5], p.Unk0[:])
	binary.LittleEndian.PutUint32(data[5:], p.Token)
	return 9, nil
}

func (p *MsgDiscover) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	copy(p.Unk0[:], data[0:5])
	p.Token = binary.LittleEndian.Uint32(data[5:])
	return 9, nil
}

type MsgServerInfo struct {
	PlayersCur byte     // 0
	PlayersMax byte     // 1
	Unk2       [5]byte  // 2-6
	MapName    string   // 7-15
	Status1    byte     // 16
	Status2    byte     // 17
	Unk19      [7]byte  // 18-24
	Flags      uint16   // 25-26
	Unk27      [2]byte  // 27-28
	Unk29      [8]byte  // 29-36
	Unk37      [4]byte  // 37-40
	Token      uint32   // 41-44
	Unk45      [20]byte // 45-64
	Unk65      [4]byte  // 65-68
	ServerName string   // 69+
}

func (*MsgServerInfo) NetOp() netmsg.Op {
	return netmsg.MSG_SERVER_INFO
}

func (p *MsgServerInfo) EncodeSize() int {
	return 69 + len(p.ServerName) + 1
}

func (p *MsgServerInfo) Encode(data []byte) (int, error) {
	sz := p.EncodeSize()
	if len(data) < sz {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.PlayersCur
	data[1] = p.PlayersMax
	copy(data[2:7], p.Unk2[:])
	binenc.CStringSet0(data[7:16], p.MapName)
	data[16] = p.Status1
	data[17] = p.Status2
	copy(data[18:25], p.Unk19[:])
	binary.LittleEndian.PutUint16(data[25:], p.Flags)
	copy(data[27:29], p.Unk27[:])
	copy(data[29:37], p.Unk29[:])
	copy(data[37:41], p.Unk37[:])
	binary.LittleEndian.PutUint32(data[41:], p.Token)
	copy(data[45:65], p.Unk45[:])
	copy(data[65:69], p.Unk65[:])
	binenc.CStringSet0(data[69:], p.ServerName)
	return sz, nil
}

func (p *MsgServerInfo) Decode(data []byte) (int, error) {
	if len(data) < 70 {
		return 0, io.ErrUnexpectedEOF
	}
	p.PlayersCur = data[0]
	p.PlayersMax = data[1]
	copy(p.Unk2[:], data[2:7])
	p.MapName = binenc.CString(data[7:16])
	p.Status1 = data[16]
	p.Status2 = data[17]
	copy(p.Unk19[:], data[18:25])
	p.Flags = binary.LittleEndian.Uint16(data[25:])
	copy(p.Unk27[:], data[27:29])
	copy(p.Unk29[:], data[29:37])
	copy(p.Unk37[:], data[37:41])
	p.Token = binary.LittleEndian.Uint32(data[41:])
	copy(p.Unk45[:], data[45:65])
	copy(p.Unk65[:], data[65:69])
	p.ServerName = binenc.CString(data[69:])
	return len(data), nil
}
