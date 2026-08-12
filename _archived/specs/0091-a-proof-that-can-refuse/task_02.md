---
task: task_02
spec: 0091-a-proof-that-can-refuse
status: completed
type: backend
complexity: medium
---

# Task 02: Read what a runtime offers before asking about it

## Overview

Establishes the honest catalogue. The disposable session the proof already
creates is ensured once without any model override, and what it advertises is
kept. The read costs one ACP round trip and no tokens, and it is taken before
the selection is applied so the answer cannot be written by the question.

## Requirements

1. MUST read the advertised models and reasoning efforts from a disposable
   session ensured without a model override.
2. MUST carry that catalogue alongside the existing capabilities rather than
   replacing them; applying and observing the selection is unchanged in this
   Task.
3. MUST bind a requested model to an advertised variant using the canonical
   rule already used by capability retention, so `opus` binds `opus[1m]`.
4. MUST record when a later read advertises a model absent from this catalogue,
   which is the contaminated case.
5. MUST send no prompt. Proof stays token-free.
6. MUST NOT change any proof verdict; this Task only obtains and carries the
   catalogue.

## Subtasks

- [ ] Ensure the disposable session without an override and read its options.
- [ ] Carry the catalogue and its contamination flag.
- [ ] Bind canonical variants.

## Acceptance Criteria

- [ ] The catalogue read returns the advertised models for a runtime with no
      override applied.
- [ ] A requested `opus` binds an advertised `opus[1m]`.
- [ ] A payload advertising a model absent from the catalogue sets the
      contaminated record.
- [ ] Every proof verdict is identical to before this Task.

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_assignment.go`
- `internal/agent/selection_assignment_test.go`
- `internal/agent/selection_catalogue_characterization_test.go`
- `docs/specs/0091-a-proof-that-can-refuse/task_02.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRuntimeCatalogueReadsAdvertisedModelsWithoutAnOverride'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRuntimeCatalogueBindsCanonicalVariant'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRuntimeCatalogueRecordsAContaminatedAdvertisement'` — expected: exits 0.
- `grep -q 'RuntimeCatalogue' internal/agent/selection_assignment.go` — expected: exits 0. This string does not exist in the file before this Task.

## References

- `_prd.md` → Goal 3; Core Features, a catalogue read before the request.
- `_techspec.md` → Build Order 2; Interfaces.
- ADR-0112.

## Result

### Implementation

- Added `RuntimeCatalogue` with advertised Agent Models, reasoning efforts, and
  a contamination record. `SelectionProof` and `SessionSelectionRequest` carry
  it alongside the existing `SelectionCapabilities`.
- `ProveExactSelection` now ensures and reads the disposable Agent Session once
  without a model override before the existing requested-selection ensure and
  application flow.
- Reused the capability-retention canonical binding rule for catalogue
  membership, so an advertised `opus[1m]` binds a requested `opus`.
- Compared every later capability read with the pre-request catalogue and
  retained `Contaminated` on the proof without consulting it for any verdict.

### Focused checks

- The pre-change `rtk rg -n '^type RuntimeCatalogue' internal/agent/selection_assignment.go`
  exited 1, establishing that the catalogue type and behavior were absent.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/agent -run 'Test(RuntimeCatalogue|SelectionCatalogueCharacterization)' -count=1`
  exited 0 after the final test edit.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/agent -count=1`
  exited 0 after the final test edit.
- `rtk git diff --check` exited 0.
- The first focused test attempt used the sandbox-blocked macOS global Go cache;
  the recorded checks use the isolated `/private/tmp` cache instead.
- The authored `## Verification` commands were not run; the Daemon owns them.

### Acceptance evidence

- `TestRuntimeCatalogueReadsAdvertisedModelsWithoutAnOverride` asserts that the
  first disposable ensure has no `--model`, the following command reads the
  Agent Session, the proof carries all fixture models and efforts, and no
  prompt command occurs.
- `TestRuntimeCatalogueBindsCanonicalVariant` asserts that `opus[1m]`
  advertises canonical `opus` and does not bind an absent `haiku`.
- `TestRuntimeCatalogueRecordsAContaminatedAdvertisement` supplies a later
  capability state containing a model absent from the catalogue, observes
  `Contaminated`, and confirms the proof status remains `proven`.
- The two pre-existing characterization cases still exercise their original
  verdicts in the focused check: Claude proves its echoed unoffered model and
  Codex refuses its unoffered model with `model_not_advertised`.
