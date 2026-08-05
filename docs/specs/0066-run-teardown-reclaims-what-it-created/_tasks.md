---
schema: spec-tasks/v1
spec: 0066-run-teardown-reclaims-what-it-created
qa: task_06
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Run teardown reclaims what it created

| id      | title                                             | type    | complexity | needs            |
| ------- | ------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Terminate the tree and prove each process gone     | backend | high       | —                |
| task_02 | Classify a target's Run Branches as a set          | backend | high       | —                |
| task_03 | Stop blocking review on superseded Run Branch work | backend | medium     | task_02          |
| task_04 | Offer both debris kinds through reconcile          | backend | medium     | task_01, task_02 |
| task_05 | Synchronise the Roundfix Skill with the new contract | chore | low        | task_03, task_04 |
| task_06 | Run the final QA gate                              | qa      | high       | task_05          |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05 · 4 → task_06

Wave 1 fans out because the two debris kinds share nothing: one is a process
tree in `internal/store`, the other a set of Git refs in `internal/worktree`.

task_01 carries the assertion most worth writing first — that an unprovable
termination is reported distinctly from a proven one. ADR-0044 reclaims
orphaned locks by reading that difference, so collapsing the two would make a
host that cannot answer look like success.

task_05 is the one authorized tooling Task, bounded to
`.agents/skills/roundfix/**` and its mirror under the 2026-08-04 standing grant
for CLI-contract synchronisation.

The QA gate is authored as task_06 and depends on task_05, the graph's only
non-QA leaf.
