---
task: task_01
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
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

- [ ] Add Evidence state, kind, and head-bound fields.
- [ ] Add the typed transient error and inspection helper.
- [ ] Extend CodeRabbit check-run output parsing.
- [ ] Preserve structured signal identity and bounded details.
- [ ] Add transient and permanent failure matrices.
- [ ] Add JSON mapping and sensitive-data regression coverage.

## Acceptance Criteria

- [ ] Every Evidence value carries stable state, kind, identity, and relevant
      head without provider-specific response leakage.
- [ ] Structured skip title or summary remains available for later
      classification.
- [ ] Each approved temporary failure is discoverable through error wrapping.
- [ ] Parent Run cancellation is never classified transient.
- [ ] Authentication and invalid-request failures remain permanent.
- [ ] Evidence and errors contain no credentials or unbounded response bodies.
- [ ] Existing CodeRabbit fetch and resolution behavior remains passing.

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
