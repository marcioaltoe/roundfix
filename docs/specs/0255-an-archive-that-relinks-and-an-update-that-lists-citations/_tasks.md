---
schema: spec-tasks/v1
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — An archive that relinks and an update that lists citations

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The glossary, the command references and the skills describe the Archive Relink and History Citations | docs    | medium     | —                |
| task_02 | The owner's archive promotes cited files and relinks Live Documents                           | backend | high       | task_01          |
| task_03 | The Delivery Queue commits and proves an archive that relinked files                          | backend | medium     | task_02          |
| task_04 | Baseline update lists History Citations, promotes on request and names Stale History Citations | backend | high       | task_02          |
| task_05 | Run the final QA gate                                                                         | qa      | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05

task_01 describes ADR-0260 and the TechSpec in the glossary, the command
references and the skills first, so every code Task names its guide through
its dependencies. task_02 creates the shared Live Document classifier and
citation rewrite in `internal/spec` and uses them in the archive. task_03
recomputes that rewrite in the Delivery Queue's proof, and task_04 reuses the
classifier in the update; they declare disjoint files, so they share a wave.
The gate follows both leaves.
