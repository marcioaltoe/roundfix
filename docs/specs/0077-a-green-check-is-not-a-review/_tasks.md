---
schema: spec-tasks/v1
spec: 0077-a-green-check-is-not-a-review
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A green check is not a review

| id      | title                                          | type    | complexity | needs            |
| ------- | ---------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Make verified require proof a review ran        | backend | high       | —                |
| task_02 | Name the refusal so the stall is legible        | backend | medium     | task_01          |
| task_03 | Say what was observed instead of merging        | backend | low        | task_01, task_02 |
| task_04 | Synchronise the Roundfix Skill                  | chore   | low        | task_03          |
| task_05 | Run the final QA gate                           | qa      | medium     | task_04          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

A chain, deliberately. task_01 closes the gate on its own — including for
refusal shapes task_02 has not learned yet — so it must land first and must be
provable without task_02 existing. Building recognition first would leave the
default open while looking finished.

task_04 is the one authorized tooling Task, bounded to
`.agents/skills/roundfix/**` and its mirror under the 2026-08-04 standing grant
for CLI-contract synchronisation.

The QA gate is authored as task_05 and depends on task_04, the graph's only
non-QA leaf.
