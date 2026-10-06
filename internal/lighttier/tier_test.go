package lighttier

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/spec"
)

func TestTierFor(t *testing.T) {
	for _, tc := range []struct {
		name, complexity, taskType, refs string
		want                             Tier
	}{
		{"low", "low", "backend", "", Light},
		{"trimmed", " low ", "docs", "", Light},
		{"qa", "low", "qa", "", Standard},
		{"medium", "medium", "backend", "", Standard},
		{"high", "high", "backend", "", Standard},
		{"missing", "", "backend", "", Standard},
		{"unknown", "LOW", "backend", "", Standard},
		{"governed creates", "low", "backend", "- creates: `go.mod`", Standard},
		{"governed interface", "low", "backend", "- interface: `Makefile`", Standard},
		{"governed deletes", "low", "backend", "- deletes: `.golangci.yml`", Standard},
		{"instruction only", "low", "backend", "- instruction: `Makefile`", Light},
		{"ordinary writes", "low", "backend", "- creates: `internal/example/new.go`\n- interface: `internal/example/edit.go`\n- deletes: `internal/example/old.go`", Light},
		{"instruction and governed write", "low", "backend", "- instruction: `Makefile`\n- interface: `go.sum`", Standard},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "task_01.md")
			complexity := ""
			if tc.complexity != "" {
				complexity = fmt.Sprintf("complexity: %q\n", tc.complexity)
			}
			content := fmt.Sprintf("---\ntask: task_01\nspec: demo\nstatus: pending\ntype: %s\n%s---\n\n# Task 01: Fixture\n\n## Context\n\n%s\n\n## Verification\n\n- `true` — expected: exit 0\n", tc.taskType, complexity, tc.refs)
			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			task := spec.Task{ID: "task_01", File: "task_01.md"}
			if err := spec.ReloadTask(root, &task); err != nil {
				t.Fatal(err)
			}
			if got := TierFor(task); got != tc.want {
				t.Fatalf("TierFor=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestPlanEnabled(t *testing.T) {
	for _, tc := range []struct {
		name string
		plan Plan
		want bool
	}{
		{"zero", Plan{}, false},
		{"empty models", Plan{Models: []string{}, CeilingUSD: 10}, false},
		{"models", Plan{Models: []string{"deepseek/deepseek-v4.1-flash"}}, true},
		{"key absent still allows skip diagnostic", Plan{Models: []string{"deepseek/deepseek-v4.1-flash"}, KeyPresent: false}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.plan.Enabled(); got != tc.want {
				t.Fatalf("Enabled=%t want=%t", got, tc.want)
			}
		})
	}
}
