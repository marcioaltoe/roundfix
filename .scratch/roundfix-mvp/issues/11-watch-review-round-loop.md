---
title: "Watch: repeat CodeRabbit Rounds to a terminal outcome"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 37
  - 38
  - 39
  - 40
  - 41
  - 42
  - 43
  - 44
blocked_by:
  - 05-fetch-coderabbit-round-artifacts.md
  - 09-daemon-verification-and-batch-commits.md
  - 10-final-push-and-source-resolution.md
---

# Watch: repeat CodeRabbit Rounds to a terminal outcome

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Implement the durable foreground watch loop that waits for CodeRabbit review
status on the current PR HEAD, observes the quiet period, fetches unresolved
Review Issues, resolves Batches, verifies, commits, pushes when allowed, and
repeats until the Run reaches a documented terminal outcome.

## Acceptance criteria

- [ ] `watch` waits for CodeRabbit to review or settle the current HEAD before
      fetching issues.
- [ ] `watch` observes the configured poll interval, quiet period, review
      timeout, `max_rounds`, and `max_run_duration`.
- [ ] `Clean`, `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`, `Failed`, and
      `Stopped` terminal outcomes map to the documented exit behavior.
- [ ] `MaxRoundsReached` exits as a non-error terminal outcome and reports any
      remaining Unresolved Review Issues.
- [ ] Review timeout output offers a manual review trigger action without
      automatically posting one.
- [ ] Tests run the full loop with fake Review Source, fake Agent, fake git,
      fake verification, and fake clock boundaries.

## Blocked by

- 05-fetch-coderabbit-round-artifacts.md
- 09-daemon-verification-and-batch-commits.md
- 10-final-push-and-source-resolution.md
