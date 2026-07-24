---
task: task_10
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Add reset-plan domain values and read-only inventory interfaces.
- [x] Implement local/remote tag and paginated GitHub Release adapters.
- [x] Add deterministic framing and digest calculation.
- [x] Add CLI flag validation, text rendering, JSON rendering, and exit mapping.
- [x] Add incomplete-inventory, dirty-target, pagination, and digest-sensitivity
      tests.
- [x] Add mutation-spy tests proving the command remains read-only.

## Acceptance Criteria

- [x] A complete fixture inventory produces deterministic text and JSON plans
      with the same digest and exit code 3.
- [x] Every stable local tag, remote tag, and paginated GitHub Release appears
      exactly once with its immutable identity.
- [x] A changed target revision, tag, or Release changes the plan digest.
- [x] Dirty state, conflicting flags, malformed target versions, and incomplete
      inventory fail with the documented exit category and useful diagnostic.
- [x] Captured provider calls contain only read operations; mutation spies are
      never invoked.
- [x] The existing non-reset Release Plan behavior remains unchanged.

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

## Result

Implemented a read-only reset mode through the existing Release Plan Command.
`roundfix release plan --reset-to <stable-version>` now resolves a clean
committed `HEAD`, inventories every stable local and remote tag plus every page
of GitHub Releases, sorts the complete inventory, and returns an
`approval_required` plan with a digest and exit code 3. Text and JSON output
include each tag's bound identity and target commit and each GitHub Release's
database ID, node ID, tag, and target commit when the tag inventory resolves it.

Acceptance evidence:

- `TestReleasePlanResetTextAndJSONInventoryMatchThroughRunBoundary` invokes the
  public `RunContext(... ) int` boundary twice. The fixture's two local tags,
  three remote tags, and three GitHub Releases across two pages appear once in
  deterministic order; text and JSON contain the same digest and both exit 3.
- `TestBuildResetPlanDigestChangesWithEveryBoundInput` changes the target
  version, target revision, target commit, tag target, and GitHub Release
  identity independently; every change produces a different digest.
- Domain and CLI negative tests reject malformed reset versions, every
  conflicting range/classification flag, dirty state, unavailable remote tags,
  unavailable GitHub results, null pagination pages, and duplicate immutable
  identities without emitting a partial stdout plan.
- The reset inventory interface exposes only `Tags` and `Releases`. Captured
  GitHub calls are exactly paginated `GET` requests with `--paginate --slurp`;
  Git and GitHub mutation spies recorded zero calls.
- The complete existing `internal/releaseplan` and `internal/cli` suites pass,
  preserving non-reset range planning, classification, rendering, and exit
  behavior.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task10-go-cache go test ./internal/releaseplan -run 'Reset|Digest|Inventory'` — PASS.
- `rtk env GOCACHE=/private/tmp/roundfix-task10-go-cache go test ./internal/cli -run 'ReleasePlan.*Reset'` — PASS.
- `rtk env GOCACHE=/private/tmp/roundfix-task10-go-cache go run -buildvcs=false ./cmd/roundfix release plan --help` — PASS; help documents reset inventory, incompatible flags, exit 3, and the separate deletion authority.
- `rtk env GOCACHE=/private/tmp/roundfix-task10-go-cache make verify` — PASS.
- `rtk git diff --check` — PASS.

Follow-ups already owned by the Task Graph:

- Task 11 changes the Release Plan schema identity with the distribution-wide
  `0.0.1` version reset.
- Task 13 updates the durable release runbook and shipped Roundfix skill after
  the dependent distribution and QA tasks complete.
