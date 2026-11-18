---
schema: spec-tasks/v1
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
qa: task_03
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
---

# Tasks — The Delivery Queue brings a Pull Request that fell behind its base branch up to date

| id      | title                                                                                  | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The checking stage makes a Base Update when GitHub reports the Pull Request BEHIND     | backend | high       | —       |
| task_02 | The deliver reference, the user guide and the glossary describe the Base Update        | docs    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                  | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The graph follows `_techspec.md` → Build Order. Task 01 changes the Delivery Engine and proves the change in the engine and through the command engine with a strict local forge. Task 02 needs it: its contract test reads the next action of `pull-request-behind` through `delivery.ClassifyPark`, which Task 01 adds. The two Tasks declare no common file. Each runs the suites that consume its changed contract in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, the checking stage of the Delivery Queue and its documentation, in two implementation Tasks and the QA gate. Each Task earns its verification cycle: Task 01 carries the engine change with nine engine tests and an end-to-end command-engine test, and Task 02 carries three documents, a contract test, two glossary entries, a skill version record and the Behavior Surface Record. A third implementation Task would split the engine change from the test that proves it through the command engine, and that test alone would not earn a cycle.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on the planning branch, based on `1a8e08b4`; the measurements are in [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md).

- **The prototype passed the existing suites.** With the engine and Park Class change of Task 01 applied, the whole `./internal/delivery` package passed unchanged, so no existing Delivery Engine test is declared as changed except the class-table case.
- **The command-engine test is feasible as written.** A prototype of `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate` merged the item after a Base Update with the change and parked `delivery-error: merge pull request: merge pull request: X Pull request #7 is not mergeable: the base branch policy prohibits the merge.` without it.
- **The park-action contract.** `checkParkActionDocumentation` requires each listed blocker's action sentence exactly once per document, built from `ClassifyPark`; Task 02 adds the new blocker to that list so a reworded or missing sentence fails.
- **Shared files with Spec 0407.** Spec 0407 delivers first and may change `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `skills/testdata/owned-skill-versions.json` and the Roundfix skill's version fields. Task 02 edits `CONTEXT.md` only at the two entries it names and regenerates the other three on the tree it runs on.
- **Exact text.** Tasks 01 and 02 fix the words of the log line, the blockers, the next action, the documents and the glossary from `_techspec.md` → API Contracts 1 to 4, so the QA gate searches for fixed strings.

## Time in tests

The Delivery Engine tests use `fakeDeliveryClock`, whose sleeps advance a fake clock; the timeout subtest asserts the `checks-timeout` blocker, never elapsed time. The command-engine test waits only on the events it asserts: the child process's exit and the forge state.
