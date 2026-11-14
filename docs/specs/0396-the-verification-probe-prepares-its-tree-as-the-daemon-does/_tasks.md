---
schema: spec-tasks/v1
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The Verification probe prepares its tree as the Daemon does

| id      | title                                                                                 | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Verification probe prepares its checkout as the Daemon prepares a Task Worktree    | backend | medium     | —       |
| task_02 | Run the final QA gate                                                                 | qa      | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

Sized by the grouping criterion, the Spec holds one context, how the
authoring Verification probe prepares its tree, in one implementation Task and
the QA gate. The guide, skill reference, glossary and ADR notices join
task_01 instead of becoming their own Task: alone they would not earn a
verification cycle, and the code change already requires the Behavior Surface
Record and the skill that covers the `spec check` surface. task_01 runs the
suites that consume its changed contracts (`./internal/cli`,
`./internal/worktree`, `./skills` and `./internal/docscontract`) in its own
Verification, so no separate Integration Checkpoint is needed. No other Spec
is required; Specs 0392 and 0393 deliver first and touch only the shared
files named in `_prd.md` → Prerequisites.
