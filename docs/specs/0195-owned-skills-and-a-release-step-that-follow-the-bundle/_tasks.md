---
schema: spec-tasks/v1
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
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
      needs: [task_03, task_06]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04, task_06]
    - id: task_06
      file: task_06.md
      needs: [task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | An installed owned skill older than the bundle is reported |
| task_02 | test | An owned skill's content cannot change under one version |
| task_03 | backend | The Roundfix skill is in every setup and has a dispatch trigger |
| task_04 | docs | Every release checks skills and guides first |
| task_05 | qa | Run the final QA gate |
| task_06 | chore | The Roundfix skill declares the version its merged content needs |
