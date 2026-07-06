---
task: task_04
spec: 0013-codex-runtime-hygiene
status: pending
type: backend
complexity: medium
---

# Task 04: Verified-clean codex on the codex-acp spawn path

## Overview

Stop Runs from spawning a quarantine-blocked codex: when Roundfix launches codex
through `codex-acp`, resolve a verified-clean binary (respecting a configured
codex path) so an agent loop's repeated execs no longer intermittently trip
XProtect. Reuses the task_02 inspector; never silently spawns a known-quarantined
binary without surfacing the risk.

## Requirements

1. MUST resolve the codex for a codex-acp launch through the configured codex
   path then PATH, and on macOS prefer a binary that passes the hygiene
   inspection.
2. MUST NOT silently spawn a known-quarantined codex — the risk MUST be surfaced
   when no clean binary is available.
3. MUST inspect once per Run's codex resolution, not per exec, to avoid adding
   latency to every codex call.
4. MUST leave non-macOS spawning behavior unchanged.

## Subtasks

- [ ] Resolve codex path for codex-acp launch (config then PATH)
- [ ] On macOS, prefer a hygiene-passing binary via the task_02 inspector
- [ ] Surface the risk when only a quarantined binary is available
- [ ] Cache the resolution per Run; tests for clean-over-quarantined selection

## Acceptance Criteria

- [ ] A codex-acp launch resolves the configured clean codex over a quarantined PATH entry.
- [ ] When only a quarantined codex is available, the risk is surfaced rather than silently spawned.
- [ ] Codex hygiene is inspected once per Run resolution, not on every exec.
- [ ] Non-macOS spawn behavior is unchanged.

## Verification

- `rtk go test ./internal/agent/ -run Codex` — expected: spawn-resolution tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Verified-clean codex
on spawn, Build Order 4. ADR-0032.
