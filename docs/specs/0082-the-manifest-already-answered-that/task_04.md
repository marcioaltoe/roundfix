---
task: task_04
spec: 0082-the-manifest-already-answered-that
status: completed
type: backend
complexity: high
---

# Task 04: Refresh one repository from its own manifest

## Overview

Composes the previous two slices into the command a maintainer actually runs:
`roundfix baseline update` reads the manifest, plans a managed refresh, and
applies it only against a confirmed Plan Digest. This is the first slice that
delivers the PRD's headline outcome end to end — a stale repository reaches the
current catalog without answering a settled question — and it is demoable on a
fixture repository on its own, before skills are involved.

## Requirements

1. MUST add a `baseline update` command that resolves its complete input from
   the stored manifest and prompts for nothing.
2. MUST refuse, without writing, when the repository has no manifest, and direct
   the maintainer to first adoption.
3. MUST refuse, without writing, when the current catalog requires a decision the
   manifest does not carry, naming each such decision; an explicit
   suggestion-adopting flag instead takes the catalog's suggested value for
   exactly those decisions and reports every value adopted that way.
4. MUST present the plan and write nothing when no confirmation is supplied, and
   apply only against the exact current Plan Digest otherwise.
5. MUST support confirming the digest computed in the same invocation, for a
   scripted sweep, and confirming a digest reviewed in a previous invocation;
   the two forms MUST be mutually exclusive.
6. MUST emit a structured result naming the prior and current catalog identity,
   the file changes, the retention evidence and any clause delta, the warnings,
   and the approved Plan Digest, in both text and JSON.
7. MUST use exit categories consistent with the existing Baseline command family
   so a sweep can branch on them without parsing prose.
8. MUST be idempotent: refreshing an already-current repository reports zero file
   changes, and a second refresh immediately after a first reports zero file
   changes.
9. MUST NOT invoke a semantic analyzer or spawn an ACP runtime on any path.

## Subtasks

- [ ] Add the command, its flags, and their mutual exclusions.
- [ ] Wire manifest resolution to managed-refresh plan assembly.
- [ ] Implement the refusal paths: no manifest, new decisions, no confirmation.
- [ ] Implement suggestion adoption and its reporting.
- [ ] Emit the result document in text and JSON.
- [ ] Prove idempotence and every exit category by test.

## Acceptance Criteria

- [ ] Running the command against a fixture repository with a stale catalog and a
      complete manifest rewrites the managed artifacts and exits successfully,
      having asked nothing.
- [ ] Running it a second time immediately after reports zero file changes.
- [ ] Running it against a repository with no manifest writes nothing and exits
      in the action-required category naming adoption as the next action.
- [ ] Running it against a manifest missing a catalog-required decision writes
      nothing, exits in the action-required category, and names that decision id.
- [ ] Running the same case with the suggestion-adopting flag applies and lists
      every adopted value in the result.
- [ ] Running it without confirmation presents a Plan Digest and leaves the
      fixture repository byte-identical.
- [ ] Supplying both confirmation forms together is rejected as invalid input.
- [ ] No ACP runtime process is required for any of the above.

## Context

- interface: `internal/cli/baseline_plan_test.go`
- interface: `internal/cli/baseline_apply.go`
- interface: `internal/baseline/apply.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix baseline update --help` — expected: exits 0 and prints usage naming the update command.
- `go test ./internal/cli/ -run 'BaselineUpdate' -v 2>&1 | grep -q '^--- PASS: .*BaselineUpdate'` — expected: exits 0, proving the command's cases exist and pass.
- `go test ./internal/cli/ -run 'BaselineUpdate' -v 2>&1 | grep -q -i 'idempot'` — expected: exits 0, proving the idempotence case ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the task_01 corpus passing unchanged.

## References

- `_techspec.md` → Build Order 5; API Contracts; Data Models.
- `_prd.md` → Core Features 1, 4, 6, 9; User Stories 1 and 5; Goals 1 and 4; Success Metrics.
- ADR-0066, ADR-0068, ADR-0071, ADR-0073, ADR-0099.

## Result

Implemented the non-interactive `roundfix baseline update` slice. The command
resolves the stored Setup Manifest, reports or explicitly adopts catalog
suggestions for newly required decisions, builds a managed-refresh Plan, and
either presents it without writes or applies it against the exact current Plan
Digest. Text and JSON use the `roundfix/baseline-update-result/v1` result shape
with prior/current catalog identities, file changes, retention evidence, clause
delta, warnings, new/adopted decisions, Plan Digest, approved digest, and the
apply status matrix.

The suggestion-adoption flow exposed a lower-layer preservation defect when a
decision adds a managed root block: whitespace-only renderer separators were
counted as repository-authored regions. The validator now excludes only those
separator-only regions while continuing to digest every non-managed region
that contains authored content. The focused positive regression and the
existing prose-change, hand-edited-marker, and unaccounted-clause negative
tests all pass.

Focused checks used the worktree-local `.gocache` because the sandbox denied
the default macOS Go build cache:

- `gofmt -w` over the changed Go files in two focused invocations — both exit 0.
- Ten individual `go test ./internal/cli -run '^TestBaselineUpdate...$' -count=1` invocations — each exit 0. They covered apply, idempotence, no manifest, missing decisions, suggestion adoption, both confirmation forms, stale confirmation, help, exit categories, and the no-ACP path.
- Five individual managed-refresh checks in `./internal/baseline` — `TestManagedRefreshPreservationAllowsNewManagedRootBlock`, `TestManagedRefreshPreservesNonManagedRegionDigests`, `TestManagedRefreshApplyRejectsChangedNonManagedRegion`, `TestManagedRefreshBlocksHandEditedManagedMarker`, and `TestManagedRefreshUnaccountedClauseStillBlocks` — each exit 0.
- `go test ./internal/cli -run '^TestBaselineDocumentationContract$' -count=1` — exit 0.
- `git diff --check` — exit 0.

Acceptance evidence:

1. `TestBaselineUpdateAppliesManifestPlanAndReportsJSON` applies a stale
   manifest-backed fixture without an input stream or analyzer and observes
   managed guide plus Setup Manifest rewrites, exit 0, catalog identities,
   file changes, retention, warnings, and the approved digest.
2. `TestBaselineUpdateIdempotenceReportsZeroFileChanges` immediately invokes
   the command again, observes zero file changes and verified idempotence, and
   proves the repository tree remains byte-identical after the second call.
3. `TestBaselineUpdateNoManifestRequiresAdoptionWithoutWrites` observes exit 3,
   category `adoption`, a first-adoption next action, and an unchanged tree.
4. `TestBaselineUpdateNewDecisionRequiresActionWithoutWrites` removes
   `secondbrain.enabled` from the manifest and observes exit 3, the named
   decision, and an unchanged tree.
5. `TestBaselineUpdateAdoptsEverySuggestedDecision` repeats that case with
   `--adopt-suggested --yes`, observes exit 0, and records
   `secondbrain.enabled=true` in `adoptedSuggestions`.
6. `TestBaselineUpdatePresentsPlanAndConfirmsPreviousDigest` observes the text
   Plan, digest, file/retention evidence, and byte-identical tree without
   confirmation, then applies by passing that reviewed digest in a second
   invocation. `TestBaselineUpdateRejectsDigestOtherThanCurrentPlanWithoutWrites`
   proves any other digest exits 3 without writes.
7. `TestBaselineUpdateRejectsMutuallyExclusiveConfirmationForms` observes exit
   2 and a structured `invalid` result when `--yes` and `--confirm-plan` appear
   together.
8. The assembled success cases exercise manifest resolution, `BuildPlan`, and
   `ApplyPlan` without an analyzer seam. The update implementation imports no
   analyzer or ACP package, while
   `TestBaselineUpdateExitCategoriesAndNoACPRuntimeDependency` also proves
   cancellation maps to exit 130.

The authored `## Verification` commands were not run; the Daemon owns them and
Task settlement. Skill refresh remains outside this diff for Task 05.
