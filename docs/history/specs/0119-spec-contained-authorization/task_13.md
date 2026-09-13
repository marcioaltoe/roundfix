---
task: task_13
spec: 0119-spec-contained-authorization
status: completed
type: backend
complexity: high
---

# Task 13: Ask the grant which operations it permits

## Overview

The typed `operations` list is parsed and never consulted: `Permits` has no
production caller, so a record authorizing nothing reaches Implement dispatch
exactly like the through-merge record. Core Feature 5 promises those approvals
stay distinguishable, and an unread permission is not a distinction. This slice
makes each action boundary ask the ancestor grant for the operation it needs
before performing it.

## Requirements

1. MUST resolve the required operation from the ancestor grant at each boundary
   the TechSpec's operation-authority table names, and MUST refuse before the
   action rather than reporting beside a completed one.
2. MUST refuse strict authoring validation for a Spec that carries an
   executable Task Graph while its consuming grant does not permit `implement`,
   so a grant authorizing nothing cannot look dispatchable.
3. MUST refuse Implement dispatch before the first Agent Session opens when
   `implement` is absent, leaving no Agent work, shell execution, Task
   Worktree, commit or push behind.
4. MUST refuse the Daemon's Task commit when `commit` is absent, and the Run's
   push when `push` is absent, as independent decisions rather than one
   combined check.
5. MUST apply the same contract through the public Settle command.
6. MUST treat an absent operations list as permitting nothing, so every
   boundary refuses on it, and MUST name both the operation required and the
   record read in each refusal.
7. MUST NOT refuse an action a record permits: the through-merge record used by
   this Spec must continue to authorize implementation, commit and push exactly
   as it does today, proven against the current record rather than a fixture
   alone.
8. MUST NOT invent a CLI refusal for `pull_request`, `merge` or `release`; those
   have no command boundary in this Spec and remain recorded policy enforced by
   the changed-path audit.

## Subtasks

- [ ] Resolve the ancestor grant's operations at each named boundary.
- [ ] Refuse strict authoring validation without `implement`.
- [ ] Refuse dispatch, Task commit and push on their own operations.
- [ ] Apply the same contract through Settle.
- [ ] Prove the current record still authorizes its permitted actions.

## Acceptance Criteria

- [ ] A consuming grant without `operations` fails strict authoring validation
      for a Spec with an executable Task Graph; the same Spec passes with
      `implement` present.
- [ ] Implement dispatch refuses before any Agent Session opens when `implement`
      is absent, and the refusal names the operation and the record; no Run,
      Agent, Task Worktree, commit or push side effect exists afterwards.
- [ ] A grant permitting `implement` but not `commit` settles its Task
      unresolved instead of writing a commit; a grant permitting `commit` but
      not `push` commits and does not push.
- [ ] Settle refuses on a missing `commit` operation through the public command.
- [ ] This Spec's current record still authorizes implementation, commit and
      push, proven against the record in the repository.
- [ ] No CLI path refuses `pull_request`, `merge` or `release`, and a search
      finds no such refusal.

## Context

- interface: `internal/authorization/authorization.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/settle.go`

## Verification

- `grep -rn 'Permits(' internal/ --exclude='*_test.go' | grep -qv 'func (resolution AuthorizationResolution) Permits'` — a production caller exists; today only the definition does.
- `grep -q 'func TestStrictCheckRefusesMissingImplementAuthority' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestStrictCheckRefusesMissingImplementAuthority$'` — strict validation refuses a Task Graph whose grant permits no implementation.
- `grep -q 'func TestDispatchRefusesMissingImplementAuthority' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestDispatchRefusesMissingImplementAuthority$'` — dispatch refuses before any Agent Session, with no side effect left behind.
- `grep -q 'func TestCommitAndPushAuthorityAreSeparate' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestCommitAndPushAuthorityAreSeparate$'` — commit and push refuse on their own operations.
- `grep -q 'func TestSettleRefusesMissingCommitAuthority' internal/cli/settle_test.go && go test -count=1 ./internal/cli -run '^TestSettleRefusesMissingCommitAuthority$'` — the public command applies the same contract.
- `grep -q 'func TestCurrentRecordPermitsItsDeclaredOperations' internal/authorization/authorization_test.go && go test -count=1 ./internal/authorization -run '^TestCurrentRecordPermitsItsDeclaredOperations$'` — the repository's own record still authorizes implementation, commit and push.

## References

- `_prd.md` → Core Features 5; Goals 2; User Stories 1.
- `_techspec.md` → Implementation Design: Where operation authority is enforced; Coverage Map Core Feature 5.
- `qa/qa-report-2026-09-09.md` → F-002.

## Result

Implemented operation-specific authority checks against the ancestor grant.
Implement preflight resolves the committed Spec record before profile proof or
Run creation and carries that immutable resolution into the Daemon. Strict
authoring, Task dispatch, Task and QA commits, Spec Run push, and Settle now ask
for only the operation at their own boundary. Every refusal names the required
operation and the record read. Settle repeats the check at its commit seam, and
commit refusal becomes an unresolved Task outcome without calling the
Committer. No check was added for `pull_request`, `merge`, or `release`.

Focused red evidence before the production change:

- `rtk go test ./internal/speccheck ./internal/cli` reached 1,491 passing tests
  and the two intended regressions: strict checking emitted no
  `SC-TOOLING-UNAPPROVED` finding, and public Settle exited 0 without `commit`
  authority.
- `rtk go test ./internal/daemon` reached 294 passing tests and the intended
  dispatch, commit, and push regressions: dispatch entered work, the Task
  completed without `commit`, and Final Push reached the Pusher without
  `push`.

Focused post-change checks:

- `rtk go test ./internal/authorization ./internal/speccheck ./internal/daemon ./internal/cli`
  passed 1,793 tests across all four changed package boundaries.
- `rtk go test ./internal/spec` passed 425 tests for the shared canonical
  record resolver and its existing authorization aliases.
- `rtk git diff --check` exited 0.
- `rtk rg -n 'RequireOperation\(|Permits\(' internal --glob '!**/*_test.go'`
  found production asks for `implement`, `commit`, and `push` at the strict,
  CLI, Daemon Task, QA commit, Settle, and Final Push seams.
- `rtk rg -n 'AuthorizationOperation(PullRequest|Merge|Release)' internal/cli internal/daemon internal/speccheck --glob '!**/*_test.go'`
  exited 1 with no matches, which is the expected negative search result.

Acceptance evidence:

1. `TestStrictCheckRefusesMissingImplementAuthority` exercises a real Task
   Graph: the approved record without `operations` yields one strict
   `SC-TOOLING-UNAPPROVED` finding naming `implement` and the record, while the
   same Spec with `implement` yields none.
2. `TestDispatchRefusesMissingImplementAuthority` proves the Daemon leaves no
   Agent request, Agent Session attempt, Verification call, Task Worktree,
   commit, push, Run Event, Run-state change, Task-status change, or HEAD
   change. `TestRunImplementRefusesMissingImplementAuthorityBeforeRun` proves
   the public preflight also leaves no profile probe, Run Database, Run
   Worktree, Agent work, Git status change, commit, or push.
3. `TestCommitAndPushAuthorityAreSeparate` proves `implement` without `commit`
   settles the Task unresolved with no Committer call, while `implement` plus
   `commit` creates the Task commit and independently refuses Final Push when
   `push` is absent.
4. `TestSettleRefusesMissingCommitAuthority` proves the public command exits at
   preflight before Verification, Task status mutation, staging, commit, Run
   creation, or HEAD change when `commit` is absent.
5. `TestCurrentRecordPermitsItsDeclaredOperations` reads
   `docs/specs/0119-spec-contained-authorization/_authorization.md` from this
   repository and proves its operative record permits `implement`, `commit`,
   and `push`.
6. The production negative search above finds no `pull_request`, `merge`, or
   `release` operation check in CLI, Daemon, or Spec-check paths.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains the settlement owner.
