package cli

// Suite: shared Spec folder refusal and command transcripts.
// Boundary: real Git status and CLI dispatch in disposable repositories.

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestArchiveRefusesIgnoredSpecFiles(t *testing.T) {
	t.Parallel()
	repo, env := ignoredSpecFixture(t)
	assertSpecFolderRefusal(t, repo, env, "archive", ignoredSpecRefusal("archive", "docs/specs/"+implementTestSlug))
}

func TestCancelRefusesIgnoredSpecFiles(t *testing.T) {
	t.Parallel()
	repo, env := ignoredSpecFixture(t)
	assertSpecFolderRefusal(t, repo, env, "cancel", ignoredSpecRefusal("cancel", "docs/specs/"+implementTestSlug))
}

func TestSpecFolderStatusRefusesUncommittedChangesFirst(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"archive", "cancel"} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()
			for _, change := range []string{"modified", "staged", "deleted", "untracked", "renamed"} {
				t.Run(change, func(t *testing.T) {
					t.Parallel()
					repo, env := ignoredSpecFixture(t)
					folder := filepath.Join(repo, "docs/specs", implementTestSlug)
					path := filepath.Join(folder, "note.md")
					switch change {
					case "modified", "staged":
						mustWrite(t, path, mustRead(t, path)+"\nUncommitted change.\n")
						if change == "staged" {
							gittest.Run(t, repo, "add", path)
						}
					case "deleted":
						gittest.Run(t, repo, "rm", path)
					case "untracked":
						mustWrite(t, filepath.Join(folder, "scratch.md"), "untracked bytes\n")
					case "renamed":
						gittest.Run(t, repo, "mv", path, filepath.Join(folder, "renamed.md"))
					}
					reason := fmt.Sprintf("Spec %q has changes not committed at HEAD under docs/specs/%s; commit them before %s so the Archive Record's source_revision holds the Spec", implementTestSlug, implementTestSlug, command)
					want := "Preflight failed\n\nReason:\n  " + reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix " + command + " --help' for usage.\n"
					assertSpecFolderRefusal(t, repo, env, command, want)
				})
			}
		})
	}
}

func TestSpecFolderStatusReadsTheExternalSpecsRootRepository(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"archive", "cancel"} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()
			repo, env := cancelFixture(t, 1, 0)
			externalRepo, externalRoot := newExternalSpecsRoot(t, implementTestSlug, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
			writeArchiveQAReport(t, externalRepo, spec.VerdictPass, "rows_passed: 1")
			seedIgnoredSpecFiles(t, externalRepo)
			configureExternalSpecsRoot(t, repo, externalRoot)
			before := historySnapshot(t, repo)
			status := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all", "--ignored")
			assertPathMissing(t, filepath.Join(repo, ".gitignore"))
			assertSpecFolderRefusal(t, externalRepo, env, command, ignoredSpecRefusal(command, "docs/specs/"+implementTestSlug))
			if !reflect.DeepEqual(historySnapshot(t, repo), before) {
				t.Fatal("external-root refusal changed the project repository")
			}
			if got := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all", "--ignored"); got != status {
				t.Fatalf("project Git status changed: %q -> %q", status, got)
			}
			assertPathMissing(t, filepath.Join(externalRoot, "_archived", implementTestSlug+".md"))
		})
	}
}

func ignoredSpecFixture(t *testing.T) (string, commandEnvironment) {
	t.Helper()
	repo, env := cancelFixture(t, 1, 0)
	writeArchiveQAReport(t, repo, spec.VerdictPass, "rows_passed: 1")
	mustWrite(t, filepath.Join(repo, "docs/specs", implementTestSlug, "note.md"), "committed note\n")
	seedIgnoredSpecFiles(t, repo)
	return repo, env
}

func seedIgnoredSpecFiles(t *testing.T, repo string) {
	t.Helper()
	// Ignore the directory itself too, so Git must enumerate its files.
	mustWrite(t, filepath.Join(repo, ".gitignore"), "*.log\nout/\n")
	historyCommit(t, repo, "docs: ignored Spec fixture")
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	for _, path := range []string{"run.log", "qa/evidence/out/b.log", "qa/evidence/out/a.log"} {
		mustMkdir(t, filepath.Dir(filepath.Join(folder, path)))
		mustWrite(t, filepath.Join(folder, path), "ignored evidence for "+path+"\n")
	}
}

func ignoredSpecRefusal(command, folder string) string {
	return fmt.Sprintf(`Preflight failed

Reason:
  Spec "%s" cannot %s while Git ignores 3 file(s) under %s that no commit holds: %s/qa/evidence/out/a.log, %s/qa/evidence/out/b.log, %s/run.log; remove each file or commit it with git add -f, then retry

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Usage:
  Run 'roundfix %s --help' for usage.
`, implementTestSlug, command, folder, folder, folder, folder, command)
}

func assertSpecFolderRefusal(t *testing.T, repo string, env commandEnvironment, command, wantErr string) {
	t.Helper()
	before := historySnapshot(t, repo)
	status := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all", "--ignored")
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	hash := cancelFolderHash(t, folder)
	args := []string{command, implementTestSlug}
	if command == "cancel" {
		args = append(args, "--reason", "no consumer")
	}
	var stdout, stderr bytes.Buffer
	code := runWithContext(t.Context(), args, &stdout, &stderr, env)
	cancelAssertTranscript(t, code, stdout.String(), stderr.String(), 2, "", wantErr)
	if !reflect.DeepEqual(historySnapshot(t, repo), before) {
		t.Fatal("refusal changed repository file bytes, including ignored files")
	}
	if got := cancelFolderHash(t, folder); got != hash {
		t.Fatal("refusal changed Spec folder hash")
	}
	if got := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all", "--ignored"); got != status {
		t.Fatalf("Git status changed: %q -> %q", status, got)
	}
	for _, path := range []string{"run.log", "qa/evidence/out/a.log", "qa/evidence/out/b.log"} {
		assertPathExists(t, filepath.Join(folder, path))
	}
	assertPathMissing(t, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"))
	assertNoRunDatabase(t, env.homeDir)
}

func TestSpecFolderStatusPreservesRawPaths(t *testing.T) {
	t.Parallel()
	repo, _ := ignoredSpecFixture(t)
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	tracked := filepath.Join(folder, " leading space.md")
	mustWrite(t, tracked, "committed\n")
	historyCommit(t, repo, "docs: tracked path with leading space")
	mustWrite(t, tracked, "modified\n")
	ignored := filepath.Join(folder, " space\tand\nnewline.log")
	mustWrite(t, ignored, "ignored\n")
	status, err := readSpecFolderStatus(t.Context(), repo, folder)
	if err != nil {
		t.Fatal(err)
	}
	prefix := "docs/specs/" + implementTestSlug + "/"
	if !reflect.DeepEqual(status.Changes, []string{prefix + " leading space.md"}) {
		t.Fatalf("changes = %#v; want the raw modified path", status.Changes)
	}
	wantIgnored := []string{prefix + " space\tand\nnewline.log", prefix + "qa/evidence/out/a.log", prefix + "qa/evidence/out/b.log", prefix + "run.log"}
	if !reflect.DeepEqual(status.Ignored, wantIgnored) {
		t.Fatalf("ignored = %#v, want %#v", status.Ignored, wantIgnored)
	}
}
