package cli

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestResumeAcceptsAnArchiveCommitWithRewrittenLinks(t *testing.T) {
	for _, files := range [][]string{{"_prd.md"}, {"task_01.md"}, {"qa/links.md"}, {"_prd.md", "task_01.md", "qa/links.md"}} {
		t.Run(strings.Join(files, "+"), func(t *testing.T) {
			repository, reviewedHead, archiveHead := commitLinkRewritingArchive(t, files, nil)
			item := resumeArchivedDelivery(t, repository, reviewedHead)
			if item.Blocker == delivery.BlockerReviewStale {
				t.Fatalf("resumed archive item = %+v, want link rewrites accepted", item)
			}
			if !reflect.DeepEqual(item.CandidateCommits, []string{reviewedHead, archiveHead}) {
				t.Fatalf("candidate commits = %q, want reviewed and archive heads", item.CandidateCommits)
			}
		})
	}
}

func TestResumeRefusesALinkRewritingArchiveCommitWithExtraChanges(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{"extra Task byte", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "task_01.md")
			mustWrite(t, path, mustRead(t, path)+"x")
		}},
		{"another link target", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "task_01.md")
			mustWrite(t, path, strings.ReplaceAll(mustRead(t, path), "archive-link.md", "other.md"))
		}},
		{"extra PRD byte", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "_prd.md")
			mustWrite(t, path, mustRead(t, path)+"x")
		}},
		{"another PRD link target", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "_prd.md")
			mustWrite(t, path, strings.ReplaceAll(mustRead(t, path), "archive-link.md", "other.md"))
		}},
		{"frontmatter link rewrite", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "_prd.md")
			mustWrite(t, path, strings.Replace(mustRead(t, path), "../../adr/frontmatter.md", "../../../adr/frontmatter.md", 1))
		}},
		{"added file", func(t *testing.T, dir string) {
			mustWrite(t, filepath.Join(dir, "added.md"), "extra\n")
		}},
		{"removed file", func(t *testing.T, dir string) {
			if err := os.Remove(filepath.Join(dir, "qa", "links.md")); err != nil {
				t.Fatal(err)
			}
		}},
		{"changed mode", func(t *testing.T, dir string) {
			if err := os.Chmod(filepath.Join(dir, "task_01.md"), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{"changed type", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "task_01.md")
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("_prd.md", path); err != nil {
				t.Fatal(err)
			}
		}},
		{"non-Markdown rewrite", func(t *testing.T, dir string) {
			mustWrite(t, filepath.Join(dir, "links.txt"), "[ADR](../../../adr/archive-link.md)\n")
		}},
		{"changed query", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "task_01.md")
			mustWrite(t, path, strings.ReplaceAll(mustRead(t, path), "?view=full", "?view=other"))
		}},
		{"changed fragment", func(t *testing.T, dir string) {
			path := filepath.Join(dir, "task_01.md")
			mustWrite(t, path, strings.ReplaceAll(mustRead(t, path), "#decision", "#other"))
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, reviewedHead, _ := commitLinkRewritingArchive(t, []string{"_prd.md", "task_01.md", "qa/links.md"}, tt.mutate)
			item := resumeArchivedDelivery(t, repository, reviewedHead)
			if item.Stage != store.DeliveryStageParked || item.Blocker != delivery.BlockerReviewStale {
				t.Fatalf("resumed archive item = %+v, want parked as review-stale", item)
			}
			if !reflect.DeepEqual(item.CandidateCommits, []string{reviewedHead}) {
				t.Fatalf("candidate commits = %q, want only reviewed head", item.CandidateCommits)
			}
		})
	}
}

func commitLinkRewritingArchive(t *testing.T, files []string, mutate func(*testing.T, string)) (string, string, string) {
	t.Helper()
	_, repository := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	appendArchiveUnreachableDeclarations(t, repository, []string{"a maintainer publishes the tagged release"})
	writeArchiveQAReport(t, repository, spec.VerdictPartial,
		"rows_blocked_declared: 1", "rows_blocked_finding: 0", "rows_blocked_environment: 0")
	source := filepath.Join(repository, "docs", "specs", implementTestSlug)
	if err := os.MkdirAll(filepath.Join(repository, "docs", "adr"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(repository, "docs", "adr", "archive-link.md"), "# Decision\n")
	mustWrite(t, filepath.Join(repository, "docs", "adr", "other.md"), "# Other\n")
	prd := filepath.Join(source, "_prd.md")
	mustWrite(t, prd, strings.Replace(mustRead(t, prd), "---\n", "---\nreference: '../../adr/frontmatter.md'\n", 1))
	mustWrite(t, filepath.Join(source, "links.txt"), "[ADR](../../adr/archive-link.md)\n")
	for _, file := range files {
		path := filepath.Join(source, file)
		prefix, err := filepath.Rel(filepath.Dir(path), filepath.Join(repository, "docs", "adr", "archive-link.md"))
		if err != nil {
			t.Fatal(err)
		}
		content := ""
		if file == "_prd.md" || file == "task_01.md" {
			content = mustRead(t, path)
		}
		mustWrite(t, path, content+"\n[ADR]("+filepath.ToSlash(prefix)+"?view=full#decision)\n")
	}
	gittest.Run(t, repository, "add", "-A")
	gittest.Run(t, repository, "commit", "-m", "docs: record reviewed links")
	reviewedHead := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
	result, err := spec.Archive(spec.ArchiveRequest{
		SpecsRoot: filepath.Join(repository, "docs", "specs"), BuiltInRoot: true,
		Slug: implementTestSlug, ArchivedAt: time.Date(2026, time.October, 4, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("archive Spec with links: %v", err)
	}
	if result.RewrittenLinks != len(files) {
		t.Fatalf("rewritten links = %d, want %d", result.RewrittenLinks, len(files))
	}
	if mutate != nil {
		mutate(t, result.ArchivedDir)
	}
	gittest.Run(t, repository, "add", "-A")
	gittest.Run(t, repository, "commit", "-m", "docs: archive "+implementTestSlug)
	return repository, reviewedHead, strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
}
