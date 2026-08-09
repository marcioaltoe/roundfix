---
task: task_01
spec: 0083-a-gate-that-can-say-no
status: completed
type: infra
complexity: high
---

# Task 01: Take the wrapper off the gate and prove the gate can fail

## Overview

The authoritative verification invocation runs through an output-filtering
wrapper, and on the full package set that wrapper was observed returning success
for a run the toolchain failed. This task moves the wrapper off every path whose
exit status is load-bearing and lands the regression test that makes the repair
provable rather than believed. It is verifiable on its own: after it, a
deliberately failing test propagates a non-zero exit through the gate's own
command shape.

## Requirements

1. MUST split the Makefile's toolchain variable so the authoritative invocation
   calls the toolchain directly, and a separate variable keeps the wrapper for
   targets whose output — not exit status — is the result.
2. MUST route every target the verification gate composes through the
   unfiltered variable.
3. MUST keep the wrapper available for human-facing convenience targets; this
   task removes it from the gate, not from the repository.
4. MUST add a regression test that runs the gate's own command shape against a
   throwaway module containing a test that fails while emitting output at the
   volume that concealed the real failure, and requires a non-zero exit.
5. MUST build that throwaway module in a temporary directory so the test never
   depends on this repository being red, and MUST NOT shell out to `make`, which
   would make the test depend on the developer's `make` and on repository state.
6. MUST leave the gate's name, composition, and meaning unchanged; it gains only
   the ability to fail.
7. MUST change only these repository-relative paths plus this Task file:
   `Makefile` and `internal/spec/gate_test.go`, which is new and is where the
   regression test lives. Any other changed path fails this Task.
8. MUST NOT weaken, retire, or reword any existing gate in this Task.

## Rehearsal Cases

- Case: a test that fails while emitting output at the volume that concealed the 2026-08-07 failure; Observation: the authoritative invocation exits non-zero and names the failing package.
- Case: a test that fails with short output; Observation: the authoritative invocation exits non-zero.
- Case: a temporary module whose tests all pass; Observation: the authoritative invocation exits zero, proving the repair did not make the gate always fail.
- Case: the gate's command shape routed back through the output-filtering wrapper; Observation: the regression test fails, proving it detects the defect rather than the symptom.

## Subtasks

- [x] Split the toolchain variable into an authoritative and a human-facing form.
- [x] Point every gate-composed target at the authoritative variable.
- [x] Add the regression test with its temporary failing module.
- [x] Confirm the wrapper still serves the convenience targets.
- [x] Confirm the changed-file set matches the declared boundary.

## Acceptance Criteria

- [x] The authoritative toolchain variable does not reference the wrapper.
- [x] Every target the verification gate composes resolves to the authoritative
      variable.
- [x] The regression test fails if the gate's command shape is changed back to
      route through the wrapper.
- [x] The regression test passes on a machine where the repository suite is red,
      proving it depends on its own temporary module rather than on this tree.
- [x] At least one convenience target still uses the wrapper.
- [x] The verification gate's target list and name are unchanged.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- instruction: `docs/findings/2026-08-07-the-only-gate-reports-green-on-a-red-suite.md`
- interface: `Makefile`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `grep -E '^GO[[:space:]]*:?=' Makefile | grep 'RTK' ; test $? -eq 1` — expected: exits 0, proving the authoritative toolchain variable no longer routes through the wrapper.
- `grep -q 'GO_HUMAN' Makefile` — expected: exits 0, proving the convenience wrapper variable is still declared rather than deleted.
- `make -n verify > /tmp/task_01-2.log 2>&1 && grep -q -F '$(GO) ' /tmp/task_01-2.log` — expected: exits 0, proving every verify prerequisite resolves to the authoritative toolchain variable without the wrapper.
- `go test ./internal/spec -run '^TestAuthoritativeGateReportsFailure$' -count=1 -v > /tmp/task_01-1.log 2>&1 && grep -q '^--- PASS: TestAuthoritativeGateReportsFailure' /tmp/task_01-1.log` — expected: exits 0, proving the regression test exists at its declared home and passes rather than being selected out.
- `(git diff --name-only HEAD; git ls-files --others --exclude-standard) | grep -v -E '^(Makefile|internal/spec/gate_test\.go|docs/specs/0083-a-gate-that-can-say-no/task_01\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 1 and 2; System Architecture: the Makefile's tool variables; Interfaces.
- `_prd.md` → Core Features 1 and 2; Goals 1.
- ADR-0081.

## Result

### Implementation

- Split the Makefile toolchain into direct authoritative `GO` and `GOFMT`
  variables plus the wrapped `GO_HUMAN` variable. The human-facing `version`
  target retains the wrapper.
- Added `TestAuthoritativeGateReportsFailure`. It reads and expands the
  Makefile's `test` recipe without invoking `make`, runs that command in
  temporary Go modules, and covers a 300-line failure, a short failure, and a
  passing package. A temporary masking wrapper exercises the regression
  mutation without depending on RTK internals.

### Focused checks

- Before the Makefile split,
  `rtk proxy env GOCACHE=<worktree>/.gocache go test ./internal/spec -run '^TestAuthoritativeGateReportsFailure$'`
  failed: both failing fixtures exited zero and the direct `GO := go`
  assignment was absent. The first attempt without the explicit repository
  cache stopped at the sandboxed system Go cache with `operation not
  permitted`; the rerun above reached the intended red assertion.
- After the split, the same focused test exited zero with
  `ok roundfix/internal/spec`, including the high-volume failure, short
  failure, passing-package, and masking-wrapper subtests.
- With `GO` temporarily mutated back to `$(RTK) go`, the focused regression
  exited non-zero because the two failing fixtures returned zero; restoring
  `GO := go` made the focused regression exit zero again.
- `rtk proxy env GOCACHE=<worktree>/.gocache go test ./internal/spec -run '^TestCoverageEquivalence$'`
  remained red on the repository-owned coverage regression, while the isolated
  gate regression immediately afterward exited zero. The new test therefore
  reads its verdict from its temporary modules rather than this tree.
- `rtk make -n verify` exited zero and expanded the unchanged
  `verify: fmt-check test spec-budget skills-sync-check skills-check build spec-check`
  composition to direct `gofmt`, `go test`, `go run`, and `go build` commands.
- `rtk gofmt -l internal/spec/gate_test.go` produced no paths, and
  `rtk git diff --check` exited zero.
- `rtk git -c core.fsmonitor=false status --short --untracked-files=all`
  listed only `Makefile`, this Task file, and `internal/spec/gate_test.go`.

### Acceptance evidence

1. `Makefile` assigns `GO := go` and `GOFMT := gofmt`; neither authoritative
   variable references `RTK`.
2. The `verify` dry run resolves every composed Go and formatting invocation
   through those direct variables.
3. The temporary mutation probe made `TestAuthoritativeGateReportsFailure`
   fail when `GO` was routed through the masking wrapper.
4. The repository coverage test was independently red while the new test
   passed against its temporary modules.
5. `version` resolves through `GO_HUMAN := $(RTK) go`, so the wrapper remains
   available on a convenience target.
6. The Makefile diff leaves the `verify` declaration byte-identical; its name
   and prerequisite list are unchanged.

### Daemon verification

Not run in this Agent turn. The Daemon owns every command in `## Verification`
and the Task's terminal status.
