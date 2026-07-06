---
task: task_03
spec: 0013-codex-runtime-hygiene
status: pending
type: backend
complexity: medium
---

# Task 03: Doctor Command

## Overview

Add `roundfix doctor`, the diagnosis-only support command that reports a
machine's readiness for Runs by running the shared health checks (Node, acpx,
Agent probe) plus the codex hygiene check, one line each with a next action, and
exits non-zero when a Run-breaking check fails. It mutates nothing — the
read-only sibling of the Setup Command.

## Requirements

1. MUST add a non-interactive `doctor` support command that runs every shared
   check plus the codex hygiene inspector.
2. MUST print one line per check (`node`, `acpx`, `agent`, `codex`) with its
   status and, on failure, the next action, writing requested output to stdout
   and diagnostics to stderr per the CLI contract.
3. MUST exit non-zero when any Run-breaking check fails and zero when all pass or
   are skipped; a not-applicable codex check MUST NOT cause a non-zero exit.
4. MUST perform no installs, config writes, downloads, or any other mutation.

## Subtasks

- [ ] `doctor` command parsing and dispatch
- [ ] Run shared checks + codex hygiene and format one line per check
- [ ] Exit-code policy (fail on Run-breaking failure; skip is not failure)
- [ ] Tests with a fake checker: pass set, failing set, codex-not-applicable

## Acceptance Criteria

- [ ] `roundfix doctor` prints one line per check with status and next action.
- [ ] With a quarantined codex on macOS, the codex line fails and the command exits non-zero.
- [ ] On Linux/Windows, the codex line reports not-applicable and the command does not fail on that account.
- [ ] The command mutates nothing on disk.

## Verification

- `rtk go test ./internal/cli/ -run Doctor` — expected: the Doctor command tests pass.
- `rtk go run ./cmd/roundfix doctor --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-3, 5; Core Feature 1. `_techspec.md` → Doctor Command,
Build Order 3. CONTEXT.md → Doctor Command. ADR-0032.
