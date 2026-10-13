---
schema: spec-tasks/v1
spec: 0258-authoring-that-proves-what-it-claims
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

# Tasks — Authoring that proves what it claims

| id      | title                                                                                   | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Spec check reports a TechSpec without its absence evidence, ADR alternatives or cross-read | backend | medium     | —                |
| task_02 | Spec check reports a pending Task whose consumers no Verification runs                   | backend | medium     | task_01          |
| task_03 | The authoring skills state the four proof rules                                          | docs    | medium     | —                |
| task_04 | The command reference, the user guide and the glossary name the four proofs              | docs    | low        | task_02, task_03 |
| task_05 | Run the final QA gate                                                                    | qa      | high       | task_04          |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05

task_01 and task_02 share `internal/speccheck/coherence.go`,
`internal/speccheck/citations.go` and the three corpus pins, so they run in
series. task_03 changes only the two authoring skills, their templates, the
skill contract test and the version record, so it runs beside task_01.
task_04 shares the version record with task_03 and documents the codes
task_01 and task_02 add. The gate follows the single leaf. The operator queues
this Spec after the other active Specs of the cycle (`_prd.md` →
Prerequisites).
