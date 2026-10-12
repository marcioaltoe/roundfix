// Suite: shared Go build cache Makefile wiring.
// Invariant: worktrees share a writable cache, retain a local fallback, and honor explicit GOCACHE.
// Boundary IN: the repository Makefile's test target with a compiled Go stub and isolated environment.
// Boundary OUT: real Go builds, cache trimming, and Worktree Build Cache reclamation.
package verifyselect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestMakeSharesOneGoBuildCache(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	uidCommand := exec.CommandContext(t.Context(), "id", "-u")
	uid, err := uidCommand.Output()
	if err != nil {
		t.Fatalf("id -u: %v", err)
	}
	cacheName := "roundfix-go-build-" + strings.TrimSpace(string(uid))
	goStub := testfixture.FixtureBinary(t, "go-stub", `package main
import (
	"fmt"
	"os"
)
func main() {
	if len(os.Args) < 2 || os.Args[1] != "test" {
		os.Exit(2)
	}
	if err := os.WriteFile(os.Getenv("ROUNDFIX_BUILD_CACHE_RECORD"), []byte(os.Getenv("GOCACHE")), 0600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`)
	var environment []string
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "GOCACHE", "TMPDIR", "MAKEFLAGS", "MFLAGS", "MAKEOVERRIDES", "ROUNDFIX_BUILD_CACHE_RECORD":
			continue
		}
		environment = append(environment, entry)
	}

	for _, name := range []string{"shared", "fallback", "explicit"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tmpRoot := t.TempDir()
			if name != "shared" {
				if err := os.Chmod(tmpRoot, 0555); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chmod(tmpRoot, 0700); err != nil {
						t.Error(err)
					}
				})
			}
			cases := []string{"default"}
			if name == "shared" {
				cases = []string{"first_worktree", "second_worktree"}
			} else if name == "fallback" {
				cases = []string{"directory_unwritable", "cache_unwritable"}
			} else if name == "explicit" {
				cases = []string{"environment", "command_line"}
			}
			for _, source := range cases {
				t.Run(source, func(t *testing.T) {
					if name == "fallback" && source == "cache_unwritable" {
						if err := os.Chmod(tmpRoot, 0700); err != nil {
							t.Fatal(err)
						}
						if err := os.Mkdir(filepath.Join(tmpRoot, cacheName), 0555); err != nil {
							t.Fatal(err)
						}
					}
					dir := t.TempDir()
					// CURDIR is physical even when macOS's temporary directory is a symlink.
					runDir, err := filepath.EvalSymlinks(dir)
					if err != nil {
						t.Fatal(err)
					}
					record := filepath.Join(dir, "cache-record")
					command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "test", "GO="+goStub)
					command.Dir = dir
					command.Env = append(append([]string{}, environment...), "TMPDIR="+tmpRoot+"/", "ROUNDFIX_BUILD_CACHE_RECORD="+record)
					want := filepath.Join(tmpRoot, cacheName)
					if name == "fallback" {
						want = filepath.Join(runDir, ".gocache")
					} else if name == "explicit" {
						want = filepath.Join(dir, "explicit cache")
						if source == "environment" {
							command.Env = append(command.Env, "GOCACHE="+want)
						} else {
							command.Env = append(command.Env, "GOCACHE="+filepath.Join(dir, "overridden"))
							command.Args = append(command.Args, "GOCACHE="+want)
						}
					}
					output, err := command.CombinedOutput()
					if err != nil {
						t.Fatalf("make test: %v\n%s", err, output)
					}
					got, err := os.ReadFile(record)
					if err != nil {
						t.Fatal(err)
					}
					if string(got) != want {
						t.Fatalf("Go received GOCACHE=%q, want %q\n%s", got, want, output)
					}
					if name == "shared" {
						entries, err := os.ReadDir(want)
						if err != nil {
							t.Fatalf("shared cache directory: %v", err)
						}
						if len(entries) != 0 {
							t.Fatalf("shared cache retains write probes: %v", entries)
						}
					}
				})
			}
		})
	}
}
