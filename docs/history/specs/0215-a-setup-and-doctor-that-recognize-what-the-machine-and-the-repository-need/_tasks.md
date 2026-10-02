---
schema: spec-tasks/v1
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
qa: task_05
requires: [0211-a-delivery-queue-that-finishes-without-intervention]
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
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Doctor reports the forge and Git with coded findings |
| task_02 | backend | Doctor reports the toolchain and the environment, and this repository declares what it needs |
| task_03 | backend | A required upstream skill that trails its Setup Snapshot is reported and restored |
| task_04 | backend | Setup reports the readiness lines and deliver start refuses when this machine cannot publish |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_01, task_02 and task_03 all change `internal/cli/doctor.go`, its tests
and the `doctor` guide, so they run in series. task_04 is the only Task that
edits the Roundfix Skill, after every behavior it describes exists, so the
skill version rises once.
