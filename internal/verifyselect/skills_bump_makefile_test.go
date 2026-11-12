// Suite: Catalog Skill Bump Makefile contract.
// Invariant: one catalog commit supplies every vendored skill through the existing writers.
// Boundary IN: the real Makefile, a local Git catalog, and recording Go and Roundfix stubs.
// Boundary OUT: the network, real skill restoration, and production regeneration.
package verifyselect_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/testfixture"
)

func TestTheSkillsBumpPinsTheCatalogMainCommit(t *testing.T) {
	t.Parallel()
	root, catalog, revision := newSkillsBumpFixture(t)
	command := skillsBumpCommand(t, root, catalog)
	var diagnostics bytes.Buffer
	command.Stderr = &diagnostics
	output, err := command.Output()
	if err != nil {
		t.Fatalf("bump: %v\n%s\n%s", err, output, &diagnostics)
	}
	wantFinal := "skills-third-party-bump: 2 catalog skills and every upstream Setup Snapshot pinned to " + revision
	if string(output) != wantFinal+"\n" {
		t.Fatalf("missing final line %q:\n%s", wantFinal, output)
	}
	calls := readSkillsBumpCalls(t, root)
	lines := strings.Split(strings.TrimSpace(calls), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "go build -o ") {
		t.Fatalf("first call must build Roundfix:\n%s", calls)
	}
	fields := strings.Fields(lines[0])
	checkout := filepath.Join(filepath.Dir(fields[3]), "catalog")
	build := "go build -o " + fields[3] + " ./cmd/roundfix"
	restore := "roundfix baseline skills restore --profile fixture-profile --skill alpha --skill zeta --source-dir " + checkout + " --repo . --format text"
	want := []string{
		build,
		"roundfix baseline assets sync --source-dir " + checkout + "/setups --format text",
		"roundfix baseline assets sync --source-dir " + checkout + "/setups --check --format text",
		build,
		restore,
		restore + " --confirm-plan sha256:fixture-plan",
		"go test probe -run Probe -update -count=1",
		"go test ./internal/baseline -run TestCatalogCompatibility -count=1",
		build,
		"roundfix baseline update --repo . --no-skills --no-history --yes --format text",
		"roundfix baseline update --repo . --no-skills --no-history --yes --format text",
	}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("calls =\n%s\nwant =\n%s", calls, strings.Join(want, "\n"))
	}
	owned, err := os.ReadFile(filepath.Join(root, "repo/.agents/skills/roundfix/SKILL.md"))
	if err != nil || string(owned) != "owned\n" {
		t.Fatalf("owned skill changed: %q, %v", owned, err)
	}
	retired := skillsBumpCommand(t, root, catalog)
	for index, argument := range retired.Args {
		if argument == "skills-third-party-bump" {
			retired.Args[index] = "skills-update"
		}
	}
	retiredOutput, err := retired.CombinedOutput()
	if err == nil || !strings.Contains(string(retiredOutput), "No rule to make target") {
		t.Fatalf("retired skills-update target: %v\n%s", err, retiredOutput)
	}
	if got := readSkillsBumpCalls(t, root); got != calls {
		t.Fatalf("retired target ran a stub:\n%s", got)
	}
	assertSkillsBumpTemporaryDirectoryRemoved(t, root)
}

func TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests(t *testing.T) {
	t.Parallel()
	root, catalog, _ := newSkillsBumpFixture(t)
	repo := filepath.Join(root, "repo")
	writeFile(t, repo, "embedded-asset", "before digests\n")
	output, err := skillsBumpCommand(t, root, catalog).CombinedOutput()
	if err != nil {
		t.Fatalf("bump: %v\n%s", err, output)
	}
	refreshed, err := os.ReadFile(filepath.Join(repo, "refreshed-from"))
	if err != nil {
		t.Fatal(err)
	}
	if string(refreshed) != "regenerated\n" {
		t.Fatalf("the managed refresh ran a binary that embeds %q; want the asset make baseline-digests regenerated, %q", refreshed, "regenerated\n")
	}
	assertSkillsBumpTemporaryDirectoryRemoved(t, root)
}

func TestTheSkillsBumpRefusesWhatItMustNotPin(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct {
		name        string
		want        string
		beforeBuild bool
	}{
		{name: "owned skill in lock", want: "Roundfix-owned skill roundfix", beforeBuild: true},
		{name: "branch ref", want: "full 40-character", beforeBuild: true},
		{name: "unserved commit", want: "does not serve commit ", beforeBuild: true},
		{name: "folder without lock", want: "differ: untracked-skill", beforeBuild: true},
		{name: "crafted ref", want: "full 40-character", beforeBuild: true},
		{name: "crafted catalog", want: "full 40-character", beforeBuild: true},
		{name: "option catalog", want: "must not start with -", beforeBuild: true},
		{name: "skill missing from catalog", want: "has no single skills/*/ghost tree", beforeBuild: true},
		{name: "differing folder", want: ".agents/skills/alpha differs from the catalog at "},
		{name: "entry at another ref", want: "not pinned to "},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			root, catalog, revision := newSkillsBumpFixture(t)
			command := skillsBumpCommand(t, root, catalog)
			switch scenario.name {
			case "owned skill in lock":
				writeFile(t, root, "repo/skills-lock.json", `{"skills":{"roundfix":{"source":"marcioaltoe/skills"}}}`)
			case "branch ref":
				command.Args = append(command.Args, "SKILLS_REF=main")
			case "unserved commit":
				command.Args = append(command.Args, "SKILLS_REF="+strings.Repeat("0", 40))
			case "crafted ref":
				command.Args = append(command.Args, "SKILLS_REF='; : > injected; '")
			case "crafted catalog":
				command.Args = append(command.Args, `SKILLS_CATALOG=x"; : > injected; "`)
			case "option catalog":
				command.Args = append(command.Args, "SKILLS_CATALOG=--upload-pack=sh")
			case "skill missing from catalog":
				command.Args = append(command.Args, "SKILLS_REF="+revision)
				writeFile(t, root, "repo/skills-lock.json", `{"skills":{"zeta":{"source":"marcioaltoe/skills","ref":"old"},"alpha":{"source":"marcioaltoe/skills","ref":"old"},"ghost":{"source":"marcioaltoe/skills","ref":"old"}}}`)
				writeFile(t, root, "repo/.agents/skills/ghost/SKILL.md", "ghost\n")
			case "folder without lock":
				writeFile(t, root, "repo/.agents/skills/untracked-skill/SKILL.md", "extra\n")
			case "differing folder":
				writeFile(t, root, "repo/leave-folder", "")
			case "entry at another ref":
				writeFile(t, root, "repo/leave-ref", "")
			}
			output, err := command.CombinedOutput()
			if err == nil {
				t.Fatalf("bump accepted %s:\n%s", scenario.name, output)
			}
			if _, statErr := os.Stat(filepath.Join(root, "repo", "injected")); !os.IsNotExist(statErr) {
				t.Fatalf("%s ran shell source from a make variable: %v", scenario.name, statErr)
			}
			want := scenario.want
			if scenario.name == "unserved commit" {
				want += strings.Repeat("0", 40)
			} else if !scenario.beforeBuild {
				want += revision
			}
			if !strings.Contains(string(output), "skills-third-party-bump: ") || !strings.Contains(string(output), want) {
				t.Fatalf("refusal missing %q:\n%s", want, output)
			}
			calls := readSkillsBumpCalls(t, root)
			if scenario.beforeBuild && calls != "" {
				t.Fatalf("preflight ran a stub:\n%s", calls)
			}
			if strings.Contains(calls, "go test ") || strings.Contains(calls, "baseline update ") {
				t.Fatalf("refusal reached regeneration:\n%s", calls)
			}
			assertSkillsBumpTemporaryDirectoryRemoved(t, root)
		})
	}
}

func TestTheSkillsBumpAcceptsAnAlreadyCurrentCatalog(t *testing.T) {
	t.Parallel()
	root, catalog, revision := newSkillsBumpFixture(t)
	repo := filepath.Join(root, "repo")
	writeFile(t, repo, ".agents/skills/alpha/SKILL.md", "catalog alpha\n")
	writeFile(t, repo, ".agents/skills/zeta/SKILL.md", "catalog zeta\n")
	writeFile(t, repo, "skills-lock.json", `{"skills":{"zeta":{"source":"marcioaltoe/skills","ref":"`+revision+`"},"alpha":{"source":"marcioaltoe/skills","ref":"`+revision+`"},"unrelated":{"source":"other/catalog","ref":"keep"}}}`)
	writeFile(t, repo, "already-current", "")
	command := skillsBumpCommand(t, root, catalog)
	var diagnostics bytes.Buffer
	command.Stderr = &diagnostics
	output, err := command.Output()
	if err != nil {
		t.Fatalf("already-current bump: %v\n%s\n%s", err, output, &diagnostics)
	}
	wantFinal := "skills-third-party-bump: 2 catalog skills and every upstream Setup Snapshot pinned to " + revision
	if string(output) != wantFinal+"\n" {
		t.Fatalf("missing final line %q:\n%s", wantFinal, output)
	}
	if !strings.Contains(diagnostics.String(), "vendored catalog skills already match "+revision) {
		t.Fatalf("missing already-current diagnostic:\n%s", &diagnostics)
	}
	if calls := readSkillsBumpCalls(t, root); strings.Contains(calls, "--confirm-plan") {
		t.Fatalf("an already-current catalog confirmed a restore plan:\n%s", calls)
	}
	assertSkillsBumpTemporaryDirectoryRemoved(t, root)
}

func TestTheSkillsBumpRefusesASilentZeroExitPreview(t *testing.T) {
	t.Parallel()
	root, catalog, _ := newSkillsBumpFixture(t)
	writeFile(t, filepath.Join(root, "repo"), "preview-without-changes-line", "")
	output, err := skillsBumpCommand(t, root, catalog).CombinedOutput()
	if err == nil || !strings.Contains(string(output), "restore preview exited 0; expected 3 with a Plan Digest, or 0 with no changes") {
		t.Fatalf("silent zero-exit preview: %v\n%s", err, output)
	}
	if calls := readSkillsBumpCalls(t, root); strings.Contains(calls, "--confirm-plan") {
		t.Fatalf("a refused preview confirmed a restore plan:\n%s", calls)
	}
}

func newSkillsBumpFixture(t *testing.T) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	for _, directory := range []string{"bin", "tmp", "home", "empty", "catalog", "repo"} {
		if err := os.MkdirAll(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// PATH itself is fixture-owned; only the tools this Makefile invokes are exposed.
	for _, name := range []string{"make", "git", "sh", "cat", "cp", "chmod", "dirname", "mkdir", "mktemp", "rm", "mv", "find", "sort", "comm", "awk", "sed", "diff", "jq", "shasum", "id", "date"} {
		tool := testfixture.RequireTool(t, name)
		if err := os.Symlink(tool, filepath.Join(root, "bin", name)); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, filepath.Join(root, "repo"), "init", "-q", "--template=")
	catalog := filepath.Join(root, "catalog")
	runGit(t, catalog, "init", "-q", "--template=", "-b", "main")
	writeFile(t, catalog, "skills/group/alpha/SKILL.md", "catalog alpha\n")
	writeFile(t, catalog, "skills/other/zeta/SKILL.md", "catalog zeta\n")
	writeFile(t, catalog, "setups/go.txt", "setup\n")
	runGit(t, catalog, "add", ".")
	runGit(t, catalog, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null", "commit", "-qm", "catalog")
	revision := strings.TrimSpace(runGit(t, catalog, "rev-parse", "HEAD"))
	writeFile(t, root, "repo/skills-lock.json", `{"skills":{"zeta":{"source":"marcioaltoe/skills","ref":"old"},"alpha":{"source":"marcioaltoe/skills","ref":"old"},"unrelated":{"source":"other/catalog","ref":"keep"}}}`)
	writeFile(t, root, "repo/docs/agents/setup-context.json", `{"profile":"fixture-profile"}`)
	writeFile(t, root, "repo/.agents/skills/roundfix/SKILL.md", "owned\n")
	writeFile(t, root, "repo/.agents/skills/alpha/SKILL.md", "stale\n")
	writeFile(t, root, "repo/.agents/skills/zeta/SKILL.md", "stale\n")
	writeFile(t, root, "calls", "")
	prefix := "#!/bin/sh\nset -eu\nfixture='" + strings.ReplaceAll(root, "'", "'\"'\"'") + "'\n"
	writeFile(t, root, "bin/go", prefix+`printf 'go %s\n' "$*" >> "$fixture/calls"
if [ "$1" = build ]; then
    test "$2" = -o
    cp "$fixture/roundfix-stub" "$3"
    if [ -f embedded-asset ]; then
        cp embedded-asset "$3.embedded"
    else
        printf 'absent\n' > "$3.embedded"
    fi
    chmod +x "$3"
elif [ "$1" = test ] && [ "$2" = probe ]; then
    printf 'regenerated\n' > embedded-asset
fi
`)
	writeFile(t, root, "roundfix-stub", prefix+`printf 'roundfix %s\n' "$*" >> "$fixture/calls"
case "$HOME" in "$fixture/tmp/"*/home) ;; *) echo "unsafe HOME" >&2; exit 1;; esac
test -d "$HOME"
test -z "$(find "$HOME" -mindepth 1 -print)"
case "$*" in
    'baseline assets sync '*)
        test -f "$5/go.txt"
        test "$(git -C "$5/.." remote get-url origin)" = https://github.com/marcioaltoe/skills.git
        ;;
    'baseline skills restore '*)
        source_dir=
        confirmed=false
        while [ "$#" -gt 0 ]; do
            case "$1" in
                --source-dir) source_dir=$2; shift ;;
                --confirm-plan) test "$2" = sha256:fixture-plan; confirmed=true; shift ;;
            esac
            shift
        done
        if [ "$confirmed" = false ] && [ -f already-current ]; then
            printf 'Baseline skills restore: no changes\n'
            exit 0
        fi
        if [ "$confirmed" = false ] && [ -f preview-without-changes-line ]; then
            exit 0
        fi
        if [ "$confirmed" = false ]; then
            printf 'Plan Digest: sha256:fixture-plan\n'
            exit 3
        fi
        revision=$(git -C "$source_dir" rev-parse HEAD)
        if [ ! -f leave-ref ]; then
            jq --arg revision "$revision" '.skills |= with_entries(if .value.source == "marcioaltoe/skills" then .value.ref = $revision else . end)' skills-lock.json > lock.new
            mv lock.new skills-lock.json
        fi
        for name in alpha zeta; do
            if [ "$name" = alpha ] && [ -f leave-folder ]; then continue; fi
            cp "$source_dir/skills/"*/"$name/SKILL.md" ".agents/skills/$name/SKILL.md"
        done
        ;;
    'baseline update '*)
        cp "$0.embedded" refreshed-from
        printf 'File changes: 0\n'
        ;;
    *) exit 1 ;;
esac
`)
	if err := os.Chmod(filepath.Join(root, "bin/go"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root, catalog, revision
}

func skillsBumpCommand(t *testing.T, root, catalog string) *exec.Cmd {
	t.Helper()
	command := exec.CommandContext(t.Context(), filepath.Join(root, "bin/make"), "--no-print-directory", "-f", filepath.Join(findRepositoryRoot(t), "Makefile"), "skills-third-party-bump", "SKILLS_CATALOG="+catalog, "GO="+filepath.Join(root, "bin/go"), "DERIVED_DIGEST_PATHS="+filepath.Join(root, "empty"), "BASELINE_DIGEST_STEPS=probe:Probe")
	command.Dir = filepath.Join(root, "repo")
	command.Env = []string{"PATH=" + filepath.Join(root, "bin"), "TMPDIR=" + filepath.Join(root, "tmp"), "HOME=" + filepath.Join(root, "home"), "GIT_ALLOW_PROTOCOL=file", "GIT_CONFIG_NOSYSTEM=1"}
	return command
}

func readSkillsBumpCalls(t *testing.T, root string) string {
	t.Helper()
	calls, err := os.ReadFile(filepath.Join(root, "calls"))
	if err != nil {
		t.Fatal(err)
	}
	return string(calls)
}

func assertSkillsBumpTemporaryDirectoryRemoved(t *testing.T, root string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, "tmp/rf-skills.*"))
	if err != nil || len(paths) != 0 {
		t.Fatalf("temporary bump directories remain: %v, %v", paths, err)
	}
}
