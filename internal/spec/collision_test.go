// Suite: Task Graph Wave collision evidence.
// Invariant: only independent Tasks with a shared repository file are reported.
// Boundary IN: Task Graph declarations, repository files, and local Git objects.
// Boundary OUT: authoring and Run-time collision presentation in later Tasks.
package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func TestCollisionsLearnsPathFromPackedPriorRunSettlementCommits(t *testing.T) {
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n\nconst prior = 0\n")
	gittest.Run(t, repoRoot, "add", "internal/spec/prior.go")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	commitCollisionTask(t, repoRoot, "0097-collision", "task_01", "package spec\n\nconst prior = 1\n")
	commitCollisionTask(t, repoRoot, "0097-collision", "task_02", "package spec\n\nconst prior = 2\n")
	// Packed objects, which is the ordinary state of any repository that has
	// been gc'd. git reads packs natively; the rule asks git rather than
	// parsing them itself.
	gittest.Run(t, repoRoot, "gc", "--prune=now")

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
		t.Fatalf("Collisions returned error over packed objects: %v", err)
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

func TestCollisionsUsesOnlyTheNewestSettlementPerTask(t *testing.T) {
	// A Task that settled in more than one Run is the ordinary shape once a Run
	// ends Unresolved and the next re-executes it. Unioning an older settlement
	// into current evidence would refuse a Wave that is safe.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n")
	writeCollisionFile(t, repoRoot, "internal/spec/current.go", "package spec\n")
	gittest.Run(t, repoRoot, "add", ".")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	// task_01 settled once on the shared file, then again on a file of its own.
	// Only the newer settlement is its evidence.
	commitCollisionTask(t, repoRoot, "0097-collision", "task_01", "package spec\n\nconst a = 1\n")
	writeCollisionFile(t, repoRoot, "internal/spec/current.go", "package spec\n\nconst b = 1\n")
	gittest.Run(t, repoRoot, "add", "internal/spec/current.go")
	gittest.Run(t, repoRoot, "commit", "-m", "feat: newer settlement", "-m", "Roundfix-Spec: 0097-collision\nRoundfix-Task: task_01")
	// task_02 touches only the file task_01 has stopped touching.
	commitCollisionTask(t, repoRoot, "0097-collision", "task_02", "package spec\n\nconst c = 1\n")

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want none: task_01's older settlement is superseded", collisions)
	}
}

func TestCollisionsReadsTrailersNotTheMessageBody(t *testing.T) {
	// A body that mentions another Spec must not combine with the terminal
	// Roundfix-Task trailer. git knows where the trailer block starts; a scan
	// of every line does not.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n")
	gittest.Run(t, repoRoot, "add", ".")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n\nconst a = 1\n")
	gittest.Run(t, repoRoot, "add", "internal/spec/prior.go")
	gittest.Run(t, repoRoot, "commit",
		"-m", "feat: mentions another Spec in its body",
		"-m", "Refers to Roundfix-Spec: 0097-collision in prose, which is not a trailer.",
		"-m", "Roundfix-Task: task_01")

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want none: the Spec appears in the body, not in a trailer", collisions)
	}
}

func TestCollisionsIgnoresPathsHistoryNamesButTheTreeNoLongerCarries(t *testing.T) {
	// History names deleted files. Two Tasks that both touched a file since
	// removed do not share a file now, and reporting them would refuse a Wave
	// over a path that does not exist.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n")
	gittest.Run(t, repoRoot, "add", ".")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	commitCollisionTask(t, repoRoot, "0097-collision", "task_01", "package spec\n\nconst a = 1\n")
	commitCollisionTask(t, repoRoot, "0097-collision", "task_02", "package spec\n\nconst b = 1\n")
	gittest.Run(t, repoRoot, "rm", "internal/spec/prior.go")
	gittest.Run(t, repoRoot, "commit", "-m", "remove the shared file")

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	if len(collisions) != 0 {
		t.Fatalf("Collisions = %#v, want none: the shared path was deleted before the check ran", collisions)
	}
}

func TestCollisionsReadsSettlementsOnRunBranchesNotOnlyHead(t *testing.T) {
	// A settlement commit lives on its Run Branch until integration moves it,
	// and a Run that ended Unresolved leaves it there. Reading HEAD alone would
	// miss the prior Run this source exists to read.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	writeCollisionFile(t, repoRoot, "internal/spec/prior.go", "package spec\n")
	gittest.Run(t, repoRoot, "add", ".")
	gittest.Run(t, repoRoot, "commit", "-m", "seed")

	gittest.Run(t, repoRoot, "switch", "-c", "roundfix/run-run_20260902T000000Z_abandoned")
	commitCollisionTask(t, repoRoot, "0097-collision", "task_01", "package spec\n\nconst a = 1\n")
	commitCollisionTask(t, repoRoot, "0097-collision", "task_02", "package spec\n\nconst b = 1\n")
	gittest.Run(t, repoRoot, "switch", "main")

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}

	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths:  map[string]TouchSource{"internal/spec/prior.go": TouchFromPriorRun},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v: the settlements are reachable only from the Run Branch", collisions, want)
	}
}

func TestCollisionsReadsHistoricalPathsGitReportsVerbatim(t *testing.T) {
	// A filename can contain a character that is a shell metacharacter and an
	// ordinary letter in a path. Verification candidates are tokens from a
	// command line and are filtered for those; a path Git reports is the path,
	// and filtering it there would discard real evidence.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	const shared = "internal/spec/fixture[1].go"
	writeCollisionFile(t, repoRoot, shared, "package spec\n")
	gittest.Run(t, repoRoot, "add", "--", shared)
	gittest.Run(t, repoRoot, "commit", "-m", "seed")
	for _, taskID := range []string{"task_01", "task_02"} {
		writeCollisionFile(t, repoRoot, shared, "package spec\n\nconst x = \""+taskID+"\"\n")
		gittest.Run(t, repoRoot, "add", "--", shared)
		gittest.Run(t, repoRoot, "commit", "-m", "feat: bracket fixture",
			"-m", "Roundfix-Spec: 0097-collision\nRoundfix-Task: "+taskID)
	}

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}
	collisions, err := Collisions(repoRoot, graph)
	if err != nil {
		t.Fatalf("Collisions returned error: %v", err)
	}
	want := []WaveCollision{{
		First:  "task_01",
		Second: "task_02",
		Paths:  map[string]TouchSource{shared: TouchFromPriorRun},
	}}
	if !reflect.DeepEqual(collisions, want) {
		t.Fatalf("Collisions = %#v, want %#v", collisions, want)
	}
}

func TestCollisionsTreatsAnAbsentRepositoryAsAbsentEvidence(t *testing.T) {
	// A Spec can be checked outside a working repository, and a repository can
	// have no commits yet. Neither is a read failure, and neither may fail the
	// check: prior-Run history is one of three sources.
	t.Parallel()
	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}
	for _, test := range []struct {
		name  string
		setup func(t *testing.T, repoRoot string)
	}{
		{name: "outside any repository", setup: func(*testing.T, string) {}},
		{name: "repository with no commits", setup: func(t *testing.T, repoRoot string) {
			gittest.InitRepo(t, repoRoot, "--initial-branch=main")
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repoRoot := t.TempDir()
			test.setup(t, repoRoot)
			collisions, err := Collisions(repoRoot, graph)
			if err != nil {
				t.Fatalf("Collisions failed on %s: %v", test.name, err)
			}
			if len(collisions) != 0 {
				t.Fatalf("Collisions = %#v, want none", collisions)
			}
		})
	}
}

func TestCollisionsReportsAnInspectionFailureOnAPriorRunPath(t *testing.T) {
	// A path history names can be unreadable for a reason that is not
	// not-exists — here a symlink loop. Excluding it the way a missing path is
	// excluded would drop the collision evidence silently, so it is reported
	// with the path and the Task that settled it.
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	const shared = "internal/spec/loop.go"
	writeCollisionFile(t, repoRoot, shared, "package spec\n")
	gittest.Run(t, repoRoot, "add", "--", shared)
	gittest.Run(t, repoRoot, "commit", "-m", "seed")
	for _, taskID := range []string{"task_01", "task_02"} {
		writeCollisionFile(t, repoRoot, shared, "package spec\n\nconst x = \""+taskID+"\"\n")
		gittest.Run(t, repoRoot, "add", "--", shared)
		gittest.Run(t, repoRoot, "commit", "-m", "feat: loop fixture",
			"-m", "Roundfix-Spec: 0097-collision\nRoundfix-Task: "+taskID)
	}

	// The tree now carries a cycle where the settled file was.
	absolute := filepath.Join(repoRoot, filepath.FromSlash(shared))
	partner := filepath.Join(filepath.Dir(absolute), "loop-partner.go")
	if err := os.Remove(absolute); err != nil {
		t.Fatalf("remove settled file: %v", err)
	}
	if err := os.Symlink(partner, absolute); err != nil {
		t.Fatalf("link settled file to its partner: %v", err)
	}
	if err := os.Symlink(absolute, partner); err != nil {
		t.Fatalf("link partner back: %v", err)
	}

	graph := &Graph{
		Spec:  Spec{Slug: "0097-collision"},
		Tasks: []Task{{ID: "task_01"}, {ID: "task_02"}},
	}
	_, err := Collisions(repoRoot, graph)
	if err == nil {
		t.Fatal("Collisions reported no error over an unreadable prior Run path")
	}
	// task_02 settled last, and history is read newest first, so its
	// settlement is the one being read when the inspection fails.
	for _, want := range []string{shared, "task_02"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Collisions error %q does not name %q", err, want)
		}
	}
}
