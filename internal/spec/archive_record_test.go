package spec

// Suite: Archive Records and retirement.
// Invariant: metadata survives the cut and refusals preserve the source tree.
// Boundary IN: Spec files and archive requests.
// Boundary OUT: Git provenance and command confirmation are covered in internal/cli.

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"strings"
	"testing"
)

func recordFixture(t *testing.T) (ArchiveRequest, string) {
	req, dir := archiveLinksFixture(t, true, "\n## Decisions\n\nADR-0247 and ADR-0230, ADR-0247 again.\n")
	req.SourceRevision = strings.Repeat("a", 40)
	req.RepositoryRoot = filepath.Dir(filepath.Dir(req.SpecsRoot))
	writeFile(t, filepath.Join(dir, "references", "_index.md"), "| title | path |\n| --- | --- |\n| Adopted | adopted.md |\n")
	writeFile(t, filepath.Join(dir, "references", "adopted.md"), "Source text\n")
	writeFile(t, filepath.Join(dir, "_authorization.md"), "## Sanctioned regeneration\n\n```yaml\ncommand: make example\noutputs: [generated/example.json]\n```\n")
	return req, dir
}
func readRecordResult(t *testing.T, result ArchiveResult) ArchiveRecord {
	t.Helper()
	content, err := os.ReadFile(result.RecordPath)
	if err != nil {
		t.Fatal(err)
	}
	record, err := ParseArchiveRecord(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) > ArchiveRecordTargetBytes {
		t.Fatalf("record size %d", len(content))
	}
	if _, err := os.Stat(result.SourceDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("source remains: %v", err)
	}
	return record
}
func TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	before := archiveLinksTree(t, dir)
	var size int64
	for _, text := range before {
		size += int64(len(text))
	}
	result, err := Archive(req)
	if err != nil {
		t.Fatal(err)
	}
	r := readRecordResult(t, result)
	if r.Disposition != ArchivePass || r.SourceRevision != req.SourceRevision || r.Source != "docs/specs/demo" || !reflect.DeepEqual(r.ADRs, []string{"ADR-0247", "ADR-0230"}) || !reflect.DeepEqual(r.Sources, []string{"adopted.md"}) || len(r.Regeneration) != 1 || !reflect.DeepEqual(r.Regeneration[0].Outputs, []string{"generated/example.json"}) || result.RemovedFiles != len(before) || result.RemovedBytes != size {
		t.Fatalf("record=%+v result=%+v", r, result)
	}
}
func TestArchiveRecordKeepsEveryOverrideField(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	writeFile(t, filepath.Join(dir, "task_qa.md"), taskFixture("task_qa", "QA", "failed", "qa", defaultVerificationSection))
	writeFile(t, filepath.Join(dir, "qa", "qa-report-2026-10-04.md"), "---\nverdict: fail\n---\n\n# QA\n")
	req.QAOverride = &QAArchiveOverride{Approval: "maintainer", Reason: strings.Repeat("r", 300), Revision: req.SourceRevision}
	result, err := Archive(req)
	if err != nil {
		t.Fatal(err)
	}
	r := readRecordResult(t, result)
	want := &QAArchiveOverrideRecord{"maintainer", strings.Repeat("r", 300), "fail", "failed", req.SourceRevision}
	if r.Disposition != ArchiveQAOverride || !reflect.DeepEqual(r.QAOverride, want) {
		t.Fatalf("override=%+v", r)
	}
}
func TestArchiveRecordOfASupersededSpec(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	if err := os.Remove(filepath.Join(dir, "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteSupersession(dir, "successor", "Delivered elsewhere", req.ArchivedAt); err != nil {
		t.Fatal(err)
	}
	result, err := Archive(req)
	if err != nil {
		t.Fatal(err)
	}
	r := readRecordResult(t, result)
	if r.Disposition != ArchiveSuperseded || r.SupersededBy != "successor" || r.Outcome != "Delivered elsewhere" {
		t.Fatalf("record=%+v", r)
	}
}
func TestArchiveRecordStaysWithinTheTargetSize(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	writeFile(t, filepath.Join(dir, "qa", "qa-report-2026-10-04.md"), "---\nverdict: pass\nrows_total: 1\nrows_passed: 1\nrows_blocked: 0\n---\n\n# QA\n\n## Outcome\n\n"+strings.Repeat("é", 2500)+"\n")
	result, err := Archive(req)
	if err != nil {
		t.Fatal(err)
	}
	r := readRecordResult(t, result)
	if !strings.HasSuffix(r.Outcome, "…") || r.SourceRevision != req.SourceRevision {
		t.Fatalf("record=%+v", r)
	}
}
func TestArchiveRecordRoundTrips(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	r, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: dir, Slug: req.Slug, Source: "docs/specs/demo", SourceRevision: req.SourceRevision, Archived: "2026-10-06"})
	if err != nil {
		t.Fatal(err)
	}
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParseArchiveRecord(content)
	if err != nil {
		t.Fatal(err)
	} // Empty lists normalize to YAML sequences.
	r.Regeneration = append(r.Regeneration, ArchiveRegeneration{Command: "empty", Outputs: []string{}})
	content, err = RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	got, err = ParseArchiveRecord(content)
	if err != nil {
		t.Fatal(err)
	}
	r.Unproven = []string{}
	r.Promoted = []string{}
	if !reflect.DeepEqual(r, got) {
		t.Fatalf("round trip\n%+v\n%+v", r, got)
	}
	for _, bad := range []string{strings.Replace(string(content), ArchiveRecordSchema, "unknown", 1), strings.Replace(string(content), "spec: demo\n", "", 1)} {
		if _, err := ParseArchiveRecord([]byte(bad)); err == nil {
			t.Fatal("accepted malformed record")
		}
	}
}
func TestArchiveRefusesBeforeAnyFileChanges(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"empty revision", "invalid revision", "record exists", "folder exists", "outside", "directory", "symlink", "core", "task", "report", "occupied promotion"} {
		t.Run(kind, func(t *testing.T) {
			req, dir := recordFixture(t)
			switch kind {
			case "empty revision":
				req.SourceRevision = ""
			case "invalid revision":
				req.SourceRevision = "xyz"
			case "record exists":
				writeFile(t, ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug), "occupied")
			case "folder exists":
				writeFile(t, filepath.Join(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug, "_prd.md"), "occupied")
			case "outside":
				req.Promote = []string{"../other.md"}
			case "directory":
				req.Promote = []string{"references"}
			case "symlink":
				if err := os.Symlink("references/adopted.md", filepath.Join(dir, "link.md")); err != nil {
					t.Fatal(err)
				}
				req.Promote = []string{"link.md"}
			case "core":
				req.Promote = []string{"_prd.md"}
			case "task":
				req.Promote = []string{"task_01.md"}
			case "report":
				req.Promote = []string{"qa/qa-report-2026-10-04.md"}
			case "occupied promotion":
				writeFile(t, filepath.Join(req.RepositoryRoot, "docs", "references", "adopted.md"), "occupied")
				req.Promote = []string{"references/adopted.md"}
			}
			before := archiveLinksTree(t, req.RepositoryRoot)
			if _, err := Archive(req); err == nil {
				t.Fatal("accepted refusal case")
			}
			if after := archiveLinksTree(t, req.RepositoryRoot); !reflect.DeepEqual(before, after) {
				t.Fatal("tree changed on refusal")
			}
		})
	}
}
func TestReadArchivedSpecReadsRecordAndLegacyFolder(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	root := ArchiveSpecRoot(req.SpecsRoot, true)
	if _, err := ReadArchivedSpec(root, "demo"); !errors.Is(err, ErrNotArchived) {
		t.Fatal(err)
	}
	legacy := filepath.Join(root, "demo")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(dir, legacy); err != nil {
		t.Fatal(err)
	}
	if err := stampArchiveMetadata(filepath.Join(legacy, "_prd.md"), "demo", "2026-10-06", nil, nil, "", ""); err != nil {
		t.Fatal(err)
	}
	folder, err := ReadArchivedSpec(root, "demo")
	if err != nil || folder.Form != ArchivedFolder {
		t.Fatalf("folder=%+v err=%v", folder, err)
	}
	bytes, err := RenderArchiveRecord(folder.Record)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, ArchiveRecordPath(root, "demo"), string(bytes))
	if _, err := ReadArchivedSpec(root, "demo"); err == nil {
		t.Fatal("accepted both forms")
	}
	if err := os.RemoveAll(legacy); err != nil {
		t.Fatal(err)
	}
	record, err := ReadArchivedSpec(root, "demo")
	if err != nil || record.Form != ArchivedRecord {
		t.Fatalf("record=%+v err=%v", record, err)
	}
	slugs, err := ArchivedSpecSlugs(root)
	if err != nil || !reflect.DeepEqual(slugs, []string{"demo"}) {
		t.Fatalf("slugs=%v err=%v", slugs, err)
	}
	writeFile(t, ArchiveRecordPath(root, "other"), string(bytes))
	if _, err := ReadArchivedSpec(root, "other"); err == nil {
		t.Fatal("accepted different file stem")
	}
}

func archiveFixtureRevision(t *testing.T, specsRoot string) string {
	t.Helper()
	root := filepath.Dir(filepath.Dir(specsRoot))
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "add", "-A")
	gittest.Run(t, root, "commit", "-m", "docs: archive evidence fixture")
	return strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
}
func assertArchiveEvidenceInGit(t *testing.T, specsRoot string, result ArchiveResult, before map[string]string) {
	t.Helper()
	r, err := ParseArchiveRecord([]byte(archiveTestReadFile(t, result.RecordPath)))
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(filepath.Dir(specsRoot))
	for path, want := range before {
		got := gittest.Run(t, root, "show", r.SourceRevision+":"+r.Source+"/"+path)
		if got != want {
			t.Fatalf("QA evidence changed for %s", path)
		}
	}
}
func TestArchiveRollsBackTheRecordWhenPromotionCannotBeWritten(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	req.Promote = []string{"references/adopted.md"}
	if err := os.Symlink(filepath.Join(req.RepositoryRoot, "missing-directory"), filepath.Join(req.RepositoryRoot, "docs", "references")); err != nil {
		t.Fatal(err)
	}
	before := archiveLinksTree(t, dir)
	if _, err := Archive(req); err == nil {
		t.Fatal("accepted unwritable promotion destination")
	}
	if after := archiveLinksTree(t, dir); !reflect.DeepEqual(before, after) {
		t.Fatal("source changed")
	}
	if _, err := os.Stat(ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("record remains: %v", err)
	}
}

func TestArchiveRecordRefusesAnOverrideDispositionWithoutATrueOverride(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	writeFile(t, filepath.Join(dir, "task_qa.md"), taskFixture("task_qa", "QA", "failed", "qa", defaultVerificationSection))
	writeFile(t, filepath.Join(dir, "qa", "qa-report-2026-10-04.md"), "---\nverdict: fail\n---\n\n# QA\n")
	req.QAOverride = &QAArchiveOverride{Approval: "maintainer", Reason: strings.Repeat("r", 300), Revision: req.SourceRevision}
	result, err := Archive(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(result.RecordPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "qa_override: true\n") {
		t.Fatalf("record lacks qa_override: true:\n%s", content)
	}
	forged := strings.Replace(string(content), "qa_override: true\n", "qa_override: false\n", 1)
	if _, err := ParseArchiveRecord([]byte(forged)); err == nil || !strings.Contains(err.Error(), "qa_override must be true") {
		t.Fatalf("forged override accepted: %v", err)
	}
}
