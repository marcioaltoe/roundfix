---
task: task_04
spec: 0031-decision-driven-setup-generation
status: pending
type: backend
complexity: high
---

# Task 04: Make autonomous work and Secondbrain decision-driven

## Overview

Generalize the optional-module path so autonomous work and Secondbrain resolve
through declared effects. Enabling autonomous work introduces only its three
dependent questions; disabling it or Secondbrain omits their managed guidance
without weakening Secondbrain safety or managed-removal boundaries.

## Requirements

1. MUST activate the autonomous-work module only when
   `autonomous.enabled=true` and omit its root block and guide when false.
2. MUST require `runtime.backend`, `runtime.design`, and `verification.gate`
   only after autonomous work is enabled.
3. MUST report no runtime or Verification decision finding when autonomous work
   is disabled and those stored answers are absent.
4. MUST resolve Secondbrain activation through the general Decision Plan and
   remove the Secondbrain-only branch from the execution path.
5. MUST preserve the complete Secondbrain read-only, citation, Hermes, mirror,
   raw-source, and secret-safety guidance when enabled.
6. MUST remove only proven setup-owned autonomous or Secondbrain content when a
   previously true decision becomes false.

## Subtasks

- [ ] Move autonomous module activation and dependent decisions into catalog
      effects.
- [ ] Resolve enabled and disabled autonomous branches prospectively and
      concretely.
- [ ] Route Secondbrain through the same optional-module effect mechanism.
- [ ] Remove the dedicated Secondbrain activation branch from the generator.
- [ ] Add enable/disable, dependent-question, removal, and safety regressions.
- [ ] Prove surrounding repository-authored content survives both opt-outs.

## Acceptance Criteria

- [ ] `autonomous.enabled=false` generates no autonomous root block or guide and
      asks for none of its three dependent values.
- [ ] `autonomous.enabled=true` reports exactly the missing runtime and
      Verification decisions, then generates autonomous guidance after they are
      supplied.
- [ ] Unanswered autonomous and Secondbrain decisions appear as conditional
      modules and operations in preview without repository writes.
- [ ] Secondbrain enablement still produces the compact pointer and complete
      safety guide; disablement produces neither managed artifact.
- [ ] Changing either decision from true to false removes only marked managed
      content and preserves surrounding owner bytes.
- [ ] Repeated concrete apply and audit are idempotent and deterministic for
      every tested branch.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/assets/modules/autonomous-work.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/secondbrain.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_workflow.py`
- interface: `.agents/skills/setup-context-driven/tests/test_secondbrain.py`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_autonomous_secondbrain_decisions*.py'` — expected: conditional preview, dependent questions, enable/disable, managed removal, safety, deterministic audit, and idempotency pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_secondbrain*.py'` — expected: existing Secondbrain pointer, guide, safety, language, and removal contracts remain passing.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 1, 3, 5; Core Features 2–3, 6; Success Criteria QA-07.
- `_techspec.md` → Data models: Decision effects; API contracts: Apply; Integration points: Secondbrain; Build Order 3.
- ADR-0047.
