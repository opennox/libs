package maps

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/noxtest"
)

func TestIsAllowedFile(t *testing.T) {
	var cases = []struct {
		path string
		exp  bool
	}{
		{"example.map", true},
		{"example.nxz", false},
		{"example.rul", true},
		{"example.zip", false},
		{"example.tar", false},
		{"example.tar.gz", false},
		{"user.rul", false},
		{"some.lua", true},
		{"go.mod", true},
		{"go.sum", true},
		{"some.go", true},
		{"sub/some.go", true},
		{"vendor/sub/some.go", true},
		{"LICENSE", true},
		{"README.md", true},
		{"README.txt", true},
		{"some.other", false},
		{"some.json", true},
		{"some.yaml", true},
		{"some.yml", true},
		{"some.png", true},
		{"some.jpg", true},
		{"some.mp3", true},
		{"some.ogg", true},
		{".git/config", false},
		{".git/refs/heads/fake.go", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.path, func(t *testing.T) {
			got := IsAllowedFile(c.path)
			if got != c.exp {
				t.FailNow()
			}
		})
	}
}

func copyFile(t testing.TB, dst, src string) {
	s, err := ifs.Open(src)
	must.NoError(t, err)
	defer s.Close()
	d, err := os.Create(dst)
	must.NoError(t, err)
	defer d.Close()
	_, err = io.Copy(d, s)
	must.NoError(t, err)
	err = d.Close()
	must.NoError(t, err)
}

func TestMapServer(t *testing.T) {
	dpath := noxtest.DataPath(t, "maps")
	srcdir, err := os.MkdirTemp("", "opennox-map-test-src-*")
	must.NoError(t, err)
	dstdir, err := os.MkdirTemp("", "opennox-map-test-dst-*")
	must.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(srcdir)
		_ = os.RemoveAll(dstdir)
	})
	const mname = "estate"
	for _, name := range []string{
		"estate", "!estate",
		"Estate", "!Estate",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			// Copy a well-known map to a source folder.
			err = os.MkdirAll(filepath.Join(srcdir, name), 0755)
			must.NoError(t, err)
			smpath := filepath.Join(srcdir, name, name+".map")
			copyFile(t, smpath, filepath.Join(dpath, mname, mname+".map"))

			// Start serving maps from source folder.
			srv := NewServer(slog.Default(), srcdir)
			hsrv := &http.Server{Handler: srv}
			l, err := net.Listen("tcp", ":0")
			must.NoError(t, err)
			t.Cleanup(func() {
				_ = l.Close()
			})
			t.Logf("listening on %q", l.Addr())
			go func() {
				err := hsrv.Serve(l)
				if err == http.ErrServerClosed {
					return
				}
				must.NoError(t, err)
			}()
			t.Cleanup(func() {
				_ = hsrv.Close()
			})

			ctx := context.Background()
			cli, err := NewClient(ctx, slog.Default(), l.Addr().String())
			must.NoError(t, err)
			defer cli.Close()

			err = cli.DownloadMap(ctx, dstdir, name)
			must.NoError(t, err)
		})
	}
}
