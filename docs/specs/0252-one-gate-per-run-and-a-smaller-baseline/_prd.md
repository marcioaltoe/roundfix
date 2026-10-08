---
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: active
created: 2026-10-08
surfaces: [backend, docs]
---

# One gate per Run and a smaller Baseline

The Baseline audit of 2026-10-08 measured Specs 0225–0248 in the Run Database
([the committed audit](../../references/2026-10-08-baseline-audit.md), findings
B01–B07, B15 and B19). Task agents inside Runs ran `make verify-incremental`,
which is the full `go test` suite, 74 times in 20 of 23 Specs: 284
agent-minutes, 231 s each on average. They did so because the Baseline says
"For each Task, run the selected incremental Verification ... before handoff".
Meanwhile the Daemon ran `make verify-changed` at every settlement (125 runs,
350 s on average), and those runs failed only 6 times. The guides also record
`rtk`-prefixed Verification commands that CI never runs, name models that no
Agent Selection Profile uses, state the review policy, the tracker rule and the
local-research rule two or three times, and point every session at the
docs-layout and Secondbrain guides (about 23 KB of the 69 KB of guidance)
whether the work needs them or not.

After this Spec, a Daemon-assigned Task runs focused tests and leaves the full
gate to the Daemon. The recorded Verification commands are the ones CI runs,
the autonomous-work guide says that Agent Selection Profiles choose models,
each rule is stated once, and the root tells a session when it needs the
docs-layout and Secondbrain guides. The maintainer asked for this on
2026-10-08: "Atualizar e verificar o baseline de regras para que tenhamos o
máximo de performance e resultado com o roundfix e congruente com as mudanças".

## Prerequisites

Spec 0249 is being delivered, and Specs 0250 and 0251 deliver before this one.
None of them changes a Baseline module, a Baseline guide or
`docs/agents/setup-context.json`, and this Spec's Verification reads none of
their artifacts. Spec 0253 is authored after this one and changes the
`spec-workflow`, `autonomous-work` and `go` modules and `rule.autonomous.loop`.
This Spec delivers first, and Spec 0253 rebases on it and takes the next
version when it records its modules.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. Every
  clause that is reworded keeps its identity, and the decision identifiers
  keep their names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Tests and Verification use the embedded catalog, temporary
  repositories and this repository's own files. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0256: a Spec that changes behavior a
  skill describes plans the skill change; task_01 changes `implement-task` §7
  and records its version, so the covering skill moves with the guidance.
  ADR-0257 (this Spec) decides the session
  scoping of the two tiers, the recorded values, the runtime header, the merges
  and the root reading set. ADR-0222: "A clause whose wording changes takes a
  new identity that declares `replaces` for the old one". ADR-0257 narrows
  that rule: a reworded clause keeps its identity, and `replaces` names each
  clause another one absorbs. ADR-0186: "Every clause keeps its enforcement
  level when its wording changes", and every merge here joins clauses of the
  same enforcement. ADR-0250: "A Baseline module's version names one content,
  and the record step chooses it", so each Task records its modules and names
  no version number. ADR-0238: "A Task is `light` when its `complexity` is
  `low`, its `type` is not `qa`, and it declares no Governed Path"; the header
  states that rule. ADR-0252: "`make verify-changed` is the Verification of
  Runs and of the QA gate", which is the incremental value this repository
  records. ADR-0253: "CI runs the Full Contract Run on every push to main and in
  the release workflow before any build or publication", beside `make verify`,
  the gate value this repository records. ADR-0244: "A Spec declares the domain terms it introduces"; task_01
  revises **Incremental Verification**. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization covers the
  Baseline modules, the decision catalog, the guide and root templates, the
  rendered guides and Setup Manifest, two owned skills and the tests the
  historical record governs: "considere autorizado a ajustar todas as skills
  se necessário", "Autorizar os dois", the standing "Concedo", and the
  2026-10-08 grants for the governed paths this Spec declares. No Makefile,
  `.roundfixrc.yml`, CI workflow, lint, formatter or test-runner configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0252-one-gate-per-run-and-a-smaller-baseline/_authorization.md`.
  Bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/autonomous-work.md`, `docs/agents/secondbrain.md`,
  `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`,
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/issue-tracker.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/guides/autonomous-work.md`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/root/context-workflow.md`,
  `internal/baseline/assets/templates/root/secondbrain.md`,
  `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go`,
  `skills/baseline_skill_contract_test.go`,
  `skills/implement-task/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Goals

- A Daemon-assigned Task's guidance tells its Agent to run focused tests of
  the changed packages and neither the incremental nor the repository
  Verification. Outside a Run the two tiers apply as before, and the strict
  Spec check comes before a Run.
- This repository records `make verify` and `make verify-changed`, the
  commands CI runs, and the catalog proposes no `rtk`-prefixed command.
- The autonomous-work guide says that Agent Selection Profiles choose each
  session's runtime and model, names the Light Tier rule and renders current
  preferences, while both runtime decisions stay.
- The review policy, the tracker rule and the local-research rule are each
  stated once. Every adopter's update records the removed clauses `replaced`,
  never `unaccounted`.
- The root block makes the docs-layout guide conditional on documentation
  work and the Secondbrain guide conditional on Secondbrain or authoring work.

## User Stories

1. As a Task agent in a Run, I want the guides to agree with
   `implement-task`, so that I stop spending four minutes on a suite the
   Daemon runs again.
2. As the operator, I want the recorded Verification commands to be the ones
   CI and the Daemon run, so that the guides do not name a local filter.
3. As an adopter, I want the autonomous-work guide to tell me where models
   are chosen, so that a stale model name never reads as a rule.
4. As a reader of the guides, I want each rule once, so that two copies never
   drift apart.
5. As a Run session, I want the root to tell me which guides my work needs,
   so that I do not read 23 KB of templates and Secondbrain rules a code Task
   never uses.

## Core Features

1. **Tiers scoped by session.** The core and Spec workflow clauses keep
   their identities. They scope the incremental and repository Verification to
   work outside a Run and send a Daemon-assigned turn to focused tests, because
   the Daemon runs both checks at settlement. The authoring clause names
   `--strict`. `implement-task` says the same in its Daemon-assigned handoff
   (ADR-0257).
2. **Recorded values.** The catalog's gate default and incremental suggestion
   drop `rtk`. This repository's Setup Manifest records `make verify` and
   `make verify-changed`. The user guide and the `setup-context-driven` skill
   show the same values (ADR-0257).
3. **Runtime header.** The autonomous-work guide states that Agent Selection
   Profiles choose each Agent Session's runtime, model and effort, and that a
   low-complexity Task that is not `qa` and changes no Governed Path runs on
   the Light Tier. It renders the two runtime decisions as preferences, with
   defaults `codex gpt-6.1-sol high` and `claude opus high`. Neither decision
   is retired (ADR-0257).
4. **One statement per rule.** The two core review clauses become one. The
   autonomous loop's review sentence drops what core states. The tracker
   clause absorbs the status clause, and the legacy-Spec clause drops its
   ownership tail. The core local-research prohibition absorbs the Secondbrain
   copy, and the Secondbrain consultation clause drops its repeated query
   order. The repository-owned skill-sync rule, which a CLI clause already
   covers, is deleted (ADR-0257).
5. **Root reading set.** The CONTEXT-driven root block keeps the domain guide
   mandatory. It names the docs-layout guide for documentation work and for a
   test or build step that reads a document. The Secondbrain root block names
   its guide for Secondbrain and authoring work (ADR-0257).
6. **Glossary.** `CONTEXT.md` revises **Incremental Verification**.

## Non-Goals / Out of Scope

- The audit's Spec B findings (B08–B14), B16, B17, B18 and B20, and the
  briefing migration. Spec 0253 owns B08–B14. The others stay in the audit as
  later candidates.
- Changing `loop-01`, `loop-04` or any `go` module clause.
- Retiring `runtime.backend` or `runtime.design`, changing how Roundfix
  selects a model, or changing `.roundfixrc.yml`.
- Changing the Makefile, `make verify-changed`'s duration, CI workflows or
  the Daemon's Verification behavior.
- Splitting the docs-layout record templates into their own guide.
- Renaming a clause identity, or editing the Source Baseline corpus,
  manifest or index.
- Editing any adopter repository.

## Success Metrics

1. Success Metric: in this repository's rendered guides, every sentence that
   asks for the incremental or repository Verification is scoped to work
   outside a Run. The agent-instructions and spec-routing guides each name
   focused tests for a Daemon-assigned turn, and `implement-task` §7 forbids
   both selected commands.
2. Success Metric: `docs/agents/setup-context.json` records `make verify` and
   `make verify-changed`. No catalog default or suggestion begins with `rtk`,
   and `roundfix baseline update --repo . --no-skills` reports `current`.
3. Success Metric: the rendered autonomous-work guide names
   `roundfix profiles show`, the Light Tier and the two recorded preferences.
   It names no model that `.roundfixrc.yml` does not use.
4. Success Metric: a Source Baseline adopter's Managed Refresh is ready and
   records `clause.core.request-pull-request-review`,
   `clause.spec.status-only-in-task` and
   `clause.secondbrain.prohibit-external-local-discovery` `replaced`, with no
   clause `unaccounted`. `TestNoTwoBaselineClausesShareText` passes.
5. Success Metric: the reading set of a session whose work touches no
   document under `docs/` and no Secondbrain path is `AGENTS.md` plus every
   `docs/agents/*.md` other than `docs-layout.md` and `secondbrain.md`. It is
   at least 20,480 bytes smaller than the 69,330 bytes that `AGENTS.md` and
   every `docs/agents/*.md` hold at `55de2a73`. The authoring prototype
   measured 45,435 bytes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Baseline audit of 2026-10-08, committed by this Spec's planning change
  as `docs/references/2026-10-08-baseline-audit.md`. It is the operator's
  measurement of Specs 0225–0248 in the Run Database: 74 agent-side
  full-suite runs, 284 agent-minutes, 125 Daemon `verify-changed` runs with a
  350 s mean and 6 failures, the guide sizes, and the CI commands in
  `.github/workflows/ci-verify.yml`.
- Gloaguen et al., "Evaluating AGENTS.md: Are Repository-Level Context Files
  Helpful for Coding Agents?" (arXiv:2602.11988, 2026). Agents follow
  context-file instructions. Those instructions lead to more testing and
  exploration, and cost rises by over 20 % with no significant gain in
  success. The authors conclude that context files "should describe only
  minimal requirements".
- "On the Impact of AGENTS.md Files on the Efficiency of AI Coding Agents"
  (arXiv:2601.20404, 2026) challenges a blanket removal: with an AGENTS.md,
  median agent runtime fell by 28.64 % on 124 pull requests. This Spec keeps the
  root and its mandatory guides and makes only two guides conditional.
- The Secondbrain concept
  `wiki/concepts/arquitetura-de-instrucoes-e-progressive-disclosure.md`, on
  the instruction budget, stale context and progressive disclosure.

## Unreachable Acceptance

- criterion: fewer agent-side full-suite runs and shorter Task sessions in
  Runs after release
  reason: the effect shows only in Runs that read the delivered guides, which
  start after this Spec merges, and the Run Database is outside the QA
  sandbox
  satisfied-by: the operator re-runs the audit's `run_events` query over the
  first five Specs delivered after this release and compares the count of
  agent-run `make verify-incremental` and `make verify` with the audit's 74
  runs in 23 Specs

## Glossary

- changes: **Incremental Verification**

## Decisions

- In a Daemon-assigned turn the Agent runs focused tests only; see ADR-0257.
- This repository records `make verify` and `make verify-changed`; see
  ADR-0257.
- The runtime decisions stay, with current defaults, and the guide defers to
  Agent Selection Profiles; see ADR-0257.
- A reworded clause keeps its identity, and a removed clause is named in the
  absorbing clause's `replaces`; see ADR-0257.
- The docs-layout and Secondbrain guides are read when the work needs them;
  see ADR-0257.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the clause texts, the version
changes, the derived files and the build order.
