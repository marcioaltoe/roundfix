---
task: task_03
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: pending
type: backend
complexity: high
---

# Task 03: Classify a superseded QA-report-only Run Branch

## Overview

When a QA run fails and a later run passes, the first Run's branch holds nothing
but its QA report, yet reconciliation calls it `unintegrated` and preserves it
forever. Prove supersession from Git evidence the classifier already reaches —
the Daemon's deterministic QA-commit contract plus the existing recency
comparator — and let `--apply` release those branches.

## Requirements

1. MUST add a shared probe in `internal/worktree` that reports whether every
   commit in `target..run` matches the Daemon QA-commit contract for the Run's
   Spec and every changed path lies under that Spec's `qa/` directory,
   including the archived path.
2. MUST classify a terminal Implement Run as `superseded` only when all three
   hold: the probe proves QA-report-only, and the target branch's newest QA
   report for the Spec is newer than the branch's report by the existing
   recency comparator.
3. MUST fall back to `unintegrated` with `preserve` whenever any step is
   unprovable — an unparseable commit, a missing target-side report, or a
   target report that is not newer.
4. MUST let `--apply` release a `superseded` Run, re-proving the classification
   against freshly read heads before cleanup, and record the superseding report
   path as the release reason.
5. MUST add the classification value and a `superseded` summary counter to
   `roundfix-reconcile/v1` additively, leaving the schema version unchanged.
6. MUST mirror any new field in the terminal Run snapshot so the apply-time
   freshness proof still compares complete state.
7. MUST NOT change the Run Database schema.

## Subtasks

- [ ] Add the QA-report-only probe over the QA-commit message contract and the
      changed-path check.
- [ ] Classify `superseded` at the existing ancestry-miss branch, before
      settling `unintegrated`.
- [ ] Extend the `--apply` allowance with re-proof and the release reason.
- [ ] Cover the classification table and every fallback.

## Acceptance Criteria

- [ ] A branch holding only QA-report commits, whose Spec has a newer report on
      the target branch, classifies `superseded`.
- [ ] The same branch with one non-QA commit classifies `unintegrated` and is
      preserved.
- [ ] The same branch with no target-side report, or an older one, classifies
      `unintegrated` and is preserved.
- [ ] `--apply` releases a re-proven `superseded` Run and records the
      superseding report path; a Run that stops proving it is preserved.
- [ ] `roundfix-reconcile/v1` reports the classification and the counter at the
      unchanged schema version.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/worktree/ ./internal/cli/` — expected: pass,
  including every fallback case.

## References

`_prd.md` → Goal 3, Stories 3–4, Features 4–6; `_techspec.md` → Build Order 3,
API Contracts; ADR-0053.
