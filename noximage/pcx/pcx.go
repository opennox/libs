package pcx

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	noxcolor "github.com/opennox/libs/color"
)

var (
	endiness = binary.LittleEndian
)

// PlaceholderPalette is a placeholder palette used by PCX sprite decoder for different materials.
// It is set mostly to allow encoding PCX images directly to a pair of PNG images.
// Engine doesn't use this palette and instead sets its own palette based on the object's materials.
func PlaceholderPalette() color.Palette {
	const (
		colors = 6
		part   = 255/6 + 1
		max    = 255
	)
	pal := make(color.Palette, 256)
	pal[0] = color.White
	for i := 0; i < 255; i++ {
		up := colors * (byte(i) % part)
		down := 254 - up
		switch i % colors {
		case 0:
			pal[1+i] = color.NRGBA{R: max, G: up, B: 0, A: 0xff}
		case 1:
			pal[1+i] = color.NRGBA{R: down, G: max, B: 0, A: 0xff}
		case 2:
			pal[1+i] = color.NRGBA{R: 0, G: max, B: up, A: 0xff}
		case 3:
			pal[1+i] = color.NRGBA{R: 0, G: down, B: max, A: 0xff}
		case 4:
			pal[1+i] = color.NRGBA{R: up, G: 0, B: max, A: 0xff}
		case 5:
			pal[1+i] = color.NRGBA{R: max, G: 0, B: down, A: 0xff}
		default:
			panic(i)
		}
	}
	return pal
}

var _ image.Image = (*Image)(nil)

type ImageMeta struct {
	Type  byte        `json:"type,omitempty"`
	Point image.Point `json:"point"`
}

type Image struct {
	image.Image
	ImageMeta
	// Material stores a color multiplier as a paletted image. Indexes in the palette are usually interpreted as materials.
	// Palette index value 0 indicates no material, meaning that color can be read from NRGBA directly.
	// All other palette indexes are moved by +1, comparing to the ones passed expected by the engine.
	Material *image.Paletted
}

func Decode(r io.Reader, typ byte) (*Image, error) {
	switch typ {
	case 0:
		return readTile(r, typ)
	case 2, 7:
		return readSprite2(r, typ)
	case 3, 4, 5, 6, 8:
		return readSprite(r, typ)
	default:
		return nil, fmt.Errorf("unsupported type: %d", typ)
	}
}

func DecodeHeader(r io.Reader, typ byte) (*ImageMeta, image.Point, error) {
	switch typ {
	case 3, 4, 5, 6, 8:
		meta, sz, err := readSpriteHeader(r, typ)
		return meta, sz, err
	default:
		return nil, image.Point{}, fmt.Errorf("unsupported type: %d", typ)
	}
}

func readSpriteHeader(r io.Reader, typ byte) (*ImageMeta, image.Point, error) {
	var b [17]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return nil, image.Point{}, err
	}
	width := int(endiness.Uint32(b[0:]))
	height := int(endiness.Uint32(b[4:]))
	offsX := int(endiness.Uint32(b[8:]))
	offsY := int(endiness.Uint32(b[12:]))
	offs := image.Pt(offsX, offsY)
	// one byte ignored
	if width <= 0 || width > 1024 || height <= 0 || height > 1024 {
		return nil, image.Point{}, fmt.Errorf("invalid image size: %dx%d", width, height)
	}
	return &ImageMeta{
			Type:  typ,
			Point: offs,
		}, image.Point{
			X: width,
			Y: height,
		}, nil
}

func readSprite(r io.Reader, typ byte) (*Image, error) {
	meta, sz, err := readSpriteHeader(r, typ)
	if err != nil {
		return nil, err
	}
	width := sz.X
	height := sz.Y

	rgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	img := &Image{
		Image:     rgba,
		ImageMeta: *meta,
	}
	br := bufio.NewReader(r)

	var buf []byte
	growBuf := func(n int) {
		if cap(buf) < n {
			buf = make([]byte, n)
		} else {
			buf = buf[:n]
		}
	}
	for y := 0; y < height; y++ {
		var dx int
		for x := 0; x < width; x += dx {
			op, err := br.ReadByte()
			if err == io.EOF {
				return img, io.ErrUnexpectedEOF
			} else if err != nil {
				return img, err
			}
			val, err := br.ReadByte()
			if err == io.EOF {
				return img, io.ErrUnexpectedEOF
			} else if err != nil {
				return img, err
			}
			dx = int(val)
			if op&0xF == 1 {
				continue // skip
			}
			if typ == 8 {
				switch op & 0xF {
				case 2, 7: // RGB565
					growBuf(2 * dx)
					_, err = io.ReadFull(br, buf)
					if err != nil {
						return img, err
					}
					for i := 0; i < dx; i++ {
						cl := noxcolor.RGB565(endiness.Uint16(buf[2*i:]))
						rgba.SetNRGBA(x+i, y, cl.ColorNRGBA())
					}
				default:
					return img, fmt.Errorf("invalid draw op (image type %d): 0x%x, (%d,%d)", typ, op, x, y)
				}
			} else {
				switch op & 0xF {
				case 2, 7: // RGB555
					growBuf(2 * dx)
					_, err = io.ReadFull(br, buf)
					if err != nil {
						return img, err
					}
					for i := 0; i < dx; i++ {
						cl := noxcolor.RGB555(endiness.Uint16(buf[2*i:]))
						rgba.SetNRGBA(x+i, y, cl.ColorNRGBA())
					}
				case 4:
					growBuf(dx)
					_, err = io.ReadFull(br, buf)
					if err != nil {
						return img, err
					}
					if img.Material == nil {
						img.Material = image.NewPaletted(rgba.Rect, PlaceholderPalette())
					}
					// we move it by 1, so that 0 on the indexed mask image will indicate a non-indexed pixel
					ind := (op >> 4) + 1
					for i := 0; i < dx; i++ {
						cl := buf[i]
						rgba.SetNRGBA(x+i, y, color.NRGBA{R: cl, G: cl, B: cl, A: 0xff})
						img.Material.SetColorIndex(x+i, y, ind)
					}
				case 5: // RGBA4444
					growBuf(2 * dx)
					_, err = io.ReadFull(br, buf)
					if err != nil {
						return img, err
					}
					for i := 0; i < dx; i++ {
						cl := noxcolor.RGBA4444(endiness.Uint16(buf[2*i:]))
						rgba.SetNRGBA(x+i, y, cl.ColorNRGBA())
					}
				case 6:
					// nop
				default:
					return img, fmt.Errorf("invalid draw op (image type %d): 0x%x, (%d,%d)", typ, op, x, y)
				}
			}
		}
	}
	return img, nil
}

func readSprite2(r io.Reader, typ byte) (*Image, error) {
	meta, sz, err := readSpriteHeader(r, typ)
	if err != nil {
		return nil, err
	}
	width := sz.X
	height := sz.Y

	rgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	img := &Image{
		Image:     rgba,
		ImageMeta: *meta,
	}
	br := bufio.NewReader(r)
	buf := make([]byte, 2*width)
	for y := 0; y < height; y++ {
		_, err = io.ReadFull(br, buf)
		if err != nil {
			return img, err
		}
		for x := 0; x < width; x++ {
			cl := noxcolor.RGBA5551(endiness.Uint16(buf[2*x:]))
			rgba.SetNRGBA(x, y, cl.ColorNRGBA())
		}
	}
	return img, nil
}

const (
	tileSize = 46
)

func readTile(r io.Reader, typ byte) (*Image, error) {
	const (
		width  = tileSize
		height = tileSize
	)
	buf := make([]byte, width*height)
	_, err := io.ReadFull(r, buf)
	if err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	}

	rgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	img := &Image{
		Image:     rgba,
		ImageMeta: ImageMeta{Type: typ},
	}
	var (
		dx = tileSize / 2
		n  = 1
	)
	for y := 0; y < height; y++ {
		for x := 0; x < n; x++ {
			ind := binary.LittleEndian.Uint16(buf[0:2])
			buf = buf[2:]
			rgba.SetNRGBA(x+dx, y, noxcolor.RGBA5551(ind).ColorNRGBA())
		}
		if y < tileSize/2-1 {
			n += 2
			dx--
		} else if y > tileSize/2-1 {
			n -= 2
			dx++
		}
	}
	return img, nil
}

func asNRGBA(img image.Image) *image.NRGBA {
	if rgba, ok := img.(*image.NRGBA); ok {
		return rgba
	}
	rgba := image.NewNRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Rect, img, img.Bounds().Min, draw.Src)
	return rgba
}

func Encode(img *Image) []byte {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	data := make([]byte, 17+width*height*4+2)
	endiness.PutUint32(data[0:], uint32(width))
	endiness.PutUint32(data[4:], uint32(height))
	endiness.PutUint32(data[8:], uint32(int32(img.Point.X)))
	endiness.PutUint32(data[12:], uint32(int32(img.Point.Y)))
	if img.Type == 0 {
		data[16] = 0x3
	} else {
		data[16] = img.Type
	}
	pixdata := data[17:]
	pixdata = pixdata[:0] // will append to it
	rgba := asNRGBA(img.Image)
	mat := img.Material
	const (
		modeZero = iota + 1
		modeRGB15
		modeRGBA16
		modeMask
	)
	const (
		pixMax = 0xfd
	)
	var (
		pbuf  [2]byte
		pmode = -1
		ni    = -1
	)
	addZero := func() {
		if pmode == modeZero {
			if n := pixdata[ni]; n < pixMax {
				pixdata[ni]++
				return
			}
		}
		pmode = modeZero
		pixdata = append(pixdata, 1, 1)
		ni = len(pixdata) - 1
	}
	addRGB15 := func(p noxcolor.RGB555) {
		endiness.PutUint16(pbuf[:], uint16(p))
		if pmode == modeRGB15 {
			if n := pixdata[ni]; n < pixMax {
				pixdata[ni]++
				pixdata = append(pixdata, pbuf[0], pbuf[1])
				return
			}
		}
		pmode = modeRGB15
		pixdata = append(pixdata, 2, 1, pbuf[0], pbuf[1])
		ni = len(pixdata) - 3
	}
	addRGBA16 := func(p noxcolor.RGBA4444) {
		endiness.PutUint16(pbuf[:], uint16(p))
		if pmode == modeRGBA16 {
			if n := pixdata[ni]; n < pixMax {
				pixdata[ni]++
				pixdata = append(pixdata, pbuf[0], pbuf[1])
				return
			}
		}
		pmode = modeRGBA16
		pixdata = append(pixdata, 5, 1, pbuf[0], pbuf[1])
		ni = len(pixdata) - 3
	}
	addMask := func(op, cl byte) {
		md := modeMask + int(op)
		if pmode == md {
			if n := pixdata[ni]; n < pixMax {
				pixdata[ni]++
				pixdata = append(pixdata, cl)
				return
			}
		}
		pmode = md
		pixdata = append(pixdata, op, 1, cl)
		ni = len(pixdata) - 2
	}
	for y := 0; y < height; y++ {
		pmode = 0 // compression resets on each row
		for x := 0; x < width; x++ {
			i := y*rgba.Stride + x*4
			r := rgba.Pix[i+0]
			g := rgba.Pix[i+1]
			b := rgba.Pix[i+2]
			a := rgba.Pix[i+3]
			if mat != nil {
				if ind := mat.ColorIndexAt(x, y); ind != 0 {
					addMask((ind<<4)|4, r)
					continue
				}
			}
			if a == 0x00 {
				addZero()
			} else if a == 0xff {
				// RGB555
				addRGB15(noxcolor.RGB555Color(r, g, b))
			} else {
				// RGBA4444
				addRGBA16(noxcolor.RGBA4444Color(r, g, b, a))
			}
		}
	}
	pixdata = append(pixdata, 0, 0) // end
	return data[:17+len(pixdata)]
}
