package cli

// Suite: Archive Relink command transcripts and refusal atomicity.
// Invariant: the archive preserves committed Live Document bytes outside citations.
// Boundary IN: CLI dispatch in committed temporary repositories and homes.
// Boundary OUT: judge transport rejects every request.

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

const archiveRelinkOwner = "0001-owner"
const archiveRelinkSource = "docs/specs/0001-owner/references/2026-01-01-source.md"
const archiveRelinkPromotion = "docs/references/2026-01-01-source.md"
const archiveRelinkRecord = "docs/history/specs/0001-owner.md"
const archiveRelinkHandoff = "docs/handoffs/2026-01-03-notes.md"

func archiveRelinkTree(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = mustRead(t, path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func archiveRelinkFixture(t *testing.T, inbound bool) (string, commandEnvironment, *specJudgeTransport) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	old := filepath.Join(repo, "docs/specs", implementTestSlug)
	owner := filepath.Join(repo, "docs/specs", archiveRelinkOwner)
	if err := os.Rename(old, owner); err != nil {
		t.Fatal(err)
	}
	for path, content := range archiveRelinkTree(t, owner) {
		mustWrite(t, filepath.Join(owner, path), strings.ReplaceAll(content, implementTestSlug, archiveRelinkOwner))
	}
	mustMkdir(t, filepath.Dir(filepath.Join(repo, archiveRelinkSource)))
	mustWrite(t, filepath.Join(repo, archiveRelinkSource), "# Source\n\nPreserve these bytes, including CRLF.\r\n[PRD](../_prd.md)\n")
	if inbound {
		writeImplementSpec(t, repo, "0002-secondary", []implementSeed{{id: "task_01", status: string(spec.StatusPending)}})
		prd := filepath.Join(repo, "docs/specs/0002-secondary/_prd.md")
		mustWrite(t, prd, mustRead(t, prd)+"\n[Source](../0001-owner/references/2026-01-01-source.md#lesson)\n")
		task := filepath.Join(repo, "docs/specs/0002-secondary/task_01.md")
		mustWrite(t, task, mustRead(t, task)+"\nRead `"+archiveRelinkSource+"`.\n")
		backlog := filepath.Join(repo, "docs/backlog/2026-01-02-follow-up.md")
		mustMkdir(t, filepath.Dir(backlog))
		mustWrite(t, backlog, "---\ntype: feat\nstatus: open\n---\n\n# Follow-up\n\n[Owner](../specs/0001-owner/_prd.md)\n")
		mustMkdir(t, filepath.Dir(filepath.Join(repo, archiveRelinkHandoff)))
		mustWrite(t, filepath.Join(repo, archiveRelinkHandoff), "# Notes\n\n[Source](../specs/0001-owner/references/2026-01-01-source.md)\n")
	}
	commitArchiveFixture(t)
	fake := &specJudgeTransport{t: t}
	env := commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
	env.dependencies.judgeTransport = fake
	return repo, env, fake
}

func TestArchiveRelinksLiveCitationsTranscript(t *testing.T) {
	t.Parallel()
	repo, env, fake := archiveRelinkFixture(t, true)
	before := archiveRelinkTree(t, repo)
	revision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	files, size := 0, 0
	for path, content := range before {
		if strings.HasPrefix(path, "docs/specs/0001-owner/") {
			files++
			size += len(content)
		}
	}
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner}, &out, &errout, env)
	wantOut := fmt.Sprintf("archived 0001-owner -> %s; removed %d file(s) (%d bytes) kept in Git at %.12s; promoted 1 file(s) to docs/references/; relinked 3 file(s)\n", archiveRelinkRecord, files, size, revision)
	wantErr := "roundfix: warning: archive leaves " + archiveRelinkHandoff + ":3 naming removed " + archiveRelinkSource + " unrelinked\n"
	if code != exitOK || out.String() != wantOut || errout.String() != wantErr {
		t.Fatalf("exit=%d stdout=%q want=%q stderr=%q want=%q", code, &out, wantOut, &errout, wantErr)
	}
	if got := mustRead(t, filepath.Join(repo, archiveRelinkPromotion)); got != strings.ReplaceAll(before[archiveRelinkSource], "[PRD](../_prd.md)", "PRD (`docs/specs/0001-owner/_prd.md@"+revision+"`)") {
		t.Fatalf("promoted bytes differ from the rewritten source: %q", got)
	}
	assertPathMissing(t, filepath.Join(repo, "docs/specs/0001-owner"))
	changes := map[string][2]string{
		"docs/specs/0002-secondary/_prd.md":    {"../0001-owner/references/2026-01-01-source.md#lesson", "../../references/2026-01-01-source.md#lesson"},
		"docs/specs/0002-secondary/task_01.md": {archiveRelinkSource, archiveRelinkPromotion},
		"docs/backlog/2026-01-02-follow-up.md": {"../specs/0001-owner/_prd.md", "../history/specs/0001-owner.md"},
	}
	for path, replacement := range changes {
		committed := gittest.Run(t, repo, "show", "HEAD:"+path)
		want := strings.ReplaceAll(committed, replacement[0], replacement[1])
		if got := mustRead(t, filepath.Join(repo, path)); got != want {
			t.Fatalf("relinked bytes %s=%q want=%q", path, got, want)
		}
	}
	if got := mustRead(t, filepath.Join(repo, archiveRelinkHandoff)); got != before[archiveRelinkHandoff] {
		t.Fatal("handoff changed")
	}
	record := archiveRecordAt(t, filepath.Join(repo, archiveRelinkRecord))
	if !reflect.DeepEqual(record.Promoted, []string{archiveRelinkPromotion}) || !reflect.DeepEqual(record.Relinked, []string{"docs/backlog/2026-01-02-follow-up.md", "docs/specs/0002-secondary/_prd.md", "docs/specs/0002-secondary/task_01.md"}) {
		t.Fatalf("record=%+v", record)
	}
	if fake.calls != 0 {
		t.Fatalf("provider requests=%d", fake.calls)
	}
}

func TestArchiveRefusesACollidingPromotion(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"existing", "duplicate basename"} {
		t.Run(kind, func(t *testing.T) {
			repo, env, fake := archiveRelinkFixture(t, true)
			want := "promotion destination"
			if kind == "existing" {
				mustMkdir(t, filepath.Dir(filepath.Join(repo, archiveRelinkPromotion)))
				mustWrite(t, filepath.Join(repo, archiveRelinkPromotion), "existing\n")
			} else {
				path := "docs/specs/0001-owner/other/2026-01-01-source.md"
				mustMkdir(t, filepath.Dir(filepath.Join(repo, path)))
				mustWrite(t, filepath.Join(repo, path), "other\n")
				prd := filepath.Join(repo, "docs/specs/0002-secondary/_prd.md")
				mustWrite(t, prd, mustRead(t, prd)+"\n[other](../0001-owner/other/2026-01-01-source.md)\n")
				want = "duplicate promotion destination"
			}
			commitArchiveFixture(t)
			before := archiveRelinkTree(t, repo)
			var out, errout bytes.Buffer
			code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner}, &out, &errout, env)
			if code != exitPreflight || out.Len() != 0 || !strings.Contains(errout.String(), want) {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &errout)
			}
			if !reflect.DeepEqual(before, archiveRelinkTree(t, repo)) {
				t.Fatal("refusal changed files")
			}
			assertPathExists(t, filepath.Join(repo, "docs/specs/0001-owner"))
			if fake.calls != 0 {
				t.Fatalf("provider requests=%d", fake.calls)
			}
		})
	}
}

func TestArchiveRefusesADirtyLiveDocument(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"unstaged", "staged", "deleted", "staged deletion", "removed citation", "staged new citation"} {
		t.Run(kind, func(t *testing.T) {
			repo, env, fake := archiveRelinkFixture(t, true)
			path := "docs/specs/0002-secondary/_prd.md"
			if kind == "staged new citation" {
				path = "docs/references/new.md"
				mustMkdir(t, filepath.Dir(filepath.Join(repo, path)))
				mustWrite(t, filepath.Join(repo, path), archiveRelinkSource+"\n")
				gittest.Run(t, repo, "add", "--", path)
			} else if kind == "staged deletion" {
				gittest.Run(t, repo, "rm", "--", path)
			} else if kind == "deleted" {
				if err := os.Remove(filepath.Join(repo, path)); err != nil {
					t.Fatal(err)
				}
			} else {
				content := mustRead(t, filepath.Join(repo, path)) + "\nUncommitted edit.\n"
				if kind == "removed citation" {
					content = strings.ReplaceAll(content, "[Source](../0001-owner/references/2026-01-01-source.md#lesson)", "Source removed.")
				}
				mustWrite(t, filepath.Join(repo, path), content)
				if kind == "staged" {
					gittest.Run(t, repo, "add", "--", path)
				}
			}
			before := archiveRelinkTree(t, repo)
			var out, errout bytes.Buffer
			code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner}, &out, &errout, env)
			want := fmt.Sprintf("Spec %q cannot archive: %s has changes not committed at HEAD; commit them so the archive can relink it", archiveRelinkOwner, path)
			if code != exitPreflight || out.Len() != 0 || !strings.Contains(errout.String(), want) {
				t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, &out, &errout, want)
			}
			if !reflect.DeepEqual(before, archiveRelinkTree(t, repo)) {
				t.Fatal("dirty refusal changed files")
			}
			if fake.calls != 0 {
				t.Fatalf("provider requests=%d", fake.calls)
			}
		})
	}
}

func TestArchiveWithoutInboundCitationsIsUnchanged(t *testing.T) {
	t.Parallel()
	repo, env, fake := archiveRelinkFixture(t, false)
	revision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	old, err := spec.BuildArchiveRecord(spec.ArchiveRecordInput{SpecDir: filepath.Join(repo, "docs/specs", archiveRelinkOwner), Slug: archiveRelinkOwner, Source: "docs/specs/0001-owner", SourceRevision: revision, Archived: time.Now().Format("2006-01-02")})
	if err != nil {
		t.Fatal(err)
	}
	wantRecord, err := spec.RenderArchiveRecord(old)
	if err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner}, &out, &errout, env)
	path := filepath.Join(repo, archiveRelinkRecord)
	if code != exitOK || out.String() != archiveConfirmation(t, repo, path, false) || errout.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &errout)
	}
	if got := mustRead(t, path); got != string(wantRecord) {
		t.Fatalf("record bytes differ\n%s\nwant\n%s", got, wantRecord)
	}
	if fake.calls != 0 {
		t.Fatalf("provider requests=%d", fake.calls)
	}
}

func TestArchivePlanListsRelinks(t *testing.T) {
	t.Parallel()
	repo, env, fake := archiveRelinkFixture(t, true)
	before := archiveRelinkTree(t, repo)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner, "--plan"}, &out, &errout, env)
	want := "promotes references/2026-01-01-source.md to " + archiveRelinkPromotion + "\n" +
		"relinks docs/backlog/2026-01-02-follow-up.md: 1 citation(s)\n" +
		"relinks docs/specs/0002-secondary/_prd.md: 1 citation(s)\n" +
		"relinks docs/specs/0002-secondary/task_01.md: 1 citation(s)\n" +
		"leaves " + archiveRelinkHandoff + ":3 naming " + archiveRelinkSource + " unrelinked\n" +
		"promote with: roundfix archive 0001-owner --promote <path>\n"
	if code != exitOK || errout.Len() != 0 || !strings.HasSuffix(out.String(), want) {
		t.Fatalf("exit=%d stdout=%q want suffix=%q stderr=%q", code, &out, want, &errout)
	}
	if !reflect.DeepEqual(before, archiveRelinkTree(t, repo)) {
		t.Fatal("plan changed files")
	}
	if fake.calls != 0 {
		t.Fatalf("provider requests=%d", fake.calls)
	}
}

func TestArchiveRelinkCombinesManualPromotionsAndPreservesHistory(t *testing.T) {
	t.Parallel()
	repo, env, fake := archiveRelinkFixture(t, true)
	manual := "docs/specs/0001-owner/references/aaa-manual.md"
	mustWrite(t, filepath.Join(repo, manual), "manual bytes\n")
	history := "docs/history/specs/old.md"
	mustMkdir(t, filepath.Dir(filepath.Join(repo, history)))
	mustWrite(t, filepath.Join(repo, history), archiveRelinkSource+"\n")
	commitArchiveFixture(t)
	before := archiveRelinkTree(t, repo)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner, "--promote", "references/aaa-manual.md", "--promote", "references/2026-01-01-source.md"}, &out, &errout, env)
	if code != exitOK {
		t.Fatalf("exit=%d stderr=%q", code, &errout)
	}
	r := archiveRecordAt(t, filepath.Join(repo, archiveRelinkRecord))
	if !reflect.DeepEqual(r.Promoted, []string{archiveRelinkPromotion, "docs/references/aaa-manual.md"}) {
		t.Fatalf("promoted=%v", r.Promoted)
	}
	if mustRead(t, filepath.Join(repo, "docs/references/aaa-manual.md")) != before[manual] || mustRead(t, filepath.Join(repo, history)) != before[history] {
		t.Fatal("manual bytes or History bytes changed")
	}
	if fake.calls != 0 {
		t.Fatalf("provider requests=%d", fake.calls)
	}
}

func TestArchiveRelinkAllowsAnUnrelatedStagedLiveDocument(t *testing.T) {
	t.Parallel()
	repo, env, _ := archiveRelinkFixture(t, true)
	path := "docs/references/unrelated-new.md"
	mustMkdir(t, filepath.Dir(filepath.Join(repo, path)))
	mustWrite(t, filepath.Join(repo, path), "# New unrelated document\n")
	gittest.Run(t, repo, "add", "--", path)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner}, &out, &errout, env)
	if code != exitOK || mustRead(t, filepath.Join(repo, path)) != "# New unrelated document\n" {
		t.Fatalf("exit=%d stderr=%q", code, &errout)
	}
}

func TestArchiveRelinkKeepsExistingPromotionRefusalWording(t *testing.T) {
	t.Parallel()
	repo, env, _ := archiveRelinkFixture(t, true)
	before := archiveRelinkTree(t, repo)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner, "--promote", "./_prd.md"}, &out, &errout, env)
	if code != exitPreflight || !strings.Contains(errout.String(), `cannot promote core artifact "./_prd.md"`) || !reflect.DeepEqual(before, archiveRelinkTree(t, repo)) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &errout)
	}
}
