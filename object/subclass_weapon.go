package object

import (
	"github.com/opennox/libs/enum"
)

var WeaponClassNames = []string{
	"FLAG",
	"QUIVER",
	"BOW",
	"CROSSBOW",
	"ARROW",
	"BOLT",
	"CHAKRAM",
	"SHURIKEN",
	"SWORD",
	"LONG_SWORD",
	"GREAT_SWORD",
	"MACE",
	"AXE",
	"OGRE_AXE",
	"HAMMER",
	"STAFF",
	"STAFF_SULPHOROUS_FLARE",
	"STAFF_SULPHOROUS_SHOWER",
	"STAFF_LIGHTNING",
	"STAFF_FIREBALL",
	"STAFF_TRIPLE_FIREBALL",
	"STAFF_FORCE_OF_NATURE",
	"STAFF_DEATH_RAY",
	"STAFF_OBLIVION_HALBERD",
	"STAFF_OBLIVION_HEART",
	"STAFF_OBLIVION_WIERDLING",
	"STAFF_OBLIVION_ORB",
}

func (c SubClass) AsWeapon() WeaponClass {
	return WeaponClass(c)
}

func ParseWeaponClass(s string) (WeaponClass, error) {
	return enum.Parse[WeaponClass]("weapon class", s, WeaponClassNames)
}

func ParseWeaponClassSet(s string) (WeaponClass, error) {
	return enum.ParseSet[WeaponClass]("weapon class", s, WeaponClassNames)
}

var _ enum.Enum[WeaponClass] = WeaponClass(0)

type WeaponClass uint32

const (
	WeaponFlag                   = WeaponClass(1 << iota) // 0x1
	WeaponQuiver                                          // 0x2
	WeaponBow                                             // 0x4
	WeaponCrossbow                                        // 0x8
	WeaponArrow                                           // 0x10
	WeaponBolt                                            // 0x20
	WeaponChakram                                         // 0x40
	WeaponShuriken                                        // 0x80
	WeaponSword                                           // 0x100
	WeaponLongSword                                       // 0x200
	WeaponGreatSword                                      // 0x400
	WeaponMace                                            // 0x800
	WeaponAxe                                             // 0x1000
	WeaponOgreAxe                                         // 0x2000
	WeaponHammer                                          // 0x4000
	WeaponStaff                                           // 0x8000
	WeaponStaffSulphorousFlare                            // 0x10000
	WeaponStaffSulphorousShower                           // 0x20000
	WeaponStaffLightning                                  // 0x40000
	WeaponStaffFireball                                   // 0x80000
	WeaponStaffTripleFireball                             // 0x100000
	WeaponStaffForceOfNature                              // 0x200000
	WeaponStaffDeathRay                                   // 0x400000
	WeaponStaffOblivionHalberd                            // 0x800000
	WeaponStaffOblivionHeart                              // 0x1000000
	WeaponStaffOblivionWierdling                          // 0x2000000
	WeaponStaffOblivionOrb                                // 0x4000000
)

func (c WeaponClass) Has(c2 WeaponClass) bool {
	return c&c2 != 0
}

func (c WeaponClass) HasAny(c2 WeaponClass) bool {
	return c&c2 != 0
}

func (c WeaponClass) Split() []WeaponClass {
	return enum.SplitBits(c)
}

func (c WeaponClass) String() string {
	return enum.StringBits(c, WeaponClassNames)
}

func (c WeaponClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}
