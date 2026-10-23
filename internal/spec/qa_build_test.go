package spec

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Suite: final refusal report provenance.
// Invariant: a resolved audited commit is quoted once beside the verdict.
// Boundary IN: auditor evidence and a precondition refusal.
// Boundary OUT: report selection and settlement.

func TestPreconditionRefusalReportRecordsItsAuditedCommit(t *testing.T) {
	t.Parallel()
	const commit = "1234567e89012345678901234567890123456789"
	refusal := PreconditionRefusal{CheckName: "strict", Reason: "SC-REQUIREMENT-CONTRADICTORY"}
	for _, tt := range []struct{ name, commit string }{
		{"resolved", " \t" + commit + "\n "},
		{"blank", " \t\n "},
		{"unresolved", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var content strings.Builder
			if err := WritePreconditionRefusalReport(&content, refusal, AuditorEvidence{AuditedCommit: tt.commit}); err != nil {
				t.Fatalf("WritePreconditionRefusalReport: %v", err)
			}
			prefix := "---\nverdict: fail\n"
			wantKeys := 0
			if strings.TrimSpace(tt.commit) != "" {
				prefix += "build: " + strconv.Quote(commit) + "\n"
				wantKeys = 1
			}
			prefix += "auditing_binary:"
			if !strings.HasPrefix(content.String(), prefix) || strings.Count(content.String(), "\nbuild:") != wantKeys {
				t.Fatalf("report must start with %q and contain %d build keys:\n%s", prefix, wantKeys, content.String())
			}
			var withoutCommit strings.Builder
			if err := WritePreconditionRefusalReport(&withoutCommit, refusal, AuditorEvidence{}); err != nil {
				t.Fatalf("WritePreconditionRefusalReport without commit: %v", err)
			}
			stripped := strings.Replace(content.String(), "build: "+strconv.Quote(commit)+"\n", "", 1)
			if stripped != withoutCommit.String() {
				t.Fatalf("audited commit changed other refusal report bytes:\n%s", content.String())
			}
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-10-22.md"), content.String())
			report, err := ReadQAReport(specDir)
			if err != nil {
				t.Fatalf("ReadQAReport: %v", err)
			}
			if report.Verdict != VerdictFail || report.RowsBlockedPrecondition != 1 || report.Precondition != refusal {
				t.Fatalf("read refusal = %+v, want unchanged refusal %+v", report, refusal)
			}
		})
	}
}
