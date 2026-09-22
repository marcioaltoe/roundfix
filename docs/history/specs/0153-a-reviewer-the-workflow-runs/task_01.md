---
task: task_01
spec: 0153-a-reviewer-the-workflow-runs
status: completed
type: backend
complexity: medium
---

# Task 01: The review record and its writer

## Overview

A pre-PR review today leaves no trace the workflow can read. This Task adds the
record that binds a review to the candidate it examined.

## Requirements

1. MUST carry the repository, the base commit and the head commit.
2. MUST carry the effective provider and the `Source` that selected it, read
   from the resolved Pre-PR Review Policy.
3. MUST carry an execution outcome that is exactly one of: reviewed, findings,
   blocked, omitted.
4. MUST carry the findings when the outcome is findings, and the reason when the
   outcome is blocked.
5. MUST make a record that names one head say nothing about another: the head is
   part of the record, not implied by where it is stored.
6. MUST NOT change how the Pre-PR Review Policy is resolved or defaulted.

## Subtasks

- [ ] Add the record type and its writer.
- [ ] Add unit tests for each outcome and its required fields.

## Acceptance Criteria

- [ ] Each of the four outcomes round-trips with its required fields.
- [ ] A record without a head is refused rather than written.
- [ ] Policy resolution is untouched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReviewRecord" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewRecord"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — What the record carries

## Result

Implemented a JSON review record that names the repository, base commit, head
commit, effective provider, policy source, and one of the four specified
outcomes. The constructor copies provider provenance from the resolved
`config.PrePRReview` value. The writer validates the complete record before it
emits bytes; findings and blocked outcomes require their corresponding detail.

Focused checks:

- `GOCACHE=/private/tmp/roundfix-task01-go-cache go test -count=1 -run '^TestReviewRecordRoundTripsEachOutcome$' ./internal/cli` — passed; reviewed, findings, blocked, and omitted records encoded and decoded with their required fields.
- `GOCACHE=/private/tmp/roundfix-task01-go-cache go test -count=1 -run '^TestReviewRecordWithoutHeadIsRefusedBeforeWriting$' ./internal/cli` — passed; the writer returned the missing-head error and wrote zero bytes.
- `GOCACHE=/private/tmp/roundfix-task01-go-cache go test -count=1 -run '^TestReviewRecordRefusesInvalidOutcomeDetails$' ./internal/cli` — passed; unknown outcomes and missing findings/reasons were refused before writing.
- `gofmt -d internal/cli/review.go internal/cli/review_test.go` — no output; both new Go files are formatted.
- `git diff --check -- docs/specs/0153-a-reviewer-the-workflow-runs/task_01.md` — passed.
- `git diff --name-only -- internal/config/config.go` — no output; policy resolution and defaulting are untouched.

Acceptance evidence:

- Each outcome round-trips with required fields: `TestReviewRecordRoundTripsEachOutcome` covers all four outcome constants, policy provenance, candidate commits, findings, and blocked reason.
- A record without a head is refused: `TestReviewRecordWithoutHeadIsRefusedBeforeWriting` also proves validation happens before any output.
- Policy resolution is untouched: the implementation only reads the resolved `config.PrePRReview` value and the config source has no diff.

The declared Verification command was not run; the Daemon owns that check.
