---
schema: spec-tasks/v1
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A session starts with an Update Check, and the update mode applies it

| id      | title                                                                     | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Project Config chooses notify or auto in update.mode                      | backend | medium     | —       |
| task_02 | roundfix update check reports a newer release and the Baseline plan       | backend | high       | task_01 |
| task_03 | The Roundfix skill's update mode applies the plan it showed               | docs    | medium     | task_02 |
| task_04 | The Baseline tells each session to run the Update Check before its first task | docs | medium     | task_03 |
| task_05 | Run the final QA gate                                                     | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series. task_02 reads task_01's key and shares with it the
Behavior Surface Record; task_03 runs task_02's command and both change the
Roundfix `SKILL.md` pair, its version record and `CONTEXT.md`; task_04's
clause names task_03's mode. Each implementation Task runs the suites that
consume its changed contract in its own Verification, so no separate
Integration Checkpoint is needed. No other Spec is required; Specs 0360, 0362
and 0363 deliver first and edit other sections of the shared files named in
`_prd.md` → Prerequisites.
