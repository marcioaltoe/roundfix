---
task: task_02
spec: 0052-claude-adapter-standardization
status: pending
type: backend
complexity: high
---

# Task 02: Make advertised model identifiers opaque

## Overview

Stop parsing bracketed advertised Agent Model identifiers as Roundfix's
`canonical[effort]` variant encoding when the adapter advertises an
independent reasoning control. After this slice, `opus[1m]` and `opus` both
select the 1M-context Opus with an explicit reasoning effort, a
context-window annotation can no longer masquerade as a reasoning effort, and
selection diagnostics list only selectable identifiers. Verifiable entirely
through the capability-parsing and assignment-planning unit seams.

## Requirements

1. MUST resolve the independent reasoning option before parsing the
   advertised model list, so model parsing can depend on its presence.
2. MUST, when an independent reasoning control is advertised, parse a
   bracketed identifier as an opaque model-managed entry: the advertised
   value stays selectable, the bracket-stripped prefix becomes its canonical
   alias, and no reasoning effort is extracted.
3. MUST keep the existing `canonical[effort]` variant parsing for adapters
   that advertise no independent reasoning control.
4. MUST let a requested model match either the canonical alias or the raw
   advertised value, so an identifier copied from a diagnostic always
   matches.
5. MUST reject a selection whose reasoning effort is not advertised by the
   independent control — a context-window annotation such as `1m` — with the
   advertised reasoning efforts in the message.
6. MUST render each advertised model in selection failures as its canonical
   form, appending the raw advertised value when the two differ, so no
   printed identifier is unselectable by construction.
7. SHOULD keep the existing dedup fail-closed behavior when a canonical
   alias collides with another advertised entry.

## Subtasks

- [ ] Reorder capability projection so the reasoning option is known before
      the model loop, covering both the config-options and session-snapshot
      entry points.
- [ ] Implement the opaque parse rule in the model-capability parser.
- [ ] Extend canonical matching to include the raw advertised value.
- [ ] Render dual-form identifiers in the unsupported-selection diagnostic.
- [ ] Add a capability fixture with bracketed identifiers plus an independent
      reasoning option, and assignment-planning cases for `opus`, `opus[1m]`,
      and the rejected `1m` effort; keep the existing variant-encoding
      fixture passing unchanged.

## Acceptance Criteria

- [ ] With an advertised list containing `opus[1m]` and an independent
      `effort` option advertising `xhigh`, planning `opus` with `xhigh` and
      `opus[1m]` with `xhigh` both produce the independent encoding applying
      `xhigh`.
- [ ] Planning `opus` with reasoning effort `1m` against that fixture is
      rejected, and the failure message names the advertised reasoning
      efforts.
- [ ] The failure message's advertised-model list contains no identifier
      that would itself be rejected if requested verbatim.
- [ ] A fixture with `future-model[high]` and no independent reasoning
      option still produces the model-variant encoding.
- [ ] Both capability entry points (config options and session snapshot)
      apply the same rule.

## Context

- interface: `internal/agent/selection_capabilities.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/selection_capabilities_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/agent/ -run 'TestSelectionCapabilities|TestPlanSelectionAssignment|TestParseSession'` — expected: pass, including the new opaque-identifier cases.
- `go test -count=1 ./internal/agent/` — expected: full package passes, proving no regression in the variant and model-managed encodings.

## References

`_prd.md` → User Stories 3–4, Core Features 4–5; `_techspec.md` → Build
Order 2, Interfaces: parseModelCapability; ADR-0079, ADR-0055.
