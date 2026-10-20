---
task: task_02
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: completed
type: backend
complexity: medium
---

# Task 02: A timed-out profile proof gets a third attempt after a backoff

## Overview

Today a profile proof whose setup times out is tried once more at once, and a
second timeout parks a Delivery Queue item (`_prd.md` → Sources 1). This Task
makes three attempts in all, waiting 5 and then 10 seconds through an
injected sleeper, and makes the final error name the attempts, the deadline
and the config key. It is verifiable on its own through the CLI's injected
runner and sleeper, with no real wait.

## Requirements

1. MUST answer the Secondbrain Inbox Entry of 2026-10-20,
   [`references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md`](references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md)
   (`_prd.md` → Sources 1), by implementing API Contract 2 of `_techspec.md`
   in `proveProfileSelectionWithRetry` (ADR-0298), with
   `profileProofRetryWaits` and `profileProofTimeoutError` as
   `_techspec.md` → Interfaces sketches them.
2. MUST add the sleeper `profileProofSleep func(context.Context,
   time.Duration) error` to `commandDependencies`. Its default waits on a
   timer, returns `ctx.Err()` when the context ends first and starts no
   goroutine; the retry reads it through `commandDependenciesForContext`.
3. MUST keep the retry's scope as API Contract 2 states it:
   - only an error for which `profileProofSetupTimedOut` holds is retried,
     at most twice, after `profileProofRetryWaits[0]` (5 s) and then
     `profileProofRetryWaits[1]` (10 s);
   - any other error returns at once, unchanged;
   - a sleeper error, which a cancellation during the wait produces, returns
     at once with no further attempt and keeps today's unclassified
     cancellation report;
   - the third timeout is classified `temporary`, its text is `profile proof
     setup timed out on 3 attempts of <deadline> each: <last error>`, and its
     next action is exactly `rerun the command when load drops, or raise
     defaults.profile_proof_setup_timeout; the configured profile was not
     shown to be wrong`;
   - a retry that passes writes nothing to standard output or standard error.
4. MUST change these existing tests in
   `internal/cli/profile_proof_retry_test.go`, each a declared break, and
   inject a recording `profileProofSleep` in every test of the file:
   - `TestProfileProofRetriesATimedOutSetupOnce` also asserts exactly one
     wait, `profileProofRetryWaits[0]`;
   - `TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce` runs
     with a context carrying the test's command dependencies and asserts the
     same single wait;
   - `TestProfileProofReportsASecondTimeoutAsTemporary` is renamed
     `TestProfileProofReportsAThirdTimeoutAsTemporary`. It asserts three
     attempts, both waits in order, the `temporary` classification and the
     next action above in JSON mode, and, in a text-mode subtest, the empty
     standard output, exit `2` and the exact two standard error lines of API
     Contract 2.
5. MUST add these tests to the same file, each top-level test calling
   `t.Parallel()` first:
   - `TestProfileProofPassesOnTheThirdAttempt`: setup times out on the first
     two attempts and passes on the third; `profiles validate` exits `0`
     with three attempts, both waits in order and empty standard error;
   - `TestProfileProofStopsWhenCancelledDuringTheBackoff`: the sleeper
     cancels the command context and returns its error; there is one
     attempt, no classification and the error names `context canceled`.
6. MUST re-record the coverage record, because the rename removes a recorded
   test name: run
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   before handing back. The Daemon never regenerates a record, and the record
   is never hand-edited.
7. MUST search the code and tests for the old text `timed out twice` and
   change every occurrence; the guides and skills change in task_04.

## Subtasks

- [ ] Add the sleeper dependency and the attempts with their waits.
- [ ] Change the final error, its next action and its text.
- [ ] Update the three existing tests and add the two new ones.
- [ ] Re-record the coverage record.

## Acceptance Criteria

- [ ] A setup that times out once or twice passes on a later attempt after waits of 5 s and then 10 s, and writes nothing.
- [ ] A third timeout is `temporary` and its message and next action match API Contract 2 byte for byte.
- [ ] A cancellation during a wait starts no further attempt.
- [ ] Every other failure is still reported after one attempt.

## Consumers

- `profileProofTimeoutError`'s text, its next action and the attempt count — consumers: `./internal/cli`
- the coverage record's test names — consumers: `./internal/spec`

## Context

- instruction: `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- instruction: `internal/cli/profile_preflight.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/profile_proof_retry_test.go`
- interface: `docs/references/coverage-record.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestProfileProofPassesOnTheThirdAttempt|TestProfileProofReportsAThirdTimeoutAsTemporary|TestProfileProofStopsWhenCancelledDuringTheBackoff|TestProfileProofRetriesATimedOutSetupOnce|TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce|TestProfileProofDoesNotRetryARejectedSelection|TestProfileProofDoesNotRetryACancelledCommand|TestProfileProofPassesTheConfiguredSetupDeadline)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProfileProofPassesOnTheThirdAttempt TestProfileProofReportsAThirdTimeoutAsTemporary TestProfileProofStopsWhenCancelledDuringTheBackoff TestProfileProofRetriesATimedOutSetupOnce TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce TestProfileProofDoesNotRetryARejectedSelection TestProfileProofDoesNotRetryACancelledCommand; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -q 'TestProfileProofReportsASecondTimeoutAsTemporary' docs/references/coverage-record.json internal/cli/profile_proof_retry_test.go; then printf 'the renamed test is still recorded or present\n' >&2; exit 1; fi; if grep -rq 'timed out twice' internal --include='*.go'; then printf 'old timeout text remains\n' >&2; exit 1; fi; GOOS=windows go build -buildvcs=false ./internal/cli || exit 1; GOOS=windows go vet ./internal/cli || exit 1; go test -count=1 ./internal/spec -run '^TestCoverageEquivalence$'` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the attempt, backoff and cancellation tests pass with the existing retry tests, the old name and text are gone, the package and its tests compile for Windows, and the coverage record matches the tests.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 2; API Contract 4; Invariants 3-4; Interfaces; Current behavior; Testing Approach; Build Order 2
- ADR-0298; ADR-0050; ADR-0114; ADR-0259; ADR-0271; ADR-0285; ADR-0250

## Result

Implemented the task_02 retry slice. Profile proof setup timeouts now receive
up to three attempts, with injected waits of 5 s and 10 s. The default sleeper
uses a timer and context cancellation without starting a goroutine. The final
temporary error reports the attempt count and effective setup deadline, and
the next action matches API Contract 2. Other errors and sleeper errors return
immediately.

Focused evidence:

- Red starting point: after adding the regression tests,
  `go test ./internal/cli -run '^TestProfileProofPassesOnTheThirdAttempt$' -count=1`
  exited 1 because the retry schedule and sleeper dependency did not exist.
- `go test ./internal/cli -run '^(TestProfileProof|TestOperationalPreflightPassesAfterAFallbackProof)' -count=1`
  exited 0 after implementation.
- After the final code/test edits,
  `go test -race ./internal/cli -run '^(TestProfileProof|TestOperationalPreflightPassesAfterAFallbackProof)' -count=1`
  exited 0.
- Acceptance criterion 1: the one-timeout and third-attempt success tests
  assert the recorded waits in order, successful JSON responses and empty
  stderr. Operational Preflight also asserts a silent fallback recovery.
- Acceptance criterion 2: `TestProfileProofReportsAThirdTimeoutAsTemporary`
  asserts three attempts, both waits, exact JSON error and next action,
  and the complete two-line text stderr, empty stdout and exit 2.
- Acceptance criterion 3:
  `TestProfileProofStopsWhenCancelledDuringTheBackoff` asserts one attempt,
  one wait and an unclassified context cancellation. The existing cancelled
  command test also exercises the default sleeper with a cancelled context
  and with an immediately expired timer, without a real wait.
- Acceptance criterion 4: the rejected-selection and cancelled-command tests
  assert one attempt and no waits. The retry returns non-timeout errors
  directly, preserving the existing error and classification.
- A read-only scan found no `timed out twice` in Go code or tests. Inspection
  of the generated coverage record confirmed the third-timeout name and both
  new test names, with the second-timeout name absent.
- `git diff --check` exited 0.

Coverage regeneration and follow-up:

- Ran the required
  `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`.
  The test body printed PASS and generated
  `docs/references/coverage-record.json`, but the process exited 1 because
  suiteguard reported `modified: docs/references/coverage-record.json`.
- Root cause inspection: `internal/spec/coverage_test.go` writes the record
  without calling `suiteguard.DeclareSanctionedRegeneration`; suiteguard
  requires that in-process command registration even though this Spec's
  authorization already declares the command and output. Fixing that
  generator/guard integration is outside this Task's declared slice.
- Retained the generated record without hand edits. Its additional entries
  reflect tests already present in the checkout, including task_01's tests.

Task status is unchanged by this Agent. Declared Verification, Windows
checks and repository Verification remain Daemon-owned and were not run.
Guide and skill wording remains task_04's slice.
