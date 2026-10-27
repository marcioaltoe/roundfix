---
schema: spec-tasks/v1
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04]
---

# Tasks — The skills lock matches what is installed, and a Run without a QA gate settles green

| id      | title                                                                                      | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | The guides, skills and glossary describe the lock rewrite and the gate-less settlement      | docs    | medium     | —                |
| task_02 | The bundle install rewrites the lock hash of each installed owned skill                     | backend | medium     | task_01          |
| task_03 | Every Task settles on the repository Verification, with or without a QA gate                | backend | medium     | task_01          |
| task_04 | A Run without a QA gate checks its integrated tree before it ends Clean                     | backend | high       | task_03          |
| task_05 | Run the final QA gate                                                                       | qa      | high       | task_02, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05

task_01 writes the guides, skills and glossary first, so each code Task names
its command's guide through that dependency; it records the owned-skill
versions and the Behavior Surface fingerprint last. task_02 touches only the
Baseline lock function and the two install callers, so it runs beside
task_03. task_04 follows task_03 because both edit
`internal/daemon/task_engine.go` and
`internal/cli/implement_settlement_checks_test.go`. Each code Task runs the
existing settlement, implement and baseline update tests that consume its
contract, so no separate Integration Checkpoint is needed. No other Spec is
required.
