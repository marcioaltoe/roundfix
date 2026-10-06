// Suite: advisory model tier planning
// Invariant: only authored, bounded, graph-listed non-QA Tasks reach the judge.
// Boundary IN: PlanSpec with disposable Spec files
// Boundary OUT: service transport and dispatch
package judge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestModelTierPlansNonQATasks(t *testing.T) {
	q := loadQuestions(t)
	dir := t.TempDir()
	write := func(name, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("_tasks.md", "---\ngraph:\n  nodes:\n    - file: task_02.md\n    - file: task_01.md\n    - file: task_03.md\n---\n")
	for _, name := range []string{"task_01.md", "task_02.md", "task_unlisted.md"} {
		write(name, "---\ntask: "+name+"\nstatus: STATUS_SENTINEL\ncomplexity: COMPLEXITY_SENTINEL\ntype: backend\n---\n# Task\n\n## Requirements\nThe author reads the evidence for the work and keeps the rule.\n## Result\nRESULT_SENTINEL\n### Details\nRESULT_CHILD_SENTINEL\n## Recorded paths\nRECORDED_SENTINEL\n## Carry-forward provenance\nPROVENANCE_SENTINEL\n## Verification\n- `test-command`\n")
	}
	write("task_03.md", "---\ntype: qa\n---\n# QA gate\nThe author reads the evidence for the work.\n")
	// Missing PRD and TechSpec must not affect the tasks stage.
	plan, err := PlanSpec(q, dir, dir, "tasks")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Pending) != 2 || len(plan.Skipped) != 0 || len(plan.SkippedArtifacts) != 0 {
		t.Fatalf("plan=%+v", plan)
	}
	for i, p := range plan.Pending {
		want := []string{"task_02.md", "task_01.md"}[i]
		if p.Kind != "model-tier" || filepath.Base(p.Artifact) != want {
			t.Fatalf("judgment=%+v", p)
		}
		var state struct {
			TaskFile string `json:"task_file"`
			Task     string `json:"task"`
		}
		if err := json.Unmarshal(p.state, &state); err != nil {
			t.Fatal(err)
		}
		if state.TaskFile != want {
			t.Fatalf("Task identity lost: %s", p.state)
		}
		for _, sentinel := range []string{"STATUS_SENTINEL", "COMPLEXITY_SENTINEL", "RESULT_SENTINEL", "RESULT_CHILD_SENTINEL", "RECORDED_SENTINEL", "PROVENANCE_SENTINEL"} {
			if strings.Contains(state.Task, sentinel) {
				t.Errorf("leaked %s", sentinel)
			}
		}
		if !strings.Contains(state.Task, "test-command") || !strings.Contains(state.Task, "type: backend") {
			t.Fatalf("authored text lost: %s", state.Task)
		}
	}
	q.ModelTier.TaskMaxChars = 60
	plan, err = PlanSpec(q, dir, dir, "tasks")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range plan.Pending {
		var state struct{ Task string }
		if err := json.Unmarshal(p.state, &state); err != nil {
			t.Fatal(err)
		}
		if utf8.RuneCountInString(state.Task) != 60 {
			t.Fatalf("unbounded state: %s", p.state)
		}
	}
}

func TestModelTierRequiresTaskGraph(t *testing.T) {
	q := loadQuestions(t)
	if _, err := PlanSpec(q, t.TempDir(), t.TempDir(), "tasks"); err == nil {
		t.Fatal("missing graph accepted")
	}
}

func TestModelTierChoicePolicy(t *testing.T) {
	q := loadQuestions(t)
	if q.ModelTier.Question.Type != "choice" || q.ModelTier.TaskMaxChars <= 0 || len(q.ModelTier.Question.Criteria) != 3 {
		t.Fatalf("catalog=%+v", q.ModelTier)
	}
	for _, tier := range []string{"light", "standard", "heavy"} {
		confidence := 0.2
		a := answer{Type: "choice", Choice: &tier, Confidence: &confidence, Probabilities: map[string]float64{"light": 0.3, "standard": 0.3, "heavy": 0.4}}
		c := call{Status: 200, Model: q.Transports[1].RequestModel, Answers: map[string]answer{q.ModelTier.QuestionID: a}}
		_, outcome, reason, _ := evaluate(q, PendingJudgment{Kind: "model-tier"}, c)
		if outcome != "suggested" || reason != "" {
			t.Fatalf("%s: %s %s", tier, outcome, reason)
		}
		delete(a.Probabilities, "heavy")
		_, outcome, reason, _ = evaluate(q, PendingJudgment{Kind: "model-tier"}, c)
		if outcome != "skipped" || reason != "unreadable answer" {
			t.Fatalf("invalid choice: %s %s", outcome, reason)
		}
	}
}
