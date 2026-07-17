---
spec: 0041-agent-selection-runtime-readiness
status: active
created: 2026-07-17
surfaces: [cli, backend, docs]
---

# Agent Selection Runtime Readiness

Roundfix now owns Agent Selection Profiles, but the generated defaults and the
effective ACP adapter can disagree about how an exact selection is expressed.
Dogfooding with Codex CLI `0.144.5`, ACPX `0.12.0`, and a setup-generated bare
`codex-acp` override that resolved to the legacy
`@zed-industries/codex-acp` `0.16.0` exposed four connected failures:

- the generated `codex / gpt-5.6-sol / high` Preferred Selection was advertised
  at the model level but rejected when Roundfix sent `reasoning_effort=high`;
- the generated `codex / gpt-5.6-terra / max` fallback was not advertised;
- the canonical `--agent codex` invocation replaced category profiles with the
  legacy `gpt-5.5 / xhigh` tuple and could duplicate the configured fallback;
- Doctor proved only the legacy runtime default while reporting the repository
  ready.

The runtime evidence does not mean Sol lacks configurable reasoning. A second
disposable-session proof with the current official
`@agentclientprotocol/codex-acp` `1.1.4` advertised Sol, Terra, Luna, GPT-5.5,
and an independent `reasoning_effort` option containing `high`, `xhigh`, and
the other supported values. The same session accepted
`gpt-5.6-sol / high`. ACPX `0.12.0` already uses this official adapter as its
built-in Codex command, but Roundfix Setup had masked it with the legacy bare
override. The primary defect is therefore adapter selection and provenance,
followed by Roundfix's lack of one shared exact-capability proof.

Agent Selection Runtime Readiness gives setup, profile management, Doctor, and
operational commands one capability-driven proof boundary. The canonical Codex
default remains `gpt-5.6-sol / high`; Roundfix must never reinterpret that as
model-managed reasoning. If the installed stack cannot prove the exact tuple,
setup refuses to write the profile and names the incompatible component and
next action. A Run may automatically activate a separately proven configured
fallback only under ADR-0050's pre-prompt notification contract.

## Goals

- Preserve `codex / gpt-5.6-sol / high` as the canonical Preferred Selection
  for `general`, `backend`, `qa`, and `review`.
- Replace the unproven generated Terra fallback with the separately proven
  `codex / gpt-5.5 / xhigh` baseline for those profiles.
- Keep the official `gpt-5.6-terra` and `gpt-5.6-luna` identifiers valid for
  explicit profiles and advisory recommendations; the current official
  adapter advertises them, while each configured tuple still requires exact
  proof in the user's effective environment.
- Verify the effective Codex adapter's official package lineage and supported
  version before Setup writes an ACPX override or profile configuration.
- Discover how the installed ACP adapter represents model and reasoning
  controls and prove the requested Agent Selection through those advertised
  capabilities.
- Reuse the same proof result in `roundfix setup`, `profiles configure`,
  `profiles validate`, `doctor`, and operational Preflight Validation.
- Refuse setup configuration when either the Preferred Selection or its
  required distinct Fallback Chain cannot be proven exactly.
- Make profile-driven commands use category profiles by default and reject a
  bare or partial one-Run Agent Selection override.
- Produce actionable diagnostics that name the runtime, model, effort,
  adapter command/version when available, advertised controls, affected
  categories, and next safe action.

## User Stories

1. As a developer running setup, I want Roundfix to prove the generated
   Preferred Selection and fallback before writing them, so setup cannot create
   a configuration that blocks the first Run.
2. As a user selecting Sol at `high`, I want Roundfix to preserve that exact
   intent instead of silently running Sol at its model-managed default.
3. As a user with a newer adapter, I want Roundfix to use the adapter's
   advertised control shape, so a newly supported tuple works without a
   Roundfix allowlist release.
4. As a user with an older adapter, I want a bounded failure naming the missing
   capability and update action, so I do not repeatedly edit valid model names.
5. As a user relying on category profiles, I want ordinary commands to omit
   one-Run overrides, so backend, frontend, QA, and review use their intended
   selections.
6. As a user requesting an override, I want to provide a complete exact tuple,
   so a runtime-only flag cannot silently inherit a legacy model.
7. As a developer running Doctor, I want readiness to cover all effective
   required profiles and fallbacks, so `doctor: ok` predicts operational
   Preflight Validation.
8. As a user configuring only one proven tuple, I want the configurator to
   explain that a distinct proven fallback is still required without changing
   my config or choosing another model for me.
9. As a user with a stale `codex-acp` override, I want Setup and Doctor to name
   the legacy adapter and the supported replacement, so a current ACPX install
   is not incorrectly treated as incompatible.

## Core Features

1. **Adapter identity proof.** Roundfix resolves the effective adapter command,
   proves its official package lineage and supported version, and rejects a
   stale or ambiguous override before profile proof or persistence.
2. **Capability-driven assignment.** A disposable Agent Session supplies the
   authoritative advertised model/configuration controls. Roundfix derives one
   assignment plan for the requested tuple: independent model and reasoning
   controls, an advertised model variant that represents both, or unsupported.
3. **Exact effective proof.** Successful command exit alone is insufficient.
   Roundfix verifies that the session's effective advertised state corresponds
   to the requested model and reasoning. It never reads a private Codex cache
   or runtime-owned configuration as proof.
4. **Shared readiness service.** Setup, profile configuration/validation,
   Doctor, and operational preflight call the same selection prover and receive
   the same structured result, diagnostics, cleanup behavior, and tuple
   deduplication.
5. **Safe generated profiles.** Built-in and generated profiles use Sol/high
   plus GPT-5.5/xhigh for `general`, `backend`, `qa`, and `review`. Setup writes
   no profile until every distinct effective tuple required by the generated
   profiles passes exact proof. Frontend retains its explicit preferred and
   fallback policy but is subject to the same proof.
6. **Official model handling.** Sol, Terra, and Luna remain accepted official
   identifiers and may appear in the top-five advisory ranking. The supported
   official adapter advertises all three; Terra and Luna are not generated
   operational defaults because the accepted product default remains
   Sol/high with GPT-5.5/xhigh fallback.
7. **Complete override grammar.** `--agent`, `--model`, and
   `--reasoning-effort` form an all-or-none one-Run override for Agent-starting
   commands. Providing only a subset exits `2` before proof or mutation and
   explains that omitting all three uses category profiles.
8. **Profile-aware Doctor.** Doctor proves the effective adapter and profiles
   for the required categories, deduplicates shared tuples, reports every
   affected category, and does not substitute the legacy `defaults.agent`
   runtime.
9. **Distinct-fallback diagnostics.** Configuration continues to require a
   non-empty Fallback Chain with no duplicate of the Preferred Selection. When
   only one tuple is proven or authorized, it reports the missing distinct
   proof and leaves User Config or Project Config byte-identical.
10. **Capability diagnostics.** A failed proof differentiates unsupported or
    legacy adapter, model not advertised, reasoning control not advertised,
    model/reasoning variant not advertised, and adapter rejection. It
    identifies the installed commands and versions it actually observed and
    gives one deterministic next action.

## User Experience

The normal profile-driven commands become:

```text
roundfix implement --spec <slug> --qa --detach
roundfix resolve --pr <number> --no-input
roundfix watch --source coderabbit --pr <number> --until-clean
```

A deliberate one-Run override is complete:

```text
roundfix implement --spec <slug> \
  --agent codex \
  --model gpt-5.6-sol \
  --reasoning-effort high
```

`roundfix setup` proves the effective adapter and proposed generated profiles
before writing either scope. When a legacy bare override resolves to
`@zed-industries/codex-acp` `0.16.0`, it fails without changing the target
config and names the supported official adapter update or migration action. On
`@agentclientprotocol/codex-acp` `1.1.4`, Sol/high is proven without replacing
`high` with model-managed reasoning. Future advertised tuples can work without
a Roundfix Model Catalog release when the adapter remains supported.

Doctor adds one deterministic Agent Selection Profile readiness result. It
lists proven tuple/category counts on success. On failure it names the first
failed exact tuple, every category that references it, the relevant advertised
capability evidence, and the update or profile-configuration command. It still
runs the remaining readiness checks.

## Non-Goals / Out of Scope

- Modifying Codex, ACPX, `codex-acp`, or runtime-owned model metadata.
- Reading `~/.codex/models_cache.json` or other private caches as production
  compatibility authority.
- Guaranteeing that a model remains available after successful preflight.
- Automatically selecting a recommendation-ranked model.
- Treating model-managed reasoning as equivalent to an explicitly requested
  `high`, `xhigh`, `max`, or other non-empty value.
- Allowing profiles without a distinct configured fallback.
- Switching Agent Selection after the first Agent prompt begins.
- Changing the benchmark ranking beyond correcting claims that confuse a
  recommendation with proven availability.
- Implementing Repository Skill Set readiness; Spec 0036 owns that independent
  Doctor check.

## Success Metrics

- The generated required Codex profiles contain Sol/high as Preferred and
  GPT-5.5/xhigh as fallback; no generated profile uses Terra/max.
- Setup and Doctor reject a legacy or unproven adapter override and report the
  effective command, package lineage/version, and deterministic next action.
- Setup writes zero config bytes when any generated Preferred Selection or
  fallback is not exactly proven.
- A fake adapter exposing independent reasoning and a fake adapter exposing
  model variants both prove the same canonical selection through one service.
- A fake adapter exposing Sol without a `high` control fails every consumer
  with the same classification and next action.
- The official adapter fixture advertises Sol, Terra, and Luna plus independent
  reasoning values, and proves Sol/high. A legacy fixture reports Terra/Luna as
  unavailable without changing their catalog validity.
- A bare `--agent codex`, bare `--model`, or bare `--reasoning-effort` exits
  `2` without creating a Run, Session, worktree, or config change.
- Omitting all override flags resolves the profile for each Task, QA, or review
  category; a complete triple overrides only the Preferred Selection under the
  documented fallback contract.
- Doctor and `profiles validate` deduplicate and prove the same effective
  tuples and report the same failure categories.
- Every disposable Agent Session is closed on success, rejection,
  cancellation, and malformed capability evidence.

## Decisions

- Use capability-driven exact proof rather than CLI-only validation or a fixed
  version matrix.
- Pin and verify the supported official Codex ACP adapter package; do not write
  a bare PATH override solely because an executable named `codex-acp` exists.
- Keep Sol/high as the desired default; do not replace it with model-managed
  reasoning merely to pass an older adapter.
- Use GPT-5.5/xhigh as the generated general Codex fallback.
- Block setup persistence when the exact generated profile is unsupported;
  operational fallback remains automatic only after notification and only
  before Agent work under ADR-0050.
- Require all three one-Run override fields together; a bare `--agent` is not a
  runtime-only override.
- See [ADR-0055](../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md).
- See the [validation record](references/validation.md) for the Secondbrain,
  primary-source, and disposable-session evidence.

## Dependencies

- This Spec must be implemented before Spec 0036 so Doctor skill readiness can
  append its independent result to the profile-aware Doctor contract without
  recreating selection proof.
- Spec 0039 may rely on the corrected canonical Roundfix Skill commands, but
  its Review Source Evidence and Detached Run outcome design is otherwise
  independent.

## Open Questions

None.
