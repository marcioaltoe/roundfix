// Boundary: real local Git merges with the shipped skill layout and shell regeneration.
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

var realSkillConflictPaths = []string{".agents/skills/example/SKILL.md", "skills/example/SKILL.md"}

func newRealSkillConflictFixture(t *testing.T, authorConflict bool) (conflictFixture, string) {
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
		path := filepath.Join(repo, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	skill := func(version, author, itemBody, mainBody string) string {
		return "---\nname: example\ndescription: Example skill\nmetadata:\n  category: workflow\n  version: " + version + "\n  author: " + author + "\n  source: https://example.com/skills\nversion: " + version + "\n---\n\n# Example\n\n" + itemBody + "\n\nShared body\n\nMore shared body\n\n" + mainBody + "\n"
	}
	writeSkills := func(version, author, itemBody, mainBody string) {
		t.Helper()
		for _, name := range realSkillConflictPaths {
			write(name, skill(version, author, itemBody, mainBody))
		}
	}
	// Outside the worktree, so even an aborted merge cannot hide a command invocation.
	invocations := filepath.Join(root, "regenerations")
	command := "printf 'ran\\n' >> '" + invocations + "'; for skill in .agents/skills/example/SKILL.md skills/example/SKILL.md; do test \"$(grep -c '^ *version: 0.1.32$' \"$skill\")\" = 2 || exit 1; done; for skill in .agents/skills/example/SKILL.md skills/example/SKILL.md; do sed 's/version: 0.1.32/version: 0.1.33/' \"$skill\" > \"$skill.tmp\" && mv \"$skill.tmp\" \"$skill\" || exit 1; done; printf 'regenerated\\n' > skills/testdata/owned-skill-versions.json"
	write(".roundfixrc.yml", "delivery:\n  derived_paths:\n    - paths: [skills/testdata/owned-skill-versions.json]\n      lines: {paths: ['.agents/skills/*/SKILL.md', 'skills/*/SKILL.md'], match: '^ *version: '}\n      regenerate: |\n        "+command+"\n")
	writeSkills("0.1.30", "Example", "Base item body", "Base main body")
	write("skills/testdata/owned-skill-versions.json", "base\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "base")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, repo, "config", "merge.conflictStyle", "diff3")
	gittest.Run(t, repo, "checkout", "-b", "feat/item")
	writeSkills("0.1.31", "Example", "Item body edit", "Base main body")
	write("skills/testdata/owned-skill-versions.json", "item\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "item")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "main")
	author := "Example"
	if authorConflict {
		author = "Default author"
	}
	writeSkills("0.1.32", author, "Base item body", "Default body edit")
	write("skills/testdata/owned-skill-versions.json", "default\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "default")
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "checkout", "feat/item")
	cfg := roundconfig.Builtin()
	cfg.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	return conflictFixture{workflow: &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: cfg, HomeDir: home, GitRoot: repo, UserConfigPath: filepath.Join(home, ".roundfix", "config.yml")}}, repo: repo, head: head}, invocations
}

func TestARealSkillVersionHunkIsMergedAndRegeneratedOnePatchAboveMain(t *testing.T) {
	t.Parallel()
	fixture, invocations := newRealSkillConflictFixture(t, false)
	// Prove Git puts both version fields and the shared source line in one hunk.
	if _, err := fixture.workflow.git.RunGit(t.Context(), fixture.repo, "-c", "merge.conflictStyle=merge", "merge", "--no-ff", "--no-commit", "origin/main"); err == nil {
		t.Fatal("fixture did not conflict")
	}
	content, err := os.ReadFile(filepath.Join(fixture.repo, realSkillConflictPaths[0]))
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	if strings.Count(text, "<<<<<<<") != 1 || strings.Count(text, "=======") != 1 || strings.Count(text, ">>>>>>>") != 1 {
		t.Fatalf("expected exactly one hunk: %s", text)
	}
	hunk := text[strings.Index(text, "<<<<<<<"):strings.Index(text, ">>>>>>>")]
	for _, side := range strings.Split(hunk, "=======") {
		if strings.Count(side, "version: ") != 2 || !strings.Contains(side, "  source: https://example.com/skills\nversion: ") || !strings.Contains(side, "  author: Example\n") {
			t.Fatalf("expected both versions around author/source: %s", side)
		}
	}
	gittest.Run(t, fixture.repo, "merge", "--abort")
	result := resolveFixture(t, fixture)
	if result.Head == "" || result.Head == fixture.head || len(result.SourcePaths) != 0 || len(result.Regenerated) != 1 {
		t.Fatalf("result=%+v", result)
	}
	for _, name := range realSkillConflictPaths {
		content, err := os.ReadFile(filepath.Join(fixture.repo, name))
		if err != nil || strings.Count(string(content), "version: 0.1.33") != 2 || !strings.Contains(string(content), "Item body edit") || !strings.Contains(string(content), "Default body edit") || strings.Contains(string(content), "<<<<<<<") {
			t.Fatalf("%s=%q err=%v", name, content, err)
		}
	}
	if content, err := os.ReadFile(invocations); err != nil || string(content) != "ran\n" {
		t.Fatalf("regenerations=%q err=%v", content, err)
	}
	if content, err := os.ReadFile(filepath.Join(fixture.repo, "skills/testdata/owned-skill-versions.json")); err != nil || string(content) != "regenerated\n" {
		t.Fatalf("record=%q err=%v", content, err)
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

func TestAVersionHunkWhoseInterveningLineDiffersAbortsTheMerge(t *testing.T) {
	t.Parallel()
	fixture, invocations := newRealSkillConflictFixture(t, true)
	result := resolveFixture(t, fixture)
	if !reflect.DeepEqual(result.SourcePaths, realSkillConflictPaths) || len(result.Regenerated) != 0 {
		t.Fatalf("result=%+v", result)
	}
	if _, err := os.Stat(invocations); !os.IsNotExist(err) {
		t.Fatalf("regeneration ran: %v", err)
	}
	if head := strings.TrimSpace(gittest.Run(t, fixture.repo, "rev-parse", "HEAD")); head != fixture.head {
		t.Fatalf("head=%s want=%s", head, fixture.head)
	}
	if dirty := gittest.Run(t, fixture.repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%q", dirty)
	}
	if _, err := os.Stat(filepath.Join(fixture.repo, ".git", "MERGE_HEAD")); !os.IsNotExist(err) {
		t.Fatalf("merge still active: %v", err)
	}
}
