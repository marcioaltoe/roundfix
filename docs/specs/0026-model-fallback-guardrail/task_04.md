---
task: task_04
spec: 0026-model-fallback-guardrail
status: completed
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

- [x] Update the canonical Roundfix Skill with the guardrail contract and
      orchestrator relay rule, and sync the embedded copy.
- [x] Update the README Agent selection guidance with both flows.
- [x] Add the Fallback Selection glossary entry to CONTEXT.md.

## Acceptance Criteria

- [x] The Roundfix Skill names the orchestrator relay rule and the
      explicit-flags re-run recipe, and the embedded copy has zero drift.
- [x] The README documents the interactive confirmation and the
      non-interactive report.
- [x] CONTEXT.md defines Fallback Selection.

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

## Result

Every guidance surface now describes the confirmation-gated Fallback
Selection contract. The canonical and embedded Roundfix Skills name the exact
interactive prompt, exit-2 non-interactive report, explicit model/effort
re-run shape, model-managed rendering, and the hard rule that an orchestrating
agent must relay the decision to the human user and never confirm
autonomously. The README documents both user flows, and the glossary defines
Fallback Selection with its same-runtime, one-Run-only boundary.

Verification:

- `rtk make skills-sync-check`: passed — canonical and embedded skill bundles
  have zero drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check`: passed — all 14
  shipped skills validated.
- `rtk make verify`: passed — 1,070 tests across 19 packages, skill sync and
  validation checks, and the CLI build.

Acceptance evidence:

1. `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` contain
   the orchestrator relay prohibition and explicit `--model` plus
   `--reasoning-effort` re-run recipe; the sync gate passed.
2. `README.md` documents the interactive confirmation, token-cost caveat,
   decline behavior, and prompt-free exit-2 report for no-input, detached,
   and non-interactive stderr contexts.
3. `CONTEXT.md` defines Fallback Selection as a proven same-runtime
   alternative that applies to one Run after explicit human confirmation and
   never changes configuration.

Follow-ups: none.
