package object

import (
	"github.com/opennox/libs/enum"
)

var BookClassNames = []string{
	"SPELL_BOOK",
	"FIELD_GUIDE",
	"ABILITY_BOOK",
}

func (c SubClass) AsBook() BookClass {
	return BookClass(c)
}

func ParseBookClass(s string) (BookClass, error) {
	return enum.Parse[BookClass]("book class", s, BookClassNames)
}

func ParseBookClassSet(s string) (BookClass, error) {
	return enum.ParseSet[BookClass]("book class", s, BookClassNames)
}

var _ enum.Enum[BookClass] = BookClass(0)

type BookClass uint32

const (
	BookSpell      = BookClass(1 << iota) // 0x1
	BookFieldGuide                        // 0x2
	BookAbility                           // 0x4
)

func (c BookClass) Has(c2 BookClass) bool {
	return c&c2 != 0
}

func (c BookClass) HasAny(c2 BookClass) bool {
	return c&c2 != 0
}

func (c BookClass) Split() []BookClass {
	return enum.SplitBits(c)
}

func (c BookClass) String() string {
	return enum.StringBits(c, BookClassNames)
}

func (c BookClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
