---
schema: spec-tasks/v1
spec: 0002-acpx-migration
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
---

# Tasks — ACPX Migration

| id      | title                                                                  | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | acpx invocation core: commands, NDJSON stream, exit-code mapping       | backend | high       | —                |
| task_02 | Agent Session lifecycle: ensure, set-mode, cancel, EndSession          | backend | medium     | task_01          |
| task_03 | Probe and Preflight Validation with the acpx version pin               | backend | low        | task_01          |
| task_04 | Run wiring: named session in both engines, close on every terminal path | backend | medium    | task_02, task_03 |
| task_05 | Cutover: SDK removal, module shrink, parity gate, gated integration test | backend | high      | task_04          |
| task_06 | Docs and skill: acpx pin, Node prerequisite, handoff item closed       | docs    | low        | task_05          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05 · 5 → task_06
