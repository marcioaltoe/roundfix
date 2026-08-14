---
task: task_07
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 07: Prove the tree exited, not only its owner

## Overview

Stopping a Run proves the registered owner exited. The processes that owner
spawned are what survived for three days. This slice makes Force Stop prove the
tree, using the process-table reader the residue inventory built.

## Requirements

1. MUST prove the exit of the processes a Run spawned, not only of its registered
   owner.
2. MUST fail rather than report success when a descendant survives, naming the
   surviving process.
3. MUST leave the command's output and exit status unchanged when the tree does
   exit.
4. MUST NOT terminate a process outside the Run's own spawn lineage.

## Subtasks

- [ ] Prove the tree's exit rather than the owner's.
- [ ] Name a surviving descendant in the failure.
- [ ] Cover a surviving descendant and a clean tree.

## Acceptance Criteria

- [ ] A Run whose descendant survives makes the stop fail, naming that process.
- [ ] A Run whose tree exits reports success exactly as before.
- [ ] No process outside the Run's lineage is touched.

## Verification

- `go test -count=1 ./internal/cli -run 'TestForceStopProvesTheTree' -v > /tmp/0103-t07.log 2>&1; s=$?; grep -q '^--- PASS: TestForceStopProvesTheTree' /tmp/0103-t07.log || { cat /tmp/0103-t07.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t07.log` — expected: exits 0, refusing a vacuous run.
- `grep -c '^--- PASS' /tmp/0103-t07.log > /tmp/0103-t07-n.txt; test "$(cat /tmp/0103-t07-n.txt)" -ge 2 || { echo "expected the surviving-descendant and clean-tree cases, got $(cat /tmp/0103-t07-n.txt)"; cat /tmp/0103-t07.log; exit 1; }` — expected: exits 0, proving both directions are exercised.

## Context

- interface: `internal/cli/cli.go`

## References

`_techspec.md` → Build Order 7; API Contracts. `_prd.md` → Core Feature 6;
Goal 2; Non-Goals, processes started outside Roundfix. ADR-0127.
