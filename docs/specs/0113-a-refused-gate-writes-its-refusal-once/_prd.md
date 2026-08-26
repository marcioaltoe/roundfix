---
spec: 0113-a-refused-gate-writes-its-refusal-once
status: active
created: 2026-08-25
surfaces: [backend, docs]
---

# A refused gate writes its refusal once

When the QA gate refuses at a precondition, it writes a report whose Results
table is empty — correctly, because it stopped before building the matrix.
The mechanical stage of every later run then reads that report and refuses:

```text
QA-REPORT-SHAPE
location: qa/qa-report-2026-08-14.md:1
detail:   Results table has no report rows
fix:      Materialize every planned QA row with one terminal status.
```

So a gate that refused for a good reason leaves behind an artifact that blocks
every future gate run on the same Spec, and the prescribed fix — materialize
every planned row — is impossible for a run that never built a matrix. This
creates a loop that is unbreakable without manual intervention (artifact
deletion) and manual edge-case handling.

Measured on Spec 0103 on 2026-08-14 and on Spec 0094 on 2026-08-12/13. The
gate writes a structurally invalid report while doing the right thing — stopping
at a strict Spec check failure. The only exit was deleting the empty report,
which is evidence removal.

A third measurement, on 2026-08-26, shows the same deadlock reached through a
report that is not empty. Spec 0098's gate wrote `qa-report-2026-08-25.md` with
verdict `fail`, correctly recording two `Trust-Damage` findings as `fail` rows.
Both findings were then fixed by that Spec's two corrective Tasks. The next gate
run (`run_20260826T004155Z_8a013e2b7e48654e`) refused anyway: the mechanical
stage read the superseded report and raised `QA-REPORT-SHAPE` against its `fail`
rows, blocking the very run that carried the fix. The stage reads
`previousReportPath` (`internal/daemon/task_engine.go:2073`), so newer evidence
never supersedes older evidence.

Two properties of that measurement matter for the design below. First, the
blocking rows were correct when written — this is not a malformed artifact but a
valid one outliving its run, so deleting it would again be evidence removal.
Second, the same run's findings included rows named `Case`, `82-line function vs
80`, and `` `sort()` vs `toSorted()` ``, which are cells of an evidence table in
the report's prose, not rows of its Results table. The shape detector parses any
markdown table in the report as the Results matrix, so prose evidence a gate
writes to justify a row becomes a row itself, and each one is a fresh blocker.

## Project Constraints

- Identifier strategy: applicable — QA Report, precondition, and terminal row
  are glossary terms this Spec clarifies in code. The closing node checks
  whether the work introduced or changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is gate logic and report writing. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0096 asks for machine facts proof, ADR-0117 places checks at defect origin
  before spending an Agent turn, which is exactly the behavior this Spec's
  precondition refusal implements. This Spec changes what the gate writes when
  that precondition fires. No accepted ADR governs the precondition-refusal
  report shape, which is why the rule this Spec adds is new. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work changes QA gate behavior and its report writing in
  production Go and its tests, creating or editing no linter, formatter,
  test-runner, build, or skill configuration. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. A gate that refuses at a precondition writes a report whose shape its own
   contract accepts.
2. A subsequent run does not inherit a refusal from a previous run's
   precondition failure.
3. The contract between gate refusal and mechanical stage is written and
   verified.

## User Stories

1. As a Supervisor running a Spec whose gate encounters a precondition
   failure, I want the gate to write a valid report, so that a subsequent run
   does not inherit an error about the previous run's state.
2. As a Spec author defining strict preconditions, I want the gate's refusal
   to be recorded correctly, so that the reason for the refusal is auditable
   and does not block future runs.
3. As a maintainer reading a gate refusal, I want to know what check caused it,
   so that I can decide whether to fix the Spec or fix the tree.

## Core Features

1. **A precondition-refused gate writes one terminal row.** When the gate stops
   at a precondition check, it writes a QA Report with exactly one terminal row
   that records the precondition that caused the refusal, rather than an empty
   Results table. The row is `status: blocked` and `provenance: precondition`.
2. **The mechanical stage reads the newest report.** The mechanical stage that
   validates reports before subsequent gate runs reads only the newest report
   in the QA directory, not every report, so a superseded refusal cannot block
   the run that supersedes it.
3. **Gate refusal precondition is recorded in the report.** The terminal row
   that records the precondition refusal includes the check name and the reason
   it failed, so a reader understands what was checked and why.
4. **Only the Results table is read as results.** The shape detector parses rows
   from the report's Results table alone. A markdown table a gate writes in its
   prose to justify a row — an evidence matrix, a per-case comparison — is not
   parsed as a result and cannot block a later run. Measured on Spec 0098 on
   2026-08-26, where three prose evidence cells each became a separate blocker.

## Non-Goals / Out of Scope

- Changing what checks the gate runs, or the order.
- Changing the strict precondition itself; this Spec accepts it and ensures the
  gate can survive its own refusal.
- Changing Review Source behavior or any other external contract.

## Success Metrics

- Spec 0103, which failed at a strict precondition check, can retry without
  deadlock. Its gate passes when conditions improve.
- A gate refusal report is structurally valid per the QA Report contract,
  proven by the mechanical stage accepting it without SC-REPORT-SHAPE errors.
- The newest-report-only read is proven by measuring that a previous run's
  refusal does not block a subsequent run when fresh conditions are met.

## Decisions

- A precondition refusal records the failure as a terminal `blocked` row rather
  than failing without a report, because the refusal itself is the observed
  state worth recording for audit.
- The mechanical stage reads newest-only, not filtered-by-status, because the
  gate is the authority on its own result and newer evidence supersedes older
  evidence.

## Open Questions

- Whether every precondition refusal should be recorded in the report, or only
  certain ones (e.g., strict Spec check failures but not others). The default
  is to record every precondition refusal that causes the gate to stop.
