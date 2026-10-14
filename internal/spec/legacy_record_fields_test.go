package spec

// Suite: Legacy Archive Record fields.
// Invariant: legacy fallbacks preserve recorded values and active archive behavior.
// Boundary IN: disposable legacy folders and Git history.
// Boundary OUT: conversion records and archived Task completion.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func legacyFieldsFixture(t *testing.T, metadata string) LegacyConversionRequest {
	t.Helper()
	req := legacyFixture(t, "status: archived\n", "")
	writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "_prd.md"), "---\nspec: demo\nstatus: archived\n"+metadata+"---\n\n# Demo\n\nPreserved outcome.\n")
	return req
}

func TestLegacyOverrideWithoutApprovalRecordsNotRecorded(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, fields string }{
		{"absent", ""},
		{"empty", "qa_override_approval: \"\"\nqa_override_reason: \"\"\n"},
		{"blank", "qa_override_approval: \" \t \"\nqa_override_reason: \" \"\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := legacyFieldsFixture(t, "qa_override: true\n"+tc.fields)
			c := legacyPlan(t, req)
			want := &QAArchiveOverrideRecord{Approval: LegacyNotRecorded, Reason: LegacyNotRecorded}
			if c.Record.Disposition != ArchiveQAOverride || !reflect.DeepEqual(c.Record.QAOverride, want) {
				t.Fatalf("record = %+v, override = %+v, want %+v", c.Record, c.Record.QAOverride, want)
			}
		})
	}
}

func TestLegacySourcesFallBackToAdoptedReferences(t *testing.T) {
	t.Parallel()
	for _, index := range []string{"absent", "empty"} {
		t.Run(index, func(t *testing.T) {
			t.Parallel()
			req := legacyFieldsFixture(t, "")
			refs := filepath.Join(req.ArchiveRoot, req.Slug, "references")
			writeFile(t, filepath.Join(refs, "z-backlog.md"), "---\ntype: fix\nstatus: promoted\ncreated: 2026-01-01\n---\n\n# Backlog\n")
			writeFile(t, filepath.Join(refs, "a-finding.md"), "---\nstatus: done\ncreated_at: 2026-01-02\n---\n\n# Finding\n")
			writeFile(t, filepath.Join(refs, "note.md"), "# Plain note\n")
			writeFile(t, filepath.Join(refs, "ignored.txt"), "---\nstatus: done\ncreated_at: 2026-01-02\n---\n")
			writeFile(t, filepath.Join(refs, "nested", "finding.md"), "---\nstatus: done\ncreated_at: 2026-01-02\n---\n")
			writeFile(t, filepath.Join(refs, "sequence.md"), "---\n- status\n- created_at\n---\n")
			if index == "empty" {
				writeFile(t, filepath.Join(refs, "_index.md"), "| path |\n| --- |\n")
			}
			c := legacyPlan(t, req)
			if want := []string{"a-finding.md", "z-backlog.md"}; !reflect.DeepEqual(c.Record.Sources, want) {
				t.Fatalf("sources = %v, want %v", c.Record.Sources, want)
			}
		})
	}
}

func TestLegacyCreatedFallsBackToTheTechSpec(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, prd, tech, fallback, want string }{
		{"date", "", "created: 2026-01-02\n", "2025-12-31", "2026-01-02"},
		{"timestamp", "", "created: 2026-01-02T12:34:56Z\n", "2025-12-31", "2026-01-02"},
		{"PRD wins", "created: 2026-01-01\n", "created: 2026-01-02\n", "2025-12-31", "2026-01-01"},
		{"invalid date", "", "created: 2026-02-30\n", "2025-12-31", "2025-12-31"},
		{"short value", "", "created: unknown\n", "2025-12-31", "2025-12-31"},
		{"no TechSpec", "", "", "2025-12-31", "2025-12-31"},
		{"no recorded date", "", "", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := legacyFieldsFixture(t, tc.prd)
			if tc.tech != "" {
				writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "_techspec.md"), "---\n"+tc.tech+"---\n\n# Design\n")
			}
			req.Created = tc.fallback
			c := legacyPlan(t, req)
			if c.Record.Created != tc.want {
				t.Fatalf("created = %q, want %q", c.Record.Created, tc.want)
			}
			folder, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
			folderWant := tc.want
			if folderWant == tc.fallback {
				folderWant = "" // The folder reader does not read Git or the request.
			}
			if err != nil || folder.Record.Created != folderWant {
				t.Fatalf("folder created = %q, err = %v, want %q", folder.Record.Created, err, folderWant)
			}
		})
	}
}

func TestLegacyOverrideKeepsRecordedValues(t *testing.T) {
	t.Parallel()
	for _, status := range []string{"", "pending", "failed"} {
		t.Run("QA status "+status, func(t *testing.T) {
			t.Parallel()
			req := legacyFieldsFixture(t, "qa_override: true\nqa_override_approval: \"  maintainer\\tapproval  \"\nqa_override_reason: |\n  Kept reason.\n  Second line.\nqa_override_qa_outcome: fail\nqa_override_qa_task_status: \""+status+"\"\nqa_override_revision: "+strings.Repeat("b", 40)+"\n")
			c := legacyPlan(t, req)
			want := &QAArchiveOverrideRecord{Approval: "  maintainer\tapproval  ", Reason: "Kept reason.\nSecond line.\n", QAOutcome: "fail", QATaskStatus: status, Revision: strings.Repeat("b", 40)}
			if !reflect.DeepEqual(c.Record.QAOverride, want) {
				t.Fatalf("override = %+v, want %+v", c.Record.QAOverride, want)
			}
			// The completion reader's contract depends on the recorded status alone.
			c.Record.QATask = "task_qa"
			if ArchivedTaskCompleted(c.Record, "task_qa") != (status == "") || !ArchivedTaskCompleted(c.Record, "task_01") {
				t.Fatalf("completion changed for status %q", status)
			}
		})
	}
	for _, tc := range []struct{ name, fields, approval, reason string }{
		{"approval only", "qa_override_approval: \" maintainer \"\n", " maintainer ", LegacyNotRecorded},
		{"reason only", "qa_override_reason: \" accepted \"\n", LegacyNotRecorded, " accepted "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := legacyPlan(t, legacyFieldsFixture(t, "qa_override: true\n"+tc.fields))
			if c.Record.QAOverride.Approval != tc.approval || c.Record.QAOverride.Reason != tc.reason {
				t.Fatalf("override = %+v", c.Record.QAOverride)
			}
		})
	}
}

func TestLegacySourcesPreferTheIndex(t *testing.T) {
	t.Parallel()
	req := legacyFieldsFixture(t, "")
	refs := filepath.Join(req.ArchiveRoot, req.Slug, "references")
	writeFile(t, filepath.Join(refs, "adopted.md"), "---\ntype: fix\nstatus: promoted\ncreated: 2026-01-01\n---\n")
	writeFile(t, filepath.Join(refs, "_index.md"), "| title | path |\n| --- | --- |\n| Last | `z-note.md` |\n| First | `a-note.md` |\n")
	c := legacyPlan(t, req)
	if want := []string{"z-note.md", "a-note.md"}; !reflect.DeepEqual(c.Record.Sources, want) {
		t.Fatalf("sources = %v, want %v", c.Record.Sources, want)
	}
}

func TestFindLegacyCreationReadsTheOldestAddingCommit(t *testing.T) {
	t.Parallel()
	for _, location := range []string{"active then archived", "archive only"} {
		t.Run(location, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			gittest.InitRepo(t, root, "-b", "main")
			active := "docs/specs/demo/_prd.md"
			archived := "docs/history/specs/demo/_prd.md"
			first := active
			if location == "archive only" {
				first = archived
			}
			writeFile(t, filepath.Join(root, first), "---\nspec: demo\n---\n\n# Demo\n")
			gittest.Run(t, root, "add", "--all")
			gittest.Run(t, root, "commit", "--date=2026-01-02T12:00:00+00:00", "-m", "Create Spec")
			want := strings.TrimSpace(gittest.Run(t, root, "show", "-s", "--format=%as", "HEAD"))
			if location == "active then archived" {
				if err := os.MkdirAll(filepath.Dir(filepath.Join(root, archived)), 0o755); err != nil {
					t.Fatal(err)
				}
				gittest.Run(t, root, "mv", active, archived)
				gittest.Run(t, root, "commit", "--date=2026-02-03T12:00:00+00:00", "-m", "Archive Spec")
			}
			got, err := FindLegacyCreation(t.Context(), root, "docs/specs", "docs/history/specs", "demo")
			if err != nil || got != want {
				t.Fatalf("created = %q, err = %v, want %q", got, err, want)
			}
			got, err = FindLegacyCreation(t.Context(), root, "docs/specs", "docs/history/specs", "missing")
			if err != nil || got != "" {
				t.Fatalf("missing created = %q, err = %v", got, err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if _, err := FindLegacyCreation(ctx, root, "docs/specs", "docs/history/specs", "demo"); !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation = %v", err)
			}
		})
	}
	t.Run("Git failure", func(t *testing.T) {
		t.Parallel()
		_, err := FindLegacyCreation(t.Context(), t.TempDir(), "docs/specs", "docs/history/specs", "demo")
		if err == nil || !strings.Contains(err.Error(), "read legacy creation") || errors.Unwrap(err) == nil {
			t.Fatalf("Git failure = %v", err)
		}
	})
	t.Run("unsafe path", func(t *testing.T) {
		t.Parallel()
		if _, err := FindLegacyCreation(t.Context(), t.TempDir(), "../specs", "docs/history/specs", "demo"); err == nil {
			t.Fatal("unsafe path accepted")
		}
	})
}

func TestActiveArchiveRecordKeepsItsFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, metadata, created string
		indexed, override       bool
	}{
		{name: "no legacy fallbacks", metadata: "qa_override: true\nqa_override_reason: \" \"\n", override: true},
		{name: "recorded fields", metadata: "created: 2026-01-01\n", created: "2026-01-01", indexed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			folder := t.TempDir()
			writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\n"+tc.metadata+"---\n\n# Demo\n\nPreserved outcome.\n\n## Decisions\n\nADR-0247\n")
			writeFile(t, filepath.Join(folder, "_techspec.md"), "---\ncreated: 2026-02-02\n---\n")
			writeFile(t, filepath.Join(folder, "references", "adopted.md"), "---\nstatus: done\ncreated_at: 2026-01-01\n---\n")
			want := ArchiveRecord{Spec: "demo", Title: "Demo", Created: tc.created, Archived: "2026-03-01", Source: "docs/specs/demo", SourceRevision: strings.Repeat("a", 40), Disposition: ArchiveNoQA, Unproven: []string{"row 1"}, Promoted: []string{"docs/references/lesson.md"}, ADRs: []string{"ADR-0247"}, Outcome: "Preserved outcome."}
			if tc.override {
				want.Disposition = ArchiveQAOverride
				want.QAOverride = &QAArchiveOverrideRecord{Reason: " "}
			}
			if tc.indexed {
				writeFile(t, filepath.Join(folder, "references", "_index.md"), "| path |\n| --- |\n| z.md |\n| a.md |\n")
				want.Sources = []string{"z.md", "a.md"}
			}
			got, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder, Slug: "demo", Source: want.Source, SourceRevision: want.SourceRevision, Archived: want.Archived, Unproven: want.Unproven, Promoted: want.Promoted})
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("record = %+v, err = %v, want %+v", got, err, want)
			}
			rendered, err := RenderArchiveRecord(got)
			if err != nil {
				t.Fatal(err)
			}
			expected, err := RenderArchiveRecord(want)
			if err != nil || !bytes.Equal(rendered, expected) {
				t.Fatalf("active record bytes changed: %v", err)
			}
		})
	}
}
