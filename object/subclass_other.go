package object

import (
	"github.com/opennox/libs/enum"
)

var OtherClassNames = []string{
	"HEAVY", "LAVA", "GATE", "VISIBLE_OBELISK", "INVISIBLE_OBELISK", "TECH", "LOTD",
	"USEABLE", "CHEST_NW", "CHEST_NE", "CHEST_SE", "CHEST_SW", "STONE_DOOR",
}

func (c SubClass) AsOther() OtherClass {
	return OtherClass(c)
}

func ParseOtherClass(s string) (OtherClass, error) {
	return enum.Parse[OtherClass]("other class", s, OtherClassNames)
}

func ParseOtherClassSet(s string) (OtherClass, error) {
	return enum.ParseSet[OtherClass]("other class", s, OtherClassNames)
}

var _ enum.Enum[OtherClass] = OtherClass(0)

type OtherClass uint32

const (
	OtherHeavy            = OtherClass(1 << iota) // 0x1
	OtherLava                                     // 0x2
	OtherGate                                     // 0x4
	OtherVisibleObelisk                           // 0x8
	OtherInvisibleObelisk                         // 0x10
	OtherTech                                     // 0x20
	OtherLOTD                                     // 0x40
	OtherUseable                                  // 0x80
	OtherChestNW                                  // 0x100
	OtherChestNE                                  // 0x200
	OtherChestSE                                  // 0x400
	OtherChestSW                                  // 0x800
	OtherStoneDoor                                // 0x1000
)

func (c OtherClass) Has(c2 OtherClass) bool {
	return c&c2 != 0
}

func (c OtherClass) HasAny(c2 OtherClass) bool {
	return c&c2 != 0
}

func (c OtherClass) Split() []OtherClass {
	return enum.SplitBits(c)
}

func (c OtherClass) String() string {
	return enum.StringBits(c, OtherClassNames)
}

func (c OtherClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
