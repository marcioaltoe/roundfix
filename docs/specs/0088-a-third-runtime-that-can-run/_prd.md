---
spec: 0088-a-third-runtime-that-can-run
status: active
created: 2026-08-08
surfaces: [backend, cli, docs]
---

# A third ACP Runtime that can actually run

`CONTEXT.md` lists OpenCode as a supported ACP Runtime through `opencode acp`.
Measured on 2026-08-08, no Roundfix Run can start on it. Three independent
defects stand in the way, and the cost is concrete: the Codex quota was
exhausted until 2026-08-12 and the Anthropic weekly limit stood at 66% with two
days to reset, leaving a freshly subscribed runtime with full capacity and no
route to it.

This is a bug fix. It changes no product promise — it delivers one the glossary
already makes.

The three defects, each measured and recorded in
[references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md](references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md):

1. **The capability projection reads a large catalog as malformation.** OpenCode
   advertises 417 model values; the projection discards any option above 64
   values, which leaves no model state and no Exact Agent Selection Proof.
2. **A reasoning effort cannot be applied before the Run's first prompt.** The
   `effort` option is per-model and exists only once a queue owner process holds
   the selected model, which acpx starts on the first prompt. Every `set effort`
   Roundfix issues during a token-free proof hits a transient agent on the
   runtime default and answers ACP `-32602`.
3. **Profile readiness never looks at optional Agent Work Categories.** A
   configured `data` profile that fails still leaves `roundfix doctor` printing
   `profiles: ok`, and its ACP Runtime is never named. This defect cost a
   session of misdiagnosis before the first two were even visible.

## Project Constraints

- Identifier strategy: not applicable — this Spec introduces no new entity or
  resource identity. Agent Selection tuples, Agent Work Category keys, and ACP
  config option ids are all assigned by existing contracts or by the adapter.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the ACP Runtime boundary is local
  stdio through acpx, with no network surface and no HTTP contract of Roundfix's
  own. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0050 keeps Fallback Chains inactive
  until after Run creation, so Preflight still proves every configured tuple and
  substitutes none; ADR-0069 pins Baseline semantic analysis to a Codex-only
  retry pair proven through Exact Agent Selection Proof and is unaffected,
  because this Spec neither changes how a Codex selection is proven nor routes
  that analysis to another runtime; ADR-0081 makes the copies and digests
  rewritten by `make skills-sync` sanctioned fallout of the authorized skill
  edit rather than separate targets; ADR-0091, ADR-0096, and ADR-0097 govern the
  authored QA gate and its evidence; ADR-0104 requires at least one acceptance
  row to rest on evidence this Spec did not author, which the adopted
  measurement supplies. New decisions land as ADR-0105, ADR-0106, and ADR-0107.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the roundfix skill and this repository's
  Roundfix configuration are mutated. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-08-the-third-runtime-gets-a-profile-and-a-skill-entry.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md` and `.roundfixrc.yml`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A configured `opencode` Agent Selection passes Exact Agent Selection Proof
   without weakening what proof means for any other runtime.
2. A Roundfix Run executes a real Task on an `opencode-go` model, with captured
   evidence.
3. A readiness command never reports `ok` while a configured Agent Selection
   Profile is failing.
4. A maintainer who configures an unusable `opencode` reasoning effort is told
   what to write instead, before any Run starts.

## Core Features

- **Relevance-bounded capability retention.** The projection accepts a large
  advertised option and retains the current value, every value the requested
  Agent Selection references, and a bounded diagnostic window, recording the
  retained count against the advertised total. A model that is genuinely not
  advertised still fails as an unsupported Agent Selection, not as invalid
  capability evidence. See ADR-0105.
- **Model-managed reasoning for OpenCode.** Configuration refuses a non-empty
  `reasoning_effort` on the `opencode` runtime and names the empty value as the
  repair, so the proof and the Run both skip the config set that cannot succeed.
  See ADR-0106.
- **Readiness that covers what is configured.** Agent Selection Profile
  Readiness and adapter readiness resolve every Agent Work Category the
  effective configuration defines, not only the required five. See ADR-0107.
- **A reachable route in this repository.** One Agent Work Category selects
  `opencode` in `.roundfixrc.yml`, which is what makes the acceptance Run
  possible.

## Non-Goals / Out of Scope

- A Model Catalog or reasoning-effort choice list for `opencode` in the
  interactive picker or the TUI. Offering 417 adapter-enumerated models through
  a picker is a separate design question, and the repository-configured
  allowlist it probably needs was considered and deferred.
- Changing the Fallback Chain activation contract. ADR-0050 stands; the
  separately filed defect about fallbacks not activating on quota exhaustion is
  not this Spec's subject.
- Recording per-session usage or cost for any runtime. Filed separately on
  2026-08-08 and still open.
- Any change to how Codex or Claude selections are proven or applied.

## Decisions

- Reasoning effort on `opencode` is model-managed and refused when non-empty,
  rather than proven from observed state without applying it or applied after a
  prompt. See ADR-0106.
- The capability cap bounds retention, not advertisement. See ADR-0105.
- Readiness covers every configured category; categories that merely inherit
  `general` contribute no distinct tuple and are not enumerated. See ADR-0107.
- The Spec is implemented by hand rather than by a Roundfix Run, because no ACP
  Runtime has capacity to execute its Task Graph while the runtime this Spec
  repairs is the only one with quota.
