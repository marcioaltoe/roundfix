---
schema: spec-tasks/v1
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A transient profile-proof timeout does not park the Delivery Queue

| id      | title                                                        | type    | complexity | needs            |
| ------- | ------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | The profile proof's setup deadline comes from config          | backend | medium     | —                |
| task_02 | A timed-out profile proof gets a third attempt after a backoff | backend | medium     | task_01          |
| task_03 | Cleanup of an absent disposable Agent Session is silent       | backend | medium     | task_01          |
| task_04 | Document the setup deadline and the third attempt             | docs    | low        | task_01, task_02 |
| task_05 | Run the final QA gate                                         | qa      | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05

task_01 carries the deadline from config to the agent's setup context and
documents the key.
task_02 adds the attempts and the backoff in the CLI, and task_03 silences
the cleanup of an absent Session in the agent; they share no file and run
together. task_04 documents what task_01 and task_02 deliver. The gate needs
every leaf. No other Spec is required.
