---
schema: spec-tasks/v1
spec: 0056-profiles-configure-merge-semantics
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
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Profiles configure merge semantics

| id      | title                                             | type    | complexity | needs   |
| ------- | ------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Record how the writer behaves today                | test    | medium     | —       |
| task_02 | Derive the effective change set                    | backend | medium     | task_01 |
| task_03 | Merge by category instead of replacing the map     | backend | high       | task_02 |
| task_04 | Summarize the effective change before writing      | backend | medium     | task_03 |
| task_05 | Make a refusal distinguishable by exit code        | backend | medium     | task_04 |
| task_06 | Document the merge contract                        | docs    | low        | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 ·
6 → task_06

The graph is a chain. Tasks 02–05 all edit the same two files — the profile
config document writer and the configure command — so parallel waves would put
two Agent Sessions in the same code.

task_01 changes no behavior on purpose. It records the current writer's output
for real config shapes so every later slice is measured against what works
today. Captured after any behavior change, the corpus would encode the new
behavior as the baseline and prove nothing.
