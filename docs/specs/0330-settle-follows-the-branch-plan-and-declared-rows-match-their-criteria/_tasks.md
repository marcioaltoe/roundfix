---
schema: spec-tasks/v1
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
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

# Tasks — Settle follows the branch's plan and the Task commit's rules, and a declared row matches the criterion it names

| id      | title                                                                             | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | QA acceptance matches each declared-blocked row to the criterion a declaration names | backend | high       | —       |
| task_02 | Settle runs from the user checkout                                                | backend | medium     | task_01 |
| task_03 | Settle follows the Spec plan its Run's branch carries                             | backend | high       | task_02 |
| task_04 | Settle commits what a Task commit keeps                                           | backend | medium     | task_03 |
| task_05 | Run the final QA gate                                                             | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

Each Task is a vertical slice through one rule's code and the documents that
describe it. All four write `CONTEXT.md`, the Roundfix `SKILL.md` version
pair, the Roundfix `settle` reference, `skills/testdata/owned-skill-versions.json`
and `docs/references/behavior-surfaces.json`, and task_02 to task_04 all
write `internal/cli/settle.go` and the `settle` guide, so they run in series;
task_03's advance test also uses task_01's matching. No other Spec is
required.
