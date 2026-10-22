---
spec: 0301-the-light-tier-yields-to-the-category-profile
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# The Light Tier yields to the category profile when it cannot serve the requested access

Minimal PRD for a bug fix, minted by `write-techspec`. With
`defaults.agent_full_access: true` and the Light Tier on, every light-eligible
Task fails before any Agent work, because the light runtime, OpenCode, offers
no full-access mode and the Daemon treats that refusal as a Task failure
instead of falling to the category's profile. On 2026-10-21 that left a
Fluxus Run Unresolved, and both of the maintainer's machines now run with the
Light Tier off. After this Spec a light Task whose runtime cannot serve the
requested access skips the Light Tier with a warning and runs on its category
profile, any other refusal of the light session before Agent work moves the
Task to its category profile, and `roundfix doctor` warns about the
combination. The TechSpec holds the design.

## Sources

One Backlog Entry, written by this Spec's planning change directly into
`references/` and adopted at birth (see `references/_index.md`):

1. [`references/2026-10-21-the-light-tier-fails-a-task-under-agent-full-access.md`](references/2026-10-21-the-light-tier-fails-a-task-under-agent-full-access.md),
   type `fix`, written from the Fluxus report of 2026-10-21 (roundfix
   0.104.0, Run `run_20261021T223522Z_9aa628a4520889af` of Fluxus Spec 0113):
   every light-eligible Task failed in about 15 seconds with `requested access
   policy "full-access" failed predicate "runtime_has_access_mode" for runtime
   "opencode"`.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` held the QA refusal report entry and the
documentation contract cost entry, which other cycle AP Specs own, and
`docs/findings/` held none.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0) while
the cycle AO Specs 0299 and 0300 are active, and delivers after them. Neither
changes `internal/daemon/agent_session_owner.go`, `internal/lighttier`,
`internal/cli/doctor.go`, the Light Tier guides or the `setup` and `runtime`
references. Spec 0300 raises the Roundfix skill's version, so this Spec's
documentation Task raises it over whatever the tree holds when it starts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change adds a reason code and an event phase to an existing
  Task event, and a Doctor check line. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; light and category sessions keep their
  existing keys and adapters, and every test uses fake runners, fake runtime
  factories and temporary homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0314: "A Light Tier that cannot
  serve the Run yields to the category profile", which refines the accepted
  decision that a light tier runs low-complexity Tasks on an open model, and
  keeps the decisions that opening an Agent Session is not Agent work, that a
  lost rollout is runtime infrastructure, that ACP Runtimes run Batches with
  explicit full access, and that a full-access request selects the mode the
  Agent Session advertises, all of which the glossary entries Light Tier,
  Fallback Selection and Agent Work Started carry. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0301-the-light-tier-yields-to-the-category-profile/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config change.

## Goals

- A light Task under a full-access request whose light runtime offers no
  full-access mode never fails for that reason: it runs on its category
  profile, with a warning and a `light_tier_skipped` event whose reason code
  is `access_mode_unavailable`.
- A light session refused before Agent work for any other reason moves the
  Task to its category profile, with a warning and an event naming the cause.
- The three existing skip reasons, the selection fallback between light
  models, the escalation after a failed Verification and the rule that no
  fallback follows Agent work keep their behavior.
- `roundfix doctor` warns when the Light Tier is on, full access is
  requested and no light runtime offers a full-access mode.
- The guides, the Roundfix skill references and the glossary describe the new
  skip reason, the yield and the Doctor line.

## Core Features

1. The access skip. While deriving a light Task's profile, the Daemon keeps
   only the light selections whose runtime, built with the Run's runtime
   factory, offers a mode for the requested access policy. When none remains,
   it skips the Light Tier with the reason code `access_mode_unavailable`,
   before the key, spend log and ceiling checks.
2. The pre-work yield. A light selection's failure before
   `AGENT_WORK_STARTED` that is not a selection failure, a lost rollout, a
   provider system error or a cancellation records the light attempt as
   failed, publishes the Task event phase `light_tier_yielded` and the Agent
   Selection fallback notification, prints a warning, and opens the
   category's Preferred Selection.
3. The Doctor line. `light-tier:` is printed after `recommendations:` only
   when the effective light model list is not empty and
   `defaults.agent_full_access` is `true`; it is `warn` with
   `DR-LIGHT-ACCESS-MODE` when no light runtime offers a full-access mode,
   otherwise `ok`, and it never changes the exit code.
4. Documentation. The configuration guide's Light Tier section, the Doctor
   guide, the Roundfix `runtime` and `setup` references and the glossary's
   **Light Tier** entry describe the skip, the yield and the Doctor line.

## Non-Goals / Out of Scope

- Giving OpenCode a full-access mode, or running a light session without the
  requested access.
- Changing `roundfix profiles validate`, which proves configured profiles
  only.
- Changing the selection fallback between light models, the escalation after
  a failed Verification, the Light Spend Log or the ceiling.
- Changing the full-access mode choice from the advertised modes or the
  degraded policy it reports.

## Success Metrics

1. Success Metric: with a fake runner whose session setup applies the
   runtime's access-policy check and the real static runtime capabilities, a
   `low` Task under a full-access request completes on its category's
   Preferred Selection, no light session is opened, one `light_tier_skipped`
   event carries `access_mode_unavailable`, and one warning names the
   runtime and the policy.
2. Success Metric: under the runtime-default policy the same Task still runs
   on the first light model, and with a fake runtime factory that gives only
   the second light model a full-access mode, the Task runs on that model.
3. Success Metric: through `roundfix implement` with
   `defaults.agent_full_access: true` in Project Config and a light model in
   User Config, the Task's work request goes to its category runtime and
   stderr carries the skip warning; the Run is not Unresolved.
4. Success Metric: a light session whose runtime factory fails, whose setup
   fails with an access-policy refusal, or whose prompt fails before any
   Agent output completes on the category's Preferred Selection with one
   `light_tier_yielded` event and one warning; a failure after Agent output,
   and a non-light Task's setup failure, still fail the Task.
5. Success Metric: the skip reasons `key_missing`, `spend_unreadable` and
   `ceiling_reached`, the light dispatch order, the escalation and the
   no-fallback-after-work rule keep their tests passing unchanged.
6. Success Metric: Doctor prints `light-tier: warn (DR-LIGHT-ACCESS-MODE:
   …)` with the configured combination, prints no `light-tier:` line when
   full access is off or the light model list is empty, keeps every other
   line byte-identical, and exits 0.
7. Success Metric: after release, a Fluxus Run with
   `defaults.agent_full_access: true` and the Light Tier back on completes its
   light-eligible Tasks on their category profiles.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- The Fluxus report of 2026-10-21, recorded in the adopted Backlog Entry:
  the refusal text `requested access policy "full-access" failed predicate
  "runtime_has_access_mode" for runtime "opencode"; recovery: disable full
  access or select a runtime that supports it`, which the characterization
  reproduced byte for byte and which the new skip must prevent.
- OpenCode's ACP documentation (https://opencode.ai/v2/docs/cli/acp, read on
  2026-10-22): its session config options list `mode` with the values
  "Visible primary agents.", and its permission requests offer "Allow once",
  "Always allow" and "Reject"; it names no full-access session mode, which is
  why OpenCode's runtime declares no full-access candidate.
- The Agent Client Protocol's Session Modes page
  (https://agentclientprotocol.com/protocol/session-modes): a mode passed to
  `session/set_mode` "Must be one of the modes listed in `availableModes`".

## Unreachable Acceptance

- criterion: Success Metric 7, after release, a Fluxus Run with
  `defaults.agent_full_access: true` and the Light Tier back on completes its
  light-eligible Tasks on their category profiles.
  reason: only Fluxus's repository and the maintainer's machines hold that
  configuration and its Runs, and no Verification or QA row may read another
  repository, reach OpenRouter or a real ACP runtime, or read the real
  `~/.roundfix`.
  satisfied-by: the operator's notice to Fluxus after release and the next
  Fluxus Run with `openrouter.light_models` restored in User Config.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2` with two `go test -overlay` probes that
  wrote nothing to the repository. Through the Daemon's Task cycle, with the
  real static runtime capabilities (`agent.RuntimeFor` with full access) and
  a fake runner whose session setup returns the runtime's access-policy
  check, a `low` Task failed with `Agent failed: requested access policy
  "full-access" failed predicate "runtime_has_access_mode" for runtime
  "opencode"; recovery: disable full access or select a runtime that
  supports it`, after one light session setup and no prompt, with no skip
  event; under the runtime-default policy it completed on the light model;
  with a runtime factory that refuses `opencode`, it failed with the factory's
  error. Through `roundfix implement` with `defaults.agent_full_access: true`
  in Project Config, the light Task's work request reached `opencode` with the
  requested policy `full-access`, which fails the same check. The TechSpec
  records the code it read.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "light tier
  opencode full access agent_full_access fallback category profile" --all
  --files --min-score 0.3` and `qmd query "fluxus light tier falha
  full-access opencode runtime_has_access_mode" --all --files --min-score
  0.3`. Both returned Roundfix's own mirrored documents
  (the mirrored decision records on the light tier, on full access and on
  the advertised full-access mode, and
  `projects/roundfix/mirror/docs/user-guide/configuration.md`) and unrelated
  Fluxus notes; none records a decision on this combination. They confirmed
  that the light tier decision promised a fall to the category profile for a light session
  that cannot start, which this Spec makes true. No Inbox Entry was written,
  because the briefing forbids writing to the Secondbrain.
- Exa: OpenCode's ACP and agents documentation and the Agent Client
  Protocol's Session Modes page above. Challenging: OpenCode's Build agent
  already allows tools and access outside the workspace by default
  (https://opencode.ai/v2/docs/agents/), so a light session might do the work
  without a mode; but that is OpenCode's permission configuration, not a
  mode Roundfix can request and prove, and the full-access decision forbids
  running when the
  requested access cannot be applied.

## Glossary

- changes: **Light Tier**

## Decisions

- 2026-10-22 (questions of 15:25 BRT, answered by about 18:30 BRT): cycle AP
  scope "Correções de adotantes (Recomendado)" among the four chosen, order
  "Correções, follow-up, Python (Recomendado)" (option text: "Segue a ordem
  que você definiu para o follow-up: primeiro o que trava outros projetos
  (Light Tier no Fluxus), depois o bom funcionamento, depois as features."),
  Light Tier "Pular no planejamento + escalar (Recomendado)" (option text:
  "Quando o acesso pedido é full-access e o runtime light não o atende, a
  Task nem escolhe o Light Tier (evento light_tier_skipped, motivo
  access_mode_unavailable). Qualquer outra recusa da sessão light antes do
  trabalho escala para o perfil normal da categoria. O doctor avisa a
  combinação."), and authorization "same as cycle AO".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AP briefing for Spec 0301 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- "The light runtime advertises no mode that serves it" is decided at
  planning from the runtime's declared full-access candidates, the same
  check the session setup applies first, because no session exists yet; a
  light selection whose runtime does offer a mode stays in the profile.
- The access check runs before the key, spend log and ceiling checks: it is a
  configuration fact and needs no file read.
- "Before Agent work" is the existing `AGENT_WORK_STARTED` boundary of
  decision that opening an Agent Session is not Agent work, so a prompt that fails before any Agent output also yields.
- The yield skips the remaining light models, which share the light runtime;
  a selection start failure still falls to the next light model.
- `roundfix profiles validate` is unchanged; Doctor alone warns, from
  configuration, without opening a session.
- One new ADR, ADR-0314, records the rule; the glossary revises **Light
  Tier**.

## Open Questions

None.
