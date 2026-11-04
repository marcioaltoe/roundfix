// Suite: Claude token-bounded segments through the public review command.
// Boundary IN: disposable Git commits; OUT: prompts, refusals and persisted records.
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/preflight"
)

func TestReviewSegmentsAClaudeCandidateOverItsBudget(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		bytes        []int
		contextLines int
		ends         []int
	}{
		{"packs by tokens rather than bytes", []int{475_000, 475_000, 475_000}, 20, []int{2, 3}},
		{"Spec context changes the first boundary", []int{490_000, 490_000, 490_000}, 1_100, []int{1, 3}},
		{"Spec context triggers segmentation", []int{475_000, 475_000}, 1_100, []int{1, 2}},
		{"later commit fits without first segment context", []int{100_000, 980_000}, 1_100, []int{1, 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
				{result: agent.ExecuteResult{Message: "No findings.", ACPSessionID: "claude-segments"}},
				{result: agent.ExecuteResult{Message: "No findings.", ACPSessionID: "claude-segments"}},
			}}}
			fixture := newReviewCommandFixture(t, "claude", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
			commitScopeFile(t, fixture, "docs/specs/example/_prd.md", "# Example\n\n## Decisions\n\n"+strings.Repeat("Keep this decision.\n", tc.contextLines))
			commitScopeFile(t, fixture, "docs/specs/example/_techspec.md", "# Technical design\n"+strings.Repeat("Keep this design.\n", tc.contextLines*7/11))
			heads := []string{fixture.baseCommit}
			for i, size := range tc.bytes {
				heads = append(heads, commitScopeFile(t, fixture, fmt.Sprintf("part%d.txt", i+1), strings.Repeat("x", size)))
			}
			code, record, stderr := fixture.run(t)
			if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(record.Segments) != 2 {
				t.Fatalf("exit=%d segments=%+v reason=%q stderr=%q", code, record.Segments, record.Reason, stderr)
			}
			if runner.prepareCalls != 1 || len(runner.prepared) != 1 || len(runner.requests) != 2 || runner.endCalls != 1 {
				t.Fatalf("prepares=%d prompts=%d ends=%d", runner.prepareCalls, len(runner.requests), runner.endCalls)
			}
			var persisted []struct {
				EstimatedPromptTokens int `json:"estimatedPromptTokens"`
			}
			if err := json.Unmarshal(reviewPromptBoundRecordFields(t, fixture)["segments"], &persisted); err != nil {
				t.Fatal(err)
			}
			if len(persisted) != 2 {
				t.Fatalf("persisted segments=%+v", persisted)
			}
			totalBytes, totalTokens, start := 0, 0, 0
			for i, request := range runner.requests {
				end := tc.ends[i]
				segment := record.Segments[i]
				header := fmt.Sprintf("Review Segment %d of 2: commits %s..%s of the candidate %s..%s; answer for this segment only.\n", i+1, heads[start], heads[end], heads[0], heads[len(heads)-1])
				if !strings.HasPrefix(request.Prompt, header) || request.Session != runner.prepared[0].Session || request.Access != agent.SessionAccessReadOnly || request.Runtime.ID != "claude" {
					t.Fatalf("segment %d=%+v want %s..%s; prompt tokens=%d runtime=%s session match=%t", i+1, segment, heads[start], heads[end], estimateReviewPromptTokens(request.Prompt), request.Runtime.ID, request.Session == runner.prepared[0].Session)
				}
				if strings.Contains(request.Prompt, "--- BEGIN SPEC CONTEXT: example ---") != (i == 0) {
					t.Fatalf("Spec context on segment %d", i+1)
				}
				diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", heads[start], heads[end], "--")
				if err != nil {
					t.Fatal(err)
				}
				tokens := estimateReviewPromptTokens(request.Prompt)
				if segment.Base != heads[start] || segment.Head != heads[end] || segment.DiffBytes != len(diff) || tokens > claudeReviewPromptBudget || persisted[i].EstimatedPromptTokens != tokens {
					t.Fatalf("segment=%+v tokens=%d persisted=%d", segment, tokens, persisted[i].EstimatedPromptTokens)
				}
				if tc.name == "packs by tokens rather than bytes" && i == 0 && segment.DiffBytes <= reviewDiffBound {
					t.Fatal("first segment must exceed the Codex byte bound")
				}
				totalBytes += len(diff)
				totalTokens += tokens
				start = end
			}
			if record.DiffBytes != totalBytes || record.EstimatedPromptTokens != totalTokens || record.Lineage.Round != 1 || len(record.Lineage.ACPSessionIDs) != 1 || record.Lineage.ACPSessionIDs[0] != "claude-segments" {
				t.Fatalf("totals or lineage: %+v", record)
			}
		})
	}
}

func TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
	for i := 1; i <= 5; i++ {
		commitScopeFile(t, fixture, fmt.Sprintf("part%d.txt", i), strings.Repeat("x", 700_000))
	}
	code, record, stderr := fixture.run(t)
	want := fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) needs 5 review segments; a review allows at most %d", record.EstimatedPromptTokens, len(record.OmittedPaths), reviewSegmentLimit)
	assertClaudeSegmentRefusal(t, code, record, stderr, want)
	if record.EstimatedPromptTokens <= claudeReviewPromptBudget || runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 || runner.endCalls != 0 {
		t.Fatalf("estimate=%d probes=%d prepares=%d prompts=%d ends=%d", record.EstimatedPromptTokens, runner.probeCalls, runner.prepareCalls, runner.preparedCalls, runner.endCalls)
	}
}

func TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"round two", "oversized commit", "oversized commit hidden by range"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: sessionReviewFinding}}}}}
			fixture := newReviewCommandFixture(t, "claude", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
			if name == "round two" {
				code, _, stderr := fixture.run(t)
				if code != exitRunFailed {
					t.Fatalf("seed exit=%d stderr=%q", code, stderr)
				}
			}
			probes, prepares, prompts, ends := runner.probeCalls, runner.prepareCalls, runner.preparedCalls, runner.endCalls
			switch name {
			case "round two":
				for i := 1; i <= 3; i++ {
					commitScopeFile(t, fixture, fmt.Sprintf("part%d.txt", i), strings.Repeat("x", 475_000))
				}
			case "oversized commit":
				commitScopeFile(t, fixture, "large.txt", strings.Repeat("x", claudeReviewContextTokens))
			case "oversized commit hidden by range":
				commitScopeFile(t, fixture, "rewrite.txt", strings.Repeat("x", 600_000))
				commitScopeFile(t, fixture, "rewrite.txt", strings.Repeat("y", 600_000))
				commitScopeFile(t, fixture, "next.txt", strings.Repeat("z", 600_000))
			}
			code, record, stderr := fixture.run(t)
			want := fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) exceeds the claude review budget of %d tokens, half of its %d-token context window", record.EstimatedPromptTokens, len(record.OmittedPaths), claudeReviewPromptBudget, claudeReviewContextTokens)
			assertClaudeSegmentRefusal(t, code, record, stderr, want)
			if len(record.Segments) != 0 || record.AnswerPath != "" || record.EstimatedPromptTokens <= claudeReviewPromptBudget || (name == "round two" && record.Lineage.Round != 2) {
				t.Fatalf("refusal record=%+v", record)
			}
			if runner.probeCalls != probes || runner.prepareCalls != prepares || runner.preparedCalls != prompts || runner.endCalls != ends {
				t.Fatal("size refusal reached the provider")
			}
		})
	}
}

func TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "No findings."}},
		{result: agent.ExecuteResult{Message: "No findings."}},
		{result: agent.ExecuteResult{Message: "No findings."}},
	}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	heads := commitReviewSegments(t, fixture, 3)
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(record.Segments) != 3 || record.EstimatedPromptTokens != 0 {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	fields := reviewPromptBoundRecordFields(t, fixture)
	if _, ok := fields["estimatedPromptTokens"]; ok {
		t.Fatal("Codex record contains a token estimate")
	}
	var segments []map[string]json.RawMessage
	if err := json.Unmarshal(fields["segments"], &segments); err != nil {
		t.Fatal(err)
	}
	if len(segments) != 3 {
		t.Fatalf("persisted segments=%d", len(segments))
	}
	for i, fields := range segments {
		if _, ok := fields["estimatedPromptTokens"]; ok {
			t.Fatalf("Codex segment %d contains a token estimate", i+1)
		}
		diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", heads[i], heads[i+1], "--")
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf(`{"base":%q,"head":%q,"diffBytes":%d}`, heads[i], heads[i+1], len(diff))
		got, err := json.Marshal(record.Segments[i])
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("Codex segment=%s, want %s", got, want)
		}
	}
}

func assertClaudeSegmentRefusal(t *testing.T, code int, record reviewRecord, stderr, reason string) {
	t.Helper()
	if record.Lineage.Round == 1 {
		diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), record.Repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", record.BaseCommit, record.HeadCommit, "--")
		if err != nil {
			t.Fatal(err)
		}
		want := estimateReviewPromptTokens(buildReviewPrompt(record.BaseCommit, record.HeadCommit, diff, record.OmittedPaths))
		if record.EstimatedPromptTokens != want {
			t.Fatalf("refusal estimate=%d, whole unsegmented prompt estimate=%d", record.EstimatedPromptTokens, want)
		}
	}
	if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || !strings.HasPrefix(record.Reason, reason) || !strings.HasPrefix(stderr, "roundfix: review blocked: "+reason) || !strings.HasSuffix(stderr, "\n") {
		t.Fatalf("exit=%d reason=%q stderr=%q, want prefix %q", code, record.Reason, stderr, reason)
	}
}
