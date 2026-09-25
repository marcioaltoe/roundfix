---
task: task_01
spec: 0172-a-qa-gate-that-tells-the-truth
status: completed
type: backend
complexity: high
---

# Task 01: A report starts pending and a hollow pass is refused

## Overview

`mechanicalQAReportContent` in `internal/daemon/task_engine.go` seeds every non-blocking QA Report with `verdict: pass` and a Results table that has a header and no row. When the QA Agent stops without filling it, or writes `pass` without adding a row, `QAReportEligibility` in `internal/spec/qa.go` accepts it, so QA settlement, `roundfix archive`, `roundfix settle` and `roundfix qa-report accept` all accept a gate that measured nothing. The report is read from the Spec's `qa/` directory, which the QA Agent writes, and the decision is read by the Daemon's settlement and by every maintainer command that archives or settles the Spec.

## Requirements

1. MUST add `VerdictPending = "pending"` to `internal/spec/qa.go`, read it in `readQAReport` as a supported verdict, and make `QAReportEligibility` refuse it with the existing message form `newest QA Report verdict is "pending"; expected "pass"`.
2. MUST record on `QAReport` whether the report is hollow, as a field whose zero value means not hollow: the report has a `## Results` heading, no data row in any table between that heading and the next heading of depth one or two (tables under its sub-headings count), and no data row in any other table whose header has a `Status` cell. Header and separator rows are never data rows. A report with no `## Results` heading is never hollow.
3. MUST make `QAReportEligibility` refuse a hollow `pass`, and a hollow `partial` that would otherwise be eligible, with an error containing `records no QA row`, and MUST keep every existing refusal's message and precedence.
4. MUST make `mechanicalQAReportContent` write `verdict: pending` for a non-blocking mechanical result, keeping `fail` for a blocking one and the precondition refusal report unchanged.
5. MUST make `settleQAVerdict` return the eligibility error, and settle a refused `pass` or `partial` QA Task with reason `QA verdict <verdict> not accepted: <cause>`; `fail`, `pending`, `missing` and `unreadable` keep the reason `QA verdict <verdict>`.
6. MUST make `matchesQAReportCommitMessage` in `internal/worktree/worktree.go` accept a QA commit whose verdict is `pending`, because the Daemon now writes one.
7. MUST state in `.agents/skills/qa-gate/SKILL.md` that the seeded report starts pending (the phrase `seeded report starts pending`) and that a report which records no QA row never settles the gate (the phrase `records no QA row`), and regenerate `skills/qa-gate/SKILL.md` with `make skills-sync`.
8. MUST state in `docs/user-guide/commands.md` (archive, settle and `qa-report accept`) and in the QA Report entry of `CONTEXT.md` that a `pending` verdict is never accepted and that a report which records no QA row is refused, using the phrase `records no QA row`.
9. MUST update the tests this change invalidates, and only those: the non-blocking case of `TestWriteMechanicalQAReportRecordsTheRefusal` now expects `pending`, and the non-qualifying partial row of `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` now expects the `not accepted` reason. `internal/spec/archive_test.go` stays byte-identical, and `TestArchivedPassCorpusRemainsArchiveEligible` stays green.
10. MUST put new daemon tests in `internal/daemon/qa_hollow_report_test.go`, keeping each Task's tests in a file of its own.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A pending report is readable and never accepted; an untouched seed settles the QA Task `failed` with verdict `pending`.
- [ ] A hollow `pass` is refused by settlement, `archive`, `settle` and `qa-report accept`, and a hollow otherwise-eligible `partial` is refused.
- [ ] A `pass` with a Results row, a matrix recorded under a `Status` table, and a report without a `## Results` section are accepted, and every archived `pass` Spec stays archive-eligible.
- [ ] A refused `pass` settles with a reason that names its cause, and a `pending` QA commit counts as a QA-report-only commit.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/qa.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/cli/qa_report.go`
- interface: `internal/cli/settle.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReadQAReportReadsAPendingVerdict|TestQAReportEligibilityRefusesAPendingVerdict|TestQAReportEligibilityRefusesAHollowPass|TestQAReportEligibilityRefusesAHollowPartial|TestQAReportEligibilityAcceptsAPassWithAResultsRow|TestQAReportEligibilityAcceptsAMatrixUnderAStatusTable|TestQAReportEligibilityKeepsAReportWithoutAResultsSection|TestArchivedPassCorpusRemainsArchiveEligible|TestMechanicalQAReportSeedStartsPending|TestQASettlementRefusesAnUntouchedSeed|TestQASettlementRefusesAHollowPass|TestQASettlementAcceptsAPassWithAResultsRow|TestQAReportOnlyBranchAcceptsAPendingQACommit|TestQAReportAcceptRefusesAPendingReport|TestQAReportAcceptRefusesAHollowPass|TestArchiveRefusesAHollowPass|TestSettleQATaskRefusesAHollowPass)$" ./internal/spec ./internal/daemon ./internal/worktree ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReadQAReportReadsAPendingVerdict TestQAReportEligibilityRefusesAPendingVerdict TestQAReportEligibilityRefusesAHollowPass TestQAReportEligibilityRefusesAHollowPartial TestQAReportEligibilityAcceptsAPassWithAResultsRow TestQAReportEligibilityAcceptsAMatrixUnderAStatusTable TestQAReportEligibilityKeepsAReportWithoutAResultsSection TestArchivedPassCorpusRemainsArchiveEligible TestMechanicalQAReportSeedStartsPending TestQASettlementRefusesAnUntouchedSeed TestQASettlementRefusesAHollowPass TestQASettlementAcceptsAPassWithAResultsRow TestQAReportOnlyBranchAcceptsAPendingQACommit TestQAReportAcceptRefusesAPendingReport TestQAReportAcceptRefusesAHollowPass TestArchiveRefusesAHollowPass TestSettleQATaskRefusesAHollowPass; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "records no QA row" docs/user-guide/commands.md && grep -q "records no QA row" CONTEXT.md && grep -q "seeded report starts pending" .agents/skills/qa-gate/SKILL.md && grep -q "records no QA row" .agents/skills/qa-gate/SKILL.md && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither phrase is documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — The pending seed and the hollow report

## Result

Implemented the pending seed and shared hollow-report refusal. QA Report
parsing now records whether a Results-bearing report has any QA data row,
eligibility rejects hollow `pass` and otherwise-eligible `partial` reports,
and the Daemon preserves the refusal cause when settling those verdicts. The
same decision now reaches archive, settle and `qa-report accept`; pending QA
commits remain recognizable as QA-report-only work.

Documentation now states the pending and hollow-report contract in the QA
skill, all three maintainer command sections, and the QA Report glossary entry.
`make skills-sync` regenerated the shipped skill copy, and the repository's
required digest regeneration reported no derived changes.

Focused checks:

- Red signal: `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache go test ./internal/spec -run '^TestReadQAReportReadsAPendingVerdict$'` failed to compile before implementation because `VerdictPending` and `QAReport.Hollow` did not exist.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache go test ./internal/spec` — passed.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache go test ./internal/daemon` — passed after correcting the untouched-seed fixture to preserve the real Daemon seed.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache go test ./internal/cli -run '^(TestQAReportAcceptRefusesAPendingReport|TestQAReportAcceptRefusesAHollowPass|TestArchiveRefusesAHollowPass|TestSettleQATaskRefusesAHollowPass)$'` — passed.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache go test ./internal/worktree -run '^TestQAReportOnlyBranchAcceptsAPendingQACommit$'` — passed.
- `make skills-sync` and `diff -r .agents/skills/qa-gate skills/qa-gate` — passed; canonical and shipped QA skills are byte-identical.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache make baseline-digests` — passed and reported no derived changes.
- `GOCACHE=/private/tmp/roundfix-task-0172-01-gocache make verify-incremental` — the first full run passed format, vet, the full Go suite, skill checks and build. A post-edit rerun reached one unrelated existing CLI failure: `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree` observed `Failed` instead of `Unresolved`; its immediate isolated `-count=1` rerun passed. Task 01 does not change temporary Verification flow.
- `git diff -- internal/spec/archive_test.go` — no output; the required archive corpus test file remains byte-identical.

Acceptance evidence:

- Pending readability/refusal and untouched-seed settlement are covered by `TestReadQAReportReadsAPendingVerdict`, `TestQAReportEligibilityRefusesAPendingVerdict`, `TestMechanicalQAReportSeedStartsPending`, and `TestQASettlementRefusesAnUntouchedSeed`.
- Hollow `pass` and qualifying `partial` refusals are covered at the shared decision, Daemon settlement, archive, settle and `qa-report accept` boundaries by the new named regression tests.
- Positive compatibility is covered by the Results-row, external Status-table and no-Results tests; the full incremental suite also ran `TestArchivedPassCorpusRemainsArchiveEligible` without changing `internal/spec/archive_test.go`.
- Refused-pass reason propagation and pending QA-only commit recognition are covered by `TestQASettlementRefusesAHollowPass` and `TestQAReportOnlyBranchAcceptsAPendingQACommit`.

The authored `## Verification` command was not run; Daemon Verification owns
that command and terminal Task settlement. No follow-up work was identified
inside this Task's slice; the unrelated full-suite observation above remains
recorded for the Daemon handoff rather than broadening this diff.
