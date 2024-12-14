package sdl

import (
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/client/seat/opengl"
	"github.com/opennox/libs/env"
	noxlog "github.com/opennox/libs/log"
)

var (
	debugGpad = os.Getenv("NOX_DEBUG_GPAD") == "true"
)

var _ seat.Seat = &Window{}

// New creates a new SDL window which implements a Seat interface.
func New(log *slog.Logger, title string, sz image.Point) (*Window, error) {
	log = noxlog.WithSystem(log, "sdl")
	// That hint won't work if it is called after sdl.Init
	sdl.SetHint(sdl.HINT_WINDOWS_DPI_AWARENESS, "permonitorv2")
	// TODO: if we ever decide to use multiple windows, this will need to be moved elsewhere; same for sdl.Quit
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER); err != nil {
		return nil, fmt.Errorf("SDL Initialization failed: %w", err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_CORE, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL core: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL forward: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL attribute: %w", err)
	}

	win, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(sz.X), int32(sz.Y),
		sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("SDL Window creation failed: %w", err)
	}
	h := &Window{
		log: log,
		win: win, prevSz: sz,
	}
	h.SetScreenMode(seat.Windowed)
	if err := h.initGL(); err != nil {
		_ = win.Destroy()
		sdl.Quit()
		return nil, err
	}
	return h, nil
}

type Window struct {
	log      *slog.Logger
	win      *sdl.Window
	gl       opengl.Window
	prevPos  image.Point
	prevSz   image.Point
	textInp  bool
	mode     seat.ScreenMode
	rel      bool
	mpos     image.Point
	onResize []func(sz image.Point)
	onInput  []func(ev seat.InputEvent)
}

func (win *Window) Close() error {
	if win.win == nil {
		return nil
	}
	win.gl.Close()
	err := win.win.Destroy()
	win.win = nil
	win.onResize = nil
	win.onInput = nil
	sdl.Quit()
	return err
}

func (win *Window) NewSurface(sz image.Point, filter bool) seat.Surface {
	return win.gl.NewSurface(sz, filter)
}

func (win *Window) Clear() {
	win.gl.Clear()
}

func (win *Window) ScreenSize() image.Point {
	w, h := win.win.GLGetDrawableSize()
	return image.Point{
		X: int(w), Y: int(h),
	}
}

func (win *Window) screenPos() image.Point {
	x, y := win.win.GetPosition()
	return image.Point{
		X: int(x), Y: int(y),
	}
}

func (win *Window) displayRect() sdl.Rect {
	disp, err := win.win.GetDisplayIndex()
	if err != nil {
		win.log.Warn("can't get display index", "err", err)
		return sdl.Rect{}
	}
	rect, err := sdl.GetDisplayBounds(disp)
	if err != nil {
		win.log.Warn("can't get display bounds", "err", err)
		return sdl.Rect{}
	}
	return rect
}

func (win *Window) ScreenMaxSize() image.Point {
	rect := win.displayRect()
	return image.Point{
		X: int(rect.W), Y: int(rect.H),
	}
}

func (win *Window) setSize(sz image.Point) {
	win.log.Info("window", "size", sz)
	win.win.SetSize(int32(sz.X), int32(sz.Y))
}

func (win *Window) center() {
	win.win.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
}

func (win *Window) ResizeScreen(sz image.Point) {
	if win.mode != seat.Windowed {
		return
	}
	win.setSize(sz)
	win.prevSz = sz
}

func (win *Window) setRelative(rel bool) {
	if win.rel == rel {
		return
	}
	win.rel = rel
	win.win.SetGrab(rel)
	sdl.SetRelativeMouseMode(rel)
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
	switch mode {
	case seat.Windowed:
		win.win.SetFullscreen(0)
		win.win.SetResizable(true)
		win.win.SetBordered(true)
		win.setSize(win.prevSz)
		if win.prevPos != (image.Point{}) {
			win.win.SetPosition(int32(win.prevPos.X), int32(win.prevPos.Y))
		} else {
			win.center()
		}
		if env.IsDevMode() || env.IsE2E() {
			sdl.ShowCursor(sdl.ENABLE)
		} else {
			sdl.ShowCursor(sdl.DISABLE)
		}
		win.setRelative(false)
	case seat.Fullscreen:
		win.win.SetResizable(false)
		win.win.SetBordered(false)
		win.setSize(win.ScreenMaxSize())
		win.win.SetFullscreen(uint32(sdl.WINDOW_FULLSCREEN_DESKTOP))
		sdl.ShowCursor(sdl.DISABLE)
		win.setRelative(true)
	case seat.Borderless:
		win.win.SetFullscreen(0)
		win.win.SetResizable(false)
		win.win.SetBordered(true)
		win.setSize(win.ScreenMaxSize())
		win.center()
		sdl.ShowCursor(sdl.DISABLE)
		win.setRelative(false)
	}
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
	gtx, err := win.win.GLCreateContext()
	if err != nil {
		return fmt.Errorf("OpenGL creation failed: %w", err)
	}
	err = win.win.GLMakeCurrent(gtx)
	if err != nil {
		return fmt.Errorf("OpenGL bind failed: %w", err)
	}
	sdl.GLSetSwapInterval(0)
	if err := win.gl.Init(win.log); err != nil {
		return err
	}
	return nil
}

func (win *Window) Present() {
	win.win.GLSwap()
}
