---
schema: spec-tasks/v1
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
qa: task_04
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
---

# Tasks — A profile proof whose session answered passes when only its close fails

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The guides, the Roundfix skill and the glossary describe the close of a profile proof                   | docs    | medium     | —       |
| task_02 | A proof whose session answered passes when only its close fails, and the session is recorded            | backend | high       | task_01 |
| task_03 | The next proof closes a recorded session, and the CLI reports the close truthfully                      | backend | high       | task_02 |
| task_04 | Run the final QA gate                                                                                   | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 writes the guides, the Roundfix skill references and the glossary
first, so each code Task names its guide through that dependency; it records
the owned-skill version and the Behavior Surface fingerprint last. task_03
follows task_02 because both edit `internal/agent/selection_assignment.go`,
`internal/agent/fallback.go` and `internal/agent/acpx_runner.go`, its
reaper reads the records task_02 writes, and its report reads the
`SelectionProof.Unclosed` field task_02 adds. Each code Task runs the proof,
probe, analyzer and CLI suites that consume its contract, so no separate
Integration Checkpoint is needed. No other Spec is required.
