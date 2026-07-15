---
task: task_05
spec: 0030-context-driven-agent-instructions
status: pending
type: docs
complexity: medium
---

# Task 05: Add optional Secondbrain guidance and setup orchestration

## Overview

Complete the user-facing setup flow by making `SKILL.md` audit-first and decision-aware, and by adding the optional Secondbrain module. The slice is usable when a fresh agent can inspect a repository, reuse its durable decisions, ask only unresolved questions, preview managed changes, and generate concise English guidance.

## Requirements

1. MUST rewrite the skill workflow to run audit before asking questions, load stored decision codes, and ask one question at a time only for unresolved decisions.
2. MUST present detected profile/module composition, blocking findings, planned managed changes, and optional cleanup information without dumping full generated documents by default.
3. MUST require confirmation before apply and MUST explain that only setup-owned boundaries will change.
4. MUST add Secondbrain as an explicit opt-in module with a compact mandatory root pointer and a detailed supporting guide.
5. MUST include Secondbrain index-first query order, `qmd` usage, project-mirror guidance, file citation requirements, read-only boundaries, Hermes escalation, and secret-exposure prohibitions.
6. MUST omit Secondbrain files and pointers when the module is disabled and remove only previously managed Secondbrain artifacts when it is later disabled.
7. MUST keep generated repository content in English and preserve canonical domain terms from the target repository.
8. SHOULD keep `SKILL.md` concise by routing detailed asset contracts and variant guidance to bundled resources.

## Subtasks

- [ ] Rewrite the setup flow around audit, stored decisions, preview, and explicit apply.
- [ ] Define question routing for missing and incompatible decision codes.
- [ ] Add the portable Secondbrain module, root pointer, detailed guide, and safety rules.
- [ ] Add validation for Secondbrain pointer/guide consistency and English generated output.
- [ ] Add end-to-end tests for first setup, unchanged rerun, new-decision migration, and Secondbrain enable/disable.
- [ ] Ensure the skill names the selected canonical skill setup and explains optional extra-skill reporting without offering removal.

## Acceptance Criteria

- [ ] An unchanged rerun requires no repeated confirmation for stored compatible decisions.
- [ ] A newly required decision yields one precise question and persists the answer for later runs.
- [ ] Apply is never invoked before the user sees the managed change summary and confirms it.
- [ ] Secondbrain opt-in produces one concise root pointer and one complete read-only guide.
- [ ] Secondbrain opt-out produces neither artifact, and disabling it removes only marked managed content.
- [ ] Generated guidance forbids reading or exposing secrets and forbids writes to the Secondbrain, raw sources, and project mirrors.
- [ ] All generated content is English and the root agent instructions remain an index rather than a duplicated manual.

## Context

- instruction: `.agents/skills/setup-context-driven/SKILL.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `docs/agents/secondbrain.md`
- interface: `AGENTS.md`

## Verification

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow*.py'` — expected: stored-decision reuse, question routing, preview, and explicit apply behavior pass.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_secondbrain*.py'` — expected: enable, disable, pointer, guide, safety, and language contracts pass.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → User Stories 2, 7; Core Features 8, 12–14; User Experience; Decisions.
- `_techspec.md` → System Architecture: SKILL orchestration; Data Models: decisions; Coverage Map; Build Order 4–5.
- ADR-0046.
