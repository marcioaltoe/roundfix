---
task: task_03
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: completed
type: backend
complexity: high
---

# Task 03: Settle follows the Spec plan its Run's branch carries

## Overview

Settle loads the Spec from the surface it settles, so a plan the Run's branch
amended after the Run started is never read: Quaestor's QA Task stayed
refused against the old PRD until Quaestor moved the kept Run Worktree by
hand. This Task makes settle compare the surface with the Run's branch and,
when the branch changed the Spec folder, fast-forward a Run Worktree the
branch is ahead of or refuse any other drift, naming the files; and writes
the rule into the settle usage, the Roundfix `settle` reference, the
`settle` guide and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md`](references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md)
   of 2026-10-27, its settle half, by implementing API Contracts 6 and 10 of
   `_techspec.md` and ADR-0343.
2. MUST add the branch comparison to `internal/cli/settle.go` as API
   Contract 6 states, as the last preflight step, after the commit-authority
   check and every other refusal, for a surface that belongs to a kept Run:
   act only on files the branch changed under the Spec folder since its
   merge base with the surface's HEAD; fast-forward a Run Worktree whose
   HEAD is that merge base with `git merge --ff-only`, print
   `Run Worktree advanced: <old short> -> <new short> (<branch>)` on stderr
   and reload the Spec and Task from the surface; refuse every other such
   drift, or a fast-forward Git refuses, with exit `2` and nothing written.
   A branch ahead only outside the Spec folder keeps today's behaviour.
3. MUST add `internal/cli/settle_plan_test.go`, through CLI dispatch on the
   implement workspace fixture with a kept Run Worktree, with:
   - `TestSettleAdvancesTheRunWorktreeToTheBranchPlan`: a failed QA Task
     whose report in the Run Worktree has two declared rows; the branch then
     adds a second `criterion:` line to the PRD's one declaration; settle
     exits `0`, stderr shows `Run Worktree advanced: ` before
     `Run Worktree removed: `, and the user branch holds the settled QA Task;
   - `TestSettleRefusesWhenTheBranchChangeOverlapsTheSurface`: the branch
     edits the Task file the surface also changed uncommitted; exit `2`
     naming that file, surface HEAD and files unchanged;
   - `TestSettleRefusesASpecChangedOnADivergedBranch`: the Run Branch has a
     commit of its own and the branch amends the PRD; exit `2` naming
     `_prd.md`, nothing written;
   - `TestSettleRefusesASpecChangedAfterATaskWorktreeForked`: the same drift
     with a kept Task Worktree surface; exit `2`;
   - `TestSettleKeepsTheAdvanceWhenVerificationFails`: an advance followed by
     a failing Verification exits `1`, the Run Worktree stays at the branch's
     commit and the Task status is unchanged.
4. MUST keep `TestSettleNamesNoRemovalWhenIntegrationIsPending` in
   `internal/cli/settle_test.go` unchanged and passing: its branch gains a
   commit outside the Spec folder, so settle still reports
   `integration pending`.
5. MUST add API Contract 10's sentence, verbatim, to `settleUsage`, the
   `## Settle Command` section of `.agents/skills/roundfix/references/settle.md`
   and `docs/user-guide/commands/settle.md`, each with the
   `Run Worktree advanced:` line and the note that a later Verification
   failure keeps the advance; then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the `roundfix` version pair.
6. MUST revise **Run Worktree** in `CONTEXT.md` through `domain-modeling`
   with "settle fast-forwards it to the Run's branch when that branch changed
   the Spec folder", keeping its `_Avoid_` line.
7. MUST, as the last step after the last guide edit, re-record the Behavior
   Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit no mirror, version field or record.
8. MUST NOT rewrite, rebase or merge any Run commit other than by a
   fast-forward, advance a Task Worktree, change the authorization settle
   reads, Task Carry-Forward, the Daemon, any text inside a
   `### QA settlement` section, the Makefile or the CI workflows. Each new
   top-level test calls `t.Parallel()` first.

## Subtasks

- [ ] Compare the surface with the Run's branch; fast-forward or refuse on Spec-folder drift.
- [ ] Add the five CLI tests and keep the pending-integration test passing.
- [ ] Write the rule into the usage, reference, guide and glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] A Run Worktree the branch is ahead of is fast-forwarded before settle reads the plan, and the QA Task settles against the amended PRD.
- [ ] A diverged branch, a Task Worktree or an overlapping change is refused with nothing written; a branch ahead outside the Spec folder behaves as before.
- [ ] The usage, reference, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- Settle preflight — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0343-settle-follows-the-branch-plan-and-a-declared-row-matches-the-criterion-it-names.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/carryforward.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/settle_test.go`
- creates: `internal/cli/settle_plan_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/settle.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestSettleAdvancesTheRunWorktreeToTheBranchPlan|TestSettleRefusesWhenTheBranchChangeOverlapsTheSurface|TestSettleRefusesASpecChangedOnADivergedBranch|TestSettleRefusesASpecChangedAfterATaskWorktreeForked|TestSettleKeepsTheAdvanceWhenVerificationFails|TestSettleNamesNoRemovalWhenIntegrationIsPending)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettleAdvancesTheRunWorktreeToTheBranchPlan TestSettleRefusesWhenTheBranchChangeOverlapsTheSurface TestSettleRefusesASpecChangedOnADivergedBranch TestSettleRefusesASpecChangedAfterATaskWorktreeForked TestSettleKeepsTheAdvanceWhenVerificationFails TestSettleNamesNoRemovalWhenIntegrationIsPending; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/cli/settle.go | grep -qF -- 'settle fast-forwards a Run Worktree the branch is ahead of and refuses any other drift.' || { echo 'settle usage lacks the branch rule' >&2; exit 1; }; go test -count=1 -timeout 15m -run 'Settle|CarryForward' ./internal/cli` — expected: exit 0. Before this Task the five new tests do not exist, so the first pass check fails. Afterwards the advance and the three refusals hold, the pending-integration case is unchanged, and every settle and carry-forward test of the package passes.
- `for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md docs/user-guide/commands/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'settle fast-forwards a Run Worktree the branch is ahead of and refuses any other drift.' || { printf '%s lacks the branch rule\n' "$file" >&2; exit 1; }; grep -q 'Run Worktree advanced: ' "$file" || { printf '%s lacks the advance line\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'settle fast-forwards it to the Run' || { echo 'CONTEXT.md lacks the revised **Run Worktree**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Worktree**:' || { echo 'CONTEXT.md lacks **Run Worktree**' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference or guide states the rule, so the first phrase check fails. Afterwards the skill and its mirror agree, the versions are recorded and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 3 and 6; Core Features 3 and 6; Success Metrics 5 and 7; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior 6, 7 and 9; API Contracts 6 and 10; Declared breaks 4 and 6; Invariant 3; Testing Approach; Build Order 3
- ADR-0343


## Result

Implemented the Task 03 branch-plan slice. Initial inspection found no branch
comparison between commit-authority preflight and Verification; the five
required regression tests were absent. The existing status-only change in
this Task file was preserved.

- Acceptance criterion 1: the final preflight step compares surface HEAD and
  the Run's target branch at their merge base, scoped to branch changes in
  this Spec folder. An ancestor Run Worktree advances through `git merge
  --ff-only`, reports both short revisions on stderr, then reloads the Spec
  and Task. The new QA regression uses one declaration and two declared rows;
  adding the second criterion on the user branch allows settlement against
  that amended plan, with the advance printed before removal.
- Acceptance criterion 2: CLI regressions assert exit 2 and unchanged HEAD,
  Task, changed-file bytes and Git status for overlapping local work,
  diverged Run commits and a Task Worktree. A Verification failure after an
  advance exits 1, retains the branch commit and leaves Task status unchanged.
  Authorization remains the Run-start record. The existing pending-integration
  test was left byte-identical and passes for branch changes outside the Spec.
- Acceptance criterion 3: usage, the canonical settle reference and the user
  guide contain API Contract 10 verbatim, the advance transcript and retained
  advance on Verification failure. The Run Worktree definition includes the
  requested rule and preserves its `_Avoid_` line. Generated skill mirrors
  and both roundfix version fields are now 0.1.154; the Behavior Surface record
  was regenerated after the last guide edit.

Focused checks and required record commands:

- `go test ./internal/cli -run 'TestSettle(AdvancesTheRun|RefusesWhenTheBranch|RefusesASpec|KeepsTheAdvance)' -count=1 -timeout 5m`
  — exit 0, all five new CLI regressions exercised.
- `go test ./internal/cli -run '^TestSettleNamesNoRemovalWhenIntegrationIsPending$' -count=1 -timeout 3m`
  — exit 0.
- `make skills-sync` — exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; recorder raised and recorded the roundfix version pair.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  — exit 0; last regeneration after guide edits.
- `git diff --check` — exit 0.

The Task's declared Verification and repository Verification were not run;
those remain Daemon-owned. No Task status, other Task file, Task Graph,
Carry-Forward implementation, Daemon, authorization reader or QA settlement
section was changed. No commit, push or pull request was made.
