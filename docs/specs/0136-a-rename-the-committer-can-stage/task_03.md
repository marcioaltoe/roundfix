---
task: task_03
spec: 0136-a-rename-the-committer-can-stage
status: pending
type: backend
complexity: low
---

# Task 03: Carry a resolvable record path into the unresolved result

## Overview

When a record's location cannot be resolved, the unresolved result labels the
record with the Spec-relative fragment rather than a path under the configured
Spec Root, and the refusal message prints that label. A reader following it
looks at a path that does not exist.

## Requirements

1. MUST report a record path that resolves under the configured Spec Root when
   the location could not be resolved.
2. MUST use the path the resolver derived when one was derived before the
   failure.
3. MUST keep every existing reason code and field, including the
   unavailable-revision reason and the unresolvable-Spec-Root reason.
4. MUST keep the outcome unresolved, so nothing becomes permitted.

## Subtasks

- [ ] Carry the derived path into the unresolved result.
- [ ] Compose the path from the Spec Root when none was derived.
- [ ] Prove the reason codes and fields are unchanged.

## Acceptance Criteria

- [ ] An unresolved record reports a path under the configured Spec Root, not a
      bare slug-relative fragment.
- [ ] The unavailable-revision reason and the unresolvable-Spec-Root reason keep
      their codes and fields.
- [ ] Both cases remain unresolved.

## Context

- interface: `internal/spec/authorization.go`

## Verification

- `grep -q 'func TestUnresolvedRecordNamesAPathUnderTheSpecRoot' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestUnresolvedRecordNamesAPathUnderTheSpecRoot$'` — the reported path resolves under the Spec Root; this fails today.
- `grep -q 'func TestUnresolvedRecordNamesAPathUnderTheSpecRoot' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestSpecAuthorizationNamesUnavailableRevision$'` — the unavailable-revision reason is unchanged.
- `grep -q 'func TestUnresolvedRecordNamesAPathUnderTheSpecRoot' internal/spec/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestOperationAuthorityReportsUnresolvableSpecRoot$'` — the Spec Root reason is unchanged.

## References

- `_prd.md` → Goal 3; Core Feature 5; Declared intentional breaks 3.
- `_techspec.md` → Implementation Design: Name the record a reader can open;
  Testing Approach observation 6; Build Order 3.
