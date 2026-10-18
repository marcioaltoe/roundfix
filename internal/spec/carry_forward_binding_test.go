// Suite: Task Binding Text in planning files.
// Invariant: only the named sections and the Task's own graph declaration bind it.
// Boundary IN: committed planning and Task bytes.
// Boundary OUT: Git drift detection and carry-forward mutation.
package spec

import (
	"strings"
	"testing"
)

func TestPlanningBindingReadsTheTaskNodeAndRow(t *testing.T) {
	t.Parallel()
	manifest := "---\nschema: spec-tasks/v1\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: [task_02, task_03]\n    - id: task_02\n      file: task_02.md\n---\n\n# Tasks\n\n| id | title | type | complexity | needs |\n| --- | --- | --- | --- | --- |\n| task_01 | Build | backend | high | task_02, task_03 |\n| task_02 | Other | backend | low | — |\n"
	binding, whole, err := PlanningBinding("_tasks.md", []byte(manifest), "task_01", nil)
	if err != nil || whole {
		t.Fatalf("binding = %q, whole=%v, err=%v", binding, whole, err)
	}
	for _, tc := range []struct {
		name, content  string
		changed, whole bool
	}{
		{name: "another node and row", content: strings.ReplaceAll(manifest, "task_02.md", "other.md")},
		{name: "another row title", content: strings.ReplaceAll(manifest, "| Other |", "| Corrected |")},
		{name: "prose", content: manifest + "\nAmended prose.\n"},
		{name: "node file", content: strings.ReplaceAll(manifest, "task_01.md", "renamed.md"), changed: true},
		{name: "dependency order", content: strings.ReplaceAll(manifest, "[task_02, task_03]", "[task_03, task_02]"), changed: true},
		{name: "own row", content: strings.ReplaceAll(manifest, "| Build |", "| Changed |"), changed: true},
		{name: "missing node", content: strings.ReplaceAll(manifest, "id: task_01", "id: task_04"), whole: true},
		{name: "fenced row", content: strings.ReplaceAll(manifest, "| task_01 | Build | backend | high | task_02, task_03 |\n", "~~~markdown\n| task_01 | Build | backend | high | task_02, task_03 |\n~~~\n"), whole: true},
		{name: "missing row", content: strings.ReplaceAll(manifest, "| task_01 | Build | backend | high | task_02, task_03 |\n", ""), whole: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, full, err := PlanningBinding("_tasks.md", []byte(tc.content), "task_01", nil)
			if err != nil || full != tc.whole {
				t.Fatalf("whole=%v err=%v, want whole=%v", full, err, tc.whole)
			}
			if !full && (got != binding) != tc.changed {
				t.Fatalf("binding changed=%v, want %v; got %q vs %q", got != binding, tc.changed, got, binding)
			}
		})
	}
}

func TestPlanningBindingReadsTheSectionsTheReferencesName(t *testing.T) {
	t.Parallel()
	planning := "# Plan\n\n## Project Constraints\nTooling.\n\n## User Stories\nStory.\n\n### API Contracts\nContract.\n\n#### Invariants\nRule.\n\n##### Detail\nStill in invariant.\n\n## Goals\nGoal.\n\n~~~markdown\n## Phantom\n~~~\n\n## Goals\nSecond goal.\n\n## Decisions\nChoice.\n"
	stories := planning[strings.Index(planning, "## User Stories"):strings.Index(planning, "## Goals")]
	goals := planning[strings.Index(planning, "## Goals"):strings.Index(planning, "## Decisions")]
	contracts := planning[strings.Index(planning, "### API Contracts"):strings.Index(planning, "## Goals")]
	for _, tc := range []struct {
		name, refs, want string
		whole            bool
	}{
		{name: "plural and numbered units", refs: "User Story 1-3, 5 and 7; API Contract 6; Invariant 2", want: stories + contracts + planning[strings.Index(planning, "#### Invariants"):strings.Index(planning, "## Goals")]},
		{name: "colon and case", refs: "api contract 6: details", want: contracts},
		{name: "file order and duplicate units", refs: "Goal 1; User Story 4; Goals", want: stories + goals},
		{name: "fenced heading", refs: "Phantom", whole: true},
		{name: "level one is not a binding", refs: "Plan", whole: true},
		{name: "level five is not a binding", refs: "Detail", whole: true},
		{name: "unicode ranges", refs: "User Story 1–3, 5 and 7—9", want: stories},
		{name: "unresolved unit", refs: "Goals; Absent", whole: true},
		{name: "missing bullet", whole: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, file := range []string{"_prd.md", "_techspec.md"} {
				task := "## References\n"
				if tc.refs != "" {
					task += "- `" + file + "` → " + tc.refs + "\n"
				}
				got, full, err := PlanningBinding(file, []byte(planning), "task_01", []byte(task))
				if err != nil || full != tc.whole || (!full && got != tc.want) {
					t.Fatalf("%s binding=%q whole=%v err=%v; want %q whole=%v", file, got, full, err, tc.want, tc.whole)
				}
			}
		})
	}
	t.Run("a fenced References bullet does not bind", func(t *testing.T) {
		t.Parallel()
		task := "## References\n\n~~~markdown\n- `_prd.md` → Goals\n~~~\n"
		_, whole, err := PlanningBinding("_prd.md", []byte(planning), "task_01", []byte(task))
		if err != nil || !whole {
			t.Fatalf("whole=%v err=%v, want whole", whole, err)
		}
	})
	t.Run("a wrapped References bullet binds and excludes later sections", func(t *testing.T) {
		t.Parallel()
		task := "## References\n\n- `_prd.md` → User Story 1;\n  API Contract 6\n\n## Other\n- `_prd.md` → Absent\n"
		got, whole, err := PlanningBinding("_prd.md", []byte(planning), "task_01", []byte(task))
		if err != nil || whole || got != stories+contracts {
			t.Fatalf("binding=%q whole=%v err=%v", got, whole, err)
		}
	})
	t.Run("Git text line endings do not change the binding", func(t *testing.T) {
		t.Parallel()
		task := []byte("## References\n- `_prd.md` → Goals\n")
		got, whole, err := PlanningBinding("_prd.md", []byte(strings.ReplaceAll(planning, "\n", "\r\n")), "task_01", task)
		if err != nil || whole || got != goals {
			t.Fatalf("binding=%q whole=%v err=%v", got, whole, err)
		}
	})

}
