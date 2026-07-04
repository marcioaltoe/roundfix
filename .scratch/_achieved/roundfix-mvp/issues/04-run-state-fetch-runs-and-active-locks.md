---
title: "Run state: record Fetch Runs and lock Active Runs"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 11
  - 12
  - 13
  - 14
  - 48
  - 49
blocked_by:
  - 02-preflight-config-and-artifact-directory.md
  - 03-preflight-git-pr-and-push-safety.md
---

# Run state: record Fetch Runs and lock Active Runs

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add the central Run Database behavior needed to create durable Run history and
prevent conflicting Active Runs. The first externally visible path should be a
tracked Fetch Run that records history without starting an Agent, committing, or
pushing.

## Acceptance criteria

- [x] Roundfix opens or creates the central Run Database under Roundfix Home and
      applies migrations before creating any Run.
- [x] `fetch` creates a Fetch Run only after Preflight Validation succeeds.
- [x] A successful Fetch Run reaches the `Fetched` terminal outcome and releases
      any Active Run lock it acquired.
- [x] A second Active Run for the same Head Repository and PR Head Branch is
      rejected with the existing run identifier and state.
- [x] Terminal outcomes, including `Fetched` and `Stopped`, release the logical
      Active Run lock.
- [x] Preflight Validation failures do not create Run records or Active Run
      locks.

## Blocked by

- 02-preflight-config-and-artifact-directory.md
- 03-preflight-git-pr-and-push-safety.md
