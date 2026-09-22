---
task: task_01
spec: 0153-a-reviewer-the-workflow-runs
status: pending
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
