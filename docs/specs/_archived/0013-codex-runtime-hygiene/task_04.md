---
task: task_04
spec: 0013-codex-runtime-hygiene
status: completed
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

- [x] Resolve codex path for codex-acp launch (config then PATH)
- [x] On macOS, prefer a hygiene-passing binary via the task_02 inspector
- [x] Surface the risk when only a quarantined binary is available
- [x] Cache the resolution per Run; tests for clean-over-quarantined selection

## Acceptance Criteria

- [x] A codex-acp launch resolves the configured clean codex over a quarantined PATH entry.
- [x] When only a quarantined codex is available, the risk is surfaced rather than silently spawned.
- [x] Codex hygiene is inspected once per Run resolution, not on every exec.
- [x] Non-macOS spawn behavior is unchanged.

## Verification

- `rtk go test ./internal/agent/ -run Codex` — expected: spawn-resolution tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Verified-clean codex
on spawn, Build Order 4. ADR-0032.

## Result

Implemented verified-clean codex resolution on the codex-acp Run spawn path.
On macOS, Roundfix now resolves `CODEX_PATH` first and PATH second, inspects
candidate binaries with the task_02 codex hygiene inspector, passes the selected
clean binary to acpx through `CODEX_PATH`, and fails before spawning acpx when
no clean codex is available. Non-Darwin runs return the existing acpx arguments
without running codex probes.

Evidence:

- Configured clean codex preference:
  `TestACPXRunCodexUsesConfiguredCleanPathOnDarwin` passed and asserts acpx
  receives the configured clean `CODEX_PATH` while a quarantined PATH entry is
  available.
- Clean PATH fallback:
  `TestACPXRunCodexFallsBackToCleanPathWhenConfiguredPathIsQuarantined` passed
  and asserts a quarantined configured path falls back to the clean PATH codex.
- Quarantined-only failure:
  `TestACPXRunCodexSurfacesQuarantinedPathWithoutSpawning` passed and asserts
  the error names the quarantine risk and reinstall next action, with no acpx
  invocation file created.
- Once-per-Run resolution:
  `TestACPXRunCodexInspectsOncePerSessionResolution` passed and asserts two
  acpx execs in the same session use one quarantine probe and one acceptance
  probe.
- Non-Darwin behavior:
  `TestACPXRunCodexLeavesNonDarwinSpawnUnchanged` passed and asserts no codex
  probes run and the acpx invocation arguments remain unchanged.
- Required focused gate: `rtk go test ./internal/agent/ -run Codex` passed with
  8 tests in 1 package.
- Required full suite: `rtk go test ./...` passed with 785 tests in 18 packages.
- Repository gate: `rtk make verify` passed, including full tests, Roundfix
  skill check, and build.
