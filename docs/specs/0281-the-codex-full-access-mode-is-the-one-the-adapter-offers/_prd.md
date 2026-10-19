---
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
status: active
created: 2026-10-19
surfaces: [backend, cli, docs]
---

# The Codex full-access mode is the one the adapter offers

A bug fix, minted by `write-techspec` with this minimal PRD.

`defaults.agent_full_access: true` and `--agent-full-access` ask the ACP
Runtime for its full-access mode. For Codex, Roundfix sends the fixed mode
`full-access`. codex-acp 2.2.2 refuses that mode with ACP -32602 and accepts
`agent-full-access`, so every Agent Session that requests full access fails
before any work. The Codex sandbox preset is skipped too, because it is sent
only after a mode literally named `full-access`. The `.roundfixrc.yml` comment
that tells the operator to keep the setting false still cites codex-acp
1.1.9.

After this Spec:

- Roundfix selects the full-access mode from the modes the Agent Session
  advertises, by a fixed preference per runtime.
- It sends the Codex sandbox preset on that path whenever the session offers
  it.
- When the session advertises no full-access mode, Roundfix reports a
  degraded effective access policy that names the advertised modes, instead
  of failing silently or trying a mode the adapter does not have.
- The `qa-gate` skill tells a gate whose daemon-backed browser tool cannot
  start inside the sandbox to drive the browser through the chrome-devtools
  MCP server.

## Sources

One Secondbrain Inbox Entry, captured manually from the Conexus repository on
2026-10-18, is adopted verbatim into this Spec's `references/` (see
`references/_index.md`). The lead moves the Inbox Entry to `_triaged/` after
this plan merges; this Spec never writes to the Secondbrain.

1. `references/2026-10-18-agent-full-access-pede-um-modo-que-o-codex-acp-recusa.md`,
   from `inbox/roundfix/2026-10-18-agent-full-access-pede-um-modo-que-o-codex-acp-recusa.md`.
   - Conexus Spec 0002, Roundfix 0.87.0, codex-acp 2.2.2, acpx 0.19.6: the QA
     gate blocked every browser row three times, because the agent-browser
     daemon cannot start inside the Codex sandbox.
   - A direct `acpx codex set-mode` test refused `full-access`,
     `danger-full-access`, `full-auto`, `yolo`, `auto` and `default` with ACP
     -32602, and accepted `agent-full-access` and `read-only`.
   - Changing `~/.codex/config.toml` to `danger-full-access` changed nothing:
     the next QA report was byte-identical and cost 10.3 million tokens.
   - Its addendum of 2026-10-18 16:15 BRT: the chrome-devtools MCP server runs
     outside the Codex sandbox and drove the browser rows, so only the
     agent-browser tool is blocked.

No open Backlog Entry or unresolved Finding shares this context:
`docs/backlog/` holds only `2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`
(Spec 0282) and the two Git Locator and storage report entries (Spec 0283),
and `docs/findings/` holds none.

## Prerequisites

None. Specs 0282 and 0283, authored beside this one, change other files, and
this Spec delivers first in the J → K → L order.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one effective access policy text and one Run event
  status are added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the change reads the `acpx sessions show`
  output Roundfix already reads, and every test uses the fake acpx runner and
  temporary homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the full-access selection follows
  recorded decisions. ADR-0293 (this Spec): "A full-access request selects
  the mode the Agent Session advertises". ADR-0011, which it amends: "ACP
  Runtimes can run Batches with explicit full access". ADR-0017: "Roundfix
  drives ACP Runtimes through acpx". ADR-0189: "An owned skill's version names
  its content, and its minimum is the bundle". ADR-0285: "Every test file
  compiles for Windows, and a repository contract proves it". ADR-0125: "A
  spawned fixture is a compiled binary, never a written script". ADR-0104: "A
  Spec accepts on evidence it did not author". Set no obligation here, because this Spec changes none of their rules: ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0256 (a release waits for the skills that describe a changed surface); ADR-0276 (a row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ approval of 2026-10-19 through AskUserQuestion, option "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0281-the-codex-full-access-mode-is-the-one-the-adapter-offers/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/profiles.md`, `.roundfixrc.yml`,
  `internal/cli/cli_test.go`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go
  module, lint, formatter or Baseline module change, and
  `defaults.agent_full_access` stays `false` in `.roundfixrc.yml`.

## Goals

- A full-access request on codex-acp 2.2.2 opens the Agent Session in
  `agent-full-access`, and the proof records the effective access policy
  `full-access`.
- An adapter that still advertises `full-access` and a `sandbox_mode` option
  keeps today's mode and sandbox preset.
- Claude Code keeps `bypassPermissions`.
- A session that advertises no full-access mode is reported degraded with the
  advertised modes, and no mode it lacks is ever sent.

## User Stories

1. As an operator who sets `defaults.agent_full_access: true` with Codex, I
   want the session to open with the adapter's own full-access mode, so that
   a Run does not fail every Task before any work.
2. As the operator of an adapter that offers no full-access mode, I want the
   proof, the Doctor and the Run events to name the modes it offered, so that
   I know why full access is missing.
3. As a QA gate inside a sandboxed Codex session, I want to know which
   browser driver still works, so that browser rows are not blocked when a
   driver outside the sandbox exists.

## Core Features

1. Full-access mode discovery.
   - Each runtime names its full-access candidates in preference order:
     Codex `agent-full-access`, then `full-access`; Claude Code
     `bypassPermissions`.
   - The advertised modes are the values of the session's `mode` config
     option, read from the `acpx sessions show` output Roundfix already
     reads.
   - The first candidate the session advertises is sent with `set-mode`, in
     the disposable proof and in the Run.
2. The Codex sandbox preset on the discovered path.
   - After the mode is set on a Codex session, `set sandbox_mode
     danger-full-access` is sent when the session advertises `sandbox_mode`.
   - When it does not, no preset is sent and the policy is not degraded.
3. The degraded report.
   - No advertised candidate: no `set-mode`, the effective access policy
     `full-access (degraded: no advertised mode grants full access; advertised
     modes: <values>)`, the Run event status `full_access_mode_not_advertised`
     and a warning naming the Agent Session.
4. Documentation.
   - ADR-0011 gains an amendment note, and the configuration guide, the
     Roundfix `profiles` reference and the glossary describe the selection.
   - The `.roundfixrc.yml` comment says why the setting stays false here.
   - The `qa-gate` skill gains a heading on the browser driver inside a
     sandbox.

## Non-Goals / Out of Scope

- Turning on `defaults.agent_full_access` in this repository.
- Reading the older `modes.availableModes` field, or calling
  `session/set_config_option` for the mode.
- Choosing a mode by name or description outside the candidate lists.
- Adding full access to OpenCode, Cursor or a command override.
- Changing the pinned adapter versions, the selection proof of model and
  reasoning effort, or how a refused advertised mode fails.
- Fixing the agent-browser daemon inside the sandbox.

## Success Metrics

1. Success Metric: with the fake acpx runner advertising `read-only`,
   `workspace-write`, `agent` and `agent-full-access` and no `sandbox_mode`,
   a Codex full-access proof sends `set-mode agent-full-access`, sends no
   sandbox preset, and records the effective access policy `full-access`.
2. Success Metric: with the fake advertising `full-access` and a
   `sandbox_mode` option, a Codex Run sends `set-mode full-access` and then
   `set sandbox_mode danger-full-access`. With the fake advertising
   `bypassPermissions`, a Claude Run sends `set-mode bypassPermissions`.
3. Success Metric: with the fake advertising only `read-only`,
   `workspace-write` and `agent`, the proof sends no `set-mode` and records
   `full-access (degraded: no advertised mode grants full access; advertised
   modes: read-only, workspace-write, agent)`. The Run publishes
   `full_access_mode_not_advertised` and continues.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write,
read on 2026-10-19.

- agentclientprotocol/codex-acp, `README.md`, Runtime options
  (https://github.com/agentclientprotocol/codex-acp/blob/main/README.md):
  "`INITIAL_AGENT_MODE` - initial mode id: `read-only`, `workspace-write`,
  `agent`, or `agent-full-access`."
- agentclientprotocol/codex-acp, `src/AgentMode.ts`
  (https://github.com/agentclientprotocol/codex-acp/blob/main/src/AgentMode.ts):
  the mode `agent-full-access`, named "Full access", has the approval policy
  `"never"` and the sandbox mode `"danger-full-access"`. The adapter
  advertises its modes as a config option with `id: MODE_CONFIG_ID` (`"mode"`),
  `category: "mode"` and the description "Approval and sandboxing preset for
  the session".
- Agent Client Protocol, "Session Config Options"
  (https://agentclientprotocol.com/protocol/v1/session-config-options): "If an
  Agent provides `configOptions`, Clients SHOULD use them instead of the
  `modes` field. Modes will be removed in a future version of the protocol."
  The same page says categories "MUST NOT be required for correctness".

## Validation

The decisions were checked on 2026-10-19 against these sources.

- Secondbrain.
  - Read `wiki/index.md`, then ran two `qmd query` runs from the Secondbrain
    root:
    - "codex-acp full-access mode agent_full_access sandbox"
    - "agent-browser chrome-devtools MCP sandbox QA gate browser"
  - The first returned the mirror of ADR-0011
    (`projects/roundfix/mirror/docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`),
    which agrees with this repository and names `full-access` for Codex, and
    `wiki/sources/digest-roundfix-claude-agent-sdk-permissoes-e-anexos-2026-10-04.md`,
    which records that acpx passes only some session options to the Claude
    adapter. That supports reading what the session advertises over assuming
    what the adapter accepts.
  - The second returned the agent-browser skill mirror
    (`projects/skills/mirror/skills/06-review-repair/agent-browser/SKILL.md`)
    and no decision on sandboxed browsers, so the `qa-gate` text rests on the
    adopted source's addendum.
- Exa.
  - Supporting: the ACP "Session Config Options" and "Session Modes" pages,
    and the codex-acp README and `src/AgentMode.ts`, quoted above.
  - Challenging: the ACP "v2 migration" page
    (https://agentclientprotocol.com/protocol/v2/migration) removes
    `session/set_mode` and tells Clients to change modes with
    `session/set_config_option`. acpx 0.19.6 still sends `set-mode` as
    `session/set_mode`, so this Spec keeps `set-mode` and records the move to
    a config option as later work if acpx changes.
  - No research Inbox Entry was captured; the sources are the adapter's and
    the protocol's own documentation.
- Measurement. The acpx 0.19.6 session record of a local Codex session holds a
  `mode` config option with the values `read-only`, `workspace-write`,
  `agent` and `agent-full-access`, and no `sandbox_mode` option. A local
  Claude session's `mode` option holds `default`, `acceptEdits`, `plan`,
  `auto` and `bypassPermissions`. The record keeps `config_options` and no
  `modes` field.

## Glossary

- adds: **Full-Access Mode**

## Decisions

- 2026-10-19 (about 16:20 BRT): how to decompose the cycle — "Três Specs
  (Recomendado)".
- 2026-10-19 (about 16:20 BRT): how Spec J selects the mode — "Descobrir pelos
  availableModes (Recomendado)", option text "Escolher o modo pela lista da
  sessão ACP (hoje 'agent-full-access'), manter o preset sandbox
  danger-full-access, atualizar o comentário do .roundfixrc gerado. NÃO ligo
  agent_full_access neste repo; o QA usa o chrome-devtools MCP."
- 2026-10-19 (about 16:20 BRT): which authority covers the Governed Paths —
  "Mexer em Governed Paths declarados", option text "'Concedo' para os
  caminhos que cada Spec declarar (módulos Baseline, registros, skills,
  testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
  específica)."
- 2026-10-19 (about 16:20 BRT): costs, failures and authoring — "Regras do
  ciclo AH/AI (Recomendado)", "Regras AE–AH e parar na proposta
  (Recomendado)", "Opus + strict check, enfileirar sem esperar você
  (Recomendado)".
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present. Decisions were taken from the lead's
cycle AJ briefing of 2026-10-19 for Spec 0281 and from the code and evidence
characterized in the TechSpec.

Taken from those sources, with no further interview:

- The advertised modes come from the session's `mode` config option, which
  Roundfix already reads, and the preference is an exact list per runtime.
  See ADR-0293.
- A session without an advertised candidate is reported degraded, as the
  briefing asked, instead of refused. See ADR-0293.
- The Codex sandbox preset is sent when the session advertises
  `sandbox_mode`; codex-acp 2.x folds the sandbox into its mode. See
  ADR-0293.
- The `qa-gate` guidance names the chrome-devtools MCP server as the driver
  to prefer when a daemon-backed browser tool cannot start in the sandbox,
  from the adopted source's addendum.

## Open Questions

None.
