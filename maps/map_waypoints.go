package maps

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/opennox/libs/binenc"
	"github.com/opennox/libs/types"
)

func init() {
	RegisterSection(&Waypoints{})
}

type WaypointLink struct {
	ID    uint32
	Flags byte
}

func (w *WaypointLink) MarshalBinary() ([]byte, error) {
	data := make([]byte, 5)
	binary.LittleEndian.PutUint32(data[0:], w.ID)
	data[4] = w.Flags
	return data, nil
}

func (w *WaypointLink) Decode(r *binenc.Reader) error {
	var ok bool
	w.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Flags, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type Waypoint struct {
	ID    uint32
	Pos   types.Pointf
	Name  string
	Flags uint32
	Links []WaypointLink
}

func (w *Waypoint) EncodingSize() int {
	return 13 + len(w.Name) + 5 + 5*len(w.Links)
}

func (w *Waypoint) MarshalBinary() ([]byte, error) {
	data := make([]byte, 13+len(w.Name)+5, w.EncodingSize())
	i := 0
	binary.LittleEndian.PutUint32(data[i:], w.ID)
	i += 4
	binary.LittleEndian.PutUint32(data[i:], math.Float32bits(w.Pos.X))
	i += 4
	binary.LittleEndian.PutUint32(data[i:], math.Float32bits(w.Pos.Y))
	i += 4
	data[i] = byte(len(w.Name))
	i++
	i += copy(data[13:], w.Name)
	binary.LittleEndian.PutUint32(data[i:], w.Flags)
	i += 4
	data[i] = byte(len(w.Links))
	for _, l := range w.Links {
		b, err := l.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (w *Waypoint) Decode(r *binenc.Reader) error {
	var ok bool
	w.ID, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Pos, ok = r.ReadPointF32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Name, ok = r.ReadString8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Flags, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sz, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < int(sz); i++ {
		var l WaypointLink
		if err := l.Decode(r); err != nil {
			return err
		}
		w.Links = append(w.Links, l)
	}
	return nil
}

type Waypoints struct {
	Waypoints []Waypoint
}

func (*Waypoints) MapSection() string {
	return "WayPoints"
}

func (sect *Waypoints) MarshalBinary() ([]byte, error) {
	data := make([]byte, 6)
	binary.LittleEndian.PutUint16(data[0:], 4)
	binary.LittleEndian.PutUint32(data[2:], uint32(len(sect.Waypoints)))
	for _, w := range sect.Waypoints {
		b, err := w.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return data, nil
}

func (sect *Waypoints) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *Waypoints) Decode(r *binenc.Reader) error {
	*sect = Waypoints{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 4 {
		return fmt.Errorf("unsupported version of waypoints section: %d", vers)
	}
	n, ok := r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Waypoints = make([]Waypoint, 0, n)
	for i := 0; i < int(n); i++ {
		var w Waypoint
		if err := w.Decode(r); err != nil {
			return err
		}
		sect.Waypoints = append(sect.Waypoints, w)
	}
	return nil
}
