---
schema: spec-tasks/v1
spec: 0057-baseline-capability-evidence-and-retention
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
    - id: task_06
      file: task_06.md
      needs: [task_05]
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_01]
    - id: task_09
      file: task_09.md
      needs: [task_08]
    - id: task_10
      file: task_10.md
      needs: [task_01]
    - id: task_11
      file: task_11.md
      needs: [task_05, task_09, task_10]
    - id: task_12
      file: task_12.md
      needs: [task_07, task_11]
---

# Tasks — Baseline capability evidence and retention

| id      | title                                            | type    | complexity | needs                     |
| ------- | ------------------------------------------------ | ------- | ---------- | ------------------------- |
| task_01 | Record how planning behaves today                 | test    | medium     | —                         |
| task_02 | Resolve symlinked executables without running them| backend | medium     | task_01                   |
| task_03 | Carry probe evidence into the divergence          | backend | medium     | task_02                   |
| task_04 | Render the probe and group by requirement strength| backend | high       | task_03                   |
| task_05 | Satisfy a portable Verification role from the repo| backend | medium     | task_04                   |
| task_06 | Offer a read-only capability re-check             | backend | high       | task_05                   |
| task_07 | Add the remediate-and-re-run outcome              | backend | medium     | task_06                   |
| task_08 | Account for clauses when a Profile drifts         | backend | high       | task_01                   |
| task_09 | Render the clause-level delta before apply        | backend | medium     | task_08                   |
| task_10 | Warn only about carriers nobody manages           | backend | high       | task_01                   |
| task_11 | Report the result as a status matrix              | backend | medium     | task_05, task_09, task_10 |
| task_12 | Document the evidence and retention contract      | docs    | medium     | task_07, task_11          |

Waves: 1 → task_01 · 2 → task_02, task_08, task_10 · 3 → task_03, task_09 ·
4 → task_04 · 5 → task_05 · 6 → task_06, task_11 · 7 → task_07 · 8 → task_12

The graph fans out by file ownership rather than by feature. Tasks 02–07 all
edit the profile alignment surface and form one chain; 08 and 09 edit plan
resolution and the retention contract; 10 edits carrier classification. Those
three groups touch disjoint files, so they run concurrently without putting two
Agent Sessions in the same code. Task 11 is the join.

task_01 changes no behavior and lands first. It records today's plan outcomes
and diagnostics for real repository shapes, so the retention gate in task_08 —
the one slice that turns a completing path into a stopping one — is measured
against what works today. Captured later, the corpus would encode the new
behavior and prove nothing.

The retention chain (08, 09) is isolated on purpose. It carries this Spec's
only regression risk, and a stall there leaves the capability-evidence chain
and carrier classification free to land.
