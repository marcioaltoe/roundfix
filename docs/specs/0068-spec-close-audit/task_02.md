---
task: task_02
spec: 0068-spec-close-audit
status: pending
type: backend
complexity: high
---

# Task 02: Classify every survivor with its evidence

## Overview

The audit's core: enumerate the branches and worktrees a Spec's cycle produced
and classify each survivor as backing an open Pull Request, holding
unintegrated work, reclaimable residue, or unclassifiable and therefore
preserved. Every classification carries the evidence that produced it.

Verifiable on its own through fixture repositories, one per kind.

## Requirements

1. MUST enumerate the branches and worktrees associated with a Spec, resolving
   its Runs from the Run Database rather than requiring the caller to know them.
2. MUST classify each survivor as `pull-request`, `pending`, `residue`, or
   `preserved`, and MUST attach the evidence string that produced it.
3. MUST classify a worktree with no matching Run as `preserved`, never as
   residue. Scratch worktrees live outside the Run Database, and assuming
   residue there is how work gets deleted.
4. MUST attach the exact reclaim command to a `residue` survivor and MUST NOT
   execute it. The audit reports; the operator reclaims.
5. MUST NOT touch, lock, or reclaim anything, and MUST never classify a branch
   or worktree belonging to an Active Run as residue, per ADR-0052.
6. MUST open no network connection and write no file.

## Subtasks

- [ ] Add the package with its kinds, survivor model, and result.
- [ ] Enumerate branches, worktrees, and the Spec's Runs.
- [ ] Classify each survivor and attach its evidence.
- [ ] Add the Active Run guard before the classifier that could violate it.
- [ ] Add one fixture repository per kind.

## Acceptance Criteria

- [ ] A branch backing an open Pull Request classifies `pull-request` with
      evidence naming the Pull Request.
- [ ] A branch holding unintegrated commits classifies `pending` with evidence
      naming what is unintegrated.
- [ ] A merged branch with no Pull Request classifies `residue` and carries the
      exact reclaim command as a string.
- [ ] A worktree with no matching Run classifies `preserved`, never `residue`.
- [ ] An Active Run's branch and worktree are never classified `residue`,
      proven by an injected active fixture.
- [ ] Every survivor in every fixture carries a non-empty evidence string.
- [ ] The package opens no transport and writes no file, proven by a
      repository-wide check rather than by inspection.

## Context

- interface: `internal/store`
- interface: `internal/gittest`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/specaudit -count=1 -run 'Classif|Survivor|Active|Evidence' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the classification tests ran and passed.
- `go test ./internal/specaudit -count=1` — expected: exit 0.
- `if grep -rn "net/http\|os.Create\|os.WriteFile\|os.RemoveAll" internal/specaudit --include="*.go" | grep -v "_test.go" | grep -q .; then exit 1; fi`
  — expected: exit 0; the audit opens no transport and removes nothing.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 4, 5 and 6; Decisions. Core Feature 4 is
  satisfied here by classifying a pushed-and-merged scratch worktree as
  `residue` with its reclaim command; the reclaiming stays the operator action
  the Decisions require.
- `_techspec.md` → Interfaces; Build Order 2.
- ADR-0052.
