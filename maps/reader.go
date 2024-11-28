package maps

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/opennox/noxcrypt"

	"github.com/opennox/opennox-lib/binenc"
)

type Reader struct {
	cr      *crypt.Reader
	r       io.Reader
	err     error
	limited bool
	m       *Map
}

func NewReader(r io.Reader) (*Reader, error) {
	cr, err := crypt.NewReader(r, crypt.MapKey)
	if err != nil {
		return nil, err
	}
	rd := &Reader{cr: cr, r: cr}
	if err := rd.readHeader(); err != nil {
		return nil, err
	}
	return rd, nil
}

func (r *Reader) Header() Header {
	return r.m.Header()
}

func (r *Reader) Map() *Map {
	return r.m
}

func (r *Reader) Info() *Info {
	if r.m == nil {
		return nil
	}
	return &r.m.Info
}

func (r *Reader) error() error {
	return r.err
}

// readBytes reads bytes to completely fill the buffer.
// Read error must be checked with error method.
func (r *Reader) readBytes(p []byte) {
	if r.error() != nil {
		return
	}
	_, err := io.ReadFull(r.r, p)
	if err != nil {
		r.err = err
	}
}

// readAlignedBytes reads bytes to completely fill the buffer. Read operation will be aligned with the crypt block size.
// Read error must be checked with error method.
func (r *Reader) readAlignedBytes(p []byte) {
	if r.error() != nil {
		return
	}
	if r.limited {
		r.err = errors.New("trying to align a limited reader")
		return
	}
	n, err := r.cr.ReadAligned(p)
	if err != nil {
		r.err = err
		return
	} else if n != len(p) {
		r.err = io.ErrUnexpectedEOF
		return
	}
}

func (r *Reader) readU8() byte {
	var b [1]byte
	r.readBytes(b[:])
	return b[0]
}

func (r *Reader) readU16() uint16 {
	var b [2]byte
	r.readBytes(b[:])
	return binary.LittleEndian.Uint16(b[:])
}

func (r *Reader) readU32() uint32 {
	var b [4]byte
	r.readBytes(b[:])
	return binary.LittleEndian.Uint32(b[:])
}

func (r *Reader) readU64() uint64 {
	var b [8]byte
	r.readBytes(b[:])
	return binary.LittleEndian.Uint64(b[:])
}

func (r *Reader) readI32() int32 {
	return int32(r.readU32())
}

func (r *Reader) readAlignedU32() uint32 {
	var b [4]byte
	r.readAlignedBytes(b[:])
	return binary.LittleEndian.Uint32(b[:])
}

func (r *Reader) readStringFixed(max int) string {
	b := make([]byte, max)
	r.readBytes(b)
	if r.error() != nil {
		return ""
	}
	i := bytes.IndexByte(b, 0)
	if i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func (r *Reader) readString8() string {
	n := r.readU8()
	if r.error() != nil {
		return ""
	}
	return r.readStringFixed(int(n))
}

func (r *Reader) readAlignedString8() string {
	if r.error() != nil {
		return ""
	}
	if r.limited {
		r.err = errors.New("trying to align a limited reader")
		return ""
	}
	s := r.readString8()
	if r.error() != nil {
		return ""
	}
	if err := r.cr.Align(); err != nil {
		r.err = err
	}
	return s
}

func (r *Reader) readHeader() error {
	r.m = &Map{}
	r.m.magic = r.readU32()
	if err := r.error(); err != nil {
		return fmt.Errorf("cannot read magic: %w", err)
	}
	switch r.m.magic {
	case MagicOld:
		// nop
	case Magic:
		r.m.crc = r.readAlignedU32()
		if err := r.error(); err != nil {
			return fmt.Errorf("cannot read crc: %w", err)
		}
	default:
		return fmt.Errorf("unsupported magic: 0x%x", r.m.magic)
	}
	r.m.wallOffX = r.readU32()
	r.m.wallOffY = r.readU32()
	if err := r.error(); err != nil {
		return fmt.Errorf("cannot read wall offset: %w", err)
	}
	return nil
}

func (r *Reader) ReadInfo() (*Info, error) {
	sect, err := r.nextSection()
	if err != nil {
		return nil, err
	} else if sect != "MapInfo" {
		return nil, fmt.Errorf("unexpected section: %q", sect)
	}
	data, err := io.ReadAll(r.r)
	if err != nil {
		return nil, err
	}
	if err = r.m.Info.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return r.Info(), err
}

func (r *Reader) nextSection() (string, error) {
	if r.limited {
		if _, err := io.Copy(io.Discard, r.r); err != nil {
			return "", err
		}
		r.limited = false
		r.err = nil
		r.r = r.cr
	}
	sect := r.readAlignedString8()
	if err := r.error(); err != nil {
		return "", err
	}
	size := r.readU64()
	if err := r.error(); err != nil {
		return sect, err
	}
	r.r = io.LimitReader(r.cr, int64(size))
	r.limited = true
	return sect, nil
}

func (r *Reader) ReadSections() error {
	var buf bytes.Buffer
	for {
		sect, err := r.nextSection()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		buf.Reset()
		if _, err := io.Copy(&buf, r.r); err != nil {
			return err
		}
		rd := binenc.NewReader(buf.Bytes())
		switch sect {
		default:
			r.m.Unknown = append(r.m.Unknown, RawSection{
				Name: sect,
				Data: buf.Bytes(),
			})
			continue
		case "MapInfo":
			if err := r.m.Info.Decode(rd); err != nil {
				return err
			}
		case "MapIntro":
			r.m.Intro = new(MapIntro)
			if err := r.m.Intro.Decode(rd); err != nil {
				return err
			}
		case "AmbientData":
			r.m.Ambient = new(AmbientData)
			if err := r.m.Ambient.Decode(rd); err != nil {
				return err
			}
		case "WallMap":
			r.m.Walls = new(WallMap)
			if err := r.m.Walls.Decode(rd); err != nil {
				return err
			}
		case "FloorMap":
			r.m.Floor = new(FloorMap)
			if err := r.m.Floor.Decode(rd); err != nil {
				return err
			}
		case "ScriptObject":
			r.m.Script = new(Script)
			if err := r.m.Script.Decode(rd); err != nil {
				return err
			}
		case "ScriptData":
			r.m.ScriptData = new(ScriptData)
			if err := r.m.ScriptData.Decode(rd); err != nil {
				return err
			}
		case "SecretWalls":
			r.m.SecretWalls = new(SecretWalls)
			if err := r.m.SecretWalls.Decode(rd); err != nil {
				return err
			}
		case "WindowWalls":
			r.m.WindowWalls = new(WindowWalls)
			if err := r.m.WindowWalls.Decode(rd); err != nil {
				return err
			}
		case "DestructableWalls":
			r.m.DestructableWalls = new(DestructableWalls)
			if err := r.m.DestructableWalls.Decode(rd); err != nil {
				return err
			}
		case "WayPoints":
			r.m.Waypoints = new(Waypoints)
			if err := r.m.Waypoints.Decode(rd); err != nil {
				return err
			}
		}
		if n := rd.Remaining(); n > 0 {
			return fmt.Errorf("trailing %s data: [%d]", sect, n)
		}
	}
}

func (r *Reader) ReadSectionsRaw() ([]RawSection, error) {
	var out []RawSection
	for {
		sect, err := r.nextSection()
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return out, err
		}
		data, err := io.ReadAll(r.r)
		if err != nil {
			return out, err
		}
		out = append(out, RawSection{Name: sect, Data: data})
		if err != nil {
			return out, err
		}
	}
}
