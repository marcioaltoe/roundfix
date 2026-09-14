---
task: task_03
spec: 0136-a-rename-the-committer-can-stage
status: completed
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

## Result

Implementation:

- The record-location resolver now returns the Spec-relative authorization path
  it derived before a later project-repository or revision failure.
- The unresolved-result builder uses that derived path and otherwise composes
  `<configured Spec Root>/<slug>/_authorization.md`. It leaves the existing
  unresolved outcome and reason construction unchanged.
- `TestUnresolvedRecordNamesAPathUnderTheSpecRoot` covers both branches against
  real temporary Git repositories.

Focused checks:

- `rtk go test ./internal/spec -count=1 -run UnresolvedRecordNames` reproduced
  the defect before the production edit: both cases reported the bare
  `asking-spec/_authorization.md` fragment.
- `rtk /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test ./internal/spec -count=1 -run UnresolvedRecordNames`
  passed after the edit.
- `rtk /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test ./internal/spec -count=1 -run 'Test(UnresolvedRecordNamesAPathUnderTheSpecRoot|SpecAuthorizationNamesUnavailableRevision|OperationAuthorityReportsUnresolvableSpecRoot)$'`
  passed, covering the reported paths, both reason contracts, and both
  unresolved outcomes together.
- `rtk /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go test ./internal/spec -count=1`
  passed.
- `rtk make verify-incremental GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`
  exited non-zero: one `internal/agent` `acpx` probe was killed and four
  `internal/daemon` Task-cycle tests timed out during the parallel run. The
  failed Agent test passed alone, and the four Daemon tests passed together at
  their package boundary; the full incremental gate was not rerun or treated as
  passing.
- Go 1.26.7 `gofmt -l` on both changed Go files and `rtk git diff --check`
  produced no output.

Acceptance evidence:

- The new resolver test observes `docs/specs/asking-spec/_authorization.md`
  after a derived-path revision failure and the path composed beneath a
  non-Git configured Spec Root when derivation fails first.
- `TestSpecAuthorizationNamesUnavailableRevision` retains reason code
  `unavailable_revision`, field `revision`, and the requested revision value;
  `TestOperationAuthorityReportsUnresolvableSpecRoot` retains reason code
  `unreadable_record` and field `spec_root`.
- Both new cases assert `AuthorizationUnresolved`; the unavailable-revision
  regression also continues to assert that `implement` is not permitted.

The Task's declared `## Verification` commands were not run; the Daemon owns
them.

## Carry-forward provenance

- Source Run: `run_20260914T134813Z_e15c7ee834fbe8bf`
- Source commit: `54ab556d85692f0dfd2c5476892dfbe75bbec320`
