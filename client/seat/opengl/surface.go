package opengl

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/noximage"
)

func (win *Window) NewSurface(sz image.Point, filter bool) seat.Surface {
	s := &Surface{win: win, sz: sz}
	gl.GenTextures(1, &s.tex)
	gl.BindTexture(gl.TEXTURE_2D, s.tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	if filter {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(s.sz.X), int32(s.sz.Y), 0, gl.BGRA, gl.UNSIGNED_SHORT_1_5_5_5_REV, nil)
	return s
}

type Surface struct {
	win *Window
	sz  image.Point
	tex uint32
}

func (s *Surface) rect() image.Point {
	return s.sz
}

func (s *Surface) Update(img *noximage.Image16) {
	if s.sz != img.Size() {
		panic("invalid image size")
	}
	gl.BindTexture(gl.TEXTURE_2D, s.tex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(s.sz.X), int32(s.sz.Y), 0, gl.BGRA, gl.UNSIGNED_SHORT_1_5_5_5_REV, gl.Ptr(img.Pix))
}

func (s *Surface) Size() image.Point {
	return s.sz
}

func (s *Surface) Draw(vp image.Rectangle) {
	gl.Viewport(int32(vp.Min.X), int32(vp.Min.Y), int32(vp.Dx()), int32(vp.Dy()))
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.tex)
	//gl.BindVertexArray(s.win.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.win.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.win.ebo)
	gl.UseProgram(s.win.prog)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (s *Surface) Destroy() {
	gl.DeleteTextures(1, &s.tex)
	s.tex = 0
	s.win = nil
}
