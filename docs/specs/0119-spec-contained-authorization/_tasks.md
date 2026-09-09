---
schema: spec-tasks/v1
spec: 0119-spec-contained-authorization
qa: task_09
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_02, task_05]
    - id: task_07
      file: task_07.md
      needs: [task_03, task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
    - id: task_09
      file: task_09.md
      needs: [task_08]
---

# Tasks — Spec-contained authority and trusted Verification

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Characterize today's grant resolution and governed set | test | high | — |
| task_02 | Read a typed grant by role and keep legacy records resolving | backend | high | task_01 |
| task_03 | Refuse a cited record that is not an operative grant | backend | medium | task_01, task_02 |
| task_04 | Extend the governed set to the owned shipped templates | backend | low | task_01 |
| task_05 | Audit the consuming commit against the grant's ancestor | backend | high | task_02, task_03, task_04 |
| task_06 | Execute authored commands only on committed provenance | backend | high | task_02, task_05 |
| task_07 | Make the Spec-contained record canonical in the Baseline | docs | medium | task_03, task_06 |
| task_08 | Teach the authoring skills the record's placement | docs | medium | task_07 |
| task_09 | Run the final QA gate | qa | high | task_08 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05 · 5 → task_06 · 6 → task_07 · 7 → task_08 · 8 → task_09.
