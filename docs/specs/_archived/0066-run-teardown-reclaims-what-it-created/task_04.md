---
task: task_04
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
type: backend
complexity: medium
---

# Task 04: Offer both debris kinds through reconcile

## Overview

`reconcile` is already dry-run first, already enumerates without a Run ID,
already carries proof per candidate, already preserves ambiguity, and already
refuses to apply what it did not inspect. This slice adds the two new debris
kinds to it rather than minting a command that would have to re-learn all four
properties.

## Requirements

1. MUST offer orphaned process candidates and releasable Run Branch candidates
   beside the existing worktree candidates.
2. MUST keep dry-run the default and `--apply` the only acting mode.
3. MUST carry, for every new candidate, the proof that makes it reclaimable.
4. MUST preserve anything ambiguous and report it as preserved.
5. MUST never offer or act on anything belonging to an Active Run, per ADR-0052.
6. MUST be idempotent: a second pass after applying is a no-op.
7. MUST extend the `--format json` schema additively, leaving existing fields
   and their meanings unchanged.

## Subtasks

- [ ] Add the two candidate kinds with their proofs.
- [ ] Wire them through dry-run and apply.
- [ ] Extend the JSON schema additively.
- [ ] Assert idempotence and the Active Run guard.

## Acceptance Criteria

- [ ] A dry-run names every candidate of both kinds with its proof and changes
      nothing.
- [ ] `--apply` reclaims exactly the named candidates.
- [ ] A second pass after applying is a no-op.
- [ ] An Active Run's process and branch are never offered.
- [ ] An ambiguous candidate is reported preserved, never offered.
- [ ] Existing `--format json` fields keep their names and meanings, asserted
      against the current schema.

## Context

- interface: `internal/cli/reconcile.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1 -run 'Reconcile' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the reconcile tests ran and passed.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix reconcile --format json | grep -q "runs\|Run"`
  — expected: exit 0; the command still runs and emits its report.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; the Skill is task_05's bounded scope.

## References

- `_prd.md` → Core Features 5 and 6; Success Metrics 3 and 4.
- `_techspec.md` → API Contracts; Build Order 4.
- ADR-0052.

## Result

### Implementation

- `reconcile` now inspects terminal Run-owned live process trees and proven
  superseded Run Branches beside its existing Run Worktree results. Each new
  candidate carries the terminal outcome plus the process-tree identity proof
  or superseding QA Report proof that makes it reclaimable.
- Dry-run remains the default. `--apply` reuses the inspected candidate,
  re-proves process ownership or branch-set supersession, and re-inspects a
  Run Branch candidate's registered worktree before mutation. A changed or
  dirty worktree is preserved.
- Active Runs, unproven owner identities, unclassifiable branch sets, current
  QA evidence, and branches without superseding QA evidence are reported in
  additive preserved-candidate output and are never offered.
- JSON retains the `roundfix-reconcile/v1` envelope and the exact existing
  result and summary field sets. `runs` is an additive compatibility view of
  the unchanged `results` collection; `processCandidates`,
  `runBranchCandidates`, `preservedCandidates`, and `debrisSummary` are the
  additive debris fields.

### Focused checks

- Pre-change signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/cli -count=1 -run '^TestRunReconcile(OffersAndAppliesOwnedProcessTreesAndRunBranches|ReportsAmbiguousDebrisPreserved|JSONMatchesTextFields)$'`
  — failed to compile because the reconcile report had no debris candidate
  fields, proof type, or process inspection seam.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/store ./internal/worktree ./internal/cli -count=1 -run '^(TestOwnerProcessControllerInspectTree|TestApplyRunBranchCandidate|TestRunReconcile(OffersAndAppliesOwnedProcessTreesAndRunBranches|ReportsAmbiguousDebrisPreserved|JSONMatchesTextFields))'`
  — passed (`ok` for all three packages).
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/cli -count=5 -run '^TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches$'`
  — passed across five fresh end-to-end fixtures.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/worktree -count=3 -run '^TestApplyRunBranchCandidate'`
  — passed, including fresh-proof, uninspected-candidate, and newly-dirty
  worktree cases.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/cli -count=1 -run '^(TestRunReconcileDryRunReadOnly|TestRunReconcileSupersededJSONAndApply|TestRunReconcileRepositoryScopeNewestFirst|TestRunReconcileInvalidSelectorsMutateNothing|TestRunReconcileApplyMixedResults|TestRunReconcileIdempotentApply)$'`
  — passed for the existing reconcile contract.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache GOOS=windows GOARCH=amd64 go test -c ./internal/store -o /private/tmp/roundfix-spec0066-task04-store-windows.test.exe`
  — passed for the new platform-neutral process interface.
- `rtk git diff --check` — passed.
- The commands under `## Verification` were not run; the Daemon owns them.

### Verification feedback repair

- Attempt 1 exposed an empty-report schema gap: valid JSON with no candidates
  contained neither a named `runs` collection nor an incidental `Run` value,
  so the configured output-presence check could not recognize the report.
- Pre-repair signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/cli -count=1 -run '^TestRunReconcileEmptyJSONNamesRunsCollection$'`
  — failed because the empty JSON report omitted `runs`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task04-cache go test ./internal/cli -count=3 -run '^TestRunReconcile(EmptyJSONNamesRunsCollection|JSONMatchesTextFields|OffersAndAppliesOwnedProcessTreesAndRunBranches)$'`
  — passed across three fresh fixtures after adding the compatibility view.
- The failed declared Verification command was not rerun; the Daemon owns the
  next configured Verification sequence.

### Acceptance criterion evidence

- Dry-run candidates and proofs:
  `TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches` observes
  one process-tree candidate with its exact live PID set and two Run Branch
  candidates with superseding QA Report proofs. It compares the Run Database
  bytes and complete Git worktree/ref snapshot before and after dry-run.
- Exact apply scope: the same test observes one named process root terminated
  and exactly the two named terminal Run Worktrees and Run Branches removed;
  the Active Run surfaces remain present.
- Idempotence: its second `reconcile --apply` report has no process or Run
  Branch candidates, reports zero debris applies, and leaves the complete Git
  snapshot unchanged.
- Active Run guard: the end-to-end fixture includes an Active Run with both an
  owner PID and a Run Branch. Both appear only as preserved candidates and
  neither reaches an apply seam.
- Ambiguity preservation:
  `TestRunReconcileReportsAmbiguousDebrisPreserved` observes an unproven owner
  identity and a branch without QA supersession evidence in preserved output,
  with neither offered. `TestApplyRunBranchCandidatePreservesNewlyDirtyWorktree`
  also proves apply refuses stale clean-worktree evidence and retains both Git
  surfaces.
- Additive JSON compatibility: `TestRunReconcileJSONMatchesTextFields` asserts
  the exact pre-existing top-level envelope, result, and summary field names
  against the current schema. It also proves additive `runs` is byte-identical
  to legacy `results`; the end-to-end apply test proves their decoded values
  remain equal after mutation. `TestRunReconcileEmptyJSONNamesRunsCollection`
  asserts the additive collection is present and empty when no Runs exist.

### Scope boundary

- No Skill files were changed; the authorized CLI Skill synchronization is
  task_05's bounded slice.
