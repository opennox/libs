package maps

import (
	"encoding"
	"errors"
	"fmt"
	"image"
	"math"
	"reflect"

	"golang.org/x/exp/slices"

	"github.com/opennox/libs/binenc"
)

const (
	Magic    = 0xFADEFACE
	MagicOld = 0xFADEBEEF
)

var mapSections = make(map[string]reflect.Type)

func RegisterSection(sect Section) {
	name := sect.MapSection()
	if _, ok := mapSections[name]; ok {
		panic("already registered: " + name)
	}
	mapSections[name] = reflect.TypeOf(sect).Elem()
}

type RawSection struct {
	Name string
	Data []byte
}

func (sect RawSection) Supported() bool {
	_, ok := mapSections[sect.Name]
	return ok
}

func (sect RawSection) Decode() (Section, error) {
	typ, ok := mapSections[sect.Name]
	if !ok {
		return nil, errors.New("unsupported section: " + sect.Name)
	}
	s := reflect.New(typ).Interface().(Section)
	rd := binenc.NewReader(sect.Data)
	err := s.Decode(rd)
	if err != nil {
		return s, err
	}
	if n := rd.Remaining(); n > 0 {
		return s, fmt.Errorf("trailing %s data: [%d]", sect.Name, n)
	}
	return s, nil
}

type Section interface {
	MapSection() string
	Decode(r *binenc.Reader) error
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

var sectionOrder = []string{
	"MapInfo",
	"WallMap",
	"FloorMap",
	"SecretWalls",
	"DestructableWalls",
	"WayPoints",
	"DebugData",
	"WindowWalls",
	"GroupData",
	"ScriptObject",
	"AmbientData",
	"Polygons",
	"MapIntro",
	"ScriptData",
	"ObjectTOC",
	"ObjectData",
}

// SectionOrder returns an order in which the section should be written.
// It returns math.MaxInt for unknown sections, so they sort last.
func SectionOrder(name string) int {
	i := slices.Index(sectionOrder, name)
	if i < 0 {
		return math.MaxInt
	}
	return i
}

// SortSections sorts a slice of sections according to SectionOrder.
func SortSections(arr []Section) {
	slices.SortFunc(arr, func(a, b Section) int {
		i, j := SectionOrder(a.MapSection()), SectionOrder(b.MapSection())
		if i < j {
			return -1
		} else if i > j {
			return +1
		}
		return 0
	})
}

// SortRawSections sorts a slice of sections according to SectionOrder.
func SortRawSections(arr []RawSection) {
	slices.SortFunc(arr, func(a, b RawSection) int {
		i, j := SectionOrder(a.Name), SectionOrder(b.Name)
		if i < j {
			return -1
		} else if i > j {
			return +1
		}
		return 0
	})
}

type Info struct {
	Filename string `json:"name"`
	Size     int    `json:"size"`
	MapInfo
}

type Map struct {
	Info

	magic    uint32
	crc      uint32
	wallOffX uint32
	wallOffY uint32

	Intro             *MapIntro
	Ambient           *AmbientData
	Walls             *WallMap
	Floor             *FloorMap
	Script            *Script
	ScriptData        *ScriptData
	SecretWalls       *SecretWalls
	WindowWalls       *WindowWalls
	DestructableWalls *DestructableWalls
	Waypoints         *Waypoints
	Unknown           []RawSection
}

func (m *Map) Header() Header {
	return Header{
		Magic: m.magic,
		Offs: image.Point{
			X: int(m.wallOffX),
			Y: int(m.wallOffY),
		},
	}
}

func (m *Map) CRC() uint32 {
	return m.crc
}

// GridBoundingBox returns a bounding box for all walls and tiles on the map.
// Returned rectangle uses grid coordinates, not pixel coordinates.
func (m *Map) GridBoundingBox() image.Rectangle {
	var r image.Rectangle
	r.Min.X = -1
	r.Min.Y = -1
	for _, w := range m.Walls.Walls {
		p := w.Pos
		if r.Min.X == -1 || r.Min.X > int(p.X) {
			r.Min.X = int(p.X)
		}
		if r.Min.Y == -1 || r.Min.Y > int(p.Y) {
			r.Min.Y = int(p.Y)
		}
		if r.Max.X < int(p.X) {
			r.Max.X = int(p.X)
		}
		if r.Max.Y < int(p.Y) {
			r.Max.Y = int(p.Y)
		}
	}
	// TODO: tiles
	return r
}
