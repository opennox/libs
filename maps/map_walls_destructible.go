package maps

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/opennox/opennox-lib/binenc"
)

func init() {
	RegisterSection(&DestructableWalls{})
}

type DestructableWall struct {
	Pos image.Point
}

func (w *DestructableWall) MarshalBinary() ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data[0:], uint32(int32(w.Pos.X)))
	binary.LittleEndian.PutUint32(data[4:], uint32(int32(w.Pos.Y)))
	return data, nil
}

func (w *DestructableWall) Decode(r *binenc.Reader) error {
	var ok bool
	w.Pos, ok = r.ReadPointI32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type DestructableWalls struct {
	Walls []DestructableWall
}

func (*DestructableWalls) MapSection() string {
	return "DestructableWalls"
}

func (sect *DestructableWalls) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4, 4+8*len(sect.Walls))
	binary.LittleEndian.PutUint16(data[0:], 1)
	binary.LittleEndian.PutUint16(data[2:], uint16(len(sect.Walls)))
	for _, w := range sect.Walls {
		b, err := w.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *DestructableWalls) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *DestructableWalls) Decode(r *binenc.Reader) error {
	*sect = DestructableWalls{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 1 {
		return fmt.Errorf("unsupported version of destructable walls section: %d", vers)
	}
	n, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Walls = make([]DestructableWall, 0, n)
	for i := 0; i < int(n); i++ {
		var w DestructableWall
		if err := w.Decode(r); err != nil {
			return err
		}
		sect.Walls = append(sect.Walls, w)
	}
	return nil
}
