---
task: task_01
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: pending
type: backend
complexity: high
---

# Task 01: The Daemon records an Evidence Snapshot for every row that can carry

## Overview

The mechanical stage can carry a QA row only when the report that established
it recorded `evidence_snapshots`. Nothing writes that key today, so no row has
ever carried. This Task adds the recorder in `internal/speccheck` and calls it
from the Daemon's QA stage after the Agent turn, before the QA Report commit.
Each closed QA Report then records, at the audited head, the digests of every
declared input of each row that can carry. A report whose rows cannot carry
gets no key, and an Agent-written key never survives.

## Requirements

1. MUST add `speccheck.RecordEvidenceSnapshots` and `speccheck.EvidenceRecord`
   with the signatures the TechSpec's Interfaces section shows, in the new
   file `internal/speccheck/evidence_record.go`.
2. MUST record exactly the rows the TechSpec's "Recording the snapshot"
   section qualifies, with the snapshot shape of its Data Models section:
   `head`, then each declared input in declaration order with its files sorted
   by path and their SHA-256 digests read from the Git blobs at `head`. It
   MUST reuse `buildEvidenceSnapshots` and `parseMechanicalReport` rather than
   a second parser or a second digest reader.
3. MUST add `speccheck.AlwaysObserved` and the input kind
   `EvidenceCommitRange` (`commit_range`). `AlwaysObserved` returns the three
   reasons the TechSpec's "Always-observed rows" section lists. The kind and
   the three reasons MUST be constants in `internal/speccheck/report.go`. The
   recorder MUST never record a row it names.
4. MUST remove any existing `evidence_snapshots` key and its indented
   continuation lines. It MUST write the new block before the closing `---`
   only when a row qualifies, and it MUST leave every other byte of the report
   unchanged. On a Git read error it MUST still strip the key, and return the
   error.
5. MUST call the recorder in `runQAGate` after `settleQAVerdict` and before
   `settleTask`, on the settled report path, with the `auditedHead` resolved
   before the Agent turn. The recorded key MUST be in the QA Report commit.
6. MUST publish `daemon.qa` phase `evidence_snapshots` with outcome
   `recorded`, `none`, `skipped` or `error` and the payload of the TechSpec's
   API Contract 4. A report path outside the Run Worktree is `skipped`. A Git
   read error is `error`, and the gate still settles and commits. A failure to
   write the report is returned as an infrastructure error of the QA step, and
   a cancelled context publishes the stop.
7. MUST NOT change `Carriable`, `resolveCarriedRows`, the seeded report's
   bytes, any verdict or eligibility rule, or `internal/speccheck/mechanical_test.go`.
8. MUST put the new tests in the two test files this Task creates. The
   `speccheck` tests run over temporary Git repositories. The Daemon tests use
   the existing task-cycle fixture.

## Subtasks

- [ ] Add the input kind, the always-observed predicate and the recorder.
- [ ] Rewrite only the `evidence_snapshots` key of a report.
- [ ] Call the recorder in the QA stage and publish its event.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A passing row that declares only repository inputs is recorded under its
      row identifier with the head and the digests Git holds at that head.
- [ ] Rows that failed, are blocked, declare no inputs, declare a
      non-repository or `commit_range` input, name the repository Verification
      or the Pull Request row in their provenance, or name an input absent at
      the head are not recorded.
- [ ] An Agent-written `evidence_snapshots` value is replaced, and it is
      removed when no row qualifies.
- [ ] Every byte outside the key is unchanged, and `spec.ReadQAReportFile`
      reads the same verdict and typed counts before and after.
- [ ] A recorded report, read by `RunMechanicalStage` at the same head as the
      previous report, carries the recorded row. The writer and the existing
      reader agree.
- [ ] `AlwaysObserved` names each of its three reasons and reports false for
      any other row.
- [ ] Through the QA stage, the QA Report commit carries `evidence_snapshots`
      whose head is the audited head, and the `evidence_snapshots` event is
      published with outcome `recorded` and the row count.
- [ ] Through the QA stage, an Agent-written key on a report with no
      qualifying row is absent from the committed report, and the event
      outcome is `none`.

## Context

- interface: `internal/speccheck/report.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/speccheck/evidence_record.go`
- creates: `internal/speccheck/evidence_record_test.go`
- creates: `internal/daemon/qa_evidence_snapshot_test.go`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- instruction: `docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md`
- instruction: `docs/adr/0097-a-qa-row-carries-forward-only-on-declared-unmoved-evidence.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEvidenceRecordSnapshotsEveryQualifyingPassRow|TestEvidenceRecordSkipsRowsThatCannotCarry|TestEvidenceRecordReplacesAnAgentWrittenKey|TestEvidenceRecordKeepsEveryOtherByte|TestEvidenceRecordRoundTripsThroughTheMechanicalStage|TestAlwaysObservedNamesItsReason)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEvidenceRecordSnapshotsEveryQualifyingPassRow TestEvidenceRecordSkipsRowsThatCannotCarry TestEvidenceRecordReplacesAnAgentWrittenKey TestEvidenceRecordKeepsEveryOtherByte TestEvidenceRecordRoundTripsThroughTheMechanicalStage TestAlwaysObservedNamesItsReason; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the six named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead|TestQAGateStripsAnAgentWrittenSnapshotWhenNoRowQualifies|TestWriteMechanicalQAReportRecordsTheRefusal|TestWriteMechanicalQAReportWritesThePreconditionRefusal)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQAGateCommitsTheEvidenceSnapshotAtTheAuditedHead TestQAGateStripsAnAgentWrittenSnapshotWhenNoRowQualifies TestWriteMechanicalQAReportRecordsTheRefusal TestWriteMechanicalQAReportWritesThePreconditionRefusal; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails. The two existing byte-exact report tests run unedited in the same command.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 1; Core Feature 1; Core Feature 4; Success Metric 4; Success Metric 5
- [_techspec.md](_techspec.md) — Interfaces; Data Models; Recording the snapshot; Always-observed rows; API Contract 1; API Contract 2; API Contract 4; Testing Approach 1; Build Order 1
- ADR-0194; ADR-0195; ADR-0097; ADR-0096; ADR-0080; ADR-0057
