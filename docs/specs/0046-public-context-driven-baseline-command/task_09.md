---
task: task_09
spec: 0046-public-context-driven-baseline-command
status: pending
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

- [ ] Implement strict apply request and Plan Document validation.
- [ ] Integrate immutable backups with the transaction boundary.
- [ ] Apply exact postimages and run Baseline-owned verification.
- [ ] Render text/JSON success, stale, refusal, rollback, and recovery results.
- [ ] Add real-CLI apply, cross-clone, and empty-reapply tests.

## Acceptance Criteria

- [ ] An incorrect confirmation digest causes zero writes and an actionable non-zero result.
- [ ] A stale consulted input or mutation target causes zero writes and requests a new plan.
- [ ] An approved plan applies in another matching clone and fails in an unrelated lineage.
- [ ] Root carriers have exact immutable backups included in the applied plan.
- [ ] Postimage failure rolls back all changed paths or reports incomplete rollback as exit 1.
- [ ] Baseline verification never runs repository formatter, Verification, dependency, database, or network commands.
- [ ] Empty reapply produces no managed-file delta.

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
