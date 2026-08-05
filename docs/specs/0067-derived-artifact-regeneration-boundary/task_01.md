---
task: task_01
spec: 0067-derived-artifact-regeneration-boundary
status: completed
type: backend
complexity: high
---

# Task 01: Declare one owner per derived path and prove it exhaustive

## Overview

`DERIVED_DIGEST_PATHS` lists directories, and the natural reading is that the
sanctioned command owns everything beneath them. It does not. This slice
replaces that inference with a record read at the path, and a test that makes
an unclassified artifact a failure.

Exhaustiveness is the deliverable. Three artifacts already inherited the
ambiguity; the test is what stops the next one.

## Requirements

1. MUST place an ownership record in each directory under the derived scan,
   naming exactly one owner: `sanctioned`, `dedicated`, or `frozen`.
2. MUST require a `reason` on every record, and a `command` carrying the exact
   invocation on every `dedicated` record.
3. MUST classify the frozen parity corpus as `frozen`, recording that
   repointing it was tried on 2026-07-30 and reverted, so the next reader stops
   before trying again.
4. MUST classify the plan-characterization and catalog-diagnostic corpora as
   `dedicated`, each carrying its exact flag, because both flag names differ
   from their test names and have been guessed wrong.
5. MUST add a test that fails when any path under the derived scan resolves to
   zero records or more than one.
6. MUST NOT change any artifact's content or any digest value. The derived tree
   is byte-identical after this Task.

## Subtasks

- [ ] Define the record shape and its validation.
- [ ] Write one record per directory under the derived scan.
- [ ] Add the exhaustiveness test.
- [ ] Assert the derived tree is byte-identical.

## Acceptance Criteria

- [ ] Every path under the derived scan resolves to exactly one ownership
      record.
- [ ] An artifact added without a record fails the exhaustiveness test, proven
      by a fixture that adds one and asserts the failure.
- [ ] Every record carries a reason; every `dedicated` record carries its exact
      command.
- [ ] The parity corpus is `frozen` and its reason names the 2026-07-30 revert.
- [ ] No digest value and no artifact content changed, asserted over the whole
      derived tree.

## Context

- instruction: `docs/findings/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md`
- instruction: `docs/findings/2026-08-01-characterization-corpus-is-outside-the-regeneration-command.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'Ownership|Exhaustive|Derived' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the ownership tests ran and passed.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "^internal/baseline/(assets|testdata)/" | grep -v "_ownership" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only ownership records were added under the derived tree,
  and no artifact content moved.

## References

- `_prd.md` → Core Features 1, 3 and 5; Decisions.
- `_techspec.md` → Interfaces; Build Order 1.
- ADR-0081, ADR-0085.

## Result

### Implementation

- Added strict ownership-record parsing and validation for the three allowed
  owners. Every record requires a reason, and `dedicated` records also require
  a command.
- Added one ownership record at every `DERIVED_DIGEST_PATHS` root, narrower
  records for the plan-characterization, parity, ADR-lifecycle, and legacy-v2
  corpora, and a file sidecar for the catalog-diagnostic corpus. The nearest
  directory record governs a subtree; an exact file sidecar governs that file.
- Added an exhaustive read-only scan driven by the Makefile's actual
  `DERIVED_DIGEST_PATHS`, including negative fixtures for zero and multiple
  matching records.

### Focused checks

- Red signal before adding records:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go test ./internal/baseline -run '^TestDerivedOwnershipIsExhaustive$' -count=1`
  failed because `assets/setups` resolved to zero ownership records.
- After implementation,
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go test ./internal/baseline -run '^TestDerivedOwnership' -count=1`
  exited 0.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go test ./internal/baseline -run '^(TestCatalogCompatibility|TestFormatterComposition|TestReadoptionCompatibilityMaintainedFixture|TestBaselineCompatibilityCorpus)$' -count=1`
  exited 0; the existing catalog, formatter, Source Baseline, and parity
  compatibility checks still accept the unchanged artifacts and digests.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go vet ./internal/baseline`
  exited 0.
- `rtk git diff --name-only HEAD -- internal/baseline/assets internal/baseline/testdata`
  exited 0 with no output; all changes under the derived tree are new
  `_ownership` metadata, not edits to tracked artifacts.

### Acceptance evidence

1. `TestDerivedOwnershipIsExhaustive` reads the Makefile scan roots and passed
   only after every scanned path resolved to one record.
2. `TestDerivedOwnershipRejectsUnclassifiedArtifact` adds an unclassified
   artifact to an in-memory derived scan and asserts the zero-record failure;
   `TestDerivedOwnershipRejectsMultipleRecords` covers the ambiguous case.
3. `TestDerivedOwnershipValidatesRecords` covers the owner enum, mandatory
   reason, and mandatory dedicated command. The real-tree scan validates every
   record, and `TestDerivedOwnershipDeclaresKnownBoundaries` asserts both exact
   dedicated invocations.
4. `TestDerivedOwnershipDeclaresKnownBoundaries` asserts that the parity corpus
   is frozen and that its reason contains `2026-07-30`, `tried`, and `reverted`.
5. `TestDerivedOwnershipValidationPreservesWholeTree` snapshots every scanned
   byte before and after validation. The focused compatibility checks passed,
   and the tracked derived-tree diff is empty.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and terminal settlement.
