---
schema: spec-tasks/v1
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
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
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — Acceptance rules that match the authorization and the queue

| id      | title                                                                                  | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A partial blocked only by the Spec's own network prohibition qualifies                 | backend | medium     | —                |
| task_02 | Spec check reports a QA requirement that pins a byte invariant to a commit             | backend | medium     | —                |
| task_03 | The Baseline guides and the glossary name both acceptance rules                        | docs    | medium     | task_01, task_02 |
| task_04 | The settlement skills, write-tasks and the user guide state both acceptance rules      | docs    | medium     | task_01, task_02 |
| task_05 | Run the final QA gate                                                                  | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05

task_01 changes only `internal/spec/qa.go` and the three policy test files,
and task_02 only `internal/speccheck` and the corpus pins, so they run
together. task_03 changes the Baseline modules, their derived files, the two
repository guides and `CONTEXT.md`; task_04 changes the owned skills, their
mirrors and version record, the user guide and the Behavior Surface Record.
They share no file and both document what task_01 and task_02 deliver, so
they run together after both. The gate follows the two leaves. The operator
rebases this Spec after Specs 0266 to 0269 (`_prd.md` → Prerequisites).
