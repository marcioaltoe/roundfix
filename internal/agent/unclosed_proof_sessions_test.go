// Suite: recorded proof session recovery.
// Invariant: a later runner closes only recorded sessions before its own proof.
package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestTheNextProofClosesARecordedSessionFirst(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 2}))
	proof, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h))
	if err != nil || proof.Unclosed == nil {
		t.Fatalf("proof=%+v err=%v", proof, err)
	}
	session := proof.Unclosed.Session
	path := proofRecordPath(h, session)
	h.setEnv(fakeACPXExitBy, "")
	nextProofRunner(h)
	start := len(readJSONInvocations(t, h.invocationsPath))
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	calls := readJSONInvocations(t, h.invocationsPath)[start:]
	assertRecordedCloseBeforeEnsure(t, calls, session)
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("record retained: %v", err)
	}
}

func TestARecordedSessionAcpxReportsAbsentIsForgotten(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	path := seedProofRecord(t, h, "roundfix-preflight-absent", time.Unix(1, 0))
	h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 4}))
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("absent session retained: %v", err)
	}
	assertRecordedCloseBeforeEnsure(t, readJSONInvocations(t, h.invocationsPath), "roundfix-preflight-absent")
}

func TestAFailedReapKeepsTheRecordAndStops(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	// Lexical order deliberately opposes age order.
	newer := seedProofRecord(t, h, "roundfix-preflight-a-newer", time.Unix(2, 0))
	older := seedProofRecord(t, h, "roundfix-preflight-z-older", time.Unix(1, 0))
	h.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"sessions close": 2}))
	var warnings []string
	h.runner.warnf = func(f string, a ...any) { warnings = append(warnings, fmt.Sprintf(f, a...)) }
	proof, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h))
	if err != nil || proof.Status != SelectionProofStatusProven {
		t.Fatalf("reap changed verdict: %+v %v", proof, err)
	}
	for _, path := range []string{older, newer} {
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
	}
	calls := readJSONInvocations(t, h.invocationsPath)
	assertRecordedCloseBeforeEnsure(t, calls, "roundfix-preflight-z-older")
	for _, call := range calls {
		if slices.Contains(call, "roundfix-preflight-a-newer") {
			t.Fatalf("closed newer record after failure: %v", call)
		}
	}
	want := fmt.Sprintf("recorded disposable Agent Session %q (%s) was not closed within %s:", "roundfix-preflight-z-older", older, h.runner.cleanupTimeout)
	if len(warnings) != 2 || !strings.HasPrefix(warnings[0], want) || !strings.HasSuffix(warnings[0], "; the next profile proof retries it") {
		t.Fatalf("warnings=%v", warnings)
	}
}

func TestTheReapRunsOncePerRunner(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	session := "roundfix-preflight-later"
	path := seedProofRecord(t, h, session, time.Unix(1, 0))
	start := len(readJSONInvocations(t, h.invocationsPath))
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	for _, call := range readJSONInvocations(t, h.invocationsPath)[start:] {
		if slices.Contains(call, session) {
			t.Fatalf("reaped twice: %v", call)
		}
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("record disappeared: %v", err)
	}
}

func TestAnUnreadableRecordIsReportedAndKept(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	path := seedProofRecord(t, h, "roundfix-preflight-unreadable", time.Unix(1, 0))
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	var warnings []string
	h.runner.warnf = func(f string, a ...any) { warnings = append(warnings, fmt.Sprintf(f, a...)) }
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "{" {
		t.Fatalf("unreadable record changed: %s %v", data, err)
	}
	if len(warnings) != 1 || !strings.HasPrefix(warnings[0], "unreadable Unclosed Proof Session record "+path+": ") || !strings.HasSuffix(warnings[0], "; remove it after closing its session by hand") {
		t.Fatalf("warnings=%v", warnings)
	}
}

func nextProofRunner(h *fakeACPXHarness) {
	previous := h.runner
	h.runner = &ACPXRunner{Command: previous.Command, Environment: previous.Environment, Now: previous.Now, cleanupTimeout: previous.cleanupTimeout, codexSpawn: codexSpawnDependencies{goos: "linux"}}
}
func proofRecordPath(h *fakeACPXHarness, session string) string {
	return filepath.Join(environmentValue(h.runner.baseEnv(), "HOME"), ".roundfix", "unclosed-proof-sessions", session+".json")
}
func seedProofRecord(t *testing.T, h *fakeACPXHarness, session string, at time.Time) string {
	t.Helper()
	runtime := proofCloseRequest(h).Runtime
	path, err := h.runner.writeUnclosedProofSession(unclosedProofSessionRecord{Schema: "roundfix/unclosed-proof-session/v1", Session: session, WorkDir: h.gitRoot, RecordedAt: at, Cause: "close rejected", Runtime: unclosedProofRuntime{ID: runtime.ID, Protocol: runtime.Protocol, Command: runtime.Command}})
	if err != nil {
		t.Fatal(err)
	}
	return path
}
func assertRecordedCloseBeforeEnsure(t *testing.T, calls [][]string, session string) {
	t.Helper()
	closeIndex, ensureIndex := -1, -1
	for i, call := range calls {
		if fakeACPXCommandKey(call) == "sessions close" && slices.Contains(call, session) {
			closeIndex = i
		}
		if fakeACPXCommandKey(call) == "sessions ensure" && ensureIndex < 0 {
			ensureIndex = i
		}
	}
	if closeIndex < 0 || ensureIndex < 0 || closeIndex >= ensureIndex {
		t.Fatalf("recorded close must precede ensure: %v", calls)
	}
}

func TestARecordedSessionReapHonorsCancellation(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	path := seedProofRecord(t, h, "roundfix-preflight-cancelled", time.Unix(1, 0))
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, _ = h.runner.ProveExactSelection(ctx, proofCloseRequest(h))
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(h.invocationsPath); err == nil && len(data) > 0 {
		t.Fatalf("cancelled proof started commands: %s", data)
	}
}

func TestARecordedSessionIsReapedBeforeAFallbackCandidate(t *testing.T) {
	t.Parallel()
	h := proofCloseHarness(t)
	session := "roundfix-preflight-fallback-record"
	path := seedProofRecord(t, h, session, time.Unix(1, 0))
	selection, ok, err := h.runner.probeFallbackCandidate(t.Context(), proofCloseRequest(h).Runtime, "future-model", []string{"high"}, h.gitRoot)
	if err != nil || !ok || selection.Model != "future-model" {
		t.Fatalf("candidate=%+v ok=%t err=%v", selection, ok, err)
	}
	assertRecordedCloseBeforeEnsure(t, readJSONInvocations(t, h.invocationsPath), session)
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("record retained: %v", err)
	}
	// Both proof paths share the once guard.
	later := "roundfix-preflight-after-fallback"
	seedProofRecord(t, h, later, time.Unix(2, 0))
	start := len(readJSONInvocations(t, h.invocationsPath))
	if _, err := h.runner.ProveExactSelection(t.Context(), proofCloseRequest(h)); err != nil {
		t.Fatal(err)
	}
	for _, call := range readJSONInvocations(t, h.invocationsPath)[start:] {
		if slices.Contains(call, later) {
			t.Fatalf("exact proof repeated fallback reap: %v", call)
		}
	}
}
