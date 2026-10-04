---
schema: spec-tasks/v1
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The Roundfix Skill and the command guides describe the item binary and the migration check |
| task_02 | backend | roundfix migrate --check says whether this binary would change the Run Database |
| task_03 | backend | Project Config declares the repository's Roundfix item build |
| task_04 | backend | The queue owner runs an item's steps with the binary the item builds |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05

task_02 and task_03 change the `migrate` command and the configuration
model, whose guides task_01 writes, so they follow it; they share no declared
file. task_04 calls the item binary's `migrate --check` and reads the
declaration, so it follows both, and its `deliver` guide comes from task_01
through them.
