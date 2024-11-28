package eval

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/traefik/yaegi/interp"

	eudeval "github.com/noxworld-dev/noxscript/eud/v171/eval"
	nseval3 "github.com/noxworld-dev/noxscript/ns/v3/eval"
	ns3vm "github.com/noxworld-dev/noxscript/ns/v3/vm"
	ns4 "github.com/noxworld-dev/noxscript/ns/v4"
	nseval4 "github.com/noxworld-dev/noxscript/ns/v4/eval"

	"github.com/opennox/libs/script"
	"github.com/opennox/libs/script/eval/imports"
	"github.com/opennox/libs/script/eval/stdlib"
)

var useSymbols []Exports

type Exports struct {
	Symbols interp.Exports
	Context func(ctx context.Context, g script.Game) context.Context
}

// Register additional libraries for the interpreter.
func Register(m Exports) {
	useSymbols = append(useSymbols, m)
}

func init() {
	script.RegisterVM(script.VMRuntime{
		Name: "eval", Title: "Go interpreter",
		NewMap: func(log *slog.Logger, g script.Game, maps string, name string) (script.VM, error) {
			vm := NewVM(g, VMOptions{
				Log:     log,
				MapsDir: maps,
			})
			err := vm.ExecFile(name)
			if os.IsNotExist(err) {
				log.Warn("no go files")
				return vm, nil // still run the empty script for commands
			} else if err != nil {
				return nil, err
			}
			log.Info("go scripts loaded")
			return vm, nil
		},
	})
}

type scriptPrinter struct {
	buf bytes.Buffer
	p   script.Printer
}

func (p *scriptPrinter) Write(data []byte) (int, error) {
	p.buf.Write(data)
	for p.buf.Len() > 0 {
		data := p.buf.Bytes()
		i := bytes.IndexByte(data, '\n')
		if i < 0 {
			break
		}
		p.p.Print(string(p.buf.Next(i + 1)))
	}
	return len(data), nil
}

var _ script.VM = (*VM)(nil)

type VM struct {
	log      *slog.Logger
	sclog    *slog.Logger
	fs       FS
	vm       *interp.Interpreter
	printers []*scriptPrinter
	defs     bool

	curExports map[string]reflect.Value

	onFrame    func()
	onEvent    func(typ script.EventType)
	onEventStr func(typ string)
}

type VMOptions struct {
	Log       *slog.Logger
	ScriptLog *slog.Logger
	MapsDir   string
	ModsDir   string
}

func NewVM(g script.Game, opts VMOptions) *VM {
	log := opts.Log
	if log == nil {
		log = slog.Default()
	}
	mapsDir := opts.MapsDir
	modsDir := opts.ModsDir
	if modsDir == "" {
		modsDir = filepath.Join(mapsDir, "..", "mods", "pkg", "mod") // follows go mod structure for /mods/
	}
	stdout := &scriptPrinter{p: g.Console(false)}
	stderr := &scriptPrinter{p: g.Console(true)}
	sclog := opts.ScriptLog
	if sclog == nil {
		sclog = slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case "time":
					return slog.Attr{}
				}
				return a
			},
		}))
	}
	vm := &VM{
		log:   log,
		sclog: sclog,
		fs:    newModFS(log, mapsDir, modsDir),
	}
	vm.vm = interp.New(interp.Options{
		GoPath:               gopath,
		BuildTags:            []string{"script"},
		Stdin:                bytes.NewReader(nil),
		Stdout:               stdout,
		Stderr:               stderr,
		Args:                 []string{},
		SourcecodeFilesystem: vm.fs,
	})
	vm.addPrinters(stdout, stderr)
	vm.initPackages(g)
	return vm
}

func (vm *VM) addPrinters(list ...*scriptPrinter) {
	for _, p := range list {
		vm.printers = append(vm.printers, p)
	}
}

func (vm *VM) flushPrinters() {
	for _, p := range vm.printers {
		if p.buf.Len() == 0 {
			continue
		}
		p.p.Print(p.buf.String())
		p.buf.Reset()
	}
}

func importPathFor(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath()
}

func (vm *VM) initPackages(g script.Game) {
	ctx := context.Background()

	vm.vm.Use(stdlib.Symbols)
	vm.vm.Use(imports.Symbols)
	vm.vm.Use(nseval3.Symbols)
	vm.vm.Use(nseval4.Symbols)
	vm.vm.Use(eudeval.Symbols)
	ctx = script.WithGame(ctx, g)

	for _, m := range useSymbols {
		if m.Symbols != nil {
			vm.vm.Use(m.Symbols)
		}
		if m.Context != nil {
			ctx = m.Context(ctx, g)
		}
	}
	// These strange assignments below statically check that types of the function and override are exactly the same.
	slogDefault := slog.Default
	slogWith := slog.With
	slogLog := slog.Log
	slogLogAttrs := slog.LogAttrs
	slogDebug := slog.Debug
	slogInfo := slog.Info
	slogWarn := slog.Warn
	slogError := slog.Error
	slogDebugContext := slog.DebugContext
	slogInfoContext := slog.InfoContext
	slogWarnContext := slog.WarnContext
	slogErrorContext := slog.ErrorContext
	slogDefault = func() *slog.Logger {
		return vm.sclog
	}
	slogWith = func(args ...any) *slog.Logger {
		return vm.sclog.With(args...)
	}
	slogLog = func(ctx context.Context, level slog.Level, msg string, args ...any) {
		vm.sclog.Log(ctx, level, msg, args...)
	}
	slogLogAttrs = func(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
		vm.sclog.LogAttrs(ctx, level, msg, attrs...)
	}
	slogDebug = func(msg string, args ...any) {
		vm.sclog.Debug(msg, args...)
	}
	slogInfo = func(msg string, args ...any) {
		vm.sclog.Info(msg, args...)
	}
	slogWarn = func(msg string, args ...any) {
		vm.sclog.Warn(msg, args...)
	}
	slogError = func(msg string, args ...any) {
		vm.sclog.Error(msg, args...)
	}
	slogDebugContext = func(ctx context.Context, msg string, args ...any) {
		vm.sclog.DebugContext(ctx, msg, args...)
	}
	slogInfoContext = func(ctx context.Context, msg string, args ...any) {
		vm.sclog.InfoContext(ctx, msg, args...)
	}
	slogWarnContext = func(ctx context.Context, msg string, args ...any) {
		vm.sclog.WarnContext(ctx, msg, args...)
	}
	slogErrorContext = func(ctx context.Context, msg string, args ...any) {
		vm.sclog.ErrorContext(ctx, msg, args...)
	}
	vm.vm.Use(interp.Exports{
		importPathFor((*slog.Logger)(nil)) + "/slog": {
			"Default":      reflect.ValueOf(slogDefault),
			"With":         reflect.ValueOf(slogWith),
			"Log":          reflect.ValueOf(slogLog),
			"LogAttrs":     reflect.ValueOf(slogLogAttrs),
			"Debug":        reflect.ValueOf(slogDebug),
			"Info":         reflect.ValueOf(slogInfo),
			"Warn":         reflect.ValueOf(slogWarn),
			"Error":        reflect.ValueOf(slogError),
			"DebugContext": reflect.ValueOf(slogDebugContext),
			"InfoContext":  reflect.ValueOf(slogInfoContext),
			"WarnContext":  reflect.ValueOf(slogWarnContext),
			"ErrorContext": reflect.ValueOf(slogErrorContext),
		},
	})

	// Rename map entry below when refactoring!
	getGame := script.Runtime
	getGame = func() script.Game {
		return g
	}
	vm.vm.Use(interp.Exports{
		importPathFor((*script.Game)(nil)) + "/script": {
			"Runtime": reflect.ValueOf(getGame),
		},
	})
	getCtx := context.Background
	getCtx = func() context.Context {
		return ctx
	}
	vm.vm.Use(interp.Exports{
		importPathFor((*context.Context)(nil)) + "/context": {
			"Background": reflect.ValueOf(getCtx),
		},
	})
	// TODO: properly virtualize
	script.SetRuntime(g)
	if v, ok := g.(ns3vm.Game); ok {
		ns3vm.SetRuntime(v.NoxScriptVM())
	}
	if v, ok := g.(ns4.Game); ok {
		ns4.SetRuntime(v.NoxScript())
	}
}

func (vm *VM) initMain() {
	if _, err := vm.vm.Eval(`
package main

import (
	"fmt"

	"github.com/opennox/libs/types"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/script"

	ns3 "github.com/noxworld-dev/noxscript/ns/v3"
	ns4 "github.com/noxworld-dev/noxscript/ns/v4"
	eud "github.com/noxworld-dev/noxscript/eud/v171"
)

var Game = script.Runtime()
`); err != nil {
		panic(err)
	}
}

func (vm *VM) exportByName(name string) (reflect.Value, bool) {
	rv, ok := vm.curExports[name]
	return rv, ok
}

func maybeSetExport[T any](vm *VM, p *T, name string) {
	v, ok := vm.exportByName(name)
	if !ok || v.Kind() != reflect.TypeOf(p).Elem().Kind() {
		return
	}
	*p, _ = v.Interface().(T)
}

func (vm *VM) checkExports() {
	if vm.onFrame == nil {
		maybeSetExport(vm, &vm.onFrame, "OnFrame")
	}
	if vm.onEvent == nil && vm.onEventStr == nil {
		maybeSetExport(vm, &vm.onEvent, "OnEvent")
		maybeSetExport(vm, &vm.onEventStr, "OnEvent")
	}
}

func (vm *VM) Exec(s string) (reflect.Value, error) {
	if !vm.defs {
		vm.initMain()
		vm.defs = true
	}
	defer vm.checkExports()
	defer vm.flushPrinters()
	return vm.vm.Eval(s)
}

func (vm *VM) ExecFile(pkg string) error {
	defer vm.checkExports()
	defer vm.flushPrinters()
	log := vm.log.With("pkg", pkg)
	_, err := vm.vm.EvalPath(pkg)
	if err != nil {
		if strings.Contains(err.Error(), "no Go files") {
			return os.ErrNotExist
		}
		log.Error("cannot load package", "err", err)
		return err
	}
	vm.curExports = vm.vm.Symbols(pkg)[pkg]
	if len(vm.curExports) == 0 {
		log.Warn("no exports from go; wrong package name?")
	}
	return nil
}

func (vm *VM) OnFrame() {
	if vm.onFrame == nil {
		return
	}
	defer vm.flushPrinters()
	vm.onFrame()
}

func (vm *VM) OnEvent(typ script.EventType) {
	if vm.onEvent == nil && vm.onEventStr == nil {
		return
	}
	defer vm.flushPrinters()
	if vm.onEvent != nil {
		vm.onEvent(typ)
	} else if vm.onEventStr != nil {
		vm.onEventStr(string(typ))
	}
}

func (vm *VM) GetSymbol(name string, typ reflect.Type) (reflect.Value, bool, error) {
	rv, ok := vm.exportByName(name)
	if !ok {
		return reflect.Value{}, false, nil
	} else if rv.Type() != typ {
		return reflect.Value{}, false, fmt.Errorf("unexpected type of %q: expected %v, got %v", name, typ, rv.Type())
	}
	return rv, true, nil
}

func (vm *VM) Close() error {
	vm.flushPrinters()
	*vm = VM{}
	return nil
}
