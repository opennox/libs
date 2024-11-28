package object

import (
	"github.com/opennox/libs/enum"
)

var SubClassNames = [][]string{
	ArmorClassNames,
	WeaponClassNames,
	MonsterClassNames,
	FoodClassNames,
	OtherClassNames,
	BookClassNames,
	MissileClassNames,
	GeneratorClassNames,
	ExitClassNames,
}

func ParseSubClass(s string) (SubClass, error) {
	return enum.ParseMulti[SubClass]("subclass", s, SubClassNames)
}

func ParseSubClassSet(s string) (SubClass, error) {
	return enum.ParseSetMulti[SubClass]("subclass", s, SubClassNames)
}

var _ enum.Enum[SubClass] = SubClass(0)

type SubClass uint32

func (c SubClass) Has(c2 SubClass) bool {
	return c&c2 != 0
}

func (c SubClass) HasAny(c2 SubClass) bool {
	return c&c2 != 0
}

func (c SubClass) Split() []SubClass {
	return enum.SplitBits(c)
}

func (c SubClass) String() string {
	return enum.StringBitsRaw(c)
}

func (c SubClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
