// Suite: instruction-path Wave collision evidence.
// Invariant: read-only instruction context never contributes a Task touch.
// Boundary IN: Task Context and Verification declarations consumed by Collisions.
// Boundary OUT: authoring and Daemon presentation of reported collisions.
package spec

import (
	"reflect"
	"testing"
)

func TestCollisionsIgnoresASharedInstructionPath(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	const sharedPath = ".agents/skills/implement-task/SKILL.md"
	writeCollisionFile(t, repoRoot, sharedPath, "# Implement task\n")
	context := []TaskContextRef{{Kind: ContextKindInstruction, Path: sharedPath}}
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Context: context},
		{ID: "task_02", Context: context},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want none for shared instruction context", collisions)
	}
}

func TestCollisionsStillReportsASharedInterfacePath(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	const sharedPath = "internal/spec/collision.go"
	writeCollisionFile(t, repoRoot, sharedPath, "package spec\n")
	context := []TaskContextRef{{Kind: ContextKindInterface, Path: sharedPath}}
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Context: context},
		{ID: "task_02", Context: context},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths:  map[string]TouchSource{sharedPath: TouchFromContext},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsReportsAnInstructionPathReadByBothVerifications(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	const sharedPath = ".agents/skills/implement-task/SKILL.md"
	writeCollisionFile(t, repoRoot, sharedPath, "# Implement task\n")
	context := []TaskContextRef{{Kind: ContextKindInstruction, Path: sharedPath}}
	verification := []string{"test -f " + sharedPath}
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Context: context, Verification: verification},
		{ID: "task_02", Context: context, Verification: verification},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths:  map[string]TouchSource{sharedPath: TouchFromVerification},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}
