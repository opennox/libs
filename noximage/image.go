package noximage

import (
	"image"
	"image/color"

	noxcolor "github.com/opennox/libs/color"
)

var _ image.Image = &Image16{}

func NewImage16(r image.Rectangle) *Image16 {
	return &Image16{
		Pix:    make([]uint16, r.Dx()*r.Dy()),
		Stride: r.Dx(),
		Rect:   r,
	}
}

func NewImage16WithData(data []uint16, sz image.Point) *Image16 {
	return &Image16{
		Pix:    data,
		Stride: sz.X,
		Rect:   image.Rect(0, 0, sz.X, sz.Y),
	}
}

type Image16 struct {
	Pix    []uint16
	Stride int
	Rect   image.Rectangle
}

func (p *Image16) Size() image.Point {
	if p == nil {
		return image.Point{}
	}
	return p.Rect.Size()
}

func (p *Image16) ColorModel() color.Model {
	return noxcolor.ModelRGBA5551
}

func (p *Image16) Bounds() image.Rectangle {
	return p.Rect
}

func (p *Image16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *Image16) Row(y int) []uint16 {
	y = y - p.Rect.Min.Y
	if y < 0 || y >= p.Stride {
		return nil
	}
	return p.Pix[y*p.Stride : (y+1)*p.Stride : (y+1)*p.Stride]
}

func (p *Image16) SubImage(r image.Rectangle) *Image16 {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &Image16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Image16{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

func (p *Image16) At(x, y int) color.Color {
	return p.NRGBAAt(x, y)
}

func (p *Image16) RGBAAt(x, y int) color.RGBA {
	return p.RGBA5551At(x, y).ColorRGBA()
}

func (p *Image16) NRGBAAt(x, y int) color.NRGBA {
	return p.RGBA5551At(x, y).ColorNRGBA()
}

func (p *Image16) RGBA5551At(x, y int) noxcolor.RGBA5551 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return noxcolor.RGBA5551(0)
	}
	i := p.PixOffset(x, y)
	return noxcolor.RGBA5551(p.Pix[i])
}

func (p *Image16) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.set(x, y, noxcolor.ToRGBA5551Color(c))
}

func (p *Image16) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.set(x, y, noxcolor.ToRGBA5551Color(c))
}

func (p *Image16) SetNRGBA(x, y int, c color.NRGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.set(x, y, noxcolor.RGBA5551Color(c.R, c.G, c.B, c.A))
}

func (p *Image16) SetRGBA5551(x, y int, c noxcolor.RGBA5551) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.set(x, y, c)
}

func (p *Image16) set(x, y int, c noxcolor.RGBA5551) {
	i := p.PixOffset(x, y)
	p.Pix[i] = uint16(c)
}
