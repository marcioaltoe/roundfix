// Suite: shared Task work selection
// Invariant: every canonical status but completed awaits work, for backend and QA Tasks alike.
// Boundary IN: TaskAwaitsWork with parsed Task status and type.
// Boundary OUT: scheduling, QA routing, and CLI reporting, covered by their caller suites.
package daemon

import (
	"testing"

	"roundfix/internal/spec"
)

func TestTaskAwaitsWorkSelectsEveryStatusButCompleted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		status spec.Status
		want   bool
	}{
		{spec.StatusPending, true},
		{spec.StatusInProgress, true},
		{spec.StatusFailed, true},
		{spec.StatusCompleted, false},
	}
	for _, taskType := range []spec.TaskType{"backend", spec.TaskTypeQA} {
		for _, tc := range tests {
			t.Run(string(taskType)+"/"+string(tc.status), func(t *testing.T) {
				t.Parallel()
				if got := TaskAwaitsWork(spec.Task{Type: taskType, Status: tc.status}); got != tc.want {
					t.Fatalf("TaskAwaitsWork(%s, %s) = %t, want %t", taskType, tc.status, got, tc.want)
				}
			})
		}
	}
}
