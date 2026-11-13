---
task: task_02
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: pending
type: backend
complexity: medium
---

# Task 02: The delivery authorization of an archived Spec is read without a Spec Root directory

## Overview

The Delivery Queue reads an archived Spec's delivery authorization from Git
history at the commit before its archive, so it must not need the Spec Root
directory in the item worktree. After a conflict resolution that removes the
last file under the Spec Root, Git removes the directory too; today the read
fails on the strict Spec Root resolver, and below it the Spec authorization
reader runs Git inside the missing directory. This answers item 1 of the
Backlog Entry reported by quaestor on 2026-11-12
([`references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md`](references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md)).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree,
   with the archive fixture's `docs/specs` removed,
   `commandDeliveryWorkflow.Authorization` fails with `specs.root resolved to
   ..., which does not exist`, and `spec.ReadSpecAuthorization` with that root
   is `unresolved`.
2. MUST change `commandDeliveryWorkflow.Authorization` in
   `internal/cli/deliver_workflow.go` as `_techspec.md` → API Contract 3 says:
   resolve the Spec Root with `roundconfig.ResolveSpecsRootOrAbsent`, with the
   comment it names. MUST NOT change the pre-archive revision it reads.
3. MUST add `authorizationSpecsRootLocation` to `internal/spec/authorization.go`
   and call it from `resolveAuthorizationLocation` for the Spec Root, as API
   Contract 3 says: an absent Spec Root resolves through its nearest existing
   ancestor with `os.Lstat`, the missing components joined to the prefix in
   order; an existing root resolves exactly as today; nothing is created on
   disk. `authorizationRepositoryLocation` and the project-root resolution
   stay unchanged.
4. MUST create `internal/spec/authorization_absent_root_test.go` with
   `TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree`
   and add `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot` to
   `internal/cli/deliver_operator_archive_test.go`, each calling
   `t.Parallel()` first, with the fixtures and assertions `_techspec.md` →
   Testing Approach lists; MUST NOT change the existing tests in that file.
5. MUST change `docs/user-guide/commands/deliver.md` with Text B of API
   Contract 6 (b), word for word, at the anchor it names.
6. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change `ResolveSpecsRoot`, `ResolveSpecsRootOrAbsent`, any other
   caller of `ReadSpecAuthorization`, or any exit code; no test reads the
   host's home or real configuration or reaches the network.

## Subtasks

- [ ] Record the characterization.
- [ ] Read the delivery authorization with the absent-tolerant resolver and resolve an absent Spec Root through its nearest existing ancestor.
- [ ] Add the two tests.
- [ ] Add Text B to the guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The delivery authorization of an archived Spec whose Spec Root left the working tree grants `push`, `pull_request` and `merge`, and the read creates no directory.
- [ ] The Spec authorization reader grants at a revision with the Spec Root absent, and an existing Spec Root reads as before.
- [ ] The `deliver` guide says so.

## Context

- instruction: `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`
- instruction: `internal/config/config.go`
- instruction: `internal/spec/authorization_test.go`
- instruction: `internal/cli/deliver_archive_links_test.go`
- creates: `internal/spec/authorization_absent_root_test.go`
- interface: `internal/spec/authorization.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_operator_archive_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `spec.ReadSpecAuthorization`'s Spec Root resolution — consumers: `./internal/spec`, `./internal/cli`
- `commandDeliveryWorkflow.Authorization` — consumers: `./internal/cli`
- The `deliver` guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run 'Authorization|Authority' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree " || { echo 'missing pass: TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree' >&2; exit 1; }` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and every authorization test of `./internal/spec` still passes.

- `out="$(go test -count=1 -timeout 40m -v -run 'Authorization|OperatorArchive|ArchivedRetry|DeliverStart|DeliverPlan' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot " || { echo 'missing pass: TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot' >&2; exit 1; }` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and the authorization, archived retry, start and plan tests of `./internal/cli` that read Spec authorization still pass.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'The grant is read from Git history, so it needs no Spec Root in the working tree' || { echo 'deliver guide: Text B missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the sentence is absent; afterwards the guide carries it and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 2; Core Features 2 and 5; Success Metric 2; Success Metric 5
- `_techspec.md` → Current behavior 4; API Contracts 3 and 6; Invariant 3; Testing Approach; Build Order 2
- ADR-0398
- [The adopted quaestor Backlog Entry](references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md)
