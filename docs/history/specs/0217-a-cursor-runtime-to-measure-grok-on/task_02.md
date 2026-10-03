---
task: task_02
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: completed
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

## Result

Implemented the assigned login gate for Daemon Verification. Task status and
acceptance checkboxes remain unsettled; the Daemon owns the terminal verdict.

- `checkAdapter` checks Cursor (including the normalized `cursor-custom`
  runtime) after executable lookup and before lineage inspection. The status
  subprocess receives the explicit environment and is cancelled by the caller's
  context or a ten-second ceiling. A failed command or `Not logged in` on either
  output stream returns `CursorLoginRequiredError` unchanged, with no status
  output in the error or adapter evidence.
- The health checker reads `NextAction()` through the error chain. Doctor prints
  the classified refusal and the maintainer's terminal login action. Both the
  Doctor and profiles guides describe refusal before a session opens.
- Every fake Cursor executable is provisioned by the agent suite's compiled
  binary re-execution, with behavior in a non-executable sidecar. The fixture
  now supports status stdout, stderr, exit code, environment and blocking.

Acceptance evidence:

1. `TestCursorAdapterRefusedWithoutLogin` covers the refusal message on stdout
   and stderr and exit code 1. `TestCursorProofRefusesBeforeAnySessionWithoutLogin`
   covers each case and asserts that no acpx invocation file exists.
   `TestDoctorAdapterNamesTheCursorLogin` drives the Doctor command through the
   real health/adapter path with compiled fixtures for message and exit refusal;
   it asserts exit 1, the exact Surface Transcript 1 adapter line, empty stderr
   and absence of a private-output sentinel.
2. `TestCursorAdapterReadyWithLogin` accepts the logged-in fixture and returns
   only command evidence. `TestCursorLoginUsesExplicitEnvironment`,
   `TestCursorLoginHonorsCallerDeadline` and
   `TestCursorCustomAdapterAlsoRequiresLogin` additionally cover environment,
   cancellation and normalized runtime handling.
3. A Python inspection of all non-test Go files under `internal` and `cmd`
   found no Cursor credential variable, key/token flag, or login/logout
   argument references. Direct inspection of the production gate confirms the
   only Cursor subprocess argument is `status`; no status output is logged or
   persisted and no `-H` is passed to Cursor.

Focused checks:

- Red starting point: `GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/agent -run '^TestCursorAdapterReadyWithLogin$' -count=1` failed to
  build before implementation because the required error, classification and
  login-check symbols did not exist.
- Final focused check: `GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/agent ./internal/cli -run
  'TestCursor(Adapter|Proof|Login|Custom)|TestDoctorAdapterNamesTheCursorLogin|TestRunDoctorAdapterReadiness|TestDoctorAdapterCheckAggregatesEveryFailure|TestFakeAdapterRuns'
  -count=1` exited 0 for both packages. This includes existing Doctor adapter
  aggregation and compiled-fixture provisioning regressions.
- `git -c core.fsmonitor=false diff --check` exited 0. Source and guide diffs
  were inspected against the assigned scope.

The authored Verification command and repository-wide Verification were not
run in this child turn. No live Cursor command, login, credential access,
commit, push or pull request was performed. No follow-up scope was added.

## Carry-forward provenance

- Source Run: `run_20261002T190644Z_6dc4ccb1092ddf45`
- Source commit: `fe08a97db1bcefa4ecb23dd98caadf1ed49015de`
