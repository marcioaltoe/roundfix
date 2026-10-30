---
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: active
created: 2026-10-30
surfaces: [backend, cli, docs]
---

# Runs load the dispatched skills, and the Baseline follows the curated catalog

A Run Agent today never loads the skills the Baseline's skill dispatch makes
mandatory for its Task, because nothing in its Task prompt points at the
dispatch guide and a runtime keeps only skill descriptions in context until a
skill is invoked; over 60 days, 108 Claude sessions of an adopter's frontend
Tasks loaded no frontend skill. The Baseline also still pins a skills catalog
from before the maintainer's curation of 2026-10-30, requires seven skills
the catalog removed, dispatches person-started and polish skills to every
matching Run, and ships the person-started `implement-spec` skill nobody
should start any more. After this Spec the Task prompt, the QA contract and
`implement-task` tell the Agent to load only the skills whose triggers match
its Task and to name them in its Result; the Setup Snapshots follow catalog
commit `86df860` with the seven skills and `implement-spec` retired; the
dispatch asks for person-started and polish skills only on request and splits
product UI from public pages; and `implement-spec` leaves the binary, its
two unique rules moving into the Roundfix skill.

## Sources

Two Backlog Entries, written and adopted at birth by this Spec's planning
change from ai-memory messages of `default/marcioaltoe-skills` (see
`references/_index.md`):

1. [`references/2026-10-30-the-baseline-follows-the-curated-skills-catalog.md`](references/2026-10-30-the-baseline-follows-the-curated-skills-catalog.md),
   message `01a18d69-3e58-7a33-9917-70040f0b6eb5` (2026-10-30 12:24 BRT):
   the catalog removed seven skills and completed `ui-ux-pro-max`.
2. [`references/2026-10-30-the-completed-ui-ux-pro-max-updates-automatically.md`](references/2026-10-30-the-completed-ui-ux-pro-max-updates-automatically.md),
   message `01a18da4-f737-74d0-80db-5970ece8fb53` (2026-10-30 13:29 BRT),
   which corrects the first: at `86df860` the skill updates automatically
   through local patches.

The skills study of 2026-10-26 and the 60-day usage table of 2026-10-30 are
operator measurements on one machine; their figures are cited here, the files
are not copied. No open Backlog Entry or unresolved Finding shares this
context: at the authoring base `docs/backlog/` and `docs/findings/` hold
none.

## Prerequisites

None in `requires`. This Spec is authored against `35de7f74` (v0.126.0) and
delivers second in cycle AS, after Spec 0332 and before Specs 0334 and 0335,
in a serial queue. Spec 0335 also raises the Roundfix owned skill, and Spec
0334 may edit the authoring skills this Spec touches (`write-tasks`); each
regenerates the owned-skill version record, the module version record, the
derived Baseline files and `docs/references/behavior-surfaces.json` with the
sanctioned commands on its own delivery base, never with a value measured at
authoring time.

task_01 reads the local checkout `~/dev/skills`, which must hold commit
`86df8608726979064a8fe46003dfaf403ffe0087` (it is upstream's `main` there at
authoring time); task_01 stops before any write when it does not.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; the new activation bundle and trigger follow the
  catalog's dotted naming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added, and no Task or check reaches the network. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0346, "A Run Agent loads the
  skills the dispatch selects for its Task", which this Spec records and
  task_03 implements ("The Daemon does not check the list"); ADR-0284,
  "Carry-forward binds a Task to the planning text it references", whose
  input set gains the dispatch guide, and ADR-0170, "A Task already completed
  on its target is nothing to carry", which cites it and stays unchanged; ADR-0246, "A retired skill leaves the
  Baseline and the adopter removes its copy", which tasks 01 and 04 follow
  ("Roundfix never deletes an installed skill tree"); and ADR-0206, "The
  Baseline takes upstream setup names and drops a removed skill", which the
  asset sync of task_01 follows for the seven catalog removals. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AS answers of 2026-10-30 and the Governed Paths this Spec declares,
  quoted in the authorization record; and the maintainer's express
  grant of 2026-10-30 for the `OWNED_SKILLS` line of the `Makefile` only,
  which task_04 changes to retire `implement-spec`, quoted in the record's
  Makefile grant section. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog/_authorization.md`;
  bounded files:
  `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/context7-cli/references/docs.md`,
  `.agents/skills/context7-cli/SKILL.md`,
  `.agents/skills/grilling/SKILL.md`,
  `.agents/skills/implement-spec/SKILL.md`,
  `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec-delivery.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/tui-design/references/evidence.md`,
  `.agents/skills/tui-design/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `Makefile`,
  `docs/agents/setup-context.json`,
  `docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/frontend.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `internal/baseline/assets/skill-activations.json`,
  `skills/archive-spec/SKILL.md`,
  `skills/implement-spec/SKILL.md`,
  `skills/implement-task/SKILL.md`,
  `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- Every Task prompt and QA prompt of a repository with a skill dispatch guide
  names the guide and tells the Agent to load only the skills that match its
  work, and every implementation Agent names what it loaded in its Result.
- The Setup Snapshots follow catalog commit `86df860`; `design-patterns`,
  `diagnosing-bugs`, `fixing-metadata`, `game-changing-features`,
  `refactoring-analysis`, `tdd`, `to-prompt` and `implement-spec` are Retired
  Skills that no module, trigger, bundle or snapshot names.
- The dispatch starts `architectural-analysis`, `grill-with-docs`,
  `humanizer`, `writing-clearly-and-concisely` and `web-design-guidelines`
  only when the person asks, keeps `wcag-audit-patterns` for a formal audit,
  and gives product UI and public pages separate bundles.
- The binary ships twelve owned skills; the Roundfix skill holds the two rules
  only `implement-spec` stated.

## User Stories

1. As an adopter maintainer, I want a frontend Task's Agent to load the
   frontend skills the dispatch requires, so that the bundle the Baseline
   makes mandatory changes the work.
2. As an adopter maintainer, I want an Agent to load only the skills that
   match its Task, so that the context is spent on the work and not on the
   catalog.
3. As the maintainer, I want the Baseline to follow the curated catalog and to
   tell adopters which installed copies to delete, so that no adopter keeps
   using a skill I removed.
4. As an operator of an unattended Run, I want no person-started skill to be
   dispatched to a Run, so that a Run never stalls asking for a person.
5. As the maintainer, I want the queue rules of `implement-spec` kept in the
   Roundfix skill, so that retiring the skill loses nothing.

## Core Features

1. The load rule. The Spec Context Bundle lists the skill dispatch guide when
   the repository has one. The Task prompt tells the Agent, before editing, to
   read it, select only the triggers and activation bundles that match its
   Task Type, declared paths and requested work, load each selected skill, not
   load a skill only a person can start, and name the result in its Result as
   the **Loaded Skill List**. The QA contract tells the QA Agent to load the
   skills of the QA activation and of the triggers its declared surfaces
   match. `implement-task` and `qa-gate` state the same rules. The dispatch
   guide becomes a carry-forward input. Nothing checks the list (ADR-0346).
2. The catalog. The Setup Snapshots follow catalog commit `86df860`, with the
   completed `ui-ux-pro-max`; the seven skills the catalog removed become
   Retired Skills and leave every module, trigger and bundle; the TypeScript
   snapshot's activation bundles equal the catalog's; this repository
   restores the three upstream skills the new commit changed.
3. The dispatch review. The five person-started or polish skills are
   dispatched only when the person asks; `wcag-audit-patterns` only for a
   formal audit; the product UI bundle holds `interaction-design`,
   `interface-design` and `fixing-accessibility`, and a new public pages
   bundle holds `frontend-design` and `fixing-accessibility`.
4. The retirement of `implement-spec`. It becomes a Retired Skill, leaves the
   binary, its canonical copy, its mirror, the owned-version record and the
   `context-workflow` module; `archive-spec` and `write-tasks` stop naming it;
   the Roundfix skill's delivery reference states its two unique rules.
5. Documentation. The Baseline guide lists the new Retired Skills and says
   that `ui-ux-pro-max` searches need Python 3; the skills, reconcile and
   context-driven guides follow; the glossary adds **Loaded Skill List** and
   revises **Spec Context Bundle**.

## User Experience

A frontend Task of a TypeScript adopter starts. Its prompt's Spec Context
Bundle lists `docs/agents/skill-dispatch.md` under Instructions. The Agent
reads it, matches the product UI and React triggers, loads
`interaction-design`, `interface-design`, `fixing-accessibility`, `react`,
`react-best-practices` and `react-composition-patterns`, edits, and ends its
Result with `Skills loaded: interaction-design, interface-design,
fixing-accessibility, react, react-best-practices,
react-composition-patterns`.

A TypeScript adopter runs `roundfix baseline update` after the release.
Besides the refreshed guides, the preview lists `Skills retired: 3` with, for
example, `- retired tdd: no longer required by the Baseline; delete
.agents/skills/tdd and its skills-lock.json entry`; the state and exit code
are unchanged.

## Non-Goals / Out of Scope

- Deleting an adopter's installed copy of a Retired Skill: Roundfix only
  reports it (ADR-0246).
- Removing `ui-ux-pro-max`, `zustand` or any other frontend skill: the
  maintainer kept them.
- Removing `testing-boss`: the maintainer kept it.
- A Daemon check, event or gate on the Loaded Skill List.
- A Python readiness check for `ui-ux-pro-max`.
- Editing the upstream catalog: removing `implement-spec` from its presets is
  a follow-up message to `default/marcioaltoe-skills`, which the operator
  sends after the release with the text under Open Questions.
- `implement-task`, which stays: the Daemon names it in every Task's Spec
  Context Bundle.
- The synthetic `tdd` literal in a catalog test fixture, which names no
  Baseline entry.

## Success Metrics

1. Success Metric: `TestNoCatalogEntryNamesARetiredSkill`,
   `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles` and
   `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree` pass, every
   upstream Setup Snapshot names `86df860`, and this repository's Baseline
   preview reports `Skills differing: 0`.
2. Success Metric: `TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks`
   and `TestTheFrontendBundlesSeparateProductUIFromPublicPages` pass, and a
   second managed refresh of this repository reports `File changes: 0`.
3. Success Metric: `TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills`,
   `TestBuildQAPromptTellsTheGateToLoadItsDispatchedSkills`,
   `TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent` and the
   `skill dispatch guide` case of
   `TestCarryForwardRefusesATaskWhoseBindingTextChanged` pass.
4. Success Metric: `TestImplementSpecIsNotARoundfixOwnedSkill` passes and the
   binary lists twelve owned skills.
5. Success Metric: `make verify` and `make verify-docs` pass on the delivered
   tree.
6. Success Metric: after the release, the operator's rerun of the frontend
   skill measurement over Runs started after it shows frontend Task sessions
   loading frontend bundle skills, against 0 of 108 sessions before.
7. Success Metric: after the release, a TypeScript adopter's
   `roundfix baseline update` lists under `Skills retired` each of the eight
   Retired Skills it still holds.

## Acceptance evidence

The outside-evidence row rests on published documentation this Spec did not
write, read on 2026-10-30.

- Claude Code documentation, "Extend Claude with skills"
  (https://code.claude.com/docs/en/skills): "In a regular session, skill
  descriptions are loaded into context so Claude knows what's available, but
  full skill content only loads when invoked." A skill the dispatch makes
  mandatory therefore reaches a Run Agent only when something tells the Agent
  to invoke it.
- Anthropic, "Agent Skills" overview
  (https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview):
  "Progressive disclosure ensures only relevant content occupies the context
  window at any given time." Loading only the matching skills is the rule the
  runtime itself assumes.

## Unreachable Acceptance

- criterion: Success Metric 6, after the release, frontend Task sessions load frontend bundle skills
  reason: it needs Runs of an adopter started after the release and the session transcripts of the operator's machine, which no Verification or QA row may read
  satisfied-by: the operator reruns `measure-frontend-skill-loads.py` over the Runs started after the release and reports the counts to `default/marcioaltoe-roundfix` through ai-memory
- criterion: Success Metric 7, after the release, a TypeScript adopter's baseline update lists its Retired Skills
  reason: it needs another repository and a released Roundfix, which no Verification or QA row may reach
  satisfied-by: the operator's first managed refresh of a TypeScript adopter after the release, reported to `default/marcioaltoe-roundfix` through ai-memory

## Validation

The decisions were checked on 2026-10-30 against these sources.

- Characterization at `35de7f74`: `assembleTaskContextBundle` lists the PRD,
  TechSpec, Task Graph, `AGENTS.md`, the `implement-task` skill and the
  Task's Context paths, never the dispatch guide; the Task prompt's execution
  invariants never mention a skill; `CarryForwardInputs` lists the same
  standard paths plus `CONTEXT.md`; the Retired Skills are `council` and
  `the-fool`; `baseline update` lists a retired copy and never deletes it;
  `implement-spec` is owned, person-started and dispatched by
  `context-workflow`; the `Makefile`'s `OWNED_SKILLS` line, which
  `make verify` reads through `skills-sync-check`, names it.
- A prototype in a scratch worktree of `35de7f74` applied all four parts:
  the asset sync at `86df860` wrote the five snapshots and its `--check`
  exited `0`; it kept each snapshot's recorded activation bundles, so the
  TypeScript snapshot still named `tdd` and `diagnosing-bugs` until its
  bundles were set from the catalog file before the sync; the record command
  cannot drop a retired owned skill, so its entry is deleted with `jq`, as
  Spec 0241 did; this repository's preview then reported three upstream
  skills differing (`context7-cli`, `grilling`, `tui-design`), restored with
  the sanctioned restore. The parity fixtures that list an adopter's
  installed `implement-spec` copy describe a captured tree and did not
  change. With the declared test changes `go test ./internal/baseline
  ./skills` passed.
- Skills study of 2026-10-26 (operator, one machine): `tdd` asks the user to
  confirm seams, impossible in a Run; `diagnosing-bugs` contradicts
  `systematic-debugging` in the same bundle; `architectural-analysis`,
  `refactoring-analysis` and `grill-with-docs` are person-started yet
  dispatched; `web-design-guidelines` fetches its rules from a mutable URL on
  every run. Usage table of 2026-10-30, 60 days: `tdd` loaded 180 times,
  `diagnosing-bugs` 181, `implement-spec` 41, `interface-design` 3,
  `frontend-design` 4.
- Secondbrain, read only: `wiki/index.md` and
  `wiki/concepts/arquitetura-de-instrucoes-e-progressive-disclosure.md`,
  which places skills as the layer "invocados pela tarefa" below the project
  instructions; `qmd query "skills dispatch agentes carregam skills
  obrigatórias roundfix" --all --files --min-score 0.3` returned that page and
  Roundfix's own mirror. No Secondbrain file was written.
- Exa: the two pages quoted under Acceptance evidence. Challenging: a
  description listing can drop descriptions when many skills are installed,
  so the Agent must match triggers from the dispatch guide, not from the
  runtime's listing.

## Glossary

- adds: **Loaded Skill List**
- changes: **Spec Context Bundle**

## Decisions

- 2026-10-26: skills study (AskUserQuestion, recorded in the operator's
  memory) — dispatch Spec "com tudo": take `diagnosing-bugs` and `tdd` out of
  the bundles, take the dispatch off `architectural-analysis` and
  `refactoring-analysis` (person-started; they stall Runs), `writing-clearly`
  and `humanizer` only on request, split the frontend-ui-quality bundle
  (`fixing-accessibility` stays in it; `wcag-audit-patterns` only for a
  formal audit; `web-design-guidelines` only for an explicit review),
  `interface-design` for product UI × `frontend-design` for marketing pages.
- 2026-10-30: testing skills — "Vamos manter testing-boss, vou solicitar
  somente que remova tdd por enquanto."
- 2026-10-30: scope — "como os projetos só vão atualizar a baseline
  posteriormente ao release, quero já remover tudo que é irrelevante da
  curadoria e ajustar a baseline para remover implement-spec. Assim a spec do
  dispatch vira uma revisão do processo e setup de skills."
- 2026-10-30: removals — "Remover: tdd, diagnosing-bugs, to-prompt,
  refactoring-analysis, design-patterns, game-changing-features,
  fixing-metadata"
- 2026-10-30: frontend skills — "Não remova ainda ui-ux-pro-max e zustand,
  vamos manter as skills relacionadas ao frontend e melhorar a
  implement-tasks para que o uso das skills sejam mais eficientes."
- 2026-10-30: `implement-task` stays — it is the Run agent's manual; the
  Daemon names it in every Task's Context Bundle.
- 2026-10-30: plan (AskUserQuestion) — "Aprovar como está (Recommended)" —
  "Quatro Specs na ordem validada pelo Jev; a 0332 entra primeiro para que a
  fila resolva sozinha conflitos de registros gerados durante a 0333."
- 2026-10-30: authorization (AskUserQuestion) — "Mesmas regras do AR
  (Recommended)" — "Push de branches, PR dos planos, entrega e merge pela
  Delivery Queue, eu escolho o ponto de release e aviso os adotantes por
  ai-memory."
- 2026-10-30: release (AskUserQuestion) — "Uma no fim do ciclo
  (Recommended)" — "v0.127.0 depois das quatro Specs; os adotantes recebem a
  revisão de skills e as correções juntas."

Autonomous route: no person present; decisions taken from the operator's cycle AS brief for Spec 0333 of 2026-10-30 and the characterization above.

Taken from those sources, with no further interview:

- QA Tasks follow the load rule through the QA contract and `qa-gate`, but
  write no list: the QA Report's matrix stays the gate's evidence. See
  ADR-0346.
- The Daemon does not check the Loaded Skill List. See ADR-0346.
- The dispatch guide joins the carry-forward inputs, as `AGENTS.md` and
  `implement-task` are. A Task carried across this Spec's merge sees its
  `implement-task` input moved and runs again.
- The seven skills the catalog removed become Retired Skills, not only
  removed names, so a TypeScript adopter learns which copies to delete.
- `refactoring-analysis` is retired, which settles the 2026-10-26 decision
  to stop dispatching it.
- `wcag-audit-patterns` and `web-design-guidelines` keep their own triggers,
  reworded, outside every bundle.
- The rename of `implement-spec` mentions in `archive-spec` and `write-tasks`
  points at the Roundfix Spec Run and the Daemon, the actors that replaced it.

## Open Questions

- The follow-up message to `default/marcioaltoe-skills` after the release,
  which the operator sends: "Roundfix 0.127.0 retirou a skill implement-spec
  (Retired Skill; o binário não a distribui mais). Os presets go, rust,
  python e typescript do catálogo ainda a listam em setups/*.txt; o asset
  sync do Roundfix a descarta. Pedido: remover implement-spec dos presets e
  do registry na próxima curadoria."
