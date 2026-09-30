---
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: active
created: 2026-09-30
surfaces: [backend, docs]
---

# Baseline guides that describe the product as it is

The Context-Driven Baseline ships its clauses to every adopter as mandatory
guidance, and an Agent obeys them literally. An audit on 2026-09-30 found
clauses that describe an older product. The loop clause orders archive before
the pre-PR review, while the Delivery Queue reviews first. One clause fails a
tooling Task for any undeclared path, which the audit stopped doing. One
describes an `execution_approvals` record that no code reads. Four clauses cite
this repository's own Spec and ADR numbers, which name other records in an
adopter's repository. The backend guide renders one paragraph twice. The
adopted entries in [references/_index.md](references/_index.md) list each case
with its evidence.

An Agent that follows these clauses does the wrong thing in the order it is
told to, or stops for a refusal that no longer exists. This Spec makes each
clause say what the shipped product does and adds the checks that keep it so.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Clause, rule, guide
  and module identifiers are kept, and one duplicate clause entry is removed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets, local files and
  Git only; no credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0186 (this Spec) makes a Baseline
  clause state what the product does in adopter-neutral words. The clauses are
  rewritten to match decisions already in force: ADR-0104 holds Pull Request
  preparation on the outside-evidence row, ADR-0130 has the audit judge only
  Governed Paths, ADR-0166 records an undeclared ordinary path, ADR-0167 keeps
  the pre-PR Pull Request row from deciding a qualifying partial, ADR-0179
  makes an explicit empty `paths` list an operations-only grant, ADR-0014 makes
  the Daemon the verification authority, ADR-0092 separates evidence from
  intent, and ADR-0163 retires a finished orphan Review Artifact. None of them
  changes. ADR-0081 and ADR-0149 keep sanctioned regeneration outputs inside a
  grant, and ADR-0073 and ADR-0103 require this repository's managed refresh to
  apply in one transaction and converge to `current`. This Spec's gate is bound
  by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0117, ADR-0155 and ADR-0156,
  and ADR-0093 and ADR-0094 check its consistency by citation and artifact
  presence. ADR-0178 has the audit read the grant a Task commit ran under, and the
  unchanged second sentence of the bounded-execution clause stays true under
  it. Nine ADRs cite a listed ADR and do not apply, because this Spec touches
  none of their behaviors: ADR-0020 (the parsed prompt result), ADR-0038 (the
  one Verification repair), ADR-0057 (the Daemon as only writer of Task
  status), ADR-0127 (process residue) and ADR-0160 (the red repository gate)
  cite ADR-0014; ADR-0097 (carrying a QA row forward) cites ADR-0080; ADR-0168
  (the related-ADR gap) and ADR-0176 (citation checks read authored text) cite
  ADR-0093; ADR-0173 (citations a History Relocation breaks) cites ADR-0073.
  Reached only through those citations, and equally untouched: ADR-0056, ADR-0159, ADR-0170, ADR-0177.
  All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source and its generated guides ("Autorizar os
  dois"), recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/autonomous-work.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`,
  `internal/speccheck/backlog.go`, `internal/baseline/plan_test.go` (added for task_06). Sanctioned regeneration: `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An Agent that follows the loop clause performs the steps in the order the
  Delivery Queue performs them, and knows the three recovery acts the product
  offers: reopening a settled gate, Task Carry-Forward and Delivery Retry.
- A clause about authorization or evidence states the refusal the product
  makes today, and none describes a control that does not exist.
- An adopter's guide cites no record it cannot open and renders no paragraph
  twice.
- A later change that breaks any of these fails a test before it ships.

## Core Features

1. **No clause is rendered twice, and every clause keeps its force.** The
   backend module's duplicate entry is removed, so the paragraph appears once.
   A check refuses a catalog in which two clauses share their text. A second
   check records the enforcement level of every clause and fails when a clause
   appears, disappears or changes level without the record changing with it.
2. **The loop clause follows the Delivery Queue.** The clause orders the steps
   as: implement the graph with its gate, apply the pre-PR review policy,
   archive and commit the candidate, run the repository gate, open the Pull
   Request, verify the checks, merge. It names the Delivery Queue as the
   command that runs this order, `roundfix reopen` for a gate settled before a
   corrective Task, and Task Carry-Forward before a Spec amendment. A check
   compares the declared order with the order the Delivery Queue runs.
3. **Authorization and evidence clauses state today's refusals.**
   - A tooling Task fails only for a Governed Path outside its grant. An
     ordinary undeclared path is recorded and disclosed.
   - An absent or inoperative record still lets `roundfix implement` run a
     Spec that changes no Governed Path, and `roundfix deliver start` refuses
     it.
   - An explicit `paths: []` grants operations and bounds no path.
   - The outside-evidence row never stalls authoring, and the QA gate holds
     Pull Request preparation on it.
   - Authored Verification runs only on committed provenance. The clause says
     that a Run has it by construction, that no other command checks it, and
     that no approval record makes an untrusted source executable.
4. **Lifecycle wording is consistent and adopter-neutral.**
   - "Only `accepted` is active" applies to ADRs that carry lifecycle
     frontmatter, so it no longer contradicts the legacy rule beside it.
   - `deferred` is a terminal Backlog status with a required reason, in the
     documented contract and in the checker.
   - No shipped guidance cites a Spec number or an ADR number, and a check
     refuses one.
   - The repository's own guide names `docs/history/` as the tree that is not
     validated.

## Non-Goals / Out of Scope

- Implementing the provenance check that the verification clause describes as
  missing. It is recorded as the open Backlog Entry
  `docs/backlog/2026-09-30-authored-verification-runs-without-a-provenance-check.md`.
- Changing any clause's enforcement level, or removing a clause other than the
  duplicate backend entry.
- Generalizing clauses that name Roundfix commands, its branch namespace or
  its review providers. Those are deliberate: Roundfix delivers the Baseline
  and runs the loop it describes. The TechSpec lists each one.
- The conflict between the Bun module and the core dependency clause in a
  polyglot repository, and the wording of the TypeScript-stack clauses in the
  `bun`, `typescript`, `backend` and `frontend` modules. A separate
  stack-clause Spec owns them. In the `backend` module this Spec removes the
  duplicate entry and edits nothing else.
- The membership of upstream skills in a setup: the Go skills a module
  dispatches and no setup lists, a search skill that setups list and no module
  requires, and skills renamed upstream. A separate Spec refreshes the
  upstream skill snapshot and owns them.
- The text of any skill. Spec 0192 owns the owned skills' statements, including
  the loop order that two skills repeat.
- The skill version floor, outdated-skill reporting, the Roundfix skill's
  setup membership and the release step. Spec 0195 owns them.
- Rewriting the Source Baseline corpus or the parity corpus, which record past
  releases and stay byte-identical.

## Success Metrics

1. Success Metric: the rendered backend guide holds the boundary paragraph
   exactly once, and a catalog with two clauses of equal text fails the
   duplicate check.
2. Success Metric: the loop clause declares review before archive, and a
   declared order that differs from the Delivery Queue's action order fails
   the order check.
3. Success Metric: no rendered guide contains `execution_approvals`, "does not
   by itself refuse implementation", "may mutate only its bounded" or "never
   blocks the Spec", and each replaced clause carries its new statement.
4. Success Metric: no module guidance matches a Spec number or an ADR number,
   and a module that gains one fails the citation check.
5. Success Metric: every clause present before this Spec, except the duplicate
   backend entry, is present afterwards with the same enforcement level.
6. Success Metric: after the sanctioned regeneration and this repository's
   managed refresh, a second refresh reports no file change.

## Declared breaks

Adopters receive these changes on their next `roundfix baseline update`. Each
is a correction of guidance, and none changes a command.

- **Break — loop order.** An Agent that followed the old clause archived before
  the review. It now reviews first. This is the order the product already ran.
- **Break — tooling Task scope.** The clause no longer tells an Agent to fail a
  Task for an ordinary undeclared path.
- **Break — outside evidence.** The clause now says the gate holds Pull Request
  preparation on the row. The old text said it never blocks.
- **Break — approval record.** The `execution_approvals` description is
  removed. No repository could have used it, because nothing read it.
- **Clarification — everything else.** The absent-record sentence, the
  `paths: []` sentence, the ADR lifecycle scope, the `deferred` status, the
  removed citations and the de-duplicated backend paragraph add or restate
  what already held.

## Recorded limits

- The checks prove order, uniqueness, force and the absence of repository-only
  citations. They cannot prove that a sentence is true. Spec 0195 adds the
  release step that re-reads the guides against shipped behavior.
- Twenty `deferred` entries under `docs/history/backlog/` have no reason. They
  predate the rule, history is not validated, and they stay byte-identical.
- This repository's guides lose their ADR links in the rewritten clauses. The
  decisions remain in `docs/adr/`.

## Decisions

- **Say what is enforced.** A clause that asks for a control the product lacks
  names who must apply it. See ADR-0186.
- **Keep Roundfix's own names.** Commands, the Run branch namespace and the
  configurable review providers stay in the clauses that use them.
- **`deferred` joins the Backlog vocabulary** instead of rewriting twenty
  history entries to `declined`.
- **Guard in code, not in prose.** Each of the four failure classes found gets
  a test that fails when it returns.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- Adopter repositories, read through their Secondbrain mirrors on 2026-09-30.
  `projects/conexus/mirror/docs/agents/backend.md` holds the boundary
  paragraph twice. `projects/oraculum/mirror/docs/agents/autonomous-work.md`
  and `projects/vortex/mirror/docs/agents/autonomous-work.md` carry the
  archive-before-review order, and all five adopter mirrors carry `Spec 0078`.
- The Delivery Queue's recorded behavior in this session: Spec 0183 (#282)
  was reviewed and then archived by the queue on 2026-09-29.
- Published guidance on documentation drift: Google's documentation best
  practices, "Update Docs with Code"
  (<https://github.com/google/styleguide/blob/gh-pages/docguide/best_practices.md>).

## Research basis

The Secondbrain was consulted through `wiki/index.md`, the adopter mirrors
under `projects/*/mirror/docs/agents/`, the triaged entries
`inbox/roundfix/_triaged/2026-08-08-baseline-update-medido-contra-a-frota-adotada.md`
and
`inbox/roundfix/_triaged/2026-08-12-archive-e-o-guia-do-baseline-discordam-sobre-onde-a-spec-arquiva.md`,
and the query
`qmd query "baseline update adotante guias desatualizados skills versão mínima"`.
The mirrors confirmed that adopters hold the stale clauses today. The triaged
entries show the same class, a guide disagreeing with the product, was fixed
one clause at a time before. Exa located Google's and Chromium's documentation
best practices, which place a documentation change in the same change as the
behavior. They support checking guidance mechanically where a check exists.

## Technical candidate

The [_techspec.md](_techspec.md) records the clause texts, the checks, the
regeneration outputs and the build order.
