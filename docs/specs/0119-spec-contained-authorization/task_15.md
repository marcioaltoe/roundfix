---
task: task_15
spec: 0119-spec-contained-authorization
status: completed
type: backend
complexity: medium
---

# Task 15: Withdraw the execution boundary from this Spec

## Overview

The maintainer promoted the execution trust boundary in full to its own Spec on
2026-09-10, after the terminal QA gate proved it carries a performance
regression on top of two adversarial bypasses. This slice removes it from this
Spec's delivery so the record half can close. It is a removal, not a rewrite:
Task 06's delivery is superseded by a scope decision, and authored commands must
execute exactly as they did before this Spec touched them.

## Requirements

1. MUST remove the authored-command source decision and every production caller
   of it, so no path consults it before executing an authored command: the
   pre-dispatch probe, the post-Agent verification call, the Settle path, and
   the read-only checker's probing path.
2. MUST restore each of those call sites to the behavior it had before this
   Spec, so an authored command executes exactly as it did at the delivery
   target's revision. Do not leave a disabled flag, a dead parameter, or a
   check that always permits.
3. MUST NOT spawn a subprocess per authored command anywhere in the Task cycle
   after this change, which is the regression being withdrawn.
4. MUST preserve the record half completely: the typed reader in
   `internal/authorization`, the operation vocabulary and every boundary that
   asks it, the constraint reader's refusal, the Governed Path set, the
   changed-path audit's resolved reference, and the suite guard's single
   parser. Removing the execution boundary must not weaken any of them.
5. MUST remove the `SC-SOURCE-UNTRUSTED` token and scope its glossary entry to
   the code that remains, so the glossary does not define a token the CLI
   cannot emit.
6. MUST leave the Daemon package passing as one concurrent run, including the
   four Task-cycle cases that miss their deadlines today. Do not satisfy this
   by running them in isolation, and do not raise a time budget or weaken a
   deadline: the cost must go, not the signal.
7. MUST NOT change any Task status, and MUST NOT edit a sibling Task file.
   Task 06 completed; its delivery is withdrawn by this scope decision, and the
   PRD records that.

## Subtasks

- [ ] Remove the source decision and its four production call sites.
- [ ] Restore each call site's pre-Spec execution behavior.
- [ ] Remove the coined execution token and scope the glossary entry.
- [ ] Prove the record half is untouched.
- [ ] Prove the full concurrent suite passes.

## Acceptance Criteria

- [ ] No production file references the authored-command source decision, and
      the file that implemented it is gone rather than orphaned.
- [ ] No production code under `internal/daemon` or `internal/speccheck` spawns
      a process per authored command; a search for the spawn site finds none in
      the execution path.
- [ ] An authored command from a modified or untracked Spec artifact executes,
      as it did before this Spec, because no execution gate remains.
- [ ] Every record-half contract still passes: the typed reader, the operation
      boundaries, the constraint refusal, the governed set, the audit's resolved
      reference, and the suite guard's single parser.
- [ ] `CONTEXT.md` defines only the refusal code the CLI still emits.
- [ ] The Daemon package passes as one concurrent run, including the four
      Task-cycle cases named in the QA report's F-005, with no time budget
      raised and no deadline weakened.

## Context

- interface: `internal/speccheck/verification_source.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/spec_check.go`

## Verification

- `matches="$(grep -rln 'AuthorizeAuthoredCommands\|ProbeAuthoredCommands\|SC-SOURCE-UNTRUSTED' internal/ --include='*.go' || true)"; test -z "$matches" || { printf '%s\n' "$matches"; exit 1; }` — no production or test file references the withdrawn boundary; it fails today.
- `test ! -f internal/speccheck/verification_source.go` — the implementation file is removed rather than left orphaned.
- `grep -q 'SC-SOURCE-UNTRUSTED' CONTEXT.md && exit 1; grep -q 'SC-TOOLING-UNAPPROVED' CONTEXT.md` — the glossary defines the code that remains and no longer defines the one the CLI cannot emit.
- `test ! -f internal/speccheck/verification_source.go || exit 1; go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck ./internal/suiteguardcontract` — every record-half package still passes once the boundary is gone. The guard reads the removal, so this preservation check cannot pass before the work.
- `test ! -f internal/speccheck/verification_source.go || exit 1; go test -count=1 ./internal/daemon` — the package that carried the regression passes as one concurrent run, including the four Task-cycle cases that miss their deadlines today. The guard reads the removal, and the check names the package rather than the whole repository: a whole-repository gate in a Task's Verification makes the Task undispatchable on a red tree, which is the tree this Task exists to repair.

## References

- `_prd.md` → Core Features 6; Goals 3; Promoted work and known limitations; Decisions: Declared intentional breaks 3.
- `_techspec.md` → Vocabulary Contract; Coverage Map.
- `qa/qa-report-2026-09-10-02.md` → F-005.

## Result

Implementation:

- Removed the authored-command source decision, its diagnostic token, and its
  pre-dispatch, post-Agent, Settle, and read-only-check probing callers. The
  original `ProbeCommands` and direct Verifier paths execute authored commands
  again without a source gate or dead compatibility parameter.
- Removed the source-boundary tests and execution-approval fixture data. The
  public Spec Check and Settle tests now exercise commands from modified Task
  artifacts and observe that the commands run.
- Kept the authorization record reader, operation checks, constraint refusal,
  Governed Path set, resolved-reference audit, and suite-guard reader in place.
- Scoped `Grant Refusal Code` in `CONTEXT.md` to
  `SC-TOOLING-UNAPPROVED`, the refusal the CLI still emits.

Focused checks:

- `rtk go test -count=1 ./internal/cli -run '^(TestSpecCheckRunVerification|TestSettleExecutesCommandFromModifiedSource|TestSettleRefusesMissingCommitAuthority|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementInteractiveInputListsConfiguredExternalSpecRoot|TestRunImplementRefusesMissingImplementAuthorityBeforeRun)$'` — passed, 12 tests.
- `rtk go test -count=1 ./internal/daemon -run '^(TestProbeCommands|TestTaskCyclePromptContainsBundleWithoutReferencedBodies|TestDispatchRefusesMissingImplementAuthority|TestCommitAndPushAuthorityAreSeparate|TestTaskCycleSchedulesIndependentWaveWithConcurrencyCap|TestTaskCycleVerificationCapacityTwoOverlapsReadyAttemptsWithoutPermitLoss|TestTaskCycleIntegratedVerificationCapacityOneBoundsConcurrentTaskWorktrees|TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback)$'` — passed, 12 tests. The four F-005 cases ran together in this package invocation.
- `rtk go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck ./internal/suiteguardcontract -run '^(TestCurrentRecordPermitsItsDeclaredOperations|TestAuthorizationReaderClassifiesGrantState|TestAuthorizationReaderTypesPermittedOperations|TestAuthorizationReaderRefusesEscapingPaths|TestAuthorizationReaderResolvesPreservedHistoricalRecords|TestConstraintReaderCharacterizesGrantCitation|TestConstraintsResolveSpecContainedRecord|TestConstraintsRefuseNonOperativeGrant|TestConstraintsAcceptHonestProposalDeclaration|TestGovernedPath|TestGovernedSetCoversOwnedShippedTemplates|TestGovernedSetOnlyGrows|TestEveryBoundedPathIsGoverned|TestAuditReadsTheResolvedReference|TestCleanupRegenerationDiscovery|TestSanctionedRegenerationResolvesLegacyRecordsWithoutFrontmatter|TestSanctionedRegenerationReadsArchivedSpecGrants|TestSanctionedRegenerationRejectsNonOperativeRecords|TestSanctionedRegenerationSetOnlyGrows)$'` — passed, 71 tests across four packages.
- `rtk go test -count=1 ./internal/cli -run '^(TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal|TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree|TestSettleVerificationFailureKeepsHookRefusedWorkInTaskWorktree)$'` — passed, 3 tests.
- `rtk go test -count=1 ./internal/daemon -run '^(TestRefusedReportDoesNotBlockItsSuccessor|TestWriteMechanicalQAReportWritesThePreconditionRefusal|TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt|TestHookRefusalRecovery|TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged|TestTaskCycleQAReportExternalProceedsWithoutStaging)$'` — passed, 17 tests.

Acceptance evidence:

1. A post-edit source search found no
   `AuthorizeAuthoredCommands`, `ProbeAuthoredCommands`,
   `AuthoredCommandSourceRequest`, source-condition type, or
   `SC-SOURCE-UNTRUSTED` reference under `internal/`; Git records
   `internal/speccheck/verification_source.go` as deleted.
2. The removed source reader owned the extra Git subprocess boundary. The only
   Git subprocesses found in the inspected Task engine are the preserved
   changed-path audit's commit log and diff-tree reads, outside authored-command
   execution; the pre-work probe and post-Agent Verification call the Verifier
   directly.
3. `TestSpecCheckRunVerification/executes_an_opted-in_command_from_an_edited_source`
   and `TestSettleExecutesCommandFromModifiedSource` passed and independently
   observed their modified Task commands' filesystem effects.
4. The 71-test authorization preservation selection passed. A source search
   still locates `parseAuthorizationRecord`,
   `splitAuthorizationFrontmatter`, and `authorizationFrontmatterMapping` only
   in `internal/authorization/authorization.go`.
5. A glossary search finds `SC-TOOLING-UNAPPROVED` in `CONTEXT.md` and no
   withdrawn source-refusal token there.
6. The four F-005 Task-cycle cases passed together with their original
   deadlines, and a changed-line search found no timeout or deadline edit. The
   complete concurrent `./internal/daemon` command remains part of the Task's
   declared Verification and was not run in this Daemon-assigned Agent turn.

The Daemon-owned Verification commands were not run; the Daemon retains the
full-package verdict and Task settlement.
