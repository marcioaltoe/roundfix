// Suite: instruction-path authoring collision checks.
// Invariant: speccheck reports only writable shared Task paths as Wave collisions.
// Boundary IN: complete Task artifacts read through speccheck.Check.
// Boundary OUT: Daemon pre-dispatch collision refusal.
package speccheck_test

import (
	"path/filepath"
	"testing"

	"roundfix/internal/speccheck"
)

func TestWaveCollisionCheckIgnoresASharedInstructionPath(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeWaveInstructionFixture(t, "instruction")
	result, err := speccheck.Check(specsRoot, repoRoot, slug)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeWaveCollision); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeWaveCollision, findings)
	}
}

func TestWaveCollisionCheckStillReportsASharedInterfacePath(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeWaveInstructionFixture(t, "interface")
	result, err := speccheck.Check(specsRoot, repoRoot, slug)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeWaveCollision)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeWaveCollision, findings)
	}
}

func writeWaveInstructionFixture(t *testing.T, contextKind string) (string, string, string) {
	t.Helper()

	const slug = "wave-instruction"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	const sharedPath = ".agents/skills/implement-task/SKILL.md"
	writeCitationFixtureFile(t, repoRoot, sharedPath, "# Implement task\n")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: wave-instruction
status: active
created: 2026-09-28
surfaces: [backend]
---

# Wave instruction fixture
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: wave-instruction
qa: declined
qa_reason: no behavioral surface in this fixture
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
---
`)
	for _, taskID := range []string{"task_01", "task_02"} {
		writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/"+taskID+".md", `---
status: pending
type: backend
---

# Task: `+taskID+`

## Context

- `+contextKind+`: `+"`"+sharedPath+"`"+`

## Verification

- `+"`"+`go test ./internal/spec`+"`"+`
`)
	}
	return repoRoot, specsRoot, slug
}
