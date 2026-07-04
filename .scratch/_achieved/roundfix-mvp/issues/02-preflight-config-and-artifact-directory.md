---
title: "Preflight: validate config and Artifact Directory"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 7
  - 8
  - 9
  - 10
blocked_by:
  - 01-cli-command-contract-and-exit-codes.md
---

# Preflight: validate config and Artifact Directory

## Parent

.scratch/roundfix-mvp/PRD.md

## What to build

Add YAML User Config and Project Config loading to the command path, apply the
documented precedence rules, and validate the configured Artifact Directory
before any Review Source or Agent work can start. The behavior should be visible
through MVP commands and testable without live GitHub access.

## Acceptance criteria

- [x] Built-in defaults, User Config, Project Config, and CLI flags apply in the
      documented order.
- [x] Invalid YAML, invalid config values, and invalid duration values fail
      during Preflight Validation with exit `2` and concrete diagnostics.
- [x] Empty, absolute, relative, and home-relative Artifact Directory values
      resolve according to the product contract.
- [x] The Artifact Directory is created when missing, rejected when it is not a
      directory, and checked for writability before fetching.
- [x] Artifact Directory validation failures do not create Run records, markdown
      artifacts, Agent processes, commits, pushes, or Review Source calls.

## Blocked by

- 01-cli-command-contract-and-exit-codes.md
