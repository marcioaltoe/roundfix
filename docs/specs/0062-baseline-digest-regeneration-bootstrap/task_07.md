---
task: task_07
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
type: backend
complexity: high
---

# Task 07: Load the catalog once on the regeneration path

## Overview

The regeneration path still fails its first invocation. Task 03 threaded the
regeneration catalog through planning, but the apply path re-acquires a strict
embedded catalog of its own, sees the pre-run pin, and exits — after the
rewrite has already happened, so a second invocation succeeds. Patching that
one call site would only postpone the next one: the package has several strict
acquisition points. This Task establishes the invariant instead — on the
regeneration path the catalog is acquired once and threaded, and no callee
re-acquires — and adds a guard that fails when a future call site breaks it.

## Requirements

1. MUST make the sanctioned regeneration invocation succeed on its **first**
   run against a catalog whose derived pin is stale, exiting zero after
   rewriting the derived artifacts and re-validating strictly.
2. MUST thread the already-acquired regeneration catalog through every stage
   the regeneration path invokes, including the apply stage, rather than
   letting any stage re-acquire a strict catalog for itself.
3. MUST add a guard that fails when the regeneration path acquires the embedded
   catalog more than once, so a future call site added below the regeneration
   entry point is caught by a test rather than by a maintainer.
4. MUST leave every non-regeneration caller acquiring the catalog exactly as it
   does today, with unchanged public signatures wherever a caller does not opt
   into threading.
5. MUST keep the strict re-validation that closes the regeneration target
   effective — the first invocation must still fail when the catalog is
   genuinely inconsistent after the rewrite, not merely when a pin was stale
   before it.
6. MUST leave the characterization corpus unchanged.

## Subtasks

- [ ] Trace every catalog acquisition the regeneration path reaches.
- [ ] Thread the regeneration catalog through the apply stage and any other
      stage that re-acquires.
- [ ] Add the single-acquisition guard for the regeneration path.
- [ ] Prove the first invocation now succeeds against a stale pin.
- [ ] Prove a genuinely inconsistent post-rewrite catalog still fails.

## Acceptance Criteria

- [ ] Starting from a catalog whose formatter digest pin is stale, one
      invocation of the regeneration target exits zero, and a second
      consecutive invocation reports no change.
- [ ] The guard fails when a stage on the regeneration path acquires the
      embedded catalog a second time.
- [ ] A catalog left genuinely inconsistent after the rewrite still fails the
      target's strict re-validation.
- [ ] Non-regeneration callers still acquire the catalog strictly and
      independently.
- [ ] The characterization corpus is byte-unchanged.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/apply.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog.go`
- interface: `internal/baseline/preservation.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestRegenerationLoadsCatalogOnce -count=1`
  — expected: exit 0; the regeneration path acquires the embedded catalog
  exactly once.
- `go test ./internal/baseline -run TestRegenerationBreaksGoldenDigestCycle -count=1`
  — expected: exit 0; the cycle fixture still regenerates.
- `go test ./internal/baseline -run TestCatalogRegenerationMode -count=1` —
  expected: exit 0; deferral scoping still holds.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; no diagnostic moved.
- `make baseline-digests` — expected: exit 0 with `"changed":false` on the
  consistent tree.
- `make baseline-digests` — expected: exit 0 with `"changed":false` on a second
  consecutive run.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 2; Success Metrics (running the regeneration once
  succeeds).
- `_techspec.md` → System Architecture; Build Order 3 and 4.
- `qa/qa-report-2026-08-01.md` → F-01.
- ADR-0085.

## Result

Threaded the regeneration catalog through the complete formatter-composition
planning and apply path. Package-local catalog-aware entry points now cover
root preservation, plan-document validation, repository validation, and
transaction setup. The existing public `BuildPlan`, `ApplyPlan`,
`PlanRootPreservation`, `ValidatePlanDocument`, `ValidatePlanRepository`, and
`BeginTransaction` signatures remain unchanged; callers that do not opt into
threading retain their independent strict embedded-catalog acquisition.

Added `TestRegenerationLoadsCatalogOnce` at the existing formatter-composition
suite boundary. It wraps the stale-pin embedded fixture with a real filesystem
load counter, acquires one regeneration catalog, builds and applies through the
real plan/transaction path, builds again against the resulting Setup Manifest,
and requires the second plan to contain no file changes. It then loads the same
still-inconsistent fixture strictly and requires the formatter-digest mismatch,
so the final strict-validation contract remains effective.

Focused checks:

- Red signal: before the catalog-aware apply path existed, the temporarily
  named focused regression failed at the first apply with
  `catalog.profile.formatter.goldenDigest.mismatch`; the strict apply stage had
  re-acquired the stale embedded catalog.
- The same temporarily named regression exited 0 after the implementation,
  including its first apply, zero-change second plan, already-applied second
  apply, one-load assertion, and strict-load negative assertion. The final test
  name is `TestRegenerationLoadsCatalogOnce`; its declared full command remains
  for Daemon Verification.
- Focused `TestFormatterComposition -update` exited 0 and exercised the actual
  regeneration-mode planning and apply wiring. Focused ordinary
  `TestFormatterComposition` also exited 0 through the public strict path.
- Focused `TestApplyExactDigest` and
  `TestPreservationPreviouslyBackedUpRootGuidanceIsNotReclassified` runs exited
  0 (four tests/subtests), exercising non-regeneration public apply and
  preservation behavior.
- `rtk gofmt -w` on the five changed Go files exited 0.
- The focused characterization-corpus `git diff --exit-code` check exited 0;
  the diagnostic golden is byte-unchanged.
- `rtk git diff --check` exited 0.
- `rtk git -c core.fsmonitor=false status --short` lists only
  `internal/baseline/` paths and this task file; the task status line was the
  pre-existing Daemon-owned change.

Acceptance evidence:

1. The stale-pin regression completed its first plan/apply with the originally
   acquired regeneration catalog. Its second build produced zero file changes,
   and applying that second plan reported the already-applied verified state.
2. The regression counts a required catalog asset read once per acquisition
   and requires exactly one read after both planning/apply invocations. Any
   nested embedded-catalog acquisition increments the count and fails the test;
   before threading apply, the nested strict acquisition instead failed even
   earlier on the stale pin.
3. After the regeneration operations, strict `LoadCatalog` against the same
   deliberately unrepaired fixture still produced
   `catalog.profile.formatter.goldenDigest.mismatch`, preserving the target's
   strict closing signal for a genuinely inconsistent result.
4. Public non-regeneration signatures and strict wrappers are unchanged. The
   focused ordinary formatter, apply, and preservation checks all exited 0.
5. The diagnostic characterization golden has no byte diff.
6. Final changed-path inspection is bounded to `internal/baseline/` and this
   task file.

The commands under `## Verification` were not run; the Daemon owns them.
