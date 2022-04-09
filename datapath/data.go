// Package datapath implements automatic detection of Nox game data directory.
package datapath

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/noxworld-dev/opennox-lib/ifs"
)

var datadir struct {
	sync.Once
	path string
}

// getData returns the current Nox data dir.
func getData() string {
	datadir.Do(func() {
		if datadir.path == "" {
			SetData(FindData())
		}
	})
	return datadir.path
}

// Data returns the current Nox data dir. If additional args are provided, they will be joined with the data dir.
// If no data directory was set, it will try to locate it automatically using FindData.
func Data(path ...string) string {
	if len(path) == 0 {
		return getData()
	}
	args := make([]string, 0, 1+len(path))
	args = append(args, getData())
	args = append(args, path...)
	return filepath.Join(args...)
}

// SetData the Nox data dir.
func SetData(dir string) {
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	datadir.path = dir
	Log.Printf("setting data dir to: %q", dir)
}

// FindData locates Nox game data path. It returns empty string if not found.
// It does not affect data path returned by Data.
func FindData() string {
	consider := []string{
		os.Getenv("NOX_DATA"),    // takes priority
		".",                      // current dir overrides registry and other install paths
		filepath.Dir(os.Args[0]), // same for binary dir
	}
	if runtime.GOOS != "windows" {
		// search for Nox directory in all the prefixes
		// useful for XDG paths on Linux
		consider = append(consider, tryWithPrefixes("Nox")...)
	}
	// search in registry by default
	consider = append(consider, registryPaths()...)
	// prefer GoG, since it's patched and official
	consider = append(consider, gogPaths()...)
	// then try Reloaded: patched, though unofficial
	consider = append(consider, reloadedPaths()...)
	// lastly, check Origin
	consider = append(consider, originPaths()...)
	for _, path := range consider {
		if path == "" {
			continue
		}
		if !filepath.IsAbs(path) {
			// this is a workaround for Nox trying to chdir from time to time
			path = filepath.Join(workdir, path)
		}
		if CheckData(path) {
			return path
		}
	}
	return ""
}

var checkFiles = []string{
	"gamedata.bin",
	"modifier.bin",
	"monster.bin",
	"thing.bin",
}

// CheckData checks if a directory contains Nox game data.
func CheckData(path string) bool {
	if fi, err := ifs.Stat(path); err != nil || !fi.IsDir() {
		return false
	}
	for _, name := range checkFiles {
		fpath := filepath.Join(path, name)
		if fi, err := ifs.Stat(fpath); err != nil || fi.IsDir() {
			Log.Printf("cannot find required data file %q in %q", name, path)
			return false
		}
	}
	return true
}
