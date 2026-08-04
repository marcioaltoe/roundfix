---
task: task_11
spec: 0064-spec-artifact-consistency-gate
status: completed
type: infra
complexity: low
---

# Task 11: Run the budget proof as its own gate step

## Overview

task_10 gives the sweep budget a test that measures under conditions it
declares. This slice makes the required local gate execute that proof on every
run, serially, so the budget is proven rather than merely provable.

It is a separate Task because it changes `Makefile`, and an authorized tooling
Task may mutate only its bounded files. Folding it into task_10 would fail the
tooling-authority audit that the QA gate performs against each Task commit.

## Requirements

1. MUST add a `spec-budget` target that executes the sweep-budget proof under
   the conditions task_10 declared, serially rather than inside the contended
   package sweep.
2. MUST include that step in the `verify` gate so every gate run proves the
   budget, and MUST place it where a failure stops the gate.
3. MUST NOT let a pipe hide the step's exit status — no pager, no `tail`, no
   pipeline whose last command masks a failure.
4. MUST declare the target in `.PHONY` and give it a `##` help description
   matching the file's existing style.
5. MUST change exactly one file: `Makefile`. This is the complete bounded file
   list from the Tooling authority row in `_prd.md` and `_techspec.md`, granted
   by the 2026-08-02 authorization record that names this Spec. Any other path
   is out of scope — stop rather than widen it.

## Subtasks

- [ ] Add the budget-proof target with its help description.
- [ ] Include it in `verify` where a failure stops the gate.
- [ ] Add it to `.PHONY`.

## Acceptance Criteria

- [ ] The gate step exits 0 on the current tree and measures the sweep rather
      than skipping it.
- [ ] The step fails the gate when the budget is exceeded, proven by a
      temporary limit reduction reverted within the same check.
- [ ] `verify` includes the step, and a failing step stops the gate rather
      than being reported and passed over.
- [ ] The target appears in `.PHONY` and in `make help`.
- [ ] `make verify` exits 0 on the current tree.
- [ ] `Makefile` is the only file this Task changed.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `Makefile`

## Verification

- `make verify` — expected: exit 0; the full gate including the budget proof.
- `grep -q "^\.PHONY:.*spec-budget" Makefile` — expected: exit 0.
- `grep -q "^verify:.*spec-budget" Makefile` — expected: exit 0; the gate
  includes the step.
- `make help | grep -q "spec-budget"` — expected: exit 0.
- `git diff --name-only HEAD | grep -v "^Makefile$" | grep -v "^docs/specs/0064-spec-artifact-consistency-gate/task_11.md$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded file and this Task file changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Feature 1.
- `_techspec.md` → Testing Approach; Integration Points.
- `qa/qa-report-2026-08-03.md` → F-001 required repair.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.

## Result

### Implementation

- Added the `spec-budget` target with a `##` help description. Its only recipe
  command runs the dedicated `^TestCheckCorpusBudget$` selector with a fresh
  test result and `-parallel=1`; it has no pipeline or output filter that can
  replace the test process's exit status.
- Added `spec-budget` to `.PHONY` and placed it immediately after `test` in the
  `verify` prerequisites, before the remaining skill, build, and Spec checks.

### Focused checks

- `rtk make spec-budget` — exit 0 after the Makefile edit; RTK reported one
  passed test in `internal/speccheck`, with no skipped test.
- Isolated negative probe under
  `/private/tmp/roundfix-task11-budget.vzk8et/repo`: copied the current
  Makefile into an archive of `HEAD`, changed the copied `corpusBudget` from
  `time.Second` to `time.Nanosecond`, and ran `rtk make spec-budget`. Make
  exited 2; `TestCheckCorpusBudget` reported a `299.879208ms` sweep against
  `1ns` and failed. Restored `time.Second`, reran the target with exit 0, and
  used `rtk cmp` to prove the restored copied test matched the worktree test
  byte-for-byte. The temporary directory was then removed.
- `rtk make -n verify` — exit 0; the dry run printed the ordinary parallel
  test command, then the dedicated serial budget command, then the remaining
  skill, build, and Spec-check commands. This inspected the gate order without
  running the declared Verification.
- `rtk rg -n "^\\.PHONY:.*spec-budget|^verify:.*spec-budget|^spec-budget:.*##|TestCheckCorpusBudget" Makefile`
  — exit 0; matched the phony declaration, `verify` prerequisite, help-style
  target declaration, and anchored test selector.
- `rtk git -c core.fsmonitor=false diff --name-only` — exit 0; listed only
  `Makefile` and this assigned Task file. `rtk git diff --check` also exited 0.

### Acceptance-criterion evidence

1. The focused target exited 0 and reported one passed test rather than a
   skip, proving that it selected and executed the sweep-budget proof.
2. The isolated `1ns` probe made the same target exit 2 with the exceeded
   budget diagnostic. Restoring `time.Second` made it exit 0, and `cmp`
   confirmed that the temporary test mutation was fully reverted.
3. `verify` names `spec-budget` immediately after `test`; the dry run confirms
   that order. The target contains one unpiped command, and the negative probe
   confirms its non-zero test status propagates through Make.
4. The focused source check found `spec-budget` in `.PHONY` and found its `##`
   help description. The Daemon-owned `make help` Verification remains the
   authoritative rendered-help check.
5. `make verify` is a declared Verification command and was not run in this
   Daemon-assigned turn.
6. The changed-path postflight listed only the authorized `Makefile`, plus
   this assigned Task file for required evidence.
