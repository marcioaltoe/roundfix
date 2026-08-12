---
task: task_03
spec: 0091-a-proof-that-can-refuse
status: completed
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
6. MUST break the characterization case Task 01 declared for `claude`, update
   it to the new behaviour in the same commit, and rename it to
   `TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel` so the
   name records what it now proves.

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
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofRefusalNamesTheAdvertisedSet$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofRefusalNamesTheAdvertisedSet'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofKeepsTheAdapterRefusalFastPath$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofKeepsTheAdapterRefusalFastPath'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel'` — expected: exits 0. The declared break must be updated in place and renamed to the behaviour it now records, so this name cannot exist before the work: asserting the old name passed against the unchanged tree and was refused as vacuous on 2026-08-10.

## References

- `_prd.md` → Goals 1, 2 and 5.
- `_techspec.md` → Build Order 3; System Architecture.
- ADR-0050, ADR-0112.

## Result

### Implementation

- Exact Agent Selection proof now checks the requested Agent Model against the
  pre-request Runtime Catalogue before planning or applying the echoed
  capability state.
- A catalogue-owned refusal uses `ModelNotAdvertisedError`, carries the ordered
  advertised Agent Model set, and retains the existing
  `model_not_advertised` classification.
- The adapter-owned refusal remains earlier in the proof flow and keeps its
  underlying adapter error and adapter-reported advertised set.
- Advertised selections retain their existing planning and application path;
  later capability states may still mark the catalogue contaminated while the
  proof remains `proven`.
- Renamed the Claude characterization case to
  `TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel` and
  changed it to assert refusal against the honest catalogue.

### Focused checks

- Pre-change,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/agent -run 'TestProof|TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel' -count=1`
  exited 1: the new refusal tests received the post-request
  `SelectionUnsupportedError` or continued past membership, establishing the
  missing catalogue verdict.
- After the implementation and final test edits, the same focused command
  exited 0.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/agent -count=1`
  exited 0 after the final code edit.
- `rtk git diff --check` exited 0 after the final code edit.
- The authored `## Verification` commands were not run; the Daemon owns them.
- Verification feedback attempt 1 showed that the shared
  `-run '^TestProofRefuses'` selector executed only the first refusal test. The
  second and third selectors now name their required top-level tests exactly.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/agent -run 'TestProof(RefusalNamesTheAdvertisedSet|KeepsTheAdapterRefusalFastPath)$' -count=1`
  exited 0 after the selector repair.

### Acceptance evidence

- `TestProofRefusesAModelTheCatalogueDoesNotAdvertise` supplies a Claude
  capability state that echoes the invented model while the honest catalogue
  omits it, and observes `ModelNotAdvertisedError` with
  `model_not_advertised`.
- `TestProofRefusalNamesTheAdvertisedSet` compares the typed refusal's
  `Advertised` field with the honest catalogue and asserts that its message
  contains the complete ordered set.
- `TestProofKeepsTheAdapterRefusalFastPath` makes only the requested-model
  ensure fail with the adapter's measured stderr, then confirms the returned
  refusal retains both the underlying adapter error and the adapter-reported
  model set.
- `TestProofKeepsAdvertisedSelectionEncoding` proves advertised `opus` / `high`
  through `opus[1m]` with the unchanged `independent` encoding. The focused
  package check also exercises the existing encoding and reasoning corpus.
- `TestRuntimeCatalogueRecordsAContaminatedAdvertisement` now requests an
  advertised model while a later capability payload adds an unrelated absent
  model; it observes `Contaminated` on a `proven` proof.
- `TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel`
  replaces the declared-break success assertion with the required refusal and
  advertised-set assertion.
