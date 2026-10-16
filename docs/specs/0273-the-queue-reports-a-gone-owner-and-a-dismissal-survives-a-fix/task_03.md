---
task: task_03
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: pending
type: backend
complexity: high
---

# Task 03: A dismissal is accepted after a fix that leaves the finding's anchor unchanged

## Overview

`roundfix review dispose <id> --dismiss` refuses whenever HEAD differs from the
reviewed head, so a dismissal after another finding's committed fix needs a
reset round-trip (`_techspec.md` → Current behavior). This Task implements API
Contract 3: an ancestry check, an anchor comparison read from Git at both
commits and the optional `currentHead` ledger field, and documents it in the
`review` user guide. It is verifiable on its own through `roundfix review
dispose` on a temporary repository with a seeded Pre-PR Review Record.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-14,
   [a dismissal after a committed fix needs a reset](references/2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md),
   for its first symptom.
2. MUST implement API Contract 3 in the `--dismiss` branch of
   `runReviewDisposeCommand` through `proveDismissalAnchor`: the equality fast
   path unchanged; then `git merge-base --is-ancestor`, reading exit status 1
   as a head that does not descend and any other failure as
   `check dismissal ancestry: <error>`; then the recorded `anchor` only, never
   one parsed again from the finding text; then the anchored lines at the
   same 1-based numbers, byte for byte, at both commits, read through Git and
   never from the working tree, with a start line of 0 comparing the whole
   blob. Each refusal uses the exact reason API Contract 3 gives, starts
   stderr with `roundfix: review dispose refused: `, exits 2 and appends
   nothing.
3. MUST add `CurrentHead string` with the JSON key `currentHead` and
   `omitempty` as the last field of `reviewFindingDisposition`, set only on a
   dismissal accepted at a HEAD other than the reviewed head. `headCommit`
   stays the reviewed head (Invariants 5-7).
4. MUST NOT change `--fixed-by`, any other dispose refusal, the reuse of a
   findings record, the review correction proof, the Reviewer Lineage, the
   dispose help text or the `review` command.
5. MUST add, in the new `internal/cli/review_dismiss_after_fix_test.go`, each
   test calling `t.Parallel()` first:
   `TestReviewDisposeReproducesTheDismissalAfterAFixTranscript` and
   `TestReviewDisposeReproducesTheChangedAnchorTranscript`, asserting Surface
   Transcripts 4 and 5 whole, and that the ledger gained one line or none;
   `TestReviewDisposeRefusesADismissalWhenTheAnchorMoved` (a line inserted
   above the anchor); `TestReviewDisposeRefusesADismissalWhenTheAnchoredFileIsGone`;
   `TestReviewDisposeRefusesADismissalWithoutARecordedAnchor`;
   `TestReviewDisposeRefusesADismissalAtAHeadThatDoesNotDescend`;
   `TestReviewDisposeComparesTheWholeFileForALineZeroAnchor`;
   `TestReviewDisposeAtTheReviewedHeadWritesNoCurrentHead`; and
   `TestReviewCorrectionAcceptsADismissalRecordedAfterTheFix`, in which a
   parked candidate's two standing findings are disposed, one fixed by a
   commit and one dismissed after that commit, and `ProveReviewCorrection`
   accepts the item head. `TestReviewDisposeRefusesADismissalAtAMovedHead`
   MUST pass unchanged: its finding records no anchor. Every existing review
   dispose, correction and lineage test MUST pass unchanged.
6. MUST describe in `docs/user-guide/commands/review.md`, replacing the
   sentence that a dismissal requires an unchanged reviewed `HEAD` and the
   words `moved-head dismissals`, the rule of API Contract 3 with the phrase
   `anchored lines are unchanged`, the three refusals, and the `currentHead`
   field; then re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add `proveDismissalAnchor` and call it from the dismiss branch.
- [ ] Add `CurrentHead` to the disposition.
- [ ] Add the nine tests.
- [ ] Describe the rule in the `review` guide and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A dismissal at a descendant HEAD whose anchored lines are unchanged is appended once with the reviewed head and `currentHead`.
- [ ] A changed or moved anchor, a deleted anchored file, a missing anchor and a head that does not descend are refused with their reasons and an unchanged ledger.
- [ ] A dismissal at the reviewed head writes the same line as before.
- [ ] The review correction proof accepts a dismissal recorded after the fix.
- [ ] The `review` guide states the rule and no longer requires an unchanged reviewed `HEAD`.

## Context

- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`
- instruction: `internal/cli/review_validation.go`
- instruction: `internal/cli/review_lineage.go`
- instruction: `internal/cli/deliver_review_correction.go`
- instruction: `internal/cli/review_disposition_test.go`
- instruction: `internal/cli/deliver_review_correction_test.go`
- interface: `internal/cli/review.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/review_dismiss_after_fix_test.go`

## Consumers

- Disposition ledger line and `review dispose` output — consumers: `./internal/cli`
- Behavior Surface Record of the `review` guide — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReviewDisposeReproducesTheDismissalAfterAFixTranscript|TestReviewDisposeReproducesTheChangedAnchorTranscript|TestReviewDisposeRefusesADismissalWhenTheAnchorMoved|TestReviewDisposeRefusesADismissalWhenTheAnchoredFileIsGone|TestReviewDisposeRefusesADismissalWithoutARecordedAnchor|TestReviewDisposeRefusesADismissalAtAHeadThatDoesNotDescend|TestReviewDisposeComparesTheWholeFileForALineZeroAnchor|TestReviewDisposeAtTheReviewedHeadWritesNoCurrentHead|TestReviewCorrectionAcceptsADismissalRecordedAfterTheFix)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewDisposeReproducesTheDismissalAfterAFixTranscript TestReviewDisposeReproducesTheChangedAnchorTranscript TestReviewDisposeRefusesADismissalWhenTheAnchorMoved TestReviewDisposeRefusesADismissalWhenTheAnchoredFileIsGone TestReviewDisposeRefusesADismissalWithoutARecordedAnchor TestReviewDisposeRefusesADismissalAtAHeadThatDoesNotDescend TestReviewDisposeComparesTheWholeFileForALineZeroAnchor TestReviewDisposeAtTheReviewedHeadWritesNoCurrentHead TestReviewCorrectionAcceptsADismissalRecordedAfterTheFix; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'ReviewDispose|Disposition|Correction|Lineage|Reus' ./internal/cli || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- 'anchored lines are unchanged' || { printf 'review guide lacks the dismissal rule\n' >&2; exit 1; }; if tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- 'an unchanged reviewed'; then printf 'review guide still requires an unchanged reviewed HEAD\n' >&2; exit 1; fi; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the nine tests do not exist and the guide keeps the old rule, so the command fails; afterwards they and every existing dispose, disposition, correction, lineage and reuse test pass and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Core Feature 4; Success Metric 2
- `_techspec.md` → API Contract 3; API Contract 4; Surface Transcript 4; Surface Transcript 5; Invariants 5, 6 and 7; Data Models; Current behavior; Testing Approach; Build Order 3
- ADR-0279; ADR-0272; ADR-0196; ADR-0174
