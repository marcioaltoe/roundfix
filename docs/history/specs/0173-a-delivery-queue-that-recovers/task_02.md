---
task: task_02
spec: 0173-a-delivery-queue-that-recovers
status: completed
type: backend
complexity: medium
---

# Task 02: Project Config is committed on authority or refused aloud

## Overview

`diffSnapshots` in `internal/daemon/engine.go` skips `.roundfixrc.yml` for every Batch, Task and QA Report commit with no event and no reason. Spec 0155's `task_03` was authorized to change Project Config; the Agent made the edit, Verification passed, the Task settled completed, and its commit left the change behind in the Run Worktree until the QA gate failed the Spec. The changed paths come from the Worktree snapshots the Daemon takes around Agent work; the authority comes from the Spec's frozen authorization in the Task plan; the commit is read by Task integration, and the dropped-path record is read by the operator on the console and by the Supervisor on the Run Event Stream.

## Requirements

1. MUST make `diffSnapshots` a pure snapshot diff that no longer skips `.roundfixrc.yml`, and move the decision to each commit path.
2. MUST make `prepareTaskCommit` in `internal/daemon/task_engine.go` keep `.roundfixrc.yml` when `plan.Authorization` has outcome `spec.AuthorizationGranted` and its `Record.Paths` lists `.roundfixrc.yml`, leaving the existing governed-mutation checks to apply, and otherwise drop it as a `DroppedStagePath` with reason `Project Config outside the Spec's authorization` and `Lost: true`, so the Task settles failed with `Task commit lost output: .roundfixrc.yml (Project Config outside the Spec's authorization)` and creates no commit.
3. MUST make `commitBatch` in `internal/daemon/engine.go` and `commitQAReport` in `internal/daemon/task_engine.go` drop `.roundfixrc.yml` with reason `Project Config is never committed by a Batch or QA Report commit` and `Lost: false`, publish it through `publishDroppedStagePath`, and commit the remaining paths.
4. MUST make `publishDroppedStagePath` print `roundfix: Project Config .roundfixrc.yml omitted from the commit: <reason>` and publish the summary `Project Config .roundfixrc.yml omitted from the commit: <reason>.` for both reasons, keeping the payload fields `decision`, `path` and `reason`.
5. MUST keep a `.roundfixrc.yml` that is already dirty in the before snapshot out of every commit, and keep `TestGitCommitterExcludesProjectConfigFromBatchCommit` green without edits.
6. MUST state in the Project Config entry of `CONTEXT.md` and in the Daemon commit paragraph of `.agents/skills/roundfix/SKILL.md` that a Task commits Project Config only when the frozen authorization bounds `.roundfixrc.yml`, that otherwise the Task fails with the reason `Project Config outside the Spec's authorization`, and that Batch and QA Report commits never stage it and report the exclusion; regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
7. MUST put the new tests in `internal/daemon/project_config_commit_test.go`, asserting the committed paths, the Task's settled status and reason, and the published event payload by value.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Task whose frozen authorization bounds `.roundfixrc.yml` commits its change to that file.
- [ ] A Task whose authorization does not bound it settles failed with the lost-output reason, publishes one dropped-path event naming the path and reason, and creates no commit.
- [ ] A Batch commit and a QA Report commit omit it and publish their reason, and their other paths are committed.
- [ ] A Project Config change present before the Task started is committed by none of them.

## Context

- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/project_config_commit_test.go`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization|TestTaskOutsideTheAuthorizationFailsOnProjectConfig|TestBatchCommitReportsProjectConfigAsDropped|TestQAReportCommitReportsProjectConfigAsDropped|TestPreexistingProjectConfigChangeStaysOutOfTheTaskCommit|TestGitCommitterExcludesProjectConfigFromBatchCommit)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization TestTaskOutsideTheAuthorizationFailsOnProjectConfig TestBatchCommitReportsProjectConfigAsDropped TestQAReportCommitReportsProjectConfigAsDropped TestPreexistingProjectConfigChangeStaysOutOfTheTaskCommit TestGitCommitterExcludesProjectConfigFromBatchCommit; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "Project Config outside the Spec's authorization" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "Project Config outside the Spec's authorization" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the five new named tests exists and the phrase is documented nowhere, so the command fails.

## References

- [_techspec.md](_techspec.md) — Project Config in Daemon commits

## Result

Implemented the Daemon commit policy at each commit boundary. Snapshot diffing now reports every newly dirty path; Task commit preparation keeps Project Config only when the frozen authorization grants and bounds `.roundfixrc.yml`, while Batch and QA Report commits remove it from their stageable set and publish the shared exclusion message. The QA governed-mutation check ignores Project Config because QA can never commit that path. A Project Config path already present in the before snapshot remains outside every commit without producing a new exclusion event.

Documented the same rule in the `Project Config` glossary entry and the Roundfix assigned-Task commit contract. Regenerated the distributed Roundfix skill with `make skills-sync`; `make baseline-digests` reported `changed:false`.

Acceptance evidence:

- Authorized Task: `TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization` passed and asserted the completed settlement plus the exact commit path set containing `.roundfixrc.yml`.
- Unauthorized Task: `TestTaskOutsideTheAuthorizationFailsOnProjectConfig` passed and asserted failed settlement with `Task commit lost output: .roundfixrc.yml (Project Config outside the Spec's authorization)`, no commit, one console line, and the dropped-path and settlement event payloads by value.
- Batch and QA Report: `TestBatchCommitReportsProjectConfigAsDropped` and `TestQAReportCommitReportsProjectConfigAsDropped` passed and asserted each remaining path set, exact console line, summary, and dropped-path payload by value.
- Pre-existing change: all Task, Batch, and QA Report subtests in `TestPreexistingProjectConfigChangeStaysOutOfTheTaskCommit` passed and asserted `.roundfixrc.yml` was neither committed nor reported as a new drop.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization$' ./internal/daemon` — failed before the implementation because the commit contained only the Task file; this captured the original omission.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization|TestTaskOutsideTheAuthorizationFailsOnProjectConfig|TestBatchCommitReportsProjectConfigAsDropped|TestQAReportCommitReportsProjectConfigAsDropped|TestPreexistingProjectConfigChangeStaysOutOfTheTaskCommit)$' ./internal/daemon` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestResolveCycleStagesOnlyAgentTouchedPaths|TestResolveCycleDropsExecutableFileAndCommitsRemainingBatchPaths|TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths|TestQACommitDropsExecutableFileAndCommitsRemainingPaths|TestGovernedMutationDetectionUsesTheUnfilteredSnapshot|TestGitCommitterExcludesProjectConfigFromBatchCommit)$' ./internal/daemon` — passed.
- `rtk make skills-sync-check` and `rtk git diff --check` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental` — the sandboxed run reached the package tests but could not read the host process table in two force-stop integration tests; the unchanged rerun with host process-table permission passed, including vet, all package tests, skill checks, and the build.

The Task's authored `## Verification` command was not run; Daemon Verification owns it.
