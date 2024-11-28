package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opennox/noxscript/ns/asm"
	nscmd "github.com/opennox/noxscript/ns/v3/cmd"
	"github.com/opennox/noxscript/ns/v3/noxast"

	"github.com/opennox/opennox-lib/maps"
)

func init() {
	cmd := &cobra.Command{
		Use:     "noxscript command",
		Short:   "Tools for working with NoxScript",
		Aliases: []string{"ns"},
	}
	Root.AddCommand(cmd)

	cmdExtr := &cobra.Command{
		Use:   "extract input.map output.obj",
		Short: "Extract binary NoxScript file from a map",
	}
	cmd.AddCommand(cmdExtr)
	cmdExtr.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSExtract(cmd, args)
	}

	cmdRemove := &cobra.Command{
		Use:   "remove input.map",
		Short: "Removes NoxScript from a map file",
	}
	cmd.AddCommand(cmdRemove)
	cmdRemove.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSRemove(cmd, args)
	}

	cmdIns := &cobra.Command{
		Use:   "insert input.map input.obj",
		Short: "Insert or replace binary NoxScript file in a map",
	}
	cmd.AddCommand(cmdIns)
	cmdIns.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSInsert(cmd, args)
	}

	cmdDis := &cobra.Command{
		Use:   "disasm input.obj [output.txt]",
		Short: "Disassemble binary NoxScript file or map script into text assembly",
	}
	cmd.AddCommand(cmdDis)
	cmdDis.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdNSDisasm(cmd, args)
	}

	cmdDecomp := &cobra.Command{
		Use:   "decomp input.obj [output.go]",
		Short: "Decompile binary NoxScript file or map script into human-readable script",
	}
	cmd.AddCommand(cmdDecomp)
	cmdDecompNoFold := cmdDecomp.Flags().Bool("nofold", false, "do not fold the code")
	cmdDecompNoOpt := cmdDecomp.Flags().Bool("noopt", false, "do not optimize the code")
	cmdDecompPkg := cmdDecomp.Flags().String("pkg", "", "package name for the script")
	cmdDecomp.RunE = func(cmd *cobra.Command, args []string) error {
		return nscmd.Decomp(args, &noxast.Config{
			Package:       *cmdDecompPkg,
			DoNotOptimize: *cmdDecompNoOpt,
			DoNotFold:     *cmdDecompNoFold,
		})
	}
}

func cmdNSExtract(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("at least one map file expected")
	case 1, 2:
		fname := args[0]
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer f.Close()

		s, err := maps.ReadScript(f)
		if err != nil {
			return err
		} else if s == nil || len(s.Data) == 0 {
			err := errors.New("no script in map file")
			if len(args) == 2 {
				return err
			}
			log.Println(err)
			return nil
		}

		out := strings.TrimSuffix(fname, filepath.Ext(fname)) + ".obj"
		if len(args) == 2 {
			out = args[1]
		}
		log.Printf("writing %d bytes to %s\n", len(s.Data), out)
		return os.WriteFile(out, s.Data, 0644)
	default:
		dstDir := ""
		hasDst := false
		if last := args[len(args)-1]; filepath.Ext(last) == "" {
			hasDst, dstDir = true, last
			args = args[:len(args)-1]
		}
		var last error
		for _, fname := range args {
			err := func() error {
				f, err := os.Open(fname)
				if err != nil {
					return err
				}
				defer f.Close()

				s, err := maps.ReadScript(f)
				if err != nil {
					return err
				} else if s == nil || len(s.Data) == 0 {
					return nil
				}
				name := filepath.Base(fname)
				name = strings.TrimSuffix(name, filepath.Ext(name)) + ".obj"
				var out string
				if hasDst {
					out = filepath.Join(dstDir, name)
				} else {
					out = filepath.Join(filepath.Dir(fname), name)
				}
				return os.WriteFile(out, s.Data, 0644)
			}()
			if err != nil {
				log.Printf("%s: %v\n", fname, err)
				last = err
			} else {
				log.Println(fname)
			}
		}
		return last
	}
}

func cmdNSRemove(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("at least one map file expected")
	}
	sname := (&maps.Script{}).MapSection()
	sdname := (&maps.ScriptData{}).MapSection()
	var last error
	for _, fname := range args {
		if err := mapUpdateSections(fname, func(s *maps.RawSection) (bool, error) {
			switch s.Name {
			case sname:
				*s = maps.RawSection{} // remove
				return true, nil
			case sdname:
				if len(s.Data) == 0 {
					return false, nil
				}
				data, err := (&maps.ScriptData{}).MarshalBinary()
				if err != nil {
					return false, err
				}
				s.Data = data
				return true, nil
			default:
				return false, nil
			}
		}, true); err != nil {
			last = err
			if len(args) > 1 {
				log.Println(err)
			}
		}
	}
	return last
}

func cmdNSInsert(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("expected map file and a binary script file")
	}
	fmap, fscr := args[0], args[1]
	raw, hdr, err := mapReadRawSections(fmap)
	if err != nil {
		return err
	}
	sobj, err := os.ReadFile(fscr)
	if err != nil {
		return err
	}
	scr := &maps.Script{Data: sobj}
	sdata, err := scr.MarshalBinary()
	if err != nil {
		return err
	}
	found := false
	for i, s := range raw {
		if s.Name == scr.MapSection() {
			raw[i].Data = sdata
			found = true
			break
		}
	}
	if !found {
		raw = append(raw, maps.RawSection{
			Name: scr.MapSection(),
			Data: sdata,
		})
		maps.SortRawSections(raw)
	}
	return mapWriteRawSections(fmap, hdr, raw)
}

func cmdNSDisasm(cmd *cobra.Command, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return errors.New("expected one or two argument")
	}
	fname := args[0]
	var out io.Writer = os.Stdout
	if len(args) == 2 {
		f, err := os.Create(args[1])
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var raw []byte
	if filepath.Ext(fname) == ".map" {
		s, err := maps.ReadScript(f)
		if err != nil {
			return err
		}
		raw = s.Data
	} else {
		raw, err = io.ReadAll(f)
		if err != nil {
			return err
		}
	}
	_ = f.Close()

	scr, err := asm.ReadScript(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	if len(scr.Strings) != 0 {
		fmt.Fprintln(out, "STRINGS:")
		for i, s := range scr.Strings {
			fmt.Fprintf(out, "\t%d: %q\n", i, s)
		}
		fmt.Fprintln(out)
	}
	var last error
	for i, fnc := range scr.Funcs {
		fmt.Fprintf(out, "func %d: %q\n", i, fnc.Name)
		fmt.Fprintf(out, "\targs: %d, locals: %d, returns: %d\n",
			fnc.Args, len(fnc.Vars)-fnc.Args, fnc.Return)
		fmt.Fprintln(out)

		code, err := asm.Decode(fnc.Code)
		if err != nil {
			err = fmt.Errorf("cannot disasm %q: %w", fnc.Name, err)
			log.Println(err)
			last = err
			continue
		}
		_ = asm.Print(out, code)
		fmt.Fprintln(out)
	}
	return last
}
