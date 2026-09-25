// Suite: undeclared Governed Path authoring checks.
// Invariant: every pending non-QA Task declaration of a Governed Path is backed by the grant and each present tooling row.
// Boundary IN: public Spec Consistency Check, Task Graph parsing, authorization records, and repository path governance.
// Boundary OUT: Daemon changed-path auditing and authored Verification execution.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestUndeclaredGovernedContextPathIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- interface: `Makefile`", "test -f Makefile")
	fixture.writeGrant(nil, nil)
	fixture.writeArtifacts(nil, nil)

	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		t.Fatalf("Check(): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeToolingUndeclared)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeToolingUndeclared, findings)
	}
	if !hasLocation(findings[0], "docs/specs/undeclared-governed/task_01.md") ||
		!hasLocation(findings[0], "docs/specs/undeclared-governed/_authorization.md") {
		t.Fatalf("finding locations = %#v, want Task and authorization record", findings[0].Where)
	}
}

func TestUndeclaredGovernedCreatesPathIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- creates: `.roundfixrc.yml`", "true")
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(1)
}

func TestUndeclaredGovernedVerificationPathIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "", "test -f Makefile")
	fixture.writeGrant([]string{"internal/example.go"}, nil)
	fixture.writeArtifacts([]string{"internal/example.go"}, []string{"internal/example.go"})
	fixture.requireFindingCount(1)
}

func TestDeclaredGovernedPathsPass(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- interface: `Makefile`", "test -f Makefile")
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(0)
}

func TestGovernedPathMissingFromABoundedFilesRowIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- interface: `Makefile`", "true")
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"internal/example.go"})
	fixture.requireFindingCount(1)
}

func TestBoundedFilesRowPathTheRecordDoesNotGrantIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- interface: `Makefile`", "true")
	fixture.writeGrant([]string{"internal/example.go"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(1)
}

func TestGovernedPathWithoutAGrantIsRefused(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- interface: `Makefile`", "true")
	fixture.writeGrant(nil, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(1)
}

func TestInstructionContextPathIsNotAudited(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- instruction: `Makefile`", "true")
	fixture.writeGrant([]string{"internal/example.go"}, nil)
	fixture.writeArtifacts([]string{"internal/example.go"}, []string{"internal/example.go"})
	fixture.requireFindingCount(0)
}

func TestCompletedTaskIsNotAudited(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("completed", "backend", "- interface: `Makefile`", "true")
	fixture.writeGrant([]string{"internal/example.go"}, nil)
	fixture.writeArtifacts([]string{"internal/example.go"}, []string{"internal/example.go"})
	fixture.requireFindingCount(0)
}

func TestQATaskIsNotAudited(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	writeCitationFixtureFile(t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: undeclared-governed
qa: task_01
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)
	fixture.writeTask("pending", "qa", "- interface: `Makefile`", "true")
	fixture.writeGrant([]string{"internal/example.go"}, nil)
	fixture.writeArtifacts([]string{"internal/example.go"}, []string{"internal/example.go"})
	fixture.requireFindingCount(0)
}

func TestOrdinaryPathNeedsNoDeclaration(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- creates: `internal/example.go`", "true")
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(0)
}

func TestSanctionedRegenerationOutputCountsAsDeclared(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- creates: `skills/qa-gate/SKILL.md`", "test -f skills/qa-gate/SKILL.md")
	fixture.writeGrant([]string{"Makefile"}, []string{"skills/qa-gate/SKILL.md"})
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(0)
}

func TestMissingTaskGraphListsTheUndeclaredGovernedPathSkip(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	if err := os.Remove(filepath.Join(fixture.repoRoot, "docs", "specs", fixture.slug, "_tasks.md")); err != nil {
		t.Fatalf("remove fixture Task Graph: %v", err)
	}
	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		t.Fatalf("Check(): %v", err)
	}
	if !hasSkip(result, speccheck.CodeToolingUndeclared, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeToolingUndeclared)
	}
}

func TestMissingPRDListsTheUndeclaredGovernedPathSkip(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	if err := os.Remove(filepath.Join(fixture.repoRoot, "docs", "specs", fixture.slug, "_prd.md")); err != nil {
		t.Fatalf("remove fixture PRD: %v", err)
	}
	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		t.Fatalf("Check(): %v", err)
	}
	if !hasSkip(result, speccheck.CodeToolingUndeclared, "_prd.md") {
		t.Fatalf("Skipped = %#v, want %s missing _prd.md", result.Skipped, speccheck.CodeToolingUndeclared)
	}
}

func TestUndeclaredGovernedPathIsSkippedBeforeTheTaskStage(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	result, err := speccheck.CheckStage(fixture.specsRoot, fixture.repoRoot, fixture.slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	if !hasSkip(result, speccheck.CodeToolingUndeclared, "stage prd") {
		t.Fatalf("Skipped = %#v, want %s skipped before Task stage", result.Skipped, speccheck.CodeToolingUndeclared)
	}
}

type undeclaredFixture struct {
	t         *testing.T
	repoRoot  string
	specsRoot string
	slug      string
}

func newUndeclaredFixture(t *testing.T) *undeclaredFixture {
	t.Helper()
	const slug = "undeclared-governed"
	repoRoot := t.TempDir()
	fixture := &undeclaredFixture{
		t:         t,
		repoRoot:  repoRoot,
		specsRoot: filepath.Join(repoRoot, "docs", "specs"),
		slug:      slug,
	}
	writeCitationFixtureFile(t, repoRoot, "docs/agents/agent-instructions.md", "# Agent instructions\n")
	writeCitationFixtureFile(t, repoRoot, "Makefile", "all:\n\t@true\n")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: undeclared-governed
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---

# Task Graph
`)
	return fixture
}

func (fixture *undeclaredFixture) writeTask(status, taskType, context, verification string) {
	fixture.t.Helper()
	writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/task_01.md", fmt.Sprintf(`---
task: task_01
spec: undeclared-governed
status: %s
type: %s
---

# Task 01: Fixture

## Context

%s

## Verification

- %s
`, status, taskType, context, "`"+verification+"`"))
}

func (fixture *undeclaredFixture) writeGrant(paths, outputs []string) {
	fixture.t.Helper()
	var record strings.Builder
	record.WriteString("---\nstatus: approved\ngranted: 2026-09-25\naction: test undeclared paths\nconsuming: undeclared-governed\npaths:\n")
	for _, path := range paths {
		fmt.Fprintf(&record, "  - %s\n", path)
	}
	record.WriteString("operations:\n  - implement\n---\n")
	if outputs != nil {
		record.WriteString("\n## Sanctioned regeneration\n\n```yaml\ncommand: make skills-sync\noutputs:\n")
		for _, output := range outputs {
			fmt.Fprintf(&record, "  - %s\n", output)
		}
		record.WriteString("```\n")
	}
	writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_authorization.md", record.String())
}

func (fixture *undeclaredFixture) writeArtifacts(prdPaths, techSpecPaths []string) {
	fixture.t.Helper()
	row := func(paths []string) string {
		var bounded strings.Builder
		for index, path := range paths {
			if index > 0 {
				bounded.WriteString(", ")
			}
			fmt.Fprintf(&bounded, "`%s`", path)
		}
		return "Tooling authority: applicable — express maintainer authorization recorded in [_authorization.md](_authorization.md); bounded files: " + bounded.String() + ". Sanctioned regeneration: `make skills-sync`. Source: `docs/agents/agent-instructions.md`."
	}
	constraints := func(paths []string) string {
		return "- Identifier strategy: not applicable — no identifier change. Source: `docs/agents/agent-instructions.md`.\n" +
			"- Authentication and HTTP: not applicable — local files only. Source: `docs/agents/agent-instructions.md`.\n" +
			"- Active ADR obligations: not applicable — no ADR applies. Source: `docs/agents/agent-instructions.md`.\n" +
			"- " + row(paths) + "\n"
	}
	writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_prd.md", "---\nstatus: active\n---\n\n# Fixture\n\n## Project Constraints\n\n"+constraints(prdPaths)+"\n## Success Metrics\n\nNone. This fixture has no product metric.\n")
	writeCitationFixtureFile(fixture.t, fixture.repoRoot, "docs/specs/"+fixture.slug+"/_techspec.md", "# Fixture TechSpec\n\n## Project Constraints\n\n"+constraints(techSpecPaths)+"\n## API Contracts\n\nNone. This fixture has no public API.\n")
}

func (fixture *undeclaredFixture) requireFindingCount(want int) {
	fixture.t.Helper()
	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		fixture.t.Fatalf("Check(): %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeToolingUndeclared); len(findings) != want {
		fixture.t.Fatalf("%s findings = %#v, want %d", speccheck.CodeToolingUndeclared, findings, want)
	}
}
