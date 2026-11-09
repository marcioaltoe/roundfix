package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		version string
		modes   SelectCapability
	}{
		{
			version: "2.6.0",
			modes: SelectCapability{
				ID:              "mode",
				CurrentValue:    "permission-profile:localdev",
				Values:          []string{"read-only", "workspace-write", "agent", "agent-full-access", "permission-profile:localdev"},
				AdvertisedCount: 5,
			},
		},
		{
			version: "2.5.0",
			modes: SelectCapability{
				ID:              "mode",
				CurrentValue:    "agent",
				Values:          []string{"read-only", "workspace-write", "agent", "agent-full-access"},
				AdvertisedCount: 4,
			},
		},
	}
	projections := make([]SelectionCapabilities, len(tests))
	for index, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			fixture, err := os.ReadFile(filepath.Join("testdata", "codex-acp-"+tt.version+"-config-options.json"))
			if err != nil {
				t.Fatalf("read recorded config options: %v", err)
			}
			payload := []byte(`{"schema":"acpx.session.v1","acpx":{"current_model_id":"gpt-6.1-sol","config_options":` + string(fixture) + `}}`)
			got, err := ParseSessionCapabilitySnapshot(payload, AdapterEvidence{Command: "codex-acp"}, SelectionRetention{
				Model: "gpt-6.1-sol", ReasoningEffort: "high",
			})
			if err != nil {
				t.Fatalf("project recorded config options: %v", err)
			}
			if !reflect.DeepEqual(got.Modes, &tt.modes) {
				t.Fatalf("modes = %#v, want %#v", got.Modes, &tt.modes)
			}
			projections[index] = got
		})
	}
	if t.Failed() {
		return
	}
	grouped, flat := projections[0], projections[1]
	if !reflect.DeepEqual(grouped.CurrentModel, flat.CurrentModel) {
		t.Errorf("current model changed: grouped %q, flat %q", grouped.CurrentModel, flat.CurrentModel)
	}
	if !reflect.DeepEqual(grouped.Models, flat.Models) {
		t.Errorf("models changed: grouped %#v, flat %#v", grouped.Models, flat.Models)
	}
	if !reflect.DeepEqual(grouped.ReasoningOption, flat.ReasoningOption) {
		t.Errorf("reasoning option changed: grouped %#v, flat %#v", grouped.ReasoningOption, flat.ReasoningOption)
	}
	otherOptions := make([][]SelectCapability, len(projections))
	for index, projection := range projections {
		for _, option := range projection.Options {
			if option.ID != "mode" {
				otherOptions[index] = append(otherOptions[index], option)
			}
		}
	}
	if !reflect.DeepEqual(otherOptions[0], otherOptions[1]) {
		t.Errorf("non-mode options changed: grouped %#v, flat %#v", otherOptions[0], otherOptions[1])
	}
}

func TestGroupedSelectOptionsKeepTheEvidenceBounds(t *testing.T) {
	t.Parallel()

	values := func(names ...string) []any {
		entries := make([]any, 0, len(names))
		for _, name := range names {
			entries = append(entries, map[string]any{"value": name})
		}
		return entries
	}
	group := func(id string, entries []any) map[string]any {
		return map[string]any{"group": id, "options": entries}
	}
	numberedValues := func(prefix string, count int) []any {
		entries := make([]any, 0, count)
		for index := 0; index < count; index++ {
			entries = append(entries, map[string]any{"value": fmt.Sprintf("%s%02d", prefix, index)})
		}
		return entries
	}
	tests := []struct {
		name      string
		entries   []any
		current   string
		issue     string
		want      []string
		model     bool
		retention SelectionRetention
	}{
		{
			name:    "empty group beside a values group",
			entries: []any{group("a", values()), group("b", values("read-only", "agent-full-access"))},
			current: "read-only",
			want:    []string{"read-only", "agent-full-access"},
		},
		{
			name:    "every group empty",
			entries: []any{group("a", values()), group("b", values())},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "values and groups mixed",
			entries: []any{map[string]any{"value": "read-only"}, group("b", values("agent"))},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "group without an identifier",
			entries: []any{map[string]any{"options": values("read-only")}},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "group identifier above the bound",
			entries: []any{group(strings.Repeat("g", maxCapabilityIDBytes+1), values("read-only"))},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "empty group identifier",
			entries: []any{group("", values("read-only"))},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "padded group identifier",
			entries: []any{group(" a ", values("read-only"))},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "group nested in a group",
			entries: []any{group("a", []any{group("b", values("read-only"))})},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name: "value inside a group carries a group field",
			entries: []any{group("a", []any{
				map[string]any{"value": "read-only", "group": "b"},
			})},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name: "value inside a group carries an options array",
			entries: []any{group("a", []any{
				map[string]any{"value": "read-only", "options": values()},
			})},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "grouped value above the bound",
			entries: []any{group("a", values("read-only", strings.Repeat("v", maxCapabilityValueBytes+1)))},
			current: "read-only",
			issue:   CapabilityIssueInvalidOption,
		},
		{
			name:    "value repeated across groups",
			entries: []any{group("a", values("read-only")), group("b", values("read-only"))},
			current: "read-only",
			issue:   CapabilityIssueDuplicateOptionValue,
		},
		{
			name:    "current value in no group",
			entries: []any{group("a", values("read-only")), group("b", values("agent"))},
			current: "agent-full-access",
			issue:   CapabilityIssueInvalidCurrentValue,
		},
		{
			name: "flattened count above the ceiling",
			entries: []any{
				group("a", numberedValues("a-", maxAdvertisedCapabilityValues/2+1)),
				group("b", numberedValues("b-", maxAdvertisedCapabilityValues/2)),
			},
			current: "a-00",
			issue:   CapabilityIssueTooManyValues,
		},
		{
			name:      "grouped model past the retention bound",
			entries:   []any{group("a", numberedValues("m-", 70)), group("b", values("gpt-test"))},
			current:   "m-00",
			model:     true,
			retention: SelectionRetention{Model: "gpt-test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := map[string]any{
				"id": "model", "category": "model", "type": "select",
				"currentValue": "gpt-test", "options": values("gpt-test"),
			}
			options := []any{model}
			configuredValue := "gpt-test"
			if tt.model {
				model["currentValue"] = tt.current
				model["options"] = tt.entries
				configuredValue = tt.current
			} else {
				options = append(options, map[string]any{
					"id": "mode", "type": "select", "currentValue": tt.current, "options": tt.entries,
				})
			}
			payload, err := json.Marshal(map[string]any{
				"action": "config_set", "configId": "model", "value": configuredValue, "configOptions": options,
			})
			if err != nil {
				t.Fatalf("encode grouped config options: %v", err)
			}
			got, err := ParseSessionConfigOptions(payload, AdapterEvidence{Command: "codex-acp"}, tt.retention)
			if tt.issue != "" {
				var evidence *CapabilityEvidenceError
				if !errors.As(err, &evidence) || !containsString(evidence.Issues, tt.issue) {
					t.Fatalf("error = %T %v, want issue %q", err, err, tt.issue)
				}
				return
			}
			if err != nil {
				t.Fatalf("project grouped config options: %v", err)
			}
			if tt.model {
				option := got.Options[selectCapabilityIndex(got.Options, "model")]
				if option.AdvertisedCount != 71 {
					t.Errorf("advertised model count = %d, want 71", option.AdvertisedCount)
				}
				wantValues := make([]string, 0, maxRetainedCapabilityValues)
				for index := 0; index < maxRetainedCapabilityValues-1; index++ {
					wantValues = append(wantValues, fmt.Sprintf("m-%02d", index))
				}
				wantValues = append(wantValues, "gpt-test")
				if !reflect.DeepEqual(option.Values, wantValues) {
					t.Errorf("retained model values = %v, want %v", option.Values, wantValues)
				}
				found := false
				for _, model := range got.Models {
					if model.AdapterValue == "gpt-test" {
						found = true
					}
				}
				if !found {
					t.Errorf("requested model absent from projection: %#v", got.Models)
				}
				return
			}
			if got.Modes == nil || !reflect.DeepEqual(got.Modes.Values, tt.want) {
				t.Fatalf("modes = %#v, want values %v", got.Modes, tt.want)
			}
		})
	}
}
