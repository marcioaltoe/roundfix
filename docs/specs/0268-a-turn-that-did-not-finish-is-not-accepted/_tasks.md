---
schema: spec-tasks/v1
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A turn that did not finish is not accepted

| id      | title                                                                                                                              | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The glossary, the guides and the Roundfix Skill describe the Provider System Error, the refused background tasks and the Unfinished Turn | docs    | low        | —                |
| task_02 | The ACPX Runner never accepts a turn that the provider failed                                                                      | backend | medium     | —                |
| task_03 | The Daemon recovers a Provider System Error through the Fallback Chain and the queue parks an unrecovered one without counting its retry | backend | high       | task_01, task_02 |
| task_04 | Claude sessions refuse background tasks and a QA turn that leaves one is an Unfinished Turn                                        | backend | medium     | task_03          |
| task_05 | Run the final QA gate                                                                                                              | qa      | high       | task_04          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 and task_02 share no declared file. task_03 reads
`DescribeProviderSystemError` from task_02 and changes the `deliver` surface
whose skill reference task_01 writes; it produces and consumes the
`provider_system_error` Run Event in one Task. task_04 edits
`internal/agent/acpx_runner.go` after task_02 and
`internal/daemon/task_engine.go` after task_03, so it follows task_03. The gate
follows the only leaf.
