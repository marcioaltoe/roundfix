---
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: active
created: 2026-10-07
surfaces: [backend, cli, docs]
---

# A formatted QA report and a reopen for late dependencies

On 2026-10-06 an adopter repository running Roundfix 0.46.0 reported two
defects on its Spec 0047 and Spec 0046
([the adopted Backlog Entry](references/2026-10-06-an-unformatted-qa-report-breaks-the-next-run.md)).

First, a QA Run that failed committed `qa/qa-report-2026-10-06.md` without the
repository's formatter. The next Run imported that unintegrated pass into its
worktree (ADR-0194), and its repository Verification precondition, which runs
the formatter in check mode, exited 2 on the imported file. The Run wrote a
`fail` report with `rows_blocked_precondition: 1` and ran no row. Every passing
QA in that repository also needed a separate formatting commit, because the QA
Report commit is the only commit the Daemon makes after the last repository
Verification.

Second, a corrective Task already `completed` was added as a dependency of a
QA gate that had settled `completed`. `roundfix reopen` refused, because every
dependency was completed, and the operator set `status: pending` by hand.

This is a bug fix: the QA gate, the import and the reopen command keep their
contracts. The QA step now formats what it commits with a command the
repository configures, and reopen also sees a dependency the gate never
covered. This minimal PRD exists for the downstream artifact contract; the
design lives in the [_techspec.md](_techspec.md) and ADR-0249.

## Prerequisites

None. Spec 0245 is authored in parallel. If it also raises the Roundfix
Skill's version, the operator orders the queue so the later Spec raises it from
the earlier one's value; task_01 records the next free version at record time.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  Project Config gains the snake-case key `verification.format`, and the
  existing `daemon.qa` Run Event kind gains the phase `format`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; the Format Command and Git run locally, and every test uses a
  temporary repository, a temporary home and a shell-script formatter. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0249 (this Spec) decides the Format
  Command, its two run points, its revert, and the Late Dependency proof.
  ADR-0194: "It copies that commit's report and evidence byte for byte, so a
  failed pass on an unintegrated Run Branch is the next pass's previous
  report", refined by ADR-0249: the copy and the carry proof stay byte for
  byte, and only then is the import formatted. ADR-0195 decides which rows are
  always observed and is unchanged. ADR-0097 lets a row carry on unmoved
  inputs, unchanged. ADR-0059: "managed Markdown must survive the repository's selected
  formatter unchanged", a different output this Spec leaves alone. The gate is bound by ADR-0080,
  ADR-0091, ADR-0104, ADR-0156, ADR-0167 and ADR-0240, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by citation and
  receipt. ADR-0182 runs Settlement Checks before each Task commit. ADR-0229
  cites ADR-0167 but decides how an operator archive resumes a park, and ADR-0237
  decides how a Delivery Retry records a merge made outside the queue; this
  Spec changes neither. ADR-0184: "A TechSpec now declares numbered Surface
  Transcripts", answered in the TechSpec with the reason none applies. ADR-0096 decides the gate's machine stage, which
  runs before the format and is unchanged, and ADR-0210 decides the digest an
  Evidence Snapshot records, which the format never reads. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths. The maintainer authorized skill edits
  ("considere autorizado a ajustar todas as skills se necessário") and the
  Governed Paths each Spec declares ("Concedo"), and approved this Spec as the
  first of the 2026-10-07 cycle ("Pode seguir nessa ordem"). No other Governed
  Path changes; no Makefile, formatter or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A QA step never commits a file under the Spec's `qa/` directory that the
  repository's configured formatter would rewrite, unless the formatter itself
  fails, and then the failure is recorded.
- An imported prior pass is formatted before the repository Verification
  precondition sees it, so the precondition judges what the Run will commit.
- A formatter that fails, times out or changes the verdict never costs the QA
  verdict: the original bytes are restored and committed.
- `roundfix reopen` returns a completed QA gate to `pending` when its
  dependency closure gained a Task after the newest QA Report was recorded,
  whatever that Task's status.

## Core Features

1. **Format Command.** The Project Config key `verification.format` names the
   command; empty, the default, runs nothing (ADR-0249).
2. **The QA step formats its QA directory.** The Daemon runs the command over
   the imported pass after the mechanical stage and before the precondition,
   and over the QA Report commit's files under `qa/` just before that commit.
   Each run is a `daemon.qa` Run Event with phase `format`.
3. **Reopen sees a Late Dependency.** Reopen reads the Task Graph manifest at
   the commit that added the newest QA Report and reopens the gate over any
   Task the current closure adds.
4. **Docs and glossary.** `CONTEXT.md` gains **Format Command** and **Late
   Dependency**; the configuration guide, the Spec workflow guide, the reopen
   reference and the Roundfix Skill's `implement` and `settle` references
   describe both.

## Non-Goals / Out of Scope

- Formatting Task commits, Task files, other Spec artifacts or the Archive
  Record.
- Discovering the formatter from the Baseline profile, `make fmt` or the
  repository Verification.
- Excluding files from the repository Verification precondition, or stopping
  the import of a failed pass (ADR-0194 stands).
- Detecting a Late Dependency in the Task Graph loader, `implement`, `settle`
  or `deliver`.
- Checking the Format Command's executable in Doctor.

## Success Metrics

1. Success Metric: with `verification.format` set to a shell-script formatter,
   a QA pass commits its report and evidence in the formatter's output, and a
   second pass that imports an unformatted prior report passes a precondition
   that fails on unformatted files under `qa/`.
2. Success Metric: a formatter that exits non-zero, or rewrites the report's
   verdict, leaves the committed bytes equal to the unformatted ones, and the
   Run records the `format` phase with outcome `failed` or `reverted`.
3. Success Metric: with `verification.format` empty, the QA Report commit and
   the Run Event Stream are unchanged.
4. Success Metric: in a temporary repository whose newest QA Report was
   committed before `task_03` existed, `roundfix reopen` with `task_03`
   `completed` exits 0, sets the QA Task `pending`, and records
   `Dependencies added after the QA Report`.
5. Success Metric: every existing reopen, QA prior-pass, config and implement
   test passes unchanged.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The adopter's report, triaged in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-06-relatorio-de-qa-sem-formatacao-derruba-a-proxima-run.md`,
  with both Run ids and the manual formatting and reopen workarounds.
- The adopter's mirrored repository (Secondbrain
  `projects/pantheon/mirror`, synced 2026-10-06): its repository guide says
  the repository Verification only checks and names a separate writing
  formatter target; its Spec 0047 QA report records that the formatter could
  not start inside the Agent sandbox (EPERM on the package-runner cache); and
  its Spec 0046 QA Task records the hand reopen after a supervisor-completed
  corrective Task.
- The lint-staged README (github.com/lint-staged/lint-staged): a configured
  command receives the committed file list as arguments, and ignoring files
  is the task's own configuration.

## Glossary

- adds: **Format Command**
- adds: **Late Dependency**

## Decisions

- Format with a configured command over the Spec's `qa/` files only, at two
  points, reverting on failure; see ADR-0249.
- Prove a Late Dependency from the manifest at the commit that added the
  newest QA Report; see ADR-0249.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
