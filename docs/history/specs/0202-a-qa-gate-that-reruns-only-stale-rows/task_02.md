---
task: task_02
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: completed
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


## Result

Implemented this Task's carry-forward slice. The resolver accepts either
ancestry or a reachable QA recording commit that changed the establishing
report, has the establishing head as its first parent, names the report's Spec
in its trailer, has no Task trailer, and contains the exact report blob. The
changed-input comparison still compares trees across the two heads.

Every prior Results row receives one ordered disposition. Always-observed
rows are refused even with snapshots; failed rows, absent declarations,
unavailable establishing reports, unproven heads, absent snapshots, changed
inputs and differing evidence retain explicit refusal reasons. Moved inputs,
including deletions, name their paths. `Carriable` keeps its public signature
and existing refusal behavior through `carryRefusal`.

Carried rows retain establishing provenance and the original report/head
citation. The renderer preserves its previous bytes when dispositions are
absent, including the empty-provenance fallback; otherwise it appends the
separate Row carry-forward table. Mechanical events expose carried and re-run
counts, including the error payload.

### Focused implementation evidence

- Red starting check: `rtk proxy go test ./internal/speccheck -run
  '^TestRowCarry' -count=1` exited 1 because the new disposition fields,
  provenance field and reason constants did not yet exist.
- Final focused check: `GOCACHE=/private/tmp/roundfix-0202-task02-gocache rtk
  proxy go test ./internal/speccheck ./internal/daemon -run
  '^(TestRowCarry|TestQAMechanicalEvent|TestCarriable|TestMaterializeMechanicalResult|TestMechanicalStageCarriable|TestWriteMechanicalQAReport)'
  -count=1` exited 0 for both packages. The scoped cache avoids a shared
  Go-cache access restriction encountered by an earlier rerun.
- Incremental repository check: `GOCACHE=/private/tmp/roundfix-0202-task02-gocache
  rtk make verify-incremental` exited 0 after rerunning on a stable tree with
  elevated process access for the existing CLI force-stop tests. Vet, the Go
  suite, skill checks and build passed. Full output is retained at
  `/private/tmp/roundfix-0202-task02-incremental.log`. Earlier sandbox runs
  exited 2 because process-table reads were denied; the first also detected
  concurrent source edits, and the second detected the Result edit, through
  the repository mutation guard. The clean rerun had no concurrent edits.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- Acceptance 1: `TestRowCarryAcceptsAHeadItsQAReportCommitRecorded` checks
  unchanged input bytes across separate branches and retained report/head.
- Acceptance 2: `TestRowCarryRefusesAHeadOnlyATaskCommitRecorded` and
  `TestRowCarryRefusesAnUnprovenHead` each check the rendered unproven-head
  reason. Separate negative tests cover another Spec's trailer, different
  recorded bytes, a different first parent and an unreachable recording commit.
- Acceptance 3: `TestRowCarryNeverCarriesAlwaysObservedRows` separately checks
  repository Verification, Pull Request row and commit-range declarations,
  each with an existing snapshot.
- Acceptance 4: `TestRowCarryNamesTheMovedInput` and
  `TestRowCarryNamesADeletedInput` check the moved path and rendered reason.
- Acceptance 5: `TestRowCarryKeepsTheEstablishingProvenance` resolves a prior
  carried citation and checks the original provenance in the new report cell.
- Acceptance 6: `TestRowCarryRecordsOneDispositionPerPriorRow` checks Results
  order and one disposition each for a carried, failed and undeclared row,
  including `not pass` for the failed row. Separate tests check missing
  snapshots, unavailable establishing reports and non-repository inputs.
- Acceptance 7: `TestRowCarryWithoutDispositionsRendersTodaysBytes` checks
  exact pre-change renderer bytes and empty-provenance fallback. Existing
  `TestCarriable`, `TestMaterializeMechanicalResult`, the four
  `TestMechanicalStageCarriable…` tests and the `TestWriteMechanicalQAReport…`
  tests pass unchanged in the focused check. `mechanical_test.go` was not edited.
- Acceptance 8: `TestQAMechanicalEventCountsCarriedAndRerunRows` uses the
  existing task-cycle fixture and checks one carried and two re-run rows in
  the emitted mechanical event.

Task status remains Daemon-owned (`in_progress` was present on arrival).
No authored Verification command, other Task, Task Graph, commit, push or
Pull Request operation was performed. Prior-pass import remains task_03's
slice; no follow-up implementation was added here.

## Carry-forward provenance

- Source Run: `run_20261001T114323Z_551200526baffade`
- Source commit: `ff72c059ab96d3ce244b3456a514b725d1a667b3`
