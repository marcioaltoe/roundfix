---
title: "Snapshot-diff Batch commits with explicit commit skip"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] Pre-existing unrelated changes are never staged into a Batch commit
- [ ] Mid-Run user edits stay out of Batch commits
- [ ] Agent-touched paths, including assigned issue files, are staged
- [ ] Triage-only Batch: commit skipped, Batch succeeds, Run continues
- [ ] Skip decision surfaced in command output (journaled event lands with the daemon event stream issue)

## Blocked by

- 02-git-hygiene-fsmonitor-and-output.md
- 07-daemon-run-engine-resolve-cycle.md
