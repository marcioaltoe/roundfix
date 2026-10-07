package spec

// Suite: Lenient Legacy Reading.
// Invariant: only legacy projection membership and retired types are tolerated.
// Boundary IN: synthetic Spec folders in temporary directories.
// Boundary OUT: conversion records, named tolerances, and unchanged source bytes.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const legacyProjectionManifest = `---
schema: ` + manifestSchema + `
qa: task_qa
graph:
  nodes:
    - id: task_01
      file: task_01.md
    # - id: task_02
    #   file: task_02.md
    # - id: task_03
    #   file: task_03.md
    - id: task_qa
      file: task_qa.md
      needs: [task_01]
---

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
`

func legacyProjectionFixture(t *testing.T, rows string) LegacyConversionRequest {
	t.Helper()
	root := t.TempDir()
	req := LegacyConversionRequest{RepositoryRoot: root, ArchiveRoot: filepath.Join(root, "docs/history/specs"), Slug: "demo", SourceRevision: strings.Repeat("a", 40)}
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\nstatus: archived\ncreated: 2026-09-01\narchived: 2026-09-02\n---\n\n# Demo\n\nPreserved outcome.\n")
	writeFile(t, filepath.Join(folder, "_tasks.md"), legacyProjectionManifest+rows)
	writeFile(t, filepath.Join(folder, "task_01.md"), "---\ntask: task_01\nspec: demo\nstatus: completed\ntype: backend\n---\n\n# Implementation\n")
	writeFile(t, filepath.Join(folder, "task_qa.md"), "---\ntask: task_qa\nspec: demo\nstatus: completed\ntype: qa\n---\n\n# QA\n")
	return req
}

func assertLegacyProjectionConversion(t *testing.T, req LegacyConversionRequest, want []string) {
	t.Helper()
	before := legacyTree(t, req)
	c := legacyPlan(t, req)
	legacyAssertWhole(t, req, before)
	if !reflect.DeepEqual(c.Tolerated, want) || c.Record.QATask != "task_qa" {
		t.Fatalf("tolerated=%v QA Task=%q", c.Tolerated, c.Record.QATask)
	}
	for _, tolerance := range want {
		if strings.Contains(string(c.Rendered), tolerance) {
			t.Fatalf("record contains tolerance %q", tolerance)
		}
	}
	if strings.Contains(string(c.Rendered), "tolerated:") {
		t.Fatal("record contains a tolerated field")
	}
	if err := ApplyLegacyConversion(req.RepositoryRoot, c); err != nil {
		t.Fatal(err)
	}
	archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil || archived.Form != ArchivedRecord || !reflect.DeepEqual(archived.Record, c.Record) {
		t.Fatalf("archive=%+v err=%v", archived, err)
	}
	if _, err := os.Stat(filepath.Join(req.ArchiveRoot, req.Slug)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("converted folder remains: %v", err)
	}
}

func TestLegacyConversionToleratesAProjectionRowOutsideTheGraph(t *testing.T) {
	t.Parallel()
	req := legacyProjectionFixture(t, "| task_03 | Removed first | refactor | low | - |\n| task_01 | Kept | backend | low | - |\n| task_02 | Removed second | backend | low | - |\n| task_qa | QA | qa | low | task_01 |\n")
	assertLegacyProjectionConversion(t, req, []string{
		"projection row task_03 names a Task outside the graph",
		"projection row task_02 names a Task outside the graph",
	})
}

func TestLegacyConversionToleratesARetiredTaskType(t *testing.T) {
	t.Parallel()
	req := legacyProjectionFixture(t, "| task_01 | Kept | refactor | low | - |\n| task_qa | QA | qa | low | task_01 |\n")
	assertLegacyProjectionConversion(t, req, []string{`projection row task_01 has retired Task type "refactor"`})
}

func TestLegacyConversionStillRefusesAMalformedProjectionRow(t *testing.T) {
	t.Parallel()
	req := legacyProjectionFixture(t, "| task_02 | Removed | backend | low |\n")
	before := legacyTree(t, req)
	manifestPath := filepath.Join(req.ArchiveRoot, req.Slug, "_tasks.md")
	want := (ManifestError{Path: manifestPath, Reason: "projection table has a malformed Task row; each row must start with a canonical task_NN id and include title, type, complexity, and needs cells"}).Error()
	_, err := PlanLegacyConversion(req)
	if err == nil || err.Error() != "build legacy record: "+want {
		t.Fatalf("refusal=%v, want %q", err, "build legacy record: "+want)
	}
	legacyAssertWhole(t, req, before)
}

func TestActiveSpecStillRefusesWhatLegacyReadingTolerates(t *testing.T) {
	t.Parallel()
	for _, retired := range []bool{false, true} {
		t.Run(fmt.Sprintf("retired=%t", retired), func(t *testing.T) {
			taskType := "backend"
			if retired {
				taskType = "refactor"
			}
			rows := "| task_02 | Removed | backend | low | - |\n| task_01 | Kept | " + taskType + " | low | - |\n| task_qa | QA | qa | low | task_01 |\n"
			req := legacyProjectionFixture(t, rows)
			folder := filepath.Join(req.ArchiveRoot, req.Slug)
			writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\nstatus: active\n---\n\n# Demo\n")
			manifestPath := filepath.Join(folder, "_tasks.md")
			want := (ManifestError{Path: manifestPath, Reason: `projection table row names unknown Task "task_02"; remove it or add a matching graph node`}).Error()
			if retired {
				want = (ManifestError{Path: manifestPath, Reason: fmt.Sprintf(`projection row for Task "task_01" has invalid type "refactor" (allowed: %s); update the _tasks.md type cell to one allowed value`, allowedTaskTypeValues()), Err: TaskTypeError{Path: manifestPath, Value: "refactor"}}).Error()
			}
			_, loadErr := Load(req.ArchiveRoot, req.Slug)
			_, graphErr := ReadCauseGraph(req.ArchiveRoot, req.Slug)
			_, atErr := ReadCauseGraphAt(folder, req.Slug, os.ReadFile)
			_, buildErr := BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder, Slug: req.Slug})
			for name, err := range map[string]error{"Load": loadErr, "ReadCauseGraph": graphErr, "ReadCauseGraphAt": atErr, "BuildArchiveRecord": buildErr} {
				if err == nil || err.Error() != want {
					t.Errorf("%s refusal=%v, want %q", name, err, want)
				}
			}
		})
	}
}

func TestReadArchivedSpecReadsALegacyFolderLeniently(t *testing.T) {
	t.Parallel()
	req := legacyProjectionFixture(t, "| task_02 | Removed | refactor | low | - |\n| task_01 | Kept | refactor | low | - |\n| task_qa | QA | qa | low | task_01 |\n")
	before := legacyTree(t, req)
	archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil || archived.Form != ArchivedFolder || archived.Record.QATask != "task_qa" || archived.Record.Disposition != ArchiveNoQA {
		t.Fatalf("archive=%+v err=%v", archived, err)
	}
	graph, tolerated, err := ReadLegacyCauseGraph(req.ArchiveRoot, req.Slug)
	want := []string{"projection row task_02 names a Task outside the graph", `projection row task_01 has retired Task type "refactor"`}
	if err != nil || !reflect.DeepEqual(tolerated, want) || !reflect.DeepEqual(graph.Tasks, map[string]string{"task_01": "task_01.md", "task_qa": "task_qa.md"}) {
		t.Fatalf("graph=%+v tolerated=%v err=%v", graph, tolerated, err)
	}
	legacyAssertWhole(t, req, before)
}

func TestLegacyCauseGraphPreservesManifestRefusals(t *testing.T) {
	t.Parallel()
	valid := legacyProjectionManifest + "| task_01 | Kept | backend | low | - |\n"
	cases := []struct {
		name     string
		manifest string
	}{
		{name: "missing front matter", manifest: "# Graph\n"},
		{name: "invalid YAML", manifest: strings.Replace(valid, "qa: task_qa", "qa: [", 1)},
		{name: "schema", manifest: strings.Replace(valid, manifestSchema, "unknown/v1", 1)},
		{name: "requires", manifest: strings.Replace(valid, "qa: task_qa", "requires: [demo]\nqa: task_qa", 1)},
		{name: "QA declaration", manifest: strings.Replace(valid, "qa: task_qa", "qa: declined", 1)},
		{name: "no nodes", manifest: "---\nschema: " + manifestSchema + "\ngraph:\n  nodes: []\n---\n"},
		{name: "missing id", manifest: strings.Replace(valid, "id: task_01", "id: ''", 1)},
		{name: "missing file", manifest: strings.Replace(valid, "file: task_01.md", "file: ''", 1)},
		{name: "duplicate node", manifest: strings.Replace(valid, "id: task_qa", "id: task_01", 1)},
		{name: "unknown need", manifest: strings.Replace(valid, "needs: [task_01]", "needs: [task_02]", 1)},
		{name: "malformed row", manifest: valid + "| task_02 | Removed | backend | low |\n"},
		{name: "duplicate known row", manifest: valid + "| task_01 | Again | refactor | low | - |\n"},
		{name: "duplicate outside row", manifest: valid + "| task_02 | Removed | backend | low | - |\n| task_02 | Again | refactor | low | - |\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := legacyProjectionFixture(t, "")
			writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "_tasks.md"), tc.manifest)
			before := legacyTree(t, req)
			_, strictErr := ReadCauseGraph(req.ArchiveRoot, req.Slug)
			_, _, legacyErr := ReadLegacyCauseGraph(req.ArchiveRoot, req.Slug)
			if strictErr == nil || legacyErr == nil || strictErr.Error() != legacyErr.Error() {
				t.Fatalf("strict=%v legacy=%v", strictErr, legacyErr)
			}
			legacyAssertWhole(t, req, before)
		})
	}
}
