package cli

// Boundary: command entry points in disposable linked worktrees; no host home or network.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

const spellingSlug = "0001-demo"
const spellingBacklog = "docs/backlog/2026-01-01-follow-up.md"
const spellingRecord = "docs/history/specs/0001-demo.md"

type checkoutSpellingFixture struct {
	main, physical, link, home string
	revision                   string
}

func newCheckoutSpellingFixture(t *testing.T, cancel bool) checkoutSpellingFixture {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := checkoutSpellingFixture{main: filepath.Join(root, "main"), physical: filepath.Join(root, "worktrees"), link: filepath.Join(root, "link"), home: filepath.Join(root, "home")}
	for _, dir := range []string{f.main, f.physical, f.home} {
		mustMkdir(t, dir)
	}
	if err := os.Symlink(f.physical, f.link); err != nil {
		t.Fatal(err)
	}
	gittest.InitRepo(t, f.main, "--initial-branch=main")
	seeds := []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}, {id: "task_02", taskType: string(spec.TaskTypeQA), status: string(spec.StatusPending), needs: []string{"task_01"}}}
	if cancel {
		seeds[1].taskType = "backend"
	}
	writeImplementSpec(t, f.main, spellingSlug, seeds)
	mustWrite(t, filepath.Join(f.main, "docs/specs/.keep"), "")
	prd := filepath.Join(f.main, "docs/specs", spellingSlug, "_prd.md")
	mustWrite(t, prd, strings.Replace(mustRead(t, prd), "# PRD\n\n", "# PRD\n\nSee [follow-up](../../backlog/2026-01-01-follow-up.md).\n\n", 1))
	mustMkdir(t, filepath.Dir(filepath.Join(f.main, spellingBacklog)))
	mustWrite(t, filepath.Join(f.main, spellingBacklog), "---\ntype: feat\nstatus: open\n---\n\n# Follow-up\n\n[Spec](../specs/0001-demo/_prd.md)\n")
	gittest.Run(t, f.main, "add", "-A")
	gittest.Run(t, f.main, "commit", "-m", "seed checkout spelling fixture")
	f.revision = strings.TrimSpace(gittest.Run(t, f.main, "rev-parse", "HEAD"))
	gittest.Run(t, f.main, "update-ref", "refs/remotes/origin/main", f.revision)
	gittest.Run(t, f.main, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	return f
}

func (f checkoutSpellingFixture) worktree(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(f.physical, name)
	gittest.Run(t, f.main, "worktree", "add", "-b", "test/"+name, path, f.revision)
	return path
}

func (f checkoutSpellingFixture) environment(dir string) commandEnvironment {
	return commandEnvironment{homeDir: f.home, workDir: dir, environ: []string{}, dependencies: defaultCommandDependencies()}
}

type spellingTranscript struct {
	code                         int
	out, errout, record, backlog string
}

func runSpellingCommand(t *testing.T, f checkoutSpellingFixture, dir string, args ...string) spellingTranscript {
	t.Helper()
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), args, &out, &errout, f.environment(dir))
	return spellingTranscript{code: code, out: out.String(), errout: errout.String()}
}

// Build the comparison outside a subtest so each spelling can also run alone.
func spellingControl(t *testing.T, f checkoutSpellingFixture, args ...string) spellingTranscript {
	t.Helper()
	repo := f.worktree(t, "control")
	got := runSpellingCommand(t, f, repo, args...)
	if got.code != 0 {
		t.Fatalf("control transcript=%+v", got)
	}
	got.record = mustRead(t, filepath.Join(repo, spellingRecord))
	got.backlog = mustRead(t, filepath.Join(repo, spellingBacklog))
	return got
}

func TestArchiveFromEverySpellingOfTheCheckoutTranscript(t *testing.T) {
	t.Parallel()
	f := newCheckoutSpellingFixture(t, false)
	control := spellingControl(t, f, "archive", spellingSlug, "--qa-override", "--approval", "probe", "--reason", "probe")
	for _, name := range []string{"worktree root resolved", "worktree root through link", "worktree subdirectory through link", "main checkout through link"} {
		t.Run(name, func(t *testing.T) {
			repo := f.main
			dir := ""
			if name == "main checkout through link" {
				dir = filepath.Join(filepath.Dir(f.main), "main-link")
				if err := os.Symlink(f.main, dir); err != nil {
					t.Fatal(err)
				}
			} else {
				item := strings.ReplaceAll(name, " ", "-")
				repo = f.worktree(t, item)
				dir = repo
				if name != "worktree root resolved" {
					dir = filepath.Join(f.link, item)
				}
				if name == "worktree subdirectory through link" {
					dir = filepath.Join(dir, "docs")
				}
			}
			files := snapshotDirectoryFiles(t, filepath.Join(repo, "docs/specs", spellingSlug))
			size := 0
			for _, content := range files {
				size += len(content)
			}
			got := runSpellingCommand(t, f, dir, "archive", spellingSlug, "--qa-override", "--approval", "probe", "--reason", "probe")
			want := fmt.Sprintf("archived 0001-demo with QA override -> docs/history/specs/0001-demo.md; removed %d file(s) (%d bytes) kept in Git at %.12s; relinked 1 file(s)\n", len(files), size, f.revision)
			if got.code != 0 || got.out != want || got.errout != "" {
				t.Fatalf("transcript=%+v want stdout=%q", got, want)
			}
			got.record = mustRead(t, filepath.Join(repo, spellingRecord))
			got.backlog = mustRead(t, filepath.Join(repo, spellingBacklog))
			if got != control {
				t.Fatalf("got=%+v control=%+v", got, control)
			}
			if !strings.Contains(got.record, "../../backlog/2026-01-01-follow-up.md") || !strings.Contains(got.backlog, "../history/specs/0001-demo.md") {
				t.Fatalf("record or relink destinations incorrect: %+v", got)
			}
			assertPathMissing(t, filepath.Join(repo, "docs/specs", spellingSlug))
		})
	}
}

func TestCancelFromALinkedSpellingTranscript(t *testing.T) {
	t.Parallel()
	f := newCheckoutSpellingFixture(t, true)
	control := spellingControl(t, f, "cancel", spellingSlug, "--reason", "probe")
	for _, linked := range []bool{false, true} {
		name := fmt.Sprintf("linked-%t", linked)
		t.Run(name, func(t *testing.T) {
			repo := f.worktree(t, name)
			dir := repo
			if linked {
				dir = filepath.Join(f.link, name)
			}
			files := snapshotDirectoryFiles(t, filepath.Join(repo, "docs/specs", spellingSlug))
			size := 0
			for _, content := range files {
				size += len(content)
			}
			got := runSpellingCommand(t, f, dir, "cancel", spellingSlug, "--reason", "probe")
			want := fmt.Sprintf("cancelled 0001-demo -> docs/history/specs/0001-demo.md; 1 completed Task(s), 1 pending; removed %d file(s) (%d bytes) kept in Git at %.12s\n", len(files), size, f.revision)
			if got.code != 0 || got.out != want || got.errout != "" {
				t.Fatalf("transcript=%+v want stdout=%q", got, want)
			}
			got.record = mustRead(t, filepath.Join(repo, spellingRecord))
			got.backlog = mustRead(t, filepath.Join(repo, spellingBacklog))
			if got != control {
				t.Fatalf("got=%+v control=%+v", got, control)
			}
			if !strings.Contains(got.backlog, "../history/specs/0001-demo.md") {
				t.Fatalf("backlog=%q", got.backlog)
			}
			assertPathMissing(t, filepath.Join(repo, "docs/specs", spellingSlug))
		})
	}
}

func TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling(t *testing.T) {
	t.Parallel()
	for _, recordOnly := range []bool{false, true} {
		t.Run(fmt.Sprintf("record-only-%t", recordOnly), func(t *testing.T) {
			f := newCheckoutSpellingFixture(t, false)
			args := []string{"archive", spellingSlug, "--qa-override", "--approval", "probe", "--reason", "probe"}
			reason := `Spec "0001-demo" has changes not committed at HEAD under docs/specs/0001-demo; commit them before archive so the Archive Record's source_revision holds the Spec`
			if recordOnly {
				got := runSpellingCommand(t, f, f.main, args...)
				if got.code != 0 {
					t.Fatalf("seed archive=%+v", got)
				}
				gittest.Run(t, f.main, "add", "-A")
				gittest.Run(t, f.main, "commit", "-m", "archive on default branch")
				f.revision = strings.TrimSpace(gittest.Run(t, f.main, "rev-parse", "HEAD"))
				gittest.Run(t, f.main, "update-ref", "refs/remotes/origin/main", f.revision)
				args = []string{"archive", spellingSlug, "--pull-request", "7"}
				reason = "Archive Record docs/history/specs/0001-demo.md is already on origin/main; an archived record stays byte-identical"
			}
			repo := f.worktree(t, "item")
			if !recordOnly {
				prd := filepath.Join(repo, "docs/specs", spellingSlug, "_prd.md")
				mustWrite(t, prd, mustRead(t, prd)+"\nUncommitted.\n")
			}
			before := archiveRelinkTree(t, repo)
			control := runSpellingCommand(t, f, repo, args...)
			linked := runSpellingCommand(t, f, filepath.Join(f.link, "item"), args...)
			if control != linked || linked.code != 2 || linked.out != "" || linked.errout != archivePullRequestPreflight(reason) {
				t.Fatalf("control=%+v linked=%+v want reason=%q", control, linked, reason)
			}
			after := archiveRelinkTree(t, repo)
			if !reflect.DeepEqual(before, after) {
				t.Fatal("refusal changed checkout")
			}
		})
	}
}

func TestArchivePlanWithALinkedExternalSpecRootListsTheRelink(t *testing.T) {
	t.Parallel()
	f := newCheckoutSpellingFixture(t, false)
	external := f.worktree(t, "external")
	project := filepath.Join(filepath.Dir(f.main), "project")
	mustMkdir(t, project)
	gittest.InitRepo(t, project, "--initial-branch=main")
	mustWrite(t, filepath.Join(project, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", filepath.Join(f.link, "external", "docs/specs")))
	gittest.Run(t, project, "add", "-A")
	gittest.Run(t, project, "commit", "-m", "configure external Spec Root")
	before := archiveRelinkTree(t, external)
	got := runSpellingCommand(t, f, project, "archive", spellingSlug, "--plan", "--offline")
	if got.code != 0 || got.errout != "" || !strings.Contains(got.out, "relinks "+spellingBacklog+": 1 citation(s)\n") {
		t.Fatalf("transcript=%+v", got)
	}
	if !reflect.DeepEqual(before, archiveRelinkTree(t, external)) {
		t.Fatal("plan changed external checkout")
	}
}
