// Boundary: real local Git conflict recovery runs the shipped coverage command under the suite guard.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

const archivedConflictSlug = "0001-the-archived-item"
const archivedConflictCoveragePath = "docs/references/coverage-record.json"

func newArchivedItemConflictFixture(t *testing.T, sanctionCoverage bool) (conflictFixture, string, string) {
	t.Helper()
	root := t.TempDir()
	repo, remote, home := filepath.Join(root, "item"), filepath.Join(root, "origin.git"), filepath.Join(root, "home")
	copyArchivedConflictRepository(t, repo)
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, ".roundfixrc.yml"))
	if err != nil {
		t.Fatal(err)
	}
	project, err := roundconfig.ResolveConfigProposal(nil, data)
	if err != nil {
		t.Fatal(err)
	}
	var command string
	for _, declaration := range project.Delivery.DerivedPaths {
		if declaration.Matches(archivedConflictCoveragePath) {
			command = declaration.Regenerate
			break
		}
	}
	if command == "" {
		t.Fatal("coverage regeneration declaration not found")
	}
	coverage, err := os.ReadFile(filepath.Join(repo, archivedConflictCoveragePath))
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		Packages map[string][]string `json:"packages"`
	}
	if err := json.Unmarshal(coverage, &record); err != nil {
		t.Fatal(err)
	}
	names := record.Packages["roundfix/internal/spec"]
	if len(names) == 0 {
		t.Fatal("coverage record has no internal/spec tests to replace")
	}
	originalEntry := names[0]
	if strings.Count(string(coverage), strconv.Quote(originalEntry)) != 1 {
		t.Fatalf("coverage entry %q is not unique", originalEntry)
	}
	write := func(name, content string) {
		t.Helper()
		path := filepath.Join(repo, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	write("docs/specs/"+archivedConflictSlug+"/_prd.md", "# The archived item\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "base")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, repo, "checkout", "-b", "feat/item")
	write(archivedConflictCoveragePath, strings.Replace(string(coverage), strconv.Quote(originalEntry), `"TestItemSideEntry"`, 1))
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "item coverage edit")
	sourceRevision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	if err := os.RemoveAll(filepath.Join(repo, "docs/specs", archivedConflictSlug)); err != nil {
		t.Fatal(err)
	}
	sanction := "  - command: make baseline-digests\n"
	if sanctionCoverage {
		sanction = "  - command: " + strconv.Quote(command) + "\n    outputs:\n      - " + archivedConflictCoveragePath + "\n"
	}
	write("docs/history/specs/"+archivedConflictSlug+".md", fmt.Sprintf(`---
schema: roundfix/archive-record/v1
spec: %s
title: The archived item
status: archived
created: "2026-10-30"
archived: "2026-10-30"
source: docs/specs/%s
source_revision: %s
disposition: pass
qa_task: task_01
qa_report: qa/2026-10-30.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
promoted: []
regeneration:
%s---

# The archived item
`, archivedConflictSlug, archivedConflictSlug, sourceRevision, sanction))
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "archive item")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "main")
	write(archivedConflictCoveragePath, strings.Replace(string(coverage), strconv.Quote(originalEntry), `"TestDefaultSideEntry"`, 1))
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "default coverage edit")
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "checkout", "feat/item")
	cfg := roundconfig.Builtin()
	cfg.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	workflow := &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: cfg, HomeDir: home, GitRoot: repo, UserConfigPath: filepath.Join(home, ".roundfix", "config.yml")}}
	return conflictFixture{workflow: workflow, repo: repo, head: head}, command, originalEntry
}

// Copy inputs in-process so the Go test cache observes source changes. Only the
// coverage record enters from docs; every sanction belongs to the fixture itself.
func copyArchivedConflictRepository(t *testing.T, destination string) {
	t.Helper()
	source, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(source, ".roundfixrc.yml")); err == nil {
			break
		}
		parent := filepath.Dir(source)
		if parent == source {
			t.Fatal("repository root not found")
		}
		source = parent
	}
	copyFile := func(path, target string, mode os.FileMode) error {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, content, mode)
	}
	if err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		switch relative {
		case ".git", ".gocache", "bin", "docs", ".agents":
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode().Perm())
	}); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(filepath.Join(source, archivedConflictCoveragePath), filepath.Join(destination, archivedConflictCoveragePath), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAnArchivedItemsDerivedRecordConflictIsMergedAndRegenerated(t *testing.T) {
	t.Parallel()
	fixture, command, originalEntry := newArchivedItemConflictFixture(t, true)
	result, err := fixture.workflow.ResolveConflict(t.Context(), fixture.repo, archivedConflictSlug, fixture.head)
	if err != nil {
		t.Fatal(err)
	}
	if result.Head == "" || result.Head == fixture.head || len(result.SourcePaths) != 0 || !reflect.DeepEqual(result.Regenerated, []string{command}) {
		t.Fatalf("result=%+v", result)
	}
	content := gittest.Run(t, fixture.repo, "show", result.Head+":"+archivedConflictCoveragePath)
	if !strings.Contains(content, strconv.Quote(originalEntry)) {
		t.Fatalf("regenerated record did not restore %q", originalEntry)
	}
	for _, unwanted := range []string{"TestItemSideEntry", "TestDefaultSideEntry", "<<<<<<<", "=======", ">>>>>>>"} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("committed record still contains %q", unwanted)
		}
	}
	if head := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "HEAD")); head != result.Head {
		t.Fatalf("head=%s want=%s", head, result.Head)
	}
	mainHead := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "origin/main"))
	if parents := strings.Fields(gittest.Run(t, fixture.repo, "show", "-s", "--format=%P", result.Head)); !reflect.DeepEqual(parents, []string{fixture.head, mainHead}) {
		t.Fatalf("parents=%v", parents)
	}
	if commit := gittest.Run(t, fixture.repo, "show", "-s", "--format=%B", result.Head); !strings.Contains(commit, "Roundfix-Delivery: derived-merge") {
		t.Fatalf("commit=%q", commit)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%q", dirty)
	}
}

func TestAnArchivedItemWithoutTheSanctionKeepsItsDerivedRecordConflict(t *testing.T) {
	t.Parallel()
	fixture, _, _ := newArchivedItemConflictFixture(t, false)
	_, err := fixture.workflow.ResolveConflict(t.Context(), fixture.repo, archivedConflictSlug, fixture.head)
	if err == nil || !strings.Contains(err.Error(), "regenerate derived paths") {
		t.Fatalf("err=%v", err)
	}
	if head := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "HEAD")); head != fixture.head {
		t.Fatalf("head=%s want=%s", head, fixture.head)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty after abort=%q", dirty)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, ".git", "MERGE_HEAD")); !os.IsNotExist(err) {
		t.Fatalf("merge still active: %v", err)
	}
	logs, err := filepath.Glob(filepath.Join(fixture.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", archivedConflictSlug, "derived-regeneration-*.log"))
	if err != nil || len(logs) != 1 {
		t.Fatalf("regeneration logs=%v err=%v", logs, err)
	}
	content, err := os.ReadFile(logs[0])
	if err != nil || !strings.Contains(string(content), "suiteguard: modified: "+archivedConflictCoveragePath) {
		t.Fatalf("regeneration log=%s err=%v", content, err)
	}
}
