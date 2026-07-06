---
task: task_04
spec: 0012-npm-distribution
status: pending
type: infra
complexity: high
---

# Task 04: Tag-triggered release workflow

## Overview

Add the GitHub Actions workflow that turns a `v*` tag into a published release:
prove the tag, launcher version, and embedded app version agree; run the full
local gate; cross-compile the `GOOS`/`GOARCH` matrix; publish the per-platform
packages then the launcher; and upload the binaries as GitHub Release assets the
Upgrade Command can resolve. This is the one-tag-push release path.

## Requirements

1. MUST trigger on `v*` tags and fail fast, before any build, when the tag, the
   launcher `package.json` version, and the embedded app version disagree.
2. MUST run the full local verification gate (`make verify`) before building.
3. MUST cross-compile the binary for all five targets via a `GOOS`/`GOARCH`
   matrix on one runner, staging each binary into its per-platform package.
4. MUST publish every `@roundfix/cli-*` package before the `roundfix` launcher,
   and MUST upload the same binaries as GitHub Release assets named per the
   task_03 scheme so `selectPlatformAsset` resolves them.

## Subtasks

- [ ] `v*` trigger and version-agreement guard job
- [ ] `make verify` gate job
- [ ] Matrix cross-compile staging binaries into per-platform packages
- [ ] Publish order (platform packages → launcher) and Release-asset upload
- [ ] Documented required secrets/tokens (npm publish, GitHub release)

## Acceptance Criteria

- [ ] A `v*` tag with agreeing versions publishes all five per-platform packages and the launcher and uploads matching Release assets in one run.
- [ ] A tag whose versions disagree fails the workflow before publishing anything.
- [ ] Release-asset names match the task_03 fixture scheme (verified by that test staying green).
- [ ] The workflow runs `make verify` and does not publish when it fails.

## Verification

- Workflow lint/parse (e.g. `actionlint .github/workflows/release.yml`) — expected: no errors.
- `rtk go test ./...` — expected: full suite passes (including the task_03 compatibility test).
- A dry-run or `workflow_dispatch` on a test tag (documented in the runbook) — expected: guard, verify, and build stages succeed to the publish boundary.

## References

`_prd.md` → User Stories 4-5; Core Feature 3. `_techspec.md` → Release workflow,
Version agreement, Build Order 5. ADR-0031. Depends on the widened skill bundle
(task_06) so the released binary carries the owned set and `make verify` covers
its sync/check.
