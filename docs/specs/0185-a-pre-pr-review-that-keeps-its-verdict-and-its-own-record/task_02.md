---
task: task_02
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
status: completed
type: backend
complexity: medium
---

# Task 02: Each checkout keeps its own pre-PR review record

## Overview

`roundfix review` writes its record `pre-pr-review.json` and its answer `pre-pr-review-answer.txt` at the repository's Artifact Directory. Every checkout of a repository shares that directory: the main checkout, linked worktrees and each Delivery Queue item worktree. On 2026-09-29 a review in one checkout replaced another checkout's record and answer. The replaced record could no longer be reused or disposed, and its evidence was lost.

This Task gives each checkout its own record directory under the Artifact Directory. Every reader looks only there. The record's contents do not change, and the disposition ledger does not move. The files are written and read only by the Pre-PR Review Command and `roundfix review dispose` in the same checkout. The Delivery Queue keeps reading the verdict from the command's standard output.

## Requirements

1. MUST add `reviewCheckoutDir(artifactDir, checkoutRoot string) string` to `internal/cli/review.go`, returning `<artifactDir>/pre-pr-review/<key>`. `<key>` is the first 16 lowercase hex characters of the SHA-256 of `filepath.Clean(checkoutRoot)`, and `checkoutRoot` is the root `preflight.InspectGit` reports for the checkout.
2. MUST make `persistReviewRecord`, `persistReviewAnswer`, `removeReviewAnswer`, `reusableReviewRecord` and `runReviewDisposeCommand` use the current checkout's directory. The two persist functions MUST create it with mode `0755` before creating their temporary file inside it.
3. MUST keep the disposition ledger and its lock at `<artifactDir>/pre-pr-review-dispositions.jsonl`, and MUST keep the existing refusal of a record whose `repository` differs from the checkout.
4. MUST NOT read, write or remove `<artifactDir>/pre-pr-review.json` or `<artifactDir>/pre-pr-review-answer.txt`.
5. MUST update only the existing tests this change invalidates, and MUST keep their names: those that read or write the record or the answer at the old paths in `internal/cli/review_test.go`, `internal/cli/review_archived_spec_test.go` and `internal/cli/review_disposition_test.go` switch to `reviewCheckoutDir`. It MUST name each updated test in the Result, and MUST NOT rename or remove a top-level test or change an exported function signature.
6. MUST state in the `roundfix review` section of `docs/user-guide/commands.md` and in `.agents/skills/roundfix/SKILL.md` that the record and the answer live under `pre-pr-review/<checkout key>/` in the Artifact Directory and that Roundfix `never reads another checkout's record`, then MUST run `make skills-sync` so `skills/roundfix/SKILL.md` matches.
7. MUST add a **Pre-PR Review Record** entry to `CONTEXT.md` defining it as the record of one Pre-PR Review, kept per checkout, whose verdict is read from the reviewer's final message.
8. MUST put the new tests in `internal/cli/review_record_checkout_test.go`. They MUST use a temporary repository with a `git worktree add` checkout sharing one configured Artifact Directory, and the review command's fake runner, never a real reviewer.

## Subtasks

- [ ] Resolve one record directory per checkout.
- [ ] Write, reuse, clear and dispose only in that directory.
- [ ] Leave the old shared files untouched.
- [ ] Update the tests that named the old paths.
- [ ] Update the guide, the skill and its mirror, and the glossary.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A review writes its record and answer under `pre-pr-review/<key>/` for its checkout, and nothing at the old shared paths.
- [ ] After a review in the second checkout, the first checkout's record and answer bytes are unchanged. The first checkout still reuses its findings record without calling the runner, and still disposes a finding.
- [ ] `roundfix review dispose` in a checkout without a record is refused with `pre-pr review record does not exist`, even though the other checkout has one.
- [ ] A matching findings record at the old shared location is neither reused nor removed.
- [ ] The updated existing tests keep their names and stay green.
- [ ] The guide, the skill and its mirror, and the glossary describe the record per checkout, and `make skills-sync-check` passes.

## Context

- interface: `internal/cli/review.go`
- interface: `internal/cli/review_test.go`
- interface: `internal/cli/review_archived_spec_test.go`
- interface: `internal/cli/review_disposition_test.go`
- creates: `internal/cli/review_record_checkout_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewRecordsLiveInTheCheckoutsOwnDirectory|TestReviewInOneCheckoutLeavesAnotherCheckoutsRecord|TestReviewDisposeReadsOnlyItsCheckoutsRecord|TestReviewIgnoresARecordAtTheSharedLocation|TestReviewKeepsTheRawAnswer|TestReviewRemovesAStaleAnswerFile|TestReviewDisposeDismissesAFindingWithEvidence|TestReviewRecordsTheSpecsACandidateArchives)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewRecordsLiveInTheCheckoutsOwnDirectory TestReviewInOneCheckoutLeavesAnotherCheckoutsRecord TestReviewDisposeReadsOnlyItsCheckoutsRecord TestReviewIgnoresARecordAtTheSharedLocation TestReviewKeepsTheRawAnswer TestReviewRemovesAStaleAnswerFile TestReviewDisposeDismissesAFindingWithEvidence TestReviewRecordsTheSpecsACandidateArchives; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests do not exist, so the command fails.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "pre-pr-review/<checkout key>/" || { printf 'missing phrase in %s: %s\n' "$file" "pre-pr-review/<checkout key>/" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "never reads another checkout's record" || { printf 'missing phrase in %s: %s\n' "$file" "never reads another checkout's record" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the phrases are absent.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Pre-PR Review Record**" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "**Pre-PR Review Record**" >&2; exit 1; }` — expected: exit 0; before this Task the entry is absent.

## References

- [_prd.md](_prd.md) — Goals 4–5; Core Feature 4; Success Metric 5
- [_techspec.md](_techspec.md) — A review record per checkout; Interfaces; API Contract 2; Testing Approach 3; Build Order 2
- [references/2026-09-29-one-review-record-is-shared-by-every-checkout.md](references/2026-09-29-one-review-record-is-shared-by-every-checkout.md)
- ADR-0174; ADR-0153; ADR-0165; ADR-0169

## Result

Implemented one checkout-local Pre-PR Review Record directory at
`pre-pr-review/<checkout key>/` in the Artifact Directory. The checkout key is
the first 16 lowercase hexadecimal characters of the SHA-256 of the cleaned
checkout root. Record persistence, answer persistence, stale-answer removal,
reuse and `roundfix review dispose` now resolve that directory. Both persistence
paths create it with mode `0755` before creating their temporary files. The
shared disposition ledger and lock remain at the Artifact Directory root, and
the repository mismatch guard remains in place.

The four new command tests use one temporary Git repository, a linked checkout
created with `git worktree add`, one configured Artifact Directory and the fake
review runner:

- `TestReviewRecordsLiveInTheCheckoutsOwnDirectory` proves the exact key,
  directory mode, local record and answer paths, distinct checkout directories,
  and absence of the old shared files.
- `TestReviewInOneCheckoutLeavesAnotherCheckoutsRecord` proves byte isolation,
  checkout-local reuse without another runner call, and successful disposition
  from the first checkout.
- `TestReviewDisposeReadsOnlyItsCheckoutsRecord` proves the missing-record
  refusal when only the other checkout has a record and proves no ledger append.
- `TestReviewIgnoresARecordAtTheSharedLocation` proves a matching legacy record
  is not reused and both legacy record and answer remain byte-identical.

Existing top-level tests with path-specific setup or assertions were updated
without renaming them: `TestReviewKeepsTheRawAnswer`,
`TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached`,
`TestReviewRemovesAStaleAnswerFile`,
`TestReviewRecordsEmptySkippedSpecsAsAList`,
`TestReviewRecordsNoArchivedSpecForAnActiveSpec`, and
`TestReviewRecordListsEachFindingWithAnIdentity`. The shared review and
disposition fixtures now use `reviewCheckoutDir`; this preserves the existing
names and checkout-local behavior exercised by
`TestReviewRecordsTheSpecsACandidateArchives` and
`TestReviewDisposeDismissesAFindingWithEvidence`.

The command guide, canonical Roundfix skill and generated skill mirror now say
that the record and answer live under `pre-pr-review/<checkout key>/` and that
Roundfix never reads another checkout's record. `CONTEXT.md` now defines the
**Pre-PR Review Record** as one review's per-checkout record whose verdict comes
from the reviewer's final message.

Focused-check evidence:

- Before the implementation, the focused new-test build failed because
  `reviewCheckoutDir` did not exist.
- `go test -count=1 -run '^TestReview(RecordsLiveInTheCheckoutsOwnDirectory|InOneCheckoutLeavesAnotherCheckoutsRecord|DisposeReadsOnlyItsCheckoutsRecord|IgnoresARecordAtTheSharedLocation)$' ./internal/cli` passed.
- `go test -count=1 -run '^TestReview(KeepsTheRawAnswer|KeepsNoAnswerWhenTheReviewerWasNotReached|RemovesAStaleAnswerFile|RecordsEmptySkippedSpecsAsAList|RecordsNoArchivedSpecForAnActiveSpec|RecordListsEachFindingWithAnIdentity|DisposeDismissesAFindingWithEvidence|RecordsTheSpecsACandidateArchives)$' ./internal/cli` passed.
- `make skills-sync` passed and regenerated `skills/roundfix/SKILL.md` from the
  canonical skill; `cmp` confirmed the two files match.
- `make baseline-digests` passed with `changed:false`; no derived artifact
  changed.
- The first sandboxed `make verify-incremental` run reached the complete suite
  but two unrelated force-stop integration tests could not read the process
  table (`operation not permitted`). The rerun with process-tree permission
  passed `go vet`, every Go package, skill synchronization and checks, and the
  build.
- The final phrase search, old-shared-path inspection and `git diff --check`
  passed. The remaining old shared paths are negative assertions only.

The authored `## Verification` commands were not run; the Daemon owns them and
Task settlement.
