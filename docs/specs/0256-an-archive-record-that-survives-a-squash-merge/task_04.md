---
task: task_04
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
type: backend
complexity: high
---

# Task 04: The Delivery Queue records the pull request in a Pull Request Record Commit

## Overview

The Delivery Queue archives a Spec before it opens the pull request, so its
archive commit cannot carry the number. This Task adds the Pull Request
Record Commit. At the start of the Checking stage, the queue stores the
opened pull request's number in the record through
`roundfix archive <slug> --pull-request <n>`. It proves that the commit
changes nothing else, commits it and returns to Publishing to push it. It is
verifiable on its own: the engine with a fake recorder, the command
workflow against real Git, and an end-to-end run against the fake forge whose
merged head carries the number.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12
   [an Archive Record points at a branch commit a squash merge drops](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)
   on the delivery side, implementing `_techspec.md` → API Contract 4 and API
   Contract 5 exactly.
2. MUST add `delivery.PullRequestRecorder` and `delivery.PullRequestRecord`
   as in `_techspec.md` → Interfaces, as an `EngineDependencies` field. The
   command engine MUST wire it. An engine built without one skips the step,
   as an engine without a conflict resolver skips conflict resolution.
3. MUST run the step at the start of `checkCandidate`, before the first check
   read. On a recorded commit it appends the head to the item's candidate
   commits, writes the log line of API Contract 4, and returns to Publishing
   without the repository gate. An inexact commit, or a returned head equal
   to the candidate head, parks `review-stale`. An error parks through the
   existing error path as `delivery-error`. Nothing merges while the record
   at the candidate head lacks the number.
4. MUST make the command workflow's `RecordPullRequest` idempotent and safe
   to resume, as API Contract 4 states. It returns an empty result for a
   Legacy Archive Folder and for a record that already names the number. It
   adopts an exact child commit of the candidate head found at the checkout
   `HEAD`, and commits an uncommitted record that is the only working-tree
   change without rewriting it. Otherwise it runs the Archive Command through
   the step executable and proves the result. The commit subject MUST be
   `docs: record pull request #<n> in the <slug> Archive Record`.
5. MUST add these tests, each top-level test calling `t.Parallel()` first or
   carrying a `// Sequential: <reason>` comment, with no network connection:
   - in `internal/delivery`: `TestDeliveryRecordsThePullRequestBeforeChecks`,
     `TestDeliveryRetryAdoptsAPullRequestRecordCommit`,
     `TestDeliveryParksAnInexactPullRequestRecordCommit` and
     `TestDeliveryParksWhenThePullRequestCannotBeRecorded`;
   - in `internal/cli`: `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord`,
     `TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath`,
     `TestDeliveryWorkflowAdoptsAnUncommittedRecord` and
     `TestCommandDeliveryEngineRecordsThePullRequest`. The last runs a queue
     item to merge against the fake forge and reads `pull_request` from the
     record at the merged head.
6. MUST keep `TestDeliveryArchiveStageCommitsTheArchiveRecord`,
   `TestDeliveryAcceptsAnExactRetirement` and
   `TestDeliveryEngineMergesAQueueWithoutAnOperator` passing. An existing test
   whose expected stages or candidate commits move because of the new commit
   MAY change only to that new contract. The Result MUST name each one and
   what moved.
7. MUST prove each new gate can fail. The Result MUST record two sabotages,
   restored afterwards: the proof accepting a commit that also changes another
   path, and the step skipped at Checking. For each it names the test that
   failed.
8. MUST NOT change the archive stage, its exact-archive proof, the Baseline
   delivery order, the pull request boundary or any blocker string. It MUST
   NOT push any ref other than the item branch.

## Subtasks

- [ ] Add the recorder dependency and the Checking-stage step.
- [ ] Implement the command workflow's record, adopt and commit paths.
- [ ] Implement the exactness proof.
- [ ] Add the engine, workflow and end-to-end tests.
- [ ] Record the sabotages and any existing test that moved.

## Acceptance Criteria

- [ ] A queue run against the fake forge merges a head whose record names the opened pull request.
- [ ] An interruption between opening the pull request and recording it produces exactly one Pull Request Record Commit after retry.
- [ ] A commit that changes anything but the record's number parks `review-stale`.
- [ ] A failed recording parks the item, and nothing merges without the number.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `docs/adr/0237-a-delivery-retry-records-a-merge-made-outside-the-queue.md`
- instruction: `internal/cli/deliver_archive_record_test.go`
- instruction: `internal/delivery/engine_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_test.go`
- creates: `internal/delivery/record_pull_request_test.go`
- creates: `internal/cli/deliver_record_pull_request_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryRecordsThePullRequestBeforeChecks|TestDeliveryRetryAdoptsAPullRequestRecordCommit|TestDeliveryParksAnInexactPullRequestRecordCommit|TestDeliveryParksWhenThePullRequestCannotBeRecorded|TestDeliveryEngineMergesAQueueWithoutAnOperator)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryRecordsThePullRequestBeforeChecks TestDeliveryRetryAdoptsAPullRequestRecordCommit TestDeliveryParksAnInexactPullRequestRecordCommit TestDeliveryParksWhenThePullRequestCannotBeRecorded TestDeliveryEngineMergesAQueueWithoutAnOperator; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four engine tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath|TestDeliveryWorkflowAdoptsAnUncommittedRecord|TestCommandDeliveryEngineRecordsThePullRequest|TestDeliveryArchiveStageCommitsTheArchiveRecord|TestDeliveryAcceptsAnExactRetirement)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath TestDeliveryWorkflowAdoptsAnUncommittedRecord TestCommandDeliveryEngineRecordsThePullRequest TestDeliveryArchiveStageCommitsTheArchiveRecord TestDeliveryAcceptsAnExactRetirement; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four workflow tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; User Stories 1 and 3; Core Feature 2; Success Metric 1
- [_techspec.md](_techspec.md) — Interfaces; API Contract 4; API Contract 5; Testing Approach; Risks & Considerations; Build Order 4
- ADR-0261; ADR-0223; ADR-0229; ADR-0237; ADR-0232


## Result

Implemented the delivery side of the 2026-10-12 Backlog Entry,
[an Archive Record points at a branch commit a squash merge drops](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md),
under TechSpec API Contracts 4 and 5.

The engine accepts the optional `PullRequestRecorder` dependency, wired by
`newCommandDeliveryEngine`, at the start of Checking. An exact new head is
appended to candidate commits, emits the contracted recording log line and
returns directly to Publishing. The existing gate, archive proof, pull request
boundary and blocker strings are unchanged. A missing dependency skips the
step; an inexact or unchanged returned head parks `review-stale`; a recording
error reaches the existing `delivery-error` path.

The workflow reads the candidate's record and skips Legacy Archive Folders
and an already matching number. It adopts a clean, exact single child commit,
or commits the sole uncommitted record change without rewriting its bytes.
For a fresh write it runs the Archive Command through `stepExecutable`, proves
the working record and staged tree, commits the required subject and proves
the resulting commit. The proof requires exactly the record path, an unchanged
regular-file mode and equal fields except the number and the renderer's exact
outcome shortening. It also compares the complete YAML mapping and body to
refuse unknown metadata or added prose. Body comparison accounts for the
existing Git runner stripping trailing newlines from object reads.

### Acceptance evidence

| Criterion | Focused evidence |
| --- | --- |
| Fake forge merges a head whose record names the opened pull request | `TestCommandDeliveryEngineRecordsThePullRequest` uses the real command-engine constructor, real local Git pushes and a strict local fake `gh`. It opens pull request 7, builds a squash commit in the bare origin, reads its Archive Record and asserts number 7. The queue records three candidate heads and exactly one record commit. Only `main` and the item branch exist in the forge; its merge updates `main` locally, and the queue pushes only the item branch. |
| Interruption produces exactly one record commit on resume | `TestDeliveryRetryAdoptsAPullRequestRecordCommit` resumes Checking with an existing pull request, appends the adopted head once and republishes without Implement, archive or gate replay. `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord` simulates interruption after the commit but before queue persistence: the old candidate adopts the same commit, the updated candidate returns an empty result, and `rev-list --count` remains 1. `TestDeliveryWorkflowAdoptsAnUncommittedRecord` covers staged and unstaged interrupted writes, preserving alternate YAML quoting byte for byte. |
| A commit changing anything else parks review-stale | `TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath` rejects another path, title, source revision, arbitrary outcome, wrong number, mode, grandchild, merge commit, unknown field, extra prose and changed body title. `TestDeliveryParksAnInexactPullRequestRecordCommit` asserts `review-stale`, unchanged candidate commits and no checks or merge, including a purported exact result equal to the candidate. |
| Failed recording parks and never merges without the number | `TestDeliveryParksWhenThePullRequestCannotBeRecorded` asserts `delivery-error`, no check read and zero merge effects. `TestDeliveryRecordsThePullRequestBeforeChecks` asserts recording before checks, one pull request, two pushes, the prescribed log and no second gate. The real command-engine test independently reads the number from the merged record. |

Additional focused tests cover an absent recorder through the unchanged engine
regression, Legacy Archive Folders, uncommitted unrelated/staged changes,
missing or wrong numbers, changed fields and modes, and the exact outcome
shortening produced at the record's size limit.

### Focused checks

- Initial red signal: `go test -count=1 -run '^TestDeliveryRecordsThePullRequestBeforeChecks$' ./internal/delivery`
  exited 1 before implementation because the recorder type and engine hook
  did not exist.
- `go test -count=1 -v -run '^TestDelivery(RecordsThePullRequest|RetryAdoptsAPullRequest|ParksAnInexactPullRequest|ParksWhenThePullRequest|EngineMergesAQueueWithoutAnOperator|EngineResumesWithoutDoubleEffects)' ./internal/delivery`
  exited 0. All four required new engine tests and both existing regressions
  ran and passed.
- `go test -count=1 -v -run '^(TestDeliveryWorkflow|TestDeliveryRecordProof|TestCommandDeliveryEngineRecordsThePullRequest|TestDeliveryArchiveStageCommitsTheArchiveRecord|TestDeliveryAcceptsAnExactRetirement|TestDeliveryRefusesAnInexactRetirement|TestDeliveryKeepsTheLegacyExactMove)' ./internal/cli`
  exited 0. The four required new workflow/command tests, additional refusal,
  legacy and shortening tests, and existing archive regressions ran and passed.
- `gofmt` formatted the four Go files. `git diff --check` exited 0.

Required existing tests stayed unchanged:

- `TestDeliveryArchiveStageCommitsTheArchiveRecord`: no stage or expectation
  moved; archive still commits only its exact retirement.
- `TestDeliveryAcceptsAnExactRetirement`: no expectation moved; record and
  promotion proofs retain their previous contract.
- `TestDeliveryEngineMergesAQueueWithoutAnOperator`: no stage or candidate
  expectation moved; its engine has no recorder and retains its prior order.

### Sabotage evidence

Both mutations used temporary `go test -overlay` source replacements under
`/tmp/roundfix-task04-sabotage`, leaving the checkout's production files intact.
The overlays were removed afterwards, and both affected tests were rerun
against the normal source.

1. Removed the record proof's requirement that the diff contain exactly the
   record path. Command:
   `go test -overlay /tmp/roundfix-task04-sabotage/proof.json -count=1 -v -run '^TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath$/extra_path$' ./internal/cli`.
   Exit 1: `TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath/extra_path`
   failed because the unrelated-path commit was returned with `Exact:true`.
2. Disabled the recorder block at Checking. Command:
   `go test -overlay /tmp/roundfix-task04-sabotage/checking.json -count=1 -v -run '^TestDeliveryRecordsThePullRequestBeforeChecks$' ./internal/delivery`.
   Exit 1: `TestDeliveryRecordsThePullRequestBeforeChecks` failed because the
   observed actions reached checks and merge without recording or republishing.

### Verification feedback: attempt 1

Inspected the Daemon's diagnostic artifact for the rejected
`make verify-changed` attempt. The new fake-forge test violated two existing
fixture contracts: it wrote an executable script outside the frozen inventory,
and its process-wide environment changes added a Sequential Test beyond the
CLI package's ceiling. These were fixture failures, not changes to the
recording or exactness contract.

Repaired only `TestCommandDeliveryEngineRecordsThePullRequest`: it now calls
`t.Parallel()` first and runs the real constructor scenario in a child test
process with private PATH, HOME and isolated Git configuration. The fake `gh`
uses the existing `writeScriptFixture` helper, a symlink to the compiled test
binary and a non-executable shell sidecar. Its real Git pushes, squash merge,
record-number assertions, one-record-commit assertion and ref inventory are
preserved. No production behavior, fixture inventory or package ceiling was
changed in this feedback repair.

Focused checks after the repair:

- `go test -count=1 -v -run '^TestCommandDeliveryEngineRecordsThePullRequest$' ./internal/cli`
  exited 0; both the parent and isolated child scenario passed.
- `go test -count=1 -v -run '^(TestNoTestWritesAnExecutableOutsideTheResidue|TestEveryTestInAParallelTestPackageRunsInParallel)$' ./internal/testfixture`
  exited 0; both repository guards passed and the CLI Sequential Test count
  returned to its existing ceiling of 11.
- `go test -count=1 -v -run '^(TestCommandDeliveryEngineRecordsThePullRequest|TestDeliveryWorkflow|TestDeliveryRecordProof|TestDeliveryArchiveStageCommitsTheArchiveRecord|TestDeliveryAcceptsAnExactRetirement)' ./internal/cli`
  exited 0; the revised command-engine scenario passed concurrently with the
  workflow, shortening and required archive regression tests.
- `gofmt` formatted the changed test file; `git diff --check` exited 0.

The declared Task Verification and `make verify-changed` were not rerun by
this Agent. The Daemon owns the next full Verification and Task settlement.

### Handoff boundary

Only this Task's two production files, two new test files and this Result
were changed. The pre-existing daemon status edit is preserved. No other
Task, Task Graph, archive-stage implementation, gate, or blocker string was
changed. Test repositories and the fake forge are disposable and use no
network connection. No commit, push or pull request was made for this Run.
The declared Verification commands, selected repository Verification and
Task settlement remain with the Daemon; no terminal Task verdict is claimed.

## Carry-forward provenance

- Source Run: `run_20261012T140955Z_23b372ef5c35c71c`
- Source commit: `fccf7bcfe7d0274bea01376a657d609d0f523994`
