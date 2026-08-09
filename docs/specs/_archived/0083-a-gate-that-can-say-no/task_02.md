---
task: task_02
spec: 0083-a-gate-that-can-say-no
status: completed
type: test
complexity: medium
---

# Task 02: Give the coverage contract a semantic owner and make it true

## Overview

The live coverage contract sits inside an archived Spec, so the sanctioned
repair for a renamed test is exactly the edit the archived-Spec rule forbids.
This task moves the record to a semantic owner outside `docs/specs/`, points the
test at the single new path, and brings the record back into agreement with the
suite — which is what makes the tree green for the first time in this Spec.

## Requirements

1. MUST move the coverage record out of the archived Spec to
   `docs/references/coverage-record.json` using a rename that preserves its
   history, not a copy and not a stub.
2. MUST reduce the test's two path constants to one naming the new owner, so
   there is no longer a fallback that can silently resolve to a Spec directory.
3. MUST bring the record into agreement with the current suite so the coverage
   invariant passes, using the sanctioned regeneration path rather than
   hand-editing entries.
4. MUST leave the archived Spec otherwise byte-identical; the authorized move is
   the only permitted change under `docs/specs/_archived/`.
5. MUST keep the invariant's teeth: a renamed or removed test must still fail
   the check after this task.
6. MUST change only these repository-relative paths plus this Task file:
   `internal/spec/coverage_test.go`,
   `docs/specs/_archived/0071-verification-cost/coverage-record.json` (removal
   by move), and `docs/references/coverage-record.json` (its destination). Any
   other changed path fails this Task.

## Subtasks

- [x] Move the record to its semantic owner, preserving history.
- [x] Collapse the two path constants into one naming the new owner.
- [x] Regenerate the record so it agrees with the current suite.
- [x] Prove the invariant still fails on a removed test.
- [x] Confirm the archived Spec has no change beyond the authorized move.

## Acceptance Criteria

- [x] The coverage record exists at `docs/references/coverage-record.json` and
      no longer exists under `docs/specs/_archived/`.
- [x] The test resolves its record through one path, and that path is outside
      `docs/specs/`.
- [x] The coverage invariant passes against the current suite.
- [x] Deleting or renaming any recorded test still makes the invariant fail,
      proven by observation rather than asserted.
- [x] No file under `docs/specs/_archived/` changed except the moved record.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- instruction: `docs/findings/2026-08-07-a-live-contract-lives-inside-an-archived-spec.md`
- interface: `internal/spec/coverage_test.go`

## Verification

- `test -f docs/references/coverage-record.json` — expected: exits 0.
- `test -e docs/specs/_archived/0071-verification-cost/coverage-record.json ; test $? -eq 1` — expected: exits 0, proving the record no longer lives in the archived Spec.
- `if grep -n 'docs/specs' internal/spec/coverage_test.go; then exit 1; fi` — expected: exits 0 and prints nothing, proving no remaining path constant points into the Spec tree.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v > /tmp/task_02-1.log 2>&1 && grep -q '^--- PASS: TestCoverageEquivalence' /tmp/task_02-1.log` — expected: exits 0, proving the invariant passes at its new home.
- `(git diff --name-only HEAD; git ls-files --others --exclude-standard) | grep -v -E '^(internal/spec/coverage_test\.go|docs/references/coverage-record\.json|docs/specs/_archived/0071-verification-cost/coverage-record\.json|docs/specs/0083-a-gate-that-can-say-no/task_02\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 3 and 4; Data Models.
- `_prd.md` → Core Feature 6; Goal 4.

## Result

### Implementation

- Moved the existing record with `git mv` from the archived Spec to
  `docs/references/coverage-record.json`. The combined diff detects an 84%
  rename after regeneration, so the destination retains the source history
  rather than introducing a copy or stub.
- Replaced the archived/active fallback constants with the single
  `coverageRecordPath` constant and removed fallback resolution.
- Re-recorded the destination through the test harness's
  `-update-coverage-record` flag. No record entry was edited by hand.

### Focused checks

- Pre-change signal:
  `rtk proxy env GOCACHE=<worktree>/.gocache go test -buildvcs=false ./internal/spec -run TestCoverageEquivalence -count=1`
  failed with the stale recorded CLI identities and current-suite additions.
- Sanctioned regeneration:
  `rtk proxy env GOCACHE=<worktree>/.gocache go test -buildvcs=false ./internal/spec -run TestCoverageEquivalence -count=1 -args -update-coverage-record`
  passed (`ok roundfix/internal/spec`, 7.667s).
- Current-suite check after regeneration:
  `rtk proxy env GOCACHE=<worktree>/.gocache go test -buildvcs=false ./internal/spec -run TestCoverageEquivalence -count=1`
  passed (`ok roundfix/internal/spec`, 4.418s).
- Negative observation: a disposable archive of `HEAD` received the current
  coverage test and record, then renamed
  `TestCompareCoverageRecordsReportsMissingTest` to
  `TestCompareCoverageRecordsReportsMissingTestRenamed`. Running the same
  focused check in that scratch tree failed with both the renamed addition and
  the missing recorded identity. The scratch tree was then removed.
- `rtk rg --files docs | rtk rg 'coverage-record\\.json$'` printed only
  `docs/references/coverage-record.json`.
- `rtk rg -n 'coverageRecord.*Path|docs/specs' internal/spec/coverage_test.go`
  printed only the new constant and its use; it found no Spec path.
- `rtk git diff --exit-code HEAD -- docs/specs/_archived
  ':(exclude)docs/specs/_archived/0071-verification-cost/coverage-record.json'`
  exited 0, proving every other archived-Spec byte matches `HEAD`.
- `rtk git diff --name-status HEAD` reported only the 84% record rename,
  `internal/spec/coverage_test.go`, and this Task file. The complete status
  audit reported no untracked path.

The first pre-change probe without the repository-local `GOCACHE` did not
reach the suite because the sandbox denied the global Go cache. The repeated
probe above used the repository's documented local cache and produced the red
signal.

The Daemon-owned commands under `## Verification` and the repository-wide
gate were not run in this Agent turn.

### Verification feedback — attempt 1

- Inspected the Daemon diagnostic artifact for the failed absence check. It
  confirmed that the forbidden path was absent, while `grep -c` returned its
  standard no-match status and the following assertion required a match
  status. The implementation satisfied the criterion; the authored command's
  exit-status assertion contradicted its stated zero-count expectation.
- Replaced that line with the repository's portable conditional absence form:
  it reports matching lines and exits 1 when the forbidden path exists, and
  exits 0 without output when the path is absent.
- Did not rerun either the failed command or its replacement. The Daemon owns
  the single full Verification retry.
