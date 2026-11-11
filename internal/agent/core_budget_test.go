package agent

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/corebudget"
)

const fakeCoreBudgetPath = "ROUNDFIX_FAKE_CORE_BUDGET_PATH"

type fakeCoreBudgetEnvironment struct {
	Command       string
	GoCores       string
	RoundfixCores string
	RoundfixSet   bool
}

func recordFakeCoreBudgetEnvironment() error {
	path := os.Getenv(fakeCoreBudgetPath)
	if path == "" {
		return nil
	}
	roundfixCores, roundfixSet := os.LookupEnv(corebudget.Variable)
	row := fakeCoreBudgetEnvironment{
		Command:       fakeACPXCommandKey(os.Args[1:]),
		GoCores:       os.Getenv(corebudget.GoVariable),
		RoundfixCores: roundfixCores,
		RoundfixSet:   roundfixSet,
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

func TestAnAgentSessionReceivesTheCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name                 string
		cores                int
		wantGo, wantRoundfix string
	}{{"budget three", 3, "3", "3"}, {"no budget", 0, "12", ""}} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			h := newFakeACPXHarness(t)
			path := filepath.Join(h.gitRoot, "core-budget.jsonl")
			h.setEnv(fakeCoreBudgetPath, path)
			h.setEnv(corebudget.GoVariable, "12")
			// Isolate the fixture from any budget carried by the test process.
			var env []string
			for _, entry := range h.runner.Environment {
				if !strings.HasPrefix(entry, corebudget.Variable+"=") {
					env = append(env, entry)
				}
			}
			h.runner.Environment = env
			if row.cores > 0 {
				h.setEnv(corebudget.Variable, "19")
			}
			_, err := h.runner.Run(t.Context(), ExecuteRequest{
				Runtime: selectedRuntime(RuntimeSpec{ID: "claude", Protocol: ProtocolACP}),
				RunID:   "run-acpx", Cores: row.cores,
				GitRoot: h.gitRoot, Session: SessionRef{Name: "core-budget"}, Prompt: "prompt", StopGrace: h.stopGrace,
			}, newCaptureSink(""))
			if err != nil {
				t.Fatal(err)
			}
			commands := map[string]bool{}
			for _, line := range strings.Split(strings.TrimSpace(readFile(t, path)), "\n") {
				var got fakeCoreBudgetEnvironment
				if err := json.Unmarshal([]byte(line), &got); err != nil {
					t.Fatal(err)
				}
				if got.GoCores != row.wantGo || got.RoundfixCores != row.wantRoundfix {
					t.Errorf("%s environment = %q, %q; want %q, %q", got.Command, got.GoCores, got.RoundfixCores, row.wantGo, row.wantRoundfix)
				}
				if got.RoundfixSet != (row.cores > 0) {
					t.Errorf("%s Roundfix variable present = %t, want %t", got.Command, got.RoundfixSet, row.cores > 0)
				}
				commands[got.Command] = true
			}
			for _, command := range []string{"sessions ensure", "prompt"} {
				if !commands[command] {
					t.Errorf("missing %s environment", command)
				}
			}
		})
	}
}
