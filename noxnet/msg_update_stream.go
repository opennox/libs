package noxnet

import (
	"encoding/binary"
	"errors"
	"image"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/noxnet/netmsg"
)

const decodeUpdateStream = false

func init() {
	if decodeUpdateStream {
		netmsg.Register(&MsgUpdateStream{}, true)
	}
	netmsg.Register(&MsgNewAlias{}, false)
}

type MsgUpdateStream struct {
	ID      UpdateID
	Pos     image.Point
	Flags   byte
	Unk4    byte // anim frame?
	Unk5    byte
	Objects []ObjectUpdate
}

func (*MsgUpdateStream) NetOp() netmsg.Op {
	return netmsg.MSG_UPDATE_STREAM
}

func (p *MsgUpdateStream) EncodeSize() int {
	n := 0
	n += p.ID.EncodeSize()
	n += 5
	if p.Flags&0x80 != 0 {
		n++
	}
	n++
	for _, obj := range p.Objects {
		n += obj.EncodeSize()
	}
	n += 3 // FIXME?
	return n
}

func (p *MsgUpdateStream) Encode(data []byte) (int, error) {
	if len(data) < p.EncodeSize() {
		return 0, io.ErrShortBuffer
	}
	off := 0
	n, err := p.ID.Encode(data[off:])
	if err != nil {
		return n, err
	}
	off += n

	binary.BigEndian.PutUint16(data[off+0:], uint16(p.Pos.X))
	binary.BigEndian.PutUint16(data[off+2:], uint16(p.Pos.Y))
	data[off+4] = p.Flags
	off += 5

	if p.Flags&0x80 != 0 {
		data[off] = p.Unk4
		off++
	}

	if len(p.Objects) == 0 {
		data[off+0] = 0
		data[off+1] = 0
		data[off+3] = 0
		off += 3
		return off, nil
	}
	panic("TODO")
}

func (p *MsgUpdateStream) decodeHeader(data []byte) (int, error) {
	left := data

	id, n, err := decodeUpdateID(left)
	if err != nil {
		return 0, err
	}
	left = left[n:]
	p.ID = id

	p.Pos.X = int(binary.LittleEndian.Uint16(left[0:2]))
	p.Pos.Y = int(binary.LittleEndian.Uint16(left[2:4]))
	p.Flags = left[4]
	left = left[5:]
	p.Unk4 = 0
	if p.Flags&0x80 != 0 {
		p.Unk4 = left[0]
		left = left[1:]
	}
	p.Unk5 = left[0]
	left = left[1:]
	return len(data) - len(left), nil
}

func (p *MsgUpdateStream) Decode(data []byte) (int, error) {
	left := data

	n, err := p.decodeHeader(left)
	if err != nil {
		return 0, err
	}
	left = left[n:]

	p.Objects = nil
	for len(left) != 0 {
		var u ObjectUpdate
		n, err = u.Decode(left, p.Pos)
		if err == io.EOF {
			left = left[n:]
			break
		} else if err != nil {
			return 0, err
		}
		left = left[n:]
		p.Objects = append(p.Objects, u)
	}
	return len(data) - len(left), nil
}

func decodeUpdateID(data []byte) (UpdateID, int, error) {
	if len(data) < 1 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	alias := data[0]
	var u UpdateID
	if alias != 0xff {
		u = &UpdateAlias{}
	} else {
		u = &UpdateObjectID{}
	}
	n, err := u.Decode(data)
	if err != nil {
		return nil, 0, err
	}
	return u, n, err
}

type UpdateID interface {
	isUpdateID()
	binenc.Encoded
}

type UpdateAlias struct {
	Alias byte
}

func (*UpdateAlias) isUpdateID() {}

func (u *UpdateAlias) EncodeSize() int {
	return 1
}

func (u *UpdateAlias) Encode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrShortBuffer
	}
	if u.Alias == 0xff {
		return 0, errors.New("not an alias")
	}
	data[0] = u.Alias
	return 1, nil
}

func (u *UpdateAlias) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	alias := data[0]
	if alias == 0xff {
		return 0, errors.New("not an alias")
	}
	u.Alias = alias
	return 1, nil
}

type UpdateObjectID struct {
	ID   NetCode
	Type uint16
}

func (*UpdateObjectID) isUpdateID() {}

func (u *UpdateObjectID) EncodeSize() int {
	return 5
}

func (u *UpdateObjectID) Encode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrShortBuffer
	}
	data[0] = 0xff
	binary.LittleEndian.PutUint16(data[1:3], uint16(u.ID))
	binary.LittleEndian.PutUint16(data[3:5], u.Type)
	return 5, nil
}

func (u *UpdateObjectID) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, io.ErrUnexpectedEOF
	}
	alias := data[0]
	if alias != 0xff {
		return 0, errors.New("not an object update")
	}
	u.ID = NetCode(binary.LittleEndian.Uint16(data[1:3]))
	u.Type = binary.LittleEndian.Uint16(data[3:5])
	return 5, nil
}

type ObjectUpdate struct {
	ID      UpdateID
	Pos     image.Point
	Complex *ComplexObjectUpdate
}

func (*ObjectUpdate) EncodeSize() int {
	panic("TODO")
}

type ComplexObjectUpdate struct {
	Unk0 byte
	Unk1 byte
	Unk2 byte
}

func (p *ObjectUpdate) Decode(data []byte, par image.Point) (int, error) {
	left := data
	if len(left) < 3 {
		return 0, io.ErrUnexpectedEOF
	}
	if data[0] == 0 && data[1] == 0 && data[2] == 0 {
		return 3, io.EOF
	}
	alias := left[0]
	left = left[1:]
	rel := true
	if alias == 0 {
		alias = left[0]
		left = left[1:]
		rel = false
	}
	isComplex := false
	if alias != 0xff {
		p.ID = &UpdateAlias{alias}
		// FIXME: cannot check if it's a complex object without a map
	} else {
		id := NetCode(binary.LittleEndian.Uint16(left[0:2]))
		typ := binary.LittleEndian.Uint16(left[2:4])
		left = left[4:]
		p.ID = &UpdateObjectID{ID: id, Type: typ}
		isComplex = objectTypeIsComplex(typ)
	}
	if !rel {
		x := binary.LittleEndian.Uint16(left[0:2])
		y := binary.LittleEndian.Uint16(left[2:4])
		left = left[4:]
		p.Pos = image.Point{X: int(x), Y: int(y)}
	} else {
		dx := left[0]
		dy := left[1]
		left = left[2:]
		p.Pos = par.Add(image.Point{X: int(dx), Y: int(dy)})
	}
	if !isComplex {
		p.Complex = nil
		return len(data) - len(left), nil
	}
	unk0 := left[0]
	unk1 := left[1]
	unk2 := left[2]
	left = left[3:]
	p.Complex = &ComplexObjectUpdate{
		Unk0: unk0,
		Unk1: unk1,
		Unk2: unk2,
	}
	return len(data) - len(left), nil
}

type MsgNewAlias struct {
	Alias    UpdateAlias
	ID       UpdateObjectID
	Deadline Timestamp
}

func (*MsgNewAlias) NetOp() netmsg.Op {
	return netmsg.MSG_NEW_ALIAS
}

func (*MsgNewAlias) EncodeSize() int {
	return 9
}

func (m *MsgNewAlias) Encode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrShortBuffer
	}
	data[0] = m.Alias.Alias
	binary.LittleEndian.PutUint16(data[1:3], uint16(m.ID.ID))
	binary.LittleEndian.PutUint16(data[3:5], m.ID.Type)
	binary.LittleEndian.PutUint32(data[5:9], uint32(m.Deadline))
	return 9, nil
}

func (m *MsgNewAlias) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Alias.Alias = data[0]
	id := NetCode(binary.LittleEndian.Uint16(data[1:3]))
	typ := binary.LittleEndian.Uint16(data[3:5])
	m.ID = UpdateObjectID{
		ID: id, Type: typ,
	}
	m.Deadline = Timestamp(binary.LittleEndian.Uint32(data[5:9]))
	return 9, nil
}
