---
schema: spec-tasks/v1
spec: 0284-the-linux-lane-reads-the-repository-history
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

# Tasks — The Linux Lane reads the repository history

| id      | title                                                                      | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Linux Lane mounts the repository object store read-only                 | backend | medium     | —       |
| task_02 | Run the final QA gate                                                      | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 is one vertical slice: the mount, the skips, the script, their tests,
the `deliver` guide, the Roundfix `deliver` reference, the glossary, the skill
version and the Behavior Surface Record. The gate follows it. No other Spec is
required.
