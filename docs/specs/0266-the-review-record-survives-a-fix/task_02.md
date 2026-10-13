---
task: task_02
spec: 0266-the-review-record-survives-a-fix
status: pending
type: backend
complexity: medium
---

# Task 02: The review record is keyed and compared by the resolved checkout

## Overview

`roundfix review` keys its record by the checkout root Git reports, the
resolved path, while the Delivery Retry's correction proof keys it by the item
worktree path the Delivery Queue recorded, which under a symbolic-linked
`worktree.location` is the link spelling, so the proof finds no record
(`_techspec.md` → Current behavior). This Task names the review directory from
the resolved checkout path and routes every repository comparison over a
review record or a disposition through one helper. It is verifiable on its
own in `internal/cli`: under a worktree root behind a symbolic link, review,
dispose `--fixed-by` and the proof called with the link spelling end with the
proof accepted.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [the review correction proof misses its record under a symlinked worktree root](references/2026-10-13-the-review-correction-proof-misses-its-record-under-a-symlinked-root.md),
   by implementing `_techspec.md` → API Contract 1 in `reviewCheckoutDir`
   with `roundconfig.CanonicalizeResolvedPath`, keeping the 16-character
   SHA-256 digest and the `pre-pr-review/<key>` layout.
2. MUST implement `_techspec.md` → API Contract 2: add the unexported
   `sameReviewCheckout(left, right string) bool` in `internal/cli/review.go`
   and use it at every repository comparison the sweep lists in
   `internal/cli/review.go`, `internal/cli/review_lineage.go`,
   `internal/cli/review_validation.go` and
   `internal/cli/deliver_review_correction.go`, including
   `reviewDispositionsAtHead` and the proof's record check: the twelve sites
   that compare a `Repository` field with `==` or `!=` or through
   `filepath.Clean`. None may remain in those four files; the emptiness check
   in `validateReviewRecord` is not a comparison and stays. Recorded `repository` values and ledger lines MUST NOT be
   rewritten.
3. MUST add `internal/cli/review_resolved_checkout_test.go` with the four
   tests `_techspec.md` → Testing Approach names for task_02, each calling
   `t.Parallel()` first, building a worktree root `<tmp>/link` that is a
   symbolic link to `<tmp>/real` and adding the item worktree with
   `git worktree add <tmp>/link/repo-slug/item`:
   - `TestReviewCheckoutDirIsOneForEverySpellingOfACheckout`: the link and the
     resolved spelling name one directory, equal to the digest of the resolved
     path (Invariant 1), and a sibling worktree names another (Invariant 2).
   - `TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot`, with
     subtests `explicit-artifact-directory` and `default-artifact-directory`:
     `roundfix review` run in the item worktree with a fake reviewer raising
     one finding on `review.txt`, a commit changing only
     `docs/history/specs/0225-example/_prd.md`,
     `roundfix review dispose F1 --fixed-by <fix>`, then
     `ProveReviewCorrection` called with the link spelling, the slug
     `0225-example`, the reviewed head and the fix returns `Accepted: true`
     and no error.
   - `TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot`: the
     same flow with a fix that changes `review.txt`, which no Task of the
     archived Spec declares, is refused with a reason containing
     `changed path review.txt is outside the named archived Specs` (ADR-0165);
     the reason task_04 widens keeps that text as its prefix.
   - `TestTheCorrectionProofRefusesARecordOfAnotherCheckout`: a findings
     record naming a sibling worktree, placed in the item's review directory,
     is refused with a reason containing `review record must name the item worktree`.
4. MUST keep every existing test passing unchanged, including
   `deliver_review_correction_test.go`, `review_record_checkout_test.go`,
   `review_lineage_test.go` and `review_disposition_test.go`; a test that
   breaks is a defect in this Task, not a contract change.
5. MUST run the new tests once with `-race` (ADR-0259), give no wait a real
   time bound (ADR-0271), and MUST NOT change any command's output, flag,
   field or exit code, `internal/config`, or `internal/delivery`.

## Subtasks

- [ ] Resolve the checkout path in `reviewCheckoutDir`.
- [ ] Add `sameReviewCheckout` and replace every repository comparison of the sweep.
- [ ] Add the four tests under a symbolic-linked worktree root.

## Acceptance Criteria

- [ ] The proof accepts a docs-only disposed fix when called with the link spelling, for both Artifact Directory modes.
- [ ] A code fix and another checkout's record are still refused.
- [ ] No lexical repository comparison remains in the four review files.

## Consumers

- Pre-PR Review Record directory key and the review correction proof — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `internal/config/config.go`
- instruction: `internal/cli/deliver_review_correction_test.go`
- instruction: `internal/cli/review_record_checkout_test.go`
- instruction: `internal/cli/review_test.go`
- instruction: `docs/user-guide/commands/review.md`
- instruction: `docs/user-guide/commands/deliver.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_lineage.go`
- interface: `internal/cli/review_validation.go`
- interface: `internal/cli/deliver_review_correction.go`
- creates: `internal/cli/review_resolved_checkout_test.go`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestReviewCheckoutDirIsOneForEverySpellingOfACheckout|TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot|TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot|TestTheCorrectionProofRefusesARecordOfAnotherCheckout)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewCheckoutDirIsOneForEverySpellingOfACheckout TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot/explicit-artifact-directory TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot/default-artifact-directory TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot TestTheCorrectionProofRefusesARecordOfAnotherCheckout; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the tests do not exist, so their PASS lines are missing and the command fails; afterwards they pass under the race detector.
- `grep -q 'func sameReviewCheckout(' internal/cli/review.go || { printf 'sameReviewCheckout is missing\n' >&2; exit 1; }; matches="$(grep -nE '[.]Repository *(==|!=)|filepath[.]Clean[(][A-Za-z.]*Repository[)]' internal/cli/review.go internal/cli/review_lineage.go internal/cli/review_validation.go internal/cli/deliver_review_correction.go)"; test -z "$matches" || { printf 'lexical repository comparisons remain:\n%s\n' "$matches" >&2; exit 1; }; go test -count=1 -run 'Review|Correction|Dispos|Lineage' ./internal/cli || exit 1` — expected: exit 0. Before this Task the helper is missing and the lexical comparisons remain, so the command fails; afterwards the review, correction, disposition and lineage tests of the consuming CLI suite pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 3; Core Feature 1; Core Feature 2; Success Metric 1; Success Metric 3
- `_techspec.md` → Current behavior; API Contract 1; API Contract 2; Invariants 1-2; Invariant 4; Testing Approach; Build Order 2
- ADR-0272; ADR-0174; ADR-0233; ADR-0165; ADR-0268; ADR-0259; ADR-0271
