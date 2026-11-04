---
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
status: active
created: 2026-11-04
surfaces: [backend, docs]
type: fix
---

# A retry after a QA precondition refusal runs the whole QA gate

When a QA gate refuses at a precondition, the Daemon commits a Precondition
Refusal report on the Run Branch, and the next Run of the same Spec imports
that report as its previous QA report, because the import that hands a failed
QA pass to the next pass does not tell a refusal from a pass that built its
matrix. On 2026-11-03 Pantheon's retry imported a refusal whose tables its
formatter check rejected, so the next gate refused again at the same
precondition, without an Agent session, and committed the old refusal beside
a new one. After this Spec a later Run never imports a Precondition Refusal,
so a retry after one runs the whole gate again, precondition first and then
the Agent session; the refusal report keeps going through the configured
Format Command before its commit, as every QA Report does, and the Roundfix
skill's implement reference says what the import brings in, what it never
brings in, and how to start a pass with no previous report.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md`](references/2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md),
   reported by pantheon, ai-memory message `01a1a32d` of 2026-11-03, copied
   verbatim: the two Runs of Pantheon's Spec 0059 on roundfix 0.129.1, the
   commits `73d704a` and `a2c3e68`, and the manual workaround.

No open Backlog Entry or unresolved Finding on `main` shares this context
(the only open entry, about the waiting declaration in the Baseline
contracts, belongs to Spec 0346).

## Prerequisites

None in `requires`. Authored against `6a0ed904` (v0.130.1), where no Spec is
active. Specs 0346, 0347 and 0348 are authored in parallel and deliver after
this Spec, in the order 0345, 0348, 0347, 0346. `_techspec.md` → Build Order
names the files they may share.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or file name is
  minted; the import's existing `refused` outcome gains one reason text.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the import and the Format Command read and write only
  the Run Worktree and the Run Database. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0194 brings the newest unintegrated
  QA Report commit of the Spec into the next pass, ADR-0194: "a failed pass on
  an unintegrated Run Branch is the next pass's previous report"; this Spec
  keeps that for a pass that built its matrix and excludes a Precondition
  Refusal. ADR-0210, which refines ADR-0194, records one digest per declared
  input in the Evidence Snapshot, ADR-0210: "This refines ADR-0194 without
  superseding it"; a refusal records no snapshot, and nothing here changes the
  digest. ADR-0249 formats the QA
  directory at the import and commit boundaries, ADR-0249: "An empty value
  runs nothing, and that is today's behavior."; this Spec keeps that rule for
  the refusal report and rejects again the alternative ADR-0249 rejected, not
  importing a failed pass, except for a refusal, which is not a pass.
  ADR-0132 makes a refused gate write a report, ADR-0132: "it emits one
  terminal row naming the refusal that produced it"; the refusing Run still
  writes and commits that report, and only a later Run stops importing it. ADR-0355, written by
  this Spec, records the exclusion. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AX answers of 2026-11-04 quoted in the authorization record, which keep the
  cycle AW rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the Roundfix skill's `SKILL.md` pair and
  its canonical `implement` reference. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint, formatter or Baseline module change.

## Goals

- A Run that follows a Precondition Refusal of the same Spec runs its QA gate
  in full: the repository Verification precondition, then the Agent session
  when the precondition passes.
- No later Run imports a Precondition Refusal report, settles its QA Task
  from one, or commits one again.
- A failed or passing QA pass that built its matrix is still imported byte
  for byte, as today.
- The Precondition Refusal report the Daemon writes goes through the
  configured Format Command before its commit, under the same restore rules.
- The Roundfix skill's implement reference and the configuration guide tell
  an adopter what the import brings in, what it never brings in, how to start
  a pass with no previous report, and that a repository whose Verification
  checks the formatting of QA files configures the Format Command.

## Core Features

1. Refusal never imported. Before a QA pass builds its matrix, the import
   refuses a candidate QA Report commit when any QA Report the commit adds or
   modifies records a Precondition Refusal (`rows_blocked_precondition`
   greater than zero). The import's event says `refused` with the reason
   `precondition refusal: <report path>`, copies nothing into the Run
   Worktree, and the pass runs with no previous report.
2. The Pantheon shape is impossible. A Run whose earlier Run Branch holds a
   refusal commit runs the precondition and, when it passes, the Agent
   session; its QA Report commit holds no file of the refusal commit.
3. Pass import kept. A candidate whose reports record no Precondition Refusal
   is imported exactly as today, failed passes included.
4. The refusal report is formatted. The Format Command runs over the
   Precondition Refusal report immediately before its QA Report commit, and a
   failing formatter restores its bytes and records the outcome, as for every
   QA Report. This is the current behavior, which the Spec pins with a test.
5. The rules in writing. The implement reference (canonical and mirror), the
   configuration guide, `CONTEXT.md` and ADR-0355 state them.

## Non-Goals / Out of Scope

- Not importing a failed QA pass that built its matrix: ADR-0194 keeps it,
  ADR-0194: "a failed pass on an unintegrated Run Branch is the next pass's
  previous report".
- Carrying the refusal commit forward or superseding it: Task Carry-Forward
  carries only Task settlement commits, so a QA Report commit never leaves its
  Run Branch; this Spec changes no carry-forward rule.
- A new command or flag to discard a prior report, or deleting a Run Branch
  automatically.
- Rendering the Daemon's report tables in a particular formatter's style, or
  choosing a formatter for a repository that configures none.
- Changing how a Precondition Refusal is written, named or settled in the Run
  that refuses.
- `roundfix doctor` or `spec check` findings about the Format Command.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, any Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: a Run whose earlier Run Branch holds a recorded QA Report
   commit adding a Precondition Refusal publishes the import outcome
   `refused` with reason `precondition refusal: <report path>`, leaves no file
   of that commit in the Run Worktree, runs the repository Verification and
   the Agent session, settles its QA Task from its own report and commits no
   file of the refusal commit.
2. Success Metric: a candidate commit whose newest report built its matrix
   but which also adds an older report recording a Precondition Refusal is
   refused the same way.
3. Success Metric: every existing import, carry and format test of
   `internal/daemon` passes unchanged, among them the byte-for-byte import of
   a failed pass and the row carried from an unintegrated failed pass.
4. Success Metric: with a Format Command configured, the committed
   Precondition Refusal report of a refused repository Verification holds the
   formatter's bytes and a `commit` format event `formatted`; with a failing
   formatter it holds the Daemon's bytes and the event `failed`.
5. Success Metric: the implement reference (both copies), the configuration
   guide and `CONTEXT.md` state that a Precondition Refusal is never imported,
   that the Format Command covers the refusal report, what else the import
   never brings in, and how to start a pass with no previous report; each
   mirror equals its canonical copy.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Pantheon's report, quoted in the adopted reference and written before this
  Spec existed: on roundfix 0.129.1 Run `run_20261103T204416Z_73e61e8cd94f990d`
  recommitted (`a2c3e68`) the refusal of Run
  `run_20261103T200657Z_a437bac8bed64b17` (`73d704a`) without an Agent
  session, and its precondition failed on `oxfmt --check` over
  `qa/qa-report-2026-11-03.md`.
- The Bazel user manual, "Commands and Options", option
  `--cache_test_results` (https://bazel.build/docs/user-manual, read through
  Exa on 2026-11-04): "If this option is set to 'auto' (the default) then
  Bazel will only rerun a test if any of the following conditions applies:
  Bazel detects changes in the test or its dependencies; the test is marked
  as external; multiple test runs were requested with --runs_per_test; the
  test failed." A build tool that reuses earlier results never reuses a
  failure as the next run's result; the next run executes.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Reading at `6a0ed904`: `findRecordedPriorQAPass`
  (`internal/daemon/qa_prior_pass.go`) selects the newest Daemon QA Report
  commit of the Spec that a recorded Run Branch holds and the Run's history
  lacks, whatever the verdict of its report; `copyPriorQAPass` refuses only a
  differing path, a report not newer than the tree's newest, one dated after
  today and one the shape detector refuses. `runQAGate`
  (`internal/daemon/task_engine.go`) runs the Format Command over the files
  of every QA Report commit, the refusal report included, through
  `formatQADirectory`.
- Characterization at `6a0ed904` with a temporary test, removed afterwards,
  in a temporary repository: a recorded Run Branch commit holding a refusal
  written by `spec.WritePreconditionRefusalReport`, and a repository
  Verification that fails when a file under the Spec's `qa/` directory holds
  the delimiter row `| - | --- |`, as Pantheon's `oxfmt --check` did. The next
  `TaskCycle` imported the refusal, refused at the precondition with no Agent
  prompt, settled `fail`, and its QA Report commit held the imported refusal
  and a new one. A second temporary test showed that a refused repository
  Verification with a Format Command configured commits the formatter's bytes
  and publishes `QA format commit formatted.`.
- Pantheon's Project Config, read once with `gh api` on 2026-11-04 to explain
  the unformatted report and not kept: it sets no `verification` key, so no
  Format Command ran on either Run.
- With a prototype of Core Feature 1 in a scratch copy, the same `TaskCycle`
  refused the import, reached the Agent and settled `pass`, and the 62
  existing import, carry and format tests of `internal/daemon` passed.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix QA gate precondition refusal retry reimport report formatter" --all --files --min-score 0.3`,
  which returned only Roundfix and skills project mirrors: the Archive
  Records of Specs 0113, 0131, 0202, 0235 and 0244, ADR-0132, the `qa-gate`
  and `roundfix` skill copies and the `qa-report` guide. They confirm that
  the refusal report must exist (ADR-0132) and that the Format Command was
  the repair chosen for an imported failed pass (Spec 0244, ADR-0249); none
  speaks of importing a refusal. The page
  `wiki/concepts/agendas-e-reexecucao-em-motores-de-workflow-durable.md`
  (re-execution in durable workflow engines) was read and decides nothing
  here.
- Exa: the Bazel passage quoted above, also found on the Bazel Knowledge Hub
  page "Test Caching" (https://bazel.virtuslab.com/book/1~2~1/): "A failing
  test is always re-executed on the next invocation. This prevents a flake
  from getting stuck in cache".

## Glossary

- adds: **QA Report Import** — the Daemon step that brings the newest unintegrated QA Report commit of a Spec into the next pass
- changes: **Precondition Refusal** — a later Run never imports it

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AX brief for Spec 0345 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04), verbatim:

- 0345: "Rodar QA de novo + format (Recommended)" — option text: "Retry após
  recusa de precondição descarta o report de recusa e roda o QA gate inteiro
  (precondição, depois agente); o report de precondição passa pelo
  verification.format antes do commit; implement.md documenta a importação
  de QA Report no retry."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text:
  "Uma minor (a mudança da Baseline e o runs list fora de repo são feat);
  notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e o
  contorno que podem largar."

Taken from those answers, ADR-0194, ADR-0249 and the reading above, with no
further interview:

- A Precondition Refusal is recognized by its `rows_blocked_precondition`
  count, which only the refusal writer sets above zero.
- The import refuses the whole candidate commit when any report it adds or
  modifies records a refusal, and does not search an older candidate: the
  pass then has no previous report and re-runs every row, the safe direction.
- The refusal commit is not carried forward, now or before: Task
  Carry-Forward reads only commits with a `Roundfix-Task` trailer of a Task
  settled `completed`, and a QA Report commit has neither. It stays on its
  Run Branch as history.
- The refusal report already goes through the Format Command at its commit;
  Pantheon's report stayed unformatted because Pantheon configures none. The
  Spec pins the behavior with a test and documents the setting instead of
  changing the format step.
- Starting a pass with no previous report needs no new command: the import
  reads a commit only through its recorded Run Branch, so removing that Run
  Branch, a destructive Git operation the reference names as such, is enough.
- Two implementation Tasks in series (the import and the pinned format path;
  the documents) and a QA gate.

## Open Questions

None.
