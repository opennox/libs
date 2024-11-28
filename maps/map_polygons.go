package maps

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/opennox/opennox-lib/binenc"
	"github.com/opennox/opennox-lib/types"
)

func init() {
	RegisterSection(&Polygons{})
}

type PolygonPoint struct {
	ID  uint32
	Pos types.Pointf
}

func (w *PolygonPoint) EncodingSize() int {
	return 12
}

func (w *PolygonPoint) MarshalBinary() ([]byte, error) {
	data := make([]byte, 12)
	binary.LittleEndian.PutUint32(data[0:], w.ID)
	binary.LittleEndian.PutUint32(data[4:], math.Float32bits(w.Pos.X))
	binary.LittleEndian.PutUint32(data[8:], math.Float32bits(w.Pos.Y))
	return data, nil
}

func (w *PolygonPoint) Decode(r *binenc.Reader) error {
	var ok bool
	w.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Pos.X, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Pos.Y, ok = r.ReadF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type Polygon struct {
	Name         string
	AmbientLight types.RGB
	MinimapGroup byte
	Points       []uint32
	PlayerEnter  *ScriptHandler
	MonsterEnter *ScriptHandler
	Flags        uint32
}

func (w *Polygon) EncodingSize() int {
	return 1 + len(w.Name) + 3 + 1 + 2 + 4*len(w.Points) + w.PlayerEnter.EncodingSize() + w.MonsterEnter.EncodingSize() + 4
}

func (w *Polygon) MarshalBinary(vers uint16) ([]byte, error) {
	if len(w.Name) > 0xff {
		return nil, fmt.Errorf("pilygon name is too long: %q", w.Name)
	}
	data := make([]byte, 0, w.EncodingSize())
	data = append(data, uint8(len(w.Name)))
	data = append(data, w.Name...)
	data = append(data, w.AmbientLight.R, w.AmbientLight.G, w.AmbientLight.B)
	data = append(data, w.MinimapGroup)
	data = binary.LittleEndian.AppendUint16(data, uint16(len(w.Points)))
	for _, p := range w.Points {
		data = binary.LittleEndian.AppendUint32(data, p)
	}
	if vers >= 2 {
		b, err := w.PlayerEnter.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)

		b, err = w.MonsterEnter.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	if vers >= 4 {
		data = binary.LittleEndian.AppendUint32(data, w.Flags)
	}
	return data, nil
}

func (w *Polygon) Decode(r *binenc.Reader, vers uint16) error {
	*w = Polygon{}
	var ok bool
	w.Name, ok = r.ReadString8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	rgb, ok := r.ReadU24()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.AmbientLight = types.RGB{R: rgb[0], G: rgb[1], B: rgb[2]}
	w.MinimapGroup, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	var sz1 uint16
	sz1, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Points = make([]uint32, 0, sz1)
	for i := 0; i < int(sz1); i++ {
		id, ok := r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
		w.Points = append(w.Points, id)
	}
	if vers >= 2 {
		w.PlayerEnter = new(ScriptHandler)
		if err := w.PlayerEnter.Decode(r); err != nil {
			return err
		}
		w.MonsterEnter = new(ScriptHandler)
		if err := w.MonsterEnter.Decode(r); err != nil {
			return err
		}
	}
	if vers >= 4 {
		w.Flags, ok = r.ReadU32()
		if !ok {
			return io.ErrUnexpectedEOF
		}
	}
	return nil
}

type Polygons struct {
	Vers     uint16
	Points   []PolygonPoint
	Polygons []Polygon
}

func (*Polygons) MapSection() string {
	return "Polygons"
}

func (sect *Polygons) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, 2+4+12*len(sect.Points)+4)
	data = binary.LittleEndian.AppendUint16(data, sect.Vers)
	data = binary.LittleEndian.AppendUint32(data, uint32(len(sect.Points)))
	for _, p := range sect.Points {
		b, err := p.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	data = binary.LittleEndian.AppendUint32(data, uint32(len(sect.Polygons)))
	for _, p := range sect.Polygons {
		b, err := p.MarshalBinary(sect.Vers)
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *Polygons) Decode(r *binenc.Reader) error {
	*sect = Polygons{}

	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers > 4 {
		return fmt.Errorf("unsupported polygons data version: %d", vers)
	}
	sect.Vers = vers

	sz, ok := r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Points = make([]PolygonPoint, 0, sz)
	for i := 0; i < int(sz); i++ {
		var p PolygonPoint
		if err := p.Decode(r); err != nil {
			return err
		}
		sect.Points = append(sect.Points, p)
	}
	sz, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Polygons = make([]Polygon, 0, sz)
	for i := 0; i < int(sz); i++ {
		var p Polygon
		if err := p.Decode(r, vers); err != nil {
			return err
		}
		sect.Polygons = append(sect.Polygons, p)
	}
	return nil
}

func (sect *Polygons) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}
