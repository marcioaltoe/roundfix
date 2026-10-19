---
task: task_02
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
status: pending
type: backend
complexity: medium
---

# Task 02: The Delivery Base, the archive warning and the review base read the Default Branch Head

## Overview

The Git Locator was not the only reader of the default branch that answered
from a stale local branch or refused without origin/HEAD (`_techspec.md` →
Current behavior, "Other default-branch callers, swept").

- The QA Delivery Base is unresolved when origin/HEAD is unset on a work
  branch.
- The Archive Command checks `source_revision` against the local branch on a
  name match, and against nothing without origin/HEAD.
- The pre-PR review's default base falls back to the local `main`.

This Task makes all three read the Default Branch Head that task_01 added. It
is also the Integration Checkpoint that runs the `cli` and `daemon` suites
against task_01's detection step. It is verifiable on its own through
temporary repositories without origin/HEAD.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-a-git-locator-in-a-fresh-record-fails-the-check-in-ci.md`
   (`_prd.md` → Sources 1), whose expected behavior is a check that "does not
   trust a stale local branch over the fetched remote". It does so by
   implementing API Contracts 4 to 6 of `_techspec.md` (ADR-0295).
2. MUST make `qaDeliveryBase` in `internal/daemon/qa_delivery_base.go` merge
   `HEAD` with the `Ref` of `preflight.ResolveDefaultBranchHead` (API Contract
   4). `found=false` leaves the base unresolved and a Git error is returned,
   as today.
3. MUST make `archiveDefaultRef` in `internal/cli/archive.go` return the
   head's `Ref`, displayed as `origin/<name>` for a remote-tracking ref and
   `<name>` for a local one (API Contract 5). `warnArchiveSourceRevision` and
   `recordArchiveCommandPullRequest` keep every text, including both "the
   default branch is unknown" texts for `found=false`.
4. MUST make `resolveReviewBaseCommit` in `internal/cli/review.go` use the
   head's `Ref` when `--base` is omitted (API Contract 6). `found=false` keeps
   the `refs/heads/main` then `refs/heads/master` fallback and its error text,
   and `--base <ref>` is unchanged.
5. MUST leave the callers the TechSpec names as unchanged byte-identical:
   - `internal/cli/implement.go`;
   - `internal/worktree/worktree.go`;
   - `internal/cli/deliver.go`;
   - `internal/cli/deliver_workflow.go`;
   - `internal/cli/deliver_linux_lane.go`.
6. MUST add these tests. Each calls `t.Parallel()` first, compiles for
   Windows, uses temporary repositories and reaches no network.
   - `internal/daemon/qa_every_run_audit_test.go`:
     `TestQADeliveryBaseReadsTheFetchedDefaultWithoutOriginHead`. The fixture
     has `refs/remotes/origin/main`, no origin/HEAD and no local `main`, and
     HEAD on a work branch. The base resolves to the merge base with
     `origin/main`.
   - `internal/cli/archive_pull_request_test.go`:
     `TestArchiveChecksTheRevisionAgainstOriginWithoutOriginHead`. It uses
     `archivePullRequestFixture` with `refs/remotes/origin/HEAD` deleted and
     asserts the whole transcript. Its warning is the one
     `TestArchiveOnABranchWarnsAboutASquashDroppedRevision` expects, naming
     `origin/main`, instead of "the default branch is unknown".
   - `internal/cli/review_merge_base_test.go`:
     `TestReviewDefaultBaseIsTheFetchedDefaultBranch`. On a work branch
     without origin/HEAD, local `main` is one commit behind
     `refs/remotes/origin/main`. The review record's `baseTipCommit` is the
     `origin/main` tip, and its `baseCommit` is the merge base with it.

   Every existing test MUST keep its name and pass, including
   `TestQADeliveryBasePrefersTheRemoteTrackingBranch`, the archive pull
   request tests and `TestReviewCommandDefaultsBaseToMain`.
7. MUST document the change:
   - `docs/user-guide/commands/archive.md` and
     `.agents/skills/roundfix/references/archive.md` each say that the
     command warns when `source_revision` "is not an ancestor of the Default
     Branch Head".
   - `docs/user-guide/commands/review.md` and
     `.agents/skills/roundfix/references/review.md` each replace "Roundfix
     uses the repository's default branch" with "Roundfix uses the Default
     Branch Head".

   No text inside any `### QA settlement` section may change.
8. MUST revise **Delivery Base** in `CONTEXT.md` through domain-modeling, with
   ADR-0295's link, so it reads "The merge base of the audited head and the
   Default Branch Head."
9. MUST run these commands, because this Task edits user guide pages and the
   Roundfix references. No record or mirror is hand-edited.
   Before them, raise both version fields of each edited owned skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`)
   by one patch over the value the tree holds when this Task starts,
   because the record command refuses changed content under an unchanged
   version. The Agent runs these commands itself before handing back: the
   Daemon only runs the Verification and never regenerates a mirror or a
   record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

## Subtasks

- [ ] Make the QA Delivery Base, the archive default reference and the review base read the resolver.
- [ ] Add the three tests.
- [ ] Document the warning and the review base in the guides, the references and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Without origin/HEAD, the QA Delivery Base resolves, the archive warning names `origin/main`, and the review base is the `origin/main` tip over a stale local `main`.
- [ ] The callers named unchanged stay byte-identical.
- [ ] The guides, the references and the glossary name the Default Branch Head.

## Consumers

- `preflight.DetectDefaultBranch` and `preflight.ResolveDefaultBranchHead` — consumers: `./internal/daemon`, `./internal/cli`, `./internal/worktree`

## Context

- instruction: `docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md`
- instruction: `docs/adr/0169-the-pre-pr-review-diffs-the-candidate-from-its-merge-base.md`
- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- interface: `internal/daemon/qa_delivery_base.go`
- interface: `internal/daemon/qa_every_run_audit_test.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/archive_pull_request_test.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_merge_base_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestQADeliveryBaseReadsTheFetchedDefaultWithoutOriginHead|TestQADeliveryBasePrefersTheRemoteTrackingBranch|TestQADeliveryBaseIsTheMergeBaseWithTheDefaultBranch)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestQADeliveryBaseReadsTheFetchedDefaultWithoutOriginHead ' || { printf 'missing pass: TestQADeliveryBaseReadsTheFetchedDefaultWithoutOriginHead\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestArchiveChecksTheRevisionAgainstOriginWithoutOriginHead|TestReviewDefaultBaseIsTheFetchedDefaultBranch|TestReviewCommandDefaultsBaseToMain)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveChecksTheRevisionAgainstOriginWithoutOriginHead TestReviewDefaultBaseIsTheFetchedDefaultBranch TestReviewCommandDefaultsBaseToMain; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/archive.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'is not an ancestor of the Default Branch Head' || { printf 'missing archive warning rule in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Roundfix uses the Default Branch Head' || { printf 'missing review base rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'The merge base of the audited head and the Default Branch Head.' || { printf 'missing revised glossary term Delivery Base\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Delivery Base**:' || { printf 'missing glossary term Delivery Base\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/worktree || exit 1; go test -count=1 -run 'QA|Delivery|Mechanical|Settle' ./internal/daemon || exit 1; go test -count=1 -run 'Archive|Review|Deliver|Implement|Preflight|History|Default' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the Delivery Base, archive and review tests, the `worktree`, `daemon` and `cli` suites that consume task_01's detection, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 2; Core Feature 3; Core Feature 5; Success Metric 3; Non-Goals; Glossary
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 6; API Contract 10; Invariant 2; Current behavior; Testing Approach; Build Order 2
- ADR-0295; ADR-0276; ADR-0169; ADR-0261; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285
