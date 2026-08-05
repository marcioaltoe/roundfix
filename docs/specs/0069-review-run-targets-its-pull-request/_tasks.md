---
schema: spec-tasks/v1
spec: 0069-review-run-targets-its-pull-request
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

# Tasks — A Review Run targets its Pull Request

| id      | title                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Compare the checkout against its Pull Request   | backend | medium     | —       |
| task_02 | Keep the target while the Run writes            | backend | high       | task_01 |
| task_03 | Synchronise the Roundfix Skill                  | chore   | low        | task_02 |
| task_04 | Run the final QA gate                           | qa      | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

A chain. task_01 closes the failure that costs a round — a Run acting on the
wrong branch — and it closes it alone, because `preflight.Run` already holds
both values and needs only the comparison. task_02 closes the one that costs a
Run, and it is the larger slice: a recorded target, a re-read at every write
boundary, and a terminal outcome that did not exist.

Authoring disproved the PRD's premise that a late mid-Run check exists to be
moved. `checkout branch mismatch` is not in the tree; the nearest live code is
Round artifact reuse. task_02 therefore builds the guard rather than relocating
one, and its complexity is `high` for that reason.

task_03 is the one authorized tooling Task, bounded to
`.agents/skills/roundfix/**` and its mirror under the 2026-08-04 standing grant
for CLI-contract synchronisation, whose covered list names Spec 0069. The PRD's
tooling-authority row was corrected from `not applicable` in the same authoring
pass.

The QA gate is authored as task_04 and depends on task_03, the graph's only
non-QA leaf.
