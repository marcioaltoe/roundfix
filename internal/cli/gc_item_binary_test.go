// Suite: repository Item Binary retention.
// Invariant: GC removes only merged or unqueued Item Binaries and preserves evidence.
// Boundary IN: GC command, temporary homes, artifacts, and repository Delivery Queues.
// Boundary OUT: Run retention and Worktree Build Caches, covered by existing suites.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

type itemBinaryFixture struct {
	loaded roundconfig.Loaded
	paths  []string
}

func newItemBinaryFixture(t *testing.T) itemBinaryFixture {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	artifacts := filepath.Join(home, "artifacts")
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\n", artifacts))
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	q, err := s.CreateDeliveryQueue(t.Context(), repo, []string{"merged", "parked"})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range q.Items {
		item.Stage = store.DeliveryStageMerged
		if item.SpecSlug == "parked" {
			item.Stage = store.DeliveryStageParked
		}
		if err := s.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, slug := range []string{"merged", "unqueued", "parked"} {
		path := filepath.Join(artifacts, "delivery", slug, "item-binary")
		mustMkdir(t, path)
		mustWrite(t, filepath.Join(path, "copy"), "copy bytes")
		mustWrite(t, filepath.Join(filepath.Dir(path), "repository-gate.log"), "keep")
		paths = append(paths, path)
	}
	return itemBinaryFixture{loaded: loaded, paths: paths}
}

func assertItemBinaryTranscript(t *testing.T, f itemBinaryFixture, dry bool, output string) {
	t.Helper()
	action, bytes, title := "removed", "reclaimed", "Removed"
	if dry {
		action, bytes, title = "reclaimable", "reclaimable", "Reclaimable"
	}
	want := fmt.Sprintf("  Item Binaries %s: 2\n  Item Binary bytes %s: 20\n  Item Binaries kept: 1\n  %s Item Binaries:\n    %s\n    %s\n", action, bytes, title, f.paths[0], f.paths[1])
	offset := strings.Index(output, "  Item Binaries ")
	if offset < 0 || output[offset:] != want {
		t.Fatalf("Item Binary report = %q, want suffix %q", output, want)
	}
}

func TestGCDryRunListsReclaimableItemBinaries(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	before := snapshotBuildCacheTree(t, f.loaded.HomeDir)
	code, out, diagnostics := runBuildCacheGC(t, true)
	if code != exitOK || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	assertItemBinaryTranscript(t, f, true, out)
	if after := snapshotBuildCacheTree(t, f.loaded.HomeDir); before != after {
		t.Fatal("dry-run changed home bytes or paths")
	}
	for _, path := range f.paths {
		assertPathExists(t, path)
	}
}

func TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	code, out, diagnostics := runBuildCacheGC(t, false)
	if code != exitOK || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, diagnostics, out)
	}
	assertItemBinaryTranscript(t, f, false, out)
	assertPathMissing(t, f.paths[0])
	assertPathMissing(t, f.paths[1])
	assertPathExists(t, f.paths[2])
	for _, path := range f.paths {
		assertFileContent(t, filepath.Join(filepath.Dir(path), "repository-gate.log"), "keep")
	}
}

func TestGCItemBinaryRemovalRechecksQueueState(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	s := openDeliverReleaseStore(t, t.Context(), f.loaded.HomeDir)
	scan, err := scanItemBinaries(t.Context(), s, f.loaded)
	if err != nil || len(scan.paths) != 2 {
		t.Fatalf("scan=%+v err=%v", scan, err)
	}
	q, _, err := s.DeliveryQueue(t.Context(), f.loaded.GitRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range q.Items {
		if item.SpecSlug == "merged" {
			item.Stage = store.DeliveryStageRunning
			if err := s.UpdateDeliveryQueueItem(t.Context(), f.loaded.GitRoot, item); err != nil {
				t.Fatal(err)
			}
		}
	}
	// Report bytes must come from what is present at removal, not the preview.
	mustWrite(t, filepath.Join(f.paths[1], "later"), "new")
	removed, err := removeItemBinaries(t.Context(), s, f.loaded, scan)
	if err != nil || len(removed.paths) != 1 || removed.paths[0] != f.paths[1] || removed.bytes != 13 || removed.kept != 2 {
		t.Fatalf("removal=%+v err=%v", removed, err)
	}
	assertPathExists(t, f.paths[0])
	assertPathExists(t, f.paths[2])
	assertPathMissing(t, f.paths[1])
}

func TestGCItemBinariesKeepEveryUnmergedStage(t *testing.T) {
	t.Parallel()
	for _, stage := range []store.DeliveryStage{store.DeliveryStageQueued, store.DeliveryStageRunning, store.DeliveryStageReviewing, store.DeliveryStageArchiving, store.DeliveryStageGating, store.DeliveryStagePublishing, store.DeliveryStageChecking, store.DeliveryStageMerging, store.DeliveryStageParked} {
		t.Run(string(stage), func(t *testing.T) {
			t.Parallel()
			f := newItemBinaryFixture(t)
			s := openDeliverReleaseStore(t, t.Context(), f.loaded.HomeDir)
			q, _, err := s.DeliveryQueue(t.Context(), f.loaded.GitRoot)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range q.Items {
				item.Stage = stage
				if err := s.UpdateDeliveryQueueItem(t.Context(), f.loaded.GitRoot, item); err != nil {
					t.Fatal(err)
				}
			}
			report, err := sweepItemBinaries(t.Context(), s, f.loaded, false)
			if err != nil || len(report.paths) != 1 || report.paths[0] != f.paths[1] || report.kept != 2 {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			assertPathExists(t, f.paths[0])
			assertPathExists(t, f.paths[2])
		})
	}
}

func TestGCItemBinariesKeepSymlinksAndEscapingAncestors(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	outside := t.TempDir()
	mustMkdir(t, filepath.Join(outside, "item-binary"))
	mustWrite(t, filepath.Join(outside, "item-binary", "valuable"), "keep outside")
	root := filepath.Dir(filepath.Dir(f.paths[0]))
	mustMkdir(t, filepath.Join(root, "linked"))
	if err := os.Symlink(filepath.Join(outside, "item-binary"), filepath.Join(root, "linked", "item-binary")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	before := snapshotBuildCacheTree(t, outside)
	for _, dry := range []bool{true, false} {
		code, out, diagnostics := runBuildCacheGC(t, dry)
		if code != exitOK || diagnostics != "" || !strings.Contains(out, "  Item Binaries kept: 3\n") {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, diagnostics)
		}
	}
	if after := snapshotBuildCacheTree(t, outside); before != after {
		t.Fatal("GC followed an escaping path")
	}
	for _, path := range []string{filepath.Join(root, "linked", "item-binary"), filepath.Join(root, "escape")} {
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("symlink changed: %s: %v", path, err)
		}
	}
}

func TestGCItemBinaryRemovalRechecksThePath(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	s := openDeliverReleaseStore(t, t.Context(), f.loaded.HomeDir)
	scan, err := scanItemBinaries(t.Context(), s, f.loaded)
	if err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	mustWrite(t, filepath.Join(outside, "valuable"), "keep")
	if err := os.RemoveAll(f.paths[0]); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, f.paths[0]); err != nil {
		t.Fatal(err)
	}
	report, err := removeItemBinaries(t.Context(), s, f.loaded, scan)
	if err != nil || len(report.paths) != 1 || report.paths[0] != f.paths[1] || report.kept != 2 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	assertFileContent(t, filepath.Join(outside, "valuable"), "keep")
	info, err := os.Lstat(f.paths[0])
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("symlink changed: %v", err)
	}
}

func TestGCItemBinariesWithoutADatabaseAndWithANewQueue(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	artifacts := filepath.Join(home, "artifacts")
	loaded := roundconfig.Loaded{HomeDir: home, GitRoot: repo, Config: roundconfig.Builtin()}
	loaded.Config.Defaults.ArtifactDir = artifacts
	path := filepath.Join(artifacts, "delivery", "new", "item-binary")
	mustMkdir(t, path)
	mustWrite(t, filepath.Join(path, "copy"), "bytes")
	before := snapshotBuildCacheTree(t, home)
	scan, err := sweepItemBinaries(t.Context(), nil, loaded, true)
	if err != nil || len(scan.paths) != 1 || scan.bytes != 5 {
		t.Fatalf("scan=%+v err=%v", scan, err)
	}
	if after := snapshotBuildCacheTree(t, home); before != after {
		t.Fatal("dry-run created a database or artifacts")
	}
	s := openDeliverReleaseStore(t, t.Context(), home)
	if _, err := s.CreateDeliveryQueue(t.Context(), repo, []string{"new"}); err != nil {
		t.Fatal(err)
	}
	report, err := removeItemBinaries(t.Context(), nil, loaded, scan)
	if err != nil || len(report.paths) != 0 || report.kept != 1 {
		t.Fatalf("new queue removal=%+v err=%v", report, err)
	}
	assertPathExists(t, path)
}

func TestGCItemBinariesLeaveTheReportUnchangedWhenAbsent(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	for _, path := range f.paths {
		if err := os.RemoveAll(path); err != nil {
			t.Fatal(err)
		}
	}
	code, out, diagnostics := runBuildCacheGC(t, true)
	if code != exitOK || diagnostics != "" || strings.Contains(out, "Item Binar") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, diagnostics)
	}
	var rendered bytes.Buffer
	printItemBinaryReport(&rendered, itemBinaryReport{}, true)
	if rendered.Len() != 0 {
		t.Fatalf("absent Item Binaries added output: %q", rendered.String())
	}
}

func TestGCItemBinarySweepReportsSuccessfulRemovalsAfterAPathError(t *testing.T) {
	t.Parallel()
	f := newItemBinaryFixture(t)
	root := filepath.Dir(filepath.Dir(f.paths[0]))
	if err := os.Symlink("loop", filepath.Join(root, "loop")); err != nil {
		t.Fatal(err)
	}
	code, out, diagnostics := runBuildCacheGC(t, false)
	if code != exitRunFailed || !strings.Contains(diagnostics, "roundfix: gc failed:") || !strings.Contains(diagnostics, "loop") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, diagnostics)
	}
	assertItemBinaryTranscript(t, f, false, out)
	assertPathMissing(t, f.paths[0])
	assertPathMissing(t, f.paths[1])
	assertPathExists(t, f.paths[2])
}
