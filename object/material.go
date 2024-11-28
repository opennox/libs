package object

import (
	"github.com/opennox/libs/enum"
)

var MaterialNames = []string{
	"FLESH", "CLOTH", "ANIMAL_HIDE", "WOOD", "METAL", "STONE", "EARTH", "LIQUID",
	"GLASS", "PAPER", "SNOW", "MUD", "MAGIC", "DIAMOND", "NONE",
}

var _ enum.Enum[Material] = Material(0)

func ParseMaterial(s string) (Material, error) {
	return enum.Parse[Material]("material", s, MaterialNames)
}

func ParseMaterialSet(s string) (Material, error) {
	return enum.ParseSet[Material]("material", s, MaterialNames)
}

type Material uint32

const (
	MaterialFlesh = Material(1 << iota)
	MaterialCloth
	MaterialAnimalHide
	MaterialWood
	MaterialMetal
	MaterialStone
	MaterialEarth
	MaterialLiquid
	MaterialGlass
	MaterialPaper
	MaterialSnow
	MaterialMud
	MaterialMagic
	MaterialDiamond
	MaterialNone
)

func (c Material) Has(c2 Material) bool {
	return c&c2 != 0
}

func (c Material) HasAny(c2 Material) bool {
	return c&c2 != 0
}

func (c Material) Split() []Material {
	return enum.SplitBits(c)
}

func (c Material) String() string {
	return enum.StringBits(c, MaterialNames)
}

func (c Material) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
