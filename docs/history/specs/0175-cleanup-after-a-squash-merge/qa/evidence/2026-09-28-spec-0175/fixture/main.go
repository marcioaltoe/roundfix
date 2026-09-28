package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"

	_ "modernc.org/sqlite"
)

const (
	slug   = "0172-reconcile-merged-spec"
	taskID = "task_05"
)

type scenario struct {
	Home       string `json:"home"`
	Repository string `json:"repository"`
	RunID      string `json:"runID,omitempty"`
	Worktree   string `json:"worktree,omitempty"`
	Branch     string `json:"branch,omitempty"`
	Commit     string `json:"commit,omitempty"`
}

func main() {
	if len(os.Args) != 3 {
		fatalf("usage: fixture <merged-record|merged-fallback|incomplete|ordinary|other-record|legacy|staging|hold-staging> <root>")
	}
	root, err := filepath.Abs(os.Args[2])
	must(err)
	must(os.MkdirAll(root, 0o755))

	var result scenario
	switch os.Args[1] {
	case "merged-record":
		result = setupMerged(root, "record")
	case "merged-fallback":
		result = setupMerged(root, "fallback")
	case "incomplete":
		result = setupMerged(root, "incomplete")
	case "ordinary":
		result = setupMerged(root, "ordinary")
	case "other-record":
		result = setupMerged(root, "other")
	case "legacy":
		result = setupLegacy(root)
	case "staging":
		result = setupStaging(root)
	case "hold-staging":
		result = holdStaging(root)
	default:
		fatalf("unknown fixture %q", os.Args[1])
	}
	must(json.NewEncoder(os.Stdout).Encode(result))
}

func setupMerged(root string, mode string) scenario {
	ctx := context.Background()
	home := filepath.Join(root, "home")
	repo := filepath.Join(root, "repo")
	location := filepath.Join(root, "run-worktrees")
	initRepo(repo)
	start := git(repo, "rev-parse", "HEAD")

	runStore, err := store.Open(ctx, home)
	must(err)
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: repo, LocalBranch: "delivery/" + slug,
		HeadSHA: start, SpecSlug: slug, Agent: "codex",
	})
	must(err)
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repo, Location: location, RunID: run.ID, HeadSHA: start,
	})
	must(err)
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	must(err)

	var taskFour string
	for index := 1; index <= 4; index++ {
		id := fmt.Sprintf("task_%02d", index)
		write(filepath.Join(ref.Path, fmt.Sprintf("feature-%02d.txt", index)), "delivered\n")
		git(ref.Path, "add", fmt.Sprintf("feature-%02d.txt", index))
		git(ref.Path, "commit", "-m", "feat: settle "+id, "-m", "Roundfix-Spec: "+slug+"\nRoundfix-Task: "+id)
		taskFour = git(ref.Path, "rev-parse", "HEAD")
	}
	git(repo, "reset", "--hard", taskFour)

	write(filepath.Join(ref.Path, "feature-05.txt"), "redone implementation\n")
	git(ref.Path, "add", "feature-05.txt")
	git(ref.Path, "commit", "-m", "feat: settle "+taskID, "-m", "Roundfix-Spec: "+slug+"\nRoundfix-Task: "+taskID)
	taskCommit := git(ref.Path, "rev-parse", "HEAD")
	if mode == "ordinary" {
		write(filepath.Join(ref.Path, "unrepresented.txt"), "unique run content\n")
		git(ref.Path, "add", "unrepresented.txt")
		git(ref.Path, "commit", "-m", "feat: keep unique run content")
		taskCommit = git(ref.Path, "rev-parse", "HEAD")
	}
	qaPath := filepath.Join(ref.Path, "docs", "specs", slug, "qa", "qa-report-2026-09-28.md")
	write(qaPath, "---\nverdict: fail\n---\n")
	git(ref.Path, "add", filepath.ToSlash(filepath.Join("docs", "specs", slug, "qa", "qa-report-2026-09-28.md")))
	git(ref.Path, "commit", "-m", "docs: qa report for "+slug+" (fail)", "-m", "Roundfix-Spec: "+slug)

	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped)
	must(err)
	must(runStore.Close())

	archiveSlug := slug
	if mode == "other" {
		archiveSlug = "0174-another-spec"
	}
	status := "completed"
	if mode == "incomplete" {
		status = "pending"
	}
	writeArchivedSpec(repo, archiveSlug, status)
	mergedHead := git(repo, "rev-parse", "HEAD")

	if mode != "fallback" {
		runStore, err = store.Open(ctx, home)
		must(err)
		queue, err := runStore.CreateDeliveryQueue(ctx, repo, []string{archiveSlug})
		must(err)
		item := queue.Items[0]
		item.Stage = store.DeliveryStageMerged
		item.Branch = "roundfix/delivery-" + archiveSlug
		item.CandidateCommits = []string{completed.Run.HeadSHA, mergedHead}
		item.PullRequestNumber = "259"
		item.MergeCommit = mergedHead
		must(runStore.UpdateDeliveryQueueItem(ctx, repo, item))
		must(runStore.Close())
	}

	return scenario{Home: home, Repository: repo, RunID: run.ID, Worktree: ref.Path, Branch: ref.Branch, Commit: taskCommit}
}

func writeArchivedSpec(repo string, archiveSlug string, taskStatus string) {
	root := filepath.Join(repo, "docs", "history", "specs", archiveSlug)
	must(os.MkdirAll(filepath.Join(root, "qa"), 0o755))
	var nodes strings.Builder
	for index := 1; index <= 5; index++ {
		id := fmt.Sprintf("task_%02d", index)
		fmt.Fprintf(&nodes, "    - id: %s\n      file: %s.md\n      needs: []\n", id, id)
		status := "completed"
		if id == taskID {
			status = taskStatus
		}
		write(filepath.Join(root, id+".md"), fmt.Sprintf("---\ntask: %s\nspec: %s\nstatus: %s\ntype: backend\ncomplexity: low\n---\n\n# %s\n\n## Verification\n\n- `true`\n", id, archiveSlug, status, id))
	}
	write(filepath.Join(root, "_tasks.md"), "---\nschema: spec-tasks/v1\nspec: "+archiveSlug+"\ngraph:\n  nodes:\n"+nodes.String()+"---\n\n# Task Graph\n")
	write(filepath.Join(root, "qa", "qa-report-2026-09-29.md"), "---\nverdict: pass\n---\n")
	git(repo, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", archiveSlug)))
	git(repo, "commit", "-m", "feat: record merged Spec")
}

func setupLegacy(root string) scenario {
	ctx := context.Background()
	home := filepath.Join(root, "home")
	repo := filepath.Join(root, "repo")
	location := filepath.Join(root, "run-worktrees")
	initRepo(repo)
	start := git(repo, "rev-parse", "HEAD")
	source := filepath.Join(root, "removed-source")
	git(repo, "worktree", "add", "-b", "feature/legacy-source", source, start)

	runStore, err := store.Open(ctx, home)
	must(err)
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: source, LocalBranch: "main", HeadSHA: start,
		SpecSlug: "reconcile-spec", Agent: "codex",
	})
	must(err)
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: source, Location: location, RunID: run.ID, HeadSHA: start,
	})
	must(err)
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	must(err)
	_, err = runStore.CompleteRun(ctx, run.ID, store.StateStopped)
	must(err)
	must(runStore.Close())

	database, err := sql.Open("sqlite", "file:"+store.DatabasePath(home))
	must(err)
	_, err = database.ExecContext(ctx, `UPDATE runs SET repository_root = '' WHERE id = ?`, run.ID)
	must(err)
	must(database.Close())
	git(repo, "worktree", "remove", source)

	foreign := filepath.Join(root, "foreign")
	initRepo(foreign)
	foreignHead := git(foreign, "rev-parse", "HEAD")
	runStore, err = store.Open(ctx, home)
	must(err)
	foreignRun, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: foreign, LocalBranch: "main", HeadSHA: foreignHead,
		SpecSlug: "foreign-spec", Agent: "codex",
	})
	must(err)
	foreignRef, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: foreign, Location: filepath.Join(root, "foreign-worktrees"), RunID: foreignRun.ID, HeadSHA: foreignHead,
	})
	must(err)
	_, err = runStore.SetRunWorkDir(ctx, foreignRun.ID, foreignRef.Path)
	must(err)
	_, err = runStore.CompleteRun(ctx, foreignRun.ID, store.StateStopped)
	must(err)
	must(runStore.Close())

	return scenario{Home: home, Repository: repo, RunID: run.ID, Worktree: ref.Path, Branch: ref.Branch, Commit: foreignRun.ID}
}

func setupStaging(root string) scenario {
	home := filepath.Join(root, "home")
	repo := filepath.Join(root, "repo")
	initRepo(repo)
	head := git(repo, "rev-parse", "HEAD")
	legacyRoot := filepath.Join(root, "roundfix-carry-forward-legacy")
	must(os.MkdirAll(legacyRoot, 0o755))
	legacy := filepath.Join(legacyRoot, "worktree")
	git(repo, "worktree", "add", "--detach", legacy, head)
	git(repo, "worktree", "lock", "--reason", "initializing", legacy)
	dead, err := runworktree.AddCarryForwardStaging(context.Background(), repo, head, root)
	must(err)
	return scenario{Home: home, Repository: repo, Worktree: legacy, Branch: dead.Worktree}
}

func holdStaging(root string) scenario {
	home := filepath.Join(root, "home")
	repo := filepath.Join(root, "repo")
	if _, err := os.Stat(filepath.Join(repo, ".git")); err != nil {
		initRepo(repo)
	}
	head := git(repo, "rev-parse", "HEAD")
	staging, err := runworktree.AddCarryForwardStaging(context.Background(), repo, head, root)
	must(err)
	result := scenario{Home: home, Repository: repo, Worktree: staging.Worktree}
	must(json.NewEncoder(os.Stdout).Encode(result))
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	must(staging.Remove(context.Background()))
	os.Exit(0)
	return scenario{}
}

func initRepo(repo string) {
	must(os.MkdirAll(repo, 0o755))
	git(repo, "init", "-b", "main")
	git(repo, "config", "user.name", "Roundfix QA")
	git(repo, "config", "user.email", "qa@example.invalid")
	git(repo, "config", "commit.gpgsign", "false")
	write(filepath.Join(repo, "README.md"), "fixture\n")
	git(repo, "add", "README.md")
	git(repo, "commit", "-m", "chore: initialize fixture")
}

func git(dir string, args ...string) string {
	command := exec.Command("git", args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, output)
	}
	return strings.TrimSpace(string(output))
}

func write(path string, value string) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	must(os.WriteFile(path, []byte(value), 0o644))
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
