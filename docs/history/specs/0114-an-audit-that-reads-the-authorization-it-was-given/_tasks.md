---
schema: spec-tasks/v1
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
qa: task_06
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_07
      file: task_07.md
      needs: [task_05]
    - id: task_06
      file: task_06.md
      needs: [task_01, task_05, task_07]
---

# Tasks — An audit that reads the authorization it was given

| id      | title                                                     | type    | complexity | needs            |
| ------- | --------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Stop refusing the row the template teaches                  | backend | low        | —                |
| task_02 | Tell a governed path from an ordinary one                   | backend | high       | —                |
| task_03 | Hold the governed set to what was actually protected        | backend | medium     | task_02          |
| task_04 | Resolve a command's outputs from the tree                   | backend | high       | —                |
| task_05 | Let the audit judge the grant it was given                  | backend | high       | task_03, task_04 |
| task_07 | Teach the Daemon's fixture what governed now means          | test    | low        | task_05          |
| task_06 | Run the final QA gate                                       | qa      | high       | task_01, task_05, task_07 |

Wave plan: `1 → task_01, task_02, task_04 · 2 → task_03 · 3 → task_05 · 4 → task_06`.

task_05 is where the two repairs meet, which is why it lands alone: a failure
there is about their composition rather than about either one. task_01 is
independent and stays that way, so the author-facing refusal is gone whether or
not the audit work converges.

The shared regeneration reader this Spec planned as its first slice was already
delivered by Spec 0103's task_11 on 2026-08-14, which pointed the changed-path
audit at `suiteguardcontract.ParseSanctionedRegenerations`. The pre-work gate
audit found it passing before any Run was spent on it, and the node was dropped
rather than reopened as work already done.

task_07 was added on 2026-08-15 from the gate's own finding. task_05's regression
command covered the audit's neighbours and not its caller, so the Daemon fixture
that still asserted the old behaviour was found by the gate rather than by the
Task that changed it.
