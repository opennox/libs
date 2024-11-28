package pcx

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/noxtest"
)

func newPaletteImage(pal color.Palette) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, 256, 20), pal)
	for x := 0; x < img.Rect.Max.X; x++ {
		for y := 0; y < img.Rect.Max.Y; y++ {
			img.SetColorIndex(x, y, byte(x))
		}
	}
	return img
}

func TestDecodePalette(t *testing.T) {
	path := noxtest.DataPath(t, "default.pal")

	f, err := os.Open(path)
	must.NoError(t, err)
	defer f.Close()

	pal, err := DecodePalette(f)
	must.NoError(t, err)
	writePNG(t, "default.png", newPaletteImage(pal))
	must.Eq(t, DefaultPalette(), pal, must.Sprintf("\n%#v", pal))
}

func TestPlaceholderPalette(t *testing.T) {
	pal := PlaceholderPalette()
	writePNG(t, "palette.png", newPaletteImage(pal))
	byColor := make(map[color.Color]byte)
	for i, c := range pal {
		if i2, ok := byColor[c]; ok {
			t.Fatalf("collision between %d and %d", i, i2)
		}
		byColor[c] = byte(i)
	}
}

func writePNG(t testing.TB, path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
