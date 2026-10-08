package speccheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestADeferredBacklogEntryIsTerminal(t *testing.T) {
	t.Parallel()
	if !terminalBacklogStatus("deferred") {
		t.Fatal("deferred Backlog status is not terminal")
	}
}

func TestADeferredBacklogEntryLeftActiveIsReported(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	const path = "docs/backlog/2026-09-30-deferred.md"
	if err := os.MkdirAll(filepath.Join(repo, "docs/backlog"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, path), []byte("---\nstatus: deferred\nreason: no longer needed\n---\n\n# Deferred\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var result Result
	if err := detectBacklogPromotion(&result, repo); err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings = %v, want one terminal-entry finding", result.Findings)
	}
	finding := result.Findings[0]
	if finding.Code != CodeBacklogUnmoved || finding.Severity != SeverityError ||
		!strings.Contains(finding.Summary, `declares terminal status "deferred"`) ||
		!strings.Contains(finding.Fix, "docs/history/backlog/") {
		t.Fatalf("unexpected deferred finding: %+v", finding)
	}
	if len(finding.Where) != 1 || finding.Where[0].Path != path || finding.Where[0].Line != 1 {
		t.Fatalf("finding location = %v, want status line in %s", finding.Where, path)
	}
}
