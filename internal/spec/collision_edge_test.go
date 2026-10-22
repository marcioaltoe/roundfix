// Suite: resolving edges for independent Tasks.
// Invariant: the named edge orders the pair without a cycle.
// Boundary IN: repository files and local Git history on every ref, isolated
// with gittest in temporary repositories; no real ~/.roundfix, other repo or network.
// Boundary OUT: diagnostic and gate presentation.
package spec

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestWaveCollisionNamesTheEdgeThatOrdersThePair(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, first, second, edge string }{
		{"original manifest order", "task_02", "task_05", "add a `needs` edge from Task task_05 to Task task_02"},
		{"reversed manifest order", "task_05", "task_02", "add a `needs` edge from Task task_02 to Task task_05"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			gittest.InitRepo(t, root)
			writeCollisionFile(t, root, "pkg/versions_test.go", "package pkg\n")
			const slug = "wave-edge"
			dir := "docs/specs/" + slug + "/"
			writeCollisionFile(t, root, dir+"_prd.md", "---\nstatus: active\n---\n\n# Wave edge\n")
			manifest := "---\nschema: spec-tasks/v1\nqa: declined\nqa_reason: edge fixture\ngraph:\n  nodes:\n"
			for _, id := range []string{"task_01", tc.first, tc.second} {
				needs := ""
				context := ""
				if id != "task_01" {
					needs = "task_01"
					context = "\n## Context\n\n- interface: `pkg/versions_test.go`\n"
				}
				manifest += fmt.Sprintf("    - id: %s\n      file: %s.md\n      needs: [%s]\n", id, id, needs)
				writeCollisionFile(t, root, dir+id+".md", "---\nstatus: pending\ntype: backend\n---\n\n# "+id+"\n\n## Verification\n\n- `true`\n"+context)
			}
			manifest += "---\n"
			writeCollisionFile(t, root, dir+"_tasks.md", manifest)
			graph, err := Load(filepath.Join(root, "docs", "specs"), slug)
			if err != nil {
				t.Fatal(err)
			}
			collisions, err := Collisions(root, graph)
			if err != nil || len(collisions) != 1 {
				t.Fatalf("Collisions = %#v, %v", collisions, err)
			}
			collision := collisions[0]
			if collision.First != tc.first || collision.Second != tc.second {
				t.Fatalf("pair = %#v", collision)
			}
			if got := collision.ResolvingEdge(); got != tc.edge {
				t.Fatalf("edge = %q, want %q", got, tc.edge)
			}
			if got := collision.SharedPathList(); got != "pkg/versions_test.go (declared context)" {
				t.Fatalf("paths = %q", got)
			}
			manifest = strings.Replace(manifest,
				"file: "+collision.Second+".md\n      needs: [task_01]",
				"file: "+collision.Second+".md\n      needs: [task_01, "+collision.First+"]", 1)
			writeCollisionFile(t, root, dir+"_tasks.md", manifest)
			graph, err = Load(filepath.Join(root, "docs", "specs"), slug)
			if err != nil {
				t.Fatalf("load after named edge: %v", err)
			}
			collisions, err = Collisions(root, graph)
			if err != nil || len(collisions) != 0 {
				t.Fatalf("after edge: Collisions = %#v, %v", collisions, err)
			}
			nodes := make([]manifestNode, len(graph.Tasks))
			for i, task := range graph.Tasks {
				nodes[i] = manifestNode{ID: task.ID, Needs: task.Needs}
			}
			order, err := topologicalOrder(nodes)
			if err != nil || len(order) != len(nodes) {
				t.Fatalf("topologicalOrder = %v, %v", order, err)
			}
		})
	}
}

func TestWaveCollisionSharedPathListSortsPathsWithTheirSources(t *testing.T) {
	t.Parallel()
	collision := WaveCollision{Paths: map[string]TouchSource{
		"z.go": TouchFromVerification,
		"a.go": mergedTouchSource(TouchFromPriorRun, TouchFromContext),
	}}
	const want = "a.go (declared context; prior Run settlement commit), z.go (verification command)"
	if got := collision.SharedPathList(); got != want {
		t.Fatalf("SharedPathList = %q, want %q", got, want)
	}
}
