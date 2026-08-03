---
task: task_04
spec: 0072-qa-is-a-task-not-a-flag
status: pending
type: docs
complexity: medium
---

# Task 04: Author the gate decision in the owned skills

## Overview

`write-tasks` becomes where the gate decision is made: every decomposition
either emits the gate as a terminal `qa` node depending on every leaf, or
records `qa: declined` with a reason in the manifest — one of the two,
always, for post-contract Specs. `qa-gate` documents that it runs as the
authored terminal node rather than as a per-run request. Skill mirrors sync
via `make skills-sync`, and the deterministic digest fallout regenerates per
ADR-0081.

## Requirements

1. MUST add the authoring rule to `.agents/skills/write-tasks/SKILL.md` and
   its task template reference: the decomposition declares the gate or
   declines it with a reason; a post-contract graph with neither is a
   defect the skill refuses to produce.
2. MUST describe the emission shape: a `task_NN.md` of `type: qa`, terminal,
   depending on every leaf, named by `qa:` in the manifest frontmatter.
3. MUST update `.agents/skills/qa-gate/SKILL.md` to state the gate runs as
   the graph's authored terminal node, with the per-run request form
   removed.
4. MUST regenerate the `skills/` mirrors with `make skills-sync` and let
   the sanctioned digest fallout land in the same change (ADR-0081).
5. MUST change only the bounded files of the recorded authorization: the
   two owned skills, their mirrors, and the deterministic digest fallout.

## Subtasks

- [ ] Write the include-or-decline rule and emission shape in `write-tasks`.
- [ ] Update the task template reference with the `qa` node and manifest
      declaration.
- [ ] Update `qa-gate` to its authored-node contract.
- [ ] Run `make skills-sync` and regenerate digests.

## Acceptance Criteria

- [ ] `write-tasks` names the exactly-one-of-two rule and the terminal-node
      emission shape.
- [ ] `qa-gate` no longer describes a per-run request.
- [ ] `make skills-sync-check` passes.
- [ ] `git status --porcelain` shows no path outside `.agents/skills/`,
      `skills/`, the ADR-0081 digest fallout paths, and this task file.

## Context

- instruction: `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`

## Verification

- `grep -q "qa: declined" .agents/skills/write-tasks/SKILL.md || grep -qr "qa: declined" .agents/skills/write-tasks/references/`
  — expected: exit 0; the decline form is documented.
- `grep -rq "terminal" .agents/skills/write-tasks/SKILL.md` — expected:
  exit 0; the emission shape is stated.
- `make skills-sync-check` — expected: exit 0.
- `go test ./skills -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Feature 1; Goals (declared when authored; declining
  recorded once).
- `_techspec.md` → Build Order 4; Project Constraints (Tooling authority).
