---
task: task_02
spec: 0091-a-proof-that-can-refuse
status: pending
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

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRuntimeCatalogueReadsAdvertisedModelsWithoutAnOverride'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRuntimeCatalogueBindsCanonicalVariant'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestRuntimeCatalogue' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRuntimeCatalogueRecordsAContaminatedAdvertisement'` — expected: exits 0.
- `grep -q 'RuntimeCatalogue' internal/agent/selection_assignment.go` — expected: exits 0. This string does not exist in the file before this Task.

## References

- `_prd.md` → Goal 3; Core Features, a catalogue read before the request.
- `_techspec.md` → Build Order 2; Interfaces.
- ADR-0112.
