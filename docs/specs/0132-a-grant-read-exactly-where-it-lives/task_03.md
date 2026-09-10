---
task: task_03
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: backend
complexity: high
---

# Task 03: Derive the record path from the resolved Spec Root

## Overview

The operation resolver assembles the consuming record's path from the constant
`docs/specs`, so a configured Spec Root elsewhere is read in the wrong place and
Implement and Settle refuse valid work unless the record is duplicated into the
code repository. This slice takes the resolved root as an input. It is
verifiable alone: an external root resolves, and the default root is unchanged.

## Requirements

1. MUST derive the consuming record's path and revision from the resolved Spec
   Root rather than from a constant, at every reader that needs them.
2. MUST resolve a Spec Root that lies outside the code repository, so Implement
   dispatch and Settle read the record where the repository says Specs live.
3. MUST keep the default root's behavior identical, including the path it reads
   and the revision it resolves.
4. MUST NOT require a duplicate record inside the code repository for an
   external root to work.
5. MUST report an unresolvable root as unresolved rather than as a refusal or a
   grant, so a misconfiguration is distinguishable from a missing authority.
6. MUST update only the characterization rows this Task intentionally moves.

## Subtasks

- [ ] Thread the resolved Spec Root into the operation resolver.
- [ ] Resolve path and revision from that root.
- [ ] Prove the default root is unchanged.
- [ ] Prove an unresolvable root reports unresolved.

## Acceptance Criteria

- [ ] A Spec Root outside the code repository resolves its consuming record, and
      operation authority answers from it; it refused before this Task.
- [ ] The default root resolves the identical path and revision it resolves
      today.
- [ ] No duplicate record inside the code repository is needed for the external
      case.
- [ ] An unresolvable root reports unresolved, distinguishable from both refused
      and granted.

## Context

- interface: `internal/spec/authorization.go`
- interface: `internal/cli/implement.go`

## Verification

- `grep -q 'func TestOperationAuthorityResolvesExternalSpecRoot' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestOperationAuthorityResolvesExternalSpecRoot$'` — an external Spec Root resolves its record; this fails today.
- `grep -q 'func TestOperationAuthorityResolvesExternalSpecRoot' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestOperationAuthorityDefaultRootUnchanged$'` — the default root's path and revision did not move.
- `grep -q 'func TestOperationAuthorityResolvesExternalSpecRoot' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec ./internal/cli -run 'Authorization|Authority'` — the public dispatch and settle paths agree with the resolver.

## References

- `_prd.md` → Goals 2; User Stories 2; Core Features 2.
- `_techspec.md` → Implementation Design: Paths derived from the resolved root; Build Order 3.
