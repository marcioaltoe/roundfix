// Suite: Spec Check resumed Task selection
// Invariant: only Tasks a Run still executes are probed; every omitted command is named without affecting refusal or provenance.
// Boundary IN: CLI dispatch, temporary Git repositories, real command probes, and text/JSON reports.
// Boundary OUT: Daemon scheduling and probe classification, covered by their own suites.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

type specCheckStatusSeed struct {
	status   spec.Status
	taskType spec.TaskType
	commands []string
}

type specCheckStatusDocument struct {
	Verification struct {
		Commands  []specCheckVerificationCommandReport `json:"commands"`
		NotProbed []struct {
			Task    string `json:"task"`
			Command string `json:"command"`
			Reason  string `json:"reason"`
		} `json:"notProbed"`
	} `json:"verification"`
}

func newSpecCheckStatusWorkspace(t *testing.T, seeds []specCheckStatusSeed) (string, string) {
	t.Helper()
	home, repo := newSpecCheckWorkspace(t, "clean")
	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\nspec: clean\n")
	for i, seed := range seeds {
		if seed.taskType == spec.TaskTypeQA {
			fmt.Fprintf(&manifest, "qa: task_%02d\n", i+1)
		}
	}
	manifest.WriteString("graph:\n  nodes:\n")
	for i, seed := range seeds {
		id := fmt.Sprintf("task_%02d", i+1)
		var needs []string
		if seed.taskType == spec.TaskTypeQA {
			for j := 0; j < i; j++ {
				needs = append(needs, fmt.Sprintf("task_%02d", j+1))
			}
		}
		fmt.Fprintf(&manifest, "    - id: %s\n      file: %s.md\n      needs: [%s]\n", id, id, strings.Join(needs, ", "))
		taskType := seed.taskType
		commands := seed.commands
		if taskType == "" {
			taskType = "backend"
		}
		if taskType == spec.TaskTypeQA {
			commands = spec.DerivedQAVerification("clean")
		}
		var task strings.Builder
		fmt.Fprintf(&task, "---\ntask: %s\nspec: clean\nstatus: %s\ntype: %s\ncomplexity: low\n---\n\n# Task: Exercise status selection\n\n## Consumers\n\nNone. This fixture exercises command reporting.\n\n## Verification\n\n", id, seed.status, taskType)
		for _, command := range commands {
			fmt.Fprintf(&task, "- `%s`\n", command)
		}
		mustWrite(t, filepath.Join(repo, "docs/specs/clean", id+".md"), task.String())
	}
	manifest.WriteString("---\n\n# Status fixture Task Graph\n")
	mustWrite(t, filepath.Join(repo, "docs/specs/clean/_tasks.md"), manifest.String())
	gitImplement(t, repo, "add", "docs/specs/clean")
	gitImplement(t, repo, "commit", "-m", "test: seed Task statuses")
	return home, repo
}

func TestSpecCheckLeavesCompletedTasksAlone(t *testing.T) {
	t.Parallel()
	const missingTool = "roundfix-verification-binary-that-does-not-exist"
	t.Run("a resumed Spec probes only the Tasks still to run", func(t *testing.T) {
		marker := filepath.Join(t.TempDir(), "completed-command-ran")
		touch := fmt.Sprintf("touch %q", marker)
		_, _ = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusCompleted, commands: []string{touch, missingTool}},
			{status: spec.StatusPending, commands: []string{"test -f pending-output"}},
			{status: spec.StatusInProgress, commands: []string{"test -f in-progress-output"}},
			{status: spec.StatusFailed, commands: []string{"test -f failed-output"}},
		})
		code, output, _ := probeProvenanceCLI(t, "text")
		if code != exitOK {
			t.Fatalf("exit=%d, want 0; output=%s", code, output)
		}
		verdict := "No findings. Authored Verification commands executed: 3; not probed (completed Tasks): 2."
		if !strings.HasPrefix(output, "Spec clean\n"+verdict+"\n") {
			t.Fatalf("verdict line missing: %s", output)
		}
		wantLines := fmt.Sprintf("Verification tree: HEAD\n- task_01: not probed — %q (Task completed)\n- task_01: not probed — %q (Task completed)\n", touch, missingTool) +
			"- task_02: honest — \"test -f pending-output\" (exited non-zero before work)\n" +
			"- task_03: honest — \"test -f in-progress-output\" (exited non-zero before work)\n" +
			"- task_04: honest — \"test -f failed-output\" (exited non-zero before work)\n"
		if _, lines, found := strings.Cut(output, "Verification tree: HEAD\n"); !found || "Verification tree: HEAD\n"+lines != wantLines {
			t.Fatalf("ordered command lines differ; output=%s; want=%s", output, wantLines)
		}
		code, output, _ = probeProvenanceCLI(t, "json")
		if code != exitOK {
			t.Fatalf("JSON exit=%d, want 0; output=%s", code, output)
		}
		var doc specCheckStatusDocument
		if err := json.Unmarshal([]byte(output), &doc); err != nil {
			t.Fatal(err)
		}
		got := doc.Verification
		if len(got.NotProbed) != 2 {
			t.Fatalf("notProbed=%+v, want two commands", got.NotProbed)
		}
		for i, command := range []string{touch, missingTool} {
			entry := got.NotProbed[i]
			if entry.Task != "task_01" || entry.Command != command || entry.Reason != "completed" {
				t.Fatalf("notProbed[%d]=%+v", i, entry)
			}
		}
		wantCommands := []specCheckVerificationCommandReport{
			{Task: "task_02", Command: "test -f pending-output", Verdict: "honest"},
			{Task: "task_03", Command: "test -f in-progress-output", Verdict: "honest"},
			{Task: "task_04", Command: "test -f failed-output", Verdict: "honest"},
		}
		if !reflect.DeepEqual(got.Commands, wantCommands) {
			t.Fatalf("commands=%+v, want %+v", got.Commands, wantCommands)
		}
		if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("completed command ran: %v", err)
		}
	})
	for _, tc := range []struct{ name, command, verdict string }{
		{"a pending vacuous command still fails", "true", "vacuous"},
		{"a pending unknown command still fails", missingTool, "unknown"},
		{"a pending malformed command still fails", "if", "malformed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
				{status: spec.StatusCompleted, commands: []string{"true"}},
				{status: spec.StatusPending, commands: []string{tc.command}},
			})
			code, output, _ := probeProvenanceCLI(t, "json")
			var doc specCheckStatusDocument
			if err := json.Unmarshal([]byte(output), &doc); err != nil {
				t.Fatal(err)
			}
			if code != exitRunFailed || len(doc.Verification.Commands) != 1 || len(doc.Verification.NotProbed) != 1 {
				t.Fatalf("exit=%d report=%s", code, output)
			}
			command := doc.Verification.Commands[0]
			if command.Task != "task_02" || command.Command != tc.command || command.Verdict != tc.verdict {
				t.Fatalf("command=%+v, want pending %s", command, tc.verdict)
			}
		})
	}
	t.Run("a pending QA gate is still probed", func(t *testing.T) {
		_, _ = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusCompleted, commands: []string{"true"}},
			{status: spec.StatusPending, taskType: spec.TaskTypeQA},
		})
		code, output, doc := probeProvenanceCLI(t, "json")
		want := []specCheckVerificationCommandReport{{Task: "task_02", Command: spec.DerivedQAVerification("clean")[0], Verdict: "honest"}}
		if code != exitOK || !reflect.DeepEqual(doc.Verification.Commands, want) {
			t.Fatalf("exit=%d report=%s", code, output)
		}
	})
	t.Run("every Task completed", func(t *testing.T) {
		_, _ = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusCompleted, commands: []string{"true"}},
			{status: spec.StatusCompleted, taskType: spec.TaskTypeQA},
		})
		code, output, _ := probeProvenanceCLI(t, "text")
		want := "Spec clean\nNo findings. Authored Verification commands executed: 0; not probed (completed Tasks): 2.\n"
		if code != exitOK || !strings.HasPrefix(output, want) || strings.Contains(output, "No authored Verification commands.") {
			t.Fatalf("exit=%d report=%s", code, output)
		}
		code, output, _ = probeProvenanceCLI(t, "json")
		var doc specCheckStatusDocument
		if err := json.Unmarshal([]byte(output), &doc); err != nil {
			t.Fatal(err)
		}
		if code != exitOK || len(doc.Verification.Commands) != 0 || len(doc.Verification.NotProbed) != 2 {
			t.Fatalf("exit=%d report=%s", code, output)
		}
	})
	t.Run("an uncommitted completed Task source is still named", func(t *testing.T) {
		_, repo := newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusCompleted, commands: []string{"true"}},
		})
		path := filepath.Join(repo, "docs/specs/clean/task_01.md")
		mustWrite(t, path, strings.Replace(mustRead(t, path), "`true`", "`false`", 1))
		code, output, _ := probeProvenanceCLI(t, "text")
		want := "Uncommitted Verification source: docs/specs/clean/task_01.md (modified)\n"
		if code != exitOK || !strings.Contains(output, want) {
			t.Fatalf("exit=%d report=%s", code, output)
		}
		code, output, doc := probeProvenanceCLI(t, "json")
		sources := []specCheckVerificationSource{{Path: "docs/specs/clean/task_01.md", State: "modified"}}
		if code != exitOK || !reflect.DeepEqual(doc.Verification.Uncommitted, sources) {
			t.Fatalf("exit=%d report=%s", code, output)
		}
	})
	t.Run("without a completed Task the JSON has no notProbed key", func(t *testing.T) {
		_, _ = newSpecCheckStatusWorkspace(t, []specCheckStatusSeed{
			{status: spec.StatusPending, commands: []string{"test -f pending-output"}},
		})
		code, output, _ := probeProvenanceCLI(t, "json")
		var doc map[string]json.RawMessage
		if err := json.Unmarshal([]byte(output), &doc); err != nil {
			t.Fatal(err)
		}
		var verification map[string]json.RawMessage
		if err := json.Unmarshal(doc["verification"], &verification); err != nil {
			t.Fatal(err)
		}
		if _, exists := verification["notProbed"]; exists || code != exitOK {
			t.Fatalf("exit=%d report=%s", code, output)
		}
	})
}
