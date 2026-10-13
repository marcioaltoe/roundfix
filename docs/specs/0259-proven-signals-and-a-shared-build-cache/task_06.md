---
task: task_06
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
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
