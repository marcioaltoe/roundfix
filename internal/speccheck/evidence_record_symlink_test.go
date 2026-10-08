package speccheck

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordEvidenceSnapshotsRefusesASymlinkedReport(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "victim.txt")
	if err := os.WriteFile(outside, []byte("untouched\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	qa := filepath.Join(repo, "docs", "specs", "demo", "qa")
	if err := os.MkdirAll(qa, 0o755); err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(qa, "qa-report-2026-10-01.md")
	if err := os.Symlink(outside, report); err != nil {
		t.Fatal(err)
	}
	_, err := RecordEvidenceSnapshots(context.Background(), repo, report, "head")
	var fileErr *EvidenceReportFileError
	if !errors.As(err, &fileErr) {
		t.Fatalf("err = %v, want EvidenceReportFileError", err)
	}
	if got, _ := os.ReadFile(outside); string(got) != "untouched\n" {
		t.Fatalf("symlink target was written: %q", got)
	}
}

func TestRecordEvidenceSnapshotsRefusesAReportUnderASymlinkedDirectory(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	outsideDir := t.TempDir()
	victim := filepath.Join(outsideDir, "qa-report-2026-10-01.md")
	if err := os.WriteFile(victim, []byte("---\nverdict: pass\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := filepath.Join(repo, "docs", "specs", "demo")
	if err := os.MkdirAll(spec, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideDir, filepath.Join(spec, "qa")); err != nil {
		t.Fatal(err)
	}
	_, err := RecordEvidenceSnapshots(context.Background(), repo, filepath.Join(spec, "qa", "qa-report-2026-10-01.md"), "head")
	var fileErr *EvidenceReportFileError
	if !errors.As(err, &fileErr) {
		t.Fatalf("err = %v, want EvidenceReportFileError", err)
	}
	if got, _ := os.ReadFile(victim); string(got) != "---\nverdict: pass\n---\n" {
		t.Fatalf("target under symlinked directory was written: %q", got)
	}
}
