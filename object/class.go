package object

import (
	"github.com/opennox/libs/enum"
)

var ClassNames = []string{
	"MISSILE", "MONSTER", "PLAYER", "OBSTACLE", "FOOD", "EXIT", "KEY",
	"DOOR", "INFO_BOOK", "TRIGGER", "TRANSPORTER", "HOLE", "WAND", "FIRE",
	"ELEVATOR", "ELEVATOR_SHAFT", "DANGEROUS", "MONSTERGENERATOR", "READABLE", "LIGHT", "SIMPLE",
	"COMPLEX", "IMMOBILE", "VISIBLE_ENABLE", "WEAPON", "ARMOR", "NOT_STACKABLE", "TREASURE",
	"FLAG", "CLIENT_PERSIST", "CLIENT_PREDICT", "PICKUP",
}

var goClassNames = []string{
	"ClassMissile",
	"ClassMonster",
	"ClassPlayer",
	"ClassObstacle",
	"ClassFood",
	"ClassExit",
	"ClassKey",
	"ClassDoor",
	"ClassInfoBook",
	"ClassTrigger",
	"ClassTransporter",
	"ClassHole",
	"ClassWand",
	"ClassFire",
	"ClassElevator",
	"ClassElevatorShaft",
	"ClassDangerous",
	"ClassMonsterGenerator",
	"ClassReadable",
	"ClassLight",
	"ClassSimple",
	"ClassComplex",
	"ClassImmobile",
	"ClassVisibleEnable",
	"ClassWeapon",
	"ClassArmor",
	"ClassNotStackable",
	"ClassTreasure",
	"ClassFlag",
	"ClassClientPersist",
	"ClassClientPredict",
	"ClassPickup",
}

var _ enum.Enum[Class] = Class(0)

func ParseClass(s string) (Class, error) {
	return enum.Parse[Class]("class", s, ClassNames)
}

func ParseClassSet(s string) (Class, error) {
	return enum.ParseSet[Class]("class", s, ClassNames)
}

type Class uint32

const (
	ClassMissile          = Class(1 << iota) // 0x1
	ClassMonster                             // 0x2
	ClassPlayer                              // 0x4
	ClassObstacle                            // 0x8
	ClassFood                                // 0x10
	ClassExit                                // 0x20
	ClassKey                                 // 0x40
	ClassDoor                                // 0x80
	ClassInfoBook                            // 0x100
	ClassTrigger                             // 0x200
	ClassTransporter                         // 0x400
	ClassHole                                // 0x800
	ClassWand                                // 0x1000
	ClassFire                                // 0x2000
	ClassElevator                            // 0x4000
	ClassElevatorShaft                       // 0x8000
	ClassDangerous                           // 0x10000
	ClassMonsterGenerator                    // 0x20000
	ClassReadable                            // 0x40000
	ClassLight                               // 0x80000
	ClassSimple                              // 0x100000
	ClassComplex                             // 0x200000
	ClassImmobile                            // 0x400000
	ClassVisibleEnable                       // 0x800000
	ClassWeapon                              // 0x1000000
	ClassArmor                               // 0x2000000
	ClassNotStackable                        // 0x4000000
	ClassTreasure                            // 0x8000000
	ClassFlag                                // 0x10000000
	ClassClientPersist                       // 0x20000000
	ClassClientPredict                       // 0x40000000
	ClassPickup                              // 0x80000000
)

const (
	MaskUnits   = ClassPlayer | ClassMonster
	MaskTargets = ClassMonsterGenerator | MaskUnits
)

func (c Class) Has(c2 Class) bool {
	return c&c2 != 0
}

func (c Class) HasAny(c2 Class) bool {
	return c&c2 != 0
}

func (c Class) Split() []Class {
	return enum.SplitBits(c)
}

func (c Class) String() string {
	return enum.StringBits(uint32(c), ClassNames)
}

func (c Class) GoString() string {
	return enum.StringBits(uint32(c), goClassNames)
}

func (c Class) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
