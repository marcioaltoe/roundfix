---
task: task_02
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: high
---

# Task 02: Expose bounded ACP selection capabilities

## Overview

Add the public ACP response seam and bounded data model needed to observe a
Session's advertised model and configuration controls. This prefactoring slice
must reject malformed or contradictory evidence without reading ACPX Session
files or private Codex caches.

## Requirements

1. MUST acquire model and configuration evidence from documented ACP
   `configOptions` responses exposed through a machine-readable ACPX boundary.
2. MUST project only current model, advertised model values, select-option IDs,
   current values, option values, and bounded adapter identity evidence.
3. MUST represent independent reasoning controls, advertised model variants,
   and explicit model-managed state without inferring capability from model
   marketing names.
4. MUST reject missing, malformed, duplicate, ambiguous, or contradictory
   capability evidence with a typed failure.
5. MUST apply deterministic ordering and bounds to every diagnostic list.
6. MUST not read `~/.acpx/sessions`, Codex model caches, environment secrets,
   or other runtime-private persistence as production authority.
7. MUST provide deterministic official and legacy adapter fixtures for later
   exact-selection tests.

## Subtasks

- [x] Add bounded capability and advertised-option values.
- [x] Add the public ACPX machine-readable response seam.
- [x] Parse documented Session configuration state.
- [x] Validate ordering, bounds, ambiguity, and malformed evidence.
- [x] Add official and legacy adapter capability fixtures.
- [x] Assert that production acquisition never opens private runtime files.

## Acceptance Criteria

- [x] An official-adapter fixture exposes Sol, Terra, Luna, GPT-5.5, and the
      advertised reasoning values in stable order.
- [x] A legacy-adapter fixture exposes Sol without inventing a reasoning
      control or unadvertised Terra/Luna values.
- [x] Independent options and model-variant fixtures produce unambiguous
      canonical capability projections.
- [x] Duplicate option IDs, invalid current values, ambiguous variants, and
      missing required state fail with bounded typed evidence.
- [x] Capability diagnostics contain no private paths, environment values, raw
      Session records, or unbounded adapter output.
- [x] A filesystem guard proves production capability acquisition does not
      access ACPX or Codex private persistence.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/acp_stream.go`
- interface: `internal/agent/sessions.go`

## Verification

- `rtk go test ./internal/agent -run 'Test(SelectionCapabilities|ParseSessionConfigOptions|CapabilityEvidence)' -count=1` — expected: official, legacy, independent-option, model-variant, malformed, ambiguous, and bounded-output cases pass.
- `rtk go test ./internal/agent -run 'TestCapabilityAcquisitionDoesNotReadPrivateRuntimeState' -count=1` — expected: the production path uses only public ACP/ACPX responses.
- `rtk go test -race ./internal/agent -run 'Test(SelectionCapabilities|ParseSessionConfigOptions|CapabilityEvidence|CapabilityAcquisition)' -count=1` — expected: capability acquisition and parsing are race-free.

## References

- `_prd.md` → User Stories 3 and 4; Core Features 2, 3, and 10; Non-Goals.
- `_techspec.md` → Data Model; Capability Acquisition; Testing Approach; Build
  Order 2.
- `references/validation.md` → official and legacy adapter evidence.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → advertised capabilities are authoritative.

## Result

Added a bounded `SelectionCapabilities` projection for documented ACP select
configuration state. It preserves advertised model and per-option value order,
sorts select options by ID, sorts and bounds typed diagnostic issue codes,
represents base values as explicit model-managed state, and recognizes
reasoning only from independent advertised controls or explicit bracketed
model variants. Display names, descriptions, boolean controls, raw Session
identifiers, and unknown response fields do not enter the projection.

Added `AcquireSelectionCapabilities`, which calls ACPX's public strict-JSON
`set <config-id> <value>` boundary and parses the complete returned
`configOptions` state. ACPX command failures expose a typed bounded error while
retaining the underlying error only for programmatic inspection. The
production path has no filesystem dependency for ACPX Session files, Codex
caches, or runtime-private metadata.

Acceptance evidence:

- `TestSelectionCapabilitiesOfficialAndLegacyFixtures` proves the official
  fixture exposes Sol, Terra, Luna, GPT-5.5, and all six advertised reasoning
  values in source order, while the legacy fixture exposes only Sol and no
  reasoning control.
- `TestSelectionCapabilitiesIndependentAndVariantOptions` proves independent
  reasoning controls and explicit bracketed model variants yield canonical,
  unambiguous projections with explicit model-managed base values.
- `TestParseSessionConfigOptionsRejectsInvalidEvidence` covers malformed JSON,
  missing options/model state, duplicate IDs and values, invalid current
  values, contradictory responses, ambiguous variants, and malformed variants
  through `CapabilityEvidenceError`.
- `TestCapabilityEvidenceIsBounded` and
  `TestCapabilityEvidenceAcquisitionFailureIsBounded` prove diagnostics are
  sorted, capped, and exclude private paths, environment values, raw Session
  records, and unbounded ACPX output.
- `TestCapabilityAcquisitionDoesNotReadPrivateRuntimeState` makes `HOME` and
  `XDG_CONFIG_HOME` unusable as directories, still acquires the fixture through
  one public ACPX JSON invocation, and asserts that no private runtime path is
  passed to ACPX.

Verification:

- `rtk go test ./internal/agent -run 'Test(SelectionCapabilities|ParseSessionConfigOptions|CapabilityEvidence)' -count=1`: passed, 18 tests.
- `rtk go test ./internal/agent -run 'TestCapabilityAcquisitionDoesNotReadPrivateRuntimeState' -count=1`: passed, 1 test.
- `rtk go test -race ./internal/agent -run 'Test(SelectionCapabilities|ParseSessionConfigOptions|CapabilityEvidence|CapabilityAcquisition)' -count=1`: passed, 19 tests.
- `rtk make verify`: passed with 1,602 Go tests, 79 setup-context-driven tests,
  Roundfix skill synchronization, and the CLI build.

Follow-up: deterministic assignment planning and exact Session application
remain Task 03's slice.
