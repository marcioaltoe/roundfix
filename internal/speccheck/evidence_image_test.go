// Suite: readable QA evidence
// Invariant: tracked evidence accepts text and images by content, never by name.
// Boundary IN: mechanical stage over real temporary Git repositories
// Boundary OUT: QA Agent and report commit
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestEvidenceAcceptsTrackedImages(t *testing.T) {
	t.Parallel()
	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
	paths := []string{evidenceScratchReportPath}
	for _, image := range []struct{ name, content string }{
		{"screen.png", "\x89PNG\r\n\x1a\n\x00image"},
		{"screen.jpg", "\xff\xd8\xff\xe0\x00image"},
		{"diagram.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"></svg>\n"},
		{"capture", "\x89PNG\r\n\x1a\n\x00image"},
		{"probe-bound.txt", strings.Repeat("a", speccheck.EvidenceProbeBytes) + "\x00"},
	} {
		path := evidenceScratchRoot + "/" + image.name
		writeMechanicalFile(t, repoRoot, path, image.content)
		paths = append(paths, path)
	}
	commitMechanicalFiles(t, repoRoot, "record image evidence", paths...)
	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})
	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalEvidencePath)
}

func TestEvidenceRefusesNonImageBinaries(t *testing.T) {
	t.Parallel()
	for _, sample := range []struct{ name, content string }{
		{"program", "\x7fELF\x00fixture"},
		{"state.db", "SQLite format 3\x00fixture"},
		{"report.pdf", "%PDF-1.7\n\x00fixture"},
		{"program.png", "\x7fELF\x00fixture"},
		{"last-probe-byte.bin", strings.Repeat("a", speccheck.EvidenceProbeBytes-1) + "\x00"},
	} {
		t.Run(sample.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := newMechanicalGitRepo(t)
			path := evidenceScratchRoot + "/" + sample.name
			writeMechanicalFile(t, repoRoot, evidenceScratchReportPath, mechanicalEvidenceReport())
			writeMechanicalFile(t, repoRoot, path, sample.content)
			commitMechanicalFiles(t, repoRoot, "record non-image evidence", evidenceScratchReportPath, path)
			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: repoRoot, ReportPath: evidenceScratchReportPath})
			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalEvidencePath)
			if len(findings) != 1 || findings[0].File != path ||
				findings[0].Detail != "tracked evidence path "+path+" is a non-image binary" {
				t.Fatalf("findings = %#v, want one named non-image binary %s", findings, path)
			}
			wantFix := "Remove the gate's scratch state from the Spec evidence directory and keep only text, PNG, JPEG or SVG artifacts."
			if findings[0].Fix != wantFix {
				t.Fatalf("fix = %q, want %q", findings[0].Fix, wantFix)
			}
		})
	}
}

func TestReadableEvidenceDecidesByContent(t *testing.T) {
	t.Parallel()
	for _, sample := range []struct {
		name     string
		probe    string
		readable bool
	}{
		{"empty", "", true},
		{"text", "command output\n", true},
		{"SVG", "<svg></svg>\n", true},
		{"PNG", "\x89PNG\r\n\x1a\n\x00image", true},
		{"JPEG", "\xff\xd8\xff\xe0\x00image", true},
		{"ELF", "\x7fELF\x00fixture", false},
		{"SQLite", "SQLite format 3\x00fixture", false},
		{"PDF with NUL", "%PDF-1.7\n\x00fixture", false},
		{"truncated PNG with NUL", "\x89PNG\r\n\x1a\x00", false},
		{"truncated JPEG with NUL", "\xff\xd8\x00", false},
		{"PNG after leading NUL", "\x00\x89PNG\r\n\x1a\n", false},
		{"JPEG after leading NUL", "\x00\xff\xd8\xff", false},
		{"unknown binary", "\x00fixture", false},
	} {
		t.Run(sample.name, func(t *testing.T) {
			t.Parallel()
			if got := speccheck.ReadableEvidence([]byte(sample.probe)); got != sample.readable {
				t.Fatalf("ReadableEvidence(%q) = %t, want %t", sample.probe, got, sample.readable)
			}
		})
	}
}
