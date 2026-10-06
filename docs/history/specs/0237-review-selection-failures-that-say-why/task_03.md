---
task: task_03
spec: 0237-review-selection-failures-that-say-why
status: completed
type: backend
complexity: medium
---

# Task 03: The review names the failing step and retries a failure before the prompt once

## Overview

`roundfix review` puts the step and the adapter's message that task_02 keeps
into its record, and retries a failure during Agent Selection once on the same
selection when the review prompt was not yet sent, before any configured
fallback. A failure after the prompt is never retried and no failure selects
`none`. It answers the Backlog Entry "A pre-PR review agent selection failure
says nothing actionable" of 2026-10-06. Verifiable on its own through the
command entry point with the fake review runner.

## Requirements

1. MUST add `reviewSelectionRetry` and the record fields `failedStep`,
   `adapterMessage` and `selectionRetries` of `_techspec.md` → Interfaces,
   each omitted when empty, keeping older records readable.
2. MUST, in `runConfiguredReviewSession`, retry per Invariants 7 and 8: at
   most once per selection; a `PrepareSession` error that
   `reviewSelectionCanFallback` accepts is retried after ending that session
   by preparing the same selection again; a `RunPrepared` error that
   `agent.DescribeProtocolFailure` places with `FromPrompt` true and
   `PromptSent` false is retried by calling `RunPrepared` again on the
   prepared session, both in the fresh-selection loop and on round two's
   resumed session. A failed preparation of round two's resumed session keeps
   today's path to a fresh session. Context cancellation and deadline are
   never retried.
3. MUST write the stderr notice of API Contract 4 before each retry and
   append `{selection, step, message}` to `selectionRetries`, using
   `session setup` and an empty message when the failure is not placed.
4. MUST make a configured fallback activate only after the selection's retry
   also failed before the prompt, with today's notice, and keep a failure
   after the prompt, or one not placed, as a review failure with no retry and
   no fallback.
5. MUST, in `classifyReviewCommandResult`, set `failedStep` and, when
   non-empty, `adapterMessage` for a runtime failure that
   `agent.DescribeProtocolFailure` places, keep the reason
   `review runtime failure: <err>`, and append
   ` (after one automatic retry before the prompt)` when the failing selection
   was retried, per API Contract 3.
6. MUST add the tests named in Verification to a new
   `internal/cli/review_selection_retry_test.go` through
   `newReviewCommandFixture` and the fake `reviewCommandRunner`: a
   preparation selection failure then success exits `0` with two preparations
   of the same selection, one retry recorded and no fallback notice; a
   prompt-process failure with `PromptSent` false then success exits `0` with
   two prompt calls and one retry; a failure with `PromptSent` true blocks
   after one prompt call with `failedStep` `session/prompt`, the adapter's
   message and no retry; two failures before the prompt activate a configured
   fallback, or block with the retry suffix when none is configured; a
   generic preparation error blocks without a retry; and round two's resumed
   session retries a prompt-process failure before `session/prompt`. No case
   records `none` or a reviewed outcome for a failure.
7. MUST update exactly these existing expectations, which pin a fallback after
   one failure, to give the preferred selection two failures and expect one
   more preparation, changing nothing else in them:
   `TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt`,
   `TestReviewRoundTwoContinuesTheRecordedFallbackSelection` and
   `TestReviewValidatorUsesTheSuccessfulFallbackSelection`.
8. MUST NOT change the provider policy, the review profiles, the diff and
   token bounds, the verdict classification, the lineage rules or the daemon
   packages.

## Subtasks

- [ ] Add the record fields and the retry record.
- [ ] Retry once per selection before the prompt in both session paths.
- [ ] Keep fallback after the retry and no retry after the prompt.
- [ ] Fill the step, the message and the retry suffix when classifying.
- [ ] Add the retry tests and update the three fallback tests.

## Acceptance Criteria

- [ ] A transient failure before the prompt is retried once on the same
      selection and the review completes with the retry recorded.
- [ ] A failure after the prompt is never retried, blocks, and names
      `session/prompt` and the adapter's message.
- [ ] A fallback activates only after the retry failed; no failure records
      `none` or a pass.

## Context

- creates: `internal/cli/review_selection_retry_test.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_test.go`
- interface: `internal/cli/review_session_test.go`
- interface: `internal/cli/review_convention_validator_test.go`
- instruction: `docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReviewRetriesASelectionFailureBeforeThePromptOnce|TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt|TestReviewNeverRetriesAfterSessionPromptWasSent|TestReviewActivatesTheFallbackAfterTheRetryFails|TestReviewBlocksAfterTheRetryFailsWithoutFallback|TestReviewDoesNotRetryAGenericPreparationError|TestReviewRoundTwoRetriesAResumedPromptBeforeSessionPrompt|TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt|TestReviewRoundTwoContinuesTheRecordedFallbackSelection|TestReviewValidatorUsesTheSuccessfulFallbackSelection)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewRetriesASelectionFailureBeforeThePromptOnce TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt TestReviewNeverRetriesAfterSessionPromptWasSent TestReviewActivatesTheFallbackAfterTheRetryFails TestReviewBlocksAfterTheRetryFailsWithoutFallback TestReviewDoesNotRetryAGenericPreparationError TestReviewRoundTwoRetriesAResumedPromptBeforeSessionPrompt TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt TestReviewRoundTwoContinuesTheRecordedFallbackSelection TestReviewValidatorUsesTheSuccessfulFallbackSelection; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven retry tests do not exist, so the command fails; after it the three updated fallback tests pass with them.

## References

- `_prd.md` → Core Feature 2; Core Feature 3; User Stories 1-3; Success Metric 3; Success Metric 4; Goals
- `_techspec.md` → Interfaces; Invariants 7 and 8; Data Models; API Contract 3; API Contract 4; Testing Approach 2; Testing Approach 3; Build Order 3; Vocabulary Contract
- ADR-0242
- ADR-0050
- ADR-0153
- ADR-0197

## Result

Implemented the Task's review-command slice. The record now omits empty
`failedStep`, `adapterMessage` and `selectionRetries` fields. Preparation
failures eligible for selection fallback close their session and retry the
same selection once. Placed prompt-process failures before `session/prompt`
retry on the prepared session, including round two's resumed session. Both
paths share one retry budget per selection, emit the specified stderr notice,
and retain the retry in the persisted record.

A configured fallback activates only after that selection's retry also fails
before the prompt. Failures after the prompt, unplaced prompt failures,
generic preparation errors, cancellation and deadlines do not retry or
activate fallback. Classification preserves the runtime-failure reason, adds
the placed step and adapter message, and adds the retry suffix only when the
failing selection was retried. Existing provider, verdict, bounds and lineage
policies remain in place; no daemon package changed.

### Acceptance evidence

| Criterion | Implementation and focused-check evidence |
| --- | --- |
| A transient failure before the prompt retries once on the same selection and completes with the retry recorded. | `TestReviewRetriesASelectionFailureBeforeThePromptOnce` checks two preparations of the same runtime/session with cleanup between them; `TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt` checks two prompt calls on one prepared session. Both check exit 0, the persisted retry and no fallback notice. `TestReviewRoundTwoRetriesAResumedPromptBeforeSessionPrompt` checks the same behavior on the continued session. All passed in the focused review suite. |
| A failure after the prompt never retries, blocks, and names the step and message. | `TestReviewNeverRetriesAfterSessionPromptWasSent` checks one prompt call, blocked outcome, `session/prompt`, the adapter message and no retry/fallback. `TestReviewFailureAfterRetryDoesNotActivateFallbackAfterPrompt` covers an after-prompt failure on the retry. Both passed. |
| Fallback activates only after the retry fails; failures never record omission or a reviewed outcome. | `TestReviewActivatesTheFallbackAfterTheRetryFails` covers preparation and prompt-process failures followed by a successful fallback. `TestReviewBlocksAfterTheRetryFailsWithoutFallback` checks that the final available selection blocks after two failures and records the retry suffix. Generic-error, cancellation, deadline and unplaced-prompt cases block without retry/fallback. Additional checks cover the shared preparation/prompt budget and avoid attributing an earlier selection's retry to an unretried fallback failure. All passed. |

Added the seven requested named tests and four additional boundary tests in
`internal/cli/review_selection_retry_test.go`. Updated only the three existing
fallback expectations named in Requirement 7 to supply two preferred-selection
failures and account for the extra preparation/session cleanup.

### Checks run

- Before-change signal: inspection found immediate fallback after a single
  preparation failure and an unconditional return after one `RunPrepared`
  call. A temporary Go overlay restored that original session logic, adding
  only the new record schema so the regression test could compile:
  `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -overlay /private/tmp/roundfix-task03-before-overlay.json -run '^TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt$' ./internal/cli`
  exited 1 with a blocked record and one prompt call, as expected.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -run 'TestReview(Retries|NeverRetries|ActivatesTheFallback|BlocksAfterTheRetry|DoesNotRetry|RoundTwoRetries|SelectionRetry)' ./internal/cli`:
  exit 0 after correcting the fixture's invalid empty fallback configuration.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -run '^TestReview' ./internal/cli`:
  exit 0, 116.058 seconds, including the new tests and all three adjusted
  fallback tests.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk make verify-incremental`:
  the initial sandboxed run exited 2 because force-stop integration tests
  could not read the process table and because this Agent edited the new
  test file while repository guards were running. The rerun with process
  access enabled and files held unchanged exited 0, including formatting,
  vet, tests, skill checks and build.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.

### Scope note

Requirement 6's literal case of a review with no configured fallback cannot
reach the public command with today's valid configuration:
`internal/config/profiles.go` requires at least one fallback in an Agent
Selection Profile. Changing that policy is explicitly outside Requirement 8.
The requested `TestReviewBlocksAfterTheRetryFailsWithoutFallback` therefore
exercises the final selection with no fallback remaining, through the public
command and the supplied fixture, and checks the same terminal retry/block
behavior. A follow-up can clarify the Spec wording; no profile policy changed.

Task status, Subtask and Acceptance Criteria settlement, and the declared
Verification remain Daemon-owned. The declared Verification command was not
run. No commit, push or Pull Request was made.
