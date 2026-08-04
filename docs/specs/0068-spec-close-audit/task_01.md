---
task: task_01
spec: 0068-spec-close-audit
status: pending
type: backend
complexity: high
---

# Task 01: Resolve integration by content when the branch is gone

## Overview

`gh pr merge --delete-branch` destroys the target branch every squash merge
resolves by name, so `inspectTerminalRun` degrades that Run to `unknown`
forever. The name is the handle the merge destroys; the content is what
integration actually means. This slice decides by content instead.

It leads the graph because it is the only change to existing behaviour and the
only one that can lose work if it is wrong.

## Requirements

1. MUST resolve a terminal Run whose recorded target branch no longer exists by
   comparing the Run Branch against the synced default branch.
2. MUST resolve `safe` only when no file is unique to the Run Branch and every
   shared implementation file matches. Anything else MUST resolve
   `unintegrated`.
3. MUST prove integration only, never disprove it. A false `safe` is the one
   error that loses work, so ambiguity resolves against reclaiming.
4. MUST NOT touch an Active Run or a live writer on any path, per ADR-0052.
5. MUST leave every other `reconcile` refusal exactly as it is today, including
   the dirty-worktree and unresolvable-branch cases.
6. MUST record the evidence that produced the resolution in the existing reason
   field, so a reader learns why rather than only what.

## Subtasks

- [ ] Add the content comparison at the missing-target-branch path.
- [ ] Resolve `safe` on full match, `unintegrated` otherwise.
- [ ] Record the deciding evidence in the reason.
- [ ] Add the deleted-target fixture and the Active Run guard.

## Acceptance Criteria

- [ ] A terminal Run whose target branch was deleted after a squash merge, and
      whose content matches the default branch, resolves `safe` where today it
      resolves `unknown`.
- [ ] The same Run with one file unique to the Run Branch resolves
      `unintegrated`, not `safe`.
- [ ] The resolution's reason names the evidence that decided it.
- [ ] An Active Run fixture is untouched and never resolved by this path.
- [ ] Every existing `reconcile` refusal still holds, asserted by the existing
      reconcile tests passing unchanged.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/gittest`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/worktree -count=1 -run 'Reconcil|Integration|Content|Terminal' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the reconciliation tests ran and passed.
- `go test ./internal/worktree ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metric 2.
- `_techspec.md` → Build Order 1; Risks & Considerations.
- ADR-0052, ADR-0053.
