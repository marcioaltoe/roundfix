---
task: task_01
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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

- [ ] Add the shared prober with its verdict type.
- [ ] Make the Daemon's pre-work probe call it.
- [ ] Cover ordering, the three verdicts, and an empty command list.

## Acceptance Criteria

- [ ] A command that exits zero against the unchanged tree returns vacuous.
- [ ] A command that exits non-zero returns honest, not unknown.
- [ ] A command that could not run returns unknown carrying its cause.
- [ ] Verdicts come back in the order the commands were given.
- [ ] The Daemon's existing probe tests pass unedited.

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
