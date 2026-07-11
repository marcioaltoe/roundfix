---
task: task_04
spec: 0026-model-fallback-guardrail
status: pending
type: docs
complexity: low
---

# Task 04: Ship fallback guardrail guidance

## Overview

Align every guidance surface with the fallback guardrail: the Roundfix Skill
teaches the confirmation contract — including the rule that an orchestrating
agent relays the confirmation to the human user and never confirms
autonomously — the README documents the flows, and the glossary gains the
Fallback Selection term. The slice is verifiable through the Skill sync and
validation gates.

## Requirements

1. MUST document in the canonical Roundfix Skill and its embedded copy,
   together with zero sync drift: the fallback offer on selection failure,
   the interactive confirmation, the non-interactive report with the
   explicit-flags re-run recipe, and the hard rule that an agent
   orchestrating Roundfix must relay the fallback confirmation to the human
   user and never decide autonomously.
2. MUST document the interactive and non-interactive fallback flows in the
   README's Agent selection guidance.
3. MUST add a Fallback Selection entry to the CONTEXT.md glossary consistent
   with the shipped behavior.
4. MUST keep all guidance truthful to the implemented prompts, reports, and
   exit codes.

## Subtasks

- [ ] Update the canonical Roundfix Skill with the guardrail contract and
      orchestrator relay rule, and sync the embedded copy.
- [ ] Update the README Agent selection guidance with both flows.
- [ ] Add the Fallback Selection glossary entry to CONTEXT.md.

## Acceptance Criteria

- [ ] The Roundfix Skill names the orchestrator relay rule and the
      explicit-flags re-run recipe, and the embedded copy has zero drift.
- [ ] The README documents the interactive confirmation and the
      non-interactive report.
- [ ] CONTEXT.md defines Fallback Selection.

## Verification

- `rtk make skills-sync-check` - expected: canonical and embedded Skill
  bundles have zero drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` - expected: every
  shipped Skill passes validation.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → User Story 3; Core Feature 5.
- `_techspec.md` → Build Order 4.
- ADR-0041; `docs/agents/skill-governance.md`.
