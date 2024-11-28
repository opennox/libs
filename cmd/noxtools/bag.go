package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/opennox/opennox-lib/bag"
	"github.com/opennox/opennox-lib/noximage/pcx"
	"github.com/opennox/opennox-lib/things"
)

func init() {
	cmd := &cobra.Command{
		Use:     "videobag command",
		Short:   "Tools related to Nox video.bag and video.idx files",
		Aliases: []string{"bag"},
	}
	Root.AddCommand(cmd)
	fBag := cmd.PersistentFlags().StringP("bag", "b", "video.bag", "path to video.bag file")
	fIdx := cmd.PersistentFlags().StringP("idx", "i", "", "path to video.idx file (empty means auto)")

	cmdJSON := &cobra.Command{
		Use:     "idx2json [--bag video.bag] [--idx video.idx] [--out ./out.json]",
		Short:   "Reads Nox video.idx and dumps metadata stored there as JSON",
		Aliases: []string{"i2j"},
	}
	cmd.AddCommand(cmdJSON)
	{
		fOut := cmdJSON.Flags().StringP("out", "o", "video.idx.json", "output path for images or archive")
		cmdJSON.RunE = func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			f, err := bag.OpenWithIndex(*fBag, *fIdx)
			if err != nil {
				return err
			}
			defer f.Close()

			list, err := f.Segments()
			if err != nil {
				return err
			}

			w, err := os.Create(*fOut)
			if err != nil {
				return err
			}
			defer w.Close()

			enc := json.NewEncoder(w)
			enc.SetIndent("", "\t")
			if err = enc.Encode(list); err != nil {
				return err
			}
			return w.Close()
		}
	}

	cmdExtract := &cobra.Command{
		Use:     "extract [--bag video.bag] [--idx video.idx] [--out ./out] [file ...]",
		Short:   "Extracts images from Nox video.bag file",
		Aliases: []string{"e"},
	}
	cmd.AddCommand(cmdExtract)
	{
		fOut := cmdExtract.Flags().StringP("out", "o", "", "output path for images or archive")
		fJSON := cmdExtract.Flags().BoolP("json", "j", false, "write additional image metadata as JSON")
		fZIP := cmdExtract.Flags().BoolP("zip", "z", false, "write files into a ZIP archive")
		fMode := cmdExtract.Flags().StringP("mode", "m", "", "extraction mode (imgs, objs)")
		fAll := cmdExtract.Flags().BoolP("all", "a", false, "include unclassified images")
		cmdExtract.RunE = func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			f, err := bag.OpenWithIndex(*fBag, *fIdx)
			if err != nil {
				return err
			}
			defer f.Close()

			e := newExtractor(*fOut, f)
			defer e.Close()
			e.json = *fJSON

			switch ext := filepath.Ext(*fOut); ext {
			case ".zip", ".gz":
				if err := e.Compress(ext); err != nil {
					return err
				}
			default:
				if *fZIP {
					if err := e.Compress(".zip"); err != nil {
						return err
					}
				}
			}
			switch *fMode {
			case "", "obj", "objs":
				fname := filepath.Join(filepath.Dir(*fBag), "thing.bin")
				r, err := things.Open(fname)
				if err != nil {
					return err
				}
				defer r.Close()
				if err := e.ExtractObjects(r, *fAll, args...); err != nil {
					return err
				}
			case "img", "imgs":
				if err := e.ExtractImages(args...); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported mode: %q", *fMode)
			}
			return e.Close()
		}
	}
	cmdDecompress := &cobra.Command{
		Use:   "decompress [--bag video.bag] [--idx video.idx]",
		Short: "Decompress Nox video.bag file and update video.idx accordingly",
	}
	cmd.AddCommand(cmdDecompress)
	{
		cmdDecompress.RunE = func(cmd *cobra.Command, args []string) error {
			return bag.Decompress(*fBag, *fIdx)
		}
	}
	cmdReplace := &cobra.Command{
		Use:   "replace [--bag video.bag] [--idx video.idx] ",
		Short: "Replace one or more sprites in the Nox video.bag file",
	}
	cmd.AddCommand(cmdReplace)
	{
		fSprites := cmdReplace.Flags().StringSliceP("sprite", "s", nil, "sprite in the format: index1=image1.png")
		cmdReplace.RunE = func(cmd *cobra.Command, args []string) error {
			var list []bag.Replacement
			for _, s := range *fSprites {
				sub := strings.SplitN(s, "=", 3)
				if len(sub) != 2 {
					return fmt.Errorf("sprites should be in the format: index1=image1.png; got: %q", s)
				}
				ind, err := strconv.Atoi(sub[0])
				if err != nil {
					return err
				}
				f, err := os.Open(sub[1])
				if err != nil {
					return err
				}
				img, _, err := image.Decode(f)
				_ = f.Close()
				if err != nil {
					return err
				}
				list = append(list, bag.Replacement{
					ImageInd: ind,
					Image:    img,
					Point:    nil, // TODO: allow overriding the point
				})
			}
			if len(list) == 0 {
				return errors.New("no sprites to replace")
			}
			cmd.SilenceUsage = true
			return bag.ReplaceSprites(*fBag, *fIdx, list)
		}
	}
}

func newExtractor(out string, b *bag.File) *extractor {
	return &extractor{b: b, out: out}
}

type extractor struct {
	b       *bag.File
	out     string
	json    bool
	closers []func() error
	zw      *zip.Writer
	tw      *tar.Writer
	buf     bytes.Buffer
	seen    map[string]int
	used    map[int]struct{}
}

func (e *extractor) Close() error {
	var last error
	for i := len(e.closers) - 1; i >= 0; i-- {
		if err := e.closers[i](); err != nil {
			last = err
		}
	}
	e.closers = nil
	return last
}

func (e *extractor) writeFile(name string, r io.Reader, sz int64) error {
	switch {
	case e.zw != nil:
		w, err := e.zw.Create(name)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, r)
		return err
	case e.tw != nil:
		err := e.tw.WriteHeader(&tar.Header{
			Name: name,
			Size: sz,
		})
		if err != nil {
			return err
		}
		_, err = io.Copy(e.tw, r)
		return err
	default:
		f, err := os.Create(filepath.Join(e.out, name))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
		return f.Close()
	}
}

func (e *extractor) Compress(ext string) error {
	switch ext {
	case "zip", ".zip":
		zpath := e.out
		if filepath.Ext(zpath) != ".zip" {
			zpath = filepath.Join(zpath, "video.bag.zip")
		}
		zf, err := os.Create(zpath)
		if err != nil {
			return err
		}
		e.closers = append(e.closers, zf.Close)

		e.zw = zip.NewWriter(zf)
		e.closers = append(e.closers, e.zw.Close)
		return nil
	case "gz", ".gz", ".tar.gz":
		zpath := e.out
		if filepath.Ext(zpath) != ".gz" {
			zpath = filepath.Join(zpath, "video.tar.gz")
		}
		zf, err := os.Create(zpath)
		if err != nil {
			return err
		}
		e.closers = append(e.closers, zf.Close)
		zw := gzip.NewWriter(zf)
		e.closers = append(e.closers, zw.Close)
		e.tw = tar.NewWriter(zw)
		e.closers = append(e.closers, e.tw.Close)
		return nil
	default:
		return fmt.Errorf("unsupported compression: %q", ext)
	}
}

func (e *extractor) ExtractImages(names ...string) error {
	e.seen = make(map[string]int)
	if len(names) == 0 {
		imgs, err := e.b.Images()
		if err != nil {
			return err
		}
		var last error
		for _, img := range imgs {
			if err := e.processImageRec("", img); err != nil {
				last = err
				log.Println(img.Name, last)
			}
		}
		return last
	}
	var last error
	for _, name := range names {
		img, err := e.b.ImageByName(name)
		if err != nil {
			last = err
			log.Println(name, last)
			continue
		} else if img == nil {
			last = errors.New("not found")
			log.Println(name, last)
			continue
		}
		if err = e.processImageRec("", img); err != nil {
			last = err
			log.Println(name, last)
		}
	}
	return last
}

func (e *extractor) processImageRec(base string, img *bag.ImageRec) error {
	im, err := img.Decode()
	if err != nil {
		return err
	}
	if base == "" {
		base = strings.TrimSuffix(img.Name, path.Ext(img.Name))
	}
	key := strings.ToLower(base)
	if n := e.seen[key]; n > 0 {
		base += "_" + strconv.Itoa(n)
	}
	e.seen[key]++
	if e.json {
		e.buf.Reset()
		enc := json.NewEncoder(&e.buf)
		enc.SetIndent("", "\t")
		if err := enc.Encode(im.ImageMeta); err != nil {
			return err
		}
		if err := e.writeFile(base+".json", &e.buf, int64(e.buf.Len())); err != nil {
			return err
		}
	}
	e.buf.Reset()
	err = png.Encode(&e.buf, im.Image)
	if err != nil {
		return err
	}
	if err := e.writeFile(base+".png", &e.buf, int64(e.buf.Len())); err != nil {
		return err
	}
	if im.Material != nil {
		e.buf.Reset()
		err = png.Encode(&e.buf, im.Material)
		if err != nil {
			return err
		}
		if err := e.writeFile(base+"_mat.png", &e.buf, int64(e.buf.Len())); err != nil {
			return err
		}
	}
	return nil
}

func (e *extractor) processFrames(name string, frames []things.ImageRef, row int) (image.Point, error) {
	if row <= 0 {
		row = len(frames)
	}
	imgByInd, err := e.b.Images()
	if err != nil {
		return image.Point{}, err
	}
	var (
		last     error
		rect     image.Rectangle
		imgs     = make([]*pcx.Image, 0, len(frames))
		material = 0
		good     = 0
	)
	for _, ref := range frames {
		if ref.Ind < 0 {
			return image.Point{}, fmt.Errorf("external image refs not supported")
		}
		e.used[ref.Ind] = struct{}{}
		img, err := imgByInd[ref.Ind].Decode()
		if err != nil {
			last = err
			imgs = append(imgs, &pcx.Image{})
			continue
		}
		if img.Material != nil {
			material++
		}
		good++
		imgs = append(imgs, img)
		rect = rect.Union(img.Bounds().Add(img.Point))
	}

	base := name
	key := strings.ToLower(base)
	if n := e.seen[key]; n > 0 {
		base += "_" + strconv.Itoa(n)
	}
	e.seen[key]++

	sz := rect.Size()
	cols := len(frames) / row
	if len(frames)%row != 0 {
		cols++
	}
	fullRect := image.Rect(0, 0, sz.X*row, sz.Y*cols)
	full := image.NewRGBA(fullRect)
	for i, img := range imgs {
		if img.Image == nil {
			continue
		}
		xi, yi := i%row, i/row
		zp := image.Point{X: xi * sz.X, Y: yi * sz.Y}.Add(img.Point).Sub(rect.Min)
		draw.Draw(full, image.Rectangle{Min: zp, Max: zp.Add(img.Bounds().Size())}, img.Image, image.Point{}, draw.Src)
	}

	e.buf.Reset()
	err = png.Encode(&e.buf, full)
	if err != nil {
		return image.Point{}, err
	}
	if err := e.writeFile(base+".png", &e.buf, int64(e.buf.Len())); err != nil {
		return image.Point{}, err
	}

	if material == good {
		mat := image.NewPaletted(fullRect, imgs[0].Material.Palette)
		for i, img := range imgs {
			if img.Material == nil {
				continue
			}
			xi, yi := i%row, i/row
			zp := image.Point{X: xi * sz.X, Y: yi * sz.Y}.Add(img.Point).Sub(rect.Min)
			draw.Draw(mat, image.Rectangle{Min: zp, Max: zp.Add(img.Bounds().Size())}, img.Material, image.Point{}, draw.Src)
		}

		e.buf.Reset()
		err = png.Encode(&e.buf, mat)
		if err != nil {
			return image.Point{}, err
		}
		if err := e.writeFile(base+"_mat.png", &e.buf, int64(e.buf.Len())); err != nil {
			return image.Point{}, err
		}
	}
	return rect.Min, last
}

func (e *extractor) processImage(name string, ref *things.ImageRef) error {
	if ref.Ind < 0 {
		return fmt.Errorf("external image refs not supported")
	}
	e.used[ref.Ind] = struct{}{}
	imgByInd, err := e.b.Images()
	if err != nil {
		return err
	}
	return e.processImageRec(name, imgByInd[ref.Ind])
}

func (e *extractor) ExtractObjects(thg *things.Reader, all bool, names ...string) error {
	d, err := thg.ReadAll()
	if err != nil {
		return err
	}
	var want map[string]bool
	if len(names) != 0 {
		want = make(map[string]bool)
		for _, name := range names {
			want[name] = true
		}
	}
	var last error
	e.seen = make(map[string]int)
	e.used = make(map[int]struct{})
	for _, v := range d.Images {
		if v.Name == "" || (want != nil && !want[v.Name]) {
			continue
		}
		const pref = "images/"
		log.Println(v.Name)
		switch {
		case v.Ani != nil:
			_, err := e.processFrames(pref+v.Name, v.Ani.Frames, 0)
			if err != nil {
				last = err
				log.Println(v.Name, last)
				continue
			}
		case v.Img != nil:
			if err := e.processImage(pref+v.Name, v.Img); err != nil {
				last = err
				log.Println(v.Name, last)
				continue
			}
		}
	}
	for _, th := range d.Things {
		if th.Name == "" || (want != nil && !want[th.Name]) {
			continue
		}
		const pref = "things/"
		log.Println(th.Name)
		switch dr := th.Draw.(type) {
		case things.BaseDraw:
			if err := e.processImage(pref+th.Name, &dr.Img); err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.StaticDraw:
			if err := e.processImage(pref+th.Name, &dr.Img); err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.WeaponDraw:
			if err := e.processImage(pref+th.Name, &dr.Img); err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.ArmorDraw:
			if err := e.processImage(pref+th.Name, &dr.Img); err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.StaticRandomDraw:
			_, err := e.processFrames(pref+th.Name, dr.Imgs, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.DoorDraw:
			_, err := e.processFrames(pref+th.Name, dr.Imgs, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.AnimateDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.GlyphDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.WeaponAnimateDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.ArmorAnimateDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.FlagDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.SphericalShieldDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.SummonEffectDraw:
			_, err := e.processFrames(pref+th.Name, dr.Anim.Frames, 0)
			if err != nil {
				last = err
				log.Println(th.Name, last)
				continue
			}
		case things.PlayerDraw:
			anims := maps.Keys(dr.Anims)
			slices.Sort(anims)
			for _, aname := range anims {
				a := dr.Anims[aname]
				parts := maps.Keys(a.Parts)
				slices.Sort(parts)
				for _, pname := range parts {
					p := a.Parts[pname]
					frames := make([]things.ImageRef, 0, len(p)*int(a.FramesPerDir))
					for _, dir := range p {
						frames = append(frames, dir...)
					}
					_, err := e.processFrames(fmt.Sprintf(pref+"%s/%s/%s", th.Name, pname, aname), frames, int(a.FramesPerDir))
					if err != nil {
						last = err
						log.Println(th.Name, last)
						continue
					}
				}
			}
		case things.MonsterDraw:
			for _, v := range dr.Anims {
				frames := make([]things.ImageRef, 0, int(v.FramesPerDir))
				for _, dir := range v.Frames {
					frames = append(frames, dir...)
				}
				_, err := e.processFrames(fmt.Sprintf(pref+"%s/%s", th.Name, v.Type), frames, int(v.FramesPerDir))
				if err != nil {
					last = err
					log.Println(th.Name, last)
					continue
				}
			}
		case things.MaidenDraw:
			for _, v := range dr.Anims {
				frames := make([]things.ImageRef, 0, int(v.FramesPerDir))
				for _, dir := range v.Frames {
					frames = append(frames, dir...)
				}
				_, err := e.processFrames(fmt.Sprintf(pref+"%s/%s", th.Name, v.Type), frames, int(v.FramesPerDir))
				if err != nil {
					last = err
					log.Println(th.Name, last)
					continue
				}
			}
		case things.ConditionalAnimateDraw:
			for i, v := range dr.Anims {
				_, err := e.processFrames(fmt.Sprintf(pref+"%s/%d", th.Name, i), v.Frames, 0)
				if err != nil {
					last = err
					log.Println(th.Name, last)
					continue
				}
			}
		case things.MonsterGeneratorDraw:
			for i, v := range dr.Anims {
				var name string
				switch i {
				case 0:
					name = "NORMAL"
				case 1:
					name = "DAMAGE_1"
				case 2:
					name = "DAMAGE_2"
				case 3:
					name = "DEAD"
				case 4:
					name = "MONSTER"
				default:
					name = strconv.Itoa(i)
				}
				_, err := e.processFrames(fmt.Sprintf(pref+"%s/%s", th.Name, name), v.Frames, 0)
				if err != nil {
					last = err
					log.Println(th.Name, last)
					continue
				}
			}
		}
	}
	if want == nil && all {
		imgs, err := e.b.Images()
		if err != nil {
			return err
		}
		for i, img := range imgs {
			if _, ok := e.used[i]; ok {
				continue
			}
			base := strings.TrimSuffix(img.Name, path.Ext(img.Name))
			if err := e.processImageRec(fmt.Sprintf("unclassified/%d_%s", i, base), img); err != nil {
				last = err
				log.Println(img.Name, last)
			}
		}
	}
	return last
}
