---
task: task_01
spec: 0210-evidence-snapshots-that-stay-small
status: pending
type: backend
complexity: medium
---

# Task 01: An Evidence Snapshot records one digest per declared input and carries by it

## Overview

The Daemon's Evidence Snapshot writes one entry per file a declared input
matches, so a glob such as `**` writes every tracked file into the QA Report.
This Task keeps the matched files in memory and persists, per declared input,
one line holding its ref, its matched-file count and one SHA-256 digest of
Go's `h1:` summary of those files. The carry proof compares that pair per
input, names the inputs that moved, and still reads a report recorded in the
per-file form by converting its list into the same pair.

## Requirements

1. MUST add the `Count` and `SHA256` fields to `EvidenceSnapshot` and the
   unexported `evidenceInputDigest` and `evidenceSnapshotPair` with the shapes
   of the TechSpec's Interfaces section. The digest MUST be the lowercase hex
   SHA-256 of the summary defined in "The input digest": per matched file in
   ascending byte order of path, the lowercase hex SHA-256 of its content, two
   spaces, the path and `\n`. A path containing `\n` MUST make the input
   unresolved.
2. MUST make `RecordEvidenceSnapshots` render each input as one flow-style
   mapping line of `ref`, `count` and `sha256`, as Data Models shows. Its
   qualification rules, its coverage check over the in-memory files and its
   treatment of every byte outside the key MUST stay as they are.
3. MUST make `mechanicalEvidenceSnapshots` read both the new entry and the
   per-file entry (`ref` with `files`). `evidenceSnapshotPair` MUST derive the
   pair from a valid per-file list and MUST refuse an entry that holds both a
   file list and a recorded `count` or `sha256`.
4. MUST make `carryRefusal` follow the five steps of "Carrying by digest".
   A changed path under an input, or a differing pair, MUST yield
   `re-run: input moved: ` followed by the moved refs in declaration order
   joined by `, `, through the existing `CarryReasonInputMoved` constant.
   An invalid snapshot MUST yield `no evidence snapshot`, and an uncovered
   cited evidence path MUST yield `evidence differs`.
5. MUST NOT change `Carriable`'s signature, `resolveCarriedRows`, the prior
   pass import, the QA stage, the seeded report or any carry reason constant,
   and MUST NOT edit `internal/speccheck/mechanical_test.go`,
   `internal/speccheck/qa_row_carry_test.go` or
   `internal/daemon/qa_prior_pass_test.go`. Their tests MUST pass unedited.
6. MUST add the tests of the TechSpec's Testing Approach 1 to the new file
   `internal/speccheck/evidence_digest_test.go`, one top-level test per
   behavior, each computing its expected digest independently of the
   production function. The 2,500-file test MUST assert that the block's line
   count equals one plus three per recorded row plus one per declared input,
   and that the input's line carries `count: 2500`.
7. MUST update only these existing assertions to the new shape:
   `TestEvidenceRecordSnapshotsEveryQualifyingPassRow` and its decode helper
   in `internal/speccheck/evidence_record_test.go`;
   `TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead` in
   `internal/daemon/qa_evidence_snapshot_test.go`; and
   `assertTwoPassSnapshot` in `internal/daemon/qa_two_pass_carry_test.go`.
   Each MUST assert the count and the digest, not a file list.

## Subtasks

- [ ] Add the input digest and the snapshot pair.
- [ ] Render one line per input in the recorder.
- [ ] Read both entry shapes and carry by the pair, naming moved refs.
- [ ] Add the new tests and move the four existing assertions to the new shape.

## Acceptance Criteria

- [ ] A row whose glob matches 2,500 committed files records one line for that
      input with `count: 2500`, and the block's line count equals the TechSpec's
      formula.
- [ ] The recorded digest equals an independently built `h1` summary digest
      for files committed in non-sorted order.
- [ ] A new-form record carries through `RunMechanicalStage` at an unchanged
      head. After files under `src/` change, are added and are removed, the
      row re-runs with exactly `input moved: src/**`.
- [ ] A recorded count or digest that differs from the head's re-runs the row
      with `input moved: <ref>`.
- [ ] A per-file record of a two-file glob carries, and one missing a file
      re-runs with `input moved: <ref>`. An unsorted list, or an entry holding
      both `files` and `sha256`, re-runs with `no evidence snapshot`. In every
      case `RunMechanicalStage` returns no error and no finding.
- [ ] `Carriable` accepts an established snapshot holding only the pair
      against a current snapshot with files, and refuses a differing pair.
- [ ] The governed `TestCarriable` and `TestMechanicalStageCarriable…` tests,
      every `TestRowCarry…` test, the other `TestEvidenceRecord…` tests,
      `TestQAGateCarriesARowFromAnUnintegratedFailedPass` and the
      `TestPriorQAPass…` tests pass unedited.
- [ ] The QA Report commit and both two-pass tests carry the one-line entry.
      The second pass imports the first pass's new-form report and carries
      its unmoved row.

## Context

- interface: `internal/speccheck/report.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/speccheck/evidence_record.go`
- interface: `internal/speccheck/evidence_record_test.go`
- interface: `internal/daemon/qa_evidence_snapshot_test.go`
- interface: `internal/daemon/qa_two_pass_carry_test.go`
- creates: `internal/speccheck/evidence_digest_test.go`
- instruction: `internal/speccheck/mechanical_test.go`
- instruction: `internal/speccheck/qa_row_carry_test.go`
- instruction: `internal/daemon/qa_prior_pass_test.go`
- instruction: `docs/adr/0210-an-evidence-snapshot-records-one-digest-per-declared-input.md`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- instruction: `docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEvidenceRecordWritesOneLinePerInputWhateverItMatches|TestEvidenceInputDigestFollowsTheSortedSummary|TestCarryComparesTheRecordedDigestPerInput|TestCarryNamesTheMovedRefNotItsFiles|TestCarryReadsAPerFileSnapshotAsItsDigest|TestCarriableAcceptsARecordedDigestAgainstTheCurrentFiles|TestEvidenceRecordSnapshotsEveryQualifyingPassRow)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEvidenceRecordWritesOneLinePerInputWhateverItMatches TestEvidenceInputDigestFollowsTheSortedSummary TestCarryComparesTheRecordedDigestPerInput TestCarryNamesTheMovedRefNotItsFiles TestCarryReadsAPerFileSnapshotAsItsDigest TestCarriableAcceptsARecordedDigestAgainstTheCurrentFiles TestEvidenceRecordSnapshotsEveryQualifyingPassRow; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestCarriable|TestMechanicalStageCarriable.*|TestRowCarry.*|TestEvidenceRecord.*)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCarriable TestMechanicalStageCarriableCarriesUnchangedEvidenceWithCitation TestRowCarryNamesTheMovedInput TestRowCarryNamesADeletedInput TestEvidenceRecordRoundTripsThroughTheMechanicalStage TestEvidenceRecordWritesOneLinePerInputWhateverItMatches; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the governed and Spec 0202 carry tests run unedited beside the new recorder test, which does not exist before this Task.
- `out="$(go test -count=1 -v -run "^(TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead|TestQAGateStripsAnAgentWrittenSnapshotWhenNoRowQualifies|TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport|TestTwoGatePassesReRunARowWhoseInputMoved|TestQAGateCarriesARowFromAnUnintegratedFailedPass|TestPriorQAPass.*)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport TestTwoGatePassesReRunARowWhoseInputMoved TestQAGateCarriesARowFromAnUnintegratedFailedPass TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && ! grep -q 'path: snapshot-input.txt' internal/daemon/qa_evidence_snapshot_test.go && grep -q 'count' internal/daemon/qa_two_pass_carry_test.go` — expected: exit 0; before this Task the snapshot test still asserts the per-file entry and the two-pass helper decodes no count, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1-3; Core Features 1-4; Success Metrics 1-4
- [_techspec.md](_techspec.md) — Interfaces; The input digest; Data Models; Recording the snapshot; Carrying by digest; API Contract 1; API Contract 2; Integration Points; Testing Approach 1; Testing Approach 2; Build Order 1
- ADR-0210; ADR-0194; ADR-0195; ADR-0097
