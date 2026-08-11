---
task: task_05
spec: 0003-dogfood-polish
status: completed
type: backend
complexity: low
---

# Task 05: Stop by spec target

## Overview

`roundfix stop` accepts a run id, `--pr`, or head-repo/branch — but a stuck
spec Run must be hunted by run id. Add `--spec <slug>`: resolve the Active Run
for the current repository's spec work target and stop it, mirroring the
`--pr` selector's shape. Verifiable through buffer-captured stop tests over a
real temp store.

## Requirements

1. MUST add `--spec <slug>` to the stop command: resolves the Active Run for
   the spec target `("spec", "<git_root>#<slug>")` where the git root comes
   from the current working directory's repository, then stops it exactly as
   stop-by-run-id does.
2. MUST keep all stop selectors mutually exclusive with the existing error
   shape when combined or absent.
3. MUST report "no Active Run" for the spec target with an actionable message
   naming the target, mirroring the `--pr` not-found shape.
4. MUST update the stop usage/help text truthfully.

## Subtasks

- [x] `--spec` flag and target resolution from the current repository
- [x] Mutual-exclusion and not-found error shapes
- [x] Help text
- [x] Buffer-captured tests over a temp store

## Acceptance Criteria

- [x] Stopping an Active implement Run by `--spec` releases its lock and
      settles `Stopped`, asserted via the store.
- [x] `--spec` + `--pr` (and `--spec` + run id) fail with the existing
      mutual-exclusion message; unknown slug fails with the named target.
- [x] Help lists the new selector; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix stop --help` — expected: `--spec` listed, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 5. `_techspec.md` → API Contracts,
Build Order 5. Dogfood finding 13. ADR-0016 (spec target key shape).

## Result

- `roundfix stop --spec <slug>` now resolves the Active Run for the current
  repository's Spec target via the existing ADR-0016 key shape:
  `("spec", "<git_root>#<slug>")`. The CLI uses `roundconfig.Load`'s resolved
  `GitRoot` and `store.ActiveSpecRun`, then completes the same run id with
  `StateStopped` through the existing stop path.
- `TestRunStopBySpecStopsMatchingActiveRun` creates an Active implement Run in
  a temp Run Database, runs `stop --spec` from a repository subdirectory,
  asserts the run is `Stopped`, and proves the target lock is released by
  creating another implement Run for the same Spec target.
- `TestRunStopSpecSelectorRejectsOtherSelectors` covers `--spec` with `--pr`
  and `--spec` with a positional run id; both fail with the existing
  `cannot be combined` selector error shape. `TestRunStopBySpecReportsMissingActiveRun`
  asserts the unknown slug error names the repository and Spec target.
- `TestRunStopHelpListsSpecSelector` and
  `rtk go run ./cmd/roundfix stop --help` both show
  `roundfix stop --spec <slug>` and the `--spec` option.
- Red signal before implementation: the focused stop tests failed because
  `--spec` was an unknown flag and stop help omitted the selector.
- Verification passed: `rtk go test ./internal/cli/` reported
  `Go test: 150 passed in 1 packages`; `rtk go run ./cmd/roundfix stop --help`
  exited 0 and listed `--spec`; `rtk go test ./...` reported
  `Go test: 466 passed in 16 packages`.
