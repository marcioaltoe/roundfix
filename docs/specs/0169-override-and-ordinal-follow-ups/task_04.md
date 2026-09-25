---
task: task_04
spec: 0169-override-and-ordinal-follow-ups
status: pending
type: backend
complexity: low
---

# Task 04: One instruction for overrides, one finding per conflict

## Overview

Corrective Task from the pre-PR review of 2026-09-25. The archive-spec skill's Precondition 2 still says to record the override in the stamped frontmatter (`qa_override: true`), contradicting the new Steps; the override field lists in the archive-spec, qa-gate and Roundfix skills omit `qa_override_qa_task_status`. And a same-Spec duplicate beside a fulfilled claim is reported twice, with a message and fix written for two different Specs.

## Requirements

1. MUST reword every archive-spec passage that describes hand-stamping, including Precondition 2 and the paragraph opening with `qa_override: true`, to say the override is performed only through `roundfix archive <slug> --qa-override --approval <source> --reason <text>`.
2. MUST add `qa_override_qa_task_status` to the override field lists in the archive-spec, qa-gate and Roundfix skills, and regenerate the mirrors with `make skills-sync`.
3. MUST skip the same-Spec pair when either path is already held on the tree, and give a same-Spec duplicate its own message and a fix that says to renumber one Task's `creates:` path.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] No skill instructs hand-stamping; every field list names the new field.
- [ ] A same-Spec duplicate beside a fulfilled claim yields exactly one finding, whose fix names renumbering.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/ordinal.go`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce|TestSameSpecDuplicateFixNamesRenumbering)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce TestSameSpecDuplicateFixNamesRenumbering; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && ! grep -q "record it in the stamped frontmatter" .agents/skills/archive-spec/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/archive-spec/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/qa-gate/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task the tests do not exist and Precondition 2 still hand-stamps, so the command fails.

## References

- [_techspec.md](_techspec.md) — The stamp and the guidance
