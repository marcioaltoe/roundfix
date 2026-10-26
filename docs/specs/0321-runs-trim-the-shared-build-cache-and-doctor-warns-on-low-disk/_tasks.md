---
schema: spec-tasks/v1
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: [task_03]
    - id: task_02
      file: task_02.md
      needs: [task_03]
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — Runs trim the shared build cache, and doctor warns on low disk

| id      | title                                                                 | type    | complexity | needs                     |
| ------- | --------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Implement Runs and the queue gate trim the declared Shared Build Caches | backend | high       | task_03                   |
| task_02 | Doctor warns when a worktree or cache filesystem is low on disk        | backend | medium     | task_03                   |
| task_03 | Document the automatic trims and the disk-space warning                | docs    | medium     | —                         |
| task_04 | Run the final QA gate                                                  | qa      | high       | task_01, task_02, task_03 |

Waves: 1 → task_03 · 2 → task_01, task_02 · 3 → task_04

task_01 changes the `gc`, `implement` and `deliver` command code and task_02
the `doctor` command code, whose guides and skill references task_03 writes,
so both need task_03, as a Task naming a CLI surface names its guide through a
Task it depends on. The three share no file: task_01 and task_02 touch
disjoint Go files and tests, and task_03 documents, records the skill and the
Behavior Surface fingerprint only, with a Verification that reads documents
and records. The gate follows all three. No other Spec is required.
