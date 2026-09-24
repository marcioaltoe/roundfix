package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"

	_ "modernc.org/sqlite"
)

type output struct {
	Mode                   string `json:"mode"`
	Root                   string `json:"root"`
	Home                   string `json:"home"`
	Main                   string `json:"main"`
	Linked                 string `json:"linked,omitempty"`
	Other                  string `json:"other,omitempty"`
	RunID                  string `json:"runId,omitempty"`
	OtherRunID             string `json:"otherRunId,omitempty"`
	ResolvableRunID        string `json:"resolvableRunId,omitempty"`
	UnresolvableRunID      string `json:"unresolvableRunId,omitempty"`
	SharedArtifactRoot     string `json:"sharedArtifactRoot,omitempty"`
	RecordedRepositoryRoot string `json:"recordedRepositoryRoot,omitempty"`
	CommonDirBefore        string `json:"commonDirBefore,omitempty"`
	WorktreeListBefore     string `json:"worktreeListBefore,omitempty"`
	WorktreeListAfter      string `json:"worktreeListAfter,omitempty"`
	ResolvableRoot         string `json:"resolvableRoot,omitempty"`
	UnresolvableRoot       string `json:"unresolvableRoot"`
	SchemaMatchesFresh     bool   `json:"schemaMatchesFresh,omitempty"`
	WindowRows             int    `json:"windowRows"`
}

func main() {
	mode := flag.String("mode", "", "fixture mode")
	root := flag.String("root", "", "absolute scratch root")
	flag.Parse()
	if *mode == "" || !filepath.IsAbs(*root) {
		fatalf("--mode and an absolute --root are required")
	}

	var (
		result output
		err    error
	)
	switch *mode {
	case "setup-live":
		result, err = setupLive(*root)
	case "setup-reconcile":
		result, err = setupReconcile(*root)
	case "seed-window":
		result, err = seedWindow(*root)
	case "inspect-live":
		result, err = inspectLive(*root)
	case "setup-migration":
		result, err = setupMigration(*root)
	case "inspect-migration":
		result, err = inspectMigration(*root)
	default:
		err = fmt.Errorf("unknown mode %q", *mode)
	}
	if err != nil {
		fatalf("%v", err)
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fatalf("encode result: %v", err)
	}
	fmt.Println(string(encoded))
}

func seedWindow(root string) (output, error) {
	homeDir := filepath.Join(root, "home")
	linkedRoot := filepath.Join(root, "linked")
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		return output{}, err
	}
	defer db.Close()
	cutoff := time.Now().UTC().Add(3 * time.Hour).Truncate(time.Minute)
	if _, err := db.Exec(
		`INSERT INTO run_windows (git_root, cutoff_at, created_at) VALUES (?, ?, ?)`,
		linkedRoot,
		cutoff.Unix(),
		time.Now().UTC().Add(-time.Hour).Unix(),
	); err != nil {
		return output{}, err
	}
	return output{Mode: "seed-window", Root: root, Home: homeDir, Linked: linkedRoot, WindowRows: 1}, nil
}

func setupReconcile(root string) (output, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	linkedRoot := filepath.Join(root, "linked")
	otherRoot := filepath.Join(root, "other")
	homeDir := filepath.Join(root, "home")
	if err := initRepository(mainRoot); err != nil {
		return output{}, err
	}
	if _, err := git(mainRoot, "worktree", "add", "-b", "feature/linked", linkedRoot); err != nil {
		return output{}, err
	}
	if err := initRepository(otherRoot); err != nil {
		return output{}, err
	}
	head, err := git(linkedRoot, "rev-parse", "HEAD")
	if err != nil {
		return output{}, err
	}
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return output{}, err
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: linkedRoot, LocalBranch: "feature/linked",
		HeadSHA: strings.TrimSpace(head), SpecSlug: "qa-0162-reconcile", Agent: "codex",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped); err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	other, err := runStore.CreateRunSkippingActiveLock(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: otherRoot, LocalBranch: "main",
		SpecSlug: "qa-0162-reconcile-other", Agent: "codex",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	if _, err := runStore.CompleteRun(ctx, other.ID, store.StateStopped); err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	if err := runStore.Close(); err != nil {
		return output{}, err
	}
	return output{
		Mode: "setup-reconcile", Root: root, Home: homeDir, Main: mainRoot,
		Linked: linkedRoot, Other: otherRoot, RunID: run.ID, OtherRunID: other.ID,
		RecordedRepositoryRoot: run.RepositoryRoot,
	}, nil
}

func setupLive(root string) (output, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	linkedRoot := filepath.Join(root, "linked")
	otherRoot := filepath.Join(root, "other")
	homeDir := filepath.Join(root, "home")
	if err := initRepository(mainRoot); err != nil {
		return output{}, err
	}
	if _, err := git(mainRoot, "worktree", "add", "-b", "feature/linked", linkedRoot); err != nil {
		return output{}, err
	}
	if err := initRepository(otherRoot); err != nil {
		return output{}, err
	}
	sharedRoot, err := roundconfig.ResolveArtifactDirectory("", mainRoot, homeDir)
	if err != nil {
		return output{}, fmt.Errorf("resolve shared Artifact Root: %w", err)
	}
	head, err := git(linkedRoot, "rev-parse", "HEAD")
	if err != nil {
		return output{}, err
	}

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return output{}, fmt.Errorf("open Run Database: %w", err)
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     linkedRoot,
		LocalBranch: "feature/linked",
		HeadSHA:     strings.TrimSpace(head),
		ArtifactDir: sharedRoot,
		SpecSlug:    "qa-0162-linked",
		Agent:       "codex",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, fmt.Errorf("create linked-worktree Run: %w", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped); err != nil {
		_ = runStore.Close()
		return output{}, fmt.Errorf("complete linked-worktree Run: %w", err)
	}
	other, err := runStore.CreateRunSkippingActiveLock(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     otherRoot,
		LocalBranch: "main",
		ArtifactDir: filepath.Join(root, "other-artifacts"),
		SpecSlug:    "qa-0162-other",
		Agent:       "codex",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, fmt.Errorf("create unrelated Run: %w", err)
	}
	if _, err := runStore.CompleteRun(ctx, other.ID, store.StateStopped); err != nil {
		_ = runStore.Close()
		return output{}, fmt.Errorf("complete unrelated Run: %w", err)
	}
	if err := runStore.Close(); err != nil {
		return output{}, fmt.Errorf("close Run Database: %w", err)
	}

	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		return output{}, fmt.Errorf("open Run Database for legacy window: %w", err)
	}
	cutoff := time.Now().UTC().Add(3 * time.Hour).Truncate(time.Minute)
	if _, err := db.Exec(
		`INSERT INTO run_windows (git_root, cutoff_at, created_at) VALUES (?, ?, ?)`,
		linkedRoot,
		cutoff.Unix(),
		time.Now().UTC().Add(-time.Hour).Unix(),
	); err != nil {
		_ = db.Close()
		return output{}, fmt.Errorf("seed worktree-keyed Run Window: %w", err)
	}
	if err := db.Close(); err != nil {
		return output{}, fmt.Errorf("close legacy-window database: %w", err)
	}
	artifactDir := filepath.Join(sharedRoot, run.ID)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		return output{}, fmt.Errorf("create Run artifact directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "marker.txt"), []byte("linked Run artifact\n"), 0o644); err != nil {
		return output{}, fmt.Errorf("write Run artifact: %w", err)
	}
	commonDir, err := git(linkedRoot, "rev-parse", "--git-common-dir")
	if err != nil {
		return output{}, err
	}
	worktreeListBefore, err := git(mainRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return output{}, err
	}
	if _, err := git(mainRoot, "worktree", "remove", linkedRoot); err != nil {
		return output{}, err
	}
	worktreeListAfter, err := git(mainRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return output{}, err
	}

	return output{
		Mode:                   "setup-live",
		Root:                   root,
		Home:                   homeDir,
		Main:                   mainRoot,
		Linked:                 linkedRoot,
		Other:                  otherRoot,
		RunID:                  run.ID,
		OtherRunID:             other.ID,
		SharedArtifactRoot:     sharedRoot,
		RecordedRepositoryRoot: run.RepositoryRoot,
		CommonDirBefore:        strings.TrimSpace(commonDir),
		WorktreeListBefore:     worktreeListBefore,
		WorktreeListAfter:      worktreeListAfter,
	}, nil
}

func inspectLive(root string) (output, error) {
	homeDir := filepath.Join(root, "home")
	mainRoot := filepath.Join(root, "main")
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		return output{}, err
	}
	defer db.Close()
	var rows int
	if err := db.QueryRow(`SELECT COUNT(*) FROM run_windows WHERE git_root IN (?, ?)`, mainRoot, filepath.Join(root, "linked")).Scan(&rows); err != nil {
		return output{}, fmt.Errorf("count Run Window rows: %w", err)
	}
	return output{Mode: "inspect-live", Root: root, Home: homeDir, Main: mainRoot, WindowRows: rows}, nil
}

func setupMigration(root string) (output, error) {
	ctx := context.Background()
	mainRoot := filepath.Join(root, "main")
	linkedRoot := filepath.Join(root, "linked")
	homeDir := filepath.Join(root, "home")
	if err := initRepository(mainRoot); err != nil {
		return output{}, err
	}
	if _, err := git(mainRoot, "worktree", "add", "-b", "feature/migration", linkedRoot); err != nil {
		return output{}, err
	}
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return output{}, err
	}
	resolvable, err := runStore.CreateRunSkippingActiveLock(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: linkedRoot, LocalBranch: "feature/migration", SpecSlug: "qa-0162-migration-resolvable",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	unresolvable, err := runStore.CreateRunSkippingActiveLock(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: mainRoot, LocalBranch: "main", SpecSlug: "qa-0162-migration-unresolvable",
	})
	if err != nil {
		_ = runStore.Close()
		return output{}, err
	}
	if err := runStore.Close(); err != nil {
		return output{}, err
	}
	unresolvableRoot := filepath.Join(root, "unresolvable")
	if err := os.MkdirAll(unresolvableRoot, 0o755); err != nil {
		return output{}, err
	}
	if err := os.WriteFile(filepath.Join(unresolvableRoot, ".git"), []byte("not a gitdir pointer\n"), 0o644); err != nil {
		return output{}, err
	}
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		return output{}, err
	}
	statements := []struct {
		query string
		args  []any
	}{
		{query: `UPDATE runs SET git_root = ? WHERE id = ?`, args: []any{unresolvableRoot, unresolvable.ID}},
		{query: `ALTER TABLE runs DROP COLUMN repository_root`},
		{query: `PRAGMA user_version = 16`},
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement.query, statement.args...); err != nil {
			_ = db.Close()
			return output{}, fmt.Errorf("build v16 migration fixture: %w", err)
		}
	}
	if err := db.Close(); err != nil {
		return output{}, err
	}
	return output{
		Mode: "setup-migration", Root: root, Home: homeDir, Main: mainRoot, Linked: linkedRoot,
		ResolvableRunID: resolvable.ID, UnresolvableRunID: unresolvable.ID, UnresolvableRoot: unresolvableRoot,
	}, nil
}

func inspectMigration(root string) (output, error) {
	ctx := context.Background()
	homeDir := filepath.Join(root, "home")
	freshHome := filepath.Join(root, "fresh-home")
	mainRoot := filepath.Join(root, "main")
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir))
	if err != nil {
		return output{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT repository_root FROM runs ORDER BY spec_slug`)
	if err != nil {
		return output{}, err
	}
	defer rows.Close()
	roots := make([]string, 0, 2)
	for rows.Next() {
		var rootValue string
		if err := rows.Scan(&rootValue); err != nil {
			return output{}, err
		}
		roots = append(roots, rootValue)
	}
	if err := rows.Err(); err != nil {
		return output{}, err
	}
	if len(roots) != 2 {
		return output{}, fmt.Errorf("migrated Run count = %d, want 2", len(roots))
	}
	migratedSchema, err := schema(db)
	if err != nil {
		return output{}, err
	}
	fresh, err := store.Open(ctx, freshHome)
	if err != nil {
		return output{}, err
	}
	if err := fresh.Close(); err != nil {
		return output{}, err
	}
	freshDB, err := sql.Open("sqlite", "file:"+store.DatabasePath(freshHome))
	if err != nil {
		return output{}, err
	}
	defer freshDB.Close()
	freshSchema, err := schema(freshDB)
	if err != nil {
		return output{}, err
	}
	return output{
		Mode: "inspect-migration", Root: root, Home: homeDir, Main: mainRoot,
		ResolvableRoot: roots[0], UnresolvableRoot: roots[1], SchemaMatchesFresh: migratedSchema == freshSchema,
	}, nil
}

func schema(db *sql.DB) (string, error) {
	rows, err := db.Query(`SELECT type, name, tbl_name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' ORDER BY type, name`)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var builder strings.Builder
	for rows.Next() {
		var kind, name, tableName, statement string
		if err := rows.Scan(&kind, &name, &tableName, &statement); err != nil {
			return "", err
		}
		fmt.Fprintf(&builder, "%s|%s|%s|%s\n", kind, name, tableName, statement)
	}
	return builder.String(), rows.Err()
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
	if _, err := git(root, "config", "commit.gpgSign", "false"); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		return err
	}
	if _, err := git(root, "add", "README.md"); err != nil {
		return err
	}
	_, err := git(root, "commit", "-m", "initial fixture")
	return err
}

func git(workDir string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = workDir
	combined, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(combined)))
	}
	return string(combined), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
