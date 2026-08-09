---
schema: spec-tasks/v1
spec: 0085-what-an-agent-reads-before-it-decides
qa: task_08
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
      needs: [task_01]
    - id: task_07
      file: task_07.md
      needs: [task_04]
    - id: task_08
      file: task_08.md
      needs: [task_05, task_06, task_07]
---

# Tasks — What an Agent reads before it decides

| id      | title                                                     | type    | complexity | needs                     |
| ------- | --------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Record the archive paths and the conditional clause today | test    | medium     | —                         |
| task_02 | Give the archive layout one owner                         | backend | medium     | task_01                   |
| task_03 | Move every consumer onto the resolver                     | backend | high       | task_02                   |
| task_04 | Relocate the retired artifacts under one root             | docs    | high       | task_03                   |
| task_05 | Give every retired ADR a status and a forward pointer     | docs    | medium     | task_04                   |
| task_06 | Make the Secondbrain consultation unconditional           | infra   | medium     | task_01                   |
| task_07 | Exclude the single archive root from review               | infra   | low        | task_04                   |
| task_08 | Run the final QA gate                                     | qa      | high       | task_05, task_06, task_07 |

Waves: 1 → task_01 · 2 → task_02, task_06 · 3 → task_03 · 4 → task_04 · 5 → task_05, task_07 · 6 → task_08

Both parallel waves are file-disjoint. Wave 2 splits `internal/spec` (task_02)
from the Baseline catalog and its two rendered guides (task_06); wave 5 splits
`docs/adr/` (task_05) from `.coderabbit.yaml` (task_07).

Tasks 06 and 07 are the two authorized tooling slices. Each names its bounded
files, and neither may touch a path the authorization does not list.
