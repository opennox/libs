package datapath

import "github.com/noxworld-dev/opennox-lib/log"

var reloadedLog = log.New("reloaded")

// reloadedPaths returns possible path of Nox Reloaded installation.
func reloadedPaths() []string {
	out := tryWithPrefixes(
		`Program Files (x86)\Nox`,
		`Program Files\Nox`,
	)
	for _, path := range out {
		reloadedLog.Printf("potential install path: %q", path)
	}
	return out
}
