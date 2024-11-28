package maps_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	crypt "github.com/opennox/noxcrypt"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/shoenig/test/must"

	"github.com/opennox/opennox-lib/ifs"
	"github.com/opennox/opennox-lib/maps"
	"github.com/opennox/opennox-lib/noxtest"
)

var casesMapInfo = []maps.Info{
	{
		Filename: "con01a",
		Size:     167256,
		MapInfo: maps.MapInfo{
			Format:      2,
			Summary:     "Second Warrior Mission ",
			Description: "Town of Ix\t",
			Author:      "Bryan Hansen/John Lee",
			Email:       "bhansen@westwood.com",
			Date:        "Friday, January 7 2000",
			Flags:       0x3,
			MinPlayers:  2,
			MaxPlayers:  16,
			Trailing: maps.MapInfoCompat{
				Email: "m",
				Date:  "9\x009",
			},
		},
	},
	{
		Filename: "estate",
		Size:     59064,
		MapInfo: maps.MapInfo{
			Format:      2,
			Summary:     "Death Match suitable for 2 - 8 players",
			Description: "Outdoor woodland setting for killer Deathmatches",
			Version:     "2.1",
			Author:      "Bryan Hansen",
			Copyright:   "Copyright 1999 Westwood Studios.  All rights reserved.",
			Date:        "Monday, January 3 2000",
			Flags:       0x34,
			MinPlayers:  2,
			MaxPlayers:  8,
			Trailing: maps.MapInfoCompat{
				Email: "hansen@westwood.com",
				Date:  "999",
			},
		},
	},
	{
		Filename: "g_castle",
		Size:     475264,
		MapInfo: maps.MapInfo{
			Format:        3,
			Author:        "John Lee/Bryan Hansen",
			Author2:       "Phil Robb",
			Date:          "Monday, July 17 2000",
			Flags:         0x2,
			MinPlayers:    2,
			MaxPlayers:    16,
			QuestIntro:    "QIntro.dat:GauntletCastleText",
			QuestGraphics: "WizardChapterBegin2",
			Trailing: maps.MapInfoCompat{
				Date: "00",
			},
		},
	},
	{
		Filename: "g_mines",
		Size:     652432,
		MapInfo: maps.MapInfo{
			Format:        3,
			Author:        "John Lee",
			Author2:       "Phil Robb",
			Date:          "Tuesday, July 18 2000",
			Flags:         0x2,
			MinPlayers:    2,
			MaxPlayers:    16,
			QuestIntro:    "QIntro.dat:GauntletMinesText",
			QuestGraphics: "WarriorChapterBegin8",
		},
	},
	{
		Filename: "so_brin",
		Size:     12368,
		MapInfo: maps.MapInfo{
			Format:      2,
			Summary:     "Brin Social Map",
			Description: "Social map set in Brin Farm",
			Author:      "Jeremiah Cohn",
			Date:        "Monday, January 3 2000",
			Flags:       0x80000000,
			MinPlayers:  2,
			MaxPlayers:  16,
			Trailing: maps.MapInfoCompat{
				Summary:     "p",
				Description: "rary",
				Date:        "\x009",
			},
		},
	},
	{
		Filename: "war01a",
		Size:     341312,
		MapInfo: maps.MapInfo{
			Format:     2,
			Summary:    "Warrior Chapter 1a",
			Author:     "Eric Beaumont",
			Date:       "Saturday, January 8 2000",
			Flags:      0x1,
			MinPlayers: 2,
			MaxPlayers: 16,
			Trailing: maps.MapInfoCompat{
				Summary: " map) ",
				Date:    "99",
			},
		},
	},
}

func TestReadFileInfo(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	for _, m := range casesMapInfo {
		t.Run(m.Filename, func(t *testing.T) {
			info, err := maps.ReadMapInfo(filepath.Join(path, m.Filename))
			must.NoError(t, err)
			must.Eq(t, m, *info)
		})
	}
}

func TestReadFile(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	for _, m := range casesMapInfo {
		t.Run(m.Filename, func(t *testing.T) {
			mp, err := maps.ReadMap(filepath.Join(path, m.Filename))
			must.NoError(t, err)
			for _, s := range mp.Unknown {
				t.Logf("unknwon section: %q [%d]", s.Name, len(s.Data))
			}
			if mp.Script != nil {
				if len(mp.Script.Data) == 0 {
					t.Logf("script [%d]", len(mp.Script.Data))
				} else {
					sc, err := asm.ReadScript(bytes.NewReader(mp.Script.Data))
					must.NoError(t, err)
					t.Logf("script [%d]: %d funcs, %d strings", len(mp.Script.Data), len(sc.Funcs), len(sc.Strings))
				}
			}
		})
	}
}

func TestMapSections(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	list, err := os.ReadDir(path)
	must.NoError(t, err)
	for _, fi := range list {
		if !fi.IsDir() {
			continue
		}
		fname := filepath.Join(path, fi.Name(), fi.Name()+".map")
		if _, err := ifs.Stat(fname); os.IsNotExist(err) {
			continue
		}
		t.Run(strings.ToLower(fi.Name()), func(t *testing.T) {
			f, err := ifs.Open(fname)
			must.NoError(t, err)
			defer f.Close()
			rd, err := maps.NewReader(f)
			must.NoError(t, err)
			sect, err := rd.ReadSectionsRaw()
			must.NoError(t, err)
			for _, s := range sect {
				if !s.Supported() {
					t.Logf("skip section: %q", s.Name)
					continue
				}
				t.Run(s.Name, func(t *testing.T) {
					d, err := s.Decode()
					must.NoError(t, err)
					//t.Logf("%#v", d)
					data, err := d.MarshalBinary()
					must.NoError(t, err)
					must.Eq(t, s.Data, data, must.Sprintf("%q", s.Data))
				})
			}
		})
	}
}

func TestMapWrite(t *testing.T) {
	path := noxtest.DataPath(t, maps.Dir)
	list, err := os.ReadDir(path)
	must.NoError(t, err)
	for _, fi := range list {
		if !fi.IsDir() {
			continue
		}
		fname := filepath.Join(path, fi.Name(), fi.Name()+".map")
		if _, err := ifs.Stat(fname); os.IsNotExist(err) {
			continue
		}
		t.Run(strings.ToLower(fi.Name()), func(t *testing.T) {
			f, err := ifs.Open(fname)
			must.NoError(t, err)
			defer f.Close()

			var exp bytes.Buffer

			rd, err := maps.NewReader(io.TeeReader(f, &exp))
			must.NoError(t, err)
			sect, err := rd.ReadSectionsRaw()
			must.NoError(t, err)
			t.Logf("map crc: 0x%x", rd.Map().CRC())

			var got buffer
			wr, err := maps.NewWriter(&got, rd.Header())
			must.NoError(t, err)
			err = wr.WriteRawSections(sect)
			must.NoError(t, err)
			err = wr.Close()
			must.NoError(t, err)
			if bexp, bgot := exp.Bytes(), got.Bytes(); !bytes.Equal(bexp, bgot) {
				must.Eq(t, decodeMapBytes(bexp), decodeMapBytes(bgot))
				must.Eq(t, bexp, bgot)
			}
		})
	}
}

func decodeMapBytes(data []byte) []byte {
	cr, err := crypt.NewReader(bytes.NewReader(data), crypt.MapKey)
	if err != nil {
		panic(err)
	}
	raw, err := io.ReadAll(cr)
	if err != nil {
		panic(err)
	}
	return raw
}

type buffer struct {
	data []byte
}

func (b *buffer) Bytes() []byte {
	return b.data
}

func (b *buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *buffer) WriteAt(p []byte, off int64) (int, error) {
	n := copy(b.data[off:], p)
	if n != len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}
