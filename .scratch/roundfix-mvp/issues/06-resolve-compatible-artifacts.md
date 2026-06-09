---
title: "Resolve: select downloaded Compatible Artifacts"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 18
  - 19
  - 20
  - 21
blocked_by:
  - 05-fetch-coderabbit-round-artifacts.md
---

# Resolve: select downloaded Compatible Artifacts

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Make `resolve` operate only on downloaded Review Issue artifacts that match the
requested Open Pull Request and optional Round selector. This slice establishes
the command boundary: `resolve` consumes Compatible Artifacts and never fetches
new Review Source issues.

## Acceptance criteria

- [x] `resolve` rejects closed, missing, or mismatched pull requests during
      Preflight Validation.
- [x] Without a Round selector, `resolve` selects all downloaded Unresolved
      Review Issues across all Compatible Artifact Rounds for the pull request.
- [x] With a Round selector, `resolve` limits selection to that Compatible
      Artifact Round.
- [x] If no Compatible Artifacts exist, `resolve` exits `2` before creating a Run
      and tells the user to run `fetch` or `watch`.
- [x] `resolve` does not call the Review Source fetch path.
- [x] Tests cover compatible, incompatible, missing, and Round-filtered artifact
      sets.

## Blocked by

- 05-fetch-coderabbit-round-artifacts.md
