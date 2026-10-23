package daemon

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

// Suite: final mechanical QA report provenance.
// Invariant: only final reports record the commit the gate audited.
// Boundary IN: mechanical results, auditor evidence, and the Run Worktree HEAD.
// Boundary OUT: report settlement.

func TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal(t *testing.T) {
	t.Parallel()
	const commit = "1234567e89012345678901234567890123456789"
	refusal := spec.PreconditionRefusal{CheckName: speccheck.GatePreconditionCheck, Reason: "missing constraint"}
	for _, tt := range []struct {
		name   string
		result speccheck.MechanicalResult
		commit string
		final  bool
	}{
		{"refusal", speccheck.MechanicalResult{Blocking: true, PreconditionRefused: true, PreconditionRefusal: refusal}, " \t" + commit + "\n ", true},
		{"blocking", speccheck.MechanicalResult{Blocking: true}, " \t" + commit + "\n ", true},
		{"pending seed", speccheck.MechanicalResult{}, commit, false},
		{"unresolved refusal", speccheck.MechanicalResult{Blocking: true, PreconditionRefused: true, PreconditionRefusal: refusal}, " \t\n ", false},
		{"unresolved blocking", speccheck.MechanicalResult{Blocking: true}, " \t\n ", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			evidence := spec.AuditorEvidence{AuditedCommit: tt.commit}
			content, err := mechanicalQAReportContent(tt.result, evidence)
			if err != nil {
				t.Fatalf("mechanicalQAReportContent: %v", err)
			}
			verdict := spec.VerdictPending
			if tt.result.Blocking {
				verdict = spec.VerdictFail
			}
			prefix := "---\nverdict: " + verdict + "\n"
			wantKeys := 0
			if tt.final {
				prefix += "build: " + strconv.Quote(commit) + "\n"
				wantKeys = 1
			}
			prefix += "auditing_binary:"
			if !strings.HasPrefix(string(content), prefix) || strings.Count(string(content), "\nbuild:") != wantKeys {
				t.Fatalf("report must start with %q and contain %d build keys:\n%s", prefix, wantKeys, content)
			}
			evidence.AuditedCommit = ""
			withoutCommit, err := mechanicalQAReportContent(tt.result, evidence)
			if err != nil {
				t.Fatalf("mechanicalQAReportContent without commit: %v", err)
			}
			stripped := strings.Replace(string(content), "build: "+strconv.Quote(commit)+"\n", "", 1)
			if stripped != string(withoutCommit) {
				t.Fatalf("audited commit changed other report bytes:\n%s", content)
			}
		})
	}
	t.Run("written refusal names the Run Worktree HEAD", func(t *testing.T) {
		fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
		engine := fixture.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, &engineFakeCommitter{}, fixture.worktree)
		plan := fixture.qaPlan()
		head := strings.TrimSpace(runGitForTest(t, plan.WorkDir, "rev-parse", "HEAD"))
		path, err := engine.writeMechanicalQAReport(context.Background(), plan, speccheck.MechanicalResult{Blocking: true, PreconditionRefused: true, PreconditionRefusal: refusal})
		if err != nil {
			t.Fatalf("writeMechanicalQAReport: %v", err)
		}
		content := readFileForTest(t, filepath.Join(plan.WorkDir, filepath.FromSlash(path)))
		prefix := "---\nverdict: fail\nbuild: " + strconv.Quote(head) + "\nauditing_binary:"
		if !strings.HasPrefix(content, prefix) || strings.Count(content, "\nbuild:") != 1 {
			t.Fatalf("written refusal must name HEAD %q once:\n%s", head, content)
		}
		report, err := spec.ReadQAReport(plan.Spec.Dir)
		if err != nil {
			t.Fatalf("ReadQAReport: %v", err)
		}
		if report.Precondition != refusal {
			t.Fatalf("read refusal = %+v, want %+v", report.Precondition, refusal)
		}
	})
	t.Run("unresolved HEAD still returns auditor evidence", func(t *testing.T) {
		engine := &Engine{}
		evidence, err := engine.resolveQAAuditorEvidence(context.Background(), TaskPlan{WorkDir: t.TempDir()})
		if err != nil {
			t.Fatalf("resolveQAAuditorEvidence: %v", err)
		}
		if evidence.AuditedCommit != "" || evidence.Binary != engine.auditor() {
			t.Fatalf("unresolved HEAD evidence = %+v", evidence)
		}
	})
}
