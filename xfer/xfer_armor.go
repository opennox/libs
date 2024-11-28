package xfer

import (
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	Register(&Armor{})
}

type Modifier struct {
	Name string
}

type Armor struct {
	Vers      uint16
	Object    Object
	Modifiers []*Modifier
	Health    uint16
	Val13     byte
	Val14     uint32
	Sub       []Xfer
}

func (*Armor) XferType() Type {
	return "ArmorXfer"
}

func (x *Armor) DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error {
	*x = Armor{}
	var ok bool
	x.Vers, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x.Vers > 62 {
		return fmt.Errorf("unsupported armor xfer vers: %v", x.Vers)
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
	if x.Vers >= 41 {
		x.Health, ok = r.ReadU16()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	if x.Vers == 61 {
		x.Val13, ok = r.ReadU8()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	} else if x.Vers >= 62 {
		x.Val14, ok = r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	x.Sub, err = xferReadSub(reg, x.Vers, int(x.Object.SubN), r)
	return err
}
