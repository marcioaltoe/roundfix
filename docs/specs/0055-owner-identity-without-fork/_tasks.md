---
schema: spec-tasks/v1
spec: 0055-owner-identity-without-fork
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Owner identity without fork

| id      | title                                                    | type    | complexity | needs                            |
| ------- | -------------------------------------------------------- | ------- | ---------- | -------------------------------- |
| task_01 | Read owner start identity from the kernel, never a fork    | backend | high       | —                                |
| task_02 | Separate an unreadable identity from a proven mismatch     | backend | medium     | task_01                          |
| task_03 | Mark and warn about a Run created without reuse protection | backend | medium     | —                                |
| task_04 | Fix Stop argument order and add the supervised exit        | backend | medium     | task_02                          |
| task_05 | Document the diagnostics, the marker, and the exit         | docs    | low        | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05

`task_01` touches `internal/store/process*.go`; `task_03` touches
`internal/store/store.go` and `internal/cli/implement.go`. They are
file-disjoint, so ADR-0026's serialized integration has no conflict to resolve
between them.
