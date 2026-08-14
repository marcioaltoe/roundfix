---
schema: spec-tasks/v1
spec: 0075-typed-docs-backlog
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A typed backlog for triage intent

| id      | title                                            | type    | complexity | needs            |
| ------- | ------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Give the layout a home for typed intent           | backend | medium     | —                |
| task_02 | Re-record the corpus and declare what moved       | test    | medium     | task_01          |
| task_03 | Adopt the contract in this repository             | docs    | low        | task_02          |
| task_04 | Name the backlog vocabulary in the glossary       | docs    | low        | task_01          |
| task_05 | Run the final QA gate                             | qa      | medium     | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_04 · 3 → task_03 · 4 → task_05

task_01 carries the whole contract because the clauses are one edit to one
module asset; splitting them would leave the guide half-written between Tasks.

task_02 is separate and typed `test` because the corpus is a characterization
gate: re-recording it is sanctioned but must be **reviewed**, and only layout
content and digests may move. Anything else moving means the edit leaked, and
that is a distinct piece of work from writing the clauses.

The maintainer settled on 2026-08-05 that Baseline module assets are product
content rather than protected tooling, so no Task here is a tooling Task and no
grant is cited. Recorded in
`docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.

The QA gate is authored as task_05 and depends on both non-QA leaves.
