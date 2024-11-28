package maps

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterSection(&WindowWalls{})
}

type WindowWall struct {
	Pos image.Point
}

func (w *WindowWall) MarshalBinary() ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data[0:], uint32(int32(w.Pos.X)))
	binary.LittleEndian.PutUint32(data[4:], uint32(int32(w.Pos.Y)))
	return data, nil
}

func (w *WindowWall) Decode(r *binenc.Reader) error {
	var ok bool
	w.Pos, ok = r.ReadPointI32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type WindowWalls struct {
	Walls []WindowWall
}

func (*WindowWalls) MapSection() string {
	return "WindowWalls"
}

func (sect *WindowWalls) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4, 4+8*len(sect.Walls))
	binary.LittleEndian.PutUint16(data[0:], 2)
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

func (sect *WindowWalls) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *WindowWalls) Decode(r *binenc.Reader) error {
	*sect = WindowWalls{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 2 {
		return fmt.Errorf("unsupported version of window walls section: %d", vers)
	}
	n, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Walls = make([]WindowWall, 0, n)
	for i := 0; i < int(n); i++ {
		var w WindowWall
		if err := w.Decode(r); err != nil {
			return err
		}
		sect.Walls = append(sect.Walls, w)
	}
	return nil
}
