package runcause

import (
	"reflect"
	"testing"

	"roundfix/internal/spec"
)

func TestCorrectiveTasksListsTheTasksAboveTheQAGate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, qa string
		want     []string
	}{
		{"above gate", "task_03", []string{"task_04", "task_10"}},
		{"no gate", "", nil}, {"invalid gate", "qa", nil}, {"last task", "task_10", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			graph := spec.CauseGraph{QATaskID: tc.qa, Tasks: map[string]string{"task_01": "task_01.md", "task_03": "task_03.md", "task_04": "task_04.md", "task_10": "task_10.md", "unrelated": "unrelated.md"}}
			if got := CorrectiveTasks(graph); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("CorrectiveTasks = %v, want %v", got, tc.want)
			}
		})
	}
}
