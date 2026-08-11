---
schema: spec-tasks/v1
spec: 0062-baseline-digest-regeneration-bootstrap
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
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
---

# Tasks — Baseline digest regeneration bootstrap

| id      | title                                            | type    | complexity | needs   |
| ------- | ------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Capture the diagnostic characterization corpus   | test    | medium     | —       |
| task_02 | Give the catalog loader a regeneration mode      | backend | high       | task_01 |
| task_03 | Break the regeneration cycle on the update path  | backend | high       | task_02 |
| task_04 | Re-validate strictly after regeneration          | infra   | medium     | task_03 |
| task_05 | Say what the regenerator cannot supply           | backend | low        | task_04 |
| task_06 | Document the regeneration contract               | docs    | low        | task_05 |
| task_07 | Load the catalog once on the regeneration path   | backend | high       | task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 ·
6 → task_06 · 7 → task_07

The graph is a chain. Tasks 02, 03, and 05 all edit the catalog validation
surface in the same package, and task_04 depends on the regeneration path
existing before it can re-validate it. Parallel waves would put two Agent
Sessions in the same files.

task_01 is deliberately first and changes no behavior: it captures today's
diagnostics as goldens so every later task is measured against the current
truth. Captured after a behavior change, the corpus would encode the new
behavior and prove nothing.
