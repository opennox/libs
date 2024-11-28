package modifiers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	crypt "github.com/noxworld-dev/noxcrypt"

	"github.com/opennox/libs/ifs"
)

type File struct {
	Weapons       []WeaponOrArmor
	Armor         []WeaponOrArmor
	Enchants      []Effect
	Effectiveness []Effect
	Materials     []Effect
	Other         []Unknown `json:",omitempty"`
}

func (f *File) WeaponByName(name string) *WeaponOrArmor {
	if f == nil {
		return nil
	}
	for _, v := range f.Weapons {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (f *File) ArmorByName(name string) *WeaponOrArmor {
	if f == nil {
		return nil
	}
	for _, v := range f.Armor {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (f *File) EnchantByName(name string) *Effect {
	if f == nil {
		return nil
	}
	for _, v := range f.Enchants {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

type Unknown struct {
	Name    string
	Entries []*Entry
}

type ColorIndex int

type ColorSlot struct {
	Color   Color  `json:",omitempty"`
	Comment string `json:",omitempty"`
}

type WeaponOrArmor struct {
	Name        string   `nox:",name"`
	Desc        string   `nox:"DESC" json:",omitempty"`
	Class       []string `nox:"CLASSUSE" json:",omitempty"`
	Armor       float64  `nox:"ARMOR_VALUE" json:",omitempty"`
	Durability  int      `nox:"DURABILITY" json:",omitempty"`
	ReqStrength int      `nox:"REQUIRED_STRENGTH" json:",omitempty"`
	Range       float64  `nox:"RANGE" json:",omitempty"`
	DamageType  string   `nox:"DAMAGE_TYPE" json:",omitempty"`
	DamageMin   int      `nox:"DAMAGE_MIN" json:",omitempty"`
	DamageCoeff float64  `nox:"DAMAGE_COEFFICIENT" json:",omitempty"`

	Colors        []ColorSlot `nox:",colors" json:",omitempty"`
	Effectiveness ColorIndex  `nox:"EFFECTIVENESS,color"`
	Material      ColorIndex  `nox:"MATERIAL,color"`
	PriEnchant    ColorIndex  `nox:"PRIMARYENCHANTMENT,color"`
	SecondEnchant ColorIndex  `nox:"SECONDARYENCHANTMENT,color"`

	Unknown []KeyValue `nox:",unknown" json:",omitempty"`
}

type AllowTerm struct {
	Add bool
	Val string
}

type AllowList struct {
	Base  string
	Terms []AllowTerm `json:",omitempty"`
}

type Effect struct {
	Name         string `nox:",name"`
	Desc         string `nox:"DESC" json:",omitempty"`
	PriDesc      string `nox:"PRIMARYDESC" json:",omitempty"`
	SecondDesc   string `nox:"SECONDARYDESC" json:",omitempty"`
	IdentifyDesc string `nox:"IDENTIFYDESC" json:",omitempty"`
	Price        int    `nox:"WORTH" json:",omitempty"`
	Color        Color  `nox:"COLOR" json:",omitempty"`

	AllowWeapons AllowList `nox:"ALLOWED_WEAPONS" json:",omitempty"`
	AllowArmor   AllowList `nox:"ALLOWED_ARMOR" json:",omitempty"`
	AllowPos     []string  `nox:"ALLOWED_POSITION" json:",omitempty"`

	Attack          string `nox:"ATTACKEFFECT" json:",omitempty"`
	AttackPreHit    string `nox:"ATTACKPREHITEFFECT" json:",omitempty"`
	AttackPreDamage string `nox:"ATTACKPREDAMAGEEFFECT" json:",omitempty"`
	Defend          string `nox:"DEFENDEFFECT" json:",omitempty"`
	DefendCollide   string `nox:"DEFENDCOLLIDEEFFECT" json:",omitempty"`

	Engage    string `nox:"ENGAGEEFFECT" json:",omitempty"`
	Update    string `nox:"UPDATEEFFECT" json:",omitempty"`
	Disengage string `nox:"DISENGAGEEFFECT" json:",omitempty"`

	Unknown []KeyValue `nox:",unknown" json:",omitempty"`
}

func (s *Unknown) EntryByName(name string) *Entry {
	if s == nil {
		return nil
	}
	for _, e := range s.Entries {
		if e.Name == name {
			return e
		}
	}
	return nil
}

type Entry struct {
	Name      string
	KeyValues []KeyValue
}

func (e *Entry) GetKeyValue(key string) (KeyValue, bool) {
	if e == nil {
		return KeyValue{}, false
	}
	for _, kv := range e.KeyValues {
		if kv.Key == key {
			return kv, true
		}
	}
	return KeyValue{}, false
}

func (e *Entry) GetValue(key string) (Value, bool) {
	if e == nil {
		return nil, false
	}
	for _, kv := range e.KeyValues {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return nil, false
}

type KeyValue struct {
	Key     string
	Value   Value
	Comment string `json:",omitempty"`
}

type Value interface {
	isValue()
}

type String string

func (String) isValue() {}

type Int int

func (Int) isValue() {}

type Float float64

func (Float) isValue() {}

type Color struct {
	R, G, B int
}

func (Color) isValue() {}

func ReadFile(path string) (*File, error) {
	f, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := crypt.NewReader(f, crypt.ModifierBin)
	if err != nil {
		return nil, err
	}
	return Parse(r)
}

func replaceSpace(s []byte) []byte {
	if !bytes.ContainsAny(s, " \t\n\r") {
		return s
	}
	var (
		out   []byte
		first = true
	)
	for _, b := range s {
		switch b {
		default:
			out = append(out, b)
			first = true
		case ' ', '\t', '\n', '\r':
			if first {
				out = append(out, ' ')
				first = false
			}
		}
	}
	return out
}

func parsePlusList(s string) []string {
	sub := strings.Split(s, "+")
	for i := range sub {
		sub[i] = strings.TrimSpace(sub[i])
	}
	return sub
}

func parseAllowList(s string) AllowList {
	if s == "" {
		return AllowList{}
	}
	var sub []string
	for len(s) > 0 {
		i := strings.IndexAny(s, "+-")
		if i < 0 {
			s = strings.TrimSpace(s)
			if s != "" {
				sub = append(sub, s)
			}
			break
		}
		sub = append(sub, strings.TrimSpace(s[:i]), string(s[i]))
		s = s[i+1:]
	}
	if len(sub) == 0 {
		return AllowList{}
	}
	l := AllowList{Base: sub[0]}
	sub = sub[1:]
	for i := 0; i < len(sub); i += 2 {
		l.Terms = append(l.Terms, AllowTerm{Add: sub[i] == "+", Val: sub[i+1]})
	}
	return l
}

var (
	reflStrSlice  = reflect.TypeOf((*[]string)(nil)).Elem()
	reflAllowList = reflect.TypeOf(AllowList{})
)

func Parse(r io.Reader) (*File, error) {
	sc := bufio.NewScanner(r)
	f := new(File)

	var buf bytes.Buffer
	// Scan for sections
	for sc.Scan() {
		line := sc.Bytes()
		line = bytes.TrimSpace(line)
		line = bytes.Trim(line, "\x00")
		if len(line) == 0 {
			continue
		}
		name := string(line)
		var entries []*Entry

		// Scan for entries
		for sc.Scan() {
			line := sc.Bytes()
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			sline := string(line)
			if sline == "END" {
				break
			}
			e := &Entry{Name: sline}
			entries = append(entries, e)

			// Scan for key-value pairs
			buf.Reset()
			for sc.Scan() {
				line := sc.Bytes()
				tline := bytes.TrimSpace(line)
				if len(tline) == 0 {
					continue
				} else if bytes.Equal(tline, []byte("END")) {
					break
				}
				buf.Write(replaceSpace(line))
			}
			data := buf.Bytes()
			for len(data) > 0 {
				i := bytes.IndexByte(data, ';')
				line := data
				if i >= 0 {
					line = data[:i]
					data = data[i+1:]
				} else {
					data = nil
				}
				line = bytes.TrimSpace(line)
				if len(line) == 0 {
					continue
				}
				i = bytes.IndexByte(line, '=')
				if i < 0 {
					return nil, fmt.Errorf("invalid key-value pair: %q", string(line))
				}
				key, val := line[:i], line[i+1:]
				key = bytes.TrimSpace(key)
				val = bytes.TrimSpace(val)
				comment := ""
				if i := bytes.LastIndexByte(val, '('); i >= 0 && bytes.HasSuffix(val, []byte{')'}) {
					comment = string(bytes.TrimSpace(val[i+1 : len(val)-1]))
					val = bytes.TrimSpace(val[:i])
				}
				sval := string(val)
				var v Value = String(sval)
				if vi, err := strconv.Atoi(sval); err == nil {
					v = Int(vi)
				} else if vf, err := strconv.ParseFloat(strings.TrimSuffix(sval, "f"), 64); err == nil {
					v = Float(vf)
				} else if sub := strings.SplitN(sval, " ", 4); len(sub) == 3 {
					r, err1 := strconv.Atoi(sub[0])
					g, err2 := strconv.Atoi(sub[1])
					b, err3 := strconv.Atoi(sub[2])
					if err1 == nil && err2 == nil && err3 == nil {
						v = Color{r, g, b}
					}
				}
				e.KeyValues = append(e.KeyValues, KeyValue{
					Key:     string(key),
					Value:   v,
					Comment: comment,
				})
			}
		}
		var parr reflect.Value
		switch name {
		default:
			f.Other = append(f.Other, Unknown{Name: name, Entries: entries})
			continue
		case "WEAPON_DEFINITIONS":
			parr = reflect.ValueOf(&f.Weapons)
		case "ARMOR_DEFINITIONS":
			parr = reflect.ValueOf(&f.Armor)
		case "ENCHANTMENT":
			parr = reflect.ValueOf(&f.Enchants)
		case "EFFECTIVENESS":
			parr = reflect.ValueOf(&f.Effectiveness)
		case "MATERIAL":
			parr = reflect.ValueOf(&f.Materials)
		}
		arr := parr.Elem()
		rt := arr.Type().Elem()
		fields := make(map[string]reflect.StructField)
		clfields := make(map[string]bool)
		var fname, fcolors, funk reflect.StructField
		for i := 0; i < rt.NumField(); i++ {
			fld := rt.Field(i)
			if !fld.IsExported() {
				continue
			}
			tag := fld.Tag.Get("nox")
			if tag == "" || tag == "-" {
				continue
			}
			sub := strings.Split(tag, ",")
			if sub[0] != "" {
				fields[sub[0]] = fld
				if len(sub) > 1 && sub[1] == "color" {
					clfields[sub[0]] = true
				}
			} else if len(sub) > 1 {
				switch sub[1] {
				case "name":
					fname = fld
				case "colors":
					fcolors = fld
				case "unknown":
					funk = fld
				}
			}
		}
		for _, e := range entries {
			rv := reflect.New(rt).Elem()
			rv.FieldByIndex(fname.Index).Set(reflect.ValueOf(e.Name))
			for fname, fc := range fields {
				if clfields[fname] {
					rv.FieldByIndex(fc.Index).SetInt(-1)
				}
			}
			var rcolors reflect.Value
			if fcolors.Index != nil {
				rcolors = rv.FieldByIndex(fcolors.Index)
			}
			for _, kv := range e.KeyValues {
				fld, ok := fields[kv.Key]
				if !ok {
					if strings.HasPrefix(kv.Key, "COLOR") {
						fld = fcolors
						c, ok := kv.Value.(Color)
						if ind, err := strconv.Atoi(kv.Key[5:]); ok && err == nil {
							if rcolors.Len() < ind {
								arr := reflect.MakeSlice(rcolors.Type(), ind, ind)
								reflect.Copy(arr, rcolors)
								rcolors.Set(arr)
							}
							rcolors.Index(ind - 1).Set(reflect.ValueOf(ColorSlot{Color: c, Comment: kv.Comment}))
							continue
						}
					}
					fld = funk
					fv := rv.FieldByIndex(fld.Index)
					fv.Set(reflect.Append(fv, reflect.ValueOf(kv)))
					continue
				}
				fv := rv.FieldByIndex(fld.Index)
				if clfields[kv.Key] {
					if s, ok := kv.Value.(String); ok && strings.HasPrefix(string(s), "COLOR") {
						if ind, err := strconv.Atoi(string(s[5:])); err == nil {
							fv.Set(reflect.ValueOf(ColorIndex(ind - 1)))
							continue
						}
					}
				}
				if s, ok := kv.Value.(String); ok {
					if fv.Type() == reflStrSlice {
						fv.Set(reflect.ValueOf(parsePlusList(string(s))))
					} else if fv.Type() == reflAllowList {
						fv.Set(reflect.ValueOf(parseAllowList(string(s))))
					} else {
						fv.Set(reflect.ValueOf(kv.Value).Convert(fv.Type()))
					}
				} else {
					fv.Set(reflect.ValueOf(kv.Value).Convert(fv.Type()))
				}
			}
			arr = reflect.Append(arr, rv)
		}
		parr.Elem().Set(arr)
	}
	return f, nil
}
