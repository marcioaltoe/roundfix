---
schema: spec-tasks/v1
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
---

# Tasks — The regeneration contract runs when its inputs change

| id      | title                                                                                     | type    | complexity | needs                     |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | verify-select lists every contract with -all and names the contracts it leaves out         | backend | medium     | —                         |
| task_02 | TestRegenerationIsDeclared runs when an input its regeneration reads changes               | backend | medium     | —                         |
| task_03 | make verify-contracts runs every contract, on every push to main and before every release | infra   | medium     | task_01                   |
| task_04 | The glossary and the repository rules describe the Full Contract Run                      | docs    | low        | —                         |
| task_05 | Run the final QA gate                                                                     | qa      | high       | task_02, task_03, task_04 |

Waves: 1 → task_01, task_02, task_04 · 2 → task_03 · 3 → task_05

task_03's target calls the `-all` mode that task_01 adds, so it follows
task_01. task_01 changes `internal/verifyselect/verifyselect.go` and
`contracts_test.go`, task_02 changes `contracts_repository_test.go` and the
regeneration test's header, task_03 creates `full_contract_run_test.go` and
changes the `Makefile` and the workflows, and task_04 changes `CONTEXT.md` and
the repository rules, so no two Tasks of one wave share a declared file. The
gate follows every leaf.
