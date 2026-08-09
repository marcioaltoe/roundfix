---
task: task_06
spec: 0093-a-spec-that-validates-itself
status: pending
type: infra
complexity: high
---

# Task 06: Let the QA gate read product instead of paperwork

## Overview

Measured on Spec 0090, eight of the gate's sixteen rows audited artifacts rather
than product, and the finding that failed the Run came from one of those eight —
reached after 461 tool calls and two context compactions. Every rule those rows
apply is now decided by the checker during authoring. This Task removes them from
the gate's matrix and keeps what a file read cannot settle.

## Requirements

1. MUST remove from the `qa-gate` matrix contract every row whose rule the
   checker now decides, naming each removed rule and the checker rule that runs
   it.
2. MUST keep every row that needs judgement or a running surface: the Spec's
   goals exercised through the surfaces a user reaches, with captured evidence.
3. MUST keep the post-commit rows that no authoring check can answer — which
   paths a Task actually touched against its bounded list — and state that they
   run as commands rather than judgements.
4. MUST NOT reduce what the loop detects. A rule may move; a rule may not
   disappear. Any rule with no checker equivalent stays in the gate.
5. MUST state in the skill that a clean authoring check is a precondition of the
   gate, not a substitute for it.
6. MUST record in the commit message which half of the standing grant it serves
   and how reliability was preserved, per that grant's obligations.

## Subtasks

- [ ] Map each governance row to the checker rule that replaces it.
- [ ] Remove the mapped rows; keep the unmapped ones.
- [ ] State the precondition relationship.

## Acceptance Criteria

- [ ] Every removed row has a named checker rule running its check.
- [ ] No rule is removed without an equivalent.
- [ ] The post-commit rows remain, described as commands.
- [ ] The goal and surface rows are unchanged.

## Rehearsal Cases

- Case: a Spec whose PRD omits a Project Constraints row; Observation: the
  authoring check reports it, and the gate has no row for it.
- Case: a Task whose commit touched a path outside its bounded list;
  Observation: the gate still reports it, because no authoring check can see a
  commit that does not exist yet.
- Case: a Spec whose goals do not work through the CLI; Observation: the gate
  reports it, unchanged from today.

## Bounded scope

Covered by the standing grant at
`docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`.
This Task may create or modify only:

- `.agents/skills/qa-gate/SKILL.md`
- `skills/qa-gate/SKILL.md`, as `make skills-sync` fallout under ADR-0081
- `docs/specs/0093-a-spec-that-validates-itself/task_06.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'spec check' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the gate names the authoring check it now depends on. This string does not exist in the file before this Task.
- `grep -q 'bounded' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the post-commit path audit was kept rather than removed with the rest.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0.
- `test -z "$(git diff --name-only -- . ':!.agents/skills/qa-gate/SKILL.md' ':!skills/qa-gate/SKILL.md' ':!docs/specs/0093-a-spec-that-validates-itself/task_06.md')"` — expected: exits 0, proving no path outside the bounded list moved.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`

## References

- `_prd.md` → Goal 3.
- `_techspec.md` → Build Order 6.
- ADR-0081, ADR-0091, ADR-0096, ADR-0117.
