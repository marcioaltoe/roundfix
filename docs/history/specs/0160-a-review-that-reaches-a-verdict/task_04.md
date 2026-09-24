---
task: task_04
spec: 0160-a-review-that-reaches-a-verdict
status: completed
type: docs
complexity: low
---

# Task 04: Keep the shipped skill and the guide true

## Overview

The Roundfix skill says `claude` is refused and that only an exact `No findings` passes; both become false with this Spec.

## Requirements

1. MUST state in `.agents/skills/roundfix/SKILL.md` that `claude` runs, that the verdict is classified by substance, that the raw answer is kept at `answerPath`, and that a delivery's Spec decisions are reviewed.
2. MUST keep `coderabbit` described as refused.
3. MUST describe the same in `docs/user-guide/commands.md`.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both skill files and the guide describe the new behavior.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -q "answerPath" .agents/skills/roundfix/SKILL.md && grep -q "answerPath" skills/roundfix/SKILL.md && grep -q "answerPath" docs/user-guide/commands.md && ! grep -q "refuses to execute them for now" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implementation:

- Updated the canonical Roundfix skill to document the running `claude` provider,
  substance-based verdict classification, raw answers retained at `answerPath`,
  and review against a delivery's Spec Decisions and rejected alternatives.
- Kept `coderabbit` documented as refused because no supported local CodeRabbit
  review surface is installed or specified.
- Updated the user guide with the same behavior and regenerated the distributed
  `skills/roundfix` mirror using `make skills-sync`.

Focused checks:

- `make skills-sync` — completed; the distributed mirror was regenerated from
  the canonical skill.
- `git diff --check` — passed.
- The declared Verification command, including `go run -buildvcs=false
  ./cmd/roundfix skills check`, was not run because the Daemon owns that gate.

Acceptance evidence:

- Both skill copies and the guide contain the new provider, classification,
  `answerPath`, and Spec-review behavior at `.agents/skills/roundfix/SKILL.md:232-250`,
  `skills/roundfix/SKILL.md:232-250`, and `docs/user-guide/commands.md:145-163`.
- The final skill check remains pending Daemon Verification.
