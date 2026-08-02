---
task: task_11
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: medium
---

# Task 11: Report the result as a status matrix

## Overview

A successful apply currently reads as an update being complete, when what was
proven is that postimages were written. Retention, alignment, repository
Verification, and idempotence are separate facts and some of them may not have
run at all. This Task reports five axes separately and reserves completion
language for the case that earns it.

## Requirements

1. MUST report the final result as five separate axes: approved postimages,
   semantic retention, profile alignment, repository Verification, and
   idempotence.
2. MUST report each axis as verified or not run, so an axis that never
   executed is never read as passing.
3. MUST use completion language only when semantic retention is verified and
   the idempotence check passed.
4. MUST derive each axis from the evidence the run actually produced, not from
   the absence of an error.
5. MUST carry the same five axes in machine output, additively.
6. MUST leave the transaction, apply, and digest confirmation behavior
   unchanged.

## Subtasks

- [ ] Report the five axes separately.
- [ ] Distinguish verified from not run on each.
- [ ] Gate completion language on retention and idempotence.
- [ ] Add the axes to machine output additively.

## Acceptance Criteria

- [ ] The final result shows all five axes, each verified or not run.
- [ ] A run where the idempotence check did not execute shows it as not run,
      not as verified.
- [ ] Completion language appears only when retention is verified and
      idempotence passed.
- [ ] A run with verified postimages but unverified retention does not read as
      complete.
- [ ] Machine output carries the same five axes, with every prior field
      unchanged.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run '^TestResultStatusMatrix$' -count=1 -v | grep -q -- "--- PASS: TestResultStatusMatrix"` — expected:
  exit 0; five axes, each verified or not run.
- `go test ./internal/baseline -run '^TestCompletionLanguageRequiresRetention$' -count=1 -v | grep -q -- "--- PASS: TestCompletionLanguageRequiresRetention"`
  — expected: exit 0; verified postimages alone never read as complete.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 7; Core Features 10; User Experience.
- `_techspec.md` → API Contracts; Build Order 11.

## Result

### Implementation

- Successful apply results now carry an additive `statusMatrix` with approved
  postimages, semantic retention, profile alignment, repository Verification,
  and idempotence. Every axis is exactly `verified` or `not run`.
- The matrix derives approved postimages from the transaction's exact
  postimage evidence, semantic retention from a nonempty fully accounted
  `ClauseDelta`, profile alignment from the digest-bound Plan and Setup
  Manifest identity, and idempotence from the exact already-applied preimage
  branch. Repository Verification remains `not run` because apply does not
  execute the recommended command.
- Human result text renders the same five axes. It uses `Baseline update
  complete` only when semantic retention and idempotence are both `verified`;
  first apply and retention-free reapply results do not use completion
  language.
- The existing result fields, apply state, transaction, digest confirmation,
  verified-postimage list, warnings, and recommendations remain in place. The
  strict result codec accepts the new field and rejects any axis value outside
  `verified` and `not run`.

### Focused checks

- Pre-change signal: `rtk go test ./internal/baseline -run
  'TestResultStatusMatrix/first_apply' -count=1` reached the expected compile
  failure because the status-matrix types did not exist.
- `rtk gofmt -w internal/baseline/plan.go internal/baseline/apply.go
  internal/baseline/plan_json.go internal/baseline/apply_test.go` — exit 0.
- `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task11-gocache rtk go test
  ./internal/baseline -run
  "StatusMatrix|CompletionLanguage|ApplyExactDigest|EmptyReapply|PlanDocumentStrictCodecs"'`
  — 8 tests passed. The task-scoped cache was required because the sandbox
  denied the default macOS Go cache.
- `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task11-gocache rtk go test
  ./internal/baseline -run
  "InstructionHierarchyPreservesPlanAndResultSchemas|BaselineVerification"'`
  — 2 tests passed; existing result-schema fields and the no-execution
  repository Verification contract remain covered.
- The commands declared under `## Verification` were not run; the Daemon owns
  those checks.

### Acceptance evidence

- Five axes: `TestResultStatusMatrix` checks all five human labels and the five
  machine fields with their exact statuses.
- Not-run distinction: the first apply reports idempotence as `not run`, while
  the exact already-applied branch reports it as `verified`; repository
  Verification remains `not run` in both.
- Completion gate: `TestCompletionLanguageRequiresRetention` checks that
  retention-free reapply and retention-only first apply do not read as
  complete, while verified retention plus verified idempotence does.
- Unverified retention: the first greenfield apply has verified postimages and
  `not run` semantic retention without completion language.
- Additive machine output: `TestResultStatusMatrix` marshals the result, checks
  the prior schema, operation, state, message, digest, postimage, warning, and
  recommendation fields, and verifies the additive `statusMatrix`; it also
  checks rejection of an invalid axis status.
- Scope: the final changed-path audit reports only `internal/baseline/` paths
  and this Task file.
