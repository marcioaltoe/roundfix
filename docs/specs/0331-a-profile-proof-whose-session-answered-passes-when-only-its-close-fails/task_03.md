---
task: task_03
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
status: completed
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

## Result

Implemented API Contracts 5–7 within this Task's slice. Task status remains
Daemon-owned. The initial working tree already contained the Daemon's status
change in this Task file; no other pre-existing changed paths were present.

- Recorded-session recovery: both exact proofs and fallback candidate probes
  run one shared `sync.Once` reaping pass before opening their disposable
  session. Records sort by `recorded_at`, then session name. Close commands
  use the recorded runtime and working directory, under the current proof's
  close deadline. Successful or proven-absent closes delete their records;
  an already-missing record is accepted. A failed close warns and stops,
  retaining the record and later records. Unreadable records warn and remain.
  Cancellation stops the pass; its outcome does not decide the proof verdict.
- Close-deadline retry: before changing production code,
  `rtk proxy go test ./internal/cli -run '^TestACloseDeadlineIsNotRetriedAsASetupTimeout$' -count=1 -v`
  characterized **three calls** and `profile proof setup timed out on 3 attempts
  of 30s each` for a fake prover returning `AgentSessionCleanupError` wrapping
  `context.DeadlineExceeded`. The test now requires one call, no backoff,
  classification `session_cleanup_failed`, and no `setup timed out` text.
  Error-tree traversal excludes cleanup subtrees while preserving independent
  setup deadlines in joined errors. Existing third-attempt setup timeout and
  third-attempt success checks retain their previous behavior.
- Warning reporting: passed proofs copy `Unclosed.Warning` into the optional
  JSON `warning` field and print `   warning: <text>` immediately after their
  reference lines. The command test checks text and JSON with exit zero.

Focused evidence for the acceptance criteria:

1. Recovery and retention: all five required agent tests passed. Additional
   fake-harness tests passed for a cancelled proof and for a fallback candidate
   reaping before ensure, followed by an exact proof sharing the once guard.
   Every new top-level test starts with `t.Parallel()`; no new test uses a
   real-time budget or a real adapter.
2. Retry classification: the characterization test passed under its new
   one-call contract, and the joined-error table passed for wrapped cleanup,
   independent setup deadlines, and cancelled command contexts. The existing
   `TestProfileProofReportsAThirdTimeoutAsTemporary` and
   `TestProfileProofPassesOnTheThirdAttempt` passed.
3. Warning output: `TestProfilesValidateReportsAnUnclosedProofSession` passed
   for both output modes, including the literal JSON field and text placement.

Final focused command (exit 0, both packages):

```sh
rtk proxy go test ./internal/agent ./internal/cli -run '^(TestTheNextProofClosesARecordedSessionFirst|TestARecordedSession.*|TestAFailedReapKeepsTheRecordAndStops|TestTheReapRunsOncePerRunner|TestAnUnreadableRecordIsReportedAndKept|TestACloseDeadlineIsNotRetriedAsASetupTimeout|TestProfilesValidateReportsAnUnclosedProofSession|TestProfileProofTimeoutDistinguishesJoinedSetupAndCleanupDeadlines|TestProfileProofReportsAThirdTimeoutAsTemporary|TestProfileProofPassesOnTheThirdAttempt)$' -count=1
```

Temporary `go test -overlay` probes, with altered files outside the repository,
confirmed behavioral failures when reaping was removed (the next-close,
absence, stop-on-failure and unreadable-record tests), the once guard was
removed (the once-per-runner test), age ordering was reversed (the failed-reap
test), warning propagation was removed (both command output modes), and broad
deadline retry was restored (three calls and `temporary` instead of cleanup).
Each probe compiled and failed at the intended behavioral assertion; no mutant
was retained. `rtk proxy git diff --check` exited 0.

The declared Verification, including the Windows build, and repository
Verification were not run; they remain for the Daemon. No Task Graph, other
Task, tooling configuration, commit, push, or Pull Request was changed or made.

### Verification feedback repair — attempt 1

Inspected the Daemon's `batch-002-attempt-1.log`. The configured incremental
Verification reported a failure in
`TestTheProofCloseDeadlineFollowsTheSetupDeadline/write_error`: the reaper
added a store-read warning before the existing record-write recovery warning
when `.roundfix` was a regular file. Such a root cannot contain recorded
sessions. The reaper now recognizes that case before reading the records
subdirectory, leaving the writer to report the failed record and manual close
command. Actual unreadable records and store-read errors retain their warnings.
No existing test assertions were changed.

Focused repair check, exit 0:

```sh
rtk proxy go test ./internal/agent -run '^(TestTheProofCloseDeadlineFollowsTheSetupDeadline|TestTheNextProofClosesARecordedSessionFirst|TestAFailedReapKeepsTheRecordAndStops|TestAnUnreadableRecordIsReportedAndKept)$' -count=1 -v
```

All four tests passed, including all three close-deadline subcases. This checks
the reported warning-count regression, successful recovery, failed-close
retention and unreadable-record reporting. The Daemon's full configured
Verification remains outstanding; neither it nor the Task's declared
Verification was rerun in this repair turn. Task status remains unchanged.
