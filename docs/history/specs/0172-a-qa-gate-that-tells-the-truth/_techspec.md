---
spec: 0172-a-qa-gate-that-tells-the-truth
status: active
created: 2026-09-25
surfaces: [backend, cli, docs]
---

# A QA gate that tells the truth

## Executive Summary

Seed QA Reports as `pending`, refuse a hollow `pass` or `partial` in the shared
acceptance decision, allocate a report where the reader looks, keep a
first-run deterministic failure across a temporary retry, classify an
unobserved Verification on the stream, and let `roundfix events` survive a
record it cannot project.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0015, ADR-0038, ADR-0056,
  ADR-0057, ADR-0080, ADR-0091, ADR-0093, ADR-0096, ADR-0097, ADR-0104,
  ADR-0111, ADR-0117, ADR-0130, ADR-0132, ADR-0135, ADR-0148, ADR-0155,
  ADR-0156, ADR-0159 and ADR-0160 hold; ADR-0020 and ADR-0127 do not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `internal/spec/archive_test.go`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `docs/references/coverage-record.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The pending seed and the hollow report

`internal/spec/qa.go` adds `VerdictPending = "pending"`. `readQAReport` reads it
as a supported verdict, and records on `QAReport` whether the report is hollow:
it has a `## Results` heading, no data row in any table between that heading
and the next heading of depth one or two, and no data row in any other table
whose header has a `Status` cell. A header row and a separator row are never
data rows. A report with no `## Results` heading is never hollow. The zero value
of the new field means not hollow, so a `QAReport` built in code keeps its
current eligibility.

`QAReportEligibility` refuses `pending` with the existing verdict message, and
refuses a hollow `pass` or a hollow otherwise-eligible `partial` with an error
containing `records no QA row`. Every existing refusal keeps its message and
precedence.

`mechanicalQAReportContent` in `internal/daemon/task_engine.go` writes
`verdict: pending` for a non-blocking result; a blocking result still writes
`fail`, and the precondition refusal report is unchanged. `settleQAVerdict`
returns the eligibility error beside the verdict, and the QA Task's reason for a
refused `pass` or `partial` becomes `QA verdict <verdict> not accepted: <cause>`;
`fail`, `pending`, `missing` and `unreadable` keep `QA verdict <verdict>`.
`matchesQAReportCommitMessage` in `internal/worktree/worktree.go` accepts a
`pending` QA commit, because the Daemon now writes one.

The qa-gate skill states that the seeded report starts pending and that a report
which records no QA row never settles the gate; `docs/user-guide/commands.md`
(archive, settle and `qa-report accept`) and the QA Report entry of
`CONTEXT.md` say the same.

## The allocator

`writeMechanicalQAReport` lists the existing `qa-report-<date>*.md` names of
its date, takes the highest numeric sequence (an unsuffixed report counts as
zero) and creates the next one: unsuffixed when none exists, `-01` after an
unsuffixed report, one above the highest otherwise. It re-reads on an exclusive
create collision. When a report dated after today already exists,
`NewestQAReport` would never select the new one, so allocation fails with an
error naming that report and writes nothing. The qa-gate skill's naming sentence
says one above the highest existing suffix.

## The temporary retry

`retainCollectedVerificationFailures` in `internal/daemon/task_engine.go`
treats a command the retry ended temporarily like one it could not observe: its
first-run deterministic failure is neither dropped from the merged set nor
replaced by the retry's temporary failure. The outcome keeps the retry's
`TemporaryFailure`, so `executeTask` still settles the Task failed without a
repair turn, and `taskVerificationFailureReason` names the first-run failure and
its diagnostic path.

## The unknown classification

`publishUnknownFailure` in `internal/daemon/engine.go` writes
`classification: verification_unknown`, `command`, `reason` (the cause's error
text, or `reason unavailable`) and `diagnostic_path` (or `unavailable`) on both
the `failed` and the `verdict` event, overriding a request's preset
classification: a precondition Verification the runner could not observe is
unknown, not a repository that was red on entry. `publishFailedCommand` and
`publishVerdict` for a command verdict are unchanged. The projection in
`internal/runevent/stream.go` already reads this shape.

## The event stream

`publishPreWorkProbeFindings` adds `commands`, the vacuous command strings in
probe order, beside the retained `probed_commands`. The vacuous projection reads
`commands`, and when it is absent derives it from the `probed_commands` entries
whose verdict is `passed`, so journals written before the fix project. In
`internal/cli/events.go`, a `runevent.ProjectStreamEvent` error skips that
record, writes one stderr warning naming its cursor, event kind and the error,
and continues in replay and in follow; a write or store error still exits `1`.
`docs/user-guide/commands.md`, the Roundfix skill's Supervisor Run Event Stream
section and the Supervisor Run Event Stream entry of `CONTEXT.md` describe the
warning, and the first two name the two classified record shapes.

## API Contracts

1. A QA Report may carry `verdict: pending`; it is readable and never accepted.
2. `QAReportEligibility` refuses a hollow `pass` or `partial` with an error
   containing `records no QA row`.
3. A refused `pass` or `partial` settles the QA Task with reason
   `QA verdict <verdict> not accepted: <cause>`.
4. A new QA Report's suffix is one above the highest of its date.
5. A `roundfix-events/v1` verification record of an unobserved outcome carries
   `classification: verification_unknown`, `command`, `reason` and
   `diagnostic_path`; a vacuous pre-work record carries
   `classification: verification_vacuous` and `commands`.
6. `roundfix events` skips an unprojectable record with a stderr warning naming
   its cursor, and exits `0` when nothing else fails.

## Coverage Map

- Goal 1 → The pending seed and the hollow report; API Contracts 1-3.
- Goal 2 → The allocator; API Contract 4.
- Goal 3 → The temporary retry.
- Goal 4 → The unknown classification; API Contract 5.
- Goal 5 → The event stream; API Contracts 5-6.
- Core Feature 1 → The pending seed and the hollow report.
- Core Feature 2 → The allocator.
- Core Feature 3 → The temporary retry.
- Core Feature 4 → The unknown classification.
- Core Feature 5 → The event stream.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.
- API Contracts 1-3 → The pending seed and the hollow report.
- API Contract 4 → The allocator.
- API Contract 5 → The unknown classification, The event stream.
- API Contract 6 → The event stream.

## Integration Points

- **Spec 0167.** Recorded the temporary-retry limit this Spec closes.
- **Spec 0139.** Has the Daemon run the QA gate's repository Verification
  through the same attempt publisher, so its unobserved outcomes gain the
  classification too.

## Testing Approach

1. **Pending and hollow.** An untouched seed settles `failed` with verdict
   `pending`; a hollow `pass` is refused by settlement, `archive`, `settle` and
   `qa-report accept`; a `pass` with a Results row, a matrix under a `Status`
   table, and a report without a `## Results` section are accepted;
   `TestArchivedPassCorpusRemainsArchiveEligible` stays green unchanged.
2. **Allocator.** A gap below the highest suffix is skipped, a fresh date starts
   unsuffixed, an unsuffixed report is followed by `-01`, a later-dated report
   refuses; after each allocation `NewestQAReport` returns the new path.
3. **Temporary retry.** The merged outcome keeps A's first-run failure and the
   Task's reason names its diagnostic; a command that passes on the retry drops
   its first-run failure; the existing retry tests stay green.
4. **Unknown classification.** Published payloads go through
   `runevent.ProjectStreamEvent` and are asserted field by field, for an
   ordinary and a precondition request, against a deterministic failure that
   projects unclassified.
5. **Event stream.** A published vacuous event projects its commands; a legacy
   journal projects from `probed_commands`; a record with neither key still
   fails projection; `roundfix events` warns and continues.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The pending seed and the hollow report (depends on: none).
2. The allocator (depends on: 1).
3. The temporary retry (depends on: 2).
4. The unknown classification (depends on: 3).
5. The event stream (depends on: 4).
6. Terminal QA (depends on: 1, 2, 3, 4, 5, 7).
7. The journal consumer corpus harness (depends on: 5).

## Risks & Considerations

- **One file, five Tasks.** Tasks 1, 2, 3 and 5 edit
  `internal/daemon/task_engine.go`, and Tasks 1, 4 and 5 edit
  `docs/user-guide/commands.md`, so the graph is a chain; each Task still puts
  its new daemon tests in a test file of its own, and only Task 1 edits
  `internal/daemon/task_engine_test.go`.
- **A stale skill.** Each guidance change ships in the Task that changes the
  behavior it describes.
