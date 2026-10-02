package speccheck_test

import (
	"roundfix/internal/speccheck"
	"testing"
)

func TestDeletesPathIsNotAnUnresolvedReference(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- deletes: `old/file.txt`", "true")
	result, err := speccheck.Check(fixture.specsRoot, fixture.repoRoot, fixture.slug)
	if err != nil {
		t.Fatal(err)
	}
	if findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved); len(findings) != 0 {
		t.Fatalf("unresolved = %#v", findings)
	}
}

func TestDeletingAGovernedPathNeedsItsGrant(t *testing.T) {
	t.Parallel()
	fixture := newUndeclaredFixture(t)
	fixture.writeTask("pending", "backend", "- deletes: `Makefile`", "true")
	fixture.writeGrant(nil, nil)
	fixture.writeArtifacts(nil, nil)
	fixture.requireFindingCount(1)
	fixture.writeGrant([]string{"Makefile"}, nil)
	fixture.writeArtifacts([]string{"Makefile"}, []string{"Makefile"})
	fixture.requireFindingCount(0)
}

func TestDeletesCLISurfaceNeedsAGuideAndDeletionIsNotAGuide(t *testing.T) {
	t.Parallel()
	fixture := newCLISurfaceFixture(t)
	fixture.writeTasks([]cliSurfaceTask{{id: "task_01", context: "- deletes: `internal/cli/flags.go`\n- deletes: `docs/user-guide/old.md`"}})
	fixture.requireFindingCount(1)
	fixture.writeTasks([]cliSurfaceTask{{id: "task_01", context: "- deletes: `internal/cli/flags.go`\n- interface: `.agents/skills/roundfix/SKILL.md`"}})
	fixture.requireFindingCount(0)
}
