---
task: task_03
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
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

## Result

Implementation evidence:

- The operation resolver now accepts the resolved Spec Root separately from
  the project repository. It resolves the Spec repository root, its committed
  revision, and the record path relative to that repository, while retaining
  the project repository and delivery target as distinct location fields.
- Implement passes its resolved checkout Spec Root to the resolver. Settle
  passes the Spec Root for the selected recovery surface, so both public paths
  read the same external record their Task Graph came from.
- The external-root CLI fixtures now remove the project-repository copy of the
  record. Their authority comes only from the committed external Spec
  repository.
- The Task 01 characterization moved only its operation-resolution row from an
  unreadable default-root record to the committed record under the external
  Spec Root. Its Task 04 citation-refusal row remains unchanged.

Focused-check evidence:

- The first focused regression run was blocked by the sandboxed Go build
  cache. Its unchanged authorized retry reached the intended pre-change red
  signal: all three new operation-authority cases failed to compile because
  `ReadSpecAuthorization` had no resolved Spec Root input.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260911T103803Z_ebaa4c2da86a935e.task_03/.gocache
  /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
  test ./internal/spec ./internal/speccheck -count=1 -run
  'Test(OperationAuthority|ExternalSpecRootResolutionCharacterization)'`
  passed four tests. The external record uses a non-default `specs/<slug>`
  path, reads the external repository's committed HEAD despite a dirty
  worktree record, and needs no project copy. The default case preserves the
  existing `docs/specs/<slug>/_authorization.md` path and project revision;
  the non-Git root case returns `unresolved` with field `spec_root`.
- `rtk go test ./internal/cli -count=1 -run
  'TestRun(Implement|Settle)UsesConfiguredExternalSpecRoot|TestRunSettleNoCommitPrintsNoCommitPathsOrSharedWarning'`
  passed three tests after the project-copy fixture was removed, exercising
  Implement and Settle through their public command runners.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260911T103803Z_ebaa4c2da86a935e.task_03/.gocache
  /Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
  test ./internal/daemon -count=1 -run
  'Test(RunDispositionCharacterizationStoppedRunLeavesTasksPending|TaskCycleExecutesAgentVerifySettleCommitContract)$'`
  passed two default-root Daemon callers.
- `rtk make fmt-check build
  GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
  GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`
  exited 0.
  The same formatter listed none of this Task's changed Go files, and
  `rtk git diff --check` exited 0.
- The required `rtk make verify-incremental` first stopped at Go 1.27.1
  formatter differences in two unchanged CLI test files. The unchanged retry
  with the repository-compatible Go 1.26.7 toolchain advanced into the test
  stage, where sandbox policy blocked an existing `api.github.com` request;
  automatic approval review rejected that network egress. This command has no
  terminal verdict from this turn.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.

Acceptance-criterion evidence:

- External Spec Root authority: `TestOperationAuthorityResolvesExternalSpecRoot`
  and the focused public Implement and Settle cases passed against a separate
  committed Git repository.
- Default-root compatibility: `TestOperationAuthorityDefaultRootUnchanged`
  asserted the identical canonical path and exact delivery revision, and the
  focused Daemon callers passed.
- No duplicate record: the external resolver test asserts the project path is
  absent, while the public command fixtures delete their former project copy
  before dispatch.
- Unresolvable root: `TestOperationAuthorityReportsUnresolvableSpecRoot`
  distinguishes `unresolved` from both `refused` and `granted` and records the
  `unreadable_record` reason at `spec_root`.
