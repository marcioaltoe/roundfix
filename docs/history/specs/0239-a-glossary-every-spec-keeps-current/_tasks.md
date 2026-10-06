---
schema: spec-tasks/v1
spec: 0239-a-glossary-every-spec-keeps-current
qa: task_05
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
---

# Tasks — A glossary every Spec keeps current

| id      | title                                                                                         | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Spec Consistency Check and the Archive Command hold a Glossary Declaration to the glossary | backend | high       | —                         |
| task_02 | The Baseline's domain clauses say the glossary and the ADRs are the base                      | backend | medium     | —                         |
| task_03 | The authoring, QA and Roundfix skills write and check the Glossary Declaration                | docs    | medium     | task_01                   |
| task_04 | The glossary gains the terms dropped since 2026-10-02 and this Spec's own                     | docs    | medium     | —                         |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_02, task_03, task_04 |

Waves: 1 → task_01, task_02, task_04 · 2 → task_03 · 3 → task_05

task_03 describes the findings and the refusal task_01 ships, so it follows
task_01. task_01, task_02 and task_04 share no declared file. task_04 is the
docs Task that writes every term this Spec's Glossary Declaration adds or
changes. The gate follows every leaf.
