package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

func defaultSpecsRoot(gitRoot string) string {
	return filepath.Join(gitRoot, "docs", "specs")
}

func writeSpecDir(t *testing.T, specsRoot string, slug string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		writeFile(t, filepath.Join(specsRoot, slug, name), content)
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

func TestNormalizeStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw  string
		want string
	}{
		{raw: "pending", want: "pending"},
		{raw: "in_progress", want: "in_progress"},
		{raw: "completed", want: "completed"},
		{raw: "failed", want: "failed"},
		{raw: "done", want: "completed"},
		{raw: "in-progress", want: "in_progress"},
		{raw: "in progress", want: "in_progress"},
		{raw: "  done\t", want: "completed"},
		{raw: "finished", want: "finished"},
		{raw: "Done", want: "Done"},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			if got := NormalizeStatus(tt.raw); got != tt.want {
				t.Fatalf("NormalizeStatus(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func manifestFixture(schema string, nodes string) string {
	return manifestFixtureWithBody(schema, nodes, "")
}

func manifestFixtureWithBody(schema string, nodes string, body string) string {
	return manifestFixtureWithQA(schema, "", nodes, body)
}

func manifestFixtureWithQA(schema string, declaration string, nodes string, body string) string {
	return fmt.Sprintf(`---
schema: %s
spec: demo
%sgraph:
  nodes:
%s---

# Tasks — Demo
%s`, schema, declaration, nodes, body)
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
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", diamondSpecFiles())

	// Kahn with manifest-order tiebreak: task_01 is the only root; task_03
	// beats task_02 because the manifest lists it first.
	want := []string{"task_01", "task_03", "task_02", "task_04"}
	for attempt := 0; attempt < 20; attempt++ {
		graph, err := Load(specsRoot, "demo")
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

	graph, err := Load(specsRoot, "demo")
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
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	verification := md(`## Verification

- 'go test ./internal/spec/' — expected: all tests pass.
- a bullet without a command is skipped
- 'go build ./...' — expected: builds cleanly.

## References

- 'go vet ./...' outside the Verification section is not a command.
`)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
		"task_01.md": taskFixture("task_01", "Build the parser", "in_progress", "docs", verification),
	})

	graph, err := Load(specsRoot, "demo")
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
	if task.File != filepath.Join("demo", "task_01.md") {
		t.Errorf("File = %q, want Spec Root-relative task path", task.File)
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

func TestTaskTypeCanonicalValuesLoadThroughTaskGraph(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	taskTypes := []string{"backend", "frontend", "data", "infra", "docs", "test", "chore", "qa"}
	files := map[string]string{"_prd.md": prdFixture("active")}
	var nodes strings.Builder
	var table strings.Builder
	table.WriteString("\n| id      | title   | type    | complexity | needs |\n")
	table.WriteString("| ------- | ------- | ------- | ---------- | ----- |\n")
	for index, taskType := range taskTypes {
		id := fmt.Sprintf("task_%02d", index+1)
		needs := "[]"
		projectedNeeds := "—"
		if taskType == "qa" {
			needs = "[task_01, task_02, task_03, task_04, task_05, task_06, task_07]"
			projectedNeeds = "task_01, task_02, task_03, task_04, task_05, task_06, task_07"
		}
		nodes.WriteString(fmt.Sprintf("    - id: %s\n      file: %s.md\n      needs: %s\n", id, id, needs))
		table.WriteString(fmt.Sprintf("| %s | Fixture | %s | low | %s |\n", id, taskType, projectedNeeds))
		files[id+".md"] = taskFixture(id, "Fixture", "pending", taskType, defaultVerificationSection)
	}
	files["_tasks.md"] = manifestFixtureWithQA("spec-tasks/v1", "qa: task_08\n", nodes.String(), table.String())
	writeSpecDir(t, specsRoot, "demo", files)

	graph, err := Load(specsRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(graph.Tasks) != len(taskTypes) {
		t.Fatalf("len(Tasks) = %d, want %d", len(graph.Tasks), len(taskTypes))
	}
	for index, task := range graph.Tasks {
		if string(task.Type) != taskTypes[index] {
			t.Fatalf("Task %q Type = %q, want %q", task.ID, task.Type, taskTypes[index])
		}
	}
}

func TestTaskTypeRejectsInvalidFrontmatterValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		taskFile  string
		wantValue string
	}{
		{
			name:      "missing",
			taskFile:  strings.Replace(taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection), "type: backend\n", "", 1),
			wantValue: "",
		},
		{
			name:      "empty",
			taskFile:  taskFixture("task_01", "Fixture", "pending", "", defaultVerificationSection),
			wantValue: "",
		},
		{
			name:      "leading whitespace",
			taskFile:  taskFixture("task_01", "Fixture", "pending", `" backend"`, defaultVerificationSection),
			wantValue: " backend",
		},
		{
			name:      "trailing whitespace",
			taskFile:  taskFixture("task_01", "Fixture", "pending", `"backend "`, defaultVerificationSection),
			wantValue: "backend ",
		},
		{
			name:      "mixed case",
			taskFile:  taskFixture("task_01", "Fixture", "pending", "Backend", defaultVerificationSection),
			wantValue: "Backend",
		},
		{
			name:      "unknown",
			taskFile:  taskFixture("task_01", "Fixture", "pending", "api", defaultVerificationSection),
			wantValue: "api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			writeSpecDir(t, specsRoot, "demo", map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": tt.taskFile,
			})

			_, err := Load(specsRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want Task Type validation error")
			}
			var taskErr TaskFileError
			if !errors.As(err, &taskErr) {
				t.Fatalf("error = %v, want TaskFileError", err)
			}
			var taskTypeErr TaskTypeError
			if !errors.As(err, &taskTypeErr) {
				t.Fatalf("error = %v, want TaskTypeError", err)
			}
			wantPath := filepath.Join(specsRoot, "demo", "task_01.md")
			message := err.Error()
			for _, want := range []string{
				wantPath,
				fmt.Sprintf("%q", tt.wantValue),
				"backend, frontend, data, infra, docs, test, chore, qa",
				"frontmatter",
				"type",
			} {
				if !strings.Contains(message, want) {
					t.Fatalf("error %q does not contain %q", message, want)
				}
			}
		})
	}
}

func TestLoadQAGateContract(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", "qa: task_03\n", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
`, `
| id      | title   | type    | complexity | needs   |
| ------- | ------- | ------- | ---------- | ------- |
| task_01 | Prepare | backend | low        | —       |
| task_02 | Build   | backend | low        | task_01 |
| task_03 | QA      | qa      | low        | task_02 |
`),
		"task_01.md": taskFixture("task_01", "Prepare", "completed", "backend", defaultVerificationSection),
		"task_02.md": taskFixture("task_02", "Build", "completed", "backend", defaultVerificationSection),
		"task_03.md": taskFixture("task_03", "QA", "pending", "qa", defaultVerificationSection),
	})

	graph, err := Load(specsRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if graph.QATaskID != "task_03" {
		t.Fatalf("QATaskID = %q, want %q", graph.QATaskID, "task_03")
	}
	if graph.QADeclined {
		t.Fatal("QADeclined = true, want false")
	}
	if graph.QAReason != "" {
		t.Fatalf("QAReason = %q, want empty", graph.QAReason)
	}
}

func TestLoadQADeclinedContract(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		declaration string
		wantReason  string
		wantError   string
	}{
		{
			name:        "reason recorded",
			declaration: "qa: declined\nqa_reason: no behavioral surface\n",
			wantReason:  "no behavioral surface",
		},
		{
			name:        "reason required",
			declaration: "qa: declined\n",
			wantError:   "qa_reason",
		},
		{
			name:        "reason cannot be empty",
			declaration: "qa: declined\nqa_reason:\n",
			wantError:   "qa_reason",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			writeSpecDir(t, specsRoot, "demo", map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", tt.declaration, `    - id: task_01
      file: task_01.md
      needs: []
`, ""),
				"task_01.md": taskFixture("task_01", "Docs", "pending", "docs", defaultVerificationSection),
			})

			graph, err := Load(specsRoot, "demo")
			if tt.wantError != "" {
				if err == nil {
					t.Fatal("Load succeeded, want QA declaration error")
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("error %q does not contain %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if graph.QATaskID != "" || !graph.QADeclined || graph.QAReason != tt.wantReason {
				t.Fatalf("QA declaration = (%q, %t, %q), want (%q, %t, %q)", graph.QATaskID, graph.QADeclined, graph.QAReason, "", true, tt.wantReason)
			}
		})
	}
}

func TestLoadRejectsExplicitNullQAReasonWithoutDeclaration(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", "qa_reason:\n", `    - id: task_01
      file: task_01.md
      needs: []
`, ""),
		"task_01.md": taskFixture("task_01", "Docs", "pending", "docs", defaultVerificationSection),
	})

	_, err := Load(specsRoot, "demo")
	var gateErr QAGateError
	if !errors.As(err, &gateErr) {
		t.Fatalf("Load error = %T %v, want QAGateError", err, err)
	}
	if !strings.Contains(gateErr.Reason, "qa_reason requires a qa: declaration") {
		t.Fatalf("QAGateError reason = %q", gateErr.Reason)
	}
}

func TestLoadWrapsManifestDecodeErrors(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", "qa: [task_01]\n", `    - id: task_01
      file: task_01.md
      needs: []
`, ""),
		"task_01.md": taskFixture("task_01", "Docs", "pending", "docs", defaultVerificationSection),
	})

	_, err := Load(specsRoot, "demo")

	if err == nil {
		t.Fatal("Load succeeded, want manifest decode error")
	}
	var manifestErr ManifestError
	if !errors.As(err, &manifestErr) {
		t.Fatalf("Load error = %T %v, want ManifestError", err, err)
	}
	if !strings.Contains(err.Error(), "decode manifest frontmatter") {
		t.Fatalf("Load error = %v, want manifest decode context", err)
	}
}

func TestLoadRejectsInvalidQAGateShape(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		declaration string
		nodes       string
		tasks       map[string]string
		wantTask    string
	}{
		{
			name: "qa type without declaration",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "QA", "pending", "qa", defaultVerificationSection),
			},
			wantTask: "task_01",
		},
		{
			name:        "declaration names another type",
			declaration: "qa: task_01\n",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "Build", "pending", "backend", defaultVerificationSection),
			},
			wantTask: "task_01",
		},
		{
			name:        "declined graph has gate node",
			declaration: "qa: declined\nqa_reason: no behavioral surface\n",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "QA", "pending", "qa", defaultVerificationSection),
			},
			wantTask: "task_01",
		},
		{
			name:        "gate must be unique",
			declaration: "qa: task_03\n",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01]
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "Build", "completed", "backend", defaultVerificationSection),
				"task_03.md": taskFixture("task_03", "QA", "pending", "qa", defaultVerificationSection),
				"task_04.md": taskFixture("task_04", "Other QA", "pending", "qa", defaultVerificationSection),
			},
			wantTask: "task_04",
		},
		{
			name:        "gate is not terminal",
			declaration: "qa: task_02\n",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "Build", "completed", "backend", defaultVerificationSection),
				"task_02.md": taskFixture("task_02", "QA", "pending", "qa", defaultVerificationSection),
				"task_03.md": taskFixture("task_03", "Append", "pending", "backend", defaultVerificationSection),
			},
			wantTask: "task_03",
		},
		{
			name:        "gate does not cover leaf",
			declaration: "qa: task_03\n",
			nodes: `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01]
`,
			tasks: map[string]string{
				"task_01.md": taskFixture("task_01", "Build A", "completed", "backend", defaultVerificationSection),
				"task_02.md": taskFixture("task_02", "Build B", "pending", "backend", defaultVerificationSection),
				"task_03.md": taskFixture("task_03", "QA", "pending", "qa", defaultVerificationSection),
			},
			wantTask: "task_02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			files := map[string]string{
				"_prd.md":   prdFixture("active"),
				"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", tt.declaration, tt.nodes, ""),
			}
			for name, content := range tt.tasks {
				files[name] = content
			}
			writeSpecDir(t, specsRoot, "demo", files)

			_, err := Load(specsRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want QA gate validation error")
			}
			var gateErr QAGateError
			if !errors.As(err, &gateErr) {
				t.Fatalf("error = %v, want QAGateError", err)
			}
			if !strings.Contains(err.Error(), "validate qa gate") {
				t.Fatalf("error %q does not identify QA gate validation", err)
			}
			if !strings.Contains(err.Error(), tt.wantTask) {
				t.Fatalf("error %q does not name %q", err, tt.wantTask)
			}
		})
	}
}

func TestLoadInvalidatesSettledQAGateAfterTaskAppend(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	specDir := filepath.Join(specsRoot, "demo")
	files := map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixtureWithQA("spec-tasks/v1", "qa: task_03\n", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
`, ""),
		"task_01.md": taskFixture("task_01", "Prepare", "completed", "backend", defaultVerificationSection),
		"task_02.md": taskFixture("task_02", "Build", "completed", "backend", defaultVerificationSection),
		"task_03.md": taskFixture("task_03", "QA", "completed", "qa", defaultVerificationSection),
	}
	writeSpecDir(t, specsRoot, "demo", files)
	if _, err := Load(specsRoot, "demo"); err != nil {
		t.Fatalf("Load before append: %v", err)
	}

	writeFile(t, filepath.Join(specDir, "_tasks.md"), manifestFixtureWithQA("spec-tasks/v1", "qa: task_03\n", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_03
      file: task_03.md
      needs: [task_04]
`, ""))
	writeFile(t, filepath.Join(specDir, "task_04.md"), taskFixture("task_04", "Correct", "pending", "backend", defaultVerificationSection))

	_, err := Load(specsRoot, "demo")
	if err == nil {
		t.Fatal("Load after append succeeded, want stale gate error")
	}
	var stale StaleGateError
	if !errors.As(err, &stale) {
		t.Fatalf("error = %v, want StaleGateError", err)
	}
	if stale.QATaskID != "task_03" || len(stale.TaskIDs) != 1 || stale.TaskIDs[0] != "task_04" {
		t.Fatalf("StaleGateError = %+v, want task_03 invalidated by task_04", stale)
	}
	for _, want := range []string{"task_04", "gate result is invalidated"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err, want)
		}
	}
}

func TestTaskTypeProjectionMustMatchTaskFile(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixtureWithBody("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`, `
| id      | title   | type     | complexity | needs |
| ------- | ------- | -------- | ---------- | ----- |
| task_01 | Fixture | frontend | low        | —     |
`),
		"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
	})

	_, err := Load(specsRoot, "demo")
	if err == nil {
		t.Fatal("Load succeeded, want Task Type projection mismatch")
	}
	var projectionErr TaskTypeProjectionError
	if !errors.As(err, &projectionErr) {
		t.Fatalf("error = %v, want TaskTypeProjectionError", err)
	}
	message := err.Error()
	for _, want := range []string{
		filepath.Join(specsRoot, "demo", "_tasks.md"),
		filepath.Join(specsRoot, "demo", "task_01.md"),
		"frontend",
		"backend",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("error %q does not contain %q", message, want)
		}
	}
}

func TestTaskTypeProjectionTablePresenceValidatesRows(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		table    string
		contains []string
	}{
		{
			name: "empty table still requires task rows",
			table: `
| id      | title   | type     | complexity | needs |
| ------- | ------- | -------- | ---------- | ----- |
`,
			contains: []string{"projection table has no row", "task_01"},
		},
		{
			name: "malformed task id rejected",
			table: `
| id      | title   | type     | complexity | needs |
| ------- | ------- | -------- | ---------- | ----- |
| taks_01 | Fixture | backend  | low        | —     |
`,
			contains: []string{"malformed Task row", "task_NN"},
		},
		{
			name: "duplicate task id rejected",
			table: `
| id      | title   | type     | complexity | needs |
| ------- | ------- | -------- | ---------- | ----- |
| task_01 | Fixture | backend  | low        | —     |
| task_01 | Fixture | backend  | low        | —     |
`,
			contains: []string{"defines Task \"task_01\" more than once"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			writeSpecDir(t, specsRoot, "demo", map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixtureWithBody("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`, tt.table),
				"task_01.md": taskFixture("task_01", "Fixture", "pending", "backend", defaultVerificationSection),
			})

			_, err := Load(specsRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want projection table validation error")
			}
			for _, want := range tt.contains {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error %q does not contain %q", err, want)
				}
			}
		})
	}
}

func TestLoadParsesOptionalTaskContext(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	contextSection := md(`## Context

- instruction: '.agents/skills/golang-testing/SKILL.md'
- interface: 'internal/spec/task.go'
- interface: 'internal/spec/task.go'

`)
	writeSpecDir(t, specsRoot, "demo", map[string]string{
		"_prd.md": prdFixture("active"),
		"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
`),
		"task_01.md": taskFixture("task_01", "No context", "pending", "backend", defaultVerificationSection),
		"task_02.md": taskFixture("task_02", "With context", "pending", "backend", contextSection+defaultVerificationSection),
	})

	graph, err := Load(specsRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(graph.Tasks[0].Context) != 0 {
		t.Fatalf("missing ## Context should leave no refs, got %+v", graph.Tasks[0].Context)
	}
	got := graph.Tasks[1].Context
	want := []TaskContextRef{
		{Kind: ContextKindInstruction, Path: ".agents/skills/golang-testing/SKILL.md"},
		{Kind: ContextKindInterface, Path: "internal/spec/task.go"},
	}
	if len(got) != len(want) {
		t.Fatalf("Context = %+v, want %+v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("Context[%d] = %+v, want %+v", index, got[index], want[index])
		}
	}
}

func TestLoadRejectsInvalidTaskContext(t *testing.T) {
	t.Parallel()
	tooMany := strings.Builder{}
	tooMany.WriteString("## Context\n\n")
	for index := 0; index < maxTaskContextRefs+1; index++ {
		tooMany.WriteString(fmt.Sprintf("- interface: 'internal/file_%02d.go'\n", index))
	}
	tests := []struct {
		name    string
		context string
		want    string
	}{
		{
			name:    "unknown label",
			context: "## Context\n\n- source: 'internal/spec/task.go'\n\n",
			want:    "expected",
		},
		{
			name:    "absolute path",
			context: "## Context\n\n- interface: '/tmp/outside.go'\n\n",
			want:    "repository-relative",
		},
		{
			name:    "escaping path",
			context: "## Context\n\n- interface: '../outside.go'\n\n",
			want:    "inside the repository",
		},
		{
			name:    "unclean path",
			context: "## Context\n\n- interface: 'internal/../task.go'\n\n",
			want:    "clean",
		},
		{
			name:    "too many unique entries",
			context: tooMany.String(),
			want:    "more than 50",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			writeSpecDir(t, specsRoot, "demo", map[string]string{
				"_prd.md": prdFixture("active"),
				"_tasks.md": manifestFixture("spec-tasks/v1", `    - id: task_01
      file: task_01.md
      needs: []
`),
				"task_01.md": taskFixture("task_01", "Invalid context", "pending", "backend", md(tt.context)+defaultVerificationSection),
			})

			_, err := Load(specsRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want Task Context validation error")
			}
			var taskErr TaskFileError
			if !errors.As(err, &taskErr) {
				t.Fatalf("error = %v, want TaskFileError", err)
			}
			var contextErr TaskContextError
			if !errors.As(err, &contextErr) {
				t.Fatalf("error = %v, want TaskContextError", err)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err, tt.want)
			}
		})
	}
}

func TestLoadReturnsTypedValidationErrors(t *testing.T) {
	t.Parallel()
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
				"task_01.md": taskFixture("task_01", "Fixture", "finished", "backend", defaultVerificationSection),
			},
			check: func(t *testing.T, err error) {
				var taskErr TaskFileError
				if !errors.As(err, &taskErr) {
					t.Fatalf("error = %v, want TaskFileError", err)
				}
				if !strings.Contains(err.Error(), `"finished"`) {
					t.Errorf("message %q does not name the unsupported status", err)
				}
				if !strings.Contains(err.Error(), "pending, in_progress, completed, failed") {
					t.Errorf("message %q does not name allowed statuses", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			if tt.files != nil {
				writeSpecDir(t, specsRoot, "demo", tt.files)
			}
			_, err := Load(specsRoot, "demo")
			if err == nil {
				t.Fatal("Load succeeded, want a typed validation error")
			}
			tt.check(t, err)
		})
	}
}

func TestLoadAcceptsValidGraph(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "demo", diamondSpecFiles())

	graph, err := Load(specsRoot, "demo")
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

func TestLoadUsesExplicitExternalSpecRoot(t *testing.T) {
	t.Parallel()
	specsRoot := filepath.Join(t.TempDir(), "external-specs")
	writeSpecDir(t, specsRoot, "demo", diamondSpecFiles())

	graph, err := Load(specsRoot, "demo")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if graph.Spec.Dir != filepath.Join(specsRoot, "demo") {
		t.Errorf("Spec.Dir = %q, want external Spec Root path", graph.Spec.Dir)
	}
	if len(graph.Tasks) == 0 {
		t.Fatal("Load returned no tasks")
	}
	if graph.Tasks[0].File != filepath.Join("demo", "task_01.md") {
		t.Errorf("Task.File = %q, want Spec Root-relative path", graph.Tasks[0].File)
	}

	active, err := ListActive(specsRoot)
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(active) != 1 || active[0].Dir != filepath.Join(specsRoot, "demo") {
		t.Fatalf("ListActive = %+v, want demo from external Spec Root", active)
	}
}

func TestQAGateLegacyArchivedManifestsLoadUnchanged(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	archivedRoot := filepath.Join(filepath.Dir(testFile), "..", "..", "docs", "specs", "_archived")
	tempSpecsRoot := defaultSpecsRoot(t.TempDir())
	entries, err := os.ReadDir(archivedRoot)
	if err != nil {
		t.Fatalf("read archived Spec root: %v", err)
	}

	manifestCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifestPath := filepath.Join(archivedRoot, entry.Name(), "_tasks.md")
		before, err := os.ReadFile(manifestPath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatalf("read archived manifest %q: %v", entry.Name(), err)
		}
		manifestCount++
		nodes, _, _, _, err := loadManifestNodes(manifestPath)
		if err != nil {
			t.Fatalf("parse archived manifest %q: %v", entry.Name(), err)
		}
		files := map[string]string{
			"_prd.md":   prdFixture("active"),
			"_tasks.md": string(before),
		}
		for _, node := range nodes {
			content, err := os.ReadFile(filepath.Join(archivedRoot, entry.Name(), node.File))
			if err != nil {
				t.Fatalf("read archived Task %q from Spec %q: %v", node.ID, entry.Name(), err)
			}
			files[node.File] = string(content)
		}
		writeSpecDir(t, tempSpecsRoot, entry.Name(), files)

		graph, err := Load(tempSpecsRoot, entry.Name())
		if err != nil {
			t.Fatalf("Load archived Spec %q: %v", entry.Name(), err)
		}
		if graph.QATaskID != "" || graph.QADeclined || graph.QAReason != "" {
			t.Fatalf("archived Spec %q gained QA declaration state: (%q, %t, %q)", entry.Name(), graph.QATaskID, graph.QADeclined, graph.QAReason)
		}
		for _, task := range graph.Tasks {
			if task.Type == TaskTypeQA {
				t.Fatalf("archived Spec %q unexpectedly has QA Task %q", entry.Name(), task.ID)
			}
		}

		copiedManifest := filepath.Join(tempSpecsRoot, entry.Name(), "_tasks.md")
		after, err := os.ReadFile(copiedManifest)
		if err != nil {
			t.Fatalf("re-read loaded manifest %q: %v", entry.Name(), err)
		}
		if string(after) != string(before) {
			t.Fatalf("Load changed archived manifest %q", entry.Name())
		}
	}
	if manifestCount == 0 {
		t.Fatal("no archived Task Graph manifests found")
	}
}

func TestListActiveFiltersInactiveArchivedAndNonSpecDirectories(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "0002-later", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, specsRoot, "0001-early", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, specsRoot, "0003-shipped", map[string]string{"_prd.md": prdFixture("shipped")})
	writeSpecDir(t, specsRoot, "_archived", map[string]string{"_prd.md": prdFixture("active")})
	writeSpecDir(t, specsRoot, "0004-broken", map[string]string{"_prd.md": "no frontmatter"})
	if err := os.MkdirAll(filepath.Join(specsRoot, "0005-no-prd"), 0o755); err != nil {
		t.Fatalf("create directory: %v", err)
	}
	writeFile(t, filepath.Join(specsRoot, "notes.md"), "# not a spec\n")

	specs, err := ListActive(specsRoot)
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
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	writeSpecDir(t, specsRoot, "0001-active", map[string]string{"_prd.md": prdFixture("active")})
	if err := os.MkdirAll(filepath.Join(specsRoot, "0002-missing-prd"), 0o755); err != nil {
		t.Fatalf("create missing PRD fixture: %v", err)
	}
	writeSpecDir(t, specsRoot, "0003-broken-frontmatter", map[string]string{"_prd.md": "no frontmatter\n"})
	writeSpecDir(t, specsRoot, "0004-archived-status", map[string]string{"_prd.md": prdFixture("archived")})
	writeSpecDir(t, specsRoot, "_archived", map[string]string{"0005-old/_prd.md": prdFixture("broken")})

	specs, skipped, err := ListActiveDetailed(specsRoot)
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

	simple, err := ListActive(specsRoot)
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(simple) != 1 || simple[0].Slug != specs[0].Slug {
		t.Fatalf("ListActive = %+v, want active list unchanged from detailed result %+v", simple, specs)
	}
}

func TestListActiveWithoutSpecsRootReturnsNothing(t *testing.T) {
	t.Parallel()
	specs, err := ListActive(t.TempDir())
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(specs) != 0 {
		t.Fatalf("ListActive = %+v, want none", specs)
	}
}
