package spec

// Suite: Legacy Archive Record summaries.
// Invariant: legacy summaries preserve paragraph text and outcome precedence.
// Boundary IN: disposable legacy folders and active Spec inputs.
// Boundary OUT: conversion record bodies and archived outcome readers.

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func legacySummaryFixture(t *testing.T, body string) LegacyConversionRequest {
	t.Helper()
	req := legacyFixture(t, "status: archived\n", "")
	writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "_prd.md"), "---\nspec: demo\nstatus: archived\ncreated: 2026-09-01\n---\n\n# Demo\n\n"+body)
	req.SourceRevision = legacyCommit(t, req.RepositoryRoot, "docs: legacy summary fixture")
	return req
}

func assertLegacySummary(t *testing.T, req LegacyConversionRequest, want string) LegacyConversion {
	t.Helper()
	before := legacyTree(t, req)
	c := legacyPlan(t, req)
	legacyAssertWhole(t, req, before)
	if c.Record.Outcome != want {
		t.Fatalf("outcome = %q, want %q", c.Record.Outcome, want)
	}
	const title = "\n# Demo\n\n"
	_, body, found := strings.Cut(string(c.Rendered), title)
	if !found || body != want+"\n" {
		t.Fatalf("record body = %q, want %q", body, want+"\n")
	}
	return c
}

func TestLegacySummaryRemovesBlockQuoteMarkers(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, body, want string
	}{
		{"two lines", "> a\n> b\n\nA later paragraph.\n", "a b"},
		{"nested", "   > > a\n >   > >b\n", "a b"},
		{"no following space", ">a\n>>b\n", "a b"},
		{"quote paragraph break", "> a\n>\n> b\n", "a"},
		{"plain paragraph", "a\nb\n", "a b"},
		{"four spaces", "    > a\n    > b\n", "> a > b"},
		{"inline greater than", "> a > b\n", "a > b"},
		{"nested four spaces", ">     > a\n", "> a"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertLegacySummary(t, legacySummaryFixture(t, tc.body), tc.want)
		})
	}
}

func TestLegacySummaryFallsBackToTheFirstSection(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, body, want string
	}{
		{"plain", "## Problem\n\nThe first paragraph\ncontinues here.\n\nThe second paragraph.\n\n## Goals\n\nLater section.\n", "The first paragraph continues here."},
		{"quoted", "## Problem\n\n> The first paragraph\n>> continues here.\n", "The first paragraph continues here."},
		{"title wins", "Title paragraph.\n\n## Problem\n\nSection paragraph.\n", "Title paragraph."},
		{"empty first section", "## Problem\n\n## Goals\n\nLater section.\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := legacySummaryFixture(t, tc.body)
			c := assertLegacySummary(t, req, tc.want)
			folder, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
			if err != nil || folder.Form != ArchivedFolder || folder.Record.Outcome != tc.want {
				t.Fatalf("folder = %+v, err = %v", folder, err)
			}
			if err := ApplyLegacyConversion(req.RepositoryRoot, c); err != nil {
				t.Fatal(err)
			}
			record, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
			if err != nil || record.Form != ArchivedRecord || record.Record.Outcome != tc.want {
				t.Fatalf("record = %+v, err = %v", record, err)
			}
		})
	}
}

func TestLegacySummaryKeepsTheQAReportOutcome(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, prd, outcome, reason, want string
	}{
		{"title paragraph", "> PRD intent.\n", "> QA outcome\n   > > wins.\n", "", "QA outcome wins."},
		{"section fallback", "## Problem\n\n> PRD intent.\n", "> QA outcome\n   > > wins.\n", "", "QA outcome wins."},
		{"empty outcome", "## Problem\n\n> PRD intent.\n", "", "", "PRD intent."},
		{"supersession wins", "> PRD intent.\n", "> QA outcome\n> wins.\n", "Delivered elsewhere.", "Delivered elsewhere."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := legacySummaryFixture(t, tc.prd)
			folder := filepath.Join(req.ArchiveRoot, req.Slug)
			writeFile(t, filepath.Join(folder, "qa", "qa-report-2026-09-02.md"), "---\nverdict: pass\n---\n\n# QA\n\n## Outcome\n\nOld outcome.\n")
			writeFile(t, filepath.Join(folder, "qa", "qa-report-2026-09-03.md"), "---\nverdict: pass\n---\n\n# QA\n\n## Outcome\n\n"+tc.outcome+"\n## Evidence\n\nOther evidence.\n")
			if tc.reason != "" {
				if _, err := WriteSupersession(folder, "successor", tc.reason, time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)); err != nil {
					t.Fatal(err)
				}
			}
			req.SourceRevision = legacyCommit(t, req.RepositoryRoot, "docs: legacy QA summary fixture")
			c := assertLegacySummary(t, req, tc.want)
			if c.Record.QAReport != "qa-report-2026-09-03.md" {
				t.Fatalf("QA report = %q", c.Record.QAReport)
			}
		})
	}
}

func TestActiveArchiveRecordKeepsItsSummaryReading(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, prd, outcome, want string
	}{
		{"PRD quote", "> a\n> b\n", "", "> a > b"},
		{"QA quote", "PRD intent.\n", "> a\n> b\n", "> a > b"},
		{"no section fallback", "## Problem\n\nSection paragraph.\n", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			folder := t.TempDir()
			writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\nstatus: active\n---\n\n# Demo\n\n"+tc.prd)
			if tc.outcome != "" {
				writeFile(t, filepath.Join(folder, "qa", "qa-report-2026-09-03.md"), "---\nverdict: pass\n---\n\n# QA\n\n## Outcome\n\n"+tc.outcome)
			}
			r, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder, Slug: "demo"})
			if err != nil {
				t.Fatal(err)
			}
			if r.Outcome != tc.want {
				t.Fatalf("active outcome = %q, want %q", r.Outcome, tc.want)
			}
			content, err := RenderArchiveRecord(r)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ParseArchiveRecord(content)
			if err != nil || parsed.Outcome != tc.want {
				t.Fatalf("parsed outcome = %q, err = %v, want %q", parsed.Outcome, err, tc.want)
			}
		})
	}
}
