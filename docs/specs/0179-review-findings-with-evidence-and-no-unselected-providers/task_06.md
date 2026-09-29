---
task: task_06
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: completed
type: backend
complexity: low
---

# Task 06: A disposition is reserved atomically

## Overview

Corrective Task from the pre-PR review of 2026-09-29. `roundfix review dispose` in `internal/cli/review.go` checks the dispositions ledger for an existing entry and then appends, with no lock or atomic reservation between the two steps; two concurrent invocations for the same finding can both pass the duplicate check and append two dispositions, breaking the at-most-one-disposition contract.

## Requirements

1. MUST make the duplicate check and the append one critical section per ledger file, guarded by an exclusive advisory file lock (for example `flock` on a sibling lock file in the Artifact Directory), held only for the read-check-append and released on every path.
2. MUST keep the existing exit codes, messages and ledger format unchanged; the loser of a race exits `2` with the existing "already disposed" refusal and appends nothing.
3. MUST put the new tests in a file of their own.

## Subtasks

- [x] Implement the requirements above.
- [x] Add a named test for each acceptance criterion.

## Acceptance Criteria

- [x] Many concurrent `review dispose` invocations for the same finding (helper subprocesses) append exactly one ledger entry, and every other invocation exits `2` with the existing refusal.
- [x] Concurrent dispositions of two different findings both succeed.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_dispose_lock_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestConcurrentDisposeOfOneFindingAppendsOnce|TestConcurrentDisposeOfTwoFindingsBothSucceed)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestConcurrentDisposeOfOneFindingAppendsOnce TestConcurrentDisposeOfTwoFindingsBothSucceed; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented an exclusive advisory lock on a sibling
`pre-pr-review-dispositions.jsonl.lock` file. The authoritative duplicate
check and append now run under one lock, and every acquired lock is unlocked
and closed on success, duplicate refusal, read failure, or append failure.
Unix uses `flock`; Windows uses `LockFileEx`. The existing ledger format,
success output, refusal text, and exit codes are unchanged.

Focused-check evidence:

- Acceptance criterion 1: `GOCACHE=/tmp/roundfix-task06-gocache go test
  -count=1 -run '^TestConcurrentDisposeOfOneFindingAppendsOnce$'
  ./internal/cli` exited 0. The test starts 32 synchronized helper processes,
  observes one exit 0, 31 exits 2 with the existing `already has a
  disposition` refusal, and one ledger entry.
- Acceptance criterion 2: `GOCACHE=/tmp/roundfix-task06-gocache go test
  -count=1 -run '^TestConcurrentDisposeOfTwoFindingsBothSucceed$'
  ./internal/cli` exited 0. Both helper processes exit 0 and the ledger
  contains one entry for each finding.
- Regression coverage: `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1
  -run '^TestReviewDispose' ./internal/cli` exited 0.
- Platform build: `GOCACHE=/tmp/roundfix-task06-gocache GOOS=windows
  GOARCH=amd64 go build ./internal/cli` exited 0. A preceding Windows
  `go test -c` attempt was blocked by pre-existing Unix-only code in
  `internal/cli/implement_test.go` (`syscall.Mkfifo`, `Setpgid`, and `Kill`).
- `git diff --check` exited 0.

The Daemon-owned `## Verification` command was not run.
