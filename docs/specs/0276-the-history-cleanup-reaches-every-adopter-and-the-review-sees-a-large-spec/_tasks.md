---
schema: spec-tasks/v1
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
qa: task_05
requires:
  - 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The history cleanup reaches every adopter, and the review sees a large Spec

| id      | title                                                                                                        | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The history cleanup writes a successor's lineage line, refuses only an unsucceeded ADR and repairs live links | backend | high       | —       |
| task_02 | The pre-PR review omits retired history and generated files, reviews a large candidate in segments, and exempts the QA Task from C2 | backend | high | task_01 |
| task_03 | Spec check reports a Task that edits an owned skill without its version record                              | backend | medium     | task_02 |
| task_04 | The archive plan has an offline form, and the Delivery Queue runs a fresh copy of each item binary          | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                                        | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

Every implementation Task raises the Roundfix skill version and re-records
`skills/testdata/owned-skill-versions.json` and the Behavior Surface Record,
and task_01 and task_04 both edit the Roundfix `archive` reference, so the four
run in series. task_04 changes `internal/cli/deliver_item_binary_test.go`,
which Spec 0275 rewrites, and extends ADR-0281, so the Spec requires 0275. The
gate follows the only leaf.
