---
task: task_03
spec: 0012-npm-distribution
status: completed
type: test
complexity: low
---

# Task 03: Upgrade-asset compatibility test pinning the asset-name scheme

## Overview

Pin the release-asset naming to the Upgrade Command's platform selection with a
Go test, so the workflow's asset names and `selectPlatformAsset` can never drift
apart. The test encodes the task_01 mapping's asset-name scheme and asserts the
existing selection resolves each platform's asset.

## Requirements

1. MUST add a Go test that, for every mapping row, constructs the release-asset
   name the workflow will produce and asserts `selectPlatformAsset` resolves it
   for that row's `GOOS`/`GOARCH`.
2. MUST assert checksum-style assets are ignored and that the wrong-platform
   asset is never selected.
3. MUST express the asset-name scheme as the single fixture the workflow will
   also use, so a scheme change breaks this test first.
4. MUST NOT change `selectPlatformAsset` behavior — this task only locks the
   contract.

## Subtasks

- [x] Test fixture of asset names per mapping row (Go tokens)
- [x] Assert `selectPlatformAsset` resolves each platform's asset
- [x] Assert checksum/other assets are skipped and no cross-platform match
- [x] Reference the fixture as the workflow's naming source of truth

## Acceptance Criteria

- [x] For each of the five targets, the constructed asset name resolves via `selectPlatformAsset` for the matching `GOOS`/`GOARCH`.
- [x] A checksum asset and a wrong-platform asset are not selected.
- [x] The asset-name scheme lives in one fixture the release workflow reuses.

## Verification

- `rtk go test ./internal/cli/ -run PlatformAsset` — expected: the compatibility test passes.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 3; Core Feature 5. `_techspec.md` → Asset-name
compatibility, Build Order 3. ADR-0031.
