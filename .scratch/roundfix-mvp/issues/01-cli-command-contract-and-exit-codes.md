---
title: "CLI: enforce MVP command contract and exit codes"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 2
  - 3
  - 4
  - 50
  - 62
blocked_by: []
---

# CLI: enforce MVP command contract and exit codes

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Turn the reserved `fetch`, `resolve`, and `watch` command names into real MVP
command surfaces with shared option parsing, command-specific required inputs,
Preflight Validation hooks, and documented exit code behavior. This slice does
not need live Review Source, Agent, Run Database, or git mutations yet; it must
create the visible command contract that later slices plug into.

## Acceptance criteria

- [ ] `fetch`, `resolve`, and `watch` accept their MVP flags and reject invalid
      command input with the Preflight Validation exit code.
- [ ] Unknown commands, invalid flag values, missing required inputs under
      `--no-input`, and setup failures produce actionable stderr and exit `2`.
- [ ] Non-error help and version behavior still exits `0`.
- [ ] Command tests cover stdout, stderr, and exit codes through the public CLI
      runner seam.
- [ ] No Run state, artifacts, Agent process, Review Source call, commit, or
      push is attempted when command parsing or early Preflight Validation
      fails.

## Blocked by

None - can start immediately.
