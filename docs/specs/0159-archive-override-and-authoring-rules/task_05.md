---
task: task_05
spec: 0159-archive-override-and-authoring-rules
status: pending
type: backend
complexity: medium
---

# Task 05: The override waives what normal archive would refuse

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The override refuses whenever the newest QA Report qualifies, without looking at the QA Task's status, so a Spec whose QA Task is failed or pending while a stale report says `pass` can be archived neither normally nor with the override, contradicting ADR-0154. And an unreadable report stores an absolute local path in `qa_override_qa_outcome`.

## Requirements

1. MUST refuse the override only when a normal archive of the same Spec would succeed: every Task completed and the newest report eligible.
2. MUST accept the override for a failed or pending QA Task whatever the newest report says, still requiring every non-QA Task completed.
3. MUST record an unreadable report's outcome with a path relative to the Spec folder, never an absolute path.
4. MUST keep the archive-spec skill, its mirror and `docs/user-guide/commands.md` true to these rules, regenerating the mirror with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A failed QA Task with a `pass` report archives under the override.
- [ ] The override is refused when a normal archive would succeed.
- [ ] The recorded outcome of an unreadable report holds no absolute path.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `.agents/skills/archive-spec/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveQAOverrideAcceptsAFailedQATaskWithAPassReport|TestArchiveQAOverrideRefusedOnlyWhenNormalArchiveSucceeds|TestArchiveQAOverrideRecordsARelativeOutcome)$" ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveQAOverrideAcceptsAFailedQATaskWithAPassReport TestArchiveQAOverrideRefusedOnlyWhenNormalArchiveSucceeds TestArchiveQAOverrideRecordsARelativeOutcome; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The override
