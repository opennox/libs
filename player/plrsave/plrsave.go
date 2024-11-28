package plrsave

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"

	crypt "github.com/opennox/noxcrypt"
)

// Section is a common interface for player save file sections.
type Section interface {
	// ID returns a unique section ID used in binary encoding.
	ID() uint32
	// SectName returns a human-friendly section name.
	SectName() string
	encoding.BinaryUnmarshaler
}

var sections = make(map[uint32]reflect.Type)

// RegisterSection registers a section for automated decoding.
func RegisterSection(s Section) {
	id := s.ID()
	if _, ok := sections[id]; ok {
		panic("already registered")
	}
	sections[id] = reflect.TypeOf(s).Elem()
}

func sectionByID(id uint32) Section {
	rt, ok := sections[id]
	if !ok {
		return &UnknownSect{Ind: id}
	}
	return reflect.New(rt).Interface().(Section)
}

func readSectionSize(cr *crypt.Reader) (int, error) {
	size, err := cr.ReadI32()
	if err == io.EOF {
		return 0, io.ErrUnexpectedEOF
	} else if err != nil {
		return 0, err
	}
	if err = cr.Align(); err != nil {
		return 0, err
	}
	if size < 0 {
		return 0, errors.New("invalid section size")
	}
	return int(size), nil
}

// UnknownSect represents a section that is not yet supported.
type UnknownSect struct {
	Ind  uint32
	Data []byte
}

func (s *UnknownSect) ID() uint32 {
	return s.Ind
}

func (s *UnknownSect) SectName() string {
	return fmt.Sprintf("Unk%d", s.Ind)
}

func (s *UnknownSect) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}

// ReadAll reads all player save file sections.
func ReadAll(r io.Reader) ([]Section, error) {
	sr, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	var out []Section
	for {
		sect, err := sr.ReadSection()
		if err == io.EOF {
			break
		}
		if sect != nil {
			out = append(out, sect)
		}
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

// ReadInfo reads save file info.
func ReadInfo(r io.Reader) (*FileInfo, error) {
	sr, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	for {
		id, data, err := sr.ReadSectionRaw()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		} else if err != nil {
			return nil, err
		}
		if id == FileInfoID {
			var info FileInfo
			err := info.UnmarshalBinary(data)
			return &info, err
		}
	}
}

// NewReader creates a new player save file reader.
func NewReader(r io.Reader) (*Reader, error) {
	cr, err := crypt.NewReader(r, crypt.SaveKey)
	if err != nil {
		return nil, err
	}
	return &Reader{cr: cr}, nil
}

// Reader is a player save file reader.
type Reader struct {
	cr *crypt.Reader
}

// ReadSectionRaw reads next save file section. It returns its ID and raw data.
func (r *Reader) ReadSectionRaw() (uint32, []byte, error) {
	id, err := r.cr.ReadU32()
	if err != nil {
		return 0, nil, err
	}
	err = r.cr.Align()
	if err != nil {
		return id, nil, err
	}
	size, err := readSectionSize(r.cr)
	if err != nil {
		return id, nil, err
	}
	data := make([]byte, size)
	_, err = io.ReadFull(r.cr, data)
	return id, data, err
}

// ReadSection reads next save file section.
func (r *Reader) ReadSection() (Section, error) {
	id, data, err := r.ReadSectionRaw()
	if err != nil {
		return nil, err
	}
	sect := sectionByID(id)
	err = sect.UnmarshalBinary(data)
	return sect, err
}
