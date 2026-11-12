---
schema: spec-tasks/v1
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
qa: task_04
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
---

# Tasks — Roundfix reads a version 2 skills lock and keeps its version

| id      | title                                                                                         | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Every skills lock reader accepts versions 1 and 2 and hashes in the lock's own version         | backend | medium     | —       |
| task_02 | Doctor and the Baseline update name roundfix upgrade for an unsupported lock version           | backend | low        | task_01 |
| task_03 | The guides, Roundfix skill references and glossary describe skills lock versions 1 and 2       | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                                         | qa      | low        | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 runs, in its own Verification, the lock, restore, retired-skill,
reconcile and owned-skill tests of `./internal/baseline`, all of `./skills`
and `./internal/skillhash`, and the Doctor, skills, lock and release-plan
tests of `./internal/cli`, which consume the changed readers and hashes, so
no separate Integration Checkpoint is needed; task_02 runs the same
`./internal/cli` set over its next actions. The docs come last so their text
describes the delivered next actions. No other Spec is required; the queue
delivers Specs 0384 and 0385 after this one.
