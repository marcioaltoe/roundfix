// Suite: disposable profile proof close.
// Invariant: the setup verdict survives an unproven close.
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	h.runner.cleanupTimeout = 100 * time.Millisecond
	h.setEnv(fakeACPXBlockCmd, "sessions close")
	var warnings []string
	h.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
	proof, err := h.runner.ProveExactSelection(context.Background(), proofCloseRequest(h))
	if err != nil || proof.Status != SelectionProofStatusProven || proof.Unclosed == nil {
		t.Fatalf("passing setup lost its verdict: proof=%+v error=%v", proof, err)
	}
	if proof.Unclosed.Deadline != h.runner.cleanupTimeout || !strings.Contains(proof.Unclosed.Warning, context.DeadlineExceeded.Error()) {
		t.Fatalf("blocked close deadline not reported: %+v", proof.Unclosed)
	}
	assertProofCloseRecord(t, h, proof.Unclosed, warnings)

}

func proofCloseHarness(t *testing.T) *fakeACPXHarness {
	t.Helper()
	h := newFakeACPXHarness(t)
	models := []string{"future-model", "future-model[high]", "future-model[xhigh]"}
	h.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show":                      sessionCapabilitySnapshotFixture(t, "future-model", models, "", "", nil),
		"set model value=future-model":       selectionStateFixture(t, "model", "future-model", "future-model", models, "", "", nil),
		"set model value=future-model[high]": selectionStateFixture(t, "model", "future-model[high]", "future-model[high]", models, "", "", nil),
	}))
	return h
}

func proofCloseRequest(h *fakeACPXHarness) ProbeRequest {
	return ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "future-model", ReasoningEffort: "high"}, WorkDir: h.gitRoot}
}

func TestTheProofCloseDeadlineFollowsTheSetupDeadline(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"recorded", "no HOME", "write error"} {
		t.Run(mode, func(t *testing.T) {
			h := proofCloseHarness(t)
			h.runner.cleanupTimeout = 0
			h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 2}))
			// Keep adapter discovery valid when testing the record's HOME failure.
			var warnings []string
			h.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
			if mode == "write error" {
				home := environmentValue(h.runner.baseEnv(), "HOME")
				if err := os.WriteFile(filepath.Join(home, ".roundfix"), []byte("blocked"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			request := proofCloseRequest(h)
			request.SetupTimeout = 2 * time.Hour
			proof, err := h.runner.ProveExactSelection(context.Background(), request)
			if err != nil || proof.Unclosed == nil {
				t.Fatalf("proof=%+v err=%v", proof, err)
			}
			if mode == "no HOME" {
				// Exercise the writer without changing the setup's adapter HOME.
				warnings = nil
				h.setEnv("HOME", "")
				proof.Unclosed = h.runner.recordUnclosedProofSession(request.Runtime, proof.Unclosed.Session, h.gitRoot, request.SetupTimeout, errors.New("close rejected"))
			}
			if len(warnings) != 1 || !strings.Contains(warnings[0], "within 2h0m0s") {
				t.Fatalf("warnings=%v", warnings)
			}
			if mode == "recorded" {
				assertProofCloseRecord(t, h, proof.Unclosed, warnings)
			} else if !strings.Contains(warnings[0], "; it could not be recorded (") || !strings.Contains(warnings[0], "; close it with: acpx ") || !strings.Contains(warnings[0], "sessions close") {
				t.Fatalf("missing recovery command: %v", warnings)
			}
		})
	}
}

func TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails(t *testing.T) {
	t.Parallel()
	for _, entry := range []string{"exact proof", "probe"} {
		t.Run(entry, func(t *testing.T) {
			h := newFakeACPXHarness(t)
			h.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
			h.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"sessions show":               sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
				"set model value=gpt-5.6-sol": selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol"}, "reasoning_effort", "medium", []string{"medium", "high"}),
			}))
			h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set reasoning_effort": 2, "sessions close": 2}))
			var warnings []string
			h.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
			request := ProbeRequest{Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"}, WorkDir: h.gitRoot}
			var err error
			if entry == "probe" {
				err = h.runner.Probe(context.Background(), request)
			} else {
				_, err = h.runner.ProveExactSelection(context.Background(), request)
			}
			var rejected *SelectionRejectedError
			var cleanup *AgentSessionCleanupError
			if !errors.As(err, &rejected) || errors.As(err, &cleanup) {
				t.Fatalf("want setup rejection only, got %T %v", err, err)
			}
			const want = `apply Agent Selection "codex"/"gpt-5.6-sol"/"high" during set reasoning_effort: adapter rejected selection: acquire Agent Selection capabilities through acpx: command failed; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection`
			if err.Error() != want {
				t.Fatalf("setup diagnosis changed: %v", err)
			}
			assertProofCloseRecord(t, h, nil, warnings)
			assertLastInvocationClosesDisposable(t, h)
		})
	}
}

func TestATimedOutProofStepIsNamed(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	started := filepath.Join(h.gitRoot, "catalogue-started")
	h.setEnv(fakeACPXStarted, started)
	h.setEnv(fakeACPXBlockCmd, "sessions show")
	ctx := newTestDeadlineContext()
	defer ctx.expire()
	result := make(chan error, 1)
	go func() { _, err := h.runner.ProveExactSelection(ctx, proofCloseRequest(h)); result <- err }()
	waitForFile(t, "catalogue blocked", started, result)
	ctx.expire()
	err := receiveError(t, result)
	if !errors.Is(err, context.DeadlineExceeded) || !strings.HasPrefix(err.Error(), `profile proof step "read runtime catalogue": `) {
		t.Fatalf("unnamed setup deadline: %v", err)
	}
	assertLastInvocationClosesDisposable(t, h)
}

func TestAFallbackCandidateKeepsItsResultWhenItsCloseFails(t *testing.T) {
	// Sequential: ProbeFallback resolves its working directory from process state through t.Chdir.
	for _, outcome := range []string{"functional", "rejected", "error"} {
		t.Run(outcome, func(t *testing.T) {
			h := newFakeACPXHarness(t)
			t.Chdir(h.gitRoot)
			exits := map[string]int{"sessions close": 2}
			if outcome == "rejected" {
				exits["sessions ensure"] = 2
			}
			if outcome == "error" {
				h.runner.Command = filepath.Join(h.gitRoot, "missing-acpx")
			}
			h.setEnv(fakeACPXExitBy, mustJSONForTest(t, exits))
			var warnings []string
			h.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
			selection, ok, err := h.runner.ProbeFallback(context.Background(), RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "failed"}, FallbackCandidateSet{Models: []string{"working"}, Efforts: []string{"high"}})
			if outcome == "functional" {
				if err != nil || !ok || selection != (FallbackSelection{Model: "working", ReasoningEffort: "high"}) {
					t.Fatalf("candidate lost: %+v %t %v", selection, ok, err)
				}
			} else {
				var cleanup *AgentSessionCleanupError
				if ok || selection != (FallbackSelection{}) || errors.As(err, &cleanup) {
					t.Fatalf("close replaced candidate result: %+v %t %v", selection, ok, err)
				}
				if outcome == "error" && err == nil {
					t.Fatal("candidate error was lost")
				}
				if outcome == "rejected" && err != nil {
					t.Fatalf("rejection became error: %v", err)
				}
			}
			assertProofCloseRecord(t, h, nil, warnings)
		})
	}
}

func TestAnUnclosedProofSessionRecordHoldsNoEnvironmentValue(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	sentinel := strings.Join([]string{"private", "fixture", "value"}, "-")
	request := proofCloseRequest(h)
	request.Runtime.OpenRouterKeyVariable = "ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY"
	h.setEnv(request.Runtime.OpenRouterKeyVariable, sentinel)
	h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 2}))
	var warnings []string
	h.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
	proof, err := h.runner.ProveExactSelection(context.Background(), request)
	if err != nil || proof.Unclosed == nil {
		t.Fatalf("proof=%+v err=%v", proof, err)
	}
	data := assertProofCloseRecord(t, h, proof.Unclosed, warnings)
	if strings.Contains(string(data), sentinel) || !strings.Contains(string(data), request.Runtime.OpenRouterKeyVariable) {
		t.Fatal("record leaked environment value or lost variable name")
	}
}

func assertProofCloseRecord(t *testing.T, h *fakeACPXHarness, unclosed *UnclosedProofSession, warnings []string) []byte {
	t.Helper()
	dir := filepath.Join(environmentValue(h.runner.baseEnv(), "HOME"), ".roundfix", "unclosed-proof-sessions")
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("want one record: entries=%v err=%v", entries, err)
	}
	path := filepath.Join(dir, entries[0].Name())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]json.RawMessage
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if len(record) != 6 {
		t.Fatalf("record fields=%v", record)
	}
	var session, schema, workDir string
	_ = json.Unmarshal(record["session"], &session)
	_ = json.Unmarshal(record["schema"], &schema)
	_ = json.Unmarshal(record["work_dir"], &workDir)
	if schema != "roundfix/unclosed-proof-session/v1" || workDir != h.gitRoot || entries[0].Name() != session+".json" {
		t.Fatalf("record metadata=%s", data)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], fmt.Sprintf("disposable Agent Session %q was not closed within ", session)) || !strings.HasSuffix(warnings[0], "; recorded in "+path+" for the next profile proof to close") {
		t.Fatalf("warnings=%v", warnings)
	}
	if unclosed != nil && (unclosed.Session != session || unclosed.WorkDir != workDir || unclosed.Warning != warnings[0]) {
		t.Fatalf("unclosed=%+v warnings=%v", unclosed, warnings)
	}
	for _, check := range []struct {
		path string
		mode os.FileMode
	}{{dir, 0o700}, {path, 0o600}} {
		info, err := os.Stat(check.path)
		if err != nil || info.Mode().Perm() != check.mode {
			t.Fatalf("private mode: path=%s info=%v err=%v", check.path, info, err)
		}
	}
	return data
}
