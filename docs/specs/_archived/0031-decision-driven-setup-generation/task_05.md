---
task: task_05
spec: 0031-decision-driven-setup-generation
status: completed
type: backend
complexity: high
---

# Task 05: Render decision values and enforce semantic audit

## Overview

Make enum and string decisions observable in the guidance that owns them. The
same renderer will produce expected digests for apply and audit, so a manifest
cannot remain clean when its domain, runtime, Verification, or language content
disagrees with its durable answers.

## Requirements

1. MUST select approved single-context or multi-context domain templates from
   `domain.layout` and render no guidance for the unselected layout.
2. MUST render `runtime.backend`, `runtime.design`, and `verification.gate`
   into autonomous guidance when that module is active.
3. MUST enforce the `language.generated` policy through template selection and
   managed-content validation; the current accepted value remains `English`.
4. MUST allow templates to reference only declared bindings and reject a
   missing concrete value or undeclared token before any write.
5. MUST reject control characters, line breaks, and ownership-marker syntax in
   inline decision values and choose an inline-code delimiter the value cannot
   close.
6. MUST calculate managed digests from rendered content and make audit derive
   its expected content from the same renderer as apply.
7. MUST report stale or contradictory managed content with actionable findings
   and the corresponding refresh or removal preview.

## Subtasks

- [x] Add value-selected domain templates for both supported layouts.
- [x] Implement declared token binding and safe Markdown inline rendering.
- [x] Bind backend runtime, design runtime, and Verification values to the
      autonomous guide.
- [x] Resolve the English language policy through the Decision Plan.
- [x] Derive expected digests and semantic audit from rendered artifacts.
- [x] Add non-default, unsafe-input, drift, and audit/apply agreement tests.

## Acceptance Criteria

- [x] Single-context and multi-context answers produce distinct approved domain
      guidance and a clean audit for the selected value.
- [x] The exact stored backend runtime, design runtime, and Verification values
      appear in active autonomous guidance and its rendered digest inventory.
- [x] Newline, control-character, or ownership-marker injection returns invalid
      input, performs no write, and names the unsafe decision.
- [x] Values containing Markdown backticks render without breaking out of their
      inline-code boundary.
- [x] Editing rendered decision content causes audit to report drift and preview
      the same refresh that apply performs.
- [x] The QA-07 alternate values pass every rendering and omission assertion,
      then a second apply produces no file change.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/assets/decisions.json`
- interface: `.agents/skills/setup-context-driven/assets/templates/index.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`
- interface: `docs/specs/0030-context-driven-agent-instructions/qa/evidence/2026-07-16-manual-variations-03/qa-07-rust-alternate-decision-behavior.json`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_decision_rendering*.py'` — expected: domain selection, scalar rendering, language policy, rendered digests, unsafe-input rejection, semantic drift, and QA-07 replay pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — expected: stable findings and read-only audit behavior remain passing with rendered artifacts.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 1, 3, 5; Core Features 4, 6; Success Criteria QA-07 and per-decision macro coverage.
- `_techspec.md` → Interfaces; Data models: Decision effects, Managed artifacts and templates; API contracts: Semantic audit; Build Order 4.
- ADR-0047.

## Result

- Implemented value-selected domain guidance for single-context and
  multi-context layouts, with only the selected guide emitted and audited.
- Implemented declared render bindings for autonomous runtime and Verification
  values, including safe Markdown inline-code rendering for backtick values.
- Enforced unsafe inline decision rejection for control characters, newlines,
  and ownership-marker syntax before any write.
- Made managed digests and audit expectations use the same rendered artifact
  path as apply, so edited rendered content reports drift and previews the same
  refresh apply performs.
- Added rendering coverage for distinct domain layouts, autonomous scalar
  rendering and manifest digests, language policy validation, unsafe input,
  backtick-safe rendering, semantic audit drift, and QA-07 alternate values.

Evidence:

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_decision_rendering*.py'` — passed, 6 tests.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — passed, 7 tests.
- `rtk git diff --check` — passed.
- `rtk make verify` — passed: setup-context-driven Python suite, `rtk go test ./...`, asset catalog load for both skill copies, `roundfix skills check`, and Go build.
