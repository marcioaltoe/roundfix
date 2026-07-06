---
task: task_01
spec: 0013-codex-runtime-hygiene
status: pending
type: backend
complexity: medium
---

# Task 01: Extract shared read-only health checks from the Setup Command

## Overview

Prepare the ground for the Doctor Command by extracting the Setup Command's
read-only checks (Node, pinned acpx, configured Agent probe) into a shared unit
both Setup and Doctor can consume. Setup keeps its preparing/mutating steps and
its exact reported output; this is a prefactor with no user-visible change.

## Requirements

1. MUST extract the Node, acpx-version, and Agent-probe checks into a shared
   health-check unit returning a structured per-check result (name, status,
   detail, next action).
2. MUST have the Setup Command consume the shared checks with byte-stable
   reported output and unchanged mutating steps (install acpx, write config).
3. MUST NOT add a new command or user-facing behavior in this task.
4. MUST keep the Setup Command's existing tests passing unmodified except where
   they legitimately move to cover the shared unit.

## Subtasks

- [ ] Define the shared check result type and checker seam
- [ ] Move Node/acpx/agent-probe checks into the shared unit
- [ ] Rewire Setup to consume the shared checks, output byte-stable
- [ ] Confirm Setup tests still pass; add shared-unit coverage

## Acceptance Criteria

- [ ] The Node, acpx, and agent-probe checks live in a shared unit returning structured results.
- [ ] `roundfix setup` output and mutating behavior are unchanged (existing tests pass).
- [ ] No new command or flag is introduced by this task.

## Verification

- `rtk go test ./internal/cli/ -run Setup` — expected: existing Setup tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 3 (foundation). `_techspec.md` → System Architecture,
Build Order 1, Interfaces: `HealthChecker`. CONTEXT.md → Setup Command, Doctor
Command.
