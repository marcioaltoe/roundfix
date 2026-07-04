---
title: "Snapshot-diff Batch commits with explicit commit skip"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 1
  - 2
  - 3
blocked_by:
  - 02-git-hygiene-fsmonitor-and-output.md
  - 07-daemon-run-engine-resolve-cycle.md
---

# Snapshot-diff Batch commits with explicit commit skip

## Parent

.scratch/roundfix-event-journal/04-daemon-run-engine-prd.md

## What to build

Commit strategy from the Daemon Run engine PRD, fixing two production defects: the engine captures worktree status
before starting the Agent and after it finishes, and the Batch commit stages only paths changed between snapshots
(Agent-touched, legitimately including assigned Review Issue files). Pre-existing or mid-Run user changes are never
staged, even if they slipped past Preflight Validation. When nothing changed, the engine skips the commit and the
Batch still succeeds — no more nothing-to-commit failures on triage-only Batches.

## Acceptance criteria

- [x] Pre-existing unrelated changes are never staged into a Batch commit
- [x] Mid-Run user edits stay out of Batch commits
- [x] Agent-touched paths, including assigned issue files, are staged
- [x] Triage-only Batch: commit skipped, Batch succeeds, Run continues
- [x] Skip decision surfaced in command output (journaled event lands with the daemon event stream issue)

## Blocked by

- 02-git-hygiene-fsmonitor-and-output.md
- 07-daemon-run-engine-resolve-cycle.md

## Comments

**2026-06-10 (agent):** The engine now snapshots the worktree (`daemon.WorktreeSnapshotter`, real implementation
`GitWorktreeSnapshotter` over `git status --porcelain=v1 -z` with fsmonitor disabled and stdout-only parsing,
rename/copy second fields handled) at Batch start and again just before the commit, and stages exactly
`diffSnapshots(before, after)` — paths dirty after the Batch that were not already dirty at Batch start, sorted,
with `.roundfixrc.yml` excluded as defense in depth. Anything dirty at Batch start — pre-existing user work or
mid-Run edits from earlier in the Run — never reaches a Batch commit. `CommitRequest` gained `Paths` (exact-stage
mode, `git add -- <paths>`); the sweep-everything mode remains only as the legacy fallback when no paths are given.
Empty diff → commit skipped, Batch succeeds, Run continues; `BatchOutcome.CommitSkipped` records it and the
command output prints "Batch commit skipped: Batch NNN made no worktree changes." (the journaled event arrives
with issue 12). The snapshotter is a new engine dependency threaded through the CLI collaborators (tests use a
fake; engine fixtures use scripted snapshots). Tests: fake-level staging proof (user-wip excluded, agent code +
assigned issue file staged), triage-only skip with surfaced output, and a real-git integration test proving a
pre-existing tracked edit and untracked user file stay out of the commit and remain in the worktree while only
`agent-fix.go` lands. Verification: `rtk go vet ./...` clean, `rtk go test ./...` 187 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
