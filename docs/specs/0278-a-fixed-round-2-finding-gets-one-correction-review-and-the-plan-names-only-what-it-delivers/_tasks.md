---
schema: spec-tasks/v1
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_02, task_04]
---

# Tasks — A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers

| id      | title                                                                                   | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | A fixed round-2 finding gets one Correction Review, and the retry admits its fix         | backend | high       | —                |
| task_02 | The Spec Consistency Check refuses a Verification path the plan does not provide         | backend | high       | task_01          |
| task_04 | Every Task settlement reports a stale Setup Manifest                                     | test    | low        | —                |
| task_05 | Run the final QA gate                                                                    | qa      | high       | task_02, task_04 |

Waves: 1 → task_01, task_04 · 2 → task_02 · 3 → task_05

task_01 and task_02 each raise the Roundfix skill version and re-record
`skills/testdata/owned-skill-versions.json` and the Behavior Surface Record, so
they run in series. task_04 changes only
`internal/docscontract/publicdocs_test.go`, shares no file with them, and runs
beside task_01. The gate follows both leaves. There is no task_03: the
QA-delivery judgment it planned was measured and dropped before the Run
(ADR-0287), and the remaining ids are kept.
