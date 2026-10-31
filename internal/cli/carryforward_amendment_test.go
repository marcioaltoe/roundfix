// Suite: Task Carry-Forward across planning amendments.
// Invariant: unchanged Binding Text carries; changed contracts refuse the whole set.
// Boundary IN: the Reconcile Command, real Git staging and Run settlement evidence.
// Boundary OUT: Daemon Verification execution and delivery queue orchestration.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText(t *testing.T) {
	t.Parallel()
	fixture := newPlanningAmendmentFixture(t, []implementSeed{{id: "task_01"}}, "")
	dir := filepath.Join(fixture.repoDir, "docs", "specs", implementTestSlug)
	for _, file := range []string{"_prd.md", "_techspec.md"} {
		path := filepath.Join(dir, file)
		mustWrite(t, path, strings.Replace(mustRead(t, path), "## Project Constraints\n", "## Project Constraints\n\n- Tooling authority: corrective Task authorized.\n", 1))
	}
	path := filepath.Join(dir, "_tasks.md")
	manifest := strings.Replace(mustRead(t, path), "      file: task_01.md\n", "      file: task_01.md\n    - id: task_02\n      file: task_02.md\n      needs: [task_01]\n", 1)
	mustWrite(t, path, manifest+"| task_02 | Corrective work | backend | low | task_01 |\n")
	mustWrite(t, filepath.Join(dir, "task_02.md"), implementTaskContent(implementTestSlug, implementSeed{id: "task_02"}))
	gitImplement(t, fixture.repoDir, "add", "docs/specs")
	gitImplement(t, fixture.repoDir, "commit", "-m", "docs: amend constraints and add corrective Task")
	amended := []string{}
	for _, file := range []string{"_prd.md", "_tasks.md", "_techspec.md"} {
		amended = append(amended, filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, file)))
	}

	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--format=json"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("inspect exit=%d stderr=%s stdout=%s", code, &stderr, &stdout)
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.CarryForwards) != 1 {
		t.Fatalf("candidates=%+v", report.CarryForwards)
	}
	candidate := report.CarryForwards[0]
	if candidate.Action != carryForwardReadyAction || candidate.InputsMoved || len(candidate.MovedInputs) != 0 {
		t.Fatalf("candidate=%+v", candidate)
	}
	var wire struct {
		CarryForwards []struct {
			AmendedInputs []string `json:"amendedInputs"`
		} `json:"carryForwards"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &wire); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(wire.CarryForwards[0].AmendedInputs, amended) {
		t.Fatalf("amendedInputs=%v, want %v", wire.CarryForwards[0].AmendedInputs, amended)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--carry-forward"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("carry exit=%d stderr=%s stdout=%s", code, &stderr, &stdout)
	}
	wantLines := fmt.Sprintf("  moved-inputs: []\n  amended-inputs: %v\n", amended)
	if !strings.Contains(stdout.String(), wantLines) {
		t.Fatalf("text missing %q: %s", wantLines, &stdout)
	}
	task := mustRead(t, implementTaskPath(fixture.repoDir, "task_01"))
	for _, want := range []string{"status: completed", "## Carry-forward provenance", fixture.run.ID, fixture.commits["task_01"]} {
		if !strings.Contains(task, want) {
			t.Fatalf("carried Task missing %q: %s", want, task)
		}
	}
	if got := mustRead(t, filepath.Join(fixture.repoDir, "src", "task_01.txt")); got != "task_01 settled\n" {
		t.Fatalf("implementation=%q", got)
	}
	if !strings.Contains(mustRead(t, path), "task_02") || !strings.Contains(mustRead(t, filepath.Join(dir, "_prd.md")), "corrective Task authorized") {
		t.Fatal("carry lost the amendment")
	}
}

func TestCarryForwardRefusesATaskWhoseBindingTextChanged(t *testing.T) {
	t.Parallel()
	specDir := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug))
	for _, tc := range []struct{ name, path, context string }{
		{name: "named section", path: specDir + "/_techspec.md"},
		{name: "Task file", path: specDir + "/task_01.md"},
		{name: "root instructions", path: "AGENTS.md"},
		{name: "domain context", path: "CONTEXT.md"},
		{name: "implementation skill", path: ".agents/skills/implement-task/SKILL.md"},
		{name: "skill dispatch guide", path: "docs/agents/skill-dispatch.md"},
		{name: "Context path", path: "shared.txt"},
		{name: "explicit planning Context path", path: specDir + "/_prd.md", context: specDir + "/_prd.md"},
		{name: "another Spec planning file", path: "docs/specs/other/_prd.md", context: "docs/specs/other/_prd.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture := newPlanningAmendmentFixture(t, []implementSeed{{id: "task_01"}, {id: "task_02"}}, tc.context)
			path := filepath.Join(fixture.repoDir, filepath.FromSlash(tc.path))
			content := mustRead(t, path)
			if tc.name == "named section" {
				content = strings.Replace(content, "Keep the contract.", "Change the contract.", 1)
			} else if tc.name == "explicit planning Context path" {
				content = strings.Replace(content, "## Project Constraints\n", "## Project Constraints\n\nAmended authority.\n", 1)
			} else {
				content += "\nAmended input.\n"
			}
			mustWrite(t, path, content)
			gitImplement(t, fixture.repoDir, "add", tc.path)
			gitImplement(t, fixture.repoDir, "commit", "-m", "docs: amend a bound input")
			head := itemRecoveryHead(t, fixture.repoDir)
			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"}, &stdout, &stderr)
			if code != exitPreflight {
				t.Fatalf("exit=%d stderr=%s stdout=%s", code, &stderr, &stdout)
			}
			var report reconcileReport
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if len(report.CarryForwards) != 2 {
				t.Fatalf("candidates=%+v", report.CarryForwards)
			}
			want := "Task task_01 declared input(s) moved: " + tc.path
			if report.CarryForwards[0].RefusalReason != want {
				t.Fatalf("refusal=%q, want %q", report.CarryForwards[0].RefusalReason, want)
			}
			if !strings.Contains(stderr.String(), "carry-forward refused the whole set") {
				t.Fatalf("stderr=%s", &stderr)
			}
			if itemRecoveryHead(t, fixture.repoDir) != head {
				t.Fatal("refused set moved HEAD")
			}
			for _, id := range []string{"task_01", "task_02"} {
				if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", id+".txt")); !os.IsNotExist(err) {
					t.Fatalf("refused %s implementation stat=%v", id, err)
				}
			}
			if !bytes.Contains(stdout.Bytes(), []byte(`"amendedInputs":[]`)) {
				t.Fatalf("empty amendedInputs missing: %s", &stdout)
			}
		})
	}
}

func newPlanningAmendmentFixture(t *testing.T, seeds []implementSeed, extraContext string) carryForwardFixture {
	t.Helper()
	home, repo := newImplementWorkspace(t, seeds)
	dir := filepath.Join(repo, "docs", "specs", implementTestSlug)
	mustWrite(t, filepath.Join(dir, "_prd.md"), mustRead(t, filepath.Join(dir, "_prd.md"))+"\n## Goals\n\nKeep the goal.\n")
	mustWrite(t, filepath.Join(dir, "_techspec.md"), "# TechSpec\n\n## Project Constraints\n\nExisting authority.\n\n## API Contracts\n\nKeep the contract.\n")
	manifest := mustRead(t, filepath.Join(dir, "_tasks.md")) + "\n| id | title | type | complexity | needs |\n| --- | --- | --- | --- | --- |\n"
	for _, seed := range seeds {
		manifest += fmt.Sprintf("| %s | Build | backend | low | — |\n", seed.id)
		path := implementTaskPath(repo, seed.id)
		task := mustRead(t, path) + "\n## Context\n\n- interface: `shared.txt`\n"
		if extraContext != "" {
			task += "- interface: `" + extraContext + "`\n"
		}
		task += "\n## References\n\n- `_prd.md` → Goal 1\n- `_techspec.md` → API Contract 6\n"
		mustWrite(t, path, task)
	}
	mustWrite(t, filepath.Join(dir, "_tasks.md"), manifest)
	for _, path := range []string{"AGENTS.md", "CONTEXT.md", ".agents/skills/implement-task/SKILL.md", "docs/agents/skill-dispatch.md", "shared.txt", "docs/specs/other/_prd.md"} {
		mustMkdir(t, filepath.Dir(filepath.Join(repo, path)))
		mustWrite(t, filepath.Join(repo, path), "Unchanged input.\n")
	}
	gitImplement(t, repo, "add", "-A")
	gitImplement(t, repo, "commit", "-m", "docs: declare settled Task binding")
	workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
	fixture := createDeliveryRetryRun(t, workflow, repo, "ma/widget-flow", seeds, store.StateUnresolved)
	fixture.homeDir = home
	return fixture
}
