---
spec: 0253-authoring-rules-that-stop-qa-reruns
status: active
created: 2026-10-08
surfaces: [backend, docs]
---

# Authoring rules that stop QA reruns

The Baseline audit of 2026-10-08 measured Specs 0225–0248 in the Run Database
([the committed audit](../../references/2026-10-08-baseline-audit.md),
findings B08–B14 and B16, and its section on the authoring briefing). The QA
gate ran about 1.5 times per Spec. Five reruns came from Surface Transcript
lines the implementing test never asserted, such as a final Plan Digest line,
`exit status 1`, indentation or a `Usage` block, while the code was correct.
Four of the last eleven Specs needed a QA Archive Override because an
outside-evidence row pointed at an operator-only log, a lookup the
authorization forbade or a host the sandbox denied. Stale wording in a second
guide failed QA or CI four times. Tests that read a host key or bound a long
macOS socket path parked a Run or forced a corrective Task. Specs authored in
parallel needed hand fixes after a sibling merged an ADR, and a corrective Task
needed a hand-added `needs` edge before `roundfix reopen` saw it. The Baseline
also still tells a retirement to move whole Review Artifacts and handoffs into
`docs/history/`, which the next sanitize deletes.

After this Spec, the Baseline and the authoring skills turn each of those
classes into a rule an author meets before the Run, the repository-owned guide
carries the Roundfix-only authoring rules the operator kept outside the
repository, and a retirement writes reduced history directly. The maintainer
asked for this on 2026-10-08: "Atualizar e verificar o baseline de regras para
que tenhamos o máximo de performance e resultado com o roundfix e congruente
com as mudanças", and approved "Auditoria e depois Spec (Recommended)", with the
audit's Specs delivered without further consultation within the bound of four
implementation Tasks. Of history the maintainer said "Não quero nada no
histórico que não seja relevante para o secondbrain".

## Prerequisites

Spec 0252 (`0252-one-gate-per-run-and-a-smaller-baseline`) delivers first. It
changes the `core`, `spec-workflow`, `autonomous-work`, `secondbrain` and
`context-workflow` modules, the guides and `docs/agents/specific-repository.md`.
This Spec rebases on its merge and edits the state its TechSpec describes. Its
Task Graph names 0252 in `requires`, so the Delivery Queue does not start this
Spec before 0252 merges. No Verification of this Spec reads an artifact 0252
creates: every clause this Spec changes is one 0252 leaves byte-identical, and
each module Task records the next version from the value on its starting
commit.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. Every
  reworded clause keeps its identity, and the one new clause takes a new
  identity in the Go module. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Tests and Verification use the embedded catalog, temporary
  repositories and this repository's own files. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0258 (this Spec) decides the
  authoring rules and the reduced retirement. ADR-0257: "A clause whose wording
  changes keeps its identity", so every Baseline rule here extends a clause in
  place except the new Go clause. ADR-0186: "Every clause keeps its
  enforcement level when its wording changes", and "It never cites a Spec
  number or an ADR number"; no new clause text does. ADR-0250: "A Baseline
  module's version names one content, and the record step chooses it", so each
  module Task records its modules and names no version number. ADR-0248: a
  retired Finding or Backlog Entry keeps "its front matter unchanged, its
  title, its first paragraph and the revision that holds its full text", and
  "Retired Review Artifacts and handoffs have no reader, so they are removed";
  this Spec applies both at retirement. ADR-0254: "So is a unit whose paths an
  existing History Full Tag does not hold", which is why a full entry written
  into history after the tag would be refused. ADR-0244: "A Spec declares the
  domain terms it introduces"; task_04 revises **Reduced History Entry**.
  ADR-0251 reads a Legacy Archive Folder leniently and keeps a
  failed QA's verdict; this Spec changes neither, and the sanitize it describes
  keeps its contract.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization covers the
  Baseline modules, the rendered guides and Setup Manifest, the owned skills
  and the repository-owned guide: "considere autorizado a ajustar todas as
  skills se necessário", "Autorizar os dois", the standing "Concedo", and the
  2026-10-08 grants for the governed paths this Spec declares. No Makefile,
  `.roundfixrc.yml`, CI workflow, lint, formatter or test-runner configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0253-authoring-rules-that-stop-qa-reruns/_authorization.md`.
  Bounded files: `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/concrete-contracts.md`,
  `docs/agents/autonomous-work.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- An outside-evidence row names a source the QA gate can check without the
  network, and a criterion that needs the network or an open Pull Request is
  declared under Unreachable Acceptance when it is authored.
- The Task that implements a Surface Transcript asserts every line of its
  standard output and standard error and its exit code in one named test that
  its Verification runs.
- A Spec that depends on an unmerged Spec names it in `requires`.
- The autonomous loop re-checks a Spec authored in parallel after rebasing,
  adds a corrective Task to the gate's `needs` before reopening it, sweeps old
  wording and characterizes behavior before changing it, and Go tests stay
  hermetic.
- The repository-owned guide states the suite guard rule and the Roundfix-only
  authoring rules: the `### QA settlement` section, clause retention, skill and
  module records, derived files and the measured Governed Path set.
- A retired Finding or Backlog Entry is written to history in reduced form,
  and a finished Review Artifact or a confirmed handoff is deleted.

## User Stories

1. As a QA gate, I want every outside-evidence source reachable without the
   network before the Run starts, so that I never end `partial` on a log only
   the operator holds.
2. As the operator, I want the implementing test to assert a transcript whole,
   so that a missing `exit status 1` line fails the Task, not the QA gate.
3. As the Delivery Queue, I want a Spec's prerequisites in `requires`, so that
   I start it only after the Spec it depends on has merged.
4. As a Supervisor, I want the loop to tell me to re-check a Spec authored
   beside another and to add a corrective Task to the gate's `needs`, so that
   neither needs a hand fix.
5. As a Task agent, I want the rules that hold only for Roundfix in its own
   guide, so that I meet them without the operator's private briefing.
6. As the maintainer, I want history to hold only what the Secondbrain needs,
   so that a retirement never adds a whole review or handoff to `docs/history/`.

## Core Features

1. **Reachable outside evidence.** The outside-evidence clause names the
   sources the QA gate can check without the network and the sources it never
   rests on, and sends a network or Pull Request criterion to Unreachable
   Acceptance at authoring (ADR-0258).
2. **Whole transcripts and declared prerequisites.** The Task Graph clause
   requires one named test that asserts a Surface Transcript whole. The
   grouping clause requires `requires` and a Build Order entry for a
   dependency on an unmerged Spec (ADR-0258).
3. **Authoring skills.** `write-tasks` and the `write-techspec` concrete
   contracts guide say the same, name the lines authors miss, and keep one
   line on backticks in Verification commands.
4. **Loop and Go rules.** The loop clause adds the parallel re-check and the
   corrective `needs` edge. The class clause adds the wording sweep and
   characterization before change. A new Go clause keeps tests hermetic
   (ADR-0258).
5. **Repository rules.** `docs/agents/specific-repository.md` gains the suite
   guard rule and the briefing rules on the QA settlement section, clause
   retention, owned skills, derived files, sanctioned regeneration and the
   Governed Path probe.
6. **Reduced retirement.** The docs-layout clauses write a retired Finding or
   Backlog Entry as a reduced entry and delete a finished Review Artifact or
   confirmed handoff. The user guide follows (ADR-0258).
7. **Glossary.** `CONTEXT.md` revises **Reduced History Entry**.

## Non-Goals / Out of Scope

- Spec 0252's findings (B01–B07, B15, B19) and the audit's B17, B18 and B20.
- Changing History Relocation, `roundfix history sanitize`, `roundfix
  reopen`, `roundfix spec check`, the Delivery Queue or any production Go file.
- A command that writes a reduced entry; the guide states the form.
- New clause identities in a module a Source Baseline holds, or any edit to
  the Source Baseline corpus, manifest, index or retention transition.
- Trimming the operator's briefing, which lives outside the repository.
- Editing any adopter repository.

## Success Metrics

1. Success Metric: this repository's `docs/agents/spec-routing.md` carries the
   reachable-source, whole-transcript and `requires` sentences, each once, and
   a Source Baseline adopter's Managed Refresh records the three reworded
   clauses `retained` with no clause `unaccounted`.
2. Success Metric: this repository's `docs/agents/autonomous-work.md` carries
   the parallel re-check, the corrective `needs` edge, the wording sweep and
   characterization before change. `docs/agents/go.md` carries the hermetic
   test clause as a `mandatory` bullet.
3. Success Metric: `write-tasks`, the `write-techspec` concrete contracts
   guide and `docs/agents/specific-repository.md` carry the rules of Core
   Features 3 and 5. Both skills record a raised version, and each mirror
   matches its canonical copy.
4. Success Metric: the rendered docs-layout guide names no move of a new
   retirement into `docs/history/reviews/` or `docs/history/handoffs/`, and a
   Finding written in the guide's reduced form is one the History Sanitize
   Command plans no change for.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Baseline audit of 2026-10-08, committed by Spec 0252's planning change as
  `docs/references/2026-10-08-baseline-audit.md`. It is the operator's
  measurement of Specs 0225–0248: the five transcript reruns, the four
  overrides in Specs 0236–0248, the four stale-wording failures, the hermeticity
  and suite guard incidents and the parallel-authoring fixes. The QA gate reads
  it without the network.
- Luo, Hariri, Eloussi and Marinov, "An Empirical Analysis of Flaky Tests"
  (FSE 2014). Of the order-dependent flaky tests they studied, "Many Test Order
  Dependency flaky tests (47%) are caused by dependency on external
  resources".
- Gruber et al., "An Empirical Study of Flaky Tests in Python" (ICST 2021,
  arXiv:2101.09077): "Another 28 % were caused by test infrastructure
  problems", and 13 % mostly by network and randomness APIs.
- The Secondbrain concept `wiki/concepts/ciclo-de-vida-de-decisoes-e-raiz-docs-history.md`,
  read during authoring, on `docs/history/` as a ledger of identities.

## Unreachable Acceptance

- criterion: fewer QA reruns from transcript drift and fewer QA Archive
  Overrides from unreachable evidence in Specs delivered after release
  reason: the effect shows only in Runs that read the delivered guides, which
  start after this Spec merges, and the Run Database is outside the QA
  sandbox
  satisfied-by: the operator counts, over the first five Specs delivered after
  this release, the QA reruns whose only defect was a transcript line and the
  archives with `qa_override: true`, and compares them with the audit's five
  reruns in 0219–0247 and four overrides in 0236–0248

## Glossary

- changes: **Reduced History Entry**

## Decisions

- Each rerun class becomes a rule in an existing clause, and hermetic Go tests
  get a new clause; see ADR-0258.
- A retirement writes a reduced entry or deletes, and the sanitize keeps its
  contract for older history; see ADR-0258.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the clause changes, the skill and
guide texts, the version changes, the derived files and the build order.
