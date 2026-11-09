package worktree

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"roundfix/internal/store"
)

func TestAnUnregisteredRunWorktreeHoldingOnlyIgnoredFilesIsSafe(t *testing.T) {
	t.Parallel()
	for _, empty := range []bool{false, true} {
		t.Run(map[bool]string{false: "ignored files", true: "empty"}[empty], func(t *testing.T) {
			run := newUnregisteredIgnoredRun(t, empty)
			result, err := InspectTerminalRunMerged(context.Background(), run, nil)
			if err != nil || result.State != ReconciliationSafe || result.Reason != reconciliationReasonWorktreeIgnoredOnly {
				t.Fatalf("inspection = %+v, %v", result, err)
			}
			want := []string(nil)
			if !empty {
				want = []string{".storage/postgres/18/docker/base/1"}
			}
			if result.evidence == nil || result.evidence.worktreePresent || result.evidence.runBranchPresent || !slices.Equal(result.evidence.ignoredPaths, want) {
				t.Fatalf("inspection evidence = %+v, want absent Git surfaces and paths %v", result.evidence, want)
			}
			runStore := &recordingTerminalRunStore{}
			if err := ApplyTerminalRun(context.Background(), runStore, result); err != nil {
				t.Fatal(err)
			}
			if len(runStore.requests) != 0 {
				t.Fatalf("unexpected Run Database records: %+v", runStore.requests)
			}
			if _, err := os.Lstat(run.WorkDir); !os.IsNotExist(err) {
				t.Fatalf("removed directory stat = %v", err)
			}
			next, err := InspectTerminalRunMerged(context.Background(), run, nil)
			if err != nil || next.State != ReconciliationReleased {
				t.Fatalf("next inspection = %+v, %v", next, err)
			}
		})
	}
}

func TestAnUnregisteredRunWorktreeKeepsTodaysUnknown(t *testing.T) {
	t.Parallel()
	cases := map[string]func(*testing.T, *store.Run){
		"unignored file": func(t *testing.T, run *store.Run) {
			mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, "notes.txt"), "keep\n")
		},
		"Git metadata": func(t *testing.T, run *store.Run) {
			mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, ".git"), "gitdir: missing\n")
		},
		"Git metadata directory": func(t *testing.T, run *store.Run) {
			if err := os.Mkdir(filepath.Join(run.WorkDir, ".git"), 0o755); err != nil {
				t.Fatal(err)
			}
		},
		"dangling Git metadata symlink": func(t *testing.T, run *store.Run) {
			if err := os.Symlink(filepath.Join(run.WorkDir, "missing"), filepath.Join(run.WorkDir, ".git")); err != nil {
				t.Fatal(err)
			}
		},
		"Run Branch": func(t *testing.T, run *store.Run) {
			gitWorktreeTest(t, run.GitRoot, "branch", BranchName(run.ID))
		},
		"tracked path matching ignore": func(t *testing.T, run *store.Run) {
			mustWriteWorktreeTest(t, filepath.Join(run.GitRoot, "tracked.txt"), "tracked\n")
			gitWorktreeTest(t, run.GitRoot, "add", "tracked.txt")
			gitWorktreeTest(t, run.GitRoot, "commit", "-m", "track path")
			mustWriteWorktreeTest(t, filepath.Join(run.GitRoot, ".gitignore"), ".storage/\ntracked.txt\n")
			gitWorktreeTest(t, run.GitRoot, "add", ".gitignore")
			gitWorktreeTest(t, run.GitRoot, "commit", "-m", "ignore tracked path")
			mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, "tracked.txt"), "keep\n")
		},
		"wrong directory name": func(t *testing.T, run *store.Run) {
			other := filepath.Join(filepath.Dir(run.WorkDir), "another-run")
			if err := os.Rename(run.WorkDir, other); err != nil {
				t.Fatal(err)
			}
			run.WorkDir = other
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			run := newUnregisteredIgnoredRun(t, false)
			change(t, &run)
			result, err := InspectTerminalRunMerged(context.Background(), run, nil)
			if err != nil || result.State != ReconciliationUnknown || result.Reason != reconciliationReasonWorktreeUnregistered {
				t.Fatalf("inspection = %+v, %v", result, err)
			}
		})
	}
}

func TestApplyRefusesAnIgnoredOnlyWorktreeWhoseFilesChanged(t *testing.T) {
	t.Parallel()
	run := newUnregisteredIgnoredRun(t, false)
	result, err := InspectTerminalRunMerged(context.Background(), run, nil)
	if err != nil || result.State != ReconciliationSafe {
		t.Fatalf("inspection = %+v, %v", result, err)
	}
	mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, ".storage", "new"), "new\n")
	runStore := &recordingTerminalRunStore{}
	if err := ApplyTerminalRun(context.Background(), runStore, result); err == nil {
		t.Fatal("apply accepted a changed path list")
	}
	if _, err := os.Stat(run.WorkDir); err != nil || len(runStore.requests) != 0 {
		t.Fatalf("refusal changed retained surfaces: stat=%v records=%+v", err, runStore.requests)
	}
}

func TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	want := []string{filepath.Join(root, "fiscus-0035-implementation-20260908-392a816f"), filepath.Join(root, "fiscus-2f7556b8")}
	for _, name := range append(slices.Clone(want), filepath.Join(root, "notes"), filepath.Join(root, "repo-0123abcd")) {
		if err := os.Mkdir(name, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mustWriteWorktreeTest(t, filepath.Join(root, "repo-0123abcd", "keep"), "keep\n")
	link := filepath.Join(root, "link-0123abcd")
	if err := os.Symlink(filepath.Join(root, "notes"), link); err != nil {
		t.Fatal(err)
	}
	got, err := EmptyRepositoryDirectories(root)
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("empty directories = %v, %v; want %v", got, err, want)
	}
	newFile := filepath.Join(want[0], "new-entry")
	mustWriteWorktreeTest(t, newFile, "keep\n")
	if err := RemoveEmptyRepositoryDirectory(root, want[0]); err == nil {
		t.Fatal("removed a repository directory that gained an entry")
	}
	if err := os.Remove(newFile); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{filepath.Join(root, "notes"), filepath.Join(root, "repo-0123abcd"), link, filepath.Join(root, "..", "outside-0123abcd"), root} {
		if err := RemoveEmptyRepositoryDirectory(root, name); err == nil {
			t.Fatalf("removed invalid candidate %q", name)
		}
	}
	for _, name := range want {
		if err := RemoveEmptyRepositoryDirectory(root, name); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"notes", "repo-0123abcd", "link-0123abcd"} {
		if _, err := os.Lstat(filepath.Join(root, name)); err != nil {
			t.Fatalf("preserved %s: %v", name, err)
		}
	}
	if got, err := EmptyRepositoryDirectories(filepath.Join(root, "absent")); err != nil || len(got) != 0 {
		t.Fatalf("absent root = %v, %v", got, err)
	}
}

func newUnregisteredIgnoredRun(t *testing.T, empty bool) store.Run {
	t.Helper()
	repo := initWorktreeRepo(t)
	repo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	mustWriteWorktreeTest(t, filepath.Join(repo, ".gitignore"), ".storage/\n")
	gitWorktreeTest(t, repo, "add", ".gitignore")
	gitWorktreeTest(t, repo, "commit", "-m", "ignore runtime storage")
	path := filepath.Join(t.TempDir(), "fiscus-2f7556b8", "ignored-run")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if !empty {
		if err := os.MkdirAll(filepath.Join(path, ".storage", "postgres", "18", "docker", "base"), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteWorktreeTest(t, filepath.Join(path, ".storage", "postgres", "18", "docker", "base", "1"), "volume\n")
	}
	return store.Run{ID: "ignored-run", Kind: store.KindImplement, State: store.StateFailed, GitRoot: repo, WorkDir: path, LocalBranch: "main"}
}

func TestIgnoredOnlyReleaseDoesNotFollowSymbolicLinks(t *testing.T) {
	t.Parallel()
	run := newUnregisteredIgnoredRun(t, false)
	outside := t.TempDir()
	keep := filepath.Join(outside, "notes.txt")
	mustWriteWorktreeTest(t, keep, "keep\n")
	if err := os.Symlink(outside, filepath.Join(run.WorkDir, ".storage", "link")); err != nil {
		t.Fatal(err)
	}
	// Whitespace must survive the stdin and stdout NUL framing.
	mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, ".storage", "a file"), "ignored\n")
	result, err := InspectTerminalRunMerged(context.Background(), run, nil)
	want := []string{".storage/a file", ".storage/link", ".storage/postgres/18/docker/base/1"}
	if err != nil || result.State != ReconciliationSafe || result.evidence == nil || !slices.Equal(result.evidence.ignoredPaths, want) {
		t.Fatalf("inspection = %+v, %v; want paths %q", result, err, want)
	}
	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(keep); err != nil || string(data) != "keep\n" {
		t.Fatalf("symlink target changed: %q, %v", data, err)
	}
}

func TestIgnoredOnlyApplyReprovesGitMetadataAndIgnoreRules(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"metadata", "branch", "ignore rules"} {
		t.Run(change, func(t *testing.T) {
			run := newUnregisteredIgnoredRun(t, false)
			result, err := InspectTerminalRunMerged(context.Background(), run, nil)
			if err != nil || result.State != ReconciliationSafe {
				t.Fatalf("inspection = %+v, %v", result, err)
			}
			switch change {
			case "metadata":
				mustWriteWorktreeTest(t, filepath.Join(run.WorkDir, ".git"), "keep\n")
			case "branch":
				gitWorktreeTest(t, run.GitRoot, "branch", BranchName(run.ID))
			case "ignore rules":
				mustWriteWorktreeTest(t, filepath.Join(run.GitRoot, ".gitignore"), "")
			}
			runStore := &recordingTerminalRunStore{}
			if err := ApplyTerminalRun(context.Background(), runStore, result); err == nil {
				t.Fatal("apply accepted stale Git proof")
			}
			if _, err := os.Stat(run.WorkDir); err != nil || len(runStore.requests) != 0 {
				t.Fatalf("refusal changed surfaces: stat=%v records=%+v", err, runStore.requests)
			}
		})
	}
}
