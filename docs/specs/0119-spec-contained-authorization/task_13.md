---
task: task_13
spec: 0119-spec-contained-authorization
status: pending
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
