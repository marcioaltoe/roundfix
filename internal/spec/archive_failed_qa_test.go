package spec

// Suite: Failed QA Archive Records.
// Invariant: legacy failure stays recorded without granting QA completion.
// Boundary IN: synthetic Spec folders and rendered Archive Records.
// Boundary OUT: conversion metadata, parser refusals, and preserved active files.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func failedQALegacyFixture(t *testing.T) LegacyConversionRequest {
	t.Helper()
	req := legacyProjectionFixture(t, "| task_01 | Implementation | backend | low | - |\n| task_qa | QA | qa | low | task_01 |\n")
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	writeFile(t, filepath.Join(folder, "qa", "qa-report-2026-09-01.md"), "---\nverdict: pass\n---\n\n# Older QA\n")
	writeFile(t, filepath.Join(folder, "qa", "qa-report-2026-09-02.md"), "---\nverdict: fail\n---\n\n# Newest QA\n\n## Outcome\n\nThe gate failed.\n")
	return req
}

func TestLegacyConversionOfAFailedQAIsFailedQA(t *testing.T) {
	t.Parallel()
	req := failedQALegacyFixture(t)
	before := legacyTree(t, req)
	conversion, err := PlanLegacyConversion(req)
	if err != nil {
		t.Fatal(err)
	}
	legacyAssertWhole(t, req, before)
	record := conversion.Record
	if record.Disposition != ArchiveFailedQA || record.QAVerdict != VerdictFail || record.QAReport != "qa-report-2026-09-02.md" || record.QAOverride != nil {
		t.Fatalf("record=%+v", record)
	}
	if strings.Contains(string(conversion.Rendered), "qa_override:") {
		t.Fatalf("record carries an override:\n%s", conversion.Rendered)
	}
	if err := ApplyLegacyConversion(req.RepositoryRoot, conversion); err != nil {
		t.Fatal(err)
	}
	archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil || archived.Form != ArchivedRecord || !reflect.DeepEqual(archived.Record, record) {
		t.Fatalf("archive=%+v err=%v", archived, err)
	}
}

func TestFailedQARecordRoundTrips(t *testing.T) {
	t.Parallel()
	req := failedQALegacyFixture(t)
	record, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: filepath.Join(req.ArchiveRoot, req.Slug), Slug: req.Slug, Legacy: true, Source: "docs/history/specs/demo", SourceRevision: req.SourceRevision})
	if err != nil {
		t.Fatal(err)
	}
	content, err := RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseArchiveRecord(content)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Disposition != ArchiveFailedQA || parsed.QAVerdict != VerdictFail || parsed.QAReport != "qa-report-2026-09-02.md" || parsed.QAOverride != nil {
		t.Fatalf("record=%+v", parsed)
	}
	if !reflect.DeepEqual(archiveRecordMapping(record), archiveRecordMapping(parsed)) {
		t.Fatalf("metadata changed: before=%+v after=%+v", record, parsed)
	}
	rerendered, err := RenderArchiveRecord(parsed)
	if err != nil || string(rerendered) != string(content) {
		t.Fatalf("render changed: %s err=%v", rerendered, err)
	}
}

func TestFailedQARecordRefusesAnotherVerdictOrAnOverride(t *testing.T) {
	t.Parallel()
	record := ArchiveRecord{Disposition: ArchiveFailedQA, QAVerdict: VerdictFail, QAReport: "qa-report-demo.md"}
	content, err := RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, old, replacement, want string
	}{
		{"pass verdict", "qa_verdict: fail\n", "qa_verdict: pass\n", "archive record disposition failed-qa requires qa_verdict fail"},
		{"empty verdict", "qa_verdict: fail\n", "qa_verdict: \"\"\n", "archive record disposition failed-qa requires qa_verdict fail"},
		{"absent verdict", "qa_verdict: fail\n", "", "archive record disposition failed-qa requires qa_verdict fail"},
		{"empty report", "qa_report: qa-report-demo.md\n", "qa_report: \"\"\n", "archive record disposition failed-qa requires qa_report"},
		{"absent report", "qa_report: qa-report-demo.md\n", "", "archive record disposition failed-qa requires qa_report"},
		{"true override", "qa_verdict: fail\n", "qa_verdict: fail\nqa_override: true\n", "archive record disposition failed-qa cannot carry qa_override"},
		{"false override", "qa_verdict: fail\n", "qa_verdict: fail\nqa_override: false\n", "archive record disposition failed-qa cannot carry qa_override"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bad := strings.Replace(string(content), tc.old, tc.replacement, 1)
			if bad == string(content) {
				t.Fatal("fixture mutation did not apply")
			}
			if _, err := ParseArchiveRecord([]byte(bad)); err == nil || err.Error() != tc.want {
				t.Fatalf("refusal=%v, want %q", err, tc.want)
			}
		})
	}
}

func TestArchivedTaskCompletedExcludesTheQATaskOfAFailedQA(t *testing.T) {
	t.Parallel()
	record := ArchiveRecord{Disposition: ArchiveFailedQA, QATask: "task_qa"}
	if ArchivedTaskCompleted(record, record.QATask) {
		t.Fatal("failed QA Task counted as completed")
	}
	if !ArchivedTaskCompleted(record, "task_01") {
		t.Fatal("implementation Task excluded from completion")
	}
	if ArchivedTaskCompleted(record, "") {
		t.Fatal("empty Task counted as completed")
	}
}

func TestArchiveStillRefusesAFailingQAWithoutAnOverride(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	writeFile(t, filepath.Join(dir, "qa", "qa-report-2026-10-04.md"), "---\nverdict: fail\n---\n\n# Failed QA\n")
	before := archiveLinksTree(t, req.RepositoryRoot)
	_, err := Archive(req)
	const want = `no passing QA verdict: newest QA Report verdict is "fail"; expected "pass"`
	if err == nil || err.Error() != want {
		t.Fatalf("refusal=%v, want %q", err, want)
	}
	if after := archiveLinksTree(t, req.RepositoryRoot); !reflect.DeepEqual(before, after) {
		t.Fatal("repository files changed on refusal")
	}
	if _, err := os.Stat(ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, req.BuiltInRoot), req.Slug)); !os.IsNotExist(err) {
		t.Fatalf("archive record exists: %v", err)
	}
}
