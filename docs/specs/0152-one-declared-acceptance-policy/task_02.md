---
task: task_02
spec: 0152-one-declared-acceptance-policy
status: pending
type: backend
complexity: low
---

# Task 02: Archive calls it

## Overview

Archive holds the policy the repository documented. This Task moves it behind
the one decision without changing what archive does.

## Requirements

1. MUST replace archive's inline verdict judgement with a call to the one
   decision.
2. MUST keep archive accepting exactly the reports it accepts today and refusing
   exactly the reports it refuses today, with the same reasons.
3. MUST stay inside the bounded path `internal/spec/archive.go`.
4. MUST NOT change any other archive precondition.

## Subtasks

- [ ] Call the one decision from archive.
- [ ] Pin archive's accepted and refused shapes with tests.

## Acceptance Criteria

- [ ] A `pass` archives, including one carrying an environment-blocked row.
- [ ] A qualifying `partial` archives, and every other shape refuses with
      today's reason.
- [ ] No archive precondition other than the verdict judgement is touched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestArchiveAppliesTheOneEligibilityDecision" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestArchiveAppliesTheOneEligibilityDecision"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation
