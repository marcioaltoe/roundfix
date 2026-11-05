---
schema: spec-tasks/v1
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — A numbered file collision between queued Specs parks before the Run

| id      | title                                                                                         | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Project Config declares the numbered files under delivery.numbered_files                       | backend | medium     | —                |
| task_02 | The queue parks a colliding item before its Run and a retry refuses while it collides          | backend | medium     | —                |
| task_03 | CONTEXT.md, ADR-0367, the guides and the deliver reference state the rule                       | docs    | medium     | task_01, task_02 |
| task_04 | Delivery Revalidation reads the numbered files of the item, the default branch and earlier items | backend | high       | task_03          |
| task_05 | Run the final QA gate                                                                         | qa      | medium     | task_04          |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 and task_02 share no file: task_01 changes `internal/config` and the
Behavior Surface Record (the `config keys` fingerprint reads the YAML tags),
task_02 changes `internal/delivery`. task_03 writes the documents: it needs
task_02, because the park action contract reads the next action from
`ClassifyPark`, and task_01, because both re-record the Behavior Surface
Record. task_04 changes the command workflow, a CLI surface, and needs
task_03 so that the `deliver` guide and reference describing it come first;
it also calls `(Delivery).IsNumberedFile` and fills
`Revalidation.NumberedFiles`. Each implementation Task runs the existing
suites that consume its changed contract in its own Verification, so no
separate Integration Checkpoint is needed. No other Spec is required.
