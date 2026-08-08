---
task: task_02
spec: 0082-the-manifest-already-answered-that
status: completed
type: backend
complexity: high
---

# Task 02: Refresh managed regions without asking who owns the prose

## Overview

Adds the third instruction-preservation mode: a managed refresh that regenerates
only Baseline-owned regions and leaves every other byte identical. This is the
slice that removes the supervised analyzer from the update path — in the mode,
planning produces no source baseline and no decision skeleton, which are already
the two conditions every caller uses to skip classification. It is verifiable on
its own through planning alone, before any new command exists.

## Requirements

1. MUST add a managed-refresh preservation mode accepted by root-preservation
   planning alongside the existing greenfield and preservation modes.
2. MUST produce, in that mode, a ready preservation plan carrying no source
   baseline entries and no decision skeleton, so no classification input is ever
   required and no semantic analyzer is invoked.
3. MUST still detect and report blocking carriers; an unsafe root carrier blocks
   a managed refresh exactly as it blocks adoption.
4. MUST treat a hand-edited managed marker as blocking rather than as a warning,
   because the mode's guarantee depends on the markers being trustworthy.
5. MUST plan no root instruction backup in this mode, per ADR-0100, and instead
   let the plan's preimages carry the preservation proof.
6. MUST leave the greenfield and preservation modes behaviorally unchanged.
7. MUST prove, by test, that every byte outside a managed marker is identical
   before and after a managed-refresh plan is applied — including authored prose
   and repository-rule blocks — asserting on region digests rather than a golden
   file.
8. MUST leave retention accounting on this path untouched, so an unaccounted
   managed clause still blocks planning.

## Subtasks

- [x] Add the mode and accept it in root-preservation planning.
- [x] Return a ready plan with no source baseline and no decision skeleton.
- [x] Make hand-edited managed markers blocking in this mode.
- [x] Suppress root backup planning in this mode.
- [x] Build the mixed fixture: managed markers, authored prose, repository rules.
- [x] Assert non-managed region digests survive a plan and apply unchanged.
- [x] Assert an unaccounted managed clause still blocks in this mode.

## Acceptance Criteria

- [x] A managed-refresh plan on an adopted repository reports a ready state with
      zero source baseline entries and no decision skeleton.
- [x] Applying a managed-refresh plan against a stale catalog leaves every
      non-managed region byte-identical, proven by digest comparison.
- [x] A repository with a hand-edited managed marker blocks in managed-refresh
      mode and names the offending path.
- [x] A managed-refresh plan emits no root backup postimage.
- [x] A managed-refresh plan whose transition leaves a managed clause unaccounted
      does not become applicable and names the unaccounted clause.
- [x] The task_01 corpus still passes, proving adoption is unchanged.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'ManagedRefresh' -v > /tmp/task_02-1.log 2>&1 && grep -q '^--- PASS: .*ManagedRefresh' /tmp/task_02-1.log` — expected: exits 0, proving the new mode's cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'ManagedRefresh.*Preserv|Preserv.*ManagedRefresh' -v > /tmp/task_02-2.log 2>&1 && grep -q '^--- PASS: ' /tmp/task_02-2.log` — expected: exits 0, proving the byte-identity assertion ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the task_01 corpus passing unchanged.

## References

- `_techspec.md` → Build Order 2 and 3; Interfaces: `PreservationModeManagedRefresh`; Testing Approach.
- `_prd.md` → Core Features 2 and 3; User Stories 3 and 6; Goals 2 and 5.
- ADR-0058, ADR-0069, ADR-0070, ADR-0099, ADR-0100.

## Result

### Implementation

- Added `managed-refresh` as a third root-preservation mode. It returns before
  source-baseline, Decision Document, and backup construction while retaining
  repository inventory and blocking-carrier checks.
- Bound managed-marker trust to the adopted Setup Manifest. A changed, missing,
  malformed, or duplicated adopted marker blocks planning and reports its path;
  catalog drift alone does not make unchanged adopted marker bytes unsafe.
- Added the managed-refresh mode to portable plans only on this path. Existing
  greenfield and preservation plan JSON remains unchanged because the field is
  omitted for those modes.
- Replaced root-backup enforcement only for managed-refresh plans with a
  preimage-bound region-digest proof. Planning and first apply both compare the
  ordered SHA-256 digests of every byte region outside setup-owned markers;
  Setup Manifest bytes remain the sole whole-file managed exception.
- Added filesystem-backed plan/apply coverage with managed markers, CRLF
  authored root prose, authored guide prose, and a repository-rule block. The
  suite also rejects a re-digested plan that attempts to change authored bytes.

### Focused checks

- The pre-change focused run failed at compilation because
  `PreservationModeManagedRefresh` was undefined, establishing the red signal.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-02-final-gocache go test ./internal/baseline -run '^(TestManagedRefreshPlanNeedsNoClassificationInputOrBackup|TestManagedRefreshUnsafeRootCarrierStillBlocks|TestManagedRefreshPreservesNonManagedRegionDigests|TestManagedRefreshApplyRejectsChangedNonManagedRegion|TestManagedRefreshBlocksHandEditedManagedMarker|TestManagedRefreshUnaccountedClauseStillBlocks)$' -count=1`
  exited 0: `ok roundfix/internal/baseline`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-02-final-gocache go test ./internal/baseline -run '^(TestGreenfieldPlanBacksUpWithoutImport|TestPreservationRequiresEveryDisposition|TestPreservationPlanAcceptsCompleteDecisionDocument|TestDecisionDocumentSkeletonDoesNotProposeManagedSemanticVersionBytes|TestRootBackupIdentityRejectsCollisions|TestBaselinePlanCharacterization|TestBaselineRetentionAccountingCharacterizationCorpus|TestPlanDocumentStrictCodecs)$' -count=1`
  exited 0: `ok roundfix/internal/baseline`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-02-final-cli-gocache go test ./internal/cli -run '^(TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus|TestBaselinePlanCommandEmitsPortableJSONAndNormalizesDecisionFiles)$' -count=1`
  exited 0: `ok roundfix/internal/cli`.
- `rtk git diff --check` exited 0.
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance evidence

- `TestManagedRefreshPlanNeedsNoClassificationInputOrBackup` observes ready
  state, zero Source Baseline entries, a nil Decision Skeleton, and zero root
  backups. `TestManagedRefreshUnsafeRootCarrierStillBlocks` proves the existing
  unsafe-root finding remains blocking.
- `TestManagedRefreshPreservesNonManagedRegionDigests` builds and applies a
  stale-catalog plan, round-trips its portable mode field, asserts no backup
  ledger entry or backup postimage, and compares every non-managed region digest
  before and after apply. The mixed fixture includes authored prose and a
  repository-rule block.
- `TestManagedRefreshApplyRejectsChangedNonManagedRegion` proves apply repeats
  the preservation check instead of trusting a self-declared mode.
- `TestManagedRefreshBlocksHandEditedManagedMarker` reports
  `docs/agents/agent-instructions.md` with
  `baseline.preservation.managed-marker.modified` and produces no applicable
  plan.
- `TestManagedRefreshUnaccountedClauseStillBlocks` produces no applicable plan,
  retains the `classification` action category, and names
  `clause.core.keep-follow-ups-outside-slice` as `unaccounted`.
- The focused baseline and CLI characterization checks named above exercise the
  Task 01 corpus without changing its recorded goldens, while the selected
  greenfield and preservation cases retain their prior backup and disposition
  behavior.
