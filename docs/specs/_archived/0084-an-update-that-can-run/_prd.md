---
spec: 0084-an-update-that-can-run
status: active
created: 2026-08-08
surfaces: [backend, cli, docs]
---

# A Baseline Command update that runs on the repositories that already exist

Spec 0082 shipped `roundfix baseline update` so a maintainer could refresh a
repository's Context-Driven Baseline without re-answering every setup question.
Measured against the eight repositories that have actually adopted a Baseline,
it refuses to run on six of them. Two of those refusals are this Spec's subject:
the command reads a legitimately current managed region as damage and stops
before it plans anything. Three more stop on structural clauses the catalog
silently stopped emitting. One stops on a Baseline Profile its Setup Manifest
names and its checkout no longer contains.

None of these are recoverable by the maintainer, because every one of them
blocks before the command produces the plan that would repair it. The outcome
this Spec buys is the one 0082 promised and did not deliver: a maintainer sweeps
the fleet, and each repository either applies its refresh or reports a condition
a human can act on.

## Project Constraints

- Identifier strategy: not applicable — this Spec introduces no new entity or
  resource identity; managed artifacts keep the Setup Manifest identities they
  already carry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the Baseline Command is a local
  filesystem operation with no network surface and no HTTP contract. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0058 keeps retention accounting
  fail-closed and this Spec does not weaken it; ADR-0070 governs root backups
  and ADR-0100 replaces them with preimage proof on the managed-refresh path,
  which this Spec preserves; ADR-0074 and ADR-0078 govern repository-rule
  ownership and are unaffected, because no repository rule changes owner here;
  ADR-0081 governs sanctioned digest regeneration; ADR-0091 keeps the authored
  QA gate inside the Task Graph and ADR-0096 keeps that gate proving machine
  facts before spending an Agent turn; ADR-0095 keeps the Secondbrain inbox the
  single capture door, which this Spec's consultation clause reads from and does
  not widen; ADR-0097 governs QA row carry-forward; ADR-0099 keeps retention
  accounting mechanical. New decisions land as ADR-0101, ADR-0102, ADR-0103, and
  ADR-0104. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected Baseline catalog assets and two
  repo-owned authorial skills are mutated. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-07-restore-structural-clauses.md`,
  `docs/workflow/authorizations/2026-08-08-the-brain-is-a-source-not-an-archive.md`,
  `docs/workflow/authorizations/2026-08-08-evidence-from-outside-the-spec.md`,
  `docs/workflow/authorizations/2026-08-08-glossary-currency-clause.md`, and
  `docs/workflow/authorizations/2026-08-08-the-skill-ships-with-the-cli-change.md`.
  Bounded files: `internal/baseline/assets/modules/*.json`,
  `internal/baseline/assets/retention/**`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/domain.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, and
  `.agents/skills/roundfix/SKILL.md`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. Every repository that has adopted a Baseline can be refreshed by one
   non-interactive command, or reports a condition naming what a human must do.
2. A refresh converges: applying it once makes the next run report the
   repository current, with no further changes proposed.
3. A managed region whose bytes are not the ones the Setup Manifest recorded is
   reported by path before it is replaced, and its replacement never removes a
   line the Baseline cannot reproduce without naming that line.
4. The catalog emits every structural Normative Clause the Source Baseline
   accounts for, so retention accounting stops blocking on the catalog's own
   omissions.
5. A Spec's acceptance rests, in at least one named row, on evidence the Spec did
   not author.

## User Stories

1. As a maintainer, I want to run one update command in each repository of my
   fleet, so that a Baseline change reaches every project without me answering
   the same setup questions nine times.
2. As a maintainer, I want a repository whose managed regions moved ahead of its
   Setup Manifest to refresh anyway, so that adopting the Baseline before the
   update command existed is not a permanent block.
3. As a maintainer, I want to see which managed regions carried bytes the Setup
   Manifest did not record, and which lines the refresh removed, so that I can
   approve the refresh knowing what it replaces.
4. As a maintainer, I want a second run of the update to tell me the repository
   is current, so that I can tell a finished sweep from an unfinished one.
5. As a maintainer, I want a repository whose recorded Baseline Profile no longer
   resolves to name the missing profile and the action that restores it, so that
   a broken repository is diagnosable without reading Go source.
6. As a Supervisor authoring a Spec, I want the Secondbrain treated as a source
   for strategy, sibling-project decisions, and literature, so that a Spec is
   informed by what the ecosystem already learned rather than by this repository
   alone.
7. As a Supervisor authoring a Task Graph, I want at least one acceptance row
   backed by evidence from outside the Spec, so that a gate cannot pass by
   confirming the requirement that produced it.

## Core Features

1. **Managed regions are classified, not presumed damaged.** A managed region
   whose bytes differ from the Setup Manifest's recorded digest is classified as
   *unrecorded* and reported by path and managed identity. Planning continues.
   Only a region that cannot be read, parsed, or matched to a managed identity
   remains blocking.
2. **A refresh names what it removes.** For every unrecorded managed region, the
   presented plan reports each line present on disk that the refreshed rendering
   does not reproduce. A refresh that removes no such line is reported as
   carrying none.
3. **Approval covers replacement.** Applying a refresh still requires the
   existing explicit approval of a Plan Digest; the classification changes what
   the maintainer is shown, never whether a human approves.
4. **The update converges.** An applied refresh republishes the Setup Manifest so
   the recorded digests describe the bytes on disk. A second run against an
   unchanged catalog reports the repository current and proposes no change.
5. **Structural clauses are emitted again.** The catalog declares the fourteen
   structural Normative Clauses the Source Baseline accounts for, so retention
   accounting on an unchanged repository no longer reports them unaccounted.
6. **An unresolved Baseline Profile reports its own repair.** A Setup Manifest
   naming a Baseline Profile the checkout cannot resolve reports the profile
   identity, where the command looked for it, and the action that restores it.
7. **The Secondbrain is a consultation source.** The Baseline obliges consulting
   the Secondbrain when a Spec is authored or an approach is chosen, for prior
   decisions from sibling projects, literature, and technical knowledge the
   repository does not hold. An unreachable Secondbrain is a reported condition,
   never a blocker.
8. **Acceptance carries outside evidence.** The Baseline obliges every Spec to
   rest at least one named acceptance row on evidence originating outside its own
   artifacts, and to record that evidence's origin.
9. **The glossary is checked when work closes.** The Baseline obliges observing,
   at the close of a Spec, feature, refactor, or fix, whether the work introduced,
   changed, or retired a term the project glossary should carry, and updating it
   when it did. Human interaction is not required.

## User Experience

The command's text output gains one block in the presented plan, listing each
unrecorded managed region by path and managed identity, and under it the lines
the refresh removes. Its JSON result gains the same information as structured
fields. No flag is added and no prompt is introduced; a repository that would
have blocked now presents a plan the maintainer approves with the flags that
already exist.

## Non-Goals / Out of Scope

- Interactive repair of a hand-edited managed region. The refresh replaces it and
  reports what it replaced; editing it back is the maintainer's act.
- Reconstructing the historical catalog that rendered an unrecorded region. The
  command classifies against what it can render now, never against lineage it
  cannot observe.
- Changing the greenfield or preservation modes, the semantic analyzer, or the
  classification flow.
- Changing which decisions a profile asks, or adding new project decisions.
- Restoring a Baseline Profile a checkout is missing. The command reports it;
  authoring it is separate work.
- Automatic Secondbrain consultation performed by the tool. The obligation is on
  the session authoring a Spec, not on the binary.
- The remaining baseline reform items already queued for Spec 0085.

## Success Metrics

- Eight of eight adopted repositories in the maintainer's fleet reach either an
  applied refresh or a reported condition naming a human action; zero reach a
  state that blocks before planning for a reason the maintainer cannot act on.
- Running the update twice against an unchanged catalog leaves the second run
  reporting `current` with zero file changes, in every repository the first run
  applied.
- Zero authored lines are lost across the fleet sweep: for every applied refresh,
  the removed-line report accounts for every line the diff removed.

## Decisions

- A managed region's trust comes from its marker and the plan's preimage, not
  from a digest recorded on adoption day. See ADR-0101.
- An unrecorded managed region is refreshed and named rather than blocking. See
  ADR-0102.
- An applied refresh republishes the Setup Manifest, so the update converges on
  a second run. See ADR-0103.
- A Spec's acceptance must include at least one row resting on evidence the Spec
  did not author. See ADR-0104.
- Spec 0082's Task 02 Requirement 4 — "MUST treat a hand-edited managed marker as
  blocking rather than as a warning" — is superseded by this Spec. Its other
  seven requirements stand.
- The clauses committed outside a Spec pipeline on 2026-08-07 and 2026-08-08
  remain in scope for Spec 0085. This Spec takes only the three clauses whose
  authorizations name the defect it fixes.

## Open Questions

- Whether the fourteen restored structural clauses need retention accounting
  entries for identity moves, or restore cleanly to their owning modules. The
  default until measured is that they restore cleanly and the accounting
  directory is untouched.
