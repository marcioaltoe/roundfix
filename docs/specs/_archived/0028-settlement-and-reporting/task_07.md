---
task: task_07
spec: 0028-settlement-and-reporting
status: completed
type: backend
complexity: medium
---

# Task 07: Diagnose missing adapter binaries at preflight

## Overview

Turn the opaque adapter spawn failure into an actionable preflight diagnosis: Roundfix resolves the adapter binary acpx will spawn for the selected ACP Runtime (from acpx's own configuration, falling back to built-in defaults), checks it exists on PATH during the probe, and fails with the binary name and its install command — mirroring the existing acpx-missing message. The Doctor Command reports the same check. Field-proven need: a missing/broken adapter today surfaces only as a raw acpx stderr tail after model probing.

## Requirements

1. MUST resolve the adapter command for the selected runtime from the acpx configuration's agents map when present and readable, falling back to the built-in defaults (codex → `codex-acp`, claude → `claude-code-acp`, opencode → `opencode`); stdio command overrides resolve to the override.
2. MUST check the resolved adapter binary during the existing probe, before the disposable Agent Session, and fail preflight with `<adapter> is required but was not found on PATH; install it with: <command>` using a built-in install-hint map for the known adapters and a generic hint otherwise.
3. MUST degrade silently to the built-in defaults on absent or malformed acpx configuration — never fail preflight because the config file itself is unreadable.
4. MUST add one `adapter: ok|failed` line for the configured agent to the Doctor Command, with the same next action on failure.

## Subtasks

- [ ] Adapter-command resolution with acpx-config parsing and built-in fallbacks
- [ ] Probe extension with the install-hint error text
- [ ] Doctor line for the adapter check
- [ ] Tests: resolution against fixture config files (present, missing, malformed, override); probe failure text for a missing adapter; doctor line in both states

## Acceptance Criteria

- [ ] With the adapter binary absent from PATH, preflight fails naming the binary and its install command before any Agent Session is created
- [ ] A malformed acpx config falls back to defaults without failing preflight by itself
- [ ] `roundfix doctor` prints an adapter line with `ok` when present and `failed` plus a next action when absent
- [ ] The full test suite passes

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/cli/doctor.go`

## Verification

- `grep -q "resolveAdapterCommand" internal/agent/acpx_runner.go` — expected: exit 0
- `go test ./internal/agent/... ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 5, User Story 6, Core Feature 6; `_techspec.md` → Build Order 7, Interfaces (resolveAdapterCommand), Integration Points (acpx configuration), Risks (acpx config drift), Decisions (no parallel adapter registry).

## Result

- Added adapter-command resolution in `internal/agent/acpx_runner.go`: stdio overrides resolve first, then `~/.acpx/config.json` `agents.<runtime>.command`, then built-in defaults.
- Added the adapter PATH check before disposable Agent Session setup, with install hints for known adapters and a generic fallback hint for unknown commands.
- Added `adapter: ok|failed` Doctor output through the health checker; Doctor skips the Agent probe when the adapter check already failed.
- Acceptance evidence:
  - `TestACPXProbeMissingAdapterNamesInstallCommandBeforeSession` asserts a missing `codex-acp` fails with its install command after only the `acpx --version` probe.
  - `TestResolveAdapterCommandUsesConfigFallbacksAndOverrides` covers present config, missing config, malformed config, and stdio override resolution.
  - `TestACPXProbeMalformedConfigFallsBackToDefaultAdapter` proves malformed acpx config does not fail preflight by itself.
  - `TestRunDoctorReportsReadinessChecks` and `TestRunDoctorReportsAdapterFailureWithNextAction` cover Doctor `adapter: ok` and `adapter: failed` output with next action.
- Verification:
  - `rtk grep -q "resolveAdapterCommand" internal/agent/acpx_runner.go` — passed.
  - `rtk go test ./internal/agent/... ./internal/cli/...` — passed, 556 tests.
  - `rtk go build -buildvcs=false ./...` — passed.
  - `rtk make verify` — passed: `go test ./...` 1235 tests, skill check passed, build passed.
