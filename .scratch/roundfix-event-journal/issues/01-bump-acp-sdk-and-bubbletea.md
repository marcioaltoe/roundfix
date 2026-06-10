---
title: "Deps: bump acp-go-sdk to v0.13.5 and bubbletea to v2.0.7"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] `coder/acp-go-sdk` at v0.13.5; bubbletea at v2.0.7 with compatible companions
- [x] All existing tests pass with no behavior change
- [x] `rtk go test ./...` and `rtk go run ./cmd/roundfix --help` green
- [x] Single isolated commit for the bumps

## Blocked by

None - can start immediately

## Comments

**2026-06-10 (agent):** Bumped `coder/acp-go-sdk` v0.6.3 → v0.13.5 and `charm.land/bubbletea/v2` v2.0.2 → v2.0.7,
with companion patches (lipgloss v2.0.3, x/ansi v0.11.7, colorprofile v0.4.3, ultraviolet, go-runewidth,
golang.org/x/sys). `go mod tidy` dropped `charm.land/bubbles/v2` — it was an unused indirect dependency; nothing
imports it. Resolved three SDK breaking changes in `internal/agent/acp_runner.go`: `FileSystemCapability` →
`FileSystemCapabilities`; client method `KillTerminalCommand` → `KillTerminal` (request/response types renamed the
same way); removed `session/set_model` migrated to `session/set_config_option` against the model-category select
option advertised in the NewSession response (new `setSessionModel`/`modelConfigOption` helpers, same error
surface). Verification: `rtk go test ./...` 148 passed in 14 packages (same as pre-bump baseline), `rtk go vet ./...`
clean, `rtk go run ./cmd/roundfix --help` green.
