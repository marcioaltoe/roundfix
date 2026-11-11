---
task: task_06
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
type: backend
complexity: low
---

# Task 06: The Linux Lane builds the budget entries from the corebudget names

## Overview

Corrective Task 1 of 2. The QA gate of 2026-11-11
(`qa/qa-report-2026-11-11.md`, finding F01, row D01) found that
`internal/cli/deliver_linux_lane.go` appends the two budget entries to the
Docker argument list with the literal names `GOMAXPROCS` and
`ROUNDFIX_VERIFICATION_CORES`, while task_05 Requirement 6 requires every
non-test file outside `internal/corebudget` to take those names from that
package. The behavior is right; the names are duplicated.

## Requirements

1. MUST answer finding F01 of the QA Report of 2026-11-11: build the two
   `-e` entries of the Linux Lane from `corebudget.GoVariable` and
   `corebudget.Variable`, keeping their order (`GOMAXPROCS` first), their
   position in the Docker argument list and the condition that writes them
   only when the budget is above zero.
2. MUST leave no non-test Go file outside `internal/corebudget` that writes
   either name as a string literal followed by `=`.
3. MUST keep every Linux Lane test passing, with no test changed.

## Subtasks

- [x] Replace the two literals with the corebudget constants.
- [x] Run the Linux Lane tests.

## Acceptance Criteria

- [x] The Linux Lane takes both names from `internal/corebudget`, the sweep finds no literal outside it, and the Linux Lane tests pass.

## Context

- interface: `internal/cli/deliver_linux_lane.go`
- instruction: `internal/corebudget/corebudget.go`

## Consumers

- The Linux Lane Docker argument list — consumers: `./internal/cli`

## Verification

- `if grep -rnE '"(GOMAXPROCS|ROUNDFIX_VERIFICATION_CORES)=' --include='*.go' internal cmd | grep -v '_test[.]go:' | grep -v '^internal/corebudget/'; then echo 'a literal budget variable name remains outside internal/corebudget' >&2; exit 1; fi; go test -count=1 -timeout 20m -run '^(TestLinuxLane|TestGateWithoutTheLane|TestGateSkipsTheLane)' ./internal/cli` — expected: exit 0. Before this Task the sweep prints the line of `internal/cli/deliver_linux_lane.go` that appends both literals and exits 1; afterwards it prints nothing and the Linux Lane tests pass.

## References

- `qa/qa-report-2026-11-11.md` → F01, row D01
- task_05 → Requirement 6
- `_techspec.md` → Build Order

## Result

F01 (QA Report of 2026-11-11, row D01) answered. The Linux Lane's two budget
`-e` entries are now built from `corebudget.GoVariable` and
`corebudget.Variable` instead of the duplicated string literals; behavior,
order (`GOMAXPROCS` first), position in the Docker argument list and the
`request.Cores > 0` guard are unchanged. Only `internal/cli/deliver_linux_lane.go`
changed: one import line and one `append` line.

Implementation:
- `internal/cli/deliver_linux_lane.go` — added the `roundfix/internal/corebudget`
  import and replaced `"GOMAXPROCS="+cores` / `"ROUNDFIX_VERIFICATION_CORES="+cores`
  with `corebudget.GoVariable+"="+cores` / `corebudget.Variable+"="+cores`.

Focused checks (this session, after the last edit):
- Pre-change signal: `grep -rnE '"(GOMAXPROCS|ROUNDFIX_VERIFICATION_CORES)=' --include='*.go' internal cmd` (excluding `_test.go` and `internal/corebudget/`) printed
  `internal/cli/deliver_linux_lane.go:180`, and exited 0 — the duplicated literal F01 names.
- Post-change sweep: the same filter now prints nothing and exits 1 (grep found
  no match), so no non-test Go file outside `internal/corebudget` writes either
  name followed by `=` (Requirement 2).
- `gofmt -l internal/cli/deliver_linux_lane.go` — no output (formatted).
- `go build ./internal/cli` — exits 0.
- `go test -count=1 -timeout 20m -run 'LinuxLane' ./internal/cli` — `ok roundfix/internal/cli 5.321s`, all Linux Lane tests PASS. `TestLinuxLaneRunsTheChangedPackagesInAContainer` asserts the full Docker argument list including `-e GOMAXPROCS=4 -e ROUNDFIX_VERIFICATION_CORES=4` at their original position, which covers Requirement 1's names, order, position and positive-budget condition (Requirement 3: no test changed).

Acceptance evidence: the Linux Lane now takes both names from
`internal/corebudget` (import + symbol references above), the literal sweep
finds nothing outside `internal/corebudget`, and the Linux Lane tests pass with
no test modified (`git status --short` shows only `internal/cli/deliver_linux_lane.go`
and this task file).

Not run (Daemon-owned): the declared `## Verification` command and the full
suite; the Daemon runs them at settlement.

Follow-ups: none.
