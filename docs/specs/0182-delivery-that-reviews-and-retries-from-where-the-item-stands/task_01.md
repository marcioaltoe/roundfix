---
task: task_01
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: pending
type: backend
complexity: medium
---

# Task 01: The pre-PR review diffs the candidate from its merge base

## Overview

`roundfix review` diffs the tip of the selected base ref against the current head, so every commit the default branch gained after the candidate was cut reads as reverted by the candidate, and a main-side edit of an archived Spec parks a delivery item as `corrective-spec-required`. This Task resolves the merge base of the head and the base tip once in `runReviewCommand` and uses it for the diff, the changed-Spec list, the prompt's base commit, the record's `baseCommit` and the reuse key, keeping the tip as `baseTipCommit`. The base ref and head come from the user's repository through Git; the diff and prompt go to the configured reviewer in a read-only session, and the record lands in the Artifact Directory that `roundfix deliver` and `review dispose` read.

## Requirements

1. MUST add `resolveReviewMergeBase(ctx, gitRoot, tip, head string, runner preflight.GitRunner) (string, error)` in `internal/cli/review.go`, running `git merge-base <tip> <head>` through the given runner, and call it in `runReviewCommand` immediately after `resolveReviewBaseCommit`, before the record is created, the reuse lookup, the provider switch and any readiness probe.
2. MUST fail with exit `2` through `printReviewCommandFailure`, before any reviewer or readiness call, when Git fails or returns no merge base; the message MUST name both commits, say the head shares no history with the base, and end with `pass --base <ref>`.
3. MUST pass the merge base as the base commit to `newReviewRecord`, `reusableReviewRecord`, `runConfiguredReviewSession` (and through it `reviewCandidateDiff`, `reviewCandidateSpecContexts` and `buildReviewPrompt`) and every `runReviewSession` caller, and MUST NOT change the prompt's wording.
4. MUST add `BaseTipCommit string` with JSON name `baseTipCommit` and `omitempty` to `reviewRecord`, set it to the resolved tip, keep it out of the reuse comparison, and keep `validateReviewRecord` accepting a record without it.
5. MUST update the review section of `docs/user-guide/commands.md` and the Pre-PR review section of `.agents/skills/roundfix/SKILL.md` to contain the phrase `merge base of the current head and the selected base`, name `baseTipCommit`, state that a findings verdict stands while the base branch moves, and state that a head with no shared history exits `2`; both MUST drop the phrase `from the selected base to the current head`. Then MUST run `make skills-sync` so `skills/roundfix/SKILL.md` matches.
6. MUST add the sentence `The pre-PR review diffs the candidate from it.` to the **Delivery Base** entry of `CONTEXT.md`.
7. MUST put the new tests in `internal/cli/review_merge_base_test.go`, against real Git repositories through the public review command with the package's fake reviewer runner, reusing the existing review fixtures, and MUST keep every existing review test green without renaming any.

## Subtasks

- [ ] Resolve the merge base once and thread it through every base consumer.
- [ ] Record the tip as `baseTipCommit` and keep reuse keyed on the merge base.
- [ ] Refuse a head with no shared history before any reviewer call.
- [ ] Update the guide, the skill and its mirror, and the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] With the default branch advanced after the candidate was cut, the reviewer's diff contains only the candidate's change, `baseCommit` is `git merge-base <tip> HEAD` and `baseTipCommit` is the tip.
- [ ] A main-side edit of an archived Spec leaves `archivedSpecs` empty, while a candidate's own edit of an archived Spec is listed.
- [ ] A findings record written before the default branch moved is reused after it moves, with no reviewer call.
- [ ] A base ref sharing no history with `HEAD` exits `2`, names both commits and makes no reviewer call.
- [ ] The guide, the skill and its mirror, and the glossary state the merge-base rule, and `make skills-sync-check` passes.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_merge_base_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0169-the-pre-pr-review-diffs-the-candidate-from-its-merge-base.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewDiffsTheCandidateFromItsMergeBase|TestReviewListsOnlyTheArchivedSpecsTheCandidateChanged|TestReviewReusesAFindingsVerdictAfterTheBaseBranchMoves|TestReviewRefusesABaseThatSharesNoHistory)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewDiffsTheCandidateFromItsMergeBase TestReviewListsOnlyTheArchivedSpecsTheCandidateChanged TestReviewReusesAFindingsVerdictAfterTheBaseBranchMoves TestReviewRefusesABaseThatSharesNoHistory; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "merge base of the current head and the selected base" || { printf 'missing phrase in %s: %s\n' "$file" "merge base of the current head and the selected base" >&2; exit 1; }; grep -q "baseTipCommit" "$file" || exit 1; ! tr -s '[:space:]' ' ' < "$file" | grep -qF -- "from the selected base to the current head" || exit 1; done; make skills-sync-check` — expected: exit 0; before this Task the phrase and `baseTipCommit` are absent, so the command fails.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "The pre-PR review diffs the candidate from it." || { printf 'missing phrase in %s: %s\n' CONTEXT.md "The pre-PR review diffs the candidate from it." >&2; exit 1; }` — expected: exit 0; before this Task the sentence is absent.

## References

- [_prd.md](_prd.md) — Goals 1–2; Core Features 1–2; Success Metrics 1–2
- [_techspec.md](_techspec.md) — The review diffs from the merge base; API Contract 1; Testing Approach 1; Build Order 1
- ADR-0169; ADR-0153; ADR-0165
