package cli

// Boundary: command dispatch, real Git and links in disposable repositories and homes.
// Invariant: an escaping path is refused before any repository or outside file changes.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"docs/history", "Spec folder", "promotion", "offline plan"} {
		t.Run(name, func(t *testing.T) {
			f := newCheckoutSpellingFixture(t, false)
			repo := f.main
			outside := filepath.Join(filepath.Dir(repo), "outside")
			mustMkdir(t, outside)
			mustWrite(t, filepath.Join(outside, "sentinel.md"), "outside bytes\n")
			folder := filepath.Join(repo, "docs/specs", spellingSlug)
			link := filepath.Join(repo, "docs/history")
			target := outside
			args := []string{"archive", spellingSlug, "--qa-override", "--approval", "probe", "--reason", "probe"}
			reason := fmt.Sprintf("resolve Archive Root for archive relink: path %q is outside the repository %q", filepath.Join(link, "specs"), repo)
			switch name {
			case "Spec folder":
				target = filepath.Join(outside, spellingSlug)
				if err := os.Rename(folder, target); err != nil {
					t.Fatal(err)
				}
				link = folder
				reason = fmt.Sprintf("path %q is outside the repository %q", folder, repo)
			case "promotion":
				historyWrite(t, filepath.Join(folder, "references/notes.md"), "# Notes\n\nPreserved knowledge.\n")
				link = filepath.Join(repo, "docs/references")
				args = append(args, "--promote", "references/notes.md")
				reason = fmt.Sprintf("resolve promotion destination: path %q is outside the repository %q", filepath.Join(link, "notes.md"), repo)
			case "offline plan":
				args = []string{"archive", spellingSlug, "--plan", "--offline"}
				reason = fmt.Sprintf("path %q is outside the repository %q", filepath.Join(link, "specs"), repo)
			}
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}
			historyCommit(t, repo, "seed escaping archive link")
			assertEscapingCommandRefusal(t, f, outside, archivePullRequestPreflight(reason), args...)
		})
	}
}

func TestCancelRefusesALinkThatLeavesTheRepositoryTranscript(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"docs/backlog", "docs/history"} {
		t.Run(name, func(t *testing.T) {
			f := newCheckoutSpellingFixture(t, true)
			repo := f.main
			outside := filepath.Join(filepath.Dir(repo), "outside")
			mustMkdir(t, outside)
			mustWrite(t, filepath.Join(outside, "sentinel.md"), "outside bytes\n")
			link := filepath.Join(repo, filepath.FromSlash(name))
			reason := fmt.Sprintf("path %q is outside the repository %q", filepath.Join(link, "specs"), repo)
			if name == "docs/backlog" {
				// The returned source has no existing destination or tracked Live Document.
				if err := os.RemoveAll(link); err != nil {
					t.Fatal(err)
				}
				folder := filepath.Join(repo, "docs/specs", spellingSlug, "references")
				historyWrite(t, filepath.Join(folder, "_index.md"), "| path | type |\n| --- | --- |\n| 2026-10-20-entry.md | backlog |\n")
				historyWrite(t, filepath.Join(folder, "2026-10-20-entry.md"), "---\ntype: fix\nstatus: promoted\nspec: 0001-demo\n---\n\n# Entry\n")
				reason = fmt.Sprintf("cancellation path %q is outside the repository", filepath.Join(link, "2026-10-20-entry.md"))
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Fatal(err)
			}
			historyCommit(t, repo, "seed escaping cancel link")
			assertEscapingCommandRefusal(t, f, outside, cancelRefusal(reason), "cancel", spellingSlug, "--reason", "probe")
		})
	}
}

func assertEscapingCommandRefusal(t *testing.T, f checkoutSpellingFixture, outside, wantErr string, args ...string) {
	t.Helper()
	folder, err := filepath.EvalSymlinks(filepath.Join(f.main, "docs/specs", spellingSlug))
	if err != nil {
		t.Fatal(err)
	}
	beforeSpec := snapshotDirectoryFiles(t, folder)
	beforeOutside := snapshotDirectoryFiles(t, outside)
	status := gittest.Run(t, f.main, "status", "--porcelain", "--untracked-files=all")
	got := runSpellingCommand(t, f, f.main, args...)
	if got.code != 2 || got.out != "" || got.errout != wantErr {
		t.Fatalf("transcript=%+v; want exit=2 stdout=empty stderr=%q", got, wantErr)
	}
	if !reflect.DeepEqual(beforeSpec, snapshotDirectoryFiles(t, folder)) {
		t.Fatal("refusal changed Spec bytes")
	}
	if !reflect.DeepEqual(beforeOutside, snapshotDirectoryFiles(t, outside)) {
		t.Fatal("refusal changed outside files")
	}
	if got := gittest.Run(t, f.main, "status", "--porcelain", "--untracked-files=all"); got != status {
		t.Fatalf("refusal changed Git status: %q -> %q", status, got)
	}
	assertPathExists(t, filepath.Join(f.main, "docs/specs", spellingSlug))
	assertPathMissing(t, filepath.Join(f.main, spellingRecord))
}

func TestArchiveKeepsLinksThatStayInsideTranscript(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"docs/history", "linked checkout", "external Spec Root"} {
		t.Run(name, func(t *testing.T) {
			f := newCheckoutSpellingFixture(t, false)
			args := []string{"archive", spellingSlug, "--qa-override", "--approval", "probe", "--reason", "probe"}
			controlRepo := f.worktree(t, "control")
			repo := f.worktree(t, "linked")
			controlDir, dir := controlRepo, repo
			record := spellingRecord
			switch name {
			case "docs/history":
				target := filepath.Join(repo, "docs/hist-real")
				mustMkdir(t, target)
				if err := os.Symlink(target, filepath.Join(repo, "docs/history")); err != nil {
					t.Fatal(err)
				}
			case "linked checkout":
				dir = filepath.Join(f.link, "linked")
			case "external Spec Root":
				// Both project repositories use the same external-root layout; only its
				// spelling differs. Their external worktrees share the committed Spec.
				controlDir = linkedRootProject(t, f, "control-project", filepath.Join(controlRepo, "docs/specs"))
				dir = linkedRootProject(t, f, "linked-project", filepath.Join(f.link, "linked", "docs/specs"))
				record = "docs/specs/_archived/0001-demo.md"
			}
			control := runSpellingCommand(t, f, controlDir, args...)
			got := runSpellingCommand(t, f, dir, args...)
			if control.code != 0 || control.errout != "" || control.out == "" {
				t.Fatalf("control=%+v", control)
			}
			if name == "external Spec Root" {
				// External destinations are absolute in stdout. Assert that
				// exact spelling before substituting only the fixture path.
				controlPath := filepath.ToSlash(filepath.Join(controlRepo, record))
				linkedPath := filepath.ToSlash(filepath.Join(f.link, "linked", record))
				if !strings.Contains(control.out, " -> "+controlPath+"; removed ") ||
					!strings.Contains(got.out, " -> "+linkedPath+"; removed ") {
					t.Fatalf("external record destinations: got=%+v control=%+v", got, control)
				}
				got.out = strings.Replace(got.out, " -> "+linkedPath+"; removed ", " -> "+controlPath+"; removed ", 1)
			}
			if got != control {
				t.Fatalf("transcript=%+v control=%+v", got, control)
			}
			if got, want := mustRead(t, filepath.Join(repo, record)), mustRead(t, filepath.Join(controlRepo, record)); got != want {
				t.Fatalf("record bytes=%q control=%q", got, want)
			}
			if got, want := mustRead(t, filepath.Join(repo, spellingBacklog)), mustRead(t, filepath.Join(controlRepo, spellingBacklog)); got != want {
				t.Fatalf("relinked bytes=%q control=%q", got, want)
			}
			assertPathMissing(t, filepath.Join(repo, "docs/specs", spellingSlug))
		})
	}
}

func linkedRootProject(t *testing.T, f checkoutSpellingFixture, name, root string) string {
	t.Helper()
	project := filepath.Join(filepath.Dir(f.main), name)
	mustMkdir(t, project)
	gittest.InitRepo(t, project, "--initial-branch=main")
	mustWrite(t, filepath.Join(project, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", root))
	historyCommit(t, project, "configure external Spec Root")
	return project
}
