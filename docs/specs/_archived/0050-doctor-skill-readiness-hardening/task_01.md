---
task: task_01
spec: 0050-doctor-skill-readiness-hardening
status: completed
type: chore
complexity: low
---

# Task 01: Add the authorized Unicode collation dependency

## Overview

Add the one module required for the shared Unicode path-collation contract.
This is an isolated protected-tooling slice: only the expressly authorized
module files and this Task file may change.

## Requirements

1. MUST add `golang.org/x/text` at version `v0.40.0` through Go tooling.
2. MUST change only `go.mod`, `go.sum`, and this Task file.
3. MUST preserve every existing direct and indirect requirement unless Go
   tooling changes its classification as part of adding the authorized module.
4. MUST leave the module graph verifiable without adding source code in this
   Task.

## Subtasks

- [x] Capture the pre-existing changed-file set.
- [x] Add the authorized module through the Go toolchain.
- [x] Inspect the resulting manifest and checksum changes.
- [x] Prove the exact selected version and module integrity.
- [x] Run the protected-file postflight.

## Acceptance Criteria

- [x] `golang.org/x/text` resolves exactly to `v0.40.0`.
- [x] `go.mod` and `go.sum` are internally consistent and pass module
      verification.
- [x] No path other than `go.mod`, `go.sum`, and this Task file changes in this
      Task.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/go.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- interface: `go.mod`
- interface: `go.sum`

## Verification

- `rtk go list -m -f '{{.Path}} {{.Version}}' golang.org/x/text` — expected:
  exactly `golang.org/x/text v0.40.0`.
- `rtk go mod verify` — expected: every downloaded module is verified.

## References

- `_prd.md` → Project Constraints; Goals; Core Feature 2; Decisions.
- `_techspec.md` → Project Constraints; Integration Points; Build Order 1;
  Decisions.

## Result

Added `golang.org/x/text` at `v0.40.0` through `rtk go get`. Go's module
selection also raised the existing indirect `golang.org/x/sync` requirement
from `v0.20.0` to `v0.22.0`, the minimum selected by
`golang.org/x/text@v0.40.0`. No source code changed.

Verification:

- `rtk go list -m -f '{{.Path}} {{.Version}}' golang.org/x/text`: passed
  with exactly `golang.org/x/text v0.40.0`.
- `rtk go mod verify`: passed with `all modules verified`.
- `rtk go mod download -json golang.org/x/text@v0.40.0`: passed and reported
  the matching module and `go.mod` checksums plus the upstream
  `refs/tags/v0.40.0` origin.
- Protected-file postflight: passed. The pre-existing changed-file set was
  empty; the resulting unstaged set contains only `go.mod`, `go.sum`, and this
  Task file, with no staged or untracked paths.

Acceptance evidence:

- Exact version: the focused module query returned
  `golang.org/x/text v0.40.0`.
- Module consistency and integrity: `rtk go mod verify` exited successfully
  and verified every downloaded module.
- Bounded changes: `rtk git -c core.fsmonitor=false diff --name-only` listed
  only the three authorized paths; the cached-diff and untracked-file
  inspections were empty.

Follow-ups: none.
