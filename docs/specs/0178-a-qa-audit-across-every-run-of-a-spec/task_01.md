---
task: task_01
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: pending
type: backend
complexity: high
---

# Task 01: Every Run's Task commits reach the audit

## Overview

`qaMechanicalRequest` in `internal/daemon/task_engine.go` builds the QA
mechanical stage's request from `plan.HeadSHA`, the checkout head when the
current Run started. `mechanicalTaskCommits` reads `git log <HeadSHA>..HEAD`
and keeps only the newest commit of each Task, so a Task commit made by an
earlier Run of the same Spec, and an older commit of a Task, never reach the
authorization audit. The same revision is the delivery target the stage reads
each grant at, so a grant widened on the Spec branch before the current Run
authorizes the Tasks after it. The commits and grants are read from the Run
Worktree's Git history, which every Run of the Spec wrote; the audit's result is
read by the QA gate and, through the seeded QA Report, by the QA Agent and the
maintainer who decides whether the Spec ships.

## Requirements

1. MUST add `internal/daemon/qa_delivery_base.go` with
   `qaDeliveryBase(ctx context.Context, plan TaskPlan) (string, bool, error)`.
   It reads the current branch of `plan.WorkDir` with
   `git symbolic-ref --quiet --short HEAD` (empty when detached), calls
   `preflight.DetectDefaultBranch`, resolves `refs/remotes/origin/<name>` and,
   when that ref is absent, `refs/heads/<name>`, and returns
   `git merge-base <ref> HEAD` run in `plan.WorkDir` with `true`. An
   undetermined default branch, a ref that resolves under neither name, or no
   merge base (exit `1`) MUST return `("", false, nil)`; any other Git failure
   MUST return an error wrapped with `%w` that names the operation.
2. MUST make `qaMechanicalRequest`, when `plan.HeadSHA` is set, call
   `qaDeliveryBase`. With a resolved base it MUST pass the base to
   `speccheck.ResolveMechanicalAuthorization`, set `DeliveryTargetRevision` to
   it, and read Task commits from it. Unresolved, it MUST keep `plan.HeadSHA`
   for all three and set the new `MechanicalRequest.TaskCommitsFromRunStart`
   field to `true`. A Git error from `qaDeliveryBase` MUST be returned, never
   read as unresolved.
3. MUST make `mechanicalTaskCommits` take the range start as a parameter and
   keep every commit of each non-QA Task of the Spec in `<start>..HEAD`, not
   only the first per Task in log order, still dropping a commit that intersects
   no governed path, and return the kept commits oldest first.
4. MUST add `TaskCommitsFromRunStart bool` to `speccheck.MechanicalRequest` and
   the exported constant `MechanicalSkipEarlierRunTaskCommits = "Task commits of earlier Runs"`
   to `internal/speccheck/mechanical.go`. When the field is set,
   `detectMechanicalAuthPaths` MUST record the skip
   `(DetectorMechanicalAuthPaths, MechanicalSkipEarlierRunTaskCommits)` and
   still audit the commits it was given. The authorizing-revision rule and the
   bounded-path rule MUST stay unchanged.
5. MUST update `TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit` in
   `internal/daemon/task_engine_test.go`, which this change invalidates: its
   fixture commits the Task commits on a work branch created after the initial
   commit, sets `refs/remotes/origin/main` and the `refs/remotes/origin/HEAD`
   symbolic ref, and its delivery target assertion expects the Delivery Base,
   which is still the initial commit. Its name and its two subtests stay.
6. MUST put new tests in `internal/daemon/qa_every_run_audit_test.go`. Each
   fixture is a disposable repository with `refs/remotes/origin/main` and the
   `refs/remotes/origin/HEAD` symbolic ref, a Spec whose PRD cites a grant, and
   Task commits carrying `Roundfix-Spec` and `Roundfix-Task` trailers on a work
   branch; each mechanical case runs `speccheck.RunMechanicalStage` on the
   request `qaMechanicalRequest` built:
   - `TestQADeliveryBaseIsTheMergeBaseWithTheDefaultBranch`: after two work
     branch commits, the base equals the fork point, and `true`;
   - `TestQADeliveryBasePrefersTheRemoteTrackingBranch`: with local `main`
     behind `origin/main`, the base is the merge base with `origin/main`;
   - `TestQADeliveryBaseIsUnresolvedWithoutADefaultBranch`: with no
     `origin/HEAD` and a current branch that is not `main` or `master`, it
     returns `""`, `false` and no error;
   - `TestQAMechanicalRequestAuditsATaskCommitFromAnEarlierRun`: task_01's
     commit changes an ungranted governed path, `plan.HeadSHA` is that commit,
     task_02 is committed after it; the result carries a `QA-AUTH-PATHS` finding
     naming task_01's commit;
   - `TestQAMechanicalRequestAuditsEveryCommitOfATask`: task_01 has two commits
     and only the older changes an ungranted governed path; the finding names the
     older commit;
   - `TestQAMechanicalRequestLeavesDefaultBranchCommitsUnaudited`: a commit
     carrying this Spec's trailers that is already on `origin/main` before the
     fork is not in `request.TaskCommits`;
   - `TestQAMechanicalRequestRefusesAGrantWidenedOnTheSpecBranch`: a commit on
     the work branch adds a governed path to the grant, and a later task_02
     commit changes that path; the result carries a `QA-AUTH-PATHS` finding for
     it;
   - `TestQAMechanicalRequestAcceptsAGrantLandedOnTheDefaultBranch`: the same
     widening is committed on `main`, `origin/main` is moved to it and `main` is
     merged into the work branch before task_02; no finding;
   - `TestQAMechanicalRequestFallsBackToTheRunStartHeadWithASkip`: with no
     default branch, `DeliveryTargetRevision` is `plan.HeadSHA`, `TaskCommits`
     holds only the commits after it, and the result's skips include
     `Task commits of earlier Runs`.
7. MUST add a **Delivery Base** entry to `CONTEXT.md` that contains the phrase
   `merge base of the audited head and the repository default branch`, says the
   mechanical stage reads a Spec's Task commits and its grant from it, and keeps
   the glossary free of Spec references.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Task commit from an earlier Run, and an older commit of a Task, reach
      the authorization audit.
- [ ] A grant widened on the Spec branch does not authorize; the same grant
      landed on the default branch does.
- [ ] Without a default branch the stage keeps the Run start head and records
      `Task commits of earlier Runs`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `internal/daemon/qa_delivery_base.go`
- creates: `internal/daemon/qa_every_run_audit_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestQADeliveryBaseIsTheMergeBaseWithTheDefaultBranch|TestQADeliveryBasePrefersTheRemoteTrackingBranch|TestQADeliveryBaseIsUnresolvedWithoutADefaultBranch|TestQAMechanicalRequestAuditsATaskCommitFromAnEarlierRun|TestQAMechanicalRequestAuditsEveryCommitOfATask|TestQAMechanicalRequestLeavesDefaultBranchCommitsUnaudited|TestQAMechanicalRequestRefusesAGrantWidenedOnTheSpecBranch|TestQAMechanicalRequestAcceptsAGrantLandedOnTheDefaultBranch|TestQAMechanicalRequestFallsBackToTheRunStartHeadWithASkip|TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQADeliveryBaseIsTheMergeBaseWithTheDefaultBranch TestQADeliveryBasePrefersTheRemoteTrackingBranch TestQADeliveryBaseIsUnresolvedWithoutADefaultBranch TestQAMechanicalRequestAuditsATaskCommitFromAnEarlierRun TestQAMechanicalRequestAuditsEveryCommitOfATask TestQAMechanicalRequestLeavesDefaultBranchCommitsUnaudited TestQAMechanicalRequestRefusesAGrantWidenedOnTheSpecBranch TestQAMechanicalRequestAcceptsAGrantLandedOnTheDefaultBranch TestQAMechanicalRequestFallsBackToTheRunStartHeadWithASkip TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "merge base of the audited head and the repository default branch"` — expected: exit 0; before this Task none of the nine new tests exists and the glossary carries no Delivery Base entry, so the command fails.

## References

- [_techspec.md](_techspec.md) — The Delivery Base; Every Run's Task commits
