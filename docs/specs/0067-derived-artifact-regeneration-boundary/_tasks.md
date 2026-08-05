---
schema: spec-tasks/v1
spec: 0067-derived-artifact-regeneration-boundary
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
      needs: [task_02, task_03]
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_05
      file: task_05.md
      needs: [task_06]
---

# Tasks — Derived artifact regeneration boundary

| id      | title                                              | type    | complexity | needs            |
| ------- | -------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Declare one owner per derived path and prove it exhaustive | backend | high | —          |
| task_02 | Prove every declared step actually rewrites its artifacts | test | medium | task_01          |
| task_03 | Make the sanctioned command cover what it claims    | infra   | medium     | task_01          |
| task_04 | Say what a human must do when the command cannot    | backend | low        | task_02, task_03 |
| task_06 | Make every record state what is measurably true      | backend | medium     | task_04          |
| task_05 | Run the final QA gate                               | qa      | medium     | task_06          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_06 · 5 → task_05

task_01 leads because exhaustiveness is what stops the next occurrence. Three
artifacts have already inherited this ambiguity, and a fourth instance happened
on 2026-08-05 while authoring the Spec queue: the flag was guessed from a test
name and the guess was wrong.

task_03 is the one authorized tooling Task, bounded to exactly `Makefile` under
the 2026-08-02 grant, which names 0067 for "the regeneration step list and
derived path scan".

The QA gate is authored as task_05 and depends on task_04, the graph's only
non-QA leaf.

## Corrective Task from the 2026-08-05 gate

Two findings, both real, both from under-specification.

F-001: task_03 named "the diagnostic characterization corpus" where Core
Feature 2 requires the sanctioned command to cover everything it claims — and
there are two dedicated corpora, not one. One sanctioned run still left the
gate red.

F-002: the records assert things the measurement contradicts. The
catalog-diagnostic record still said it was outside `BASELINE_DIGEST_STEPS`
after task_03 put it there, and the parity record claims nothing regenerates
the corpus while the sanctioned run rewrites two files beneath it.

That second half reaches past this Spec: the PRD asserts the parity corpus is
frozen, and the evidence says it is not. task_06 makes the records state what
is measurably true and deliberately does **not** decide whether that state is
desirable — freezing it for real would change what the sanctioned command
rewrites, which is a product decision the maintainer owns.

One corrective Task, within the ceiling `docs/agents/autonomous-work.md` sets.
task_05 returns to pending because a gate whose dependency closure grew is
invalid by construction.
