package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
)

const (
	systemAttr = "sys"
)

var (
	defText = NewTextHandler(os.Stderr)
	defHnd  = NewMulti(defText)
)

func init() {
	slog.SetDefault(slog.New(NewHandler(defHnd)))
}

func DefaultHandler() *MultiHandler {
	return defHnd
}

func DefaultTextHandler() *TextHandler {
	return defText
}

type Logger struct {
	*slog.Logger
}

func (l Logger) Printf(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l Logger) Print(args ...interface{}) {
	l.Info(fmt.Sprintln(args...))
}

func (l Logger) Println(args ...interface{}) {
	l.Info(fmt.Sprintln(args...))
}

func (l Logger) WithSystem(name string) *Logger {
	l2 := l
	l2.Logger = WithSystem(l.Logger, name)
	return &l2
}

func New(name string) *Logger {
	if name == "" {
		name = "main"
	}
	log := slog.Default()
	log = WithSystem(log, name)
	return NewSlog(log)
}

func NewSlog(log *slog.Logger) *Logger {
	if log == nil {
		log = slog.Default()
	}
	return &Logger{log}
}

func WithSystem(log *slog.Logger, name string) *slog.Logger {
	return log.With(systemAttr, name)
}

func Printf(format string, args ...interface{}) {
	Logger{slog.Default()}.Printf(format, args...)
}

func Println(args ...interface{}) {
	Logger{slog.Default()}.Println(args...)
}

func AddHandler(h Handler) {
	defHnd.AddHandler(h)
}

func RemoveHandler(h Handler) {
	defHnd.RemoveHandler(h)
}

type Handler interface {
	Enabled(ctx context.Context, level slog.Level) bool
	Handle(ctx context.Context, r slog.Record) error
}

func NewHandler(h Handler) slog.Handler {
	return &slogHandler{h: h}
}

type slogHandler struct {
	h     Handler
	attrs []slog.Attr
}

func (m *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.h.Enabled(ctx, level)
}

func (m *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	r2 := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r2.AddAttrs(m.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		r2.AddAttrs(a)
		return true
	})
	return m.h.Handle(ctx, r2)
}

func (m *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := slices.Clip(m.attrs)
	out = append(out, attrs...)
	return &slogHandler{h: m.h, attrs: out}
}

func (m *slogHandler) WithGroup(name string) slog.Handler {
	return m // TODO
}
