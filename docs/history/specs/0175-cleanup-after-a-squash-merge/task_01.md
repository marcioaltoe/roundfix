---
task: task_01
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: high
---

# Task 01: Prove a merged Spec's Run against the merged head

## Overview

`inspectTerminalRun` in `internal/worktree/worktree.go` proves a terminal Run
integrated only by ancestry against its target, by a QA-report-only branch, or,
when the target branch was deleted, by `inspectDeletedTargetRunByContent`. That
last check compares the whole Run tree with the current default branch through
`compareRunContentToDefault`. A squash merge rewrites every commit, archiving
moves the Spec folder, and `main` keeps moving. The proof therefore fails for
every squash-merged Spec: the Spec 0172 Run kept "48 Run-only files, 42
differing shared files", and the Spec 0164 Run kept "11 Run-only files, 35
differing shared files". This Task adds the merged-head proof described in the
TechSpec. It reads committed Git objects of the Run's own repository and, when
given, merge records the caller took from the Run Database. The Reconcile
Command, the implement-preflight prune and the delivery owner read its result
to decide whether a Run Worktree and Run Branch may be removed. A wrong
`safe` or `superseded` therefore destroys work, and every uncertainty must
preserve.

## Requirements

1. MUST add `internal/worktree/merged_head.go` with the `MergedHead` type and
   `InspectTerminalRunMerged(ctx, run, merged []MergedHead)` of the TechSpec,
   and keep `InspectTerminalRun(ctx, run)` with its signature, behaving as
   `InspectTerminalRunMerged` with no records.
2. MUST choose the merged head H exactly as the TechSpec's "Choosing the merged
   head H" says. The record source is used only when the records for the Run's
   Spec slug name exactly one distinct head and that commit exists locally.
   Otherwise the default-branch source is used, and only when the Spec is
   archived at the default branch head.
3. MUST run the proof for a present target only after ancestry and
   `supersedingQAReport` both fail and a source is usable. Without a usable
   source, the result and reason stay exactly today's `unintegrated` and
   `Run Branch is not integrated into the target branch`.
4. MUST replace the whole-tree comparison of `inspectDeletedTargetRunByContent`
   with the proof, keeping its order: representation first, then the
   archived-Spec requirement of the default-branch source.
5. MUST settle each commit of `git rev-list <RunHead> ^<H>` as a Task commit, a
   QA Report commit or another commit, as the TechSpec's "Representation" says:
   - Task status is read from the Task file at H through the Spec's `_tasks.md`
     at H, under the active or the archived Spec path, with
     `spec.CarryForwardStatus`.
   - A QA Report commit is decided by `supersedingQAReport`.
   - The files of the other commits are compared with `git diff --name-only -z
     --no-renames`, restricted by pathspec to those files.
6. MUST never re-examine a Task commit whose Task is not `completed` at H, a Task
   commit whose `Roundfix-Spec` trailer names another Spec, or an unsuperseded
   QA Report commit as another commit. Each of them preserves the Run.
7. MUST produce the states and reason texts of the TechSpec's outcome table,
   bounded by `boundedReconciliationReason`, keeping `Run Branch content is
   fully represented on default branch "<name>"` and today's content-count
   reason byte-identical for the default-branch source. `TargetHead` MUST keep
   today's meaning: the target head for a present target, and H for a deleted
   target. `ClassifySupersededBranch` in `internal/daemon/reconcile.go` reads it
   unchanged.
8. MUST keep the records in `terminalRunReconciliationEvidence` and make
   `revalidateTerminalRunApply` re-inspect with them, so `ApplyTerminalRun`
   releases a proven Run and refuses when a record or a head changed since
   inspection.
9. MUST add `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
   with the lifecycle front matter of `docs/agents/docs-layout.md` and `status:
   accepted`. It MUST state the following:
   - A terminal Run of a merged Spec is released when every commit its Run
     Branch holds and the merged head lacks is represented there (the phrase
     `merged head`).
   - A Delivery Queue merge record comes first, and the default branch carrying
     the archived Spec is the fallback.
   - The delivery owner releases a Spec's Runs after it records a merge.
   - `--discard-superseded` stays the separate act ADR-0115 requires.
   - It extends ADR-0053's proof set without superseding ADR-0053 or ADR-0115.
10. MUST update the Run Worktree Reconciliation entry of `CONTEXT.md` so that
    `safe` and `superseded` name the merged-head proof (the phrase `merged
    head`).
11. MUST put every new test in `internal/worktree/merged_head_test.go`, using
    real Git fixtures, and MUST NOT remove or rename an existing test. An
    existing test in `internal/worktree/worktree_test.go` may change only when
    it pinned the whole-tree comparison, and the Result MUST name it.

## Subtasks

- [ ] Implement the merged-head proof and its revalidation.
- [ ] Record ADR-0161 and update the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] The Spec 0172 shape is `superseded` against a merge record and against a default branch carrying the archived Spec, and `ApplyTerminalRun` removes its Run Worktree and Run Branch.
- [ ] The Spec 0164 shape is released although the default branch later changed other files and a file the Run touched.
- [ ] An unrepresented commit, a Task not completed at H, a Task commit of another Spec, a QA Report H does not supersede, a record for another Spec and a dirty Run Worktree each keep the Run, with a reason naming the cause.
- [ ] Every existing deleted-target and superseded test named in the Verification stays green.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`
- creates: `internal/worktree/merged_head.go`
- creates: `internal/worktree/merged_head_test.go`
- creates: `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMergedHeadRecordReleasesARunContainedInTheMergedHead|TestMergedHeadRecordSupersedesARedoneTaskAndAFailedQAReport|TestMergedHeadDefaultBranchSupersedesARedoneTaskAndAFailedQAReport|TestMergedHeadDefaultBranchComparesOnlyTheRunsChangedFiles|TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved|TestMergedHeadRefusesAnUnrepresentedCommit|TestMergedHeadRefusesATaskNotCompletedAtTheMergedHead|TestMergedHeadRefusesATaskCommitOfAnotherSpec|TestMergedHeadRefusesAQAReportTheMergedHeadDoesNotSupersede|TestMergedHeadIgnoresARecordForAnotherSpec|TestMergedHeadFallsBackWhenRecordsDisagree|TestMergedHeadKeepsADirtyRunWorktree|TestMergedHeadKeepsTheUnintegratedReasonWithoutASource|TestMergedHeadApplyReleasesTheRunWorktreeAndBranch|TestMergedHeadApplyRefusesAChangedRecord|TestInspectTerminalRunSafeWhenTargetDeletedAfterSquashMerge|TestInspectTerminalRunRequiresArchivedEvidence|TestInspectTerminalRunUnintegrated|TestInspectTerminalRunUnintegratedWhenDeletedTargetHasRunOnlyFile|TestInspectTerminalRunUnintegratedWhenDeletedTargetHasDifferentSharedFile|TestInspectTerminalRunUnintegratedWhenDeletedTargetRetainsRunDeletedFile|TestInspectTerminalRunUnintegratedWhenDeletedTargetContentComparisonFails|TestInspectTerminalRunUnknownWhenDeletedTargetDefaultBranchCannotBeResolved|TestInspectTerminalRunClassifiesSupersededQAReport|TestApplyTerminalRunSuperseded|TestApplyTerminalRunStaleHeads|TestSupersededBranchRefusesAnUnreachableCommit|TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks)$" ./internal/worktree ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergedHeadRecordReleasesARunContainedInTheMergedHead TestMergedHeadRecordSupersedesARedoneTaskAndAFailedQAReport TestMergedHeadDefaultBranchSupersedesARedoneTaskAndAFailedQAReport TestMergedHeadDefaultBranchComparesOnlyTheRunsChangedFiles TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved TestMergedHeadRefusesAnUnrepresentedCommit TestMergedHeadRefusesATaskNotCompletedAtTheMergedHead TestMergedHeadRefusesATaskCommitOfAnotherSpec TestMergedHeadRefusesAQAReportTheMergedHeadDoesNotSupersede TestMergedHeadIgnoresARecordForAnotherSpec TestMergedHeadFallsBackWhenRecordsDisagree TestMergedHeadKeepsADirtyRunWorktree TestMergedHeadKeepsTheUnintegratedReasonWithoutASource TestMergedHeadApplyReleasesTheRunWorktreeAndBranch TestMergedHeadApplyRefusesAChangedRecord TestInspectTerminalRunSafeWhenTargetDeletedAfterSquashMerge TestInspectTerminalRunRequiresArchivedEvidence TestInspectTerminalRunUnintegrated TestInspectTerminalRunUnintegratedWhenDeletedTargetHasRunOnlyFile TestInspectTerminalRunUnintegratedWhenDeletedTargetHasDifferentSharedFile TestInspectTerminalRunUnintegratedWhenDeletedTargetRetainsRunDeletedFile TestInspectTerminalRunUnintegratedWhenDeletedTargetContentComparisonFails TestInspectTerminalRunUnknownWhenDeletedTargetDefaultBranchCannotBeResolved TestInspectTerminalRunClassifiesSupersededQAReport TestApplyTerminalRunSuperseded TestApplyTerminalRunStaleHeads TestSupersededBranchRefusesAnUnreachableCommit TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "^status: accepted" docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md && tr -s '[:space:]' ' ' < docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md | grep -qF -- "merged head" && grep -q "ADR-0115" docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "merged head"` — expected: exit 0; before this Task none of the new named tests exists, ADR-0161 is absent and `CONTEXT.md` never says `merged head`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Decision: when cleanup runs, and what a squash merge proves; The merged-head proof
- `_prd.md` → Core Feature 1; Success Metrics 1-2
- `_techspec.md` → API Contracts 1-2

## Result

Implemented the merged-head proof without changing the Task status or Task Graph. `InspectTerminalRunMerged` now selects one locally available Delivery Queue head for the Run's exact Spec slug, otherwise falls back to an archived Spec at the default-branch head. The proof classifies each Run-only commit once as a completed Task, a superseded QA Report, or other changed files; uncertainty preserves the Run. `InspectTerminalRun` delegates with no records, and apply revalidation retains the records plus the private proof head so a changed record or moved fallback head refuses cleanup.

Recorded ADR-0161 with accepted lifecycle metadata and updated the Run Worktree Reconciliation glossary entry to name the merged-head `safe` and `superseded` proofs.

Focused evidence by acceptance criterion:

- Spec 0172 shape and apply release: `TestMergedHeadRecordSupersedesARedoneTaskAndAFailedQAReport`, `TestMergedHeadDefaultBranchSupersedesARedoneTaskAndAFailedQAReport`, `TestMergedHeadRecordReleasesARunContainedInTheMergedHead`, and `TestMergedHeadApplyReleasesTheRunWorktreeAndBranch` passed against real Git fixtures in `go test -count=1 -run '^TestMergedHead' ./internal/worktree`.
- Spec 0164 shape and changed-file restriction: `TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved` and `TestMergedHeadDefaultBranchComparesOnlyTheRunsChangedFiles` passed in the same focused run. The former changes a file after its completed Task was delivered; the latter advances an unrelated default-branch file.
- Preservation cases: the separate real-Git tests for an unrepresented commit, incomplete Task, another-Spec Task trailer, unsuperseded QA Report, another-Spec record, disagreeing records, a missing record commit, dirty worktree, and no usable source passed. `TestMergedHeadApplyRefusesAChangedRecord` and `TestMergedHeadApplyRefusesAMovedFallbackHead` passed and kept both Git surfaces.
- Existing behavior: `go test -count=1 ./internal/worktree` passed, and the focused daemon command for `TestSupersededBranchRefusesAnUnreachableCommit` plus `TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks` passed. The focused legacy deleted-target, archived-evidence, QA-supersession, and stale-apply matrix also passed.

Existing tests changed because they pinned the retired whole-tree comparison:

- `TestInspectTerminalRunUnintegratedWhenDeletedTargetContentComparisonFails` now creates one represented ordinary commit so its injected `git diff` failure reaches the new path-restricted comparison.
- `TestApplyRunBranchCandidateRefusesACandidateWithoutEvidence` retains its name but now expects release: its squash-merged QA-only fixture has positive merged-head evidence where the whole-tree comparison previously returned none. The separate ADR-0115 Branch Disposition proof still runs before cleanup.

Focused commands and outcomes:

- `GOCACHE=/tmp/roundfix-task-0175-01-gocache go test -run '^$' ./internal/worktree` — passed (compile-only). The first attempt without the Task-scoped cache was blocked by sandbox access to the macOS Go build cache.
- `GOCACHE=/tmp/roundfix-task-0175-01-gocache go test -count=1 -run '^TestMergedHead' ./internal/worktree` — passed.
- `GOCACHE=/tmp/roundfix-task-0175-01-gocache go test -count=1 ./internal/worktree` — passed.
- `GOCACHE=/tmp/roundfix-task-0175-01-gocache go test -count=1 -run '^(TestSupersededBranchRefusesAnUnreachableCommit|TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks)$' ./internal/daemon` — passed.
- `GOCACHE=/tmp/roundfix-task-0175-01-gocache go vet ./internal/worktree` — passed.
- `GOCACHE=/tmp/roundfix-task-0175-01-gocache make verify-incremental` — passed with process-table permission. The sandboxed attempt reached the full test suite but its two owner-process integration tests were blocked by `operation not permitted`; the permitted rerun passed vet, all Go packages, skill checks, and the build.
- `git diff --check` — passed.

The authored `## Verification` command was not run; Daemon Verification owns it.

## Carry-forward provenance

- Source Run: `run_20260928T154312Z_63679c6540b34603`
- Source commit: `8cbd476b715288435afa2803371d3bc139c64daf`
