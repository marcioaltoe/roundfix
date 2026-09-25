---
spec: 0172-a-qa-gate-that-tells-the-truth
status: active
created: 2026-09-25
surfaces: [backend, cli, docs]
---

# A QA gate that tells the truth

A QA gate's verdict and the Verification signal on the Run Event Stream must
mean exactly what they say. Five defects make them say more, or less, than was
measured:

- The Daemon seeds every QA Report as `verdict: pass` with an empty Results
  table. When the QA Agent stops without filling it, or writes `pass` without a
  single row, settlement, `archive`, `settle` and `qa-report accept` all accept
  it. A Fiscus Spec settled its gate on a 669-byte report written in 30 seconds.
- The QA Report allocator takes the first free `-NN` suffix of the day while
  the reader picks the highest, so after a gap the fresh report loses to an
  older one and settlement reads a stale verdict.
- In independent Verification, a command that fails deterministically on the
  first run and temporarily on the exclusive retry loses its deterministic
  failure: the Task's reason names only the temporary one, and the operator
  re-runs instead of fixing an observed defect.
- An unobserved Verification is published without its classification, so the
  Supervisor stream reads "we did not find out" as "the work is wrong".
- The vacuous pre-work Verification event carries `probed_commands` while the
  projection requires `commands`, so `roundfix events` aborts for the whole Run.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; QA Report names keep
  their date and numeric-suffix form. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0015 reads the verdict from the QA
  Report and counts a missing or unreadable one as a failure; ADR-0080 owns
  verdict semantics and the typed blocked-row counts; ADR-0091 keeps the gate a
  Task node; ADR-0096 says no stage may make a verdict more permissive;
  ADR-0132 calls a Results table with no rows malformed; ADR-0014 and ADR-0057
  keep Verification and Task status Daemon-owned; ADR-0038 allows one
  Verification repair, ADR-0056 makes a temporary failure on the exclusive
  retry terminal and ADR-0159 collects every deterministic failure of an
  independent Verification; ADR-0111 separates an unobserved Verification from
  a verdict and ADR-0135 reports an absent diagnostic as a state; ADR-0160
  keeps the red repository gate's entry and settlement rules, which the
  unknown classification of an unobserved precondition leaves unchanged;
  ADR-0148 owns the vacuous pre-work probe. ADR-0020 and ADR-0127 cite ADR-0014
  but govern the Agent prompt result and process residue, which this Spec does
  not touch, so they do not apply. ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps
  a path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved the 2026-09-25
  efficiency sequence in chat on 2026-09-25; the Go test file rides the
  standing grant of 2026-09-21 for governed source and the skill files ride the
  standing grant of 2026-09-18 for keeping the shipped skills true to the CLI,
  recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `internal/spec/archive_test.go`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `docs/references/coverage-record.json`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A QA Report nobody filled never passes.
- The newest QA Report is the one settlement reads.
- A deterministic Verification failure is never hidden behind a temporary one.
- An unobserved Verification reaches the Run Event Stream as unobserved.
- `roundfix events` survives a vacuous pre-work probe and any other record it
  cannot project.

## Core Features

1. **A report starts pending and a hollow pass is refused.** The Daemon seeds a
   non-blocking QA Report with `verdict: pending`, a readable verdict that is
   never accepted. A `pass` or `partial` whose `## Results` section carries no
   row, and which carries no row in any other table with a `Status` column, is
   refused by every acceptance decision: QA settlement, `archive`, `settle` and
   `qa-report accept`. A refused `pass` or `partial` settles the QA Task with a
   reason that names why.
2. **The allocator writes where the reader looks.** A new QA Report takes the
   suffix one above the highest existing suffix of its date, and a report dated
   after today stops the allocation instead of hiding the new one.
3. **A temporary retry keeps the deterministic failure.** A command that failed
   deterministically on the first run and temporarily on the exclusive retry
   keeps its first-run failure and diagnostic in the Task's reason; the Task
   still settles failed without a repair turn, as ADR-0056 requires.
4. **An unobserved Verification is classified.** The shared attempt publisher
   emits `classification: verification_unknown` with `command`, `reason` and
   `diagnostic_path` on both events of an unobserved outcome; a command verdict
   keeps its current classification.
5. **The event stream survives.** The vacuous pre-work event carries `commands`,
   the projection also reads journals written before the fix, and `roundfix
   events` skips a record it cannot project with a warning on stderr instead of
   aborting.

## Non-Goals / Out of Scope

- A fourth Task status or a new Run outcome for an unobserved Verification.
- Changing when a readable `pass` with rows, a declared-only `partial` or a
  `fail` is accepted or refused.
- Changing the derived QA Verification command the `qa` Task renders; the
  hollow-report check lives in the in-process acceptance decision.
- Rewriting archived QA Reports or any archived Spec.

## Success Metrics

1. A seeded QA Report left untouched by the QA Agent settles the QA Task
   `failed` with verdict `pending`, and a `pass` with no QA row is refused by
   settlement, `archive`, `settle` and `qa-report accept`, while every archived
   `pass` Spec in `docs/history/specs/` stays archive-eligible.
2. With reports `qa-report-D.md` and `qa-report-D-02.md` present, the next
   report of date D is `qa-report-D-03.md` and `NewestQAReport` returns it.
3. First run A deterministic plus B temporary, retry A temporary: the Task's
   reason names A's first-run failure and first-run diagnostic path.
4. Every event of an unobserved Verification projects with classification
   `verification_unknown`, its command, reason and diagnostic path; a
   deterministic command failure projects with no classification.
5. `roundfix events` on a Run whose journal holds a vacuous pre-work event, or
   a malformed record, exits `0` and emits every projectable record.

## Recorded limits

- A QA Report with no `## Results` section keeps its current eligibility: it is
  the shape written before the Daemon seeded reports, and the archived corpus
  depends on it. The mechanical stage's shape detector already reports a
  missing matrix on the gate path.

## Decisions

- **Refuse in the acceptance decision, not in the reader.** `ReadQAReport`
  keeps reading a hollow report, so status displays and the event record still
  show what the report says; `QAReportEligibility` is the one decision every
  accepting caller shares, so one check closes all four doors.
- **A row is a row of the matrix.** The Results section, including tables under
  its sub-headings, plus any table with a `Status` column: the archived Spec
  0094 recorded its matrix under a second heading beside an empty mechanical
  Results table, and it stays eligible.
- **Keep ADR-0056 terminal.** A temporary failure on the exclusive retry stays
  terminal; only the reason changes, so the fix restores a fact without adding a
  repair turn.
- **Warn, do not abort.** A Supervisor following a Run loses every later record
  when one record cannot be projected; a warning on stderr keeps stdout JSONL
  and names the record that was skipped.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an untouched seed that passes,
an allocation below the highest suffix, a temporary reason where a
deterministic one was observed, an unclassified unknown event, or an aborted
stream would each pass a happy-path test.

The outside-evidence row rests on QA Reports this Spec did not write: every
archived `pass` Spec in `docs/history/specs/` stays archive-eligible, and every
newest `pass` report in the Fiscus mirror at
`/Users/marcio/dev/secondbrain/projects/fiscus/mirror/docs/history/specs/`
is still accepted by the built `roundfix qa-report accept`.

## Research basis

The Fiscus hollow reports were captured in the secondbrain inbox
(`inbox/roundfix/_triaged/2026-09-18-um-relatorio-de-qa-sem-matriz-passa-na-verificacao-canonica.md`):
two reports of 669 and 581 bytes, each with a `## Results` header and no row,
settled Spec 0046's gate. The vacuous-event mismatch was captured in
`inbox/roundfix/_triaged/2026-09-17-evento-de-verificacao-vacua-usa-chave-que-a-projecao-nao-le.md`.
The temporary-retry reproduction is a recorded limit of Spec 0167's second
pre-PR review. The adopted sources are indexed in
[references/_index.md](references/_index.md).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
