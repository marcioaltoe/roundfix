---
title: "TUI: collect input and monitor Live Run state"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 51
  - 52
  - 53
  - 54
  - 55
blocked_by:
  - 01-cli-command-contract-and-exit-codes.md
  - 05-fetch-coderabbit-round-artifacts.md
  - 08-agent-runtime-batch-resolution.md
  - 11-watch-review-round-loop.md
---

# TUI: collect input and monitor Live Run state

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add the MVP terminal UI for Interactive Input and Live Run monitoring. The UI
should collect missing command parameters only when allowed, then show Run
metadata, Review Issues, streaming Agent output, verification output, budget
state, git state, and keybindings while active work runs.

## Acceptance criteria

- [ ] Missing required parameters open Interactive Input unless `--no-input` is
      set.
- [ ] Interactive Input suggests the current or remembered pull request and the
      configured or remembered Agent when available.
- [ ] Full Preflight Validation runs after Interactive Input and before waits,
      fetches, Agent starts, commits, or pushes.
- [ ] Live Run View shows repo, pull request, PR Head Branch, Review Source,
      Agent, HEAD, pipeline state, budget state, git state, and keybindings.
- [ ] The Review Issue sidebar groups issues by Round, severity, status, file,
      and line while the console streams Agent and verification output.
- [ ] TUI model and snapshot tests cover input defaults, missing input errors,
      sidebar grouping, streaming console output, and status strips.

## Blocked by

- 01-cli-command-contract-and-exit-codes.md
- 05-fetch-coderabbit-round-artifacts.md
- 08-agent-runtime-batch-resolution.md
- 11-watch-review-round-loop.md
