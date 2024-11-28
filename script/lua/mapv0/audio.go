package mapv0

import (
	lua "github.com/yuin/gopher-lua"

	"github.com/opennox/libs/script"
)

type metaAudio struct {
	Audio *lua.LTable
}

func (vm *api) initMetaAudio() {
	vm.meta.Audio = vm.newMeta("Audio")
}

func (vm *api) initAudio() {
	// Nox.Audio.Effect("name", pos)
	vm.newFuncOn(vm.meta.Audio, "Effect", func(name string, pos script.Positioner) {
		vm.g.AudioEffect(name, pos)
	})
}
