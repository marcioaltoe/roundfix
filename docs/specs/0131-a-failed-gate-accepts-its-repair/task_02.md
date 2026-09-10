---
task: task_02
spec: 0131-a-failed-gate-accepts-its-repair
status: completed
type: backend
complexity: low
---

# Task 02: Accept a failed gate above incomplete dependencies

## Overview

Narrow the gate-staleness rule to a gate whose recorded verdict is a pass, so a
Spec whose QA gate failed accepts corrective Tasks and can be checked and
dispatched. The slice is verifiable alone: the failed case loads, the completed
case still refuses, and no Task status is touched.

## Requirements

1. MUST invalidate a settled gate above incomplete dependencies only when the
   gate's recorded status is completed.
2. MUST load a Task Graph whose gate status is exactly failed even when
   dependencies in the gate's closure are not completed.
3. MUST leave a gate status that is absent, malformed, or any value other than
   completed or failed on its existing path, so acceptance is granted on a
   positive observation and never on missing evidence.
4. MUST keep the stale-gate error's type, message and exported identity
   unchanged, so every caller and fixture that matches on them keeps working.
5. MUST keep requiring that the gate depends on every non-QA leaf, for both
   verdicts.
6. MUST NOT read or write any Task status beyond the existing load-time read;
   the Daemon remains the only writer.
7. MUST update only the characterization row this Task intentionally changes.

## Subtasks

- [ ] Restrict the staleness condition to a completed gate.
- [ ] Prove the failed case loads and the completed case still refuses.
- [ ] Prove leaf coverage still applies to both verdicts.
- [ ] Update the one characterization row this change moves.

## Acceptance Criteria

- [ ] A Spec whose gate status is failed loads with an incomplete dependency in
      the gate's closure, where it refused before this Task.
- [ ] A Spec whose gate status is completed still refuses with the same error
      identity and the same dependency named.
- [ ] A Spec whose gate omits a non-QA leaf still refuses for both verdicts.
- [ ] Loading writes no file: the Spec directory is byte-identical before and
      after.
- [ ] The repository's own active and archived Spec corpus still loads.

## Context

- interface: `internal/spec/spec.go`

## Verification

- `grep -q 'func TestFailedGateLoadsAboveIncompleteDependencies' internal/spec/spec_test.go && go test -count=1 ./internal/spec -run '^TestFailedGateLoadsAboveIncompleteDependencies$'` — the failed gate loads; this fails today.
- `grep -q 'func TestFailedGateLoadsAboveIncompleteDependencies' internal/spec/spec_test.go || exit 1; go test -count=1 ./internal/spec -run '^(TestLoadRejectsAppendedTaskUnderSettledGate|TestGateStalenessCharacterizesEachVerdict)$'` — the completed-gate refusal and the characterization still pass alongside the change.
- `grep -q 'func TestFailedGateLoadsAboveIncompleteDependencies' internal/spec/spec_test.go || exit 1; go test -count=1 ./internal/spec` — the package suite passes with the change present.
- `grep -q 'func TestRepositorySpecCorpusStillLoads' internal/spec/spec_test.go && go test -count=1 ./internal/spec -run '^TestRepositorySpecCorpusStillLoads$'` — every active and archived Spec in this repository still loads, evidence written by earlier work under the contract this Spec changes.

## References

- `_prd.md` → Goals 1-2; Core Features 1-4; Decisions: Declared intentional breaks.
- `_techspec.md` → Implementation Design; Build Order 2; Risks & Considerations.
- ADR-0057, ADR-0091, ADR-0117.

## Result

Narrowed load-time gate staleness to a QA Task recorded `completed`. A QA
Task recorded `failed` now loads above incomplete dependencies, while the
existing leaf-coverage validation still runs before the verdict-specific
staleness check. Task parsing, status ownership, and `StaleGateError` remain
unchanged.

Focused checks:

- Before the loader edit, `rtk go test -count=1 ./internal/spec -run
  'TestFailedGateLoadsAboveIncompleteDependencies/loads_without_writing'`
  reached the real filesystem-backed loader and failed with
  `StaleGateError`, naming QA Task `task_02` and incomplete Task `task_01`.
  The first sandboxed attempt could not read the Go build cache; the unchanged
  rerun with cache access produced this behavior signal.
- After the edit, `rtk go test -count=1 ./internal/spec -run
  '^(TestFailedGateLoadsAboveIncompleteDependencies|TestGateStalenessCharacterizesEachVerdict|TestQAGateLeafCoverageAppliesToBothVerdicts|TestLoadReturnsTypedValidationErrors|TestRepositorySpecCorpusStillLoads|TestLoadInvalidatesSettledQAGateAfterTaskAppend)$'`
  passed 27 focused tests.
- `rtk git diff --check` passed.
- `rtk make verify-incremental` exited 2 at `fmt-check` because unchanged
  files `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. They are
  outside this Task's slice and were not edited.

Acceptance evidence:

- AC1: `TestFailedGateLoadsAboveIncompleteDependencies/loads without writing
  Spec files` loads a failed QA Task above pending `task_01` and observes
  `task_02` as the graph's QA Task. The characterization's failed-verdict row
  is the only existing row whose expected answer changed.
- AC2: the completed-verdict characterization still requires
  `errors.As(..., StaleGateError)` and the exact `task_02`/`task_01`
  identity. `TestLoadInvalidatesSettledQAGateAfterTaskAppend` also preserves
  the exported error identity, invalidated dependency, and existing message
  fragments. The error definition and construction were not edited.
- AC3: `TestQAGateLeafCoverageAppliesToBothVerdicts` runs the same uncovered
  `task_02` fixture with `completed` and `failed` gate statuses; both
  return `QAGateError` whose reason names the leaf-coverage rule and
  `task_02`.
- AC4: the failed-gate test recursively snapshots every fixture directory
  entry and file byte before and after `Load`; the maps remain equal.
- AC5: `TestRepositorySpecCorpusStillLoads` loads every materialized active
  Task Graph and copies each archived manifest plus its real Task files into an
  active temporary Spec before loading it. Active PRD-only planning folders
  without `_tasks.md` are outside this loader corpus.

Additional requirement evidence: `TestLoadReturnsTypedValidationErrors`
continues to reject unparseable frontmatter and unsupported status values, so
missing or unknown gate evidence does not gain acceptance. The production diff
only narrows the verdict guard; it adds no Task status read or write.

The commands under `## Verification` were not run. The Agent preserved the
Daemon-authored Task status and changed no Task Graph, sibling Task, commit,
push, or pull request.
