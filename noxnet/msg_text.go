package noxnet

import (
	"encoding/binary"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterMessage(&MsgText{}, true)
}

const (
	TextTeam      = TextFlags(0x1)
	TextUTF8      = TextFlags(0x2)
	TextFlag4     = TextFlags(0x2)
	TextLocalized = TextFlags(0x8)
	TextNotice    = TextFlags(0x10)
	TextFlag20    = TextFlags(0x20)
	TextFlag40    = TextFlags(0x40)
	TextExt       = TextFlags(0x80) // unused in vanilla
)

type TextFlags byte

func (f TextFlags) Has(f2 TextFlags) bool {
	return f&f2 != 0
}

type MsgText struct {
	NetCode uint16
	Flags   TextFlags
	PosX    uint16
	PosY    uint16
	Size    byte
	Dur     uint16
	Data    []byte // max: 510
}

func (*MsgText) NetOp() Op {
	return MSG_TEXT_MESSAGE
}

func (p *MsgText) Text() string {
	if p.Flags.Has(TextUTF8) || p.Flags.Has(TextLocalized) {
		data, sz := p.Data, int(p.Size)
		if len(data) > sz {
			data = data[:sz]
		}
		return binenc.CString(data)
	}
	// UTF-16
	data, sz := p.Data, 2*int(p.Size)
	if len(data) > sz {
		data = data[:sz]
	}
	return binenc.CString16(data)
}

func (p *MsgText) Payload() []byte {
	if !p.Flags.Has(TextUTF8) && !p.Flags.Has(TextLocalized) {
		i := binenc.CLen16(p.Data)
		if i+2 < len(p.Data) {
			return p.Data[i+2:]
		}
	} else {
		i := binenc.CLen(p.Data)
		if i+1 < len(p.Data) {
			return p.Data[i+1:]
		}
	}
	return nil
}

func (p *MsgText) dataSize() int {
	if p.Size == 0 && len(p.Data) != 0 {
		return len(p.Data)
	}
	sz := int(p.Size)
	if !p.Flags.Has(TextUTF8) && !p.Flags.Has(TextLocalized) {
		sz *= 2
	}
	data := p.Data
	if len(data) > sz {
		data = data[:sz]
	}
	return len(data)
}

func (p *MsgText) EncodeSize() int {
	return 10 + p.dataSize()
}

func (p *MsgText) Encode(data []byte) (int, error) {
	dsz := p.dataSize()
	if len(data) < 10+dsz {
		return 0, io.ErrShortBuffer
	}
	binary.LittleEndian.PutUint16(data[0:2], p.NetCode)
	data[2] = byte(p.Flags)
	binary.LittleEndian.PutUint16(data[3:5], p.PosX)
	binary.LittleEndian.PutUint16(data[5:7], p.PosY)
	data[7] = p.Size
	binary.LittleEndian.PutUint16(data[8:10], p.Dur)
	copy(data[10:10+dsz], p.Data)
	return 10 + dsz, nil
}

func (p *MsgText) Decode(data []byte) (int, error) {
	if len(data) < 10 {
		return 0, io.ErrUnexpectedEOF
	}
	p.NetCode = binary.LittleEndian.Uint16(data[0:2])
	p.Flags = TextFlags(data[2])
	p.PosX = binary.LittleEndian.Uint16(data[3:5])
	p.PosY = binary.LittleEndian.Uint16(data[5:7])
	p.Size = data[7]
	p.Dur = binary.LittleEndian.Uint16(data[8:10])
	data = data[10:]
	sz := int(p.Size)
	if !p.Flags.Has(TextUTF8) && !p.Flags.Has(TextLocalized) {
		sz *= 2
	}
	if sz > len(data) {
		sz = len(data)
	}
	p.Data = make([]byte, sz)
	copy(p.Data, data)
	return 10 + sz, nil
}
