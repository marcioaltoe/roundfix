package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestRunArchiveMovesCompletedSpecAndStampsMetadata(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Wire the widget API", status: string(spec.StatusCompleted), needs: []string{"task_01"}},
	})
	writeArchiveQAReport(t, repoDir, spec.VerdictPass)
	withNoEngineCollaborators(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected archive exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	wantStdout := "archived " + implementTestSlug + " -> docs/specs/_archived/" + implementTestSlug + "\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout %q, got %q", wantStdout, stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected archive stderr empty, got %q", stderr.String())
	}
	assertPathMissing(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	archivedDir := filepath.Join(repoDir, "docs", "specs", "_archived", implementTestSlug)
	assertPathExists(t, archivedDir)
	assertPathExists(t, filepath.Join(archivedDir, "task_01.md"))
	prd := mustRead(t, filepath.Join(archivedDir, "_prd.md"))
	for _, want := range []string{
		"status: archived",
		"source_slug: " + implementTestSlug,
		"archived:",
	} {
		if !strings.Contains(prd, want) {
			t.Fatalf("expected archived PRD to contain %q, got:\n%s", want, prd)
		}
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveUsesConfiguredExternalSpecRoot(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should stay active", status: string(spec.StatusCompleted)},
	})
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	writeImplementSpecAtRoot(t, externalRoot, implementTestSlug, []implementSeed{
		{id: "task_01", title: "Archive external task", status: string(spec.StatusCompleted)},
	})
	reportPath := filepath.Join(externalRoot, implementTestSlug, "qa", "qa-report-2026-07-06.md")
	mustMkdir(t, filepath.Dir(reportPath))
	mustWrite(t, reportPath, implementQAReport(spec.VerdictPass))
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	withNoEngineCollaborators(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected archive exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	archivedDir := filepath.Join(externalRoot, "_archived", implementTestSlug)
	wantStdout := "archived " + implementTestSlug + " -> " + filepath.ToSlash(archivedDir) + "\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout %q, got %q", wantStdout, stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected archive stderr empty, got %q", stderr.String())
	}
	assertPathMissing(t, filepath.Join(externalRoot, implementTestSlug))
	assertPathExists(t, filepath.Join(archivedDir, "task_01.md"))
	assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveRefusesIncompleteTask(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Wire the widget API", status: string(spec.StatusPending), needs: []string{"task_01"}},
	})
	writeArchiveQAReport(t, repoDir, spec.VerdictPass)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected archive refusal exit %d, got %d", exitPreflight, code)
	}
	if stdout.String() != "" {
		t.Fatalf("expected refusal stdout empty, got %q", stdout.String())
	}
	for _, want := range []string{
		"Task \"task_02\" is \"pending\"",
		"archive requires every Task to be \"completed\"",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	assertPathMissing(t, filepath.Join(repoDir, "docs", "specs", "_archived", implementTestSlug))
	prd := mustRead(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_prd.md"))
	if strings.Contains(prd, "status: archived") || strings.Contains(prd, "source_slug: "+implementTestSlug) {
		t.Fatalf("expected active PRD left unstamped, got:\n%s", prd)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveRefusesMissingOrFailingQA(t *testing.T) {
	tests := []struct {
		name       string
		verdict    string
		wantStderr []string
	}{
		{
			name:    "missing QA Report",
			verdict: "",
			wantStderr: []string{
				"no passing QA verdict",
				"no QA Report found",
			},
		},
		{
			name:    "failing QA Report",
			verdict: spec.VerdictFail,
			wantStderr: []string{
				"no passing QA verdict",
				"newest QA Report verdict is \"fail\"",
				"expected \"pass\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
			})
			if tt.verdict != "" {
				writeArchiveQAReport(t, repoDir, tt.verdict)
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected archive refusal exit %d, got %d", exitPreflight, code)
			}
			if stdout.String() != "" {
				t.Fatalf("expected refusal stdout empty, got %q", stdout.String())
			}
			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
			assertPathMissing(t, filepath.Join(repoDir, "docs", "specs", "_archived", implementTestSlug))
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunArchiveHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"archive", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected help exit 0, got %d", code)
	}
	if stderr.String() != "" {
		t.Fatalf("expected help stderr empty, got %q", stderr.String())
	}
	for _, want := range []string{
		"Usage:",
		"roundfix archive <slug>",
		"archive creates no Run and",
		"never pushes",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected help to contain %q, got %q", want, stdout.String())
		}
	}
}

func writeArchiveQAReport(t *testing.T, repoDir string, verdict string) {
	t.Helper()
	reportPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "qa", "qa-report-2026-07-06.md")
	mustMkdir(t, filepath.Dir(reportPath))
	mustWrite(t, reportPath, implementQAReport(verdict))
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path %s to exist: %v", path, err)
	}
}

func assertPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected path %s to be missing", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat path %s: %v", path, err)
	}
}

func withNoEngineCollaborators(t *testing.T) {
	t.Helper()
	old := newEngineCollaborators
	newEngineCollaborators = func() engineCollaborators {
		t.Fatal("archive must not create Run engine collaborators")
		return engineCollaborators{}
	}
	t.Cleanup(func() {
		newEngineCollaborators = old
	})
}
