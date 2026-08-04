---
schema: spec-tasks/v1
spec: 0068-spec-close-audit
qa: task_07
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
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_05, task_06]
---

# Tasks — Spec close audit

| id      | title                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Resolve integration by content when the branch is gone | backend | high   | —                |
| task_02 | Classify every survivor with its evidence          | backend | high       | task_01          |
| task_03 | Report what the Spec claims but has not delivered   | backend | medium     | task_02          |
| task_04 | Expose the audit as `roundfix spec audit`          | backend | medium     | task_02, task_03 |
| task_05 | Replay the session that motivated this Spec         | test    | medium     | task_03, task_04 |
| task_06 | Synchronise the Roundfix Skill with the new command | chore   | low        | task_04          |
| task_07 | Run the final QA gate                               | qa      | high       | task_05, task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 ·
5 → task_05, task_06 · 6 → task_07

task_01 leads because it is the only change to existing behaviour and the only
one that can lose work if it is wrong. Its content check proves integration and
never disproves it; everything downstream classifies on top of that asymmetry.

task_06 is the one authorized tooling Task. Its bounded files are exactly
`.agents/skills/roundfix/**` and the `skills/roundfix/**` mirror, under the
2026-08-04 standing grant for CLI-contract synchronisation. It is separate from
task_04 because an authorized tooling Task may mutate only its bounded files.

The QA gate is authored as task_07 and depends on both non-QA leaves.
