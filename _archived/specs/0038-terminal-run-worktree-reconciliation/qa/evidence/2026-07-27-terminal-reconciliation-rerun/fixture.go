//go:build qa

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

type fixture struct {
	RunID      string `json:"runId"`
	State      string `json:"state"`
	Repository string `json:"repository"`
	Worktree   string `json:"worktree,omitempty"`
	Branch     string `json:"branch,omitempty"`
}

type manifest struct {
	Root     string             `json:"root"`
	Home     string             `json:"home"`
	Database string             `json:"database"`
	Location string             `json:"location"`
	Repos    map[string]string  `json:"repos"`
	Fixtures map[string]fixture `json:"fixtures"`
}

func main() {
	if len(os.Args) < 2 {
		fatalf("usage: fixture seed <root> | fixture inspect <home> <run-id>")
	}
	switch os.Args[1] {
	case "seed":
		if len(os.Args) != 3 {
			fatalf("usage: fixture seed <root>")
		}
		seed(os.Args[2])
	case "inspect":
		if len(os.Args) != 4 {
			fatalf("usage: fixture inspect <home> <run-id>")
		}
		inspect(os.Args[2], os.Args[3])
	default:
		fatalf("unknown command %q", os.Args[1])
	}
}

func seed(root string) {
	ctx := context.Background()
	root = mustAbs(root)
	must(os.MkdirAll(root, 0o755))
	home := filepath.Join(root, "home")
	location := filepath.Join(root, "worktrees")
	must(os.MkdirAll(home, 0o755))
	must(os.MkdirAll(location, 0o755))

	m := manifest{
		Root:     root,
		Home:     home,
		Database: store.DatabasePath(home),
		Location: location,
		Repos:    map[string]string{},
		Fixtures: map[string]fixture{},
	}
	for _, name := range []string{"matrix", "integration", "failure", "selector", "other"} {
		repo := filepath.Join(root, name+"-repo")
		initRepo(repo)
		m.Repos[name] = repo
	}

	m.Fixtures["matrix_safe"] = createWorktreeRun(
		ctx, home, m.Repos["matrix"], location, "ma/qa-target", store.StateFailed, "safe",
	)
	m.Fixtures["matrix_dirty"] = createWorktreeRun(
		ctx, home, m.Repos["matrix"], location, "ma/qa-target", store.StateStopped, "dirty",
	)
	m.Fixtures["matrix_unintegrated"] = createWorktreeRun(
		ctx, home, m.Repos["matrix"], location, "ma/qa-target", store.StateUnresolved, "unintegrated",
	)
	m.Fixtures["matrix_unknown"] = createWorktreeRun(
		ctx, home, m.Repos["matrix"], location, "missing-target", store.StateTimedOut, "safe",
	)
	m.Fixtures["matrix_released"] = createMetadataRun(
		ctx, home, m.Repos["matrix"], "ma/qa-target", store.KindImplement, store.StateClean,
	)

	for name, state := range map[string]string{
		"integration_pending":    store.StateIntegrationPending,
		"integration_failed":     store.StateFailed,
		"integration_stopped":    store.StateStopped,
		"integration_timedout":   store.StateTimedOut,
		"integration_unresolved": store.StateUnresolved,
	} {
		m.Fixtures[name] = createWorktreeRun(
			ctx, home, m.Repos["integration"], location, "ma/qa-target", state, "safe",
		)
	}
	parent := m.Fixtures["integration_failed"]
	taskPath := parent.Worktree + "-task_01"
	taskBranch := parent.Branch + "-task_01"
	git(m.Repos["integration"], "worktree", "add", "-b", taskBranch, taskPath, "HEAD")
	m.Fixtures["integration_task_worktree"] = fixture{
		State:      "TaskWorktree",
		Repository: m.Repos["integration"],
		Worktree:   taskPath,
		Branch:     taskBranch,
	}

	m.Fixtures["failure_locked"] = createWorktreeRun(
		ctx, home, m.Repos["failure"], location, "ma/qa-target", store.StateFailed, "safe",
	)
	git(m.Repos["failure"], "worktree", "lock", m.Fixtures["failure_locked"].Worktree)

	m.Fixtures["selector_active"] = createMetadataRun(
		ctx, home, m.Repos["selector"], "ma/qa-target", store.KindImplement, "",
	)
	m.Fixtures["selector_review"] = createReviewRun(
		ctx, home, m.Repos["selector"], store.StateFailed,
	)
	m.Fixtures["other_terminal"] = createMetadataRun(
		ctx, home, m.Repos["other"], "ma/qa-target", store.KindImplement, store.StateFailed,
	)

	content, err := json.MarshalIndent(m, "", "  ")
	must(err)
	content = append(content, '\n')
	must(os.WriteFile(filepath.Join(root, "manifest.json"), content, 0o644))
	fmt.Print(string(content))
}

func inspect(home, runID string) {
	ctx := context.Background()
	reader, err := store.OpenReader(ctx, home)
	must(err)
	defer func() {
		must(reader.Close())
	}()
	run, found, err := reader.Run(ctx, runID)
	must(err)
	if !found {
		fatalf("Run %q not found", runID)
	}
	events, err := reader.RunEventsAfter(ctx, runID, 0, 100)
	must(err)
	output := struct {
		Run    store.Run            `json:"run"`
		Events []store.JournalEvent `json:"events"`
	}{Run: run, Events: events}
	content, err := json.MarshalIndent(output, "", "  ")
	must(err)
	fmt.Println(string(content))
}

func initRepo(repo string) {
	must(os.MkdirAll(repo, 0o755))
	git(repo, "init", "--initial-branch=main")
	git(repo, "config", "user.name", "Roundfix QA")
	git(repo, "config", "user.email", "roundfix-qa@example.com")
	git(repo, "config", "commit.gpgsign", "false")
	must(os.WriteFile(filepath.Join(repo, "README.md"), []byte("# QA fixture\n"), 0o644))
	git(repo, "add", "README.md")
	git(repo, "commit", "-m", "seed fixture")
	git(repo, "checkout", "-b", "ma/qa-target")
}

func createWorktreeRun(
	ctx context.Context,
	home, repo, location, targetBranch, terminalState, mode string,
) fixture {
	runStore, err := store.Open(ctx, home)
	must(err)
	defer func() {
		must(runStore.Close())
	}()
	head := strings.TrimSpace(gitOutput(repo, "rev-parse", "HEAD"))
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repo,
		LocalBranch: targetBranch,
		HeadSHA:     head,
		SpecSlug:    "reconcile-" + strings.ToLower(terminalState),
		Agent:       "codex",
		OwnerPID:    os.Getpid(),
	})
	must(err)
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repo,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  head,
	})
	must(err)
	ref.Path = mustAbs(ref.Path)
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	must(err)
	completed, err := runStore.CompleteRun(ctx, run.ID, terminalState)
	must(err)

	switch mode {
	case "safe":
	case "dirty":
		must(os.WriteFile(filepath.Join(ref.Path, "README.md"), []byte("# tracked dirt\n"), 0o644))
		must(os.WriteFile(filepath.Join(ref.Path, "untracked.txt"), []byte("untracked dirt\n"), 0o644))
	case "unintegrated":
		must(os.WriteFile(filepath.Join(ref.Path, "unique.txt"), []byte("unique run work\n"), 0o644))
		git(ref.Path, "add", "unique.txt")
		git(ref.Path, "commit", "-m", "unique run work")
	default:
		fatalf("unknown fixture mode %q", mode)
	}
	return fixture{
		RunID:      completed.Run.ID,
		State:      completed.Run.State,
		Repository: repo,
		Worktree:   ref.Path,
		Branch:     ref.Branch,
	}
}

func createMetadataRun(
	ctx context.Context,
	home, repo, branch, kind, terminalState string,
) fixture {
	runStore, err := store.Open(ctx, home)
	must(err)
	defer func() {
		must(runStore.Close())
	}()
	request := store.CreateRunRequest{
		Kind:        kind,
		GitRoot:     repo,
		LocalBranch: branch,
		HeadSHA:     strings.TrimSpace(gitOutput(repo, "rev-parse", "HEAD")),
		SpecSlug:    "metadata-" + strings.ToLower(strings.ReplaceAll(terminalState, " ", "-")),
		Agent:       "codex",
		OwnerPID:    os.Getpid(),
	}
	run, err := runStore.CreateRun(ctx, request)
	must(err)
	if terminalState != "" {
		completed, completeErr := runStore.CompleteRun(ctx, run.ID, terminalState)
		must(completeErr)
		run = completed.Run
	}
	return fixture{RunID: run.ID, State: run.State, Repository: repo}
}

func createReviewRun(ctx context.Context, home, repo, terminalState string) fixture {
	runStore, err := store.Open(ctx, home)
	must(err)
	defer func() {
		must(runStore.Close())
	}()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "ma/qa-target",
		BaseRepository: "owner/project",
		PRNumber:       "38",
		GitRoot:        repo,
		LocalBranch:    "ma/qa-target",
		HeadSHA:        strings.TrimSpace(gitOutput(repo, "rev-parse", "HEAD")),
		ArtifactDir:    filepath.Join(repo, ".roundfix", "reviews"),
		Agent:          "codex",
		OwnerPID:       os.Getpid(),
	})
	must(err)
	completed, err := runStore.CompleteRun(ctx, run.ID, terminalState)
	must(err)
	return fixture{RunID: completed.Run.ID, State: completed.Run.State, Repository: repo}
}

func git(dir string, args ...string) {
	_ = gitOutput(dir, args...)
}

func gitOutput(dir string, args ...string) string {
	cmd := exec.Command("git", append([]string{
		"-c", "user.name=Roundfix QA",
		"-c", "user.email=roundfix-qa@example.com",
		"-c", "commit.gpgsign=false",
	}, args...)...)
	cmd.Dir = dir
	content, err := cmd.CombinedOutput()
	if err != nil {
		fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, content)
	}
	return string(content)
}

func mustAbs(path string) string {
	absolute, err := filepath.Abs(path)
	must(err)
	resolved, err := filepath.EvalSymlinks(absolute)
	if err == nil {
		return resolved
	}
	if !errors.Is(err, os.ErrNotExist) {
		must(err)
	}
	return absolute
}

func must(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
