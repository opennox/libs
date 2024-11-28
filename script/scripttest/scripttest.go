package scripttest

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"testing"
	"time"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/script"
	"github.com/opennox/libs/types"
)

type Object struct {
	SID    int
	IDVal  string
	PosVal types.Pointf
}

func (v *Object) ScriptID() int {
	return v.SID
}

func (v *Object) ObjScriptID() int {
	return v.SID
}

func (v *Object) String() string {
	panic("implement me")
}

func (v *Object) GetObject() script.Object {
	panic("implement me")
}

func (v *Object) Class() object.Class {
	panic("implement me")
}

func (v *Object) Owner() script.Object {
	panic("implement me")
}

func (v *Object) SetOwner(owner script.ObjectWrapper) {
	panic("implement me")
}

func (v *Object) Destroy() {
	panic("implement me")
}

func (v *Object) OnTriggerActivate(fnc func(obj script.Object)) {
	panic("implement me")
}

func (v *Object) OnTriggerDeactivate(fnc func()) {
	panic("implement me")
}

func (v *Object) ObjectType() script.ObjectType {
	panic("implement me")
}

func (v *Object) Z() float32 {
	panic("implement me")
}

func (v *Object) SetZ(z float32) {
	panic("implement me")
}

func (v *Object) Push(vec types.Pointf, force float32) {
	panic("implement me")
}

func (v *Object) PushTo(p types.Pointf) {
	panic("implement me")
}

func (v *Object) IsEnabled() bool {
	panic("implement me")
}

func (v *Object) ID() string {
	return v.IDVal
}

func (v *Object) Type() script.ObjectType {
	panic("implement me")
}

func (v *Object) Pos() types.Pointf {
	return v.PosVal
}

func (v *Object) SetPos(p types.Pointf) {
	v.PosVal = p
}

func (v *Object) Enable(val bool) {
	panic("implement me")
}

func (v *Object) Delete() {
	panic("implement me")
}

type testPrinter struct {
	t   testing.TB
	lvl string
	w   io.Writer
}

func (p testPrinter) Print(text string) {
	if p.w != nil {
		fmt.Fprintf(p.w, "[%s] %s", p.lvl, text)
	}
	p.t.Logf("[%s] %s", p.lvl, text)
}

type Game struct {
	T   testing.TB
	Log bytes.Buffer
}

func (g *Game) BlindPlayers(blind bool) {
	panic("implement me")
}

func (g *Game) CinemaPlayers(v bool) {
	panic("implement me")
}

func (g *Game) OnPlayerJoin(fnc func(p script.Player)) {
	panic("implement me")
}

func (g *Game) OnPlayerLeave(fnc func(p script.Player)) {
	panic("implement me")
}

func (g *Game) ObjectGroupByID(id string) *script.ObjectGroup {
	panic("implement me")
}

func (g *Game) WaypointGroupByID(id string) *script.WaypointGroup {
	panic("implement me")
}

func (g *Game) WallAt(pos types.Pointf) script.Wall {
	panic("implement me")
}

func (g *Game) WallNear(pos types.Pointf) script.Wall {
	panic("implement me")
}

func (g *Game) WallAtGrid(pos image.Point) script.Wall {
	panic("implement me")
}

func (g *Game) WallGroupByID(id string) *script.WallGroup {
	panic("implement me")
}

func (g *Game) Console(error bool) script.Printer {
	lvl := "info"
	if error {
		lvl = "error"
	}
	return testPrinter{t: g.T, lvl: lvl, w: &g.Log}
}

func (g *Game) Frame() int {
	return 1
}

func (g *Game) Time() time.Duration {
	return time.Second
}

func (g *Game) Players() []script.Player {
	panic("implement me")
}

func (g *Game) HostPlayer() script.Player {
	panic("implement me")
}

func (g *Game) WaypointByID(id string) script.Waypoint {
	panic("implement me")
}

func (g *Game) WallByPos(pos image.Point) script.Wall {
	panic("implement me")
}

func (g *Game) AudioEffect(name string, pos script.Positioner) {
	panic("implement me")
}

func (g *Game) Global() script.Printer {
	return testPrinter{t: g.T, lvl: "global"}
}

func (g *Game) ObjectTypeByID(id string) script.ObjectType {
	return nil
}

func (g *Game) ObjectByID(id string) script.Object {
	switch id {
	case "Frog":
		return &Object{PosVal: types.Pointf{1, 2}, IDVal: id}
	case "Cookie":
		return &Object{PosVal: types.Pointf{3, 4}, IDVal: id}
	}
	return nil
}
