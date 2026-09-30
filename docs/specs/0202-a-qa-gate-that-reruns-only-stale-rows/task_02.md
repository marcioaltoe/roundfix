---
task: task_02
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: pending
type: backend
complexity: high
---

# Task 02: A later pass carries unmoved rows and names why every other row re-runs

## Overview

With snapshots recorded (task_01), the mechanical stage still refuses a row
whose establishing head is not an ancestor of the current head, which is every
row of a failed pass in the Delivery Queue. It also carries rows whose truth
depends on the repository Verification, the Pull Request row or Task commits,
and it says nothing about the rows it did not carry. This Task accepts an
establishing head proven by the QA Report commit that recorded it, and refuses
always-observed rows. It records one Carry Disposition per prior row in the
seeded report, and keeps a carried row's establishing provenance.

## Requirements

1. MUST accept an establishing head when either proof of the TechSpec's
   "Carrying across an unintegrated pass" section holds: ancestry, or a
   recording commit reachable from any ref. The recording commit's first
   parent is the head. It carries a `Roundfix-Spec` trailer equal to the Spec
   directory holding the report and no `Roundfix-Task` trailer, and it holds
   the report path with the blob the establishing report's bytes hash to.
2. MUST refuse every row `speccheck.AlwaysObserved` names, even when an older
   report holds a snapshot for it.
3. MUST add `speccheck.CarryDisposition` and `MechanicalResult.Dispositions`,
   with exactly one disposition per row of the previous report, in Results
   order. A refused row's `Reason` MUST be one reason from the TechSpec's
   Vocabulary Contract, and a moved input names the moved paths. Every reason
   and the `Row carry-forward` heading MUST be constants in
   `internal/speccheck/report.go`.
4. MUST add `CarriedRow.Provenance` from the establishing row's provenance
   cell. `WriteMechanicalResult` MUST write it as the carried row's provenance
   and keep `report and head retained` when it is empty.
5. MUST make `WriteMechanicalResult` append the `## Row carry-forward` section
   of the TechSpec's "The seeded report" section only when `Dispositions` is
   non-empty. A result without dispositions MUST render exactly today's bytes.
6. MUST keep `Carriable`'s signature and every refusal it makes. It becomes
   `carryRefusal(...) == ""` over a new internal function that returns the
   reason. `ReportRow.AncestryVerified` keeps its name, and its comment names
   both proofs.
7. MUST add `carried_rows` and `rerun_rows` to the payload of the
   `daemon.qa` `mechanical` event.
8. MUST NOT edit `internal/speccheck/mechanical_test.go`. `TestCarriable`,
   `TestMaterializeMechanicalResult` and the four
   `TestMechanicalStageCarriable…` tests MUST pass unedited, including the
   refusal of a non-ancestor head that no QA Report commit recorded.
9. MUST put the new tests in the two test files this Task creates, over
   temporary Git repositories. The event test uses the existing task-cycle
   fixture.

## Subtasks

- [ ] Prove an establishing head by its recording commit.
- [ ] Refuse always-observed rows and return a disposition for every prior row.
- [ ] Render carried provenance and the `## Row carry-forward` section.
- [ ] Count carried and re-run rows in the mechanical event.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A row whose establishing head lies on another branch carries when a QA
      Report commit on that head recorded the establishing report with the
      same bytes and its declared inputs are byte-identical.
- [ ] The same row is refused `re-run: establishing head unproven` when the
      recording commit carries a `Roundfix-Task` trailer, and when no commit
      recorded the report on that head.
- [ ] Rows naming the repository Verification or the Pull Request row, and
      rows declaring `commit_range`, are refused with their
      `always observed: …` reason, even with a valid snapshot.
- [ ] A row whose declared input changed is refused with
      `re-run: input moved: <path>` naming that path.
- [ ] A carried row's provenance cell in the seeded report equals the
      establishing row's provenance.
- [ ] Every row of the previous report has exactly one disposition, and a
      failed row reads `re-run: not pass`.
- [ ] A result with no dispositions renders the same bytes as before this
      Task, and the existing byte-exact report tests pass unedited.
- [ ] The `mechanical` event payload counts carried and re-run rows.

## Context

- interface: `internal/speccheck/mechanical.go`
- interface: `internal/speccheck/report.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/speccheck/qa_row_carry_test.go`
- creates: `internal/daemon/qa_row_carry_event_test.go`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- instruction: `docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md`
- instruction: `docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRowCarryAcceptsAHeadItsQAReportCommitRecorded|TestRowCarryRefusesAHeadOnlyATaskCommitRecorded|TestRowCarryRefusesAnUnprovenHead|TestRowCarryNeverCarriesAlwaysObservedRows|TestRowCarryNamesTheMovedInput|TestRowCarryKeepsTheEstablishingProvenance|TestRowCarryRecordsOneDispositionPerPriorRow|TestRowCarryWithoutDispositionsRendersTodaysBytes|TestCarriable|TestMaterializeMechanicalResult|TestMechanicalStageCarriableCarriesUnchangedEvidenceWithCitation|TestMechanicalStageCarriableRefusesNonAncestorEstablishingHead|TestMechanicalStageCarriableRefusesChangedDeclaredInput|TestMechanicalStageCarriablePreservesOriginalEstablishingCitation)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRowCarryAcceptsAHeadItsQAReportCommitRecorded TestRowCarryRefusesAHeadOnlyATaskCommitRecorded TestRowCarryRefusesAnUnprovenHead TestRowCarryNeverCarriesAlwaysObservedRows TestRowCarryNamesTheMovedInput TestRowCarryKeepsTheEstablishingProvenance TestRowCarryRecordsOneDispositionPerPriorRow TestRowCarryWithoutDispositionsRendersTodaysBytes TestCarriable TestMaterializeMechanicalResult TestMechanicalStageCarriableCarriesUnchangedEvidenceWithCitation TestMechanicalStageCarriableRefusesNonAncestorEstablishingHead TestMechanicalStageCarriableRefusesChangedDeclaredInput TestMechanicalStageCarriablePreservesOriginalEstablishingCitation; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the eight new tests exists, so the command fails. The six existing carry and materialization tests run unedited in the same command.
- `out="$(go test -count=1 -v -run "^(TestQAMechanicalEventCountsCarriedAndRerunRows|TestWriteMechanicalQAReportRecordsTheRefusal|TestWriteMechanicalQAReportWritesThePreconditionRefusal)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQAMechanicalEventCountsCarriedAndRerunRows TestWriteMechanicalQAReportRecordsTheRefusal TestWriteMechanicalQAReportWritesThePreconditionRefusal; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new event test does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 3; Goal 4; User Story 1; User Story 3; Core Feature 3; Core Feature 4; Core Feature 5; Success Metric 1; Success Metric 2; Success Metric 5
- [_techspec.md](_techspec.md) — Interfaces; Carrying across an unintegrated pass; Always-observed rows; The seeded report; API Contract 3; API Contract 4; Vocabulary Contract; Testing Approach 2; Build Order 2
- ADR-0194; ADR-0195; ADR-0097; ADR-0080; ADR-0096
