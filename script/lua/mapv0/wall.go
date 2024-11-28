package mapv0

import (
	"image"

	lua "github.com/yuin/gopher-lua"

	"github.com/opennox/libs/script"
	"github.com/opennox/libs/types"
)

type metaWall struct {
	Wall      *lua.LTable
	WallGroup *lua.LTable
}

func (vm *api) newWall(v script.Wall) lua.LValue {
	if v == nil {
		return lua.LNil
	}
	return &lua.LUserData{Value: v, Metatable: vm.meta.Wall}
}

func (vm *api) newWallGroup(v *script.WallGroup) lua.LValue {
	if v == nil {
		return lua.LNil
	}
	return &lua.LUserData{Value: v, Metatable: vm.meta.WallGroup}
}

func (vm *api) initMetaWall() {
	vm.meta.Wall = vm.newMeta("Wall")
	vm.meta.WallGroup = vm.newMeta("WallGroup")

	vm.registerObjMethod("GridPos", func(obj script.GridPositioner) (x, y int) {
		pos := obj.GridPos()
		return pos.X, pos.Y
	})
}

func (vm *api) initWall() {
	// Wall(x, y number)
	vm.meta.Wall.RawSetString("__call", vm.s.NewFunction(func(s *lua.LState) int {
		x := int(s.CheckNumber(2))
		y := int(s.CheckNumber(3))
		wl := vm.g.WallAtGrid(image.Point{X: x, Y: y})
		s.Push(vm.newWall(wl))
		return 1
	}))
	// Nox.WallAt(obj Object)
	// Nox.WallAt(x, y number)
	vm.root.RawSetString("WallAt", vm.s.NewFunction(func(s *lua.LState) int {
		var wl script.Wall
		switch s.Get(1).(type) {
		case lua.LNumber:
			x := float32(s.CheckNumber(1))
			y := float32(s.CheckNumber(2))
			wl = vm.g.WallAt(types.Pointf{X: x, Y: y})
		default:
			obj2, ok := s.CheckUserData(1).Value.(script.Positioner)
			if !ok {
				return 0
			}
			wl = vm.g.WallAt(obj2.Pos())
		}
		s.Push(vm.newWall(wl))
		return 1
	}))
	// Nox.WallNear(obj Object)
	// Nox.WallNear(x, y number)
	vm.root.RawSetString("WallNear", vm.s.NewFunction(func(s *lua.LState) int {
		var wl script.Wall
		switch s.Get(1).(type) {
		case lua.LNumber:
			x := float32(s.CheckNumber(1))
			y := float32(s.CheckNumber(2))
			wl = vm.g.WallNear(types.Pointf{X: x, Y: y})
		default:
			obj2, ok := s.CheckUserData(1).Value.(script.Positioner)
			if !ok {
				return 0
			}
			wl = vm.g.WallNear(obj2.Pos())
		}
		s.Push(vm.newWall(wl))
		return 1
	}))
	// Wall[key]
	vm.setIndexFunction(vm.meta.Wall, nil)
	// Wall[key] = v
	vm.setSetIndexFunction(vm.meta.Wall, nil)
}

func (vm *api) initWallGroup() {
	// WallGroup(id string)
	vm.meta.WallGroup.RawSetString("__call", vm.s.NewFunction(func(s *lua.LState) int {
		id := s.CheckString(2)
		g := vm.g.WallGroupByID(id)
		s.Push(vm.newWallGroup(g))
		return 1
	}))
	// WallGroup[key]
	vm.meta.WallGroup.RawSetString("__index", vm.s.NewFunction(func(s *lua.LState) int {
		val := s.CheckUserData(1).Value
		switch s.Get(2).(type) {
		case lua.LNumber:
			i := int(s.CheckNumber(2))
			if i <= 0 {
				return 0
			}
			i--
			obj, ok := val.(*script.WallGroup)
			if !ok {
				return 0
			}
			list := obj.Walls()
			if i >= len(list) {
				return 0
			}
			w := list[i]
			s.Push(vm.newWall(w))
			return 1
		default:
		}
		key := s.CheckString(2)
		if v, ok := vm.indexInterfaceV0(val, key); ok {
			s.Push(v)
			return 1
		}
		obj, ok := val.(*script.WallGroup)
		if !ok {
			return 0
		}
		_ = obj
		switch key {
		default:
			s.Push(s.RawGet(vm.meta.WallGroup, lua.LString(key)))
			return 1
		}
	}))
	// WallGroup[key] = v
	vm.setSetIndexFunction(vm.meta.WallGroup, nil)
}
