package types

import (
	"fmt"
	"image"
	"math"
)

func Ptf(x, y float32) Pointf {
	return Pointf{X: x, Y: y}
}

func Point2f(p Point) Pointf {
	return Pointf{
		X: float32(p.X),
		Y: float32(p.Y),
	}
}

type Point = image.Point

type Pointf struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (p Pointf) Pos() Pointf {
	return p
}

func (p Pointf) Point() Point {
	return Point{
		X: int(p.X),
		Y: int(p.Y),
	}
}

func (p Pointf) Add(p2 Pointf) Pointf {
	p.X += p2.X
	p.Y += p2.Y
	return p
}

func (p Pointf) Sub(p2 Pointf) Pointf {
	p.X -= p2.X
	p.Y -= p2.Y
	return p
}

func (p Pointf) Mul(v float32) Pointf {
	p.X *= v
	p.Y *= v
	return p
}

func (p Pointf) Div(v float32) Pointf {
	p.X /= v
	p.Y /= v
	return p
}

func (p Pointf) Len() float64 {
	x, y := float64(p.X), float64(p.Y)
	return math.Hypot(x, y)
}

func (p Pointf) Normalize() Pointf {
	return p.Div(float32(p.Len()))
}

type Rectf struct {
	Min Pointf
	Max Pointf
}

func (r *Rectf) IsEmpty() bool {
	return r.Max.X <= r.Min.X || r.Max.Y <= r.Min.Y
}

func (r Rectf) Canon() Rectf {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func RectFromPointsf(p1, p2 Pointf) Rectf {
	return Rectf{Min: p1, Max: p2}.Canon()
}

type RGB struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
}

func (c RGB) HexString() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
