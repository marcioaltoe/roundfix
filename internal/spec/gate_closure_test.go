package spec

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

// Suite: manifest-byte closure and late-dependency invalidation.
// Boundary: recorded manifest bytes and atomic QA Task replacement.
func TestQAGateClosureFromManifestBytes(t *testing.T) {
	t.Parallel()
	const manifest = "---\nschema: spec-tasks/v1\nqa: task_qa\ngraph:\n  nodes:\n    - id: task_qa\n      file: task_qa.md\n      needs: [task_02, task_01]\n    - id: task_02\n      file: task_02.md\n      needs: [task_01]\n    - id: task_01\n      file: task_01.md\n---\n"
	id, closure, err := QAGateClosure("sample/_tasks.md", []byte(manifest))
	if err != nil || id != "task_qa" || !slices.Equal(closure, []string{"task_01", "task_02"}) {
		t.Fatalf("closure = %q %q %v", id, closure, err)
	}
	for _, tc := range []struct{ name, content string }{
		{"invalid YAML", "---\n[\n---\n"},
		{"unknown dependency", string(bytes.ReplaceAll([]byte(manifest), []byte("needs: [task_01]"), []byte("needs: [absent]")))},
		{"missing gate", string(bytes.ReplaceAll([]byte(manifest), []byte("qa: task_qa"), []byte("qa: absent")))},
		{"cycle", string(bytes.ReplaceAll([]byte(manifest), []byte("needs: [task_01]"), []byte("needs: [task_qa]")))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := QAGateClosure("sample/_tasks.md", []byte(tc.content)); err == nil {
				t.Fatal("want invalid manifest refusal")
			}
		})
	}
}

func TestReopenGateForLateDependenciesRecordsTheLine(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "task_qa.md")
	original := []byte("---\ntask: task_qa\nstatus: completed\ntype: qa\n---\n\n## Result\n\nPrior evidence.\n")
	if err := os.WriteFile(path, original, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := ReopenGateForLateDependencies(path, "qa/qa-report-2026-10-06.md", []string{"task_03", "task_04"}, time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := string(bytes.Replace(original, []byte("status: completed"), []byte("status: pending"), 1)) + "\n## Invalidation\n\n- Date: `2026-10-07`\n- QA Report: `qa/qa-report-2026-10-06.md`\n- Dependencies added after the QA Report: `task_03`, `task_04`\n"
	if string(got) != want {
		t.Fatalf("Task = %q, want %q", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode = %v", info.Mode())
	}
}
