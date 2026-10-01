---
task: task_03
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: completed
type: backend
complexity: high
---

# Task 03: A second review reads the delta, and a third closes only on dispositions

## Overview

Every `roundfix review` today reviews the full merge-base diff and knows
nothing of the last round. This Task places each review in a Reviewer Lineage
using the checkout's own record:

- **Round 2** sends only the diff since round 1, and lists round 1's findings
  with their validation and disposition.
- **A third review** never calls the reviewer. It closes the lineage as
  `ceiling-closed` when every standing round-2 finding has an operator
  disposition. Otherwise it prints a blocked result and leaves the round-2
  record in place.

The Delivery Queue advances on `ceiling-closed`. Session continuation is
task_04; here every round opens and closes its own session, as today.

The lineage reads only the current checkout's record and the shared
disposition ledger, and it writes nothing to the ledger.

## Requirements

1. MUST add `internal/cli/review_lineage.go` with `decideReviewLineage`,
   implementing `_techspec.md` → The Reviewer Lineage (the table) exactly,
   including a blocked prior record that repeats its own round, and a legacy
   record without `lineage` that counts as round 1.
2. MUST add the top-level `lineage` field of `_techspec.md` → Data Models,
   with `round`, `previousHead`, `previousFindings`, `previousDispositions` and
   `reviewedHead`. It MUST add the outcome `ceiling-closed` with its
   validation rule. Older records MUST still read.
3. MUST build the round-2 prompt with `buildRoundTwoPrompt` as the TechSpec
   states. The delta is `git diff --no-ext-diff --no-textconv --no-color
   <previousHead> <head> --`. The prompt MUST NOT contain the full merge-base
   diff. The previous dispositions are the ledger entries for the prior
   record's repository and head.
4. MUST keep anchor validation on the full merge-base diff at the current head
   in every round (Invariant 6).
5. MUST implement `closeAtCeiling`. Either of two dispositions counts:
   `dismissed` with evidence, or `fixed` with a `fixedBy` that
   `git merge-base --is-ancestor <fixedBy> <head>` accepts.
   - When every standing round-2 finding has one, it persists a
     `ceiling-closed` record and exits `0` (Surface Transcript 6).
   - Otherwise it prints, without persisting, the blocked record of Surface
     Transcript 5, and exits `2`.
   - In both cases the runner is never called, and the answer file is never
     written or removed.
6. MUST map `ceiling-closed` to `delivery.ReviewOutcomeReviewed` in
   `deliveryReviewResult`, and leave every other mapping and park reason
   unchanged.
7. MUST put the new tests in `internal/cli/review_lineage_test.go`. They use
   temporary repositories with commits for each round, the review fixture's
   configuration and a fake runner that records each prompt and counts calls.

## Subtasks

- [ ] Decide the lineage round from the checkout's record.
- [ ] Build the round-2 prompt from the delta and the round-1 findings.
- [ ] Close or block at the ceiling without a runner call.
- [ ] Map `ceiling-closed` in the Delivery Queue adapter.
- [ ] Add a test for each acceptance criterion, with each negative case separate.

## Acceptance Criteria

- [ ] `TestReviewLineageDecidesEachRound`: one subtest per row of the lineage
      table, including another provider, a moved merge base after a rebase, a
      head that does not descend, a blocked round 2 that repeats round 2, and
      a legacy record.
- [ ] `TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings`: the
      prompt holds a line added after round 1 and not a line only the
      round-1 diff added. It holds `Previous head: <head 1>`, and each round-1
      finding with its validation and its disposition, or `no disposition`.
      The record reports `round` 2 and `previousHead`.
- [ ] `TestReviewRoundTwoAnchorsAgainstTheFullCandidateDiff`: a round-2
      finding anchored on a line only round 1's diff added stands.
- [ ] `TestReviewCeilingBlocksWithoutCallingTheReviewer` reproduces Surface
      Transcript 5 without the session fields task_04 adds. The runner is called zero times, and the round-2 record and
      answer bytes are unchanged afterwards.
- [ ] `TestReviewCeilingClosesOnDispositions` reproduces Surface Transcript 6,
      without the session fields task_04 adds, after `roundfix review dispose F1 --fixed-by <fix>`. With a `fixedBy`
      the head does not contain, it blocks instead.
- [ ] `TestDeliveryReviewResultAdvancesACeilingClosedRecord`: the mapping
      yields Reviewed, and `blocked` still yields Blocked.
- [ ] task_01's and task_02's tests,
      `TestReviewRecordRoundTripsEachOutcome` and
      `TestDeliveryReviewResultAdvancesDismissedFindings` still pass.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_lineage.go`
- creates: `internal/cli/review_lineage_test.go`
- interface: `internal/cli/deliver_workflow.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewLineageDecidesEachRound|TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings|TestReviewRoundTwoAnchorsAgainstTheFullCandidateDiff|TestReviewCeilingBlocksWithoutCallingTheReviewer|TestReviewCeilingClosesOnDispositions|TestDeliveryReviewResultAdvancesACeilingClosedRecord|TestReviewValidatorDismissesADaemonSettlementAsC2|TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne|TestReviewRecordRoundTripsEachOutcome|TestDeliveryReviewResultAdvancesDismissedFindings)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewLineageDecidesEachRound TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings TestReviewRoundTwoAnchorsAgainstTheFullCandidateDiff TestReviewCeilingBlocksWithoutCallingTheReviewer TestReviewCeilingClosesOnDispositions TestDeliveryReviewResultAdvancesACeilingClosedRecord TestReviewValidatorDismissesADaemonSettlementAsC2 TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne TestReviewRecordRoundTripsEachOutcome TestDeliveryReviewResultAdvancesDismissedFindings; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 4–5; User Stories 4–5; Core Features 6 and 8; Success Metrics 4 and 6; Declared breaks
- [_techspec.md](_techspec.md) — The Reviewer Lineage; Data Models; Invariants 4–6; API Contracts 1, 2 and 5; Surface Transcripts 5 and 6; Testing Approach 3; Build Order 3
- ADR-0197; ADR-0153; ADR-0165; ADR-0169; ADR-0174


## Result

Implemented the Reviewer Lineage using only the current checkout's record and
read-only access to the disposition ledger. Round 2 prompts carry the delta,
round-1 findings, validation and dispositions; anchor validation still uses the
full merge-base candidate diff. A blocked review repeats its recorded round
with its previous data. Descendants after round 2 close only with qualifying
operator dispositions, without preparing or calling a reviewer or changing the
answer file. Blocked ceiling calls print their result without replacing the
round-2 record. `ceiling-closed` advances the Delivery Queue as Reviewed.
Session continuation remains task_04's slice.

Focused-check evidence:

- Red starting point: `GOCACHE=/tmp/roundfix-task03-cache rtk proxy go test ./internal/cli -run '^TestReviewLineageDecidesEachRound$' -count=1`
  exited 1 because the lineage types, record field and decision function did
  not exist. The same focused check exited 0 after implementation.
- `GOCACHE=/tmp/roundfix-task03-cache rtk proxy go test ./internal/cli -count=1 -run 'TestReview(Lineage|RoundTwo|Ceiling)|TestDeliveryReviewResultAdvancesACeiling'`
  exited 0 with the six authored acceptance tests present.
- After the final source changes,
  `GOCACHE=/tmp/roundfix-task03-cache rtk proxy go test ./internal/cli -count=1 -run 'TestReview|TestDeliveryReviewResult'`
  exited 0, including existing task_01/task_02 review tests, record outcome
  round trips and findings-dismissed delivery mapping.
- `GOCACHE=/tmp/roundfix-task03-cache rtk make verify-incremental` initially
  exited 2: sandbox process-table access was denied, and suite guards caught
  source edits made while the check was running. With source files frozen,
  the same command rerun with process inspection allowed exited 0: formatting,
  vet, package tests, skill synchronization/checks and build all passed.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

| Criterion | Implementation and focused evidence |
| --- | --- |
| Lineage table | `TestReviewLineageDecidesEachRound` covers absence, provider/checkout changes, an actual rebase moving the merge base, a non-descendant head, same-head reuse, blocked round retries, legacy records, round-1 outcomes, round 2 and an already closed ceiling. |
| Round-2 prompt | `TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings` checks an added delta line, absence of the original-only diff line, previous head, standing and validation-dismissed findings, fixed/dismissed/missing dispositions, copied record data and separate session closure after each round. |
| Full-diff anchoring | `TestReviewRoundTwoAnchorsAgainstTheFullCandidateDiff` keeps a finding standing on the original round-1 line although that line is absent from the round-2 prompt's delta. |
| Blocked ceiling | `TestReviewCeilingBlocksWithoutCallingTheReviewer` checks exit 2, the exact reason and diagnostic, zero probe/preparation/prompt calls, and byte-identical round-2 record and answer. |
| Disposition closure | `TestReviewCeilingClosesOnDispositions` has separate contained-fix, evidence-dismissal and uncontained-fix cases. It uses the public dispose command, checks persisted closure, preserves answer and ledger bytes, makes zero reviewer calls, and checks another descendant still refers to the original round-2 head. |
| Delivery mapping | `TestDeliveryReviewResultAdvancesACeilingClosedRecord` checks Reviewed for ceiling closure and preserves Blocked and its reason. |
| Prior behavior | The broader focused selection passes task_01/task_02's tests, `TestReviewRecordRoundTripsEachOutcome` and `TestDeliveryReviewResultAdvancesDismissedFindings`. |

Additional checks cover a blocked round-2 retry's original delta and rejection
of malformed ceiling records, including missing evidence, wrong head/text and
duplicate dispositions.

The declared Verification command was not run. Task status and settlement
remain Daemon-owned; no other Task file or Task Graph manifest was edited, and
no commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261001T145917Z_95f0b38476c6ebcc`
- Source commit: `f7104e5a5ee5509aff6f39467d7c46b56fe1cbe5`
