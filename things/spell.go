package things

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
)

var (
	_ json.Marshaler   = SpellFlags(0)
	_ json.Unmarshaler = (*SpellFlags)(nil)
	_ yaml.Marshaler   = SpellFlags(0)
	_ yaml.Unmarshaler = (*SpellFlags)(nil)
)

const (
	SpellFlagUnk1       = SpellFlags(0x1)        // 1
	SpellDuration       = SpellFlags(0x2)        // 2
	SpellTargeted       = SpellFlags(0x4)        // 4
	SpellAtLocation     = SpellFlags(0x8)        // 8
	SpellMobsCanCast    = SpellFlags(0x10)       // 16
	SpellOffensive      = SpellFlags(0x20)       // 32
	SpellFlagUnk7       = SpellFlags(0x40)       // 64
	SpellFlagUnk8       = SpellFlags(0x80)       // 128
	SpellInstant        = SpellFlags(0x100)      // 256
	SpellDefensive      = SpellFlags(0x200)      // 512
	SpellFlagUnk11      = SpellFlags(0x400)      // 1024
	SpellFlagUnk12      = SpellFlags(0x800)      // 2048
	SpellSummonMain     = SpellFlags(0x1000)     // 4096
	SpellSummonCreature = SpellFlags(0x2000)     // 8192
	SpellMarkMain       = SpellFlags(0x4000)     // 16384
	SpellMarkNumber     = SpellFlags(0x8000)     // 32768
	SpellGotoMarkMain   = SpellFlags(0x10000)    // 65536
	SpellGotoMarkNumber = SpellFlags(0x20000)    // 131072
	SpellCanCounter     = SpellFlags(0x40000)    // 262144
	SpellCantHoldCrown  = SpellFlags(0x80000)    // 524288
	SpellFlagUnk21      = SpellFlags(0x100000)   // 1048576
	SpellCantTargetSelf = SpellFlags(0x200000)   // 2097152
	SpellNoTrap         = SpellFlags(0x400000)   // 4194304
	SpellNoMana         = SpellFlags(0x800000)   // 8388608
	SpellClassAny       = SpellFlags(0x1000000)  // 16777216
	SpellClassWizard    = SpellFlags(0x2000000)  // 33554432
	SpellClassConjurer  = SpellFlags(0x4000000)  // 67108864
	SpellFlagUnk28      = SpellFlags(0x8000000)  // 134217728
	SpellFlagUnk29      = SpellFlags(0x10000000) // 268435456
	SpellFlagUnk30      = SpellFlags(0x20000000) // 536870912
	SpellFlagUnk31      = SpellFlags(0x40000000) // 1073741824
	SpellFlagUnk32      = SpellFlags(0x80000000) // 2147483648
)

type SpellFlags uint32

func (f SpellFlags) string() string {
	switch f {
	case SpellDuration:
		return "DURATION"
	case SpellTargeted:
		return "TARGETED"
	case SpellAtLocation:
		return "AT_LOCATION"
	case SpellMobsCanCast:
		return "MOBS_CAN_CAST"
	case SpellOffensive:
		return "OFFENSIVE"
	case SpellInstant:
		return "INSTANT"
	case SpellDefensive:
		return "DEFENSIVE"
	case SpellSummonMain:
		return "SUMMON_SPELL"
	case SpellSummonCreature:
		return "SUMMON_CREATURE"
	case SpellMarkMain:
		return "MARK_SPELL"
	case SpellMarkNumber:
		return "MARK_NUMBER"
	case SpellGotoMarkMain:
		return "GOTO_MARK_SPELL"
	case SpellGotoMarkNumber:
		return "GOTO_MARK_NUMBER"
	case SpellCanCounter:
		return "CAN_COUNTER"
	case SpellCantHoldCrown:
		return "CANT_HOLD_CROWN"
	case SpellCantTargetSelf:
		return "CANT_TARGET_SELF"
	case SpellNoTrap:
		return "NO_TRAP"
	case SpellNoMana:
		return "NO_MANA"
	case SpellClassAny:
		return "CLASS_ANY"
	case SpellClassWizard:
		return "CLASS_WIZARD"
	case SpellClassConjurer:
		return "CLASS_CONJURER"
	}
	return ""
}

func (f SpellFlags) String() string {
	if f == 0 {
		return ""
	}
	arr := f.Split()
	if len(arr) > 1 {
		str := make([]string, 0, len(arr))
		for _, v := range arr {
			str = append(str, v.String())
		}
		return strings.Join(str, " | ")
	}
	if s := f.string(); s != "" {
		return s
	}
	return fmt.Sprintf("SpellFlags(%d)", int(f))
}

func (f SpellFlags) Has(f2 SpellFlags) bool {
	return f&f2 != 0
}

func (f SpellFlags) Split() []SpellFlags {
	var out []SpellFlags
	for i := 0; i < 32; i++ {
		v := SpellFlags(1 << i)
		if f&v != 0 {
			out = append(out, v)
		}
	}
	return out
}

func (f SpellFlags) MarshalJSON() ([]byte, error) {
	if f == 0 {
		return []byte("0"), nil
	}
	arr := f.Split()
	out := make([]interface{}, 0, len(arr))
	for _, v := range f.Split() {
		if s := v.string(); s != "" {
			out = append(out, s)
		} else {
			out = append(out, uint(v))
		}
	}
	if len(out) == 1 {
		return json.Marshal(out[0])
	}
	return json.Marshal(out)
}

func (f SpellFlags) MarshalYAML() (interface{}, error) {
	if f == 0 {
		return 0, nil
	}
	arr := f.Split()
	out := make([]interface{}, 0, len(arr))
	for _, v := range f.Split() {
		if s := v.string(); s != "" {
			out = append(out, s)
		} else {
			out = append(out, uint(v))
		}
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return out, nil
}

func (f *SpellFlags) parseText(s string) error {
	switch s {
	case "DURATION":
		*f = SpellDuration
	case "TARGET_FOE", "TARGETED":
		*f = SpellTargeted
	case "TARGET_POINT", "AT_LOCATION":
		*f = SpellAtLocation
	case "MOBS_CAN_CAST":
		*f = SpellMobsCanCast
	case "CANCELS_PROTECT", "OFFENSIVE":
		*f = SpellOffensive
	case "INSTANT":
		*f = SpellInstant
	case "AUTO_TRACK", "DEFENSIVE":
		*f = SpellDefensive
	case "SUMMON_SPELL":
		*f = SpellSummonMain
	case "SUMMON", "SUMMON_CREATURE":
		*f = SpellSummonCreature
	case "MARK_SPELL":
		*f = SpellMarkMain
	case "MARK_NUMBER":
		*f = SpellMarkNumber
	case "GOTO_MARK_SPELL":
		*f = SpellGotoMarkMain
	case "GOTO_MARK_NUMBER":
		*f = SpellGotoMarkNumber
	case "CAN_COUNTER":
		*f = SpellCanCounter
	case "CANT_HOLD_CROWN":
		*f = SpellCantHoldCrown
	case "CANT_TARGET_SELF":
		*f = SpellCantTargetSelf
	case "NO_TRAP":
		*f = SpellNoTrap
	case "NO_MANA":
		*f = SpellNoMana
	case "COMMON_USE", "CLASS_ANY":
		*f = SpellClassAny
	case "WIS_USE", "CLASS_WIZARD":
		*f = SpellClassWizard
	case "CON_USE", "CLASS_CONJURER":
		*f = SpellClassConjurer
	default:
		return fmt.Errorf("unknown spell flag: %q", s)
	}
	return nil
}

func (f *SpellFlags) unmarshalJSON(data []byte) (bool, error) {
	var v uint32
	if err := json.Unmarshal(data, &v); err == nil {
		*f = SpellFlags(v)
		return true, nil
	}
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}
	err = f.parseText(s)
	return true, err
}

func (f *SpellFlags) UnmarshalJSON(data []byte) error {
	if ok, err := f.unmarshalJSON(data); ok {
		return err
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v := SpellFlags(0)
	for _, a := range arr {
		var f2 SpellFlags
		if _, err := f2.unmarshalJSON(a); err != nil {
			return err
		}
		v |= f2
	}
	*f = v
	return nil
}

func (f *SpellFlags) unmarshalYAML(n *yaml.Node) (bool, error) {
	var v uint32
	if err := n.Decode(&v); err == nil {
		*f = SpellFlags(v)
		return true, nil
	}
	var s string
	err := n.Decode(&s)
	if err != nil {
		return false, err
	}
	err = f.parseText(s)
	return true, err
}

func (f *SpellFlags) UnmarshalYAML(n *yaml.Node) error {
	if ok, err := f.unmarshalYAML(n); ok {
		return err
	}
	var arr []SpellFlags
	if err := n.Decode(&arr); err != nil {
		return err
	}
	v := SpellFlags(0)
	for _, a := range arr {
		v |= a
	}
	*f = v
	return nil
}

type Spell struct {
	ID          string             `json:"name" yaml:"name"`
	Effect      string             `json:"effect,omitempty" yaml:"effect,omitempty"`
	Icon        *ImageRef          `json:"icon,omitempty" yaml:"icon,omitempty"`
	IconEnabled *ImageRef          `json:"icon_enabled,omitempty" yaml:"icon_enabled,omitempty"`
	ManaCost    int                `json:"mana_cost" yaml:"mana_cost"`
	Price       int                `json:"price" yaml:"price"`
	Flags       SpellFlags         `json:"flags" yaml:"flags"`
	Phonemes    []spell.Phoneme    `json:"phonemes,omitempty" yaml:"phonemes,flow,omitempty"`
	Title       strman.ID          `json:"title,omitempty" yaml:"title,omitempty"`
	Desc        strman.ID          `json:"desc,omitempty" yaml:"desc,omitempty"`
	CastSound   string             `json:"cast_sound,omitempty" yaml:"cast_sound,omitempty"`
	OnSound     string             `json:"on_sound,omitempty" yaml:"on_sound,omitempty"`
	OffSound    string             `json:"off_sound,omitempty" yaml:"off_sound,omitempty"`
	Missiles    *MissilesSpellConf `json:"missiles,omitempty" yaml:"missiles,omitempty"`
}

func ReadSpellsYAML(path string) ([]Spell, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// compatibility
	var out []Spell
	if err := yaml.Unmarshal(data, &out); err == nil {
		return out, nil
	}
	// new format - individual objects split as YAML documents
	out = nil
	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var sp Spell
		err := dec.Decode(&sp)
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return out, err
		}
		out = append(out, sp)
	}
}

func WriteSpellsYAML(path string, list []Spell) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := yaml.NewEncoder(f)
	for _, sp := range list {
		if err := enc.Encode(sp); err != nil {
			return err
		}
	}
	if err := enc.Close(); err != nil {
		return err
	}
	return f.Close()
}

func (f *Reader) ReadSpells() ([]Spell, error) {
	if err := f.seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	ok, err := f.skipUntil("SPEL")
	if !ok {
		return nil, err
	}
	return f.ReadSpellsSect()
}

// SkipSpellsSect skips one SPEL section. It assumes the reader is already positioned at the section start.
func (f *Reader) SkipSpellsSect() error {
	n, err := f.readU32()
	if err != nil {
		return err
	} else if n <= 0 {
		return nil
	}
	for i := 0; i < int(n); i++ {
		if err := f.skipBytes8(); err != nil {
			return err
		}
		if err = f.skip(1 + 2); err != nil {
			return err
		}
		sz, err := f.readU8()
		if err != nil {
			return err
		}
		for j := 0; j < int(sz); j++ {
			v, err := f.readU8()
			if err != nil {
				return err
			} else if v >= 9 {
				return fmt.Errorf("invalid spell code: %d", v)
			}
		}
		for k := 0; k < 2; k++ {
			if err := f.skipImageRef(); err != nil {
				return err
			}
		}
		if err = f.skip(4); err != nil {
			return err
		}
		if err := f.skipBytes8(); err != nil {
			return err
		}
		if err := f.skipBytes16(); err != nil {
			return err
		}
		if err := f.skipBytes8(); err != nil {
			return err
		}
		if err := f.skipBytes8(); err != nil {
			return err
		}
		if err := f.skipBytes8(); err != nil {
			return err
		}
	}
	return nil // no END here
}

// ReadSpellsSect reads one SPEL section. It assumes the reader is already positioned at the section start.
func (f *Reader) ReadSpellsSect() ([]Spell, error) {
	n, err := f.readU32()
	if err != nil {
		return nil, err
	} else if n <= 0 {
		return nil, nil
	}
	out := make([]Spell, 0, n)
	for i := 0; i < int(n); i++ {
		id, err := f.readString8()
		if err != nil {
			return out, err
		}
		mana, err := f.readU8()
		if err != nil {
			return out, err
		}
		price, err := f.readU16()
		if err != nil {
			return out, err
		}
		sz, err := f.readU8()
		if err != nil {
			return out, err
		}
		phon := make([]spell.Phoneme, 0, sz)
		for j := 0; j < int(sz); j++ {
			v, err := f.readU8()
			if err != nil {
				return out, err
			} else if v >= spell.PhonMax {
				return out, fmt.Errorf("invalid phoneme: %d", v)
			}
			phon = append(phon, spell.Phoneme(v))
		}
		im1, err := f.readImageRef()
		if err != nil {
			return out, err
		}
		im2, err := f.readImageRef()
		if err != nil {
			return out, err
		}
		fl, err := f.readU32()
		if err != nil {
			return out, err
		}
		title, err := f.readString8()
		if err != nil {
			return out, err
		}
		desc, err := f.readString16()
		if err != nil {
			return out, err
		}
		cast, err := f.readString8()
		if err != nil {
			return out, err
		}
		if cast == "NULL" {
			cast = ""
		}
		on, err := f.readString8()
		if err != nil {
			return out, err
		}
		if on == "NULL" {
			on = ""
		}
		off, err := f.readString8()
		if err != nil {
			return out, err
		}
		if off == "NULL" {
			off = ""
		}
		out = append(out, Spell{
			ID:          id,
			Icon:        im1,
			IconEnabled: im2,
			ManaCost:    int(mana),
			Price:       int(price),
			Flags:       SpellFlags(fl),
			Title:       strman.ID(title),
			Desc:        strman.ID(desc),
			Phonemes:    phon,
			CastSound:   cast,
			OnSound:     on,
			OffSound:    off,
		})
	}
	return out, nil // no END here
}
