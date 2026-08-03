---
task: task_04
spec: 0074-git-spawn-economy
status: pending
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

- [ ] Combine the resolution queries in `repository.go`.
- [ ] Combine the sequence in `assets_sync_git.go`.
- [ ] Pin the positional output parsing with a test.
- [ ] Prove the error cases still distinguish.

## Acceptance Criteria

- [ ] One resolution costs measurably fewer spawns, proven by a test
      counting invocations through the package's runner seam.
- [ ] Combined-output parsing has a test pinning the order.
- [ ] All existing resolution and error-case tests pass unmodified.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`
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
