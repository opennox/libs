package datapath

import "github.com/noxworld-dev/opennox-lib/log"

var originLog = log.New("origin")

// originPaths returns possible path of Origin Nox installation.
func originPaths() []string {
	out := tryWithPrefixes(
		`Program Files (x86)\Origin Games\Nox`,
		`Program Files\Origin Games\Nox`,
	)
	for _, path := range out {
		originLog.Printf("potential install path: %q", path)
	}
	return out
}
