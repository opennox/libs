package object

import (
	"github.com/opennox/libs/enum"
)

var MissileClassNames = []string{
	"MISSILE_COUNTERSPELL",
	"MAGIC",
}

func (c SubClass) AsMissile() MissileClass {
	return MissileClass(c)
}

func ParseMissileClass(s string) (MissileClass, error) {
	return enum.Parse[MissileClass]("missile class", s, MissileClassNames)
}

func ParseMissileClassSet(s string) (MissileClass, error) {
	return enum.ParseSet[MissileClass]("missile class", s, MissileClassNames)
}

var _ enum.Enum[MissileClass] = MissileClass(0)

type MissileClass uint32

const (
	MissileMissileCounterSpell = MissileClass(1 << iota) // 0x1
	MissileMagic                                         // 0x2
)

func (c MissileClass) Has(c2 MissileClass) bool {
	return c&c2 != 0
}

func (c MissileClass) HasAny(c2 MissileClass) bool {
	return c&c2 != 0
}

func (c MissileClass) Split() []MissileClass {
	return enum.SplitBits(c)
}

func (c MissileClass) String() string {
	return enum.StringBits(c, MissileClassNames)
}

func (c MissileClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
