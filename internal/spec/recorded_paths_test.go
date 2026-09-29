package spec

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestUndeclaredTaskPathsExcludesDeclaredGovernedAndTheTaskFile(t *testing.T) {
	t.Parallel()
	task := Task{Context: []TaskContextRef{
		{Kind: ContextKindInterface, Path: "internal/store/delivery.go"},
		{Kind: ContextKindCreates, Path: "internal/store/new.go"},
		{Kind: ContextKindInstruction, Path: "docs/agents/go.md"},
	}}
	committed := []string{
		"internal/store/delivery_test.go",
		"docs/specs/demo/task_01.md",
		"internal/store/new.go",
		"Makefile",
		"internal/store/delivery.go",
		"internal/store/delivery_test.go",
	}

	got := UndeclaredTaskPaths(task, "docs/specs/demo/task_01.md", committed, func(path string) bool {
		return path == "Makefile"
	})

	want := []string{"internal/store/delivery_test.go"}
	if !slices.Equal(got, want) {
		t.Fatalf("UndeclaredTaskPaths() = %v, want %v", got, want)
	}
}

func TestUndeclaredTaskPathsRecordsAnEditedInstructionPath(t *testing.T) {
	t.Parallel()
	task := Task{Context: []TaskContextRef{{Kind: ContextKindInstruction, Path: "docs/agents/go.md"}}}

	got := UndeclaredTaskPaths(task, "docs/specs/demo/task_01.md", []string{"docs/agents/go.md"}, func(string) bool { return false })

	if !slices.Equal(got, []string{"docs/agents/go.md"}) {
		t.Fatalf("UndeclaredTaskPaths() = %v, want the edited instruction path", got)
	}
}

func TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte(t *testing.T) {
	t.Parallel()
	taskPath := filepath.Join(t.TempDir(), "task_01.md")
	base := []byte("---\nstatus: completed\n---\n\n# Task\n\n## Result\n\nEvidence without a final newline")
	original := append(append([]byte(nil), base...), []byte("\n## Recorded paths\n\nold prose\n\n- `old.go`\n")...)
	if err := os.WriteFile(taskPath, original, 0o640); err != nil {
		t.Fatalf("write Task fixture: %v", err)
	}

	if err := RecordTaskPaths(taskPath, []string{"internal/store/delivery_test.go", "notes.txt"}); err != nil {
		t.Fatalf("RecordTaskPaths: %v", err)
	}

	updated, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("read updated Task: %v", err)
	}
	if !bytes.HasPrefix(updated, base) {
		t.Fatalf("RecordTaskPaths changed bytes before the section:\n%s", updated)
	}
	if bytes.Contains(updated, []byte("old prose")) || bytes.Contains(updated, []byte("old.go")) {
		t.Fatalf("RecordTaskPaths kept the replaced section:\n%s", updated)
	}
	wantPaths := []string{"internal/store/delivery_test.go", "notes.txt"}
	if got := RecordedTaskPaths(updated); !slices.Equal(got, wantPaths) {
		t.Fatalf("RecordedTaskPaths() = %v, want %v", got, wantPaths)
	}
	info, err := os.Stat(taskPath)
	if err != nil {
		t.Fatalf("stat updated Task: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o640 {
		t.Fatalf("updated Task mode = %o, want 640", got)
	}
}

func TestRecordTaskPathsWithNoPathsRemovesTheSection(t *testing.T) {
	t.Parallel()
	taskPath := filepath.Join(t.TempDir(), "task_01.md")
	base := []byte("---\nstatus: completed\n---\n\n# Task\n")
	original := append(append([]byte(nil), base...), []byte("\n## Recorded paths\n\nThe Daemon recorded these paths, which this Task changed without declaring them in `## Context`.\n\n- `old.go`\n")...)
	if err := os.WriteFile(taskPath, original, 0o644); err != nil {
		t.Fatalf("write Task fixture: %v", err)
	}

	if err := RecordTaskPaths(taskPath, nil); err != nil {
		t.Fatalf("RecordTaskPaths: %v", err)
	}

	updated, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("read updated Task: %v", err)
	}
	if !bytes.Equal(updated, base) {
		t.Fatalf("RecordTaskPaths with no paths wrote:\n%s\nwant exact base:\n%s", updated, base)
	}
	if got := RecordedTaskPaths(updated); len(got) != 0 {
		t.Fatalf("RecordedTaskPaths() = %v, want none", got)
	}
}

func TestRecordTaskPathsRefusesAnUnrecordablePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
	}{
		{name: "backtick", path: "bad`path.go"},
		{name: "carriage return", path: "bad\rpath.go"},
		{name: "line feed", path: "bad\npath.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			taskPath := filepath.Join(t.TempDir(), "task_01.md")
			original := []byte("---\nstatus: completed\n---\n\n# Task\n")
			if err := os.WriteFile(taskPath, original, 0o644); err != nil {
				t.Fatalf("write Task fixture: %v", err)
			}

			err := RecordTaskPaths(taskPath, []string{tt.path})

			if err == nil {
				t.Fatal("RecordTaskPaths succeeded, want an error")
			}
			if !strings.Contains(err.Error(), "cannot be recorded") {
				t.Fatalf("RecordTaskPaths error = %q, want an unrecordable-path error", err)
			}
			updated, readErr := os.ReadFile(taskPath)
			if readErr != nil {
				t.Fatalf("read Task after refusal: %v", readErr)
			}
			if !bytes.Equal(updated, original) {
				t.Fatalf("RecordTaskPaths mutated the Task on refusal:\n%s", updated)
			}
		})
	}
}

func TestARecordedSectionLeavesContextAndCarryForwardInputsUnchanged(t *testing.T) {
	t.Parallel()
	taskFile := "docs/specs/demo/task_01.md"
	base := taskFixture("task_01", "Record paths", "completed", "backend", `## Context

- instruction: 'docs/agents/go.md'
- interface: 'internal/store/delivery.go'
- creates: 'internal/store/delivery_test.go'

`+defaultVerificationSection)
	withRecord := base + `
## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in '## Context'.

- 'internal/store/extra_test.go'
`

	withoutDocument, err := parseTaskDocument([]byte(base), taskFile)
	if err != nil {
		t.Fatalf("parse Task without recorded paths: %v", err)
	}
	withDocument, err := parseTaskDocument([]byte(withRecord), taskFile)
	if err != nil {
		t.Fatalf("parse Task with recorded paths: %v", err)
	}
	if !reflect.DeepEqual(withDocument.Context, withoutDocument.Context) {
		t.Fatalf("Context with record = %#v, want %#v", withDocument.Context, withoutDocument.Context)
	}

	withoutInputs, err := CarryForwardInputs("docs/specs/demo", taskFile, []byte(base))
	if err != nil {
		t.Fatalf("CarryForwardInputs without record: %v", err)
	}
	withInputs, err := CarryForwardInputs("docs/specs/demo", taskFile, []byte(withRecord))
	if err != nil {
		t.Fatalf("CarryForwardInputs with record: %v", err)
	}
	if !slices.Equal(withInputs, withoutInputs) {
		t.Fatalf("CarryForwardInputs with record = %v, want %v", withInputs, withoutInputs)
	}
}
