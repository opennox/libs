package eval

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/script"
	"github.com/opennox/libs/script/scripttest"
)

func TestScriptDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "opennox_eval_*")
	must.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	writeFile := func(path string, src string) {
		path = filepath.Join(dir, path)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		must.NoError(t, err)
		err = os.WriteFile(path, []byte(src), 0644)
		must.NoError(t, err)
	}

	// Map script for "themap"
	writeFile("maps/themap/script.go", `
package themap

import (
	"log/slog"

	"othermap"

	"test.org/vend"
	"test.org/share"
	"test.org/themod"
)

var Cnt int

func init() { 
	println("init map")
	Cnt++
}
func Do() { 
	println("call map")
	slog.Info("call", "src", "map")
	othermap.Do()
	vend.Do()
	share.Do()
	themod.Do()
	Cnt++
}
func OnFrame() {
	println("frame")
	Cnt++
}

func Panic1() {
	panic2()
}
func panic2() {
	func(){
		panic3()
	}()
}
func panic3() {
	panic("test")
}
`)

	// Vendored package "test.org/vend" in "themap"
	writeFile("maps/themap/vendor/test.org/vend/script.go", `
package vend

func init() { println("init vendor") }
func Do() { println("call vendor") }
`)

	// Second map script "othermap"
	writeFile("maps/othermap/script.go", `
package othermap

func init() { println("init other map") }
func Do() { println("call other map") }
`)

	// Shared module "test.org/themod" in mods
	writeFile("mods/pkg/mod/test.org/themod/script.go", `
package themod

func init() { println("init module") }
func Do() { println("call module") }
`)

	// Legacy way to share package "test.org/share"
	writeFile("maps/test.org/share/script.go", `
package share

func init() { println("init shared") }
func Do() { println("call shared") }
`)
	g := &scripttest.Game{T: t}
	vm := NewVM(g, VMOptions{
		MapsDir: filepath.Join(dir, "maps"),
	})
	err = vm.ExecFile("themap")
	must.NoError(t, err)

	rv, ok := vm.exportByName(`Do`)
	must.True(t, ok)
	fnc, ok := rv.Interface().(func())
	must.True(t, ok)
	fnc()

	_, err = vm.Exec(`println("builtin")`)
	must.NoError(t, err)
	_, err = vm.Exec(`fmt.Println("fmt")`)
	must.NoError(t, err)
	_, err = vm.Exec(`Game.Global().Print("global")`)
	must.NoError(t, err)
	vm.OnFrame()

	rv, ok = vm.exportByName(`Cnt`)
	must.True(t, ok)
	p, ok := rv.Interface().(int)
	must.True(t, ok)
	must.EqOp(t, 3, p)

	one, err := script.GetVMSymbol[func()](vm, "Do")
	must.NoError(t, err)
	one()

	cnt, err := script.GetVMSymbolPtr[int](vm, "Cnt")
	must.NoError(t, err)
	must.EqOp(t, 4, *cnt)

	panic1, err := script.GetVMSymbol[func()](vm, "Panic1")
	must.NoError(t, err)
	func() {
		defer func() {
			recover()
		}()
		panic1()
	}()
	must.EqOp(t, strings.TrimSpace(`
[info] init other map
[info] init vendor
[info] init shared
[info] init module
[info] init map
[info] call map
[error] level=INFO msg=call src=map
[info] call other map
[info] call vendor
[info] call shared
[info] call module
[info] builtin
[info] fmt
[info] frame
[info] call map
[error] level=INFO msg=call src=map
[info] call other map
[info] call vendor
[info] call shared
[info] call module
[error] /src/themap/script.go:43:2: panic: themap.panic3(...)
[error] /src/themap/script.go:39:3: panic: themap.panic2.func(...)
[error] /src/themap/script.go:38:2: panic: themap.panic2.func(...)
[error] /src/themap/script.go:35:2: panic: themap.Panic1(...)
`), strings.TrimSpace(g.Log.String()))
}
