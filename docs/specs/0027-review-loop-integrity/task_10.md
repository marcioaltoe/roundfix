---
task: task_10
spec: 0027-review-loop-integrity
status: pending
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

- [ ] Compute per-Run counts from the Run's own issue set alongside the existing cumulative scan
- [ ] Update the report renderer with the two labeled summary lines
- [ ] Add reason suffixes to failed/unresolved/invalid issue lines
- [ ] Rendering tests with fixtures where per-Run and cumulative counts differ, and with reasons present and absent

## Acceptance Criteria

- [ ] A fixture with prior-round artifacts on disk yields different per-Run and cumulative counts, each labeled unambiguously
- [ ] Failed issue lines show the terminal reason; issues without reasons render cleanly
- [ ] Existing report consumers (tests asserting the old summary) are updated, not appended around
- [ ] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/rounds/rounds.go`

## Verification

- `go test ./internal/cli/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 4, User Story 8, Core Feature 9, Open Questions (invalid reason lines); `_techspec.md` → Build Order 9, API Contracts (report format).
