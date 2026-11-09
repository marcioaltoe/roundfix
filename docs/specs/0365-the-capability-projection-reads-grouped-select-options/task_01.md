---
task: task_01
spec: 0365-the-capability-projection-reads-grouped-select-options
status: pending
type: backend
complexity: medium
---

# Task 01: The capability projection reads a grouped select option as its groups' values

## Overview

codex-acp 2.6.0 advertises its `mode` select as groups of values, and the
capability projection reads only a flat list, so every Codex profile fails
its proof with `invalid_current_value, invalid_select_option`. This Task
makes `projectSelectCapability` flatten a grouped select's values in
advertised order under the bounds of a flat select, keeps every other shape
refused (ADR-0375), and pins the two payloads recorded by the operator as
fixtures. It answers the operator record of 2026-11-09
([references](references/2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md)).

## Requirements

1. MUST copy `references/codex-acp-2.6.0-config-options.json` and
   `references/codex-acp-2.5.0-config-options.json` of this Spec, byte for
   byte, to `internal/agent/testdata/codex-acp-2.6.0-config-options.json`
   and `internal/agent/testdata/codex-acp-2.5.0-config-options.json`. Tests
   read only the `testdata` copies, never the Spec folder.
2. MUST add to `internal/agent/selection_capabilities.go` the group fields of
   `acpSelectValuePayload` and the helper `flattenSelectValues`, exactly as
   `_techspec.md` → Interfaces shows, and make `projectSelectCapability`
   read its value entries from that helper before its existing loop, so a
   flat list follows API Contract 1 and a grouped list follows API
   Contracts 2 and 3. The issue constants, `retainAdvertisedValues`,
   `ParseSessionConfigOptions`, `ParseSessionCapabilitySnapshot` and the
   `Modes` selection keep their behavior for every flat payload.
3. MUST add `internal/agent/selection_capabilities_grouped_test.go` with
   `TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads`, which calls
   `t.Parallel()` first, reads both fixtures, wraps each as
   `{"schema":"acpx.session.v1","acpx":{"current_model_id":"gpt-6.1-sol","config_options":<fixture>}}`
   and parses it with `ParseSessionCapabilitySnapshot`, adapter
   `AdapterEvidence{Command: "codex-acp"}` and retention
   `SelectionRetention{Model: "gpt-6.1-sol", ReasoningEffort: "high"}`, in
   two subtests named `2.6.0` and `2.5.0`, asserting:
   - no error for either payload;
   - for 2.6.0, `Modes` equals `SelectCapability{ID: "mode", CurrentValue: "permission-profile:localdev", Values: []string{"read-only", "workspace-write", "agent", "agent-full-access", "permission-profile:localdev"}, AdvertisedCount: 5}`;
   - for 2.5.0, `Modes` equals `SelectCapability{ID: "mode", CurrentValue: "agent", Values: []string{"read-only", "workspace-write", "agent", "agent-full-access"}, AdvertisedCount: 4}`;
   - `CurrentModel`, `Models`, `ReasoningOption` and every element of
     `Options` whose `ID` is not `mode` are `reflect.DeepEqual` between the
     two projections.
4. MUST add to the same file `TestGroupedSelectOptionsKeepTheEvidenceBounds`,
   which calls `t.Parallel()` first and runs one subtest per case. Each case
   is a `config_set` payload for `configId` `model` holding a valid flat
   `model` select (current `gpt-test`) and a `mode` select built from the
   case, except the last case, which groups the `model` select itself. Each
   refusal asserts a `*CapabilityEvidenceError` whose `Issues` contain the
   named issue; each acceptance asserts no error and the exact `Values`:
   - `empty group beside a values group`: groups `a` (no values) and `b`
     (`read-only`, `agent-full-access`), current `read-only` → accepted,
     `Modes.Values` = `read-only`, `agent-full-access`;
   - `every group empty`: groups `a` and `b`, both with no values →
     `invalid_select_option`;
   - `values and groups mixed`: a value entry `read-only` then group `b`
     (`agent`) → `invalid_select_option`;
   - `group without an identifier`: one group with no `group` field
     (`read-only`) → `invalid_select_option`;
   - `group identifier above the bound`: one group whose `group` is 65 bytes
     (`read-only`) → `invalid_select_option`;
   - `group nested in a group`: group `a` holding a group `b`
     (`read-only`) → `invalid_select_option`;
   - `value repeated across groups`: groups `a` (`read-only`) and `b`
     (`read-only`) → `duplicate_option_value`;
   - `current value in no group`: groups `a` (`read-only`) and `b`
     (`agent`), current `agent-full-access` → `invalid_current_value`;
   - `flattened count above the ceiling`: groups `a` with 2,049 values and
     `b` with 2,048 values (each under the ceiling, 4,097 together) →
     `too_many_option_values`;
   - `grouped model past the retention bound`: the `model` select grouped as
     `a` with 70 values `m-00` … `m-69` and `b` with the single value
     `gpt-test`, current `m-00`, parsed with
     `SelectionRetention{Model: "gpt-test"}` → accepted, the `model`
     option's `AdvertisedCount` is 71 and `Models` contains `gpt-test`.
5. MUST keep passing, unchanged: `TestSelectionCapabilitiesOfficialAndLegacyFixtures`,
   `TestParseSessionConfigOptionsRejectsInvalidEvidence`,
   `TestSelectionCapabilitiesProjectsTheSessionModeOption`, the
   `TestCharacterizationInvariant*` tests and the `TestRetention*` tests.
6. MUST NOT change any issue constant or its text, the capability error
   message, `applyFullAccess`, `chooseFullAccessMode`, any guide or skill,
   or any file outside this Task's declarations.

## Subtasks

- [ ] Copy the two recorded payloads into `internal/agent/testdata`.
- [ ] Add the group fields and `flattenSelectValues`, and call it from `projectSelectCapability`.
- [ ] Add the recorded-payload test and the bounds table.

## Acceptance Criteria

- [ ] The recorded 2.6.0 payload projects with the five `mode` values in advertised order, and every other option equals the 2.5.0 projection.
- [ ] Every shape outside the flat and grouped lists, and every flat bound applied to the flattened values, still refuses with its named issue.

## Context

- instruction: `docs/adr/0375-a-grouped-select-option-projects-its-groups-values-and-capability-evidence-stays-fail-closed.md`
- instruction: `docs/adr/0105-a-capability-projection-bounds-what-it-keeps-not-what-a-runtime-advertises.md`
- instruction: `docs/specs/0365-the-capability-projection-reads-grouped-select-options/references/codex-acp-2.6.0-config-options.json`
- instruction: `docs/specs/0365-the-capability-projection-reads-grouped-select-options/references/codex-acp-2.5.0-config-options.json`
- instruction: `internal/agent/selection_capabilities_test.go`
- instruction: `internal/agent/selection_capabilities_characterization_test.go`
- interface: `internal/agent/selection_capabilities.go`
- creates: `internal/agent/testdata/codex-acp-2.6.0-config-options.json`
- creates: `internal/agent/testdata/codex-acp-2.5.0-config-options.json`
- creates: `internal/agent/selection_capabilities_grouped_test.go`

## Consumers

- `projectSelectCapability` and `SelectionCapabilities` — consumers: `./internal/agent`

## Verification

- `cmp internal/agent/testdata/codex-acp-2.6.0-config-options.json docs/specs/0365-the-capability-projection-reads-grouped-select-options/references/codex-acp-2.6.0-config-options.json && cmp internal/agent/testdata/codex-acp-2.5.0-config-options.json docs/specs/0365-the-capability-projection-reads-grouped-select-options/references/codex-acp-2.5.0-config-options.json` — expected: exit 0. Before this Task the fixtures do not exist, so `cmp` fails; afterwards they are the recorded bytes.
- `for t in TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads TestGroupedSelectOptionsKeepTheEvidenceBounds; do grep -q "func $t(" internal/agent/selection_capabilities_grouped_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestSelectionCapabilitiesProjectTheRecordedCodexACPPayloads|TestGroupedSelectOptionsKeepTheEvidenceBounds|TestSelectionCapabilitiesOfficialAndLegacyFixtures|TestParseSessionConfigOptionsRejectsInvalidEvidence|TestSelectionCapabilitiesProjectsTheSessionModeOption|TestCharacterizationInvariant.*|TestRetention.*)$' ./internal/agent` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the recorded 2.6.0 payload projects, every bound refuses by name, and the existing projection tests pass unchanged.

## References

- `_prd.md` → Goals; User Stories 1 and 3; Core Features 1-3; Success Metrics 1, 2 and 4; Acceptance evidence
- `_techspec.md` → Interfaces; API Contracts 1-3; Testing Approach; Build Order 1
- ADR-0375, ADR-0105
