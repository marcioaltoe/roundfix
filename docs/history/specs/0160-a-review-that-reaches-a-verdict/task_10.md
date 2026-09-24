---
task: task_10
spec: 0160-a-review-that-reaches-a-verdict
status: completed
type: docs
complexity: low
---

# Task 10: Describe the whole-answer verdict contract

## Overview

QA finding F-001 of 2026-09-24 (second gate): Task 09 made a no-findings verdict account for the whole answer and read terminal `Findings: none`, `Findings: n/a` and `Findings: no findings` as no findings, but the shipped Roundfix skill, its mirror and `docs/user-guide/commands.md` still say one substantive verdict is sufficient.

## Requirements

1. MUST state in `.agents/skills/roundfix/SKILL.md` and `docs/user-guide/commands.md` that a no-findings verdict passes only when it accounts for the whole answer: it is the only content line, or the last one after a preamble with no list item, heading or `path:line` reference; anything else beside it blocks.
2. MUST state that `Findings: none`, `Findings: n/a` and `Findings: no findings` with nothing after them are read as no findings.
3. MUST state that verdict-shaped lines after a `Findings:` header do not create a conflict, and that Specs dropped by the total context bound are named in the record.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The skill, its mirror and the guide describe the whole-answer rule and the findings-none forms.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -q "whole answer" .agents/skills/roundfix/SKILL.md && grep -q "whole answer" skills/roundfix/SKILL.md && grep -q "whole answer" docs/user-guide/commands.md && grep -q "Findings: n/a" .agents/skills/roundfix/SKILL.md && grep -q "Findings: n/a" docs/user-guide/commands.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task no surface names the whole-answer rule, so the command fails.

## References

- [_techspec.md](_techspec.md) — Classification

## Result

- Updated the canonical Roundfix skill and user guide to state that a no-findings verdict passes only when it accounts for the whole answer: it is the only content line, or the last line after an unstructured preamble; list items, headings, `path:line` references, or any other content beside it block.
- Documented `Findings: none`, `Findings: n/a`, and `Findings: no findings` with nothing after them as no-findings forms, and documented that verdict-shaped lines after a `Findings:` header are findings text rather than a conflict.
- Documented that Specs dropped by the total context bound are named in `skippedSpecs`; regenerated the distributed mirror with `make skills-sync` (exit 0).
- Focused checks: `rtk git diff --check` passed; a post-sync `rtk rg` contract probe found the whole-answer rule, all three findings-none forms, verdict-shaped-line rule, and dropped-Spec record rule in `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, and `docs/user-guide/commands.md`.
- Daemon-owned Verification was not run, as required for this assigned Task execution mode.
