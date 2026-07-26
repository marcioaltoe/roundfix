---
task: task_02
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: backend
complexity: high
---

# Task 02: Unify CodeRabbit evidence classification

## Overview

Replace separate pre-fetch and Merge-Ready interpretations with one CodeRabbit
evidence hierarchy. Both watch phases receive the same result for the same
head, including explicit skip, current-head approval, unresolved-thread, stale
signal, and no-signal cases.

## Requirements

1. MUST expose one evidence operation for an Open Pull Request and expected
   head.
2. MUST apply the TechSpec evidence hierarchy in stable precedence order.
3. MUST classify explicit structured skip before successful or pending signals.
4. MUST require the expected head for every accepted check, status, or review.
5. MUST accept a current-head CodeRabbit `APPROVED` review as verified only when
   zero unresolved CodeRabbit threads remain.
6. MUST classify other current-head reviews as reviewed, never verified
   approval.
7. MUST publish changed observations once and suppress duplicate unchanged
   polling events.

## Subtasks

- [ ] Implement the shared evidence operation.
- [ ] Unify check-run, commit-status, review, and thread precedence.
- [ ] Route pre-fetch and Merge-Ready checks through the classifier.
- [ ] Enforce expected-head and approval requirements.
- [ ] Publish bounded changed-observation events.
- [ ] Add stale-head, unresolved-thread, and precedence tables.

## Acceptance Criteria

- [ ] Pre-fetch and Merge-Ready return byte-equivalent Evidence for the same
      fixture.
- [ ] Explicit skip wins over success or pending signals and preserves reason.
- [ ] Current-head approval with zero unresolved threads produces verified
      `review_approval` Evidence.
- [ ] Commented, stale-head, or unresolved-thread reviews cannot prove
      Merge-Ready.
- [ ] No usable signal produces pending without guessing.
- [ ] One changed observation creates one event; an unchanged poll creates
      none.
- [ ] Existing Review Source authentication and repository boundaries remain
      unchanged.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/reviewsource/reviewsource.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/coderabbit/coderabbit_test.go`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/reviewsource/... -run 'Test.*(EvidenceHierarchy|Approval|Skipped|ExpectedHead|UnresolvedThread|Precedence)' -count=1`
  — expected: one classifier produces every accepted and refused state.
- `rtk go test ./internal/watch ./internal/runevent -run 'Test.*(ReviewEvidence|ReviewStatusEvent|UnchangedEvidence)' -count=1`
  — expected: both watch phases share Evidence and events deduplicate unchanged
  polls.

## References

- `_prd.md` → Goal 1; User Stories 1 and 7; Core Features 1–3; Success Metrics.
- `_techspec.md` → API Contracts: Review Source evidence hierarchy; Integration
  Points: GitHub through `gh`; Build Order 2.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  accepted Review Source Evidence.
