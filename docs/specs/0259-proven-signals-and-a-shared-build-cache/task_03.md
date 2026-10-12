---
task: task_03
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
type: infra
complexity: medium
---

# Task 03: make shares one Go build cache that the sandbox can write

## Overview

The Makefile defaults `GOCACHE` to the worktree's `.gocache`, so every Run,
Task and item worktree builds its own cache of up to 8 GB. A probe under the
Codex workspace sandbox showed `$TMPDIR` and `/tmp` writable and the Go
default cache denied (`_techspec.md` → Current behavior). This Task makes
`make` use one cache under the temporary directory and keeps the worktree
cache only as the fallback. It is verifiable on its own: the Makefile's `test`
target, run with a stub Go, exports the expected path in each case.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   [each worktree keeps its own multi-gigabyte Go build cache](references/2026-10-09-each-worktree-keeps-its-own-go-build-cache.md),
   by sharing one cache between worktrees.
2. MUST change only the `GOCACHE` default of the Makefile, following
   `_techspec.md` → API Contract 7: computed once, `roundfix-go-build-<uid>`
   under `$TMPDIR` without its trailing slash or under `/tmp`, used when the
   directory can be created and a file created and removed in it, otherwise
   `$(CURDIR)/.gocache`; an explicit `GOCACHE` from the environment or the
   command line is exported unchanged. The bounded file of the cycle grant for
   this Task is `Makefile`, and `docs/agents/specific-repository.md` for the
   guide line; no other Makefile target, no CI workflow and no ignore file
   changes.
3. MUST add `TestMakeSharesOneGoBuildCache` with subtests `shared`,
   `fallback` and `explicit`. Each runs the repository Makefile's `test`
   target with `GO` set to a stub that records the `GOCACHE` it receives, in a
   temporary directory, with `GOCACHE` removed from the child environment
   unless the case sets it: `shared` sets `TMPDIR` to a writable temporary
   directory and expects `<TMPDIR>/roundfix-go-build-<uid>`, which exists
   afterwards; `fallback` sets `TMPDIR` to a directory without write
   permission and expects `<directory of the run>/.gocache`; `explicit` sets
   `GOCACHE` and expects it unchanged. The test lives beside the existing
   Makefile wiring tests, in a package that already installs the suite guard.
4. MUST add one line to `docs/agents/specific-repository.md`, outside every
   setup marker, naming the shared `roundfix-go-build-<uid>` cache, its
   fallback, that an explicit `GOCACHE` wins, and that a sandboxed Agent sets
   nothing.
5. MUST NOT remove or rewrite the `.gocache` entry of `.gitignore`, which the
   fallback still needs.

## Subtasks

- [ ] Change the `GOCACHE` default with its writability fallback.
- [ ] Add the Makefile test with its three cases.
- [ ] Add the repository guide line.

## Acceptance Criteria

- [ ] With `GOCACHE` unset and a writable `TMPDIR`, `make test` passes the shared cache to Go.
- [ ] With an unwritable `TMPDIR`, it passes the worktree's `.gocache`.
- [ ] An explicit `GOCACHE` reaches Go unchanged.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `internal/verifyselect/full_contract_run_test.go`
- interface: `Makefile`
- creates: `internal/verifyselect/build_cache_makefile_test.go`
- interface: `docs/agents/specific-repository.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'roundfix-go-build-' || { printf 'the repository guide does not name the shared cache\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestMakeSharesOneGoBuildCache$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMakeSharesOneGoBuildCache TestMakeSharesOneGoBuildCache/shared TestMakeSharesOneGoBuildCache/fallback TestMakeSharesOneGoBuildCache/explicit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/verifyselect || exit 1` — expected: exit 0. Before this Task the guide line and the test are absent, so the command fails; afterwards the three cases and the existing Makefile wiring tests pass.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 6; Success Metric 4; Acceptance evidence
- `_techspec.md` → Current behavior; API Contract 7; Invariant 7; Testing Approach; Build Order 3
- ADR-0266


## Result

Implemented the assigned cache slice. The Makefile computes its default once,
using `roundfix-go-build-<uid>` under `$TMPDIR` without the trailing slash or
under `/tmp` when `TMPDIR` is empty. It selects that directory only after
`mkdir -p`, a unique file creation and file removal succeed, otherwise exports
`$(CURDIR)/.gocache`. Environment and command-line `GOCACHE` values bypass the
probe and are exported unchanged. The repository guide names the default,
fallback and override, and says a sandboxed Agent sets nothing.

`TestMakeSharesOneGoBuildCache` exercises the repository Makefile's `test`
target with a compiled Go stub from `testfixture.FixtureBinary`, temporary
working directories and a child environment cleared of inherited `GOCACHE`
and Make overrides. The existing package suite guard covers it.

### Focused evidence

- Before the Makefile edit,
  `rtk proxy go test -race -count=1 -v -run '^TestMakeSharesOneGoBuildCache/(shared|fallback|explicit)$' ./internal/verifyselect`
  exited 1: both shared-worktree cases received their own `.gocache` instead
  of the expected shared cache. The fallback and explicit cases passed.
- After the final implementation and test edits,
  `rtk proxy go test -race -count=1 -v -run '^(TestMakeSharesOneGoBuildCache|TestVerifyContractsRunsEveryContract|TestVerifyChangedRunsTheSelectedContracts|TestMainAndReleaseRunEveryContract)$' ./internal/verifyselect`
  exited 0. All selected tests passed; the suite guard reported its measured
  cost and no repository boundary violation or Helper Residue.

| Acceptance criterion | Focused evidence |
| --- | --- |
| Writable `TMPDIR` passes the shared cache to Go | `shared/first_worktree` and `shared/second_worktree` both received the same `<TMPDIR>/roundfix-go-build-<uid>` despite different working directories. `TMPDIR` had a trailing slash; the cache existed and retained no probe files. |
| Unwritable `TMPDIR` passes the worktree fallback | `fallback/directory_unwritable` received the run directory's `.gocache` with `TMPDIR` mode `0555`. `fallback/cache_unwritable` also received that fallback when the shared directory existed with mode `0555` under a writable `TMPDIR`, proving the file-creation check. |
| Explicit `GOCACHE` reaches Go unchanged | `explicit/environment` preserved a path containing a space despite an unwritable `TMPDIR`; `explicit/command_line` preserved the command-line path over a conflicting environment value. |

### Execution boundary

The executor could not initialize the configured sandbox because its writable
worktree root contains a symlink. Reviewed escalation allowed the focused
checks; this turn does not claim a fresh sandbox probe. No declared Task
Verification, repository Verification or full suite ran. Status and settlement
remain Daemon-owned. No other Task, Task Graph, CI workflow, Makefile target or
ignore entry was edited; no commit, push or pull request was made.


### Verification feedback — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T221735Z_f00baa2097cf9b98/verification/batch-003-attempt-1.log`.
The Daemon's `make verify-changed` attempt did not pass. Its diagnostic names
`TestBudgetExceededRunRecordsItsReason` in `internal/daemon` and
`TestRemoveRegisteredWorktreeWaitsForTheAdminLock` in `internal/worktree`;
`internal/verifyselect` passed in that attempt. This is recorded Daemon evidence,
not an Agent rerun of the configured gate.

Inspection found that the budget test requires a first Task settlement after
real Git operations within a fixed 500 ms Run budget. The worktree test gives
Git common-directory discovery and lock acquisition one 250 ms context; its
Git runner wraps the subprocess exit error without mapping it to the expired
context. These are unchanged paths outside Task 03. Concurrent scheduling
pressure is consistent with both diagnostics, but the precise cause during
the full gate is not established by an isolated pass.

Focused feedback checks:

- `rtk proxy go test -count=1 -v -run '^(TestBudgetExceededRunRecordsItsReason|TestRemoveRegisteredWorktreeWaitsForTheAdminLock)$' ./internal/daemon ./internal/worktree`
  exited 0: both diagnostic-named tests passed in isolation. No repeat-until-pass
  loop or test relaxation was used.
- With `task_cache_root="${TMPDIR:-/tmp}"`,
  `GOCACHE="${task_cache_root%/}/roundfix-go-build-$(id -u)" rtk proxy go test -race -count=1 -v -run '^(TestMakeSharesOneGoBuildCache|TestVerifyContractsRunsEveryContract|TestVerifyChangedRunsTheSelectedContracts|TestMainAndReleaseRunEveryContract)$' ./internal/verifyselect`
  exited 0. All three cache acceptance cases and the adjacent Makefile wiring
  tests passed using the shared cache for the outer Go invocation; the suite
  guard reported no boundary violation or Helper Residue.

No defect requiring an implementation change within Task 03 was demonstrated.
Only this Result was updated during feedback. Follow-up outside this slice:
make the budget fixture independent of Git scheduling within its fixed budget,
and investigate context-error reporting during worktree common-directory
discovery. The unsuccessful full gate remains unresolved evidence for the
Daemon's next configured Verification; isolated checks do not replace it. The
Agent did not rerun declared Task Verification or the configured gate, edit
status, change either failing package, commit, push or open a pull request.

## Carry-forward provenance

- Source Run: `run_20261012T221735Z_f00baa2097cf9b98`
- Source commit: `460365109b9366306162e97487bcb2de107b9ef1`
