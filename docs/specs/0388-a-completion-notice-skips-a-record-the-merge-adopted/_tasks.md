---
schema: spec-tasks/v1
spec: 0388-a-completion-notice-skips-a-record-the-merge-adopted
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — A completion notice skips a record the merge adopted

| id      | title                                                                                                      | type | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------------------------------- | ---- | ---------- | ------- |
| task_01 | A completion notice skips a record the merge adopted, and every guarantee of the section is pinned          | docs | low        | —       |
| task_02 | Run the final QA gate                                                                                      | qa   | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 rewrites one section of one skill reference, syncs its mirror, raises
and records the skill's version, and adds the pins and the executable check; it
runs the whole `./skills` package that consumes the shipped skill and the
version record in its own Verification, so no separate Integration Checkpoint
is needed. The gate follows it. No other Spec is required; Spec 0389 delivers
after this one and also raises the `roundfix` version.
