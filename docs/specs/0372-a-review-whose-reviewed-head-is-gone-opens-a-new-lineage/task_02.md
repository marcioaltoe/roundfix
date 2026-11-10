---
task: task_02
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
status: pending
type: backend
complexity: medium
---

# Task 02: A recorded commit proven gone restarts the lineage and is refused by name

## Overview

On 2026-11-10 Conexus reported that `roundfix review` blocked every later
review of a checkout whose recorded reviewed head had been pruned after a
squash merge
([`references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md`](references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md),
reported by conexus): `decideReviewLineage` returns every `git merge-base
--is-ancestor` failure other than exit 1 as `check review lineage ancestry`.
This Task adds one proof that a commit a review record or the disposition
ledger names is gone, and makes the four ancestry checks that read such a
commit restart or refuse by name when it is, while every other Git failure
still blocks.

## Requirements

1. MUST add `internal/cli/review_recorded_commit.go` with the type
   `reviewCommitPresence` and its three values `reviewCommitNotProvenGone`,
   `reviewCommitGone` and `reviewCommitHistoryIncomplete`, and the functions
   of `_techspec.md` → Interfaces. `reviewRecordedCommitPresence(ctx, root,
   commit)` returns, in this order:
   - `reviewCommitNotProvenGone` when `commit` is not exactly 40 or 64
     lowercase hexadecimal digits;
   - `reviewCommitNotProvenGone` when `git rev-parse
     --is-shallow-repository` in `root` fails or prints neither `true` nor
     `false`;
   - `reviewCommitHistoryIncomplete`, without asking Git about `commit`,
     when `git config --get-regexp
     '^(remote\..+\.promisor|extensions\.partialclone)$'` exits 0 and prints
     an `extensions.partialclone` key or a promisor whose value is `true`;
     `reviewCommitNotProvenGone` when it exits with a status other than 0 or
     1;
   - after `git cat-file -e <commit>`: `reviewCommitGone` when it exits 1
     with empty standard output and standard error and the repository is not
     shallow; `reviewCommitHistoryIncomplete` for the same answer in a
     shallow repository; `reviewCommitNotProvenGone` for every other
     answer.

   Each Git command runs as `git -C <root> -c core.fsmonitor=false …` through
   `exec.CommandContext` with standard output and standard error captured
   apart. `reviewHistoryIncompleteError(commit)` returns the error text
   `a shallow or partial clone cannot prove that the recorded commit <commit> is gone; fetch the full history and try again`.
   `reviewedHeadGoneRefusal(ctx, root, reviewed)` returns
   `reviewed head "<reviewed>" is no longer in the repository; run roundfix review to open a new Reviewer Lineage`
   (the commit quoted as Go `%q` prints it) for `reviewCommitGone`,
   `reviewHistoryIncompleteError(reviewed)` for
   `reviewCommitHistoryIncomplete`, and `nil` otherwise.
2. MUST, in `internal/cli/review_lineage.go` `decideReviewLineage`, when the
   ancestry check of `prior.HeadCommit` fails with a status other than 1:
   for `reviewCommitGone` set `Restart` to
   `the reviewed head <prior head> is no longer in the repository` and
   return the round-1 plan with no prior and no error; for
   `reviewCommitHistoryIncomplete` return the error
   `check review lineage ancestry: ` followed by
   `reviewHistoryIncompleteError(<prior head>)`; otherwise return the error
   it returns today, unchanged. The ancestry check of `prior.BaseCommit`
   behaves the same with the reason
   `the recorded merge base <prior base> is no longer in the repository`.
   The proof runs only after such a failure, never before the ancestry
   check, and the comment that other Git failures must not reset the
   ceiling stays true.
3. MUST, in `reviewDispositionValidAtHead`, return `false, nil` when an
   ancestry check fails with a status other than 1 and
   `reviewRecordedCommitPresence` proves the disposition's fixing commit
   `reviewCommitGone`; every other failure keeps its
   `check disposition fixing commit ancestry` error.
4. MUST, in `internal/cli/review.go`, call `reviewedHeadGoneRefusal` for the
   record's reviewed head when the `review dispose --fixed-by` ancestry check
   of the reviewed head fails, and print its error as the refusal when it is
   not `nil`, before the existing
   `fixing commit … does not descend from reviewed head …` refusal; and in
   `proveDismissalAnchor`, when the ancestry check fails with a status other
   than 1, return `reviewedHeadGoneRefusal`'s error when it is not `nil`,
   before the existing `check dismissal ancestry` error. Every other refusal
   text stays byte-identical.
5. MUST, in `internal/cli/deliver_review_correction.go`
   `ProveReviewCorrection`, when the ancestry check of the parked candidate
   fails with a status other than 1: refuse with
   `parked candidate <candidate> is no longer in the repository` for
   `reviewCommitGone`; return the error `check review correction ancestry: `
   followed by `reviewHistoryIncompleteError(<candidate>)` for
   `reviewCommitHistoryIncomplete`; otherwise the error it returns today.
6. MUST add `internal/cli/review_recorded_commit_test.go` with the five tests
   of `_techspec.md` → Testing Approach, each calling `t.Parallel()` first:
   - `TestAReviewWhoseReviewedHeadIsGoneOpensANewLineage`: run a first review
     with the fake reviewer answering `No findings.`; check out a new branch
     from `main`, commit a file, delete `feature/review`, run
     `git reflog expire --expire=now --all` and `git gc --prune=now`, and
     assert that `git cat-file -e` of the old head exits 1 with no output;
     assert `decideReviewLineage` returns the restart reason of Requirement
     2, round 1 and no error, `reviewLineageNotice` returns
     `Review: <slug> head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: the reviewed head <prior head> is no longer in the repository`,
     and a second `roundfix review --base <base>` exits 0 with an empty
     stderr and a round-1 lineage.
   - `TestARecordedMergeBaseThatIsGoneOpensANewLineage`: a prior record
     whose head is an ancestor of the candidate and whose base is a full
     object ID the repository never held; assert the base restart reason.
   - `TestARecordedCommitThatIsNotProvenGoneStillBlocks`, with subtests
     `shallow clone` (a `--depth 1` clone of a three-commit source over
     `file://`, recording its first commit: the error is
     `check review lineage ancestry: a shallow or partial clone cannot prove that the recorded commit <commit> is gone; fetch the full history and try again`),
     `partial clone` (a `--filter=blob:none` clone over `file://` from a
     source with `uploadpack.allowFilter` true, whose remote URL is then set
     to a missing directory, recording an absent full ID: the same
     incomplete-history error), `object store error` (a `--no-local` clone
     whose `objects/info/alternates` names a missing directory, recording an
     absent full ID: the error starts with
     `check review lineage ancestry: git merge-base --is-ancestor <commit>`),
     `abbreviated commit` (a 12-digit prefix: the same prefix as
     `object store error`), and `Git failure on a present commit` (a fake
     `preflight.GitRunner` failing every `merge-base` with an error
     `ancestry unavailable` while the recorded head is present: the error
     starts with `check review lineage ancestry: ancestry unavailable`); in
     every subtest `Restart` stays empty. The candidate head of every
     subtest is the source's last commit, taken before any alternate is
     broken.
   - `TestReviewDisposeRefusesAReviewedHeadThatIsGone`, subtests `dismiss`
     and `fixed-by`: a findings record whose head is an absent full ID and a
     committed fix; `review dispose F1 --dismiss --evidence <text>` and
     `review dispose F1 --fixed-by <fix>` each exit 2 with stdout empty and
     stderr exactly `roundfix: review dispose refused: ` plus the refusal of
     Requirement 1, and the dismissal appends nothing to the ledger.
   - `TestAReviewCorrectionNamesACommitThatIsGone`, subtests `parked
     candidate` (the parked candidate an absent full ID: the refusal names
     `parked candidate <candidate> is no longer in the repository`) and
     `fixing commit` (the first fixed disposition's fixing commit an absent
     full ID: the refusal names `finding <id> requires exactly one valid
     disposition`); `ProveReviewCorrection` returns no error in either.

   The test file names no identifier this Task adds: it drives the existing
   functions and commands only, so the QA gate can compile it against the
   Delivery Base and record the behavior before this Task.
7. MUST move the last assertion of
   `TestReviewLineageRestartNamesItsReason/git_failures_stay_errors`
   (`internal/cli/review_main_merge_lineage_test.go`): a recorded head
   `strings.Repeat("f", 40)` now yields no error, round 1 and the restart
   reason `the reviewed head <that head> is no longer in the repository`,
   with a comment naming this Spec. This is the declared break: that case is
   the reported defect. Its five fake failures on present commits stay
   errors, unchanged, as do `TestReviewDisposeReportsAncestryFailureWithoutAppending`
   and the `Git failure` case of `TestReviewLineageNoticeIsAdvisory`.
8. MUST NOT change the review record, Reviewed Head Record or ledger formats,
   any other refusal or error text, `preflight.InspectGit` or its
   detached-HEAD refusal, any guide, skill or reference (task_01 owns them),
   or any file outside `## Context`. No test reaches the network, starts an
   Agent session, reads the real `~/.roundfix` or writes outside its
   temporary directories, and every test file compiles for Windows.
   `docs/agents/cli.md`, the CLI guide this Task follows, and the `review`
   and `deliver` guides, which task_01 updates, are read only.

## Subtasks

- [ ] Add the presence proof and the two refusal helpers.
- [ ] Call them from the lineage decision, the disposition validity check, `review dispose` and the correction proof.
- [ ] Add the five tests and move the pinned assertion.

## Acceptance Criteria

- [ ] A review whose recorded head was pruned runs round 1 of a new lineage, and `deliver status` names the reason.
- [ ] A shallow clone, a partial clone, a broken object store, an abbreviated commit and a Git failure on a present commit still block.
- [ ] `review dispose` and the correction proof refuse a gone commit by name.

## Context

- creates: `internal/cli/review_recorded_commit.go`
- creates: `internal/cli/review_recorded_commit_test.go`
- interface: `internal/cli/review_lineage.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/deliver_review_correction.go`
- interface: `internal/cli/review_main_merge_lineage_test.go`
- instruction: `internal/cli/review_lineage_notice.go`
- instruction: `internal/cli/review_dismiss_after_fix_test.go`
- instruction: `internal/cli/deliver_review_lineage_notice_test.go`
- instruction: `internal/cli/deliver_review_correction_test.go`
- instruction: `internal/cli/review_lineage_test.go`
- instruction: `internal/cli/review_disposition_test.go`
- instruction: `docs/agents/cli.md`
- instruction: `docs/user-guide/commands/review.md`
- instruction: `docs/user-guide/commands/deliver.md`

## Consumers

- `decideReviewLineage`, read by `roundfix review` and by `reviewLineageNotice` (`deliver status`, `deliver retry`), and `ProveReviewCorrection`, read by the Delivery Retry — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAReviewWhoseReviewedHeadIsGoneOpensANewLineage|TestARecordedMergeBaseThatIsGoneOpensANewLineage|TestARecordedCommitThatIsNotProvenGoneStillBlocks|TestReviewDisposeRefusesAReviewedHeadThatIsGone|TestAReviewCorrectionNamesACommitThatIsGone)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAReviewWhoseReviewedHeadIsGoneOpensANewLineage TestARecordedMergeBaseThatIsGoneOpensANewLineage TestARecordedCommitThatIsNotProvenGoneStillBlocks/shallow_clone TestARecordedCommitThatIsNotProvenGoneStillBlocks/partial_clone TestARecordedCommitThatIsNotProvenGoneStillBlocks/object_store_error TestARecordedCommitThatIsNotProvenGoneStillBlocks/abbreviated_commit TestARecordedCommitThatIsNotProvenGoneStillBlocks/Git_failure_on_a_present_commit TestReviewDisposeRefusesAReviewedHeadThatIsGone/dismiss TestReviewDisposeRefusesAReviewedHeadThatIsGone/fixed-by TestAReviewCorrectionNamesACommitThatIsGone/parked_candidate TestAReviewCorrectionNamesACommitThatIsGone/fixing_commit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the five tests exists, so their pass lines are missing; afterwards every test and subtest passes.
- `grep -qF -- 'is no longer in the repository' internal/cli/review_main_merge_lineage_test.go || { echo 'the pinned assertion still expects an error' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'Review|Lineage|Dispos|Dismiss|Correction|Deliver' ./internal/cli && GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the pinned assertion still expects an error, so the first check fails; afterwards it expects the restart, and every review, lineage, disposition, dismissal, correction and deliver test of the package passes with the package vetting for Windows.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 3; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Interfaces; API Contract 1; API Contract 2; API Contract 3; Testing Approach; Build Order 2
- ADR-0197; ADR-0272; ADR-0279
