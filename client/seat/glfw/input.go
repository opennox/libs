package glfw

import (
	"image"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/types"
)

type inputData struct {
	prevX, prevY float64
}

func (win *Window) ReplaceInputs(cfg seat.InputConfig) seat.InputConfig {
	oldCfg := win.onInput
	win.onInput = cfg
	return oldCfg
}

func (win *Window) OnInput(fnc func(ev seat.InputEvent)) {
	win.onInput = append(win.onInput, fnc)
}

func (win *Window) SetTextInput(enable bool) {
	if win.textInp == enable {
		return
	}
	win.textInp = enable
	//if enable {
	//	win.win.
	//	sdl.StartTextInput()
	//} else {
	//	sdl.StopTextInput()
	//}
}

func (win *Window) InputTick() {
	glfw.PollEvents()
}

func (win *Window) inputEvent(ev seat.InputEvent) {
	for _, fnc := range win.onInput {
		fnc(ev)
	}
}

func (win *Window) processQuitEvent(_ *glfw.Window) {
	win.inputEvent(seat.WindowClosed)
}

func (win *Window) processFocusEvent(_ *glfw.Window, focused bool) {
	if focused {
		win.inputEvent(seat.WindowFocused)
	} else {
		win.inputEvent(seat.WindowUnfocused)
	}
}

//
//func (win *Window) processTextEditingEvent(ev *sdl.TextEditingEvent) {
//	win.inputEvent(&seat.TextEditEvent{
//		Text: ev.GetText(),
//	})
//}
//
//func (win *Window) processTextInputEvent(ev *sdl.TextInputEvent) {
//	text := ev.GetText()
//	if sdl.GetModState()&sdl.KMOD_CTRL != 0 && len(text) == 1 && strings.ToLower(text) == "v" {
//		return // ignore "V" from Ctrl-V
//	}
//	win.inputEvent(&seat.TextInputEvent{
//		Text: text,
//	})
//}

func (win *Window) processKeyboardEvent(_ *glfw.Window, key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
	if win.textInp && act == glfw.Press && mods&glfw.ModControl != 0 && key == glfw.KeyV {
		text := glfw.GetClipboardString()
		if text == "" {
			return
		}
		win.inputEvent(&seat.TextInputEvent{
			Text: text,
		})
		return
	}
	k, ok := keymap[key]
	if !ok {
		win.log.Debug("unknown key code", "key", key)
		return
	}
	win.inputEvent(&seat.KeyboardEvent{
		Key:     k,
		Pressed: act == glfw.Press,
	})
}

func (win *Window) processMouseButtonEvent(_ *glfw.Window, btn glfw.MouseButton, act glfw.Action, mods glfw.ModifierKey) {
	pressed := act == glfw.Press
	var button seat.MouseButton
	switch btn {
	case glfw.MouseButtonLeft:
		button = seat.MouseButtonLeft
	case glfw.MouseButtonRight:
		button = seat.MouseButtonRight
	case glfw.MouseButtonMiddle:
		button = seat.MouseButtonMiddle
	default:
		return
	}
	win.inputEvent(&seat.MouseButtonEvent{
		Button:  button,
		Pressed: pressed,
	})
}

func (win *Window) processMotionEvent(_ *glfw.Window, x float64, y float64) {
	win.inputEvent(&seat.MouseMoveEvent{
		Relative: win.rel,
		Pos:      image.Point{X: int(x), Y: int(y)},
		Rel:      types.Pointf{X: float32(x - win.input.prevX), Y: float32(y - win.input.prevY)},
	})
	win.input.prevX, win.input.prevY = x, y
}

func (win *Window) processWheelEvent(_ *glfw.Window, xoff float64, yoff float64) {
	win.inputEvent(&seat.MouseWheelEvent{
		Wheel: int(yoff),
	})
}
