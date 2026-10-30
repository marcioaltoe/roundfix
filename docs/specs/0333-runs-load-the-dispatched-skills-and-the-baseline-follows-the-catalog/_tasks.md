---
schema: spec-tasks/v1
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
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

# Tasks — Runs load the dispatched skills, and the Baseline follows the curated catalog

| id      | title                                                                                   | type    | complexity | needs            |
| ------- | --------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Baseline follows catalog 86df860 and retires the seven skills it removed            | backend | high       | —                |
| task_02 | The dispatch asks for person-started and polish skills only on request                  | backend | medium     | task_01          |
| task_03 | Run Agents load only the skills the dispatch selects and name them in their Result      | backend | high       | task_02          |
| task_04 | implement-spec retires and its two rules move into the Roundfix skill                   | backend | high       | task_03          |
| task_05 | Run the final QA gate                                                                   | qa      | high       | task_04          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01 moves the snapshots and the Retired Skills first, so the later Tasks
build on the new catalog. The graph is a chain: task_02 edits the modules,
the activation file, the derived files and this repository's dispatch guide
task_01 edits; task_03 names the dispatch guide task_02 renders and edits the
Roundfix skill, the owned-skill version record and the Behavior Surface
record task_01 edits; task_04 edits all of them again. Each code Task runs the whole package whose contract it
changes, so no separate Integration Checkpoint is needed. No other Spec is
required.
