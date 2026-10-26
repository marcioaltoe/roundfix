---
schema: spec-tasks/v1
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — Runs are not lost to an authorization row or a provider at capacity

| id      | title                                                                                   | type    | complexity | needs                              |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------------------------- |
| task_01 | A Network-Prohibited Row names what the authorization forbids                           | backend | medium     | —                                  |
| task_02 | A Provider System Error falls through the Fallback Chain during any Task or QA turn     | backend | high       | —                                  |
| task_03 | The Baseline clauses name the authorization row's subject                               | backend | medium     | —                                  |
| task_04 | The skills, guides and glossary say both rules                                          | docs    | medium     | —                                  |
| task_05 | Run the final QA gate                                                                   | qa      | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01, task_02, task_03, task_04 · 2 → task_05

The four implementation Tasks share no declared file and none reads another's
output: task_01 changes `internal/spec/qa.go` and the partial-policy tests,
task_02 changes `internal/daemon/agent_session_owner.go` and the provider
recovery tests, task_03 changes the two Baseline modules, their derived files,
`docs/agents/` and the clause test, and task_04 changes skills, user guides,
`CONTEXT.md`, the skill contract test and the Behavior Surface record. Skill
edits do not move the Baseline's derived files (measured: `make
baseline-digests` after a skill edit reported `"changed":false`). The gate
follows all four. No other Spec is required.
