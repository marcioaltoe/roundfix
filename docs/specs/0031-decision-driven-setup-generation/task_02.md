---
task: task_02
spec: 0031-decision-driven-setup-generation
status: pending
type: backend
complexity: high
---

# Task 02: Expose conditional setup preview

## Overview

Deliver the first observable Decision Plan path through the real CLI. A
maintainer can select a profile in an empty repository and receive the profile,
skill setup, unresolved decisions, definite operations, and conditional
operations without allowing audit or blocked apply to write any file.

## Requirements

1. MUST resolve profile entry decisions and activated dependent decisions to a
   deterministic fixed point.
2. MUST expose the selected profile, canonical setup, and active or conditional
   modules in text and JSON results.
3. MUST produce definite and conditional `plannedChanges`, with every
   conditional operation naming its decision ID and triggering value.
4. MUST let `audit --profile` preview a missing-manifest repository instead of
   returning before profile resolution.
5. MUST return the same selection and preview from blocked apply while
   retaining the existing exit-code precedence and JSON envelope.
6. MUST perform no writes when the manifest or any required decision is
   unresolved, invalid, or missing.

## Subtasks

- [ ] Implement fixed-point resolution for entry and dependent decisions.
- [ ] Build prospective artifact plans for definite and conditional branches.
- [ ] Add additive selection and planned-change fields to text and JSON output.
- [ ] Route missing-manifest audit through prospective profile resolution.
- [ ] Reuse the prospective plan in blocked apply without a second planner.
- [ ] Add real-process preview, determinism, exit-code, and no-write tests.

## Acceptance Criteria

- [ ] The QA-03 TypeScript/Bun first-run command returns exit code `3`, all
      unresolved entry decisions, and a non-empty managed-change preview.
- [ ] A conditional artifact reports `state: conditional` and the exact
      decision/value that would activate it; definite artifacts omit a
      condition.
- [ ] Audit and blocked apply report semantically identical selection and
      planned operations for the same profile and answers.
- [ ] Repeating preview against unchanged input produces identical semantic
      output.
- [ ] File snapshots taken before missing-manifest audit and unanswered apply
      match byte-for-byte afterward.
- [ ] Existing text/JSON finding fields and exit codes remain compatible.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`
- interface: `docs/specs/0030-context-driven-agent-instructions/qa/evidence/2026-07-16-manual-variations-03/qa-03-first-apply-questions-and-preview.json`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview*.py'` — expected: selection, definite/conditional preview, audit/apply agreement, deterministic output, exit codes, and no-write behavior pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — expected: the existing read-only audit contract remains passing.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 2, 3, 5; Core Features 5–6; Success Criteria QA-03.
- `_techspec.md` → Interfaces: DecisionPlan; API contracts: Read-only preview and Apply; Build Order 2.
- ADR-0047.
