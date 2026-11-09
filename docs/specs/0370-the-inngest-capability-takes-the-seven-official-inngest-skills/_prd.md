---
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
status: active
created: 2026-11-09
surfaces: [backend, docs]
type: fix
---

# The Inngest capability takes the seven official Inngest skills

On 2026-11-09 the skills catalog `marcioaltoe/skills` removed the `inngest` skill and replaced it with seven official skills from `inngest/inngest-skills`, but the Baseline's `backend` module still requires and dispatches `inngest`, both TypeScript Setup Snapshots still pin it, and the Catalog Skill Bump stops at the mismatch. A TypeScript adopter therefore keeps a skill whose catalog path no longer exists and never receives its replacements. After this Spec the `backend` module requires and dispatches `inngest-setup`, `inngest-durable-functions`, `inngest-steps`, `inngest-events`, `inngest-flow-control`, `inngest-cli` and `inngest-middleware` by the catalog's suggested mapping, every upstream Setup Snapshot and this repository's vendored catalog skills follow catalog commit `e45466d2`, and `inngest` becomes a Retired Skill, so an approved `roundfix baseline update` installs the seven, deletes an adopter's `inngest` copy it proves unedited and keeps an edited one.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md`](references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md),
   ai-memory message 01a1c01f-b6ad, reported by marcioaltoe-skills on
   2026-11-09, `nature: measured`.

At the authoring base `docs/backlog/` holds one open entry, the disk-full
tests reading a host path, which Spec 0368 adopts and which shares no context
with this one, and `docs/findings/` holds no unresolved Finding about the
Baseline's skills.

## Prerequisites

None in `requires`. Authored against `cfc08683` (v0.137.2 plus the 0366 and
0367 plans) and rebased on `a60928fa`, which adds the 0368 and 0369 plans. The queue delivers 0366, 0367, 0368 and 0369 first. They may
change `CONTEXT.md` (other entries), `docs/references/behavior-surfaces.json`,
`skills/testdata/owned-skill-versions.json` and the version fields of
`.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`; this Spec
edits only the **Retired Skill** entry of `CONTEXT.md` and the `baseline`
reference of the Roundfix skill, and re-records the version record and the
Behavior Surface Record by their commands, so a rebase reruns the commands
instead of merging bytes. No other Spec of the cycle changes
`internal/baseline/assets/**`, the Setup Snapshots, `skills-lock.json` or the
vendored skills.

task_01 and the QA gate read the catalog from the operator's clone
`/Users/marcio/.roundfix-operator/scratch/spec-0370/skills`, which must hold
commit `e45466d2b9e15480eb286a4c19b513494a867bce` before the Run starts
(ADR-0360: "A Run without network takes the commit from a clone the operator
prepares and names in the Task"). task_01 stops and reports when it does not.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the seven skill names and their trigger ids follow the existing
  `trigger.<module>.<skill>` form. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is Baseline content,
  a Go set and tests that read local files; no Task, test or QA row reaches
  the network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0380 (new, in the tree with this
  plan) records the rule, ADR-0380: "A skill the catalog replaces is followed
  to its replacements, and the replaced skill becomes a Retired Skill." It
  amends ADR-0206, whose removed-skill rule this Spec narrows, ADR-0206: "It
  is never followed to a skill with a similar purpose under another name";
  task_02 writes the amendment notice. ADR-0350 deletes a Retired Skill's
  copy the lock proves unedited, ADR-0350: "It is `delete` when every copy
  it finds is provably the one the repository installed", which is how an
  adopter's `inngest` leaves. ADR-0360 makes the Catalog Skill Bump the only
  writer of the snapshots and the vendored skills, ADR-0360: "The Catalog
  Skill Bump is the only writer." ADR-0191 keeps the
  asset sync the snapshots' writer, ADR-0191: "A setup snapshot follows its
  upstream by name". ADR-0364 checks a required external skill only for
  presence, ADR-0364: "A required external skill is checked only for
  presence", so the seven appear as missing until installed. ADR-0204 builds
  the composed `go-cli-typescript-bun` snapshot from the upstream setups by
  name, ADR-0204: "A composed profile takes a setup composed from upstream
  setups by name", and the sync recomposes it. ADR-0322 imports the published
  Python preset, ADR-0322: "A Python profile imports the published upstream
  Python preset", which the bump moves to the same Catalog Revision. ADR-0345
  names the regeneration commands the suite guard accepts, ADR-0345: "A
  derived-record writer names its command to the suite guard". Kept and
  not changed, each because it cites a listed ADR: ADR-0219 binds a profile
  draft to the closest built-in profile; ADR-0221 and ADR-0241 hold a
  required skill and a repository profile to the embedded Setup Snapshots,
  which the bump moves; ADR-0224 holds this repository's vendored skills to
  the snapshot and the lock, which the bump keeps in step; and ADR-0264 never
  replaces a vendored skill it cannot prove older, which the bump's restore
  applies. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA batch 2 answers of 2026-11-09 quoted in the authorization record
  grant the Baseline source, its snapshots, fixtures and goldens, the
  vendored skill the bump moves and the owned skill reference this Spec
  edits; no Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0370-the-inngest-capability-takes-the-seven-official-inngest-skills/_authorization.md`;
  bounded files: `.agents/skills/context7-cli/references/setup.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `skills/roundfix/SKILL.md`.

## Goals

- A repository whose profile selects the `backend` module is required to
  hold the seven `inngest-*` skills, and its skill dispatch guide names each
  one with its own trigger.
- Every upstream Setup Snapshot and this repository's vendored catalog skills
  name catalog commit `e45466d2`, and the two TypeScript snapshots pin the
  seven skills and no `inngest`.
- An adopter with an unedited, lock-recorded `inngest` loses it on an
  approved `roundfix baseline update`, and one with an edited copy keeps it
  with the paths to delete by hand.
- ADR-0206, the glossary, the `baseline` command guide and the Roundfix
  skill's `baseline` reference say that `inngest` was replaced and retired.

## User Stories

1. As a TypeScript adopter that builds with Inngest, I want the Baseline to
   install and dispatch the official Inngest skills, so that my agent loads
   current guidance for functions, steps, events, flow control, setup, the
   CLI and middleware.
2. As an adopter that has the old `inngest` installed, I want the approved
   update to remove the copy it installed and keep one I edited, so that I
   neither keep a skill whose catalog path is gone nor lose my own changes.
3. As the maintainer, I want the Catalog Skill Bump to pass again at the
   catalog's current `main`, so that the next bump is the routine command
   ADR-0360 describes.

## Core Features

1. The `backend` module requires the seven `inngest-*` skills in place of
   `inngest` and dispatches each through one trigger,
   `trigger.backend.<skill>`, by the catalog's suggested mapping; its
   version rises to the next one by the record command.
2. `inngest` joins the Retired Skills, so no module, bundle or snapshot may
   name it and the approved update deletes an adopter's unedited copy.
3. The Catalog Skill Bump moves the five upstream Setup Snapshots, this
   repository's vendored catalog skills and their lock entries to catalog
   commit `e45466d2`; the asset-sync parity row, the derived Baseline files
   and the Setup Manifest follow through their commands.
4. Documentation: ADR-0380 (written by this plan) and its amendment notice
   in ADR-0206; the **Retired Skill** entry of `CONTEXT.md`; the `baseline`
   command guide and the Roundfix skill's `baseline` reference.

## User Experience

A TypeScript adopter's rendered `docs/agents/skill-dispatch.md` lists, in
place of the `inngest` entry:

```text
- `inngest-cli`:
  - `trigger.backend.inngest-cli`: Installing or running the Inngest CLI and its local Dev Server.
- `inngest-durable-functions`:
  - `trigger.backend.inngest-durable-functions`: Building or changing Inngest functions, their triggers, retries, and schedules.
- `inngest-events`:
  - `trigger.backend.inngest-events`: Designing Inngest events, fan-out, and idempotent event handling.
- `inngest-flow-control`:
  - `trigger.backend.inngest-flow-control`: Configuring Inngest concurrency, throttling, rate limits, debounce, or batching.
- `inngest-middleware`:
  - `trigger.backend.inngest-middleware`: Adding Inngest middleware for logging, tracing, error tracking, or dependency injection.
- `inngest-setup`:
  - `trigger.backend.inngest-setup`: Setting up Inngest in a project: the SDK, client configuration, environment variables, and the serve endpoint.
- `inngest-steps`:
  - `trigger.backend.inngest-steps`: Building or changing the steps of Inngest functions: step runs, sleeps, waits for events, and invocations.
```

Its `roundfix baseline update` preview lists the seven under
`Skills missing` and, for an unedited copy,
`- retired inngest: no longer required by the Baseline; the approved update deletes .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry`;
for an edited copy it prints
`- retired inngest: kept, its copy differs from its skills-lock.json computedHash; delete .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry by hand if the repository no longer uses it`.
Both lines are the existing Retired Skill messages of ADR-0350.

## Non-Goals / Out of Scope

- Requiring the Inngest skills only when a repository declares Inngest; the
  `backend` module required `inngest` unconditionally and keeps that shape
  (ADR-0380 records why).
- Installing the other eight skills `inngest/inngest-skills` publishes, which
  the catalog does not list.
- Following any other removed skill to a replacement, or changing the
  `capability.optional.inngest` probe.
- Repairing the Catalog Skill Bump's own order (it refreshes the Setup
  Manifest with a binary built before `make baseline-digests`); task_01
  refreshes it again with a fresh build, and the operator records the defect
  apart.
- Any change to the Makefile, the workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the frozen parity fixtures other than the
  sanctioned asset-sync row, or the `### QA settlement` section of any skill.

## Success Metrics

1. Success Metric: the embedded `backend` module requires and dispatches
   exactly the seven `inngest-*` skills with the triggers above and never
   `inngest`, and the restore profiles of `standard-typescript-monorepo` and
   `go-cli-typescript-monorepo` require and pin all seven at
   `skills/05-implementation-loop/<skill>`.
2. Success Metric: all five upstream Setup Snapshots and every catalog entry
   of `skills-lock.json` name `e45466d2`, none names `c8b2e85`, and a second
   Catalog Skill Bump at the same commit changes nothing.
3. Success Metric: on an adopted repository with an unedited, lock-recorded
   `inngest`, the preview plans its deletion and the approved apply removes
   `.agents/skills/inngest`, `.claude/skills/inngest` and the lock entry and
   nothing else; with an edited copy the update keeps it and changes no byte.
4. Success Metric: the new tests fail on the Delivery Base and pass on the
   delivered tree, and this Spec passes its own strict check.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The marcioaltoe-skills report quoted in the adopted reference, written by
  the catalog project before this Spec existed: "`inngest` → removed.
  Replacements, all at skills/05-implementation-loop/<slug>, update mode
  auto: `inngest-setup`, `inngest-durable-functions`, `inngest-steps`,
  `inngest-events`, `inngest-flow-control`, `inngest-cli`,
  `inngest-middleware`", with its suggested trigger mapping.
- The published `inngest/inngest-skills` README
  (https://github.com/inngest/inngest-skills, read through Exa on
  2026-11-09), whose "Available Skills" table lists `inngest-setup` ("Set up
  Inngest in a TypeScript project"), `inngest-events`,
  `inngest-durable-functions`, `inngest-steps`, `inngest-flow-control`
  ("Concurrency limits, throttling, rate limiting, debounce, priority,
  batching"), `inngest-middleware` ("Middleware lifecycle, dependency
  injection, built-in middleware") and `inngest-cli` ("Install and use the
  Inngest CLI and Dev Server").

## Unreachable Acceptance

- criterion: Adopter migration: after the patch release that carries this
  Spec, a TypeScript adopter that has `inngest` installed runs
  `roundfix baseline update --yes` and ends with the seven `inngest-*`
  skills installed from the catalog and `inngest` deleted, or kept and listed
  when its copy was edited.
  reason: installing the seven reaches GitHub and the adopter is another
  repository, which no Verification or QA row may use.
  satisfied-by: the operator runs the update in one TypeScript adopter after
  the patch release and records its `Skills installed` and `Skills deleted`
  lines through ai-memory, recorded session evidence.

## Validation

The decisions were checked on 2026-11-09 against these sources.

- Prototype on a scratch copy of `cfc08683`: the module edit, the Retired
  Skill entry, the module record command,
  `make skills-third-party-bump SKILLS_CATALOG=<local clone> SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce`
  (exit 0, "33 catalog skills and every upstream Setup Snapshot pinned"),
  the asset-sync parity transform, `make baseline-digests` and a managed
  refresh built after it. The bump changed 19 files besides the module and
  the Go set: the five snapshots, `skills-lock.json`, one vendored
  `context7-cli` file, the profile golden digest, the formatter golden, three
  catalog records, four plan goldens and the Setup Manifest; the parity
  transform and the second `make baseline-digests` changed only the
  `typescript` row of `asset-sync.json` and the corpus manifest.
- The same prototype passed `go test ./internal/baseline ./skills`,
  `-tags docscontract ./internal/docscontract`, `-run 'Baseline|Skill|Doctor'
  ./internal/cli` and `make repo-test`; dropping `inngest` from the Retired
  Skills, dropping one slug from the module and changing one trigger each
  failed a named new test.
- The bump refreshes the Setup Manifest with a binary built before
  `make baseline-digests` rewrites the profile golden digest, so the recorded
  catalog digest is stale until a managed refresh with a fresh build; this
  Spec's task_01 runs that refresh.
- Reading at `cfc08683`: `loadRestoreProfile` takes the external skills from
  the profile's snapshot and the required ones from its modules;
  `InstalledRetiredSkills` inspects every Retired Skill whatever the profile;
  `TestNoCatalogEntryNamesARetiredSkill` pins the Retired Skill list;
  `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` compares the
  snapshots with the sanctioned asset-sync parity row.
- Catalog compare `c8b2e85...e45466d2` (`gh api`, read only): eight commits;
  `setups/typescript.txt` is the only setup list that changed.
- Published sources (Exa): the `inngest/inngest-skills` README.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "inngest skill baseline backend module skills catalog" --all --files --min-score 0.3`
  and
  `qmd query "retired skill baseline update deletes unedited copy replacement" --all --files --min-score 0.3`;
  they returned this repository's mirrored ADRs, Spec records 0200, 0208
  and 0317 and the catalog's `backend-docs` skill map; no page decides
  otherwise.
- Pending ai-memory messages: none read by this author; the operator relayed
  the marcioaltoe-skills message.

## Glossary

- changes: **Retired Skill** — a skill the catalog replaces with skills it names becomes a Retired Skill, and the module that required it requires the replacements.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BA batch 2 brief for Spec 0370 of 2026-11-09 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-09 ~11:20 BRT), verbatim:

- Question "O catálogo marcioaltoe/skills trocou `inngest` por 7 skills
  oficiais `inngest-*`. A Baseline do Roundfix ainda instala e despacha
  `inngest` (módulo backend, setups typescript e go-cli-typescript-bun,
  fixtures). Mudança de Baseline exige `baseline update` nos adotantes. Como
  seguir?" → "Spec F agora (Recommended)" — option text: "Spec fix na
  Baseline: a capacidade opcional Inngest instala as 7 skills e despacha
  pelos gatilhos sugeridos (funções/steps, eventos, flow control, setup, CLI,
  middleware); fixtures atualizadas. Release logo após o merge; os adotantes
  pegam pelo Update Check."
- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto
  fica para amanhã."

Taken from those and the evidence, with no further interview:

- `inngest` becomes a Retired Skill rather than an upstream-removed name the
  adopter deletes by hand, because ADR-0350 already deletes a copy proven
  unedited and keeps an edited one; ADR-0380 records the rule and amends
  ADR-0206, and task_02 writes the notice.
- The snapshots move through the Catalog Skill Bump at `e45466d2`, the
  repository's only writer of them (ADR-0360), from the operator's local
  clone; the other skills the bump moves are routine catalog changes.
- The trigger texts are the catalog's suggested mapping, one trigger per
  skill.
- Two implementation Tasks and the QA gate, in series: the Baseline change
  with its tests, then the ADR notice, glossary and guides.
- The release is the operator's patch right after the merge; no Task
  releases.

## Open Questions

None.
