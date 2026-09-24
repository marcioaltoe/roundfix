---
task: task_07
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: medium
---

# Task 07: Spec context that never blocks and stays bounded

## Overview

Corrective Task from the pre-PR review of 2026-09-24. A candidate touching a Spec whose PRD has no `## Decisions` section, or that lacks one of `_prd.md` and `_techspec.md`, is blocked as a "runtime failure" before any reviewer is called — true of several PRDs already on `main`. Spec folders are looked up only under the Spec Root, so a Spec archived within the candidate is missed. Spec context is appended without any bound.

## Requirements

1. MUST skip, without blocking, a changed Spec that has no `## Decisions` section or lacks a PRD or TechSpec, and record its slug as skipped.
2. MUST also collect Spec folders the candidate adds or changes under the resolved archive root, reading their files at `HEAD`.
3. MUST cap the carried context per Spec and in total, truncating with a visible marker and recording that truncation in the record.
4. MUST keep reporting a local Spec-reading error, if any remains, with its own reason rather than as a runtime failure.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A Spec without Decisions is skipped and the review proceeds.
- [ ] A Spec archived within the candidate is carried.
- [ ] Oversized context is truncated and the truncation recorded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewSkipsASpecWithoutDecisions|TestReviewReadsAnArchivedSpecFromTheCandidate|TestReviewBoundsSpecContext|TestReviewPromptCarriesSpecDecisions)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewSkipsASpecWithoutDecisions TestReviewReadsAnArchivedSpecFromTheCandidate TestReviewBoundsSpecContext TestReviewPromptCarriesSpecDecisions; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task three of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Spec-aware prompt
