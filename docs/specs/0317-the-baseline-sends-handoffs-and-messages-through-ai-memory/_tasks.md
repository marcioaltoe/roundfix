---
schema: spec-tasks/v1
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
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
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — The Baseline sends handoffs and messages through ai-memory

| id      | title                                                                                        | type    | complexity | needs                              |
| ------- | -------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | Every upstream setup follows the ai-memory commit and the core module dispatches its skills  | backend | high       | —                                  |
| task_02 | The guidance sends cross-project work through ai-memory and retires the inbox clauses        | backend | high       | task_01                            |
| task_03 | The follow-up command reads no Secondbrain inbox and its keys are deprecated                 | backend | medium     | task_02                            |
| task_04 | The authoring skills, the glossary and the history guide follow the new channel              | docs    | medium     | task_03                            |
| task_05 | Run the final QA gate                                                                        | qa      | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 swaps the core module's skills and imports the upstream setups, and
task_02 edits the same core module, module record and derived files, so it
needs task_01. task_03 writes `CONTEXT.md`, the owned-skill version record
and the Behavior Surface Record, which task_04 writes again, so they run in
series; task_03 follows task_02 only to keep the shared records in one line.
The gate follows every Task. No other Spec is required.
