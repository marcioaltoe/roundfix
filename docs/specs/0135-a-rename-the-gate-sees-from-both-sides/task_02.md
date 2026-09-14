---
task: task_02
spec: 0135-a-rename-the-gate-sees-from-both-sides
status: pending
type: backend
complexity: low
---

# Task 02: Report an unresolvable revision as an unavailable revision

## Overview

Resolving a record's location resolves a revision before the record is read,
and a failure there returns early reporting an unreadable record against the
Spec Root field. An unknown delivery target is therefore described as a missing
Spec Root. The distinct reason for this condition already exists and is already
returned by the direct reader; this slice uses it on the branch that resolves
the location.

## Requirements

1. MUST report an unresolvable revision with the existing unavailable-revision
   reason, naming the revision that could not be resolved rather than the Spec
   Root.
2. MUST keep a genuinely unresolvable Spec Root reporting its own existing
   reason and field.
3. MUST keep the outcome unresolved and fail-closed in both cases: this is a
   diagnostic contract, so nothing becomes permitted.
4. MUST NOT change the direct reader, which already returns the
   unavailable-revision reason for a bad revision.

## Subtasks

- [ ] Classify the failed revision resolution with the existing reason.
- [ ] Prove the Spec Root branch keeps its own reason.
- [ ] Prove both outcomes stay unresolved.

## Acceptance Criteria

- [ ] Resolving a record against an unknown delivery revision reports the
      unavailable-revision reason and names that revision; today it reports an
      unreadable record against the Spec Root.
- [ ] An unresolvable Spec Root still reports its existing reason and field.
- [ ] Both cases remain unresolved, so neither permits a grant.

## Context

- interface: `internal/spec/authorization.go`

## Verification

- `grep -q 'func TestSpecAuthorizationNamesUnavailableRevision' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestSpecAuthorizationNamesUnavailableRevision$'` — an unknown delivery revision is named as such; this fails today.
- `grep -q 'func TestSpecAuthorizationNamesUnavailableRevision' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestOperationAuthorityReportsUnresolvableSpecRoot$'` — the Spec Root branch keeps its own reason.
- `grep -q 'func TestSpecAuthorizationNamesUnavailableRevision' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestOperationAuthorityResolvesExternalSpecRoot$'` — external Spec Root resolution is untouched.

## References

- `_prd.md` → Goal 2; Core Features 3-4; Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: An unresolvable revision says so;
  Testing Approach observation 4; Build Order 2.
