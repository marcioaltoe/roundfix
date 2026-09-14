---
task: task_02
spec: 0135-a-rename-the-gate-sees-from-both-sides
status: completed
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

## Result

Implemented:

- The composed Spec authorization resolver now carries a typed reason out of
  location resolution. A failed delivery revision reports the existing
  `unavailable_revision` code with field `revision` and the normalized revision
  value; other location failures retain the existing unreadable-record default.
- Added a real-Git regression through `ReadSpecAuthorization` that asserts the
  revision diagnostic, the unresolved outcome, and refusal of `implement`.
- The direct `ReadAuthorization` implementation under `internal/authorization`
  remains unchanged.

Focused-check evidence:

- Before the production edit,
  `rtk proxy go test ./internal/spec -run 'NamesUnavailableRevision'` exited 1:
  the new regression observed `unreadable_record`, field `spec_root`, and the
  Spec Root value for `refs/heads/unavailable`.
- After the production edit,
  `rtk go test ./internal/spec -run 'NamesUnavailableRevision'` exited 0 with
  one passing test.
- `rtk go test ./internal/spec -run 'Test(SpecAuthorizationNamesUnavailableRevision|OperationAuthorityReportsUnresolvableSpecRoot|OperationAuthorityResolvesExternalSpecRoot)$'`
  exited 0 with three passing tests, covering the changed branch and both
  preserved branches together.
- `rtk go test ./internal/spec` exited 0 with 434 passing tests.
- `rtk gofmt -d internal/spec/authorization.go internal/spec/authorization_test.go`
  exited 0 with no output.
- `rtk make verify-incremental` exited 2 at `fmt-check` before tests because
  unchanged `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. Separate
  `rtk git diff --exit-code` checks for both paths exited 0, so this Task did
  not alter them.

Acceptance-criterion evidence:

1. `TestSpecAuthorizationNamesUnavailableRevision` changed from red to green
   and requires `unavailable_revision`, field `revision`, and value
   `refs/heads/unavailable`.
2. `TestOperationAuthorityReportsUnresolvableSpecRoot` passes in the combined
   focused run and still requires an unresolved `unreadable_record` on field
   `spec_root`.
3. The new revision regression requires `AuthorizationUnresolved` and rejects
   `Permits(implement)`; the preserved Spec Root regression also requires
   `AuthorizationUnresolved`.

The Daemon-owned commands under `## Verification` were not run.
