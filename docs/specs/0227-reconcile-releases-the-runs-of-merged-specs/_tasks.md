---
schema: spec-tasks/v1
spec: 0227-reconcile-releases-the-runs-of-merged-specs
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: [task_04]
    - id: task_02
      file: task_02.md
      needs: [task_04]
    - id: task_03
      file: task_03.md
      needs: [task_01, task_04]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Merge evidence releases a terminal Run of a merged Spec whose files diverged later |
| task_02 | backend | The post-merge cleanup proves a squash merge onto a moved default branch |
| task_03 | backend | Reconcile lists and releases the item branches of merged Specs |
| task_04 | docs | The Roundfix Skill and the reconcile and deliver guides describe merge evidence and item branches |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_04 · 2 → task_01, task_02 · 3 → task_03 · 4 → task_05

task_04 writes, from the TechSpec, the guides of the two command surfaces the
other Tasks change, so it runs first. task_01 changes the Run proof in the
worktree package and task_02 the cleanup proof in the delivery workflow; they
share no declared file and run together. task_03 reuses task_01's merge
evidence for item branches, so it follows task_01. The gate follows all four.
