// Suite: the Wave collision finding's resolving edge.
// Invariant: the finding and durable refusal carry the same concrete remedy.
// Boundary IN: repository files and local Git history on every ref, built with
// gittest in a temporary repository; no real ~/.roundfix, other repo or network.
// Boundary OUT: Daemon orchestration.
package speccheck_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

func TestWaveCollisionFindingNamesTheResolvingEdge(t *testing.T) {
	t.Parallel()
	const slug = "wave-edge"
	root := t.TempDir()
	gittest.InitRepo(t, root)
	dir := "docs/specs/" + slug + "/"
	writeCitationFixtureFile(t, root, "pkg/versions_test.go", "package pkg\n")
	writeCitationFixtureFile(t, root, dir+"_prd.md", "---\nstatus: active\n---\n\n# Wave edge\n")
	writeCitationFixtureFile(t, root, dir+"_tasks.md", `---
schema: spec-tasks/v1
qa: declined
qa_reason: diagnostic fixture
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_01]
---
`)
	for _, id := range []string{"task_01", "task_02", "task_05"} {
		body := "---\nstatus: pending\ntype: backend\n---\n\n# " + id + "\n\n## Verification\n\n- `true`\n"
		if id != "task_01" {
			body += "\n## Context\n\n- interface: `pkg/versions_test.go`\n"
		}
		writeCitationFixtureFile(t, root, dir+id+".md", body)
	}
	result, err := speccheck.Check(filepath.Join(root, "docs", "specs"), root, slug)
	if err != nil {
		t.Fatal(err)
	}
	findings := findingsWithCode(result, speccheck.CodeWaveCollision)
	if len(findings) != 1 {
		t.Fatalf("findings = %#v", findings)
	}
	finding := findings[0]
	const summary = "Task task_02 and Task task_05 can run in the same Wave and share repository files: pkg/versions_test.go (declared context); add a `needs` edge from Task task_05 to Task task_02"
	if finding.Summary != summary {
		t.Fatalf("summary = %q, want %q", finding.Summary, summary)
	}
	const fix = "Give Task task_05 a `needs` edge on Task task_02 in docs/specs/wave-edge/_tasks.md."
	if finding.Fix != fix {
		t.Fatalf("fix = %q, want %q", finding.Fix, fix)
	}
	if got := speccheck.RefusalReason(finding); got != "SC-WAVE-COLLISION: "+summary {
		t.Fatalf("RefusalReason = %q", got)
	}
	wantLocations := []speccheck.Location{{Path: dir + "task_02.md", Line: 1}, {Path: dir + "task_05.md", Line: 1}, {Path: dir + "_tasks.md", Line: 1}}
	if finding.Severity != speccheck.SeverityError || !reflect.DeepEqual(finding.Where, wantLocations) {
		t.Fatalf("finding = %#v", finding)
	}
}
