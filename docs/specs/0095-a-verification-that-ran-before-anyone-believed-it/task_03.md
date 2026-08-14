---
task: task_03
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 03: Execute every authored command before any Run

## Overview

The slice that answers the Spec's question. `spec check --run-verification` runs
each authored Verification command in the disposable tree and reports what it
did, so a command that cannot run, cannot fail, or fails on correct work is
visible before a Run starts rather than after one is spent.

## Requirements

1. MUST execute every authored Verification command of every requested Spec in a
   disposable tree and report one line per command with its verdict.
2. MUST exit non-zero when any command is vacuous, matching what the Daemon does
   with the same command.
3. MUST stay opt-in, so `spec check` without the flag executes nothing and keeps
   its current speed.
4. MUST report that commands were not executed when the flag is absent, so a
   clean check is not read as covering what it never ran.
5. MUST report a command it could not run as unknown, distinctly from one that
   failed.
6. MUST state that the tree it ran against is `HEAD`, since a later Task's real
   pre-work tree includes its predecessors' work.

## Subtasks

- [ ] Add the flag and wire the checkout to the prober.
- [ ] Report one line per command with its verdict.
- [ ] Report the unexecuted state when the flag is absent.
- [ ] Cover a vacuous command, an honest one, and an unknown one end to end.

## Acceptance Criteria

- [ ] A Spec whose Task carries a vacuous command exits non-zero and names that
      command.
- [ ] A Spec whose commands all fail against `HEAD` exits zero.
- [ ] Without the flag, nothing executes and the output says so.
- [ ] A command that cannot run reports unknown rather than vacuous or honest.
- [ ] The report names `HEAD` as the tree it ran against.

## Rehearsal Cases

- Case: a Task with one vacuous and one honest command; Observation: the vacuous
  one is named, the honest one is reported as failing correctly, and the exit is
  non-zero.
- Case: the same Spec without the flag; Observation: nothing executes and the
  output states that commands were not run.
- Case: a command naming a binary the tree does not have; Observation: unknown,
  with its cause.

## Verification

- `go test -count=1 ./internal/cli -run 'TestSpecCheckRunVerification' -v > /tmp/0095-t03.log 2>&1; s=$?; grep -q '^--- PASS: TestSpecCheckRunVerification' /tmp/0095-t03.log || { cat /tmp/0095-t03.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0095-t03.log` — expected: exits 0, refusing a vacuous run.
- `go build -buildvcs=false -o /tmp/0095-t03-roundfix ./cmd/roundfix && /tmp/0095-t03-roundfix spec check --help 2>&1 | grep -q 'run-verification'` — expected: exits 0, proving the flag reaches the built command's own help rather than only its tests.

## Context

- interface: `internal/cli/spec_check.go`

## References

`_techspec.md` → Build Order 3; API Contracts; Risks: the probe runs against
`HEAD`, and a partial detector reads as a gate. `_prd.md` → Core Feature 1;
Goal 1; User Story 1. ADR-0124, ADR-0111.
