---
schema: spec-tasks/v1
spec: 0093-a-spec-that-validates-itself
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_08
      file: task_08.md
      needs: [task_04]
    - id: task_07
      file: task_07.md
      needs: [task_06, task_08]
---

# Tasks — A Spec that validates itself

| id      | title                                                  | type    | complexity | needs            |
| ------- | ------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Record that a false citation passes every check today  | test    | medium     | —                |
| task_02 | Read a cited decision against the claim made about it  | backend | high       | task_01          |
| task_03 | Let a caller ask for one authoring stage               | backend | medium     | task_01          |
| task_04 | Surface the stage through the command                  | backend | low        | task_02, task_03 |
| task_05 | Wire the checker into PRD and TechSpec authoring        | infra   | medium     | task_04          |
| task_06 | Let the QA gate read product instead of paperwork      | infra   | high       | task_05          |
| task_08 | Run the citation detector in the authoring stages      | backend | medium     | task_04          |
| task_07 | Run the final QA gate                                  | qa      | medium     | task_06, task_08 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05 · 5 → task_06, task_08 · 6 → task_07

Wave 2 is file-disjoint: task_02 owns the new detector and its test, task_03
owns the checker's entry point and its test.

Tasks 05 and 06 are serialised rather than parallel even though they edit
different skills: both end by running `make skills-sync`, which rewrites every
generated copy, so running them in one wave would put two Tasks on the same
generated files.

The QA gate is deliberately the smallest in recent Specs. That is the point of
the Spec: what it can prove by reading files, it proves during authoring, and
the gate keeps only what a file read cannot settle.

Task 08 is corrective, added on 2026-08-09 after the gate's F-001: the stage
registry never carried the citation detector, so the authoring scopes the skills
tell an author to run returned `No findings` for a claim the unscoped run
rejects. It is file-disjoint from task_06, which owns only the gate skill.
