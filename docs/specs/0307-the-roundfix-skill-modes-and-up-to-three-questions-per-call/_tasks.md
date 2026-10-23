---
schema: spec-tasks/v1
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
qa: task_04
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — The Roundfix skill's follow-up, status and unattended modes, and up to three questions per call

| id      | title                                                                                           | type    | complexity | needs                     |
| ------- | ----------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Baseline clause allows up to three independent questions per call                          | backend | medium     | —                         |
| task_02 | The Roundfix skill reads its Skill Mode from the request and reports status                    | docs    | medium     | —                         |
| task_03 | The follow-up and unattended modes gather every decision and record each answer where it is read | docs    | medium     | task_01, task_02          |
| task_04 | Run the final QA gate                                                                           | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04

task_01 changes the Baseline core module, its derived files and this
repository's guides, and proves itself in `internal/baseline`, the package that
consumes the catalog. task_02 changes only the Roundfix skill, its mirror, the
layout test, a new skills test and the glossary, so the two share no file and
run in the same Wave. task_03 extends the modes reference and the skills test
task_02 creates, and its test reads the clause task_01 writes, so it needs
both. The gate follows every Task. No other Spec is required; Spec 0306's
follow-up command is named by the skill text only.
