---
task: task_11
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
type: test
complexity: low
---

# Task 11: Make the assembled tree pass its own gate

## Overview

`TestRunCommandHelp` pins the `reconcile` synopsis by exact string, and Task 10
changed that synopsis to name the two acts Tasks 05 and 06 added. The contract
lives in `internal/cli/cli_test.go`, outside Task 10's bounded scope, so Task 10
could not update it and the QA gate returned `fail` a third time.

This is the fourth Task this Spec has needed for the same reason: a public
change whose surrounding contract sat outside every boundary. The first three
were minted one instance at a time, and each rerun of the gate found the next
one. This Task closes the class instead of the instance — its acceptance is the
repository gate itself, so it cannot settle while any contract this Spec broke
is still pinned to the pre-Spec text.

## Requirements

1. MUST update every contract assertion that pins text this Spec changed, so
   each states what the surface now says.
2. MUST keep every such assertion exact. Do not relax a pinned string to a
   substring or a regex to make it pass; the point of pinning is that an
   unannounced change fails.
3. MUST NOT change production code, help copy, or behaviour. If an assertion
   and the surface genuinely disagree about what is correct, stop and report
   which one is wrong rather than editing the test to match.

## Subtasks

- [ ] Update the pinned `reconcile` synopsis.
- [ ] Confirm no other contract pins text this Spec changed.

## Acceptance Criteria

- [ ] `TestRunCommandHelp` passes against the help copy Task 10 wrote.
- [ ] Every updated assertion still compares exact text.
- [ ] `git diff --name-only` lists only test files and this Task file.

## Bounded scope

This Task may create or modify only:

- `internal/cli/cli_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_11.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunCommandHelp$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunCommandHelp'` — expected: exits 0. The assertion pins the pre-Spec synopsis and fails against the unchanged tree.

Asserting that the help still contains `--carry-forward` is deliberately absent:
Task 10 already delivered that copy, so the check passes before this Task does
anything. The pre-work probe refused it on 2026-08-11. A command that guards a
predecessor's work cannot prove this Task's own.
`make verify` is deliberately not declared here. The Daemon recognises it as the
repository gate and refuses to dispatch a Task while it is red — "repository not
green on entry" — so that a pre-existing failure is never attributed to the Task
being started. That refusal fired on 2026-08-11 and is correct; it simply means
the gate cannot be a Task's own proof at the moment the gate is failing. The
Run-level gate still proves the assembled tree after this Task settles.

## References

- `_prd.md` → Goals 3 and 4.
- `task_10.md` → the help copy this contract pins.

## Result

Updated the `reconcile` row in `TestRunCommandHelp` to pin the complete synopsis
Task 10 added, including the mutually exclusive `--discard-superseded` and
`--carry-forward` acts. The assertion still uses the existing exact string
literal; no matcher, production code, help copy, or behavior changed.

Acceptance evidence:

- Help contract: before the edit,
  `rtk proxy rg -n -C 20 'TestRunCommandHelp|reconcile \\[--apply\\]|discard-superseded|carry-forward' internal/cli/cli_test.go internal/cli/cli.go`
  showed the pre-Spec synopsis in `internal/cli/cli_test.go:432` and the Task 10
  synopsis in `internal/cli/cli.go:51` and `internal/cli/cli.go:5296`. After the
  edit, `rtk proxy rg -n -F 'roundfix reconcile [run-id] [--apply | --discard-superseded | --carry-forward] [--format <text|json>]' internal/cli/cli.go internal/cli/cli_test.go`
  found the same complete literal at all three locations.
- Focused behavior: `rtk proxy env GOCACHE=/private/tmp/roundfix-task11-go-cache rtk go test ./internal/cli -run '^TestRunCommandHelp$/^reconcile$' -count=1`
  exited 0 and reported `2 passed in 1 packages`.
- Exactness and sweep: a repository-wide `*_test.go` search for the three help
  strings Task 10 removed — the old synopsis and the two old aligned option
  rows — exited 1 with no matches. Inspection of `rtk git diff --
  internal/cli/cli_test.go` showed only the exact expected synopsis literal
  changed; the `strings.Contains` assertion and surrounding test are unchanged.
- Bounded scope: `rtk git diff --name-only` listed only
  `internal/cli/cli_test.go` and this Task file.

Follow-up: `docs/user-guide/commands.md:616` still contains the pre-Spec
synopsis. It is not a contract assertion and is outside this Task's bounded
scope, so this Task did not edit it.

The command under `## Verification` and the repository gate were not run; they
remain owned by the Daemon.
