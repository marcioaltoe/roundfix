// Suite: Sealed Baseline ACP classification
// Invariant: semantic output enters review only after exact selection, sealed execution, complete validation, and proven cleanup.
// Boundary IN: preferred/fallback supervision, private work directories, canonical prompt bytes, cancellation, and manual fallback.
// Boundary OUT: model quality, CLI interaction, Run persistence, repository mutation, and final human approval.

package baselineacp

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/baseline"
)

func TestPreferredFallbackValidPreferredPreventsFallback(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		outputs: map[string]fakeSealedOutput{
			PreferredModel: {output: analyzerProposalJSON(t, snapshot)},
		},
	}

	proposal, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("classify with preferred selection: %v", err)
	}
	if len(proposal.Dispositions) != len(snapshot.Entries) {
		t.Fatalf("accepted proposal is incomplete: %+v", proposal)
	}
	if got, want := runtime.provedModels(), []string{PreferredModel, FallbackModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("eager proofs\nwant: %v\ngot:  %v", want, got)
	}
	if got, want := runtime.attemptedModels(), []string{PreferredModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("attempts\nwant: %v\ngot:  %v", want, got)
	}
}

func TestPreferredFallbackInvalidPreferredUsesIdenticalSnapshotBytes(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		outputs: map[string]fakeSealedOutput{
			PreferredModel: {output: []byte("not json")},
			FallbackModel:  {output: analyzerProposalJSON(t, snapshot)},
		},
	}

	if _, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot); err != nil {
		t.Fatalf("classify with fallback selection: %v", err)
	}
	if got, want := runtime.attemptedModels(), []string{PreferredModel, FallbackModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("attempts\nwant: %v\ngot:  %v", want, got)
	}
	inputs := runtime.promptInputs()
	if len(inputs) != 2 || !reflect.DeepEqual(inputs[0], inputs[1]) {
		t.Fatalf("preferred and fallback inputs differ:\npreferred: %q\nfallback:  %q", inputs[0], inputs[1])
	}
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(inputs[0], canonical) {
		t.Fatalf("attempt input is not the canonical Analysis Snapshot:\nwant: %s\ngot:  %s", canonical, inputs[0])
	}
}

func TestPreferredFallbackRejectsToolTimeoutAndOversizedOutput(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	tests := []struct {
		name      string
		preferred fakeSealedOutput
	}{
		{
			name:      "tool use",
			preferred: fakeSealedOutput{toolUsed: true},
		},
		{
			name:      "timeout",
			preferred: fakeSealedOutput{err: context.DeadlineExceeded},
		},
		{
			name: "oversized output",
			preferred: fakeSealedOutput{
				output: []byte(strings.Repeat("x", baseline.ClassificationProposalMaxBytes+1)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runtime := &fakeSealedRuntime{
				outputs: map[string]fakeSealedOutput{
					PreferredModel: test.preferred,
					FallbackModel:  {output: analyzerProposalJSON(t, snapshot)},
				},
			}
			if _, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot); err != nil {
				t.Fatalf("fallback after rejected preferred attempt: %v", err)
			}
			if got, want := runtime.attemptedModels(), []string{PreferredModel, FallbackModel}; !reflect.DeepEqual(got, want) {
				t.Fatalf("attempts\nwant: %v\ngot:  %v", want, got)
			}
		})
	}
}

func TestManualClassificationFallbackWhenSelectionsUnavailableOrInvalid(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		proofErrors: map[string]error{
			PreferredModel: errors.New("preferred unavailable"),
		},
		outputs: map[string]fakeSealedOutput{
			FallbackModel: {output: []byte(`{"schemaVersion":"invalid"}`)},
		},
	}

	proposal, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("manual classification fallback: %v", err)
	}
	if got, want := runtime.attemptedModels(), []string{FallbackModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("attempts\nwant: %v\ngot:  %v", want, got)
	}
	if len(proposal.Dispositions) != len(snapshot.Entries) {
		t.Fatalf("manual proposal is incomplete: %+v", proposal)
	}
	for _, disposition := range proposal.Dispositions {
		if disposition.Disposition != "repository-rules" && disposition.Disposition != "rejected" {
			t.Fatalf("manual proposal has unsupported destination: %+v", disposition)
		}
	}
}

func TestPreferredFallbackDoesNotStartAfterProofCleanupFailure(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		proofErrors: map[string]error{
			PreferredModel: &agent.AgentSessionCleanupError{
				Session: "proof",
				Err:     errors.New("still active"),
			},
		},
		outputs: map[string]fakeSealedOutput{
			FallbackModel: {output: analyzerProposalJSON(t, snapshot)},
		},
	}

	proposal, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("proof cleanup failure should return manual classification: %v", err)
	}
	if got, want := runtime.provedModels(), []string{PreferredModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("fallback proof started after cleanup failure\nwant: %v\ngot:  %v", want, got)
	}
	if len(runtime.attemptedModels()) != 0 {
		t.Fatalf("analysis started after cleanup failure: %v", runtime.attemptedModels())
	}
	if len(proposal.Dispositions) != len(snapshot.Entries) {
		t.Fatalf("proof cleanup failure did not return complete manual proposal: %+v", proposal)
	}
}

func TestPreferredFallbackRejectsToolUseAndStopsAfterCleanupFailure(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		outputs: map[string]fakeSealedOutput{
			PreferredModel: {
				toolUsed: true,
				err:      &agent.AgentSessionCleanupError{Session: "sealed", Err: errors.New("still active")},
			},
			FallbackModel: {output: analyzerProposalJSON(t, snapshot)},
		},
	}

	proposal, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("cleanup failure should return manual classification: %v", err)
	}
	if got, want := runtime.attemptedModels(), []string{PreferredModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("fallback activated without proven cleanup\nwant: %v\ngot:  %v", want, got)
	}
	if len(proposal.Dispositions) != len(snapshot.Entries) {
		t.Fatalf("cleanup failure did not return complete manual proposal: %+v", proposal)
	}
}

func TestAnalyzerCancellationClosesActiveSessionWithoutFallback(t *testing.T) {
	snapshot := analyzerTestSnapshot(t)
	runtime := &cancelingSealedRuntime{started: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := NewAnalyzer(runtime).Classify(ctx, snapshot)
		result <- err
	}()
	<-runtime.started
	cancel()

	err := <-result
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %v, want context.Canceled", err)
	}
	if !runtime.cleaned {
		t.Fatal("active sealed session was not cleaned before cancellation returned")
	}
	if got, want := runtime.attemptedModels(), []string{PreferredModel}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cancellation activated fallback\nwant: %v\ngot:  %v", want, got)
	}
}

func TestAnalyzerNoPersistenceOrRepositoryMutation(t *testing.T) {
	checkout := t.TempDir()
	sentinel := filepath.Join(checkout, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	snapshot := analyzerTestSnapshot(t)
	runtime := &fakeSealedRuntime{
		outputs: map[string]fakeSealedOutput{
			PreferredModel: {output: analyzerProposalJSON(t, snapshot)},
		},
		checkPrivateDirectories: true,
		forbiddenPath:           checkout,
	}

	if _, err := NewAnalyzer(runtime).Classify(context.Background(), snapshot); err != nil {
		t.Fatalf("sealed classification: %v", err)
	}
	content, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "unchanged\n" {
		t.Fatalf("classification mutated checkout sentinel: %q", content)
	}
	if runtime.nonEmptyDirectory || runtime.nonPrivateDirectory || runtime.checkoutExposed {
		t.Fatalf("sealed directory violation: nonempty=%t nonprivate=%t checkout=%t",
			runtime.nonEmptyDirectory, runtime.nonPrivateDirectory, runtime.checkoutExposed)
	}
	if got, want := runtime.uniqueWorkDirectoryCount(), 3; got != want {
		t.Fatalf("proofs and attempt used %d private directories, want %d fresh directories", got, want)
	}
}

func TestRevisionManualFallback(t *testing.T) {
	snapshot := analyzerRevisionSnapshot(t)
	runtime := &fakeSealedRuntime{
		proofErrors: map[string]error{
			PreferredModel: errors.New("preferred unavailable"),
			FallbackModel:  errors.New("fallback unavailable"),
		},
	}
	proposal, err := NewAnalyzer(runtime).Revise(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("manual revision fallback: %v", err)
	}
	if !proposal.Manual || len(proposal.Changes) != 0 {
		t.Fatalf("fallback proposal = %+v, want unchanged manual next action", proposal)
	}
}

func TestScopedRevisionProposal(t *testing.T) {
	snapshot := analyzerRevisionSnapshot(t)
	output, err := json.Marshal(baseline.RevisionProposal{
		SchemaVersion:  baseline.RevisionProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Area:           snapshot.Area,
		Changes: []baseline.DecisionValue{
			{ID: snapshot.AllowedDecisionIDs[0], Value: false},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	runtime := &fakeSealedRuntime{
		outputs: map[string]fakeSealedOutput{
			PreferredModel: {output: output},
		},
	}
	proposal, err := NewAnalyzer(runtime).Revise(context.Background(), snapshot)
	if err != nil {
		t.Fatalf("scoped semantic revision: %v", err)
	}
	if proposal.Manual || len(proposal.Changes) != 1 {
		t.Fatalf("accepted proposal = %+v", proposal)
	}
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		t.Fatal(err)
	}
	inputs := runtime.promptInputs()
	if len(inputs) != 1 || !reflect.DeepEqual(inputs[0], canonical) {
		t.Fatalf("revision attempt did not receive exact canonical snapshot")
	}
}

func TestACPXReadOnlyArguments(t *testing.T) {
	runtime := agent.RuntimeSpec{
		ID:              "codex",
		Protocol:        agent.ProtocolACP,
		Model:           PreferredModel,
		ReasoningEffort: RequiredReasoningEffort,
		FullAccessMode:  "full-access",
	}
	ensure, err := agent.SealedEnsureArguments(runtime, "roundfix-baseline-test", "/private/empty")
	if err != nil {
		t.Fatal(err)
	}
	prompt, err := agent.SealedPromptArguments(runtime, "roundfix-baseline-test", "/private/empty")
	if err != nil {
		t.Fatal(err)
	}
	for name, args := range map[string][]string{"ensure": ensure, "prompt": prompt} {
		joined := strings.Join(args, "\x00")
		for _, required := range []string{
			"--cwd\x00/private/empty",
			"--deny-all",
			"--non-interactive-permissions\x00fail",
			"--allowed-tools\x00",
			"--max-turns\x001",
			"--no-terminal",
			"--timeout\x00120",
			"--ttl\x00120",
		} {
			if !strings.Contains(joined, required) {
				t.Fatalf("%s args omit sealed flag %q: %#v", name, required, args)
			}
		}
		for _, forbidden := range []string{"--approve-all", "--approve-reads", "danger-full-access", "full-access"} {
			if strings.Contains(joined, forbidden) {
				t.Fatalf("%s args expose forbidden access %q: %#v", name, forbidden, args)
			}
		}
	}
	if strings.Contains(strings.Join(prompt, "\x00"), checkoutPathForTest) {
		t.Fatalf("sealed args expose checkout path %q: %#v", checkoutPathForTest, prompt)
	}
}

const checkoutPathForTest = "/repository/checkout"

func analyzerRevisionSnapshot(t *testing.T) baseline.RevisionSnapshot {
	t.Helper()
	plan := baseline.PlanDocument{
		SchemaVersion: baseline.PlanSchemaVersion,
		PlanDigest:    "sha256:fixture",
		Decisions: []baseline.DecisionValue{
			{ID: "spec.scaffold", Value: true},
		},
	}
	snapshot, err := baseline.NewRevisionSnapshot(
		plan,
		baseline.RevisionAreaDivergences,
		"disable spec scaffolding",
	)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

type fakeSealedOutput struct {
	output   []byte
	toolUsed bool
	err      error
}

type recordedSealedCall struct {
	kind    string
	model   string
	workDir string
	input   []byte
}

type fakeSealedRuntime struct {
	mu                      sync.Mutex
	proofErrors             map[string]error
	outputs                 map[string]fakeSealedOutput
	calls                   []recordedSealedCall
	checkPrivateDirectories bool
	forbiddenPath           string
	nonEmptyDirectory       bool
	nonPrivateDirectory     bool
	checkoutExposed         bool
}

func (runtime *fakeSealedRuntime) ProveProfileSelection(
	_ context.Context,
	request agent.ProbeRequest,
) (agent.SelectionProof, error) {
	runtime.record("proof", request.Runtime.Model, request.WorkDir, nil)
	if err := runtime.proofErrors[request.Runtime.Model]; err != nil {
		return agent.SelectionProof{}, err
	}
	return agent.SelectionProof{
		Runtime:         request.Runtime.ID,
		Model:           request.Runtime.Model,
		ReasoningEffort: request.Runtime.ReasoningEffort,
		Status:          agent.SelectionProofStatusProven,
	}, nil
}

func (runtime *fakeSealedRuntime) RunSealedPrompt(
	_ context.Context,
	request agent.SealedPromptRequest,
) (agent.SealedPromptResult, error) {
	runtime.record("attempt", request.Runtime.Model, request.WorkDir, request.Input)
	output := runtime.outputs[request.Runtime.Model]
	return agent.SealedPromptResult{
		Output:   append([]byte(nil), output.output...),
		ToolUsed: output.toolUsed,
	}, output.err
}

func (runtime *fakeSealedRuntime) record(kind, model, workDir string, input []byte) {
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	runtime.calls = append(runtime.calls, recordedSealedCall{
		kind: kind, model: model, workDir: workDir, input: append([]byte(nil), input...),
	})
	if runtime.forbiddenPath != "" &&
		(strings.Contains(workDir, runtime.forbiddenPath) || bytesContain(input, []byte(runtime.forbiddenPath))) {
		runtime.checkoutExposed = true
	}
	if !runtime.checkPrivateDirectories {
		return
	}
	info, err := os.Stat(workDir)
	if err != nil || !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		runtime.nonPrivateDirectory = true
	}
	entries, err := os.ReadDir(workDir)
	if err != nil || len(entries) != 0 {
		runtime.nonEmptyDirectory = true
	}
}

func (runtime *fakeSealedRuntime) provedModels() []string {
	return runtime.modelsFor("proof")
}

func (runtime *fakeSealedRuntime) attemptedModels() []string {
	return runtime.modelsFor("attempt")
}

func (runtime *fakeSealedRuntime) modelsFor(kind string) []string {
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	var models []string
	for _, call := range runtime.calls {
		if call.kind == kind {
			models = append(models, call.model)
		}
	}
	return models
}

func (runtime *fakeSealedRuntime) promptInputs() [][]byte {
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	var inputs [][]byte
	for _, call := range runtime.calls {
		if call.kind == "attempt" {
			inputs = append(inputs, append([]byte(nil), call.input...))
		}
	}
	return inputs
}

func (runtime *fakeSealedRuntime) uniqueWorkDirectoryCount() int {
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	unique := make(map[string]struct{}, len(runtime.calls))
	for _, call := range runtime.calls {
		unique[call.workDir] = struct{}{}
	}
	return len(unique)
}

type cancelingSealedRuntime struct {
	mu      sync.Mutex
	calls   []recordedSealedCall
	started chan struct{}
	cleaned bool
}

func (runtime *cancelingSealedRuntime) ProveProfileSelection(
	_ context.Context,
	request agent.ProbeRequest,
) (agent.SelectionProof, error) {
	runtime.mu.Lock()
	runtime.calls = append(runtime.calls, recordedSealedCall{kind: "proof", model: request.Runtime.Model})
	runtime.mu.Unlock()
	return agent.SelectionProof{
		Runtime:         request.Runtime.ID,
		Model:           request.Runtime.Model,
		ReasoningEffort: request.Runtime.ReasoningEffort,
		Status:          agent.SelectionProofStatusProven,
	}, nil
}

func (runtime *cancelingSealedRuntime) RunSealedPrompt(
	ctx context.Context,
	request agent.SealedPromptRequest,
) (agent.SealedPromptResult, error) {
	runtime.mu.Lock()
	runtime.calls = append(runtime.calls, recordedSealedCall{kind: "attempt", model: request.Runtime.Model})
	runtime.mu.Unlock()
	close(runtime.started)
	<-ctx.Done()
	runtime.cleaned = true
	return agent.SealedPromptResult{}, ctx.Err()
}

func (runtime *cancelingSealedRuntime) attemptedModels() []string {
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	var models []string
	for _, call := range runtime.calls {
		if call.kind == "attempt" {
			models = append(models, call.model)
		}
	}
	return models
}

func analyzerTestSnapshot(t *testing.T) baseline.AnalysisSnapshot {
	t.Helper()
	content := []byte("preserve this repository rule\n")
	contentSum := sha256.Sum256(content)
	carrierSum := sha256.Sum256(content)
	provenance := map[string]any{"markerState": "unmarked"}
	entryID := analyzerSourceEntryID(
		t,
		"AGENTS.md",
		"unmarked-span",
		0,
		len(content),
		hex.EncodeToString(contentSum[:]),
		provenance,
	)
	source := baseline.ReadoptionSourceBaseline{
		ID:               "baseline.readoption." + strings.Repeat("b", 64),
		DeclaredIdentity: "unconfigured",
		Compatibility:    "incompatible",
		Digest:           strings.Repeat("b", 64),
		CarrierCount:     1,
		EntryCount:       1,
		ByteCount:        len(content),
		Entries: []baseline.ReadoptionSourceEntry{{
			ID:                   entryID,
			Path:                 "AGENTS.md",
			Carrier:              "AGENTS.md",
			Kind:                 "unmarked-span",
			Start:                0,
			End:                  len(content),
			Digest:               hex.EncodeToString(contentSum[:]),
			CarrierDigest:        hex.EncodeToString(carrierSum[:]),
			SourceBytes:          content,
			Encoding:             "base64",
			StructuralProvenance: provenance,
		}},
	}
	snapshot, err := baseline.NewAnalysisSnapshot(source)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func analyzerSourceEntryID(
	t *testing.T,
	sourcePath string,
	kind string,
	start int,
	end int,
	digest string,
	provenance map[string]any,
) string {
	t.Helper()
	provenanceJSON, err := json.Marshal(provenance)
	if err != nil {
		t.Fatal(err)
	}
	identity := sha256.New()
	for _, value := range []string{
		sourcePath,
		kind,
		strconv.Itoa(start),
		strconv.Itoa(end),
		digest,
		string(provenanceJSON),
	} {
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(value)))
		_, _ = identity.Write(length[:])
		_, _ = identity.Write([]byte(value))
	}
	return "source-entry." + hex.EncodeToString(identity.Sum(nil))
}

func analyzerProposalJSON(t *testing.T, snapshot baseline.AnalysisSnapshot) []byte {
	t.Helper()
	proposal, err := baseline.ManualClassificationProposal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func bytesContain(data, fragment []byte) bool {
	return strings.Contains(string(data), string(fragment))
}
