---
task: task_05
spec: 0093-a-spec-that-validates-itself
status: pending
type: infra
complexity: medium
---

# Task 05: Wire the checker into PRD and TechSpec authoring

## Overview

`write-tasks` already ends by running the checker; this Task gives `write-prd`
and `write-techspec` the same closing step, scoped to what each stage can
decide. An author fixes the artifact while it is still open, for the price of a
sub-second command, instead of meeting the same finding at a gate.

## Requirements

1. MUST add a closing validation step to `write-prd` and to `write-techspec`,
   each running the checker scoped to its own stage.
2. MUST make an error-level finding block the stage: neither skill may report
   completion while one stands.
3. MUST state, in each skill, which classes the checker does **not** decide at
   that stage, so a clean result is not read as full coverage.
4. MUST NOT change what either skill decides about its artifact's content; this
   Task adds a verification step and nothing else.
5. MUST record in the commit message which half of the standing grant it serves
   and how reliability was preserved, per that grant's obligations.

## Subtasks

- [ ] Add the closing step to `write-prd`.
- [ ] Add the closing step to `write-techspec`.
- [ ] Name the classes each stage does not decide.

## Acceptance Criteria

- [ ] Both skills instruct the author to run the stage-scoped checker before
      reporting.
- [ ] Both state that an error-level finding blocks the report.
- [ ] Both name what the stage does not cover.
- [ ] The generated copies match their sources.

## Bounded scope

Covered by the standing grant at
`docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`.
This Task may create or modify only:

- `.agents/skills/write-prd/SKILL.md`
- `.agents/skills/write-techspec/SKILL.md`
- `skills/write-prd/SKILL.md` and `skills/write-techspec/SKILL.md`, as
  `make skills-sync` fallout under ADR-0081
- `docs/specs/0093-a-spec-that-validates-itself/task_05.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'spec check' .agents/skills/write-prd/SKILL.md` — expected: exits 0. This string does not exist in the file before this Task.
- `grep -q 'spec check' .agents/skills/write-techspec/SKILL.md` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0, proving the generated copies match their sources.
- `test -z "$(git diff --name-only -- . ':!.agents/skills/write-prd/SKILL.md' ':!.agents/skills/write-techspec/SKILL.md' ':!skills/write-prd/SKILL.md' ':!skills/write-techspec/SKILL.md' ':!docs/specs/0093-a-spec-that-validates-itself/task_05.md')"` — expected: exits 0, proving no path outside the bounded list moved.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 5.
- ADR-0081, ADR-0117.
