---
schema: spec-tasks/v1
spec: 0226-an-archived-spec-keeps-its-links
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
      needs: [task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The Roundfix Skill and the archive guide describe the link rewrite |
| task_02 | backend | The archive rewrites outward links and refuses a broken one |
| task_03 | backend | A resumed archive commit accepts the link rewrites |
| task_04 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_02 changes the command whose wording task_01 writes, so it follows it.
task_03 calls the match function task_02 adds to the archive, so it follows
task_02. The Tasks share no declared file.
