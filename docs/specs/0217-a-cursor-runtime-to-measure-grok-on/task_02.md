---
task: task_02
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: pending
type: backend
complexity: medium
---

# Task 02: A Cursor selection without the maintainer's login is refused, never logged in

## Overview

`cursor-agent` authenticates with the maintainer's own Cursor login, and on
2026-10-01 it was not logged in: `status` printed `Not logged in` and
`session/new` was refused. Every proof, Run, sealed session and Doctor line
passes through `checkAdapter`. This Task makes `checkAdapter` ask
`cursor-agent status` for runtime `cursor` and refuse with
`cursor_login_required` and a next action when there is no login, and makes
the Doctor `adapter:` line print that next action. Roundfix never runs a
login and never touches a Cursor credential (ADR-0217).

## Requirements

1. MUST add `internal/agent/cursor_login.go` with `CursorLoginRequired`,
   `CursorLoginRequiredError` (`Error`, `Classification`, `NextAction`) and
   `checkCursorLogin`, as the TechSpec Interfaces state. `checkCursorLogin`
   MUST run `<executable> status` with the explicit environment, bounded by
   the caller's context and 10 seconds, and MUST return the error when the
   command fails or its combined output contains `Not logged in`. It MUST
   keep no part of the output.
2. MUST call `checkCursorLogin` from `checkAdapter` for runtime `cursor`
   (without `-custom`), after the executable lookup and before any lineage
   contract, so `ProveExactSelection`, the Run's prompt path, sealed sessions
   and Doctor all receive the error unchanged.
3. MUST make `runtimeHealthChecker.Adapter` in `internal/cli/health.go` take
   the next action from an error's `NextAction()` when it has one, as it takes
   an install command from `InstallCommand()`, so the `adapter:` line prints
   Surface Transcript 1 of the TechSpec.
4. MUST provide every fake `cursor-agent` through the compiled test binary
   the `internal/agent` fake adapter provisioning re-executes (ADR-0125),
   extended with a `status` answer and exit code, never as a written script.
5. MUST add to `docs/user-guide/commands/doctor.md`, in the `adapter:`
   bullet, that a Cursor runtime is checked with `cursor-agent status` and
   that no login fails the line with `cursor_login_required` and
   `next: run cursor-agent login in a terminal yourself`; and to
   `docs/user-guide/commands/profiles.md` that `profiles validate` and
   `profiles configure` refuse a Cursor tuple without the login, as
   `cursor_login_required`, before any session opens.
6. MUST NOT run `login`, `logout`, `--api-key`, `--auth-token` or `-H` on
   `cursor-agent`, read `CURSOR_API_KEY` or `CURSOR_AUTH_TOKEN`, or print,
   log or store the status output.

## Subtasks

- [ ] Add the login check and its classified error.
- [ ] Call it from the adapter check and carry its next action to Doctor.
- [ ] Prove refusal, readiness, refusal before any session, and the Doctor line.
- [ ] Describe the check in the Doctor and profiles guides.

## Acceptance Criteria

- [ ] With a fake `cursor-agent` whose `status` prints `Not logged in`, and
      with one whose `status` exits `1`, the adapter check returns
      `CursorLoginRequiredError`, proof refuses before any acpx call, and
      Doctor prints the classified line with
      `next: run cursor-agent login in a terminal yourself`.
- [ ] With a fake whose `status` reports a login, the adapter check passes.
- [ ] No production Go file names a Cursor credential variable, a key or
      token flag, or a `login` argument.

## Context

- creates: `internal/agent/cursor_login.go`
- creates: `internal/agent/cursor_login_test.go`
- creates: `internal/cli/doctor_cursor_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/cli/health.go`
- instruction: `internal/cli/doctor.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/profiles.md`
- instruction: `internal/cli/doctor_test.go`
- instruction: `docs/adr/0125-a-spawned-fixture-is-a-compiled-binary-never-a-written-script.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCursorAdapterRefusedWithoutLogin|TestCursorAdapterReadyWithLogin|TestCursorProofRefusesBeforeAnySessionWithoutLogin|TestDoctorAdapterNamesTheCursorLogin)$' ./internal/agent ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCursorAdapterRefusedWithoutLogin TestCursorAdapterReadyWithLogin TestCursorProofRefusesBeforeAnySessionWithoutLogin TestDoctorAdapterNamesTheCursorLogin; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in "docs/user-guide/commands/doctor.md|cursor_login_required" "docs/user-guide/commands/doctor.md|run cursor-agent login in a terminal yourself" "docs/user-guide/commands/profiles.md|cursor_login_required"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; if grep -rnE --include='*.go' -e 'CURSOR_API_KEY|CURSOR_AUTH_TOKEN|--api-key|--auth-token|[^:]"(login|logout)"' internal cmd | grep -v '_test.go:'; then printf 'a production file handles a Cursor credential or login\n' >&2; exit 1; fi` — expected: exit 0; before this Task the named tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Stories 2-3; Core Feature 3; Success Metric 3; Project Constraints (Authentication and HTTP)
- [_techspec.md](_techspec.md) — Interfaces; The login check; Surface Transcript 1; API Contract 3; Testing Approach 3; Build Order 2
- ADR-0217; ADR-0050; ADR-0125
