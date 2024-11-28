package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/opennox/libs/maps"
)

func init() {
	cmd := &cobra.Command{
		Use:     "map command",
		Short:   "Tools for working with Nox maps",
		Aliases: []string{"m", "maps"},
	}
	Root.AddCommand(cmd)

	cmdCompress := &cobra.Command{
		Use:   "compress mapdir",
		Short: "Compresses a Nox/OpenNox map to ZIP archive",
	}
	cmd.AddCommand(cmdCompress)
	cmdCompressFormat := cmdCompress.Flags().StringP("format", "f", "zip", "format to use (only zip for now)")
	cmdCompressOut := cmdCompress.Flags().StringP("out", "o", "", "output file name")
	cmdCompress.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("one file path expected")
		}
		return cmdMapCompress(cmd, args[0], *cmdCompressOut, *cmdCompressFormat)
	}
}

func cmdMapCompress(cmd *cobra.Command, in, out string, format string) error {
	fi, err := os.Stat(in)
	if err != nil {
		return err
	}
	isDir := fi.IsDir()
	if out == "" {
		base := filepath.Base(in)
		if !isDir {
			base = strings.TrimSuffix(base, filepath.Ext(base))
		}
		out = base + "." + format
	}
	switch format {
	default:
		return fmt.Errorf("unsupported format: %s", format)
	case "nxz":
		if isDir {
			in = filepath.Join(in, filepath.Base(in)+maps.Ext)
		}
		// FIXME: support NXZ encoding
		return fmt.Errorf("NXZ encoding is not supported yet")
	case "zip":
		if !isDir {
			in = filepath.Dir(in)
		}
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()

		if err = maps.CompressMap(f, nil, in); err != nil {
			return err
		}
		return f.Close()
	}
}

func mapReadRawSections(fname string) ([]maps.RawSection, *maps.Header, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	rd, err := maps.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	hdr := rd.Map().Header()
	raw, err := rd.ReadSectionsRaw()
	return raw, &hdr, err
}

func mapWriteRawSections(fname string, hdr *maps.Header, raw []maps.RawSection) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	wr, err := maps.NewWriter(f, *hdr)
	if err != nil {
		return err
	}
	if err = wr.WriteRawSections(raw); err != nil {
		return err
	}
	if err = wr.Close(); err != nil {
		return err
	}
	return f.Close()
}

func mapRemoveSections(fname string, sections []string, verbose bool) error {
	return mapUpdateSections(fname, func(s *maps.RawSection) (bool, error) {
		if !slices.Contains(sections, s.Name) {
			return false, nil
		}
		*s = maps.RawSection{}
		return true, nil
	}, verbose)
}

func mapUpdateSections(fname string, fnc func(s *maps.RawSection) (bool, error), verbose bool) error {
	raw, hdr, err := mapReadRawSections(fname)
	if err != nil {
		return err
	}
	changed := 0
	out := make([]maps.RawSection, 0, len(raw))
	for _, s := range raw {
		upd, err := fnc(&s)
		if err != nil {
			return err
		}
		if upd {
			changed++
		}
		if !upd || s.Name != "" {
			out = append(out, s)
		}
	}
	if changed == 0 {
		if verbose {
			fmt.Println("no sections updated")
		}
		return nil
	}
	return mapWriteRawSections(fname, hdr, out)
}
