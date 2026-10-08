package spec

// Suite: Legacy archive conversion.
// Invariant: conversion preserves metadata and Git-backed bytes; planning and refusals preserve the tree.
// Boundary IN: disposable repositories created with gittest.
// Boundary OUT: filesystem records, promotions, and first-parent Git history.

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

func legacyFixture(t *testing.T, stamp, report string) LegacyConversionRequest {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	req := LegacyConversionRequest{RepositoryRoot: root, ArchiveRoot: filepath.Join(root, "docs/history/specs"), Slug: "demo"}
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\ncreated: 2026-09-01\n"+stamp+"---\n\n# Demo\n\nA preserved outcome.\n\n## Decisions\n\nADR-0247\n")
	writeFile(t, filepath.Join(folder, "notes.md"), "Measurements\n\x00binary bytes\n")
	if report != "" {
		writeFile(t, filepath.Join(folder, "qa/qa-report-2026-09-02.md"), report)
	}
	req.SourceRevision = legacyCommit(t, root, "Archive demo (#42)")
	return req
}
func legacyCommit(t *testing.T, root, subject string) string {
	t.Helper()
	gittest.Run(t, root, "add", "--all")
	gittest.Run(t, root, "commit", "-m", subject)
	return strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
}
func legacyPlan(t *testing.T, req LegacyConversionRequest) LegacyConversion {
	t.Helper()
	c, err := PlanLegacyConversion(req)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func legacyTree(t *testing.T, req LegacyConversionRequest) map[string]string {
	t.Helper()
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	tree := map[string]string{}
	for name, content := range archiveLinksTree(t, folder) {
		relative, err := filepath.Rel(folder, name)
		if err != nil {
			t.Fatal(err)
		}
		tree[filepath.ToSlash(relative)] = content
	}
	return tree
}
func legacyAssertWhole(t *testing.T, req LegacyConversionRequest, before map[string]string) {
	t.Helper()
	if got := legacyTree(t, req); !reflect.DeepEqual(got, before) {
		t.Fatal("legacy folder changed")
	}
	if _, err := os.Lstat(filepath.Join(req.ArchiveRoot, req.Slug+".md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("record exists: %v", err)
	}
}
func TestLegacyArchiveFoldersSkipRecords(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\n", "")
	writeFile(t, filepath.Join(req.ArchiveRoot, "record.md"), "a record")
	writeFile(t, filepath.Join(req.ArchiveRoot, "empty/notes.md"), "no PRD")
	writeFile(t, filepath.Join(req.ArchiveRoot, "aaa/_prd.md"), "status: active")
	got, err := LegacyArchiveFolders(req.ArchiveRoot)
	if err != nil || !reflect.DeepEqual(got, []string{"aaa", "demo"}) {
		t.Fatalf("folders=%v err=%v", got, err)
	}
	missing, err := LegacyArchiveFolders(filepath.Join(req.RepositoryRoot, "missing"))
	if err != nil || len(missing) != 0 {
		t.Fatalf("missing=%v err=%v", missing, err)
	}
}
func TestLegacyDeliveryNamesTheCommitThatRetiredTheSpec(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"removed from Spec Root", "added directly", "merge retirement", "newest retirement"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			gittest.InitRepo(t, root, "-b", "main")
			active := filepath.Join(root, "docs/specs/demo/_prd.md")
			archived := filepath.Join(root, "docs/history/specs/demo/_prd.md")
			if kind != "added directly" {
				writeFile(t, active, "original")
				legacyCommit(t, root, "Create")
			}
			if kind == "merge retirement" {
				gittest.Run(t, root, "checkout", "-b", "feature")
			}
			writeFile(t, archived, "archived")
			if kind != "added directly" {
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
			}
			commit := legacyCommit(t, root, "Retire demo (#123)")
			if kind == "newest retirement" {
				writeFile(t, active, "reintroduced")
				legacyCommit(t, root, "Reintroduce")
				if err := os.Remove(active); err != nil {
					t.Fatal(err)
				}
				commit = legacyCommit(t, root, "Retire again (#123)")
			}
			if kind == "merge retirement" {
				gittest.Run(t, root, "checkout", "main")
				gittest.Run(t, root, "merge", "--no-ff", "feature", "-m", "Retire merge (#123)")
				commit = strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
			}
			writeFile(t, filepath.Join(root, "later.md"), "later")
			legacyCommit(t, root, "Later (#999)")
			got, err := FindLegacyDelivery(context.Background(), root, "docs/specs", "docs/history/specs", "demo")
			date := strings.TrimSpace(gittest.Run(t, root, "show", "-s", "--format=%as", commit))
			if err != nil || got != (LegacyDelivery{Commit: commit, PullRequest: "123", Date: date}) {
				t.Fatalf("delivery=%+v err=%v", got, err)
			}
		})
	}
}
func TestLegacyDeliveryIgnoresARelocation(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	old := filepath.Join(root, "old/archive/demo/_prd.md")
	writeFile(t, old, "same bytes")
	legacyCommit(t, root, "Old archive (#7)")
	writeFile(t, filepath.Join(root, "docs/history/specs/demo/_prd.md"), "same bytes")
	if err := os.Remove(old); err != nil {
		t.Fatal(err)
	}
	legacyCommit(t, root, "Relocate (#8)")
	got, err := FindLegacyDelivery(context.Background(), root, "docs/specs", "docs/history/specs", "demo")
	if err != nil || got != (LegacyDelivery{}) {
		t.Fatalf("delivery=%+v err=%v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FindLegacyDelivery(ctx, root, "docs/specs", "docs/history/specs", "demo"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}
func TestLegacyConversionWritesARecordGitStillHolds(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\narchived: 2026-09-02\n", "---\nverdict: pass\nrows_total: 1\nrows_passed: 1\nrows_blocked: 0\n---\n\n# QA\n")
	req.Delivery = LegacyDelivery{Commit: req.SourceRevision, PullRequest: "42", Date: "2026-09-03"}
	req.Promote = []string{"docs/history/specs/demo/notes.md"}
	before := legacyTree(t, req)
	c := legacyPlan(t, req)
	legacyAssertWhole(t, req, before)
	if len(c.Rendered) > ArchiveRecordTargetBytes || c.Record.Source != "docs/history/specs/demo" || c.Record.SourceRevision != req.SourceRevision || c.Record.Disposition != ArchivePass || c.Record.Archived != "2026-09-02" || c.Record.PullRequest != "42" || c.Record.DeliveryCommit != req.SourceRevision {
		t.Fatalf("record=%+v size=%d", c.Record, len(c.Rendered))
	}
	var size int64
	for name, content := range before {
		size += int64(len(content))
		if got := gittest.Run(t, req.RepositoryRoot, "show", req.SourceRevision+":"+c.Folder+"/"+name); got != content {
			t.Fatalf("Git differs for %s", name)
		}
	}
	if c.Bytes != size || len(c.Files) != len(before) {
		t.Fatalf("inventory=%+v", c)
	}
	if err := ApplyLegacyConversion(req.RepositoryRoot, c); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(req.RepositoryRoot, filepath.FromSlash(c.RecordPath)))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseArchiveRecord(raw)
	if err != nil || !reflect.DeepEqual(parsed, c.Record) {
		t.Fatalf("parsed=%+v err=%v", parsed, err)
	}
	if _, err := os.Stat(filepath.Join(req.ArchiveRoot, req.Slug)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("folder remains: %v", err)
	}
	promoted, err := os.ReadFile(filepath.Join(req.RepositoryRoot, "docs/references/notes.md"))
	if err != nil || string(promoted) != before["notes.md"] {
		t.Fatalf("promotion=%q err=%v", promoted, err)
	}
	for name, content := range before {
		if got := gittest.Run(t, req.RepositoryRoot, "show", parsed.SourceRevision+":"+parsed.Source+"/"+name); got != content {
			t.Fatalf("removed Git blob differs for %s", name)
		}
	}
}
func TestLegacyConversionOfAFolderWithoutQAIsNoQA(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\n", "")
	c := legacyPlan(t, req)
	if c.Record.Disposition != ArchiveNoQA || c.Record.QAVerdict != "" || c.Record.QAReport != "" {
		t.Fatalf("record=%+v", c.Record)
	}
}
func TestLegacyConversionKeepsEveryOverrideField(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\nqa_override: true\nqa_override_approval: maintainer\nqa_override_reason: Missing QA accepted\nqa_override_qa_outcome: absent\nqa_override_qa_task_status: pending\nqa_override_revision: "+strings.Repeat("a", 40)+"\n", "")
	c := legacyPlan(t, req)
	want := &QAArchiveOverrideRecord{Approval: "maintainer", Reason: "Missing QA accepted", QAOutcome: "absent", QATaskStatus: "pending", Revision: strings.Repeat("a", 40)}
	if c.Record.Disposition != ArchiveQAOverride || !reflect.DeepEqual(c.Record.QAOverride, want) || !bytes.Contains(c.Rendered, []byte("qa_override: true")) {
		t.Fatalf("record=%+v", c.Record)
	}
}
func TestLegacyConversionAcceptsAPreStampFolder(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: active\n", "")
	req.Delivery = LegacyDelivery{Commit: req.SourceRevision, Date: "2026-09-02"}
	c := legacyPlan(t, req)
	if c.Record.Disposition != ArchiveNoQA || c.Record.Archived != req.Delivery.Date {
		t.Fatalf("record=%+v", c.Record)
	}
}
func TestLegacyConversionPromotionRefusals(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"outside", "absolute", "traversal", "directory", "missing", "symlink", "symlink parent", "core", "task", "report", "duplicate", "occupied", "destination symlink"} {
		t.Run(kind, func(t *testing.T) {
			req := legacyFixture(t, "status: archived\n", "")
			folder := filepath.Join(req.ArchiveRoot, req.Slug)
			candidate := "docs/history/specs/demo/notes.md"
			switch kind {
			case "outside":
				candidate = "other.md"
			case "absolute":
				candidate = filepath.Join(folder, "notes.md")
			case "traversal":
				candidate = "docs/history/specs/demo/../demo/notes.md"
			case "directory":
				candidate = "docs/history/specs/demo/qa"
				if err := os.MkdirAll(filepath.Join(folder, "qa"), 0o755); err != nil {
					t.Fatal(err)
				}
			case "missing":
				candidate = "docs/history/specs/demo/missing.md"
			case "symlink":
				candidate = "docs/history/specs/demo/link.md"
				if err := os.Symlink("notes.md", filepath.Join(folder, "link.md")); err != nil {
					t.Fatal(err)
				}
			case "symlink parent":
				candidate = "docs/history/specs/demo/link/notes.md"
				if err := os.Symlink(".", filepath.Join(folder, "link")); err != nil {
					t.Fatal(err)
				}
			case "core":
				candidate = "docs/history/specs/demo/_prd.md"
			case "task":
				candidate = "docs/history/specs/demo/task_01.md"
				writeFile(t, filepath.Join(folder, "task_01.md"), "task")
			case "report":
				candidate = "docs/history/specs/demo/qa/qa-report-extra.md"
				writeFile(t, filepath.Join(folder, "qa/qa-report-extra.md"), "report")
			case "duplicate":
				writeFile(t, filepath.Join(folder, "nested/notes.md"), "other")
				req.Promote = []string{"docs/history/specs/demo/nested/notes.md"}
			case "occupied":
				writeFile(t, filepath.Join(req.RepositoryRoot, "docs/references/notes.md"), "occupied")
			case "destination symlink":
				if err := os.Symlink(folder, filepath.Join(req.RepositoryRoot, "docs/references")); err != nil {
					t.Fatal(err)
				}
			}
			req.Promote = append(req.Promote, candidate)
			before := legacyTree(t, req)
			if _, err := PlanLegacyConversion(req); err == nil {
				t.Fatal("promotion accepted")
			}
			legacyAssertWhole(t, req, before)
		})
	}
}
func TestLegacyConversionLeavesTheFolderWholeOnFailure(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\n", "")
	writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "a.md"), "first copy")
	req.Promote = []string{"docs/history/specs/demo/a.md", "docs/history/specs/demo/notes.md"}
	c := legacyPlan(t, req)
	before := legacyTree(t, req)
	// The second copy fails after the record and the first copy have been written.
	writeFile(t, filepath.Join(req.RepositoryRoot, "docs/references/notes.md"), "another writer")
	if err := ApplyLegacyConversion(req.RepositoryRoot, c); err == nil {
		t.Fatal("apply accepted occupied destination")
	}
	legacyAssertWhole(t, req, before)
	if _, err := os.Stat(filepath.Join(req.RepositoryRoot, "docs/references/a.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("first promotion remains: %v", err)
	}
	if got := archiveTestReadFile(t, filepath.Join(req.RepositoryRoot, "docs/references/notes.md")); got != "another writer" {
		t.Fatal("occupied destination changed")
	}
}
func TestNoQARecordRoundTrips(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\n", "")
	c := legacyPlan(t, req)
	rendered, err := RenderArchiveRecord(c.Record)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseArchiveRecord(rendered)
	if err != nil || !reflect.DeepEqual(c.Record, parsed) || parsed.Disposition != ArchiveNoQA {
		t.Fatalf("parsed=%+v err=%v", parsed, err)
	}
	if !ArchivedTaskCompleted(parsed, "task_01") {
		t.Fatal("no-qa record did not retain task completion contract")
	}
}

func TestLegacyConversionRefusesMetadataThatDoesNotRoundTrip(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\nqa_override: true\nqa_override_approval: maintainer\n", "")
	writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, SupersessionFilename), "---\nsuperseded_by: successor\ndate: 2026-09-02\nreason: Delivered elsewhere\n---\n\nDelivered elsewhere.\n")
	// Both independent metadata branches exist, but the parser cannot retain
	// override metadata under a superseded disposition. Refuse rather than lose it.
	before := legacyTree(t, req)
	if _, err := PlanLegacyConversion(req); err == nil || !strings.Contains(err.Error(), "round-trip") {
		t.Fatalf("round-trip refusal: %v", err)
	}
	legacyAssertWhole(t, req, before)
}

func TestLegacyConversionPreservesTheRendererOutcomeBudget(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\n", "")
	prd := filepath.Join(req.ArchiveRoot, req.Slug, "_prd.md")
	writeFile(t, prd, "---\nspec: demo\nstatus: archived\n---\n\n# Demo\n\n"+strings.Repeat("Measured outcome. ", 300)+"\n")
	c := legacyPlan(t, req)
	if len(c.Rendered) > ArchiveRecordTargetBytes || c.Record.Outcome == "" {
		t.Fatalf("record size=%d outcome=%q", len(c.Rendered), c.Record.Outcome)
	}
	rendered, err := RenderArchiveRecord(c.Record)
	if err != nil || !bytes.Equal(rendered, c.Rendered) {
		t.Fatalf("budgeted round trip: %v", err)
	}
}
