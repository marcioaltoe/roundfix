---
title: "Preflight: validate git, PR, and push safety inputs"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 1
  - 5
  - 6
  - 35
  - 36
blocked_by:
  - 01-cli-command-contract-and-exit-codes.md
---

# Preflight: validate git, PR, and push safety inputs

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add the git and pull request discovery needed before Roundfix can act on an Open
Pull Request. The user-facing behavior should detect repository state, branch
state, upstream state, pull request metadata, and dirty worktree blockers before
any long-running work or mutation starts.

## Acceptance criteria

- [ ] Commands can detect git root, current branch, current HEAD, upstream
      remote, upstream branch, unpushed commit count, and dirty worktree status.
- [ ] Commands can resolve Open Pull Request metadata, PR Head Branch, and Head
      Repository from explicit input or safe inference.
- [ ] Dirty worktree blockers outside the Artifact Directory fail with exit `2`,
      list changed paths and statuses, and describe the required user action.
- [ ] Auto-push preflight rejects missing upstream information and never plans a
      force-push.
- [ ] Tests cover fake git and fake pull request metadata without live network
      calls.

## Blocked by

- 01-cli-command-contract-and-exit-codes.md
