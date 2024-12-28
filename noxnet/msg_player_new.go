package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/types"
)

func init() {
	netmsg.Register(&MsgNewPlayer{}, false)
	netmsg.Register(&MsgPlayerRespawn{}, false)
}

type PlayerColors struct {
	Hair     types.RGB // 0-2
	Skin     types.RGB // 3-5
	Mustache types.RGB // 6-8
	Goatee   types.RGB // 9-11
	Beard    types.RGB // 12-14
	Pants    byte      // 15
	Shirt1   byte      // 16
	Shirt2   byte      // 17
	Shoes1   byte      // 18
	Shoes2   byte      // 19
}

func (*PlayerColors) EncodeSize() int {
	return 20
}

func (p *PlayerColors) Encode(data []byte) (int, error) {
	if len(data) < 20 {
		return 0, io.ErrShortBuffer
	}
	data[0] = p.Hair.R
	data[1] = p.Hair.G
	data[2] = p.Hair.B

	data[3] = p.Skin.R
	data[4] = p.Skin.G
	data[5] = p.Skin.B

	data[6] = p.Mustache.R
	data[7] = p.Mustache.G
	data[8] = p.Mustache.B

	data[9] = p.Goatee.R
	data[10] = p.Goatee.G
	data[11] = p.Goatee.B

	data[12] = p.Beard.R
	data[13] = p.Beard.G
	data[14] = p.Beard.B

	data[15] = p.Pants
	data[16] = p.Shirt1
	data[17] = p.Shirt2
	data[18] = p.Shoes1
	data[19] = p.Shoes2
	return 20, nil
}

func (p *PlayerColors) Decode(data []byte) (int, error) {
	if len(data) < 20 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Hair.R = data[0]
	p.Hair.G = data[1]
	p.Hair.B = data[2]

	p.Skin.R = data[3]
	p.Skin.G = data[4]
	p.Skin.B = data[5]

	p.Mustache.R = data[6]
	p.Mustache.G = data[7]
	p.Mustache.B = data[8]

	p.Goatee.R = data[9]
	p.Goatee.G = data[10]
	p.Goatee.B = data[11]

	p.Beard.R = data[12]
	p.Beard.G = data[13]
	p.Beard.B = data[14]

	p.Pants = data[15]
	p.Shirt1 = data[16]
	p.Shirt2 = data[17]
	p.Shoes1 = data[18]
	p.Shoes2 = data[19]
	return 20, nil
}

type PlayerInfo struct {
	PlayerName     string       // 0-49
	Unk50          uint32       // 50-53
	Unk54          uint32       // 54-57
	Unk58          uint32       // 58-61
	Unk62          uint32       // 62-65
	PlayerClass    byte         // 66
	IsFemale       byte         // 67
	Colors         PlayerColors // 68-87
	Unk88          byte         // 88
	PlayerNameSuff string       // 89-96
}

func (*PlayerInfo) EncodeSize() int {
	return 97
}

func (p *PlayerInfo) Encode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrShortBuffer
	}
	binenc.CStringSet16(data[0:50], p.PlayerName)
	binary.LittleEndian.PutUint32(data[50:54], p.Unk50)
	binary.LittleEndian.PutUint32(data[54:58], p.Unk54)
	binary.LittleEndian.PutUint32(data[58:62], p.Unk58)
	binary.LittleEndian.PutUint32(data[62:66], p.Unk62)
	data[66] = p.PlayerClass
	data[67] = p.IsFemale
	if _, err := p.Colors.Encode(data[68:88]); err != nil {
		return 0, err
	}
	data[88] = p.Unk88
	binenc.CStringSet16(data[89:97], p.PlayerNameSuff)
	return 97, nil
}

func (p *PlayerInfo) Decode(data []byte) (int, error) {
	if len(data) < 97 {
		return 0, io.ErrUnexpectedEOF
	}
	p.PlayerName = binenc.CString16(data[0:50])
	p.Unk50 = binary.LittleEndian.Uint32(data[50:54])
	p.Unk54 = binary.LittleEndian.Uint32(data[54:58])
	p.Unk58 = binary.LittleEndian.Uint32(data[58:62])
	p.Unk62 = binary.LittleEndian.Uint32(data[62:66])
	p.PlayerClass = data[66]
	p.IsFemale = data[67]
	if _, err := p.Colors.Decode(data[68:88]); err != nil {
		return 0, err
	}
	p.Unk88 = data[88]
	p.PlayerNameSuff = binenc.CString16(data[89:97])
	return 97, nil
}

type MsgNewPlayer struct {
	NetCode    NetCode       // 0-1
	PlayerInfo               // 2-98
	Lessons    uint16        // 99-100
	Unk101     uint16        // 101-102
	Armor      uint32        // 103-106
	Weapon     uint32        // 107-110
	Unk111     uint32        // 111-114
	Unk115     byte          // 115
	Unk116     byte          // 116
	Unk117     byte          // 117
	Unk118     binenc.String // 118-127
}

func (*MsgNewPlayer) NetOp() netmsg.Op {
	return netmsg.MSG_NEW_PLAYER
}

func (*MsgNewPlayer) EncodeSize() int {
	return 128
}

func (p *MsgNewPlayer) Encode(data []byte) (int, error) {
	if len(data) < 128 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], uint16(p.NetCode))
	_, err := p.PlayerInfo.Encode(data[2:99])
	if err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(data[99:101], p.Lessons)
	binary.LittleEndian.PutUint16(data[101:103], p.Unk101)
	binary.LittleEndian.PutUint32(data[103:107], p.Armor)
	binary.LittleEndian.PutUint32(data[107:111], p.Weapon)
	binary.LittleEndian.PutUint32(data[111:115], p.Unk111)
	data[115] = p.Unk115
	data[116] = p.Unk116
	data[117] = p.Unk117
	p.Unk118.Encode(data[118:128])
	return 128, nil
}

func (p *MsgNewPlayer) Decode(data []byte) (int, error) {
	if len(data) < 128 {
		return 0, io.ErrUnexpectedEOF
	}
	p.NetCode = NetCode(binary.LittleEndian.Uint16(data[0:2]))
	_, err := p.PlayerInfo.Decode(data[2:99])
	if err != nil {
		return 0, err
	}
	p.Lessons = binary.LittleEndian.Uint16(data[99:101])
	p.Unk101 = binary.LittleEndian.Uint16(data[101:103])
	p.Armor = binary.LittleEndian.Uint32(data[103:107])
	p.Weapon = binary.LittleEndian.Uint32(data[107:111])
	p.Unk111 = binary.LittleEndian.Uint32(data[111:115])
	p.Unk115 = data[115]
	p.Unk116 = data[116]
	p.Unk117 = data[117]
	p.Unk118.Decode(data[118:128])
	return 128, nil
}

type MsgPlayerRespawn struct {
	NetCode NetCode
	Unk2    uint32
	Unk6    byte
	Unk7    byte
}

func (*MsgPlayerRespawn) NetOp() netmsg.Op {
	return netmsg.MSG_PLAYER_RESPAWN
}

func (*MsgPlayerRespawn) EncodeSize() int {
	return 8
}

func (p *MsgPlayerRespawn) Encode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], uint16(p.NetCode))
	binary.LittleEndian.PutUint32(data[2:6], p.Unk2)
	data[6] = p.Unk6
	data[7] = p.Unk7
	return 8, nil
}

func (p *MsgPlayerRespawn) Decode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}
	p.NetCode = NetCode(binary.LittleEndian.Uint16(data[0:2]))
	p.Unk2 = binary.LittleEndian.Uint32(data[2:6])
	p.Unk6 = data[6]
	p.Unk7 = data[7]
	return 8, nil
}
