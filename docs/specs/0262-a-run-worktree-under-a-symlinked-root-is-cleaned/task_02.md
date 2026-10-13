---
task: task_02
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: completed
type: backend
complexity: medium
---

# Task 02: Worktree paths compare through their longest existing prefix

## Overview

`samePath` compares `canonicalPath` of both sides, and `canonicalPath` falls
back to the lexical path when the directory is gone. Under a symbolic-linked
worktree location, an item worktree whose directory was removed no longer
matches the resolved path Git keeps, so `CleanupItem` refuses and `UseItem`
leaves the stale registration (`_techspec.md` → Current behavior, the sweep).
This Task makes `canonicalPath` resolve the longest existing prefix by reusing
the review-root helper of `internal/config`. It is verifiable on its own in
`internal/worktree` and `internal/config`.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [a Run Worktree under a symlinked worktree root is never cleaned](references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md),
   for every comparison of a recorded and a registered worktree path, by
   implementing `_techspec.md` → API Contract 2.
2. MUST rename `canonicalizeResolvedPath` in `internal/config/config.go` to the
   exported `CanonicalizeResolvedPath`, keep its behavior, update its only
   caller, and make `canonicalPath` in `internal/worktree/worktree.go`
   delegate to it. No second copy of the algorithm may be written.
3. MUST add, in `internal/worktree/worktree_path_test.go`, each test calling
   `t.Parallel()` first:
   - `TestCanonicalPathResolvesTheExistingPrefix` with subtests `existing`,
     `missing-leaf`, `missing-parent` and `no-existing-prefix`; the first three
     use a directory reached through a symbolic link and expect the resolved
     prefix with the missing components kept verbatim.
   - `TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot`: an item
     worktree created with `CreateItem` under a linked location, its directory
     removed from disk, is cleaned up: the registration is pruned and the item
     branch deleted, with no error.
   - `TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot`: the same gone
     item worktree is recreated at its recorded path on its branch.
4. MUST run the whole `internal/worktree` and `internal/config` packages. A
   test that compared a missing path with its unresolved spelling by accident
   may be corrected only to the resolved value, and the Task MUST name each one
   it corrects in its Result; no assertion may be weakened.
5. MUST run the new tests once with `-race` (ADR-0259).

## Subtasks

- [ ] Export the config helper and delegate `canonicalPath` to it.
- [ ] Add the comparison and item worktree tests.

## Acceptance Criteria

- [ ] `canonicalPath` of an existing path is unchanged, and of a missing one keeps the resolved prefix.
- [ ] A gone item worktree under a linked location is cleaned up and recreated.

## Consumers

- Worktree path comparison through samePath for item worktree use and cleanup and item branch inspection — consumers: `./internal/cli`
- Review root containment through CanonicalizeResolvedPath — consumers: `./internal/config`
- checkpoint: task_03

## Context

- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `internal/worktree/item_branch_reconcile.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/config/config.go`
- creates: `internal/worktree/worktree_path_test.go`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestCanonicalPathResolvesTheExistingPrefix|TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot|TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCanonicalPathResolvesTheExistingPrefix TestCanonicalPathResolvesTheExistingPrefix/existing TestCanonicalPathResolvesTheExistingPrefix/missing-leaf TestCanonicalPathResolvesTheExistingPrefix/missing-parent TestCanonicalPathResolvesTheExistingPrefix/no-existing-prefix TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three tests do not exist, so the command fails.
- `grep -q 'func CanonicalizeResolvedPath(' internal/config/config.go || { printf 'the config helper is not exported\n' >&2; exit 1; }; go test -count=1 ./internal/worktree ./internal/config || exit 1` — expected: exit 0. Before this Task the helper is unexported, so the command fails; afterwards both whole packages pass.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 3; Success Metric 3
- `_techspec.md` → Current behavior; API Contract 2; Invariant 3; Testing Approach; Build Order 2; Risks & Considerations
- ADR-0268; ADR-0259

## Result

Implemented API Contract 2 by exporting the existing config helper as
`CanonicalizeResolvedPath`, updating both calls in its sole config consumer,
`reviewRootInsideHistory`, and delegating worktree `canonicalPath` to it. The
config algorithm is unchanged; no second implementation was introduced.
Every `samePath` consumer now uses the resolved existing prefix when a
worktree directory is missing.

Acceptance evidence:

- Existing paths retain their resolved value, and missing paths retain the
  resolved prefix: `TestCanonicalPathResolvesTheExistingPrefix` passes with
  `existing`, `missing-leaf`, `missing-parent` and `no-existing-prefix`. Expected
  resolved paths come from `filepath.EvalSymlinks` of the existing fixture,
  independently of `canonicalPath`; missing component case and spaces are
  preserved. The relative case has no existing named prefix and keeps its
  clean spelling.
- Gone item worktrees under a linked location are cleaned up and recreated:
  `TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot` confirms that
  cleanup removes the stale Git registration and item branch without error.
  `TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot` confirms the
  recorded path exists again, has the recorded branch and its committed HEAD
  and contents, and has exactly one Git registration. Both create the item
  with `CreateItem`, confirm Git's resolved registration before removing the
  directory, and then call the public cleanup or reuse operation. Every new
  top-level test and comparison subtest calls `t.Parallel()` first.

Focused implementation checks:

- Before the production change:
  `rtk proxy go test -count=1 -run 'Test(CanonicalPathResolvesTheExistingPrefix|CleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot|UseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot)' ./internal/worktree`
  exited 1. Missing-leaf and missing-parent retained the unresolved spelling;
  cleanup refused the registered path; reuse failed because the gone worktree
  was still registered.
- After the production change:
  `rtk proxy go test -race -count=1 -v -run 'Test(CanonicalPathResolvesTheExistingPrefix|CleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot|UseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot)' ./internal/worktree`
  exited 0; all three tests and all four comparison subtests passed with no
  race report.
- `rtk proxy go test -shuffle=on -count=1 ./internal/worktree ./internal/config`
  exited 0 for both whole packages. No existing test or assertion was changed.
- `rtk proxy git diff --check` exited 0.

The execution sandbox could not start with the symlinked workspace root;
focused checks and scoped file writes used elevated execution. No tooling
configuration changed. The Task's declared Verification commands and the
repository gates remain for the Daemon. Task status and authored checkboxes
were preserved, and no other Task or Task Graph file was edited. No commit,
push or pull request was created. Consumer checks remain in task_03.
