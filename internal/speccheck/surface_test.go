// Suite: undocumented CLI surface authoring checks.
// Invariant: every pending non-QA Task that names a CLI surface also names a guide in its dependency ancestry.
// Boundary IN: public Spec Consistency Check, parsed Task Context, and Task Graph dependencies.
// Boundary OUT: guide content quality and Daemon Verification execution.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestCLISurfaceWithoutAGuideIsReported(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- interface: `internal/cli/flags.go`",
	}})

	result := fixture.check()
	findings := findingsWithCode(result, speccheck.CodeCLIUndocumented)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeCLIUndocumented, findings)
	}
	if findings[0].Severity != speccheck.SeverityGap {
		t.Errorf("severity = %q, want %q", findings[0].Severity, speccheck.SeverityGap)
	}
	if len(findings[0].Where) != 1 || findings[0].Where[0].Path != "docs/specs/cli-surface/task_01.md" || findings[0].Where[0].Line != 12 {
		t.Errorf("Where = %#v, want CLI surface declaration at task_01.md:12", findings[0].Where)
	}
}

func TestStrictPromotesAnUndocumentedCLISurface(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- creates: `cmd/roundfix/new_command.go`",
	}})

	result := fixture.check()
	speccheck.PromoteGaps(&result)
	findings := findingsWithCode(result, speccheck.CodeCLIUndocumented)
	if len(findings) != 1 || findings[0].Severity != speccheck.SeverityError {
		t.Fatalf("strict %s findings = %#v, want one error", speccheck.CodeCLIUndocumented, findings)
	}
}

func TestCLISurfaceNamingItsGuidePasses(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- interface: `internal/cli/flags.go`\n- interface: `.agents/skills/roundfix/SKILL.md`",
	}})

	fixture.requireFindingCount(0)
}

func TestCLISurfaceWhoseDependencyNamesTheGuidePasses(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{
		{id: "task_01", context: "- creates: `docs/user-guide/commands.md`"},
		{id: "task_02", needs: []string{"task_01"}, context: "- interface: `internal/example.go`"},
		{id: "task_03", needs: []string{"task_02"}, context: "- interface: `internal/cli/flags.go`"},
	})

	fixture.requireFindingCount(0)
}

func TestGuideOnlyInADependentTaskDoesNotCount(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{
		{id: "task_01", context: "- interface: `internal/cli/flags.go`"},
		{id: "task_02", needs: []string{"task_01"}, context: "- creates: `docs/agents/cli.md`"},
	})

	result := fixture.check()
	findings := findingsWithCode(result, speccheck.CodeCLIUndocumented)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "task_01.md") {
		t.Fatalf("%s findings = %#v, want only task_01", speccheck.CodeCLIUndocumented, findings)
	}
}

func TestCLIContractTestIsASurface(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- interface: `internal/cli/cli_test.go`",
	}})

	fixture.requireFindingCount(1)
}

func TestOrdinaryCLITestIsNotASurface(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- interface: `internal/cli/flags_test.go`",
	}})

	fixture.requireFindingCount(0)
}

func TestInstructionGuideDoesNotDocumentCLISurface(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{
		id:      "task_01",
		context: "- interface: `internal/cli/flags.go`\n- instruction: `.agents/skills/roundfix/SKILL.md`",
	}})

	fixture.requireFindingCount(1)
}

func TestMissingTaskGraphListsTheCLISurfaceSkip(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{id: "task_01", context: "- interface: `internal/cli/flags.go`"}})
	if err := os.Remove(filepath.Join(fixture.repoRoot, "docs", "specs", fixture.slug, "_tasks.md")); err != nil {
		t.Fatalf("remove fixture Task Graph: %v", err)
	}

	result := fixture.check()
	if !hasSkip(result, speccheck.CodeCLIUndocumented, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeCLIUndocumented)
	}
}

func TestMissingPRDListsTheCLISurfaceSkip(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{id: "task_01", context: "- interface: `internal/cli/flags.go`"}})
	if err := os.Remove(filepath.Join(fixture.repoRoot, "docs", "specs", fixture.slug, "_prd.md")); err != nil {
		t.Fatalf("remove fixture PRD: %v", err)
	}

	result := fixture.check()
	if !hasSkip(result, speccheck.CodeCLIUndocumented, "_prd.md") {
		t.Fatalf("Skipped = %#v, want %s missing _prd.md", result.Skipped, speccheck.CodeCLIUndocumented)
	}
}

type cliSurfaceTask struct {
	id      string
	needs   []string
	context string
}

type cliSurfaceFixture struct {
	t         *testing.T
	repoRoot  string
	specsRoot string
	slug      string
}

func newCLISurfaceFixture(t *testing.T) *cliSurfaceFixture {
	t.Helper()
	const slug = "cli-surface"
	repoRoot := t.TempDir()
	fixture := &cliSurfaceFixture{
		t:         t,
		repoRoot:  repoRoot,
		specsRoot: filepath.Join(repoRoot, "docs", "specs"),
		slug:      slug,
	}
	writeCitationFixtureFile(t, repoRoot, "docs/agents/agent-instructions.md", "# Agent instructions\n")
	constraints := "- Identifier strategy: not applicable — no identifier change. Source: `docs/agents/agent-instructions.md`.\n" +
		"- Authentication and HTTP: not applicable — local files only. Source: `docs/agents/agent-instructions.md`.\n" +
		"- Active ADR obligations: not applicable — no ADR applies. Source: `docs/agents/agent-instructions.md`.\n" +
		"- Tooling authority: not applicable — fixture edits ordinary files only. Source: `docs/agents/agent-instructions.md`.\n"
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", "---\nstatus: active\n---\n\n# Fixture\n\n## Project Constraints\n\n"+constraints+"\n## Success Metrics\n\nNone. This fixture has no product metric.\n")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_techspec.md", "# Fixture TechSpec\n\n## Project Constraints\n\n"+constraints+"\n## API Contracts\n\nNone. This fixture has no public API.\n")
	return fixture
}

func (fixture *cliSurfaceFixture) writeTasks(tasks []cliSurfaceTask) {
	fixture.t.Helper()
	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\nspec: cli-surface\ngraph:\n  nodes:\n")
	for _, task := range tasks {
		fmt.Fprintf(&manifest, "    - id: %s\n      file: %s.md\n      needs: [%s]\n", task.id, task.id, strings.Join(task.needs, ", "))
		writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/"+task.id+".md", fmt.Sprintf(`---
task: %s
spec: cli-surface
status: pending
type: backend
---

# Task: Fixture

## Context

%s

## Verification

- %s
`, task.id, task.context, "`true`"))
	}
	manifest.WriteString("---\n\n# Task Graph\n")
	writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_tasks.md", manifest.String())
}

func (fixture *cliSurfaceFixture) check() speccheck.Result {
	fixture.t.Helper()
	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		fixture.t.Fatalf("Check(): %v", err)
	}
	return result
}

func (fixture *cliSurfaceFixture) requireFindingCount(want int) {
	fixture.t.Helper()
	result := fixture.check()
	if findings := findingsWithCode(result, speccheck.CodeCLIUndocumented); len(findings) != want {
		fixture.t.Fatalf("%s findings = %#v, want %d", speccheck.CodeCLIUndocumented, findings, want)
	}
}
