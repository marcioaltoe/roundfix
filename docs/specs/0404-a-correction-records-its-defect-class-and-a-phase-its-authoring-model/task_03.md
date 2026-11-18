---
task: task_03
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: completed
type: backend
complexity: high
---

# Task 03: The Delivery Queue reads the Correction Signs at its merge and parks an undeclared one where the repository requires it

## Overview

Roundfix learns, without a declaration, that a delivery of the Delivery
Queue was corrected: a commit by hand on the Item Branch, a fixed review
finding and a Corrective Task are Correction Signs. Each time an item
reaches its merge, the Queue Owner reads the signs of the candidate head and
appends them to the Spec's corrections record. With the Project Config key
`delivery.require_defect_class` on at the default branch head, a sign that
no standing declaration names parks the item as `defect-class-undeclared`,
of Park Class `finding`, and `roundfix deliver status` shows each such sign
on a `Corrections:` line before the merge. With the key off or absent the
step never parks, never fails and prints nothing, so a repository that does
not turn the requirement on delivers as it does today.

## Requirements

1. MUST export `runcause.CorrectiveTasks` over the existing rule that tells
   a Corrective Task, with no change to `roundfix runs causes`, Cause Class
   or the Classification Signatures (`_techspec.md` → invariant 10).
2. MUST add the Correction Sign reader as a method of the delivery workflow
   in `internal/cli`, as `_techspec.md` → API Contract 3 states: the default
   branch head from `preflight.ResolveDefaultBranchHead`; the commit signs of
   `git rev-list --first-parent D..H` less a commit with a `Roundfix-Spec`
   trailer or `Roundfix-Delivery: derived-merge`, a merge whose other parents
   are ancestors of `D`, a commit the workflow's `archiveCommitIsExact`,
   `archiveRetirementIsExact` or `pullRequestRecordCommitIsExact` accepts,
   and the fixing commit of a finding sign; the finding signs of the fixed
   dispositions whose reviewed head is in `D..H`; the Corrective Task signs
   of the Spec's Task Graph, from the working tree or from the Archive
   Record's `source_revision` in the local objects, never by a fetch; and
   the subject rule, cut at 72 characters on one line. The Queue Owner, the
   status view and task_04's per-Spec report call the same method (ADR-0422).
3. MUST add `delivery.require_defect_class` to `internal/config` as API
   Contract 4 states: a boolean, `false` when absent, refused by name when
   not a boolean, and removed from User Config before validation with the
   one warning line, as `skills` is in `internal/config/authoring.go`.
4. MUST add to `internal/delivery` the blocker
   `BlockerDefectClassUndeclared`, the `CorrectionSignGate` interface and
   `CorrectionSignCheck` of `_techspec.md` → Interfaces, the field
   `CorrectionSigns` of `EngineDependencies`, and the step of API Contract
   5 in `mergeCandidate`, after the `unauthorized` check and the candidate
   head and before the merge intent: a nil dependency skips it; the key is
   read with `roundconfig.DeliveryConfigAtCommit` at `D`; one signs line is
   appended with the count of dismissed dispositions; with the key on a
   failure is a step error and an undeclared sign parks the item; with the
   key off the step returns no error and prints nothing, a failure to read
   the signs is written as `unavailable` and a failure to write is dropped
   (`_techspec.md` → invariant 4). `ClassifyPark` gives the blocker the
   class `finding` and the next action of API Contract 5, which names the
   five classes and `not-a-correction`. The workflow wires the gate in
   `internal/cli/deliver_workflow.go`. A Delivery Retry needs no change.
5. MUST add the `Corrections:` line of API Contract 6 to
   `roundfix deliver status`, after the `Review:` lines and before the
   `Owner:` line, in queue order, only with the key on; with the key off,
   when `D` does not resolve or when the key or an item's signs cannot be
   read, nothing is printed for it, and the help of `roundfix deliver` does
   not change.
6. MUST write these new tests, each in a new file, each calling
   `t.Parallel()` first:
   - `internal/runcause`: `TestCorrectiveTasksListsTheTasksAboveTheQAGate`.
   - `internal/config`: `TestRequireDefectClassIsAProjectConfigBoolean`
     (absent, `true`, `false`, not a boolean, and in User Config with its
     warning).
   - `internal/delivery`: `TestTheMergeStepParksAnUndeclaredCorrectionSign`
     (key on with one undeclared sign; on with every sign covered,
     including by a `not-a-correction` answer; off with a sign; a nil
     dependency; an error with the key on and with the key off), and
     `TestADefectClassParkRetriesToTheStageItsHeadSupports` (the parked
     item, after its declaration, returns through the Delivery Retry and
     merges).
   - `internal/cli`: `TestCorrectionSignsOfAnItemBranchLikePullRequest892`
     (a fixture Item Branch with three Task commits, a correction, a
     correction of the plan, the QA Report commit, a correction, the
     archive commit and the Pull Request Record Commit, whose signs are
     exactly the three corrections; a merge of the default branch and a
     derived merge that are no sign; a fixed finding whose fixing commit is
     not a second sign; a Corrective Task; and a subject of 300 characters
     with a line that reads as an instruction, printed cut at 72 characters
     on one line); `TestDeliverStatusShowsUndeclaredCorrectionSigns` (key
     on, and key off or absent with the bytes of the Delivery Base);
     `TestTheQueueOwnerRecordsTheSignsOfAnItemBeforeItsMerge` (the signs
     line in the record after a delivery with the key off, still read after
     a new Delivery Queue replaced the one that delivered it and after the
     Item Branch was released).
   Every test creates its own repositories, refs and home directory, and
   reads no ref, tag or history of this repository.
7. MUST change these existing tests' expectations, and no other, declaring
   each in the Result: `TestParkClassesAndNextActions` in
   `internal/delivery/park_class_test.go` gains the row of the new blocker
   with class `finding` and its next action; and
   `parkFindingBlockerCell` in
   `internal/docscontract/command_documentation_test.go` gains the new
   blocker, with its next action among the documented ones. `_tasks.md` →
   What the Tasks stage settled lists what the measurement found. The step
   is built so that no delivery test of `internal/cli` changes: with the key
   absent it never parks, never fails and prints nothing. Before changing
   any test the Task runs `go test ./internal/cli ./internal/delivery
   ./internal/config ./internal/runcause` and `go test -tags docscontract
   ./internal/docscontract` on its base and on its change, and declares
   every other test whose expectation it had to change. It MUST NOT touch
   `internal/cli/cli_test.go` or `internal/docscontract/publicdocs_test.go`,
   and renames or removes no test.
8. MUST describe the slice (ADR-0256): in
   `.agents/skills/roundfix/references/deliver.md`, the new blocker in the
   `finding` row with its next action, and the `Corrections:` status line;
   the same in `docs/user-guide/commands/deliver.md`; and in
   `docs/user-guide/configuration.md`, `delivery.require_defect_class`, what
   it does, that it belongs to Project Config, and that a Roundfix older
   than the release that carries it refuses a Project Config that holds it,
   so a repository turns it on after every machine that runs its Roundfix
   is upgraded.
9. MUST add **Correction Sign** to `CONTEXT.md` through domain-modeling.
10. MUST regenerate, never hand-edit, with the commands of
    `_authorization.md` → Sanctioned regeneration, in this order:
    `make skills-sync`; the `-record-skill-versions` step
    (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
    and the `-record-skill-coverage` step
    (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
11. MUST NOT write `delivery.require_defect_class` into `.roundfixrc.yml`
    or any other Project Config of this repository, and MUST change no
    Governed Path outside `_authorization.md` `paths`.

## Subtasks

- [ ] Export the Corrective Tasks of a Task Graph.
- [ ] Add the Correction Sign reader with its fixture of pull request #892.
- [ ] Add the Project Config key with its User Config warning.
- [ ] Add the gate, the blocker, its Park Class and the merge step, and wire them.
- [ ] Add the `Corrections:` status line.
- [ ] Describe the blocker, the status line and the key, add the glossary term, and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] For the fixture of pull request #892 the Correction Signs are exactly the three corrections; a merge of the default branch and the derived merge are no sign; a fixed finding and a Corrective Task are each one.
- [ ] With the key on, an item with one undeclared sign parks with a Park line that names the sign, the five classes, `not-a-correction`, the declaring command and the Delivery Retry, and `roundfix deliver status` showed the sign before the park; after the declaration a Delivery Retry merges it; an item whose signs are all covered merges without a park.
- [ ] With the key off or absent the same item merges, `roundfix deliver status` prints the bytes of the Delivery Base, and the record holds the item's signs line.
- [ ] A commit subject of 300 characters is printed cut at 72 characters on one line.
- [ ] The key in User Config is ignored with one warning, and a value that is not a boolean is refused by name.

## Context

- instruction: `docs/adr/0422-a-commit-without-a-mark-of-a-run-or-of-the-queue-on-the-item-branch-is-a-correction-sign.md`
- instruction: `docs/adr/0423-an-undeclared-correction-sign-parks-the-item-at-its-merge-as-a-finding.md`
- instruction: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- interface: `internal/runcause/report.go`
- creates: `internal/runcause/corrective_tasks_test.go`
- interface: `internal/config/config.go`
- interface: `internal/config/authoring.go`
- creates: `internal/config/require_defect_class_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/correction_sign_step_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- creates: `internal/cli/correction_signs.go`
- creates: `internal/cli/correction_signs_test.go`
- creates: `internal/cli/deliver_corrections_status_test.go`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `delivery.EngineDependencies`, `ClassifyPark` and the merge stage — consumers: `./internal/delivery`, `./internal/cli`
- Project Config loading and its keys — consumers: `./internal/config`, `./internal/cli`
- `runcause.CorrectiveTasks` over the existing Corrective Task rule — consumers: `./internal/runcause`
- The `finding` row of the delivery guide and reference, and the Behavior Surface Record — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestCorrectiveTasksListsTheTasksAboveTheQAGate|TestRequireDefectClassIsAProjectConfigBoolean|TestTheMergeStepParksAnUndeclaredCorrectionSign|TestADefectClassParkRetriesToTheStageItsHeadSupports|TestCorrectionSignsOfAnItemBranchLikePullRequest892|TestDeliverStatusShowsUndeclaredCorrectionSigns|TestTheQueueOwnerRecordsTheSignsOfAnItemBeforeItsMerge|TestParkClassesAndNextActions)$' ./internal/runcause ./internal/config ./internal/delivery ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCorrectiveTasksListsTheTasksAboveTheQAGate TestRequireDefectClassIsAProjectConfigBoolean TestTheMergeStepParksAnUndeclaredCorrectionSign TestADefectClassParkRetriesToTheStageItsHeadSupports TestCorrectionSignsOfAnItemBranchLikePullRequest892 TestDeliverStatusShowsUndeclaredCorrectionSigns TestTheQueueOwnerRecordsTheSignsOfAnItemBeforeItsMerge TestParkClassesAndNextActions/defect-class-undeclared; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/runcause ./internal/config ./internal/delivery ./internal/cli` — expected: exit 0. Before this Task the blocker, the key, the reader and the named tests do not exist, so the command fails; afterwards every named test passes and the four packages that read the queue, the configuration and the Corrective Tasks pass.
- `grep -qF -- 'defect-class-undeclared' .agents/skills/roundfix/references/deliver.md && grep -qF -- 'defect-class-undeclared' docs/user-guide/commands/deliver.md && grep -qF -- 'require_defect_class' docs/user-guide/configuration.md && grep -F -- '**Correction Sign**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Correction Sign\*\*(:| —( |$))' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task neither guide names the blocker or the key and the glossary lacks the term; afterwards they do, and the documentation contracts, the skill tests, the mirror and the records pass.

## References

- `_prd.md` → Goals (fifth bullet); User Stories 6 and 7; Core Features 7, 8, 9, 10 and 22; Success Metrics 5, 6, 7 and 15
- `_techspec.md` → Correction Sign reader; Merge step and setting; Interfaces (invariants 4 and 10); API Contracts 3, 4, 5, 6 and 11; Testing Approach; Build Order 3
- ADR-0422, ADR-0423, ADR-0417, ADR-0420, ADR-0288, ADR-0408, ADR-0392, ADR-0395, ADR-0398, ADR-0214


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Implemented only this Task's Correction Sign reader, Project Config key, merge step, status lines, documentation and glossary term. The approved design supplied the behavior; no design-authoring skill or person-only skill was activated. The Daemon retains status and declared Verification ownership.

Acceptance evidence:

- Sign identity and exclusions: `commandDeliveryWorkflow.CorrectionSigns` is shared by the owner and status and available for task_04. It reads only local Git objects and refs, uses first-parent commit signs, real trailer metadata, default-parent ancestry and the queue's exact archive/retirement/PR-record proofs. Fixed dispositions suppress their fixing commit; dismissed dispositions are counted separately. `runcause.CorrectiveTasks` exports the existing numeric-above-QA rule without changing causes classification. The real-Git #892 fixture returns exactly the three corrections; separate cases cover default and derived merges, one fixed finding, and active/archived Corrective Tasks from `source_revision`.
- Required declarations and retry: the gate runs after authorization and candidate-head resolution, before merge intent. The setting comes from `DeliveryConfigAtCommit` at the resolved default head. Undeclared signs park as `defect-class-undeclared`, class `finding`; the next action names the signs, all five classes, `not-a-correction`, declaration and retry. A fixture runs status, parks, declares through the actual CLI, retries through `checking` and merges with one counted retry. Moved-head retry returns through `reviewing`; covered signs and `not-a-correction` are accepted. Nil dependency and on/off errors are exercised.
- Default-off compatibility and persistence: explicit false and absent keys preserve the exact Delivery Base status bytes. A real engine merge with the key absent writes a snapshot silently; replacing the queue and deleting the Item Branch preserve and leave readable the same record. Read failures are recorded as unavailable while off; record failures are ignored while off and returned while on. No resolvable default means no requirement and no snapshot. Status reads without writing, omits unreadable records, and appears after review notices before owner status.
- Subject boundary: a 300-character instruction-like subject becomes one quoted line of 72 characters ending in an ellipsis. Controls, line/paragraph separators and double quotes are sanitized. The second message line is not rendered.
- Config validation: absent/true/false, non-booleans and immutable commit reads are covered. User Config values, including invalid ones, are removed before validation with exactly the specified warning. The configuration guide documents Project Config ownership and upgrading every machine before adding the key. No Project Config in this repository was edited.

Focused checks:

- Clean base at `92c3e409`, isolated detached checkout `/tmp/roundfix-0404-task03-base`: `go test ./internal/cli ./internal/delivery ./internal/config ./internal/runcause` passed (CLI 272.372s; delivery 3.898s; config 1.653s; runcause 0.440s). Before isolation, the first Run-checkout measurement was rejected by suiteguard because implementation edits occurred while it was running; that result is not counted as a product failure or a clean baseline.
- Base `go test -tags docscontract ./internal/docscontract` passed (10.743s), before documentation edits.
- `go test ./internal/cli ./internal/delivery ./internal/config ./internal/runcause -run 'Test(CorrectionSignsOf|DeliverStatusShowsUndeclared|TheQueueOwnerRecords|CorrectiveTasksLists|RequireDefectClassIs|TheMergeStepParks|ADefectClassPark|ParkClassesAndNextActions)'` passed (CLI 2.950s; delivery 1.013s; config 1.207s; runcause 0.318s). Earlier new-fixture failures were repaired at their source: a missing fake PR caused the retry to wait for checks, and a Corrective Task needed to be added to its QA gate's dependencies.
- Final post-change `go test ./internal/cli ./internal/delivery ./internal/config ./internal/runcause` passed (CLI 257.317s; the other three packages reused unchanged successful test-cache entries). An earlier post-change measurement also passed (CLI 255.726s; delivery 3.275s; config 1.659s; runcause 0.274s), before the two additional reader-boundary probes below.
- Two additional regression probes first failed for a fixed finding recorded in an earlier worktree and for an unreadable default ref. The reader now trusts the repository-scoped ledger plus D..H membership without imposing the current checkout path; default-resolution failure now skips both requirement and recording. After those repairs, the focused four-package command above passed again (CLI 3.208s; other packages cached), and the final package measurement passed.
- Final post-change `go test -tags docscontract ./internal/docscontract` passed (13.083s). Its first post-change run exposed the glossary's required final `_Avoid_` list; that list was added to Correction Sign without changing any existing expectation.
- `go test ./skills -run 'Test(EveryOwnedSkillVersionIsRecorded|RoundfixReferenceIndex|RoundfixSkill|Deliver)'` passed (0.285s). Direct byte comparisons of the canonical/mirrored Roundfix entry and delivery reference also passed.
- Changed-path postflight found 25 paths, including the seven requested new tests, all within this Task's declared interfaces/creates plus its own file. The Task status remains `in_progress`, as supplied by the Daemon.
- `git diff --check` passed.

Existing expectations changed, and no others: `TestParkClassesAndNextActions` in `internal/delivery/park_class_test.go` gains the new blocker/class/next-action row; `parkFindingBlockerCell` and its documented-action list in `internal/docscontract/command_documentation_test.go` gain that blocker and action. No existing test was renamed or removed. `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `_tasks.md` and other Task files were not edited.

Sanctioned regeneration, executed in the authored order:

1. `make skills-sync` passed.
2. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` passed. Both canonical and mirrored Roundfix skill version fields changed from `0.1.264` to `0.1.265` through the recorder.
3. `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` passed, regenerating `docs/references/behavior-surfaces.json`. The three steps were confirmed again in this order after the final source/glossary adjustments; the version stayed `0.1.265`.

The Task's two declared Verification commands were not run; the Daemon owns them and settlement. No commit, push, PR, Spec announcement, other-repository read or memory write was performed.

## Carry-forward provenance

- Source Run: `run_20261118T051745Z_d63008af9d14015a`
- Source commit: `e834903982eebcf81f6fc17a6f1449b3edd260dd`
