---
task: task_01
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: completed
type: backend
complexity: high
---

# Task 01: A recorded worktree path is resolved within its Worktree Root

## Overview

`recordedWorktree` refuses any recorded path that crosses a symbolic link, so
every terminal Run under a symbolic-linked worktree location is `unknown` with
`recorded Run Worktree path is unsafe` and nothing reclaims it
(`_techspec.md` → Current behavior). This Task replaces that refusal with the
Worktree Root rule: links are followed only above the repository directory,
which, with the worktree directory, must be real. It is verifiable on its own
in `internal/worktree`: a Run under a linked location reconciles, applies and
prunes, while a link at the repository or worktree directory stays unsafe.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [a Run Worktree under a symlinked worktree root is never cleaned](references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md),
   by implementing `_techspec.md` → API Contract 1 in `recordedWorktree`
   through an unexported helper, `resolveWithinWorktreeRoot`. The reason
   strings, the reconciliation `Path` and every later proof MUST stay as they
   are. `pathContainsSymlink` MUST keep refusing a recorded Git root.
2. MUST add, in `internal/worktree/worktree_root_test.go`, the five tests that
   `_techspec.md` → Testing Approach names for task_01, each calling
   `t.Parallel()` first and building its worktree location as a symbolic link
   to a real temporary directory with `Create`:
   - `TestInspectTerminalRunResolvesASymlinkedWorktreeRoot`: a clean Run whose
     branch is contained in its target is `safe`, its `Path` is the recorded
     spelling, and its reason is not `recorded Run Worktree path is unsafe`.
   - `TestApplyTerminalRunReleasesAWorktreeUnderASymlinkedRoot`: `ApplyTerminalRun`
     removes the Run Worktree and deletes the Run Branch.
   - `TestPruneTerminalReportPrunesUnderASymlinkedRoot`: with the linked
     location passed as `location`, the Run Worktree and a kept Task Worktree
     of that Run are pruned and reported.
   - `TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot`, with
     subtests `repository-directory-link`, `worktree-directory-link` and
     `dangling-worktree-link`: each recorded path is `unknown` with reason
     `recorded Run Worktree path is unsafe`, and the real worktree and Run
     Branch still exist afterwards.
   - `TestInspectTerminalRunTreatsAMissingWorktreeUnderASymlinkedRootAsAbsent`:
     after the Run Worktree is removed with Git and the Run Branch deleted, the
     Run is `released`.
3. MUST keep `TestInspectTerminalRunUnsafePath` and
   `TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot` unchanged
   and passing; they pin the refusals that stay.
4. MUST add **Worktree Root** to `CONTEXT.md` through `domain-modeling`, next
   to **Run Worktree**, linked to ADR-0268: the directory two levels above a
   recorded worktree path, the `worktree.location` it was created under, inside
   which symbolic links are followed; `Avoid` synonyms such as worktree
   directory or Roundfix Home.
5. MUST run the new tests once with `-race` (ADR-0259) and MUST NOT change any
   other test, the reconcile report, a flag or an exit code.

## Subtasks

- [ ] Add the Worktree Root resolution helper and use it in `recordedWorktree`.
- [ ] Add the five tests on a linked worktree location.
- [ ] Add the glossary term.

## Acceptance Criteria

- [ ] A terminal Run under a linked location is classified, applied and pruned as under a real directory.
- [ ] A link at the repository or worktree directory stays unsafe and removes nothing.
- [ ] The existing unsafe-path and symlinked Git root tests pass unchanged.

## Consumers

- Run Worktree Reconciliation of a recorded path through InspectTerminalRun, InspectTerminalRunMerged, ApplyTerminalRun and PruneTerminalReport — consumers: `./internal/cli`, `./internal/daemon`
- checkpoint: task_03

## Context

- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/worktree/retained_vanished_checkout_test.go`
- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/worktree_root_test.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestInspectTerminalRunResolvesASymlinkedWorktreeRoot|TestApplyTerminalRunReleasesAWorktreeUnderASymlinkedRoot|TestPruneTerminalReportPrunesUnderASymlinkedRoot|TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot|TestInspectTerminalRunTreatsAMissingWorktreeUnderASymlinkedRootAsAbsent|TestInspectTerminalRunUnsafePath|TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestInspectTerminalRunResolvesASymlinkedWorktreeRoot TestApplyTerminalRunReleasesAWorktreeUnderASymlinkedRoot TestPruneTerminalReportPrunesUnderASymlinkedRoot TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot/repository-directory-link TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot/worktree-directory-link TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot/dangling-worktree-link TestInspectTerminalRunTreatsAMissingWorktreeUnderASymlinkedRootAsAbsent TestInspectTerminalRunUnsafePath TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the five new tests do not exist, so their PASS lines are missing and the command fails; afterwards they and the two pinned refusal tests pass under the race detector.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Worktree Root**:' || { printf 'CONTEXT.md does not define **Worktree Root**\n' >&2; exit 1; }; grep -q 'func resolveWithinWorktreeRoot(' internal/worktree/worktree.go || { printf 'the Worktree Root helper is missing\n' >&2; exit 1; }; go test -count=1 ./internal/worktree || exit 1` — expected: exit 0. Before this Task the term and the helper are absent, so the command fails; afterwards the whole worktree package passes.

## References

- `_prd.md` → Goals; User Story 1; User Story 3; Core Feature 1; Core Feature 2; Core Feature 6; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contract 1; Invariants 1-2; Invariant 4; Testing Approach; Build Order 1
- ADR-0268; ADR-0053; ADR-0259; ADR-0244

## Result

Implemented the assigned Worktree Root slice for Daemon Verification.
`recordedWorktree` now calls `resolveWithinWorktreeRoot` after its existing
path-format checks. The helper uses `Lstat` to require real repository and
worktree directories, treats missing directories as absent, resolves links
at or above the Worktree Root, and requires the resolved worktree path to
match the resolved root joined with the original two directory names.
Registration matching then proceeds as before. The recorded reconciliation
`Path`, reason constants, release proofs and recorded Git root validation
are unchanged; `pathContainsSymlink` still guards that Git root.

Added the five authored public API tests in
`internal/worktree/worktree_root_test.go`. Each top-level test and each
unsafe-link subtest starts with `t.Parallel()`. The shared fixture calls
`Create` with a temporary location linked to a real temporary directory and
records `ref.Path` without resolving its spelling. The pruning test retains a
Task Worktree and asserts both exact reported refs and removal of both
worktrees and branches. The unsafe-link tests also exercise refused apply
and preservation by pruning. Added **Worktree Root** beside **Run Worktree**
in `CONTEXT.md`, linked to ADR-0268 with the required avoided synonyms.

### Focused checks

- Before the implementation,
  `go test -count=1 -run '^TestInspectTerminalRunResolvesASymlinkedWorktreeRoot$' ./internal/worktree`
  exited 1: the new test observed `unknown` with
  `recorded Run Worktree path is unsafe` instead of `safe` for a merged clean
  Run under the linked location.
- After the implementation,
  `go test -count=1 -race -v -run '^(Test.*Symlinked.*|TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot|TestInspectTerminalRunUnsafePath)$' ./internal/worktree`
  exited 0. All five new tests, all three unsafe-link subtests and the two
  pinned refusal tests passed under the race detector.
- After strengthening the pruning-report assertion,
  `go test -count=1 -race -shuffle=on -v -run '^(Test.*Symlinked.*|TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot|TestInspectTerminalRunUnsafePath)$' ./internal/worktree`
  exited 0 on the final source. Shuffle seed: `1791903057831861000`.
  The same seven top-level tests and their subtests passed.

### Acceptance evidence

1. A terminal Run under a linked location is classified, applied and pruned
   as under a real directory: the new inspection, apply and prune tests
   passed, preserving the recorded report path and removing the Run and kept
   Task Worktrees and their branches. The missing-worktree test additionally
   proved `released` after Git removed the Run Worktree and branch.
2. A repository or worktree directory link stays unsafe and removes nothing:
   `repository-directory-link`, `worktree-directory-link` and
   `dangling-worktree-link` passed with `unknown` and the existing unsafe-path
   reason; apply refused, pruning returned no refs, and the real worktree,
   Run Branch and recorded path remained.
3. Existing refusals remain: `TestInspectTerminalRunUnsafePath` and
   `TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot` passed
   unchanged in both race runs. Their source files were not edited.

The executor cannot start its sandbox with this workspace's symbolic-linked
writable root, so file operations and focused checks used approved execution
outside that sandbox. Tests used isolated temporary repositories and roots.

The Task status remains Daemon-owned. No declared Verification command,
repository Verification, full package suite, commit, push or pull request
was run. No other Task or Task Graph file was edited. Missing-prefix path
comparison and consumer coverage remain with task_02 and task_03.
