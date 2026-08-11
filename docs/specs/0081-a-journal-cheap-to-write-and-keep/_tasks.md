---
schema: spec-tasks/v1
spec: 0081-a-journal-cheap-to-write-and-keep
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
      needs: [task_04]
    - id: task_08
      file: task_08.md
      needs: [task_07]
    - id: task_10
      file: task_10.md
      needs: [task_04, task_05]
    - id: task_09
      file: task_09.md
      needs: [task_06, task_08, task_10]
---

# Tasks — A journal cheap to write and cheap to keep

| id      | title                                              | type     | complexity | needs            |
| ------- | -------------------------------------------------- | -------- | ---------- | ---------------- |
| task_01 | Measure what the journal actually costs             | test     | medium     | —                |
| task_02 | Give the writer one transaction discipline          | data     | high       | task_01          |
| task_03 | Stop scanning the journal at every Run start        | data     | medium     | task_02          |
| task_04 | Amortize appends across a batch                     | data     | high       | task_03          |
| task_05 | Read only what the reader uses                      | data     | medium     | task_04          |
| task_06 | Make cockpit cost track new events                  | frontend | medium     | task_05          |
| task_07 | Prove parallel Runs at the pre-raise timeout        | test     | medium     | task_04          |
| task_08 | Decide the retention shape on the measurement       | docs     | medium     | task_07          |
| task_10 | Measure and record the before that the after is compared against | test | high | task_04, task_05 |
| task_09 | Run the final QA gate                               | qa       | medium     | task_06, task_08 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05,
task_07 · 6 → task_06, task_08 · 7 → task_09

**This graph is a chain where the TechSpec's Build Order reads parallel, and
that is deliberate.** Steps 2 through 5 are logically independent of one
another, and every one of them edits the store's journal and connection code.
Task Worktrees editing one file in the same wave conflict at integration, so
the graph pays waves to buy integration correctness. Stating it here means a
reader comparing the graph against the Build Order does not read the
serialization as an accident. The only genuine parallelism left is wave 5,
where the read projection meets the parallel-Run proof — the proof adds its own
test file and touches nothing the projection touches.

**task_01 is not ceremony.** The Spec's own decision is that the
retention-shape question waits behind a measurement, and the answer may
legitimately be that no payload is shed at all. A graph that let the decision
start before the proof would invite settling the ADR-0008 question from the
armchair, which is the one outcome this Spec forbids. Every later Task cites
the baseline task_01 records.

**No Task in this graph carries protected tooling.** The work is Go, SQL,
Go-side config defaults, and a user-guide document; this repository's own
`.roundfixrc.yml` is out of scope and defaults ship in code.

**task_08 is where the Spec may end smaller than its title.** The retention
shape is decided on the re-measurement, and concluding that no payload is shed
— leaving ADR-0008 untouched — is a correct outcome rather than a shortfall.
It is a node so that conclusion has an owner and a recorded reasoning either
way.

The QA gate is authored as task_09 and depends on task_06 and task_08, the
graph's only non-QA leaves.
