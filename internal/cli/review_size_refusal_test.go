// Boundary: size refusals through the public review command and disposable Git history.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewSizeRefusalNamesTheLargestFiles(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, provider                         string
		round                                  int
		correction, segments, twoFiles, hidden bool
	}{
		{name: "Codex single commit", provider: "codex", round: 1},
		{name: "Claude single commit and round one budget", provider: "claude", round: 1},
		{name: "Codex segment limit", provider: "codex", round: 1, segments: true},
		{name: "Claude segment limit", provider: "claude", round: 1, segments: true},
		{name: "Codex round two byte bound", provider: "codex", round: 2},
		{name: "Claude round two budget", provider: "claude", round: 2},
		{name: "Codex Correction Review byte bound", provider: "codex", round: 2, correction: true},
		{name: "Claude Correction Review byte bound", provider: "claude", round: 2, correction: true},
		{name: "two files", provider: "codex", round: 1, twoFiles: true},
		{name: "Codex oversized commit hidden by range", provider: "codex", round: 1, hidden: true},
		{name: "Claude oversized commit hidden by range", provider: "claude", round: 1, hidden: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{
				{result: agent.ExecuteResult{Message: sessionReviewFinding}},
				{result: agent.ExecuteResult{Message: sessionReviewFinding}},
			}}
			fixture := newReviewCommandFixture(t, tc.provider, runner)
			writeReviewCommandProfileConfig(t, fixture.repository, tc.provider, fixture.artifactDir, tc.provider, tc.provider)
			// Pin the omission in the declaration base and leave no config delta.
			configBytes, err := os.ReadFile(filepath.Join(fixture.repository, ".roundfixrc.yml"))
			if err != nil {
				t.Fatal(err)
			}
			config := strings.Replace(string(configBytes), "  provider: "+tc.provider+"\n", "  provider: "+tc.provider+"\n  omit_paths: ['omitted.txt']\n", 1)
			fixture.baseCommit = commitScopeFile(t, fixture, ".roundfixrc.yml", config)
			// Keep the fixture's finding anchor in the candidate after the new base.
			commitScopeFile(t, fixture, "review.txt", "first reviewed line\nsecond reviewed line\n")
			measuredBase := fixture.baseCommit
			if tc.round == 2 {
				commitScopeFile(t, fixture, "old.txt", strings.Repeat("o", 800_000))
				code, prior, stderr := fixture.run(t)
				if code != exitRunFailed {
					t.Fatalf("seed round one exit=%d reason=%q stderr=%q", code, prior.Reason, stderr)
				}
				measuredBase = prior.HeadCommit
				if tc.correction {
					commitScopeFile(t, fixture, "review.txt", "first corrected line\nsecond corrected line\n")
					code, prior, stderr = fixture.run(t)
					if code != exitRunFailed || prior.Lineage.Round != 2 {
						t.Fatalf("seed round two exit=%d reason=%q stderr=%q", code, prior.Reason, stderr)
					}
					measuredBase = prior.HeadCommit
				}
			}
			probes, prepares, prompts := runner.probeCalls, runner.prepareCalls, runner.preparedCalls
			type payload struct {
				path string
				size int
			}
			files := []payload{
				{"zzz.txt", 300_000}, {"tie-b.txt", 350_000},
				{"dados gerados/ação.json", 450_000}, {"tie-a.txt", 350_000}, {"aaa.txt", 200_000},
			}
			if tc.segments {
				files = []payload{
					{"zzz.txt", 650_000}, {"tie-b.txt", 670_000},
					{"dados gerados/ação.json", 700_000}, {"tie-a.txt", 670_000}, {"aaa.txt", 660_000},
				}
			}
			if tc.twoFiles {
				files = []payload{{"zzz.txt", 500_000}, {"dados gerados/ação.json", 600_000}}
				// This case has exactly two reviewed sections.
				fixture.baseCommit = commitScopeFile(t, fixture, "review.txt", "two-file base\n")
				measuredBase = fixture.baseCommit
			}
			for i, file := range files {
				if err := os.MkdirAll(filepath.Dir(filepath.Join(fixture.repository, file.path)), 0o755); err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(fixture.repository, file.path), strings.Repeat(string(rune('a'+i)), file.size))
				gittest.Run(t, fixture.repository, "add", "--", file.path)
				if tc.segments || tc.round == 2 || tc.hidden {
					gittest.Run(t, fixture.repository, "commit", "-m", "add "+file.path)
				}
			}
			if tc.hidden {
				// The rewrite alone is over both providers' bounds, while the
				// refusal must describe the whole candidate's smaller section.
				commitScopeFile(t, fixture, "dados gerados/ação.json", strings.Repeat("r", 600_000))
			}
			commitScopeFile(t, fixture, "omitted.txt", strings.Repeat("g", 2_000_000))
			head := gittest.Run(t, fixture.repository, "rev-parse", "HEAD")
			if tc.correction {
				disposeCorrectionFinding(t, "F1", "--fixed-by", strings.TrimSpace(head))
			}
			code, record, stderr := fixture.run(t)
			reason := fmt.Sprintf("review diff too large: %d bytes after omitting %d path(s) exceeds the review bound of %d bytes", record.DiffBytes, len(record.OmittedPaths), reviewDiffBound)
			if tc.provider == "claude" && !tc.correction {
				reason = fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) exceeds the claude review budget of %d tokens, half of its %d-token context window", record.EstimatedPromptTokens, len(record.OmittedPaths), claudeReviewPromptBudget, claudeReviewContextTokens)
			}
			if tc.segments {
				if tc.provider == "claude" {
					reason = fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) needs 5 review segments; a review allows at most %d", record.EstimatedPromptTokens, len(record.OmittedPaths), reviewSegmentLimit)
				} else {
					reason = fmt.Sprintf("review diff too large: %d bytes after omitting %d path(s) needs 5 review segments; a review allows at most %d", record.DiffBytes, len(record.OmittedPaths), reviewSegmentLimit)
				}
			}
			reason += reviewSizeRefusalTestSuffix(t, fixture, measuredBase, record.HeadCommit, record.OmittedPaths)
			if code != 2 || record.Outcome != reviewOutcomeBlocked || record.Reason != reason || stderr != "roundfix: review blocked: "+reason+"\n" {
				t.Fatalf("exit=%d reason=%q stderr=%q; want %q", code, record.Reason, stderr, reason)
			}
			if runner.probeCalls != probes || runner.prepareCalls != prepares || runner.preparedCalls != prompts {
				t.Fatal("refusal reached the provider")
			}
			if strings.Contains(record.Reason, "omitted.txt") || strings.Contains(record.Reason, "old.txt") {
				t.Fatalf("refusal lists a file outside the measured scope: %q", record.Reason)
			}
			if tc.round == 2 && (record.Lineage.Round != 2 || record.Lineage.Correction != tc.correction) {
				t.Fatalf("wrong refusal lineage: %+v", record.Lineage)
			}
		})
	}
}

// The oracle measures each path separately rather than splitting the combined patch.
func reviewSizeRefusalTestSuffix(t *testing.T, fixture reviewCommandFixture, base, head string, omitted []reviewOmittedPath) string {
	t.Helper()
	git := preflight.ExecGitRunner{}
	names, err := git.RunGit(t.Context(), fixture.repository, "diff", "--name-only", "-z", base, head, "--")
	if err != nil {
		t.Fatal(err)
	}
	excluded := make(map[string]bool)
	for _, path := range omitted {
		excluded[path.Path] = true
	}
	type fileSize struct {
		path  string
		bytes int
	}
	var files []fileSize
	var paths []string
	for _, path := range strings.Split(names, "\x00") {
		if path != "" && !excluded[path] {
			paths = append(paths, path)
		}
	}
	for i, path := range paths {
		diff := gittest.Run(t, fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", base, head, "--", ":(literal)"+path)
		// ExecGitRunner trims only the end of the combined patch. Earlier
		// sections retain their final newline; measure each raw Git output.
		if i == len(paths)-1 {
			diff = strings.TrimRight(diff, "\n")
		}
		files = append(files, fileSize{path, len(diff)})
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].bytes != files[j].bytes {
			return files[i].bytes > files[j].bytes
		}
		return files[i].path < files[j].path
	})
	var descriptions []string
	for _, file := range files[:min(3, len(files))] {
		descriptions = append(descriptions, fmt.Sprintf("%s (%d bytes)", file.path, file.bytes))
	}
	return "; largest files in the diff: " + strings.Join(descriptions, ", ") + "; pre_pr_review.omit_paths in the base commit's Project Config leaves generated files out of the review"
}
