---
schema: spec-tasks/v1
spec: 0266-the-review-record-survives-a-fix
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The review record survives a fix

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The guides, the skill and the glossary name the resolved key and the Reviewed Head Record | docs    | low        | —       |
| task_02 | The review record is keyed and compared by the resolved checkout                          | backend | medium     | task_01 |
| task_03 | Each reviewed head keeps its record, and the proof reads the candidate's                  | backend | medium     | task_02 |
| task_04 | A code correction inside the archived Spec's declared paths returns to review             | backend | medium     | task_03 |
| task_05 | Run the final QA gate                                                                     | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 writes the review and deliver contract into the guides, the Roundfix
skill and `CONTEXT.md` first, so each code Task changes a CLI surface whose
guide a Task it depends on already names. task_02, task_03 and task_04 all
edit `internal/cli/deliver_review_correction.go`, and task_02 and task_03 also
`internal/cli/review.go`, so they run in series. Every consumer of the changed
record key, proof and declared-path reader is in `./internal/cli` or
`./internal/spec`, which the Tasks run in their Verification. The gate follows
task_04.
