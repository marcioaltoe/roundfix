---
schema: spec-tasks/v1
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A parked item recovers with a live Queue Owner and its completed Tasks

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The item records its Failed Run and a Delivery Retry carries the Tasks of a QA-Only Failed Run | backend | high       | —       |
| task_02 | A Delivery Retry ends with a proven live Queue Owner or parks the item again and fails         | backend | high       | task_01 |
| task_03 | Status reports a queue nobody advances and the Run of each running item, under ADR-0450        | backend | high       | task_02 |
| task_04 | The deliver guide and the Roundfix skill describe the recovery lines                           | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order one step per Task: steps 1 to 4 are task_01 to task_04 and step 5, the gate, is task_05. The Tasks run in series because each pair shares a declared file: `CONTEXT.md` is changed by task_01, task_02 and task_03; `internal/cli/deliver.go` by task_02 and task_03; the carry-forward fixture `internal/cli/deliver_retry_runs_test.go` may be changed by task_01 and task_02; and task_04 describes the lines task_02 and task_03 print. No Task is joined to another: each changes a different contract of the Delivery Queue (the recorded Run and the carry-forward source, the retry's end, the status lines and their decision record, the documentation) and each pays its verification cycle with tests of its own.

task_01 to task_03 are `backend`: their outcome is the behavior of the Delivery Engine and of the `deliver` commands. task_04 is `docs`: its outcome is the guide page and the skill reference, with their regenerated records.

## What the Tasks stage settled beyond the TechSpec

- **The ADR-0450 citation of the glossary entries.** API Contract 8 has each revised entry cite ADR-0450, and Build Order 3 writes ADR-0450. task_01 and task_02 therefore write the glossary text of their entries without a link to a file that does not exist yet, and task_03 adds the ADR-0450 link to the four entries (QA-Only Failed Run, Task Carry-Forward, Delivery Retry and Queue Owner) in the change that creates the record.
- **The reconcile Behavior Surface.** `internal/cli/reconcile.go` is a source of the Behavior Surface `command reconcile` in `docs/references/skill-coverage.json`, covered by the Roundfix skill's `reconcile` reference, which the authorization record does not grant. task_01 adds one field to `reconcileRunSelection` that only a Delivery Retry sets, so the Reconcile Command's help, flags, output and exit codes do not change (Core Feature 3). The PRD's `## Skills` records the surface as unchanged with that reason, and task_01 writes the dated Coverage Review in the map's `command reconcile` entry. That entry is the only change to the map; the TechSpec's note that the map does not change still holds for every other surface.
- **The guide in task_01.** task_01 changes `internal/cli` sources of the `deliver` commands, and the Tasks stage requires such a Task to name the command's guide itself or through a Task it needs (write-tasks: "Each Task that changes a command declares the one command file it changes"; `roundfix spec check` reports `SC-CLI-UNDOCUMENTED` otherwise). task_01 has no Task before it, so it documents its own slice in `docs/user-guide/commands/deliver.md`, the QA-Only Failed Run in the retry section, and re-records the guide page's fingerprint with the `-record-skill-coverage` step, which `SC-SURFACE-RECORD-UNTASKED` requires of the Task that edits a guide page. task_02 and task_03 name the guide through task_01. The three granted skill files, `make skills-sync` and the skill version record stay in task_04 as decided, and task_04 runs `-record-skill-coverage` again for its own guide edit. This departs from `_techspec.md` → Build Order, which keeps every guide edit in step 4; the behavior, the contracts and the order are unchanged.
- **Test and subtest names.** The TechSpec names the tests; each Task names the subtests of its controls and of the further retry cases, so the Verification and the gate check them by name. The existing test `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` is renamed `TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner` in task_03.
- **The ADR number.** 0450. The default branch ends at decision record 0449 (checked on 2026-11-23). When 0450 is taken before task_03 runs, the operator renames the file, its two notice links and its export line and amends task_03; invariant 7 keeps the decision record test independent of the number.
- **The Governed Paths.** Measured again over every file the five Tasks declare: the three paths of `_authorization.md` are the only governed ones, and only task_04 declares them, so the record is unchanged.

## Shared files

Spec 0417 is reserved and not in the repository. If it changes an owned skill, a user-guide page or the glossary first, the two Specs share the Roundfix skill's version fields, the owned-skill version record, the Behavior Surface Record and `CONTEXT.md`; task_04 regenerates the shared records on the tree it runs on and writes "the next version", never a number. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Delivery

The authorization record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. No Task performs the release, edits `CHANGELOG.md`, a version file or the npm manifest, or sends a message; the operator cuts the release and tells quaestor and the Praetor.

## Time in tests

No new test waits on elapsed time. The retry tests end on the fake owner start's return and the queue read that follows it; the status tests read Runs and Run Events appended to a temporary Run Database; the built-binary test waits on the binary's exit.
