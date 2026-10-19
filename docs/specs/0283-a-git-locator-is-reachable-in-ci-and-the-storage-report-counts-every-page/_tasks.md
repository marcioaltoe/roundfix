---
schema: spec-tasks/v1
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
qa: task_04
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
---

# Tasks — A Git Locator is reachable in CI, and the storage report counts every page

| id      | title                                                                                   | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Default Branch Head, and a Git Locator is reachable in a CI checkout                | backend | high       | —       |
| task_02 | The Delivery Base, the archive warning and the review base read the Default Branch Head | backend | medium     | task_01 |
| task_03 | The storage report measures one Storage Snapshot and groups every page                  | backend | high       | task_02 |
| task_04 | Run the final QA gate                                                                   | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 reads the resolver task_01 adds and is its Integration Checkpoint for
the `cli`, `daemon` and `worktree` suites. task_01 and task_02 both change the
`archive` guide and reference. All three implementation Tasks raise the
Roundfix skill version and re-record `skills/testdata/owned-skill-versions.json`,
the Behavior Surface Record and `CONTEXT.md`, so they run in series. The gate
follows the last of them. No other Spec is required.
