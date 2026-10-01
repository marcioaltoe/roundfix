---
status: approved
granted: 2026-10-01
action: make six intermittently failing tests deterministic by changing only tests and test helpers — event-bound waits, a fake-clock Run Budget, script fixtures that re-execute the compiled test binary with a guard on written executables, fixture processes that end with their test binary, and an in-memory Assets Sync template
consuming: 0213-a-test-suite-that-does-not-flake
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0213

On 2026-10-01 the maintainer authorized the next cycle of unattended work with
"Tudo, de A a F". This Spec is item B of that cycle: a test suite that does not
flake. The authority extends the broad autonomy for authoring, corrective work
and delivery through merge that the maintainer granted on 2026-09-29.

## What is not governed

The set was measured with `GovernedPath` on the authoring branch at
`30cd147c`, through a `go test -overlay` probe that wrote nothing to the
repository. These paths are ordinary:

- `internal/store/journal_batch_test.go`
- `internal/daemon/task_budget_renewal_test.go`
- `internal/cli/adapter_floor_test.go`
- `internal/cli/implement_test.go`
- `internal/cli/carryforward_hooks_test.go`
- `internal/cli/settle_test.go`
- `internal/cli/detach_test.go`
- the new `internal/cli/script_fixture_test.go`
- the new `internal/cli/implement_detach_teardown_test.go`
- the new `internal/testfixture/written_executable_test.go`
- `internal/baseline/assets_sync_test.go`
- this Spec's directory
- ADR-0213

The governed `internal/cli/cli_test.go` is not touched. Its written executable
stays in the guard's residual inventory.

## Limits

- No production code change, and no change to a product timeout, command
  output or exit code.
- No edit to the Makefile, to lint, formatter or test-runner configuration,
  to a CI workflow or to `go.mod`.
- No retry, sleep or widened deadline in place of a fix.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No process is killed except by a PID the test itself started or recorded.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
