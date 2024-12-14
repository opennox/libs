package glfw

import (
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/client/seat/opengl"
	"github.com/opennox/libs/env"
	noxlog "github.com/opennox/libs/log"
)

var (
	debug     = os.Getenv("GLFW_DEBUG") == "true"
	debugGpad = os.Getenv("NOX_DEBUG_GPAD") == "true"
)

var _ seat.Seat = &Window{}

// New creates a new SDL window which implements a Seat interface.
func New(log *slog.Logger, title string, sz image.Point) (*Window, error) {
	log = noxlog.WithSystem(log, "glfw")
	// TODO: if we ever decide to use multiple windows, this will need to be moved elsewhere; same for glfw.Terminate
	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("cannot init GLFW: %w", err)
	}
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	if debug {
		glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	}
	glfw.WindowHint(glfw.RedBits, 5)
	glfw.WindowHint(glfw.GreenBits, 5)
	glfw.WindowHint(glfw.BlueBits, 5)
	glfw.WindowHint(glfw.AlphaBits, 1)
	win, err := glfw.CreateWindow(sz.X, sz.Y, title, nil, nil)
	if err != nil {
		panic(err)
	}
	h := &Window{
		log: log,
		win: win, prevSz: sz,
	}
	win.SetMouseButtonCallback(h.processMouseButtonEvent)
	win.SetCursorPosCallback(h.processMotionEvent)
	win.SetScrollCallback(h.processWheelEvent)
	win.SetKeyCallback(h.processKeyboardEvent)
	win.SetCloseCallback(h.processQuitEvent)
	win.SetFocusCallback(h.processFocusEvent)
	win.SetFramebufferSizeCallback(func(_ *glfw.Window, width int, height int) {
		sz := image.Pt(width, height)
		log.Info("framebuffer", "size", sz)
		for _, fnc := range h.onResize {
			fnc(sz)
		}
	})
	win.SetRefreshCallback(func(w *glfw.Window) {
		log.Debug("refresh")
	})
	//win.SetCharCallback(func(_ *glfw.Window, c rune) {
	//
	//})
	h.SetScreenMode(seat.Windowed)
	if err := h.initGL(); err != nil {
		win.Destroy()
		glfw.Terminate()
		return nil, err
	}
	return h, nil
}

type Window struct {
	log      *slog.Logger
	win      *glfw.Window
	gl       opengl.Window
	prevPos  image.Point
	prevSz   image.Point
	textInp  bool
	mode     seat.ScreenMode
	rel      bool
	mpos     image.Point
	onResize []func(sz image.Point)
	onInput  []func(ev seat.InputEvent)
	input    inputData
}

func (win *Window) Close() error {
	if win.win == nil {
		return nil
	}
	win.gl.Close()
	win.win.Destroy()
	win.win = nil
	win.onResize = nil
	win.onInput = nil
	glfw.Terminate()
	return nil
}

func (win *Window) NewSurface(sz image.Point, filter bool) seat.Surface {
	win.win.MakeContextCurrent()
	return win.gl.NewSurface(sz, filter)
}

func (win *Window) Clear() {
	win.gl.Clear()
}

func (win *Window) ScreenSize() image.Point {
	w, h := win.win.GetFramebufferSize()
	return image.Pt(w, h)
}

func (win *Window) screenPos() image.Point {
	x, y := win.win.GetPos()
	return image.Pt(x, y)
}

func (win *Window) monitor() *glfw.Monitor {
	mon := win.win.GetMonitor()
	if mon == nil {
		mon = glfw.GetPrimaryMonitor()
	}
	return mon
}

func (win *Window) ScreenMaxSize() image.Point {
	mon := win.monitor()
	m := mon.GetVideoMode()
	return image.Pt(m.Width, m.Height)
}

func (win *Window) ResizeScreen(sz image.Point) {
	if win.mode != seat.Windowed {
		return
	}
	win.log.Info("window resize", "size", sz)
	win.win.SetSize(sz.X, sz.Y)
	win.prevSz = sz
}

func (win *Window) setRelative(rel bool) {
	if win.rel == rel {
		return
	}
	win.rel = rel
	if rel {
		win.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else if env.IsDevMode() || env.IsE2E() {
		win.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		win.win.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

func (win *Window) SetScreenMode(mode seat.ScreenMode) {
	if win.mode == mode {
		return
	}
	if win.mode == seat.Windowed {
		// preserve size and pos, so we can restore them later
		win.prevSz = win.ScreenSize()
		win.prevPos = win.screenPos()
	}
	var (
		monitor *glfw.Monitor
		pos     image.Point
		sz      image.Point
		refresh = 60
		rel     bool
	)
	mon := win.monitor()
	switch mode {
	case seat.Windowed:
		monitor = nil // windowed
		pos = win.prevPos
		sz = win.prevSz
		if pos == (image.Point{}) {
			md := mon.GetVideoMode()
			pos = image.Pt((md.Width-sz.X)/2, (md.Height-sz.Y)/2)
		}
		rel = false
	case seat.Fullscreen:
		monitor = mon // fullscreen
		m := mon.GetVideoMode()
		pos = image.Point{}
		sz = image.Pt(m.Width, m.Height)
		rel = true
	case seat.Borderless:
		monitor = mon // fullscreen
		m := mon.GetVideoMode()
		pos = image.Point{}
		sz = image.Pt(m.Width, m.Height)
		refresh = m.RefreshRate
		rel = false
	}
	name := ""
	if monitor != nil {
		name = monitor.GetName()
	}
	win.log.Info("set window", "name", name, "size", sz, "pos", pos, "rate", refresh)
	win.win.SetMonitor(monitor, pos.X, pos.Y, sz.X, sz.Y, refresh)
	win.setRelative(rel)
	win.mode = mode
}

// SetGamma sets screen gamma parameter.
func (win *Window) SetGamma(v float32) {
	win.gl.SetGamma(v)
}

func (win *Window) OnScreenResize(fnc func(sz image.Point)) {
	win.onResize = append(win.onResize, fnc)
}

func (win *Window) initGL() error {
	win.win.MakeContextCurrent()
	glfw.SwapInterval(0)
	if err := win.gl.Init(win.log); err != nil {
		return err
	}
	return nil
}

func (win *Window) Present() {
	win.win.SwapBuffers()
}
