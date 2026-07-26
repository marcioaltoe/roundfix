---
task: task_01
spec: 0049-baseline-preservation-idempotency
status: completed
type: backend
complexity: high
---

# Task 01: Make Baseline Preservation converge after semantic redistribution

## Overview

Deliver the complete Preservation correction from source inventory through
portable Plan assembly. The observable result is a first approved semantic
redistribution that removes consumed source bytes from the live `AGENTS.md`,
followed by empty unchanged reruns and the same convergence after a later
unmarked addition.

## Requirements

1. MUST classify only new nonblank unmarked root bytes, excluding
   setup-managed blocks and already retained repository-owned rule carriers.
2. MUST back up the complete pre-migration root carrier once for each distinct
   content identity before consuming any classified root byte.
3. MUST rebuild the live `AGENTS.md` from active owned regions after every
   source entry has a valid approved disposition.
4. MUST preserve exact disposition, retention, preimage, postimage, rollback,
   recovery, and backup validation.
5. MUST make an unchanged compatible Preservation produce zero file changes
   and no new backup.
6. MUST repeat the same review, backup, consumption, and convergence when a
   maintainer later adds new unmarked root bytes.
7. MUST keep public command, Plan, Result, confirmation, and exit-code schemas
   unchanged.

## Subtasks

- [x] Add failing inventory and end-to-end Plan tests for the reproduced
      duplicate-classification and second-backup behavior.
- [x] Separate complete root backup identity from unmarked classification
      entries.
- [x] Derive consumed root paths only after complete disposition validation
      or verified prior-apply evidence.
- [x] Rebuild consumed `AGENTS.md` postimages from active owned regions.
- [x] Recognize the prior buggy-apply state from a compatible Manifest and
      verified content-addressed backup without reclassifying its stale root
      payload.
- [x] Avoid duplicating source bytes already present in their approved
      semantic or residual owner.
- [x] Cover a later unmarked addition and a third zero-change Plan.
- [x] Preserve existing documentation claims with a contract assertion where
      applicable.
- [x] Run focused tests and the complete repository Verification.

## Acceptance Criteria

- [x] The first approved Preservation stores the exact original `AGENTS.md` in
      `AGENTS.<digest>.md`.
- [x] Every consumed source byte is absent from the live root and present in
      exactly its approved semantic or residual owner.
- [x] Repairing a prior buggy apply does not duplicate source bytes already
      present in their approved owner.
- [x] If a prior backed-up payload and a new unmarked addition coexist, only
      the new addition enters the Source Baseline.
- [x] Setup-managed blocks and already retained repository-owned rules do not
      appear as new Source Baseline Entries.
- [x] A fresh unchanged Preservation has no backup entry and no changed
      postimage.
- [x] Adding one new unmarked rule presents only that new source for
      classification and produces a backup for the new complete root identity.
- [x] Applying the later migration removes the new source from the live root,
      retains its approved owner, and makes the next Plan empty.
- [x] Missing dispositions, stale bytes, malformed markers, backup collisions,
      or invalid owners still fail before mutation.
- [x] Public schemas and existing rollback and recovery tests remain
      unchanged and pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/cli.md`
- instruction: `docs/agents/go.md`
- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/apply_test.go`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestPreservationSemanticRedistributionConverges|TestPreservationLaterUnmarkedAdditionConverges|TestPreservationPreviouslyBackedUpRootGuidanceIsNotReclassified|TestPreservationPreviouslyBackedUpRootGuidanceExposesOnlyLaterAddition|TestRootPreservation'` — expected: root inventory, first migration, backed-up migration repair, later addition, backup, and empty-rerun contracts pass.
- `rtk go test -count=1 ./internal/baseline` — expected: every Baseline unit and integration test passes.
- `rtk go test -count=1 ./...` — expected: all repository tests pass.
- `rtk go build -buildvcs=false ./...` — expected: all Roundfix packages build.

## References

- `_prd.md` → Goals; User Stories 1–3; Core Features 1–5.
- `_techspec.md` → System Architecture; Interfaces; Testing Approach; Build
  Order 1–4.
- ADR-0058; ADR-0064; ADR-0070; ADR-0071; ADR-0073; ADR-0074; ADR-0078.

## Result

Completed on 2026-07-26.

- Preservation now classifies only new nonblank unmarked bytes while complete
  carrier bytes continue to determine immutable backup identity.
- Approved root sources are removed from the live `AGENTS.md`; prior buggy
  applies are repaired from compatible Manifest plus verified backup evidence
  without reclassifying or duplicating their retained rules.
- A later unmarked addition remains independently reviewable, and unchanged
  reruns converge to an empty Change Plan.
- The live repository dogfood applied Plan Digest
  `sha256:0cffe761a6f33de76ca6dcc92865efcf75e1e65b0c2f5fba7d5a0069b5cce41f`
  with only the new backup and `AGENTS.md` update. The next Preservation used
  Plan Digest
  `sha256:93818af95d2c3c120699d790ea7c35b5c80c1be2dcf068d8262c9ab683cc86b1`
  and reported no file changes.

Verification evidence:

- Focused Preservation tests: 6 passed.
- `rtk go test -count=1 ./internal/baseline`: 419 passed.
- `rtk go test -count=1 ./...`: 2,369 passed in 22 packages.
- `rtk go build -buildvcs=false ./...`: passed.
- `rtk make verify`: passed, including 2,369 repository tests, 4 skill contract
  tests, `roundfix skills check`, and the final binary build.
