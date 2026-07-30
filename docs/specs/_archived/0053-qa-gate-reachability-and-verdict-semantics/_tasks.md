---
schema: spec-tasks/v1
spec: 0053-qa-gate-reachability-and-verdict-semantics
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — QA gate reachability and verdict semantics

| id      | title                                                        | type    | complexity | needs                            |
| ------- | ------------------------------------------------------------ | ------- | ---------- | -------------------------------- |
| task_01 | Distinguish environment-blocked rows in the verdict contract  | backend | medium     | —                                |
| task_02 | State the Pull Request fact and a collision-safe report name  | backend | low        | task_01                          |
| task_03 | Classify a superseded QA-report-only Run Branch               | backend | high       | —                                |
| task_04 | Keep superseded branches out of automatic integration         | backend | medium     | task_03                          |
| task_05 | Align the Skill pairs, the guides, and the authorized digests | docs    | medium     | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01, task_03 · 2 → task_02, task_04 · 3 → task_05

`task_01` and `task_03` are file-disjoint by design (`internal/spec` plus
`internal/cli/archive_test.go` against `internal/worktree` plus
`internal/cli/cli_test.go`), so ADR-0026's serialized integration has no
conflict to resolve between them.
