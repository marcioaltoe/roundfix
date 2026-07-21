package agent

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSelectionCapabilitiesOfficialAndLegacyFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fixture       string
		adapter       AdapterEvidence
		wantModel     string
		wantModels    []ModelCapability
		wantReasoning *SelectCapability
	}{
		{
			name:      "official adapter exposes advertised models and independent reasoning",
			fixture:   officialAdapterCapabilityFixture(),
			adapter:   AdapterEvidence{Command: "npx -y @agentclientprotocol/codex-acp", Package: CodexAdapterPackage, Version: "1.1.4"},
			wantModel: "gpt-5.6-sol",
			wantModels: []ModelCapability{
				{AdapterValue: "gpt-5.6-sol", CanonicalModel: "gpt-5.6-sol", ModelManaged: true},
				{AdapterValue: "gpt-5.6-terra", CanonicalModel: "gpt-5.6-terra", ModelManaged: true},
				{AdapterValue: "gpt-5.6-luna", CanonicalModel: "gpt-5.6-luna", ModelManaged: true},
				{AdapterValue: "gpt-5.5", CanonicalModel: "gpt-5.5", ModelManaged: true},
			},
			wantReasoning: &SelectCapability{
				ID:           "reasoning_effort",
				CurrentValue: "high",
				Values:       []string{"low", "medium", "high", "xhigh", "max", "ultra"},
			},
		},
		{
			name:      "legacy adapter exposes only advertised Sol model",
			fixture:   legacyAdapterCapabilityFixture(),
			adapter:   AdapterEvidence{Command: "codex-acp", Package: "@zed-industries/codex-acp", Version: "0.16.0"},
			wantModel: "gpt-5.6-sol",
			wantModels: []ModelCapability{
				{AdapterValue: "gpt-5.6-sol", CanonicalModel: "gpt-5.6-sol", ModelManaged: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSessionConfigOptions([]byte(tt.fixture), tt.adapter)
			if err != nil {
				t.Fatalf("parse fixture: %v", err)
			}
			if got.CurrentModel != tt.wantModel {
				t.Fatalf("current model = %q, want %q", got.CurrentModel, tt.wantModel)
			}
			if !reflect.DeepEqual(got.Models, tt.wantModels) {
				t.Fatalf("models = %#v, want %#v", got.Models, tt.wantModels)
			}
			if !reflect.DeepEqual(got.ReasoningOption, tt.wantReasoning) {
				t.Fatalf("reasoning option = %#v, want %#v", got.ReasoningOption, tt.wantReasoning)
			}
			if got.Adapter != tt.adapter {
				t.Fatalf("adapter = %#v, want %#v", got.Adapter, tt.adapter)
			}
		})
	}
}

func TestParseSessionCapabilitySnapshot(t *testing.T) {
	t.Parallel()

	snapshot := sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "xhigh", []string{"low", "medium", "high", "xhigh"})
	got, err := ParseSessionCapabilitySnapshot([]byte(snapshot), AdapterEvidence{Command: "adapter"})
	if err != nil {
		t.Fatalf("parse Agent Session capability snapshot: %v", err)
	}
	if got.CurrentModel != "gpt-5.5" || got.ReasoningOption == nil || got.ReasoningOption.CurrentValue != "xhigh" {
		t.Fatalf("unexpected capabilities: %#v", got)
	}
}

func TestParseSessionCapabilitySnapshotRejectsInvalidEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
		issue   string
	}{
		{name: "unknown schema", payload: `{"schema":"acpx.session.v2","acpx":{}}`, issue: CapabilityIssueMalformedResponse},
		{name: "missing options", payload: `{"schema":"acpx.session.v1","acpx":{"current_model_id":"gpt-5.5"}}`, issue: CapabilityIssueMissingOptions},
		{name: "contradictory model", payload: sessionCapabilitySnapshotFixture(t, "gpt-5.4", []string{"gpt-5.5"}, "", "", nil), issue: CapabilityIssueInvalidCurrentValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseSessionCapabilitySnapshot([]byte(tt.payload), AdapterEvidence{Command: "adapter"})
			var evidence *CapabilityEvidenceError
			if !errors.As(err, &evidence) || !containsString(evidence.Issues, tt.issue) {
				t.Fatalf("error = %T %v, want issue %q", err, err, tt.issue)
			}
		})
	}
}

func TestSelectionCapabilitiesIndependentAndVariantOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fixture       string
		wantModels    []ModelCapability
		wantReasoning bool
	}{
		{
			name:    "independent reasoning option",
			fixture: officialAdapterCapabilityFixture(),
			wantModels: []ModelCapability{
				{AdapterValue: "gpt-5.6-sol", CanonicalModel: "gpt-5.6-sol", ModelManaged: true},
				{AdapterValue: "gpt-5.6-terra", CanonicalModel: "gpt-5.6-terra", ModelManaged: true},
				{AdapterValue: "gpt-5.6-luna", CanonicalModel: "gpt-5.6-luna", ModelManaged: true},
				{AdapterValue: "gpt-5.5", CanonicalModel: "gpt-5.5", ModelManaged: true},
			},
			wantReasoning: true,
		},
		{
			name:    "advertised model variants",
			fixture: modelVariantCapabilityFixture(),
			wantModels: []ModelCapability{
				{AdapterValue: "future-model", CanonicalModel: "future-model", ModelManaged: true},
				{AdapterValue: "future-model[high]", CanonicalModel: "future-model", ReasoningEffort: "high"},
				{AdapterValue: "future-model[xhigh]", CanonicalModel: "future-model", ReasoningEffort: "xhigh"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSessionConfigOptions([]byte(tt.fixture), AdapterEvidence{Command: "adapter", Package: "example/adapter", Version: "1.0.0"})
			if err != nil {
				t.Fatalf("parse fixture: %v", err)
			}
			if !reflect.DeepEqual(got.Models, tt.wantModels) {
				t.Fatalf("models = %#v, want %#v", got.Models, tt.wantModels)
			}
			if (got.ReasoningOption != nil) != tt.wantReasoning {
				t.Fatalf("reasoning option presence = %t, want %t", got.ReasoningOption != nil, tt.wantReasoning)
			}
		})
	}
}

func TestPlanSelectionAssignment(t *testing.T) {
	t.Parallel()

	official, err := ParseSessionConfigOptions([]byte(officialAdapterCapabilityFixture()), AdapterEvidence{Command: "official-adapter"})
	if err != nil {
		t.Fatalf("parse official capabilities: %v", err)
	}
	variant, err := ParseSessionConfigOptions([]byte(modelVariantCapabilityFixture()), AdapterEvidence{Command: "variant-adapter"})
	if err != nil {
		t.Fatalf("parse variant capabilities: %v", err)
	}
	legacy, err := ParseSessionConfigOptions([]byte(legacyAdapterCapabilityFixture()), AdapterEvidence{Command: "legacy-adapter"})
	if err != nil {
		t.Fatalf("parse legacy capabilities: %v", err)
	}

	tests := []struct {
		name           string
		runtime        RuntimeSpec
		capabilities   SelectionCapabilities
		want           SelectionAssignment
		classification string
	}{
		{
			name:         "official Sol high uses independent control",
			runtime:      RuntimeSpec{ID: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"},
			capabilities: official,
			want: SelectionAssignment{
				Runtime:         "codex",
				Model:           "gpt-5.6-sol",
				ReasoningEffort: "high",
				AdapterModel:    "gpt-5.6-sol",
				ReasoningKey:    "reasoning_effort",
				ReasoningValue:  "high",
				Encoding:        SelectionEncodingIndependent,
			},
		},
		{
			name:         "official GPT 5.5 xhigh uses independent control",
			runtime:      RuntimeSpec{ID: "codex", Model: "gpt-5.5", ReasoningEffort: "xhigh"},
			capabilities: official,
			want: SelectionAssignment{
				Runtime:         "codex",
				Model:           "gpt-5.5",
				ReasoningEffort: "xhigh",
				AdapterModel:    "gpt-5.5",
				ReasoningKey:    "reasoning_effort",
				ReasoningValue:  "xhigh",
				Encoding:        SelectionEncodingIndependent,
			},
		},
		{
			name:         "advertised variant preserves canonical tuple",
			runtime:      RuntimeSpec{ID: "codex", Model: "future-model", ReasoningEffort: "high"},
			capabilities: variant,
			want: SelectionAssignment{
				Runtime:         "codex",
				Model:           "future-model",
				ReasoningEffort: "high",
				AdapterModel:    "future-model[high]",
				Encoding:        SelectionEncodingModelVariant,
			},
		},
		{
			name:         "empty effort is explicit model managed intent",
			runtime:      RuntimeSpec{ID: "codex", Model: "gpt-5.6-sol"},
			capabilities: legacy,
			want: SelectionAssignment{
				Runtime:      "codex",
				Model:        "gpt-5.6-sol",
				AdapterModel: "gpt-5.6-sol",
				Encoding:     SelectionEncodingModelManaged,
			},
		},
		{
			name:           "missing model is unsupported",
			runtime:        RuntimeSpec{ID: "codex", Model: "gpt-missing", ReasoningEffort: "high"},
			capabilities:   official,
			classification: SelectionModelNotAdvertised,
		},
		{
			name:           "missing independent control is unsupported",
			runtime:        RuntimeSpec{ID: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"},
			capabilities:   legacy,
			classification: SelectionReasoningControlNotAdvertised,
		},
		{
			name:           "missing advertised variant is unsupported",
			runtime:        RuntimeSpec{ID: "codex", Model: "future-model", ReasoningEffort: "ultra"},
			capabilities:   variant,
			classification: SelectionModelVariantNotAdvertised,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := PlanSelectionAssignment(tt.runtime, tt.capabilities)
			if tt.classification == "" {
				if err != nil {
					t.Fatalf("plan selection: %v", err)
				}
				if got != tt.want {
					t.Fatalf("assignment = %#v, want %#v", got, tt.want)
				}
				return
			}
			if err == nil {
				t.Fatal("expected unsupported selection")
			}
			var unsupported *SelectionUnsupportedError
			if !errors.As(err, &unsupported) {
				t.Fatalf("error type = %T, want *SelectionUnsupportedError", err)
			}
			if unsupported.Classification() != tt.classification {
				t.Fatalf("classification = %q, want %q", unsupported.Classification(), tt.classification)
			}
		})
	}
}

func TestParseSessionConfigOptionsRejectsInvalidEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fixture string
		issue   string
	}{
		{name: "malformed JSON", fixture: `{`, issue: CapabilityIssueMalformedResponse},
		{name: "missing config options", fixture: `{"action":"config_set","configId":"model","value":"gpt-5.5"}`, issue: CapabilityIssueMissingOptions},
		{name: "duplicate option ids", fixture: duplicateOptionIDCapabilityFixture(), issue: CapabilityIssueDuplicateOptionID},
		{name: "duplicate option values", fixture: duplicateOptionValueCapabilityFixture(), issue: CapabilityIssueDuplicateOptionValue},
		{name: "invalid current value", fixture: invalidCurrentValueCapabilityFixture(), issue: CapabilityIssueInvalidCurrentValue},
		{name: "contradictory set response", fixture: contradictoryCapabilityFixture(), issue: CapabilityIssueContradictoryResponse},
		{name: "ambiguous variants", fixture: ambiguousModelVariantCapabilityFixture(), issue: CapabilityIssueAmbiguousModelVariant},
		{name: "malformed nested variant", fixture: malformedModelVariantCapabilityFixture(), issue: CapabilityIssueMalformedModelValue},
		{name: "missing model state", fixture: missingModelCapabilityFixture(), issue: CapabilityIssueMissingModel},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseSessionConfigOptions([]byte(tt.fixture), AdapterEvidence{Command: "adapter"})
			if err == nil {
				t.Fatal("expected invalid capability evidence")
			}
			var evidenceErr *CapabilityEvidenceError
			if !errors.As(err, &evidenceErr) {
				t.Fatalf("error type = %T, want *CapabilityEvidenceError", err)
			}
			if evidenceErr.Classification() != CapabilityEvidenceInvalid {
				t.Fatalf("classification = %q, want %q", evidenceErr.Classification(), CapabilityEvidenceInvalid)
			}
			if !containsString(evidenceErr.Issues, tt.issue) {
				t.Fatalf("issues = %#v, want %q", evidenceErr.Issues, tt.issue)
			}
		})
	}
}

func TestCapabilityEvidenceAcquisitionFailureIsBounded(t *testing.T) {
	harness := newFakeACPXHarness(t)
	privatePath := "/Users/example/.codex/models_cache.json"
	secret := "API_TOKEN=do-not-report"
	t.Setenv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{"set model": 2}))
	t.Setenv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"set model": privatePath + "\n" + secret + "\n" + strings.Repeat("raw-adapter-output", 1000),
	}))

	_, err := harness.runner.AcquireSelectionCapabilities(context.Background(), CapabilityAcquisitionRequest{
		Runtime:  RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol"},
		Session:  SessionRef{Name: "roundfix-capability", WorkDir: harness.gitRoot},
		ConfigID: "model",
		Value:    "gpt-5.6-sol",
		Adapter:  AdapterEvidence{Command: "adapter"},
	})
	if err == nil {
		t.Fatal("expected capability acquisition failure")
	}
	var acquisitionErr *CapabilityAcquisitionError
	if !errors.As(err, &acquisitionErr) {
		t.Fatalf("error type = %T, want *CapabilityAcquisitionError", err)
	}
	for _, forbidden := range []string{privatePath, secret, "raw-adapter-output"} {
		if strings.Contains(err.Error(), forbidden) {
			t.Fatalf("acquisition error exposed private adapter output %q: %q", forbidden, err)
		}
	}
}

func TestCapabilityEvidenceIsBounded(t *testing.T) {
	t.Parallel()

	privatePath := "/Users/example/.acpx/sessions/private.json"
	secret := "ROUND_FIX_SECRET=do-not-report"
	raw := map[string]any{
		"action":        "config_set",
		"configId":      "model",
		"value":         "missing-current",
		"sessionPath":   privatePath,
		"environment":   secret,
		"rawSession":    strings.Repeat("unbounded-session-record", 1000),
		"configOptions": []any{},
	}
	payload := []byte(mustJSONForTest(t, raw))
	_, err := ParseSessionConfigOptions(payload, AdapterEvidence{Command: strings.Repeat("adapter-output", 1000)})
	if err == nil {
		t.Fatal("expected bounded evidence failure")
	}
	var evidenceErr *CapabilityEvidenceError
	if !errors.As(err, &evidenceErr) {
		t.Fatalf("error type = %T, want *CapabilityEvidenceError", err)
	}
	if len(evidenceErr.Issues) > MaxCapabilityDiagnosticIssues {
		t.Fatalf("issues = %d, want at most %d", len(evidenceErr.Issues), MaxCapabilityDiagnosticIssues)
	}
	if !sortStringsAreStable(evidenceErr.Issues) {
		t.Fatalf("issues are not in stable order: %#v", evidenceErr.Issues)
	}
	message := err.Error()
	for _, forbidden := range []string{privatePath, secret, "rawSession", "unbounded-session-record", strings.Repeat("adapter-output", 2)} {
		if strings.Contains(message, forbidden) {
			t.Fatalf("error exposed forbidden capability evidence %q: %q", forbidden, message)
		}
	}
}

func TestCapabilityAcquisitionDoesNotReadPrivateRuntimeState(t *testing.T) {
	harness := newFakeACPXHarness(t)
	t.Setenv(fakeACPXStdoutBy, mustJSONForTest(t, map[string]string{
		"set model": selectionStateFixture(t, "model", "gpt-5.6-sol", "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna", "gpt-5.5"}, "reasoning_effort", "high", []string{"low", "medium", "high", "xhigh", "max", "ultra"}),
	}))

	guard := filepath.Join(t.TempDir(), "private-runtime-state-is-not-a-directory")
	if err := os.WriteFile(guard, []byte("guard"), 0o600); err != nil {
		t.Fatalf("write private-state guard: %v", err)
	}
	t.Setenv("HOME", guard)
	t.Setenv("XDG_CONFIG_HOME", guard)

	runtime := RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol"}
	got, err := harness.runner.AcquireSelectionCapabilities(context.Background(), CapabilityAcquisitionRequest{
		Runtime:  runtime,
		Session:  SessionRef{Name: "roundfix-capability", WorkDir: harness.gitRoot},
		ConfigID: "model",
		Value:    runtime.Model,
		Adapter:  AdapterEvidence{Command: "npx -y @agentclientprotocol/codex-acp", Package: CodexAdapterPackage, Version: "1.1.4"},
	})
	if err != nil {
		t.Fatalf("acquire capabilities with private runtime paths guarded: %v", err)
	}
	if got.CurrentModel != runtime.Model {
		t.Fatalf("current model = %q, want %q", got.CurrentModel, runtime.Model)
	}

	want := [][]string{{
		"--cwd", harness.gitRoot,
		"--format", "json",
		"--json-strict",
		"codex", "set", "model", "gpt-5.6-sol", "-s", "roundfix-capability",
	}}
	invocations := readJSONInvocations(t, harness.invocationsPath)
	if !reflect.DeepEqual(invocations, want) {
		t.Fatalf("capability invocations = %#v, want %#v", invocations, want)
	}
	for _, invocation := range invocations {
		joined := strings.Join(invocation, " ")
		if strings.Contains(joined, ".acpx/sessions") || strings.Contains(joined, ".codex") {
			t.Fatalf("capability acquisition referenced private runtime state: %q", joined)
		}
	}
}

func officialAdapterCapabilityFixture() string {
	return `{"action":"config_set","configId":"reasoning_effort","value":"high","resumed":true,"configOptions":[{"id":"model","name":"Model","description":"ignored","category":"model","type":"select","currentValue":"gpt-5.6-sol","options":[{"value":"gpt-5.6-sol","name":"Sol"},{"value":"gpt-5.6-terra","name":"Terra"},{"value":"gpt-5.6-luna","name":"Luna"},{"value":"gpt-5.5","name":"GPT-5.5"}]},{"id":"reasoning_effort","name":"Reasoning effort","type":"select","currentValue":"high","options":[{"value":"low"},{"value":"medium"},{"value":"high"},{"value":"xhigh"},{"value":"max"},{"value":"ultra"}]},{"id":"brave_mode","name":"Brave mode","type":"boolean","currentValue":true}],"acpxRecordId":"private-record-id","agentSessionId":"private-agent-session-id"}`
}

func legacyAdapterCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.6-sol","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.6-sol","options":[{"value":"gpt-5.6-sol","name":"Sol"}]}]}`
}

func modelVariantCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"future-model[high]","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"future-model[high]","options":[{"value":"future-model"},{"value":"future-model[high]"},{"value":"future-model[xhigh]"}]}]}`
}

func duplicateOptionIDCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.5","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.5","options":[{"value":"gpt-5.5"}]},{"id":"model","category":"model","type":"select","currentValue":"gpt-5.5","options":[{"value":"gpt-5.5"}]}]}`
}

func duplicateOptionValueCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.5","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.5","options":[{"value":"gpt-5.5"},{"value":"gpt-5.5"}]}]}`
}

func invalidCurrentValueCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"unadvertised","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"unadvertised","options":[{"value":"gpt-5.5"}]}]}`
}

func contradictoryCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.5","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.6-sol","options":[{"value":"gpt-5.6-sol"},{"value":"gpt-5.5"}]}]}`
}

func ambiguousModelVariantCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.5[xhigh]","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.5[xhigh]","options":[{"value":"gpt-5.5[xhigh]"},{"value":"gpt-5.5[xhigh ]"}]}]}`
}

func malformedModelVariantCapabilityFixture() string {
	return `{"action":"config_set","configId":"model","value":"gpt-5.5[xhigh][max]","configOptions":[{"id":"model","category":"model","type":"select","currentValue":"gpt-5.5[xhigh][max]","options":[{"value":"gpt-5.5[xhigh][max]"}]}]}`
}

func missingModelCapabilityFixture() string {
	return `{"action":"config_set","configId":"reasoning_effort","value":"high","configOptions":[{"id":"reasoning_effort","type":"select","currentValue":"high","options":[{"value":"high"}]}]}`
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func sortStringsAreStable(values []string) bool {
	for index := 1; index < len(values); index++ {
		if values[index-1] > values[index] {
			return false
		}
	}
	return true
}
