---
schema: spec-tasks/v1
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
qa: task_03
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
---

# Tasks — The Makefile follows the Go build cache set with go env

| id      | title                                                                                  | type  | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------- | ----- | ---------- | ------- |
| task_01 | The Makefile uses the Go build cache set with go env and falls back small with a warning | infra | medium     | —       |
| task_02 | The guide, README, glossary and ADR-0370 describe the cache set with go env             | docs  | low        | task_01 |
| task_03 | Run the final QA gate                                                                  | qa    | low        | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 holds the `Makefile` change, every test that runs the `Makefile` and
the setup manifest regeneration together, because the tests read the
`Makefile` and the manifest records its final bytes. task_02 documents the
behavior task_01 delivers, so it runs after it. task_01's Verification runs
the whole `./internal/verifyselect` package, which holds every test that
consumes the `Makefile`, so no separate Integration Checkpoint is needed. No
other Spec is required; Spec 0383 delivers first and may change `CONTEXT.md`
and the setup manifest in other places (`_prd.md` → Prerequisites).
