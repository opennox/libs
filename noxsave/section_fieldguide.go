package noxsave

import (
	"encoding/binary"
	"fmt"
	"io"
)

func init() {
	Register(&FieldGuides{})
}

type FieldGuides struct {
	Guides []string
}

func (*FieldGuides) SectionID() SectionID {
	return SectFieldGuide
}

func (s *FieldGuides) EncodeSize() int {
	n := 2 + 1
	if len(s.Guides) == 0 {
		return n
	}
	n++
	for _, g := range s.Guides {
		n += 1 + len(g)
	}
	return n
}

func (s *FieldGuides) Encode(data []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FieldGuides) Decode(data []byte) (int, error) {
	left := data
	if len(left) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	vers := binary.LittleEndian.Uint16(left[:2])
	left = left[2:]
	if vers > 1 {
		return 0, fmt.Errorf("unsupported section version: %d", vers)
	}
	*s = FieldGuides{}
	if len(left) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	if left[0] == 0 {
		return len(data) - len(left), nil
	}
	left = left[1:]

	if len(left) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	cnt := left[0]
	left = left[1:]

	for range cnt {
		if len(left) < 1 {
			return 0, io.ErrUnexpectedEOF
		}
		ssz := left[0]
		left = left[1:]
		if len(left) < int(ssz) {
			return 0, io.ErrUnexpectedEOF
		}
		name := string(left[:ssz])
		left = left[ssz:]
		s.Guides = append(s.Guides, name)
	}
	return len(data) - len(left), nil
}
