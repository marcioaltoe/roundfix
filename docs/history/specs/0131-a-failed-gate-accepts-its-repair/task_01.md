---
task: task_01
spec: 0131-a-failed-gate-accepts-its-repair
status: completed
type: test
complexity: low
---

# Task 01: Characterize the loader's answer for each gate verdict

## Overview

Record what the Task Graph loader answers today for a settled QA gate standing
above an incomplete dependency, separately for a completed gate and a failed
one. The slice is verifiable alone: it states the present contract, including
the answer the next Task changes.

## Requirements

1. MUST record that a completed gate above one incomplete dependency refuses,
   naming the error identity and the dependency the error reports.
2. MUST record that a failed gate above one incomplete dependency refuses today
   with the same error identity, which is the answer the next Task changes.
3. MUST record that a gate covering every dependency loads for both verdicts,
   so the characterization separates staleness from coverage.
4. MUST NOT change loader behavior, and MUST NOT weaken or delete an existing
   assertion.

## Subtasks

- [ ] Build Spec fixtures for a completed and a failed gate above an incomplete dependency.
- [ ] Record the refusal each one produces today, with its error identity.
- [ ] Record the loading case for both verdicts.

## Acceptance Criteria

- [ ] A case asserts the completed gate's refusal and the dependency it names.
- [ ] A case asserts the failed gate's refusal today, marked as the behavior the
      next Task changes.
- [ ] A case asserts that both verdicts load when every dependency is completed.
- [ ] Every assertion that existed in the file before this Task still runs.

## Context

- interface: `internal/spec/spec.go`

## Verification

- `grep -q 'func TestGateStalenessCharacterizesEachVerdict' internal/spec/spec_test.go && go test -count=1 ./internal/spec -run '^TestGateStalenessCharacterizesEachVerdict$'` — the per-verdict characterization executes; its absence cannot pass.
- `grep -q 'func TestGateStalenessCharacterizesEachVerdict' internal/spec/spec_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestLoadRejectsAppendedTaskUnderSettledGate'` — the pre-existing stale-gate assertion still runs and passes beside the new characterization.

## References

- `_prd.md` → Success Metrics; Decisions: Regression locks.
- `_techspec.md` → Testing Approach observations 1-3; Build Order 1.

## Result

Added one table-driven loader characterization over real temporary Spec files.
It varies the QA gate verdict and its sole dependency's status without changing
loader source or any pre-existing assertion.

Acceptance evidence:

- AC1: `completed gate rejects an incomplete dependency` requires
  `errors.As` to expose `StaleGateError` and requires its fields to name QA Task
  `task_02` and incomplete dependency `task_01`.
- AC2: `failed gate currently rejects an incomplete dependency` records the
  same typed error and dependency fields. A source comment marks this as the
  case Task 02 changes to a successful load.
- AC3: the `completed` and `failed` gate cases with completed `task_01` both
  require `Load` to return without error and expose `task_02` as the QA Task.
- AC4: no existing line in `internal/spec/spec_test.go` was changed. Focused
  check `rtk go test -count=1 ./internal/spec` passed with 396 tests, exercising
  the complete package suite alongside the new cases. The first sandboxed run
  could not access the Go build cache; the unchanged rerun with cache access
  produced the recorded result.
- Incremental check `rtk make verify-incremental` stopped at `fmt-check`
  because unchanged files `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. They are outside
  this Task's slice and were not changed.

Follow-up: the second declared Verification command selects
`TestLoadRejectsAppendedTaskUnderSettledGate`, while the pre-existing regression
is named `TestLoadInvalidatesSettledQAGateAfterTaskAppend`. This Task leaves the
authored Verification and existing test name unchanged.

Limits: the commands under `## Verification` were not run because the Daemon
owns them. No Task status, Task Graph, sibling Task, loader behavior, commit,
push, or pull request was changed by this implementation turn.
