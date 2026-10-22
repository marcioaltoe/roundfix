---
schema: spec-tasks/v1
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A subsumed test selector is flagged and a QA refusal names its commit

| id      | title                                                                                                         | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Spec Consistency Check reports a test selector its own invocation runs whole                              | backend | medium     | —                |
| task_02 | The spec guide, the Roundfix spec reference, the write-tasks skill and the glossary state the subsumed selector rule | docs    | low        | task_01          |
| task_03 | A QA Report the Daemon closes without an Agent records the commit it audited                                  | backend | medium     | —                |
| task_04 | The qa-gate refusal template, the Context-driven development guide and the glossary name the audited commit   | docs    | low        | task_02, task_03 |
| task_05 | Run the final QA gate                                                                                         | qa      | high       | task_04          |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05

task_01 and task_03 change disjoint files (`internal/speccheck`, and
`internal/spec` with `internal/daemon`), so they share Wave 1. task_02
documents the code task_01 ships and follows it. task_04 documents the key
task_03 writes and shares `CONTEXT.md`, `skills/testdata/owned-skill-versions.json`
and `docs/references/behavior-surfaces.json` with task_02, so it follows both.
task_01 and task_03 each prove their contracts in the packages that consume
them (`internal/speccheck`; `internal/spec` and `internal/daemon`), so no
separate Integration Checkpoint is needed. The gate follows the one leaf,
task_04. No other Spec is required.
