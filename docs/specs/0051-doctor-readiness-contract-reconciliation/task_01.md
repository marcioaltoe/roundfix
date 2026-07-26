---
task: task_01
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
type: chore
complexity: low
---

# Task 01: Tidy authorized Go module metadata

## Overview

Restore the module graph to the exact state produced by Go tooling after the
approved text-collation dependency was introduced. This slice changes no
runtime behavior and is complete when the module files are tidy, verified, and
the protected mutation stays inside its exact allowlist.

## Requirements

1. MUST run Go tooling to produce the module metadata; hand editing is
   forbidden.
2. MUST keep every selected dependency version unchanged while classifying
   `golang.org/x/text` as a direct requirement and removing stale checksums.
3. MUST limit all Task mutations to `go.mod`, `go.sum`, and
   `docs/specs/0051-doctor-readiness-contract-reconciliation/task_01.md`.
4. MUST NOT change a tool configuration, ignore file, script, plugin
   declaration, or any other repository path.

## Subtasks

- [ ] Capture the current `go mod tidy -diff` delta as the red signal.
- [ ] Run `go mod tidy` with the repository's selected Go toolchain.
- [ ] Confirm the resolved module versions did not change.
- [ ] Prove the protected changed-file allowlist and module integrity.

## Acceptance Criteria

- [ ] `golang.org/x/text v0.40.0` is a direct requirement selected by Go
      tooling.
- [ ] Stale sums are removed and required sums are present without a version
      upgrade or downgrade.
- [ ] `go mod tidy -diff` exits successfully with no output.
- [ ] `go mod verify` succeeds.
- [ ] Newly changed paths are limited to the two authorized module files and
      this Task file.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/go.md`
- interface: `go.mod`
- interface: `go.sum`

## Verification

- `rtk go mod tidy -diff` — expected: exit `0` with no diff.
- `rtk go mod verify` — expected: every downloaded module is verified.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != "go.mod" && path != "go.sum" && path != "docs/specs/0051-doctor-readiness-contract-reconciliation/task_01.md") {print; bad=1}} END {exit bad}'` — expected: no out-of-allowlist path.

## References

- `_prd.md` → Core Features 5; User Story 5; Success Metrics.
- `_techspec.md` → Module and naming hygiene; Build Order 1.
