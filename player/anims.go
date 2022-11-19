package player

type AnimType string

func (v AnimType) String() string {
	return string(v)
}

type AnimPart string

func (v AnimPart) String() string {
	return string(v)
}

const (
	AnimAxeStrike            = AnimType("AXE_STRIKE")
	AnimBerserkerCharge      = AnimType("BERSERKER_CHARGE")
	AnimCastSpell            = AnimType("CAST_SPELL")
	AnimChakramStrike        = AnimType("CHAKRAM_STRIKE")
	AnimConcentrate          = AnimType("CONCENTRATE")
	AnimDead                 = AnimType("DEAD")
	AnimDie                  = AnimType("DIE")
	AnimDust                 = AnimType("DUST")
	AnimDodgeLeft            = AnimType("DODGE_LEFT")
	AnimDodgeRight           = AnimType("DODGE_RIGHT")
	AnimElectricZap          = AnimType("ELECTRIC_ZAP")
	AnimElectrocuted         = AnimType("ELECTROCUTED")
	AnimFall                 = AnimType("FALL")
	AnimGetUp                = AnimType("GET_UP")
	AnimGreatSwordBlockDown  = AnimType("GREAT_SWORD_BLOCK_DOWN")
	AnimGreatSwordBlockLeft  = AnimType("GREAT_SWORD_BLOCK_LEFT")
	AnimGreatSwordBlockRight = AnimType("GREAT_SWORD_BLOCK_RIGHT")
	AnimGreatSwordIdle       = AnimType("GREAT_SWORD_IDLE")
	AnimGreatSwordStrike     = AnimType("GREAT_SWORD_STRIKE")
	AnimGreatSwordParry      = AnimType("GREAT_SWORD_PARRY")
	AnimHammerStrike         = AnimType("HAMMER_STRIKE")
	AnimHarpoonThrow         = AnimType("HARPOONTHROW")
	AnimIdle                 = AnimType("IDLE")
	AnimJump                 = AnimType("JUMP")
	AnimLaugh                = AnimType("LAUGH")
	AnimLongSwordStrike      = AnimType("LONG_SWORD_STRIKE")
	AnimMaceStrike           = AnimType("MACE_STRIKE")
	AnimPickup               = AnimType("PICKUP")
	AnimPoint                = AnimType("POINT")
	AnimPunchLeft            = AnimType("PUNCH_LEFT")
	AnimPunchRight           = AnimType("PUNCH_RIGHT")
	AnimPunchRightHook       = AnimType("PUNCH_RIGHT_HOOK")
	AnimRaiseShield          = AnimType("RAISE_SHIELD")
	AnimRecoil               = AnimType("RECOIL")
	AnimRecoilShield         = AnimType("RECOIL_SHIELD")
	AnimRecoilForward        = AnimType("RECOIL_FORWARD")
	AnimRecoilBackward       = AnimType("RECOIL_BACKWARD")
	AnimRun                  = AnimType("RUN")
	AnimRunningJump          = AnimType("RUNNING_JUMP")
	AnimSit                  = AnimType("SIT")
	AnimShootBow             = AnimType("SHOOT_BOW")
	AnimShootCrossbow        = AnimType("SHOOT_CROSSBOW")
	AnimSleep                = AnimType("SLEEP")
	AnimSneak                = AnimType("SNEAK")
	AnimStaffBlock           = AnimType("STAFF_BLOCK")
	AnimStaffSpellBlast      = AnimType("STAFF_SPELL_BLAST")
	AnimStaffStrike          = AnimType("STAFF_STRIKE")
	AnimStaffThrust          = AnimType("STAFF_THRUST")
	AnimSwordStrike          = AnimType("SWORD_STRIKE")
	AnimTalk                 = AnimType("TALK")
	AnimTaunt                = AnimType("TAUNT")
	AnimTrip                 = AnimType("TRIP")
	AnimWalk                 = AnimType("WALK")
	AnimWalkAndDrag          = AnimType("WALK_AND_DRAG")
	AnimWarcry               = AnimType("WARCRY")
)

const (
	PartAxe                    = AnimPart("AXE")
	PartBow                    = AnimPart("BOW")
	PartChainCoif              = AnimPart("CHAIN_COIF")
	PartChainLeggins           = AnimPart("CHAIN_LEGGINGS")
	PartChainTunic             = AnimPart("CHAIN_TUNIC")
	PartChakram                = AnimPart("CHAKRAM")
	PartConjurerHelm           = AnimPart("CONJURER_HELM")
	PartCrossbow               = AnimPart("CROSSBOW")
	PartGreatSword             = AnimPart("GREAT_SWORD")
	PartHammer                 = AnimPart("HAMMER")
	PartKiteShield             = AnimPart("KITE_SHIELD")
	PartLeatherARMBANDS        = AnimPart("LEATHER_ARMBANDS")
	PartLeatherArmoredBoots    = AnimPart("LEATHER_ARMORED_BOOTS")
	PartLeatherBoots           = AnimPart("LEATHER_BOOTS")
	PartLeatherHelm            = AnimPart("LEATHER_HELM")
	PartLeatherLeggings        = AnimPart("LEATHER_LEGGINGS")
	PartLeatherTunic           = AnimPart("LEATHER_TUNIC")
	PartLongSword              = AnimPart("LONG_SWORD")
	PartMace                   = AnimPart("MACE")
	PartMedievalCloak          = AnimPart("MEDIEVAL_CLOAK")
	PartMedievalPants          = AnimPart("MEDIEVAL_PANTS")
	PartMedievalShirt          = AnimPart("MEDIEVAL_SHIRT")
	PartNaked                  = AnimPart("NAKED")
	PartOgreAxe                = AnimPart("OGRE_AXE")
	PartOrnateHelm             = AnimPart("ORNATE_HELM")
	PartPlateArms              = AnimPart("PLATE_ARMS")
	PartPlateBoots             = AnimPart("PLATE_BOOTS")
	PartPlateBreast            = AnimPart("PLATE_BREAST")
	PartPlateHelm              = AnimPart("PLATE_HELM")
	PartPlateLeggings          = AnimPart("PLATE_LEGGINGS")
	PartQuiver                 = AnimPart("QUIVER")
	PartRoundShield            = AnimPart("ROUND_SHIELD")
	PartShuriken               = AnimPart("SHURIKEN")
	PartStaff                  = AnimPart("STAFF")
	PartStaffDeathRay          = AnimPart("STAFF_DEATH_RAY")
	PartStaffFireball          = AnimPart("STAFF_FIREBALL")
	PartStaffForceOfNature     = AnimPart("STAFF_FORCE_OF_NATURE")
	PartStaffLIGHTNING         = AnimPart("STAFF_LIGHTNING")
	PartStaffOblivionHalberd   = AnimPart("STAFF_OBLIVION_HALBERD")
	PartStaffOblivionHeart     = AnimPart("STAFF_OBLIVION_HEART")
	PartStaffOblivionOrb       = AnimPart("STAFF_OBLIVION_ORB")
	PartStaffOblivionWierdling = AnimPart("STAFF_OBLIVION_WIERDLING")
	PartStaffSulphorousFlare   = AnimPart("STAFF_SULPHOROUS_FLARE")
	PartStaffSulphorousShower  = AnimPart("STAFF_SULPHOROUS_SHOWER")
	PartStaffTripleFireball    = AnimPart("STAFF_TRIPLE_FIREBALL")
	PartStreetPants            = AnimPart("STREET_PANTS")
	PartStreetShirt            = AnimPart("STREET_SHIRT")
	PartStreetSneakers         = AnimPart("STREET_SNEAKERS")
	PartSword                  = AnimPart("SWORD")
	PartWizardHelm             = AnimPart("WIZARD_HELM")
	PartWizardRobe             = AnimPart("WIZARD_ROBE")
)