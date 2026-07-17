---
task: task_05
spec: 0030-context-driven-agent-instructions
status: completed
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

- [x] Rewrite the setup flow around audit, stored decisions, preview, and explicit apply.
- [x] Define question routing for missing and incompatible decision codes.
- [x] Add the portable Secondbrain module, root pointer, detailed guide, and safety rules.
- [x] Add validation for Secondbrain pointer/guide consistency and English generated output.
- [x] Add end-to-end tests for first setup, unchanged rerun, new-decision migration, and Secondbrain enable/disable.
- [x] Ensure the skill names the selected canonical skill setup and explains optional extra-skill reporting without offering removal.

## Acceptance Criteria

- [x] An unchanged rerun requires no repeated confirmation for stored compatible decisions.
- [x] A newly required decision yields one precise question and persists the answer for later runs.
- [x] Apply is never invoked before the user sees the managed change summary and confirms it.
- [x] Secondbrain opt-in produces one concise root pointer and one complete read-only guide.
- [x] Secondbrain opt-out produces neither artifact, and disabling it removes only marked managed content.
- [x] Generated guidance forbids reading or exposing secrets and forbids writes to the Secondbrain, raw sources, and project mirrors.
- [x] All generated content is English and the root agent instructions remain an index rather than a duplicated manual.

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

## Result

Rewrote `setup-context-driven` as an audit-first orchestration skill and added the optional Secondbrain integration as a portable module. The workflow now tells agents to audit before questions, reuse stored manifest decisions, ask only unresolved decision codes one at a time, preview profile/modules/findings/planned changes, require confirmation before `apply`, and explain that only setup-owned manifest fields and managed Markdown boundaries can change. The Secondbrain module is selected only when `secondbrain.enabled=true`; disabling it removes only previously marked setup-owned Secondbrain content.

Evidence by acceptance criterion:

- Unchanged reruns reuse stored compatible decisions: `test_unchanged_rerun_reuses_stored_compatible_decisions`.
- A newly required decision produces one precise `decision.required` finding and persists the answer: `test_new_required_decision_routes_one_question_and_persists_answer`.
- Apply is gated behind preview and confirmation in the skill workflow: `test_skill_workflow_requires_preview_and_confirmation_before_apply`.
- Secondbrain opt-in creates one compact root pointer and one read-only guide: `test_secondbrain_opt_in_creates_compact_pointer_and_read_only_guide`.
- Secondbrain opt-out creates no artifact, and disable removes only marked content: `test_secondbrain_opt_out_creates_no_pointer_or_guide` and `test_disabling_secondbrain_removes_only_marked_managed_content`.
- Secret safety and read-only Secondbrain boundaries are generated and audited: `test_secondbrain_audit_reports_missing_safety_rule`.
- Generated content stays English and the root block remains an index pointer: `test_secondbrain_generated_content_is_english_and_root_is_index_only`.

Verification:

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow*.py'` — passed, 4 tests.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_secondbrain*.py'` — passed, 5 tests.
- `rtk git diff --check` — passed.
- Regression checks also passed: `test_assets*.py` (6 tests), `test_audit*.py` (6 tests), `test_apply*.py` (8 tests), and `test_skills*.py` (6 tests).
- `rtk make verify` — passed; Go tests passed in 19 packages, Roundfix skill check passed, and the binary built.
