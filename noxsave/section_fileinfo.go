package noxsave

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/types"
)

func init() {
	Register(&FileInfo{})
}

type PlayerInfo struct {
	Hair     types.RGB
	Skin     types.RGB
	Mustache types.RGB
	Goatee   types.RGB
	Beard    types.RGB
	Pants    byte
	Shirt1   byte
	Shirt2   byte
	Shoes1   byte
	Shoes2   byte
	Name     string
	Class    byte
	IsFemale byte
}

func (p *PlayerInfo) EncodeSize() int {
	return 21 + len(p.Name)*2 + 2
}

func (p *PlayerInfo) Encode(data []byte) (int, error) {
	if len(data) < p.EncodeSize() {
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

	off := 20
	data[off] = byte(len(p.Name))
	off++
	off += binenc.CStringSet16(data[off:], p.Name)
	data[off+0] = p.Class
	data[off+1] = p.IsFemale
	off += 2
	return off, nil
}

func (p *PlayerInfo) Decode(data []byte) (int, error) {
	left := data
	if len(left) < 21 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Hair.R = left[0]
	p.Hair.G = left[1]
	p.Hair.B = left[2]

	p.Skin.R = left[3]
	p.Skin.G = left[4]
	p.Skin.B = left[5]

	p.Mustache.R = left[6]
	p.Mustache.G = left[7]
	p.Mustache.B = left[8]

	p.Goatee.R = left[9]
	p.Goatee.G = left[10]
	p.Goatee.B = left[11]

	p.Beard.R = left[12]
	p.Beard.G = left[13]
	p.Beard.B = left[14]

	p.Pants = left[15]
	p.Shirt1 = left[16]
	p.Shirt2 = left[17]
	p.Shoes1 = left[18]
	p.Shoes2 = left[19]

	sz := int(left[20])
	left = left[21:]
	if len(left) < sz*2+2 {
		return 0, io.ErrUnexpectedEOF
	}
	p.Name = binenc.CString16(left[:sz*2])
	left = left[sz*2:]
	p.Class = left[0]
	p.IsFemale = left[1]
	left = left[2:]
	return len(data) - len(left), nil
}

type FileInfo struct {
	Val0    uint32
	Path    string
	Val2    string
	Time    SystemTime
	Player  PlayerInfo
	Val3    byte
	MapName string
	Val4    byte
}

func (*FileInfo) SectionID() SectionID {
	return SectFileInfo
}

func (s *FileInfo) EncodeSize() int {
	//TODO implement me
	panic("implement me")
}

func (s *FileInfo) Encode(data []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FileInfo) Decode(data []byte) (int, error) {
	left := data
	if len(left) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(left[:2])
	left = left[2:]
	if vers > 12 {
		return 0, fmt.Errorf("unsupported section version: %d", vers)
	}
	*s = FileInfo{}

	if len(left) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	s.Val0 = binary.LittleEndian.Uint32(left[:4])
	left = left[4:]

	if len(left) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	sz := int(binary.LittleEndian.Uint16(left[:2]))
	left = left[2:]
	if len(left) < sz {
		return 0, io.ErrUnexpectedEOF
	}
	s.Path = binenc.CString(left[:sz])
	left = left[sz:]

	if len(left) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	sz = int(left[0])
	left = left[1:]
	if len(left) < sz {
		return 0, io.ErrUnexpectedEOF
	}
	s.Val2 = binenc.CString(left[:sz])
	left = left[sz:]

	n, err := s.Time.Decode(left)
	if err != nil {
		return 0, err
	}
	left = left[n:]

	n, err = s.Player.Decode(left)
	if err != nil {
		return 0, err
	}
	left = left[n:]
	if len(left) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	s.Val3 = left[0]
	left = left[1:]

	if vers < 11 {
		return len(data) - len(left), nil
	}
	if len(left) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	sz = int(left[0])
	left = left[1:]
	if len(left) < sz {
		return 0, io.ErrUnexpectedEOF
	}
	s.MapName = binenc.CString(left[:sz])
	left = left[sz:]
	if vers < 12 {
		return len(data) - len(left), nil
	}
	s.Val4 = left[0]
	left = left[1:]
	return len(data) - len(left), nil
}
