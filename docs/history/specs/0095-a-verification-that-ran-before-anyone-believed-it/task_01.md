---
task: task_01
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 01: Extract the prober the Daemon already runs

## Overview

The Daemon classifies every authored Verification command against the unchanged
tree before opening an Agent Session. That loop is wrapped in Run bookkeeping the
authoring caller has no use for. This slice separates the two so both callers
share one classification, which is what stops a checker from approving what the
probe later refuses.

## Requirements

1. MUST expose a prober that takes a working directory, an ordered command list,
   a verifier, and a per-command output destination, and returns one verdict per
   command in the order given.
2. MUST classify a command that passes against the unchanged tree as vacuous, one
   that fails as honest, and one that could not run as unknown with its cause.
3. MUST NOT read or write run state, acquire verification capacity, or resolve
   artifact paths; those stay with the Daemon.
4. MUST leave the Daemon's existing pre-work probe behaviour identical, including
   its run-state transition, capacity acquisition and probe log paths.
5. MUST leave the Daemon's existing probe tests passing without editing them; a
   test that has to change is a signal the extraction moved behaviour rather than
   code.

## Subtasks

- [x] Add the shared prober with its verdict type.
- [x] Make the Daemon's pre-work probe call it.
- [x] Cover ordering, the three verdicts, and an empty command list.

## Acceptance Criteria

- [x] A command that exits zero against the unchanged tree returns vacuous.
- [x] A command that exits non-zero returns honest, not unknown.
- [x] A command that could not run returns unknown carrying its cause.
- [x] Verdicts come back in the order the commands were given.
- [x] The Daemon's existing probe tests pass unedited.

## Verification

- `grep -rq 'func ProbeCommands' internal/ && grep -rq 'CommandVerdict' internal/` — expected: exits 0, proving the shared prober and its verdict type exist. Fails today, where neither is declared.
- `go test -count=1 ./internal/daemon -run 'TestProbeCommands' -v > /tmp/0095-t01.log 2>&1; s=$?; grep -q '^--- PASS: TestProbeCommands' /tmp/0095-t01.log || { cat /tmp/0095-t01.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing prober test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0095-t01.log` — expected: exits 0, refusing a vacuous run.
- `git diff --name-only HEAD -- internal/daemon/task_engine_test.go > /tmp/0095-t01-edit.txt; test ! -s /tmp/0095-t01-edit.txt && grep -rq 'ProbeCommands' internal/daemon/task_engine.go || { echo 'existing probe tests were edited, or the Daemon does not call the shared prober:'; cat /tmp/0095-t01-edit.txt; exit 1; }` — expected: exits 0, proving the Daemon delegates to the prober and its existing tests were left alone. Both clauses are one command because either alone passes on a tree where the other did not happen.

## Context

- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 1; Interfaces: `ProbeCommands`, `CommandVerdict`;
Testing Approach, the shared prober. `_prd.md` → Core Feature 1. ADR-0124,
ADR-0111, ADR-0014.

## Result

### Implementation

- Added `ProbeCommands` and `CommandVerdict`. The prober invokes the supplied
  verifier once per command, preserves command order, delegates each diagnostic
  path to `outputFor`, and owns no Run state, capacity, or artifact-path policy.
- `verifyTaskPreWork` now keeps its existing Run-state transition and
  Verification Capacity acquisition around `ProbeCommands`, and still constructs
  the same numbered probe log paths before publishing or settling the result.
- Added the source-matched prober suite without changing
  `internal/daemon/task_engine_test.go`.

### Acceptance evidence

| Criterion | Focused evidence |
| --- | --- |
| Exit zero is vacuous | `TestProbeCommands/classifies_ordered_command_outcomes` asserts the zero-error verifier result sets only `Vacuous`. |
| Observed non-zero is honest | The same subtest returns `VerificationCommandError` and asserts `Vacuous`, `Unknown`, and `Cause` remain unset. |
| A command that cannot run is unknown with its cause | The same subtest returns `VerificationUnknownError`, asserts `Unknown`, verifies `errors.Is` reaches the runner cause, and verifies the command and diagnostic path are retained. |
| Verdict order matches command order | The same subtest compares every verdict and verifier request with the ordered input and checks the zero-based `outputFor` destinations. |
| Existing Daemon probe behavior is unchanged | `GOCACHE=/tmp/roundfix-0095-t01-gocache rtk go test ./internal/daemon -run '^(TestPreWorkProbe|TestVerificationProbeCharacterization)' -count=1` passed 9 tests. `rtk git status --short -- internal/daemon/task_engine_test.go` produced no changed path. |

### Focused checks

- Red signal: `GOCACHE=/tmp/roundfix-0095-t01-gocache rtk go test ./internal/daemon -run '^TestProbeCommands$' -count=1` failed to compile before the implementation because `ProbeCommands` was undefined.
- After implementation, the same focused command passed 3 test executions,
  including the empty-command-list case.
- The final combined prober check, `GOCACHE=/tmp/roundfix-0095-t01-gocache rtk go test ./internal/daemon -run '^(TestProbeCommands|TestPreWorkProbe|TestVerificationProbeCharacterization)' -count=1`, passed 12 test executions.
- `rtk git diff --check` passed.
- The first `GOCACHE=/tmp/roundfix-0095-t01-gocache rtk make verify-incremental`
  attempt hit the unrelated 120-second context deadline in
  `internal/store/TestJournalMeasurementHarness` while another suite process was
  active. The test passed in isolation, and the fresh incremental rerun then
  passed all tests, skill checks, and the build; `internal/daemon` passed in
  3.782 seconds on that rerun.
