---
task: task_02
spec: 0150-a-reopen-that-cannot-be-raced
status: pending
type: backend
complexity: medium
---

# Task 02: Follow a symlinked Task path to its target

## Overview

`replaceTaskFile` creates its temporary file beside the Task path and renames
onto it. Where that path is a symlink, the rename replaces the link with a
regular file: the real Task keeps `status: completed` and the manifest now
points at a copy.

`SetStatus` used `os.WriteFile` before Spec 0149 and wrote through the link.
Confirmed with `git log -S "func replaceTaskFile"`, which shows the helper first
appearing in Spec 0149's own commit — this is a regression that Spec introduced
while making the write atomic, so it is repaired here.

## Requirements

1. MUST resolve the final Task path before choosing the temporary file's
   directory and the rename target, so the replacement lands on the file a
   symlink points at.
2. MUST leave a symlinked Task path a symlink after the write.
3. MUST keep the write atomic.
4. MUST NOT relax either confinement check Spec 0149 enforces: the lexical check
   against the Spec directory and the resolved check against the Spec Root.

## Subtasks

- [ ] Resolve the final path in the replacement helper.
- [ ] Add a test that reopens through a symlinked Task path.
- [ ] Confirm the confinement checks still refuse an escaping path.

## Acceptance Criteria

- [ ] A fixture whose Task path is a symlink to a file inside the Spec Root has
      the target rewritten to `pending` with the invalidation recorded.
- [ ] That path is still a symlink afterwards.
- [ ] The escaping-path refusals Spec 0149 shipped still refuse.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`

## Verification

- `out="$(go test -count=1 -run "^TestReopenGateWritesThroughASymlinkedTaskPath" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

- [_techspec.md](_techspec.md) — The symlinked path
