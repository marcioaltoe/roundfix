---
task: task_10
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 10: Plan the release history reset

## Overview

Add a read-only Release Plan mode that inventories the complete historical
surface before any operator considers cleanup. The command must make approval
requirements explicit and expose no deletion capability.

## Requirements

1. MUST add `roundfix release plan --reset-to v0.0.1` with text and JSON output
   through the existing `Run() int` CLI boundary.
2. MUST make `--reset-to` mutually exclusive with `--from`, `--to`, `--impact`,
   and `--reason`.
3. MUST require a clean committed target and inventory every local and remote
   stable tag plus every GitHub Release through read-only provider interfaces.
4. MUST exhaust paginated GitHub Release results and fail closed when any
   local, remote, or GitHub inventory is unavailable or incomplete.
5. MUST sort inventory deterministically and compute a plan digest over target
   version, target revision, complete ordered tags, and complete ordered
   releases.
6. MUST return `approval_required` with exit code 3 and name every immutable
   identity and target commit available in both output formats.
7. MUST preserve stdout for requested plan output and stderr for diagnostics.
8. MUST NOT expose or invoke tag deletion, Release deletion, changelog writes,
   version writes, or any other mutation from the reset-plan interface.

## Subtasks

- [ ] Add reset-plan domain values and read-only inventory interfaces.
- [ ] Implement local/remote tag and paginated GitHub Release adapters.
- [ ] Add deterministic framing and digest calculation.
- [ ] Add CLI flag validation, text rendering, JSON rendering, and exit mapping.
- [ ] Add incomplete-inventory, dirty-target, pagination, and digest-sensitivity
      tests.
- [ ] Add mutation-spy tests proving the command remains read-only.

## Acceptance Criteria

- [ ] A complete fixture inventory produces deterministic text and JSON plans
      with the same digest and exit code 3.
- [ ] Every stable local tag, remote tag, and paginated GitHub Release appears
      exactly once with its immutable identity.
- [ ] A changed target revision, tag, or Release changes the plan digest.
- [ ] Dirty state, conflicting flags, malformed target versions, and incomplete
      inventory fail with the documented exit category and useful diagnostic.
- [ ] Captured provider calls contain only read operations; mutation spies are
      never invoked.
- [ ] The existing non-reset Release Plan behavior remains unchanged.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/adr/0048-release-plan-is-read-only-and-confirmation-gated.md`
- instruction: `docs/adr/0065-release-plan-exposes-a-read-only-reset-mode.md`
- interface: `internal/releaseplan/model.go`
- interface: `internal/releaseplan/build.go`
- interface: `internal/cli/releaseplan_command.go`
- interface: `internal/cli/releaseplan_git_source.go`

## Verification

- `rtk go test ./internal/releaseplan -run 'Reset|Digest|Inventory'` — expected: complete deterministic inventory, digest sensitivity, and read-only provider contracts pass.
- `rtk go test ./internal/cli -run 'ReleasePlan.*Reset'` — expected: flags, pagination, text/JSON output, stdout/stderr discipline, exit code 3, and no-mutation assertions pass.
- `rtk go run -buildvcs=false ./cmd/roundfix release plan --help` — expected: concise help truthfully documents the reset-plan contract and incompatible flags.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 4; Core Feature 15; User Story 6; Non-Goals / Out of
  Scope.
- `_techspec.md` → Implementation Design: Interfaces, Data Models, and API
  Contracts; Testing Approach; Build Order 7.
- ADR-0048 → read-only Release Plan and confirmation boundary.
- ADR-0065 → complete reset inventory and separate destructive approval.
