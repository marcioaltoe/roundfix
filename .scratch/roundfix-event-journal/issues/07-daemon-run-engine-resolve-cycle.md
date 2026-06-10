---
title: "Daemon Run engine: resolve cycle extraction with explicit dependencies"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 5
  - 6
  - 7
  - 10
  - 11
  - 12
  - 14
  - 15
  - 16
blocked_by:
  - 03-runevent-seam-runner-and-writersink.md
---

# Daemon Run engine: resolve cycle extraction with explicit dependencies

## Parent

.scratch/roundfix-event-journal/04-daemon-run-engine-prd.md

## What to build

Structural core of the Daemon Run engine PRD: a Run engine in the daemon module executing one resolve cycle over a
validated plan (per Batch: Agent → verify issue statuses → commit when auto-commit → resolve Review Source threads
for resolved/invalid), returning per-Batch outcomes and remaining Unresolved Review Issue count. Final Push is a
separate explicit engine operation; gating (no Unresolved Review Issues, auto-push on) stays in the caller per
ADR 0001. The engine updates intermediate Run states; Run creation and terminal completion stay with the caller.
Constructor takes an explicit dependencies struct (runner, verifier, committer, pusher, Review Source, rounds/store,
clock, Run Event sink); CLI package globals for orchestration die. The resolve command uses the engine. Stop
semantics move with the orchestration: no new Batch/verify/commit/push/source mutation after cancellation; stop
published to the sink; Agent worktree changes preserved.

## Acceptance criteria

- [ ] Engine proves resolve→verify→commit→source-resolution contract with fake collaborators and a temp store; no terminal, no network
- [ ] Final Push is a separate operation; caller-side gating decides when it runs
- [ ] Intermediate Run states updated during the cycle; terminal completion remains caller-owned
- [ ] CLI orchestration globals removed; resolve command delegates to the engine
- [ ] Stop tests: no post-stop daemon actions, stop event reaches sink, worktree preserved
- [ ] Existing CLI orchestration tests and fakes migrate to engine tests
- [ ] Batch failure policy unchanged: verification/Agent failure fails the Run

## Blocked by

- 03-runevent-seam-runner-and-writersink.md
