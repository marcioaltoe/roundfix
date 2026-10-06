---
spec: 0241-retire-the-fool-autoresearch-and-council
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# Retire the-fool, autoresearch and council

On 2026-10-06 the maintainer asked to retire the skills that no longer have a
use, or whose job the Jev judge took over. Asked per skill, the maintainer
answered: "remover: the-fool, autoresearch, council // manter: grilling,
grill-with-docs, write-idea, business-analyst e handoff". The adopted Backlog
Entry,
[references/2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md](references/2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md),
measured the use: across 220 Runs, `the-fool` and `autoresearch` were never
read, and `council` was read in 3 Runs.

Today the `context-workflow` module requires and dispatches `council` and
`the-fool`, so every adopter installs both and its `skill-dispatch.md` names
them. All four Setup Snapshots list them. `council` is also one of the 14
Roundfix-owned skills the binary ships, and `write-idea` runs it in its step
6. `autoresearch` is installed only in this repository (`skills-lock.json`,
`skills/recommended.txt`, `.agents/skills/autoresearch/`); no module names it.

The upstream skills catalog (`marcioaltoe/skills` at `b3c45a4`) still lists
`council` and `the-fool` in its `go`, `rust` and `typescript` setups. Spec
0208 dropped `review` and `triage` because upstream had removed them; that
rule (ADR-0206) does not cover a skill upstream still lists, and the asset
sync would put both names back into every snapshot. This chore therefore
introduces Retired Skills (ADR-0246): skills the Baseline stops requiring
while upstream may still list them. It is the third item of the maintainer's
order ("Tudo, na ordem sugerida") and is planned for v0.52.0.

## Measured readers

`git grep -n -I -E 'the-fool|autoresearch|council' -- . ':!docs/history'` on
`40a7893d` found these readers outside the three skill trees and the adopted
entry:

| Reader | Names | Disposition |
| --- | --- | --- |
| `internal/baseline/assets/modules/context-workflow.json` | `council`, `the-fool` (requiredSkills, two triggers) | removed (task_01) |
| `internal/baseline/assets/setups/{go,rust,typescript,go-cli-typescript-bun}.json` | `council` (repo-owned), `the-fool` | dropped by the asset sync (task_01) |
| `internal/baseline/assets_sync.go` | `council` in the repository-owned set | removed (task_01) |
| `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` | both, three setups | sanctioned fixture follows the snapshots (task_01) |
| formatter golden and `docs/agents/skill-dispatch.md` | both triggers | regenerated (task_01) |
| `skills/skills.go`, `Makefile` `OWNED_SKILLS`, `skills/testdata/owned-skill-versions.json` | `council` | removed (task_02) |
| `write-idea` (`SKILL.md`, `references/idea-template.md`), `write-prd` `SKILL.md`, both copies | `council` | reworded (task_02) |
| `docs/user-guide/commands/skills.md` | `council` | removed (task_02) |
| `skills-lock.json`, `skills/recommended.txt` | `autoresearch`, `the-fool` | removed (task_02) |
| `skills/repository_test.go` | `autoresearch` as a sample external skill | replaced by a kept skill (task_02) |
| `.agents/skills/typesafe-ai/SKILL.md` | an upstream cookbook URL containing `autoresearch` | unrelated; vendored, unchanged |
| `.agents/skills/council/references/archetypes.md` | `the-fool` | deleted with its tree |

No Source Baseline corpus, profile, activation bundle, `CONTEXT.md` entry or
clause names any of the three. The doctor line in `README.md`,
`docs/user-guide/usage.md` and `docs/user-guide/commands/doctor.md` counts
the owned and required skills, so it changes with them.

## Prerequisites

None that the queue must enforce. Specs 0239, 0240 and 0242 are authored in
parallel. 0239 edits authoring skills and `context-workflow` clauses, so it
may share `internal/baseline/assets/modules/context-workflow.json`,
`write-prd`, `CONTEXT.md` and `skills/testdata/owned-skill-versions.json`
with this Spec; this Spec changes only the lines that name the retired skills,
and a version raise is regenerated at merge (ADR-0233). The operator orders
the queue; this Spec's Verification does not depend on another Spec's
artifacts.

## Project Constraints

- Identifier strategy: applicable — no identifier scheme changes. Two
  dispatch triggers are removed (`trigger.context-workflow.council`,
  `trigger.context-workflow.the-fool`), the `context-workflow` module version
  rises by one, and `baseline update` gains the JSON field `skills.retired`.
  No project-owned Internal Identifier is generated. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded Baseline assets, local
  files and Git only. The asset sync reads a local clone of `~/dev/skills`;
  no Task or test clones from the network or reads a credential. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0246 (this Spec): "No module
  requires or dispatches a Retired Skill". It narrows ADR-0191, which
  otherwise holds: "A setup snapshot mirrors one upstream list, and the asset
  sync is its only writer", and ADR-0191: "Roundfix never deletes an
  installed skill tree", so adopters delete their own copies. ADR-0206 covers
  a skill upstream removed and is not changed: "A skill the catalog removes
  is dropped from every module, trigger and snapshot". Under ADR-0072 the
  frozen parity corpus changes only in the sanctioned fixture, ADR-0072:
  "New product behavior is recorded explicitly as a designed delta". Under
  ADR-0081 and ADR-0149 the derived pins follow the authorized edit,
  ADR-0081: "Derived pins regenerated by the sanctioned regeneration command
  are therefore deterministic fallout", and ADR-0149: "An authorization grant
  names only the regeneration command it sanctions". Under ADR-0058 and
  ADR-0060 every prior clause stays accounted for, ADR-0058: "must map each
  prior Normative Clause to a current rule"; no Normative Clause names a
  retired skill, so no Source Baseline retention disposition changes (the
  removed entries are a required-skill list and two dispatch triggers). This
  repository's managed refresh converges,
  ADR-0103: "a second run against an unchanged catalog reports the repository
  current". Retention accounting stays mechanical, ADR-0099: "it is a mechanical
  comparison that stays fail-closed without a model". The comparison keeps
  its meaning, ADR-0221: "Doctor now compares each required external skill
  that is installed and matches its lock". This repository stays held to its snapshot, ADR-0224: "every
  required external skill of this repository's Baseline Profile must match
  the `treeDigest` its Setup Snapshot pins". ADR-0204's composed
  snapshot is refreshed by the same sync; ADR-0219, which cites it, decides
  profile-draft binding and does not apply. ADR-0192 lets a delivery resolve
  a conflict confined to the derived paths by regeneration. ADR-0241's repository profiles read
  the same snapshots and are not changed. ADR-0189 and ADR-0233 bind the
  owned-skill version raises of `write-idea`, `write-prd` and `roundfix`.
  The QA gate follows ADR-0080: "QA verdicts distinguish
  environment-blocked rows", and ADR-0091: "required to be terminal and to
  depend on every leaf"; ADR-0096, ADR-0097, ADR-0104 and ADR-0167 bind its
  machine stage, row carry, outside evidence and pre-PR Pull Request row,
  and ADR-0117, ADR-0182, ADR-0194, ADR-0195 and ADR-0210 its stage
  ownership, settlement, row record, re-observation and evidence snapshot;
  ADR-0240 decides when its QA partial qualifies. ADR-0229 and ADR-0237 cite ADRs listed here
  but decide the Delivery Retry of an operator archive or an outside merge,
  which this Spec does not change, so neither applies. ADR-0173 cites
  ADR-0103 but decides the citations a History Relocation breaks; this Spec
  relocates no history, so neither it nor ADR-0177, reached only through
  it, applies. ADR-0093, ADR-0156, ADR-0168, ADR-0176 and ADR-0183
  check this Spec's consistency by citation and receipt, and ADR-0184 binds
  the TechSpec's Surface Transcript. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded
  in [_authorization.md](_authorization.md): the Baseline source and guides
  ("Autorizar os dois"), every skill change ("considere autorizado a ajustar
  todas as skills se necessário") and, for the Governed Paths this Spec
  declares, "Concedo" (2026-10-06). That grant covers the one-line
  `OWNED_SKILLS` edit in the `Makefile`. Bounded files:
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `docs/agents/skill-dispatch.md`, `docs/agents/setup-context.json`,
  `Makefile`, `.agents/skills/council/SKILL.md`,
  `.agents/skills/council/assets/synthesis-template.md`,
  `.agents/skills/council/references/archetypes.md`,
  `.agents/skills/council/references/debate-protocols.md`,
  `skills/council/SKILL.md`, `.agents/skills/write-idea/SKILL.md`,
  `.agents/skills/write-idea/references/idea-template.md`,
  `skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
  `skills/write-prd/SKILL.md`, `.agents/skills/the-fool/SKILL.md`,
  `.agents/skills/the-fool/references/cognitive-bias-inventory.md`,
  `.agents/skills/the-fool/references/dialectic-synthesis.md`,
  `.agents/skills/the-fool/references/evidence-audit.md`,
  `.agents/skills/the-fool/references/mode-selection-guide.md`,
  `.agents/skills/the-fool/references/pre-mortem-analysis.md`,
  `.agents/skills/the-fool/references/red-team-adversarial.md`,
  `.agents/skills/the-fool/references/socratic-questioning.md`,
  `.agents/skills/autoresearch/SKILL.md`,
  `.agents/skills/autoresearch/references/eval-guide.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

1. No module, trigger, activation bundle or Setup Snapshot names `council`
   or `the-fool`, and the asset sync against upstream `b3c45a4` keeps them
   out while every other skill still follows its upstream list.
2. The binary no longer ships `council`, and no owned skill sends the reader
   to it.
3. This repository holds no lock entry, recommended-list line or tree for
   `the-fool`, `autoresearch` or `council`, and keeps `grilling`,
   `grill-with-docs`, `write-idea`, `business-analyst` and `handoff`.
4. An adopter that still holds a retired copy learns from `baseline update`
   that it is no longer required and which paths to delete, and nothing is
   deleted for it.

## User Stories

1. As the maintainer, I refresh the snapshots from upstream and the retired
   skills stay out without editing upstream.
2. As an adopter, my next `baseline update` drops the two triggers from my
   guide and tells me which retired copies I can delete.
3. As an Agent running `write-idea`, I weigh the trade-offs in step 6 without
   a skill that no longer ships.

## Core Features

1. Core Feature: the Retired Skill set (`council`, `the-fool`); the asset
   sync drops a Retired Skill from every snapshot, and the catalog tests
   refuse a module, trigger, bundle or snapshot that names one (task_01).
2. Core Feature: `context-workflow` stops requiring and dispatching the two
   skills; snapshots, parity fixture, derived pins and this repository's
   guides are regenerated (task_01).
3. Core Feature: `council` leaves the owned bundle, its two copies, the
   owned-version record and the `Makefile` list; `write-idea` and
   `write-prd` stop naming it (task_02).
4. Core Feature: this repository drops `the-fool` and `autoresearch` from
   its lock, recommended list and `.agents/skills/`, and its doctor line and
   repository skill-set check follow (task_02).
5. Core Feature: `baseline update` lists Retired Skills the repository still
   holds under `Skills retired`, read-only; the user guide, the Roundfix
   Skill's baseline reference and `CONTEXT.md` describe it (task_03).

## Non-Goals / Out of Scope

- Editing upstream `marcioaltoe/skills` or `~/dev/skills`, or the contents of
  any vendored skill; the vendored trees are only deleted.
- Retiring `grilling`, `grill-with-docs`, `write-idea`, `business-analyst` or
  `handoff`, which the maintainer kept.
- Deleting an adopter's installed copy or editing its lock.
- Making Doctor report retired copies; Doctor keeps ignoring unrequired
  installed skills.
- Treating `autoresearch` as a Retired Skill: no adopter received it from the
  Baseline.
- `CHANGELOG.md`, which the release writes from the Release note below.

## Success Metrics

1. Success Metric: `TestNoCatalogEntryNamesARetiredSkill` passes on the
   embedded catalog, and a planted `council` or `the-fool` in a required
   list, dispatch entry, trigger, bundle or snapshot is reported.
2. Success Metric: `roundfix baseline assets sync --check --source-dir
   <clone>/setups` against upstream `b3c45a4` exits `0`, the four snapshots
   list neither name, and a sync over an upstream list that names both drops
   them, including a recorded repository-owned `council` entry.
3. Success Metric: `skills.Names()` has 13 names without `council`, the
   embedded bundle has no `council` tree, and no file of `write-idea` or
   `write-prd` contains "council".
4. Success Metric: in this repository Doctor prints `skills: ok (41
   required: 13 Roundfix-owned, 28 external)`, no lock entry, recommended
   line or `.agents/skills/` directory names `the-fool`, `autoresearch` or
   `council`, and `baseline update --repo .` reports `current`.
5. Success Metric: `baseline update` on a repository holding
   `.agents/skills/council` and a locked `.agents/skills/the-fool` lists both
   under `Skills retired` with their paths and keeps its state and exit code;
   with `--no-skills`, or without retired copies, no `Skills retired` line or
   `skills.retired` field appears.
6. Success Metric: `docs/user-guide/commands/baseline.md` gives the removal
   commands, and `CONTEXT.md` defines **Retired Skill**.

## Declared breaks

- **Break — guides.** Every adopter's next update removes
  `trigger.context-workflow.council` and `trigger.context-workflow.the-fool`
  from `docs/agents/skill-dispatch.md`, and its Setup Manifest records the new
  catalog digest.
- **Break — required skills.** Adopters no longer require `council` or
  `the-fool`; Doctor's counts drop by one owned and one external skill.
- **Break — owned bundle.** `roundfix skills list` and `roundfix skills
  install` no longer carry `council`. An installed copy stays until the
  adopter deletes it.
- **Break — `write-idea`.** Step 6 no longer runs `council`; the
  `_idea.md` template section `Council Insights` becomes `Trade-off Insights`.
- **Break — snapshots.** The snapshots now differ from the upstream setup
  lists by the two retired names, as ADR-0246 decides.
- **Break — parity fixture.** `asset-sync.json` records three setups without
  the two names (sanctioned, ADR-0072).
- **Break — update output.** The `baseline update` JSON gains an optional
  `skills.retired` array, and the text gains a `Skills retired` block only
  when a retired copy exists.

## Release note

The v0.52.0 entry of `CHANGELOG.md` carries this text, which
`docs/user-guide/commands/baseline.md` repeats:

> **Retired skills.** The Baseline no longer requires `council` or
> `the-fool`, and the binary no longer ships `council`. `roundfix baseline
> update` removes their triggers from `docs/agents/skill-dispatch.md`, never
> deletes an installed copy, and lists any copy it finds under `Skills
> retired`. To remove them, run
> `git rm -r -q --ignore-unmatch .agents/skills/council .agents/skills/the-fool .claude/skills/council .claude/skills/the-fool`,
> then `rm -rf .agents/skills/council .agents/skills/the-fool .claude/skills/council .claude/skills/the-fool`
> for untracked copies, delete the `the-fool` entry from `skills-lock.json`
> by hand, and commit. The next `roundfix baseline update` lists no retired
> skill.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- The upstream skills repository at
  `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`: `setups/go.txt`,
  `setups/rust.txt` and `setups/typescript.txt` still list
  `skills/01-discovery/council` and `skills/01-discovery/the-fool`. The QA
  gate runs the asset sync `--check` against a local clone of that commit.
- The skills CLI's own issue tracker: vercel-labs/skills issue #977
  (<https://github.com/vercel-labs/skills/issues/977>) reports that
  `skills remove` in project scope deletes the tree and links but leaves the
  `skills-lock.json` entry, which is why the removal instruction deletes the
  lock entry by hand.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "retire skill from baseline adopters keep installed copy"`
(`--all --files --min-score 0.3`). It returned this repository's mirrors of
the adopted entry, ADR-0206, ADR-0221, ADR-0224 and the 2026-10-01 Finding on
the reshaped catalog; none decides a skill the Baseline drops while upstream
keeps it. Exa found the skills CLI's command reference
(<https://github.com/vercel-labs/skills/blob/main/src/cli.ts>), whose
`remove` "Remove installed skills from agents", and issue #977 above. The
two other open Backlog Entries of 2026-10-06 stay apart: the history
clean-up changes `docs/history/`, and the CONTEXT-driven loop entry changes
the clauses that make implementations keep `CONTEXT.md` current; neither
changes the skill set.

## Decisions

- Retired Skills are a set in code; the asset sync drops them and the catalog
  tests refuse them. See ADR-0246.
- Adopters keep their copies; `baseline update` lists them and the guide
  gives the removal commands. See ADR-0246.
- `write-idea` keeps its trade-off step and does it itself.
- `autoresearch` leaves this repository only.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the measured outputs, the fixed
texts and the build order.
