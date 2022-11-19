package spell

type Spell string

const (
	INVALID                     = "SPELL_INVALID"
	ANCHOR                      = "SPELL_ANCHOR"
	ARACHNAPHOBIA               = "SPELL_ARACHNAPHOBIA"
	BLIND                       = "SPELL_BLIND"
	BLINK                       = "SPELL_BLINK"
	BURN                        = "SPELL_BURN"
	CANCEL                      = "SPELL_CANCEL"
	CHAIN_LIGHTNING_BOLT        = "SPELL_CHAIN_LIGHTNING_BOLT"
	CHANNEL_LIFE                = "SPELL_CHANNEL_LIFE"
	CHARM                       = "SPELL_CHARM"
	CLEANSING_FLAME             = "SPELL_CLEANSING_FLAME"
	CLEANSING_MANA_FLAME        = "SPELL_CLEANSING_MANA_FLAME"
	CONFUSE                     = "SPELL_CONFUSE"
	COUNTERSPELL                = "SPELL_COUNTERSPELL"
	CURE_POISON                 = "SPELL_CURE_POISON"
	DEATH                       = "SPELL_DEATH"
	DEATH_RAY                   = "SPELL_DEATH_RAY"
	DETECT_MAGIC                = "SPELL_DETECT_MAGIC"
	DETONATE                    = "SPELL_DETONATE"
	DETONATE_GLYPHS             = "SPELL_DETONATE_GLYPHS"
	DISENCHANT_ALL              = "SPELL_DISENCHANT_ALL"
	TURN_UNDEAD                 = "SPELL_TURN_UNDEAD"
	DRAIN_MANA                  = "SPELL_DRAIN_MANA"
	EARTHQUAKE                  = "SPELL_EARTHQUAKE"
	LIGHTNING                   = "SPELL_LIGHTNING"
	EXPLOSION                   = "SPELL_EXPLOSION"
	FEAR                        = "SPELL_FEAR"
	FIREBALL                    = "SPELL_FIREBALL"
	FIREWALK                    = "SPELL_FIREWALK"
	FIST                        = "SPELL_FIST"
	FORCE_FIELD                 = "SPELL_FORCE_FIELD"
	FORCE_OF_NATURE             = "SPELL_FORCE_OF_NATURE"
	FREEZE                      = "SPELL_FREEZE"
	FUMBLE                      = "SPELL_FUMBLE"
	GLYPH                       = "SPELL_GLYPH"
	GREATER_HEAL                = "SPELL_GREATER_HEAL"
	HASTE                       = "SPELL_HASTE"
	INFRAVISION                 = "SPELL_INFRAVISION"
	INVERSION                   = "SPELL_INVERSION"
	INVISIBILITY                = "SPELL_INVISIBILITY"
	INVULNERABILITY             = "SPELL_INVULNERABILITY"
	LESSER_HEAL                 = "SPELL_LESSER_HEAL"
	LIGHT                       = "SPELL_LIGHT"
	CHAIN_LIGHTNING             = "SPELL_CHAIN_LIGHTNING"
	LOCK                        = "SPELL_LOCK"
	MARK                        = "SPELL_MARK"
	MARK_1                      = "SPELL_MARK_1"
	MARK_2                      = "SPELL_MARK_2"
	MARK_3                      = "SPELL_MARK_3"
	MARK_4                      = "SPELL_MARK_4"
	MAGIC_MISSILE               = "SPELL_MAGIC_MISSILE"
	SHIELD                      = "SPELL_SHIELD"
	METEOR                      = "SPELL_METEOR"
	METEOR_SHOWER               = "SPELL_METEOR_SHOWER"
	MOONGLOW                    = "SPELL_MOONGLOW"
	NULLIFY                     = "SPELL_NULLIFY"
	MANA_BOMB                   = "SPELL_MANA_BOMB"
	PHANTOM                     = "SPELL_PHANTOM"
	PIXIE_SWARM                 = "SPELL_PIXIE_SWARM"
	PLASMA                      = "SPELL_PLASMA"
	POISON                      = "SPELL_POISON"
	PROTECTION_FROM_ELECTRICITY = "SPELL_PROTECTION_FROM_ELECTRICITY"
	PROTECTION_FROM_FIRE        = "SPELL_PROTECTION_FROM_FIRE"
	PROTECTION_FROM_MAGIC       = "SPELL_PROTECTION_FROM_MAGIC"
	PROTECTION_FROM_POISON      = "SPELL_PROTECTION_FROM_POISON"
	PULL                        = "SPELL_PULL"
	PUSH                        = "SPELL_PUSH"
	OVAL_SHIELD                 = "SPELL_OVAL_SHIELD"
	RESTORE_HEALTH              = "SPELL_RESTORE_HEALTH"
	RESTORE_MANA                = "SPELL_RESTORE_MANA"
	RUN                         = "SPELL_RUN"
	SHOCK                       = "SPELL_SHOCK"
	SLOW                        = "SPELL_SLOW"
	SMALL_ZAP                   = "SPELL_SMALL_ZAP"
	STUN                        = "SPELL_STUN"
	SUMMON_BAT                  = "SPELL_SUMMON_BAT"
	SUMMON_BLACK_BEAR           = "SPELL_SUMMON_BLACK_BEAR"
	SUMMON_BEAR                 = "SPELL_SUMMON_BEAR"
	SUMMON_BEHOLDER             = "SPELL_SUMMON_BEHOLDER"
	SUMMON_BOMBER               = "SPELL_SUMMON_BOMBER"
	SUMMON_CARNIVOROUS_PLANT    = "SPELL_SUMMON_CARNIVOROUS_PLANT"
	SUMMON_ALBINO_SPIDER        = "SPELL_SUMMON_ALBINO_SPIDER"
	SUMMON_SMALL_ALBINO_SPIDER  = "SPELL_SUMMON_SMALL_ALBINO_SPIDER"
	SUMMON_EVIL_CHERUB          = "SPELL_SUMMON_EVIL_CHERUB"
	SUMMON_EMBER_DEMON          = "SPELL_SUMMON_EMBER_DEMON"
	SUMMON_GHOST                = "SPELL_SUMMON_GHOST"
	SUMMON_GIANT_LEECH          = "SPELL_SUMMON_GIANT_LEECH"
	SUMMON_IMP                  = "SPELL_SUMMON_IMP"
	SUMMON_MECHANICAL_FLYER     = "SPELL_SUMMON_MECHANICAL_FLYER"
	SUMMON_MECHANICAL_GOLEM     = "SPELL_SUMMON_MECHANICAL_GOLEM"
	SUMMON_MIMIC                = "SPELL_SUMMON_MIMIC"
	SUMMON_OGRE                 = "SPELL_SUMMON_OGRE"
	SUMMON_OGRE_BRUTE           = "SPELL_SUMMON_OGRE_BRUTE"
	SUMMON_OGRE_WARLORD         = "SPELL_SUMMON_OGRE_WARLORD"
	SUMMON_SCORPION             = "SPELL_SUMMON_SCORPION"
	SUMMON_SHADE                = "SPELL_SUMMON_SHADE"
	SUMMON_SKELETON             = "SPELL_SUMMON_SKELETON"
	SUMMON_SKELETON_LORD        = "SPELL_SUMMON_SKELETON_LORD"
	SUMMON_SPIDER               = "SPELL_SUMMON_SPIDER"
	SUMMON_SMALL_SPIDER         = "SPELL_SUMMON_SMALL_SPIDER"
	SUMMON_SPITTING_SPIDER      = "SPELL_SUMMON_SPITTING_SPIDER"
	SUMMON_STONE_GOLEM          = "SPELL_SUMMON_STONE_GOLEM"
	SUMMON_TROLL                = "SPELL_SUMMON_TROLL"
	SUMMON_URCHIN               = "SPELL_SUMMON_URCHIN"
	SUMMON_WASP                 = "SPELL_SUMMON_WASP"
	SUMMON_WILLOWISP            = "SPELL_SUMMON_WILLOWISP"
	SUMMON_WOLF                 = "SPELL_SUMMON_WOLF"
	SUMMON_BLACK_WOLF           = "SPELL_SUMMON_BLACK_WOLF"
	SUMMON_WHITE_WOLF           = "SPELL_SUMMON_WHITE_WOLF"
	SUMMON_ZOMBIE               = "SPELL_SUMMON_ZOMBIE"
	SUMMON_VILE_ZOMBIE          = "SPELL_SUMMON_VILE_ZOMBIE"
	SUMMON_DEMON                = "SPELL_SUMMON_DEMON"
	SUMMON_LICH                 = "SPELL_SUMMON_LICH"
	SUMMON_DRYAD                = "SPELL_SUMMON_DRYAD"
	SUMMON_URCHIN_SHAMAN        = "SPELL_SUMMON_URCHIN_SHAMAN"
	SWAP                        = "SPELL_SWAP"
	TAG                         = "SPELL_TAG"
	TELEPORT_OTHER_TO_MARK_1    = "SPELL_TELEPORT_OTHER_TO_MARK_1"
	TELEPORT_OTHER_TO_MARK_2    = "SPELL_TELEPORT_OTHER_TO_MARK_2"
	TELEPORT_OTHER_TO_MARK_3    = "SPELL_TELEPORT_OTHER_TO_MARK_3"
	TELEPORT_OTHER_TO_MARK_4    = "SPELL_TELEPORT_OTHER_TO_MARK_4"
	TELEPORT_POP                = "SPELL_TELEPORT_POP"
	TELEPORT_TO_MARK_1          = "SPELL_TELEPORT_TO_MARK_1"
	TELEPORT_TO_MARK_2          = "SPELL_TELEPORT_TO_MARK_2"
	TELEPORT_TO_MARK_3          = "SPELL_TELEPORT_TO_MARK_3"
	TELEPORT_TO_MARK_4          = "SPELL_TELEPORT_TO_MARK_4"
	TELEPORT_TO_TARGET          = "SPELL_TELEPORT_TO_TARGET"
	TELEKINESIS                 = "SPELL_TELEKINESIS"
	TOXIC_CLOUD                 = "SPELL_TOXIC_CLOUD"
	TRIGGER_GLYPH               = "SPELL_TRIGGER_GLYPH"
	VAMPIRISM                   = "SPELL_VAMPIRISM"
	VILLAIN                     = "SPELL_VILLAIN"
	WALL                        = "SPELL_WALL"
	WINK                        = "SPELL_WINK"
	SUMMON_CREATURE             = "SPELL_SUMMON_CREATURE"
	MARK_LOCATION               = "SPELL_MARK_LOCATION"
	TELEPORT_TO_MARKER          = "SPELL_TELEPORT_TO_MARKER"
)