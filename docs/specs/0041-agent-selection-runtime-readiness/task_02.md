---
task: task_02
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Add bounded capability and advertised-option values.
- [ ] Add the public ACPX machine-readable response seam.
- [ ] Parse documented Session configuration state.
- [ ] Validate ordering, bounds, ambiguity, and malformed evidence.
- [ ] Add official and legacy adapter capability fixtures.
- [ ] Assert that production acquisition never opens private runtime files.

## Acceptance Criteria

- [ ] An official-adapter fixture exposes Sol, Terra, Luna, GPT-5.5, and the
      advertised reasoning values in stable order.
- [ ] A legacy-adapter fixture exposes Sol without inventing a reasoning
      control or unadvertised Terra/Luna values.
- [ ] Independent options and model-variant fixtures produce unambiguous
      canonical capability projections.
- [ ] Duplicate option IDs, invalid current values, ambiguous variants, and
      missing required state fail with bounded typed evidence.
- [ ] Capability diagnostics contain no private paths, environment values, raw
      Session records, or unbounded adapter output.
- [ ] A filesystem guard proves production capability acquisition does not
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

