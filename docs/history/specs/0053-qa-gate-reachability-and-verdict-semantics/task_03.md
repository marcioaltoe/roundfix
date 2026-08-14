---
task: task_03
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: completed
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

## Result

### Implementation

- Added a shared `internal/worktree` Git probe that requires a non-empty
  `target..run` range, the exact Daemon QA commit subject/trailer contract on
  every commit, and NUL-delimited changed paths confined to the active or
  archived Spec `qa/` directory.
- Added proof-based `superseded` classification at the ancestry miss. Both
  heads' report paths are read with `git ls-tree` and ranked through the
  existing `internal/spec` report recency comparator; any missing, malformed,
  or non-newer evidence remains `unintegrated` with preserve behavior.
- Extended apply revalidation, the terminal snapshot, cleanup journaling, text
  and JSON output, and the additive summary counter. The superseding report
  path is carried as explicit evidence and persisted as the release reason.
  No Run Database table, migration, or stored Run field changed.

### Focused checks

- `GOCACHE=.gocache go test ./internal/worktree -run 'Test(QAReportOnlyBranch|InspectTerminalRun|ApplyTerminalRun)'`
  — passed; covers active and archived QA-only branches, newer-report
  classification, non-QA/missing/older/equal fallbacks, successful release,
  forged snapshot refusal, and fresh-proof failure preservation.
- `GOCACHE=.gocache go test -count=1 ./internal/cli -run '^TestRunReconcile'`
  — passed; covers existing reconcile behavior plus superseded JSON,
  unchanged `roundfix-reconcile/v1`, summary count, apply cleanup, and
  journaled superseding-report reason.
- `GOCACHE=.gocache go test -count=1 ./internal/store -run '^TestReconcileIntegration'`
  — passed; covers reconciliation validation and event persistence without a
  schema change.
- `GOCACHE=.gocache go test -count=1 ./internal/spec -run 'Test(NewestQAReport|QAVerdictPrefersTheSameDateRerun)'`
  — passed; confirms the extracted path-list entry point preserves the
  established date/sequence ordering.
- `git diff --check` — passed.

The direct focused Go checks initially could not write the host Go cache; the
final checks above used the repository's ignored `.gocache`. The commands in
`## Verification` were not run because Daemon Verification owns them.

### Acceptance evidence

- QA-only branch plus newer target report: the classification table proves
  `superseded` for both active and archived Spec paths.
- Non-QA commit: the same table proves `unintegrated`; apply is not offered.
- Missing, older, or equal target report: each remains `unintegrated` with
  preserve behavior.
- Apply: a freshly re-proven superseded Run releases its Worktree and branch
  and journals the superseding report path; mutated proof or snapshot metadata
  preserves both Git surfaces.
- CLI contract: JSON exposes `classification: superseded`, the explicit
  report path, and `summary.superseded` while retaining schema version
  `roundfix-reconcile/v1`.
