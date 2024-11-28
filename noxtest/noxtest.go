// Package noxtest implements different test helper functions.
package noxtest

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/datapath"
)

func DataPath(t testing.TB, sub ...string) string {
	path := datapath.Data(sub...)
	if path == "" || path == "." || !datapath.Found() {
		t.Skip("cannot detect Nox path and NOX_DATA is not set")
	}
	return path
}

func WritePNG(t testing.TB, path string, img image.Image, exp string) {
	f, err := os.Create(path)
	must.NoError(t, err)
	defer f.Close()
	h := md5.New()
	err = png.Encode(io.MultiWriter(f, h), img)
	must.NoError(t, err)
	got := hex.EncodeToString(h.Sum(nil))
	if exp != "" {
		must.EqOp(t, exp, got)
	} else {
		t.Logf("%s: %s", path, got)
	}
}
