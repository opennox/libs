package cfg

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/noxtest"
)

func TestParseConfig(t *testing.T) {
	path := noxtest.DataPath(t)
	for _, name := range []string{
		"nox.cfg",
		"default.cfg",
	} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(path, name))
			if os.IsNotExist(err) {
				t.SkipNow()
			}
			must.NoError(t, err)
			data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

			file, err := Parse(bytes.NewReader(data))
			must.NoError(t, err)
			must.SliceLen(t, 2, file.Sections)
			vers, _ := file.Sections[0].Get("Version")
			must.NotEqOp(t, "", vers)
			mouse, _ := file.Sections[1].Get("MousePickup")
			must.NotEqOp(t, "", mouse)
		})
	}
}

var configCases = []struct {
	name string
	text string
	exp  string
	file []Section
}{
	{
		name: "empty",
		text: "",
	},
	{
		name: "section",
		text: `Version = 65537
# comment
VideoMode = 1024 768 16`,
		exp: `Version = 65537
# comment
VideoMode = 1024 768 16
---
`,
		file: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
		},
	},
	{
		name: "closed section",
		text: `Version = 65537
# comment
VideoMode = 1024 768 16
---
`,
		file: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
		},
	},
	{
		name: "two sections",
		text: `Version = 65537
# comment
VideoMode = 1024 768 16
---
MousePickup = Left
I+M = ToggleInventory`,
		exp: `Version = 65537
# comment
VideoMode = 1024 768 16
---
MousePickup = Left
I+M = ToggleInventory
---
`,
		file: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
			{
				{Key: "MousePickup", Value: "Left"},
				{Key: "I+M", Value: "ToggleInventory"},
			},
		},
	},
	{
		name: "two closed sections",
		text: `Version = 65537
# comment
VideoMode = 1024 768 16
---
MousePickup = Left
I+M = ToggleInventory
---
`,
		file: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
			{
				{Key: "MousePickup", Value: "Left"},
				{Key: "I+M", Value: "ToggleInventory"},
			},
		},
	},
	{
		name: "second only",
		text: `---
MousePickup = Left
I+M = ToggleInventory
---
`,
		file: []Section{
			nil,
			{
				{Key: "MousePickup", Value: "Left"},
				{Key: "I+M", Value: "ToggleInventory"},
			},
		},
	},
}

func TestConfig(t *testing.T) {
	for _, c := range configCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			file, err := Parse(strings.NewReader(c.text))
			must.NoError(t, err)
			must.Eq(t, &File{
				Sections: c.file,
			}, file)

			var buf bytes.Buffer
			err = file.WriteTo(&buf)
			must.NoError(t, err)
			exp := c.text
			if c.exp != "" {
				exp = c.exp
			}
			must.EqOp(t, exp, buf.String())
		})
	}
}

func TestConfigModify(t *testing.T) {
	const conf = `Version = 65537
# comment
VideoMode = 1024 768 16
---
MousePickup = Left
I+M = ToggleInventory
---
`
	file, err := Parse(strings.NewReader(conf))
	must.NoError(t, err)
	must.Eq(t, &File{
		Sections: []Section{
			{
				{Key: "Version", Value: "65537"},
				{Key: "VideoMode", Value: "1024 768 16", Comment: "comment"},
			},
			{
				{Key: "MousePickup", Value: "Left"},
				{Key: "I+M", Value: "ToggleInventory"},
			},
		},
	}, file)

	var buf bytes.Buffer
	err = file.WriteTo(&buf)
	must.NoError(t, err)
	must.EqOp(t, conf, buf.String())

	buf.Reset()
	file.Sections[0].Set("ServerName", "Test")
	err = file.WriteTo(&buf)
	must.NoError(t, err)
	must.EqOp(t, `Version = 65537
# comment
VideoMode = 1024 768 16
ServerName = Test
---
MousePickup = Left
I+M = ToggleInventory
---
`, buf.String())

}
