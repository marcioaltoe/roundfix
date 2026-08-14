---
schema: spec-tasks/v1
spec: 0103-a-suite-that-leaks-nothing
qa: task_09
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: []
    - id: task_07
      file: task_07.md
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: []
    - id: task_10
      file: task_10.md
      needs: [task_04]
    - id: task_09
      file: task_09.md
      needs: [task_02, task_05, task_07, task_08, task_10]
---

# Tasks — A suite that leaks nothing

| id      | title                                                      | type    | complexity | needs                                  |
| ------- | ---------------------------------------------------------- | ------- | ---------- | -------------------------------------- |
| task_01 | Compile the fixtures the suite executes                     | test    | high       | —                                      |
| task_02 | Fail in milliseconds when a fixture is already dead         | test    | medium     | task_01                                |
| task_03 | Give the suite a boundary it can assert                     | test    | medium     | —                                      |
| task_04 | Install the guard where the suite spawns                    | test    | medium     | task_03                                |
| task_05 | End the processes the detach tests prove survive            | test    | medium     | task_04                                |
| task_06 | Report what Roundfix left running                           | backend | high       | —                                      |
| task_07 | Prove the tree exited, not only its owner                   | backend | medium     | task_06                                |
| task_08 | Keep a gate's scratch state out of the evidence             | backend | low        | —                                      |
| task_10 | Let a sanctioned regeneration write what it declares         | test    | medium     | task_04                                |
| task_09 | Run the final QA gate                                       | qa      | high       | task_02, task_05, task_07, task_08, task_10 |

Wave plan: `1 → task_01, task_03, task_06, task_08 · 2 → task_02, task_04, task_07 · 3 → task_05, task_10 · 4 → task_09`.

task_10 was added on 2026-08-14, after task_04 installed the guard and two
sanctioned regeneration commands failed against it. It is a new node rather than
a widened task_04, because task_04 had already settled and reopening a completed
Task is how three Runs were spent during Spec 0094.

Waves are serialized past the dependency minimum where two Tasks would edit the
same file: task_02 follows task_01 because both change the adapter harness, and
task_05 follows task_04 because both change `internal/cli` test entry points. A
wave collision cost a full Agent turn during Spec 0094 and is cheaper to prevent
here than to diagnose there.
