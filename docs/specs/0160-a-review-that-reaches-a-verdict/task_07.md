---
task: task_07
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

## Result

- The review command now discovers changed Spec folders under both the configured Spec Root and its resolved archive root, reads usable PRD Decisions and TechSpec content from `HEAD`, records unusable folders in `skippedSpecs`, and continues to the reviewer.
- Spec context is capped at 32 KiB per Spec and 64 KiB in total. Truncated content carries a visible `[Spec context truncated]` marker, and the review record sets `specContextTruncated`.
- Genuine candidate Spec read errors remain blocking, but their record reason is `Spec context read failure` rather than a reviewer runtime failure.
- Acceptance criterion 1: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReviewSkipsASpecWithoutDecisions$' ./internal/cli` passed. The test covers a PRD without `## Decisions` and a second changed Spec without `_techspec.md`; both slugs are recorded as skipped, no Spec block is appended, and the reviewer still returns a reviewed outcome.
- Acceptance criterion 2: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReviewReadsAnArchivedSpecFromTheCandidate$' ./internal/cli` passed. The test adds a Spec under `docs/history/specs`, changes the worktree copy after committing, and proves the prompt carries the archived files from `HEAD` rather than the worktree.
- Acceptance criterion 3: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReviewBoundsSpecContext$' ./internal/cli` passed. The test observes the prompt marker and record flag for an oversized Spec, asserts the per-Spec envelope, and checks a four-context input against the total cap.
- Requirement 4 focused evidence: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReviewReportsSpecReadingFailureDistinctly$' ./internal/cli` passed; the reviewer is not called and neither the record nor stderr labels the local read error as a runtime failure.
- Post-adjustment focused sweep: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReview(SkipsASpecWithoutDecisions|ReadsAnArchivedSpecFromTheCandidate|BoundsSpecContext|ReportsSpecReadingFailureDistinctly|PromptCarriesSpecDecisions)$' ./internal/cli` passed.
- Regression checks: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 -run '^TestReview' ./internal/cli` passed with the existing GitHub integration boundary allowed after the sandboxed attempt was blocked at `api.github.com`. After the final context-bound adjustment, `GOCACHE=/tmp/roundfix-task07-gocache make verify-incremental` passed with the same network allowance; vet, all packages, the focused skill-policy tests, the Roundfix skill check, and the build succeeded.
- The Task's declared Verification command was not run; Daemon Verification remains the settlement authority.

## Carry-forward provenance

- Source Run: `run_20260924T193553Z_7f2a4b9ff0f723fb`
- Source commit: `6df23f2edf460ebbe2d60fdecf7f1f60d2e1e3e5`
