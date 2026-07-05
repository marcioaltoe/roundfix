---
task: task_05
spec: 0003-dogfood-polish
status: pending
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

- [ ] `--spec` flag and target resolution from the current repository
- [ ] Mutual-exclusion and not-found error shapes
- [ ] Help text
- [ ] Buffer-captured tests over a temp store

## Acceptance Criteria

- [ ] Stopping an Active implement Run by `--spec` releases its lock and
      settles `Stopped`, asserted via the store.
- [ ] `--spec` + `--pr` (and `--spec` + run id) fail with the existing
      mutual-exclusion message; unknown slug fails with the named target.
- [ ] Help lists the new selector; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix stop --help` — expected: `--spec` listed, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 5. `_techspec.md` → API Contracts,
Build Order 5. Dogfood finding 13. ADR-0016 (spec target key shape).
