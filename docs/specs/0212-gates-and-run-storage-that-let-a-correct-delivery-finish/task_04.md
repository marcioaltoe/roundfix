---
task: task_04
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: completed
type: backend
complexity: high
---

# Task 04: Reconcile releases the Runs a merged Spec superseded

## Overview

Reconcile keeps every Run of a Spec merged by squash: commits that changed
only the Spec's directory never match the archived copy, a dirty Run Worktree
is classified before any proof runs, and a Run Branch left without a target
branch has no merged-head fallback. This Task applies ADR-0212: when the
merged head holds the Spec archived, Spec-directory work is represented and
leftovers inside the Spec's declared or recorded scope are superseded, with
forced removal only on a re-proved leftover set; and it describes the rule in
the Roundfix Skill's `reconcile` reference and the `reconcile` command guide.

## Requirements

1. MUST make `inspectRunAtMergedHead` treat paths of non-Task, non-QA commits
   under the Spec's active or archive directory as represented when
   `<archive root>/<slug>/_prd.md` exists at the merged head, and MUST count
   them in the superseded reason of "The merged Spec's Runs".
2. MUST make `inspectTerminalRunMerged` record a dirty worktree's paths, run
   the merged-head proof, and classify `superseded` with the appended reason
   only when the proof yields `safe` or `superseded`, the Spec is archived at
   the merged head, and every dirty path is under the Spec's directories or
   declared (`interface`, `creates`, `deletes`) or recorded by a Task file of
   the archived Spec at that head; otherwise `dirty` with today's reason.
3. MUST make `cleanupTerminalRun` use `git worktree remove --force` only for a
   worktree whose evidence holds superseded dirty paths and whose fresh
   revalidation yields the same state, Run head and dirty path set; every
   other removal MUST stay as today.
4. MUST make `classifyRunBranchSet`, for an absent target branch, apply the
   merged-head proof before preserving a terminal Run, and consider a Run
   whose working root is a linked worktree of the repository.
5. MUST keep preserving any Run with an unrepresented Task commit or a dirty
   path outside the Spec's scope, and MUST NOT read GitHub.
6. MUST add the tests named in Verification to the new file
   `internal/worktree/merged_spec_leftovers_test.go`, over
   `newMergedHeadTestFixture` and `terminalRunFixture`, and MUST leave the
   existing merged-head and inspection tests unedited and passing.
7. MUST describe the rule and its reason words in
   `.agents/skills/roundfix/references/reconcile.md` and
   `docs/user-guide/commands/reconcile.md`, raise the Roundfix Skill's
   version by one patch level from the value on this Task's base in both
   front-matter fields, run `make skills-sync`, and re-record the version.

## Subtasks

- [ ] Represent Spec-directory paths at an archived merged head.
- [ ] Classify declared leftovers as superseded.
- [ ] Force removal only on a re-proved leftover set.
- [ ] Give an absent target branch the merged-head proof.
- [ ] Add the tests.
- [ ] Describe the rule in the skill reference and the command guide, and raise the version.

## Acceptance Criteria

- [ ] A Run whose non-Task commits changed the Spec directory before an
      archive move is `superseded` with the Spec-directory count.
- [ ] A dirty worktree whose changes are a declared path and a Spec-directory
      file is `superseded`, and apply removes it.
- [ ] A dirty undeclared file outside the Spec keeps the Run `dirty`, and an
      unrepresented Task commit keeps it `unintegrated`.
- [ ] Apply refuses as stale when the dirty path set changed after inspection.
- [ ] A Run Branch without its target branch is released when the merged
      head supersedes it.
- [ ] The existing merged-head and inspection tests pass unedited.
- [ ] The reconcile reference and guide carry the reason words, the mirrors
      equal their canonical files, and the raised version is recorded.

## Context

- interface: `internal/worktree/merged_head.go`
- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/merged_spec_leftovers_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/reconcile.md`
- instruction: `internal/worktree/merged_head_test.go`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/spec/recorded_paths.go`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`
- instruction: `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
- instruction: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- instruction: `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/references/2026-10-01-reconcile-releases-the-runs-of-a-squash-merged-spec.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded|TestMergedSpecRunWithDeclaredLeftoversIsSuperseded|TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty|TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated|TestApplyRefusesWhenTheLeftoverSetChanged|TestApplyRemovesASupersededLeftoverWorktree|TestAbsentTargetBranchIsReleasedOnTheMergedHead)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded TestMergedSpecRunWithDeclaredLeftoversIsSuperseded TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated TestApplyRefusesWhenTheLeftoverSetChanged TestApplyRemovesASupersededLeftoverWorktree TestAbsentTargetBranchIsReleasedOnTheMergedHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved|TestMergedHeadRefusesAnUnrepresentedCommit|TestInspectTerminalRunRequiresArchivedEvidence|TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved TestMergedHeadRefusesAnUnrepresentedCommit TestInspectTerminalRunRequiresArchivedEvidence TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing merged-head and inspection tests run unedited beside the new one, which does not exist before this Task.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/reconcile.md | grep -qF -- "uncommitted path(s) superseded by the archived Spec" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/reconcile.md "uncommitted path(s) superseded by the archived Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/reconcile.md | grep -qF -- "Spec-directory path(s) archived" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/reconcile.md "Spec-directory path(s) archived" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- "uncommitted path(s) superseded by the archived Spec" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/reconcile.md "uncommitted path(s) superseded by the archived Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- "Spec-directory path(s) archived" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/reconcile.md "Spec-directory path(s) archived" >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the reconcile reference and guide carry none of these words, so the command fails.

## References

- `_prd.md` → User Story 5; Core Features 7-8; Success Metric 5; Acceptance evidence
- `_techspec.md` → The merged Spec's Runs; API Contract 5; Testing Approach 4; Build Order 4
- ADR-0212; ADR-0053; ADR-0161; ADR-0115; ADR-0052


## Result

Implemented the assigned reconciliation slice under ADR-0212. An archived PRD
at the merged head represents non-Task, non-QA paths in the Spec's active and
archive directories. Other paths still undergo content comparison, and Task
commits still require a completed Task at that head. The superseded reason
counts unique Spec-directory paths; its count clause survives the existing
reason bound even with a Delivery Queue merge-record label.

Dirty inspection reads literal NUL-delimited porcelain paths, including both
sides of renames/copies, and runs the merged-head proof. It supersedes leftovers
only inside the Spec directories or archived Tasks' interface, creates,
deletes, or recorded scope. Instruction-only paths grant no scope. Evidence
holds the sorted unique dirty path set and merged proof head. Apply compares
fresh state, Run head, target/proof heads and the exact set; cleanup revalidates
dirty evidence again so branch-candidate cleanup obeys the same force-removal
rule. Clean removal keeps its existing arguments.

Absent-target Branch Set classification retains the existing QA proof first,
then tries merged-head representation for terminal Runs. Repository-root
normalization includes linked working checkouts. Every proof reads local Git;
no provider, GitHub or network call was added or exercised.

### Acceptance evidence

All named tests below are in the new
`internal/worktree/merged_spec_leftovers_test.go` and were exercised by the
focused package check `go test ./internal/worktree -count=1` (exit 0).

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Spec-directory commits become superseded with their count | `TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded` checks fallback and explicit merged-record reasons, including the bounded count clause. |
| Declared and Spec-directory dirty work is superseded and removable | `TestMergedSpecRunWithDeclaredLeftoversIsSuperseded` checks the appended two-path reason; `TestApplyRemovesASupersededLeftoverWorktree` checks actual worktree and branch absence after apply. |
| Unrelated dirty work and unrepresented Task commits remain preserved | `TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty` checks unrelated and instruction-only paths, refuses apply and keeps both surfaces. `TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated` also checks dirty precedence when declared leftovers coexist with the unrepresented Task. |
| Changed dirty path set refuses apply as stale | `TestApplyRefusesWhenTheLeftoverSetChanged` adds a recorded, still-scoped path after inspection; the heads and superseded state can remain unchanged while apply refuses and retains both surfaces. |
| Absent target gets merged-head release proof | `TestAbsentTargetBranchIsReleasedOnTheMergedHead` checks classification and actual candidate apply from a distinct linked working checkout with local default-branch metadata. |
| Existing merged-head and inspection tests remain unedited | The full worktree package check passed; byte comparisons against HEAD confirm `merged_head_test.go` and `worktree_test.go` are unchanged. |
| Documentation, mirrors and version agree | Both references contain the exact reason words. Canonical/mirror byte checks passed. Roundfix's two frontmatter versions moved from base 0.1.11 to 0.1.12; the generated version record contains 0.1.12. |

Additional regressions cover creates/deletes/recorded scope, a missing archived
PRD, a rename with an outside source, literal quotes/newlines in filenames,
and mixed Spec-directory/outside committed paths.

### Focused checks and regeneration

- Red signal: new Spec-directory and declared-leftover regressions, run using
  a temporary Go overlay of the two production files from HEAD, failed with
  the original `unintegrated` and `dirty` classifications respectively.
- `go test ./internal/worktree -count=1`: exit 0 after the final code/test edits.
- `make skills-sync`: exit 0; sanctioned mirrors regenerated.
- `GOCACHE=/tmp/roundfix-task04-go-cache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  exit 0; generated Roundfix 0.1.12 record. The first attempt with the default
  cache was denied by the sandbox; the task-scoped cache resolved that access.
- `make baseline-digests`: exit 0; no derived artifact changed.
- `GOCACHE=/tmp/roundfix-task04-go-cache rtk make verify-incremental`:
  elevated rerun exited 0 for format, vet, repository tests, skill validation
  and build. The first sandboxed run failed at process-table access in the
  force-stop tests, and its suite guards detected code edits made while it
  ran. The rerun had process-table access and held repository files steady.
- Read-only scope checks confirmed that every changed/untracked path belongs
  to task_04, the existing tests are byte-identical to HEAD, canonical/mirror
  files match, and both reason phrases and recorded 0.1.12 are present.
- `git -c core.fsmonitor=false diff --check`: exit 0.

### Follow-up outside this slice

An exploratory fixture with its recorded `GitRoot` equal to its own Run
`WorkDir` exposed existing cleanup behavior: removing that directory leaves
no working root for the subsequent branch deletion. This unusual metadata
case is left for a separate stable-root cleanup change. The assigned linked
working-checkout behavior is covered with a distinct linked checkout root,
matching normal Run metadata.

Task status, the authored Verification section, the Task Graph and all other
Task files remain untouched by this Agent. Declared Verification and settlement
are pending with the Daemon. No commit, push or Pull Request was performed.


### Verification Feedback — attempt 1

Inspected the Daemon diagnostic at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T100707Z_cbb39973a4b07b9b/verification/batch-004-attempt-1.log`.
The configured `make verify-changed` run recorded the worktree package as
passing. Its failure was in the existing CLI detach-fixture process tests:
`TestImplementDetachChildEndsWhenItsTestBinaryDies` and
`TestRunImplementDetachSurvivesCallerProcessGroupKill` received host
permission denials while probing or killing their recorded process groups.
This is Daemon-provided evidence, not a rerun of the configured gate.

Read the failing tests and their `killDetachFixtureGroup` and
`waitForDetachFixtureExit` helpers. Their failing operations call
`syscall.Kill` for disposable fixture process groups; they do not exercise the
archived-Spec reconciliation changes. Those CLI files belong outside this
Task's slice and remain unchanged.

Focused check:
`GOCACHE=/tmp/roundfix-task04-go-cache go test ./internal/cli -run '^(TestImplementDetachChildEndsWhenItsTestBinaryDies|TestRunImplementDetachSurvivesCallerProcessGroupKill)$' -count=1`,
run with elevated permission, exited 1. The survival test still received
`operation not permitted` for its fixture process-group probe and kill.
Escalation did not resolve this host permission restriction. No Task-scoped
code repair is supported by these diagnostics; no assertion, signal handling,
fixture, or verification configuration was weakened or changed.

Only this Result addendum changed during the feedback turn. The host
process-group permission blocker remains for the Daemon's next configured
Verification attempt. Task status and declared Verification remain
Daemon-owned; neither the configured gate nor the authored Verification
commands were rerun by this Agent. No commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261002T100707Z_cbb39973a4b07b9b`
- Source commit: `1044751b65852d415cb246be2475bf8cdd763b54`
