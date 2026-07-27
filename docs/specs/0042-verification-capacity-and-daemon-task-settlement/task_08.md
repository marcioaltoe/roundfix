---
task: task_08
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
type: docs
complexity: low
---

# Task 08: Align the protected authorial Skill pairs

## Overview

Publish the completed Daemon-owned status, Verification Capacity, temporary
failure, and Task Type-routing contracts in the two protected authorial Skills.
This is a tooling-only slice bounded to the four exact canonical/generated
files authorized by the maintainer and this Task file.

## Requirements

1. MUST align `.agents/skills/implement-task/SKILL.md` with the
   implementation-ready handoff: the Agent does not edit Task status, run the
   declared `## Verification`, or claim the terminal Task verdict.
2. MUST align `.agents/skills/roundfix/SKILL.md` with independent capacities,
   Daemon-owned settlement, observable Verification phases, exit `75`, one
   exclusive retry, and ADR-0051 Task Type-selected Agent Sessions.
3. MUST apply identical content to `skills/implement-task/SKILL.md` and
   `skills/roundfix/SKILL.md`.
4. MUST limit repository mutations to those four authorized `SKILL.md` files
   and this `task_08.md` file.
5. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use `make skills-sync-check` as read-only verification.
6. MUST leave code, tests, configuration, manifests, public docs, other
   Roundfix-owned Skills, upstream-managed Skills, and lock files unchanged.

## Subtasks

- [ ] Update the canonical `implement-task` Skill.
- [ ] Apply the identical generated `implement-task` copy.
- [ ] Update the canonical Roundfix Skill.
- [ ] Apply the identical generated Roundfix copy.
- [ ] Verify changed-file scope, byte identity, shipped Skill contracts, and
      full-gate compatibility.

## Acceptance Criteria

- [ ] No supported Skill tells an Implement Agent to run declared Task
      Verification, edit Task status, or settle its terminal verdict.
- [ ] The Roundfix Skill describes both capacities, the exit-75 retry contract,
      and per-Task Task Type routing consistently with shipped behavior.
- [ ] Each canonical/generated Skill pair is byte-identical.
- [ ] Git changed-file evidence for this Task contains only the four authorized
      Skill files and `task_08.md`.
- [ ] No other protected or upstream-managed tooling change.
- [ ] Shipped Skill validation and the complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `docs/agents/autonomous-work.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md`
  — expected: no output and exit zero.
- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/implement-task/SKILL.md" && path != "skills/implement-task/SKILL.md" && path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0042-verification-capacity-and-daemon-task-settlement/task_08.md") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the four authorized Skill files and this
  Task file.
- `rtk make skills-sync-check` — expected: every canonical/generated owned
  Skill pair has no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: every
  shipped Roundfix Skill contract passes.
- `rtk git diff --check` — expected: no whitespace errors.
- `rtk make verify` — expected: formatting, tests, Skill synchronization,
  shipped Skill validation, and build all pass.

## References

- `_prd.md` → Core Feature 10; Decisions; Project Constraints.
- `_techspec.md` → Integration Points; Build Order 8; Decisions.
- `task_02.md` → implemented Daemon status and handoff contract.
- `task_07.md` → completed operator wording and ADR-0051 alignment.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.

## Result

Pending.
