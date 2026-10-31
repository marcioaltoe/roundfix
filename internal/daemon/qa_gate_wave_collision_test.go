// Suite: a corrective Task after a completed QA gate.
// Invariant: gate, strict check, settlement and dispatch judge the current graph.
// Boundary IN: the repository's own files and local Git history on every ref.
// All files and refs are built with gittest in temporary repositories; this
// never reads real ~/.roundfix, another repository or the network, and never sleeps.
// Boundary OUT: executing Verification or an Agent and writing a QA Report.
package daemon

import (
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, need string
		settled    bool
		colliding  bool
	}{
		{name: "corrective Task beside a completed Task", need: "task_01"},
		{name: "two unordered corrective Tasks", need: "task_01", colliding: true},
		{name: "ordered corrective Task with another Run settlement", need: "task_02", settled: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			const slug = "corrective-wave"
			gittest.InitRepo(t, root, "--initial-branch=main")
			seeds := []taskSpecSeed{
				{id: "task_01"},
				{id: "task_02", needs: []string{"task_01"}},
				{id: "task_03", needs: []string{"task_02"}},
				{id: "task_04", needs: []string{"task_03"}, taskType: "qa"},
			}
			writeSpecDirForTest(t, root, slug, seeds)
			specsRoot := filepath.Join(root, "docs", "specs")
			dir := filepath.Join(specsRoot, slug)
			writeSpecDirAtRootForTestWithQA(t, specsRoot, slug, seeds, qaDeclarationForTest{taskID: "task_04"})
			if err := os.MkdirAll(filepath.Join(root, "pkg"), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, filepath.Join(root, "pkg", "versions.test.ts"), "export const version = 0;\n")
			gittest.Run(t, root, "add", ".")
			gittest.Run(t, root, "commit", "-m", "docs: plan")
			gittest.Run(t, root, "checkout", "-b", "roundfix/run-run_first")
			for _, id := range []string{"task_01", "task_02", "task_03"} {
				if id == "task_02" {
					// The declaration omits this path: only the prior Run proves the touch.
					mustWriteForTest(t, filepath.Join(root, "pkg", "versions.test.ts"), "export const version = 1;\n")
				}
				if err := spec.SetStatus(filepath.Join(dir, id+".md"), spec.StatusCompleted); err != nil {
					t.Fatal(err)
				}
				gittest.Run(t, root, "add", ".")
				gittest.Run(t, root, "commit", "-m", "feat: settle "+id, "-m", "Roundfix-Spec: "+slug+"\nRoundfix-Task: "+id)
			}
			if err := spec.SetStatus(filepath.Join(dir, "task_04.md"), spec.StatusCompleted); err != nil {
				t.Fatal(err)
			}
			gittest.Run(t, root, "add", ".")
			gittest.Run(t, root, "commit", "-m", "test: settle QA gate")
			gittest.Run(t, root, "checkout", "main")
			gittest.Run(t, root, "merge", "--ff-only", "roundfix/run-run_first")
			gittest.Run(t, root, "rm", "-r", "docs/specs/"+slug)
			gittest.Run(t, root, "commit", "-m", "docs: archive")
			gittest.Run(t, root, "revert", "--no-edit", "HEAD")

			for i := range seeds {
				seeds[i].status = "completed"
			}
			seeds = append(seeds, taskSpecSeed{id: "task_05", needs: []string{tc.need}})
			seeds[3].needs = []string{"task_03", "task_05"}
			if tc.colliding {
				seeds = append(seeds, taskSpecSeed{id: "task_06", needs: []string{"task_01"}})
				seeds[3].needs = append(seeds[3].needs, "task_06")
			}
			// Adding a corrective dependency reopens the previously completed gate.
			seeds[3].status = "pending"
			writeSpecDirAtRootForTestWithQA(t, specsRoot, slug, seeds, qaDeclarationForTest{taskID: "task_04"})
			taskPath := filepath.Join(dir, "task_05.md")
			for _, seed := range seeds[4:] {
				path := filepath.Join(dir, seed.id+".md")
				body, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				mustWriteForTest(t, path, string(body)+"\n## Context\n\n- interface: `pkg/versions.test.ts`\n\n## Consumers\n\nNone. This fixture changes a test file with no shared contract.\n")
			}
			gittest.Run(t, root, "add", ".")
			gittest.Run(t, root, "commit", "-m", "fix: declare corrective Task")
			if tc.settled {
				gittest.Run(t, root, "checkout", "-b", "roundfix/run-run_second")
				mustWriteForTest(t, filepath.Join(root, "pkg", "versions.test.ts"), "export const version = 2;\n")
				if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
					t.Fatal(err)
				}
				gittest.Run(t, root, "add", ".")
				gittest.Run(t, root, "commit", "-m", "fix: settle corrective Task", "-m", "Roundfix-Spec: "+slug+"\nRoundfix-Task: task_05")
				// Keep this settlement only on its Run Branch: every reader needs --all.
				gittest.Run(t, root, "checkout", "main")
			}
			linked := filepath.Join(t.TempDir(), "run")
			gittest.Run(t, root, "worktree", "add", "--detach", linked, "HEAD")
			t.Cleanup(func() { gittest.Run(t, root, "worktree", "remove", "--force", linked) })
			graph, err := spec.Load(specsRoot, slug)
			if err != nil {
				t.Fatal(err)
			}
			linkedRoot := filepath.Join(linked, "docs", "specs")
			linkedGraph, err := spec.Load(linkedRoot, slug)
			if err != nil {
				t.Fatal(err)
			}
			gate, err := qaGatePrecondition(TaskPlan{WorkDir: linked, SpecsRoot: linkedRoot, Spec: linkedGraph.Spec, Tasks: linkedGraph.Tasks})
			if err != nil {
				t.Fatal(err)
			}
			checked, err := speccheck.Check(specsRoot, root, slug)
			if err != nil {
				t.Fatal(err)
			}
			speccheck.PromoteGaps(&checked)
			strict := speccheck.GatePrecondition(checked)
			settlement, err := (SpecCheckSettlementChecker{}).RefusingFindings(specsRoot, root, slug)
			if err != nil {
				t.Fatal(err)
			}
			const reason = "SC-WAVE-COLLISION: Task task_05 and Task task_06 can run in the same Wave and share repository files: pkg/versions.test.ts (declared context); add a `needs` edge from Task task_06 to Task task_05"
			for _, reader := range []struct {
				name     string
				findings []speccheck.Finding
			}{
				{"gate", gate.Findings}, {"strict", strict.Findings}, {"settlement", settlement},
			} {
				if !tc.colliding {
					if len(reader.findings) != 0 {
						t.Fatalf("%s findings = %+v, want none", reader.name, reader.findings)
					}
				} else if len(reader.findings) != 1 || speccheck.RefusalReason(reader.findings[0]) != reason {
					t.Fatalf("%s findings = %+v, want only %q", reader.name, reader.findings, reason)
				}
			}
			refusal, refused := speccheck.PreconditionRefusal(gate)
			if tc.colliding && (!refused || refusal.Reason != reason) {
				t.Fatalf("refusal = %+v, %v", refusal, refused)
			}
			if !tc.colliding && (refused || gate.Blocking || strict.Blocking) {
				t.Fatalf("non-colliding graph refused: gate %+v, strict %+v", gate, strict)
			}
			dispatch := refuseTaskPlanWaveCollisions(TaskPlan{WorkDir: root, SpecsRoot: specsRoot, Spec: graph.Spec, Tasks: graph.Tasks, Concurrency: 2})
			if !tc.colliding {
				if dispatch != nil {
					t.Fatal(dispatch)
				}
			} else {
				const want = "task cycle: refuse Wave before dispatch because Tasks collide: Task task_05 and Task task_06 share repository files: pkg/versions.test.ts (declared context); add a `needs` edge from Task task_06 to Task task_05"
				if dispatch == nil || dispatch.Error() != want {
					t.Fatalf("dispatch = %v, want %s", dispatch, want)
				}
			}
			t.Logf("four readers agree with task_05 needs [%s]", tc.need)
		})
	}
}
