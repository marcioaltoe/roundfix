---
task: task_10
spec: 0027-review-loop-integrity
status: failed
type: backend
complexity: medium
---

# Task 10: Split the report into per-Run and cumulative counts with reasons

## Overview

Make the final review report answer "what did THIS Run do" without artifact archaeology: the summary separates this Run's issue counts from the pull request's cumulative counts, and every failed or unresolved issue line carries its one-line terminal reason. Applies to the resolve and watch end-of-run reports.

## Requirements

1. MUST print two labeled summary lines: this Run's counts, computed from the Run's own issue set (the resolve selection, or the watch Run's accumulated Rounds), and the pull request's cumulative counts from the existing artifact scan.
2. MUST suffix each failed and unresolved issue line with its terminal reason when present, as a single concise line.
3. MUST keep issue status lines greppable and stable: one issue per line, status first-class, reason as a suffix.
4. MUST render the Clean Unverified outcome in the report vocabulary consistently with the glossary.
5. SHOULD include invalid issues' reason suffixes the same way (per the PRD open-question default).

## Subtasks

- [x] Compute per-Run counts from the Run's own issue set alongside the existing cumulative scan
- [x] Update the report renderer with the two labeled summary lines
- [x] Add reason suffixes to failed/unresolved/invalid issue lines
- [x] Rendering tests with fixtures where per-Run and cumulative counts differ, and with reasons present and absent

## Acceptance Criteria

- [x] A fixture with prior-round artifacts on disk yields different per-Run and cumulative counts, each labeled unambiguously
- [x] Failed issue lines show the terminal reason; issues without reasons render cleanly
- [x] Existing report consumers (tests asserting the old summary) are updated, not appended around
- [x] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/rounds/rounds.go`

## Verification

- `go test ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 4, User Story 8, Core Feature 9, Open Questions (invalid reason lines); `_techspec.md` → Build Order 9, API Contracts (report format).

## Result

- Implemented the review report split in `internal/cli`: issue lines now come from this Run's issue set, while the final summary prints two labeled lines — `This Run (...)` and `Pull Request cumulative`.
- Added reason suffixes for failed, unresolved, and invalid issue lines when `terminal_reason` is present; reasons are collapsed to one line and omitted cleanly when absent.
- Rendered `Clean Unverified` with the glossary spelling in the stdout report.
- Evidence for per-Run versus cumulative counts: `TestRunResolveHonorsRoundSelector` uses prior-round artifacts on disk and now expects `This Run` counts to differ from `Pull Request cumulative`.
- Evidence for reason suffixes and Clean Unverified vocabulary: `TestPrintReviewIssueReportSplitsRunAndCumulativeCountsAndReasons` covers invalid, failed, and unresolved lines with and without reasons.
- Verification passed: `go test ./internal/cli/...` (416 tests), `go test ./...` via `make verify` (1193 tests), `go build -buildvcs=false ./...`, and `make verify`.
- Verification blocker: exact task command `go build ./...` fails before compilation with `error obtaining VCS status: exit status 128`; this environment has an invalid parent Git marker at `/Users/marcio/.git`, matching the task_09 build blocker. Because the task's exact build verification does not pass, this task is settled as failed.
