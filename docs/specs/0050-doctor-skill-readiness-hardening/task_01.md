---
task: task_01
spec: 0050-doctor-skill-readiness-hardening
status: pending
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

- [ ] Capture the pre-existing changed-file set.
- [ ] Add the authorized module through the Go toolchain.
- [ ] Inspect the resulting manifest and checksum changes.
- [ ] Prove the exact selected version and module integrity.
- [ ] Run the protected-file postflight.

## Acceptance Criteria

- [ ] `golang.org/x/text` resolves exactly to `v0.40.0`.
- [ ] `go.mod` and `go.sum` are internally consistent and pass module
      verification.
- [ ] No path other than `go.mod`, `go.sum`, and this Task file changes in this
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

