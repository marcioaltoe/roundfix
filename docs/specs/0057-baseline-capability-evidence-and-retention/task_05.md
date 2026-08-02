---
task: task_05
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: medium
---

# Task 05: Satisfy a portable Verification role from the repository

## Overview

A portable Verification role that a repository already satisfies with its own
command keeps re-appearing as an unresolved workspace divergence, so a
maintainer sees a gate they have already selected reported as missing. The
projection already distinguishes a repository command from a profile
expectation; this Task lets that mapping satisfy the role.

## Requirements

1. MUST let a portable Verification role be satisfied by mapping it to a
   declared repository command, so the divergence disappears once mapped.
2. MUST require positive evidence: a role reports satisfied only when its
   declared repository command is present, and an unmapped role keeps the
   divergence it produces today.
3. MUST record which command satisfied the role, so the mapping is auditable
   rather than implicit.
4. MUST NOT execute the mapped command; presence is established from
   declaration, not by running it.
5. MUST leave every unmapped role's divergence, message, and blocking status
   unchanged.

## Subtasks

- [ ] Accept a mapping from a portable role to a declared repository command.
- [ ] Satisfy the role only when the declared command is present.
- [ ] Record the satisfying command on the result.
- [ ] Confirm unmapped roles are untouched.

## Acceptance Criteria

- [ ] A portable role mapped to a present declared repository command reports
      satisfied and produces no divergence.
- [ ] The same role with no mapping still produces its current divergence.
- [ ] A mapping naming a command that is not declared does not satisfy the
      role, and reports why.
- [ ] The satisfying command is recorded on the result.
- [ ] No command is executed to establish satisfaction.
- [ ] The characterization corpus shows no change for repositories with no
      mapping.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestPortableVerificationRoleMapping -count=1`
  — expected: exit 0; mapped-and-present satisfies, unmapped and
  mapped-but-absent do not.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0; unmapped repositories are unchanged.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 9 and 12; Success Metrics (a mapped repository gate
  satisfies its role).
- `_techspec.md` → Coverage Map; Build Order 9.
