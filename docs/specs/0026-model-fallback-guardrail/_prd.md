---
spec: 0026-model-fallback-guardrail
status: active
created: 2026-07-11
surfaces: [cli, backend, docs]
---

# Model Fallback Guardrail

When a configured Agent Model fails Agent selection Preflight Validation — a
model the adapter does not advertise, a rejected Default Reasoning Effort, or
an adapter gap like fable before claude-code-acp shipped effort support — the
Run simply fails, and the developer must diagnose adapters by hand before any
work can start. Roundfix knows the runtime's Model Catalog and can prove which
Agent Model actually works, so it should offer the recovery itself: propose
each provider's most recent functional Agent Model at its highest functional
reasoning effort. Because switching models can consume tokens very
differently than the user planned, the fallback must never be silent or
autonomous — Roundfix reports the problem and starts work only after explicit
human confirmation.

## Goals

- A failed Agent selection produces a proven Fallback Selection offer — the
  newest functional Agent Model of the same ACP Runtime at its highest
  functional reasoning effort — instead of only an error.
- No Run ever starts on a Fallback Selection without explicit human
  confirmation, in any mode.
- Non-interactive callers get an actionable failure that names the exact
  explicit re-run, never an autonomous model decision.
- Orchestrating agents driving Roundfix through the Roundfix Skill relay the
  confirmation to the human user instead of deciding themselves.

## User Stories

1. As a developer starting an interactive Run with a non-functional Agent
   Model, I want Roundfix to report the selection failure and offer a proven
   Fallback Selection I can confirm or decline, so that I recover in one step
   without debugging adapters.
2. As a developer running non-interactive or Detached Runs, I want a failed
   selection to fail Preflight Validation with the probed Fallback Selection
   and the exact re-run command, so that nothing chooses a different model —
   and a different token cost — without me.
3. As an orchestrating agent driving Roundfix for a user, I want the Roundfix
   Skill to instruct me to relay the fallback confirmation to the user before
   proceeding, so that the model decision is never taken autonomously.
4. As a Supervisor inspecting a Run started from a confirmed fallback, I want
   the Run header and Run record to carry the effective Agent Model and
   reasoning effort, so that what actually ran is never ambiguous.

## Core Features

1. **P0 - Fallback probe.** When Agent selection fails Preflight Validation,
   Roundfix walks the failed runtime's Model Catalog newest-first and probes
   candidates with disposable Agent Sessions until one proves functional,
   then determines that model's highest functional reasoning effort (probing
   effort values highest-first; a model-managed model has none). The probe
   never crosses to another ACP Runtime and never re-proposes the failed
   selection.
2. **P0 - Interactive confirmation.** In interactive mode, Roundfix reports
   the failed selection, the proven Fallback Selection, and a token-cost
   caveat, then asks for confirmation before creating the Run. Declining
   fails Preflight Validation exactly as today. A confirmation applies to
   that Run only; configuration is never modified.
3. **P0 - Non-interactive contract.** In non-interactive contexts (no-input,
   detached, or any caller that cannot prompt), the command fails Preflight
   Validation with a report naming the failed selection, the proven Fallback
   Selection, and the exact re-run command with explicit model and
   reasoning-effort flags. Exit codes are unchanged, and no fallback flag or
   config key pre-authorizes the switch.
4. **P0 - Honest effective selection.** A Run started from a confirmed
   fallback records and displays the effective Agent Model and reasoning
   effort through the existing effective-selection surfaces.
5. **P0 - Skill guidance.** The Roundfix Skill and user documentation ship
   with the behavior: they explain the fallback offer, that an orchestrating
   agent must relay the confirmation to the human user and never confirm
   autonomously, and the explicit-flags re-run recipe for non-interactive
   Runs.
6. **P1 - No functional candidate.** When no catalog candidate proves
   functional, the failure reports the probe attempts and falls back to
   today's actionable selection error.

## User Experience

An interactive `implement`, `resolve`, or `watch` with a broken configured
model prints the selection failure, the proven Fallback Selection (model plus
effort, or model-managed), and a caveat that a different model can consume
tokens differently, then prompts for confirmation. Confirming starts the Run
whose header shows the effective selection; declining ends with the familiar
Preflight Validation failure. The same command with `--no-input` or
`--detach` prints the failure report with a copy-paste re-run command instead
of prompting.

## Non-Goals / Out of Scope

- Any autonomous or pre-authorized fallback: no `--allow-fallback` flag, no
  config key, no daemon policy that switches models without a human.
- Persisting a confirmed fallback into User Config or Project Config.
- Cross-runtime fallback (a codex failure never proposes a claude model).
- Fallback for non-selection failures: authentication, missing Node/acpx,
  worktree, or Verification failures keep their existing contracts.
- Mid-Run model replacement after a selection already succeeded.
- Retrying the originally failed Agent Model as its own fallback candidate.
- Changing built-in Default Agent Model or Default Reasoning Effort values.

## Success Metrics

- A selection failure in interactive mode yields exactly one confirmation
  prompt naming a probed, functional Fallback Selection; declining creates no
  Run and changes no state.
- A selection failure in non-interactive mode creates no Run and prints the
  exact re-run command whose explicit flags succeed when executed.
- Zero code paths start Agent work on a model the user did not configure or
  explicitly confirm.
- The Roundfix Skill documents the orchestrator relay rule, and the embedded
  copy has zero sync drift.

## Decisions

- The fallback candidate is discovered by probing the Model Catalog
  newest-first with disposable Agent Sessions, not by static claims. See
  ADR-0041.
- Non-interactive confirmation is the explicit-flags re-run; no new flag or
  config key exists. See ADR-0041.
- A confirmed fallback applies to that Run only and never mutates config.
- "Fallback Selection" enters the glossary with the shipped documentation.

## Open Questions

None.
