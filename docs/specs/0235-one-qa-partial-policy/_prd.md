---
spec: 0235-one-qa-partial-policy
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# One QA partial policy

A QA gate runs before any Pull Request exists, and often inside a Run sandbox
without network access. Its report then closes `partial` on rows nobody in the
Run can reach: the Pull Request row and outside-evidence rows whose source the
sandbox denied. The Daemon settles such a QA Task `failed`, `roundfix settle`
refuses it, and the Delivery Queue parks the item `qa-environment-partial`.
The operator then archives with `--qa-override` and retries. That happened on
nearly every delivery from 2026-10-04 to 2026-10-06
(`~/.roundfix-operator/queue-interventions.md`, entries 140 to 186), and two
adopters reported the same refusal
([references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md)).

This is a bug fix. Its causes were measured during authoring with the shipped
binary (`roundfix qa-report accept`) on this repository's archived reports and
on the Fluxus report mirrored in the Secondbrain:

1. The four callers already share one function, `spec.QAReportEligibility`
   in `internal/spec/qa.go`: Daemon settlement (`settleQAVerdict` in
   `internal/daemon/task_engine.go`), `requireSettleQAReportEligibility` in
   `internal/cli/settle.go`, `archiveUnprovenActions` in
   `internal/spec/archive.go`, and `runQAReportCommand` in
   `internal/cli/qa_report.go`. The function itself refuses what ADR-0167
   allows. A `partial` whose only unmet rows are pre-PR Pull Request rows hits
   `if report.RowsBlockedDeclared == 0` and is refused with
   `newest QA Report verdict is "partial"; expected "pass"`. Archived Specs
   0213 and 0220 and the Oraculum report have this shape.
2. `qaReportProvenanceNames` matches the provenance item `Pull Request row`
   exactly, so `Pull Request row (Requirement 10 constrains execution)`
   (archived Spec 0227) is an ordinary environment row:
   `rows_blocked_environment is 2; expected 0`.
3. The Fluxus refusal (`rows_blocked_environment is 1; expected 0` from
   `roundfix settle`, while `qa-report accept` of
   `qa-report-2026-10-05-01.md` exits 0) cannot come from the archived
   report. That report is accepted by today's binary. So settle judged
   another report, the newest one in its own surface, and its refusal did
   not say which one.
4. Outside-evidence rows the sandbox could not reach count as ordinary
   environment rows. Archived Specs 0221, 0222, 0225, 0229, 0230, 0232 and
   0233 were refused for them, beside the Pull Request row. Their causes
   were written in free text (`GitHub GraphQL Forbidden`,
   `OpenCode documentation domain denied by sandbox`).
5. The Delivery Queue's `qaEnvironmentPartial` in
   `internal/cli/deliver_workflow.go` applies a fifth predicate,
   `partial && rows_blocked_finding == 0 && rows_blocked_environment > 0`. It
   parks reports that the shared policy would accept.

## Prerequisites

None. Specs 0236 and 0237 are authored in parallel. This Spec is delivered
first. Its task_04 regenerates derived Baseline files that Spec 0236 may also
touch, and 0236 then regenerates them on top of this one.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  new names are a row status marker and a Go field and method. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network call or
  HTTP surface is added; the policy reads QA Reports on disk and in Git, and
  no test reaches the network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0240 (this Spec) decides the
  policy: "The policy has one rule, in `internal/spec`, which every caller
  reads". It extends ADR-0167, whose pre-PR row "no longer counts against a
  qualifying `partial`", supersedes ADR-0104 in part, whose gate "holds pull
  request preparation until that row is satisfied", and supersedes ADR-0229
  in part, whose queue parks `qa-environment-partial` "including when every
  environment-blocked row is one that waits for an open Pull Request".
  ADR-0080 still decides that an environment-blocked row "does not cap the
  verdict when the report records equivalent observed or supervised
  evidence", and ADR-0154 stands for what stays ineligible: a QA archive override records user authority, not a pass,
  ADR-0154: "The exception can waive the terminal QA Task's archive".
  ADR-0088, ADR-0091, ADR-0096, ADR-0097, ADR-0155, ADR-0182, ADR-0194,
  ADR-0195 and ADR-0210 bind the gate, its machine stage, its row carry, its
  matrix and its evidence. ADR-0179 bounds the Governed Paths, ADR-0187,
  ADR-0189 and ADR-0233 bind the owned skills' versions, and ADR-0184 binds
  the Surface Transcripts. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176
  and ADR-0183 check this Spec's consistency by citation and receipt. ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue, which this Spec leaves unchanged, so it does not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the qa-gate, archive-spec and Roundfix
  skills, two Baseline modules with the derived Baseline files they
  regenerate, and two generated guides are Governed Paths. The maintainer
  granted them on 2026-10-06 ("Concedo"), beside the standing grants of
  2026-09-30. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0235-one-qa-partial-policy/_authorization.md`; bounded files:
  `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`,
  `docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/archive-spec/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- The four QA settlement callers accept the same reports for the same reason
  and refuse the same reports with the same reason.
- A `partial` qualifies when its only unmet rows are pre-PR Pull Request rows,
  network-denied outside-evidence rows, and declared rows its declarations
  cover. The Delivery Queue then proceeds without parking the item.
- Any other environment-blocked row, a finding-blocked row or a skipped row
  still needs a declaration, a fix or the QA Archive Override.
- A settle refusal names the report it judged.
- Every archived report keeps reading as before.

## Core Features

1. **One policy.** `QAReportEligibility` accepts a `partial` whose unmet rows
   are only exempt environment rows and covered declared rows, refuses one
   with a skipped row, and keeps every other refusal (ADR-0240).
2. **A Pull Request row with a note.** A provenance item that starts with
   `Pull Request row` followed by a space, `:` or `(` names the row.
3. **A network-denied outside-evidence row.** Status
   `blocked (environment: network denied: <host>)` and an `outside-evidence row`
   provenance item. The qa-gate skill tells the gate to write exactly this.
4. **The queue parks only what needs the override.** `qa-environment-partial`
   requires an environment row outside the two exempt kinds.
5. **Settle names the report.** Its refusal ends with the report path.
6. **The guides say it.** The qa-gate, archive-spec and Roundfix skills, the
   Spec routing and autonomous-work guides, and the user guides state the one
   policy.

## Non-Goals / Out of Scope

- Changing the `pass` rule, a `fail`, the mechanical stage's count checks, or
  the frontmatter count fields.
- Exempting any other environment row, such as a missing Run Database item or
  an operator export, or a network-denied row that is not an outside-evidence
  row.
- Having the Daemon observe network denials itself.
- Stamping the not-reached rows into the archived `_prd.md`.
- Rewriting archived reports or Specs, `CONTEXT.md`, `CHANGELOG.md`, the
  Makefile, `go.mod` or CI.

## Success Metrics

1. Success Metric: archived reports of Specs 0213 and 0220, whose only unmet
   rows are pre-PR Pull Request rows, are accepted by `roundfix qa-report
   accept`. The shipped binary refuses them with "expected pass".
2. Success Metric: a report whose unmet rows are a pre-PR Pull Request row and
   a `blocked (environment: network denied: docs.example.com)` row with an
   `outside-evidence row` provenance is accepted by `qa-report accept`,
   `roundfix archive`, `roundfix settle` and Daemon settlement alike, and
   one with any other environment row, a finding row, a skipped row or an
   uncovered declared row is refused by all four, with the same reason.
3. Success Metric: the Delivery Queue does not classify a qualifying partial
   as `qa-environment-partial`, and still does for a report with another
   environment row.
4. Success Metric: a refused `roundfix settle` names the report path it
   judged.
5. Success Metric: every existing QA, archive, settle, delivery and Baseline
   test passes, with the declared changes only.

## Acceptance evidence

The outside-evidence rows rest on records this Spec did not produce:

- This repository's archived QA Reports, measured on 2026-10-06 with the
  shipped binary (v0.46.0, `bin/roundfix qa-report accept` on the newest
  report of each archived Spec): of the archived `partial` reports, 0213 and
  0220 are refused with `newest QA Report verdict is "partial"; expected
  "pass"` although every unmet row is a Pull Request row, 0227 is refused
  with `rows_blocked_environment is 2; expected 0` because its Pull Request
  row's provenance carries a note, and 0073, 0079 and 0179, whose only unmet
  rows are declared rows their Specs cover, are accepted. Every other
  archived `partial` report is refused.
- The Fluxus report mirrored in the Secondbrain,
  `projects/fluxus/mirror/docs/history/specs/0100-cada-loja-aponta-para-o-mesmo-fabricante-e-a-mesma-categoria/qa/`.
  Its `qa-report-2026-10-05-01.md` (two declared rows, the Pull Request row)
  is accepted by the shipped binary, and its sibling
  `qa-report-2026-10-05.md` is refused with "2, 1 outside the pre-PR Pull
  Request row". The archived `_prd.md` records that the operator overrode it
  because "o settle 0.44 recusa só pela linha do PR".
- The Oraculum report of 2026-10-06
  (`inbox/roundfix/_triaged/2026-10-06-partial-so-com-a-linha-do-pr-e-recusado-no-assentamento.md`
  in the Secondbrain): a `partial` with `rows_blocked_environment: 1`, no
  declaration, and the Pull Request row as its only blocked row. It was
  refused with "expected pass".
- JUnit's published assumption contract
  (<https://docs.junit.org/5.14.4/api/org.junit.jupiter.api/org/junit/jupiter/api/Assumptions.html>),
  read 2026-10-06: "failed assumptions do not result in a test failure;
  rather, a failed assumption results in a test being aborted". This is the
  same split between an unmet precondition of the environment and a failure
  of the work under test.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "QA partial pull request row settlement archive override"`
(`--all --files --min-score 0.3`). It returned this repository's mirrors of
ADR-0154, ADR-0167 and Spec 0181, Fiscus's autonomous-work guide (same
clause) and the Fluxus Spec 0100 report above. A second query found the
mirrored Fluxus report itself. Exa found JUnit's assumption contract. The
only open Backlog Entry that shares this context is the adopted one. The
other two open entries of 2026-10-06 belong to Specs 0236 and 0237. There is
no unresolved Finding.

## Decisions

- One rule in `internal/spec` for every caller; the verdict rule of the gate
  is unchanged. See ADR-0240.
- The network-denied marker is a status form plus the outside-evidence
  provenance, both derived from the Results rows.
- A skipped row keeps a partial from qualifying.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
