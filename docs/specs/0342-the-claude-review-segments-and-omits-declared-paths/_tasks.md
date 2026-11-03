---
schema: spec-tasks/v1
spec: 0342-the-claude-review-segments-and-omits-declared-paths
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — The Claude review segments a large candidate and omits declared paths

| id      | title                                                                                                  | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | A Review Omit Path leaves a declared file out of the pre-PR review                                     | backend | medium     | —       |
| task_02 | A Claude round-1 candidate over its budget is reviewed in Review Segments                              | backend | high       | task_01 |
| task_03 | A size refusal names the largest files and the key, and the park carries it                            | backend | medium     | task_02 |
| task_04 | The guides, the Roundfix references and the glossary describe omit paths, Claude segments and the refusal | docs    | medium     | task_03 |
| task_05 | Run the final QA gate                                                                                  | qa      | medium     | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

task_02 needs task_01 because both edit `internal/cli/review.go`. task_03
needs task_02 because both edit `internal/cli/review.go` and
`internal/cli/review_segments.go`, and its refusals wrap the Claude segment
refusals task_02 adds. task_04 needs task_03 because it documents the
behavior task_01 to task_03 build. Each Task runs the suites that consume its
changed contract in its own Verification, so no separate Integration
Checkpoint is needed. No other Spec is required.
