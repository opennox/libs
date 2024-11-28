package xfer

import (
	"fmt"
	"io"

	"github.com/opennox/opennox-lib/binenc"
)

func init() {
	Register(&Weapon{})
}

type Weapon struct {
	Vers      uint16
	Object    Object
	Modifiers []*Modifier
	Val4      byte
	Val5      byte
	Health    uint16
	Val13     byte
	Val14     uint32
	Sub       []Xfer
}

func (*Weapon) XferType() Type {
	return "WeaponXfer"
}

func (x *Weapon) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Weapon{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 64 {
		return fmt.Errorf("unsupported weapon xfer vers: %v", x.Vers)
	}
	err := x.Object.DecodeXfer(reg, x.Vers, r)
	if err != nil {
		return err
	}
	if x.Vers < 11 {
		return nil
	}
	for i := 0; i < 4; i++ {
		var m Modifier
		m.Name, ok = r.ReadString8()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		x.Modifiers = append(x.Modifiers, &m)
	}
	if x.Vers > 41 {
		// FIXME: if class is wand, but it's not a wooden staff - read charge
		if x.Vers >= 61 {
			x.Health, ok = r.ReadU16()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		}
		if x.Vers >= 42 {
			x.Health, ok = r.ReadU16()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		}
		if x.Vers == 63 {
			x.Val13, ok = r.ReadU8()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		} else if x.Vers >= 64 {
			x.Val14, ok = r.ReadU32()
			if !ok {
				return io.ErrUnexpectedEOF
			}
		}
	}
	x.Sub, err = xferReadSub(reg, x.Vers, int(x.Object.SubN), r)
	return err
}
