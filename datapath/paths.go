package datapath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/noxworld-dev/opennox-lib/ifs"
	"github.com/noxworld-dev/opennox-lib/log"
)

var Log = log.New("path")

var workdir string

func init() {
	if wd, err := os.Getwd(); err != nil {
		Log.Printf("cannot get workdir: %v", err)
	} else {
		workdir = wd
	}
}

func cleanPath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, `/`, `\`)
	}
	return strings.ReplaceAll(path, `\`, `/`)
}

func tryWithPrefixes(paths ...string) []string {
	var out []string
	for _, pref := range pathPrefixes() {
		if pref == "" {
			continue
		}
		for _, path := range paths {
			fpath := filepath.Join(pref, path)
			fpath = cleanPath(fpath)
			if _, err := ifs.Stat(fpath); err == nil {
				out = append(out, fpath)
			}
		}
	}
	return out
}

func pathPrefixes() []string {
	if runtime.GOOS == "windows" {
		return []string{`C:\`, `D:\`}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	// Linux Snapcraft installation replaces HOME variable
	if rhome := os.Getenv("SNAP_REAL_HOME"); rhome != "" {
		home = rhome
	}
	return []string{
		// Linux Snapcraft installation user common dir
		os.Getenv("SNAP_USER_COMMON"),
		// XDG environment variable
		os.Getenv("XDG_DATA_HOME"),
		// XDG default directory
		filepath.Join(home, ".local/share"),
		// Wine default paths
		filepath.Join(home, ".wine/drive_c"),
		filepath.Join(home, ".wine/drive_d"),
		// legacy Wine paths
		filepath.Join(home, ".wine/dosdevices/c:"),
		filepath.Join(home, ".wine/dosdevices/d:"),
		// TODO: these are probably from Lutris, but not sure if they are the default ones
		filepath.Join(home, "Games/gog/nox/drive_c"),
		filepath.Join(home, "Games/gog/nox/drive_d"),
		// fallback to home as a last resort
		home,
	}
}
