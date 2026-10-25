---
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: active
created: 2026-10-24
surfaces: [backend, cli, docs]
---

# The Baseline sends handoffs and messages through ai-memory

The Baseline today sends every cross-project observation to the Secondbrain
inbox, a folder of files any session may commit into, and dispatches an
upstream `handoff` skill that upstream has now replaced. After this Spec the
Secondbrain stays the read-only business knowledge base; session handoffs and
cross-project messages go through ai-memory when the session exposes its MCP
tools, and the guidance says what is missing when it does not. Every upstream
Setup Snapshot follows the upstream commit that carries the six `ai-memory-*`
skills, the core module dispatches `ai-memory-handoff` and
`ai-memory-messaging` in place of `handoff`, each retired inbox clause is
replaced by one ai-memory clause with the same force, `roundfix follow-up`
stops reading the Secondbrain inbox, and a User Config that still names its
two keys keeps loading with one warning per key.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`](references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md),
   `type: refactor`, from the maintainer's request relayed through ai-memory
   (message `01a16e17`, 2026-10-24 10:26 BRT) and his words of the same
   morning, "Lembrando que solicitei em algum momento para matar o inbox de
   secondbrain e usar somente as mensagens do ai-memory".

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. This Spec is authored against `b5e45c63` (v0.118.0 plus
the plans of Specs 0315 and 0316) and delivers after Spec 0316 and before
Specs 0307, 0309 and 0310, by the maintainer's answer quoted in
`_authorization.md`. Specs 0312 to 0316, delivered before it, also edit
`internal/baseline/assets/modules/core.json`, the module version record, the
derived Baseline files, the owned-skill version record and
`docs/references/behavior-surfaces.json`; each Task here regenerates those
files with the sanctioned commands on its own delivery base and never copies
a value measured at authoring time. Spec 0307's follow-up mode and Spec
0309's Python import were re-amended in the same planning change to build on
this Spec; neither is required by this Spec's Verification.

The import reads the local checkout `~/dev/skills`, which must hold commit
`65d5f365cace27c4b468db1723213723b10c870b` (it is upstream's `origin/main`
there at authoring time); task_01 stops before any write when it does not.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; new clause, rule, trigger and coverage identifiers follow
  the catalog's dotted naming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; Roundfix's code never calls ai-memory, and the
  follow-up command keeps reading messages only from a file the caller hands
  in. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0330, "The Secondbrain is
  read-only and ai-memory carries handoffs and messages", which this Spec
  records ("ai-memory is optional and detected at run time") and every Task
  follows; and ADR-0095, "The secondbrain inbox is the single capture door",
  which ADR-0330 supersedes and task_04 retires by deletion. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP batch 6 decisions of 2026-10-24 and the Governed Paths this Spec
  declares. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory/_authorization.md`;
  bounded files:
  `.agents/skills/ai-memory-durable-pages/LICENSE`,
  `.agents/skills/ai-memory-durable-pages/SKILL.md`,
  `.agents/skills/ai-memory-handoff/LICENSE`,
  `.agents/skills/ai-memory-handoff/SKILL.md`,
  `.agents/skills/ai-memory-learning-maintenance/LICENSE`,
  `.agents/skills/ai-memory-learning-maintenance/SKILL.md`,
  `.agents/skills/ai-memory-messaging/LICENSE`,
  `.agents/skills/ai-memory-messaging/SKILL.md`,
  `.agents/skills/ai-memory-retrieval/LICENSE`,
  `.agents/skills/ai-memory-retrieval/SKILL.md`,
  `.agents/skills/ai-memory-routing-install/LICENSE`,
  `.agents/skills/ai-memory-routing-install/SKILL.md`,
  `.agents/skills/context7-cli/SKILL.md`,
  `.agents/skills/context7-cli/references/docs.md`,
  `.agents/skills/context7-cli/references/setup.md`,
  `.agents/skills/context7-cli/references/skills.md`,
  `.agents/skills/golang-cli/assets/examples/go.mod`,
  `.agents/skills/grilling/SKILL.md`,
  `.agents/skills/handoff/SKILL.md`,
  `.agents/skills/handoff/agents/openai.yaml`,
  `.agents/skills/knowledge-workspace/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/systematic-debugging/SKILL.md`,
  `.agents/skills/typesafe-ai/SKILL.md`,
  `.agents/skills/write-idea/SKILL.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `docs/agents/agent-instructions.md`,
  `docs/agents/docs-layout.md`,
  `docs/agents/secondbrain.md`,
  `docs/agents/setup-context.json`,
  `docs/agents/skill-dispatch.md`,
  `docs/agents/spec-routing.md`,
  `internal/baseline/assets/coverage.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `skills/roundfix/SKILL.md`,
  `skills/write-idea/SKILL.md`,
  `skills/write-prd/SKILL.md`,
  `skills/write-techspec/SKILL.md`.

## Goals

- Every upstream Setup Snapshot names upstream commit `65d5f36`, holds the
  six `ai-memory-*` skills, and holds no `handoff`, `brand-guidelines` or
  `gws-*` skill.
- The core module requires and dispatches `ai-memory-handoff` and
  `ai-memory-messaging`; no module requires or dispatches `handoff`.
- The rendered guides send cross-project observations as ai-memory messages,
  consume a message only after recording it, and fall back to the response
  when ai-memory is missing; no rendered guide carries the Inbox Entry
  contract or a Secondbrain write permission.
- An adopter whose Setup Manifest holds the retired inbox clauses refreshes
  with each of them recorded `replaced`, none unaccounted.
- `roundfix follow-up` reads no Secondbrain directory, and a config that still
  names `secondbrain.root` or `secondbrain.inbox_namespace` loads with one
  deprecation warning per key.
- This repository holds its changed upstream skills at the new commit, the
  owned authoring skills read ai-memory messages instead of the inbox, and
  the glossary and guides describe the new flow.

## User Stories

1. As the maintainer, I want one cross-project channel, so that an
   observation never lands in a folder no session reads.
2. As an Agent in an adopter repository, I want to be told when ai-memory is
   missing and to hand the text to the user, so that I never write into
   another project's files in its place.
3. As an adopter that refreshes its Baseline, I want every retired inbox
   clause accounted, so that the refresh never stops on an unaccounted
   clause.
4. As the maintainer, whose User Config names `secondbrain.root`, I want
   every `roundfix` command to keep working after the upgrade, so that I can
   remove the key at my own pace.

## Core Features

1. The upstream setups. The asset sync moves `go`, `rust`, `typescript` and
   the composed `go-cli-typescript-bun` to `65d5f36`; the core module
   requires and dispatches the two ai-memory skills in place of `handoff`;
   the parity row, the derived files and this repository's Setup Manifest
   follow their sanctioned commands; this repository restores its changed
   upstream skills and the six ai-memory skills, reconciles the lock entry of
   `handoff` and deletes the retained `handoff` tree; `skills/recommended.txt`
   follows the lock.
2. The guidance. A new core rule, `rule.core.cross-project-continuity`, holds
   four clauses that replace the eight Secondbrain inbox and capture clauses
   with the same force; a new Secondbrain rule carries the research digest as
   a message to the Secondbrain's owner; the Secondbrain write prohibition
   covers `inbox/`; three context-workflow clauses and one spec-workflow
   clause stop naming the Inbox Entry.
3. The follow-up command. The `secondbrain-inbox` source is removed; the
   project name in the report comes from the `origin` remote or the checkout
   directory; `secondbrain.root` and `secondbrain.inbox_namespace` become
   deprecated keys ignored with a warning; the command guide, the
   configuration guide and the Roundfix skill reference follow.
4. Documentation. The owned skills `write-idea`, `write-prd` and
   `write-techspec` read this project's ai-memory messages and send research
   digests as messages; the glossary retires the Inbox Entry and adds
   **Cross-project Message**; ADR-0095 is retired by deletion; the history
   guide sends cross-project lessons as messages.

## User Experience

An Agent in an adopter repository that notices a defect another project owns
reads in `docs/agents/agent-instructions.md` that it must send one
self-contained message with `memory_message_send`. When its session exposes
no ai-memory tools, it tells the user "the ai-memory MCP server is missing"
and puts the message text in its answer. A session that starts in the owning
project reads its messages with `memory_message_list`, records each in a
Backlog Entry, a Finding or a Spec, and only then pops it.

The maintainer upgrades Roundfix and runs `roundfix follow-up`. Standard
error shows `config: secondbrain.root is deprecated and ignored; the
Secondbrain inbox source was retired, so remove it`; the report lists the
other sources and exits `0`.

## Non-Goals / Out of Scope

- The repository-local `docs/_inbox/` folder, its `repository-inbox` source
  and the `inbox` rows of `roundfix deliver plan`, which hold raw notes of a
  repository without a Secondbrain; they keep their behavior.
- `docs/handoffs/`, which the documentation layout keeps for session
  continuity a repository records in Git.
- Migrating pending Inbox Entries or deleting `inbox/**` in the Secondbrain:
  the operator migrates them and the Secondbrain's own project deletes them.
- Any change to the Secondbrain consultation clauses, the Hermes escalation,
  the secret-safety rule or the `secondbrain.enabled` decision.
- The Python preset, which Spec 0309 imports after this Spec.
- Editing upstream skills, which stay vendored third-party content.
- The judge paragraph of `docs/user-guide/commands/spec.md` and an example
  line of `docs/user-guide/commands/doctor.md`, which name the retired inbox
  and skill only as past facts or sample output.

## Success Metrics

1. Success Metric: `TestEveryUpstreamSetupFollowsTheAiMemoryRevision` and
   `TestTheCoreModuleDispatchesTheAiMemorySkills` pass, and the asset sync's
   `--check` exits `0` against a clone of upstream at `65d5f36`.
2. Success Metric: `TestEveryInboxClauseAnAdopterHoldsIsReplaced` passes:
   eight inbox clauses `replaced`, none unaccounted, the kept clauses
   `retained`.
3. Success Metric: `TestTheGuidesSendCrossProjectWorkThroughAiMemory` passes,
   and this repository's `docs/agents/` names neither `Inbox Entr`,
   `inbox/<destination>` nor `_triaged`.
4. Success Metric: `TestTheRetiredSecondbrainKeysLoadWithAWarning` passes,
   and `roundfix follow-up --format json` under a User Config naming both
   keys exits `0` with the two warnings and no `secondbrain-inbox` source.
5. Success Metric: `go test ./internal/baseline ./skills ./internal/config
   ./internal/cli` and the documentation contracts pass, and a second
   managed refresh of this repository reports `File changes: 0`.
6. Success Metric: after release, an adopter's next managed refresh records
   each retired inbox clause `replaced`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-24.

- The upstream commit `65d5f365cace27c4b468db1723213723b10c870b` of
  `marcioaltoe/skills` (pull request #131, "refactor: replace workspace and
  continuity skills"), read in a local clone: `setups/go.txt`,
  `setups/rust.txt` and `setups/typescript.txt` replace
  `skills/07-evidence-delivery/handoff` with the six `ai-memory-*` entries;
  `setups/python.txt` holds 52 entries with the same six and no `handoff`.
- The ai-memory documentation, `docs/ARCHITECTURE.md` of
  `akitaonrails/ai-memory` (https://github.com/akitaonrails/ai-memory/blob/main/docs/ARCHITECTURE.md),
  "MCP tool surface": `memory_message_send` "Send a directed cross-project
  message into another project's inbox … the recipient must already exist
  (fail-closed, never created)", and `memory_message_pop` "Claim ONE inbox
  message exactly once".

## Unreachable Acceptance

- criterion: Success Metric 6, after release, an adopter's next managed refresh records each retired inbox clause `replaced`
  reason: it needs another repository and a released Roundfix, which no Verification or QA row may reach
  satisfied-by: the operator's first managed refresh of an adopter after the release, reported to `default/marcioaltoe-roundfix` through ai-memory

## Validation

The decisions were checked on 2026-10-24 against these sources.

- Characterization at `b5e45c63`: the core module requires and dispatches
  `handoff` (`trigger.core.handoff`); the four upstream snapshots name
  `b3c45a45`; no snapshot holds `gws-*` or `brand-guidelines`;
  `docs/agents/secondbrain.md` renders eight inbox and capture clauses; the
  Source Baseline `baseline.standard-typescript-monorepo-0.0.1` records each
  of them; `classifySourceClauseTransition` refuses a refresh with an
  unaccounted clause unless exactly one current clause `replaces` it with the
  same force; `internal/config` decodes with `KnownFields(true)` and strips
  registered deprecated keys first.
- A prototype in a scratch worktree of `b5e45c63` ran Build Orders 1 and 2
  and the configuration half of Build Order 3. The asset sync at `65d5f36`
  first refused while the core module still required `handoff`
  (`catalog.profile.skill.outside-setup`), then wrote the four snapshots once
  the module swapped the skills; `--check` exited `0`. A clause added to a
  rule a profile requires needed a Source Baseline row
  (`catalog.sourceBaseline.required-clause.missing`), so the research digest
  clause lives in a new rule. Reformatting `core.json` broke
  `TestCatalogDiagnosticCharacterization`, so module edits keep the file's
  layout. With the declared test changes the whole suite passed
  (`go test -tags repocontract ./...`) and the documentation contracts
  passed. After the prototype, Spec 0309's import at the same commit wrote
  only `python.json` (52 skills) and left the other snapshots unchanged.
- Secondbrain, read only. Read `wiki/index.md`,
  `wiki/entities/ai-memory.md` (the 2026-10-19 decision to adopt ai-memory
  for handoffs "sem substituir o secondbrain nem o inbox", which the
  maintainer reversed on 2026-10-24) and
  `wiki/concepts/porta-de-captura-unica-do-fleet-inbox.md` (the capture door
  of ADR-0095 this Spec retires). Ran `qmd query "ai-memory messages replace
  secondbrain inbox" --all --files --min-score 0.3`. No Secondbrain file was
  written.
- Exa: the ai-memory architecture page quoted under Acceptance evidence and
  the `ai-memory-messaging` skill page
  (https://www.skills.sh/akitaonrails/ai-memory/ai-memory-messaging): a
  message is "a directed, claim-once inbox/queue" and "distinct from
  handoffs, which are same-project session continuity". Challenging: the
  recipient must already exist in ai-memory, so the fallback clause also
  covers a refused send.

## Glossary

- adds: **Cross-project Message**
- changes: **Triage**
- changes: **Follow-up Item**

## Skills

- unchanged: command history sanitize — the edit replaces one sentence of the history guide that sent cross-project lessons to the Secondbrain inbox; `roundfix history sanitize` and the archive reference that covers it do not change.
- unchanged: guide docs/user-guide/commands/history.md — the same one-sentence edit; the archive reference describes the history commands, which do not change.

## Decisions

- 2026-10-24 (about 11:15 BRT): "Spec no fim do AP (Recomendado)" — "Plano
  escrito agora em branch refactor/, entregue depois da 0310. Antes, combino
  com marcioaltoe-skills os nomes e as versões; o Secondbrain continua sendo
  a base de conhecimento, e o ai-memory vira o canal de handoff e de
  mensagens. Respondo ao Jarvis com o plano."
- 2026-10-24 (about 13:10 BRT): order "Antes da 0307 (Recomendado)" — "A 0317
  entra logo após a 0316. A 0307 é emendada para usar só o ai-memory (sem
  triagem de inbox), e a 0317 também tira a fonte secondbrain-inbox do
  comando follow-up. O setup Python vai uma versão para frente (v0.127)."
- 2026-10-24 (about 13:10 to 13:20 BRT): pending entries "Migrar as
  mensages ainda não triadas para o ai-memory e remover as demais e triadas
  e arquivadas."; release notes "Descartar todas"; other entries "Migrar as
  63, descartar tax-poc (Recomendado)"; removal "Eu migro, o secondbrain
  remove (Recomendado)"; adopter notes "Via ai-memory já (Recomendado)".

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 6 briefing for Spec 0317 of 2026-10-24 and the characterization above.

Taken from those sources, with no further interview:

- ai-memory is detected at run time, not a Baseline decision; see ADR-0330.
- The core module owns the ai-memory clauses, because handoffs and messages
  concern every adopter, with or without a Secondbrain.
- The two configuration keys become deprecated keys rather than unknown ones,
  so strict validation never refuses an existing User Config.
- This repository restores all six ai-memory skills, the family upstream
  ships together, and requires only the two the core module dispatches.
- ADR-0095 is superseded by ADR-0330 and retired by deletion, not edited.

## Open Questions

None.
