---
task: task_05
spec: 0068-spec-close-audit
status: completed
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

## Result

Added a real-Git replay fixture for the session recorded at
`docs/findings/2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md`.
The fixture creates two Supervisor scratch worktrees, a retained Run Worktree
whose target branch is deleted after a squash merge, a stale remote backup
branch, and two branches held by unmerged Pull Requests. The replay exposed
that survivor enumeration omitted remote-tracking refs, so the audit now reads
local and remote branch refs, excludes the remote default branch and its `HEAD`
alias, and reports remote residue with an exact `git push --delete` command
without executing it.

Focused checks run during implementation (the declared `## Verification`
commands remain Daemon-owned and were not run):

- Red signal: after correcting the fixture's canonical worktree path, `rtk env
  GOCACHE=/private/tmp/roundfix-spec-close-audit-task05-gocache go test
  ./internal/specaudit -run '^TestAuditReplaysMotivatingSessionResidue$'
  -count=1` failed because
  `origin/roundfix/run-run_20260731T195234Z_backup` was absent from the audit
  result.
- A second red replay reported the synced `origin/main` default target as
  reclaimable after remote enumeration was added. The replay passed after the
  audit excluded both `origin/HEAD` and the branch it resolves to.
- The same focused replay command passed after remote survivor enumeration was
  added.
- `rtk env GOCACHE=/private/tmp/roundfix-spec-close-audit-task05-gocache go
  test ./internal/specaudit -run
  '^TestAudit(ClassifiesPullRequestBranch|ClassifiesPendingBranch|ClassifiesResidueBranch|PreservesUnmatchedWorktree|PreservesActiveRunSurvivors|ReplaysMotivatingSessionResidue)$'
  -count=1` passed the replay and the five existing classification fixtures.
- `rtk env GOCACHE=/private/tmp/roundfix-spec-close-audit-task05-gocache go
  test ./internal/worktree -run
  '^TestInspectTerminalRunSafeWhenTargetDeletedAfterSquashMerge$' -count=1`
  passed task_01's deleted-target content-resolution fixture.
- `rtk env GOCACHE=/private/tmp/roundfix-spec-close-audit-task05-gocache go vet
  ./internal/specaudit` exited 0.
- `rtk git -c core.fsmonitor=false status --short` exited 0 and listed only this Task file,
  `internal/specaudit/audit.go`, and `internal/specaudit/audit_test.go`; no path
  under `docs/specs/_archived/` was modified.

Acceptance evidence:

1. `TestAuditReplaysMotivatingSessionResidue` passed and found both preserved
   scratch worktrees, the content-resolved orphaned Run Worktree, the
   reclaimable remote backup branch, and both Pull Request branches, with a
   non-empty evidence string on every returned survivor.
2. The replay fixture stores the exact source finding path in its
   `findingPath` field, and the passing test asserts that provenance.
3. Both unmerged Pull Request branches classified `pull-request`, named their
   Pull Request numbers in evidence, and carried empty reclaim strings.
4. The orphaned Run Worktree classified `residue` with evidence that its
   content is fully represented on `main`; the focused task_01 regression also
   classified the deleted-target shape `safe` rather than `unknown`.
5. The fresh changed-path inspection above found no modified archived Spec
   artifact.

Follow-ups: none discovered inside this Task's slice.
