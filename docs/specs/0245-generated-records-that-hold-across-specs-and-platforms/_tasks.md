---
schema: spec-tasks/v1
spec: 0245-generated-records-that-hold-across-specs-and-platforms
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
      needs: [task_05]
    - id: task_05
      file: task_05.md
      needs: [task_03]
---

# Tasks — Generated records that hold across Specs and platforms

| id      | title                                                                                  | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A Baseline module's version is chosen by its record step and checked against its content | backend | high       | —                |
| task_02 | The Coverage Record lists every release platform and is the same bytes on any host     | test    | high       | —                |
| task_03 | The derived declarations run both record steps, and the rule and glossary name them    | chore   | medium     | task_01, task_02 |
| task_04 | Run the final QA gate                                                                  | qa      | high       | task_05          |
| task_05 | The derived declarations cover every file a module edit regenerates                    | chore   | medium     | task_03          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_05 · 4 → task_04

task_01 works in `internal/baseline` and the Baseline assets, task_02 in
`internal/spec` and the Coverage Record; they share no declared file, so they
run in parallel. task_03 declares the commands both add and re-records the
Coverage Record after them, so it follows both. The gate follows task_03.
