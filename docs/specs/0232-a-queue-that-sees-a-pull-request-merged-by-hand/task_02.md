---
task: task_02
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
status: completed
type: backend
complexity: medium
---

# Task 02: The queue observes a merge through the recorded Pull Request or the Spec's merge evidence, and says so

## Overview

The Backlog Entry of 2026-10-05
([a Pull Request merged by hand stays parked](references/2026-10-05-a-pull-request-merged-by-hand-stays-parked.md))
records Specs 0231, 0225 and 0229 left parked after their Pull Requests were
merged by hand, two of them through Pull Requests the queue never recorded.
This Task implements the Merge Observer that task_01's Delivery Retry asks:
the recorded Pull Request first, then ADR-0232's merge evidence on the local
default branch. It wires the observer into the queue, prints the merge line,
and describes both in the deliver guide and the Roundfix Skill.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-05 named in the Overview by adding
   `commandDeliveryWorkflow.ObserveMerge` per Invariants 8 to 10: with a
   recorded Pull Request number it reads that Pull Request through
   `GitHubCLI.ViewPullRequest` in the repository checkout and reports a merge
   only when it is merged from the item's recorded branch with a non-empty
   head and merge commit; a read failure is an observer error. Otherwise it
   takes the item head (the local item branch's tip when that branch exists,
   else the newest candidate) and reports a merge when `ProveDelivery`
   proves one, with the delivery commit as the merge commit. It reports a
   closed-unmerged Pull Request only when nothing proved a merge.
2. MUST add `worktree.ProveDelivery` and `DeliveryProof` per Invariant 11, as
   an exported entry to the existing `provenDeliveryEvidence`, without
   changing what reconcile decides.
3. MUST give `commandDeliveryWorkflow` a `GitHubCLI` field that
   `newCommandDeliveryEngine` sets to the repository's `gh` reader and that
   tests replace with a scripted command runner, and pass the workflow as
   `EngineDependencies.Merges`. No test may run the real `gh`, reach GitHub,
   or read or write the real `~/.roundfix`.
4. MUST print, for a retry that recorded a merge, the line of API Contract 2
   on stdout before the existing `Retried <slug>: <blocker> -> merged` line,
   matching Surface Transcript 1; a refused closed Pull Request matches
   Surface Transcript 2. All other retry output keeps its text (API
   Contract 4).
5. MUST add the tests named in Verification, in a new test file, using
   disposable Git repositories with a bare origin and a scripted `gh`:
   the 0231 replay (recorded Pull Request merged at a head that descends from
   the candidate, item branch and worktree removed) records the item
   `merged` with that merge commit and prints the merge line, and one
   `Engine.Run` pass then leaves the item `merged` with no cleanup warning;
   the 0225 replay (no recorded Pull Request; the default branch archives the
   Spec through a squash commit the item branch lacks) records the delivery
   commit; a closed Pull Request without merge evidence is refused with the
   item unchanged; an open Pull Request without merge evidence retries as
   before; and a merged Pull Request from another branch is not a merge.
6. MUST describe in `docs/user-guide/commands/deliver.md`, beside task_01's
   rule, which evidence the retry reads (the recorded Pull Request through
   `gh`, then the Spec archived on the local default branch by a delivery
   commit the item branch lacks) and quote the merge line of API Contract 2
   verbatim with its placeholders, `Merged outside the queue: <evidence>; merge commit <sha>`,
   naming both evidence forms.
7. MUST describe the same rule, both evidence forms, the merge line and the
   closed-unmerged refusal in the retry section of
   `.agents/skills/roundfix/references/deliver.md`, under the existing retry
   text and never inside a `### QA settlement` section; run `make
   skills-sync` so the mirrors equal their canonical files; and re-record the
   Roundfix Skill's version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of `.agents/skills/roundfix/SKILL.md` to
   a free version. Running that record command is an implementation step of
   this Task, not part of its Verification: run it after the last skill edit
   and never write a digest or version by hand.
8. MUST NOT change `deliver status`, `deliver resume`, the owner's pass, the
   post-merge cleanup's proofs, the Run Database schema, any other command
   reference or guide, `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Export the merge evidence proof.
- [ ] Implement and wire the Merge Observer.
- [ ] Print the merge line.
- [ ] Add the replay tests.
- [ ] Describe the rule in the deliver guide and the Roundfix Skill, sync and record the version.

## Acceptance Criteria

- [ ] A retry of the 0231 shape records the item `merged` from the recorded
      Pull Request with the item branch gone, prints the merge line, and the
      owner's pass cleans up without a warning.
- [ ] A retry of the 0225 shape records the item `merged` from the merge
      evidence with no Pull Request recorded.
- [ ] A closed Pull Request without merge evidence keeps the item parked with
      the refusal of Surface Transcript 2.
- [ ] The deliver guide and the Roundfix Skill describe the rule and the
      merge line, each mirror equals its canonical file, and the raised
      version is recorded.

## Context

- instruction: `docs/adr/0237-a-delivery-retry-records-a-merge-made-outside-the-queue.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- interface: `internal/worktree/merged_head.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `internal/cli/deliver_merged_outside_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliverRetryRecordsAPullRequestMergedByHandAfterTheItemBranchIsGone|TestDeliverRetryRecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest|TestDeliverRetryRefusesAPullRequestClosedWithoutMerging|TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore|TestObserveMergeIgnoresAMergedPullRequestFromAnotherBranch)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverRetryRecordsAPullRequestMergedByHandAfterTheItemBranchIsGone TestDeliverRetryRecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest TestDeliverRetryRefusesAPullRequestClosedWithoutMerging TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore TestObserveMergeIgnoresAMergedPullRequestFromAnotherBranch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the five tests exists, so their pass lines are missing and the command fails.
- `for phrase in 'Merged outside the queue: <evidence>; merge commit <sha>' 'was closed without merging'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "$phrase" >&2; exit 1; }; done; grep -qF -- 'Merged outside the queue' internal/cli/deliver.go || { printf 'missing merge line in internal/cli/deliver.go\n' >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the guide nor the skill reference quotes the merge line and the CLI does not print it, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → Goals; Core Features 2-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 8 to 12; API Contract 1; API Contract 2; API Contract 4; Surface Transcript 1; Surface Transcript 2; Vocabulary Contract; Testing Approach; Build Order 2
- ADR-0237; ADR-0232; ADR-0189

## Result

Implemented the task_02 slice for Daemon Verification. Task status and the
declared Verification commands remain Daemon-owned; no commit, push or Pull
Request was made.

The command workflow now reads the recorded Pull Request through an injectable
`GitHubCLI` in the repository checkout and supplies `EngineDependencies.Merges`.
Matching merged PR metadata takes precedence, including when the item branch
is gone. Otherwise the observer chooses the local item branch tip or newest
candidate and calls the exported `worktree.ProveDelivery`, which delegates
unchanged to `provenDeliveryEvidence`. A closed-unmerged PR is reported only
after the fallback finds no proof. Retry output prints the merge evidence line
before the existing retry line.

Acceptance evidence from focused checks:

| Acceptance criterion | Implementation and observed evidence |
| --- | --- |
| 0231 recorded PR, removed item branch, merge line, owner cleanup | `TestDeliverRetryRecordsAPullRequestMergedByHandAfterTheItemBranchIsGone` passed: a CI-fix head descends from the candidate, both item branch and worktree are removed, CLI stdout matches the full merge line and retry line, the PR merge commit and newer head are stored, and one real `Engine.Run` pass leaves the item `merged` with an empty blocker. |
| 0225 archive evidence without a recorded PR | `TestDeliverRetryRecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest` passed: a squash delivery archives the Spec on local `main`, retry records that delivery commit and prints the archive evidence form, with no `gh` call. |
| Closed PR remains parked with the refusal | `TestDeliverRetryRefusesAPullRequestClosedWithoutMerging` passed: exit `2`, empty stdout, `Retry refused` and the specified reason on stderr, no owner start, and the entire persisted item unchanged. |
| Guide, Skill, mirrors and version record | The guide and retry reference describe both evidence forms, the local-only fallback, merge line and closed-unmerged refusal. `make skills-sync` succeeded. The required record command succeeded and raised both version fields in canonical and mirrored `SKILL.md` from `0.1.33` to `0.1.34`, adding the version record. A Python byte comparison confirmed both changed mirror files equal their canonical files and inspected the required output phrases. |

Focused commands and outcomes:

- `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/cli -run '^TestObserveMergeIgnores' -count=1` initially failed compilation, showing the missing observer method and `gh` field before implementation, together with a fixture config typo corrected before the replay checks.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/cli -run 'TestDeliverRetryRecords|TestDeliverRetryRefusesAPullRequest|TestDeliverRetryWithAnOpen|TestObserveMerge' -count=1` passed after fixture corrections. In addition to the five requested tests, it covers failed PR reads despite archive evidence, closed-PR archive fallback, missing branch/candidate evidence, preferring the local tip to an older candidate, and incomplete merged PR metadata.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/cli ./internal/worktree -run 'TestDeliverRetry|TestObserveMerge|TestReleaseMergedRuns|TestMergeEvidence' -count=1` passed in both packages, covering nearby retry, cleanup and reconciliation behavior.
- `rtk make skills-sync` succeeded.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` initially encountered sandbox protection on the canonical `SKILL.md`; the same authorized implementation command succeeded with elevated filesystem access. No version or digest was written by hand.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk make baseline-digests` succeeded with `changed: false`; no derived Baseline artifact remains changed.
- `rtk proxy git -c core.fsmonitor=false diff --check` passed.

Every replay uses a disposable Git repository, a bare local origin, a temporary
home and Run Database, a scripted `gh` runner, and an injected owner launch.
No real GitHub call or real `~/.roundfix` access was made. No change was made to
status/resume behavior, the owner pass, cleanup proofs, the database schema,
other Task files or `_tasks.md`. No follow-up implementation was added.

Not run: the Task's two declared Verification commands and full repository
Verification, which remain for the Daemon's settlement.

## Carry-forward provenance

- Source Run: `run_20261005T214431Z_77dc874bf3a7c619`
- Source commit: `cf4ed7d78be069353318b3136c2757018c13cde0`
