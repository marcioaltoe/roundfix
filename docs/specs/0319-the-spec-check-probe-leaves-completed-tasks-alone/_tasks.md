---
schema: spec-tasks/v1
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: [task_02]
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
---

# Tasks — The Spec Check probe leaves completed Tasks alone

| id      | title                                                                    | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Spec check probes only the Tasks a Run would still execute               | backend | medium     | task_02          |
| task_02 | Document the probe's selection rule                                      | docs    | medium     | —                |
| task_03 | Run the final QA gate                                                    | qa      | high       | task_01, task_02 |

Waves: 1 → task_02 · 2 → task_01 · 3 → task_03

task_01 changes the `spec check` CLI surface, and the guide that documents it
is written by task_02, so task_01 needs task_02. They share no file: task_01
changes code and tests only, task_02 documents, skills, the ADR, the glossary
and the Behavior Surface record only, and task_02's Verification reads
documents only. The gate follows both. No other Spec is required.
