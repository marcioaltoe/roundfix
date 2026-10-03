---
task: task_02
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: completed
type: backend
complexity: high
---

# Task 02: A retry after a post-archive correction and a new queue keep the item's work

## Overview

On 2026-10-01 the Delivery Retry of Spec 0205 refused with `archived item head
… differs from candidate head …` after the operator committed a correction on
top of the archive, and the operator opened Pull Request #329 by hand. On
2026-10-01 for Spec 0205 and on 2026-10-03 for Spec 0217, `deliver start`
created a new item from the default branch although the old item branch held
completed Tasks, and the operator merged the old branch by hand. This Task
sends a descendant archived head back to review and makes `deliver start`
continue the one item branch that holds work, or refuse naming every such
branch.

## Requirements

1. MUST answer items 3 and 4 of the Backlog Entry of 2026-10-01, "Archiving or
   requeueing a Spec loses or breaks its delivery", as `_techspec.md` → The
   archived retry and The continued item branch state.
2. MUST change the archived branch of `Engine.Retry` so that, for every
   blocker, a head that descends from the item's newest candidate head is
   appended to the candidate commits and the item returns to `reviewing`
   (invariant 4); a `qa-environment-partial` blocker MUST still require the
   override record and alone MAY anchor on the Run start head; a head that
   does not descend, or an engine without item history, MUST still refuse
   with today's reason; the
   `corrective-spec-required` branch MUST stay unchanged (invariant 5).
3. MUST add `existingItemBranches` to `internal/cli/deliver_workflow.go` as
   `_techspec.md` → Interfaces and The continued item branch state, reading
   only local refs (invariant 6).
4. MUST make `roundfix deliver start` call it for each slug after the
   prerequisite check and before the authorization and readiness checks, and
   refuse with exit `2` and the reason of Surface Transcript 2, recording no
   queue, when a slug has two or more such branches; with exactly one it MUST
   print `Continuing item branch <branch> for <slug>` on stdout after the
   queue is recorded.
5. MUST make `CreateItemBranch` record that single branch for the new item
   instead of minting a name, so the item reuses the branch and its worktree;
   a slug with no such branch MUST still get a new branch from the default
   branch.
6. MUST keep the existing retry, operator-archive, conflict and corrective-Spec
   tests passing unedited.
7. MUST add the tests named in Verification to the new files
   `internal/delivery/archived_head_retry_test.go` and
   `internal/cli/deliver_item_branch_test.go`; the CLI tests build a bare
   remote and a clone in temporary directories, use a temporary Roundfix Home,
   and never reach GitHub.
8. MUST describe the retry after a post-archive correction and the continued
   item branch in `docs/user-guide/commands/deliver.md` and the Roundfix
   Skill's `deliver` reference; MUST raise the Roundfix Skill's version by one
   patch level above the version on this Task's starting tree in both
   front-matter fields, run `make skills-sync`, and re-record the version.

## Subtasks

- [ ] Send a descendant archived head back to review.
- [ ] List a Spec's item branches with work.
- [ ] Refuse or continue at `deliver start`, and reuse the branch at item creation.
- [ ] Add the tests.
- [ ] Describe it in the guide and the Skill, and raise its version.

## Acceptance Criteria

- [ ] A retry of an archived `gate-failed` item whose head is one commit past
      its candidate returns `reviewing` with that head as the newest
      candidate; a non-descendant head and a `corrective-spec-required` park
      still refuse, and the item is unchanged after each refusal.
- [ ] `deliver start` with two item branches ahead of `origin/main` prints
      Surface Transcript 2's reason, exits `2`, and records no queue.
- [ ] With one item branch two commits ahead, the new item records that branch
      and the start prints the continuation line; with one item branch that
      has no commit beyond `origin/main`, a new branch is minted.
- [ ] The guide and the Skill describe both behaviors, the mirrors equal their
      canonical files, and the raised version is recorded.

## Context

- creates: `internal/delivery/archived_head_retry_test.go`
- creates: `internal/cli/deliver_item_branch_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/delivery/operator_archive_retry_test.go`
- instruction: `internal/delivery/corrective_spec_test.go`
- instruction: `internal/store/delivery.go`
- instruction: `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/references/2026-10-01-archiving-or-requeueing-a-spec-loses-or-breaks-its-delivery.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryOfAnArchivedItemWithACorrectionReturnsToReview|TestRetryOfAnArchivedItemRefusesAHeadThatDoesNotDescend|TestRetryOfACorrectiveSpecParkStillRefusesAMovedHead|TestRetryRefusesAnOperatorArchiveWithoutAQAOverride|TestRetryResumesAnOperatorArchivedItemAtReview|TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved|TestRetryRefusesAnArchivedItemWhoseHeadMoved)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryOfAnArchivedItemWithACorrectionReturnsToReview TestRetryOfAnArchivedItemRefusesAHeadThatDoesNotDescend TestRetryOfACorrectiveSpecParkStillRefusesAMovedHead TestRetryRefusesAnOperatorArchiveWithoutAQAOverride TestRetryResumesAnOperatorArchivedItemAtReview TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved TestRetryRefusesAnArchivedItemWhoseHeadMoved; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the four existing retry tests run unedited beside the three new ones, which do not exist before this Task.
- `out="$(go test -count=1 -v -run "^(TestDeliverStartContinuesTheItemBranchThatHoldsWork|TestDeliverStartRefusesTwoItemBranchesWithWork|TestDeliverStartIgnoresAnItemBranchWithoutWork)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverStartContinuesTheItemBranchThatHoldsWork TestDeliverStartRefusesTwoItemBranchesWithWork TestDeliverStartIgnoresAnItemBranchWithoutWork; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `for file in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "Continuing item branch" || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "item branches with commits" || { printf 'missing refusal in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && out="$(go test -count=1 -v -run "^TestEveryOwnedSkillVersionIsRecorded$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the guide nor the reference names the continuation, so the command fails.

## References

- `_prd.md` → Goals; User Stories 3-4; Core Features 4-5, 9; Success Metrics 3-4; Acceptance evidence
- `_techspec.md` → The archived retry; The continued item branch; Interfaces; API Contract 3; API Contract 4; Surface Transcript 2; Testing Approach 2; Build Order 2
- ADR-0223; ADR-0165; ADR-0170; ADR-0154

## Result

Implemented this Task's slice; Task status and declared Verification remain
Daemon-owned.

### Acceptance evidence

1. Archived retry now accepts a moved head proven to descend from the newest
   candidate for every applicable blocker and appends it before returning to
   `reviewing`. The QA environment blocker still requires its override for a
   moved head and alone can fall back to Run start. The corrective-Spec branch
   is unchanged. `TestRetryOfAnArchivedItemWithACorrectionReturnsToReview`
   checks the newest of two candidates and persisted stage/candidate history;
   the two new refusal tests check unchanged persisted items, including absent
   and unreadable history. Existing operator archive, retry, corrective-Spec
   and conflict tests were left unedited.
2. `existingItemBranches` reads local refs, restricts names to the Spec's
   16-character lowercase hex item suffix, counts commits absent from the
   configured default remote-tracking ref and sorts results. Start checks it
   after prerequisites and before authorization/readiness. The two-branch
   command test asserts exit 2, the exact sorted transcript reason, empty
   stdout and no Run Database; an unavailable remote proves no fetch, and
   readiness/owner callbacks fail if reached.
3. Start prints the continuation after queue recording. Creation rechecks
   after fetch and records the single existing branch, using the existing
   worktree provisioning path. The continuation test covers two commits,
   both a missing worktree and an existing worktree with a preserved local
   note; invalid-suffix and other-Spec branches do not create ambiguity. The
   no-work test checks a new branch at origin/main with no continuation.
   Recorded workspace recovery retains its no-fetch path.
4. Both the command guide and canonical deliver reference describe branch
   selection, refusal and post-archive correction review, including the QA
   override and corrective-Spec restrictions. Both Roundfix Skill version
   fields increased from the starting 0.1.19 to 0.1.20. `make skills-sync`
   regenerated the mirrors; the version recorder generated the final 0.1.20
   digest. A Python content/byte comparison checked the phrases, both mirrors
   and both version fields. Sanctioned baseline digest regeneration left no
   final derived artifact changes.

### Focused checks

Commands below used `GOCACHE=/private/tmp/roundfix-task02-gocache` after the
sandbox denied access to an entry in the host Go cache.

- `go test ./internal/delivery ./internal/cli -run 'TestRetryOf|TestDeliverStartContinues|TestDeliverStartRefusesTwo|TestDeliverStartIgnores' -count=1`
  — exit 0; all six new tests passed, including their subtests.
- `go test ./internal/delivery ./internal/cli -run 'TestRetry|TestOperatorArchive|TestTheArchiveStage|TestDeliveryEngineParks|Test.*Conflict|Test.*Item|Test.*DeliveryGets|TestDeliverStart' -count=1`
  — exit 0; existing recovery and item-worktree behavior also passed.
- A temporary Go overlay substituted the three original HEAD production
  files while retaining the new tests. The correction, continuation and
  ambiguity tests failed on that starting behavior (exit 1), providing a
  regression signal without changing the worktree.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; recorded the finalized 0.1.20 content. The initial record was
  regenerated from the committed registry after a final formatting edit.
- `make baseline-digests` — exit 0; final artifacts already match canonical
  sources, `changed: false`.
- `git -c core.fsmonitor=false diff --check` — exit 0.
- `make verify-incremental` — exit 0 with host permissions, on the stable
  final implementation tree; formatting, vet, all package tests, skill checks
  and build passed. The first sandbox attempt exited 2 because process-table
  access and localhost test sockets were denied; documentation edits that
  overlapped that attempt also triggered the repository-change guard. The
  permission-enabled rerun kept sources stable and passed those checks.

The declared Verification commands were not run. No other Task file or Task
Graph was edited, and no commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261003T213431Z_8d253be688342055`
- Source commit: `a9d2d8d298412572e0ebf3b6670c3f2b88f7c669`
