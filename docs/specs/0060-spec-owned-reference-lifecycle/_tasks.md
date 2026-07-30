---
schema: spec-tasks/v1
spec: 0060-spec-owned-reference-lifecycle
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
---

# Tasks — Spec-owned reference lifecycle

| id      | title                                                  | type | complexity | needs            |
| ------- | ------------------------------------------------------ | ---- | ---------- | ---------------- |
| task_01 | Own the adoption lifecycle in the authorial Skills      | docs | high       | —                |
| task_02 | Land the layout, the routing rule, and the glossary     | docs | low        | task_01          |
| task_03 | Rehearse the lifecycle end to end and prove the gates   | docs | medium     | task_01, task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

All five authorized Skill pairs move in `task_01`, deliberately as one Task.
Every Skill edit rewrites the same seven derived digest artifacts through
`make baseline-digests`, so splitting the Skill work across Tasks would make
ADR-0026's serialized cherry-pick integration conflict on files neither Task
authored. `task_02` and `task_03` touch no Skill and no derived pin.
