---
task: task_02
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: completed
type: backend
complexity: medium
---

# Task 02: The skills lock is read after the fetch

## Overview

`buildSkillsReconcilePlan` in `internal/baseline/skills_reconcile.go` and
`buildSkillsRestorePlan` in `internal/baseline/skills_restore.go` both work in
this order:

1. read `skills-lock.json` with `loadSkillsLock`;
2. acquire the source commit over the network or from `--source-dir`;
3. capture the lock's transaction preimage in
   `buildSkillsReconcileTransactionDocument` or
   `buildRestoreTransactionDocument`.

Another tool, such as `npx skills add`, can rewrite the lock during step 2.
The rewrite then becomes the preimage, and the apply-time preimage check
passes. The result computed from the old bytes overwrites the new edit without
an error.

The lock is written by the maintainer's other tools. It is read by these two
commands, and the Baseline transaction writes the plan's postimage back to it.

## Requirements

1. MUST make `buildSkillsReconcilePlan` acquire the source commit before the
   `skills-lock.json` read that feeds the plan, and make `buildSkillsRestorePlan`
   take that read after its last `acquireRestoreGroup`. An earlier read MAY stay
   only to refuse a malformed or unsafe lock before acquisition. It MUST NOT feed
   the plan digest, the lock edits or the postimage.
2. MUST give `buildSkillsReconcileTransactionDocument` and
   `buildRestoreTransactionDocument` the planned lock bytes (`lock.before`) as
   an argument. Each MUST compare them with the captured `skills-lock.json`
   state: a missing lock matches only a missing capture, and a present lock
   matches only a regular file whose `ContentIdentity` equals
   `transactionContentIdentity(lock.before)`. On a mismatch, each MUST return a
   `SkillsRestoreError` with category `action_required`, code
   `lock.changed-during-plan`, message `skills-lock.json changed during
   planning; the plan no longer describes it.` and action `Rerun the preview
   and confirm its new Plan Digest.`, and write nothing. The restore builder
   checks only when the plan edits the lock. The only callers are the two plan
   builders, and no `internal/**/testdata/*.txt` harness names these functions.
3. MUST keep the apply-time preimage revalidation unchanged, so a rewrite after
   planning still ends in `plan.confirmation.stale`. MUST keep the restore
   payload shape unchanged (ADR-0072).
4. MUST NOT add a production hook for tests. The tests inject the rewrite
   through a `git` script that the test writes into a temporary directory and
   places first on `PATH` with `t.Setenv`, so these tests do not call
   `t.Parallel`. The script copies new lock bytes over `skills-lock.json` the
   first time it runs `init --bare`, then execs the real `git`, resolved with
   `exec.LookPath` before `PATH` changes.
5. MUST add these tests to `internal/baseline/skills_lock_read_test.go`, each
   negative case in its own test:
   - `TestReconcileSkillsLockPlansFromTheLockPresentAfterTheFetch`: the rewrite
     adds a second obsolete entry, and the preview's `plannedChanges` lists
     both entries.
   - `TestReconcileSkillsLockConfirmedApplyKeepsAFetchTimeRewrite`: a Plan
     Digest previewed on the old bytes is confirmed while the script rewrites
     the lock. The apply refuses with `plan.confirmation.stale`, and
     `skills-lock.json` holds exactly the rewritten bytes.
   - `TestSkillsRestorePlansFromTheLockPresentAfterTheFetch`: the rewrite adds
     an unrelated entry, and the planned lock postimage keeps it.
   - `TestSkillsRestoreConfirmedApplyKeepsAFetchTimeRewrite`: the same as the
     reconcile apply test, for restore.
   - `TestReconcileTransactionDocumentRefusesAChangedLock` and
     `TestRestoreTransactionDocumentRefusesAChangedLock`: planned bytes that
     differ from the disk bytes produce `lock.changed-during-plan` and no
     document.
   - `TestReconcileTransactionDocumentAcceptsTheLockItPlanned` and
     `TestRestoreTransactionDocumentAcceptsTheLockItPlanned`: equal bytes
     produce a document whose lock preimage `ContentIdentity` equals the
     planned bytes' identity.
6. MUST add `lock.changed-during-plan` to the exit `3` line of both the
   `baseline skills reconcile` and the `baseline skills restore` help in
   `internal/cli/cli.go`. MUST add
   `TestBaselineSkillsReconcileHelpNamesTheLockChangeRefusal` to
   `internal/cli/baseline_skills_reconcile_test.go` and
   `TestBaselineSkillsRestoreHelpNamesTheLockChangeRefusal` to
   `internal/cli/baseline_skills_restore_test.go`, each asserting that its help
   contains `lock.changed-during-plan`.
7. MUST state in the `baseline skills restore` and `baseline skills reconcile`
   paragraphs of `docs/user-guide/commands.md`, and in the Roundfix skill's
   skill-lock paragraph in `.agents/skills/roundfix/SKILL.md`, two facts: the
   lock is read after the source is acquired, and a lock that changed during
   planning is refused with `lock.changed-during-plan` and exit `3`. Both MUST
   name `lock.changed-during-plan`. MUST regenerate `skills/roundfix/SKILL.md`
   with `make skills-sync`.
8. MUST keep these existing tests green and unchanged:
   `TestReconcileSkillsLockStalePreimageWritesNothing`,
   `TestSkillsRestoreStalePlanDoesNotMutate`,
   `TestSkillsRestoreCompatibilityMatchesMaintainedPythonShape`,
   `TestSkillsRestoreProvenanceAndPreMutationRefusals`,
   `TestSkillsRestoreOfflinePreviewApplyAndIdempotence` and
   `TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A lock rewritten during the fetch is the lock the plan describes, for
      reconcile and restore.
- [ ] A confirmed apply of a digest previewed on the old bytes refuses and
      leaves the rewritten bytes.
- [ ] A transaction builder refuses a lock that differs from the planned bytes
      and accepts the lock it planned.
- [ ] The help, the user guide and the Roundfix skill name
      `lock.changed-during-plan`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/skills_reconcile.go`
- interface: `internal/baseline/skills_restore.go`
- creates: `internal/baseline/skills_lock_read_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/baseline_skills_reconcile_test.go`
- interface: `internal/cli/baseline_skills_restore_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileSkillsLockPlansFromTheLockPresentAfterTheFetch|TestReconcileSkillsLockConfirmedApplyKeepsAFetchTimeRewrite|TestSkillsRestorePlansFromTheLockPresentAfterTheFetch|TestSkillsRestoreConfirmedApplyKeepsAFetchTimeRewrite|TestReconcileTransactionDocumentRefusesAChangedLock|TestRestoreTransactionDocumentRefusesAChangedLock|TestReconcileTransactionDocumentAcceptsTheLockItPlanned|TestRestoreTransactionDocumentAcceptsTheLockItPlanned|TestReconcileSkillsLockStalePreimageWritesNothing|TestSkillsRestoreStalePlanDoesNotMutate|TestSkillsRestoreCompatibilityMatchesMaintainedPythonShape|TestSkillsRestoreProvenanceAndPreMutationRefusals|TestSkillsRestoreOfflinePreviewApplyAndIdempotence|TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries|TestBaselineSkillsReconcileHelpNamesTheLockChangeRefusal|TestBaselineSkillsRestoreHelpNamesTheLockChangeRefusal)$" ./internal/baseline ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileSkillsLockPlansFromTheLockPresentAfterTheFetch TestReconcileSkillsLockConfirmedApplyKeepsAFetchTimeRewrite TestSkillsRestorePlansFromTheLockPresentAfterTheFetch TestSkillsRestoreConfirmedApplyKeepsAFetchTimeRewrite TestReconcileTransactionDocumentRefusesAChangedLock TestRestoreTransactionDocumentRefusesAChangedLock TestReconcileTransactionDocumentAcceptsTheLockItPlanned TestRestoreTransactionDocumentAcceptsTheLockItPlanned TestReconcileSkillsLockStalePreimageWritesNothing TestSkillsRestoreStalePlanDoesNotMutate TestSkillsRestoreCompatibilityMatchesMaintainedPythonShape TestSkillsRestoreProvenanceAndPreMutationRefusals TestSkillsRestoreOfflinePreviewApplyAndIdempotence TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries TestBaselineSkillsReconcileHelpNamesTheLockChangeRefusal TestBaselineSkillsRestoreHelpNamesTheLockChangeRefusal; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "lock.changed-during-plan" docs/user-guide/commands.md && grep -q "lock.changed-during-plan" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide names `lock.changed-during-plan`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The lock read after the fetch

## Result

Implementation:

- Reconcile now acquires the immutable source commit before reading the lock
  used for classification, lock edits, postimages and the Plan Digest. Restore
  reads that lock after its final source group acquisition.
- Both transaction-document builders receive the planned lock bytes and refuse
  a missing, non-regular or content-identity-mismatched capture with
  `lock.changed-during-plan`. The restore builder performs this comparison only
  when the plan edits `skills-lock.json`; apply-time preimage validation is
  unchanged.
- Added the eight required lock-read and transaction-document regression tests
  using a temporary `git` wrapper that rewrites `skills-lock.json` on the first
  `init --bare` and then delegates to the pre-resolved real Git executable.
- Both command helps, both user-guide command paragraphs and the canonical
  Roundfix skill document the post-acquisition read and exit `3`
  `lock.changed-during-plan` refusal. `make skills-sync` regenerated the
  distributed Roundfix skill.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '<eight new Task 02 baseline tests>' ./internal/baseline`: passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '<six named existing compatibility tests>' ./internal/baseline`: passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestBaselineSkillsReconcileHelpNamesTheLockChangeRefusal|TestBaselineSkillsRestoreHelpNamesTheLockChangeRefusal)$' ./internal/cli`: passed.
- `rtk make skills-sync`: passed; canonical and distributed Roundfix skill
  copies are byte-identical.
- `rtk git diff --check`: passed.

Acceptance evidence:

- The post-fetch planning tests passed for reconcile and restore; they prove
  that the rewritten lock drives both obsolete-entry classification and the
  planned restore lock postimage.
- Both confirmed-apply tests passed with `plan.confirmation.stale` and exact
  preservation of the fetch-time rewritten lock bytes.
- Both transaction builders passed their separate refusal and acceptance
  tests, including exact `ContentIdentity` comparison to the planned bytes and
  no document on mismatch.
- Both public help tests passed, and the user guide plus synchronized Roundfix
  skill name `lock.changed-during-plan`, the post-acquisition read and exit
  `3` behavior.

The authored `## Verification` command was not rerun; the Daemon owns that
verification and Task settlement.
