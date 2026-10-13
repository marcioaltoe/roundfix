// Suite: declared Worktree Build Caches
// Invariant: Doctor and GC observe only declared caches of proven ended worktrees.
// Boundary IN: temporary homes, recorded Runs and queue items, command output and cache paths
// Boundary OUT: Run retention policy and process signaling, covered by their own suites
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

type buildCacheFixture struct {
	loaded                roundconfig.Loaded
	ended, active         store.Run
	endedPath, activePath string
}

func newBuildCacheFixture(t *testing.T) buildCacheFixture {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	location := filepath.Join(home, "worktrees")
	artifact := filepath.Join(home, "artifacts")
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), fmt.Sprintf("defaults:\n  artifact_dir: %q\nworktree:\n  location: %q\n  caches: [.gocache]\n", artifact, location))
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	create := func(name string, ended bool) store.Run {
		run, err := runStore.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindResolve, HeadRepository: "owner/project", HeadBranch: name, BaseRepository: "owner/project", PRNumber: name, GitRoot: repo, LocalBranch: name, HeadSHA: "abc123", ArtifactDir: artifact})
		if err != nil {
			t.Fatal(err)
		}
		worktree := filepath.Join(location, "repo", run.ID)
		run, err = runStore.SetRunWorkDir(t.Context(), run.ID, worktree)
		if err != nil {
			t.Fatal(err)
		}
		mustMkdir(t, filepath.Join(worktree, ".gocache"))
		mustWrite(t, filepath.Join(worktree, ".gocache", "build"), "cache bytes")
		mustWrite(t, filepath.Join(worktree, "source.go"), "keep this")
		if ended {
			_, err = runStore.CompleteRun(t.Context(), run.ID, store.StateClean)
			if err != nil {
				t.Fatal(err)
			}
		}
		return run
	}
	ended, active := create("ended", true), create("active", false)
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	return buildCacheFixture{loaded: loaded, ended: ended, active: active, endedPath: filepath.Join(ended.WorkDir, ".gocache"), activePath: filepath.Join(active.WorkDir, ".gocache")}
}
func runBuildCacheGC(t *testing.T, dry bool) (int, string, string) {
	t.Helper()
	args := []string{"gc"}
	if dry {
		args = append(args, "--dry-run")
	}
	var out, err bytes.Buffer
	code := runCLIContext(t, t.Context(), args, &out, &err)
	return code, out.String(), err.String()
}
func assertBuildCacheTranscript(t *testing.T, fixture buildCacheFixture, dry bool, output string) {
	t.Helper()
	// Every non-elided line of transcripts 1 and 2 is exact. Only paths,
	// timestamps and the declared <bytes> placeholders can vary.
	heading, count, rows, artifacts, orphan, final := "GC complete", "Runs pruned", "Journal rows removed", "Artifact bytes reclaimed", "Orphan artifact dirs removed", "  Run artifact bytes reclaimed: 0\n  Compaction: not needed\n  Database bytes before: [0-9]+\n  Database bytes after: [0-9]+\n"
	action, byteAction, title := "removed", "reclaimed", "Removed"
	if dry {
		heading, count, rows, artifacts, orphan = "GC dry-run", "Runs eligible", "Journal rows eligible", "Artifact bytes reclaimable", "Orphan artifact dirs eligible"
		final = "  Database bytes reclaimable \\(estimated\\): 0\n"
		action, byteAction, title = "reclaimable", "reclaimable", "Reclaimable"
	}
	pattern := "^" + heading + "\n  Journal Retention: 336h\n  Cutoff: [^\n]+\n  " + count + ": 0\n  " + rows + ": 0\n  " + artifacts + ": 0\n  " + orphan + ": 0\n" +
		"  Run Retention: 30 days\n  Run Retention cutoff: [^\n]+\n  Runs " + map[bool]string{true: "removable", false: "removed"}[dry] + ": 0\n" +
		"  Runs kept past the cutoff: 0 \\(queue-referenced 0, worktree present 0, artifact root unproven 0\\)\n" +
		"  Rows " + map[bool]string{true: "removable", false: "removed"}[dry] + ": runs=0 run_events=0 run_agent_selections=0 run_token_usage=0 active_run_locks=0\n" + final +
		"  Worktree Build Caches declared: .gocache\n  Worktree Build Caches " + action + ": 1\n  Worktree Build Cache bytes " + byteAction + ": 11\n  Worktree Build Caches kept: 1\n  " + title + " Worktree Build Caches:\n    " + regexp.QuoteMeta(fixture.endedPath) + "\n$"
	if len(fixture.loaded.Config.Worktree.Caches) == 0 {
		pattern = strings.Split(pattern, "  Worktree Build Caches declared:")[0] + "$"
	}
	if !regexp.MustCompile(pattern).MatchString(output) {
		t.Fatalf("GC transcript mismatch:\n%s\npattern: %s", output, pattern)
	}
}
func TestGCDryRunListsWorktreeBuildCaches(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	before := snapshotBuildCacheTree(t, f.loaded.HomeDir)
	code, out, err := runBuildCacheGC(t, true)
	if code != exitOK || err != "" {
		t.Fatalf("exit=%d stderr=%q", code, err)
	}
	assertBuildCacheTranscript(t, f, true, out)
	if after := snapshotBuildCacheTree(t, f.loaded.HomeDir); before != after {
		t.Fatal("dry run changed home bytes or paths")
	}
}
func TestGCRemovesWorktreeBuildCachesOfEndedWorktrees(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	code, out, err := runBuildCacheGC(t, false)
	if code != exitOK || err != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, err, out)
	}
	assertBuildCacheTranscript(t, f, false, out)
	assertPathMissing(t, f.endedPath)
	assertPathExists(t, f.activePath)
	assertPathExists(t, filepath.Join(f.ended.WorkDir, "source.go"))
	code, out, err = runBuildCacheGC(t, false)
	if code != exitOK || err != "" || !strings.Contains(out, "Worktree Build Caches removed: 0\n  Worktree Build Cache bytes reclaimed: 0\n") {
		t.Fatalf("second GC: %d %q %q", code, out, err)
	}
}
func TestDoctorStorageReportsWorktreeBuildCaches(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	before := snapshotBuildCacheTree(t, f.loaded.HomeDir)
	code, out, err := runDoctorWithRealStorage(t, f.loaded)
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, err)
	}
	want := "storage: found (Runs reclaimable: 0; Run Database free bytes: 0; Worktree Build Caches reclaimable: 1 (11 bytes); next: roundfix gc)"
	if line := buildCacheStorageLine(out); line != want {
		t.Fatalf("storage line = %q, want %q", line, want)
	}
	if after := snapshotBuildCacheTree(t, f.loaded.HomeDir); before != after {
		t.Fatal("Doctor changed home bytes or paths")
	}
}
func TestDoctorStorageBoundsTheWorktreeBuildCacheWalk(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	calls := 0
	clock := time.Now()
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.gc.now = func() time.Time { calls++; return clock.Add(time.Duration(calls) * 3 * time.Second) }
	})
	started := time.Now()
	code, out, err := runDoctorWithRealStorage(t, f.loaded)
	if code != exitOK {
		t.Fatalf("exit=%d stderr=%q", code, err)
	}
	want := "storage: found (Runs reclaimable: 0; Run Database free bytes: 0; Worktree Build Caches reclaimable: 1 (at least 0 bytes); next: roundfix gc)"
	if line := buildCacheStorageLine(out); line != want {
		t.Fatalf("bounded storage = %q", line)
	}
	if time.Since(started) > doctorWorktreeBuildCacheBudget {
		t.Fatal("bounded walk exceeded Doctor budget")
	}
	assertPathExists(t, f.endedPath)
}
func TestDoctorStorageIsUnchangedWithoutDeclaredCaches(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	f.loaded.Config.Worktree.Caches = nil
	code, out, err := runDoctorWithRealStorage(t, f.loaded)
	if code != exitOK {
		t.Fatalf("exit=%d stderr=%q", code, err)
	}
	want := "storage: ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: 0)"
	if line := buildCacheStorageLine(out); line != want {
		t.Fatalf("storage = %q, want %q", line, want)
	}
}
func TestGCIsUnchangedWithoutDeclaredCaches(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	f.loaded.Config.Worktree.Caches = nil
	for _, dry := range []bool{true, false} {
		report, err := runGC(t.Context(), gcOptions{dryRun: dry}, f.loaded)
		if err != nil {
			t.Fatal(err)
		}
		var got bytes.Buffer
		printGCReport(&got, report)
		assertBuildCacheTranscript(t, f, dry, got.String())
	}
	assertPathExists(t, f.endedPath)
}
func TestGCKeepsWorktreeBuildCachesOfWorkInProgress(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	stages := map[string]string{"merged": string(store.DeliveryStageMerged)}
	for _, stage := range []store.DeliveryStage{store.DeliveryStageQueued, store.DeliveryStageRunning, store.DeliveryStageReviewing, store.DeliveryStageArchiving, store.DeliveryStageGating, store.DeliveryStagePublishing, store.DeliveryStageChecking, store.DeliveryStageMerging, store.DeliveryStageParked} {
		stages[string(stage)] = string(stage)
	}
	addBuildCacheQueue(t, f, stages)
	code, out, err := runBuildCacheGC(t, false)
	if code != exitOK || err != "" {
		t.Fatalf("exit=%d out=%q err=%q", code, out, err)
	}
	assertPathExists(t, f.activePath)
	for slug, stage := range stages {
		if stage != string(store.DeliveryStageMerged) {
			assertPathExists(t, filepath.Join(f.loaded.Config.Worktree.Location, slug, ".gocache"))
		}
	}
	assertPathMissing(t, filepath.Join(f.loaded.Config.Worktree.Location, "merged", ".gocache"))
	if !strings.Contains(out, "Worktree Build Caches removed: 2") || !strings.Contains(out, fmt.Sprintf("Worktree Build Caches kept: %d", len(stages))) {
		t.Fatalf("queue report: %q", out)
	}
}
func TestGCSkipsAWorktreeBuildCacheThatEscapesItsWorktree(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	outside := t.TempDir()
	mustWrite(t, filepath.Join(outside, "valuable"), "keep")
	if err := os.RemoveAll(f.endedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, f.endedPath); err != nil {
		t.Fatal(err)
	}
	before := snapshotBuildCacheTree(t, outside)
	for _, dry := range []bool{true, false} {
		code, out, err := runBuildCacheGC(t, dry)
		if code != exitOK || err != "" || strings.Contains(out, "    "+f.endedPath) {
			t.Fatalf("exit=%d out=%q err=%q", code, out, err)
		}
	}
	if before != snapshotBuildCacheTree(t, outside) {
		t.Fatal("symlink target changed")
	}
	if info, err := os.Lstat(f.endedPath); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("cache symlink changed: %v", err)
	}
}
func addBuildCacheQueue(t *testing.T, f buildCacheFixture, stages map[string]string) {
	t.Helper()
	s, err := store.Open(t.Context(), f.loaded.HomeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	slugs := []string{}
	for slug := range stages {
		slugs = append(slugs, slug)
	}
	queue, err := s.CreateDeliveryQueue(t.Context(), f.loaded.GitRoot, slugs)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range queue.Items {
		worktree := filepath.Join(f.loaded.Config.Worktree.Location, item.SpecSlug)
		mustMkdir(t, filepath.Join(worktree, ".gocache"))
		mustWrite(t, filepath.Join(worktree, ".gocache", "build"), "cache bytes")
		if _, _, _, err := s.RecordDeliveryQueueItemWorktree(t.Context(), f.loaded.GitRoot, item.SpecSlug, "roundfix/item-"+item.SpecSlug, worktree); err != nil {
			t.Fatal(err)
		}
		item.Stage = store.DeliveryStage(stages[item.SpecSlug])
		if err := s.UpdateDeliveryQueueItem(t.Context(), f.loaded.GitRoot, item); err != nil {
			t.Fatal(err)
		}
	}
}
func buildCacheStorageLine(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "storage:") {
			return line
		}
	}
	return ""
}
func snapshotBuildCacheTree(t *testing.T, root string) string {
	t.Helper()
	var snapshot strings.Builder
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Fprintf(&snapshot, "%s:%s\n", path, entry.Type())
		if entry.Type().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Fprintf(&snapshot, "%x\n", data)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot.String()
}

// Keep context cancellation observable without an abandoned measurement.
func TestBuildCacheMeasurementHonorsCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, _, err := measureWorktreeBuildCache(ctx, t.TempDir(), func() bool { return false })
	if err != context.Canceled {
		t.Fatalf("canceled walk error = %v", err)
	}
}

func TestBuildCacheDiscoveryFindsRetainedTaskWorktrees(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	taskCache := filepath.Join(filepath.Dir(f.ended.WorkDir), f.ended.ID+".task_04", ".gocache")
	activeTaskCache := filepath.Join(filepath.Dir(f.active.WorkDir), f.active.ID+".task_07", ".gocache")
	for _, cache := range []string{taskCache, activeTaskCache} {
		mustMkdir(t, cache)
		mustWrite(t, filepath.Join(cache, "build"), "task cache")
	}
	// Terminal Run Worktree gone, but its Task Worktree remains.
	if err := os.RemoveAll(f.ended.WorkDir); err != nil {
		t.Fatal(err)
	}
	code, out, err := runBuildCacheGC(t, false)
	if code != exitOK || err != "" || !strings.Contains(out, "Worktree Build Caches removed: 1") {
		t.Fatalf("%d %q %q", code, out, err)
	}
	assertPathMissing(t, taskCache)
	assertPathExists(t, activeTaskCache)
}

func TestBuildCacheDiscoverySkipsOutsideAndAncestorSymlinkEscapes(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	outside := t.TempDir()
	mustMkdir(t, filepath.Join(outside, "cache"))
	mustWrite(t, filepath.Join(outside, "cache", "build"), "outside")
	if err := os.Symlink(outside, filepath.Join(f.ended.WorkDir, "escape")); err != nil {
		t.Fatal(err)
	}
	f.loaded.Config.Worktree.Caches = []string{"escape/cache"}
	reader, err := store.OpenReader(t.Context(), f.loaded.HomeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	scan, err := scanWorktreeBuildCaches(t.Context(), reader, f.loaded, 0)
	if err != nil || scan.Bytes != 0 || len(scan.Reclaimable) != 0 {
		t.Fatalf("escaping ancestor scan = %#v, %v", scan, err)
	}
	f.loaded.Config.Worktree.Caches = []string{".gocache"}
	// The recorded worktree is outside the configured location.
	f.loaded.Config.Worktree.Location = t.TempDir()
	scan, err = scanWorktreeBuildCaches(t.Context(), reader, f.loaded, 0)
	if err != nil || scan.Bytes != 0 || len(scan.Reclaimable) != 0 || scan.Kept != 0 {
		t.Fatalf("outside worktree scan = %#v, %v", scan, err)
	}
	assertPathExists(t, filepath.Join(outside, "cache", "build"))
}

func TestBuildCacheRemovalRechecksQueueStateAndSharedOwnership(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	addBuildCacheQueue(t, f, map[string]string{"merged": "merged"})
	s, err := store.Open(t.Context(), f.loaded.HomeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	scan, err := scanWorktreeBuildCaches(t.Context(), s, f.loaded, 0)
	if err != nil || len(scan.Reclaimable) != 2 {
		t.Fatalf("initial scan: %#v %v", scan, err)
	}
	queue, _, err := s.DeliveryQueue(t.Context(), f.loaded.GitRoot)
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageReviewing
	if err := s.UpdateDeliveryQueueItem(t.Context(), f.loaded.GitRoot, item); err != nil {
		t.Fatal(err)
	}
	// A newly recorded active Run also owns the terminal Run's worktree.
	if _, err := s.SetRunWorkDir(t.Context(), f.active.ID, f.ended.WorkDir); err != nil {
		t.Fatal(err)
	}
	report, err := removeWorktreeBuildCaches(t.Context(), s, f.loaded, scan)
	if err != nil || len(report.Reclaimable) != 0 || report.Bytes != 0 || report.Kept != 3 {
		t.Fatalf("rechecked removal: %#v %v", report, err)
	}
	assertPathExists(t, f.endedPath)
	assertPathExists(t, filepath.Join(f.loaded.Config.Worktree.Location, "merged", ".gocache"))
}

func TestBuildCacheDiscoveryProtectsPhysicalAliases(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	alias := filepath.Join(f.loaded.Config.Worktree.Location, "alias")
	if err := os.Symlink(f.ended.WorkDir, alias); err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(t.Context(), f.loaded.HomeDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetRunWorkDir(t.Context(), f.active.ID, alias); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := runBuildCacheGC(t, false)
	if code != exitOK || errOut != "" || !strings.Contains(out, "Worktree Build Caches removed: 0") {
		t.Fatalf("alias protection: %d %q %q", code, out, errOut)
	}
	assertPathExists(t, f.endedPath)
}

func TestBuildCacheDiscoveryPreservesGitMetadataThroughAnAlias(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	gitDir := filepath.Join(f.ended.WorkDir, ".git")
	mustMkdir(t, filepath.Join(gitDir, "cache"))
	mustWrite(t, filepath.Join(gitDir, "cache", "index"), "Git metadata")
	if err := os.Symlink(gitDir, filepath.Join(f.ended.WorkDir, "build")); err != nil {
		t.Fatal(err)
	}
	f.loaded.Config.Worktree.Caches = []string{"build/cache"}
	report, err := runGC(t.Context(), gcOptions{}, f.loaded)
	if err != nil || len(report.BuildCaches.Reclaimable) != 0 {
		t.Fatalf("Git alias report = %#v %v", report.BuildCaches, err)
	}
	assertPathExists(t, filepath.Join(gitDir, "cache", "index"))
}

func TestBuildCacheDiscoveryPreservesGitMetadataBehindALinkedGitDirectory(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	gitDir := filepath.Join(f.ended.WorkDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		t.Fatal(err)
	}
	metadata := filepath.Join(f.ended.WorkDir, "metadata")
	mustMkdir(t, filepath.Join(metadata, "cache"))
	mustWrite(t, filepath.Join(metadata, "cache", "index"), "Git metadata")
	if err := os.Symlink("metadata", gitDir); err != nil {
		t.Fatal(err)
	}
	f.loaded.Config.Worktree.Caches = []string{"metadata/cache"}
	report, err := runGC(t.Context(), gcOptions{}, f.loaded)
	if err != nil || len(report.BuildCaches.Reclaimable) != 0 {
		t.Fatalf("linked Git directory report = %#v %v", report.BuildCaches, err)
	}
	assertPathExists(t, filepath.Join(metadata, "cache", "index"))
}

func TestBuildCacheBudgetReportsMeasuredLowerBound(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "first"), "cache bytes")
	// Checks: before opening, before reading, before the entry, before the
	// next batch. Expire after the first entry's size has been counted.
	checks := 0
	size, bounded, err := measureWorktreeBuildCache(t.Context(), root, func() bool { checks++; return checks >= 4 })
	if err != nil || !bounded || size != 11 {
		t.Fatalf("bounded bytes = %d, bounded=%v err=%v", size, bounded, err)
	}
}

func TestBuildCacheGCContinuesAfterAPathFailure(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	mustWrite(t, f.loaded.ProjectConfigPath, fmt.Sprintf("defaults:\n  artifact_dir: %q\nworktree:\n  location: %q\n  caches: [.gocache, loop/cache]\n", f.loaded.Config.Defaults.ArtifactDir, f.loaded.Config.Worktree.Location))
	loop := filepath.Join(f.ended.WorkDir, "loop")
	if err := os.Symlink("loop", loop); err != nil {
		t.Fatal(err)
	}
	code, out, err := runBuildCacheGC(t, false)
	if code != exitRunFailed || !strings.Contains(err, "roundfix: gc failed:") || !strings.Contains(err, "loop/cache") {
		t.Fatalf("failure exit=%d stdout=%q stderr=%q", code, out, err)
	}
	for _, want := range []string{"Worktree Build Caches declared: .gocache, loop/cache", "Worktree Build Caches removed: 1", "Worktree Build Cache bytes reclaimed: 11", "Worktree Build Caches kept: 1", "Removed Worktree Build Caches:\n    " + f.endedPath + "\n"} {
		if !strings.Contains(out, want) {
			t.Fatalf("failure report %q lacks %q", out, want)
		}
	}
	assertPathMissing(t, f.endedPath)
	assertPathExists(t, f.activePath)
	if info, err := os.Lstat(loop); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("unproven path changed: %v", err)
	}
}

func TestBuildCacheDoctorReportsUnreadablePathsAsPartial(t *testing.T) {
	t.Parallel()
	f := newBuildCacheFixture(t)
	loop := filepath.Join(f.ended.WorkDir, "loop")
	if err := os.Symlink("loop", loop); err != nil {
		t.Fatal(err)
	}
	f.loaded.Config.Worktree.Caches = []string{"loop/cache"}
	code, out, err := runDoctorWithRealStorage(t, f.loaded)
	if code != exitOK || !strings.HasPrefix(buildCacheStorageLine(out), "storage: partial (") || !strings.Contains(out, "could not read Worktree Build Caches") {
		t.Fatalf("partial Doctor exit=%d stdout=%q stderr=%q", code, out, err)
	}
	assertPathExists(t, f.endedPath)
}
