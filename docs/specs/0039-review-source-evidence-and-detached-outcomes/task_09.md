---
task: task_09
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: docs
complexity: low
---

# Task 09: Align the protected Roundfix Skill pair

## Overview

Publish Review Source Evidence and Detached outcome handling in the canonical
and generated Roundfix Skill. This tooling-only slice is restricted to the two
exact files authorized by the maintainer and this Task file.

## Requirements

1. MUST align `.agents/skills/roundfix/SKILL.md` with Review Skipped, bounded
   transient retry, unknown issue knowledge, Detached monitoring, notification
   receipts, and artifact-only Evidence inheritance.
2. MUST apply byte-identical content to `skills/roundfix/SKILL.md`.
3. MUST limit repository changes to those two authorized `SKILL.md` files and
   this `task_09.md` file.
4. MUST NOT run `make skills-sync`, because it rewrites every owned Skill
   directory; use the read-only sync check.
5. MUST leave code, tests, manifests, public docs, other owned Skills,
   upstream-managed Skills, locks, and recommendation files unchanged.

## Subtasks

- [ ] Update the canonical Roundfix Skill review-evidence contract.
- [ ] Apply the identical generated Roundfix Skill copy.
- [ ] Verify exact changed-file scope and byte identity.
- [ ] Confirm shipped Skill and full-gate compatibility.

## Acceptance Criteria

- [ ] The Skill distinguishes Review Skipped and unknown issue knowledge from
      successful zero-issue review.
- [ ] The Skill surfaces the Supervisor monitor command and exact next actions.
- [ ] Canonical and generated files are byte-identical.
- [ ] Git evidence contains only the two authorized Skill paths and this Task
      file.
- [ ] No other protected or upstream-managed Skill changes.
- [ ] Shipped Skill validation and the complete repository gate pass.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — expected: no output and exit zero.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0039-review-source-evidence-and-detached-outcomes/task_09.md") {print; bad=1}} END {exit bad}'`
  — expected: no changed path outside the authorized pair and this Task file.
- `rtk make skills-sync-check`
  — expected: every canonical/generated owned Skill pair has no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check`
  — expected: every shipped Roundfix Skill contract passes.
- `rtk git diff --check`
  — expected: no whitespace errors.
- `rtk make verify`
  — expected: formatting, tests, Skill checks, and build pass.

## References

- `_prd.md` → Goals; User Experience; Decisions; Project Constraints.
- `_techspec.md` → Build Order 9; Decisions.
- `docs/agents/spec-routing.md` → tooling authorization and changed-file
  postflight.
