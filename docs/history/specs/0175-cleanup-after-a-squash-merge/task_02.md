---
task: task_02
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: medium
---

# Task 02: Key legacy Runs from their Run Worktree

## Overview

`applyRepositoryRootMigration` in `internal/store/store.go` backfills
`repository_root` only from `git_root`, and only while migrating an older
schema. A Run created from a linked worktree that was later removed keeps an
empty key. `ListRuns` therefore never returns it from the main checkout, so
`roundfix runs list` and `roundfix reconcile` miss it. The prune guard in
`PruneTerminalReport` (`internal/worktree/worktree.go`) skips it because it
compares `run.GitRoot` with the checkout by path. On 2026-09-25, 26 of 108 Runs
on the maintainer's machine were unkeyed, and 9 of them still had registered
Run Worktrees. The key is derived from files under a recorded Run Worktree
path, which the Run Database stored when the Run started. It is written into
the Run Database, and every repository-scoped listing and cleanup reads it. A
wrong key would let one repository's cleanup act on another repository's Run.

## Requirements

1. MUST run the repository-key backfill on every write-mode `store.Open`, also
   when the schema is already current, inside the migration's write
   transaction, and never on `store.OpenReader`. The schema version MUST stay
   20.
2. MUST keep today's derivation from `git_root` when `git_root/.git` exists.
   Otherwise, when `work_dir` is set and `work_dir/.git` exists, it MUST derive
   the key with `roundconfig.RepositoryRoot(work_dir)` and accept it only when
   it differs from the cleaned `work_dir`. Every other Run stays unkeyed.
3. MUST update only rows whose `repository_root` is still empty, so a second
   open changes nothing.
4. MUST make `PruneTerminalReport` match a terminal Run by repository key:
   `run.RepositoryRoot`, or `roundconfig.RepositoryRoot(run.GitRoot)` when it
   is empty, against `roundconfig.RepositoryRoot(userRoot)`. A matching Run
   MUST be inspected with `GitRoot` set to `userRoot`, and a Run keyed to
   another repository MUST be skipped.
5. MUST keep `TestJournalConsumerCorpusReplaysEveryConsumer` green. That
   corpus is opened through `store.Open`, so the backfill runs over it.
6. MUST put new tests in `internal/store/repository_key_backfill_test.go`,
   `internal/worktree/prune_repository_key_test.go` and
   `internal/cli/reconcile_legacy_key_test.go`, and MUST NOT remove or rename
   an existing test.

## Subtasks

- [ ] Implement the re-runnable backfill and the key-based prune guard.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A write-mode open keys an unkeyed Run whose checkout was removed from its registered Run Worktree, and `ListRuns` for the main checkout returns it.
- [ ] An unkeyed Run whose Run Worktree is gone, one whose worktree has no administrative entry, and one whose worktree belongs to another repository are each left unkeyed or keyed to their own repository, never to the caller's.
- [ ] A second open changes nothing.
- [ ] The implement-preflight prune releases a safe keyed Run whose checkout was removed and skips a Run keyed to another repository.
- [ ] `roundfix reconcile` from the main checkout reports the keyed legacy Run.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/worktree/worktree.go`
- creates: `internal/store/repository_key_backfill_test.go`
- creates: `internal/worktree/prune_repository_key_test.go`
- creates: `internal/cli/reconcile_legacy_key_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOpenKeysALegacyRunFromItsRunWorktree|TestOpenLeavesALegacyRunUnkeyedWhenItsWorktreeIsGone|TestOpenRejectsAWorktreeKeyWithoutACommonDirectory|TestOpenKeysALegacyRunToItsOwnRepository|TestOpenRepositoryKeyBackfillIsIdempotent|TestOpenReaderDoesNotBackfillRepositoryKeys|TestJournalConsumerCorpusReplaysEveryConsumer|TestPruneTerminalReportReleasesARunWhoseCheckoutWasRemoved|TestPruneTerminalReportSkipsARunKeyedToAnotherRepository|TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|TestPruneTerminalReportRequiresArchivedEvidence|TestPruneTerminalReconciliationPreservesUniqueChangedBranch|TestReconcileReportsALegacyRunKeyedFromItsRunWorktree|TestReconcileFromMainAfterALinkedWorktreeIsRemoved|TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath)$" ./internal/store ./internal/worktree ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOpenKeysALegacyRunFromItsRunWorktree TestOpenLeavesALegacyRunUnkeyedWhenItsWorktreeIsGone TestOpenRejectsAWorktreeKeyWithoutACommonDirectory TestOpenKeysALegacyRunToItsOwnRepository TestOpenRepositoryKeyBackfillIsIdempotent TestOpenReaderDoesNotBackfillRepositoryKeys TestJournalConsumerCorpusReplaysEveryConsumer TestPruneTerminalReportReleasesARunWhoseCheckoutWasRemoved TestPruneTerminalReportSkipsARunKeyedToAnotherRepository TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches TestPruneTerminalReportRequiresArchivedEvidence TestPruneTerminalReconciliationPreservesUniqueChangedBranch TestReconcileReportsALegacyRunKeyedFromItsRunWorktree TestReconcileFromMainAfterALinkedWorktreeIsRemoved TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Repository keys for legacy Runs
- `_prd.md` → Core Feature 4; Success Metric 4
- `_techspec.md` → API Contract 5

## Result

Implemented a re-runnable repository-key repair inside the existing migration
write transaction. A write-mode open now examines every still-unkeyed Run,
prefers surviving `git_root` metadata, falls back to a registered `work_dir`,
and leaves absent or unproven worktrees untouched. The schema remains version
20, readers remain read-only, and a repeated open performs no row update.

`PruneTerminalReport` now compares durable repository keys. It inspects a
matching Run through the current canonical checkout and skips a Run whose
recorded key belongs to another repository.

Focused evidence by acceptance criterion:

- `TestOpenKeysALegacyRunFromItsRunWorktree` passed and observed the repaired
  Run through `ListRuns` from the main checkout.
- `TestOpenLeavesALegacyRunUnkeyedWhenItsWorktreeIsGone`,
  `TestOpenRejectsAWorktreeKeyWithoutACommonDirectory`, and
  `TestOpenKeysALegacyRunToItsOwnRepository` passed as separate negative
  cases. The foreign Run was keyed to its owning repository and excluded from
  the original repository's listing.
- `TestOpenRepositoryKeyBackfillIsIdempotent` passed with one row change on
  the repairing open and zero on the second open.
  `TestOpenReaderDoesNotBackfillRepositoryKeys` also passed.
- `TestPruneTerminalReportReleasesARunWhoseCheckoutWasRemoved` and
  `TestPruneTerminalReportSkipsARunKeyedToAnotherRepository` passed. The
  neighboring empty-branch, archived-evidence, and unique-change prune tests
  also passed.
- `TestReconcileReportsALegacyRunKeyedFromItsRunWorktree` passed through the
  public reconcile runner. The existing removed-linked-worktree and
  recorded-key-over-checkout-path reconcile tests also passed.
- `TestJournalConsumerCorpusReplaysEveryConsumer` passed with the backfill
  active during `store.Open`.

Focused commands:

- `GOCACHE=/private/tmp/roundfix-task-02-gocache go test -count=1 -run '^(TestOpenKeysALegacyRunFromItsRunWorktree|TestOpenLeavesALegacyRunUnkeyedWhenItsWorktreeIsGone|TestOpenRejectsAWorktreeKeyWithoutACommonDirectory|TestOpenKeysALegacyRunToItsOwnRepository|TestOpenRepositoryKeyBackfillIsIdempotent|TestOpenReaderDoesNotBackfillRepositoryKeys|TestJournalConsumerCorpusReplaysEveryConsumer)$' ./internal/store` — passed.
- `GOCACHE=/private/tmp/roundfix-task-02-gocache go test -count=1 -run '^(TestPruneTerminalReportReleasesARunWhoseCheckoutWasRemoved|TestPruneTerminalReportSkipsARunKeyedToAnotherRepository|TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches|TestPruneTerminalReportRequiresArchivedEvidence|TestPruneTerminalReconciliationPreservesUniqueChangedBranch)$' ./internal/worktree` — passed.
- `GOCACHE=/private/tmp/roundfix-task-02-gocache go test -count=1 -run '^(TestReconcileReportsALegacyRunKeyedFromItsRunWorktree|TestReconcileFromMainAfterALinkedWorktreeIsRemoved|TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath)$' ./internal/cli` — passed.
- `GOCACHE=/private/tmp/roundfix-task-02-gocache make verify-incremental` —
  the sandboxed attempt was blocked only when two existing force-stop tests
  could not read the process table; the unchanged command passed with the
  required host permission.

The Daemon-owned `## Verification` command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260928T154312Z_63679c6540b34603`
- Source commit: `a2cc95ed6999d84f8d4481ed4def96bff8eefe64`
