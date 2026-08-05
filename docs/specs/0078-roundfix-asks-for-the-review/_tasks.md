---
schema: spec-tasks/v1
spec: 0078-roundfix-asks-for-the-review
qa: task_06
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01, task_03]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Roundfix asks for the review

| id      | title                                        | type    | complexity | needs            |
| ------- | -------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Publish one review request, once              | backend | high       | —                |
| task_03 | Refuse a pair that would strand the Run       | backend | medium     | —                |
| task_02 | Ask at the seam where the Round pushes        | backend | medium     | task_01, task_03 |
| task_04 | Turn the flow on for this repository          | chore   | low        | task_02          |
| task_05 | Synchronise the Roundfix Skill                | chore   | low        | task_04          |
| task_06 | Run the final QA gate                         | qa      | medium     | task_05          |

Waves: 1 → task_01 · task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05 ·
5 → task_06

task_01 and task_03 are deliberately independent. The requester is provable
with nothing wired to it, and the Preflight refusal is provable with no
requester in existence — each closes its own half of the failure. task_02 is
where they meet, and it is the only Task that can produce a duplicate request,
so it carries the one-per-Round assertions.

task_04 is the Task that must not land early. It commits the `.coderabbit.yaml`
already modified in the working tree, and that file takes effect when it
reaches the default branch — landing it before the requester works is the
window where every unattended Run stalls.

task_05 is the one authorized tooling Task, bounded to
`.agents/skills/roundfix/**` and its mirror under the 2026-08-04 standing grant
for CLI-contract synchronisation, whose covered list names Spec 0078.

The QA gate is authored as task_06 and depends on task_05, the graph's only
non-QA leaf.
