---
task: task_05
spec: 0090-a-gate-that-could-have-failed
status: completed
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

## Result

Implemented the authored negative-control carrier without adding an execution
path. Task parsing now reads backticked bullet declarations from every
`## Negative Control` section in source order, stores them on `spec.Task`, and
refreshes them when the Daemon reloads a Task after the Agent turn. An absent
section yields zero declarations. The declaration count remains available as
`len(task.NegativeControl)` alongside the Task status used for its outcome.

Focused-check evidence:

- Red baseline: `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/spec -run 'TestNegativeControl(Section|Absent|Declaration|Parsing)'` failed to compile because `Task` and `taskDocument` had no `NegativeControl` field.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/spec -run 'Test(LoadParsesTaskFiles|ReloadTaskPicksUpAgentEdits|NegativeControl(SectionParsesInOrder|AbsentSectionParsesEmpty|DeclarationCountTravelsWithTaskOutcome|ParsingPreservesVerificationCommands))$'` passed.
- `rtk proxy cmp <(rtk proxy git show HEAD:internal/spec/task.go | rtk proxy sed -n '/^func parseVerificationCommands/,/^}/p') <(rtk proxy sed -n '/^func parseVerificationCommands/,/^}/p' internal/spec/task.go)` exited 0, proving the existing Verification parser block is byte-identical to `HEAD`.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

- Declared controls preserve order: `TestNegativeControlSectionParsesInOrder` passed with two distinct declarations around a skipped non-command bullet.
- No declaration remains valid: `TestNegativeControlAbsentSectionParsesEmpty` passed with a zero-length declaration list.
- The declaration count travels with the outcome: `TestNegativeControlDeclarationCountTravelsWithTaskOutcome` passed after reloading a Task whose status changed to `completed`, retaining both declarations.
- Verification parsing is preserved: the exact `HEAD` parser-block comparison exited 0; `TestNegativeControlParsingPreservesVerificationCommands`, `TestLoadParsesTaskFiles`, and `TestReloadTaskPicksUpAgentEdits` also passed.

The commands under `## Verification` were not run; the Daemon owns them.
