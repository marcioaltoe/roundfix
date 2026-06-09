---
title: "Daemon: verify and commit successful Batches"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 30
  - 31
  - 32
blocked_by:
  - 08-agent-runtime-batch-resolution.md
---

# Daemon: verify and commit successful Batches

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

After an Agent completes a Batch, make the Daemon run the configured
verification gate and create one local commit for each successful Batch when
auto-commit is enabled. This slice must preserve the rule that committing a
successful Batch is allowed even while other Unresolved Review Issues remain.

## Acceptance criteria

- [x] The Daemon runs the configured verification command after Agent completion
      and before a Batch commit.
- [x] A Batch succeeds only when every assigned Review Issue is terminal and
      verification exits successfully.
- [x] Verification failure exits as a Run failure and does not create a commit,
      push, or resolve source threads.
- [x] Auto-commit creates one local commit per successful Batch.
- [x] Successful Batch commits remain allowed while other Unresolved Review
      Issues remain.
- [x] Tests use fake git and fake verification runners to prove commit and
      failure behavior.

## Blocked by

- 08-agent-runtime-batch-resolution.md
