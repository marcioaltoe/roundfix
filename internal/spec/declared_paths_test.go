package spec

import (
	"slices"
	"testing"
)

func TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths(t *testing.T) {
	t.Parallel()
	const content = `---
task: task_01
status: completed
type: backend
---

# Task 01: Example

## Context

- instruction: docs/agents/go.md
- creates: internal/spec/new.go
- interface: internal/cli/existing.go
- interface: internal/cli/existing.go
- deletes: obsolete.go

## Recorded paths

- ` + "`internal/cli/existing.go`" + `
- ` + "`./internal/spec/recorded.go`" + `
- ` + "`internal/spec/recorded.go`" + `
`
	paths, err := DeclaredTaskPaths([]byte(content), "docs/specs/example/task_01.md")
	want := []string{"internal/cli/existing.go", "internal/spec/new.go", "internal/spec/recorded.go"}
	if err != nil || !slices.Equal(paths, want) {
		t.Fatalf("declared paths=%v err=%v, want %v", paths, err, want)
	}
	for _, tc := range []struct{ name, content string }{
		{"invalid-frontmatter", "# No frontmatter"},
		{"invalid-context", "---\nstatus: completed\ntype: backend\n---\n## Context\n- interface: ../outside.go\n"},
		{"invalid-recorded-path", "---\nstatus: completed\ntype: backend\n---\n## Recorded paths\n- `../outside.go`\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if paths, err := DeclaredTaskPaths([]byte(tc.content), "task_01.md"); err == nil || len(paths) != 0 {
				t.Fatalf("declared paths=%v err=%v, want invalid Task refused", paths, err)
			}
		})
	}
}
