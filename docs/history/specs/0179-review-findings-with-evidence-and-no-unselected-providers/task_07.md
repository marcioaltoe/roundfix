---
task: task_07
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: completed
type: backend
complexity: low
---

# Task 07: A user-supplied revision is never read as a Git option

## Overview

Corrective Task from the second pre-PR review of 2026-09-29. `roundfix review dispose --fixed-by <commit>` in `internal/cli/review.go` passes the operator's value straight to `git rev-parse`; a value such as `--output=<path>` is parsed as a Git option instead of a revision and can write to an arbitrary file.

## Requirements

1. MUST pass `--end-of-options` before every operator-supplied revision given to Git by `review dispose` (and any other `roundfix review` path that forwards an operator revision), so the value is always read as a revision.
2. MUST refuse a `--fixed-by` value that does not resolve to a commit with the existing refusal and exit `2`, including a value that starts with `-`, and write nothing.
3. MUST keep every existing test and message unchanged, and put the new tests in a file of their own.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a named test for each acceptance criterion.

## Acceptance Criteria

- [ ] `--fixed-by --output=<path>` exits `2`, creates no file at `<path>`, and appends nothing to the ledger.
- [ ] A valid descendant commit given to `--fixed-by` is still recorded.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_dispose_revision_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing|TestDisposeFixedByValidCommitIsStillRecorded)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing TestDisposeFixedByValidCommitIsStillRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

### Implementation

- `review dispose --fixed-by` now places `--end-of-options` before the
  operator-supplied revision passed to `git rev-parse --verify`.
- `review --base` applies the same guard to the other operator-supplied
  revision forwarded by `roundfix review`; later Git calls receive resolved
  commit IDs.
- Added the dedicated `review_dispose_revision_test.go` regression suite. Its
  Git boundary delegates normal operations to real local Git and models the
  option side effect deterministically, independent of the installed Git
  version.

### Focused checks

- Before the implementation change,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run '^TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing$' ./internal/cli`
  failed because the unguarded revision wrote the requested output path.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run 'Test(DisposeFixedBy|ReviewBaseOptionLike)' ./internal/cli`
  exited `0` after the change and exercised both acceptance tests plus the
  `review --base` companion regression.
- A sandboxed `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 ./internal/cli`
  run reached two unrelated force-stop integration tests and failed because
  process-table access was denied. The same command rerun with process-table
  permission exited `0` (`ok roundfix/internal/cli 80.548s`).
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache make verify-incremental`
  exited `0` with process-table permission, covering formatting, vet, all Go
  packages, skill checks, and the build.
- The Task's authored `## Verification` command was not run; Daemon
  Verification owns that command.

### Acceptance evidence

- `TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing` passed: an
  option-like `--fixed-by` value exits `2` with the existing refusal, creates
  no requested output file, and leaves the disposition ledger absent.
- `TestDisposeFixedByValidCommitIsStillRecorded` passed: a valid descendant
  commit exits `0`, is written as the ledger entry's `fixedBy`, and is emitted
  as the same JSON line on stdout.
