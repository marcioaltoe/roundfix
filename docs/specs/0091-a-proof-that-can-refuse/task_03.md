---
task: task_03
spec: 0091-a-proof-that-can-refuse
status: pending
type: backend
complexity: high
---

# Task 03: Let membership decide the verdict

## Overview

The slice that makes the proof able to refuse. A selection naming a model the
catalogue does not advertise is refused, on every runtime, with the advertised
set named. The adapter's own refusal is kept as a fast path where it exists,
because its wording is more precise than a membership verdict.

## Requirements

1. MUST refuse an Agent Selection whose model the catalogue from Task 02 does
   not advertise, on every runtime.
2. MUST name the advertised set in the refusal, so the message alone is enough
   to correct the profile.
3. MUST preserve the existing adapter-refusal path: when an adapter reports the
   model as unadvertised, its message is reported rather than a membership
   verdict.
4. MUST leave a selection whose model the catalogue advertises proving exactly
   as it does today, including its encoding and reasoning outcome.
5. MUST record the contaminated case as an advisory on a passing proof, never as
   a refusal.
6. MUST break the characterization case Task 01 declared for `claude`, and
   update it to the new behaviour in the same commit.

## Subtasks

- [ ] Decide membership against the catalogue.
- [ ] Name the advertised set in the refusal.
- [ ] Keep the adapter fast path and the passing paths unchanged.

## Acceptance Criteria

- [ ] A `claude` selection naming an unoffered model is refused.
- [ ] The refusal names the advertised set.
- [ ] A `codex` selection naming an unoffered model still reports the adapter's
      own message.
- [ ] Every advertised selection proves with the same encoding as before.
- [ ] A contaminated advertisement produces an advisory, not a refusal.

## Rehearsal Cases

- Case: `claude` / `opus-9-does-not-exist` / `high`, the exact tuple measured
  proving `passed` on 2026-08-09; Observation: proof refuses, and the message
  lists `default`, `opus[1m]`, `claude-fable-5[1m]`, `sonnet`, `haiku`.
- Case: `claude` / `opus` / `high`; Observation: proof passes with encoding
  `independent`, unchanged from before this Task.
- Case: a capability payload advertising a model absent from the catalogue;
  Observation: the proof passes and records the contaminated advisory.

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_assignment.go`
- `internal/agent/selection_assignment_test.go`
- `internal/agent/selection_catalogue_characterization_test.go`
- `docs/specs/0091-a-proof-that-can-refuse/task_03.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofRefuses' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofRefusesAModelTheCatalogueDoesNotAdvertise'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofRefuses' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofRefusalNamesTheAdvertisedSet'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofRefuses' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofKeepsTheAdapterRefusalFastPath'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestSelectionCatalogueCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSelectionCatalogueCharacterization'` — expected: exits 0, proving the declared break was updated to the new behaviour rather than left failing. A whole-package sweep would pass with the work absent; this names the case that must change.

## References

- `_prd.md` → Goals 1, 2 and 5.
- `_techspec.md` → Build Order 3; System Architecture.
- ADR-0050, ADR-0112.
