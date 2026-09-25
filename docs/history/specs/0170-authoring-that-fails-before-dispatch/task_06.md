---
task: task_06
spec: 0170-authoring-that-fails-before-dispatch
status: completed
type: backend
complexity: low
---

# Task 06: Created operands reach the collision check; a bounded path the record omits is refused

## Overview

Corrective Task from the pre-PR review of 2026-09-25. `declaredTaskTouches` in `internal/spec/collision.go` calls `TaskVerificationFiles` with `Task{Verification: task.Verification}`, dropping the Task's Context, so a Verification operand the Task declares under `creates:` but that does not exist yet never reaches the wave-collision analysis. And `SC-TOOLING-UNDECLARED` checks Tooling rows only for paths a pending Task declares: a Governed Path listed in a PRD or TechSpec `bounded files:` row but absent from the authorization record's `paths:` passes when no Task declares it, although the TechSpec requires the rows and the record to agree.

## Requirements

1. MUST pass the whole Task (Verification and Context) from `declaredTaskTouches` to `TaskVerificationFiles`, so a non-existent operand declared under `creates:` is a collision touch.
2. MUST report `SC-TOOLING-UNDECLARED` for every Governed Path listed in a present Tooling row's `bounded files:` that the authorization record does not grant, whether or not a Task declares it; a sanctioned regeneration output stays exempt, and a path already reported for a Task is reported once.
3. MUST keep every existing finding, message and test of this Spec unchanged, and keep the active corpus count of `SC-TOOLING-UNDECLARED` at 0 in `internal/docscontract/testdata/corpus-golden.json` (if the active corpus now reports one, fix the authorization gap it names rather than the golden).

## Subtasks

- [x] Implement the requirements above.
- [x] Add a named test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [x] Two Tasks that may run together and name the same not-yet-existing `creates:` operand in their Verification are reported by the wave-collision check.
- [x] A Governed Path in a `bounded files:` row that the record omits is refused with no Task declaring it; the same path granted by the record passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/collision.go`
- interface: `internal/spec/collision_test.go`
- interface: `internal/speccheck/undeclared.go`
- interface: `internal/speccheck/undeclared_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCreatedVerificationOperandIsACollisionTouch|TestBoundedRowPathTheRecordOmitsIsRefusedWithoutATask|TestBoundedRowPathTheRecordGrantsPassesWithoutATask|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/spec ./internal/speccheck ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCreatedVerificationOperandIsACollisionTouch TestBoundedRowPathTheRecordOmitsIsRefusedWithoutATask TestBoundedRowPathTheRecordGrantsPassesWithoutATask TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the three new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Undeclared Governed Paths

## Result

Implementation:

- `declaredTaskTouches` now gives `TaskVerificationFiles` the complete Task, preserving `creates:` evidence for non-existent Verification operands.
- The undeclared-path detector now compares every Governed Path in each present Tooling authority row with the selected authorization record, exempts sanctioned regeneration outputs, emits row-only omissions deterministically, and avoids adding a row-only duplicate for a path already reported from a Task.
- The shared clean Spec fixture now records the `Makefile` path and `implement` operation its existing prose and executable Task Graph already claimed. No existing assertion or finding message changed, and `internal/docscontract/testdata/corpus-golden.json` remains unchanged at `SC-TOOLING-UNDECLARED: 0`.

Focused checks:

- Before the implementation, `go test -count=1 -run '^(TestCreatedVerificationOperandIsACollisionTouch|TestBoundedRowPathTheRecordOmitsIsRefusedWithoutATask|TestBoundedRowPathTheRecordGrantsPassesWithoutATask|TestBoundedRowSanctionedRegenerationOutputPassesWithoutATask)$' ./internal/spec ./internal/speccheck` failed in the created-operand and omitted-row cases, while both acceptance companions passed.
- After the implementation, the same focused command passed both packages.
- `go test -count=1 -run '^TestSpecCheckRunVerification$' ./internal/cli` passed after the clean fixture grant was made truthful.
- `go test -count=1 ./internal/spec ./internal/speccheck` passed.
- `make verify-incremental` passed with the Task-scoped Go build cache.
- The Task's declared `## Verification` command was not run; Daemon Verification remains pending.

Acceptance evidence:

- Created Verification operand collision: `TestCreatedVerificationOperandIsACollisionTouch` passes and observes one Wave collision on `internal/generated/client.go` with `TouchFromVerification`.
- Bounded-row record agreement: `TestBoundedRowPathTheRecordOmitsIsRefusedWithoutATask` passes with one `SC-TOOLING-UNDECLARED`; `TestBoundedRowPathTheRecordGrantsPassesWithoutATask` passes with none. `TestBoundedRowSanctionedRegenerationOutputPassesWithoutATask` separately preserves the regeneration exemption.

Follow-ups: none. The change introduces no new domain term, so `CONTEXT.md` needs no update.
