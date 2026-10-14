---
task: task_02
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: completed
type: backend
complexity: medium
---

# Task 02: A legacy record states its override gaps, adopted sources and creation date

## Overview

A legacy override without approval or reason is written with empty fields, a
folder without `references/_index.md` declares `sources: []` beside a promoted
adopted source, and a PRD without a date yields `created: ""` although the
TechSpec or Git holds one (`_techspec.md` → Current behavior). This Task adds
the three fallbacks under `Legacy` and wires the creation commit through both
commands that plan conversions. It is verifiable on its own through
`PlanLegacyConversion` and through `roundfix history sanitize --apply` and
`roundfix baseline update` on temporary repositories.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [history sanitize records lose their summary, links and dates](references/2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md),
   for its third, fifth and sixth defects.
2. MUST implement API Contract 2 behind `in.Legacy`: with `qa_override: true`,
   an empty or blank approval or reason records `LegacyNotRecorded`
   (`legacy: not recorded`); a recorded value is kept byte for byte; the
   outcome, Task status and revision fields keep the legacy value, so
   `ArchivedTaskCompleted` answers as before.
3. MUST implement API Contract 3 behind `in.Legacy`: when the index is absent
   or lists no path, `sources` lists, sorted, the file names of the regular
   `.md` files directly under `references/` whose front matter holds `type`,
   `status` and `created`, or `status` and `created_at`. A file without front
   matter, a non-Markdown file and a subdirectory are not sources. An index
   that lists a path keeps today's reading alone.
4. MUST implement API Contract 4: the TechSpec fallback in `BuildArchiveRecord`
   behind `in.Legacy`; `FindLegacyCreation` in
   `internal/spec/history_sanitize.go`, running Git with the context and
   wrapping its error; `LegacyConversionRequest.Created`, applied by
   `PlanLegacyConversion` only when the folder supplied no date; and the
   wiring in `planHistoryUnit` and in the two re-plans with promotions
   (`internal/cli/history.go`, `internal/cli/baseline_update_history.go`),
   which carry `Created` from the first plan. A Git failure while reading the
   creation commit refuses the command as a delivery read does.
5. MUST NOT change the non-legacy `BuildArchiveRecord`, the review
   conventions, the Plan Digest format or any command output line (API
   Contract 6).
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. In the new `internal/spec/legacy_record_fields_test.go`:
   `TestLegacyOverrideWithoutApprovalRecordsNotRecorded`,
   `TestLegacyOverrideKeepsRecordedValues`,
   `TestLegacySourcesFallBackToAdoptedReferences` (one Backlog Entry, one
   Finding, one plain note, one `.txt` file; only the first two are listed),
   `TestLegacySourcesPreferTheIndex`, `TestLegacyCreatedFallsBackToTheTechSpec`,
   `TestFindLegacyCreationReadsTheOldestAddingCommit` (a PRD added under the
   Spec Root on one date and moved to the archive root on a later date yields
   the first date) and `TestActiveArchiveRecordKeepsItsFields`. In the new
   `internal/cli/history_legacy_fields_test.go`:
   `TestHistorySanitizeRecordsTheCreationCommitDate` and
   `TestBaselineUpdateHistoryRecordsTheCreationCommitDate`, the second through
   `--yes --promote` so the re-plan path is exercised; both parse the written
   record with `spec.ParseArchiveRecord` and assert `created` equals the
   commit's `%as` date. Every existing test MUST pass unchanged.
7. MUST describe in `docs/user-guide/commands/history.md`, in the paragraph on
   Lenient Legacy Reading, the `legacy: not recorded` value for a missing
   override approval or reason and the sources and creation-date fallbacks,
   with the phrase `legacy: not recorded`, and re-record the Behavior Surface
   Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add `LegacyNotRecorded` and the override fallback.
- [ ] Add the adopted-source fallback.
- [ ] Add the TechSpec date and `FindLegacyCreation`.
- [ ] Wire `Created` through both commands and their promotion re-plans.
- [ ] Add the nine tests.
- [ ] Describe the fields in the `history` guide and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A legacy override without approval or reason records `legacy: not recorded` in each missing field and nothing else changes.
- [ ] A folder without an index lists its adopted Backlog Entries and Findings as sources.
- [ ] A PRD without a date takes the TechSpec date, then the creation commit's date, in both commands.
- [ ] A non-legacy record keeps every field it has today.

## Consumers

- Archive Record front matter of a Legacy Archive Folder conversion — consumers: `./internal/spec`, `./internal/cli`
- `LegacyConversionRequest` and `FindLegacyCreation` — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/history_sanitize.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/baseline_update_history.go`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/spec/legacy_record_fields_test.go`
- creates: `internal/cli/history_legacy_fields_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacyOverrideWithoutApprovalRecordsNotRecorded|TestLegacyOverrideKeepsRecordedValues|TestLegacySourcesFallBackToAdoptedReferences|TestLegacySourcesPreferTheIndex|TestLegacyCreatedFallsBackToTheTechSpec|TestFindLegacyCreationReadsTheOldestAddingCommit|TestActiveArchiveRecordKeepsItsFields)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLegacyOverrideWithoutApprovalRecordsNotRecorded TestLegacyOverrideKeepsRecordedValues TestLegacySourcesFallBackToAdoptedReferences TestLegacySourcesPreferTheIndex TestLegacyCreatedFallsBackToTheTechSpec TestFindLegacyCreationReadsTheOldestAddingCommit TestActiveArchiveRecordKeepsItsFields; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Legacy|Archive|History' ./internal/spec || exit 1` — expected: exit 0. Before this Task the seven tests do not exist, so the PASS check fails.
- `out="$(go test -count=1 -v -run '^(TestHistorySanitizeRecordsTheCreationCommitDate|TestBaselineUpdateHistoryRecordsTheCreationCommitDate)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizeRecordsTheCreationCommitDate TestBaselineUpdateHistoryRecordsTheCreationCommitDate; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'History|Legacy|ArchiveRecord|BaselineUpdate|Review.*Override|ConventionC5' ./internal/cli || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/history.md | grep -qF -- 'legacy: not recorded' || { printf 'history guide lacks the override value\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the two command tests do not exist and the guide lacks the value, so the command fails; afterwards they and the history, baseline update, archive record and C5 override convention tests pass.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 2; Core Feature 4; Success Metric 1; Success Metric 3
- `_techspec.md` → API Contract 2; API Contract 3; API Contract 4; API Contract 6; Invariants 1, 4, 5 and 6; Current behavior; Testing Approach; Build Order 2
- ADR-0277; ADR-0254; ADR-0251


## Result

Implemented this Task's field-reading slice of the adopted Backlog Entry
`2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md`
(defects three, five and six). Task status and settlement remain Daemon-owned.

### Implementation and acceptance evidence

- Missing override approval or reason: `LegacyNotRecorded` fills only empty or
  blank fields under `Legacy`. Recorded whitespace, tabs and multiline text
  survive; outcome, QA Task status and revision retain their legacy values.
  `TestLegacyOverrideWithoutApprovalRecordsNotRecorded` covers absent, empty
  and blank fields; `TestLegacyOverrideKeepsRecordedValues` covers recorded
  fields, one missing field, and the unchanged `ArchivedTaskCompleted` answers
  for empty, pending and failed QA Task status.
- Adopted sources: a missing or empty index falls back to sorted regular
  Markdown files immediately under `references/` with Backlog Entry or Finding
  front matter. Plain notes, non-Markdown files, nested files and non-mapping
  front matter are excluded. `TestLegacySourcesFallBackToAdoptedReferences`
  exercises those cases; `TestLegacySourcesPreferTheIndex` proves that a
  populated index keeps its paths and order without adding fallback sources.
  Both command regressions also assert the adopted source survives conversion
  and promotion.
- Creation date: `BuildArchiveRecord` reads a valid TechSpec date only under
  `Legacy` when the PRD supplies none. `PlanLegacyConversion` applies
  `LegacyConversionRequest.Created` only when neither document supplied a
  date. `FindLegacyCreation` reads the oldest adding commit reachable from
  `HEAD` for both PRD paths, honors cancellation, and wraps Git failures.
  `planHistoryUnit` propagates that read failure through the existing delivery
  read error path; both promotion re-plans retain the first plan's `Created`.
  `TestLegacyCreatedFallsBackToTheTechSpec` covers date precedence, timestamps,
  invalid dates and absence of any date; the folder reader stops at the
  TechSpec. `TestFindLegacyCreationReadsTheOldestAddingCommit` covers a later
  archive move, direct archive additions, no matching commit, cancellation,
  unsafe paths and wrapped Git failure. Both named CLI regressions parse the
  written records and compare the Git fallback to the creating commit's `%as`,
  distinct from the later archiving commit; they also exercise PRD and TechSpec
  precedence. History sanitize is exercised with and without promotion;
  baseline update is exercised through `--yes --promote`.
- Non-legacy behavior: `TestActiveArchiveRecordKeepsItsFields` compares every
  record field and rendered bytes against the existing behavior, including
  blank override fields, no adopted-source or TechSpec-date fallback, and
  original index ordering. Existing archive, override-convention and promotion
  transcript checks listed below also pass unchanged. Production changes do
  not alter review conventions, record parsing, digest serialization or output
  statements.

All nine required top-level tests, and their subtests, start with
`t.Parallel()`. The history guide documents the three readings with the exact
`legacy: not recorded` value. The sanctioned generator re-recorded the Behavior
Surface Record; only the history guide fingerprint changed.

### Focused checks

- Before production edits:
  `go test ./internal/spec -count=1 -run '^TestLegacy(OverrideWithoutApproval|SourcesFallBack|CreatedFallsBack)'`
  exited 1, reproducing blank override fields, missing adopted sources and
  missing TechSpec dates.
- After the final test edits:
  `go test ./internal/spec -count=1 -run 'TestLegacy(Override|Sources|Created)|TestFindLegacyCreation|TestActiveArchiveRecordKeepsItsFields|^TestArchiveRecordKeepsEveryOverrideField$|^TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder$|^TestReadArchivedSpecReadsRecordAndLegacyFolder$'`
  exited 0. This selects the seven new field tests plus three existing archive
  compatibility tests.
- `go test ./internal/cli -count=1 -run 'RecordsTheCreationCommitDate$'`
  exited 0 for both new command regressions and their date/promotion cases.
- `go test ./internal/cli -count=1 -run 'RecordsTheCreationCommitDate$|^TestHistorySanitizePromotesIntoReferences$|^TestBaselineUpdatePromotesHistoryCitations$|^TestOverrideArchiveIsEligibleForConventionC5$|^TestConventionC5ValidatorHonorsEligibility$'`
  exited 0 for the new command cases, existing sanitize promotion and two C5
  checks. The `TestBaselineUpdatePromotesHistoryCitations` alternative matches
  no existing test; baseline transcript evidence is the separate check below.
- `go test ./internal/cli -count=1 -run '^TestBaselineUpdatePromotesOnlyOnRequestTranscript$'`
  exited 0, preserving the existing full baseline promotion transcript.
- Required regeneration:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 and wrote `docs/references/behavior-surfaces.json`.
- `git diff --check` exited 0.

The Task's declared Verification commands, the selected repository Verification
and the full suite were not run; they belong to the Daemon. No existing test,
Task Graph manifest, other Task file, review convention or tooling configuration
was edited. No repository commit, push or pull request was made. No follow-up
outside this Task's slice was introduced.

## Carry-forward provenance

- Source Run: `run_20261014T193608Z_0f6bbae56906b671`
- Source commit: `18fd5b24942828c78dce1de2fb793a70f31fa5e3`
