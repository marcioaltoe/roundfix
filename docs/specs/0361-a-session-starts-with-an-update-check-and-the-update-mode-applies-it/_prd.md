---
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
status: active
created: 2026-11-06
surfaces: [backend, cli, docs]
type: feat
---

# A session starts with an Update Check, and the update mode applies it

Requested by the maintainer through jarvis on 2026-11-06: every Roundfix release
reached the adopting repositories as a note in each project's ai-memory inbox,
and that morning thirty of them waited unread beside the real requests between
projects. After this Spec a session in an adopting repository starts with an
Update Check that compares the installed Roundfix with the latest published
release, lists the title of each release between them and reports the
repository's own Baseline plan, reading the network at most once a day and
saying in one line when it cannot. The Roundfix skill gains an `update` Skill
Mode that verifies, shows and applies the update, asking before any step that
deletes files, and a Project Config key chooses per repository between being
told (the default) and applying at session start.

## Sources

One request, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-06-an-update-check-at-session-start-replaces-the-release-notes-in-the-inboxes.md`](references/2026-11-06-an-update-check-at-session-start-replaces-the-release-notes-in-the-inboxes.md),
   `type: feature`, requested by the maintainer through jarvis (ai-memory
   message `01a1b0fb-92e1-7262-bc05-97272b9de2cb`, 2026-11-06 10:10 BRT), with
   the operator's note of 2026-11-06.

At the authoring base `docs/backlog/` holds one other open entry, the release
workflow and a transparency-log 409, which shares no context with this one,
and `docs/findings/` holds none. This session exposes no ai-memory tools, so no
other pending message was read; the brief carried the one this Spec adopts.

## Prerequisites

None in `requires`. Authored against `9a285d77` (v0.136.1). Three Specs deliver
before this one, in the order 0360 (in the Delivery Queue now), 0362 and 0363.
Files this Spec may share with them, kept to narrow edits named by function or
section:

- With 0360: the canonical and mirrored Roundfix `setup` reference (0360 edits
  its `disk-space:` bullet and `DR-DISK-LOW` row; this Spec adds an
  `## Update check` section), the Roundfix `SKILL.md` pair and its version
  record, `CONTEXT.md` and the Behavior Surface Record.
- With 0362: the Baseline catalog derived files (`make baseline-digests`
  outputs, `internal/baseline/module-versions.json`, the catalog snapshots and
  plan goldens, the formatter goldens) and this repository's setup manifest;
  0362 edits the `context-workflow` module, this Spec only `core`'s
  cross-project continuity rule. Each Task regenerates the derived files from
  the tree it runs on, never by hand.
- With 0363: none expected; it changes the Doctor residue check.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted beyond the command's own JSON document and cache file, which name
  their schema version. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint or credential is
  added; the release listing goes through the GitHub CLI the Upgrade Command
  already uses, with the user's own `gh` authentication. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0371 (new) records the decision:
  "A session starts with an Update Check, and the Roundfix skill applies the
  update". ADR-0320 defines the Skill Modes this Spec extends with `update`,
  and its rule of up to three questions per call governs the mode's questions
  and stays unchanged; ADR-0320: "An Agent asks at most three questions per
  call or message". ADR-0354 keeps every announcement in the user's own
  instructions, and the Update Check sends no message and replaces no
  announcement step, so it stays unchanged; ADR-0354: "and Roundfix's code
  sends no message". ADR-0338, which cites ADR-0320, does not apply: it
  governs how the `--follow-up` mode takes a handoff and settles messages,
  and the `update` mode takes no handoff and consumes no message. ADR-0353,
  which cites ADR-0338, does not apply either: it holds Waiting Items back
  from the follow-up list, which the `update` mode never reads. No accepted
  ADR is amended. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 4 answers of 2026-11-06 quoted in the authorization record, which
  grant the Baseline source, the guides and the skills when the design needs
  them; the Governed Paths are the Roundfix skill files tasks 02 and 03
  change, the CLI help test task_02 changes, and the core Baseline module, the
  Standard TypeScript profile and its formatter golden, and this repository's
  managed guide and setup manifest that task_04 changes. No Makefile, workflow, Go module, lint or
  formatter change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A session in an adopting repository learns, before its first task, whether
  a newer Roundfix exists, what each newer release is about, and whether its
  Baseline has a plan, without the operator sending any message.
- The check never holds up a session: it reads the network at most once a day
  per machine, and a lookup that fails is one line.
- One request in the session, the Roundfix skill's `update` mode, upgrades the
  binary, shows the repository's exact Baseline plan and applies it, and asks
  before any step that deletes files or touches uncommitted work.
- Each repository chooses between `notify` (the default) and `auto`.
- The same instruction reaches every adopter through the Baseline, in Claude
  Code, Codex and OpenCode alike.

## Core Features

1. The `update.mode` key: Project Config and User Config accept `notify` or
   `auto`, Project Config winning; any other value is refused at load naming
   the key; absent means `notify`.
2. The Update Check, `roundfix update check [--refresh] [--format <text|json>]`:
   read-only, exit `0` for every report; the installed Roundfix against the
   latest release, the `### ` titles of each newer release from the published
   release notes, this repository's Baseline preview summarized as current or
   as a plan with its counts, and a `Next:` line that names the update mode
   and the configured `update.mode`; the release listing cached for 24 hours in
   Roundfix Home; `--refresh` reads it now; a failed lookup is one line.
3. The `update` Skill Mode: run the check with `--refresh`, upgrade the binary
   when a release is newer, preview the repository's Baseline plan with that
   binary, apply exactly that plan by its Plan Digest, ask before deleting
   files or applying over uncommitted changes, and never commit.
4. The session-start clause: one mandatory clause of the Baseline's core
   cross-project continuity rule tells each session to run the check once
   before its first task, relay it, follow its `Next:` line in `auto`, and skip
   it in a Daemon-assigned turn.
5. The rule in writing: a command guide with the optional harness hook a
   repository may add, the configuration guide, the Roundfix skill's `setup`
   and `modes` references, the glossary and ADR-0371.

## Non-Goals / Out of Scope

- Cancelling the release notes already pending in the inboxes: the operator
  does it after the release that carries this Spec.
- Writing a SessionStart hook into any harness configuration from the
  Baseline or the CLI; the command guide shows one a repository may add.
- Showing the full release notes; the check shows titles and links.
- Committing, pushing or opening a pull request for the applied Baseline plan.
- Changing `roundfix upgrade`, its Version Freshness line on operational
  commands, `roundfix baseline update`, or the npm installation path.
- A per-repository network cache: the release listing is the same for every
  repository on a machine.

## Success Metrics

1. Success Metric: with a fake release listing whose newest releases are
   `0.136.1`, `0.136.0` and `0.135.0` and an installed `0.134.1`,
   `roundfix update check` prints `Roundfix 0.134.1 → 0.136.1 is available
   (update.mode notify).`, a `Changes since 0.134.1:` block with one line per
   newer release, newest first, each with its `### ` titles, the Baseline line
   and the notify `Next:` line, and exits `0`; with `update.mode: auto` only the
   mode and the `Next:` line change. At `9a285d77` the command does not exist
   (`roundfix: unknown command "update"`, exit `2`).
2. Success Metric: two checks within 24 hours make one release lookup; a check
   24 hours later or with `--refresh` makes another; a failed lookup prints one
   line ending `the session continues.`, exits `0` and is not retried for 24
   hours.
3. Success Metric: on an adopted temporary repository the Baseline line reads
   `Baseline: current with the installed Roundfix.`; on a repository whose
   managed guide is stale it reads `Baseline: plan ready with the installed
   Roundfix: <f> file change(s), <h> history unit(s), <k> skill change(s).`
   with the preview's counts; without a Setup Manifest it names the preview's
   state and message; outside a Git repository it reads `Baseline: not checked
   (outside a Git repository).`
4. Success Metric: `update.mode: always` in Project Config makes every command
   that loads configuration refuse with `update.mode must be notify or auto`;
   an absent key reads `notify`; Project Config `auto` over User Config `notify`
   reads `auto`. At `9a285d77` the key is refused as an unknown field.
5. Success Metric: the Roundfix skill's `modes` reference holds an `update`
   section of five numbered steps whose apply step confirms the previewed
   `planDigest`, asks before a plan with history units or Retired Skills or
   over uncommitted changes, and never commits; the skill's Modes paragraph and
   the glossary name `update` beside the three existing modes.
6. Success Metric: a managed refresh of an adopted repository writes the new
   clause into its agent instructions guide, a second refresh reports it
   current, and this repository's own guide carries it.
7. Success Metric: the command guide, the configuration guide, both copies of
   the `setup` and `modes` references, `CONTEXT.md` (**Update Check** added,
   **Skill Mode** changed) and ADR-0371 state the rule, and the documentation,
   skill and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The maintainer's request quoted in the adopted reference, written before
  this Spec existed: "Em 06/11 havia 30 dessas notas paradas (Argus 11,
  Quaestor 8, Fiscus 6, Vortex 5), misturadas aos pedidos reais entre
  projetos." and "O conteúdo \"o que muda neste repositório\" continua, só muda
  o canal."
- The Claude Code hooks reference (https://code.claude.com/docs/en/hooks, read
  through Exa on 2026-11-06): "Claude Code adds a SessionStart hook's
  plain-text stdout to Claude's context", and "SessionStart runs on every
  session, so keep these hooks fast". This is why the check prints plain text,
  bounds its lookup and caches it.
- The Codex hooks documentation (https://developers.openai.com/codex/hooks,
  read through Exa on 2026-11-06): "Project-local hooks load only when the
  project `.codex/` layer is trusted" and "new or changed hooks are marked for
  review and skipped until trusted". This is why the Baseline ships an
  instruction and not a hook.
- The update-notifier README (https://github.com/sindresorhus/update-notifier,
  read through Exa on 2026-11-06): `updateCheckInterval` "Default: `1000 * 60
  * 60 * 24` *(1 day)*" and the check "prevents any impact on your package
  startup performance". This is the interval the cache uses.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, a session in an
  adopting repository on the operator's machine shows the Update Check before
  its first task, and the Roundfix skill's `update` mode upgrades the binary
  and applies that repository's Baseline plan.
  reason: it needs the published release, the GitHub API, the real Roundfix
  Home and another repository, which no Verification or QA row may reach.
  satisfied-by: the operator's session evidence in an adopting repository after
  the release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-06 against these sources.

- Characterization at `9a285d77`:
  - `roundfix update check` prints `roundfix: unknown command "update"` and
    exits `2`; a `.roundfixrc.yml` with `update:` is refused as `field update
    not found in type config.configOverlay`;
  - `roundfix upgrade --check` reads `releases/latest` through `gh api`; with a
    temporary `HOME` and no `gh` login it fails with `upgrade failed: resolve
    latest release: ...` and exit `1`, so a lookup without credentials is a
    realistic failure;
  - the operational commands' Version Freshness line caches `releases/latest`
    for 24 hours in Roundfix Home and prints only `behind latest`;
  - the published release body of `v0.136.1` is the `CHANGELOG.md` section
    of that version, `### ` title and bullets (`gh release view v0.136.1`);
  - `roundfix baseline update --repo . --format json` on this repository
    reports `current` in 0.4 s, 0.08 s with `--no-history`, and writes nothing
    without `--yes` or `--confirm-plan`;
  - the Baseline manages Markdown guides and root blocks only
    (`docs/agents/setup-context.json` `managedArtifacts`), and no module
    renders a harness hook.
- A prototype of task_04 on a scratch copy at `9a285d77`: the clause in
  `rule.core.cross-project-continuity`, the module record command and `make
  baseline-digests` rewrote the derived files the Task declares; the managed
  refresh added the clause to `docs/agents/agent-instructions.md` and a second
  refresh reported no change; `./internal/baseline` passed after the clause was
  added to the force characterization (`_techspec.md` → Declared breaks).
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for release notes
  to adopters and for SessionStart hooks (no result above the score); the
  index led to `wiki/concepts/claude-code-hooks-e-automacao.md`, whose rule
  "prosa sugere, mecanismo garante" challenged an instruction-only design (the
  ADR records why a hook was still rejected as the Baseline's mechanism), and
  `wiki/concepts/cadencia-de-release-gatilhos-vs-lotes.md`, whose table notes
  that each release costs every adopter a note, an upgrade and a Baseline
  update.
- Exa: the Claude Code hooks reference and guide, the Codex hooks
  documentation and openai/codex issue 35306 (a project SessionStart hook
  silently skipped until trusted), the update-notifier README, and two CLIs
  that check once a day and never fail on the lookup (bitrise-cli
  `cmd/update.go`, uploads `update-check.ts`).

## Glossary

- adds: **Update Check**
- changes: **Skill Mode**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 4 brief for Spec 0361 of 2026-11-06, the maintainer's answers below and the adopted request.

The maintainer's answers (AskUserQuestion, 2026-11-06 ~12:00 BRT), verbatim:

- Question "O Jarvis repassou seu pedido de 06/11: parar de mandar as notas de
  release às caixas (havia 30 paradas) e trocar por uma verificação de
  atualização na Baseline, com /roundfix --update e escolha por repositório
  entre automático e manual. Como seguir?" → "Spec agora, padrão avisar
  (Recommended)" — option text: "Paro as notas já. A Baseline ganha uma
  checagem no início da sessão: compara a versão do Roundfix e da Baseline com
  a última publicada e mostra o que muda. /roundfix --update aplica. Uma chave
  por repositório liga o modo automático; o padrão é só avisar. Quando sair,
  cancelo as notas pendentes nas caixas."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando parar. Release assim que cada Spec ficar
  pronta; resposta ao OnionCry, ao sql-library e ao Jarvis."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

Taken from those answers, the request and the characterization, with no
further interview:

- What changes comes from the published release notes, as the `### ` titles
  of each release between the installed and the latest version, and from this
  repository's own `roundfix baseline update` preview; nothing is composed by
  hand or by a model.
- The key is `update.mode` in Project Config, over User Config, default
  `notify`: it changes no guide text, so it is not a Setup Manifest decision.
- The Baseline ships an instruction, not a harness hook; the command guide
  shows a hook a repository may add.
- The network is read at most once a day per machine, through one `gh api`
  release listing cached in Roundfix Home; a failure is cached too.
- The update mode applies the Baseline plan it showed by its Plan Digest and
  asks before history units, Retired Skill deletions or uncommitted changes;
  `auto` at session start applies no Baseline plan over uncommitted changes.
- ADR-0371 records the decision; no accepted ADR is amended. Four
  implementation Tasks in series and a QA gate.

## Open Questions

None.
