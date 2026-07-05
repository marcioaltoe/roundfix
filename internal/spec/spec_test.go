package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// md turns single quotes into backticks so markdown fixtures can be raw
// string literals, which cannot contain backtick characters.
func md(fixture string) string {
	return strings.ReplaceAll(fixture, "'", "`")
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create directory for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}

func writeSpecDir(t *testing.T, gitRoot string, slug string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		writeFile(t, filepath.Join(gitRoot, "docs", "specs", slug, name), content)
	}
}

func prdFixture(status string) string {
	return fmt.Sprintf(`---
spec: demo
status: %s
created: 2026-07-04
surfaces: [cli]
---

# Demo
`, status)
}

const defaultVerificationSection = `## Verification

- 'go test ./...' — expected: pass.
`

func taskFixture(id string, title string, status string, taskType string, tail string) string {
	number := strings.TrimPrefix(id, "task_")
	return md(fmt.Sprintf(`---
task: %s
spec: demo
status: %s
type: %s
complexity: low
---

# Task %s: %s

## Overview

Fixture task.

%s`, id, status, taskType, number, title, tail))
}

func manifestFixture(schema string, nodes string) string {
	return fmt.Sprintf(`---
schema: %s
spec: demo
graph:
  nodes:
%s---

# Tasks — Demo
`, schema, nodes)
}

const diamondNodes = `    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_03
      file: task_03.md
      needs: [task_01]
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_01
      file: task_01.md
      needs: []
`

func diamondSpecFiles() map[string]string {
	files := map[string]string{
		"_prd.md":   prdFixture("active"),
		"_tasks.md": manifestFixture("spec-tasks/v1", diamondNodes),
	}
	for _, id := range []string{"task_01", "task_02", "task_03", "task_04"} {
		files[id+".md"] = taskFixture(id, "Fixture "+id, "pending", "backend", defaultVerificationSection)
	}
	return files
}

func TestLoadReturnsTasksInDeterministicTopologicalOrder(t *testing.T) {
	gitRoot := t.TempDir()
	writeSpecDir(t, gitRoot, "demo", diamondSpecFiles())

	// Kahn with manifest-order tiebreak: task_01 is the only root; task_03
	// beats task_02 because the manifest lists it first.
	want := []string{"task_01", "task_03", "task_02", "task_04"}
	for attempt := 0; attempt < 20; attempt++ {
		graph, err := Load(gitRoot, "demo")
		if err != nil {
			t.Fatalf("Load attempt %d: %v", attempt, err)
		}
		got := make([]string, 0, len(graph.Tasks))
		for _, task := range graph.Tasks {
			got = append(got, task.ID)
		}
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Fatalf("Load attempt %d order = %v, want %v", attempt, got, want)
		}
	}

	graph, err := Load(gitRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if graph.Spec.Slug != "demo" {
		t.Errorf("Spec.Slug = %q, want %q", graph.Spec.Slug, "demo")
	}
	wantDir := filepath.Join(gitRoot, "docs", "specs", "demo")
	if graph.Spec.Dir != wantDir {
		t.Errorf("Spec.Dir = %q, want %q", graph.Spec.Dir, wantDir)
	}
}

func TestLoadParsesTaskFiles(t *testing.T) {
	gitRoot := t.TempDir()
	verification := md(`## Verification

- 'go test ./internal/spec/' — expected: all tests pass.
- a bullet without a command is skipped
- 'go build ./...' — expected: builds cleanly.

## References

- 'go vet ./...' outside the Verification section is not a command.
`)
	writeSpecDir(t, gitRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
		"task_01.md": taskFixture("task_01", "Build the parser", "in_progress", "docs", verification),
	})

	graph, err := Load(gitRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(graph.Tasks) != 1 {
		t.Fatalf("len(Tasks) = %d, want 1", len(graph.Tasks))
	}
	task := graph.Tasks[0]
	if task.ID != "task_01" {
		t.Errorf("ID = %q, want %q", task.ID, "task_01")
	}
	if task.File != filepath.Join("docs", "specs", "demo", "task_01.md") {
		t.Errorf("File = %q, want repository-relative task path", task.File)
	}
	if task.Title != "Build the parser" {
		t.Errorf("Title = %q, want %q", task.Title, "Build the parser")
	}
	if task.Status != StatusInProgress {
		t.Errorf("Status = %q, want %q", task.Status, StatusInProgress)
	}
	if task.Type != "docs" {
		t.Errorf("Type = %q, want %q", task.Type, "docs")
	}
	wantCommands := []string{"go test ./internal/spec/", "go build ./..."}
	if len(task.Verification) != len(wantCommands) {
		t.Fatalf("Verification = %v, want %v", task.Verification, wantCommands)
	}
	for index, command := range wantCommands {
		if task.Verification[index] != command {
			t.Errorf("Verification[%d] = %q, want %q", index, task.Verification[index], command)
		}
	}
}

func TestLoadReturnsTypedValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		check func(t *testing.T, err error)
	}{
		{
			name:  "missing Spec",
			files: nil,
			check: func(t *testing.T, err error) {
				var notFound SpecNotFoundError
				if !errors.As(err, &notFound) {
					t.Fatalf("error = %v, want SpecNotFoundError", err)
				}
				if notFound.Slug != "demo" {
					t.Errorf("Slug = %q, want %q", notFound.Slug, "demo")
				}
				if !strings.Contains(err.Error(), `"demo"`) {
					t.Errorf("message %q does not name the Spec", err)
				}
			},
		},
		{
			name: "inactive Spec",
			files: map[string]string{
				"_prd.md": prdFixture("draft"),
			},
			check: func(t *testing.T, err error) {
				var inactive InactiveSpecError
				if !errors.As(err, &inactive) {
					t.Fatalf("error = %v, want InactiveSpecError", err)
				}
				if inactive.Status != "draft" {
					t.Errorf("Status = %q, want %q", inactive.Status, "draft")
				}
				if !strings.Contains(err.Error(), `"draft"`) {
					t.Errorf("message %q does not name the PRD status", err)
				}
			},
		},
		{
			name: "missing manifest",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
			},
			check: func(t *testing.T, err error) {
				var manifestErr ManifestError
				if !errors.As(err, &manifestErr) {
					t.Fatalf("error = %v, want ManifestError", err)
				}
				if !strings.Contains(err.Error(), "does not exist") {
					t.Errorf("message %q does not name the missing manifest check", err)
				}
			},
		},
		{
			name: "wrong schema",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v2", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var schemaErr ManifestSchemaError
				if !errors.As(err, &schemaErr) {
					t.Fatalf("error = %v, want ManifestSchemaError", err)
				}
				if schemaErr.Schema != "spec-tasks/v2" {
					t.Errorf("Schema = %q, want %q", schemaErr.Schema, "spec-tasks/v2")
				}
				if !strings.Contains(err.Error(), "spec-tasks/v1") {
					t.Errorf("message %q does not name the expected schema", err)
				}
			},
		},
		{
			name: "empty graph",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": `---
schema: spec-tasks/v1
spec: demo
graph:
  nodes: []
---

# Tasks — Demo
`,
			},
			check: func(t *testing.T, err error) {
				var manifestErr ManifestError
				if !errors.As(err, &manifestErr) {
					t.Fatalf("error = %v, want ManifestError", err)
				}
				if !strings.Contains(err.Error(), "no nodes") {
					t.Errorf("message %q does not name the empty-graph check", err)
				}
			},
		},
		{
			name: "duplicate Task id",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_01
      file: task_02.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var manifestErr ManifestError
				if !errors.As(err, &manifestErr) {
					t.Fatalf("error = %v, want ManifestError", err)
				}
				if !strings.Contains(err.Error(), `"task_01"`) {
					t.Errorf("message %q does not name the duplicated Task", err)
				}
			},
		},
		{
			name: "unknown needs id",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: [task_99]
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var unknownNeed UnknownNeedError
				if !errors.As(err, &unknownNeed) {
					t.Fatalf("error = %v, want UnknownNeedError", err)
				}
				if unknownNeed.TaskID != "task_01" || unknownNeed.Need != "task_99" {
					t.Errorf("UnknownNeedError = %+v, want task_01 needing task_99", unknownNeed)
				}
				if !strings.Contains(err.Error(), `"task_01"`) || !strings.Contains(err.Error(), `"task_99"`) {
					t.Errorf("message %q does not name the Task and the unknown need", err)
				}
			},
		},
		{
			name: "cycle",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: [task_02]
    - id: task_02
      file: task_02.md
      needs: [task_01]
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
				"task_02.md": taskFixture("task_02", "Fixture", "pending", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var cycle CycleError
				if !errors.As(err, &cycle) {
					t.Fatalf("error = %v, want CycleError", err)
				}
				if len(cycle.TaskIDs) != 2 || cycle.TaskIDs[0] != "task_01" || cycle.TaskIDs[1] != "task_02" {
					t.Errorf("TaskIDs = %v, want [task_01 task_02] in manifest order", cycle.TaskIDs)
				}
				if !strings.Contains(err.Error(), "task_01") || !strings.Contains(err.Error(), "task_02") {
					t.Errorf("message %q does not name the cycling Tasks", err)
				}
			},
		},
		{
			name: "missing task file",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var missing MissingTaskFileError
				if !errors.As(err, &missing) {
					t.Fatalf("error = %v, want MissingTaskFileError", err)
				}
				if missing.TaskID != "task_02" {
					t.Errorf("TaskID = %q, want %q", missing.TaskID, "task_02")
				}
				if !strings.Contains(err.Error(), `"task_02"`) {
					t.Errorf("message %q does not name the Task", err)
				}
			},
		},
		{
			name: "no Verification commands",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", `## Verification

- run the tests by hand, no backticked command here.
`),
			},
			check: func(t *testing.T, err error) {
				var missing MissingVerificationError
				if !errors.As(err, &missing) {
					t.Fatalf("error = %v, want MissingVerificationError", err)
				}
				if missing.TaskID != "task_01" {
					t.Errorf("TaskID = %q, want %q", missing.TaskID, "task_01")
				}
				if !strings.Contains(err.Error(), `"task_01"`) {
					t.Errorf("message %q does not name the Task", err)
				}
			},
		},
		{
			name: "unparseable task frontmatter",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": "# Task 01: no frontmatter at all\n",
			},
			check: func(t *testing.T, err error) {
				var taskErr TaskFileError
				if !errors.As(err, &taskErr) {
					t.Fatalf("error = %v, want TaskFileError", err)
				}
				if taskErr.TaskID != "task_01" {
					t.Errorf("TaskID = %q, want %q", taskErr.TaskID, "task_01")
				}
			},
		},
		{
			name: "unsupported task status",
			files: map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Fixture", "done", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var taskErr TaskFileError
				if !errors.As(err, &taskErr) {
					t.Fatalf("error = %v, want TaskFileError", err)
				}
				if !strings.Contains(err.Error(), `"done"`) {
					t.Errorf("message %q does not name the unsupported status", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			if tt.files != nil {
				writeSpecDir(t, gitRoot, "demo", tt.files)
			}
			_, err := Load(gitRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want a typed validation error")
			}
			tt.check(t, err)
		})
	}
}

func TestLoadAcceptsValidGraph(t *testing.T) {
	gitRoot := t.TempDir()
	writeSpecDir(t, gitRoot, "demo", diamondSpecFiles())

	graph, err := Load(gitRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(graph.Tasks) != 4 {
		t.Fatalf("len(Tasks) = %d, want 4", len(graph.Tasks))
	}
	for _, task := range graph.Tasks {
		if task.Status != StatusPending {
			t.Errorf("Task %q Status = %q, want %q", task.ID, task.Status, StatusPending)
		}
		if len(task.Verification) == 0 {
			t.Errorf("Task %q has no Verification commands", task.ID)
		}
	}
}

func TestListActiveFiltersInactiveArchivedAndNonSpecDirectories(t *testing.T) {
	gitRoot := t.TempDir()
	writeSpecDir(t, gitRoot, "0002-later", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, gitRoot, "0001-early", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, gitRoot, "0003-shipped", map[string]string{"_prd.md": prdFixture("shipped")})
	writeSpecDir(t, gitRoot, "_archived", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, gitRoot, "0004-broken", map[string]string{"_prd.md": "no frontmatter"})
	if err := os.MkdirAll(filepath.Join(gitRoot, "docs", "specs", "0005-no-prd"), 0o755); err != nil {
		t.Fatalf("create directory: %v", err)
	}
	writeFile(t, filepath.Join(gitRoot, "docs", "specs", "notes.md"), "# not a spec\n")

	specs, err := ListActive(gitRoot)
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("ListActive = %+v, want two active Specs", specs)
	}
	if specs[0].Slug != "0001-early" || specs[1].Slug != "0002-later" {
		t.Errorf("slugs = [%s %s], want sorted [0001-early 0002-later]", specs[0].Slug, specs[1].Slug)
	}
	wantDir := filepath.Join(gitRoot, "docs", "specs", "0001-early")
	if specs[0].Dir != wantDir {
		t.Errorf("Dir = %q, want %q", specs[0].Dir, wantDir)
	}
}

func TestListActiveDetailedReportsSkippedSpecFolders(t *testing.T) {
	gitRoot := t.TempDir()
	writeSpecDir(t, gitRoot, "0001-active", map[string]string{"_prd.md": prdFixture("active")})
	if err := os.MkdirAll(filepath.Join(gitRoot, "docs", "specs", "0002-missing-prd"), 0o755); err != nil {
		t.Fatalf("create missing PRD fixture: %v", err)
	}
	writeSpecDir(t, gitRoot, "0003-broken-frontmatter", map[string]string{"_prd.md": "no frontmatter\n"})
	writeSpecDir(t, gitRoot, "0004-archived-status", map[string]string{"_prd.md": prdFixture("archived")})
	writeSpecDir(t, gitRoot, "_archived", map[string]string{"0005-old/_prd.md": prdFixture("broken")})

	specs, skipped, err := ListActiveDetailed(gitRoot)
	if err != nil {
		t.Fatalf("ListActiveDetailed: %v", err)
	}
	if len(specs) != 1 || specs[0].Slug != "0001-active" {
		t.Fatalf("active Specs = %+v, want only 0001-active", specs)
	}
	wantSkipped := []SkippedSpec{
		{Dir: "docs/specs/0002-missing-prd", Reason: "missing _prd.md"},
		{Dir: "docs/specs/0003-broken-frontmatter", Reason: "unreadable _prd.md frontmatter: missing YAML frontmatter opening marker"},
		{Dir: "docs/specs/0004-archived-status", Reason: `status "archived" is not active`},
	}
	if len(skipped) != len(wantSkipped) {
		t.Fatalf("skipped = %+v, want %+v", skipped, wantSkipped)
	}
	for index := range wantSkipped {
		if skipped[index] != wantSkipped[index] {
			t.Fatalf("skipped[%d] = %+v, want %+v", index, skipped[index], wantSkipped[index])
		}
	}

	simple, err := ListActive(gitRoot)
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(simple) != 1 || simple[0].Slug != specs[0].Slug {
		t.Fatalf("ListActive = %+v, want active list unchanged from detailed result %+v", simple, specs)
	}
}

func TestListActiveWithoutSpecsRootReturnsNothing(t *testing.T) {
	specs, err := ListActive(t.TempDir())
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(specs) != 0 {
		t.Fatalf("ListActive = %+v, want none", specs)
	}
}
