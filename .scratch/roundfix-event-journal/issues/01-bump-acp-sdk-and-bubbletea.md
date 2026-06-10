---
title: "Deps: bump acp-go-sdk to v0.13.5 and bubbletea to v2.0.7"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 11
blocked_by: []
---

# Deps: bump acp-go-sdk to v0.13.5 and bubbletea to v2.0.7

## Parent

.scratch/roundfix-event-journal/01-run-event-seam-prd.md

## What to build

First task of the Run Event seam PRD, isolated so SDK breakage cannot entangle seam work: upgrade `coder/acp-go-sdk`
v0.6.3 → v0.13.5 and `bubbletea` v2.0.2 → v2.0.7 (plus compatible bubbles/lipgloss patch releases). Behavior must be
identical after the bump. The ACP gap includes stabilized `session/resume`, `session/close`, and
`session_info_update`; resolve any breaking API changes inside this slice.

## Acceptance criteria

- [ ] `coder/acp-go-sdk` at v0.13.5; bubbletea at v2.0.7 with compatible companions
- [ ] All existing tests pass with no behavior change
- [ ] `rtk go test ./...` and `rtk go run ./cmd/roundfix --help` green
- [ ] Single isolated commit for the bumps

## Blocked by

None - can start immediately
