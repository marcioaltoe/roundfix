---
schema: spec-tasks/v1
spec: 0091-a-proof-that-can-refuse
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_06
      file: task_06.md
      needs: [task_02]
    - id: task_03
      file: task_03.md
      needs: [task_06]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A proof that can refuse

| id      | title                                                   | type    | complexity | needs            |
| ------- | ------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Record that a nonexistent claude model proves passed    | test    | medium     | —                |
| task_02 | Read what a runtime offers before asking about it       | backend | medium     | task_01          |
| task_06 | Teach the sibling tests the catalogue read               | test    | medium     | task_02          |
| task_03 | Let membership decide the verdict                       | backend | high       | task_06          |
| task_04 | Stop appending a close error for a session never opened | backend | low        | —                |
| task_05 | Run the final QA gate                                   | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_04 · 2 → task_02 · 3 → task_06 · 4 → task_03 · 5 → task_05

Wave 1 is file-disjoint: task_01 writes only its own characterization file,
task_04 edits `acpx_runner.go` and its existing test. The serial chain
task_02 → task_03 exists because the verdict cannot be decided against a
catalogue that is not read yet, and both edit `selection_assignment.go`.

Task 06 is corrective, minted on 2026-08-10 after the first Run: task_02
changed the acpx invocation sequence and left sixteen sibling tests asserting
the old one. It sits between task_02 and task_03 so the chain continues from a
green tree.
