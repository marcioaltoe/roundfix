---
task: task_03
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: completed
type: backend
complexity: high
---

# Task 03: A findings verdict stands for its head

## Overview

When the Delivery Queue parks an item as `review-findings`, `roundfix deliver
retry` returns it to `reviewing`, and `commandDeliveryWorkflow.Review` in
`internal/cli/deliver_workflow.go` runs `roundfix review` again at the same
head. The reviewer is asked again, so a dismissed finding either parks the item
forever or clears when the reviewer happens to answer differently. Re-asking
at an unchanged candidate until the answer changes is the permissive direction
ADR-0153 forbids. This Task makes a recorded findings verdict stand for its
repository, base, head and provider. It clears only when every finding is
dismissed with evidence at that head, reported as `findings-dismissed`, and a
changed candidate gets a fresh review. The record and ledger are read from the
Artifact Directory that `roundfix review` itself wrote, and the verdict is read
by the Delivery Queue owner, which archives and publishes on it. A reused
record must therefore never look like a clean review.

## Requirements

1. MUST add the record outcome `findings-dismissed`
   (`reviewOutcomeFindingsDismissed`), the field `dispositions` (the matched
   ledger entries, JSON `dispositions`, omitted when empty) and the field
   `reused` (JSON `reused`, omitted when false). `validateReviewRecord` MUST
   accept `findings-dismissed` only with findings text, no reason, at least one
   finding item and exactly one `dismissed` disposition per item, and MUST keep
   every existing rule.
2. MUST make `runReviewCommand` decide reuse after it resolves the base commit
   and the Artifact Directory and confirms the policy is `codex` or `claude`.
   The decision MUST come before `removeReviewAnswer`, profile resolution and
   readiness. When `pre-pr-review.json` names the same repository, base commit,
   head commit and provider with outcome `findings` or `findings-dismissed`,
   the command MUST NOT prepare, probe or prompt any Agent session. It joins the
   ledger entries whose repository, head, identity and text match an item into
   `dispositions`, and sets `reused: true`.
   - When every item has a `dismissed` entry, the outcome is
     `findings-dismissed` and the exit is `0`.
   - Otherwise the outcome is `findings` and the exit is `1`, with one stderr
     line naming each finding identity without a disposition.
   A `fixed` entry never counts toward clearing.
3. MUST NOT reuse a record whose repository, base, head or provider differs,
   or whose outcome is `reviewed`, `blocked` or `omitted`; such a run reviews
   afresh exactly as today. Under `none` or `coderabbit` the command keeps its
   current behavior.
4. MUST let `roundfix review dispose` (task_02) accept a record whose outcome
   is `findings-dismissed`, with every other rule unchanged.
5. MUST route `commandDeliveryWorkflow.Review` through a new
   `deliveryReviewResult(record reviewRecord, head string)
   (delivery.ReviewResult, error)`. It keeps the current mapping of
   `reviewed`, `findings`, `blocked` and a different head, and maps
   `findings-dismissed` at the same head to `delivery.ReviewOutcomeReviewed`.
6. MUST describe the reuse rule, the `findings-dismissed` outcome and its exit
   code in the review section of `docs/user-guide/commands.md` and the Pre-PR
   review section of `.agents/skills/roundfix/SKILL.md`, using the phrase
   `findings-dismissed`. The delivery sections of both MUST say that a retried
   item at an unchanged head advances once its findings are dismissed and that
   the reviewer is asked again only after the head changes. Regenerate
   `skills/roundfix/SKILL.md` with `make skills-sync`, and add the reuse rule
   to the Review Finding Disposition entry of `CONTEXT.md`.
7. MUST add `internal/cli/review_head_bound_test.go` with:
   - `TestReviewReportsFindingsDismissedWithoutAskingTheReviewer`
   - `TestReviewKeepsStandingFindingsWithoutAskingTheReviewer`
   - `TestReviewIgnoresADismissalOfDifferentText`
   - `TestReviewIgnoresAFixWhenClearingAHead`
   - `TestReviewAsksAgainAfterTheHeadMoves`
   - `TestReviewAsksAgainForADifferentBase`
   - `TestReviewAsksAgainForADifferentProvider`
   - `TestReviewNeverReusesACleanVerdict`
   - `TestReviewNeverReusesABlockedVerdict`
   - `TestDeliveryReviewResultAdvancesDismissedFindings`
   - `TestDeliveryReviewResultParksStandingFindings`
   Each `roundfix review` test asserts exit code, outcome and the exact number
   of Agent runner prepare and prompt calls: zero for each reuse case, one for
   each fresh-review case. The dismissal cases record their dispositions
   through `roundfix review dispose`, not by writing the ledger by hand.
8. MUST keep `TestReviewCommandExitsOneAndRecordsFindings`,
   `TestReviewRemovesAStaleAnswerFile` and
   `TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished` green and
   unchanged.

## Subtasks

- [ ] Add the outcome, the reuse decision and the delivery mapping.
- [ ] Document the rule in the review and delivery sections and the glossary.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] At an unchanged candidate, the reviewer is never asked again after a
      findings verdict.
- [ ] Only evidence-backed dismissals of every finding clear a head, and the
      result reads `findings-dismissed`, never `reviewed`.
- [ ] A new head, base or provider gets a fresh review, and the Delivery Queue
      advances a `findings-dismissed` item.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- creates: `internal/cli/review_head_bound_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewReportsFindingsDismissedWithoutAskingTheReviewer|TestReviewKeepsStandingFindingsWithoutAskingTheReviewer|TestReviewIgnoresADismissalOfDifferentText|TestReviewIgnoresAFixWhenClearingAHead|TestReviewAsksAgainAfterTheHeadMoves|TestReviewAsksAgainForADifferentBase|TestReviewAsksAgainForADifferentProvider|TestReviewNeverReusesACleanVerdict|TestReviewNeverReusesABlockedVerdict|TestDeliveryReviewResultAdvancesDismissedFindings|TestDeliveryReviewResultParksStandingFindings|TestReviewCommandExitsOneAndRecordsFindings|TestReviewRemovesAStaleAnswerFile|TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewReportsFindingsDismissedWithoutAskingTheReviewer TestReviewKeepsStandingFindingsWithoutAskingTheReviewer TestReviewIgnoresADismissalOfDifferentText TestReviewIgnoresAFixWhenClearingAHead TestReviewAsksAgainAfterTheHeadMoves TestReviewAsksAgainForADifferentBase TestReviewAsksAgainForADifferentProvider TestReviewNeverReusesACleanVerdict TestReviewNeverReusesABlockedVerdict TestDeliveryReviewResultAdvancesDismissedFindings TestDeliveryReviewResultParksStandingFindings TestReviewCommandExitsOneAndRecordsFindings TestReviewRemovesAStaleAnswerFile TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "findings-dismissed" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "findings-dismissed" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the eleven new named tests exists and no guide names `findings-dismissed`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The verdict stands for its head
- `_prd.md` → Goal 3; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contracts 1 and 4; Testing Approach 3

## Result

Implementation:

- Added the `findings-dismissed` review outcome plus omitted-when-empty
  `dispositions` and omitted-when-false `reused` record fields. Record
  validation requires findings text, no reason, at least one finding item, and
  exactly one matching evidence-backed dismissal per item.
- Added head-bound reuse before answer cleanup, Agent profile resolution, and
  readiness. Matching findings records join exact repository, head, identity,
  and text ledger entries; complete dismissals exit `0`, while standing
  findings exit `1` and name each undisposed identity without any Agent probe,
  prepare, or prompt.
- Kept moved repository/base/head/provider candidates and clean, blocked, or
  omitted records on the fresh-review path. A `fixed` disposition remains
  historical evidence and never clears the reviewed head.
- Allowed `review dispose` to read `findings-dismissed`, and routed Delivery
  Queue review records through `deliveryReviewResult`; dismissed findings at
  the expected head map to reviewed while standing findings still park.
- Documented reuse, `findings-dismissed`, exit codes, and retry behavior in the
  command guide, canonical Roundfix skill, generated skill copy, and Review
  Finding Disposition glossary entry. Regenerated the distributed skill with
  `make skills-sync`.

Focused checks:

- Before implementation, the new focused test failed to compile because
  `reviewOutcomeFindingsDismissed`, `reviewRecord.Dispositions`, and
  `reviewRecord.Reused` did not exist.
- The twelve focused head-bound and validation tests in
  `internal/cli/review_head_bound_test.go` passed. The eleven required public
  behavior tests cover every reuse, fresh-review, and delivery-mapping case;
  the additional validation test rejects malformed `findings-dismissed`
  records.
- `TestReviewCommandExitsOneAndRecordsFindings`,
  `TestReviewRemovesAStaleAnswerFile`, and
  `TestRetryReReviewsTheCurrentHeadWhenNoTaskIsUnfinished` passed unchanged in
  focused runs.
- `diff -r .agents/skills/roundfix skills/roundfix` exited `0` after skill
  regeneration.
- `GOCACHE=/tmp/roundfix-task03-gocache make verify-incremental` first reached
  only the sandbox's process-table permission boundary in two existing
  force-stop integration tests. The same incremental gate rerun with host
  process-table access passed vet, all Go packages, skill checks, and build.

Acceptance evidence:

- At an unchanged candidate, the dismissed, standing, text-mismatch, and fixed
  cases each observed zero Agent probes, prepares, and prompts.
- Only exact evidence-backed dismissals of every item produced
  `findings-dismissed`; partial, mismatched, and fixed dispositions remained
  `findings`, and invalid dismissed records were refused before writing.
- A moved head, different base, and different provider each observed exactly
  one Agent prepare and prompt. The delivery mapping advanced
  `findings-dismissed` and preserved standing findings as a parking result.

Follow-ups: none discovered within this Task's slice.
