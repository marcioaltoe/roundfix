// Suite: QA user-flow binary metadata.
// Invariant: reports preserve an optional user-flow identity, and auditor evidence identifies only builds anchored in the audited repository.
// Boundary IN: QA front matter and the audited repository's Git object database.
// Boundary OUT: QA Task settlement, owned by internal/daemon.
package spec

import (
	"context"
	"path/filepath"
	"testing"

	"roundfix/internal/app"
)

func TestReadQAReportReadsTheUserFlowBinary(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	const binary = "roundfix 0.17.0 (0123456-dirty, built 2026-09-28T12:00:00Z)"
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-09-28.md"), qaReportFixture(VerdictPass, `user_flow_binary: "`+binary+`"`))

	report, err := ReadQAReport(specDir)
	if err != nil {
		t.Fatalf("ReadQAReport: %v", err)
	}
	if report.UserFlowBinary != binary {
		t.Fatalf("UserFlowBinary = %q, want %q", report.UserFlowBinary, binary)
	}
}

func TestReadQAReportWithoutAUserFlowBinaryStaysReadable(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-09-28.md"), qaReportFixture(VerdictPass))

	report, err := ReadQAReport(specDir)
	if err != nil {
		t.Fatalf("ReadQAReport: %v", err)
	}
	if report.UserFlowBinary != "" {
		t.Fatalf("UserFlowBinary = %q, want empty", report.UserFlowBinary)
	}
}

func TestAuditorEvidenceMarksASelfAudit(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	head := commitEvidenceFile(t, root, "audited.txt", "audited\n")

	for _, commit := range []string{head, head + "-dirty"} {
		commit := commit
		t.Run(commit, func(t *testing.T) {
			t.Parallel()
			evidence := ResolveAuditorEvidence(context.Background(), root, head, app.AuditingBinary{Version: "0.17.0", Commit: commit})
			if !evidence.SelfAudit {
				t.Fatalf("SelfAudit = false for build commit %q in the audited repository", commit)
			}
		})
	}
}

func TestAuditorEvidenceOutsideTheBuildRepositoryIsNotASelfAudit(t *testing.T) {
	t.Parallel()
	root := initEvidenceRepo(t)
	head := commitEvidenceFile(t, root, "audited.txt", "audited\n")

	for _, commit := range []string{"0123456789abcdef0123456789abcdef01234567", ""} {
		commit := commit
		t.Run(commit, func(t *testing.T) {
			t.Parallel()
			evidence := ResolveAuditorEvidence(context.Background(), root, head, app.AuditingBinary{Version: "0.17.0", Commit: commit})
			if evidence.SelfAudit {
				t.Fatalf("SelfAudit = true for build commit %q absent from the audited repository", commit)
			}
		})
	}
}
