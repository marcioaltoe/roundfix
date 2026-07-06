---
task: task_03
spec: 0013-codex-runtime-hygiene
status: completed
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

- [x] `doctor` command parsing and dispatch
- [x] Run shared checks + codex hygiene and format one line per check
- [x] Exit-code policy (fail on Run-breaking failure; skip is not failure)
- [x] Tests with a fake checker: pass set, failing set, codex-not-applicable

## Acceptance Criteria

- [x] `roundfix doctor` prints one line per check with status and next action.
- [x] With a quarantined codex on macOS, the codex line fails and the command exits non-zero.
- [x] On Linux/Windows, the codex line reports not-applicable and the command does not fail on that account.
- [x] The command mutates nothing on disk.

## Verification

- `rtk go test ./internal/cli/ -run Doctor` — expected: the Doctor command tests pass.
- `rtk go run ./cmd/roundfix doctor --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-3, 5; Core Feature 1. `_techspec.md` → Doctor Command,
Build Order 3. CONTEXT.md → Doctor Command. ADR-0032.

## Result

Implemented `roundfix doctor` as a read-only support command. It loads the
configured Agent runtime, runs the shared `node`, `acpx`, `agent`, and `codex`
checks, prints one stdout line per check, appends `next: ...` for failed checks
with remediation, and exits non-zero only when a check returns `failed`.

Evidence:

- One line per check: `TestRunDoctorReportsReadinessChecks/all_checks_pass`
  passed and asserts exact `node`, `acpx`, `agent`, and `codex` stdout lines.
- Quarantined codex failure: `TestRunDoctorReportsReadinessChecks/quarantined_codex_fails_with_reinstall_action`
  passed and asserts the `codex: failed` line includes the reinstall next
  action and exits `exitRunFailed`.
- Non-darwin/not-applicable behavior:
  `TestRunDoctorReportsReadinessChecks/codex_not_applicable_does_not_fail`
  passed and asserts `codex: skipped (not-applicable on linux)` exits `exitOK`.
- No disk mutation: the Doctor command tests run in a temp workspace with a fake
  checker and assert no `.acpx`, `.roundfix`, or `.roundfixrc.yml` paths are
  created.
- Required focused gate: `rtk go test ./internal/cli/ -run Doctor` passed with
  5 tests in 1 package.
- Required help check: `rtk go run ./cmd/roundfix doctor --help` exited 0 and
  printed concise usage for the read-only Doctor command.
- Required full suite: `rtk go test ./...` passed with 780 tests in 18 packages.
- Repository gate: `rtk make verify` passed, including full tests, Roundfix
  skill check, and build.
