---
title: "Agent: run a bounded Batch through an ACP Runtime"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 27
  - 28
  - 29
  - 56
  - 57
  - 58
blocked_by:
  - 07-resolve-deduplication-and-batches.md
---

# Agent: run a bounded Batch through an ACP Runtime

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add the Agent execution path for one bounded Batch. Roundfix should probe the
selected ACP Runtime, start it through the user's local authenticated setup,
send a strict child-agent prompt with assigned issue files, stream output, and
persist Agent logs while preventing the Agent from owning commits, pushes, or
Review Source mutations.

## Acceptance criteria

- [ ] Runtime probing supports the configured Agent, command overrides, model
      options, and actionable install or authentication diagnostics.
- [ ] The child-agent prompt includes assigned issue files, required triage
      steps, verification expectations, and forbidden actions.
- [ ] Agent output streams to the command path and is persisted with the Run.
- [ ] Agent failure marks the Batch path as failed without committing, pushing,
      or resolving source threads.
- [ ] Roundfix validates that every assigned Review Issue reaches a terminal
      local status before a Batch can proceed.
- [ ] Tests use a fake ACP Runtime and prove no git or Review Source mutation is
      attempted by this slice.

## Blocked by

- 07-resolve-deduplication-and-batches.md
