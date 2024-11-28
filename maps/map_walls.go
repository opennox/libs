package maps

import (
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterSection(&WallMap{})
}

type WallPos struct {
	X, Y byte
}

type Wall struct {
	Pos      WallPos
	Dir      byte
	DirBit   byte
	Material byte
	Variant  byte
	Minimap  byte
	Modified byte
}

func (w *Wall) IsZero() bool {
	return *w == (Wall{})
}

func (w *Wall) MarshalBinary() ([]byte, error) {
	if w.IsZero() {
		return []byte{0xff}, nil
	}
	data := make([]byte, 7)
	data[0] = w.Pos.X
	data[1] = w.Pos.Y
	data[2] = w.Dir | w.DirBit
	data[3] = w.Material
	data[4] = w.Variant
	data[5] = w.Minimap
	data[6] = w.Modified
	return data, nil
}

func (w *Wall) Decode(r *binenc.Reader) error {
	*w = Wall{}
	x, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x == 0xff {
		return nil
	}
	y, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if x == 0xff {
		return nil
	}
	w.Pos = WallPos{X: x, Y: y}
	data, ok := r.ReadNext(5)
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Dir = data[0]
	w.DirBit = w.Dir & 0x80 // TODO: check in the engine
	w.Dir &= 0x7f
	w.Material = data[1]
	w.Variant = data[2]
	w.Minimap = data[3]
	w.Modified = data[4]
	return nil
}

type WallMap struct {
	Grid  GridData
	Walls []Wall
}

func (*WallMap) MapSection() string {
	return "WallMap"
}

func (sect *WallMap) MarshalBinary() ([]byte, error) {
	out := make([]byte, 18+7*len(sect.Walls)+1)
	grid, err := sect.Grid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	n := copy(out, grid)
	data := out[n:]
	for _, w := range sect.Walls {
		wdata, err := w.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(data, wdata)
		data = data[7:]
	}
	var end Wall
	wdata, err := end.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(data, wdata)
	data = data[1:]
	return out, nil
}

func (sect *WallMap) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *WallMap) Decode(r *binenc.Reader) error {
	*sect = WallMap{}
	if err := sect.Grid.Decode(r); err != nil {
		return err
	}
	for {
		var w Wall
		if err := w.Decode(r); err != nil {
			return err
		} else if w.IsZero() {
			if r.Remaining() > 0 {
				return fmt.Errorf("trailing wall data: [%d]", r.Remaining())
			}
			return nil
		}
		sect.Walls = append(sect.Walls, w)
	}
}
