---
task: task_03
spec: 0002-acpx-migration
status: completed
type: backend
complexity: low
---

# Task 03: Pin the acpx version in Probe and Preflight Validation

## Overview

Make the acpx dependency visible before any Run exists: the runner's `Probe` verifies the acpx binary is on PATH and its version equals the pinned constant, and every mismatch or absence becomes one actionable Preflight Validation message naming the exact install command. Verifiable alone through the helper-process rig and CLI preflight tests.

## Requirements

1. MUST define the pinned acpx version as a code constant in the agent package (the version current when implementation starts; upgrades are deliberate commits, never configuration).
2. MUST implement `Probe` for the acpx runner: binary present on PATH and `acpx --version` output equal to the pin; the probe never creates sessions or spawns adapters.
3. MUST surface probe failures through the existing probe call sites as Preflight Validation failures (exit 2, no Run created) with one actionable message per case: missing binary → `npm install -g acpx@<pin>`; version mismatch → the same command phrased as an upgrade/downgrade naming both versions found and required.
4. MUST keep the command override escape hatch probing through acpx as well (the override replaces the adapter, not the client).
5. MUST NOT change probe behavior for any code path still running the SDK runner in this task — the acpx probe ships dark until task_04 wires the runner in, so the existing suite passes unchanged.

## Subtasks

- [x] Pinned version constant
- [x] `Probe`: PATH check plus version equality via the fake-acpx rig
- [x] Actionable Preflight Validation messages for missing and mismatched acpx
- [x] Probe coverage for the command override case

## Acceptance Criteria

- [x] Probe tests cover: matching version (pass), missing binary (message carries the install command), mismatched version (message names found and required versions plus the command).
- [x] The messages follow the repo's one-actionable-message preflight convention and write nothing to the Run Database.
- [x] The full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/agent/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Features 1, 6. `_techspec.md` → acpx invocation mapping (Probe), API Contracts, Build Order 3, Decisions (pin is a code constant). ADR-0017.

## Result

Implemented the acpx probe slice behind the existing runner seam. The agent package now pins `acpx` at `0.12.0` (`PinnedACPXVersion`), the current npm package version resolved when this task started, and `ACPXRunner.Probe` checks only the acpx client binary with `acpx --version`. Missing or mismatched versions return one actionable install command: `npm install -g acpx@0.12.0`. Command overrides still probe the acpx client, not the override command.

Evidence:

- `TestACPXProbePassesWhenVersionMatchesPin`, `TestACPXProbeMissingBinaryNamesInstallCommand`, `TestACPXProbeMismatchedVersionNamesFoundRequiredAndInstallCommand`, and `TestACPXProbeCommandOverrideStillChecksACPXClient` passed under `rtk go test ./internal/agent/ -run TestACPXProbe`.
- `TestRunResolveACPXProbeFailureReportsActionablePreflight` passed under `rtk go test ./internal/cli/ -run TestRunResolveACPXProbeFailureReportsActionablePreflight`; it asserts exit 2, no stdout, one install command in stderr, the standard no-side-effects line, and no Run Database.
- `rtk go test ./internal/agent/ ./internal/cli/` passed: 198 tests across 2 packages.
- `rtk go test ./...` passed: 439 tests across 16 packages.
- `rtk make verify` passed: full tests, `roundfix skills check`, and build.
