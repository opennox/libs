package xfer

import (
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/types"
)

const DefaultType = Type("DefaultXfer")

func init() {
	Register(&Default{})
}

type ScriptHandler struct {
	Vers uint16
	Func string
	Val2 uint32
}

func (x *ScriptHandler) DecodeXfer(_ ObjectRegistry, r *binenc.Reader) error {
	*x = ScriptHandler{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 1 {
		return fmt.Errorf("unsupported script handler xfer vers: %v", x.Vers)
	}
	x.Func, ok = r.ReadString32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Val2, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type Object struct {
	Vers      uint16
	Extent    uint32
	ID        uint32
	Pos       types.Pointf
	Val5      byte
	Flags     uint32
	Name      string
	Team      byte
	SubN      byte
	Owned     []uint32
	Anim      uint32
	Handler11 *ScriptHandler
	DeadFrame int32
}

func (x *Object) DecodeXfer(reg ObjectRegistry, gvers uint16, r *binenc.Reader) error {
	*x = Object{}
	var ok bool
	if gvers >= 40 {
		x.Vers, ok = r.ReadU16()
		if !ok {
			return io.ErrUnexpectedEOF
		} else if x.Vers > 64 {
			return fmt.Errorf("unsupported object xfer vers: %v", x.Vers)
		}
	}
	if gvers < 40 || x.Vers < 61 {
		return x.decodeOld(gvers, r)
	}
	x.Extent, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Pos.X, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Pos.Y, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Val5, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if x.Val5 == 0 {
		return nil
	}
	x.Flags, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Name, ok = r.ReadString8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Team, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.SubN, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sz, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Owned = make([]uint32, sz)
	for i := range x.Owned {
		x.Owned[i], ok = r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	x.Anim, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if x.Vers >= 63 {
		x.Handler11 = new(ScriptHandler)
		err := x.Handler11.DecodeXfer(reg, r)
		if err != nil {
			return err
		}
		if x.Vers >= 64 {
			x.DeadFrame, ok = r.ReadI32()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		}
	}
	return nil
}

func (x *Object) decodeOld(gvers uint16, r *binenc.Reader) error {
	var ok bool
	x.Extent, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	x.Flags, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	if gvers < 40 || x.Vers < 4 {
		posX, ok := r.ReadI32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		posY, ok := r.ReadI32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		x.Pos = types.Pointf{X: float32(posX), Y: float32(posY)}
	} else {
		posX, ok := r.ReadF32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		posY, ok := r.ReadF32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		x.Pos = types.Pointf{X: posX, Y: posY}
	}
	if gvers >= 10 {
		x.Name, ok = r.ReadString8()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	if gvers >= 20 {
		x.Team, ok = r.ReadU8()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	if gvers >= 30 {
		x.SubN, ok = r.ReadU8()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	if gvers < 40 {
		return nil
	}
	x.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	// FIXME
	return fmt.Errorf("old object xfer format is not fully supported yet")
}

func xferReadSubOne(reg ObjectRegistry, gvers uint16, r *binenc.Reader) (Xfer, error) {
	if gvers < 60 {
		typ, ok := r.ReadString8()
		if !ok {
			return nil, io.ErrUnexpectedEOF
		}
		x, err := DecodeByObjectType(reg, typ, r)
		return x, err
	} else {
		typ, ok := r.ReadU16()
		if !ok {
			return nil, io.ErrUnexpectedEOF
		}
		x, err := DecodeByObjectTypeID(reg, int(typ), r)
		return x, err
	}
}

func xferReadSub(reg ObjectRegistry, gvers uint16, sz int, r *binenc.Reader) ([]Xfer, error) {
	var out []Xfer
	for i := 0; i < sz; i++ {
		x, err := xferReadSubOne(reg, gvers, r)
		if x != nil {
			out = append(out, x)
		}
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

type Default struct {
	Vers   uint16
	Object Object
	Sub    []Xfer
}

func (*Default) XferType() Type {
	return DefaultType
}

func (x *Default) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Default{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 60 {
		return fmt.Errorf("unsupported default xfer vers: %v", x.Vers)
	}
	err := x.Object.DecodeXfer(reg, x.Vers, r)
	if err != nil {
		return err
	}
	x.Sub, err = xferReadSub(reg, x.Vers, int(x.Object.SubN), r)
	return err
}
