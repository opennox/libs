package object

import (
	"github.com/opennox/libs/enum"
)

var GeneratorClassNames = []string{
	"GENERATOR_NW",
	"GENERATOR_NE",
	"GENERATOR_SE",
	"GENERATOR_SW",
}

func (c SubClass) AsGenerator() GeneratorClass {
	return GeneratorClass(c)
}

func ParseGeneratorClass(s string) (GeneratorClass, error) {
	return enum.Parse[GeneratorClass]("generator class", s, GeneratorClassNames)
}

func ParseGeneratorClassSet(s string) (GeneratorClass, error) {
	return enum.ParseSet[GeneratorClass]("generator class", s, GeneratorClassNames)
}

var _ enum.Enum[GeneratorClass] = GeneratorClass(0)

type GeneratorClass uint32

const (
	GeneratorNW = GeneratorClass(1 << iota) // 0x1
	GeneratorNE                             // 0x2
	GeneratorSE                             // 0x4
	GeneratorSW                             // 0x8
)

func (c GeneratorClass) Has(c2 GeneratorClass) bool {
	return c&c2 != 0
}

func (c GeneratorClass) HasAny(c2 GeneratorClass) bool {
	return c&c2 != 0
}

func (c GeneratorClass) Split() []GeneratorClass {
	return enum.SplitBits(c)
}

func (c GeneratorClass) String() string {
	return enum.StringBits(c, GeneratorClassNames)
}

func (c GeneratorClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
