---
title: "Final Push: push only when Review Issues are terminal"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 26
  - 33
  - 34
  - 35
  - 36
blocked_by:
  - 05-fetch-coderabbit-round-artifacts.md
  - 09-daemon-verification-and-batch-commits.md
---

# Final Push: push only when Review Issues are terminal

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add the Daemon-owned end-of-resolution mutation boundary. Roundfix should
resolve source threads for assigned terminal Review Issues only after successful
verification, and it should run Final Push only when no Unresolved Review Issues
remain and push safety requirements are satisfied.

## Acceptance criteria

- [ ] Source threads for `resolved` and `invalid` assigned issues are resolved
      only after Batch verification succeeds.
- [ ] Older `duplicated` occurrences remain local-only and do not trigger
      separate source thread resolution.
- [ ] Final Push is blocked while any Unresolved Review Issues remain.
- [ ] Final Push requires auto-push, auto-commit, known upstream target, and
      local commits not present on the target branch.
- [ ] The push sends local HEAD to the PR Head Branch without force-push.
- [ ] Tests prove no per-issue, per-Batch, or per-Round push happens.

## Blocked by

- 05-fetch-coderabbit-round-artifacts.md
- 09-daemon-verification-and-batch-commits.md
