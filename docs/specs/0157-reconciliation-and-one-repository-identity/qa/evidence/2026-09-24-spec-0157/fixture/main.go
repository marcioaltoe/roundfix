package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"

	_ "modernc.org/sqlite"
)

const fixtureSpec = "qa-0157-fixture"

type fixtureOutput struct {
	Mode                 string         `json:"mode"`
	Root                 string         `json:"root"`
	Home                 string         `json:"home"`
	Main                 string         `json:"main"`
	Linked               string         `json:"linked,omitempty"`
	RunID                string         `json:"runId,omitempty"`
	LegacyRunID          string         `json:"legacyRunId,omitempty"`
	RunWorktree          string         `json:"runWorktree,omitempty"`
	RunBranch            string         `json:"runBranch,omitempty"`
	TargetBranch         string         `json:"targetBranch,omitempty"`
	ReportName           string         `json:"reportName,omitempty"`
	ActiveReport         string         `json:"activeReport,omitempty"`
	ArchivedReport       string         `json:"archivedReport,omitempty"`
	MainArtifact         string         `json:"mainArtifact,omitempty"`
	LinkedArtifact       string         `json:"linkedArtifact,omitempty"`
	ExpectedMainArtifact string         `json:"expectedMainArtifact,omitempty"`
	LegacyArtifact       string         `json:"legacyArtifact,omitempty"`
	LegacyMarker         string         `json:"legacyMarker,omitempty"`
	Listed               []listedOutput `json:"listed,omitempty"`
}

type listedOutput struct {
	RunID       string `json:"runId"`
	GitRoot     string `json:"gitRoot"`
	ArtifactDir string `json:"artifactDir"`
	Marker      string `json:"marker,omitempty"`
}

func main() {
	mode := flag.String("mode", "", "fixture mode")
	root := flag.String("root", "", "absolute scratch root")
	flag.Parse()
	if strings.TrimSpace(*mode) == "" || !filepath.IsAbs(*root) {
		fatalf("--mode and an absolute --root are required")
	}

	var (
		output fixtureOutput
		err    error
	)
	switch *mode {
	case "unproven":
		output, err = setupReconciliation(*root, true)
	case "archived-copy":
		output, err = setupReconciliation(*root, false)
	case "identity":
		output, err = setupIdentity(*root)
	case "inspect-identity":
		output, err = inspectIdentity(*root)
	default:
		err = fmt.Errorf("unknown mode %q", *mode)
	}
	if err != nil {
		fatalf("%v", err)
	}
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fatalf("encode fixture result: %v", err)
	}
	fmt.Println(string(encoded))
}

func setupReconciliation(root string, deleteTarget bool) (fixtureOutput, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	homeDir := filepath.Join(root, "home")
	if err := initRepository(mainRoot); err != nil {
		return fixtureOutput{}, err
	}
	base, err := git(mainRoot, "rev-parse", "HEAD")
	if err != nil {
		return fixtureOutput{}, err
	}
	target := "main"
	if deleteTarget {
		target = "fix/deleted-target"
	}
	artifactDir, err := roundconfig.ResolveArtifactDirectory("", mainRoot, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	defer runStore.Close()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     mainRoot,
		LocalBranch: target,
		HeadSHA:     strings.TrimSpace(base),
		ArtifactDir: artifactDir,
		SpecSlug:    fixtureSpec,
	})
	if err != nil {
		return fixtureOutput{}, err
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: mainRoot,
		Location: filepath.Join(root, "worktrees"),
		RunID:    run.ID,
		HeadSHA:  strings.TrimSpace(base),
	})
	if err != nil {
		return fixtureOutput{}, err
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		return fixtureOutput{}, err
	}

	reportName := "qa-report-2026-07-28.md"
	activeReport := filepath.Join(ref.Path, "docs", "specs", fixtureSpec, "qa", reportName)
	if err := commitReport(ref.Path, activeReport, reportName, "fail"); err != nil {
		return fixtureOutput{}, err
	}

	if deleteTarget {
		if _, err := git(mainRoot, "branch", target, "main"); err != nil {
			return fixtureOutput{}, err
		}
		targetRoot := filepath.Join(root, "target")
		if _, err := git(mainRoot, "worktree", "add", targetRoot, target); err != nil {
			return fixtureOutput{}, err
		}
		if err := writeCommit(targetRoot, "target.txt", "delivered\n", "target work"); err != nil {
			return fixtureOutput{}, err
		}
		if _, err := git(mainRoot, "merge", "--squash", target); err != nil {
			return fixtureOutput{}, err
		}
		if _, err := git(mainRoot, "commit", "-m", "squash merge target"); err != nil {
			return fixtureOutput{}, err
		}
		if _, err := git(mainRoot, "worktree", "remove", targetRoot); err != nil {
			return fixtureOutput{}, err
		}
		if _, err := git(mainRoot, "branch", "-D", target); err != nil {
			return fixtureOutput{}, err
		}
		reportName = "qa-report-2026-07-29.md"
	}
	archivedReport := filepath.Join(mainRoot, "docs", "history", "specs", fixtureSpec, "qa", reportName)
	if err := commitReport(mainRoot, archivedReport, reportName, "pass"); err != nil {
		return fixtureOutput{}, err
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateUnresolved); err != nil {
		return fixtureOutput{}, err
	}
	return fixtureOutput{
		Mode:           map[bool]string{true: "unproven", false: "archived-copy"}[deleteTarget],
		Root:           root,
		Home:           homeDir,
		Main:           mainRoot,
		RunID:          run.ID,
		RunWorktree:    ref.Path,
		RunBranch:      ref.Branch,
		TargetBranch:   target,
		ReportName:     reportName,
		ActiveReport:   activeReport,
		ArchivedReport: archivedReport,
	}, nil
}

func setupIdentity(root string) (fixtureOutput, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	linkedRoot := filepath.Join(root, "linked")
	homeDir := filepath.Join(root, "home")
	if err := initRepository(mainRoot); err != nil {
		return fixtureOutput{}, err
	}
	if _, err := git(mainRoot, "worktree", "add", "-b", "feature/linked", linkedRoot); err != nil {
		return fixtureOutput{}, err
	}
	mainArtifact, err := roundconfig.ResolveArtifactDirectory("", mainRoot, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	linkedArtifact, err := roundconfig.ResolveArtifactDirectory("", linkedRoot, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	mainSum := sha256.Sum256([]byte(filepath.Clean(mainRoot)))
	expectedMain := filepath.Join(homeDir, ".roundfix", "artifacts", hex.EncodeToString(mainSum[:])[:16])
	legacySum := sha256.Sum256([]byte(filepath.Clean(linkedRoot)))
	legacyArtifact := filepath.Join(homeDir, ".roundfix", "artifacts", hex.EncodeToString(legacySum[:])[:16])
	if err := os.MkdirAll(legacyArtifact, 0o755); err != nil {
		return fixtureOutput{}, err
	}
	marker := filepath.Join(legacyArtifact, "earlier-run.txt")
	if err := os.WriteFile(marker, []byte("earlier worktree artifact\n"), 0o644); err != nil {
		return fixtureOutput{}, err
	}

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	current, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     linkedRoot,
		LocalBranch: "feature/linked",
		ArtifactDir: linkedArtifact,
		SpecSlug:    "qa-0157-current-identity",
	})
	if err != nil {
		return fixtureOutput{}, err
	}
	if _, err := runStore.CompleteRun(ctx, current.ID, store.StateClean); err != nil {
		return fixtureOutput{}, err
	}
	legacy, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     linkedRoot,
		LocalBranch: "feature/linked",
		ArtifactDir: legacyArtifact,
		SpecSlug:    "qa-0157-earlier-identity",
	})
	if err != nil {
		return fixtureOutput{}, err
	}
	if _, err := runStore.CompleteRun(ctx, legacy.ID, store.StateClean); err != nil {
		return fixtureOutput{}, err
	}
	if err := runStore.Close(); err != nil {
		return fixtureOutput{}, err
	}

	db, err := sql.Open("sqlite", store.DatabasePath(homeDir))
	if err != nil {
		return fixtureOutput{}, err
	}
	if _, err := db.ExecContext(ctx, `UPDATE runs SET git_root = ? WHERE id = ?`, linkedRoot, legacy.ID); err != nil {
		_ = db.Close()
		return fixtureOutput{}, err
	}
	if err := db.Close(); err != nil {
		return fixtureOutput{}, err
	}

	return fixtureOutput{
		Mode:                 "identity",
		Root:                 root,
		Home:                 homeDir,
		Main:                 mainRoot,
		Linked:               linkedRoot,
		RunID:                current.ID,
		LegacyRunID:          legacy.ID,
		MainArtifact:         mainArtifact,
		LinkedArtifact:       linkedArtifact,
		ExpectedMainArtifact: expectedMain,
		LegacyArtifact:       legacyArtifact,
		LegacyMarker:         marker,
	}, nil
}

func inspectIdentity(root string) (fixtureOutput, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	linkedRoot := filepath.Join(root, "linked")
	homeDir := filepath.Join(root, "home")
	mainArtifact, err := roundconfig.ResolveArtifactDirectory("", mainRoot, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	linkedArtifact, err := roundconfig.ResolveArtifactDirectory("", linkedRoot, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	mainSum := sha256.Sum256([]byte(filepath.Clean(mainRoot)))
	expectedMain := filepath.Join(homeDir, ".roundfix", "artifacts", hex.EncodeToString(mainSum[:])[:16])
	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return fixtureOutput{}, err
	}
	defer reader.Close()
	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{GitRoot: mainRoot, States: store.StatesAll})
	if err != nil {
		return fixtureOutput{}, err
	}
	listed := make([]listedOutput, 0, len(runs))
	for _, run := range runs {
		entry := listedOutput{RunID: run.ID, GitRoot: run.GitRoot, ArtifactDir: run.ArtifactDir}
		markerPath := filepath.Join(run.ArtifactDir, "earlier-run.txt")
		if content, readErr := os.ReadFile(markerPath); readErr == nil {
			entry.Marker = string(content)
		}
		listed = append(listed, entry)
	}
	return fixtureOutput{
		Mode:                 "inspect-identity",
		Root:                 root,
		Home:                 homeDir,
		Main:                 mainRoot,
		Linked:               linkedRoot,
		MainArtifact:         mainArtifact,
		LinkedArtifact:       linkedArtifact,
		ExpectedMainArtifact: expectedMain,
		Listed:               listed,
	}, nil
}

func initRepository(root string) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	if _, err := git(root, "init", "--initial-branch=main"); err != nil {
		return err
	}
	if _, err := git(root, "config", "user.name", "Roundfix QA"); err != nil {
		return err
	}
	if _, err := git(root, "config", "user.email", "qa@example.invalid"); err != nil {
		return err
	}
	return writeCommit(root, "README.md", "fixture\n", "initial")
}

func commitReport(workDir, absolutePath, reportName, verdict string) error {
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(absolutePath, []byte("---\nverdict: "+verdict+"\n---\n"), 0o644); err != nil {
		return err
	}
	relative, err := filepath.Rel(workDir, absolutePath)
	if err != nil {
		return err
	}
	if _, err := git(workDir, "add", relative); err != nil {
		return err
	}
	title := fmt.Sprintf("docs: qa report for %s (%s)", fixtureSpec, verdict)
	_, err = git(workDir, "commit", "-m", title, "-m", "Roundfix-Spec: "+fixtureSpec)
	_ = reportName
	return err
}

func writeCommit(workDir, relative, content, message string) error {
	path := filepath.Join(workDir, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	if _, err := git(workDir, "add", relative); err != nil {
		return err
	}
	_, err := git(workDir, "commit", "-m", message)
	return err
}

func git(workDir string, args ...string) (string, error) {
	commandArgs := append([]string{"-c", "commit.gpgSign=false"}, args...)
	command := exec.Command("git", commandArgs...)
	command.Dir = workDir
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
