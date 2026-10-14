package spec

// Suite: Legacy Link Rebase and conversion planning.
// Invariant: only relative destinations and bounded path text change when text moves.
// Boundary IN: repository-relative paths and disposable legacy Git repositories.
// Boundary OUT: rebased bytes, planned records, promoted copies and applied files.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func legacyRebaseFixture() LegacyLinkRebase {
	folder := "docs/history/specs/0001-demo"
	return LegacyLinkRebase{
		Folder: folder, RecordPath: folder + ".md",
		Files:    []string{folder + "/_prd.md", folder + "/notes.md", folder + "/references/probe data.md"},
		Promoted: map[string]string{folder + "/references/probe data.md": "docs/references/probe data.md"},
	}
}

func TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile(t *testing.T) {
	t.Parallel()
	r := legacyRebaseFixture()
	for _, tc := range []struct{ name, from, to, content, want string }{
		{"promoted summary", r.Folder + "/_prd.md", r.RecordPath,
			`See [probe](references/probe%20data.md?raw=1#result "title").`,
			`See [probe](../../references/probe%20data.md?raw=1#result "title").`},
		{"unpromoted file", r.Folder + "/_prd.md", r.RecordPath,
			"[notes](notes.md)", "[notes](0001-demo.md)"},
		{"unlisted descendant", r.Folder + "/_prd.md", r.RecordPath,
			"[old](evidence/output.txt)", "[old](0001-demo.md)"},
		{"core from promoted copy", r.Folder + "/references/probe data.md", "docs/references/probe data.md",
			"[PRD](../_prd.md#goals)", "[PRD](../history/specs/0001-demo.md#goals)"},
		{"folder", r.Folder + "/_prd.md", r.RecordPath,
			"[folder](./)", "[folder](0001-demo.md)"},
		{"reference definition", r.Folder + "/_prd.md", r.RecordPath,
			"[probe]: <references/probe%20data.md#result> 'Title'\n", "[probe]: <../../references/probe%20data.md#result> 'Title'\n"},
		{"angle destination", r.Folder + "/_prd.md", r.RecordPath,
			"![probe](<references/probe data.md>)", "![probe](<../../references/probe%20data.md>)"},
		{"path text", r.Folder + "/_prd.md", r.RecordPath,
			"`" + r.Folder + "/_prd.md`\n" + r.Folder + "/\n" + r.Folder + "/references/probe data.md\n" + r.Folder + "/notes.md\n",
			"`" + r.RecordPath + "`\n" + r.RecordPath + "\ndocs/references/probe data.md\n" + r.Folder + "/notes.md\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := r.Rebase([]byte(tc.content), tc.from, tc.to); string(got) != tc.want {
				t.Fatalf("rebase = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLegacyLinkRebaseKeepsOutsideTargets(t *testing.T) {
	t.Parallel()
	r := legacyRebaseFixture()
	for _, tc := range []struct{ name, from, to, content, want string }{
		{"summary ADR", r.Folder + "/_prd.md", r.RecordPath,
			"[ADR](../../../adr/0001-rule.md?q=1#decision)", "[ADR](../../adr/0001-rule.md?q=1#decision)"},
		{"copy ADR", r.Folder + "/references/probe data.md", "docs/references/probe data.md",
			"[ADR](../../../../adr/0001-rule.md)", "[ADR](../adr/0001-rule.md)"},
		{"unresolved outside target", r.Folder + "/_prd.md", r.RecordPath,
			"[missing](../../../references/missing.md)", "[missing](../../references/missing.md)"},
		{"neighbor folder", r.Folder + "/_prd.md", r.RecordPath,
			"[neighbor](../0001-demo-other/_prd.md)", "[neighbor](0001-demo-other/_prd.md)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := r.Rebase([]byte(tc.content), tc.from, tc.to)
			if string(got) != tc.want {
				t.Fatalf("rebase = %q, want %q", got, tc.want)
			}
			before := scanArchiveMarkdown([]byte(tc.content))[0]
			after := scanArchiveMarkdown(got)[0]
			oldTarget, oldSuffix, _ := archiveLinkTarget(tc.content[before.start:before.end], filepath.Dir(tc.from))
			newTarget, newSuffix, _ := archiveLinkTarget(string(got[after.start:after.end]), filepath.Dir(tc.to))
			if oldTarget != newTarget || oldSuffix != newSuffix {
				t.Fatalf("target or suffix changed: %q %q -> %q %q", oldTarget, oldSuffix, newTarget, newSuffix)
			}
		})
	}
}

func TestLegacyLinkRebaseLeavesOtherBytes(t *testing.T) {
	t.Parallel()
	r := legacyRebaseFixture()
	for _, tc := range []struct{ name, content string }{
		{"no links", "\x00# Text\r\n\r\n  Exact bytes.\t\n"},
		{"absolute URL", "[web](https://example.com/a?q=1#b)\n"},
		{"scheme", "[mail](mailto:person@example.com) [network](//example.com/path)\n"},
		{"root path", "[absolute](/docs/adr/rule.md)\n"},
		{"anchor", "[section](#goals) [query](?raw=1)\n"},
		{"code span", "`[PRD](../_prd.md)` and ``[probe](references/probe%20data.md)``\n"},
		{"fenced block", "```md\n[probe](references/probe%20data.md)\n```\n~~~\n[PRD](_prd.md)\n~~~\n"},
		{"leaves repository", "[outside](../../../../../outside.md) [encoded](%2E%2E/%2E%2E/%2E%2E/%2E%2E/%2E%2E/outside.md)\n"},
		{"malformed destination", "[bad](references/probe%GG.md)\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			content := []byte(tc.content)
			got := r.Rebase(content, r.Folder+"/_prd.md", r.RecordPath)
			if !bytes.Equal(got, content) || string(content) != tc.content {
				t.Fatalf("bytes changed: %q -> %q", tc.content, got)
			}
		})
	}
}

func TestLegacyConversionRebasesTheOutcomeAndPromotedCopies(t *testing.T) {
	t.Parallel()
	for _, source := range []string{"PRD", "first section", "QA", "empty QA", "supersession", "no links"} {
		t.Run(source, func(t *testing.T) {
			t.Parallel()
			summary := "See [probe](references/probe.md) and [ADR](../../../adr/rule.md)."
			want := "See [probe](../../references/probe.md) and [ADR](../../adr/rule.md)."
			body := summary + "\n"
			if source == "first section" {
				body = "## Problem\n\n" + body
			}
			if source == "no links" {
				body, summary, want = "Preserved bytes.\n", "Preserved bytes.", "Preserved bytes."
			}
			req := legacySummaryFixture(t, body)
			folder := filepath.Join(req.ArchiveRoot, req.Slug)
			writeFile(t, filepath.Join(folder, "references/probe.md"), "# Probe\n\n[PRD](../_prd.md#goals) [ADR](../../../../adr/rule.md).\n")
			writeFile(t, filepath.Join(req.RepositoryRoot, "docs/adr/rule.md"), "# Rule\n")
			if source == "QA" || source == "empty QA" {
				outcome := ""
				if source == "QA" {
					outcome = "See [probe](../references/probe.md) and [ADR](../../../../adr/rule.md)."
					summary = outcome
				}
				writeFile(t, filepath.Join(folder, "qa/qa-report-2026-09-03.md"), "---\nverdict: pass\n---\n\n# QA\n\n## Outcome\n\n"+outcome+"\n")
			}
			if source == "supersession" {
				if _, err := WriteSupersession(folder, "successor", summary, time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)); err != nil {
					t.Fatal(err)
				}
			}
			if source == "no links" {
				writeFile(t, filepath.Join(folder, "references/probe.md"), "# Probe\n\nUnchanged bytes.\n")
			}
			req.Promote = []string{"docs/history/specs/demo/references/probe.md"}
			req.SourceRevision = legacyCommit(t, req.RepositoryRoot, "docs: legacy links fixture")
			before := legacyTree(t, req)
			plan := legacyPlan(t, req)
			legacyAssertWhole(t, req, before)
			if plan.Record.Outcome != want {
				t.Fatalf("outcome = %q, want %q", plan.Record.Outcome, want)
			}
			rendered, err := RenderArchiveRecord(plan.Record)
			if err != nil || !bytes.Equal(rendered, plan.Rendered) || !strings.HasSuffix(string(rendered), want+"\n") {
				t.Fatalf("render = %q, err = %v", rendered, err)
			}
			archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
			if err != nil || archived.Record.Outcome != summary {
				t.Fatalf("folder reader outcome = %q, err = %v, want %q", archived.Record.Outcome, err, summary)
			}
			wantCopy := "# Probe\n\n[PRD](../history/specs/demo.md#goals) [ADR](../adr/rule.md).\n"
			if source == "no links" {
				wantCopy = "# Probe\n\nUnchanged bytes.\n"
				plain, err := BuildArchiveRecord(ArchiveRecordInput{Legacy: true, SpecDir: folder, Slug: req.Slug, Source: plan.Folder, SourceRevision: req.SourceRevision, Promoted: plan.Promoted})
				if err != nil {
					t.Fatal(err)
				}
				plainBytes, err := RenderArchiveRecord(plain)
				if err != nil || !bytes.Equal(plainBytes, plan.Rendered) {
					t.Fatalf("link-free record changed: %v", err)
				}
			}
			if got := plan.legacyCopies["docs/references/probe.md"]; string(got) != wantCopy {
				t.Fatalf("planned copy = %q, want %q", got, wantCopy)
			}
			if err := ApplyLegacyConversion(req.RepositoryRoot, plan); err != nil {
				t.Fatal(err)
			}
			for path, wantBytes := range map[string][]byte{plan.RecordPath: plan.Rendered, "docs/references/probe.md": []byte(wantCopy)} {
				got, err := os.ReadFile(filepath.Join(req.RepositoryRoot, path))
				if err != nil || !bytes.Equal(got, wantBytes) {
					t.Fatalf("written %s = %q, err = %v, want %q", path, got, err, wantBytes)
				}
			}
		})
	}
}
