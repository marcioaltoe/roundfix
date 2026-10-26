package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWorktreeSharedCachesLoadsAndResolvesItsTemplate(t *testing.T) {
	t.Parallel()
	home, repo := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(repo, ".git"))
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "worktree:\n  shared_caches:\n    - path: '{tmpdir}/roundfix-go-build-{uid}'\n      kind: go-build\n      max_size_mb: 20480\n      idle: 4h\n")
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	entries := loaded.Config.Worktree.SharedCaches
	if len(entries) != 1 {
		t.Fatalf("entries = %+v", entries)
	}
	entry := entries[0]
	if entry.Path != "{tmpdir}/roundfix-go-build-{uid}" || entry.Kind != "go-build" || entry.MaxSizeMB != 20480 || entry.Idle != 4*time.Hour {
		t.Fatalf("entry = %+v", entry)
	}
	resolved, err := ResolveSharedCachePath(entry.Path, "/var/tmp/", "/Users/me", 501)
	if err != nil || resolved != "/var/tmp/roundfix-go-build-501" {
		t.Fatalf("resolved = %q, %v", resolved, err)
	}
	if len(Builtin().Worktree.SharedCaches) != 0 {
		t.Fatal("built-in shared cache list must be empty")
	}
	if strings.Contains(DefaultProjectConfigYAML(), "shared_caches") {
		t.Fatal("project template must stay unchanged")
	}
	user := []byte("worktree:\n  shared_caches:\n    - path: '{home}/user-cache'\n      kind: go-build\n      max_size_mb: 1\n      idle: 2h\n")
	cfg, err := ResolveConfigProposal(user, []byte("worktree:\n  shared_caches: []\n"))
	if err != nil || len(cfg.Worktree.SharedCaches) != 0 {
		t.Fatalf("empty override = %+v, %v", cfg.Worktree.SharedCaches, err)
	}
	cfg, err = ResolveConfigProposal(user, []byte("worktree:\n  shared_caches:\n    - path: '{home}/project-cache'\n      kind: go-build\n      max_size_mb: 2\n      idle: 3h\n"))
	if err != nil || len(cfg.Worktree.SharedCaches) != 1 || cfg.Worktree.SharedCaches[0].Path != "{home}/project-cache" {
		t.Fatalf("replacement = %+v, %v", cfg.Worktree.SharedCaches, err)
	}
}

func TestWorktreeSharedCachesRejectsAnUnsafeEntry(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, kind, limit, idle string }{
		{"unknown_kind", "{home}/cache", "other", "1", "2h"},
		{"relative_path", "cache", "go-build", "1", "2h"},
		{"parent_segment", "{home}/build/../cache", "go-build", "1", "2h"},
		{"no_limit", "{home}/cache", "go-build", "0", "2h"},
		{"short_idle", "{home}/cache", "go-build", "1", "1h59m"},
		{"unknown_token", "{home}/{unknown}", "go-build", "1", "2h"},
		{"unclean_path", "{home}//cache", "go-build", "1", "2h"},
		{"invalid_idle", "{home}/cache", "go-build", "1", "invalid"},
		{"overflow_limit", "{home}/cache", "go-build", "9223372036854775807", "2h"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(repo, ".git"))
			mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), fmt.Sprintf("worktree:\n  shared_caches:\n    - path: %q\n      kind: %s\n      max_size_mb: %s\n      idle: %s\n", tc.path, tc.kind, tc.limit, tc.idle))
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err == nil || !strings.Contains(err.Error(), "worktree.shared_caches") {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
