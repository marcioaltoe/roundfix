---
task: task_03
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: backend
complexity: high
---

# Task 03: Move every consumer onto the resolver

## Overview

Every package that composes an archive path stops doing so and asks the resolver
instead, and the two hardcoded literals in the Spec checker are deleted. Nothing
moves on disk in this Task — the resolver still answers today's paths — so this
is the prefactoring that makes Task 04's relocation a one-line change rather than
five.

## Requirements

1. MUST replace every self-composed archive path in the Spec checker, the Spec
   audit, the worktree QA evidence path, and the CLI with a resolver call.
2. MUST delete the two hardcoded archive literals in the Spec checker.
3. MUST leave every observable path identical, because the resolver still
   answers the locations Task 01 recorded.
4. MUST leave no package other than the resolver's own expressing the archive
   layout.

## Subtasks

- [ ] Move the Spec checker onto the resolver and delete its two literals.
- [ ] Move the Spec audit, worktree, and CLI onto the resolver.
- [ ] Prove no other package expresses the layout.

## Acceptance Criteria

- [ ] Every consumer resolves its archive path through the one owner.
- [ ] The two literals are gone.
- [ ] Every observable path is unchanged.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/speccheck.go`
- `internal/speccheck/speccheck_test.go`
- `internal/specaudit/specaudit.go`
- `internal/specaudit/specaudit_test.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `internal/cli/archive.go`
- `internal/cli/archive_test.go`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_03.md`

## Verification

- `test -z "$(grep -rnE '"_archived"|_archived/' internal/speccheck/*.go internal/specaudit/*.go internal/worktree/*.go internal/cli/archive.go | grep -v '_test.go')"` — expected: exits 0, proving no consumer still expresses the layout itself.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck ./internal/specaudit ./internal/worktree ./internal/cli -count=1 2>&1 | grep -c '^ok' | grep -q '^4$'` — expected: exits 0, proving all four packages pass rather than some being skipped.
- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveLayoutCharacterization' -count=1 2>&1 | grep -q '^ok'` — expected: exits 0, proving the recorded paths did not move in this Task.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 3; System Architecture.
