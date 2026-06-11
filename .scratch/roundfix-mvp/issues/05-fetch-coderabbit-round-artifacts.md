---
title: "Fetch: import CodeRabbit findings as Round artifacts"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 15
  - 16
  - 17
  - 57
  - 58
blocked_by:
  - 04-run-state-fetch-runs-and-active-locks.md
---

# Fetch: import CodeRabbit findings as Round artifacts

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Implement the first Review Source path: `fetch` imports unresolved CodeRabbit
findings for an Open Pull Request, maps them into Review Issues, and persists a
Round of markdown artifacts. The command must remain a tracked Fetch Run and
must not start Agent work or mutate git.

## Acceptance criteria

- [x] `fetch` retrieves CodeRabbit review comments and review threads for the
      current pull request head and filters out non-CodeRabbit and resolved
      source threads.
- [x] Review Source identifiers, review hashes, source review identifiers, and
      source review submission timestamps are preserved in markdown artifacts.
- [x] Markdown artifacts include the MVP frontmatter fields and allowed Review
      Issue statuses.
- [x] Fetched reviewer text is treated as untrusted content and is never
      executed, shell-interpolated, or logged with credentials.
- [x] A successful fetch persists one Round, records `Fetched`, exits `0`, and
      does not start an Agent, commit, push, or resolve source threads.
- [x] Tests use fake GitHub or fake CodeRabbit responses and temporary artifact
      storage.

## Blocked by

- 04-run-state-fetch-runs-and-active-locks.md
