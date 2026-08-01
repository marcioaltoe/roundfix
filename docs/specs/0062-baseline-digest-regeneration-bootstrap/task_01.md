---
task: task_01
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
type: test
complexity: medium
---

# Task 01: Capture the diagnostic characterization corpus

## Overview

Every later Task in this Spec changes catalog validation, and the Spec's whole
promise is that validation outside regeneration keeps behaving exactly as it
does today. This Task captures that behavior as goldens before anything else
moves. It changes no product behavior — it is the gate the rest of the Spec is
measured against, and it only works if it lands first.

## Requirements

1. MUST record, for every catalog the Baseline package can load today, the
   complete set of diagnostics that load produces — code, subject, and detail —
   in a deterministic, order-stable form.
2. MUST cover at minimum the embedded catalog and every catalog fixture the
   package already exercises, so the corpus reflects the real validation
   surface rather than a sample.
3. MUST fail with a readable diff naming the changed codes when a later change
   alters any recorded diagnostic.
4. MUST be regenerable through an explicit flag so an intended change can be
   re-recorded deliberately, never silently.
5. MUST NOT change any production behavior, diagnostic, message, or exported
   API in this Task.

## Subtasks

- [ ] Enumerate the catalogs the package can load today.
- [ ] Record each catalog's diagnostics deterministically as a golden.
- [ ] Add the comparison test with a readable failure diff.
- [ ] Add the explicit regeneration flag for the corpus.

## Acceptance Criteria

- [ ] A test loads every covered catalog and compares its diagnostics against
      the recorded golden, passing on the unmodified tree.
- [ ] Deliberately altering one recorded diagnostic makes the test fail and
      name the changed code.
- [ ] The comparison is order-stable: running it repeatedly produces the same
      result without regenerating.
- [ ] The corpus can be re-recorded through an explicit flag and not otherwise.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_load.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/catalog_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; the corpus matches the current tree.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0 on a second consecutive run, proving the comparison is
  stable and does not rewrite its own golden.
- `grep -rq "TestCatalogDiagnosticCharacterization" internal/baseline` —
  expected: exit 0; the test exists at the declared name.
- `go test ./internal/baseline -count=1` — expected: exit 0; the existing suite
  is unaffected.

## References

- `_prd.md` → Core Features 4; Success Metrics (every catalog loads with the
  same diagnostics after the change).
- `_techspec.md` → Testing Approach: characterization corpus; Build Order 1.
