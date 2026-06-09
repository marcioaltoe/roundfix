package rounds

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
)

func TestPersistRoundWritesReviewIssueArtifacts(t *testing.T) {
	createdAt := time.Date(2026, 6, 9, 12, 30, 0, 0, time.UTC)
	artifactDir := t.TempDir()

	result, err := PersistRound(context.Background(), PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		HeadSHA:        "abc123",
		CreatedAt:      createdAt,
		Items: []reviewsource.ReviewItem{
			{
				Title:                   "major: handle nil cache",
				File:                    "internal/cache/cache.go",
				Line:                    42,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "Reviewer text is untrusted: `$(rm -rf /)` must stay literal.",
				SourceRef:               "thread:PRRT_1,comment:PRRC_1",
				ReviewHash:              "abc",
				SourceReviewID:          "9001",
				SourceReviewSubmittedAt: createdAt,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected round persistence, got %v", err)
	}

	if result.Round != 1 {
		t.Fatalf("expected round 1, got %d", result.Round)
	}
	if len(result.IssuePaths) != 1 {
		t.Fatalf("expected one issue path, got %#v", result.IssuePaths)
	}
	if result.RoundDir != filepath.Join(artifactDir, "reviews", "pr-123", "round-001") {
		t.Fatalf("unexpected round dir %q", result.RoundDir)
	}

	issue := readFile(t, result.IssuePaths[0])
	for _, expected := range []string{
		"source: coderabbit",
		"pr: \"123\"",
		"round: 1",
		"round_created_at: \"2026-06-09T12:30:00Z\"",
		"status: pending",
		"head_repository: owner/project",
		"head_branch: feature/review",
		"head_sha: abc123",
		"file: internal/cache/cache.go",
		"line: 42",
		"severity: major",
		"author: coderabbitai[bot]",
		"source_ref: thread:PRRT_1,comment:PRRC_1",
		"review_hash: abc",
		"duplicate_of: \"\"",
		"source_review_id: \"9001\"",
		"source_review_submitted_at: \"2026-06-09T12:30:00Z\"",
		"Reviewer text is untrusted: `$(rm -rf /)` must stay literal.",
		"- Decision: `UNREVIEWED`",
	} {
		if !strings.Contains(issue, expected) {
			t.Fatalf("expected issue artifact to contain %q, got:\n%s", expected, issue)
		}
	}

	round := readFile(t, filepath.Join(result.RoundDir, "round.md"))
	if !strings.Contains(round, "issue_count: 1") {
		t.Fatalf("expected round metadata issue count, got:\n%s", round)
	}
}

func TestPersistRoundUsesNextRoundNumber(t *testing.T) {
	artifactDir := t.TempDir()
	mustMkdir(t, filepath.Join(artifactDir, "reviews", "pr-123", "round-001"))

	result, err := PersistRound(context.Background(), PersistRequest{
		ArtifactDir:    artifactDir,
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		HeadSHA:        "abc123",
		CreatedAt:      time.Date(2026, 6, 9, 12, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected second round persistence, got %v", err)
	}
	if result.Round != 2 {
		t.Fatalf("expected round 2, got %d", result.Round)
	}
	if _, err := os.Stat(filepath.Join(artifactDir, "reviews", "pr-123", "round-002", "round.md")); err != nil {
		t.Fatalf("expected round metadata: %v", err)
	}
}

func TestAllowedStatusRejectsOutOfScopeStatuses(t *testing.T) {
	if !AllowedStatus(StatusDuplicated) {
		t.Fatal("expected duplicated to be an allowed terminal status")
	}
	if AllowedStatus("postponed") {
		t.Fatal("did not expect postponed to be an MVP status")
	}
}

func TestSelectCompatibleIssuesFindsAllUnresolvedCompatibleRounds(t *testing.T) {
	artifactDir := t.TempDir()
	createdAt := time.Date(2026, 6, 9, 12, 30, 0, 0, time.UTC)
	first := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          1,
		CreatedAt:      createdAt,
	})
	second := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          2,
		CreatedAt:      createdAt.Add(time.Minute),
	})
	replaceStatus(t, second.IssuePaths[0], StatusValid)
	wrongBranch := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "other-branch",
		Round:          3,
		CreatedAt:      createdAt.Add(2 * time.Minute),
	})
	terminal := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          4,
		CreatedAt:      createdAt.Add(3 * time.Minute),
	})
	replaceStatus(t, terminal.IssuePaths[0], StatusResolved)

	result, err := SelectCompatibleIssues(context.Background(), SelectRequest{
		ArtifactDir:    artifactDir,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
	})
	if err != nil {
		t.Fatalf("expected compatible issues, got %v", err)
	}

	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 unresolved compatible issues, got %#v", result.Issues)
	}
	if result.Issues[0].Path != first.IssuePaths[0] || result.Issues[1].Path != second.IssuePaths[0] {
		t.Fatalf("expected round 1 and 2 issue paths, got %#v", result.Issues)
	}
	if strings.Contains(result.Issues[0].Path, wrongBranch.RoundDir) {
		t.Fatalf("did not expect wrong branch artifact path in result: %#v", result.Issues)
	}
	if len(result.Rounds) != 2 || result.Rounds[0] != 1 || result.Rounds[1] != 2 {
		t.Fatalf("expected selected rounds [1 2], got %#v", result.Rounds)
	}
}

func TestSelectCompatibleIssuesHonorsRoundSelector(t *testing.T) {
	artifactDir := t.TempDir()
	createdAt := time.Date(2026, 6, 9, 12, 30, 0, 0, time.UTC)
	persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          1,
		CreatedAt:      createdAt,
	})
	second := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          2,
		CreatedAt:      createdAt.Add(time.Minute),
	})

	result, err := SelectCompatibleIssues(context.Background(), SelectRequest{
		ArtifactDir:    artifactDir,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          2,
	})
	if err != nil {
		t.Fatalf("expected round-filtered compatible issue, got %v", err)
	}

	if len(result.Issues) != 1 || result.Issues[0].Path != second.IssuePaths[0] {
		t.Fatalf("expected only round 2 issue, got %#v", result.Issues)
	}
	if len(result.Rounds) != 1 || result.Rounds[0] != 2 {
		t.Fatalf("expected selected rounds [2], got %#v", result.Rounds)
	}
}

func TestSelectCompatibleIssuesReportsMissingCompatibleArtifacts(t *testing.T) {
	artifactDir := t.TempDir()
	persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "999",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          1,
		CreatedAt:      time.Date(2026, 6, 9, 12, 30, 0, 0, time.UTC),
	})

	_, err := SelectCompatibleIssues(context.Background(), SelectRequest{
		ArtifactDir:    artifactDir,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
	})

	var noArtifacts NoCompatibleArtifactsError
	if !errors.As(err, &noArtifacts) {
		t.Fatalf("expected NoCompatibleArtifactsError, got %T %v", err, err)
	}
	if !strings.Contains(err.Error(), "roundfix fetch --source coderabbit --pr 123") {
		t.Fatalf("expected fetch guidance, got %q", err.Error())
	}
}

func TestReviewIssueFingerprintPrefersSourceRefAndFallsBackToReviewHash(t *testing.T) {
	withSource, err := ReviewIssueFingerprint(Issue{
		Path:       "issue_001.md",
		SourceRef:  "thread:PRRT_1,comment:PRRC_1",
		ReviewHash: "same-hash",
	})
	if err != nil {
		t.Fatalf("expected source ref fingerprint, got %v", err)
	}
	if withSource != "source_ref:thread:PRRT_1,comment:PRRC_1" {
		t.Fatalf("expected source ref fingerprint, got %q", withSource)
	}

	withHash, err := ReviewIssueFingerprint(Issue{
		Path:       "issue_002.md",
		ReviewHash: "same-hash",
	})
	if err != nil {
		t.Fatalf("expected review hash fingerprint, got %v", err)
	}
	if withHash != "review_hash:same-hash" {
		t.Fatalf("expected review hash fingerprint, got %q", withHash)
	}
}

func TestPlanBatchesDeduplicatesNewestOccurrences(t *testing.T) {
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	issues := []Issue{
		{
			Path:                    "round-001/issue_001.md",
			Round:                   1,
			RoundCreatedAt:          base,
			SourceRef:               "thread:dup,comment:old",
			ReviewHash:              "ignored-hash",
			SourceReviewSubmittedAt: base.Add(10 * time.Minute),
		},
		{
			Path:                    "round-002/issue_001.md",
			Round:                   2,
			RoundCreatedAt:          base.Add(time.Minute),
			SourceRef:               "thread:dup,comment:old",
			ReviewHash:              "ignored-hash",
			SourceReviewSubmittedAt: base,
		},
		{
			Path:                    "round-001/issue_002.md",
			Round:                   1,
			RoundCreatedAt:          base,
			ReviewHash:              "hash-only",
			SourceReviewSubmittedAt: base,
		},
		{
			Path:                    "round-001/issue_003.md",
			Round:                   1,
			RoundCreatedAt:          base,
			ReviewHash:              "hash-only",
			SourceReviewSubmittedAt: base.Add(time.Minute),
		},
		{
			Path:                    "round-001/issue_004.md",
			Round:                   1,
			RoundCreatedAt:          base,
			SourceRef:               "thread:a,comment:1",
			ReviewHash:              "shared-hash",
			SourceReviewSubmittedAt: base,
		},
		{
			Path:                    "round-001/issue_005.md",
			Round:                   1,
			RoundCreatedAt:          base,
			SourceRef:               "thread:b,comment:1",
			ReviewHash:              "shared-hash",
			SourceReviewSubmittedAt: base,
		},
	}

	plan, err := PlanBatches(BatchRequest{Issues: issues, BatchSize: 2})
	if err != nil {
		t.Fatalf("expected Batch plan, got %v", err)
	}

	if len(plan.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicate associations, got %#v", plan.Duplicates)
	}
	if len(plan.Batches) != 2 {
		t.Fatalf("expected 2 batches, got %#v", plan.Batches)
	}
	assigned := assignedPaths(plan.Batches)
	for _, unexpected := range []string{"round-001/issue_001.md", "round-001/issue_002.md"} {
		if containsString(assigned, unexpected) {
			t.Fatalf("did not expect older duplicate %q in assigned issues %#v", unexpected, assigned)
		}
	}
	for _, expected := range []string{"round-002/issue_001.md", "round-001/issue_003.md", "round-001/issue_004.md", "round-001/issue_005.md"} {
		if !containsString(assigned, expected) {
			t.Fatalf("expected assigned issue %q in %#v", expected, assigned)
		}
	}
}

func TestPlanBatchesFailsAmbiguousNewestDuplicate(t *testing.T) {
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	_, err := PlanBatches(BatchRequest{
		BatchSize: 1,
		Issues: []Issue{
			{
				Path:                    "round-001/issue_001.md",
				Round:                   1,
				RoundCreatedAt:          base,
				SourceRef:               "thread:dup,comment:1",
				SourceReviewSubmittedAt: base,
			},
			{
				Path:                    "round-001/issue_002.md",
				Round:                   1,
				RoundCreatedAt:          base,
				SourceRef:               "thread:dup,comment:1",
				SourceReviewSubmittedAt: base,
			},
		},
	})

	var ambiguous AmbiguousDuplicateError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("expected AmbiguousDuplicateError, got %T %v", err, err)
	}
}

func TestMarkDuplicatedAfterTerminalDefersUntilNewestResolvedOrInvalid(t *testing.T) {
	artifactDir := t.TempDir()
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	first := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          1,
		CreatedAt:      base,
		Items: []reviewsource.ReviewItem{
			reviewItem("thread:dup,comment:1", "hash-a", base),
		},
	})
	second := persistTestRound(t, artifactDir, PersistRequest{
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		Round:          2,
		CreatedAt:      base.Add(time.Minute),
		Items: []reviewsource.ReviewItem{
			reviewItem("thread:dup,comment:1", "hash-a", base.Add(time.Minute)),
		},
	})
	selection, err := SelectCompatibleIssues(context.Background(), SelectRequest{
		ArtifactDir:    artifactDir,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
	})
	if err != nil {
		t.Fatalf("select compatible issues: %v", err)
	}
	plan, err := PlanBatches(BatchRequest{Issues: selection.Issues, BatchSize: 1})
	if err != nil {
		t.Fatalf("plan batches: %v", err)
	}
	if len(plan.Duplicates) != 1 {
		t.Fatalf("expected one duplicate association, got %#v", plan.Duplicates)
	}

	marked, err := MarkDuplicatedAfterTerminal(context.Background(), plan.Duplicates)
	if err != nil {
		t.Fatalf("mark before terminal: %v", err)
	}
	if marked != 0 {
		t.Fatalf("expected no duplicate mark before newest terminal, got %d", marked)
	}
	oldBefore, err := ParseIssue(first.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse old issue before: %v", err)
	}
	if oldBefore.Status != StatusPending || oldBefore.DuplicateOf != "" {
		t.Fatalf("expected old issue to stay pending before newest terminal, got %#v", oldBefore)
	}

	if err := SetIssueStatus(second.IssuePaths[0], StatusResolved, ""); err != nil {
		t.Fatalf("mark newest resolved: %v", err)
	}
	marked, err = MarkDuplicatedAfterTerminal(context.Background(), plan.Duplicates)
	if err != nil {
		t.Fatalf("mark after terminal: %v", err)
	}
	if marked != 1 {
		t.Fatalf("expected one duplicate mark after newest terminal, got %d", marked)
	}
	oldAfter, err := ParseIssue(first.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse old issue after: %v", err)
	}
	if oldAfter.Status != StatusDuplicated {
		t.Fatalf("expected old issue duplicated, got %q", oldAfter.Status)
	}
	if oldAfter.DuplicateOf != second.IssuePaths[0] {
		t.Fatalf("expected duplicate_of newest path %q, got %q", second.IssuePaths[0], oldAfter.DuplicateOf)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assignedPaths(batches []Batch) []string {
	paths := []string{}
	for _, batch := range batches {
		for _, issue := range batch.Issues {
			paths = append(paths, issue.Path)
		}
	}
	return paths
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func persistTestRound(t *testing.T, artifactDir string, req PersistRequest) PersistResult {
	t.Helper()
	if req.Source == "" {
		req.Source = reviewsource.SourceCodeRabbit
	}
	if req.HeadSHA == "" {
		req.HeadSHA = "abc123"
	}
	if len(req.Items) == 0 {
		req.Items = []reviewsource.ReviewItem{
			{
				Title:                   "major: handle nil cache",
				File:                    "internal/cache/cache.go",
				Line:                    42,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "review body",
				SourceRef:               "thread:PRRT_1,comment:PRRC_1",
				ReviewHash:              "abc",
				SourceReviewID:          "9001",
				SourceReviewSubmittedAt: req.CreatedAt,
			},
		}
	}
	req.ArtifactDir = artifactDir
	result, err := PersistRound(context.Background(), req)
	if err != nil {
		t.Fatalf("persist test round: %v", err)
	}
	return result
}

func reviewItem(sourceRef string, reviewHash string, submittedAt time.Time) reviewsource.ReviewItem {
	return reviewsource.ReviewItem{
		Title:                   "major: handle nil cache",
		File:                    "internal/cache/cache.go",
		Line:                    42,
		Severity:                "major",
		Author:                  "coderabbitai[bot]",
		Body:                    "review body",
		SourceRef:               sourceRef,
		ReviewHash:              reviewHash,
		SourceReviewID:          "9001",
		SourceReviewSubmittedAt: submittedAt,
	}
}

func replaceStatus(t *testing.T, path string, status string) {
	t.Helper()
	content := readFile(t, path)
	content = strings.Replace(content, "status: pending", "status: "+status, 1)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write updated status %s: %v", path, err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}
