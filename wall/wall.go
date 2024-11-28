package wall

import (
	"image"

	"github.com/opennox/opennox-lib/common"
	"github.com/opennox/opennox-lib/enum"
	"github.com/opennox/opennox-lib/types"
)

const GridStep = common.GridStep

func PosToGrid(pos types.Pointf) image.Point {
	return image.Point{
		X: int(pos.X / GridStep),
		Y: int(pos.Y / GridStep),
	}
}

func GridToPos(pos image.Point) types.Pointf {
	return types.Pointf{
		X: float32(pos.X) * GridStep,
		Y: float32(pos.Y) * GridStep,
	}
}

var _ enum.Enum[Flags] = Flags(0)

type Flags byte

func (f Flags) Has(f2 Flags) bool {
	return f&f2 != 0
}

func (f Flags) HasAny(f2 Flags) bool {
	return f&f2 != 0
}

func (f Flags) Split() []Flags {
	return enum.SplitBits(f)
}

func (f Flags) String() string {
	return enum.StringBits(f, flagNames)
}

func (f Flags) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(f)
}

const (
	Flag1         = Flags(0x1)
	FlagFront     = Flags(0x2)
	FlagSecret    = Flags(0x4)
	FlagBreakable = Flags(0x8)
	FlagDoor      = Flags(0x10)
	FlagBroken    = Flags(0x20)
	FlagWindow    = Flags(0x40)
	Flag8         = Flags(0x80)
)

var flagNames = []string{
	"Flag1",
	"FlagFront",
	"FlagSecret",
	"FlagBreakable",
	"FlagDoor",
	"FlagBroken",
	"FlagWindow",
	"Flag8",
}
