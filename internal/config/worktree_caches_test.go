package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWorktreeCachesDefaultsToNone(t *testing.T) {
	t.Parallel()
	if len(Builtin().Worktree.Caches) != 0 {
		t.Fatal("built-in caches must be empty")
	}
	if !strings.Contains(DefaultProjectConfigYAML(), "caches: []") {
		t.Fatal("template must declare the empty cache list")
	}
}

func TestWorktreeCachesAcceptsRepositoryRelativeDirectories(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"user", "project"} {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(repo, ".git"))
			path := filepath.Join(repo, ".roundfixrc.yml")
			if scope == "user" {
				path = filepath.Join(home, ".roundfix", "config.yml")
				mustMkdir(t, filepath.Dir(path))
			}
			mustWrite(t, path, "worktree:\n  caches: [.gocache, build/cache, .gocache]\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(loaded.Config.Worktree.Caches, []string{".gocache", "build/cache"}) {
				t.Fatalf("caches = %v", loaded.Config.Worktree.Caches)
			}
		})
	}
	cfg, err := ResolveConfigProposal([]byte("worktree:\n  caches: [user-cache]\n"), []byte("worktree:\n  caches: []\n"))
	if err != nil || len(cfg.Worktree.Caches) != 0 {
		t.Fatalf("project empty list did not override user list: %v, %v", cfg.Worktree.Caches, err)
	}
}

func TestWorktreeCachesRefusesEscapingEntries(t *testing.T) {
	t.Parallel()
	for _, entry := range []string{"/tmp/cache", "../cache", "build/../cache", ".git", ".git/cache", "", ".", "cache/", "./cache", "cache//output", `C:\cache`, `build\cache`} {
		t.Run(entry, func(t *testing.T) {
			t.Parallel()
			_, err := ResolveConfigProposal(nil, []byte(fmt.Sprintf("worktree:\n  caches: [%q]\n", entry)))
			if err == nil || !strings.Contains(err.Error(), "worktree.caches") || !strings.Contains(err.Error(), fmt.Sprintf("%q", entry)) {
				t.Fatalf("entry %q error = %v", entry, err)
			}
		})
	}
}
