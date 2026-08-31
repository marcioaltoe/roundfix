---
status: completed
type: backend
---

# Task: A finding blocks the rows it names

A governance finding blocks the whole matrix. Measured: one finding blocked
fifteen of nineteen rows, so a round that cost a full Agent Session reported
signal about governance and nothing about function.

## Work

- A mechanical finding blocks the rows it names. A matrix row the finding does
  not name is measured rather than blocked.
- Keep the ability to block widely: a finding that names every row still blocks
  every row. What goes away is the implicit cascade, not the reach.
- Change nothing about withholding. ADR-0096 keeps the Agent Session withheld
  when a blocking machine fact is present before a matrix exists; this Task
  changes only how blame is attributed across a matrix that already exists.
- Change no verdict rule, no row contract, and no typed blocked-cause count.
- Cover a finding that names one row in a matrix of several, asserting the
  unnamed rows are measured, and a finding that names all of them, asserting it
  still blocks all.

## References

- `_prd.md` → Goal 4, User Story 3, Core Feature 4
- `_techspec.md` → Build Order 2; Interfaces: `BlockedRows`
- ADR-0080 keeps the typed blocked-cause counts this Task does not touch

## Verification
- `grep -q "func (finding MechanicalFinding) BlockedRows" internal/speccheck/mechanical.go || exit 1; grep -q "TestFindingBlocksOnlyTheRowsItNames" internal/speccheck/mechanical_test.go || exit 1; go test -count=1 ./internal/speccheck ./internal/daemon`

## Result

Mechanical findings now materialize blocked rows only for `RowHint` values
that name rows in an existing QA matrix. When at least one matrix row remains
unnamed, the result seeds the report and lets the QA Agent measure the
remaining rows. A precondition refusal, an assigned-repair failure, a finding
before any matrix exists, or findings that name every matrix row still
withhold the Agent Session.

The seeded report remains fail-closed whenever a mechanical finding or repair
failure exists, even when the matrix continues. The QA row syntax, typed
blocked-cause counts, and verdict values were not changed.

Acceptance evidence:

- `TestFindingBlocksOnlyTheRowsItNames/one_named_row_leaves_the_others_measurable`
  builds a three-row matrix with only `R01` implicated and observes one blocked
  row with `Blocking=false`.
- `TestFindingBlocksOnlyTheRowsItNames/every_named_row_can_still_block_the_matrix`
  implicates `R01`, `R02`, and `R03`, then observes all three blocked with
  `Blocking=true`.
- `TestMechanicalFindingsWithoutRowHintsBlockTheirRefusalCode` preserves the
  pre-matrix refusal behavior, and `TestMechanicalStageWithholdsAgentSession`
  preserves Daemon withholding for a blocking mechanical result.
- `TestWriteMechanicalQAReportRecordsTheRefusal/scoped_finding_keeps_a_failing_verdict_while_the_matrix_continues`
  proves scoped attribution does not turn a finding-bearing report into a pass
  or a pre-matrix refusal. Existing materialization assertions keep the
  `blocked (finding: ...)` row form and `rows_blocked_finding` count.

Focused checks:

- Before the implementation edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/speccheck -run '^TestFindingBlocksOnlyTheRowsItNames$'`
  failed because the one-row finding still set `Blocking=true`.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/speccheck -run '^(TestFindingBlocksOnlyTheRowsItNames|TestMechanical.*)$'`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/daemon -run '^Test(WriteMechanicalQAReportRecordsTheRefusal|MechanicalStageWithholdsAgentSession|MechanicalStageSeedsReportBeforeAgentSession)$'`
  passed.
- `rtk git diff --check` passed.

Verification-feedback repair, attempt 1:

- The diagnostic artifact showed that `internal/speccheck` passed while shared
  Daemon QA fixtures were withheld before their matrices. Those fixtures still
  rendered `true` for the QA Task after Roundfix began owning that Task's
  Verification.
- The shared fixture now renders `spec.DerivedQAVerification(slug)` for QA Tasks
  and keeps `true` only for non-QA Tasks. Non-hermetic authoring checks skip QA
  commands because their effective command is supplied by Roundfix; authored
  QA text remains refused by `SC-QA-VERIFICATION-AUTHORED`.
- The focused clean-Spec reproduction first failed with
  `SC-VERIFY-NON-HERMETIC`, then passed after the ownership-boundary repair:
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/daemon -run '^TestQAMechanicalRequestCarriesTheGatePrecondition$/^a_clean_Spec_reaches_the_stage_refusing_nothing$'`.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/speccheck -run '^Test(NonHermeticVerificationSkipsRoundfixOwnedQAVerification|VerifyNonHermeticRegistersAtTasksStage|AuthoredQAVerificationIsRefusedByTaskName)$'`
  passed, covering the derived-QA exclusion, the non-QA detector, and authored
  QA refusal.
- A focused rerun of every Daemon test named in the attempt-1 diagnostic passed
  with `-count=1`. The Task's declared `## Verification` command was not rerun.

The Daemon still owns Task settlement.

## Carry-forward provenance

- Source Run: `run_20260831T171256Z_df3674a688059467`
- Source commit: `b4d7640dcbf615dce2f563b7d94d13c8d786f58e`
