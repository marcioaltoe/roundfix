// Suite: Task Graph Wave collision evidence.
// Invariant: only independent Tasks with a shared repository file are reported.
// Boundary IN: Task Graph declarations, repository files, and local Git objects.
// Boundary OUT: authoring and Run-time collision presentation in later Tasks.
package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

func TestCollisionsFindsTheMeasuredShape(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeCollisionFile(t, repoRoot, "internal/speccheck/mechanical.go", "package speccheck\n")
	graph := &Graph{Tasks: []Task{
		{
			ID:           "task_01",
			Verification: []string{`grep -q "first" internal/speccheck/mechanical.go || exit 1`},
		},
		{
			ID:           "task_02",
			Verification: []string{`grep -q "second" internal/speccheck/mechanical.go || exit 1`},
		},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths: map[string]TouchSource{
			"internal/speccheck/mechanical.go": TouchFromVerification,
		},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsReturnsEveryPairAndSharedPath(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeCollisionFile(t, repoRoot, "internal/spec/first.go", "package spec\n")
	writeCollisionFile(t, repoRoot, "internal/spec/second.go", "package spec\n")
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Verification: []string{"test -f internal/spec/first.go; test -f internal/spec/second.go"}},
		{ID: "task_02", Verification: []string{"test -f internal/spec/first.go"}},
		{ID: "task_03", Verification: []string{"test -f internal/spec/first.go; test -f internal/spec/second.go"}},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{
		{
			First:  "task_01",
			Second: "task_02",
			Paths:  map[string]TouchSource{"internal/spec/first.go": TouchFromVerification},
		},
		{
			First:  "task_01",
			Second: "task_03",
			Paths: map[string]TouchSource{
				"internal/spec/first.go":  TouchFromVerification,
				"internal/spec/second.go": TouchFromVerification,
			},
		},
		{
			First:  "task_02",
			Second: "task_03",
			Paths:  map[string]TouchSource{"internal/spec/first.go": TouchFromVerification},
		},
	}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsExcludesTransitivelyOrderedTasks(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeCollisionFile(t, repoRoot, "internal/spec/shared.go", "package spec\n")
	verification := []string{"go test internal/spec/shared.go"}
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Verification: verification},
		{ID: "task_02", Needs: []string{"task_01"}, Verification: verification},
		{ID: "task_03", Needs: []string{"task_02"}, Verification: verification},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want no collision for a transitive needs chain", collisions)
	}
}

func TestCollisionsRejectsPackageSelectorsFlagsAndTestNames(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeCollisionFile(t, repoRoot, "internal/cli/cli.go", "package cli\n")
	verification := []string{"go test -run TestCommand ./internal/cli"}
	graph := &Graph{Tasks: []Task{
		{ID: "task_01", Verification: verification},
		{ID: "task_02", Verification: verification},
	}}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want package selector, flag, and test name ignored", collisions)
	}
}

func TestCollisionsLearnsPathFromDeclaredContext(t *testing.T) {
	t.Parallel()
	repoRoot := t.TempDir()
	writeCollisionFile(t, repoRoot, "internal/spec/task.go", "package spec\n")
	context := []TaskContextRef{{Kind: ContextKindInterface, Path: "internal/spec/task.go"}}
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
		Paths:  map[string]TouchSource{"internal/spec/task.go": TouchFromContext},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsLearnsPathFromPackedPriorRunSettlementCommitsWithoutCommands(t *testing.T) {
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n\nconst prior = 0\n")
	gittest.Run(t, repoRoot, "add", "internal/spec/prior.go")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	commitCollisionTask(t, repoRoot, "0097-collision", "task_01", "package spec\n\nconst prior = 1\n")
	commitCollisionTask(t, repoRoot, "0097-collision", "task_02", "package spec\n\nconst prior = 2\n")
	gittest.Run(t, repoRoot, "gc", "--prune=now")
	t.Setenv("PATH", t.TempDir())

	verification := []string{"go test ./internal/spec"}
	graph := &Graph{
		Spec: Spec{Slug: "0097-collision"},
		Tasks: []Task{
			{ID: "task_01", Verification: verification},
			{ID: "task_02", Verification: verification},
		},
	}
	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error with command lookup disabled: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths:  map[string]TouchSource{"internal/spec/prior.go": TouchFromPriorRun},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsRequiresGraphAndRepositoryRoot(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		repoRoot string
		graph    *Graph
	}{
		{name: "missing graph", repoRoot: t.TempDir()},
		{name: "missing repository root", graph: &Graph{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if _, err := Collisions(test.repoRoot, test.graph); err == nil {
				t.Fatal("Collisions returned nil error")
			}
		})
	}
}

func writeCollisionFile(t *testing.T, repoRoot, path, content string) {
	t.Helper()
	absolute := filepath.Join(repoRoot, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", path, err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func commitCollisionTask(t *testing.T, repoRoot, slug, taskID, content string) {
	t.Helper()
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", content)
	gittest.Run(t, repoRoot, "add", "internal/spec/prior.go")
	gittest.Run(t, repoRoot, "commit", "-m", "feat: collision fixture", "-m", "Roundfix-Spec: "+slug+"\nRoundfix-Task: "+taskID)
}
