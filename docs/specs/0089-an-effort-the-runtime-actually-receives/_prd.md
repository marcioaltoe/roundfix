---
spec: 0089-an-effort-the-runtime-actually-receives
status: active
created: 2026-08-09
surfaces: [backend, cli, docs]
---

# An effort the runtime actually receives

Spec 0088 made the OpenCode runtime reachable and, to keep Exact Agent Selection
Proof token-free, made it model-managed: Roundfix refuses any non-empty
`reasoning_effort` there and accepts whatever the model opens at. ADR-0106
recorded that trade honestly. Measured a day later, the trade is worse than it
looked.

Three of the four candidate models open at the **bottom** of their own advertised
range — `grok-4.5` at `low` of `low, medium, high`; `kimi-k3` and
`deepseek-v4-flash` at `low` of `low, high, max`. Only `deepseek-v4-pro` opens
above the floor, and only because it advertises no floor. Meanwhile the
published figures a maintainer selects on describe the top: the benchmark
variants are `Kimi K3 (max)` and `Grok 4.5 (high)`.

So a maintainer picks a model on its 93rd-percentile score and Roundfix runs it
at `low`. That is not a conservative default; it is a different model.

This is a bug fix in effect and a capability change in shape. It changes no
product promise — `CONTEXT.md` already defines an Agent Selection as an exact
runtime, model, and reasoning-effort tuple, and this Spec makes the third
component real on a runtime where it was inert.

## Project Constraints

- Identifier strategy: not applicable — no new entity or resource identity. ACP
  config option ids and effort values are assigned by the adapter; Agent
  Selection tuples already exist. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the ACP Runtime boundary is local
  stdio through acpx, with no network surface or HTTP contract of Roundfix's
  own. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0108 is the decision this Spec
  implements and supersedes ADR-0106, which must be left superseded rather than
  edited away; ADR-0105 keeps capability retention bounded by relevance and is
  unaffected; ADR-0107 keeps readiness over configured categories and is
  unaffected; ADR-0050 keeps Fallback Chains inactive until after Run creation,
  so Preflight still proves every configured tuple and substitutes none;
  ADR-0069 pins Baseline semantic analysis to a Codex-only retry pair behind
  Exact Agent Selection Proof and is unaffected, because this Spec changes only
  the `opencode` path and leaves Codex proof and application untouched;
  ADR-0091, ADR-0096, and ADR-0097 govern the authored QA gate and its
  evidence; ADR-0104 requires at least one acceptance row to rest on evidence
  this Spec did not author, and now blocks pull request preparation until that
  row is satisfied or carried forward. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the roundfix skill documents the refusal this
  Spec removes and must be updated with it, and this repository's Roundfix
  configuration carries the profile that proves the result. Express maintainer
  authorization:
  `docs/workflow/authorizations/2026-08-09-an-effort-the-opencode-runtime-actually-receives.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md` and `.roundfixrc.yml`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A maintainer can configure a non-empty `reasoning_effort` on the `opencode`
   runtime, and the Agent actually receives it.
2. Preflight stays token-free and never claims to have applied what it only
   saw advertised.
3. Every work turn of a Run runs at the requested effort, including the first.
4. An effort the selected model does not advertise still fails closed, before
   any Run starts.
5. Codex and Claude selection, proof, and application are unchanged.

## Core Features

- **Session warm-up.** After ensuring the Agent Session with its model, Roundfix
  issues one minimal prompt to raise the queue-owner process, applies the
  requested effort, and observes the effective value before sending work. See
  ADR-0108.
- **A proof split across two moments, each honest.** Preflight proves the model
  is advertised and current and that the requested effort is among the values
  that model advertises. The Run proves the applied value. The new
  `runtime_deferred` encoding names exactly this, kept apart from `independent`,
  which applies and observes token-free.
- **A configuration that stops lying.** `reasoning_effort` on `opencode` is
  accepted instead of refused, and an unadvertised value fails at proof rather
  than at a Run.

## Non-Goals / Out of Scope

- Changing what any other runtime does. Codex and Claude keep the `independent`
  path unchanged.
- Choosing which model or effort this repository should route to. That is a
  configuration decision the maintainer makes once the capability exists, and it
  wants measurement this Spec does not produce.
- A Model Catalog or effort picker for `opencode` in the interactive surfaces.
  Still deferred, as in Spec 0088.
- Proving that a higher effort improves Roundfix's outcomes. The adopted
  measurement is explicit that it establishes only the gap between the level
  selected and the level delivered.

## Decisions

- Roundfix warms the session rather than applying the effort after the first
  real prompt, because the opening turn is the one that most decides a Batch and
  would otherwise be the only turn at the floor. See ADR-0108.
- The warm-up costs a round trip, not a prompt's tokens: the system-prompt cache
  write happens on whichever prompt comes first.
- ADR-0106 is superseded, not deleted. Spec 0088 is archived and stays
  byte-identical.
