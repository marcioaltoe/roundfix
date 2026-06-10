---
title: "Git hygiene: disable fsmonitor and separate stdout from stderr"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 4
blocked_by: []
---

# Git hygiene: disable fsmonitor and separate stdout from stderr

## Parent

.scratch/roundfix-event-journal/04-daemon-run-engine-prd.md

## What to build

Standalone defect fix from the Daemon Run engine PRD: fsmonitor warnings (`fsmonitor_ipc__send_query`) pollute parsed
git output in real Runs. Make every Roundfix git invocation disable fsmonitor for the invocation, and parse stdout
separately from stderr everywhere — combined-output parsing is removed.

## Acceptance criteria

- [ ] All git invocations run with fsmonitor disabled
- [ ] Status/porcelain parsing reads stdout only; stderr surfaces in error context
- [ ] Test proves fsmonitor-style stderr noise does not corrupt dirty-worktree or status detection
- [ ] No combined-output parsing remains for git commands

## Blocked by

None - can start immediately
