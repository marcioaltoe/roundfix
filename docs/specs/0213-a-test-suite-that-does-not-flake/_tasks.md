---
schema: spec-tasks/v1
spec: 0213-a-test-suite-that-does-not-flake
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_03, task_04]
---

# Tasks — A test suite that does not flake

| id      | title                                                              | type | complexity | needs                     |
| ------- | ------------------------------------------------------------------ | ---- | ---------- | ------------------------- |
| task_01 | The linger and Run Budget tests wait on events, not elapsed time   | test | low        | —                         |
| task_02 | Script fixtures are the compiled test binary, and a guard holds the class | test | high | —                    |
| task_03 | Fixture processes end with the test binary that started them       | test | high       | task_02                   |
| task_04 | The Assets Sync template lives in memory                           | test | medium     | —                         |
| task_05 | Run the final QA gate                                              | qa   | high       | task_01, task_03, task_04 |

Waves: 1 → task_01, task_02, task_04 · 2 → task_03 · 3 → task_05
