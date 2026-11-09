---
spec: 0365-the-capability-projection-reads-grouped-select-options
status: active
created: 2026-11-09
surfaces: [cli]
type: fix
---

# The capability projection reads grouped select options

Since codex-acp 2.6.0 (published 2026-11-08) every Codex profile fails its Exact Agent Selection Proof, so `roundfix doctor` reports `profiles: failed` with `capability_evidence_invalid` and no Run with a Codex profile can start on a machine that updated the adapter. The adapter now advertises its `mode` select as groups of values, a shape the Agent Client Protocol defines beside the flat list, and Roundfix reads only the flat list: each group entry fails as `invalid_select_option` and the current value is never found. After this Spec the capability projection reads a grouped select as its groups' values in advertised order, under the same bounds as a flat select, so the proof passes on codex-acp 2.6.0 as it does on 2.5.0, a full-access request still finds `agent-full-access` inside its group, and every payload outside the protocol's two shapes is still refused.

## Sources

One operator record and two captured payloads, adopted at birth (see
`references/_index.md`):

1. [`references/2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md`](references/2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md),
   recorded by the operator, `type: fix`.
2. [`references/codex-acp-2.6.0-config-options.json`](references/codex-acp-2.6.0-config-options.json),
   the `config_options` of `acpx --format json --json-strict codex sessions
   show` with codex-acp 2.6.0, recorded by the operator.
3. [`references/codex-acp-2.5.0-config-options.json`](references/codex-acp-2.5.0-config-options.json),
   the same probe with codex-acp 2.5.0, recorded by the operator.

At the authoring base `docs/backlog/` holds one open entry, the strict check
naming an Unreachable Acceptance it cannot read, which Spec 0367 adopts; it
shares no context with this Spec. `docs/findings/` holds no unresolved
Finding about capability evidence.

## Prerequisites

None in `requires`. Authored against `2e59aebe` (v0.137.1). The queue
delivers this Spec first, then 0366, then 0367, with a patch release after
each merge. Spec 0366 owns the QA gate evidence binding and Spec 0367 owns
`internal/speccheck`; neither changes `internal/agent`. This Spec changes no
shared file: no guide, skill, `CONTEXT.md`, Behavior Surface Record,
version record or `internal/cli/cli_test.go`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the projection keeps the option identifiers and values the
  adapter advertises. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the change parses a local `acpx` JSON payload. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0375 (new) records the rule.
  ADR-0105: "A capability projection bounds what it keeps, not what a
  runtime advertises", whose retention and fail-closed rule now apply to the
  flattened values. ADR-0293: "A full-access request selects the mode the
  Agent Session advertises", whose advertised modes are now the flattened
  values of a grouped `mode` option. ADR-0011: "ACP Runtimes can run Batches
  with explicit full access", cited by ADR-0293, unchanged. ADR-0314: "A Light Tier that cannot serve the Run yields to the category
  profile", cites ADR-0011 and does not apply: no Light Tier planning
  changes. They are
  reached only through citations and do not apply either: ADR-0335 ("A stop
  that is not the work's fault does not cost the Run"), ADR-0274 ("A turn
  that did not finish is not accepted"), ADR-0276 ("A row the authorization
  keeps offline qualifies, and a byte invariant compares with the Delivery
  Base") and ADR-0295 ("Reachability reads the fetched default branch, and
  origin names it when origin/HEAD is unset"); no Run outcome, turn
  acceptance, QA row rule or reachability rule changes. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA answers of 2026-11-09 quoted in the authorization record; no
  Makefile, workflow, Go module, lint or formatter change and no Governed
  Path. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0365-the-capability-projection-reads-grouped-select-options/_authorization.md`;
  bounded files: none (`paths: []`).

## Goals

- The Exact Agent Selection Proof of a Codex tuple passes over the codex-acp
  2.6.0 payload, with and without a full-access request, and reports the
  same assignment it reports over the 2.5.0 payload.
- A full-access request on codex-acp 2.6.0 sets `agent-full-access`, the
  candidate the adapter lists inside its `sandbox-modes` group.
- A flat select projects exactly as it does today.
- Every payload outside the protocol's flat and grouped shapes, and every
  bound a flat select is held to, still refuses the capability evidence.

## User Stories

1. As a maintainer who updated codex-acp, I want `roundfix doctor` and
   `roundfix profiles validate` to prove my Codex profiles again, so that my
   Runs start without pinning an older adapter.
2. As a maintainer who runs with `defaults.agent_full_access: true`, I want
   the full-access mode the adapter advertises inside a group to be chosen,
   so that the Run keeps full access instead of a degraded policy.
3. As an adopter, I want a malformed adapter payload to keep failing with a
   named issue, so that a protocol change is never hidden behind a silently
   weaker session.

## Core Features

1. The capability projection flattens a grouped select option's values, in
   advertised order, and checks them with the bounds of a flat select.
2. The recorded 2.6.0 and 2.5.0 payloads become test fixtures that pin the
   projection and the proof.
3. A mixed, nested or unidentified group still refuses the evidence, and an
   empty group contributes nothing.

## User Experience

On a machine with codex-acp 2.6.0, `roundfix doctor` prints the `profiles`
line it prints today on 2.5.0, and `roundfix profiles validate` passes every
Codex tuple. When full access is requested, the effective policy reads
`full-access`. No output text changes for any payload that passes today.

## Non-Goals / Out of Scope

- Naming the offending option in a capability issue; the issue vocabulary
  and its bound stay as they are.
- The ACP protocol version 2 field names (`configId`, `groupId`), which
  Roundfix does not read.
- Choosing a Codex permission profile such as `permission-profile:localdev`
  for full access; the full-access candidates stay `agent-full-access`, then
  `full-access` (ADR-0293).
- Guides, skills and `CONTEXT.md`: the profiles reference already says the
  candidate must be advertised in the `mode` select option, which stays true.
- The Makefile, the CI workflows, `go.mod`, `go.sum` and the lint or
  formatter configuration.

## Success Metrics

1. Success Metric: the projection of the recorded 2.6.0 payload succeeds
   with the `mode` values `read-only`, `workspace-write`, `agent`,
   `agent-full-access`, `permission-profile:localdev` in that order, current
   value `permission-profile:localdev`; at `2e59aebe` it fails with
   `invalid_current_value, invalid_select_option`.
2. Success Metric: the projected models, current model, reasoning option and
   every option other than `mode` are equal between the 2.6.0 and 2.5.0
   payloads, and every existing capability test passes unchanged.
3. Success Metric: the Exact Agent Selection Proof through the fake `acpx`
   passes over both payloads, sends `set-mode agent-full-access` under a
   full-access request and no `set-mode` without one; at `2e59aebe` it fails
   over the 2.6.0 payload with `capability evidence invalid:
   invalid_current_value, invalid_select_option`.
4. Success Metric: a mixed list, a nested group, a group without or with an
   unbounded identifier, a value repeated across groups, an unadvertised
   current value, an option whose groups are all empty and a flattened count
   above the ceiling each refuse with its named issue; an empty group beside
   a non-empty one is accepted.

## Acceptance evidence

Two outside sources this Spec did not write back its acceptance:

- The two payloads recorded by the operator on 2026-11-09 from codex-acp
  2.6.0 and 2.5.0, kept in `references/` and copied byte for byte into the
  test fixtures: a measurement the Spec did not design.
- The Agent Client Protocol, protocol v1, Session Config Options
  (https://agentclientprotocol.com/protocol/v1/session-config-options, read
  through Exa on 2026-11-09): "This is either a flat array of values or an
  array of groups. Values and groups are not mixed in the same array." and
  "Each `SessionConfigSelectGroup` requires `group` (its unique string
  identifier), `name` (its human-readable label), and `options` (an array of
  `SessionConfigSelectOption` objects)." The schema's Rust and TypeScript
  bindings declare the same union, `Ungrouped(Vec<SessionConfigSelectOption>)`
  and `Grouped(Vec<SessionConfigSelectGroup>)`.

## Unreachable Acceptance

- criterion: after the patch release that carries this Spec, `roundfix
  profiles validate` on the operator's machine with codex-acp 2.6.0
  reinstalled proves every Codex tuple and reports no
  `capability_evidence_invalid`.
  reason: it needs the released binary, the real codex-acp 2.6.0 adapter, a
  Codex login and the network, which no Verification or QA row may start or
  reach; the operator pinned codex-acp 2.5.0 on the machine until the fix
  ships.
  satisfied-by: the operator's `roundfix profiles validate` run with
  codex-acp 2.6.0 after the release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-09 against these sources.

- Characterization at `2e59aebe`: a scratch test fed the recorded 2.6.0
  `config_options` to `ParseSessionCapabilitySnapshot` and got `capability
  evidence invalid: invalid_current_value, invalid_select_option`; the 2.5.0
  payload projected the `mode` values `read-only`, `workspace-write`,
  `agent`, `agent-full-access`. Through the fake `acpx` harness,
  `ProveExactSelection` for `codex gpt-6.1-sol high` failed with the same
  error over 2.6.0, with and without full access, and passed over 2.5.0 with
  `set-mode agent-full-access`.
- A scratch prototype of the flattening, removed afterwards, made both
  proofs pass over 2.6.0 with `set-mode agent-full-access`, and the whole
  `./internal/agent` package passed with it.
- Reading at `2e59aebe`: `projectSelectCapability` is the only reader of a
  select's `options`; `applyFullAccess` (`internal/agent/acpx_runner.go`)
  chooses the full-access mode from the projected `Modes.Values`, and reads
  `sandbox_mode` by presence among the projected options, so the `mode`
  option is assigned whenever full access is requested.
- Pending ai-memory messages for this project: not read by this author; the
  operator triages them.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "codex-acp grouped select options ACP config options mode" --all --files --min-score 0.3`;
  the hits were this repository's own mirrored ADR-0011, a reference and the
  changelog, and `wiki/concepts/agent-harnesses.md`; no page decides how a
  grouped select is read, so the protocol and the repository's ADRs govern.
- Exa: the ACP v1 and v2 Session Config Options pages and the schema
  bindings quoted above. Version 2 also says "Clients SHOULD hide a group
  that ends up empty", which ADR-0375 follows; the codex-acp sources found
  did not say whether the adapter omits an empty permission-profile group.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BA brief for Spec 0365 of 2026-11-09 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-09 ~08:35–08:45 BRT),
verbatim:

- Question "Defeito novo e bloqueante: o codex-acp 2.6.0 (08/11, instalado
  hoje pelo dev_upgrade) manda a opção `mode` em grupos ACP; o Roundfix só lê
  opções planas e o preflight recusa TODO perfil codex (doctor: profiles
  failed, capability_evidence_invalid). Afeta todos os adotantes que
  atualizaram. Como trato?" → "Spec C hotfix primeiro (Recommended)" —
  option text: "Spec de correção: a projeção de capacidades lê select options
  agrupadas (achata os grupos), com teste do payload real do 2.6.0. Entra na
  frente da B e da A e sai como patch v0.137.2 logo após o merge."
- Rules: question "Que autorizações você concede de saída para eu levar o
  ciclo BA até o fim sem parar (implementações, merges e releases inclusos
  como no AZ)?" → "Mesmas do AZ (Recommended)" — "Release patch assim que
  cada Spec ficar pronta, teto de 2 corretivas por Spec, merge pela fila,
  avisos ao Jarvis, sem notas de release às caixas."; "Scripts/Makefile/CI" —
  "Mexer em ~/.roundfix-operator, Makefile e .github/workflows se a migração
  ou a Spec B precisarem."; "go.mod/lint" — "Mexer em go.mod/go.sum e na
  configuração de lint/formatter se uma Task precisar."; "Avisar
  Quaestor/Fiscus" — "Mensagem ai-memory curta a cada um dizendo qual Spec
  adotou o finding (não é nota de release)."

Taken from those and the reading, with no further interview:

- Flatten the groups into one option's values, in advertised order, and
  keep every flat-list bound (ADR-0375).
- A malformed option still refuses the whole evidence even when Roundfix
  never assigns it. The `mode` option, which the measured failure carried,
  is assigned whenever full access is requested (ADR-0293), and
  `sandbox_mode` decides the Codex sandbox preset by its presence, so a
  tolerant projection would turn a malformed payload into a degraded or
  unsandboxed session instead of a named refusal (ADR-0375).
- An empty group contributes nothing, following the protocol's guidance to
  hide an empty group; an option with no value at all is still refused.
- The QA proves the fix with the recorded payloads and the fake `acpx`; the
  live check with codex-acp 2.6.0 is the one Unreachable Acceptance entry,
  satisfied by the operator after the release.
- The release is the operator's patch v0.137.2 right after the merge; no
  Task releases.

## Open Questions

None.
