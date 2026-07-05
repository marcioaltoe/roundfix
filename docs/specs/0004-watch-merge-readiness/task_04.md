---
task: task_04
spec: 0004-watch-merge-readiness
status: completed
type: backend
complexity: high
---

# Task 04: End watch on merge readiness

## Overview

Implement ADR-0019: with `--until-clean`, a Final Push no longer ends the Run
immediately — the loop polls the Review Source's status check on the pushed
head commit and ends Clean only on success with no new Review Issues, within
the existing bounds. Verifiable through clock-stepped loop tests over a
scripted check source plus a gh-adapter unit test.

## Requirements

1. MUST add a check seam to the watch package (`CheckFunc(ctx, headSHA) →
   pending|success|failure|missing`) beside the existing status/fetch/resolve
   adapters.
2. MUST implement the confirm phase: after a Final Push under
   `--until-clean`, poll the check on `poll_interval` (first check immediate,
   per task_01's ordering), bounded by the review timeout and remaining Max
   Rounds: `success` → Clean; `failure`, or settled with new Review Issues →
   next Round through the existing fetch path; `missing` → Clean plus one
   stderr note; bounds exhausted → existing `TimedOut`/`MaxRoundsReached`.
3. MUST implement the gh adapter: read the check runs for the pushed head
   commit through the existing GitHub CLI integration, filtered to the
   Review Source's check, mapped to the four states; transient gh failures
   surface as retryable poll errors within the timeout.
4. MUST leave non-`--until-clean` watch and plain resolve unchanged, and keep
   every existing outcome's exit-code mapping.
5. MUST journal the confirm phase through existing daemon event kinds
   (status/outcome), no new Run Event vocabulary.

## Subtasks

- [x] CheckFunc seam and Request wiring
- [x] Confirm phase in the loop with bounds accounting
- [x] gh check-runs adapter with app/check filtering
- [x] Scripted-check tests: success, failure→next Round, missing, timeout
- [x] CLI wiring and stderr note for the missing case

## Acceptance Criteria

- [x] Loop tests cover all four check states and both bound exhaustions with
      exact outcomes; the failure path demonstrably re-enters fetch and
      resolves new issues in the next Round.
- [x] The adapter unit test maps real gh JSON fixtures (success, failure,
      in_progress, absent) to the four states.
- [x] A Clean outcome is only reachable with check success or missing —
      asserted by a test that empties the local queue while the check stays
      pending and observes the Run keep polling.
- [x] Full suite passes; only deliberate ADR-0019 assertion changes.

## Verification

- `rtk go test ./internal/watch/ ./internal/cli/` — expected: all tests pass.
- `rtk go test -race ./internal/watch/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 2; Core Feature 2; Decisions. `_techspec.md` → Check
adapter, Watch loop changes, Risks, Build Order 4. Dogfood finding 20.
ADR-0019.

## Result

Status: completed.

Acceptance evidence:

- Loop coverage: `internal/watch` covers check success, failure, pending,
  and missing; timeout while pending; and MaxRoundsReached after a final-round
  failed check. `TestRunReentersFetchWhenHeadCheckFails` asserts the failure
  path fetches and resolves a second Round, and asserts later status/check
  calls use the updated head SHA.
- Adapter coverage: `TestHeadCheckMapsGitHubCheckRunJSON` parses gh
  check-runs JSON fixtures for success, failure, in_progress, and absent
  CodeRabbit checks, then asserts the four watch check states.
- Clean gating: `TestRunKeepsPollingWhenLocalQueueIsEmptyUntilHeadCheckSucceeds`
  empties the local queue on fetch, observes pending checks sleep and poll,
  and reaches Clean only after check success. Missing checks reach Clean and
  `TestRunWatchMissingHeadCheckPrintsCleanNote` asserts the single stderr
  note.
- Verification passed: `rtk go test ./internal/watch/ ./internal/cli/`,
  `rtk go test -race ./internal/watch/`, and `rtk go test ./...`.
