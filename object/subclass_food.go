package object

import (
	"github.com/opennox/libs/enum"
)

var FoodClassNames = []string{
	"SIMPLE", "APPLE", "JUG",
	"POTION", "HEALTH_POTION", "MANA_POTION", "CURE_POISON_POTION",
	"MUSHROOM",
	"HASTE_POTION", "INVISIBILITY_POTION", "SHIELD_POTION",
	"FIRE_PROTECT_POTION", "SHOCK_PROTECT_POTION", "POISON_PROTECT_POTION",
	"INVULNERABILITY_POTION", "INFRAVISION_POTION", "VAMPIRISM_POTION",
}

func (c SubClass) AsFood() FoodClass {
	return FoodClass(c)
}

func ParseFoodClass(s string) (FoodClass, error) {
	return enum.Parse[FoodClass]("food class", s, FoodClassNames)
}

func ParseFoodClassSet(s string) (FoodClass, error) {
	return enum.ParseSet[FoodClass]("food class", s, FoodClassNames)
}

var _ enum.Enum[FoodClass] = FoodClass(0)

type FoodClass uint32

const (
	FoodSimple                = FoodClass(1 << iota) // 0x1
	FoodApple                                        // 0x2
	FoodJug                                          // 0x4
	FoodPotion                                       // 0x8
	FoodHealthPotion                                 // 0x10
	FoodManaPotion                                   // 0x20
	FoodCurePoisonPotion                             // 0x40
	FoodMushroom                                     // 0x80
	FoodHastePotion                                  // 0x100
	FoodInvisibilityPotion                           // 0x200
	FoodShieldPotion                                 // 0x400
	FoodFireProtectPotion                            // 0x800
	FoodShockProtectPotion                           // 0x1000
	FoodPoisonProtectPotion                          // 0x2000
	FoodInvulnerabilityPotion                        // 0x4000
	FoodInfravisionPotion                            // 0x8000
	FoodVampirismPotion                              // 0x10000
)

func (c FoodClass) Has(c2 FoodClass) bool {
	return c&c2 != 0
}

func (c FoodClass) HasAny(c2 FoodClass) bool {
	return c&c2 != 0
}

func (c FoodClass) Split() []FoodClass {
	return enum.SplitBits(c)
}

func (c FoodClass) String() string {
	return enum.StringBits(c, FoodClassNames)
}

func (c FoodClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
