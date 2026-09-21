---
task: task_01
spec: 0152-one-declared-acceptance-policy
status: pending
type: backend
complexity: medium
---

# Task 01: The one eligibility decision

## Overview

Whether the newest QA Report is acceptable is decided in two places that
disagree. This Task creates the one decision both will use. It states exactly
the rule `internal/spec/archive.go` states today, so nothing that is acceptable
now stops being acceptable.

## Requirements

1. MUST provide one exported decision that takes a parsed QA Report and answers
   whether it is acceptable.
2. MUST return the reason when it is not, so a caller can report what archive
   reports today.
3. MUST accept `pass` with no blocked rows, and `partial` whose blocked rows are
   declared unreachable.
4. MUST refuse `fail`, a `partial` whose blocked rows are not declared, and a
   `pass` carrying blocked rows.
5. MUST NOT change which report is newest, or what makes a blocked row
   unreachable.

## Subtasks

- [ ] Add the decision and its reasons.
- [ ] Add unit tests for every accepted and refused shape.

## Acceptance Criteria

- [ ] Each accepted shape is accepted and each refused shape is refused.
- [ ] Each refusal names its reason.
- [ ] Newest-report selection is untouched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestQAReportEligibility" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestQAReportEligibility"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The one decision
