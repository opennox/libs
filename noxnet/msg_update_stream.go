package noxnet

import (
	"encoding/binary"
	"image"
	"io"

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

func (*MsgUpdateStream) EncodeSize() int {
	panic("TODO")
}

func (p *MsgUpdateStream) Encode(data []byte) (int, error) {
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
	if alias != 0xff {
		return UpdateAlias(alias), 1, nil
	}
	if len(data) < 4 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	id := binary.LittleEndian.Uint16(data[1:3])
	typ := binary.LittleEndian.Uint16(data[3:5])
	return UpdateObjectID{ID: id, Type: typ}, 5, nil
}

type UpdateID interface {
	isUpdateID()
}

type UpdateAlias byte

func (UpdateAlias) isUpdateID() {}

type UpdateObjectID struct {
	ID   uint16
	Type uint16
}

func (UpdateObjectID) isUpdateID() {}

type ObjectUpdate struct {
	ID      UpdateID
	Pos     image.Point
	Complex *ComplexObjectUpdate
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
	alias := left[0]
	left = left[1:]
	rel := true
	if alias == 0 {
		alias = left[0]
		left = left[1:]
		if alias == 0 && left[1] == 0 {
			return 3, io.EOF
		}
		rel = false
	}
	isComplex := false
	if alias != 0xff {
		p.ID = UpdateAlias(alias)
		// FIXME: cannot check if it's a complex object without a map
	} else {
		id := binary.LittleEndian.Uint16(left[0:2])
		typ := binary.LittleEndian.Uint16(left[2:4])
		left = left[4:]
		p.ID = UpdateObjectID{ID: id, Type: typ}
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
	data[0] = byte(m.Alias)
	binary.LittleEndian.PutUint16(data[1:3], m.ID.ID)
	binary.LittleEndian.PutUint16(data[3:5], m.ID.Type)
	binary.LittleEndian.PutUint32(data[5:9], uint32(m.Deadline))
	return 9, nil
}

func (m *MsgNewAlias) Decode(data []byte) (int, error) {
	if len(data) < 9 {
		return 0, io.ErrUnexpectedEOF
	}
	m.Alias = UpdateAlias(data[0])
	id := binary.LittleEndian.Uint16(data[1:3])
	typ := binary.LittleEndian.Uint16(data[3:5])
	m.ID = UpdateObjectID{
		ID: id, Type: typ,
	}
	m.Deadline = Timestamp(binary.LittleEndian.Uint32(data[5:9]))
	return 9, nil
}
