// Suite: declared Shared Build Cache reclamation.
// Invariant: only idle entries above a declared limit are removed.
// Boundary IN: GC command and temporary Project Config and cache files.
// Boundary OUT: engine race checks, covered by internal/sharedcache.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/sharedcache"
)

type sharedCacheFixture struct {
	root  string
	paths []string
	size  int64
}

func newSharedCacheFixture(t *testing.T, limit int) sharedCacheFixture {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	root := filepath.Join(home, "gocache")
	mustMkdir(t, filepath.Join(root, "ab"))
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), fmt.Sprintf("worktree:\n  shared_caches:\n    - path: '{home}/gocache'\n      kind: go-build\n      max_size_mb: %d\n      idle: 2h\n", limit))
	readme := sharedcache.GoBuildReadme + "\n"
	mustWrite(t, filepath.Join(root, "README"), readme)
	f := sharedCacheFixture{root: root, size: 3*1024*1024 + int64(len(readme))}
	now := time.Now()
	for i, age := range []time.Duration{10 * time.Hour, 3 * time.Hour, 0} {
		path := filepath.Join(root, "ab", fmt.Sprintf("%064x-d", i))
		if err := os.WriteFile(path, make([]byte, 1024*1024), 0o600); err != nil {
			t.Fatal(err)
		}
		stamp := now.Add(-age)
		if err := os.Chtimes(path, stamp, stamp); err != nil {
			t.Fatal(err)
		}
		f.paths = append(f.paths, path)
	}
	return f
}

func assertSharedCacheReport(t *testing.T, f sharedCacheFixture, limit int, dry bool, count int, output string) {
	t.Helper()
	action := "removed"
	if dry {
		action = "reclaimable"
	}
	want := fmt.Sprintf("  Shared Build Cache: %s\n    Size: %d bytes of %d allowed\n    Entries %s: %d (%d bytes)\n", f.root, f.size, limit*1024*1024, action, count, count*1024*1024)
	if !strings.HasSuffix(output, want) {
		t.Fatalf("output = %q, want suffix %q", output, want)
	}
}

func TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing(t *testing.T) {
	t.Parallel()
	f := newSharedCacheFixture(t, 1)
	before := snapshotBuildCacheTree(t, f.root)
	code, out, diagnostics := runBuildCacheGC(t, true)
	if code != exitOK || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	assertSharedCacheReport(t, f, 1, true, 2, out)
	if after := snapshotBuildCacheTree(t, f.root); before != after {
		t.Fatal("dry-run changed cache files")
	}
	for _, path := range f.paths {
		assertPathExists(t, path)
	}
}

func TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache(t *testing.T) {
	t.Parallel()
	f := newSharedCacheFixture(t, 1)
	code, out, diagnostics := runBuildCacheGC(t, false)
	if code != exitOK || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	assertSharedCacheReport(t, f, 1, false, 2, out)
	assertPathMissing(t, f.paths[0])
	assertPathMissing(t, f.paths[1])
	assertPathExists(t, f.paths[2])
	assertFileContent(t, filepath.Join(f.root, "README"), sharedcache.GoBuildReadme+"\n")
}

func TestGCKeepsASharedBuildCacheWithinItsLimit(t *testing.T) {
	t.Parallel()
	f := newSharedCacheFixture(t, 9)
	before := snapshotBuildCacheTree(t, f.root)
	code, out, diagnostics := runBuildCacheGC(t, false)
	if code != exitOK || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	assertSharedCacheReport(t, f, 9, false, 0, out)
	if after := snapshotBuildCacheTree(t, f.root); before != after {
		t.Fatal("within-limit cache changed")
	}
}

func TestGCReportsARefusedSharedBuildCache(t *testing.T) {
	t.Parallel()
	f := newSharedCacheFixture(t, 1)
	if err := os.Remove(filepath.Join(f.root, "README")); err != nil {
		t.Fatal(err)
	}
	code, out, diagnostics := runBuildCacheGC(t, false)
	if code != exitOK || diagnostics != "" || !strings.Contains(out, "  Shared Build Cache: "+f.root+"\n    Kept: ") {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	for _, path := range f.paths {
		assertPathExists(t, path)
	}
}
