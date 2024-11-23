// Code generated by "stringer -type=Op"; DO NOT EDIT.

package noxnet

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MSG_SERVER_CONNECT-0]
	_ = x[MSG_SERVER_ACCEPT-1]
	_ = x[MSG_CODE2-2]
	_ = x[MSG_CODE3-3]
	_ = x[MSG_CODE4-4]
	_ = x[MSG_CODE5-5]
	_ = x[MSG_CLIENT_PING-6]
	_ = x[MSG_CODE7-7]
	_ = x[MSG_CLIENT_PONG-8]
	_ = x[MSG_CODE9-9]
	_ = x[MSG_CLIENT_CLOSE-10]
	_ = x[MSG_SERVER_CLOSE-11]
	_ = x[MSG_SERVER_DISCOVER-12]
	_ = x[MSG_SERVER_INFO-13]
	_ = x[MSG_SERVER_TRY_JOIN-14]
	_ = x[MSG_PASSWORD_REQUIRED-15]
	_ = x[MSG_SERVER_PING-16]
	_ = x[MSG_SERVER_PASSWORD-17]
	_ = x[MSG_SERVER_PONG-18]
	_ = x[MSG_SERVER_ERROR-19]
	_ = x[MSG_SERVER_JOIN_OK-20]
	_ = x[MSG_SERVER_JOIN_FAIL-21]
	_ = x[MSG_CODE22-22]
	_ = x[MSG_CODE23-23]
	_ = x[MSG_CODE24-24]
	_ = x[MSG_CODE25-25]
	_ = x[MSG_CODE26-26]
	_ = x[MSG_CODE27-27]
	_ = x[MSG_CODE28-28]
	_ = x[MSG_CODE29-29]
	_ = x[MSG_CODE30-30]
	_ = x[MSG_ACCEPTED-31]
	_ = x[MSG_CLIENT_ACCEPT-32]
	_ = x[MSG_SERVER_CLOSE_ACK-33]
	_ = x[MSG_CLIENT_CLOSE_ACK-34]
	_ = x[MSG_SPEED-35]
	_ = x[MSG_PING-36]
	_ = x[MSG_CODE37-37]
	_ = x[MSG_CODE38-38]
	_ = x[MSG_TIMESTAMP-39]
	_ = x[MSG_FULL_TIMESTAMP-40]
	_ = x[MSG_NEED_TIMESTAMP-41]
	_ = x[MSG_SIMULATED_TIMESTAMP-42]
	_ = x[MSG_USE_MAP-43]
	_ = x[MSG_JOIN_DATA-44]
	_ = x[MSG_NEW_PLAYER-45]
	_ = x[MSG_PLAYER_QUIT-46]
	_ = x[MSG_SIMPLE_OBJ-47]
	_ = x[MSG_COMPLEX_OBJ-48]
	_ = x[MSG_DESTROY_OBJECT-49]
	_ = x[MSG_OBJECT_OUT_OF_SIGHT-50]
	_ = x[MSG_OBJECT_IN_SHADOWS-51]
	_ = x[MSG_OBJECT_FRIEND_ADD-52]
	_ = x[MSG_OBJECT_FRIEND_REMOVE-53]
	_ = x[MSG_RESET_FRIENDS-54]
	_ = x[MSG_ENABLE_OBJECT-55]
	_ = x[MSG_DISABLE_OBJECT-56]
	_ = x[MSG_DRAW_FRAME-57]
	_ = x[MSG_DESTROY_WALL-58]
	_ = x[MSG_OPEN_WALL-59]
	_ = x[MSG_CLOSE_WALL-60]
	_ = x[MSG_CHANGE_OR_ADD_WALL_MAGIC-61]
	_ = x[MSG_REMOVE_WALL_MAGIC-62]
	_ = x[MSG_PLAYER_INPUT-63]
	_ = x[MSG_PLAYER_SET_WAYPOINT-64]
	_ = x[MSG_REPORT_HEALTH-65]
	_ = x[MSG_REPORT_HEALTH_DELTA-66]
	_ = x[MSG_REPORT_PLAYER_HEALTH-67]
	_ = x[MSG_REPORT_ITEM_HEALTH-68]
	_ = x[MSG_REPORT_MANA-69]
	_ = x[MSG_REPORT_POISON-70]
	_ = x[MSG_REPORT_STAMINA-71]
	_ = x[MSG_REPORT_STATS-72]
	_ = x[MSG_REPORT_ARMOR_VALUE-73]
	_ = x[MSG_REPORT_GOLD-74]
	_ = x[MSG_REPORT_PICKUP-75]
	_ = x[MSG_REPORT_MODIFIABLE_PICKUP-76]
	_ = x[MSG_REPORT_DROP-77]
	_ = x[MSG_REPORT_LESSON-78]
	_ = x[MSG_REPORT_MUNDANE_ARMOR_EQUIP-79]
	_ = x[MSG_REPORT_MUNDANE_WEAPON_EQUIP-80]
	_ = x[MSG_REPORT_MODIFIABLE_WEAPON_EQUIP-81]
	_ = x[MSG_REPORT_MODIFIABLE_ARMOR_EQUIP-82]
	_ = x[MSG_REPORT_ARMOR_DEQUIP-83]
	_ = x[MSG_REPORT_WEAPON_DEQUIP-84]
	_ = x[MSG_REPORT_TREASURE_COUNT-85]
	_ = x[MSG_REPORT_FLAG_BALL_WINNER-86]
	_ = x[MSG_REPORT_FLAG_WINNER-87]
	_ = x[MSG_REPORT_DEATHMATCH_WINNER-88]
	_ = x[MSG_REPORT_DEATHMATCH_TEAM_WINNER-89]
	_ = x[MSG_REPORT_ENCHANTMENT-90]
	_ = x[MSG_REPORT_ITEM_ENCHANTMENT-91]
	_ = x[MSG_REPORT_LIGHT_COLOR-92]
	_ = x[MSG_REPORT_LIGHT_INTENSITY-93]
	_ = x[MSG_REPORT_Z_PLUS-94]
	_ = x[MSG_REPORT_Z_MINUS-95]
	_ = x[MSG_REPORT_EQUIP-96]
	_ = x[MSG_REPORT_DEQUIP-97]
	_ = x[MSG_REPORT_ACQUIRE_SPELL-98]
	_ = x[MSG_REPORT_TARGET-99]
	_ = x[MSG_REPORT_CHARGES-100]
	_ = x[MSG_REPORT_X_STATUS-101]
	_ = x[MSG_REPORT_PLAYER_STATUS-102]
	_ = x[MSG_REPORT_MODIFIER-103]
	_ = x[MSG_REPORT_STAT_MODIFIER-104]
	_ = x[MSG_REPORT_NPC-105]
	_ = x[MSG_REPORT_CLIENT_STATUS-106]
	_ = x[MSG_REPORT_ANIMATION_FRAME-107]
	_ = x[MSG_REPORT_ACQUIRE_CREATURE-108]
	_ = x[MSG_REPORT_LOSE_CREATURE-109]
	_ = x[MSG_REPORT_EXPERIENCE-110]
	_ = x[MSG_REPORT_SPELL_AWARD-111]
	_ = x[MSG_REPORT_SPELL_START-112]
	_ = x[MSG_REPORT_INVENTORY_LOADED-113]
	_ = x[MSG_TRY_DROP-114]
	_ = x[MSG_TRY_GET-115]
	_ = x[MSG_TRY_USE-116]
	_ = x[MSG_TRY_EQUIP-117]
	_ = x[MSG_TRY_DEQUIP-118]
	_ = x[MSG_TRY_TARGET-119]
	_ = x[MSG_TRY_CREATURE_COMMAND-120]
	_ = x[MSG_TRY_SPELL-121]
	_ = x[MSG_TRY_ABILITY-122]
	_ = x[MSG_TRY_COLLIDE-123]
	_ = x[MSG_FX_PARTICLEFX-124]
	_ = x[MSG_FX_PLASMA-125]
	_ = x[MSG_FX_SUMMON-126]
	_ = x[MSG_FX_SUMMON_CANCEL-127]
	_ = x[MSG_FX_SHIELD-128]
	_ = x[MSG_FX_BLUE_SPARKS-129]
	_ = x[MSG_FX_YELLOW_SPARKS-130]
	_ = x[MSG_FX_CYAN_SPARKS-131]
	_ = x[MSG_FX_VIOLET_SPARKS-132]
	_ = x[MSG_FX_EXPLOSION-133]
	_ = x[MSG_FX_LESSER_EXPLOSION-134]
	_ = x[MSG_FX_COUNTERSPELL_EXPLOSION-135]
	_ = x[MSG_FX_THIN_EXPLOSION-136]
	_ = x[MSG_FX_TELEPORT-137]
	_ = x[MSG_FX_SMOKE_BLAST-138]
	_ = x[MSG_FX_DAMAGE_POOF-139]
	_ = x[MSG_FX_LIGHTNING-140]
	_ = x[MSG_FX_ENERGY_BOLT-141]
	_ = x[MSG_FX_CHAIN_LIGHTNING_BOLT-142]
	_ = x[MSG_FX_DRAIN_MANA-143]
	_ = x[MSG_FX_CHARM-144]
	_ = x[MSG_FX_GREATER_HEAL-145]
	_ = x[MSG_FX_MAGIC-146]
	_ = x[MSG_FX_SPARK_EXPLOSION-147]
	_ = x[MSG_FX_DEATH_RAY-148]
	_ = x[MSG_FX_SENTRY_RAY-149]
	_ = x[MSG_FX_RICOCHET-150]
	_ = x[MSG_FX_JIGGLE-151]
	_ = x[MSG_FX_GREEN_BOLT-152]
	_ = x[MSG_FX_GREEN_EXPLOSION-153]
	_ = x[MSG_FX_WHITE_FLASH-154]
	_ = x[MSG_FX_GENERATING_MAP-155]
	_ = x[MSG_FX_ASSEMBLING_MAP-156]
	_ = x[MSG_FX_POPULATING_MAP-157]
	_ = x[MSG_FX_DURATION_SPELL-158]
	_ = x[MSG_FX_DELTAZ_SPELL_START-159]
	_ = x[MSG_FX_TURN_UNDEAD-160]
	_ = x[MSG_FX_ARROW_TRAP-161]
	_ = x[MSG_FX_VAMPIRISM-162]
	_ = x[MSG_FX_MANA_BOMB_CANCEL-163]
	_ = x[MSG_UPDATE_STREAM-164]
	_ = x[MSG_NEW_ALIAS-165]
	_ = x[MSG_AUDIO_EVENT-166]
	_ = x[MSG_AUDIO_PLAYER_EVENT-167]
	_ = x[MSG_TEXT_MESSAGE-168]
	_ = x[MSG_INFORM-169]
	_ = x[MSG_IMPORTANT-170]
	_ = x[MSG_IMPORTANT_ACK-171]
	_ = x[MSG_MOUSE-172]
	_ = x[MSG_INCOMING_CLIENT-173]
	_ = x[MSG_OUTGOING_CLIENT-174]
	_ = x[MSG_GAME_SETTINGS-175]
	_ = x[MSG_GAME_SETTINGS_2-176]
	_ = x[MSG_UPDATE_GUI_GAME_SETTINGS-177]
	_ = x[MSG_DOOR_ANGLE-178]
	_ = x[MSG_OBELISK_CHARGE-179]
	_ = x[MSG_PENTAGRAM_ACTIVATE-180]
	_ = x[MSG_CLIENT_PREDICT_LINEAR-181]
	_ = x[MSG_REQUEST_MAP-182]
	_ = x[MSG_CANCEL_MAP-183]
	_ = x[MSG_MAP_SEND_START-184]
	_ = x[MSG_MAP_SEND_PACKET-185]
	_ = x[MSG_MAP_SEND_ABORT-186]
	_ = x[MSG_SERVER_CMD-187]
	_ = x[MSG_SYSOP_PW-188]
	_ = x[MSG_SYSOP_RESULT-189]
	_ = x[MSG_KEEP_ALIVE-190]
	_ = x[MSG_RECEIVED_MAP-191]
	_ = x[MSG_CLIENT_READY-192]
	_ = x[MSG_REQUEST_SAVE_PLAYER-193]
	_ = x[MSG_XFER_MSG-194]
	_ = x[MSG_PLAYER_OBJ-195]
	_ = x[MSG_TEAM_MSG-196]
	_ = x[MSG_KICK_NOTIFICATION-197]
	_ = x[MSG_TIMEOUT_NOTIFICATION-198]
	_ = x[MSG_SERVER_QUIT-199]
	_ = x[MSG_SERVER_QUIT_ACK-200]
	_ = x[MSG_TRADE-201]
	_ = x[MSG_CHAT_KILL-202]
	_ = x[MSG_MESSAGES_KILL-203]
	_ = x[MSG_SEQ_IMPORTANT-204]
	_ = x[MSG_REPORT_ABILITY_AWARD-205]
	_ = x[MSG_REPORT_ABILITY_STATE-206]
	_ = x[MSG_REPORT_ACTIVE_ABILITIES-207]
	_ = x[MSG_DIALOG-208]
	_ = x[MSG_REPORT_GUIDE_AWARD-209]
	_ = x[MSG_INTERESTING_ID-210]
	_ = x[MSG_TIMER_STATUS-211]
	_ = x[MSG_REQUEST_TIMER_STATUS-212]
	_ = x[MSG_JOURNAL_MSG-213]
	_ = x[MSG_CHAPTER_END-214]
	_ = x[MSG_REPORT_ALL_LATENCY-215]
	_ = x[MSG_REPORT_FLAG_STATUS-216]
	_ = x[MSG_REPORT_BALL_STATUS-217]
	_ = x[MSG_REPORT_OBJECT_POISON-218]
	_ = x[MSG_REPORT_MONITOR_CREATURE-219]
	_ = x[MSG_REPORT_UNMONITOR_CREATURE-220]
	_ = x[MSG_REPORT_TOTAL_HEALTH-221]
	_ = x[MSG_REPORT_TOTAL_MANA-222]
	_ = x[MSG_REPORT_SPELL_STAT-223]
	_ = x[MSG_REPORT_SECONDARY_WEAPON-224]
	_ = x[MSG_REPORT_LAST_QUIVER-225]
	_ = x[MSG_INFO_BOOK_DATA-226]
	_ = x[MSG_SOCIAL-227]
	_ = x[MSG_FADE_BEGIN-228]
	_ = x[MSG_MUSIC_EVENT-229]
	_ = x[MSG_MUSIC_PUSH_EVENT-230]
	_ = x[MSG_MUSIC_POP_EVENT-231]
	_ = x[MSG_PLAYER_DIED-232]
	_ = x[MSG_PLAYER_RESPAWN-233]
	_ = x[MSG_FORGET_DRAWABLES-234]
	_ = x[MSG_RESET_ABILITIES-235]
	_ = x[MSG_RATE_CHANGE-236]
	_ = x[MSG_REPORT_CREATURE_CMD-237]
	_ = x[MSG_VOTE-238]
	_ = x[MSG_STAT_MULTIPLIERS-239]
	_ = x[MSG_GAUNTLET-240]
	_ = x[MSG_INVENTORY_FAIL-241]
}

const _Op_name = "MSG_SERVER_CONNECTMSG_SERVER_ACCEPTMSG_CODE2MSG_CODE3MSG_CODE4MSG_CODE5MSG_CLIENT_PINGMSG_CODE7MSG_CLIENT_PONGMSG_CODE9MSG_CLIENT_CLOSEMSG_SERVER_CLOSEMSG_SERVER_DISCOVERMSG_SERVER_INFOMSG_SERVER_TRY_JOINMSG_PASSWORD_REQUIREDMSG_SERVER_PINGMSG_SERVER_PASSWORDMSG_SERVER_PONGMSG_SERVER_ERRORMSG_SERVER_JOIN_OKMSG_SERVER_JOIN_FAILMSG_CODE22MSG_CODE23MSG_CODE24MSG_CODE25MSG_CODE26MSG_CODE27MSG_CODE28MSG_CODE29MSG_CODE30MSG_ACCEPTEDMSG_CLIENT_ACCEPTMSG_SERVER_CLOSE_ACKMSG_CLIENT_CLOSE_ACKMSG_SPEEDMSG_PINGMSG_CODE37MSG_CODE38MSG_TIMESTAMPMSG_FULL_TIMESTAMPMSG_NEED_TIMESTAMPMSG_SIMULATED_TIMESTAMPMSG_USE_MAPMSG_JOIN_DATAMSG_NEW_PLAYERMSG_PLAYER_QUITMSG_SIMPLE_OBJMSG_COMPLEX_OBJMSG_DESTROY_OBJECTMSG_OBJECT_OUT_OF_SIGHTMSG_OBJECT_IN_SHADOWSMSG_OBJECT_FRIEND_ADDMSG_OBJECT_FRIEND_REMOVEMSG_RESET_FRIENDSMSG_ENABLE_OBJECTMSG_DISABLE_OBJECTMSG_DRAW_FRAMEMSG_DESTROY_WALLMSG_OPEN_WALLMSG_CLOSE_WALLMSG_CHANGE_OR_ADD_WALL_MAGICMSG_REMOVE_WALL_MAGICMSG_PLAYER_INPUTMSG_PLAYER_SET_WAYPOINTMSG_REPORT_HEALTHMSG_REPORT_HEALTH_DELTAMSG_REPORT_PLAYER_HEALTHMSG_REPORT_ITEM_HEALTHMSG_REPORT_MANAMSG_REPORT_POISONMSG_REPORT_STAMINAMSG_REPORT_STATSMSG_REPORT_ARMOR_VALUEMSG_REPORT_GOLDMSG_REPORT_PICKUPMSG_REPORT_MODIFIABLE_PICKUPMSG_REPORT_DROPMSG_REPORT_LESSONMSG_REPORT_MUNDANE_ARMOR_EQUIPMSG_REPORT_MUNDANE_WEAPON_EQUIPMSG_REPORT_MODIFIABLE_WEAPON_EQUIPMSG_REPORT_MODIFIABLE_ARMOR_EQUIPMSG_REPORT_ARMOR_DEQUIPMSG_REPORT_WEAPON_DEQUIPMSG_REPORT_TREASURE_COUNTMSG_REPORT_FLAG_BALL_WINNERMSG_REPORT_FLAG_WINNERMSG_REPORT_DEATHMATCH_WINNERMSG_REPORT_DEATHMATCH_TEAM_WINNERMSG_REPORT_ENCHANTMENTMSG_REPORT_ITEM_ENCHANTMENTMSG_REPORT_LIGHT_COLORMSG_REPORT_LIGHT_INTENSITYMSG_REPORT_Z_PLUSMSG_REPORT_Z_MINUSMSG_REPORT_EQUIPMSG_REPORT_DEQUIPMSG_REPORT_ACQUIRE_SPELLMSG_REPORT_TARGETMSG_REPORT_CHARGESMSG_REPORT_X_STATUSMSG_REPORT_PLAYER_STATUSMSG_REPORT_MODIFIERMSG_REPORT_STAT_MODIFIERMSG_REPORT_NPCMSG_REPORT_CLIENT_STATUSMSG_REPORT_ANIMATION_FRAMEMSG_REPORT_ACQUIRE_CREATUREMSG_REPORT_LOSE_CREATUREMSG_REPORT_EXPERIENCEMSG_REPORT_SPELL_AWARDMSG_REPORT_SPELL_STARTMSG_REPORT_INVENTORY_LOADEDMSG_TRY_DROPMSG_TRY_GETMSG_TRY_USEMSG_TRY_EQUIPMSG_TRY_DEQUIPMSG_TRY_TARGETMSG_TRY_CREATURE_COMMANDMSG_TRY_SPELLMSG_TRY_ABILITYMSG_TRY_COLLIDEMSG_FX_PARTICLEFXMSG_FX_PLASMAMSG_FX_SUMMONMSG_FX_SUMMON_CANCELMSG_FX_SHIELDMSG_FX_BLUE_SPARKSMSG_FX_YELLOW_SPARKSMSG_FX_CYAN_SPARKSMSG_FX_VIOLET_SPARKSMSG_FX_EXPLOSIONMSG_FX_LESSER_EXPLOSIONMSG_FX_COUNTERSPELL_EXPLOSIONMSG_FX_THIN_EXPLOSIONMSG_FX_TELEPORTMSG_FX_SMOKE_BLASTMSG_FX_DAMAGE_POOFMSG_FX_LIGHTNINGMSG_FX_ENERGY_BOLTMSG_FX_CHAIN_LIGHTNING_BOLTMSG_FX_DRAIN_MANAMSG_FX_CHARMMSG_FX_GREATER_HEALMSG_FX_MAGICMSG_FX_SPARK_EXPLOSIONMSG_FX_DEATH_RAYMSG_FX_SENTRY_RAYMSG_FX_RICOCHETMSG_FX_JIGGLEMSG_FX_GREEN_BOLTMSG_FX_GREEN_EXPLOSIONMSG_FX_WHITE_FLASHMSG_FX_GENERATING_MAPMSG_FX_ASSEMBLING_MAPMSG_FX_POPULATING_MAPMSG_FX_DURATION_SPELLMSG_FX_DELTAZ_SPELL_STARTMSG_FX_TURN_UNDEADMSG_FX_ARROW_TRAPMSG_FX_VAMPIRISMMSG_FX_MANA_BOMB_CANCELMSG_UPDATE_STREAMMSG_NEW_ALIASMSG_AUDIO_EVENTMSG_AUDIO_PLAYER_EVENTMSG_TEXT_MESSAGEMSG_INFORMMSG_IMPORTANTMSG_IMPORTANT_ACKMSG_MOUSEMSG_INCOMING_CLIENTMSG_OUTGOING_CLIENTMSG_GAME_SETTINGSMSG_GAME_SETTINGS_2MSG_UPDATE_GUI_GAME_SETTINGSMSG_DOOR_ANGLEMSG_OBELISK_CHARGEMSG_PENTAGRAM_ACTIVATEMSG_CLIENT_PREDICT_LINEARMSG_REQUEST_MAPMSG_CANCEL_MAPMSG_MAP_SEND_STARTMSG_MAP_SEND_PACKETMSG_MAP_SEND_ABORTMSG_SERVER_CMDMSG_SYSOP_PWMSG_SYSOP_RESULTMSG_KEEP_ALIVEMSG_RECEIVED_MAPMSG_CLIENT_READYMSG_REQUEST_SAVE_PLAYERMSG_XFER_MSGMSG_PLAYER_OBJMSG_TEAM_MSGMSG_KICK_NOTIFICATIONMSG_TIMEOUT_NOTIFICATIONMSG_SERVER_QUITMSG_SERVER_QUIT_ACKMSG_TRADEMSG_CHAT_KILLMSG_MESSAGES_KILLMSG_SEQ_IMPORTANTMSG_REPORT_ABILITY_AWARDMSG_REPORT_ABILITY_STATEMSG_REPORT_ACTIVE_ABILITIESMSG_DIALOGMSG_REPORT_GUIDE_AWARDMSG_INTERESTING_IDMSG_TIMER_STATUSMSG_REQUEST_TIMER_STATUSMSG_JOURNAL_MSGMSG_CHAPTER_ENDMSG_REPORT_ALL_LATENCYMSG_REPORT_FLAG_STATUSMSG_REPORT_BALL_STATUSMSG_REPORT_OBJECT_POISONMSG_REPORT_MONITOR_CREATUREMSG_REPORT_UNMONITOR_CREATUREMSG_REPORT_TOTAL_HEALTHMSG_REPORT_TOTAL_MANAMSG_REPORT_SPELL_STATMSG_REPORT_SECONDARY_WEAPONMSG_REPORT_LAST_QUIVERMSG_INFO_BOOK_DATAMSG_SOCIALMSG_FADE_BEGINMSG_MUSIC_EVENTMSG_MUSIC_PUSH_EVENTMSG_MUSIC_POP_EVENTMSG_PLAYER_DIEDMSG_PLAYER_RESPAWNMSG_FORGET_DRAWABLESMSG_RESET_ABILITIESMSG_RATE_CHANGEMSG_REPORT_CREATURE_CMDMSG_VOTEMSG_STAT_MULTIPLIERSMSG_GAUNTLETMSG_INVENTORY_FAIL"

var _Op_index = [...]uint16{0, 18, 35, 44, 53, 62, 71, 86, 95, 110, 119, 135, 151, 170, 185, 204, 225, 240, 259, 274, 290, 308, 328, 338, 348, 358, 368, 378, 388, 398, 408, 418, 430, 447, 467, 487, 496, 504, 514, 524, 537, 555, 573, 596, 607, 620, 634, 649, 663, 678, 696, 719, 740, 761, 785, 802, 819, 837, 851, 867, 880, 894, 922, 943, 959, 982, 999, 1022, 1046, 1068, 1083, 1100, 1118, 1134, 1156, 1171, 1188, 1216, 1231, 1248, 1278, 1309, 1343, 1376, 1399, 1423, 1448, 1475, 1497, 1525, 1558, 1580, 1607, 1629, 1655, 1672, 1690, 1706, 1723, 1747, 1764, 1782, 1801, 1825, 1844, 1868, 1882, 1906, 1932, 1959, 1983, 2004, 2026, 2048, 2075, 2087, 2098, 2109, 2122, 2136, 2150, 2174, 2187, 2202, 2217, 2234, 2247, 2260, 2280, 2293, 2311, 2331, 2349, 2369, 2385, 2408, 2437, 2458, 2473, 2491, 2509, 2525, 2543, 2570, 2587, 2599, 2618, 2630, 2652, 2668, 2685, 2700, 2713, 2730, 2752, 2770, 2791, 2812, 2833, 2854, 2879, 2897, 2914, 2930, 2953, 2970, 2983, 2998, 3020, 3036, 3046, 3059, 3076, 3085, 3104, 3123, 3140, 3159, 3187, 3201, 3219, 3241, 3266, 3281, 3295, 3313, 3332, 3350, 3364, 3376, 3392, 3406, 3422, 3438, 3461, 3473, 3487, 3499, 3520, 3544, 3559, 3578, 3587, 3600, 3617, 3634, 3658, 3682, 3709, 3719, 3741, 3759, 3775, 3799, 3814, 3829, 3851, 3873, 3895, 3919, 3946, 3975, 3998, 4019, 4040, 4067, 4089, 4107, 4117, 4131, 4146, 4166, 4185, 4200, 4218, 4238, 4257, 4272, 4295, 4303, 4323, 4335, 4353}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
