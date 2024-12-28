package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/opennox/libs/noxsave"
	"github.com/opennox/libs/player"
)

func init() {
	cmd := &cobra.Command{
		Use:   "save command",
		Short: "Tools for working with Nox save files",
	}
	Root.AddCommand(cmd)

	cmdList := &cobra.Command{
		Use:     "list dir",
		Short:   "List information about save files",
		Aliases: []string{"l", "ls"},
	}
	cmd.AddCommand(cmdList)
	cmdList.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = []string{"."}
		} else if len(args) != 1 {
			return errors.New("one dir path expected")
		}
		return cmdSaveList(args[0])
	}

	cmdInfo := &cobra.Command{
		Use:     "info path",
		Short:   "Print information about a save file",
		Aliases: []string{"i"},
	}
	cmd.AddCommand(cmdInfo)
	cmdInfo.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("one file path expected")
		}
		return cmdSaveInfo(args[0], false)
	}
}

func cmdSaveList(dir string) error {
	if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
		return cmdSaveInfo(dir, true)
	}
	list, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var last error
	for _, fi := range list {
		lname := strings.ToLower(fi.Name())
		path := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			if lname == "save" {
				if err := cmdSaveList(path); err != nil {
					return err
				}
				continue
			}
			list2, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, fi := range list2 {
				lname := strings.ToLower(fi.Name())
				if !strings.HasSuffix(lname, ".plr") {
					continue
				}
				if err := cmdSaveInfo(filepath.Join(path, fi.Name()), true); err != nil {
					fmt.Fprintln(os.Stderr, path, err)
					last = err
				}
			}
		} else {
			if !strings.HasSuffix(lname, ".plr") {
				continue
			}
			if err := cmdSaveInfo(path, true); err != nil {
				fmt.Fprintln(os.Stderr, path, err)
				last = err
			}
		}
	}
	return last
}

func cmdSaveInfo(path string, short bool) error {
	fmt.Println(path)
	defer fmt.Println()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := noxsave.NewReader(f)
	if err != nil {
		return err
	}
	raw, err := r.ReadRawSections()
	if err != nil {
		return err
	}
	out := make([]noxsave.Section, len(raw))
	for i, s := range raw {
		out[i] = s

		sect, err := s.DecodeWith(nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "\terror:", err)
			continue
		}
		out[i] = sect
	}
	var (
		info *noxsave.FileInfo
	)
	for _, s := range out {
		switch s := s.(type) {
		case *noxsave.FileInfo:
			info = s
		}
	}
	if short {
		if info != nil {
			p := &info.Player
			fmt.Printf("\t%s - %s (%v) - %q; %d sections\n",
				info.Time.Time().Format(time.DateTime),
				p.Name, player.Class(p.Class),
				info.MapName,
				len(raw),
			)
		} else {
			fmt.Printf("\tno info; %d sections\n", len(raw))
		}
		return nil
	}
	if info != nil {
		p := &info.Player
		skin := p.Skin.HexString()
		mustache := p.Mustache.HexString()
		goatee := p.Goatee.HexString()
		beard := p.Beard.HexString()
		if mustache == skin {
			mustache = "-"
		}
		if goatee == skin {
			goatee = "-"
		}
		if beard == skin {
			beard = "-"
		}
		fmt.Printf(
			"\tTime: %v\n"+
				"\tPlayer: %q (%v)\n"+
				"\t\tSkin: %s  Hair: %s\n"+
				"\t\tMustache: %s  Goatee: %s  Beard: %s\n"+
				"\t\tPants: %d  Shirt: %d,%d  Shoes: %d,%d\n",
			info.Time.Time().Format(time.DateTime),
			p.Name, player.Class(p.Class),
			p.Skin.HexString(),
			p.Hair.HexString(),
			mustache,
			goatee,
			beard,
			p.Pants,
			p.Shirt1, p.Shirt2,
			p.Shoes1, p.Shoes2,
		)
		if info.MapName != "" {
			fmt.Printf("\tMap name: %q\n", info.MapName)
		}
		if info.Path != "" {
			fmt.Printf("\tOrig path: %q\n", info.Path)
		}
		fmt.Println()
	}
	for i, s := range raw {
		sect := out[i]
		fmt.Printf("\t%02d (%v) - %d bytes", s.ID, s.ID, len(s.Data))
		switch sect := sect.(type) {
		case *noxsave.RawSection:
			// skip
		default:
			fmt.Printf("\n\t\t%+v", sect)
		}
		fmt.Println()
	}
	return nil
}
