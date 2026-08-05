package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"roundfix/internal/daemon"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const slug = "0066-run-teardown-reclaims-what-it-created"

type fixtureRun struct {
	ID       string `json:"id"`
	Branch   string `json:"branch"`
	Worktree string `json:"worktree"`
	Report   string `json:"report"`
}

type fixtureReport struct {
	Terminal []fixtureRun `json:"terminal"`
	Active   fixtureRun   `json:"active"`
}

func main() {
	if len(os.Args) != 4 {
		fatalf("usage: setup_reconcile <home> <repo> <worktrees>")
	}
	homeDir, repoDir, location := os.Args[1], os.Args[2], os.Args[3]
	ctx := context.Background()
	head := strings.TrimSpace(git(repoDir, "rev-parse", "HEAD"))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			fatalf("close Run Database: %v", err)
		}
	}()

	report := fixtureReport{Terminal: make([]fixtureRun, 0, 4)}
	for index := 1; index <= 4; index++ {
		report.Terminal = append(report.Terminal, createRun(ctx, runStore, repoDir, location, head, index, false))
	}
	report.Active = createRun(ctx, runStore, repoDir, location, head, 5, true)

	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		fatalf("encode fixture report: %v", err)
	}
}

func createRun(
	ctx context.Context,
	runStore *store.Store,
	repoDir string,
	location string,
	head string,
	index int,
	active bool,
) fixtureRun {
	runSlug := slug
	if active {
		runSlug = "0066-active-run-guard-fixture"
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/qa-reconcile",
		HeadSHA:     head,
		SpecSlug:    runSlug,
		Agent:       "codex",
	})
	if err != nil {
		fatalf("create Run %d: %v", index, err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir,
		Location: location,
		RunID:    run.ID,
		HeadSHA:  head,
	})
	if err != nil {
		fatalf("create Run Worktree %d: %v", index, err)
	}
	resolved, err := filepath.EvalSymlinks(ref.Path)
	if err != nil {
		fatalf("resolve Run Worktree %d: %v", index, err)
	}
	if _, err := runStore.SetRunWorkDir(ctx, run.ID, resolved); err != nil {
		fatalf("record Run Worktree %d: %v", index, err)
	}

	reportPath := filepath.ToSlash(filepath.Join("docs", "specs", runSlug, "qa", fmt.Sprintf("qa-report-2026-08-%02d.md", index)))
	absoluteReport := filepath.Join(resolved, filepath.FromSlash(reportPath))
	if err := os.MkdirAll(filepath.Dir(absoluteReport), 0o755); err != nil {
		fatalf("create report directory %d: %v", index, err)
	}
	verdict := "fail"
	if active {
		verdict = "pending"
	}
	content := fmt.Sprintf("---\nverdict: %s\n---\n\n# QA fixture report\n", verdict)
	if err := os.WriteFile(absoluteReport, []byte(content), 0o644); err != nil {
		fatalf("write report %d: %v", index, err)
	}
	git(resolved, "add", reportPath)
	git(resolved, "commit", "-m", daemon.QACommitMessage(runSlug, verdict))
	if !active {
		if _, err := runStore.CompleteRun(ctx, run.ID, store.StateUnresolved); err != nil {
			fatalf("complete Run %d: %v", index, err)
		}
	}
	return fixtureRun{ID: run.ID, Branch: ref.Branch, Worktree: resolved, Report: reportPath}
}

func git(dir string, args ...string) string {
	command := exec.Command("git", args...)
	command.Dir = dir
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		fatalf("git %s: %v\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
