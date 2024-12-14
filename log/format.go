package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"
)

var _ Handler = (*TextHandler)(nil)

func NewTextHandler(w io.Writer) *TextHandler {
	h := &TextHandler{}
	h.SetLevel(slog.LevelInfo)
	h.SetWriter(w)
	return h
}

type TextHandler struct {
	w   atomic.Pointer[io.Writer]
	lvl atomic.Int32
}

func (h *TextHandler) Level() slog.Level {
	return slog.Level(h.lvl.Load())
}

func (h *TextHandler) SetLevel(level slog.Level) {
	h.lvl.Store(int32(level))
}

func (h *TextHandler) SetWriter(w io.Writer) {
	if w != nil {
		h.w.Store(&w)
	} else {
		h.w.Store(nil)
	}
}

func (h *TextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level() && h.w.Load() != nil
}

func (h *TextHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}
	p := h.w.Load()
	if p == nil {
		return nil
	}
	w := *p
	if w == nil {
		return nil
	}
	return PrintRecord(w, PrintDefault, r)
}

type PrintFlags uint32

func (f PrintFlags) Has(f2 PrintFlags) bool {
	return f&f2 != 0
}

const (
	PrintDate = PrintFlags(1 << iota)
	PrintTime
	PrintLevel
	PrintSys
)

const (
	PrintDefault = PrintDate | PrintTime | PrintLevel | PrintSys
)

type bufferWriter interface {
	io.StringWriter
	io.ByteWriter
	io.WriterTo
}

func PrintRecord(w io.Writer, f PrintFlags, r slog.Record) error {
	if w == nil {
		return nil
	}
	buf, isBuf := w.(bufferWriter)
	if !isBuf {
		buf = new(bytes.Buffer)
	}
	if f.Has(PrintDate) {
		buf.WriteString(r.Time.Format(time.DateOnly))
		buf.WriteByte(' ')
	}
	if f.Has(PrintTime) {
		buf.WriteString(r.Time.Format(time.TimeOnly))
		buf.WriteByte(' ')
	}
	if f.Has(PrintLevel) {
		lvl := "INFO"
		switch r.Level {
		case slog.LevelDebug:
			lvl = "DEBUG"
		case slog.LevelInfo:
			lvl = "INFO"
		case slog.LevelWarn:
			lvl = "WARN"
		case slog.LevelError:
			lvl = "ERROR"
		}
		buf.WriteString(lvl)
		buf.WriteByte(' ')
	}
	if f.Has(PrintSys) {
		sys := ""
		r.Attrs(func(a slog.Attr) bool {
			if a.Key == systemAttr {
				sys = a.Value.String()
			}
			return sys == ""
		})
		if sys == "" {
			sys = "main"
		}
		buf.WriteByte('[')
		buf.WriteString(sys)
		buf.WriteString("] ")
	}
	buf.WriteString(strings.TrimSuffix(r.Message, "\n"))
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case systemAttr:
			return true
		}
		buf.WriteByte(' ')
		buf.WriteString(a.Key)
		buf.WriteByte('=')
		buf.WriteString(a.Value.String())
		return true
	})
	buf.WriteByte('\n')
	if isBuf {
		return nil
	}
	_, err := buf.WriteTo(w)
	return err
}
