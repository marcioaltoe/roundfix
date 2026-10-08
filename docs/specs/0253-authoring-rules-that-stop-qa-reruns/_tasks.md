---
schema: spec-tasks/v1
spec: 0253-authoring-rules-that-stop-qa-reruns
requires:
  - 0252-one-gate-per-run-and-a-smaller-baseline
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04]
---

# Tasks — Authoring rules that stop QA reruns

| id      | title                                                                                           | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | An outside-evidence source is reachable, a transcript is asserted whole, and a prerequisite is declared | backend | medium     | —                |
| task_02 | The authoring skills and the repository guide state the rules that stop QA reruns               | docs    | medium     | —                |
| task_03 | The loop re-checks parallel Specs and sweeps old wording, and Go tests stay hermetic           | backend | medium     | task_01          |
| task_04 | A retirement writes reduced history and deletes reviews and handoffs                            | backend | medium     | task_03          |
| task_05 | Run the final QA gate                                                                           | qa      | high       | task_02, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

The Spec requires `0252-one-gate-per-run-and-a-smaller-baseline`, which
changes the same modules, guides and `docs/agents/specific-repository.md`, so
the Delivery Queue starts this Spec only after 0252 merges. task_01, task_03
and task_04 change Baseline source and share
`internal/baseline/module-versions.json`, the derived catalog files, the
profile digest and `docs/agents/setup-context.json`, so they run in series
(Wave collision). task_02 changes only owned skills and
`docs/agents/specific-repository.md`, which no other Task declares, so it runs
beside task_01. The gate follows both leaves.
