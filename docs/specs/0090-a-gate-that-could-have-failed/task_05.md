---
task: task_05
spec: 0090-a-gate-that-could-have-failed
status: pending
type: backend
complexity: medium
---

# Task 05: Carry the negative control a Task declares

## Overview

The probe delivers the observability control mechanically, and the existing
post-Agent Verification is the positive control. The negative control cannot be
manufactured — Roundfix will not mutate a maintainer's tree to fabricate a defect
— so it is authored. This Task parses a Task's `## Negative Control` section,
carries it on the Task, and records whether one was declared, so gate health
becomes a recorded fact rather than an assumption.

## Requirements

1. MUST parse a `## Negative Control` section from a Task file into a list of
   declarations, mirroring how `## Verification` is parsed.
2. MUST leave a Task with no such section valid, carrying an empty list; a Task
   without a declared negative control is a weaker gate stated honestly, per
   ADR-0110.
3. MUST record, per Task, whether a negative control was declared, so the fact
   travels with the Task's outcome.
4. MUST NOT execute or synthesise the declared control; this Task carries it and
   nothing more.
5. MUST NOT change how `## Verification` is parsed.

## Subtasks

- [ ] Parse the section.
- [ ] Carry it on the Task.
- [ ] Record its presence with the Task's outcome.

## Acceptance Criteria

- [ ] A Task declaring negative controls parses with them in order.
- [ ] A Task declaring none parses successfully with an empty list.
- [ ] The declaration count is recorded with the Task's outcome.
- [ ] Parsing `## Verification` is byte-identical in behaviour to before.

## Bounded scope

This Task may create or modify only:

- `internal/spec/spec.go`
- `internal/spec/task.go`
- `internal/spec/spec_test.go`
- `internal/spec/task_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_05.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestNegativeControl' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestNegativeControlSectionParsesInOrder'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestNegativeControl' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestNegativeControlAbsentSectionParsesEmpty'` — expected: exits 0.
- `grep -q 'NegativeControl' internal/spec/task.go` — expected: exits 0. This string does not exist in the file before this Task.

## References

- `_prd.md` → Core Features, three controls.
- `_techspec.md` → Build Order 5; Data Models.
- ADR-0110.
