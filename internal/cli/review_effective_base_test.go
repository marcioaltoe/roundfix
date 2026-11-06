// Suite: Effective Base measurement through Git and the public review command.
// Invariant: segments hold only the candidate's own change; conflicting carries keep their start.
// Boundary IN: disposable Git repositories and a fake reviewer.
// Boundary OUT: measured objects, segment records and prompts for both providers.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange(t *testing.T) {
	t.Parallel()
	fixture := newReviewCommandFixture(t, "codex", &segmentReviewRunner{})
	gittest.Run(t, fixture.repository, "checkout", "main")
	base := commitScopeFile(t, fixture, "upstream.txt", "main moved on\n")
	gittest.Run(t, fixture.repository, "checkout", "feature/review")
	gittest.Run(t, fixture.repository, "merge", "--no-edit", "main")
	merged := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	later := commitScopeFile(t, fixture, "later.txt", "own later change\n")
	gittest.Run(t, fixture.repository, "checkout", "-b", "test/conflict", fixture.headCommit)
	conflicting := commitScopeFile(t, fixture, "upstream.txt", "candidate wrote this differently\n")

	for _, tc := range []struct {
		name, from, to, want string
		carried              bool
	}{
		{"same Effective Base", merged, later, merged, true},
		{"start on the base", base, fixture.headCommit, fixture.baseCommit, true},
		{"merge of the default branch", fixture.headCommit, merged, strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", merged+"^{tree}")), true},
		{"conflicting carry", conflicting, merged, conflicting, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			measured, carried, err := reviewMeasuredFrom(t.Context(), preflight.ExecGitRunner{}, fixture.repository, base, tc.from, tc.to)
			if err != nil {
				t.Fatal(err)
			}
			if measured != tc.want || carried != tc.carried {
				t.Fatalf("measured=%s carried=%t, want %s %t", measured, carried, tc.want, tc.carried)
			}
			if tc.name == "merge of the default branch" {
				if kind := strings.TrimSpace(gittest.Run(t, fixture.repository, "cat-file", "-t", measured)); kind != "tree" {
					t.Fatalf("carried object kind=%s, want tree", kind)
				}
				if diff := gittest.Run(t, fixture.repository, "diff", measured, merged); diff != "" {
					t.Fatalf("carried tree differs from merge: %s", diff)
				}
			}
		})
	}
}

func TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase(t *testing.T) {
	t.Parallel()
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			t.Parallel()
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
				{result: agent.ExecuteResult{Message: "No findings.", ACPSessionID: "effective-base-segments"}},
				{result: agent.ExecuteResult{Message: "No findings.", ACPSessionID: "effective-base-segments"}},
				{result: agent.ExecuteResult{Message: "No findings.", ACPSessionID: "effective-base-segments"}},
			}}}
			fixture := newReviewCommandFixture(t, provider, runner)
			unit := reviewDiffBound
			if provider == "claude" {
				unit = 2 * claudeReviewPromptBudget
				writeReviewCommandProfileConfig(t, fixture.repository, provider, fixture.artifactDir, provider, provider)
			}
			item := commitScopeFile(t, fixture, "item.txt", strings.Repeat("x", unit*2/3)+"\n")
			commitScopeFile(t, fixture, "item2.txt", strings.Repeat("y", unit*2/3)+"\n")
			gittest.Run(t, fixture.repository, "checkout", "main")
			base := commitScopeFile(t, fixture, "mainwork.txt", strings.Repeat("m", unit+unit/5)+"\n")
			gittest.Run(t, fixture.repository, "checkout", "feature/review")
			gittest.Run(t, fixture.repository, "merge", "--no-edit", "main")
			merged := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			head := commitScopeFile(t, fixture, "part1.txt", strings.Repeat("z", unit*2/3)+"\n")

			code, record, stderr := runDefaultBranchReview(t)
			if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(record.Segments) != 3 {
				t.Fatalf("exit=%d segments=%+v reason=%q stderr=%q", code, record.Segments, record.Reason, stderr)
			}
			if record.BaseCommit != base || record.HeadCommit != head || len(runner.requests) != 3 {
				t.Fatalf("base=%s head=%s prompts=%d", record.BaseCommit, record.HeadCommit, len(runner.requests))
			}
			// The second segment starts at item, carrying only that prefix onto main.
			// Its expected tree contains item.txt and mainwork.txt, but not item2.txt.
			expectedTree := strings.TrimSpace(gittest.Run(t, fixture.repository, "merge-tree", "--write-tree", item, base))
			if kind := strings.TrimSpace(gittest.Run(t, fixture.repository, "cat-file", "-t", expectedTree)); kind != "tree" {
				t.Fatalf("second segment's measured origin kind=%s, want tree", kind)
			}
			wantBases := []string{base, item, merged}
			wantHeads := []string{item, merged, head}
			wantMeasured := []string{fixture.baseCommit, expectedTree, ""}
			for i, segment := range record.Segments {
				if segment.Base != wantBases[i] || segment.Head != wantHeads[i] || segment.MeasuredFrom != wantMeasured[i] {
					t.Fatalf("segment %d=%+v, want base=%s head=%s measuredFrom=%s", i+1, segment, wantBases[i], wantHeads[i], wantMeasured[i])
				}
				prompt := runner.requests[i].Prompt
				first := fmt.Sprintf("Review Segment %d of 3: commits %s..%s of the candidate %s..%s; answer for this segment only.\n", i+1, segment.Base, segment.Head, base, head)
				if !strings.HasPrefix(prompt, first) {
					t.Fatalf("segment %d has the wrong first header line", i+1)
				}
				second := fmt.Sprintf("This segment's diff starts from %s instead of %s, so it leaves out what the default branch changed.\n", wantMeasured[i], wantBases[i])
				if wantMeasured[i] != "" {
					if !strings.HasPrefix(prompt, first+second) || strings.Count(prompt, second) != 1 {
						t.Fatalf("segment %d lacks its exact measurement header", i+1)
					}
				} else if strings.Contains(prompt, "This segment's diff starts from ") {
					t.Fatalf("segment %d has an unexpected measurement header", i+1)
				}
				if strings.Contains(prompt, strings.Repeat("m", 1000)) {
					t.Fatalf("segment %d includes the default branch's own work", i+1)
				}
			}
		})
	}
}

func runDefaultBranchReview(t *testing.T) (int, reviewRecord, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", "main"}, &stdout, &stderr)
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode %q: %v (stderr %q)", stdout.String(), err, stderr.String())
	}
	return code, record, stderr.String()
}
