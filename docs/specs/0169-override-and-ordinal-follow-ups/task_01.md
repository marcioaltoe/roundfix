---
task: task_01
spec: 0169-override-and-ordinal-follow-ups
status: pending
type: backend
complexity: medium
---

# Task 01: The stamp names the waived Task; guidance matches the command

## Overview

An override of a failed or pending QA Task whose newest report says `pass` stamps `qa_override_qa_outcome: pass` and records nowhere that the QA Task was not completed; the Roundfix skill still says the override "refuses when QA already qualifies"; and the archive-spec skill's Steps still describe stamping `qa_override: true` by hand.

## Requirements

1. MUST stamp `qa_override_qa_task_status` with the QA Task's status when it is not completed, and omit it when it is completed.
2. MUST state in `.agents/skills/roundfix/SKILL.md` that the override is refused only when a normal archive would succeed, replacing the old wording.
3. MUST state in `.agents/skills/archive-spec/SKILL.md` Steps that an override is performed only through `roundfix archive --qa-override`, never by hand-stamping.
4. MUST name `qa_override_qa_task_status` in `docs/user-guide/commands.md`, and regenerate the mirror with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A failed QA Task with a `pass` report records its status; a completed QA Task records none.
- [ ] Both skills and the guide carry the current rule and the command as the only path.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveQAOverrideRecordsTheQATaskStatus|TestArchiveQAOverrideOmitsTheStatusForACompletedQATask)$" ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveQAOverrideRecordsTheQATaskStatus TestArchiveQAOverrideOmitsTheStatusForACompletedQATask; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "qa_override_qa_task_status" docs/user-guide/commands.md && grep -q "normal archive would succeed" .agents/skills/roundfix/SKILL.md && grep -q "normal archive would succeed" skills/roundfix/SKILL.md && ! grep -q "refuses when QA already qualifies" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The stamp and the guidance
