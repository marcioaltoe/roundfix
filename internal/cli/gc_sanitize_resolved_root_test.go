// Suite: sanitation of resolved Artifact Roots
// Invariant: existing ownership proofs apply to the root each Run wrote to.
// Boundary IN: public CLI, temporary Run Database and artifact directories.
// Boundary OUT: shared resolver derivation, covered by internal/store.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestGCSanitizeDryRunReachesTheDefaultRootTranscript(t *testing.T) {
	t.Parallel()
	home, repo := newResolvedSanitationWorkspace(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	run := seedResolvedSanitationRun(t, home, repo, "", "transcript", true, now.Add(-400*time.Hour))
	root := resolveGCTestArtifactRoot(t, repo, home)
	contents := "batch log"
	path := writeRunArtifact(t, root, run.ID, contents)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"gc", "sanitize"}, &stdout, &stderr)
	want := fmt.Sprintf("GC sanitation dry-run\n  Journal Retention: 336h\n  Cutoff: %s\n  Artifact Root: %s\n    Recorded as: \"\" by 1 Run(s)\n    Classification: orphaned\n    Evidence: no Active Run records this Artifact Root; terminal Runs recorded: 1; default derived from recorded checkout %q (Run %q)\n    would remove: %s (terminal Run %s completed before cutoff)\n    Directories reclaimable here: 1\n    Artifact bytes reclaimable here: %d\n  Directories reclaimable: 1\n  Artifact bytes reclaimable: %d\n  No directories removed; rerun with --apply to mutate.\n", now.Add(-336*time.Hour).Format(time.RFC3339Nano), root, repo, run.ID, path, run.ID, len(contents), len(contents))
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q\nstdout:\n%s\nwant:\n%s", code, stderr.String(), stdout.String(), want)
	}
	assertPathExists(t, path)
}

func newResolvedSanitationWorkspace(t *testing.T) (string, string) {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo)
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "seed")
	var err error
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, home, repo)
	return home, repo
}

func seedResolvedSanitationRun(t *testing.T, home, repo, recorded, name string, terminal bool, completed time.Time) store.Run {
	t.Helper()
	ctx := context.Background()
	s, err := store.Open(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	run, err := s.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, ArtifactDir: recorded, LocalBranch: "fix/" + name, HeadSHA: "abc123", SpecSlug: name})
	if err != nil {
		_ = s.Close()
		t.Fatal(err)
	}
	if terminal {
		if _, err := s.CompleteRun(ctx, run.ID, store.StateClean); err != nil {
			_ = s.Close()
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if terminal {
		setRunTimestamps(t, home, run.ID, completed, completed)
	}
	return run
}

func TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof(t *testing.T) {
	t.Parallel()
	home, repo := newResolvedSanitationWorkspace(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	old := now.Add(-400 * time.Hour)
	eligible := seedResolvedSanitationRun(t, home, repo, "", "eligible", true, old)
	active := seedResolvedSanitationRun(t, home, repo, "", "active", false, time.Time{})
	recent := seedResolvedSanitationRun(t, home, repo, "", "recent", true, now.Add(-time.Hour))
	root := resolveGCTestArtifactRoot(t, repo, home)
	removed := []string{writeRunArtifact(t, root, eligible.ID, "eligible"), writeRunArtifact(t, root, "run_absent", "absent")}
	kept := []string{writeRunArtifact(t, root, active.ID, "active"), writeRunArtifact(t, root, recent.ID, "recent")}
	for _, area := range []string{"reviews", "delivery"} {
		path := filepath.Join(root, area, "keep.log")
		mustMkdir(t, filepath.Dir(path))
		mustWrite(t, path, area)
		kept = append(kept, path)
	}
	outside := seedResolvedSanitationRun(t, home, repo, "outside-artifacts", "outside", true, old)
	kept = append(kept, writeRunArtifact(t, filepath.Join(repo, "outside-artifacts"), outside.ID, "outside"))
	// A known Run under another resolved root must not become an orphan here.
	kept = append(kept, writeRunArtifact(t, root, outside.ID, "recorded elsewhere"))

	linkedRepo := gcSanitationRepository(t, "linked")
	linkedRoot := resolveGCTestArtifactRoot(t, linkedRepo, home)
	linkedRun := seedResolvedSanitationRun(t, home, linkedRepo, "", "linked", true, old)
	target := t.TempDir()
	kept = append(kept, writeRunArtifact(t, target, linkedRun.ID, "linked target"))
	mustMkdir(t, filepath.Dir(linkedRoot))
	if err := os.Symlink(target, linkedRoot); err != nil {
		t.Fatal(err)
	}
	kept = append(kept, linkedRoot)

	sharedRepo := gcSanitationRepository(t, "shared")
	sharedRoot := resolveGCTestArtifactRoot(t, sharedRepo, home)
	sharedRun := seedResolvedSanitationRun(t, home, sharedRepo, "", "shared", true, old)
	_ = seedResolvedSanitationRun(t, home, repo, sharedRoot, "override", true, old)
	kept = append(kept, writeRunArtifact(t, sharedRoot, sharedRun.ID, "shared"))

	// Physical roots still preserve symlinked and non-directory runs/ entries.
	entryTarget := t.TempDir()
	link := filepath.Join(root, "runs", "run_link")
	if err := os.Symlink(entryTarget, link); err != nil {
		t.Fatal(err)
	}
	kept = append(kept, link)
	file := filepath.Join(root, "runs", "run_file")
	mustWrite(t, file, "keep file")
	kept = append(kept, file)
	output := runResolvedSanitationCLI(t, true)
	for _, want := range []string{"Classification: active", "Classification: outside Roundfix Home", "Classification: unsafe", "Classification: overridden", fmt.Sprintf("Run %s is recorded under Artifact Root %q", outside.ID, filepath.Join(repo, "outside-artifacts")), "Directories removed here: 2"} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q:\n%s", want, output)
		}
	}
	if got := gcReportCount(t, output, "Directories removed"); got != len(removed) {
		t.Fatalf("removed=%d want=%d: %s", got, len(removed), output)
	}
	for _, path := range removed {
		assertPathMissing(t, path)
	}
	for _, path := range kept {
		assertPathExists(t, path)
	}
	second := runResolvedSanitationCLI(t, true)
	if gcReportCount(t, second, "Directories removed") != 0 || gcReportCount(t, second, "Artifact bytes reclaimed") != 0 {
		t.Fatalf("second apply: %s", second)
	}
	for _, path := range kept {
		assertPathExists(t, path)
	}
}

func TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval(t *testing.T) {
	t.Parallel()
	home, repo := newResolvedSanitationWorkspace(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	_ = seedResolvedSanitationRun(t, home, repo, "", "known", true, now.Add(-400*time.Hour))
	root := resolveGCTestArtifactRoot(t, repo, home)
	const orphanID = "run_appears"
	path := writeRunArtifact(t, root, orphanID, "new Run")
	opens := 0
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		opener := deps.gc.openStoreReader
		deps.gc.openStoreReader = func(ctx context.Context, home string) (*store.Store, error) {
			opens++
			if opens == 2 {
				db := openGCRetentionDB(t, home)
				stamp := now.Format(time.RFC3339Nano)
				gcRetentionExec(t, db, `INSERT INTO runs (id, kind, state, git_root, repository_root, local_branch, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, orphanID, store.KindImplement, store.StateActive, repo, repo, "fix/appeared", stamp, stamp)
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}
			return opener(ctx, home)
		}
	})
	output := runResolvedSanitationCLI(t, true)
	reason := fmt.Sprintf("preserved: %s (Run %s appeared in durable Run metadata before removal)", path, orphanID)
	if opens != 2 || !strings.Contains(output, reason) || gcReportCount(t, output, "Directories removed") != 0 || gcReportCount(t, output, "Artifact bytes reclaimed") != 0 || !strings.Contains(output, "Directories removed here: 0") || !strings.Contains(output, "Artifact bytes reclaimed here: 0") || strings.Contains(output, "    removed: ") {
		t.Fatalf("opens=%d output=%s", opens, output)
	}
	assertPathExists(t, path)
}

func TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved(t *testing.T) {
	t.Parallel()
	home, repo := newResolvedSanitationWorkspace(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	run := seedResolvedSanitationRun(t, home, repo, "", "unresolved", true, now.Add(-400*time.Hour))
	root := resolveGCTestArtifactRoot(t, repo, home)
	path := writeRunArtifact(t, root, run.ID, "keep unresolved")
	db := openGCRetentionDB(t, home)
	gcRetentionExec(t, db, `UPDATE runs SET git_root = '', repository_root = '' WHERE id = ?`, run.ID)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	output := runResolvedSanitationCLI(t, true)
	reason := fmt.Sprintf("preserved because the Artifact Root of Run %q recorded as %q:", run.ID, "")
	if !strings.Contains(output, "Classification: unsafe") || !strings.Contains(output, reason) || !strings.Contains(output, "cannot be resolved") || gcReportCount(t, output, "Directories removed") != 0 {
		t.Fatalf("output=%s", output)
	}
	assertPathExists(t, path)
}

func runResolvedSanitationCLI(t *testing.T, apply bool) string {
	t.Helper()
	args := []string{"gc", "sanitize"}
	if apply {
		args = append(args, "--apply")
	}
	var stdout, stderr bytes.Buffer
	if code := runCLIContext(t, context.Background(), args, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	return stdout.String()
}

func TestGCSanitizeStopsBeforeRemovalWhenAbsenceRecheckFails(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"open", "read"} {
		t.Run(phase, func(t *testing.T) {
			t.Parallel()
			home, repo := newResolvedSanitationWorkspace(t)
			now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
			withGCNow(t, now)
			_ = seedResolvedSanitationRun(t, home, repo, "", "known", true, now.Add(-400*time.Hour))
			root := resolveGCTestArtifactRoot(t, repo, home)
			path := writeRunArtifact(t, root, "run_absent", "keep on read failure")
			opens := 0
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				opener := deps.gc.openStoreReader
				deps.gc.openStoreReader = func(ctx context.Context, home string) (*store.Store, error) {
					opens++
					if opens == 2 && phase == "open" {
						return nil, fmt.Errorf("removal reader unavailable")
					}
					reader, err := opener(ctx, home)
					if err == nil && opens == 2 {
						// A closed reader gives Store.Run a real database read error.
						if err := reader.Close(); err != nil {
							t.Fatal(err)
						}
					}
					return reader, err
				}
			})
			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, context.Background(), []string{"gc", "sanitize", "--apply"}, &stdout, &stderr)
			if code != exitRunFailed || !strings.Contains(stderr.String(), "gc sanitize failed") || stdout.Len() != 0 {
				t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if phase == "read" && !strings.Contains(stderr.String(), "recheck Run") {
				t.Fatalf("missing read error: %s", stderr.String())
			}
			assertPathExists(t, path)
		})
	}
}

func TestGCSanitizeReportsDistinctRecordedValuesInOrder(t *testing.T) {
	t.Parallel()
	home, repo := newResolvedSanitationWorkspace(t)
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	withGCNow(t, now)
	root := resolveGCTestArtifactRoot(t, repo, home)
	relative, err := filepath.Rel(repo, root)
	if err != nil {
		t.Fatal(err)
	}
	homeRelative, err := filepath.Rel(home, root)
	if err != nil {
		t.Fatal(err)
	}
	expanded := "~/" + filepath.ToSlash(homeRelative)
	for i, value := range []string{"", "", relative, expanded, root} {
		run := seedResolvedSanitationRun(t, home, repo, value, fmt.Sprintf("recorded-%d", i), true, now.Add(-400*time.Hour))
		writeRunArtifact(t, root, run.ID, "log")
	}
	output := runResolvedSanitationCLI(t, false)
	want := fmt.Sprintf("  Artifact Root: %s\n    Recorded as: %q by 2 Run(s)\n    Recorded as: %q by 1 Run(s)\n    Recorded as: %q by 1 Run(s)\n    Classification: orphaned", root, "", relative, expanded)
	if strings.Count(output, "  Artifact Root: ") != 1 || strings.Count(output, "    Recorded as:") != 3 || !strings.Contains(output, want) {
		t.Fatalf("output=%s\nwant=%s", output, want)
	}
}
