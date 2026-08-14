---
task: task_02
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 02: Give the check a tree of its own to run in

## Overview

Executing an authored Verification command means executing text a Supervisor
wrote, which can do anything a shell can. This slice gives the authoring check a
disposable tree at `HEAD` to run in, so a command that writes cannot reach the
operator's working tree, and removes it afterwards whether the check passed or
failed.

## Requirements

1. MUST materialize the repository at `HEAD` in a temporary location distinct
   from the operator's working tree.
2. MUST remove that location when the caller finishes, including when the caller
   fails or panics.
3. MUST leave the operator's working tree unchanged by anything that runs inside
   the disposable tree, including untracked files and the index.
4. MUST create the tree with the same Git hygiene the repository's isolated test
   repositories already use, since worktree creation under load has a recorded
   failure in this repository.
5. MUST fail with a named reason rather than a raw Git error when the tree cannot
   be created.

## Subtasks

- [ ] Add the disposable checkout with its cleanup.
- [ ] Apply the repository's isolation settings on creation.
- [ ] Cover isolation, cleanup on failure, and a creation refusal.

## Acceptance Criteria

- [ ] The checkout path differs from the repository root and holds the tree at
      `HEAD`.
- [ ] Writing a file inside the checkout leaves the repository's working tree and
      index unchanged.
- [ ] The checkout is gone after the caller returns, including on failure.
- [ ] A creation failure reports a named reason rather than a raw Git message.

## Rehearsal Cases

- Case: a caller that writes a file into the checkout; Observation: the
  repository's `git status` is identical before and after.
- Case: a caller that returns an error; Observation: the checkout is still
  removed.
- Case: a creation refused by Git; Observation: the error names the checkout
  rather than passing the Git text through.

## Verification

- `grep -rq 'func DisposableCheckout' internal/` — expected: exits 0, proving the helper exists. Fails today.
- `go test -count=1 ./internal/... -run 'TestDisposableCheckout' -v > /tmp/0095-t02.log 2>&1; s=$?; grep -q '^--- PASS: TestDisposableCheckout' /tmp/0095-t02.log || { cat /tmp/0095-t02.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0095-t02.log` — expected: exits 0, refusing a vacuous run.
- `d=$(grep -rl 'func DisposableCheckout' internal/); test -n "$d" || { echo 'FAIL: no file declares DisposableCheckout'; exit 1; }; grep -q 'core.fsmonitor=false' $d || { echo "FAIL: $d omits the fsmonitor isolation setting"; exit 1; }` — expected: exits 0, proving the checkout carries the isolation setting this repository's own worktree findings ask for. The declaration is proven present before the setting is grepped, so an empty file list fails rather than reading standard input.

## Context

- interface: `internal/gittest/gittest.go`

## References

`_techspec.md` → Build Order 2; Interfaces: `DisposableCheckout`; Integration
Points; Risks: executing authored commands is executing untrusted text.
`_prd.md` → Core Feature 1.
