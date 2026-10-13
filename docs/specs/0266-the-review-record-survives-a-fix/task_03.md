---
task: task_03
spec: 0266-the-review-record-survives-a-fix
status: completed
type: backend
complexity: medium
---

# Task 03: Each reviewed head keeps its record, and the proof reads the candidate's

## Overview

A checkout keeps one Pre-PR Review Record, so review round 2 at the fixed head
replaces the record of the parked candidate's findings, and the correction
proof then refuses (`_techspec.md` → Current behavior). This Task keeps a
Reviewed Head Record for every non-reused verdict, under
`heads/<head commit>.json` in the checkout's review directory, and makes the
proof read the candidate's. It is verifiable on its own in `internal/cli`:
under a worktree root behind a symbolic link, review, dispose `--fixed-by`,
review round 2 and the proof called with the link spelling end with the proof
accepted.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [a second review round overwrites the parked findings record](references/2026-10-13-review-round-two-overwrites-the-parked-findings-record.md),
   by implementing `_techspec.md` → API Contract 3 in `persistReviewRecord`
   through the unexported `reviewHeadRecordPath`: after the current record is
   written, a record that is not `reused` and whose outcome is `reviewed`,
   `findings`, `findings-dismissed` or `ceiling-closed` is written with the
   same bytes, atomically, to `heads/<head commit>.json`. A `blocked` or
   `omitted` record, a reused record and a head that is not 40 or 64
   lowercase hexadecimal characters write no head record and leave an
   existing one byte-identical.
2. MUST implement `_techspec.md` → API Contract 4 in `ProveReviewCorrection`
   through the unexported `readReviewRecordAtHead`: the candidate's head
   record, else the current record when it names the candidate. A missing
   record is a refusal with the existing reason
   `review record must name the item worktree and parked candidate with outcome findings`,
   not an error. Every other check of the proof and its reason stay as they
   are.
3. MUST keep `pre-pr-review.json` as the record the next round, reuse and
   `roundfix review dispose` read; MUST NOT add a flag, a record field or a
   gc rule.
4. MUST add `internal/cli/review_head_record_test.go` with the six tests
   `_techspec.md` → Testing Approach names for task_03, each calling
   `t.Parallel()` first, on the symbolic-linked worktree root fixture of
   task_02 where a worktree is needed:
   - `TestASecondReviewRoundKeepsTheFirstRoundsRecord`: after round 1 with a
     finding, a fix, dispose `--fixed-by` and round 2 with no findings, the
     candidate's head record still has outcome `findings` and the round-1
     `findingItems`, the fix's head record has outcome `reviewed`, and
     `pre-pr-review.json` names the fix.
   - `TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot`:
     the same flow with a docs-only fix under `docs/history/specs/0225-example/`
     ends with the proof, called with the link spelling, accepted.
   - `TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord`, with
     subtests `reused`, `blocked` and `omitted`: each leaves no new head record
     and the existing one byte-identical.
   - `TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord`: a
     current findings record of the candidate with no `heads/` directory, as
     an item parked before this change has, is accepted.
   - `TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate`: a current
     record of another head and no head record of the candidate are refused
     with the existing reason and no error.
   - `TestReviewHeadRecordPathRefusesANonCommitName`, with subtests `empty`,
     `traversal`, `uppercase` and `short`: each returns false.
5. MUST keep every existing test passing unchanged, including
   `deliver_review_correction_test.go`, `review_record_checkout_test.go`,
   `review_lineage_test.go`, `review_head_bound_test.go` and
   `review_disposition_test.go`, and task_02's tests, with one exception:
   `TestTheCorrectionProofRefusesARecordOfAnotherCheckout` in
   `internal/cli/review_resolved_checkout_test.go` replaces only the current
   record, and the proof now prefers the candidate's Reviewed Head Record. Its
   fixture MUST also write the sibling-checkout record as the candidate's
   Reviewed Head Record, so the test keeps refusing the record the proof
   reads; its assertion and refusal text stay unchanged.
6. MUST run the new tests once with `-race` (ADR-0259), give no wait a real
   time bound (ADR-0271), and MUST NOT change any command's output, flag,
   field or exit code, or `internal/delivery`.

## Subtasks

- [ ] Write the Reviewed Head Record in `persistReviewRecord`.
- [ ] Read the candidate's record in `ProveReviewCorrection`.
- [ ] Add the six tests.

## Acceptance Criteria

- [ ] Round 2 leaves the candidate's findings record readable, and the proof accepts after it.
- [ ] Reused, blocked and omitted records write no head record.
- [ ] An item parked before this change is still proved from its current record.

## Consumers

- Reviewed Head Record and the review correction proof's reader — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `internal/cli/review_lineage_test.go`
- instruction: `internal/cli/review_head_bound_test.go`
- instruction: `docs/user-guide/commands/review.md`
- instruction: `docs/user-guide/commands/deliver.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/deliver_review_correction.go`
- creates: `internal/cli/review_head_record_test.go`
- interface: `internal/cli/review_resolved_checkout_test.go`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestASecondReviewRoundKeepsTheFirstRoundsRecord|TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot|TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord|TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord|TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate|TestReviewHeadRecordPathRefusesANonCommitName)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestASecondReviewRoundKeepsTheFirstRoundsRecord TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord/reused TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord/blocked TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord/omitted TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate TestReviewHeadRecordPathRefusesANonCommitName; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the tests do not exist, so their PASS lines are missing and the command fails; afterwards they pass under the race detector.
- `grep -q 'func reviewHeadRecordPath(' internal/cli/review.go || { printf 'reviewHeadRecordPath is missing\n' >&2; exit 1; }; grep -q 'readReviewRecordAtHead(' internal/cli/deliver_review_correction.go || { printf 'the proof does not read the candidate record\n' >&2; exit 1; }; go test -count=1 -run 'Review|Correction|Dispos|Lineage' ./internal/cli || exit 1` — expected: exit 0. Before this Task neither helper exists, so the command fails; afterwards the review, correction, disposition and lineage tests of the consuming CLI suite pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 3; Success Metric 2
- `_techspec.md` → Current behavior; API Contract 3; API Contract 4; Invariant 3; Invariant 4; Data Models; Testing Approach; Build Order 3
- ADR-0272; ADR-0197; ADR-0233; ADR-0259; ADR-0271


## Result

Implemented this Task's slice for Daemon Verification. Task status, Subtasks
and Acceptance Criteria checkboxes are unchanged; settlement belongs to the
Daemon.

`persistReviewRecord` encodes once, atomically writes the current record, and
then writes the identical bytes to `heads/<head commit>.json` for a non-reused
`reviewed`, `findings`, `findings-dismissed` or `ceiling-closed` verdict.
`reviewHeadRecordPath` admits only 40 or 64 lowercase hexadecimal characters.
Both record destinations use a temporary file, mode `0644`, and rename in a
`0755` directory. Reused, blocked, omitted and invalid-head records leave
existing head records untouched.

`ProveReviewCorrection` now uses `readReviewRecordAtHead` to prefer the parked
candidate's head record. When that file is absent, the current record is used
only if it names the candidate. Missing evidence refuses with the existing
record reason and no error; other read errors and proof checks retain their
behavior. The next round, reuse and disposition still read the current record.
The sibling-checkout refusal fixture now writes the sibling record to both
locations so its unchanged assertion covers the preferred record.

### Focused evidence

Before production changes,
`go test -count=1 -run 'TestASecondReviewRoundKeeps|TestTheCorrectionProofAcceptsAfterASecondRound' ./internal/cli`
exited `1`: the candidate's head record was absent after round 2, and the
proof refused with the existing record reason.

After the final production edit:

- `go test -race -count=1 -v -run 'Test(ASecondReviewRoundKeeps|TheCorrectionProofAcceptsAfterASecondRound|AReusedBlockedOrOmitted|TheCorrectionProofReadsTheCurrent|TheCorrectionProofRefusesWithoutARecord|ReviewHeadRecordPath)' ./internal/cli`
  exited `0`; all six required tests and their subtests passed under the race
  detector (`ok roundfix/internal/cli`, 2.616 s).
- `go test -count=1 -run '^Test(Review(RecordsLive|InOneCheckout|Dispose|IgnoresARecordAtTheShared|Lineage|RoundTwo|Ceiling|BlockedRoundTwo|Reports|KeepsStanding|IgnoresA|AsksAgain|NeverReuses|RecordRefuses|CorrectionRequires)|AReviewOnlyCorrection|ACorrection|TheCorrectionProof(Finds|RefusesA)|ReviewCheckoutDirIs)' ./internal/cli`
  exited `0` (`ok roundfix/internal/cli`, 3.901 s). This focused selection covers
  the affected existing checkout, disposition, lineage, head-bound and
  correction behavior, including task_02's sibling-checkout refusal.

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Round 2 preserves the candidate's findings and the proof accepts afterwards | `TestASecondReviewRoundKeepsTheFirstRoundsRecord` compares the candidate's full record, including its original findingItems, with round 1; confirms the fix's record and current record match round 2 with identical bytes; and checks file/directory permissions. `TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot` accepts the docs-only correction using the link spelling after disposal and round 2. |
| Reused, blocked and omitted records write no head record | `TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord` checks each exclusion against an existing head and a new head, preserving the former byte-identically and leaving the latter absent while still persisting the current record. `TestReviewHeadRecordPathRefusesANonCommitName` covers all four required invalid names, accepts SHA-1/SHA-256 names, and checks nonempty invalid names preserve existing head files. |
| An older parked item is still proved from its current record | `TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord` accepts the symlinked correction with no heads directory and leaves that directory absent. `TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate` confirms another current head or no current record produces the exact existing refusal reason without an error. |

All six new top-level tests call `t.Parallel()` first. No sleep, timer or
short real-time wait budget was added. Changes are confined to `review.go`,
`deliver_review_correction.go`, the new `review_head_record_test.go`, the
permitted fixture change in `review_resolved_checkout_test.go`, and this
Result section. No public flag, record field, output schema, exit code, gc
rule or `internal/delivery` file changed. Declared Task Verification and the
repository gates were not run; the Daemon owns them. No commit, push or Pull
Request was made. No follow-up implementation is included.
