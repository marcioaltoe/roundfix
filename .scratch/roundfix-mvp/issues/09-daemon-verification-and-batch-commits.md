---
title: "Daemon: verify and commit successful Batches"
type: AFK
category: enhancement
state: ready-for-agent
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

- [ ] The Daemon runs the configured verification command after Agent completion
      and before a Batch commit.
- [ ] A Batch succeeds only when every assigned Review Issue is terminal and
      verification exits successfully.
- [ ] Verification failure exits as a Run failure and does not create a commit,
      push, or resolve source threads.
- [ ] Auto-commit creates one local commit per successful Batch.
- [ ] Successful Batch commits remain allowed while other Unresolved Review
      Issues remain.
- [ ] Tests use fake git and fake verification runners to prove commit and
      failure behavior.

## Blocked by

- 08-agent-runtime-batch-resolution.md
