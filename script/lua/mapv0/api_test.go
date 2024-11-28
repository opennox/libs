package mapv0_test

import (
	"image"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	glua "github.com/yuin/gopher-lua"

	"github.com/opennox/libs/script"
	"github.com/opennox/libs/script/lua"
	"github.com/opennox/libs/script/lua/mapv0"
)

func newGame(t testing.TB) *testGame {
	g := &testGame{t: t, frame: 1, time: time.Second}
	lvm := lua.NewVM(nil, g, "", glua.Options{
		IncludeGoStackTrace: true,
	})
	lvm.InitAPI("Nox", mapv0.New)
	g.vm = lvm
	return g
}

type testGame struct {
	script.BaseGame
	vm          *lua.VM
	t           testing.TB
	frame       int
	time        time.Duration
	walls       map[image.Point]script.Wall
	waypoints   map[string]script.Waypoint
	objectTypes map[string]script.ObjectType
	objects     map[string]script.Object
	players     []script.Player
	host        script.Player
}

func (g *testGame) Tick() {
	g.frame++
	g.time += 33 * time.Millisecond
	g.vm.OnFrame()
}

func (g *testGame) Frame() int {
	return g.frame
}

func (g *testGame) Time() time.Duration {
	return g.time
}

func (g *testGame) Exec(s string) {
	_, err := g.vm.Exec(s)
	must.NoError(g.t, err)
}
