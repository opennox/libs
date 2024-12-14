package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

func init() {
	netmsg.Register(&MsgNewPlayer{}, false)
	netmsg.Register(&MsgPlayerRespawn{}, false)
}

type PlayerInfo struct {
	PlayerName     string  // 0-49
	Unk50          uint32  // 50-53
	Unk54          uint32  // 54-57
	Unk58          uint32  // 58-61
	Unk62          uint32  // 62-65
	PlayerClass    byte    // 66
	IsFemale       byte    // 67
	Unk68          uint16  // 68-69
	Unk70          byte    // 70
	Unk71          uint16  // 71-72
	Unk73          byte    // 73
	Unk74          uint16  // 74-75
	Unk76          byte    // 76
	Unk77          uint16  // 77-78
	Unk79          byte    // 79
	Unk80          uint16  // 80-81
	Unk82          byte    // 82
	Unk83          [6]byte // 83-88
	PlayerNameSuff string  // 89-96
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
	binary.LittleEndian.PutUint16(data[68:70], p.Unk68)
	data[70] = p.Unk70
	binary.LittleEndian.PutUint16(data[71:73], p.Unk71)
	data[73] = p.Unk73
	binary.LittleEndian.PutUint16(data[74:76], p.Unk74)
	data[76] = p.Unk76
	binary.LittleEndian.PutUint16(data[77:79], p.Unk77)
	data[79] = p.Unk79
	binary.LittleEndian.PutUint16(data[80:82], p.Unk80)
	data[82] = p.Unk82
	copy(data[83:89], p.Unk83[:])
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
	p.Unk68 = binary.LittleEndian.Uint16(data[68:70])
	p.Unk70 = data[70]
	p.Unk71 = binary.LittleEndian.Uint16(data[71:73])
	p.Unk73 = data[73]
	p.Unk74 = binary.LittleEndian.Uint16(data[74:76])
	p.Unk76 = data[76]
	p.Unk77 = binary.LittleEndian.Uint16(data[77:79])
	p.Unk79 = data[79]
	p.Unk80 = binary.LittleEndian.Uint16(data[80:82])
	p.Unk82 = data[82]
	copy(p.Unk83[:], data[83:89])
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
