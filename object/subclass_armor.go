package object

import (
	"github.com/opennox/libs/enum"
)

var ArmorClassNames = []string{
	"HELMET", "SHIELD", "BREASTPLATE", "ARM_ARMOR", "PANTS", "BOOTS", "SHIRT", "LEG_ARMOR", "BACK",
}

func (c SubClass) AsArmor() ArmorClass {
	return ArmorClass(c)
}

func ParseArmorClass(s string) (ArmorClass, error) {
	return enum.Parse[ArmorClass]("armor class", s, ArmorClassNames)
}

func ParseArmorClassSet(s string) (ArmorClass, error) {
	return enum.ParseSet[ArmorClass]("armor class", s, ArmorClassNames)
}

var _ enum.Enum[ArmorClass] = ArmorClass(0)

type ArmorClass uint32

const (
	ArmorHelmet      = ArmorClass(1 << iota) // 0x1
	ArmorShield                              // 0x2
	ArmorBreastplate                         // 0x4
	ArmorArmArmor                            // 0x8
	ArmorPants                               // 0x10
	ArmorBoots                               // 0x20
	ArmorShirt                               // 0x40
	ArmorLegArmor                            // 0x80
	ArmorBack                                // 0x100
)

func (c ArmorClass) Has(c2 ArmorClass) bool {
	return c&c2 != 0
}

func (c ArmorClass) HasAny(c2 ArmorClass) bool {
	return c&c2 != 0
}

func (c ArmorClass) Split() []ArmorClass {
	return enum.SplitBits(c)
}

func (c ArmorClass) String() string {
	return enum.StringBits(c, ArmorClassNames)
}

func (c ArmorClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
