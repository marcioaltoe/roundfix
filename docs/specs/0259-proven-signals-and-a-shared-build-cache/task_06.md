---
task: task_06
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
type: test
complexity: medium
---

# Task 06: The residue scan survives an unreadable process, and the disposal lock tests wait on a synchronization point

## Overview

Corrective Task 1 of 2, after the QA Report of 2026-10-12 at `9479226c`
(Run `run_20261012T221735Z_f00baa2097cf9b98`) failed rows R04 and R10.

- QA F1: in the assembled full suite, `TestGuardReportsAndKillsHelperResidue`
  saw fixture exit 0 with
  `suiteguard: helper residue check unavailable: read process 45114 environment: short argument list for pid 45114`.
  One live process whose `kern.procargs2` does not parse
  (`parseDarwinProcessArguments` in `internal/suiteguard/residue_darwin.go`)
  makes `checkHelperResidue` give up the whole scan, so Helper Residue goes
  unreported and unkilled while that process lives.
- QA F2: `TestRemoveRegisteredWorktreeWaitsForTheAdminLock` failed with
  `git rev-parse --path-format=absolute --git-common-dir` `signal: killed`:
  its 250 ms deadline also bounds the real `git rev-parse` that runs before the
  admin lock, so under load the deadline kills Git instead of proving the wait.
  This failure stops the ten-suite measurement of Requirement 6 of task_05, the
  acceptance of F6, so this Spec fixes it. The Backlog Entry of 2026-10-12 on
  the same test, captured after this Spec's planning, is closed by the Spec
  that adopted it.

## Requirements

1. MUST classify a `kern.procargs2` parse failure of a process still in the
   process table (short process arguments, missing executable, short argument
   list) as `errResidueProcessChanging`, so `checkHelperResidue` retries it
   inside the same grace, exactly as it retries `EIO` and `EINVAL`.
2. MUST keep scanning when a process is still unreadable at the end of the
   grace: the scan writes one line
   `suiteguard: helper residue check skipped an unreadable process: pid <pid>: <error>`
   for each such process, still reports and kills every marked process it can
   read, and returns whether it found residue. It MUST NOT treat an unreadable
   process as residue or signal it. A failure to read the process table itself
   stays `suiteguard: helper residue check unavailable: <error>`.
3. MUST prove Requirements 1 and 2 with a new test
   `TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning` in
   `internal/suiteguard`, which injects one unreadable process beside one marked
   process through a seam in `residue.go` (no real unreadable process is
   needed), and requires the skipped line, the residue line for the marked
   process and its death. The Result MUST record a sabotage (restoring the
   early return) and the failure it produced.
4. MUST make `RemoveRegisteredWorktree` and `AddDetachedWorktree` in
   `internal/worktree` delegate to unexported `removeRegisteredWorktree` and
   `addDetachedWorktree` that take a `gitRunner`, keeping the exported
   signatures and behavior.
5. MUST rewrite `TestRemoveRegisteredWorktreeWaitsForTheAdminLock` and
   `TestAddDetachedWorktreeWaitsForTheAdminLock` in
   `internal/worktree/disposal_test.go`, keeping their names and
   `t.Parallel()` first. Each holds the lock with `holdDisposalAdminLock`, runs
   the operation in a goroutine with a cancellable context and a test runner
   that delegates to `execGitRunner{}`, waits for the runner's report that the
   Git common directory is resolved, cancels, and requires `context.Canceled`,
   no `worktree` command started while the lock was held, and the worktree
   unchanged. After releasing the holder it runs the same operation with the
   test's context and requires success. Neither test passes a deadline to the
   operation, sleeps, or decides an assertion on elapsed time; a bounded wait
   may only fail a test that would otherwise hang.
6. MUST NOT change `TestRemoveRegisteredWorktreeKeepsADirtyWorktree`,
   `internal/worktree/adminlock_test.go`, the Linux residue scan, or any file
   outside its Context.

## Subtasks

- [ ] Retry and then skip an unreadable process in the residue scan.
- [ ] Add the residue seam and its test, with the sabotage.
- [ ] Add the runner-taking disposal functions and rewrite the two lock tests.

## Acceptance Criteria

- [ ] An unreadable process no longer makes the residue check unavailable.
- [ ] Neither disposal lock test passes a deadline to the operation it proves.
- [ ] Both lock tests pass 50 consecutive times with `-race -parallel 16`.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `docs/adr/0126-the-suite-proves-its-own-isolation-per-package.md`
- interface: `internal/suiteguard/residue.go`
- interface: `internal/suiteguard/residue_darwin.go`
- interface: `internal/suiteguard/suiteguard_test.go`
- creates: `internal/suiteguard/residue_test.go`
- interface: `internal/worktree/disposal.go`
- interface: `internal/worktree/disposal_test.go`
- instruction: `internal/worktree/adminlock_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning|TestGuardReportsAndKillsHelperResidue|TestGuardPassesWhenEveryHelperEnded)$' ./internal/suiteguard 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning TestGuardReportsAndKillsHelperResidue TestGuardPassesWhenEveryHelperEnded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'helper residue check skipped an unreadable process' internal/suiteguard/residue.go || { printf 'the skipped-process line is missing\n' >&2; exit 1; }`
- `if grep -q 'context.WithTimeout' internal/worktree/disposal_test.go; then printf 'a disposal test still passes a deadline\n' >&2; exit 1; fi; out="$(go test -count=50 -race -parallel 16 -v -run '^(TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRemoveRegisteredWorktreeWaitsForTheAdminLock TestAddDetachedWorktreeWaitsForTheAdminLock; do test "$(printf '%s\n' "$out" | grep -c -- "--- PASS: $name ")" -eq 50 || { printf 'fewer than 50 passes: %s\n' "$name" >&2; exit 1; }; done`

## References

- QA Report `qa/qa-report-2026-10-12.md` at `9479226c` → F1, F2; rows R04, R10
- `_techspec.md` → API Contract 5; API Contract 6; Testing Approach
- ADR-0266; ADR-0126

## Result

### Implementation

The Darwin scan collects per-process environment read errors alongside readable
Helper Residue instead of abandoning the process table. All three malformed
`kern.procargs2` cases wrap `errResidueProcessChanging`; live `EIO` and `EINVAL`
retain the same classification. The common checker retries unreadable processes
inside the existing ten-second grace, prints one skipped-process line for each
still unreadable PID, and reports and kills only readable marked processes.
Process-table failures retain the unavailable diagnostic. The scan dependency
is passed through the unexported `checkHelperResidueWithScan` seam, with no
process-wide test overrides.

The exported disposal functions delegate to runner-taking unexported functions.
Both lock tests keep `t.Parallel()` first, hold the existing administration-lock
fixture, and start the operation with `context.WithCancel(t.Context())`. Their
runner delegates to `execGitRunner{}` and reports successful common-directory
resolution after real Git returns. The tests then cancel, require
`context.Canceled`, assert that no worktree command started and that the path,
HEAD and registration remain as expected, release the holder, and require the
same operation to succeed with `t.Context()`. Channel waits use `testwait.Until`
only to fail hangs; they do not cancel the operation on a deadline. Each waiter
is cancelled and joined during cleanup.

### Focused checks and acceptance evidence

| Acceptance criterion | Current evidence |
| --- | --- |
| An unreadable process no longer makes the residue check unavailable. | The named injection regression passed with exactly one skipped line, the marked child's residue line, repeated reads and the child's SIGKILL exit. Additional tests passed for a read that recovers, two unreadable processes without residue (returns false), and an unavailable process table. The temporary Darwin parser probe passed all three malformed-input classifications. |
| Neither disposal lock test passes a deadline to the operation it proves. | Source inspection confirmed `WithCancel`, no `WithTimeout` or `WithDeadline`, no sleeps or elapsed-time assertions, and `t.Parallel()` first in each test. Both revised tests passed the focused shuffled race run below, including cancellation while locked and success after release. |
| Both lock tests pass 50 consecutive times with `-race -parallel 16`. | Focused evidence is three consecutive passes per lock test with those flags and shuffled order. The declared 50-run Verification is reserved for the Daemon and has not been run by this Agent. |

Commands and outcomes:

- `go test -race -shuffle=on -count=1 -v -run 'Test(HelperResidue|DarwinResidue|ResidueMark|GuardReportsAndKillsHelperResidue|GuardPassesWhenEveryHelperEnded)' ./internal/suiteguard`
  — exit 0; all nine selected top-level tests passed, including the existing
  real process-table residue and clean-helper fixtures.
- `go test -count=3 -race -parallel 16 -shuffle=on -v -run '^(TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock|TestRemoveRegisteredWorktreeKeepsADirtyWorktree)$' ./internal/worktree`
  — exit 0; three passes each for both rewritten lock tests and the unchanged
  dirty-worktree refusal test. Shuffle seed: `1791852470747206000`.
- `go test -race -count=1 -v -overlay /tmp/task06-darwin-parser-overlay.json -run '^TestTask06DarwinParseFailuresRetryable$' ./internal/suiteguard`
  — exit 0. A temporary overlay appended a probe to the existing Darwin test
  source without modifying it in the worktree. Short process arguments,
  missing executable and short argument list each matched
  `errResidueProcessChanging` through `errors.Is`.
- `GOOS=darwin go build ./internal/suiteguard`,
  `GOOS=linux go build ./internal/suiteguard`, and
  `GOOS=windows go build ./internal/suiteguard` — each exited 0.
- `GOOS=windows go test -c -o /tmp/task06-suiteguard-windows.test ./internal/suiteguard`
  — exit 0; the shared tests also compile for the unsupported-platform build.
- `git diff --check` and a focused Python source/scope inspection — passed.
  Only this Task file and the five declared source/test files differ. All
  authored Task bytes before `## Result` match the turn-start snapshot;
  `TestRemoveRegisteredWorktreeKeepsADirtyWorktree` is byte-identical to HEAD.
  `internal/worktree/adminlock_test.go`, the Linux scan, the Task Graph and
  every other Task file are unchanged.

### Red starting point and sabotage

Before the scan behavior changed, the new seam and regression ran with the
original error branch:

`go test -race -count=1 -run '^TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning$' ./internal/suiteguard`

It exited 1 with `residue check returned false beside an unreadable process`,
missing skipped/residue diagnostics, and the unavailable line. The fixture's
lifeline closed in cleanup so the failed regression left no live helper.

After the repair, a temporary compiler overlay restored the original error
branch, including its grace retry and early `return false`, in `residue.go`:

`go test -race -count=1 -overlay /tmp/task06-residue-sabotage.json -run '^TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning$' ./internal/suiteguard`

The sabotage exited 1 and produced:

```text
suiteguard: helper residue check unavailable: read process 51696 environment: process argument memory changed during read: short argument list for pid 51696
residue check returned false beside an unreadable process
want one diagnostic "suiteguard: helper residue check skipped an unreadable process: pid 51696: process argument memory changed during read: short argument list for pid 51696\n"
want one diagnostic "suiteguard: helper residue: pid 51786: marked fixture\n"
unreadable process misclassified
```

The overlay never replaced the worktree file. The final focused shuffled race
run above used the repaired source and passed the regression, including death
of the marked fixture. Temporary overlays and compiled probes live under
`/tmp`, outside the repository diff.

### Settlement boundary

Task status, Subtasks and Acceptance Criteria checkboxes are left for the
Daemon. Neither declared Verification command nor selected repository
Verification was executed. No commit, push or pull request was created.
No follow-up changes were included.

### Verification Feedback, attempt 1

Inspected the Daemon artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261013T004148Z_5319d413bd50ad25/verification/batch-001-attempt-1.log`.
The Daemon's `make verify-changed` ended with exit 2. Its reported test failure
was `TestBudgetExceededRunRecordsItsReason` in
`internal/daemon/task_engine_test.go`; both `internal/suiteguard` and
`internal/worktree` had successful package results in that same gate run.
These are recorded Daemon results, not checks rerun by this Agent.

The failing Daemon test uses `time.Now` and a fixed 500 ms Run budget, yet
requires the first Task and its commit to settle before the budget expires.
The observed failure was zero settled Tasks and commits. Source inspection
confirmed that `TaskCycle` starts its budget watchdog before Task processing;
the failure is consistent with that wall-clock budget expiring before the
required first settlement. The diagnostic does not identify the exact
pre-settlement step that consumed the budget. Both
`internal/daemon/task_engine_test.go` and `internal/daemon/task_engine.go`
were compared byte-for-byte with HEAD and are unchanged by this Task.

No in-scope source repair is indicated by this feedback. Requirement 6 forbids
editing files outside this Task's Context, which excludes the failing Daemon
test and engine. Follow-up: separately scope a correction to
`TestBudgetExceededRunRecordsItsReason` so its first-settlement expectation
is established through synchronization instead of assuming real Git and Task
setup finish inside 500 ms. This follow-up is recorded only; no Daemon code,
other Task file, Task Graph or tooling setting was changed.

The adjacent `batch-001-probe-1.log` and `batch-001-probe-2.log` contain
pre-change refusal evidence, not post-implementation passes. They therefore
do not establish the 50-run acceptance criterion. That criterion and the
configured Verification retry remain Daemon-owned.

Focused inspection in this feedback turn checked the failing sources against
HEAD and the changed-file scope. No tests, declared Verification commands or
repository Verification were rerun. The repair-turn diff adds only this
Result addendum; the task_06 implementation and its status are preserved.

## Carry-forward provenance

- Source Run: `run_20261013T004148Z_5319d413bd50ad25`
- Source commit: `8e658da95741677763ea9cff2581edaf29e63888`
