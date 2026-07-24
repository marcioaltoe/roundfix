---
task: task_09
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 09: Apply approved Baseline Plans

## Overview

Deliver the non-interactive apply path from a portable Plan Document through
postimage verification. The command mutates only the exact approved plan,
creates immutable root backups, and reports a verified or recoverable outcome.

## Requirements

1. MUST implement `roundfix baseline apply` with the approved plan, repository,
   confirmation, format, stdout, stderr, and exit contracts.
2. MUST strictly parse the supplied plan, verify its Plan Digest, Git lineage,
   catalog identity, complete bounded preimage, and derived projections without
   recalculating or substituting another plan.
3. MUST create content-addressed backups exclusively and accept an existing
   backup only when its bytes match its full digest name.
4. MUST apply through the recoverable transaction and verify every postimage,
   carrier relationship, manifest identity, retention record, and resolved
   audit state.
5. MUST report formatter and repository Verification commands only as
   recommendations and execute neither.
6. MUST treat an empty reapply as a verified idempotent success.

## Subtasks

- [x] Implement strict apply request and Plan Document validation.
- [x] Integrate immutable backups with the transaction boundary.
- [x] Apply exact postimages and run Baseline-owned verification.
- [x] Render text/JSON success, stale, refusal, rollback, and recovery results.
- [x] Add real-CLI apply, cross-clone, and empty-reapply tests.

## Acceptance Criteria

- [x] An incorrect confirmation digest causes zero writes and an actionable non-zero result.
- [x] A stale consulted input or mutation target causes zero writes and requests a new plan.
- [x] An approved plan applies in another matching clone and fails in an unrelated lineage.
- [x] Root carriers have exact immutable backups included in the applied plan.
- [x] Postimage failure rolls back all changed paths or reports incomplete rollback as exit 1.
- [x] Baseline verification never runs repository formatter, Verification, dependency, database, or network commands.
- [x] Empty reapply produces no managed-file delta.

## Context

- instruction: `docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md`
- instruction: `docs/adr/0073-baseline-apply-uses-a-recoverable-multi-file-transaction.md`
- interface: `internal/cli/profiles_configure.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineApplyCommand|TestApplyExactDigest|TestApplyStalePreimage|TestApplyCrossClone|TestImmutableRootBackup|TestBaselineVerification|TestEmptyReapply'` — expected: public apply, safety, verification, portability, and idempotence cases pass.
- `rtk go test -count=1 ./internal/cli -run TestBaselineApplyStdoutStderrAndExitCodes` — expected: requested output stays on stdout, diagnostics stay on stderr, and exits match the public contract.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 2–3 and 7–8; Core Features 5, 7, 12–16, 19.
- `_techspec.md` → Interfaces: Workflow and Transaction; API Contracts: Automation; Build Order 5 and 7.
- ADR-0071 → exact portable plan application.
- ADR-0073 → transaction, rollback, and recovery.

## Result

Implemented the non-interactive `roundfix baseline apply` path for one strict
`roundfix/baseline-plan/v1` document and exact `--confirm-plan` digest. Apply
validates the embedded catalog and profile identities, canonical ledger and
derived file projection, Setup Manifest bytes, Upgrade Retention accounting,
clone-stable Git lineage, every bounded preimage, root-carrier backup
relationships, and content-addressed backup collisions before visible
mutation.

The recoverable transaction now runs Baseline-owned verification before its
commit phase. It verifies every exact postimage, immutable backup, carrier
relationship, Setup Manifest identity, retention record, and resolved bounded
audit state. Any failure enters the existing reverse rollback path; incomplete
rollback remains an execution failure with recovery guidance. Formatter and
repository Verification commands are returned only as sorted
`Recommendation (not run)` values. Apply invokes no dependency, database,
network, formatter, or repository Verification command.

The CLI emits text or `roundfix/baseline-result/v1` JSON while keeping requested
results on stdout and diagnostics on stderr. Exit `0` means verified apply or
verified exact reapply, exit `1` means execution, verification, output,
rollback, or recovery failure, exit `2` means invalid input or unsafe
repository state, exit `3` means approval or stale-state action is required,
and cancellation maps to exit `130`.

Acceptance evidence:

- `TestApplyExactDigest` and
  `TestBaselineApplyStdoutStderrAndExitCodes/confirmation_refusal_is_actionable_JSON`
  prove an incorrect digest returns an actionable exit `3` result without
  changing the visible repository.
- `TestApplyStalePreimage/consulted_input` and
  `TestApplyStalePreimage/mutation_target` prove both forms of bounded drift
  return a new-plan action without changing any path.
- `TestApplyCrossClone` applies one approved plan in a matching clone and
  rejects an unrelated root-commit lineage without writes.
- `TestImmutableRootBackup` proves new backups are created from exact source
  bytes and an existing backup is accepted only when its bytes match its full
  digest name.
- `TestApplyPostimageFailureRollsBack` restores the complete visible preimage
  after an injected postimage-verification failure, while
  `TestBaselineApplyStdoutStderrAndExitCodes/incomplete_rollback_is_exit_one`
  proves an incomplete rollback is rendered as exit `1`.
- `TestBaselineVerification` replaces every repository Verification
  recommendation with an observable `touch` command, proves no marker is
  created, and confirms the commands remain recommendations.
- `TestEmptyReapply` reapplies a plan containing an already-existing immutable
  backup, reports `already applied and verified`, and observes no managed-file
  delta.
- `TestBaselineApplyCommandRealCLI` builds and executes the real binary through
  plan-file handoff and receives strict verified JSON with empty stderr.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task09-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestBaselineApplyCommand|TestApplyExactDigest|TestApplyStalePreimage|TestApplyCrossClone|TestImmutableRootBackup|TestBaselineVerification|TestEmptyReapply'`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-go-cache go test -count=1
  ./internal/cli -run TestBaselineApplyStdoutStderrAndExitCodes` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-go-cache go test -count=1
  ./internal/baseline ./internal/cli` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-go-cache go vet
  ./internal/baseline ./internal/cli` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-go-cache make verify` — passed:
  1,911 Go tests in 21 packages, both 256-test setup-context-driven suites,
  embedded asset validation, Roundfix skill check, and binary build.
- `rtk git diff --check` — passed.

The isolated `GOCACHE` keeps build artifacts inside the Task Worktree sandbox.
The Daemon remains responsible for the task file's verbatim authoritative
Verification commands. No other Task file or Task Graph manifest was edited,
and no commit, push, or pull request was created.

Verification Feedback repair:

- The cross-clone test now gives the unrelated repository a distinct root
  history. Previously, two independently initialized fixtures could create
  byte-identical root commits in the same second and therefore correctly
  resolve to the same clone-stable repository identity, making the
  "unrelated" assertion timing-dependent.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-repair-go-cache go test
  -count=50 ./internal/baseline -run '^TestApplyCrossClone$'` — passed 50
  consecutive fresh runs.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-repair-go-cache go test
  -count=1 ./internal/baseline ./internal/cli -run
  'TestBaselineApplyCommand|TestApplyExactDigest|TestApplyStalePreimage|TestApplyCrossClone|TestImmutableRootBackup|TestBaselineVerification|TestEmptyReapply'`
  — passed.
