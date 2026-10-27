package cli

// Suite: Archive eligibility for a Declared QA Decline.
// Boundary: CLI dispatch against committed temporary repositories without QA Reports.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestArchiveAcceptsADeclaredQADecline(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	declareImplementQADeclined(t, repo, "no behavioral QA surface")
	commitArchiveFixture(t)
	var out, errout bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"archive", implementTestSlug}, &out, &errout)
	if code != 0 || !strings.HasPrefix(out.String(), "archived "+implementTestSlug+" -> ") || errout.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &errout)
	}
	record := mustRead(t, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"))
	if !strings.Contains(record, "disposition: no-qa\n") || strings.Contains(record, "qa_override") {
		t.Fatalf("unexpected Archive Record: %s", record)
	}
	if _, err := os.Stat(filepath.Join(repo, "docs/specs", implementTestSlug)); !os.IsNotExist(err) {
		t.Fatalf("Spec folder remains: %v", err)
	}
}

func TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	declareImplementQADeclined(t, repo, "no behavioral QA surface")
	commitArchiveFixture(t)
	before := snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--qa-override", "--approval", "maintainer", "--reason", "archive the declared decline"}, "", archivePullRequestPreflight("QA archive override is not allowed because the Spec already qualifies for normal archive"), 2)
	if _, err := os.Stat(filepath.Join(repo, "docs/specs", implementTestSlug)); err != nil {
		t.Fatalf("Spec folder changed: %v", err)
	}
	after := snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))
	if len(after) != len(before) {
		t.Fatal("refusal changed document count")
	}
	for path, content := range before {
		if !bytes.Equal(after[path], content) {
			t.Fatalf("refusal changed %s", path)
		}
	}
}

func TestArchiveRefusesADeclaredQADeclineWithAPendingTask(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}, {id: "task_02", status: string(spec.StatusPending)}})
	declareImplementQADeclined(t, repo, "no behavioral QA surface")
	commitArchiveFixture(t)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug}, "", archivePullRequestPreflight(`Task "task_02" is "pending"; archive requires every Task to be "completed"`), 2)
	if _, err := os.Stat(filepath.Join(repo, "docs/specs", implementTestSlug)); err != nil {
		t.Fatalf("Spec folder changed: %v", err)
	}
}
