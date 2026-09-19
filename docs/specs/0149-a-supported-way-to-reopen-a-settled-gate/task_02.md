---
task: task_02
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: pending
type: backend
complexity: medium
---

# Task 02: The reopen command and its refusals

## Overview

The operation itself. `roundfix reopen --spec <slug>` returns a Spec's terminal
QA Task to `pending` when the loader's stale-gate condition already holds, and
refuses otherwise.

Its refusals carry the weight: a reopen that accepts a healthy gate is not a
recovery path but a supported way to discard a verdict.

## Requirements

1. MUST exit 0 on a stale gate, writing the QA Task's status `pending` through
   `spec.SetStatus`.
2. MUST append an invalidation record to the QA Task naming the QA Report it
   invalidated, the dependency ids that made the verdict stale, and the date,
   leaving the prior `## Result` in place above it.
3. MUST refuse in preflight with exit 2 when the Spec has no terminal QA Task,
   when the QA Task is not completed, or when every dependency is completed,
   and MUST name which condition failed.
4. MUST leave the Spec directory byte-identical when it refuses.
5. MUST NOT open the QA Report for writing.
6. MUST NOT create a Run, write a Run Event Journal entry, commit or push.
7. MUST refuse unknown flags, matching the surrounding commands.

## Subtasks

- [ ] Add the command, its preflight refusals and its invalidation record.
- [ ] Register it on the public command surface.
- [ ] Add tests at the command seam for the stale case and both refusals.

## Acceptance Criteria

- [ ] A fixture whose QA Task is completed above a pending dependency exits 0,
      and `spec.Load` on that fixture afterwards returns no error.
- [ ] A fixture whose QA Task is pending refuses, naming that condition.
- [ ] A fixture whose dependencies are all completed refuses, and a digest of
      the Spec directory is unchanged across the call.
- [ ] After a successful reopen the QA Report's bytes are unchanged and the QA
      Task retains its prior Result text.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/settle.go`
- interface: `internal/spec/task.go`

## Verification

- `out="$(go test -count=1 -run "^TestReopen" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix reopen --spec"` — expected: exit 0; the command appears on the public surface. Before this Task the usage has no reopen line, so the command fails.

## References

- [_techspec.md](_techspec.md) — The command
