// Package balance implements parsing of Nox gamedata.bin files.
package balance

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/opennox/noxcrypt"

	"github.com/opennox/libs/ifs"
)

const (
	// GamedataFile is a default filename for Nox balance file.
	GamedataFile = "gamedata.bin"
)

// Tag marks the condition for using specific values.
// It usually switches values based on the game mode.
type Tag string

const (
	TagSolo  = "SOLO"  // tag for solo games
	TagArena = "ARENA" // tag for multiplayer games
)

func newFile() *File {
	return &File{
		Global: make(Config),
		Tags:   make(map[Tag]Config),
	}
}

// Config represents a key-value map used in balance files.
type Config map[string]Value

// File represents a parsed balance config file.
type File struct {
	// Global config values that does not depend on tags.
	Global Config
	// Tags contain config overrides for specific tag values.
	Tags map[Tag]Config
	// Parent is a reference to a parent config.
	// It can be used to overlay values from one config with the other one.
	Parent *File
}

// Value returns a value for a given tag and key.
//
// Tagged values always take precedence, and if it's not set, a Global config will be used.
//
// If current File contains no value for a given key, it will fallback to Parent.
func (f *File) Value(tag Tag, k string) Value {
	if f == nil {
		return nil
	}
	k = strings.ToLower(k)
	if m, ok := f.Tags[tag]; ok {
		if v, ok := m[k]; ok {
			return v
		}
	}
	if v, ok := f.Global[k]; ok {
		return v
	}
	return f.Parent.Value(tag, k)
}

// FloatDef returns a float value for a given tag and key.
// If key is not set or has a different type, default value will be used.
//
// For lookup rules, see Value.
func (f *File) FloatDef(tag Tag, k string, def float64) float64 {
	v, ok := f.Float(tag, k)
	if !ok {
		return def
	}
	return v
}

// Float returns a float value for a given tag and key.
// If key is not set or has a different type, zero and false is returned.
//
// For lookup rules, see Value.
func (f *File) Float(tag Tag, k string) (float64, bool) {
	switch v := f.Value(tag, k).(type) {
	case nil:
		return 0, false
	case Float:
		return float64(v), true
	case Array:
		if len(v) == 0 {
			return 0, false
		}
		return v[0], true
	default:
		panic("unexpected type")
	}
}

// ArrayDef returns a float slice value for a given tag and key.
// If key is not set or has a different type, default value will be used.
//
// For lookup rules, see Value.
func (f *File) ArrayDef(tag Tag, k string, def []float64) []float64 {
	v := f.Array(tag, k)
	if len(v) == 0 {
		return def
	}
	return v
}

// Array returns a float slice value for a given tag and key.
// If key is not set or has a different type, zero and false is returned.
//
// For lookup rules, see Value.
func (f *File) Array(tag Tag, k string) []float64 {
	switch v := f.Value(tag, k).(type) {
	case nil:
		return nil
	case Float:
		return []float64{float64(v)}
	case Array:
		return v
	default:
		panic("unexpected type")
	}
}

// Value is a union for value types allowed in Config.
type Value interface {
	isValue()
}

// Float value type.
type Float float64

func (Float) isValue() {}

// Array value type.
type Array []float64

func (Array) isValue() {}

// ReadBalance reads a specified balance file. It expects the path to point to gamedata.bin.
// It will also read gamedata.yml and will use it as an overlay for the base file (if any).
func ReadBalance(path string) (*File, error) {
	orig, err := readGamedata(path)
	oerr := err
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	yml, err := readGamedataYML(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if orig == nil && yml == nil {
		return nil, oerr
	}
	if orig != nil && yml != nil {
		yml.Parent = orig
		return yml, nil
	}
	if yml != nil {
		return yml, nil
	}
	return orig, nil
}

func readGamedata(path string) (*File, error) {
	f, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := crypt.NewReader(f, crypt.GameDataBin)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(r)

	file := newFile()
	for sc.Scan() {
		raw := bytes.TrimSpace(sc.Bytes())
		raw = bytes.Trim(raw, "\x00")
		if len(raw) == 0 || raw[0] == '#' {
			continue
		}
		// value start
		sub := bytes.SplitN(raw, []byte("="), 2)
		if len(sub) < 2 || len(sub[0]) == 0 {
			return nil, fmt.Errorf("invalid gamedata line: %q", string(raw))
		}
		key := bytes.TrimSpace(sub[0])
		val := bytes.TrimSpace(sub[1])
		sub = bytes.Fields(key)
		if len(sub) > 2 {
			return nil, fmt.Errorf("invalid gamedata key: %q", string(raw))
		}
		var (
			tag  Tag
			name string
		)
		if len(sub) == 1 {
			name = string(sub[0])
		} else {
			tag = Tag(sub[0])
			name = string(sub[1])
		}
		name = strings.ToLower(name)
		dst := file.Global
		if tag != "" {
			dst = file.Tags[tag]
			if dst == nil {
				dst = make(Config)
				file.Tags[tag] = dst
			}
		}
		val = append([]byte{}, val...) // copy
		for len(val) == 0 || val[len(val)-1] != ';' {
			if !sc.Scan() {
				break
			}
			raw := bytes.TrimSpace(sc.Bytes())
			raw = bytes.Trim(raw, "\x00")
			if len(raw) != 0 && raw[0] == '#' {
				// automatically fix missing ';' if values are separated by a comment
				val = append(val, ';')
				break
			}
			val = append(val, ' ')
			val = append(val, raw...)
		}
		if len(val) == 0 {
			return nil, fmt.Errorf("invalid gamedata value: empty value for key %q", string(key))
		} else if val[len(val)-1] != ';' {
			return nil, fmt.Errorf("invalid gamedata value: missing ';' after %q", string(raw))
		}
		val = val[:len(val)-1]
		sub = bytes.FieldsFunc(val, func(r rune) bool {
			return unicode.IsSpace(r) || r == ','
		})
		out := make(Array, 0, len(sub))
		for _, v := range sub {
			if n := len(v); v[n-1] == 'f' {
				v = v[:n-1]
			}
			f, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid gamedata value: expected: number or ';', got: %q, while reading line %q", string(v), string(key)+" = "+string(val))
			}
			out = append(out, f)
		}
		if len(out) == 1 {
			dst[name] = Float(out[0])
		} else {
			dst[name] = out
		}
	}
	return file, sc.Err()
}

var (
	_ yaml.Unmarshaler = (*ymlValue)(nil)
)

type ymlValue struct {
	global Value
	tags   map[Tag]Value
}

func (v *ymlValue) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		var val float64
		if err := n.Decode(&val); err != nil {
			return err
		}
		v.global = Float(val)
		return nil
	case yaml.SequenceNode:
		var arr []float64
		if err := n.Decode(&arr); err != nil {
			return err
		}
		v.global = Array(arr)
		return nil
	default:
	}
	var mval map[Tag]float64
	if err := n.Decode(&mval); err == nil {
		v.tags = make(map[Tag]Value, len(mval))
		for key, val := range mval {
			v.tags[key] = Float(val)
		}
		return nil
	}
	var marr map[Tag][]float64
	if err := n.Decode(&marr); err == nil {
		v.tags = make(map[Tag]Value, len(mval))
		for key, val := range marr {
			v.tags[key] = Array(val)
		}
		return nil
	}
	return fmt.Errorf("cannot decode value: %q", n.Value)
}

func readGamedataYML(path string) (*File, error) {
	ext := filepath.Ext(path)
	path = strings.TrimSuffix(path, ext) + ".yml"
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m map[string]*ymlValue
	if err := yaml.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	file := newFile()
	for k, v := range m {
		k = strings.ToLower(k)
		if v.global != nil {
			file.Global[k] = v.global
			continue
		}
		for tag, val := range v.tags {
			m := file.Tags[tag]
			if m == nil {
				m = make(Config)
				file.Tags[tag] = m
			}
			m[k] = val
		}
	}
	return file, nil
}
