---
schema: spec-tasks/v1
spec: 0218-the-jev-router-on-docs-and-chore-tasks
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_07
      file: task_07.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_07]
---

# Tasks — The Jev Router on docs and chore Tasks

| id      | title                                                                         | type    | complexity | needs                              |
| ------- | ----------------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | The Jev Router is one OpenCode selection that only Project Config can name    | backend | medium     | —                                  |
| task_02 | The month's Jev spend counts the judge, the router and the key's own usage    | backend | medium     | —                                  |
| task_03 | Every routed prompt runs under the shared ceiling and leaves a Judge Log line | backend | high       | task_01, task_02                   |
| task_04 | The Jev Router, measured on four replayed docs and chore Tasks                | chore   | high       | task_01, task_02, task_03          |
| task_07 | A routed prompt starts only under a monthly key limit within the ceiling | backend | medium | task_03 |
| task_05 | Run the final QA gate | qa | high | task_01, task_02, task_03, task_04, task_07 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_07 · 5 → task_05

The measurement Task 04 specified was run on the maintainer's machine by the operator on 2026-10-02, after the Agent sandbox's approval review refused to send repository content to OpenRouter; its record is `measurement/jev-router.md`.
