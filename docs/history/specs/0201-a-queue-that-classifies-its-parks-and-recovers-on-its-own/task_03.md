---
task: task_03
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: completed
type: backend
complexity: high
---

# Task 03: A retry resumes an item the operator archived after an environment-only partial

## Overview

On 2026-09-30 Spec 0192's QA gate ended `partial` only because the QA sandbox could not reach GitHub. The operator satisfied the row and archived with a QA Archive Override, and `deliver retry` refused with "archived item head differs", so the item was delivered by hand (#296). This Task parks such a Run as `qa-environment-partial`, lets a retry resume the operator-archived item at the review, and reads the delivery authorization before the archive commit instead of at `HEAD^`.

## Requirements

1. MUST make `RunSpec` read, after an unresolved Run, the newest QA Report of the Spec at the tip of the Run Branch, and set `RunResult.QAEnvironmentPartial` when its verdict is `partial`, `rows_blocked_finding` is `0` and `rows_blocked_environment` exceeds the pre-PR Pull Request rows. A missing or unreadable report MUST leave it false.
2. MUST make the engine park `qa-environment-partial` instead of `run-unresolved` when that field is set, and classify it as `environment` with the next command the TechSpec's Park Classes table states.
3. MUST make `InspectItem` set `ItemState.QAOverride` from the archived `_prd.md` frontmatter value `qa_override: true`.
4. MUST add the optional `ItemHistory` dependency (`Descends`, `RunStart`) and implement it in `internal/cli/deliver_workflow.go`.
5. MUST make `Retry` accept an archived item whose head differs from its candidate only when the blocker is `qa-environment-partial`, `QAOverride` is true, `ItemHistory` is set, and the head descends from the anchor: the last candidate commit, or the Run start of `item.RunID` when no candidate is recorded. The accepted item MUST append the head to its candidate commits and target `reviewing`. Every other archived item whose head moved MUST be refused with today's text, and `TestRetryRefusesAnArchivedItemWhoseHeadMoved` MUST pass unchanged.
6. MUST make `Archive` return `ArchiveResult.AlreadyArchived` when the reviewed head holds the archived Spec and not the active one, and make the engine move to `gating` without appending a commit.
7. MUST make `Authorization` read the grant at the parent of the newest first-parent commit that deleted the active `_prd.md`, so a commit after the archive does not hide the grant. When the archive commit is the head, the result MUST equal today's.
8. MUST describe the `qa-environment-partial` park and the operator-archived retry in `docs/user-guide/commands/deliver.md`.
9. MUST NOT change what a QA Archive Override records, what carry-forward proves, or how a non-archived item is retried.

## Subtasks

- [ ] Detect an environment-only partial from the Run Branch and park it.
- [ ] Read the QA override of the archived Spec.
- [ ] Accept the operator-archived retry and pass the archive stage through.
- [ ] Read the authorization before the archive commit.
- [ ] Describe the behavior in the delivery command guide, with a test for each criterion.

## Acceptance Criteria

- [ ] A Run Branch whose newest QA Report is an environment-only partial parks `qa-environment-partial`; a finding-blocked partial or a missing report parks `run-unresolved`.
- [ ] A retry of an item archived with a QA override after that park resumes at `reviewing`, the archive stage adds no commit, and the item reaches `gating`.
- [ ] The same retry without the override, or with a head that does not descend from the anchor, is refused with today's text.
- [ ] With an operator commit after the archive commit, the delivery authorization still reads the pre-archive grant.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/park_class.go`
- creates: `internal/delivery/operator_archive_retry_test.go`
- instruction: `internal/delivery/retry_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_operator_archive_test.go`
- creates: `docs/user-guide/commands/deliver.md`
- instruction: `internal/spec/qa.go`
- instruction: `internal/spec/archive.go`
- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial|TestRetryResumesAnOperatorArchivedItemAtReview|TestRetryRefusesAnOperatorArchiveWithoutAQAOverride|TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor|TestTheArchiveStagePassesAnAlreadyArchivedSpec|TestRetryRefusesAnArchivedItemWhoseHeadMoved)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial TestRetryResumesAnOperatorArchivedItemAtReview TestRetryRefusesAnOperatorArchiveWithoutAQAOverride TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor TestTheArchiveStagePassesAnAlreadyArchivedSpec TestRetryRefusesAnArchivedItemWhoseHeadMoved; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five new named tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch|TestInspectItemReadsTheQAOverrideOfTheArchivedSpec|TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit|TestArchiveReportsASpecAlreadyArchivedAtTheReviewedHead)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch TestInspectItemReadsTheQAOverrideOfTheArchivedSpec TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit TestArchiveReportsASpecAlreadyArchivedAtTheReviewedHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "qa-environment-partial" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "qa-environment-partial" >&2; exit 1; }` — expected: exit 0; before this Task the phrase is not in the delivery command guide, so the command fails.

## References

- `_prd.md` → Goal 3; User Story 4; Core Feature 3; Success Metric 3; Declared breaks
- `_techspec.md` → Environment-only partial and the operator-archived retry; API Contract 3; Testing Approach 3; Build Order 3
- ADR-0154; ADR-0170; ADR-0053; ADR-0158

## Result

Implemented this Task's recovery slice; settlement and declared Verification
remain owned by the Daemon. The incoming `in_progress` status was preserved.
The baseline had only the Daemon's status change in this Task file.

An unresolved Run now reads its newest QA Report from the Run Branch through
the shared QA reader. Only a partial with zero finding-blocked rows and more
environment-blocked rows than pre-PR Pull Request rows sets
`QAEnvironmentPartial`. The engine parks it as `qa-environment-partial`, with
the TechSpec's environment recovery action. Missing or unreadable evidence
retains `run-unresolved`.

Archived item inspection reads `qa_override` from PRD frontmatter. The optional
`ItemHistory` dependency proves ancestry and reads the recorded Implement Run's
starting head. A moved operator-archived candidate is accepted only for the
new park, with the override and ancestry proof; it appends the head and resumes
review. Already-archived reviewed Specs pass through archival to gating without
adding a commit. Authorization now reads the parent of the newest first-parent
active-PRD deletion, preserving the grant across later operator commits.

The delivery command guide documents the park, recovery sequence, ancestry
anchors, review/gating re-entry, and pre-archive authorization. QA Archive
Override recording, Task Carry-Forward proof, and active-Spec retry behavior
were not changed. No other Task or Task Graph file was edited; no commit, push,
or Pull Request was made.

### Acceptance evidence

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Environment-only Run Branch partial parks under the new blocker; finding-blocked or missing QA retains the old blocker | `TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial` exercises both engine outcomes and the environment action. `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch` exercises the outcome/QA-read boundary with real isolated Run Branches: qualifying partial, finding-blocked, missing, malformed, pass, pre-PR-only, and environment beyond pre-PR. It selects the newest report while the item checkout stays at its earlier head. No Agent executor is launched by this focused test. |
| Operator-archived retry resumes review, adds no archive commit, and reaches gating | `TestRetryResumesAnOperatorArchivedItemAtReview` covers candidate and Run-start anchors, records the operator head, avoids carry-forward, then advances through review and archive to gating without appending another candidate. `TestTheArchiveStagePassesAnAlreadyArchivedSpec` checks persisted gating and unchanged candidate history. `TestArchiveReportsASpecAlreadyArchivedAtTheReviewedHead` uses a real archived repository and proves unchanged HEAD and a clean worktree. `TestInspectItemReadsTheQAOverrideOfTheArchivedSpec` covers true and absent override metadata. |
| Missing override or non-descending head retains the archived-head refusal | `TestRetryRefusesAnOperatorArchiveWithoutAQAOverride` and `TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor` assert the existing refusal text and unchanged persisted item. Additional tests refuse missing/unreadable history with the same archived-head text; ancestry-read diagnostics are logged without mutating the item. `TestOperatorArchiveHistoryReadsTheRunStartAndAncestry` exercises real ancestry, reversed ancestry, unrelated histories, and a missing revision. The existing `TestRetryRefusesAnArchivedItemWhoseHeadMoved` was left byte-identical and passed in the focused delivery package check. |
| Operator commit after archival does not hide delivery authorization | `TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit` checks both archive-at-HEAD and a later commit that changes the archived grant, proving the active pre-archive grant remains the source. Its repository includes two archive transitions to check selection of the newest deletion. |

### Focused checks

- Red starting signal: `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/delivery -run TestOperatorArchiveRetryRequiresHistory -count=1` exited 1 before implementation: the new result fields, blocker and history dependency were absent.
- Final `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/delivery -count=1` exited 0 (`ok`, 1.110s), including the unchanged archived-head refusal regression.
- Final `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'Test(RunSpecReports|InspectItemReads|TheDeliveryAuthorization|ArchiveReports|OperatorArchiveHistory|InspectItemReports|ResumeAccepts|ResumeRefuses|DeliveryRunResult)' -count=1` exited 0 (`ok`, 2.615s).
- The first incremental attempt was interrupted by sandbox network policy for `cafe.github.com` and returned no check output. The same required `rtk make verify-incremental GOCACHE=/tmp/roundfix-task03-gocache` ran with elevated execution permission and exited 0: formatting, vet, package tests, skill sync/check and build passed. The elevated check was repeated after the final refusal-text adjustment and again exited 0. Output: `/tmp/roundfix-task03-incremental-elevated.log`.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The Task's declared Verification commands were not run. No follow-up slice was
added; the conflict-recovery and shipped-skill changes remain with Task 04.
