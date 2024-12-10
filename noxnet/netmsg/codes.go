package netmsg

//go:generate stringer -type=Op

type Op byte

func (op Op) Len() int {
	n, ok := opLen[op]
	if !ok || n < 0 {
		return -1
	}
	return n
}

const (
	MSG_SERVER_CONNECT                 = Op(0)   // 0x00
	MSG_SERVER_ACCEPT                  = Op(1)   // 0x01
	MSG_CODE2                          = Op(2)   // 0x02
	MSG_CODE3                          = Op(3)   // 0x03
	MSG_CODE4                          = Op(4)   // 0x04
	MSG_CODE5                          = Op(5)   // 0x05
	MSG_CLIENT_PING                    = Op(6)   // 0x06
	MSG_CODE7                          = Op(7)   // 0x07
	MSG_CLIENT_PONG                    = Op(8)   // 0x08
	MSG_CODE9                          = Op(9)   // 0x09
	MSG_CLIENT_CLOSE                   = Op(10)  // 0x0A
	MSG_SERVER_CLOSE                   = Op(11)  // 0x0B
	MSG_SERVER_DISCOVER                = Op(12)  // 0x0C
	MSG_SERVER_INFO                    = Op(13)  // 0x0D
	MSG_SERVER_TRY_JOIN                = Op(14)  // 0x0E
	MSG_PASSWORD_REQUIRED              = Op(15)  // 0x0F
	MSG_SERVER_PING                    = Op(16)  // 0x10
	MSG_SERVER_PASSWORD                = Op(17)  // 0x11
	MSG_SERVER_PONG                    = Op(18)  // 0x12
	MSG_SERVER_ERROR                   = Op(19)  // 0x13
	MSG_SERVER_JOIN_OK                 = Op(20)  // 0x14
	MSG_SERVER_JOIN_FAIL               = Op(21)  // 0x15
	MSG_CODE22                         = Op(22)  // 0x16
	MSG_CODE23                         = Op(23)  // 0x17
	MSG_CODE24                         = Op(24)  // 0x18
	MSG_CODE25                         = Op(25)  // 0x19
	MSG_CODE26                         = Op(26)  // 0x1A
	MSG_CODE27                         = Op(27)  // 0x1B
	MSG_CODE28                         = Op(28)  // 0x1C
	MSG_CODE29                         = Op(29)  // 0x1D
	MSG_CODE30                         = Op(30)  // 0x1E
	MSG_ACCEPTED                       = Op(31)  // 0x1F
	MSG_CLIENT_ACCEPT                  = Op(32)  // 0x20
	MSG_SERVER_CLOSE_ACK               = Op(33)  // 0x21
	MSG_CLIENT_CLOSE_ACK               = Op(34)  // 0x22
	MSG_SPEED                          = Op(35)  // 0x23
	MSG_PING                           = Op(36)  // 0x24
	MSG_CODE37                         = Op(37)  // 0x25
	MSG_CODE38                         = Op(38)  // 0x26
	MSG_TIMESTAMP                      = Op(39)  // 0x27
	MSG_FULL_TIMESTAMP                 = Op(40)  // 0x28
	MSG_NEED_TIMESTAMP                 = Op(41)  // 0x29
	MSG_SIMULATED_TIMESTAMP            = Op(42)  // 0x2A
	MSG_USE_MAP                        = Op(43)  // 0x2B
	MSG_JOIN_DATA                      = Op(44)  // 0x2C
	MSG_NEW_PLAYER                     = Op(45)  // 0x2D
	MSG_PLAYER_QUIT                    = Op(46)  // 0x2E
	MSG_SIMPLE_OBJ                     = Op(47)  // 0x2F
	MSG_COMPLEX_OBJ                    = Op(48)  // 0x30
	MSG_DESTROY_OBJECT                 = Op(49)  // 0x31
	MSG_OBJECT_OUT_OF_SIGHT            = Op(50)  // 0x32
	MSG_OBJECT_IN_SHADOWS              = Op(51)  // 0x33
	MSG_OBJECT_FRIEND_ADD              = Op(52)  // 0x34
	MSG_OBJECT_FRIEND_REMOVE           = Op(53)  // 0x35
	MSG_RESET_FRIENDS                  = Op(54)  // 0x36
	MSG_ENABLE_OBJECT                  = Op(55)  // 0x37
	MSG_DISABLE_OBJECT                 = Op(56)  // 0x38
	MSG_DRAW_FRAME                     = Op(57)  // 0x39
	MSG_DESTROY_WALL                   = Op(58)  // 0x3A
	MSG_OPEN_WALL                      = Op(59)  // 0x3B
	MSG_CLOSE_WALL                     = Op(60)  // 0x3C
	MSG_CHANGE_OR_ADD_WALL_MAGIC       = Op(61)  // 0x3D
	MSG_REMOVE_WALL_MAGIC              = Op(62)  // 0x3E
	MSG_PLAYER_INPUT                   = Op(63)  // 0x3F
	MSG_PLAYER_SET_WAYPOINT            = Op(64)  // 0x40
	MSG_REPORT_HEALTH                  = Op(65)  // 0x41
	MSG_REPORT_HEALTH_DELTA            = Op(66)  // 0x42
	MSG_REPORT_PLAYER_HEALTH           = Op(67)  // 0x43
	MSG_REPORT_ITEM_HEALTH             = Op(68)  // 0x44
	MSG_REPORT_MANA                    = Op(69)  // 0x45
	MSG_REPORT_POISON                  = Op(70)  // 0x46
	MSG_REPORT_STAMINA                 = Op(71)  // 0x47
	MSG_REPORT_STATS                   = Op(72)  // 0x48
	MSG_REPORT_ARMOR_VALUE             = Op(73)  // 0x49
	MSG_REPORT_GOLD                    = Op(74)  // 0x4A
	MSG_REPORT_PICKUP                  = Op(75)  // 0x4B
	MSG_REPORT_MODIFIABLE_PICKUP       = Op(76)  // 0x4C
	MSG_REPORT_DROP                    = Op(77)  // 0x4D
	MSG_REPORT_LESSON                  = Op(78)  // 0x4E
	MSG_REPORT_MUNDANE_ARMOR_EQUIP     = Op(79)  // 0x4F
	MSG_REPORT_MUNDANE_WEAPON_EQUIP    = Op(80)  // 0x50
	MSG_REPORT_MODIFIABLE_WEAPON_EQUIP = Op(81)  // 0x51
	MSG_REPORT_MODIFIABLE_ARMOR_EQUIP  = Op(82)  // 0x52
	MSG_REPORT_ARMOR_DEQUIP            = Op(83)  // 0x53
	MSG_REPORT_WEAPON_DEQUIP           = Op(84)  // 0x54
	MSG_REPORT_TREASURE_COUNT          = Op(85)  // 0x55
	MSG_REPORT_FLAG_BALL_WINNER        = Op(86)  // 0x56
	MSG_REPORT_FLAG_WINNER             = Op(87)  // 0x57
	MSG_REPORT_DEATHMATCH_WINNER       = Op(88)  // 0x58
	MSG_REPORT_DEATHMATCH_TEAM_WINNER  = Op(89)  // 0x59
	MSG_REPORT_ENCHANTMENT             = Op(90)  // 0x5A
	MSG_REPORT_ITEM_ENCHANTMENT        = Op(91)  // 0x5B
	MSG_REPORT_LIGHT_COLOR             = Op(92)  // 0x5C
	MSG_REPORT_LIGHT_INTENSITY         = Op(93)  // 0x5D
	MSG_REPORT_Z_PLUS                  = Op(94)  // 0x5E
	MSG_REPORT_Z_MINUS                 = Op(95)  // 0x5F
	MSG_REPORT_EQUIP                   = Op(96)  // 0x60
	MSG_REPORT_DEQUIP                  = Op(97)  // 0x61
	MSG_REPORT_ACQUIRE_SPELL           = Op(98)  // 0x62
	MSG_REPORT_TARGET                  = Op(99)  // 0x63
	MSG_REPORT_CHARGES                 = Op(100) // 0x64
	MSG_REPORT_X_STATUS                = Op(101) // 0x65
	MSG_REPORT_PLAYER_STATUS           = Op(102) // 0x66
	MSG_REPORT_MODIFIER                = Op(103) // 0x67
	MSG_REPORT_STAT_MODIFIER           = Op(104) // 0x68
	MSG_REPORT_NPC                     = Op(105) // 0x69
	MSG_REPORT_CLIENT_STATUS           = Op(106) // 0x6A
	MSG_REPORT_ANIMATION_FRAME         = Op(107) // 0x6B
	MSG_REPORT_ACQUIRE_CREATURE        = Op(108) // 0x6C
	MSG_REPORT_LOSE_CREATURE           = Op(109) // 0x6D
	MSG_REPORT_EXPERIENCE              = Op(110) // 0x6E
	MSG_REPORT_SPELL_AWARD             = Op(111) // 0x6F
	MSG_REPORT_SPELL_START             = Op(112) // 0x70
	MSG_REPORT_INVENTORY_LOADED        = Op(113) // 0x71
	MSG_TRY_DROP                       = Op(114) // 0x72
	MSG_TRY_GET                        = Op(115) // 0x73
	MSG_TRY_USE                        = Op(116) // 0x74
	MSG_TRY_EQUIP                      = Op(117) // 0x75
	MSG_TRY_DEQUIP                     = Op(118) // 0x76
	MSG_TRY_TARGET                     = Op(119) // 0x77
	MSG_TRY_CREATURE_COMMAND           = Op(120) // 0x78
	MSG_TRY_SPELL                      = Op(121) // 0x79
	MSG_TRY_ABILITY                    = Op(122) // 0x7A
	MSG_TRY_COLLIDE                    = Op(123) // 0x7B
	MSG_FX_PARTICLEFX                  = Op(124) // 0x7C
	MSG_FX_PLASMA                      = Op(125) // 0x7D
	MSG_FX_SUMMON                      = Op(126) // 0x7E
	MSG_FX_SUMMON_CANCEL               = Op(127) // 0x7F
	MSG_FX_SHIELD                      = Op(128) // 0x80
	MSG_FX_BLUE_SPARKS                 = Op(129) // 0x81
	MSG_FX_YELLOW_SPARKS               = Op(130) // 0x82
	MSG_FX_CYAN_SPARKS                 = Op(131) // 0x83
	MSG_FX_VIOLET_SPARKS               = Op(132) // 0x84
	MSG_FX_EXPLOSION                   = Op(133) // 0x85
	MSG_FX_LESSER_EXPLOSION            = Op(134) // 0x86
	MSG_FX_COUNTERSPELL_EXPLOSION      = Op(135) // 0x87
	MSG_FX_THIN_EXPLOSION              = Op(136) // 0x88
	MSG_FX_TELEPORT                    = Op(137) // 0x89
	MSG_FX_SMOKE_BLAST                 = Op(138) // 0x8A
	MSG_FX_DAMAGE_POOF                 = Op(139) // 0x8B
	MSG_FX_LIGHTNING                   = Op(140) // 0x8C
	MSG_FX_ENERGY_BOLT                 = Op(141) // 0x8D
	MSG_FX_CHAIN_LIGHTNING_BOLT        = Op(142) // 0x8E
	MSG_FX_DRAIN_MANA                  = Op(143) // 0x8F
	MSG_FX_CHARM                       = Op(144) // 0x90
	MSG_FX_GREATER_HEAL                = Op(145) // 0x91
	MSG_FX_MAGIC                       = Op(146) // 0x92
	MSG_FX_SPARK_EXPLOSION             = Op(147) // 0x93
	MSG_FX_DEATH_RAY                   = Op(148) // 0x94
	MSG_FX_SENTRY_RAY                  = Op(149) // 0x95
	MSG_FX_RICOCHET                    = Op(150) // 0x96
	MSG_FX_JIGGLE                      = Op(151) // 0x97
	MSG_FX_GREEN_BOLT                  = Op(152) // 0x98
	MSG_FX_GREEN_EXPLOSION             = Op(153) // 0x99
	MSG_FX_WHITE_FLASH                 = Op(154) // 0x9A
	MSG_FX_GENERATING_MAP              = Op(155) // 0x9B
	MSG_FX_ASSEMBLING_MAP              = Op(156) // 0x9C
	MSG_FX_POPULATING_MAP              = Op(157) // 0x9D
	MSG_FX_DURATION_SPELL              = Op(158) // 0x9E
	MSG_FX_DELTAZ_SPELL_START          = Op(159) // 0x9F
	MSG_FX_TURN_UNDEAD                 = Op(160) // 0xA0
	MSG_FX_ARROW_TRAP                  = Op(161) // 0xA1
	MSG_FX_VAMPIRISM                   = Op(162) // 0xA2
	MSG_FX_MANA_BOMB_CANCEL            = Op(163) // 0xA3
	MSG_UPDATE_STREAM                  = Op(164) // 0xA4
	MSG_NEW_ALIAS                      = Op(165) // 0xA5
	MSG_AUDIO_EVENT                    = Op(166) // 0xA6
	MSG_AUDIO_PLAYER_EVENT             = Op(167) // 0xA7
	MSG_TEXT_MESSAGE                   = Op(168) // 0xA8
	MSG_INFORM                         = Op(169) // 0xA9
	MSG_IMPORTANT                      = Op(170) // 0xAA
	MSG_IMPORTANT_ACK                  = Op(171) // 0xAB
	MSG_MOUSE                          = Op(172) // 0xAC
	MSG_INCOMING_CLIENT                = Op(173) // 0xAD
	MSG_OUTGOING_CLIENT                = Op(174) // 0xAE
	MSG_GAME_SETTINGS                  = Op(175) // 0xAF
	MSG_GAME_SETTINGS_2                = Op(176) // 0xB0
	MSG_UPDATE_GUI_GAME_SETTINGS       = Op(177) // 0xB1
	MSG_DOOR_ANGLE                     = Op(178) // 0xB2
	MSG_OBELISK_CHARGE                 = Op(179) // 0xB3
	MSG_PENTAGRAM_ACTIVATE             = Op(180) // 0xB4
	MSG_CLIENT_PREDICT_LINEAR          = Op(181) // 0xB5
	MSG_REQUEST_MAP                    = Op(182) // 0xB6
	MSG_CANCEL_MAP                     = Op(183) // 0xB7
	MSG_MAP_SEND_START                 = Op(184) // 0xB8
	MSG_MAP_SEND_PACKET                = Op(185) // 0xB9
	MSG_MAP_SEND_ABORT                 = Op(186) // 0xBA
	MSG_SERVER_CMD                     = Op(187) // 0xBB
	MSG_SYSOP_PW                       = Op(188) // 0xBC
	MSG_SYSOP_RESULT                   = Op(189) // 0xBD
	MSG_KEEP_ALIVE                     = Op(190) // 0xBE
	MSG_RECEIVED_MAP                   = Op(191) // 0xBF
	MSG_CLIENT_READY                   = Op(192) // 0xC0
	MSG_REQUEST_SAVE_PLAYER            = Op(193) // 0xC1
	MSG_XFER_MSG                       = Op(194) // 0xC2
	MSG_PLAYER_OBJ                     = Op(195) // 0xC3
	MSG_TEAM_MSG                       = Op(196) // 0xC4
	MSG_KICK_NOTIFICATION              = Op(197) // 0xC5
	MSG_TIMEOUT_NOTIFICATION           = Op(198) // 0xC6
	MSG_SERVER_QUIT                    = Op(199) // 0xC7
	MSG_SERVER_QUIT_ACK                = Op(200) // 0xC8
	MSG_TRADE                          = Op(201) // 0xC9
	MSG_CHAT_KILL                      = Op(202) // 0xCA
	MSG_MESSAGES_KILL                  = Op(203) // 0xCB
	MSG_SEQ_IMPORTANT                  = Op(204) // 0xCC
	MSG_REPORT_ABILITY_AWARD           = Op(205) // 0xCD
	MSG_REPORT_ABILITY_STATE           = Op(206) // 0xCE
	MSG_REPORT_ACTIVE_ABILITIES        = Op(207) // 0xCF
	MSG_DIALOG                         = Op(208) // 0xD0
	MSG_REPORT_GUIDE_AWARD             = Op(209) // 0xD1
	MSG_INTERESTING_ID                 = Op(210) // 0xD2
	MSG_TIMER_STATUS                   = Op(211) // 0xD3
	MSG_REQUEST_TIMER_STATUS           = Op(212) // 0xD4
	MSG_JOURNAL_MSG                    = Op(213) // 0xD5
	MSG_CHAPTER_END                    = Op(214) // 0xD6
	MSG_REPORT_ALL_LATENCY             = Op(215) // 0xD7
	MSG_REPORT_FLAG_STATUS             = Op(216) // 0xD8
	MSG_REPORT_BALL_STATUS             = Op(217) // 0xD9
	MSG_REPORT_OBJECT_POISON           = Op(218) // 0xDA
	MSG_REPORT_MONITOR_CREATURE        = Op(219) // 0xDB
	MSG_REPORT_UNMONITOR_CREATURE      = Op(220) // 0xDC
	MSG_REPORT_TOTAL_HEALTH            = Op(221) // 0xDD
	MSG_REPORT_TOTAL_MANA              = Op(222) // 0xDE
	MSG_REPORT_SPELL_STAT              = Op(223) // 0xDF
	MSG_REPORT_SECONDARY_WEAPON        = Op(224) // 0xE0
	MSG_REPORT_LAST_QUIVER             = Op(225) // 0xE1
	MSG_INFO_BOOK_DATA                 = Op(226) // 0xE2
	MSG_SOCIAL                         = Op(227) // 0xE3
	MSG_FADE_BEGIN                     = Op(228) // 0xE4
	MSG_MUSIC_EVENT                    = Op(229) // 0xE5
	MSG_MUSIC_PUSH_EVENT               = Op(230) // 0xE6
	MSG_MUSIC_POP_EVENT                = Op(231) // 0xE7
	MSG_PLAYER_DIED                    = Op(232) // 0xE8
	MSG_PLAYER_RESPAWN                 = Op(233) // 0xE9
	MSG_FORGET_DRAWABLES               = Op(234) // 0xEA
	MSG_RESET_ABILITIES                = Op(235) // 0xEB
	MSG_RATE_CHANGE                    = Op(236) // 0xEC
	MSG_REPORT_CREATURE_CMD            = Op(237) // 0xED
	MSG_VOTE                           = Op(238) // 0xEE
	MSG_STAT_MULTIPLIERS               = Op(239) // 0xEF
	MSG_GAUNTLET                       = Op(240) // 0xF0
	MSG_INVENTORY_FAIL                 = Op(241) // 0xF1
)

var opLen = map[Op]int{
	MSG_INCOMING_CLIENT:                0,
	MSG_CLIENT_READY:                   0,
	MSG_PLAYER_QUIT:                    2,
	MSG_SIMPLE_OBJ:                     8,
	MSG_DESTROY_OBJECT:                 2,
	MSG_OBJECT_OUT_OF_SIGHT:            2,
	MSG_OBJECT_IN_SHADOWS:              2,
	MSG_OBJECT_FRIEND_ADD:              2,
	MSG_OBJECT_FRIEND_REMOVE:           2,
	MSG_RESET_FRIENDS:                  2,
	MSG_ENABLE_OBJECT:                  2,
	MSG_DISABLE_OBJECT:                 2,
	MSG_DRAW_FRAME:                     3,
	MSG_DESTROY_WALL:                   2,
	MSG_OPEN_WALL:                      2,
	MSG_CLOSE_WALL:                     2,
	MSG_CHANGE_OR_ADD_WALL_MAGIC:       5,
	MSG_REMOVE_WALL_MAGIC:              2,
	MSG_REPORT_HEALTH_DELTA:            4,
	MSG_REPORT_PLAYER_HEALTH:           2,
	MSG_REPORT_ITEM_HEALTH:             6,
	MSG_REPORT_MANA:                    4,
	MSG_REPORT_STAMINA:                 1,
	MSG_REPORT_STATS:                   13,
	MSG_REPORT_ARMOR_VALUE:             4,
	MSG_REPORT_GOLD:                    4,
	MSG_REPORT_PICKUP:                  4,
	MSG_REPORT_MODIFIABLE_PICKUP:       8,
	MSG_REPORT_DROP:                    4,
	MSG_REPORT_LESSON:                  10,
	MSG_REPORT_MUNDANE_ARMOR_EQUIP:     6,
	MSG_REPORT_MUNDANE_WEAPON_EQUIP:    6,
	MSG_REPORT_MODIFIABLE_WEAPON_EQUIP: 10,
	MSG_REPORT_MODIFIABLE_ARMOR_EQUIP:  10,
	MSG_REPORT_ARMOR_DEQUIP:            6,
	MSG_REPORT_WEAPON_DEQUIP:           6,
	MSG_REPORT_FLAG_BALL_WINNER:        7,
	MSG_REPORT_FLAG_WINNER:             7,
	MSG_REPORT_DEATHMATCH_WINNER:       7,
	MSG_REPORT_DEATHMATCH_TEAM_WINNER:  7,
	MSG_REPORT_ENCHANTMENT:             6,
	MSG_REPORT_ITEM_ENCHANTMENT:        1,
	MSG_REPORT_Z_PLUS:                  3,
	MSG_REPORT_Z_MINUS:                 3,
	MSG_REPORT_EQUIP:                   2,
	MSG_REPORT_DEQUIP:                  2,
	MSG_REPORT_CHARGES:                 4,
	MSG_REPORT_X_STATUS:                6,
	MSG_REPORT_MODIFIER:                6,
	MSG_REPORT_STAT_MODIFIER:           7,
	MSG_REPORT_CLIENT_STATUS:           6,
	MSG_REPORT_ANIMATION_FRAME:         6,
	MSG_FX_SUMMON:                      11,
	MSG_FX_SUMMON_CANCEL:               2,
	MSG_FX_SHIELD:                      3,
	MSG_FX_BLUE_SPARKS:                 4,
	MSG_FX_YELLOW_SPARKS:               4,
	MSG_FX_CYAN_SPARKS:                 4,
	MSG_FX_VIOLET_SPARKS:               4,
	MSG_FX_EXPLOSION:                   4,
	MSG_FX_LESSER_EXPLOSION:            4,
	MSG_FX_COUNTERSPELL_EXPLOSION:      4,
	MSG_FX_THIN_EXPLOSION:              4,
	MSG_FX_TELEPORT:                    4,
	MSG_FX_SMOKE_BLAST:                 4,
	MSG_FX_DAMAGE_POOF:                 4,
	MSG_FX_LIGHTNING:                   8,
	MSG_FX_ENERGY_BOLT:                 8,
	MSG_FX_CHAIN_LIGHTNING_BOLT:        8,
	MSG_FX_DRAIN_MANA:                  8,
	MSG_FX_CHARM:                       8,
	MSG_FX_GREATER_HEAL:                8,
	MSG_FX_SPARK_EXPLOSION:             5,
	MSG_FX_DEATH_RAY:                   8,
	MSG_FX_SENTRY_RAY:                  8,
	MSG_FX_RICOCHET:                    4,
	MSG_FX_JIGGLE:                      1,
	MSG_FX_WHITE_FLASH:                 4,
	MSG_FX_DURATION_SPELL:              6,
	MSG_FX_DELTAZ_SPELL_START:          5,
	MSG_FX_VAMPIRISM:                   10,
	MSG_FX_MANA_BOMB_CANCEL:            4,
	MSG_AUDIO_EVENT:                    3,
	MSG_AUDIO_PLAYER_EVENT:             3,
	MSG_OUTGOING_CLIENT:                2,
	MSG_GAME_SETTINGS:                  19,
	MSG_GAME_SETTINGS_2:                48,
	MSG_UPDATE_GUI_GAME_SETTINGS:       59,
	MSG_DOOR_ANGLE:                     3,
	MSG_OBELISK_CHARGE:                 3,
	MSG_PENTAGRAM_ACTIVATE:             3,
	MSG_CLIENT_PREDICT_LINEAR:          13,
	MSG_PLAYER_OBJ:                     11,
	MSG_SERVER_QUIT:                    0,
	MSG_REPORT_ABILITY_AWARD:           2,
	MSG_REPORT_ABILITY_STATE:           2,
	MSG_REPORT_ACTIVE_ABILITIES:        2,
	MSG_INTERESTING_ID:                 6,
	MSG_TIMER_STATUS:                   12,
	MSG_REPORT_ALL_LATENCY:             4,
	MSG_REPORT_FLAG_STATUS:             5,
	MSG_REPORT_BALL_STATUS:             3,
	MSG_REPORT_TOTAL_HEALTH:            6,
	MSG_REPORT_TOTAL_MANA:              6,
	MSG_REPORT_SPELL_STAT:              5,
	MSG_FADE_BEGIN:                     2,
	MSG_PLAYER_DIED:                    2,
	MSG_PLAYER_RESPAWN:                 8,
	MSG_FORGET_DRAWABLES:               4,
	MSG_RESET_ABILITIES:                1,
	MSG_RATE_CHANGE:                    1,
	MSG_STAT_MULTIPLIERS:               16,
	MSG_IMPORTANT:                      -1,
	MSG_UPDATE_STREAM:                  -1, // dynamic
	MSG_TEAM_MSG:                       -1, // dynamic
	MSG_SEQ_IMPORTANT:                  -1, // dynamic
	MSG_GAUNTLET:                       -1, // dynamic
}
