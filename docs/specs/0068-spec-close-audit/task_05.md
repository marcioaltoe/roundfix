---
task: task_05
spec: 0068-spec-close-audit
status: pending
type: test
complexity: medium
---

# Task 05: Replay the session that motivated this Spec

## Overview

This Spec's first Success Metric is a replay. One session left four kinds of
residue at once — two Supervisor scratch worktrees, one orphaned Run Worktree,
one stale remote backup branch, and two branches held by unmerged Pull
Requests — and the maintainer found all four by hand with `git branch -l` and
`git worktree list`.

The audit must report all four.

## Requirements

1. MUST build a fixture reproducing that end state, recording in the fixture
   which finding it reproduces so its provenance is not folklore.
2. MUST assert all four residue kinds are reported, each with its evidence.
3. MUST assert the two branches backing unmerged Pull Requests classify
   `pull-request` and are never reported as reclaimable.
4. MUST assert the orphaned Run Worktree — the one whose target branch the
   squash merge deleted — resolves through task_01's content check rather than
   as `unknown`.
5. MUST assert the two scratch worktrees, which have no Run, classify
   `preserved` unless their branch is pushed and merged.
6. MUST NOT modify any archived Spec artifact.

## Subtasks

- [ ] Build the four-residue fixture with its provenance recorded.
- [ ] Assert each kind and its evidence.
- [ ] Assert the Pull Request branches are never reclaimable.
- [ ] Assert the orphaned worktree resolves by content.

## Acceptance Criteria

- [ ] The replay reports all four residue kinds.
- [ ] The fixture records the finding it reproduces by path.
- [ ] Neither Pull Request branch is reported as reclaimable.
- [ ] The orphaned Run Worktree resolves by content, not `unknown`.
- [ ] No file under `docs/specs/_archived/` is modified.

## Context

- instruction: `docs/findings/2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/specaudit -count=1 -run 'Replay|Session|Residue' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the replay tests ran and passed.
- `go test ./internal/specaudit ./internal/worktree -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics 1, 2 and 3.
- `_techspec.md` → Testing Approach; Build Order 5.
