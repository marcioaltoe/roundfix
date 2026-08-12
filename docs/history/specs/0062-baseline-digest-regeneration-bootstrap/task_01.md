---
task: task_01
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
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

## Result

### Implementation

- Added `TestCatalogDiagnosticCharacterization` at the existing catalog-loader
  test seam. It deterministically sorts and records each diagnostic's code,
  subject, and detail for 32 catalog inputs: the embedded and repository
  catalogs, the maintained legacy-v2 fixture, every existing catalog mutation,
  project-decision and tooling-authority variants, and the newly-required-clause
  case this Spec changes later.
- Added `testdata/catalog.diagnostics.golden.json`, containing 51 diagnostic
  records across 36 codes. The three valid catalogs are recorded with empty
  diagnostic arrays rather than omitted.
- Added a readable record-level comparison that reports the changed diagnostic
  codes and the removed/added triples.
- Added the dedicated `-update-catalog-diagnostics` flag. The golden write is
  reachable only inside that flag guard; an ordinary test run only reads and
  compares the corpus.
- Changed no production source, diagnostic, message, or exported API. Existing
  mutation and tooling-authority tables were only lifted into shared test
  fixtures so the characterization test exercises the same inputs.

### Focused-check evidence

- The first focused Go invocation could not use the host cache under
  `~/Library/Caches/go-build`; retrying with
  `GOCACHE=/private/tmp/roundfix-task01-gocache` removed that environment block.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test
  ./internal/baseline -run '^TestCatalogDiagnosticCharacterization$'
  -update-catalog-diagnostics` — exit 0 (`2 passed`) on two consecutive runs.
  Both runs produced SHA-256
  `7ec5328da2eefa03e857e9b7776a688c8a2cb7fb87fc34a7185cbc2d41618992`.
- After deliberately changing the recorded detail for
  `catalog.decision.dependency.cycle`, the ordinary focused comparison exited 1
  and reported `changed diagnostic codes:
  catalog.decision.dependency.cycle`. Regeneration through the dedicated flag
  restored the recorded hash above.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test
  ./internal/baseline -run
  '^TestCatalogDiagnosticCharacterization$/failure_diff_names_the_changed_code$'`
  — exit 0 (`2 passed`). The golden hash was identical before and after this
  no-update run.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test
  ./internal/baseline -run
  '^(TestCatalogMutation|TestToolingAuthorityCannotBeDisabled|TestToolingAuthorityAccounting)$'`
  — exit 0 (`29 passed`).
- `rtk git diff --check` — exit 0.

### Acceptance-criterion evidence

1. The focused characterization run loaded and compared all 32 enumerated
   inputs; the golden contains 29 diagnostic-bearing inputs and three valid
   inputs with empty diagnostic arrays.
2. The deliberate golden mutation failed the comparator and named
   `catalog.decision.dependency.cycle`; the in-memory negative check also
   requires both removed and added records in the diff.
3. Consecutive explicit generations and the intervening ordinary comparison
   retained the same SHA-256 shown above.
4. Source inspection found the golden writer only below the
   `updateCatalogDiagnosticCharacterization` guard, and the ordinary focused
   comparison left the golden unchanged.
5. `rtk git status --porcelain` reported only this task file,
   `internal/baseline/catalog_test.go`, and
   `internal/baseline/testdata/catalog.diagnostics.golden.json` (plus a benign
   fsmonitor IPC warning).

The Daemon-owned commands in `## Verification` were not run in this Agent turn.
