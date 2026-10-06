---
task: task_03
spec: 0237-review-selection-failures-that-say-why
status: pending
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
