package noxsave

import (
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/opennox/libs/binenc"
)

const (
	SectFileInfo    = SectionID(1)
	SectAttrib      = SectionID(2)
	SectStatus      = SectionID(3)
	SectInventory   = SectionID(4)
	SectSpellbook   = SectionID(5)
	SectEnchantment = SectionID(6)
	SectGUI         = SectionID(7)
	SectFieldGuide  = SectionID(8)
	SectJournal     = SectionID(9)
	SectGame        = SectionID(10)
	SectPAD         = SectionID(11)
	SectMusic       = SectionID(12)
)

//go:generate stringer -type=SectionID -trimprefix=Sect

type SectionID uint32

type Section interface {
	SectionID() SectionID
	binenc.Encoded
}

var registry = NewRegistry(nil)

func DefaultRegistry() *Registry {
	return registry
}

func Register(s Section) {
	registry.Register(s)
}

func NewRegistry(parent *Registry) *Registry {
	return &Registry{
		parent: parent,
		byID:   make(map[SectionID]reflect.Type),
	}
}

type Registry struct {
	parent *Registry
	byID   map[SectionID]reflect.Type
}

func (r *Registry) get(id SectionID) (reflect.Type, bool) {
	if t, ok := r.byID[id]; ok {
		return t, true
	}
	if r.parent == nil {
		return nil, false
	}
	return r.parent.get(id)
}

func (r *Registry) Get(id SectionID) Section {
	t, ok := r.get(id)
	if !ok {
		return &RawSection{ID: id}
	}
	return reflect.New(t).Interface().(Section)
}

func (r *Registry) Register(s Section) {
	id := s.SectionID()
	t := reflect.TypeOf(s).Elem()
	if _, ok := r.byID[id]; ok {
		panic("already registered")
	}
	r.byID[id] = t
}

type RawSection struct {
	ID   SectionID
	Data []byte
}

func (s *RawSection) SectionID() SectionID {
	return s.ID
}

func (s *RawSection) EncodeSize() int {
	return len(s.Data)
}

func (s *RawSection) Encode(data []byte) (int, error) {
	if len(data) < len(s.Data) {
		return 0, io.ErrShortBuffer
	}
	n := copy(data, s.Data)
	return n, nil
}

func (s *RawSection) Decode(data []byte) (int, error) {
	s.Data = slices.Clone(data)
	return len(data), nil
}

func (s *RawSection) DecodeWith(r *Registry) (Section, error) {
	if r == nil {
		r = DefaultRegistry()
	}
	sect := r.Get(s.ID)
	n, err := sect.Decode(s.Data)
	if err != nil {
		return s, fmt.Errorf("section %v: %w", s.ID, err)
	} else if n != len(s.Data) {
		return s, fmt.Errorf("section %v: partial read", s.ID)
	}
	return sect, nil
}
