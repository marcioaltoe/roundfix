package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/judge"
	"roundfix/internal/openrouterkey"
)

const fakeJudgeEnvironmentPath = "ROUNDFIX_FAKE_JUDGE_ENVIRONMENT_PATH"

type fakeJudgeEnvironment struct {
	Command string
	Names   []string
	RunID   string
}

// The fake acpx records variable names and the Run Marker, never credential values.
func recordFakeJudgeEnvironment() error {
	path := os.Getenv(fakeJudgeEnvironmentPath)
	if path == "" {
		return nil
	}
	row := fakeJudgeEnvironment{Command: fakeACPXCommandKey(os.Args[1:]), RunID: os.Getenv(judge.RunMarkerVariable)}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		row.Names = append(row.Names, name)
	}
	data, err := json.Marshal(row)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(append(data, '\n'))
	return errors.Join(writeErr, file.Close())
}

func TestNoAgentSessionReceivesAJudgeKeyVariable(t *testing.T) {
	t.Parallel()
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	variables := append(questions.KeyVariables(), "TYPESAFE_API_KEY")
	const control = "ROUNDFIX_JUDGE_ENVIRONMENT_CONTROL"
	for _, kind := range []string{"standard", "light"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			h := newFakeACPXHarness(t)
			path := filepath.Join(h.gitRoot, "environment-names.jsonl")
			h.setEnv(fakeJudgeEnvironmentPath, path)
			h.setEnv(control, fmt.Sprintf("control-%s", kind))
			h.setEnv(judge.RunMarkerVariable, "inherited-run")
			for index, variable := range variables {
				h.setEnv(variable, fmt.Sprintf("fixture-%s-%d", kind, index))
			}
			h.setEnv("OPENCODE_CONFIG_CONTENT", "")
			h.setEnv("NODE_OPTIONS", "--no-warnings")
			h.setEnv("OPENROUTER_API_KEY", "generic-fixture")
			h.setEnv(openrouterkey.Implement, "")
			runtime := RuntimeSpec{ID: "claude", Protocol: ProtocolACP}
			if kind == "light" {
				runtime = lightRuntimeForTest()
			}
			// Run prepares the session before prompting it.
			if _, err := h.run(t.Context(), runtime, "judge-key-"+kind); err != nil {
				t.Fatal(err)
			}
			session := SessionRef{Name: "judge-key-" + kind, WorkDir: h.gitRoot}
			if err := h.runner.CancelSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			if err := h.runner.EndSession(t.Context(), runtime, session); err != nil {
				t.Fatal(err)
			}
			if _, err := h.runner.ListRoundfixSessions(t.Context(), runtime, h.gitRoot); err != nil {
				t.Fatal(err)
			}
			if got := environmentValue(h.runner.Environment, judge.RunMarkerVariable); got != "inherited-run" {
				t.Fatalf("runner inherited Run Marker changed to %q", got)
			}
			rows := strings.Split(strings.TrimSpace(readFile(t, path)), "\n")
			invocations := readJSONInvocations(t, h.invocationsPath)
			if len(rows) != len(invocations) {
				t.Fatal("not every acpx invocation recorded its environment")
			}
			commands := make(map[string]bool)
			for index, line := range rows {
				var row fakeJudgeEnvironment
				if err := json.Unmarshal([]byte(line), &row); err != nil {
					t.Fatal(err)
				}
				if row.Command != fakeACPXCommandKey(invocations[index]) {
					t.Fatal("environment record does not match invocation")
				}
				commands[row.Command] = true
				switch row.Command {
				case "cancel", "sessions close", "sessions list":
				default:
					if row.RunID != "run-acpx" {
						t.Errorf("%s received Run Marker %q, want request Run ID", row.Command, row.RunID)
					}
				}
				for _, variable := range variables {
					if slices.Contains(row.Names, variable) {
						t.Errorf("%s received judge variable %s", row.Command, variable)
					}
				}
				if !slices.Contains(row.Names, control) {
					t.Errorf("%s lost control variable", row.Command)
				}
				if kind == "light" && !slices.Contains(row.Names, openrouterkey.Implement) {
					t.Errorf("%s lost implementation Stage Key variable", row.Command)
				}
				for _, entry := range h.runner.Environment {
					name, _, _ := strings.Cut(entry, "=")
					if slices.Contains(variables, name) || name == codexPathEnv || name == claudeNestedGuardEnv || kind == "light" && name == "OPENROUTER_API_KEY" {
						continue
					}
					if !slices.Contains(row.Names, name) {
						t.Errorf("%s lost unrelated variable %s", row.Command, name)
					}
				}
			}
			for _, command := range []string{"sessions ensure", "prompt", "cancel", "sessions close", "sessions list"} {
				if !commands[command] {
					t.Errorf("session did not exercise %s", command)
				}
			}
		})
	}
}
