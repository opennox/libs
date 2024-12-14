package opengl

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

//go:embed shaders/screen.vert
var shaderVert string

//go:embed shaders/screen.frag
var shaderFrag string

type Window struct {
	log       *slog.Logger
	vao       uint32
	vbo       uint32
	ebo       uint32
	prog      uint32
	frag      uint32
	vert      uint32
	gammaAttr int32
}

// SetGamma sets screen gamma parameter.
func (win *Window) SetGamma(v float32) {
	gl.UseProgram(win.prog)
	gl.Uniform1f(win.gammaAttr, v)
}

// Init the OpenGL. Window context must be set as current by the caller.
func (win *Window) Init(log *slog.Logger) error {
	win.log = log
	if err := gl.Init(); err != nil {
		return fmt.Errorf("OpenGL init failed: %w", err)
	}
	win.log.Info("OpenGL", "vers", gl.GoStr(gl.GetString(gl.VERSION)))

	gl.GenVertexArrays(1, &win.vao)
	gl.BindVertexArray(win.vao)

	gl.GenBuffers(1, &win.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, win.vbo)
	quad := []float32{
		// pos, tex
		-1, -1, 0, 1,
		-1, +1, 0, 0,
		+1, +1, 1, 0,
		+1, -1, 1, 1,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(quad)*4, gl.Ptr(quad), gl.STATIC_DRAW)

	gl.GenBuffers(1, &win.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, win.ebo)
	elems := []uint32{
		0, 1, 2,
		2, 3, 0,
	}
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elems)*4, gl.Ptr(elems), gl.STATIC_DRAW)

	if err := win.initProgram(); err != nil {
		return err
	}
	return nil
}

func (win *Window) Close() {
	gl.DeleteProgram(win.prog)
	gl.DeleteShader(win.vert)
	gl.DeleteShader(win.frag)
	gl.DeleteBuffers(1, &win.vbo)
	gl.DeleteBuffers(1, &win.ebo)
	gl.DeleteVertexArrays(1, &win.vao)
}

func (win *Window) Clear() {
	gl.Disable(gl.DEPTH_TEST)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (win *Window) compileShader(typ uint32, src string) (uint32, error) {
	if !strings.HasSuffix(src, "\x00") {
		src += "\x00"
	}
	s := gl.CreateShader(typ)
	cstr, free := gl.Strs(src)
	gl.ShaderSource(s, 1, cstr, nil)
	free()
	gl.CompileShader(s)
	var st int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &st)
	if st == gl.FALSE {
		var sz int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &sz)
		text := make([]byte, sz+1)
		gl.GetShaderInfoLog(s, sz, nil, &text[0])
		return 0, errors.New(string(text))
	}
	return s, nil
}

func (win *Window) initProgram() error {
	vert, err := win.compileShader(gl.VERTEX_SHADER, shaderVert)
	if err != nil {
		return fmt.Errorf("cannot compile vertex shader: %w", err)
	}
	frag, err := win.compileShader(gl.FRAGMENT_SHADER, shaderFrag)
	if err != nil {
		return fmt.Errorf("cannot compile vertex shader: %w", err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.BindFragDataLocation(prog, 0, gl.Str("color\x00"))
	gl.LinkProgram(prog)

	var st int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &st)
	if st == gl.FALSE {
		var sz int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &sz)
		text := make([]byte, sz+1)
		gl.GetProgramInfoLog(prog, sz, nil, &text[0])
		return errors.New(string(text))
	}
	win.prog = prog

	gl.UseProgram(prog)
	win.gammaAttr = gl.GetUniformLocation(prog, gl.Str("gamma\x00"))
	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("tex\x00")), 0)
	gl.Uniform1f(win.gammaAttr, 1.0)

	posAttr := gl.GetAttribLocation(prog, gl.Str("position\x00"))
	gl.EnableVertexAttribArray(uint32(posAttr))
	gl.VertexAttribPointer(uint32(posAttr), 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0*4))

	texAttr := gl.GetAttribLocation(prog, gl.Str("texcoord\x00"))
	gl.EnableVertexAttribArray(uint32(texAttr))
	gl.VertexAttribPointer(uint32(texAttr), 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	return nil
}
