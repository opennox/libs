package modifiers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/noxtest"
)

const noxModifier = "modifier.bin"

func TestParse(t *testing.T) {
	f, err := ReadFile(noxtest.DataPath(t, noxModifier))
	must.NoError(t, err)

	data, _ := json.MarshalIndent(f, "", "\t")
	os.WriteFile("modifiers.json", data, 0644)
}
