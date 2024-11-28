package maps

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/opennox/libs/binenc"
)

func init() {
	RegisterSection(&FloorMap{})
}

type GridData struct {
	Prefix uint16
	Var1   uint32
	Var2   uint32
	Var3   uint32
	Var4   uint32
}

func (g *GridData) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2+4*4)
	binary.LittleEndian.PutUint16(data[0:], g.Prefix)
	binary.LittleEndian.PutUint32(data[2:], g.Var1)
	binary.LittleEndian.PutUint32(data[6:], g.Var2)
	binary.LittleEndian.PutUint32(data[10:], g.Var3)
	binary.LittleEndian.PutUint32(data[14:], g.Var4)
	return data, nil
}

func (g *GridData) Decode(r *binenc.Reader) error {
	var ok bool
	g.Prefix, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	g.Var1, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	g.Var2, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	g.Var3, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	g.Var4, ok = r.ReadU32()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type Edge struct {
	Image   byte
	Variant uint16
	Edge    byte
	Dir     byte
}

func (e *Edge) MarshalBinary() ([]byte, error) {
	data := make([]byte, 1+2+2*1)
	data[0] = e.Image
	binary.LittleEndian.PutUint16(data[1:], e.Variant)
	data[3] = e.Edge
	data[4] = e.Dir
	return data, nil
}

func (e *Edge) Decode(r *binenc.Reader) error {
	var ok bool
	e.Image, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	e.Variant, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	e.Edge, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	e.Dir, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type Tile struct {
	Image   byte
	Variant uint16
	Field4  uint16
	Edges   []Edge
}

func (t *Tile) IsZero() bool {
	return t.Image == 0 && t.Variant == 0 && t.Field4 == 0 && len(t.Edges) == 0
}

func (t *Tile) size() int {
	return 6 + 5*len(t.Edges)
}

func (t *Tile) MarshalBinary() ([]byte, error) {
	data := make([]byte, t.size())
	data[0] = t.Image
	binary.LittleEndian.PutUint16(data[1:], t.Variant)
	binary.LittleEndian.PutUint16(data[3:], t.Field4)
	data[5] = byte(len(t.Edges))
	cur := data[6:]
	for _, e := range t.Edges {
		edata, err := e.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(cur, edata)
		cur = cur[5:]
	}
	return data, nil
}

func (t *Tile) Decode(r *binenc.Reader) error {
	var ok bool
	t.Image, ok = r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	t.Variant, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	// TODO: check in the engine what this is
	t.Field4, ok = r.ReadU16()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	n, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	t.Edges = make([]Edge, 0, n)
	for i := 0; i < int(n); i++ {
		var e Edge
		if err := e.Decode(r); err != nil {
			return err
		}
		t.Edges = append(t.Edges, e)
	}
	return nil
}

type FloorPos struct {
	X, Y uint16
}

type TilePair struct {
	Pos    FloorPos
	F1, F2 byte
	L, R   *Tile
}

func (p *TilePair) IsZero() bool {
	return p.Pos == (FloorPos{}) && p.F1 == 0 && p.F2 == 0 && p.L == nil && p.R == nil
}

func (p *TilePair) HasLeft() bool {
	return p.L != nil
}

func (p *TilePair) HasRight() bool {
	return p.R != nil
}

func (p *TilePair) LeftPos() FloorPos {
	return FloorPos{X: 2 * p.Pos.X, Y: 2 * p.Pos.Y}
}

func (p *TilePair) RightPos() FloorPos {
	return FloorPos{X: 2*p.Pos.X + 1, Y: 2*p.Pos.Y - 1}
}

func (p *TilePair) size() int {
	sz := 2
	if p.HasLeft() {
		sz += p.L.size()
	}
	if p.HasRight() {
		sz += p.R.size()
	}
	return sz
}

func (p *TilePair) MarshalBinary() ([]byte, error) {
	if p.IsZero() {
		return []byte{0xff, 0xff}, nil
	}
	data := make([]byte, p.size())
	data[0] = byte(p.Pos.X) | p.F1
	data[1] = byte(p.Pos.Y) | p.F2
	cur := data[2:]
	if p.HasRight() {
		data[0] |= 0x80
		tdata, err := p.R.MarshalBinary()
		if err != nil {
			return nil, err
		}
		n := copy(cur, tdata)
		cur = cur[n:]
	}
	if p.HasLeft() {
		data[1] |= 0x80
		tdata, err := p.L.MarshalBinary()
		if err != nil {
			return nil, err
		}
		n := copy(cur, tdata)
		cur = cur[n:]
	}
	return data, nil
}

func (p *TilePair) Decode(r *binenc.Reader) error {
	x, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	y, ok := r.ReadU8()
	if !ok {
		return io.ErrUnexpectedEOF
	}
	*p = TilePair{}
	if x == 0xff && y == 0xff {
		return nil
	}
	hasR := x&0x80 != 0
	hasL := y&0x80 != 0
	p.F1 = x & 0x3
	p.F2 = y & 0x3
	p.Pos = FloorPos{X: uint16(x & 0x7c), Y: uint16(y & 0x7c)}
	if hasR {
		p.R = new(Tile)
		if err := p.R.Decode(r); err != nil {
			return err
		}
	}
	if hasL {
		p.L = new(Tile)
		if err := p.L.Decode(r); err != nil {
			return err
		}
	}
	return nil
}

type FloorMap struct {
	Grid  GridData
	Tiles []TilePair
}

func (*FloorMap) MapSection() string {
	return "FloorMap"
}

func (sect *FloorMap) MarshalBinary() ([]byte, error) {
	data, err := sect.Grid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	for _, p := range sect.Tiles {
		tdata, err := p.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, tdata...)
	}
	var end TilePair
	tdata, err := end.MarshalBinary()
	if err != nil {
		return nil, err
	}
	data = append(data, tdata...)
	return data, nil
}

func (sect *FloorMap) UnmarshalBinary(data []byte) error {
	return sect.Decode(binenc.NewReader(data))
}

func (sect *FloorMap) Decode(r *binenc.Reader) error {
	*sect = FloorMap{}
	if err := sect.Grid.Decode(r); err != nil {
		return err
	} else if sect.Grid.Prefix <= 3 {
		return fmt.Errorf("unsupported floor map: 0x%x", sect.Grid.Prefix)
	}
	for {
		var p TilePair
		if err := p.Decode(r); err != nil {
			return err
		}
		if p.IsZero() {
			if r.Remaining() > 0 {
				return fmt.Errorf("trailing floor data: [%d]", r.Remaining())
			}
			return nil
		}
		sect.Tiles = append(sect.Tiles, p)
	}
}
