package spec

// Suite: cancellation records and Source Return.
// Invariant: cancellation keeps Task facts and source bytes; refusals write nothing.
// Boundary IN: temporary Spec folders, adopted sources and committed live citations.
// Boundary OUT: CLI ownership and delivery refusals belong to later Tasks.

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
)

const cancelBacklog = "---\ntype: feat # intent\nstatus: promoted # routing\ncreated: 2026-10-10\nspec: demo  # adopted owner\nreason: null\ncustom: 'preserve: this'\n---\n\n# A consumer\n\nAn unchanged observation.\n\n```yaml\nstatus: promoted\nspec: demo\n```\n"
const cancelFinding = "---\nstatus: done # routing\ncreated_at: 2026-10-10\nupdated_at: 2026-10-11 # latest evidence\nabsorbed_by: demo\ncustom: 'preserve: this'\n---\n\n# Investigation\n\nOriginal evidence stays exactly here.\n\n## Addendum — 2026-10-11 — Routed\n\nEarlier routing evidence stays.\n"

func cancelFixture(t *testing.T, tasks int) (CancelRequest, string) {
	t.Helper()
	root := t.TempDir()
	req := CancelRequest{SpecsRoot: filepath.Join(root, "docs", "specs"), Slug: "demo", Reason: "  no\n consumer  ", SourceRevision: strings.Repeat("a", 40), RepositoryRoot: root, BuiltInRoot: true, CancelledAt: time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)}
	dir := filepath.Join(req.SpecsRoot, req.Slug)
	writeFile(t, filepath.Join(dir, "_prd.md"), prdFixture("active")+"\nOriginal intent on the first line\nand continued here.\n\nSecond paragraph.\n")
	if tasks > 0 {
		var nodes strings.Builder
		for i := 1; i <= tasks; i++ {
			id := fmt.Sprintf("task_%02d", i)
			status := "pending"
			if i <= 5 {
				status = "completed"
			}
			fmt.Fprintf(&nodes, "    - id: %s\n      file: %s.md\n      needs: []\n", id, id)
			writeFile(t, filepath.Join(dir, id+".md"), taskFixture(id, "Build", status, "backend", defaultVerificationSection))
		}
		writeFile(t, filepath.Join(dir, "_tasks.md"), manifestFixtureWithQA("spec-tasks/v1", "qa: declined\nqa_reason: fixture\n", nodes.String(), ""))
	}
	cancelAdopt(t, dir, "backlog", "entry.md", cancelBacklog)
	return req, dir
}

func cancelAdopt(t *testing.T, dir, kind, name, content string) {
	t.Helper()
	path := filepath.Join(dir, "references", "_index.md")
	index, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		index = []byte("| source | type | owner | adopted date | path |\n| --- | --- | --- | --- | --- |\n")
	} else if err != nil {
		t.Fatal(err)
	}
	index = append(index, []byte(fmt.Sprintf("| original/%s | %s | demo | 2026-10-10 | `%s` |\n", name, kind, name))...)
	writeFile(t, path, string(index))
	writeFile(t, filepath.Join(dir, "references", name), content)
}

func cancelReadRecord(t *testing.T, result CancelResult) ArchiveRecord {
	t.Helper()
	content, err := os.ReadFile(result.RecordPath)
	if err != nil {
		t.Fatal(err)
	}
	record, err := ParseArchiveRecord(content)
	if err != nil {
		t.Fatal(err)
	}
	next, err := RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, next) {
		t.Fatal("cancelled record did not round-trip byte for byte")
	}
	return record
}

// Include directories, modes and symlink targets in the no-write oracle.
func cancelTree(t *testing.T, root string) map[string]string {
	t.Helper()
	tree := make(map[string]string)
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// Git's process-owned internals are not part of the cancellation write boundary.
		if rel == ".git" {
			return filepath.SkipDir
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String() + ":"
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			value += target
		} else if !entry.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += string(data)
		}
		tree[filepath.ToSlash(rel)] = value
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return tree
}

func TestCancelRecordsCompletedAndPendingTasks(t *testing.T) {
	t.Parallel()
	req, dir := cancelFixture(t, 8)
	before := archiveLinksTree(t, dir)
	var size int64
	for _, content := range before {
		size += int64(len(content))
	}
	adopted := []byte(before[filepath.Join(dir, "references", "entry.md")])
	result, err := Cancel(req)
	if err != nil {
		t.Fatal(err)
	}
	record := cancelReadRecord(t, result)
	completed := []string{"task_01", "task_02", "task_03", "task_04", "task_05"}
	pending := []string{"task_06", "task_07", "task_08"}
	if record.Disposition != ArchiveCancelled || record.CancelReason != "no consumer" || record.Archived != "2026-10-12" || record.SourceRevision != req.SourceRevision || record.Source != "docs/specs/demo" || record.Outcome != "Original intent on the first line and continued here." || !reflect.DeepEqual(record.CompletedTasks, completed) || !reflect.DeepEqual(record.PendingTasks, pending) || !reflect.DeepEqual(result.Completed, completed) || !reflect.DeepEqual(result.Pending, pending) {
		t.Fatalf("record=%+v result=%+v", record, result)
	}
	if record.QAReport != "" || record.QAVerdict != "" || !reflect.DeepEqual(record.Sources, []string{"entry.md"}) || !reflect.DeepEqual(record.ReturnedSources, []string{"docs/backlog/entry.md"}) {
		t.Fatalf("record=%+v", record)
	}
	if result.RemovedFiles != len(before) || result.RemovedBytes != size || len(result.RemovedPaths) != len(before) {
		t.Fatalf("removals=%+v want files=%d bytes=%d", result, len(before), size)
	}
	for _, path := range result.RemovedPaths {
		if _, ok := before[filepath.Join(req.RepositoryRoot, filepath.FromSlash(path))]; !ok {
			t.Fatalf("unknown removed path %q", path)
		}
	}
	if _, err := os.Lstat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Spec folder remains: %v", err)
	}
	returned, err := os.ReadFile(filepath.Join(req.RepositoryRoot, "docs", "backlog", "entry.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := bytes.Replace(adopted, []byte("status: promoted"), []byte("status: open"), 1)
	want = bytes.Replace(want, []byte("spec: demo"), []byte("spec: null"), 1)
	if !bytes.Equal(returned, want) {
		t.Fatalf("returned source changed other bytes:\n%s", returned)
	}
	if !reflect.DeepEqual(result.Returned, []ReturnedSource{{"docs/specs/demo/references/entry.md", "docs/backlog/entry.md", "open"}}) {
		t.Fatalf("returned=%+v", result.Returned)
	}
}

func TestCancelReturnsAnAdoptedBacklogEntryAsOpen(t *testing.T) {
	t.Parallel()
	for _, variant := range []string{"comments", "missing spec", "CRLF"} {
		t.Run(variant, func(t *testing.T) {
			req, dir := cancelFixture(t, 0)
			adopted := cancelBacklog
			if variant == "missing spec" {
				adopted = strings.Replace(adopted, "spec: demo  # adopted owner\n", "", 1)
			}
			if variant == "CRLF" {
				adopted = strings.ReplaceAll(adopted, "\n", "\r\n")
			}
			writeFile(t, filepath.Join(dir, "references", "entry.md"), adopted)
			want := strings.Replace(adopted, "status: promoted", "status: open", 1)
			if variant == "missing spec" {
				want = strings.Replace(want, "status: open # routing\n", "status: open # routing\nspec: null\n", 1)
			} else {
				want = strings.Replace(want, "spec: demo", "spec: null", 1)
			}
			result, err := Cancel(req)
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(filepath.Join(req.RepositoryRoot, result.Returned[0].To))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != want {
				t.Fatalf("source bytes:\ngot %q\nwant %q", got, want)
			}
		})
	}
}

func TestCancelReturnsAnAdoptedFindingAsPending(t *testing.T) {
	t.Parallel()
	for _, variant := range []string{"comments", "missing updated_at"} {
		t.Run(variant, func(t *testing.T) {
			req, dir := cancelFixture(t, 0)
			adopted := cancelFinding
			if variant == "missing updated_at" {
				adopted = strings.Replace(adopted, "updated_at: 2026-10-11 # latest evidence\n", "", 1)
			}
			cancelAdopt(t, dir, "finding", "finding.md", adopted)
			result, err := Cancel(req)
			if err != nil {
				t.Fatal(err)
			}
			record := cancelReadRecord(t, result)
			if !reflect.DeepEqual(record.ReturnedSources, []string{"docs/backlog/entry.md", "docs/findings/finding.md"}) {
				t.Fatalf("source order=%v", record.ReturnedSources)
			}
			got, err := os.ReadFile(filepath.Join(req.RepositoryRoot, "docs", "findings", "finding.md"))
			if err != nil {
				t.Fatal(err)
			}
			want := strings.Replace(adopted, "status: done", "status: pending", 1)
			want = strings.Replace(want, "absorbed_by: demo\n", "", 1)
			if variant == "missing updated_at" {
				want = strings.Replace(want, "---\n\n# Investigation", "updated_at: 2026-10-12\n---\n\n# Investigation", 1)
			} else {
				want = strings.Replace(want, "updated_at: 2026-10-11", "updated_at: 2026-10-12", 1)
			}
			want += "\n## Addendum — 2026-10-12 — Spec demo cancelled\n\nSpec demo was cancelled on 2026-10-12 before it implemented this finding: no consumer\n"
			if string(got) != want {
				t.Fatalf("finding bytes:\ngot %q\nwant %q", got, want)
			}
			if result.Returned[1].Status != "pending" {
				t.Fatalf("return=%+v", result.Returned[1])
			}
		})
	}
}

func TestCancelLeavesOtherAdoptedSourcesInTheFolder(t *testing.T) {
	t.Parallel()
	req, dir := cancelFixture(t, 0)
	cancelAdopt(t, dir, "inbox", "note.md", "Original raw evidence.\n")
	result, err := Cancel(req)
	if err != nil {
		t.Fatal(err)
	}
	record := cancelReadRecord(t, result)
	if !reflect.DeepEqual(record.Sources, []string{"entry.md", "note.md"}) || len(result.Returned) != 1 {
		t.Fatalf("record=%+v result=%+v", record, result)
	}
	for _, family := range []string{"backlog", "findings", "references"} {
		if _, err := os.Lstat(filepath.Join(req.RepositoryRoot, "docs", family, "note.md")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("other source extracted: %v", err)
		}
	}
}

func TestCancelWithoutATaskGraph(t *testing.T) {
	t.Parallel()
	for _, builtIn := range []bool{false, true} {
		t.Run(fmt.Sprint(builtIn), func(t *testing.T) {
			req, _ := cancelFixture(t, 0)
			req.BuiltInRoot = builtIn
			result, err := Cancel(req)
			if err != nil {
				t.Fatal(err)
			}
			record := cancelReadRecord(t, result)
			if len(record.CompletedTasks) != 0 || len(record.PendingTasks) != 0 || len(result.Completed) != 0 || len(result.Pending) != 0 {
				t.Fatalf("unexpected Tasks: %+v", record)
			}
			if result.RecordPath != ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, builtIn), req.Slug) {
				t.Fatal(result.RecordPath)
			}
		})
	}
}

func TestCancelRefusalsWriteNothing(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"missing folder", "missing PRD", "supersession", "blank reason", "revision length", "revision hex", "manifest", "record occupied", "folder occupied", "source missing", "source symlink", "source parent symlink", "status missing", "destination occupied", "duplicate destination", "unsafe path", "other source missing"} {
		t.Run(kind, func(t *testing.T) {
			req, dir := cancelFixture(t, 8)
			switch kind {
			case "missing folder":
				if err := os.RemoveAll(dir); err != nil {
					t.Fatal(err)
				}
			case "missing PRD":
				if err := os.Remove(filepath.Join(dir, "_prd.md")); err != nil {
					t.Fatal(err)
				}
			case "supersession":
				writeFile(t, filepath.Join(dir, SupersessionFilename), "anything")
			case "blank reason":
				req.Reason = " \n\t "
			case "revision length":
				req.SourceRevision = strings.Repeat("a", 39)
			case "revision hex":
				req.SourceRevision = strings.Repeat("z", 40)
			case "manifest":
				writeFile(t, filepath.Join(dir, "_tasks.md"), "invalid manifest")
			case "record occupied":
				writeFile(t, ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug), "occupied")
			case "folder occupied":
				writeFile(t, filepath.Join(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug, "keep.md"), "occupied")
			case "source missing", "source symlink":
				path := filepath.Join(dir, "references", "entry.md")
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if kind == "source symlink" {
					if err := os.Symlink("../_prd.md", path); err != nil {
						t.Fatal(err)
					}
				}
			case "source parent symlink":
				path := filepath.Join(dir, "references", "nested")
				if err := os.Symlink(".", path); err != nil {
					t.Fatal(err)
				}
				writeFile(t, filepath.Join(dir, "references", "_index.md"), "| type | path |\n| --- | --- |\n| backlog | nested/entry.md |\n")
			case "status missing":
				writeFile(t, filepath.Join(dir, "references", "entry.md"), "---\nspec: demo\n---\n\nOriginal\n")
			case "destination occupied":
				writeFile(t, filepath.Join(req.RepositoryRoot, "docs", "backlog", "entry.md"), "occupied")
			case "duplicate destination":
				index := filepath.Join(dir, "references", "_index.md")
				content, err := os.ReadFile(index)
				if err != nil {
					t.Fatal(err)
				}
				writeFile(t, index, string(content)+"| another | backlog | demo | 2026-10-10 | entry.md |\n")
			case "unsafe path":
				writeFile(t, filepath.Join(dir, "references", "_index.md"), "| type | path |\n| --- | --- |\n| backlog | ../_prd.md |\n")
			case "other source missing":
				cancelAdopt(t, dir, "inbox", "missing.md", "raw")
				if err := os.Remove(filepath.Join(dir, "references", "missing.md")); err != nil {
					t.Fatal(err)
				}
			}
			before := cancelTree(t, req.RepositoryRoot)
			wantError := map[string]string{
				"missing folder": "read cancellation Spec folder", "missing PRD": "read cancellation PRD",
				"supersession": "supersession", "blank reason": "non-blank reason",
				"revision length": "40-hex", "revision hex": "40-hex", "manifest": "Task Graph",
				"record occupied": "already exists", "folder occupied": "already exists",
				"source missing": "indexed cancellation path", "source symlink": "symbolic link",
				"source parent symlink": "symbolic link", "status missing": "missing status key",
				"destination occupied": "already exists", "duplicate destination": "duplicate return destination",
				"unsafe path": "unsafe", "other source missing": "indexed cancellation path",
			}[kind]
			if _, err := Cancel(req); err == nil || !strings.Contains(err.Error(), wantError) {
				t.Fatalf("expected refusal %q: %v", wantError, err)
			}
			if got := cancelTree(t, req.RepositoryRoot); !reflect.DeepEqual(got, before) {
				t.Fatalf("refusal changed tree:\nbefore=%v\nafter=%v", before, got)
			}
		})
	}
}

func TestCancelRollsBackAFailedWrite(t *testing.T) {
	t.Parallel()
	req, dir := cancelFixture(t, 8)
	cancelAdopt(t, dir, "finding", "finding.md", cancelFinding)
	destination := filepath.Join(req.RepositoryRoot, "docs", "findings")
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(destination, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(destination, 0o755); err != nil {
			t.Error(err)
		}
	})
	before := cancelTree(t, req.RepositoryRoot)
	_, err := Cancel(req)
	if err == nil || !errors.Is(err, os.ErrPermission) || !strings.Contains(err.Error(), "finding.md") {
		t.Fatalf("expected later filesystem permission failure: %v", err)
	}
	if got := cancelTree(t, req.RepositoryRoot); !reflect.DeepEqual(got, before) {
		t.Fatalf("rollback changed tree:\nbefore=%v\nafter=%v", before, got)
	}
}

func TestCancelRelinksLiveCitations(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"success", "dirty live document", "occupied promotion", "duplicate promotion", "symlink promotion", "rewrite failure"} {
		t.Run(kind, func(t *testing.T) {
			req, dir := cancelFixture(t, 0)
			root := req.RepositoryRoot
			source := "docs/specs/demo/references/shared.md"
			live := "docs/specs/secondary/_prd.md"
			original := "# Secondary\n\n[source](../demo/references/shared.md#part)\n[owner](../demo/_prd.md)\n"
			writeFile(t, filepath.Join(root, source), "Shared reference bytes.\n")
			writeFile(t, filepath.Join(root, live), original)
			second := "docs/specs/third/_prd.md"
			if kind == "rewrite failure" {
				writeFile(t, filepath.Join(root, second), "[source](../demo/references/shared.md)\n")
			}
			if kind == "duplicate promotion" {
				writeFile(t, filepath.Join(dir, "other", "shared.md"), "Other source\n")
				writeFile(t, filepath.Join(root, live), original+"[other](../demo/other/shared.md)\n")
			}
			if kind == "symlink promotion" {
				if err := os.Remove(filepath.Join(root, source)); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("entry.md", filepath.Join(root, source)); err != nil {
					t.Fatal(err)
				}
			}
			gittest.InitRepo(t, root)
			gittest.Run(t, root, "add", "-A")
			gittest.Run(t, root, "commit", "-m", "docs: cancellation fixture")
			if kind == "dirty live document" {
				writeFile(t, filepath.Join(root, live), original+"uncommitted\n")
			}
			if kind == "occupied promotion" {
				writeFile(t, filepath.Join(root, "docs", "references", "shared.md"), "occupied")
			}
			if kind == "rewrite failure" {
				path := filepath.Join(root, second)
				if err := os.Chmod(path, 0o444); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chmod(path, 0o644); err != nil {
						t.Error(err)
					}
				})
			}
			before := cancelTree(t, root)
			result, err := Cancel(req)
			if kind != "success" {
				if err == nil {
					t.Fatal("accepted relink refusal/failure")
				}
				if kind == "rewrite failure" && (!errors.Is(err, os.ErrPermission) || !strings.Contains(err.Error(), "write cancellation relink")) {
					t.Fatalf("expected rewrite failure: %v", err)
				}
				if got := cancelTree(t, root); !reflect.DeepEqual(got, before) {
					t.Fatalf("relink failure changed tree:\nbefore=%v\nafter=%v", before, got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			record := cancelReadRecord(t, result)
			if !reflect.DeepEqual(record.Relinked, []string{live}) || !reflect.DeepEqual(record.Promoted, []string{"docs/references/shared.md"}) {
				t.Fatalf("record=%+v", record)
			}
			want := "# Secondary\n\n[source](../../references/shared.md#part)\n[owner](../../history/specs/demo.md)\n"
			got, err := os.ReadFile(filepath.Join(root, live))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != want {
				t.Fatalf("live bytes=%q want=%q", got, want)
			}
			for _, link := range []string{"../../references/shared.md", "../../history/specs/demo.md"} {
				if _, err := os.Stat(filepath.Join(filepath.Dir(filepath.Join(root, live)), link)); err != nil {
					t.Fatalf("rewritten target %s: %v", link, err)
				}
			}
			promoted, err := os.ReadFile(filepath.Join(root, "docs", "references", "shared.md"))
			if err != nil || string(promoted) != "Shared reference bytes.\n" {
				t.Fatalf("promoted bytes=%q err=%v", promoted, err)
			}
		})
	}
}

func TestCancelKeepsObservedQAAndTaskGraphOrder(t *testing.T) {
	t.Parallel()
	req, dir := cancelFixture(t, 3)
	const nodes = `    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_01
      file: task_01.md
      needs: []
`
	writeFile(t, filepath.Join(dir, "_tasks.md"), manifestFixtureWithQA("spec-tasks/v1", "qa: declined\nqa_reason: fixture\n", nodes, ""))
	writeFile(t, filepath.Join(dir, "task_01.md"), taskFixture("task_01", "Build", "failed", "backend", defaultVerificationSection))
	writeFile(t, filepath.Join(dir, "task_02.md"), taskFixture("task_02", "Build", "in_progress", "backend", defaultVerificationSection))
	writeFile(t, filepath.Join(dir, "qa", "qa-report-2026-10-11.md"), "---\nverdict: fail\n---\n\n# QA\n\n## Outcome\n\nQA outcome must not replace the PRD intent.\n")
	result, err := Cancel(req)
	if err != nil {
		t.Fatal(err)
	}
	record := cancelReadRecord(t, result)
	if !reflect.DeepEqual(result.Completed, []string{"task_03"}) || !reflect.DeepEqual(result.Pending, []string{"task_01", "task_02"}) || record.QAReport != "qa-report-2026-10-11.md" || record.QAVerdict != VerdictFail || record.Outcome != "Original intent on the first line and continued here." {
		t.Fatalf("record=%+v result=%+v", record, result)
	}
}
