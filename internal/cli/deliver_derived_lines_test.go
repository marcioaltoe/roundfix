// Boundary: real local Git merges; regeneration is a shell command and origin is local.
package cli

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func newLineConflictFixture(t *testing.T, sourceConflict bool, command string) conflictFixture {
	t.Helper()
	root := t.TempDir()
	repo, remote, home := filepath.Join(root, "item"), filepath.Join(root, "origin.git"), filepath.Join(root, "home")
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
	skill := func(version, body string) string {
		return "---\nversion: " + version + "\nname: example\ndescription: Example skill\nmetadata:\n  author: Example\n  version: " + version + "\n---\n\n# Skill\n\n" + body + "\n\nEnd\n"
	}
	write("SKILL.md", skill("1.0.0", "Base body"))
	write("derived.txt", "base\n")
	write(".roundfixrc.yml", "delivery:\n  derived_paths:\n    - paths: [derived.txt]\n      lines: {paths: [SKILL.md], match: '^ *version: '}\n      regenerate: |\n        "+command+"\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "base")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	// The resolver must override a user's three-sided conflict format.
	gittest.Run(t, repo, "config", "merge.conflictStyle", "diff3")
	gittest.Run(t, repo, "checkout", "-b", "feat/item")
	write("SKILL.md", skill("1.0.1", "Item body"))
	write("derived.txt", "base\nitem\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "item")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "main")
	body := "Base body"
	if sourceConflict {
		body = "Default body"
	}
	write("SKILL.md", skill("1.0.2", body))
	write("derived.txt", "base\ndefault\n")
	write("incoming.txt", "incoming\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "default")
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "checkout", "feat/item")
	cfg := roundconfig.Builtin()
	cfg.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	return conflictFixture{workflow: &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: cfg, HomeDir: home, GitRoot: repo, UserConfigPath: filepath.Join(home, ".roundfix", "config.yml")}}, repo: repo, head: head}
}

func updateLineFixtureDefault(t *testing.T, fixture conflictFixture, name, content string) {
	t.Helper()
	gittest.Run(t, fixture.repo, "checkout", "main")
	if err := os.WriteFile(filepath.Join(fixture.repo, name), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, fixture.repo, "add", ".")
	gittest.Run(t, fixture.repo, "commit", "-m", "default fixture adjustment")
	gittest.Run(t, fixture.repo, "push", "origin", "main")
	gittest.Run(t, fixture.repo, "checkout", "feat/item")
}

func TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated(t *testing.T) {
	t.Parallel()
	for _, recordConflict := range []bool{true, false} {
		t.Run(map[bool]string{true: "record and version conflicts", false: "only version conflicts"}[recordConflict], func(t *testing.T) {
			// Check the merge takes the default versions before the shell raises both fields.
			fixture := newLineConflictFixture(t, false, "test \"$(grep -c 'version: 1.0.2' SKILL.md)\" = 2 && sed 's/version: 1.0.2/version: 1.0.3/' SKILL.md > skill.tmp && mv skill.tmp SKILL.md && printf 'regenerated\\n' > derived.txt")
			if !recordConflict {
				updateLineFixtureDefault(t, fixture, "derived.txt", "base\n")
			}
			result := resolveFixture(t, fixture)
			if result.Head == "" || result.Head == fixture.head || len(result.SourcePaths) != 0 || len(result.Regenerated) != 1 {
				t.Fatalf("result=%+v", result)
			}
			content, err := os.ReadFile(filepath.Join(fixture.repo, "SKILL.md"))
			if err != nil || strings.Count(string(content), "version: 1.0.3") != 2 || !strings.Contains(string(content), "Item body") || strings.Contains(string(content), "<<<<<<<") {
				t.Fatalf("skill=%q err=%v", content, err)
			}
			if content, err := os.ReadFile(filepath.Join(fixture.repo, "derived.txt")); err != nil || string(content) != "regenerated\n" {
				t.Fatalf("record=%q err=%v", content, err)
			}
			if content, err := os.ReadFile(filepath.Join(fixture.repo, "incoming.txt")); err != nil || string(content) != "incoming\n" {
				t.Fatalf("incoming=%q err=%v", content, err)
			}
			if parents := strings.Fields(gittest.Run(t, fixture.repo, "show", "-s", "--format=%P", "HEAD")); len(parents) != 2 || parents[0] != fixture.head {
				t.Fatalf("parents=%v", parents)
			}
			if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
				t.Fatalf("dirty=%q", dirty)
			}
		})
	}
}

func assertLineConflictAborted(t *testing.T, fixture conflictFixture) {
	t.Helper()
	if head := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "HEAD")); head != fixture.head {
		t.Fatalf("head=%s", head)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%q", dirty)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, ".git", "MERGE_HEAD")); !os.IsNotExist(err) {
		t.Fatalf("merge still active: %v", err)
	}
	if content := gittest.Run(t, fixture.repo, "show", "HEAD:SKILL.md"); !strings.Contains(content, "Item body") || !strings.Contains(content, "version: 1.0.1") {
		t.Fatalf("candidate=%q", content)
	}
}

func TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge(t *testing.T) {
	t.Parallel()
	fixture := newLineConflictFixture(t, true, "touch forbidden.txt")
	result := resolveFixture(t, fixture)
	if !reflect.DeepEqual(result.SourcePaths, []string{"SKILL.md"}) || len(result.Regenerated) != 0 {
		t.Fatalf("result=%+v", result)
	}
	assertLineConflictAborted(t, fixture)
}

func TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, command string }{
		{"body", "sed 's/Item body/Unexpected body/' SKILL.md > skill.tmp; mv skill.tmp SKILL.md"},
		{"added matching line", "printf 'version: 9.0.0\\n' >> SKILL.md"},
		{"removed matching line", "sed '/^version:/d' SKILL.md > skill.tmp; mv skill.tmp SKILL.md"},
		{"no longer matching", "sed 's/version:/release:/' SKILL.md > skill.tmp; mv skill.tmp SKILL.md"},
		{"deleted file", "rm SKILL.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newLineConflictFixture(t, false, tc.command)
			result := resolveFixture(t, fixture)
			if !reflect.DeepEqual(result.SourcePaths, []string{"regenerated SKILL.md outside delivery.derived_paths"}) {
				t.Fatalf("result=%+v", result)
			}
			assertLineConflictAborted(t, fixture)
		})
	}
}

func TestWholeFileDerivedPathsTakePrecedenceOverDeclaredLines(t *testing.T) {
	t.Parallel()
	fixture := newLineConflictFixture(t, true, "printf 'regenerated\\n' > derived.txt")
	// Whole-file precedence also applies across separate declarations.
	updateLineFixtureDefault(t, fixture, ".roundfixrc.yml", "delivery:\n  derived_paths:\n    - paths: [SKILL.md, derived.txt]\n      regenerate: printf 'regenerated\\n' > derived.txt\n    - paths: [unused.txt]\n      lines: {paths: [SKILL.md], match: '^ *version: '}\n      regenerate: touch forbidden.txt\n")
	result := resolveFixture(t, fixture)
	if result.Head == "" || len(result.SourcePaths) != 0 || len(result.Regenerated) != 1 {
		t.Fatalf("result=%+v", result)
	}
	if content, err := os.ReadFile(filepath.Join(fixture.repo, "SKILL.md")); err != nil || !strings.Contains(string(content), "Default body") {
		t.Fatalf("skill=%q err=%v", content, err)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, "forbidden.txt")); !os.IsNotExist(err) {
		t.Fatalf("line command ran despite precedence: %v", err)
	}
}
