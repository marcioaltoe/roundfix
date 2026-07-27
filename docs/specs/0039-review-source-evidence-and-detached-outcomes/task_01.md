---
task: task_01
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: backend
complexity: high
---

# Task 01: Model Review Source Evidence and transient failures

## Overview

Introduce Review Source-neutral evidence and transient-failure contracts while
extending CodeRabbit signal parsing with the structured fields needed for
classification. This slice establishes typed, head-bound observations without
yet deciding the watch lifecycle outcome.

## Requirements

1. MUST model pending, reviewing, reviewed, verified, skipped, and failed
   Evidence tied to an expected and observed head.
2. MUST use stable Review Source-neutral evidence kinds and bounded details.
3. MUST preserve signal identity, conclusion, and explicit skip reason.
4. MUST parse CodeRabbit check output title and summary without treating
   arbitrary successful text as a skip.
5. MUST provide a typed transient failure preserving the wrapped cause and
   failed operation.
6. MUST positively classify temporary DNS, reset, HTTP 429, GitHub 5xx, and
   non-parent context timeouts as transient.
7. MUST keep cancellation, authentication, validation, and malformed responses
   permanent.

## Subtasks

- [x] Add Evidence state, kind, and head-bound fields.
- [x] Add the typed transient error and inspection helper.
- [x] Extend CodeRabbit check-run output parsing.
- [x] Preserve structured signal identity and bounded details.
- [x] Add transient and permanent failure matrices.
- [x] Add JSON mapping and sensitive-data regression coverage.

## Acceptance Criteria

- [x] Every Evidence value carries stable state, kind, identity, and relevant
      head without provider-specific response leakage.
- [x] Structured skip title or summary remains available for later
      classification.
- [x] Each approved temporary failure is discoverable through error wrapping.
- [x] Parent Run cancellation is never classified transient.
- [x] Authentication and invalid-request failures remain permanent.
- [x] Evidence and errors contain no credentials or unbounded response bodies.
- [x] Existing CodeRabbit fetch and resolution behavior remains passing.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/reviewsource/reviewsource.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/coderabbit/coderabbit_test.go`
- interface: `internal/watch/watch.go`

## Verification

- `rtk go test ./internal/reviewsource/... -run 'Test.*(Evidence|Transient|CheckRunOutput|SkipSignal)' -count=1`
  — expected: evidence shape, structured output parsing, and transient
  classification matrices pass.
- `rtk go test -race ./internal/reviewsource/... -run 'Test.*(Evidence|Transient)' -count=1`
  — expected: typed evidence and error handling are race-free.

## References

- `_prd.md` → Goals 1–2; User Stories 1–2 and 7; Core Features 1 and 5;
  Success Metrics.
- `_techspec.md` → Interfaces; Data Models: CodeRabbit check-run parsing;
  Testing Approach; Build Order 1.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  head-bound Evidence.

## Result

Implemented Review Source-neutral Evidence states and kinds with expected and
observed heads, stable signal identity, conclusion, bounded detail, and explicit
reason fields. CodeRabbit check-run parsing now retains the check identity plus
bounded output title and summary while discarding the unbounded output body.

Added `TransientError` and `IsTransient` with wrapped-cause preservation and a
bounded public error string. CodeRabbit access positively wraps temporary DNS,
connection reset, HTTP 429, GitHub 5xx, and non-parent deadline failures while
leaving parent cancellation, authentication, validation, and malformed
responses permanent. GitHub CLI failures no longer include command arguments or
response bodies in their rendered error.

### Verification

- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test ./internal/reviewsource/... -run 'Test.*(Evidence|Transient|CheckRunOutput|SkipSignal)' -count=1`
  — passed: 20 tests in 2 packages.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test -race ./internal/reviewsource/... -run 'Test.*(Evidence|Transient)' -count=1`
  — passed: 17 tests in 2 packages.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk go test ./internal/reviewsource/... -count=1`
  — passed: 52 tests in 2 packages, including existing fetch and resolution
  coverage.
- `rtk git diff --check` and `rtk gofmt -d internal/reviewsource/reviewsource.go internal/reviewsource/reviewsource_test.go internal/reviewsource/coderabbit/coderabbit.go internal/reviewsource/coderabbit/coderabbit_test.go`
  — passed with no output.

### Acceptance evidence

- `TestEvidenceJSONMappingAndBounds` covers all six stable states, neutral JSON
  field names, signal identity, both heads, conclusion, reason, bounded detail,
  and absence of provider-response fields.
- `TestCheckRunOutputJSONMapping`,
  `TestSkipSignalStructuredOutputRemainsAvailable`, and
  `TestSkipSignalDoesNotInferFromArbitrarySuccessfulText` prove structured
  title and summary preservation without retaining raw check output text or
  inventing a skip from generic success text.
- `TestTransientErrorPreservesOperationAndWrappedCause` and
  `TestTransientClassificationMatrix` prove typed discovery through wrapping,
  failed-operation preservation, and every approved temporary class.
- `TestTransientParentCancellationIsPermanent` and
  `TestTransientPermanentFailureMatrix` prove parent cancellation,
  authentication, invalid requests, and malformed responses remain permanent.
- The complete `internal/reviewsource/...` suite preserves existing CodeRabbit
  fetch, resolution, comment, status, and head-check behavior.

### Follow-ups

The shared CodeRabbit Evidence hierarchy, watch lifecycle outcomes, retry
episodes, and event publication remain assigned to later Tasks in this Spec and
were not changed here.
