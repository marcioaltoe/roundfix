// Suite: command-file Wave collision evidence.
// Invariant: command Tasks collide only when they declare the same command file.
// Boundary IN: temporary Spec Task Graph fixtures and the public stage checker.
// Boundary OUT: the collision detector implementation and repository Task files.
package speccheck_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestWaveCollisionAllowsDifferentCommandFiles(t *testing.T) {
	t.Parallel()

	result := checkCommandFileCollisionFixture(t, ".agents/skills/roundfix/references/deliver.md", ".agents/skills/roundfix/references/review.md")
	if findings := findingsWithCode(result, speccheck.CodeWaveCollision); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none for different command files", speccheck.CodeWaveCollision, findings)
	}
}

func TestWaveCollisionRefusesTheSameCommandFile(t *testing.T) {
	t.Parallel()

	const commandFile = ".agents/skills/roundfix/references/deliver.md"
	result := checkCommandFileCollisionFixture(t, commandFile, commandFile)
	findings := findingsWithCode(result, speccheck.CodeWaveCollision)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one for the same command file", speccheck.CodeWaveCollision, findings)
	}
	if !strings.Contains(findings[0].Summary, commandFile) {
		t.Fatalf("finding summary = %q, want command file %q", findings[0].Summary, commandFile)
	}
}

func checkCommandFileCollisionFixture(t *testing.T, firstCommandFile, secondCommandFile string) speccheck.Result {
	t.Helper()

	const slug = "wave-collision-command-files"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: wave-collision-command-files
status: active
created: 2026-09-30
surfaces: [backend]
---

# Wave collision command files
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: wave-collision-command-files
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
	for id, commandFile := range map[string]string{"task_01": firstCommandFile, "task_02": secondCommandFile} {
		writeCitationFixtureFile(t, repoRoot, commandFile, "# Command reference\n")
		taskContent := fmt.Sprintf("---\nstatus: pending\ntype: docs\n---\n\n# Task\n\n## Context\n\n- interface: `%s`\n\n## Verification\n\n- `test -f %s`\n", commandFile, commandFile)
		writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/"+id+".md", taskContent)
	}

	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	return result
}
