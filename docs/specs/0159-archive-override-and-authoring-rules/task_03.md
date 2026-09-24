---
task: task_03
spec: 0159-archive-override-and-authoring-rules
status: pending
type: docs
complexity: medium
---

# Task 03: One settlement table and complete authoring guidance

## Overview

The QA gate, archive and Roundfix skills each describe settlement in their own words, and the task-writing skill does not describe the declarations a Task can now make.

## Requirements

1. MUST add one `### QA settlement` section, identical in `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md` and `.agents/skills/roundfix/SKILL.md`, with a row each for pass, qualifying declared partial, environment-blocked, failed, missing and override, stating what each settles and what each archives.
2. MUST describe `roundfix archive --qa-override --approval --reason` in the archive and Roundfix skills and in `docs/user-guide/commands.md`.
3. MUST describe, in `.agents/skills/write-tasks/SKILL.md` and its task template, `verification: independent`, `precondition_repairs`, `SC-ORDINAL-CLAIMED`, temporal prerequisites as named prerequisites that never invent release authority, property-shaped acceptance, test seams, narrow Spec commits, and that a newly required test class must run under the repository gate.
4. MUST add `skills/settlement_guidance_repocontract_test.go` with `TestSettlementGuidanceIsOneTable`, failing when the three sections differ or a row is missing.
5. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The contract test passes and the mirror matches the canonical skills.
- [ ] The task-writing skill names every declaration.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `docs/user-guide/commands.md`
- creates: `skills/settlement_guidance_repocontract_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSettlementGuidanceIsOneTable$" ./skills 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSettlementGuidanceIsOneTable" && for s in qa-gate archive-spec roundfix write-tasks; do diff -r ".agents/skills/$s" "skills/$s" >/dev/null || exit 1; done && for w in "verification: independent" "precondition_repairs" "SC-ORDINAL-CLAIMED"; do grep -q -- "$w" .agents/skills/write-tasks/SKILL.md || exit 1; done && grep -q -- "--qa-override" .agents/skills/archive-spec/SKILL.md && grep -q -- "--qa-override" docs/user-guide/commands.md` — expected: exit 0; before this Task the contract test does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Guidance
