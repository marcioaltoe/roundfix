---
task: task_01
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
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

- [x] Capture the current `go mod tidy -diff` delta as the red signal.
- [x] Run `go mod tidy` with the repository's selected Go toolchain.
- [x] Confirm the resolved module versions did not change.
- [x] Prove the protected changed-file allowlist and module integrity.

## Acceptance Criteria

- [x] `golang.org/x/text v0.40.0` is a direct requirement selected by Go
      tooling.
- [x] Stale sums are removed and required sums are present without a version
      upgrade or downgrade.
- [x] `go mod tidy -diff` exits successfully with no output.
- [x] `go mod verify` succeeds.
- [x] Newly changed paths are limited to the two authorized module files and
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

## Result

Go 1.26.5 produced the authorized module metadata. The initial
`rtk go mod tidy -diff` exited `1` and showed the required red delta:
`golang.org/x/text v0.40.0` moved from indirect to direct, obsolete sums for
`golang.org/x/mod v0.33.0`, `golang.org/x/sync v0.20.0`, and
`golang.org/x/tools v0.42.0` were removed, and the selected
`golang.org/x/mod v0.37.0` and `golang.org/x/tools v0.47.0` module-file sums
were added. `rtk go mod tidy` then exited `0`.

Acceptance evidence:

- `rtk go mod edit -json` reports `golang.org/x/text v0.40.0` without
  `Indirect`, proving the direct requirement.
- Pre- and post-tidy `rtk go list -m all` output selected the same module
  versions; the `go.mod` diff changes only the direct classification.
- `rtk go mod tidy -diff` exits `0` with no output after the metadata update.
- `rtk go mod verify` exits `0` with `all modules verified`.
- The declared changed-path filter exits `0`, and
  `rtk git -c core.fsmonitor=false status --short` lists only `go.mod`,
  `go.sum`, and this Task file.

Follow-ups: none.
