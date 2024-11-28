package maps

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/opennox/opennox-lib/binenc"
)

func init() {
	RegisterSection(&SecretWalls{})
}

type SecretWall struct {
	Pos       image.Point
	OpenWait  uint32
	Flags     byte
	State     byte
	OpenDelay byte
	LastOpen  uint32
	R2        uint32
}

func (w *SecretWall) MarshalBinary() ([]byte, error) {
	data := make([]byte, 23)
	binary.LittleEndian.PutUint32(data[0:], uint32(int32(w.Pos.X)))
	binary.LittleEndian.PutUint32(data[4:], uint32(int32(w.Pos.Y)))
	binary.LittleEndian.PutUint32(data[8:], w.OpenWait)
	data[12] = w.Flags
	data[13] = w.State
	data[14] = w.OpenDelay
	binary.LittleEndian.PutUint32(data[15:], w.LastOpen)
	binary.LittleEndian.PutUint32(data[19:], w.R2)
	return data, nil
}

func (w *SecretWall) Decode(r *binenc.Reader) error {
	var ok bool
	w.Pos, ok = r.ReadPointI32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.OpenWait, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.Flags, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.State, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.OpenDelay, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.LastOpen, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	w.R2, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type SecretWalls struct {
	Walls []SecretWall
}

func (*SecretWalls) MapSection() string {
	return "SecretWalls"
}

func (sect *SecretWalls) MarshalBinary() ([]byte, error) {
	data := make([]byte, 4, 4+23*len(sect.Walls))
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

func (sect *SecretWalls) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *SecretWalls) Decode(r *binenc.Reader) error {
	*sect = SecretWalls{}
	vers, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	} else if vers != 2 {
		return fmt.Errorf("unsupported version of secret walls section: %d", vers)
	}
	n, ok := r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	sect.Walls = make([]SecretWall, 0, n)
	for i := 0; i < int(n); i++ {
		var w SecretWall
		if err := w.Decode(r); err != nil {
			return err
		}
		sect.Walls = append(sect.Walls, w)
	}
	return nil
}
