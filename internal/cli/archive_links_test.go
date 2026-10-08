package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestArchiveCommandRefusesABrokenOutwardLink(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Build", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	dir := filepath.Join(repo, "docs", "specs", implementTestSlug)
	prdPath := filepath.Join(dir, "_prd.md")
	before := mustRead(t, prdPath) + "\n[broken](../../adr/missing.md)\n"
	if err := os.WriteFile(prdPath, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	withNoEngineCollaborators(t)
	var stdout, stderr bytes.Buffer
	commitArchiveFixture(t)
	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)
	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), "Preflight failed") || !strings.Contains(stderr.String(), `_prd.md:`) || !strings.Contains(stderr.String(), `"../../adr/missing.md"; fix or remove each link, then retry the archive`) {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
	if got := mustRead(t, prdPath); got != before {
		t.Fatal("PRD changed")
	}
	assertPathMissing(t, archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug))
}

func TestArchiveCommandReportsRewrittenLinks(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name              string
		override, rewrite bool
	}{
		{name: "normal", rewrite: true},
		{name: "override", override: true, rewrite: true},
		{name: "no rewrites"},
		{name: "override no rewrites", override: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Build", status: string(spec.StatusCompleted)}})
			verdict := spec.VerdictPass
			if tt.override {
				verdict = "fail"
			}
			writeArchiveQAReport(t, repo, verdict)
			if tt.rewrite {
				target := filepath.Join(repo, "docs", "adr", "decision.md")
				if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(target, []byte("decision"), 0o644); err != nil {
					t.Fatal(err)
				}
				prd := filepath.Join(repo, "docs", "specs", implementTestSlug, "_prd.md")
				if err := os.WriteFile(prd, []byte(mustRead(t, prd)+"\n[ADR](../../adr/decision.md)\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			withNoEngineCollaborators(t)
			args := []string{"archive", implementTestSlug}
			if tt.override {
				args = append(args, "--qa-override", "--approval", "maintainer", "--reason", "authorized test")
			}
			var stdout, stderr bytes.Buffer
			commitArchiveFixture(t)
			code := runCLIContext(t, context.Background(), args, &stdout, &stderr)
			want := archiveConfirmation(t, repo, archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug)+".md", tt.override)
			if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
				t.Fatalf("code=%d stdout=%q want=%q stderr=%q", code, &stdout, want, &stderr)
			}
		})
	}
}
