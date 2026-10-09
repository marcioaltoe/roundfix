// Boundary: real local merge and Verification execution; origin is a temporary bare repo.
package cli

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

type conflictFixture struct {
	workflow   *commandDeliveryWorkflow
	repo, head string
}

func newConflictFixture(t *testing.T, sourceConflict bool, command, itemConfig string) conflictFixture {
	t.Helper()
	root := t.TempDir()
	remote := filepath.Join(root, "origin.git")
	repo := filepath.Join(root, "item")
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatal(err)
	}
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("derived.txt", "base\n")
	write("source.txt", "base\n")
	write(".roundfixrc.yml", "delivery:\n  derived_paths:\n    - paths: [derived.txt]\n      regenerate: "+command+"\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "base")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, repo, "checkout", "-b", "feat/item")
	write("derived.txt", "item\n")
	if sourceConflict {
		write("source.txt", "item source\n")
	}
	if itemConfig != "" {
		write(".roundfixrc.yml", itemConfig)
	}
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "item")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "main")
	write("derived.txt", "default\n")
	if sourceConflict {
		write("source.txt", "default source\n")
	}
	// A clean incoming source change must survive the derived merge.
	write("incoming.txt", "incoming\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "default")
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "checkout", "feat/item")
	config := roundconfig.Builtin()
	config.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	return conflictFixture{workflow: &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: config, HomeDir: home, GitRoot: repo, UserConfigPath: filepath.Join(home, ".roundfix", "config.yml")}}, repo: repo, head: head}
}
func resolveFixture(t *testing.T, fixture conflictFixture) delivery.ConflictResolution {
	t.Helper()
	result, err := fixture.workflow.ResolveConflict(t.Context(), fixture.repo, "example", fixture.head)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func assertConflictAborted(t *testing.T, fixture conflictFixture) {
	t.Helper()
	if head := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "HEAD")); head != fixture.head {
		t.Fatalf("head=%s want=%s", head, fixture.head)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty after abort=%q", dirty)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, ".git", "MERGE_HEAD")); !os.IsNotExist(err) {
		t.Fatalf("merge still active: %v", err)
	}
	if content, err := os.ReadFile(filepath.Join(fixture.repo, "derived.txt")); err != nil || string(content) != "item\n" {
		t.Fatalf("candidate file=%q err=%v", content, err)
	}
}
func TestADerivedConflictIsMergedAndRegenerated(t *testing.T) {
	t.Parallel()
	fixture := newConflictFixture(t, false, "printf 'regenerated\\n' > derived.txt", "")
	result := resolveFixture(t, fixture)
	if result.Head == "" || result.Head == fixture.head || len(result.SourcePaths) != 0 {
		t.Fatalf("result=%+v", result)
	}
	content, err := os.ReadFile(filepath.Join(fixture.repo, "derived.txt"))
	if err != nil || string(content) != "regenerated\n" {
		t.Fatalf("output=%q err=%v", content, err)
	}
	commit := gittest.Run(t, fixture.repo, "show", "-s", "--format=%B", "HEAD")
	if !strings.Contains(commit, "Roundfix-Delivery: derived-merge") {
		t.Fatalf("commit=%q", commit)
	}
	if parents := strings.Fields(gittest.Run(t, fixture.repo, "show", "-s", "--format=%P", "HEAD")); len(parents) != 2 || parents[0] != fixture.head {
		t.Fatalf("parents=%v", parents)
	}
	if content, err := os.ReadFile(filepath.Join(fixture.repo, "incoming.txt")); err != nil || string(content) != "incoming\n" {
		t.Fatalf("incoming=%q err=%v", content, err)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%q", dirty)
	}
}
func TestASourceConflictAbortsTheMerge(t *testing.T) {
	t.Parallel()
	fixture := newConflictFixture(t, true, "printf 'regenerated\\n' > derived.txt", "")
	result := resolveFixture(t, fixture)
	if !reflect.DeepEqual(result.SourcePaths, []string{"source.txt"}) || len(result.Regenerated) != 0 {
		t.Fatalf("result=%+v", result)
	}
	assertConflictAborted(t, fixture)
}
func TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"source.txt", "new.txt"} {
		t.Run(path, func(t *testing.T) {
			fixture := newConflictFixture(t, false, "printf 'regenerated\\n' > derived.txt; printf 'unexpected\\n' > "+path, "")
			result := resolveFixture(t, fixture)
			if !reflect.DeepEqual(result.SourcePaths, []string{"regenerated " + path + " outside delivery.derived_paths"}) {
				t.Fatalf("result=%+v", result)
			}
			assertConflictAborted(t, fixture)
		})
	}
}
func TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch(t *testing.T) {
	t.Parallel()
	for _, itemCommand := range []string{"printf 'item-command\\n' > derived.txt", "printf 'item-command\\n' > derived.txt; touch forbidden.txt"} {
		t.Run(itemCommand, func(t *testing.T) {
			fixture := newConflictFixture(t, false, "printf 'trusted\\n' > derived.txt", "delivery:\n  derived_paths:\n    - paths: [derived.txt, forbidden.txt]\n      regenerate: "+itemCommand+"\n")
			result := resolveFixture(t, fixture)
			if result.Head == "" {
				t.Fatalf("result=%+v", result)
			}
			content, err := os.ReadFile(filepath.Join(fixture.repo, "derived.txt"))
			if err != nil || string(content) != "trusted\n" {
				t.Fatalf("command output=%q err=%v", content, err)
			}
			if _, err := os.Stat(filepath.Join(fixture.repo, "forbidden.txt")); !os.IsNotExist(err) {
				t.Fatalf("item-only command ran: %v", err)
			}
		})
	}
}

func TestConflictRecoveryRequiresTheCleanCandidate(t *testing.T) {
	t.Parallel()
	for _, dirty := range []bool{false, true} {
		t.Run(map[bool]string{false: "wrong head", true: "dirty worktree"}[dirty], func(t *testing.T) {
			fixture := newConflictFixture(t, false, "printf 'regenerated\\n' > derived.txt", "")
			head := fixture.head
			if dirty {
				if err := os.WriteFile(filepath.Join(fixture.repo, "source.txt"), []byte("user edit"), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				head = "other-head"
			}
			_, err := fixture.workflow.ResolveConflict(t.Context(), fixture.repo, "example", head)
			if err == nil || !strings.Contains(err.Error(), "clean worktree at the candidate head") {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestConflictRecoveryAbortsAFailedRegeneration(t *testing.T) {
	t.Parallel()
	fixture := newConflictFixture(t, false, "printf 'partial\\n' > derived.txt; touch undeclared.txt; exit 1", "")
	_, err := fixture.workflow.ResolveConflict(t.Context(), fixture.repo, "example", fixture.head)
	if err == nil || !strings.Contains(err.Error(), "regenerate derived paths") {
		t.Fatalf("err=%v", err)
	}
	assertConflictAborted(t, fixture)
}

func TestConflictRecoveryDoesNotHideIgnoredNewPaths(t *testing.T) {
	t.Parallel()
	fixture := newConflictFixture(t, false, "printf 'regenerated\\n' > derived.txt; touch hidden.txt", "")
	if err := os.WriteFile(filepath.Join(fixture.repo, ".git", "info", "exclude"), []byte("hidden.txt\nexisting.txt\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture.repo, "existing.txt"), []byte("existing ignored content"), 0600); err != nil {
		t.Fatal(err)
	}
	result := resolveFixture(t, fixture)
	if !reflect.DeepEqual(result.SourcePaths, []string{"regenerated hidden.txt outside delivery.derived_paths"}) {
		t.Fatalf("result=%+v", result)
	}
	assertConflictAborted(t, fixture)
	if _, err := os.Stat(filepath.Join(fixture.repo, "hidden.txt")); !os.IsNotExist(err) {
		t.Fatalf("new ignored path remains: %v", err)
	}
	if content, err := os.ReadFile(filepath.Join(fixture.repo, "existing.txt")); err != nil || string(content) != "existing ignored content" {
		t.Fatalf("existing ignored content=%q err=%v", content, err)
	}
}

func TestConflictRecoveryRunsMatchedDeclarationsInOrder(t *testing.T) {
	t.Parallel()
	fixture := newConflictFixture(t, false, "printf initial > derived.txt", "")
	gittest.Run(t, fixture.repo, "checkout", "main")
	declaration := "delivery:\n  derived_paths:\n    - paths: [derived.txt, order.txt]\n      regenerate: printf first > order.txt; printf first > derived.txt\n    - paths: [derived.txt, order.txt]\n      regenerate: printf second >> order.txt; printf final > derived.txt\n    - paths: [unmatched.txt]\n      regenerate: touch unmatched.txt\n"
	if err := os.WriteFile(filepath.Join(fixture.repo, ".roundfixrc.yml"), []byte(declaration), 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, fixture.repo, "add", ".")
	gittest.Run(t, fixture.repo, "commit", "-m", "default declarations")
	gittest.Run(t, fixture.repo, "push", "origin", "main")
	gittest.Run(t, fixture.repo, "checkout", "feat/item")
	result := resolveFixture(t, fixture)
	if len(result.Regenerated) != 2 || result.Head == "" {
		t.Fatalf("result=%+v", result)
	}
	if output, err := os.ReadFile(filepath.Join(fixture.repo, "order.txt")); err != nil || string(output) != "firstsecond" {
		t.Fatalf("order=%q err=%v", output, err)
	}
	if output, err := os.ReadFile(filepath.Join(fixture.repo, "derived.txt")); err != nil || string(output) != "final" {
		t.Fatalf("derived=%q err=%v", output, err)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, "unmatched.txt")); !os.IsNotExist(err) {
		t.Fatalf("unmatched command executed: %v", err)
	}
}
