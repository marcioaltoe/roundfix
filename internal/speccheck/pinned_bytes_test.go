// Suite: pinned QA byte invariants
// Invariant: one gap per pending QA requirement with a byte phrase and commit token.
// Boundary IN: temporary Spec artifacts through CheckStage and parsed Task declarations.
// Boundary OUT: Git history, semantic judgment and execution of authored commands.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestPinnedByteInvariantGapNamesTheRequirementAndTheCommit(t *testing.T) {
	t.Parallel()
	const requirements = "## Requirements\n\n1. MUST hold `docs/references/behavior-surfaces.json` byte-identical to its bytes at\n   `abc1234`.\n2. MUST confirm from the diff against `def5678` that two files are byte-identical.\n"
	root := writeConsumerFixture(t, "qa: task_01\n", consumerFixtureTask{id: "task_01", kind: "qa", section: requirements})
	result := checkProofFixture(t, root, speccheck.StageTasks)
	requireProofCount(t, result, speccheck.CodeQABytesPinned, 2)
	findings := findingsWithCode(result, speccheck.CodeQABytesPinned)
	for index, want := range []struct {
		line  int
		token string
	}{{13, "abc1234"}, {15, "def5678"}} {
		finding := findings[index]
		const path = "docs/specs/proof/task_01.md"
		if finding.Severity != speccheck.SeverityGap || !hasExactLocation(finding, path, want.line) || len(finding.Where) != 1 {
			t.Errorf("finding = %#v, want gap at %s:%d", finding, path, want.line)
		}
		if wantSummary := fmt.Sprintf("%s line %d pins a byte invariant to commit %s", path, want.line, want.token); finding.Summary != wantSummary {
			t.Errorf("summary = %q, want %q", finding.Summary, wantSummary)
		}
	}
	for _, phrase := range []string{"byte-identical", "byte identical", "byte identity", "byte-for-byte", "byte-equal", "identical bytes", "same bytes", "unchanged bytes"} {
		t.Run(phrase, func(t *testing.T) {
			t.Parallel()
			for _, token := range []string{"abc1234", strings.Repeat("a1", 20), "abc1234.", "abc1234..def5678", "def5678..."} {
				task := spec.Task{Type: spec.TaskTypeQA, Status: spec.StatusPending, Requirements: []spec.TaskDeclaration{{Text: strings.ToUpper(phrase) + " at " + token + " and fed9876", Line: 7}}}
				findings := speccheck.PinnedByteInvariants(task, false)
				if len(findings) != 1 || !strings.Contains(findings[0].Summary, strings.Split(strings.TrimRight(token, "."), "..")[0]) {
					t.Errorf("phrase %q token %q: findings = %#v, want one gap naming first token", phrase, token, findings)
				}
			}
		})
	}
	// A manifest-designated gate qualifies independently of its type.
	task := spec.Task{Type: spec.TaskTypeBackend, Status: spec.StatusPending, Requirements: []spec.TaskDeclaration{{Text: "byte-identical at abc1234", Line: 7}}}
	if got := speccheck.PinnedByteInvariants(task, true); len(got) != 1 {
		t.Errorf("manifest QA findings = %#v, want one gap", got)
	}
}

func TestPinnedByteInvariantSkipsUnpinnedRequirementsAndOtherTasks(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, requirement, kind, status string }{
		{"Delivery Base", "MUST compare byte-identical with the Delivery Base.", "qa", "pending"},
		{"replay token", "MUST replay the run at abc1234.", "qa", "pending"},
		{"path", "MUST hold byte-identical docs/abc1234.", "qa", "pending"},
		{"date", "MUST hold byte-identical on 20261013 and 2026-10-13.", "qa", "pending"},
		{"ADR", "MUST hold byte-identical per ADR-0276.", "qa", "pending"},
		{"digest", "MUST hold byte-identical at " + strings.Repeat("a1", 32) + ".", "qa", "pending"},
		{"all letters", "MUST hold byte-identical at abcdefa.", "qa", "pending"},
		{"all digits", "MUST hold byte-identical at 1234567.", "qa", "pending"},
		{"short", "MUST hold byte-identical at abc123.", "qa", "pending"},
		{"long", "MUST hold byte-identical at " + strings.Repeat("a1", 20) + "a.", "qa", "pending"},
		{"prefix", "MUST hold byte-identical at _abc1234 or rev-abc1234.", "qa", "pending"},
		{"uppercase token", "MUST hold byte-identical at ABC1234.", "qa", "pending"},
		{"non-QA", "MUST hold byte-identical at abc1234.", "backend", "pending"},
		{"completed QA", "MUST hold byte-identical at abc1234.", "qa", "completed"},
		{"separate requirements", "MUST hold byte-identical.\n2. MUST replay abc1234.", "qa", "pending"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			qa := "qa: task_01\n"
			if tt.kind != "qa" {
				qa = consumerDeclinedQA
			}
			root := writeConsumerFixture(t, qa, consumerFixtureTask{id: "task_01", kind: tt.kind, status: tt.status, section: "## Requirements\n\n1. " + tt.requirement + "\n"})
			requireProofCount(t, checkProofFixture(t, root, speccheck.StageTasks), speccheck.CodeQABytesPinned, 0)
		})
	}
	t.Run("other sections and commands are not requirements", func(t *testing.T) {
		t.Parallel()
		marker := filepath.Join(t.TempDir(), "executed")
		root := writeConsumerFixture(t, "qa: task_01\n", consumerFixtureTask{id: "task_01", kind: "qa",
			section:  "## Overview\n\nByte-identical at abc1234.\n\n## Requirements\n\n1. MUST compare with the Delivery Base.\n\n## Acceptance Criteria\n\n- Byte-identical at abc1234.\n\n## Result\n\nByte-identical at abc1234.\n",
			commands: []string{"printf 'byte-identical abc1234'; touch " + marker},
		})
		requireProofCount(t, checkProofFixture(t, root, speccheck.StageTasks), speccheck.CodeQABytesPinned, 0)
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("authored command marker stat = %v, want absent", err)
		}
	})
}

func TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict(t *testing.T) {
	t.Parallel()
	root := writeConsumerFixture(t, "qa: task_01\n", consumerFixtureTask{id: "task_01", kind: "qa", section: "## Requirements\n\n1. MUST keep byte-identical to Delivery Base abc1234.\n"})
	for _, stage := range []speccheck.Stage{speccheck.StagePRD, speccheck.StageTechSpec} {
		result := checkProofFixture(t, root, stage)
		requireProofCount(t, result, speccheck.CodeQABytesPinned, 0)
		if !hasSkip(result, speccheck.CodeQABytesPinned, "stage "+string(stage)) {
			t.Errorf("skips = %#v, want stage skip", result.Skipped)
		}
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		requireProofCount(t, result, speccheck.CodeQABytesPinned, 1)
		gap := speccheck.Result{Slug: result.Slug, Findings: findingsWithCode(result, speccheck.CodeQABytesPinned)}
		speccheck.PromoteGaps(&gap)
		if gap.Findings[0].Severity != speccheck.SeverityError {
			t.Errorf("strict result = %#v, want failing error", gap)
		}
	}
	if err := os.Remove(filepath.Join(root, "docs/specs/proof/_tasks.md")); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		requireProofCount(t, result, speccheck.CodeQABytesPinned, 0)
		if !hasSkip(result, speccheck.CodeQABytesPinned, "_tasks.md") {
			t.Errorf("skips = %#v, want absent Task Graph skip", result.Skipped)
		}
	}
}

func TestPinnedByteInvariantFixNamesTheDeliveryBaseAndTheSkillRule(t *testing.T) {
	t.Parallel()
	root := writeConsumerFixture(t, "qa: task_01\n", consumerFixtureTask{id: "task_01", kind: "qa", section: "## Requirements\n\n1. MUST keep byte-identical at abc1234.\n"})
	result := checkProofFixture(t, root, speccheck.StageTasks)
	requireProofCount(t, result, speccheck.CodeQABytesPinned, 1)
	const want = "Compare with the Delivery Base, the merge base of the audited head with the default branch, and keep the commit where the plan measured out of the requirement, following the write-tasks rule \"Byte invariants name the delivery base\"."
	if got := findingsWithCode(result, speccheck.CodeQABytesPinned)[0].Fix; got != want {
		t.Errorf("fix = %q, want %q", got, want)
	}
}
