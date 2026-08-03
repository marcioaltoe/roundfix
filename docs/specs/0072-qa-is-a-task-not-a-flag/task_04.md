---
task: task_04
spec: 0072-qa-is-a-task-not-a-flag
status: completed
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

## Result

### Implementation

- `write-tasks` now requires every post-contract decomposition to choose
  exactly one authored QA shape: a manifest-named terminal `qa` Task covering
  every non-QA leaf, or `qa: declined` with a non-empty `qa_reason`. It refuses
  a graph with neither or conflicting shapes while preserving proven legacy
  graphs byte-identically.
- The task artifact reference now projects an included `qa: task_04` node and
  `type: qa` task row, shows its leaf dependencies, and gives the alternative
  reasoned-decline frontmatter.
- `qa-gate` now runs only from the unique authored terminal `qa` Task named by
  the manifest. Its scope preflight checks that node and its completed
  dependencies instead of presenting the gate as a per-Run request after the
  graph.
- Regenerated the two shipped skill mirrors with `make skills-sync`, then ran
  the ADR-0081 sanctioned digest regeneration. The command updated the three
  setup pins, catalog digest and normalized snapshot, and the parity-corpus
  fixture and manifest.

### Focused checks

- Red signal: pre-change inspection showed that `write-tasks` allowed only
  implementation Task Types and had no include-or-decline rule, while the
  `qa-gate` description said to use it after the last Task or when asked.
- `rtk make skills-sync` exited 0.
- `rtk make baseline-digests` exited 0, reported all six regeneration and
  strict-validation test invocations passing, and returned
  `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":true}`.
- A focused `rtk awk` contract scan over the three canonical documents exited
  0 and reported `write_tasks_decline=1 terminal=1 refusal=1
  template_include=1 template_decline=1 qa_gate_authored=1
  old_request_form=0`.
- `rtk /usr/bin/diff -r .agents/skills/write-tasks skills/write-tasks` and
  `rtk /usr/bin/diff -r .agents/skills/qa-gate skills/qa-gate` both exited 0.
- `rtk git diff --check` exited 0.
- `rtk git diff --name-only` listed only the two canonical skills, their
  mirrors, the seven ADR-0081 derived paths, and this Task file.
- The Task's declared `## Verification` commands were not run; they remain
  Daemon-owned.

### Acceptance criteria evidence

1. The canonical `write-tasks` skill and its task template name the
   exactly-one-of-two rule, both manifest forms, the unique `type: qa` node,
   terminal placement, and coverage of every non-QA leaf; the focused contract
   scan found every required clause.
2. The canonical `qa-gate` description and scope preflight name the authored
   terminal node, and the focused scan found neither former per-Run trigger
   phrase.
3. Both focused recursive directory comparisons exited 0 after
   `make skills-sync`. The authoritative `make skills-sync-check` command is
   reserved for Daemon Verification and was not run in this Agent turn.
4. The post-regeneration changed-path inspection contains no path outside the
   authorized canonical skills and mirrors, ADR-0081 deterministic fallout,
   and this Task file.
