---
schema: spec-tasks/v1
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
qa: task_03
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
---

# Tasks — A stale user-level skill copy is reported with its refresh command

| id      | title                                                                                      | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | skills check, Doctor and the Update Check name a Stale User Skill Copy and write nothing   | backend | high       | —       |
| task_02 | The guides, the Roundfix references, the --update mode and the glossary describe the report | docs    | medium     | task_01 |
| task_03 | Run the final QA gate                                                                      | qa      | medium     | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The Tasks run in series. task_02 documents the strings task_01 prints and
asks the `--update` mode to read the field task_01 adds. task_01 runs the
suites that consume its changed surfaces in its own Verification, so no
separate Integration Checkpoint is needed. No other Spec is required; Specs
0373, 0374, 0375 and 0377 deliver first and edit other sections of the shared
files named in `_prd.md` → Prerequisites.
