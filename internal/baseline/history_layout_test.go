package baseline

// Suite: history layout discovery
// Invariant: discovery reports every and only retired file outside docs/history,
// with its source identity, without changing the repository.
// Boundary IN: legacy history trees, active retirement metadata, orphan Review
// Artifact liveness, destination occupancy, and filesystem content.
// Boundary OUT: Baseline plan serialization, relocation apply, rollback, and reporting.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/rounds"
)

func TestDiscoverHistoryLayoutLegacyArchiveShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		files map[string]string
		want  []HistoryRelocation
	}{
		{
			name: "archive directories nested in documentation trees",
			files: map[string]string{
				"docs/specs/_archived/0001-widget/_prd.md":     "nested spec\n",
				"docs/specs/_archived/0001-widget/task_01.md":  "nested task\n",
				"docs/findings/_archived/2026-08-01-widget.md": "nested finding\n",
			},
			want: []HistoryRelocation{
				{From: "docs/findings/_archived/2026-08-01-widget.md", To: "docs/history/findings/2026-08-01-widget.md", ContentIdentity: historyFixtureIdentity("nested finding\n")},
				{From: "docs/specs/_archived/0001-widget/_prd.md", To: "docs/history/specs/0001-widget/_prd.md", ContentIdentity: historyFixtureIdentity("nested spec\n")},
				{From: "docs/specs/_archived/0001-widget/task_01.md", To: "docs/history/specs/0001-widget/task_01.md", ContentIdentity: historyFixtureIdentity("nested task\n")},
			},
		},
		{
			name: "archive root outside the documentation tree",
			files: map[string]string{
				"_archived/specs/0002-gadget/_prd.md":     "root spec\n",
				"_archived/findings/2026-08-02-gadget.md": "root finding\n",
				"_archived/adr/0042-gadget.md":            "root ADR\n",
				"_archived/backlog/2026-08-02-gadget.md":  "root backlog\n",
			},
			want: []HistoryRelocation{
				{From: "_archived/adr/0042-gadget.md", To: "docs/history/adr/0042-gadget.md", ContentIdentity: historyFixtureIdentity("root ADR\n")},
				{From: "_archived/backlog/2026-08-02-gadget.md", To: "docs/history/backlog/2026-08-02-gadget.md", ContentIdentity: historyFixtureIdentity("root backlog\n")},
				{From: "_archived/findings/2026-08-02-gadget.md", To: "docs/history/findings/2026-08-02-gadget.md", ContentIdentity: historyFixtureIdentity("root finding\n")},
				{From: "_archived/specs/0002-gadget/_prd.md", To: "docs/history/specs/0002-gadget/_prd.md", ContentIdentity: historyFixtureIdentity("root spec\n")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := t.TempDir()
			historyWriteFiles(t, repo, test.files)
			before := historySnapshot(t, repo)

			got, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
			if err != nil {
				t.Fatalf("DiscoverHistoryLayout() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("DiscoverHistoryLayout() relocations = %#v, want %#v", got, test.want)
			}
			if len(collisions) != 0 {
				t.Fatalf("DiscoverHistoryLayout() collisions = %#v, want none", collisions)
			}

			again, againCollisions, err := DiscoverHistoryLayout(context.Background(), repo)
			if err != nil {
				t.Fatalf("second DiscoverHistoryLayout() error = %v", err)
			}
			if !reflect.DeepEqual(again, got) || !reflect.DeepEqual(againCollisions, collisions) {
				t.Fatalf("second discovery = (%#v, %#v), want (%#v, %#v)", again, againCollisions, got, collisions)
			}
			historyAssertUnchanged(t, repo, before)
		})
	}
}

func TestDiscoverHistoryLayoutCurrentLayoutReportsNothing(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	historyWriteFiles(t, repo, map[string]string{
		"docs/history/specs/0001-widget/_prd.md": "retired spec\n",
		"docs/history/findings/old.md":           "retired finding\n",
		"docs/history/adr/0042-old.md":           "retired ADR\n",
		"docs/history/backlog/declined.md":       "declined backlog\n",
		"docs/history/reviews/pr-42/round.md":    "retired review\n",
		"docs/adr/0043-active.md":                "---\nstatus: accepted\n---\n",
		"docs/backlog/2026-08-12-active.md":      "---\nstatus: open\n---\n",
	})
	before := historySnapshot(t, repo)

	relocations, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
	if err != nil {
		t.Fatalf("DiscoverHistoryLayout() error = %v", err)
	}
	if len(relocations) != 0 || len(collisions) != 0 {
		t.Fatalf("DiscoverHistoryLayout() = (%#v, %#v), want no findings", relocations, collisions)
	}
	historyAssertUnchanged(t, repo, before)
}

func TestDiscoverHistoryLayoutClassifiesActiveDocuments(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	historyWriteFiles(t, repo, map[string]string{
		"docs/adr/0040-retired.md":            "---\nstatus: superseded\n---\n",
		"docs/adr/0041-pending.md":            "---\nstatus: proposed\n---\n",
		"docs/backlog/2026-08-11-declined.md": "---\nstatus: declined\n---\n",
		"docs/backlog/2026-08-12-open.md":     "---\nstatus: open\n---\n",
	})
	before := historySnapshot(t, repo)

	got, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
	if err != nil {
		t.Fatalf("DiscoverHistoryLayout() error = %v", err)
	}
	want := []HistoryRelocation{
		{From: "docs/adr/0040-retired.md", To: "docs/history/adr/0040-retired.md", ContentIdentity: historyFixtureIdentity("---\nstatus: superseded\n---\n")},
		{From: "docs/backlog/2026-08-11-declined.md", To: "docs/history/backlog/2026-08-11-declined.md", ContentIdentity: historyFixtureIdentity("---\nstatus: declined\n---\n")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DiscoverHistoryLayout() relocations = %#v, want %#v", got, want)
	}
	if len(collisions) != 0 {
		t.Fatalf("DiscoverHistoryLayout() collisions = %#v, want none", collisions)
	}
	historyAssertUnchanged(t, repo, before)
}

func TestDiscoverHistoryLayoutClassifiesOrphanReviews(t *testing.T) {
	t.Parallel()

	repo, mainHead := historyGitRepo(t)
	liveHead := historyBranchCommit(t, repo, "feature/live")
	gittest.Run(t, repo, "switch", "main")

	finishedReview := filepath.Join(repo, "docs", "specs", "reviews", "pr-101")
	historyPersistRound(t, finishedReview, "feature/merged", mainHead)
	historyWriteFiles(t, repo, map[string]string{
		"docs/specs/reviews/pr-101/issues/001.md": "finished issue\n",
	})

	liveReview := filepath.Join(repo, "docs", "specs", "reviews", "pr-102")
	historyPersistRound(t, liveReview, "feature/live", liveHead)

	undecidableReview := filepath.Join(repo, "docs", "specs", "reviews", "pr-103")
	historyPersistRound(t, undecidableReview, "feature/unknown", mainHead)
	historyReplaceRoundHead(t, undecidableReview, mainHead, "")

	legacyFinishedReview := filepath.Join(repo, "docs", "specs", "_reviews", "pr-104")
	historyPersistRound(t, legacyFinishedReview, "feature/merged", mainHead)

	specOwnedReview := filepath.Join(repo, "docs", "specs", "0094-widget", "reviews")
	historyPersistRound(t, specOwnedReview, "feature/merged", mainHead)

	before := historySnapshot(t, repo)
	got, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
	if err != nil {
		t.Fatalf("DiscoverHistoryLayout() error = %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("DiscoverHistoryLayout() collisions = %#v, want none", collisions)
	}
	if len(got) != 3 {
		t.Fatalf("DiscoverHistoryLayout() reported %d relocations, want 3: %#v", len(got), got)
	}
	for _, relocation := range got {
		if relocation.From != "docs/specs/reviews/pr-101/issues/001.md" &&
			!strings.HasPrefix(relocation.From, "docs/specs/reviews/pr-101/round-001/") &&
			!strings.HasPrefix(relocation.From, "docs/specs/_reviews/pr-104/round-001/") {
			t.Fatalf("DiscoverHistoryLayout() relocated non-finished or Spec-owned Review Artifact file %q", relocation.From)
		}
		if strings.HasPrefix(relocation.From, "docs/specs/_reviews/") {
			if !strings.HasPrefix(relocation.To, "docs/history/reviews/pr-104/") {
				t.Fatalf("legacy review relocation destination = %q", relocation.To)
			}
		} else if !strings.HasPrefix(relocation.To, "docs/history/reviews/pr-101/") {
			t.Fatalf("live review relocation destination = %q", relocation.To)
		}
		if relocation.ContentIdentity == "" {
			t.Fatalf("review relocation %#v has no content identity", relocation)
		}
	}
	historyAssertUnchanged(t, repo, before)
}

func TestRetainedReviewReport(t *testing.T) {
	t.Parallel()

	repo, _ := historyGitRepo(t)
	liveHead := historyBranchCommit(t, repo, "feature/live-report")
	gittest.Run(t, repo, "switch", "main")

	liveReview := filepath.Join(repo, "docs", "specs", "reviews", "pr-201")
	historyPersistRound(t, liveReview, "feature/live-report", liveHead)
	undecidableReview := filepath.Join(repo, "docs", "specs", "reviews", "pr-202")
	historyPersistRound(t, undecidableReview, "feature/unknown", liveHead)
	historyReplaceRoundHead(t, undecidableReview, liveHead, "")

	before := historySnapshot(t, repo)
	relocations, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
	if err != nil {
		t.Fatalf("DiscoverHistoryLayout() error = %v", err)
	}
	if len(relocations) != 0 || len(collisions) != 0 {
		t.Fatalf("retained reviews changed layout decisions: relocations=%#v collisions=%#v", relocations, collisions)
	}

	moves, findings, err := planHistoryMoves(context.Background(), repo)
	if err != nil {
		t.Fatalf("planHistoryMoves() error = %v", err)
	}
	if len(moves) != 0 {
		t.Fatalf("retained review report changed move decisions = %#v, want none", moves)
	}
	want := []Finding{
		{
			Code: historyReviewLiveCode,
			Path: "docs/specs/reviews/pr-201",
		},
		{
			Code: historyReviewUndecidableCode,
			Path: "docs/specs/reviews/pr-202",
		},
	}
	if len(findings) != len(want) {
		t.Fatalf("retained review findings = %#v, want %d", findings, len(want))
	}
	for index, finding := range findings {
		if finding.Code != want[index].Code || finding.Path != want[index].Path {
			t.Errorf("retained review finding %d = %#v, want code %q path %q", index, finding, want[index].Code, want[index].Path)
		}
		if !strings.Contains(finding.Message, strings.TrimPrefix(finding.Code, "baseline.history.review.")) ||
			!strings.Contains(finding.Message, "newest Round") &&
				!strings.Contains(finding.Message, "recorded head") {
			t.Errorf("retained review finding %d lacks liveness answer or classifier reason: %#v", index, finding)
		}
	}
	historyAssertUnchanged(t, repo, before)
}

func TestDiscoverHistoryLayoutReportsCollisionWithoutHidingSiblings(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	historyWriteFiles(t, repo, map[string]string{
		"_archived/specs/0001-widget/_prd.md":    "colliding source\n",
		"_archived/specs/0001-widget/task_01.md": "movable sibling\n",
		"docs/history/specs/0001-widget/_prd.md": "occupied destination\n",
	})
	before := historySnapshot(t, repo)

	relocations, collisions, err := DiscoverHistoryLayout(context.Background(), repo)
	if err != nil {
		t.Fatalf("DiscoverHistoryLayout() error = %v", err)
	}
	wantRelocations := []HistoryRelocation{{
		From:            "_archived/specs/0001-widget/task_01.md",
		To:              "docs/history/specs/0001-widget/task_01.md",
		ContentIdentity: historyFixtureIdentity("movable sibling\n"),
	}}
	if !reflect.DeepEqual(relocations, wantRelocations) {
		t.Fatalf("DiscoverHistoryLayout() relocations = %#v, want %#v", relocations, wantRelocations)
	}
	if len(collisions) != 1 {
		t.Fatalf("DiscoverHistoryLayout() collisions = %#v, want one", collisions)
	}
	collision := collisions[0]
	if collision.From != "_archived/specs/0001-widget/_prd.md" ||
		collision.To != "docs/history/specs/0001-widget/_prd.md" ||
		!strings.Contains(collision.Reason, "already exists") {
		t.Fatalf("DiscoverHistoryLayout() collision = %#v, want both paths and occupied-destination reason", collision)
	}
	if collision.ContentIdentity != historyFixtureIdentity("colliding source\n") {
		t.Fatalf("DiscoverHistoryLayout() collision identity = %q, want source identity", collision.ContentIdentity)
	}
	historyAssertUnchanged(t, repo, before)
}

func historyWriteFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()

	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		absolute := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
			t.Fatalf("create fixture parent for %q: %v", path, err)
		}
		if err := os.WriteFile(absolute, []byte(files[path]), 0o644); err != nil {
			t.Fatalf("write fixture %q: %v", path, err)
		}
	}
}

func historyFixtureIdentity(content string) string {
	sum := sha256.Sum256([]byte(content))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func historySnapshot(t *testing.T, root string) map[string]string {
	t.Helper()

	snapshot := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relative)] = historyFixtureIdentity(string(content))
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot fixture repository: %v", err)
	}
	return snapshot
}

func historyAssertUnchanged(t *testing.T, root string, before map[string]string) {
	t.Helper()

	after := historySnapshot(t, root)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("discovery changed fixture repository: before %#v, after %#v", before, after)
	}
}

func historyGitRepo(t *testing.T) (string, string) {
	t.Helper()

	repo := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, repo, "-b", "main")
	historyWriteFiles(t, repo, map[string]string{"tracked.txt": "main\n"})
	gittest.Run(t, repo, "add", "tracked.txt")
	gittest.Run(t, repo, "commit", "-m", "main")
	return repo, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}

func historyBranchCommit(t *testing.T, repo string, branch string) string {
	t.Helper()

	gittest.Run(t, repo, "switch", "-c", branch)
	historyWriteFiles(t, repo, map[string]string{"tracked.txt": branch + "\n"})
	gittest.Run(t, repo, "add", "tracked.txt")
	gittest.Run(t, repo, "commit", "-m", "live review head")
	return strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}

func historyPersistRound(t *testing.T, reviewDir string, branch string, head string) {
	t.Helper()

	_, err := rounds.PersistRound(t.Context(), rounds.PersistRequest{
		ArtifactDir:    filepath.Dir(reviewDir),
		ReviewRoot:     reviewDir,
		Source:         "coderabbit",
		PRNumber:       strings.TrimPrefix(filepath.Base(reviewDir), "pr-"),
		HeadRepository: "example/repository",
		HeadBranch:     branch,
		HeadSHA:        head,
		Round:          1,
	})
	if err != nil {
		t.Fatalf("persist Review Artifact fixture %q: %v", reviewDir, err)
	}
}

func historyReplaceRoundHead(t *testing.T, reviewDir string, oldHead string, newHead string) {
	t.Helper()

	path := filepath.Join(reviewDir, "round-001", "round.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Review Artifact metadata fixture: %v", err)
	}
	replaced := strings.Replace(string(content), "head_sha: "+oldHead, "head_sha: \""+newHead+"\"", 1)
	if replaced == string(content) {
		t.Fatalf("Review Artifact metadata fixture %q has no %q head record", path, oldHead)
	}
	content = []byte(replaced)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write Review Artifact metadata fixture: %v", err)
	}
}
