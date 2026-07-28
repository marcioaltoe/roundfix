---
spec: 0052-claude-adapter-standardization
status: active
created: 2026-07-28
surfaces: [backend, cli, docs]
---

# Official Claude adapter and opaque model identifiers

Roundfix still resolves the claude ACP Runtime to the deprecated
`@zed-industries/claude-code-acp`, which npm renamed to
`@agentclientprotocol/claude-agent-acp` — the same rename Roundfix already
followed for Codex. The official adapter advertises exactly the model and
reasoning controls Roundfix requires, but Roundfix misreads its advertised
model identifiers (`opus[1m]`) as an embedded reasoning-effort encoding, so
the identifiers Roundfix itself prints are unselectable, a context window is
silently accepted as a reasoning effort, and Adapter Readiness proves lineage
for Codex only. The maintainer directs Roundfix to support only the official
Claude adapter, including in the Doctor and Setup Commands. Evidence and
verification live in the
[standardization finding](../../findings/2026-07-27-claude-adapter-standardization.md)
and its
[predecessor](../../findings/2026-07-26-claude-adapter-configoptions-migration.md).

## Project Constraints

- Identifier strategy: not applicable — this Spec creates no project-owned
  Internal Identifier; adapter packages, advertised Agent Model identifiers,
  and Agent Selection tuples keep their source-native identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change is local adapter
  probing, selection parsing, configuration, CLI reporting, or documentation;
  no authentication provider, credential policy, or HTTP route changes.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0037 and ADR-0039 keep Agent Model
  and reasoning selection Roundfix-owned and proven through disposable Agent
  Sessions; ADR-0049 keeps Agent Selection Profiles atomic; ADR-0050 keeps
  Fallback Chain activation post-Run-creation and notification-first;
  ADR-0055 makes advertised ACP capabilities and adapter lineage the proof
  boundary and forbids silent override edits; ADR-0079 makes advertised model
  identifiers opaque when an independent reasoning control exists and
  supersedes ADR-0040's Claude premise while preserving its explicit-empty
  semantics. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout of
  that edit in exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- The Doctor and Setup Commands prove official Claude adapter lineage the way
  they prove Codex, and a deprecated or foreign lineage fails Adapter
  Readiness with the official install action instead of failing later as a
  capability-evidence error.
- Every Agent Model identifier the adapter advertises is selectable, and a
  context-window annotation is never accepted as a reasoning effort.
- A fresh installation's built-in `frontend` profile proves against the
  official Claude adapter without hand-editing configuration.
- No supported source, test, or documentation path names the deprecated
  `@zed-industries/claude-code-acp` lineage.

## User Stories

1. As a user whose claude adapter command resolves to a deprecated lineage, I
   want the Doctor Command to fail Adapter Readiness naming the effective
   command, package classification, and official install action, so that I
   migrate before a Run fails mid-preflight with a misleading
   capability-evidence error.
2. As a user with a stale Claude adapter override, I want the Setup Command to
   diagnose it and propose migration to the official adapter — never editing
   it silently — so that my local adapter configuration reaches the supported
   lineage the same way the Codex migration works.
3. As a user configuring a profile with an identifier copied from Roundfix's
   own advertised-models diagnostic, I want that exact identifier to prove and
   run with my explicit reasoning effort, so that the diagnostic never
   advertises unselectable values.
4. As a user who mistakenly configures a context-window annotation as a
   reasoning effort, I want the selection rejected with the advertised
   reasoning efforts named, so that a context window cannot masquerade as a
   reasoning level.
5. As a new user on default profiles, I want the built-in `frontend` profile
   to prove against the official Claude adapter, so that frontend Tasks run on
   a fresh installation.
6. As an operator whose operational Preflight Validation refuses on a broken
   Preferred Selection, I want the refusal to state that Fallback Chains
   activate only after Run creation and to name the failing tuple and next
   action, so that I am not left guessing why a configured fallback did not
   rescue the Run.

## Core Features

1. The claude ACP Runtime has a complete adapter lineage contract mirroring
   Codex: official package `@agentclientprotocol/claude-agent-acp`, pinned
   minimum version `0.63.0`, recognized legacy lineages
   `@zed-industries/claude-code-acp` and `@zed-industries/claude-agent-acp`,
   and a default adapter command that resolves to the official pinned form. A
   matching executable name alone is never Adapter Readiness proof.
2. Adapter Readiness for claude reports the effective command, resolved
   package, and version as evidence; classifies a legacy lineage as a
   migration with the official install action; and rejects a version below
   the pin with an install command. The Doctor Command, Setup Command, and
   operational Preflight Validation share this contract.
3. Every adapter install hint Roundfix prints names only official
   `@agentclientprotocol` packages; the wrong-scope hint that resolves to
   `@zed-industries/claude-agent-acp` is corrected.
4. When an adapter advertises an independent reasoning-effort control,
   advertised Agent Model identifiers are opaque: `opus[1m]` is selectable
   with an explicit reasoning effort, the canonical prefix `opus` continues to
   resolve to that advertised entry, and a bracketed annotation is never
   parsed as a reasoning effort, so `opus` with reasoning effort `1m` is
   rejected naming the advertised efforts. Adapters that advertise no
   independent reasoning control keep the existing `canonical[effort]`
   variant encoding.
5. Selection failure diagnostics render selectable identifiers: when an
   advertised identifier differs from its canonical form, the message shows
   both, so any value copied from the diagnostic proves cleanly.
6. The built-in `frontend` profile's Preferred Selection becomes
   `claude / opus / xhigh`; its Fallback Chain is unchanged.
7. The pinned minimum Codex adapter version rises to `1.1.5`, and the Setup
   Command proposes that form.
8. An operational Preflight Validation refusal caused by a failed selection
   proof names the failing tuple, states that Fallback Chains activate only
   after Run creation per ADR-0050, and names the next action.
9. Operator documentation and the Roundfix Skill pair describe the Claude
   adapter contract the way they describe Codex, and no supported source,
   test, or documentation path outside `docs/findings/` and archived Specs
   references `@zed-industries/claude-code-acp`.

## User Experience

- `roundfix doctor` reports Claude adapter evidence with package and version
  the way it reports Codex, and a failed check prints the official
  `npm install -g` next action.
- `roundfix setup` on a machine holding a deprecated Claude override offers
  the migration with the diagnosed lineage; decline preserves every byte.
- A selection failure lists advertised Agent Models in selectable form and
  advertised reasoning efforts, and never prints an identifier that would be
  rejected if copied back.
- Preflight refusals distinguish "this tuple failed proof" from "no fallback
  was consulted", naming the ADR-0050 boundary and the next action.

## Non-Goals / Out of Scope

- Changing Fallback Chain activation semantics — ADR-0050's
  post-Run-creation, notification-first activation stays; this Spec improves
  only the refusal diagnostic.
- Adding adapter aliases to the Model Catalog allowlist — explicit custom
  model strings continue to be sent verbatim for Exact Agent Selection Proof.
- Any Codex adapter behavior change beyond raising the pinned minimum
  version.
- Editing user-owned acpx configuration without diagnosis and confirmation,
  per ADR-0055.
- Continued support for either `@zed-industries` lineage — both fail Adapter
  Readiness with migration guidance.
- OpenCode adapter lineage proof, which has no official pinned contract yet.

## Success Metrics

- The Doctor Command on a machine whose claude command resolves to either
  `@zed-industries` lineage fails Adapter Readiness with the official install
  action, and with the official adapter at or above the pin reports package
  and version evidence.
- `claude / opus[1m] / xhigh` and `claude / opus / xhigh` both prove and apply
  `xhigh` through the independent reasoning control; `claude / opus / 1m` is
  rejected naming the advertised efforts.
- A fresh-install `frontend` profile proves against the official adapter with
  no configuration edits.
- An adapter that genuinely uses the `canonical[effort]` variant encoding
  keeps proving through the existing path.
- Repository-wide search finds no supported path outside `docs/findings/` and
  archived Specs naming `@zed-industries/claude-code-acp`.

## Decisions

- Only `@agentclientprotocol/claude-agent-acp` is supported for the claude
  ACP Runtime; both `@zed-industries` lineages are recognized legacies that
  fail with migration guidance. Maintainer directive, 2026-07-27.
- The Claude adapter pin is `0.63.0`, the version proven empirically on
  2026-07-27; the Codex pin rises to `1.1.5`, the acpx floor.
- The built-in `frontend` Preferred Selection becomes `claude / opus / xhigh`,
  the proven tuple that yields the 1M-context Opus with explicit `xhigh`.
- Advertised model identifiers are opaque when an independent reasoning
  control exists. See
  [ADR-0079](../../adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md).

## Open Questions

None.
