---
task: task_01
spec: 0008-worktree-isolation
status: pending
type: backend
complexity: high
---

# Task 01: Build the worktree package

## Overview

Create `internal/worktree`: Run Worktree creation on a named Run Branch,
keep/remove/prune lifecycle, untracked-file provisioning, and the
porcelain-only integration protocol from ADR-0024 — reproducing the
empirically verified git matrix from the techspec as the package's test
suite. Nothing is wired into Runs yet. Verifiable alone over hermetic temp
repositories.

## Requirements

1. MUST create worktrees with `git worktree add -b roundfix/run-<id>
   <path> <headSHA>`, path under Roundfix Home
   (`worktrees/<repo-id>/<run-id>`), returning a Ref carrying run id, path,
   branch, and the owning user root; copy-list entries are copied by
   relative path with per-file stderr notes for missing sources, never
   failures.
2. MUST implement `Integrate` per ADR-0024 exactly: target branch checked
   out in the user root → `git -C <userRoot> merge --ff-only <runSHA>`;
   checked out nowhere → `git merge-base --is-ancestor` then
   `git branch -f`; refusals (overlap, divergence, non-ancestry) return a
   `pending` result with the reason and MUST leave the target branch
   unmoved; `git update-ref` is never used.
3. MUST implement `CleanupClean` (worktree remove + Run Branch delete, only
   for integrated Clean outcomes) and `PruneTerminal` (`git worktree prune`
   plus directory/branch removal exclusively for Runs the callback reports
   terminal-Clean).
4. MUST keep dirty worktrees intact on every path — `git worktree remove`
   without `--force` refusing dirty trees is relied on, not overridden.
5. MUST use the hermetic git test helpers (0003 discipline) and context-
   first git invocations through the package's own small runner.

## Subtasks

- [ ] Create with named Run Branch, Roundfix Home paths, copy-list
- [ ] Integration protocol: two cases plus pending reasons
- [ ] CleanupClean and PruneTerminal
- [ ] Verified-matrix test suite over hermetic temp repos

## Acceptance Criteria

- [ ] Creation succeeds while the source branch is checked out in the main
      repo (named-branch path proven); the worktree starts clean and
      independent.
- [ ] Integration tests reproduce the full matrix: ff on clean checkout, ff
      preserving non-overlapping dirt, pending on overlap with branch
      unmoved and user dirt intact, pending on divergence, branch-move when
      the user switched away, pending on non-ancestry; after every pending
      case the user checkout's `git status` shows no phantom staged
      entries.
- [ ] Crash simulation (`rm -rf` of a worktree) is reaped by PruneTerminal
      only for terminal-Clean run ids; kept worktrees of other runs
      survive the sweep.
- [ ] Full suite passes; no production wiring outside the new package.

## Verification

- `rtk go test ./internal/worktree/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → Core Features 1, 3; Decisions. `_techspec.md` → Interfaces,
Worktree creation, Integration protocol, Build Order 1. ADR-0023, ADR-0024.
Round-1 findings 10, 15.
