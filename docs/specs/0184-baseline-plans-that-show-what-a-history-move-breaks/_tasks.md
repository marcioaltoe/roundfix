---
schema: spec-tasks/v1
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
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
      needs: [task_01, task_02, task_03, task_05]
    - id: task_05
      file: task_05.md
      needs: [task_01]
---

# Task Graph

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Find the Relocation Citations a set of History Relocations breaks | backend | high | — |
| task_02 | Every Baseline Plan with History Relocations reports its Relocation Citations | backend | medium | task_01 |
| task_03 | The skill, the guide and the glossary describe Relocation Citations | docs | low | task_02 |
| task_04 | Run the final QA gate | qa | medium | task_01, task_02, task_03, task_05 |
| task_05 | The Windows citation open reads files synchronously | backend | low | task_01 |

Waves: 1 → task_01 · 2 → task_02, task_05 · 3 → task_03 · 4 → task_04
