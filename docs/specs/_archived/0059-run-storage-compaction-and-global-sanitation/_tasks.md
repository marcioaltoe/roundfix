---
schema: spec-tasks/v1
spec: 0059-run-storage-compaction-and-global-sanitation
qa: task_06
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
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
    - id: task_07
      file: task_07.md
      needs: [task_05]
    - id: task_06
      file: task_06.md
      needs: [task_07]
---

# Tasks — Run storage compaction and global sanitation

| id      | title                                          | type    | complexity | needs                     |
| ------- | ---------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Report where the bytes are                      | backend | medium     | —                         |
| task_02 | Compact only when nothing can be writing        | backend | high       | task_01                   |
| task_03 | Discover and classify every Artifact Root       | backend | high       | task_01                   |
| task_04 | Give every durable table a stated lifecycle     | backend | medium     | task_01                   |
| task_05 | Synchronise the Roundfix Skill                  | chore   | low        | task_02, task_03, task_04 |
| task_07 | Ship the compaction command and document it     | backend | medium     | task_05                   |
| task_06 | Run the final QA gate                           | qa      | medium     | task_07                   |

Waves: 1 → task_01 · 2 → task_02 · task_03 · task_04 · 3 → task_05 ·
4 → task_07 · 5 → task_06

task_01 leads because it is read-only and because it defines the measurement
vocabulary the other three assert against. Landing it first means nothing later
has to invent how a byte total is counted while also proving a mutation is
safe.

task_02, task_03 and task_04 are independent of each other and form the second
wave. They share task_01's report and touch different surfaces: the writer
connection, root discovery, and per-table policy.

No Verification in this graph asserts a recorded size. Each asserts a relation
that holds at any size — the preview equals the result, a second sanitation
reclaims zero, a refusal changes no byte, the totals reconcile within a
declared tolerance. On 2026-08-05 this repository spent three Runs on a test
asserting a literal `104` that a legitimate change moved to `106`, and the
stale literal masked the real diagnostic behind it. This graph is authored to
not repeat that.

task_05 is the one authorized tooling Task, bounded to
`.agents/skills/roundfix/SKILL.md` and its mirror under the 2026-08-04 record,
which names Spec 0059.

task_07 was appended on 2026-08-06 after the gate reported `fail`, closing
F-001 and F-002. The compaction guards shipped inside `internal/store` and
nothing exposed them: the TechSpec declared `roundfix gc compact` in its API
Contracts and no Task requirement ever asked for it, so the feature shipped as
an unreachable library. Naming the insertion here is what invalidates that
report.

The second lesson is recorded in task_07's Verification section. task_05
settled `completed` having changed nothing, because all four of its checks pass
most easily on an untouched repository — the shape Spec 0065's
`SC-VERIFY-WORK-INDEPENDENT` refuses, shipped the same night but authored after
this graph.

The QA gate is authored as task_06 and depends on task_07, the graph's only
non-QA leaf.
