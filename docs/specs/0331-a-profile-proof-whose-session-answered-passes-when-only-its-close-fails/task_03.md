---
task: task_03
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
status: pending
type: backend
complexity: high
---

# Task 03: The next proof closes a recorded session, and the CLI reports the close truthfully

## Overview

The first proof of each Agent runner closes the Unclosed Proof Sessions an
earlier proof recorded, so no disposable session Roundfix failed to close is
forgotten. The CLI stops counting a close deadline as a setup timeout, and
`profiles validate` reports the warning a proof carries in text and JSON.
This completes the answer to the adopted Backlog Entry of 2026-10-27 (a
profile proof fails when closing its disposable session under load): the
3×2m park it records came from a close deadline the CLI retried as a setup
timeout.

## Requirements

1. MUST characterize first (lesson 9): before changing
   `profileProofSetupTimedOut`, write
   `TestACloseDeadlineIsNotRetriedAsASetupTimeout` in
   `internal/cli/profile_proof_retry_test.go` with a fake prover that returns
   `&agent.AgentSessionCleanupError{Session: …, Err: context.DeadlineExceeded}`
   and counts its calls, observe three calls and the `profile proof setup
   timed out on 3 attempts` message on the current tree, then make it assert
   the new contract: one call, classification `session_cleanup_failed`, no
   `setup timed out` text.
2. MUST implement `_techspec.md` → API Contract 5 in the new
   `internal/agent/unclosed_proof_session_reaper.go`, reading the records
   task_02 writes, and call it from
   `ProveExactSelection` and `probeFallbackCandidate` before each opens its
   disposable session: oldest record first, close under that proof's close
   deadline with the command rebuilt from the record, delete on a close or a
   proven absence, stop at the first failure with its warning, report and
   keep an unreadable record, stop on a cancelled context, at most once per
   `ACPXRunner`, never changing the proof's verdict (Invariants 2 and 4). A
   missing file on delete is not an error.
3. MUST implement API Contract 6: `profileProofSetupTimedOut` holds only for
   a `context.DeadlineExceeded` outside every `*agent.AgentSessionCleanupError`
   in the chain, with the command's context live. The three attempts, their
   waits and the `temporary` classification of a real setup timeout stay as
   they are (Invariant 6).
4. MUST implement API Contract 7: `profileProofReport.Warning` (`json:"warning,omitempty"`)
   carries `SelectionProof.Unclosed.Warning`, and `printProfilesValidateSuccess`
   prints `   warning: <text>` after that proof's reference lines. No other
   field, line or exit code changes.
5. MUST add the tests `_techspec.md` → Testing Approach names for task_03,
   each failing when its behavior is removed:
   `internal/agent/unclosed_proof_sessions_test.go`:
   `TestTheNextProofClosesARecordedSessionFirst` (a record written by one
   harness runner's failed close; a new runner's proof issues `sessions
   close <recorded>` before its own `sessions ensure`, and the record is
   gone), `TestARecordedSessionAcpxReportsAbsentIsForgotten` (close exit `4`),
   `TestAFailedReapKeepsTheRecordAndStops` (two records, the older one's
   close exits `2`: the newer one's close is never invoked, both records
   remain, the warning names the older record),
   `TestTheReapRunsOncePerRunner` (a second proof of the same runner issues
   no reaping close), `TestAnUnreadableRecordIsReportedAndKept`;
   `internal/cli/profile_proof_close_test.go`:
   `TestProfilesValidateReportsAnUnclosedProofSession` (a fake prover's proof
   carries `Unclosed`: exit `0`, the text line and the JSON `warning`).
   Tests in `internal/agent` use the fake acpx harness, whose `HOME` is a
   temporary directory; none uses a real-time budget. Each top-level test
   calls `t.Parallel()` first or carries a `// Sequential: <reason>`
   comment.

## Subtasks

- [ ] Characterize the CLI retry of a close deadline on the current tree.
- [ ] Add the reaping pass and call it from both proofs.
- [ ] Narrow the setup-timeout rule.
- [ ] Add the `warning` field and text line.
- [ ] Add the named tests.

## Acceptance Criteria

- [ ] A recorded session is closed by the next runner's first proof and its record deleted; a failed close keeps it and stops the pass.
- [ ] A close deadline is never retried or reported as a setup timeout.
- [ ] `profiles validate` shows the warning in text and JSON.

## Context

- instruction: `docs/adr/0344-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails-and-the-next-proof-closes-it.md`
- instruction: `docs/user-guide/commands/profiles.md`
- instruction: `.agents/skills/roundfix/references/profiles.md`
- instruction: `docs/user-guide/configuration.md`
- creates: `internal/agent/unclosed_proof_session_reaper.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/fallback.go`
- interface: `internal/agent/acpx_runner.go`
- creates: `internal/agent/unclosed_proof_sessions_test.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/profile_proof_retry_test.go`
- creates: `internal/cli/profile_proof_close_test.go`

## Consumers

- The reaping pass — consumers: `./internal/agent`
- `profileProofSetupTimedOut` and the `roundfix/profiles-validate/v1` report — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheNextProofClosesARecordedSessionFirst|TestARecordedSession.*|TestAFailedReap.*|TestTheReapRunsOncePerRunner|TestAnUnreadableRecord.*|TestAProofWhoseSessionAnswered.*|TestProveExactSelection.*|TestACPXProbe.*|TestExactSelectionProofSetup.*|TestACloseDeadlineIsNotRetriedAsASetupTimeout|TestProfileProof.*|TestProfilesValidate.*)$' ./internal/agent ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheNextProofClosesARecordedSessionFirst TestARecordedSessionAcpxReportsAbsentIsForgotten TestAFailedReapKeepsTheRecordAndStops TestTheReapRunsOncePerRunner TestAnUnreadableRecordIsReportedAndKept TestACloseDeadlineIsNotRetriedAsASetupTimeout TestProfilesValidateReportsAnUnclosedProofSession TestProfileProofReportsAThirdTimeoutAsTemporary TestProfileProofPassesOnTheThirdAttempt TestProfilesValidateDeduplicatesProofsAndReportsEveryReference TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/agent ./internal/cli` — expected: exit 0. Before this Task the reaping, CLI retry and report tests do not exist, so the first pass check fails. Afterwards they pass with the existing proof, retry and validate tests that consume the changed contracts, and the packages build for Windows.

## References

- `_prd.md` → Goals; Core Features 3-5; User Experience; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior 5, 7; Interfaces; Data Models; API Contracts 5-7; Invariants 2, 4, 6; Testing Approach; Build Order 2
- ADR-0344