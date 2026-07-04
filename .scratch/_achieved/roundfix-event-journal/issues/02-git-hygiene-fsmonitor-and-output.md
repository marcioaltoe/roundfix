---
title: "Git hygiene: disable fsmonitor and separate stdout from stderr"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] All git invocations run with fsmonitor disabled
- [x] Status/porcelain parsing reads stdout only; stderr surfaces in error context
- [x] Test proves fsmonitor-style stderr noise does not corrupt dirty-worktree or status detection
- [x] No combined-output parsing remains for git commands

## Blocked by

None - can start immediately

## Comments

**2026-06-10 (agent):** Both product git runners — `preflight.ExecGitRunner.RunGit` and `daemon.runGit` — now
prepend `-c core.fsmonitor=false` to every invocation and run with separate stdout/stderr buffers. Parsed output
comes from stdout only; stderr (falling back to stdout, then the exec error) feeds the error detail. No
`CombinedOutput` git call remains in production code; remaining `CombinedOutput` sites are `gh` invocations, ACP
runtime probes, and test setup helpers. New tests in `internal/preflight/preflight_test.go` prove: checkout's
stderr chatter never reaches parsed output, failure errors carry stderr detail, the per-invocation fsmonitor
override beats repo-local `core.fsmonitor=true`, and `status --porcelain=v1 -z` returns the exact dirty record with
a noisy fsmonitor hook configured. Verification: `rtk go vet ./...` clean, `rtk go test ./...` 152 passed in 14
packages, `rtk go run ./cmd/roundfix --help` green.
