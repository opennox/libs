package object

import (
	"github.com/opennox/libs/enum"
)

var MonsterClassNames = []string{
	"SMALL_MONSTER", "MEDIUM_MONSTER", "LARGE_MONSTER", "SHOPKEEPER", "NPC", "FEMALE_NPC",
	"UNDEAD", "MONITOR", "MIGRATE", "IMMUNE_POISON", "IMMUNE_FIRE", "IMMUNE_ELECTRICITY",
	"IMMUNE_FEAR", "BOMBER", "NO_TARGET", "NO_SPELL_TARGET", "HAS_SOUL", "WARCRY_STUN",
	"LOOK_AROUND", "WOUNDED_NPC",
}

func (c SubClass) AsMonster() MonsterClass {
	return MonsterClass(c)
}

func ParseMonsterClass(s string) (MonsterClass, error) {
	return enum.Parse[MonsterClass]("monster class", s, MonsterClassNames)
}

func ParseMonsterClassSet(s string) (MonsterClass, error) {
	return enum.ParseSet[MonsterClass]("monster class", s, MonsterClassNames)
}

var _ enum.Enum[MonsterClass] = MonsterClass(0)

type MonsterClass uint32

const (
	MonsterSmall             = MonsterClass(1 << iota) // 0x1
	MonsterMedium                                      // 0x2
	MonsterLarge                                       // 0x4
	MonsterShopkeeper                                  // 0x8
	MonsterNPC                                         // 0x10
	MonsterFemaleNPC                                   // 0x20
	MonsterUndead                                      // 0x40
	MonsterMonitor                                     // 0x80
	MonsterMigrate                                     // 0x100
	MonsterImmunePoison                                // 0x200
	MonsterImmuneFire                                  // 0x400
	MonsterImmuneElectricity                           // 0x800
	MonsterImmuneFear                                  // 0x1000
	MonsterBomber                                      // 0x2000
	MonsterNoTarget                                    // 0x4000
	MonsterNoSpellTarget                               // 0x8000
	MonsterHasSoul                                     // 0x10000
	MonsterWarcryStun                                  // 0x20000
	MonsterLookAround                                  // 0x40000
	MonsterWoundedNPC                                  // 0x80000
)

func (c MonsterClass) Has(c2 MonsterClass) bool {
	return c&c2 != 0
}

func (c MonsterClass) HasAny(c2 MonsterClass) bool {
	return c&c2 != 0
}

func (c MonsterClass) Split() []MonsterClass {
	return enum.SplitBits(c)
}

func (c MonsterClass) String() string {
	return enum.StringBits(c, MonsterClassNames)
}

func (c MonsterClass) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}

var MonsterStatusNames = []string{
	"DESTROY_WHEN_DEAD", "CHECK",
	"CAN_BLOCK", "CAN_DODGE",
	"unused", "CAN_CAST_SPELLS",
	"HOLD_YOUR_GROUND", "SUMMONED",
	"ALERT", "INJURED",
	"CAN_SEE_FRIENDS", "CAN_HEAL_SELF",
	"CAN_HEAL_OTHERS", "CAN_RUN",
	"RUNNING", "ALWAYS_RUN", "NEVER_RUN",
	"BOT", "MORPHED",
	"ON_FIRE", "STAY_DEAD",
	"FRUSTRATED",
}

func ParseMonsterStatus(s string) (MonsterStatus, error) {
	return enum.Parse[MonsterStatus]("monster status", s, MonsterStatusNames)
}

func ParseMonsterStatusSet(s string) (MonsterStatus, error) {
	return enum.ParseSet[MonsterStatus]("monster status", s, MonsterStatusNames)
}

var _ enum.Enum[MonsterStatus] = MonsterStatus(0)

type MonsterStatus uint32

func (c MonsterStatus) Has(c2 MonsterStatus) bool {
	return c&c2 != 0
}

func (c MonsterStatus) HasAny(c2 MonsterStatus) bool {
	return c&c2 != 0
}

func (c MonsterStatus) Split() []MonsterStatus {
	return enum.SplitBits(c)
}

func (c MonsterStatus) String() string {
	return enum.StringBits(c, MonsterStatusNames)
}

func (c MonsterStatus) MarshalJSON() ([]byte, error) {
	return enum.MarshalJSONArray(c)
}

const (
	MonStatusDestroyWhenDead = MonsterStatus(1 << iota) // 0x1
	MonStatusCheck                                      // 0x2
	MonStatusCanBlock                                   // 0x4
	MonStatusCanDodge                                   // 0x8
	MonStatusUnused5                                    // 0x10
	MonStatusCanCastSpells                              // 0x20
	MonStatusHoldYourGround                             // 0x40
	MonStatusSummoned                                   // 0x80
	MonStatusAlert                                      // 0x100
	MonStatusInjured                                    // 0x200
	MonStatusCanSeeFriends                              // 0x400
	MonStatusCanHealSelf                                // 0x800
	MonStatusCanHealOthers                              // 0x1000
	MonStatusCanRun                                     // 0x2000
	MonStatusRunning                                    // 0x4000
	MonStatusAlwaysRun                                  // 0x8000
	MonStatusNeverRun                                   // 0x10000
	MonStatusBot                                        // 0x20000
	MonStatusMorphed                                    // 0x40000
	MonStatusOnFire                                     // 0x80000
	MonStatusStayDead                                   // 0x100000
	MonStatusFrustrated                                 // 0x200000
)
