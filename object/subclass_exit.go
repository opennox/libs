package object

import (
	"github.com/opennox/libs/enum"
)

var ExitClassNames = []string{
	"QUEST_EXIT",
	"QUEST_WARP_EXIT",
}

func (c SubClass) AsExit() ExitClass {
	return ExitClass(c)
}

func ParseExitClass(s string) (ExitClass, error) {
	return enum.Parse[ExitClass]("exit class", s, ExitClassNames)
}

func ParseExitClassSet(s string) (ExitClass, error) {
	return enum.ParseSet[ExitClass]("exit class", s, ExitClassNames)
}

var _ enum.Enum[ExitClass] = ExitClass(0)

type ExitClass uint32

const (
	ExitQuest     = ExitClass(1 << iota) // 0x1
	ExitQuestWarp                        // 0x2
)

func (c ExitClass) Has(c2 ExitClass) bool {
	return c&c2 != 0
}

func (c ExitClass) HasAny(c2 ExitClass) bool {
	return c&c2 != 0
}

func (c ExitClass) Split() []ExitClass {
	return enum.SplitBits(c)
}

func (c ExitClass) String() string {
	return enum.StringBits(c, ExitClassNames)
}

func (c ExitClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
