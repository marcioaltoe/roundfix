// Boundary: immutable default-branch config read and actual command execution
// in temporary local repositories. No GitHub or network remote is contacted.
package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch(t *testing.T) {
	for _, original := range []bool{false, true} {
		t.Run(map[bool]string{false: "item adds declaration", true: "item changes command"}[original], func(t *testing.T) {
			repo := filepath.Join(t.TempDir(), "repo")
			gittest.InitRepo(t, repo, "-b", "main")
			write := func(content string) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			trusted := "delivery:\n  derived_paths:\n    - paths: [derived.txt]\n      regenerate: printf trusted > derived.txt\n"
			if original {
				write(trusted)
			} else {
				write("watch:\n  push_remote: origin\n")
			}
			gittest.Run(t, repo, "add", ".")
			gittest.Run(t, repo, "commit", "-m", "default")
			defaultHead := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			gittest.Run(t, repo, "checkout", "-b", "feat/item")
			write("delivery:\n  derived_paths:\n    - paths: [derived.txt, forbidden.txt]\n      regenerate: touch forbidden.txt\n")
			gittest.Run(t, repo, "add", ".")
			gittest.Run(t, repo, "commit", "-m", "item config")
			config, err := DeliveryConfigAtCommit(t.Context(), preflight.ExecGitRunner{}, repo, filepath.Join(t.TempDir(), "missing-user.yml"), defaultHead)
			if err != nil {
				t.Fatal(err)
			}
			for _, declaration := range config.Delivery.DerivedPaths {
				command := exec.CommandContext(t.Context(), "sh", "-c", declaration.Regenerate)
				command.Dir = repo
				if output, err := command.CombinedOutput(); err != nil {
					t.Fatalf("regenerate: %v %s", err, output)
				}
			}
			if _, err := os.Stat(filepath.Join(repo, "forbidden.txt")); !os.IsNotExist(err) {
				t.Fatalf("item command executed: %v", err)
			}
			if original {
				content, err := os.ReadFile(filepath.Join(repo, "derived.txt"))
				if err != nil || string(content) != "trusted" {
					t.Fatalf("trusted output=%q err=%v", content, err)
				}
			} else if len(config.Delivery.DerivedPaths) != 0 {
				t.Fatalf("item declaration trusted: %+v", config.Delivery)
			}
		})
	}
}
