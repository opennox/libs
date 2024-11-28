package things

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/opennox/libs/player"
)

type Draw interface {
	isDraw()
}

type BaseDraw struct {
	Img ImageRef `json:"base"`
}

func (BaseDraw) isDraw() {}

type StaticDraw struct {
	Img ImageRef `json:"static"`
}

func (StaticDraw) isDraw() {}

type WeaponDraw struct {
	Img ImageRef `json:"weapon"`
}

func (WeaponDraw) isDraw() {}

type ArmorDraw struct {
	Img ImageRef `json:"armor"`
}

func (ArmorDraw) isDraw() {}

type StaticRandomDraw struct {
	Imgs []ImageRef `json:"random"`
}

func (StaticRandomDraw) isDraw() {}

type DoorDraw struct {
	Imgs []ImageRef `json:"door"`
}

func (DoorDraw) isDraw() {}

type AnimateDraw struct {
	Anim Animation `json:"anim"`
}

func (AnimateDraw) isDraw() {}

type GlyphDraw struct {
	Anim Animation `json:"glyph"`
}

func (GlyphDraw) isDraw() {}

type WeaponAnimateDraw struct {
	Anim Animation `json:"weapon_anim"`
}

func (WeaponAnimateDraw) isDraw() {}

type ArmorAnimateDraw struct {
	Anim Animation `json:"armor_anim"`
}

func (ArmorAnimateDraw) isDraw() {}

type FlagDraw struct {
	Anim Animation `json:"flag"`
}

func (FlagDraw) isDraw() {}

type SphericalShieldDraw struct {
	Anim Animation `json:"spherical_shield"`
}

func (SphericalShieldDraw) isDraw() {}

type SummonEffectDraw struct {
	Anim Animation `json:"summon"`
}

func (SummonEffectDraw) isDraw() {}

type ConditionalAnimateDraw struct {
	Anims []Animation `json:"cond_anim"`
}

func (ConditionalAnimateDraw) isDraw() {}

type MonsterAnimationType byte

func (v MonsterAnimationType) String() string {
	if v >= 0 && int(v) < len(monsterAnimNames) {
		return monsterAnimNames[v]
	}
	return "A_" + strconv.Itoa(int(v))
}

const (
	MonsterAnimSpecial3  = MonsterAnimationType(0) // TODO static?
	MonsterAnimMelee     = MonsterAnimationType(1)
	MonsterAnimMeleeEnd  = MonsterAnimationType(2)
	MonsterAnimRanged    = MonsterAnimationType(3)
	MonsterAnimRangedEnd = MonsterAnimationType(4)
	MonsterAnimDefend    = MonsterAnimationType(5)
	MonsterAnimDefendEnd = MonsterAnimationType(6)
	MonsterAnimCast      = MonsterAnimationType(7)
	MonsterAnimIdle      = MonsterAnimationType(8)
	MonsterAnimDie       = MonsterAnimationType(9)
	MonsterAnimDead      = MonsterAnimationType(10)
	MonsterAnimHurt      = MonsterAnimationType(11)
	MonsterAnimWalk      = MonsterAnimationType(12)
	MonsterAnimRun       = MonsterAnimationType(13)
	MonsterAnimSpecial1  = MonsterAnimationType(14)
	MonsterAnimSpecial2  = MonsterAnimationType(15)
)

var monsterAnimNames = []string{
	MonsterAnimSpecial3:  "SPECIAL_3",
	MonsterAnimMelee:     "ATTACK",
	MonsterAnimMeleeEnd:  "ATTACK_FINISH",
	MonsterAnimRanged:    "ATTACK_FAR",
	MonsterAnimRangedEnd: "ATTACK_FAR_FINISH",
	MonsterAnimDefend:    "DEFEND",
	MonsterAnimDefendEnd: "DEFEND_FINISH",
	MonsterAnimCast:      "CAST_SPELL",
	MonsterAnimIdle:      "IDLE",
	MonsterAnimDie:       "DIE",
	MonsterAnimDead:      "DEAD",
	MonsterAnimHurt:      "HURT",
	MonsterAnimWalk:      "MOVE",
	MonsterAnimRun:       "MOVE_2",
	MonsterAnimSpecial1:  "SPECIAL_1",
	MonsterAnimSpecial2:  "SPECIAL_2",
}

type MonsterAnimation struct {
	Type         MonsterAnimationType `json:"type"`
	Sound        string               `json:"sound,omitempty"`
	Field8       string               `json:"field_8,omitempty"`
	FramesPerDir byte                 `json:"frames_per_dir"`
	Field10      byte                 `json:"field_10,omitempty"`
	Kind         AnimationKind        `json:"kind"`
	Frames       [8][]ImageRef        `json:"frames"`
}

type MonsterDraw struct {
	Anims []MonsterAnimation `json:"monster"`
}

func (MonsterDraw) isDraw() {}

type MaidenDraw struct {
	Anims []MonsterAnimation `json:"maiden"`
}

func (MaidenDraw) isDraw() {}

type MonsterGeneratorDraw struct {
	Anims []Animation `json:"monster_gen"`
}

func (MonsterGeneratorDraw) isDraw() {}

type PlayerAnim struct {
	FramesPerDir byte                              `json:"frames_per_dir"`
	Field8       byte                              `json:"field_8,omitempty"`
	Field12      string                            `json:"field_12,omitempty"`
	Parts        map[player.AnimPart][8][]ImageRef `json:"parts"`
}

type PlayerDraw struct {
	Anims map[player.AnimType]*PlayerAnim `json:"player"`
}

func (PlayerDraw) isDraw() {}

type UnknownDraw struct {
	Type string `json:"unknown"`
}

func (UnknownDraw) isDraw() {}

func (f *Reader) readAnimation() (*Animation, error) {
	frames, err := f.readU8()
	if err != nil {
		return nil, err
	}
	ani := &Animation{Frames: make([]ImageRef, 0, frames)}
	v2, err := f.readU8()
	if err != nil {
		return nil, err
	}
	ani.Field = v2
	loop, err := f.readBytes8()
	if err != nil {
		return nil, err
	}
	if err := ani.Kind.UnmarshalText(loop); err != nil {
		return nil, err
	}
	for i := 0; i < int(frames); i++ {
		fr, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		ani.Frames = append(ani.Frames, *fr)
	}
	return ani, nil
}

func (f *Reader) skipAnimation() error {
	n, err := f.readU8()
	if err != nil {
		return err
	}
	if err := f.skip(1); err != nil {
		return err
	}
	if err := f.skipBytes8(); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if err := f.skipImageRef(); err != nil {
			return err
		}
	}
	return err
}

func (f *Reader) readAnimations8() ([]Animation, error) {
	n, err := f.readU8()
	if err != nil {
		return nil, err
	}
	out := make([]Animation, 0, n)
	for i := 0; i < int(n); i++ {
		ani, err := f.readAnimation()
		if err != nil {
			return nil, err
		}
		out = append(out, *ani)
	}
	return out, nil
}

func (f *Reader) skipAnimations8() error {
	n, err := f.readU8()
	if err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if err := f.skipAnimation(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Reader) readMonsterDraw() (MonsterDraw, error) {
	d := MonsterDraw{}
	for {
		sect, err := f.readSect()
		if err == io.EOF {
			return d, io.ErrUnexpectedEOF
		} else if err != nil {
			return d, err
		}
		switch sect {
		default:
			return d, fmt.Errorf("unsupported monster draw sect: %q", sect)
		case "END ":
			return d, nil
		case "STAT":
			typ, err := f.readU8()
			if err != nil {
				return d, err
			}
			snd, err := f.readString8()
			if err != nil {
				return d, err
			}
			fld8, err := f.readString8()
			if err != nil {
				return d, err
			}
			framesN, err := f.readU8()
			if err != nil {
				return d, err
			}
			fld10, err := f.readU8()
			if err != nil {
				return d, err
			}
			kind, err := f.readBytes8()
			if err != nil {
				return d, err
			}
			ani := MonsterAnimation{
				Type:         MonsterAnimationType(typ),
				Sound:        snd,
				Field8:       strings.TrimRight(fld8, "\x00"),
				FramesPerDir: framesN,
				Field10:      fld10,
			}
			if err = ani.Kind.UnmarshalText(kind); err != nil {
				return d, err
			}
			for i := 0; i < 8; i++ {
				frames := make([]ImageRef, 0, int(framesN))
				for j := 0; j < int(framesN); j++ {
					ref, err := f.readImageRef()
					if err != nil {
						return d, err
					}
					frames = append(frames, *ref)
				}
				ani.Frames[i] = frames
			}
			d.Anims = append(d.Anims, ani)
		}
	}
}

func (f *Reader) readPlayerDraw() (PlayerDraw, error) {
	d := PlayerDraw{
		Anims: make(map[player.AnimType]*PlayerAnim),
	}
	var lastAni *PlayerAnim
	for {
		sect, err := f.readSect()
		if err == io.EOF {
			return d, io.ErrUnexpectedEOF
		} else if err != nil {
			return d, err
		}
		switch sect {
		default:
			return d, fmt.Errorf("unsupported monster draw sect: %q", sect)
		case "END ":
			return d, nil
		case "STAT":
			name, err := f.readString8()
			if err != nil {
				return d, err
			}
			framesN, err := f.readU8()
			if err != nil {
				return d, err
			}
			fld8, err := f.readU8()
			if err != nil {
				return d, err
			}
			fld12, err := f.readString8()
			if err != nil {
				return d, err
			}
			ani := &PlayerAnim{
				FramesPerDir: framesN,
				Field8:       fld8,
				Field12:      fld12,
			}
			lastAni = ani
			d.Anims[player.AnimType(name)] = ani
		case "SEQU":
			name, err := f.readString8()
			if err != nil {
				return d, err
			}
			framesN := int(lastAni.FramesPerDir)
			var dirs [8][]ImageRef
			for i := 0; i < 8; i++ {
				frames := make([]ImageRef, 0, framesN)
				for j := 0; j < framesN; j++ {
					ref, err := f.readImageRef()
					if err != nil {
						return d, err
					}
					frames = append(frames, *ref)
				}
				dirs[i] = frames
			}
			if lastAni.Parts == nil {
				lastAni.Parts = make(map[player.AnimPart][8][]ImageRef)
			}
			lastAni.Parts[player.AnimPart(name)] = dirs
		}
	}
}

func (f *Reader) skipThingDraw() error {
	dname, err := f.readString8()
	if err != nil {
		return err
	}
	sectSz, err := f.readU64align()
	if err != nil {
		return err
	}
	switch dname {
	case "StaticDraw", "WeaponDraw", "ArmorDraw", "BaseDraw":
		if err := f.skipImageRef(); err != nil {
			return err
		}
	case "StaticRandomDraw", "DoorDraw":
		if err := f.skipImageRefs8(); err != nil {
			return err
		}
	case "AnimateDraw", "GlyphDraw", "WeaponAnimateDraw", "ArmorAnimateDraw",
		"FlagDraw", "SphericalShieldDraw", "SummonEffectDraw":
		if err := f.skipAnimation(); err != nil {
			return err
		}
	case "ConditionalAnimateDraw", "MonsterGeneratorDraw":
		if err := f.skipAnimations8(); err != nil {
			return err
		}
	case "AnimateStateDraw":
		if err := f.skipThingAnimStateDraw(); err != nil {
			return err
		}
	case "VectorAnimateDraw", "ReleasedSoulDraw":
		if err := f.skipThingAnimVectorDraw(); err != nil {
			return err
		}
	case "MonsterDraw", "MaidenDraw":
		if err := f.skipThingMonsterDraw(); err != nil {
			return err
		}
	case "PlayerDraw":
		if err := f.skipThingPlayerDraw(); err != nil {
			return err
		}
	case "SlaveDraw", "BoulderDraw", "ArrowDraw", "WeakArrowDraw", "HarpoonDraw":
		if err := f.skipImageRefs8(); err != nil {
			return err
		}
	default:
		if err := f.skip(int(sectSz)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Reader) readThingDraw() (Draw, error) {
	dname, err := f.readString8()
	if err != nil {
		return nil, err
	}
	sectSz, err := f.readU64align()
	if err != nil {
		return nil, err
	}
	switch dname {
	case "", "NoDraw":
		return nil, nil
	case "StaticDraw", "WeaponDraw", "ArmorDraw", "BaseDraw":
		img, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "StaticDraw":
			return StaticDraw{Img: *img}, nil
		case "WeaponDraw":
			return WeaponDraw{Img: *img}, nil
		case "ArmorDraw":
			return ArmorDraw{Img: *img}, nil
		case "BaseDraw":
			return BaseDraw{Img: *img}, nil
		}
	case "StaticRandomDraw", "DoorDraw":
		imgs, err := f.readImageRefs8()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "StaticRandomDraw":
			return StaticRandomDraw{Imgs: imgs}, nil
		case "DoorDraw":
			return DoorDraw{Imgs: imgs}, nil
		}
	case "AnimateDraw", "GlyphDraw", "WeaponAnimateDraw", "ArmorAnimateDraw",
		"FlagDraw", "SphericalShieldDraw", "SummonEffectDraw":
		anim, err := f.readAnimation()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "AnimateDraw":
			return AnimateDraw{Anim: *anim}, nil
		case "GlyphDraw":
			return GlyphDraw{Anim: *anim}, nil
		case "WeaponAnimateDraw":
			return WeaponAnimateDraw{Anim: *anim}, nil
		case "ArmorAnimateDraw":
			return ArmorAnimateDraw{Anim: *anim}, nil
		case "FlagDraw":
			return FlagDraw{Anim: *anim}, nil
		case "SphericalShieldDraw":
			return SphericalShieldDraw{Anim: *anim}, nil
		case "SummonEffectDraw":
			return SummonEffectDraw{Anim: *anim}, nil
		}
	case "ConditionalAnimateDraw", "MonsterGeneratorDraw":
		anis, err := f.readAnimations8()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "ConditionalAnimateDraw":
			return ConditionalAnimateDraw{Anims: anis}, nil
		case "MonsterGeneratorDraw":
			return MonsterGeneratorDraw{Anims: anis}, nil
		}
	case "AnimateStateDraw":
		if err := f.skipThingAnimStateDraw(); err != nil {
			return nil, err // FIXME
		}
		return UnknownDraw{Type: dname}, nil
	case "VectorAnimateDraw", "ReleasedSoulDraw":
		if err := f.skipThingAnimVectorDraw(); err != nil {
			return nil, err // FIXME
		}
		return UnknownDraw{Type: dname}, nil
	case "MonsterDraw", "MaidenDraw":
		d, err := f.readMonsterDraw()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "MonsterDraw":
			return d, nil
		case "MaidenDraw":
			return MaidenDraw(d), nil
		}
	case "PlayerDraw":
		d, err := f.readPlayerDraw()
		if err != nil {
			return nil, err
		}
		return d, nil
	case "SlaveDraw", "BoulderDraw", "ArrowDraw", "WeakArrowDraw", "HarpoonDraw":
		if err := f.skipImageRefs8(); err != nil {
			return nil, err // FIXME
		}
	default:
		if err := f.skip(int(sectSz)); err != nil {
			return nil, err
		}
	}
	return UnknownDraw{Type: dname}, nil
}
