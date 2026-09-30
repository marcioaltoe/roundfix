---
schema: spec-tasks/v1
spec: 0199-stack-rules-that-say-what-they-mean
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
      needs: [task_01, task_02, task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The Bun and TypeScript rules say what the profile runs |
| task_02 | backend | The core rules hold in a repository with several languages |
| task_03 | backend | Backend and frontend guides name their workspace, and the stated HTTP suggestion is REST |
| task_04 | qa | Run the final QA gate |
