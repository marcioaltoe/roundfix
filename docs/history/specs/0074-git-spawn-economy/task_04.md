---
task: task_04
spec: 0074-git-spawn-economy
status: completed
type: backend
complexity: medium
---

# Task 04: Combine repository resolution queries

## Overview

Repository resolution asks git one question per spawn: `repository.go`
resolves root, HEAD validity, and object format in three `rev-parse` calls
plus a `rev-list` for the root commit; `assets_sync_git.go` repeats the
pattern. `rev-parse` answers several queries in one invocation with
line-ordered output, so resolution can ask its questions in one spawn where
git's own interface allows combining them.

The census counted 3,518 `rev-parse` spawns per suite run. Not all combine —
only queries against the same repository state at the same moment may share
a spawn, and this Task combines only within those scopes, never caching
across mutations (ADR-0090).

## Requirements

1. MUST combine the multi-fact queries in `resolveRepository`
   (`internal/baseline/repository.go`) into the fewest spawns git's
   interface supports, with positional parsing pinned by a test.
2. MUST apply the same combination to the resolution sequence in
   `assets_sync_git.go`.
3. MUST NOT cache any repository fact across a mutation boundary; combined
   queries share one spawn only when they read the same immutable state.
4. MUST keep every error case distinguishable: a missing HEAD, a non-repo
   directory, and an unknown object format report exactly as they do
   today.
5. MUST keep the whole `internal/baseline` suite passing unmodified.

## Subtasks

- [x] Combine the resolution queries in `repository.go`.
- [x] Combine the sequence in `assets_sync_git.go`.
- [x] Pin the positional output parsing with a test.
- [x] Prove the error cases still distinguish.

## Acceptance Criteria

- [x] One resolution costs measurably fewer spawns, proven by a test
      counting invocations through the package's runner seam.
- [x] Combined-output parsing has a test pinning the order.
- [x] All existing resolution and error-case tests pass unmodified.
- [x] `git status --porcelain` shows no path outside `internal/baseline/`
      and this task file.

## Verification

- `go test ./internal/baseline -count=1 -run 'Repository|Resolve' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Feature 2.
- `_techspec.md` → Implementation Design (the combined rev-parse sketch);
  Risks (combined output is order-dependent).
- ADR-0090.

## Result

Implementation:

- Repository identity now reads the worktree root, object format, and verified
  HEAD through one line-ordered `rev-parse`; `rev-list` remains the separate
  root-commit query that Git cannot combine with it.
- Assets-sync checkout inspection now reads its root and verified HEAD through
  one `rev-parse`, while status and origin remain fresh reads after resolution.
- Combined resolution values are parsed positionally. Repository-resolution
  failure handling retains distinct missing-HEAD, non-repository, and unknown
  object-format errors.

Focused checks:

- Pre-change signal: the focused runner-seam test failed with `unexpected Git
  command` because production still issued separate `rev-parse` calls.
- `rtk proxy env GOCACHE=$PWD/.gocache go test ./internal/baseline -count=1 -run
  '^(TestRepositoryInspectionUsesNarrowReadOnlyGitCommands|TestRepositoryInspectionParsesCombinedResolutionPositionally|TestRepositoryInspectionCombinedResolutionErrorsRemainDistinct|TestInspectAssetsSyncCheckoutCombinesResolutionQueries)$'
  -v` — pass.
- `rtk proxy env GOCACHE=$PWD/.gocache go test ./internal/baseline -count=1 -run
  '^(TestRepositoryIdentityEquivalentClones|TestRepositoryIdentityRequiresCommittedGitWorktree|TestAssetsSyncProvenanceAndPreMutationRefusals)$'
  -v` — pass against real Git repositories.
- `git diff --check` — pass.
- `git -c core.fsmonitor=false status --porcelain` — only this task file and
  four paths under `internal/baseline/` are modified.

Acceptance evidence:

1. `TestRepositoryInspectionUsesNarrowReadOnlyGitCommands` pins two Git calls
   per successful identity resolution instead of four; the first call combines
   all three `rev-parse` facts. `TestInspectAssetsSyncCheckoutCombinesResolutionQueries`
   pins one resolution call before the independent status and origin reads.
2. `TestRepositoryInspectionParsesCombinedResolutionPositionally` swaps the
   object-format and HEAD lines and proves the second line owns the format
   position.
3. The focused real-Git checks above pass the existing clone, missing-HEAD,
   non-repository, and assets-sync refusal coverage without weakening those
   tests. The new table test asserts exact outer errors for missing HEAD,
   non-repository input, and an unknown object format.
4. The scoped status inspection lists no path outside `internal/baseline/` and
   `docs/specs/0074-git-spawn-economy/task_04.md`.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
