---
schema: spec-tasks/v1
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — The Catalog Skill Bump refreshes with a binary built after the digests

| id      | title                                                                    | type  | complexity | needs   |
| ------- | ------------------------------------------------------------------------ | ----- | ---------- | ------- |
| task_01 | The Catalog Skill Bump refreshes with a binary built after the digests   | infra | low        | —       |
| task_02 | Run the final QA gate                                                    | qa    | low        | task_01 |

Waves: 1 → task_01 · 2 → task_02

task_01 runs, in its own Verification, the Catalog Skill Bump Makefile
contract that consumes the recipe and the Setup Manifest currency test that
reads the regenerated manifest, so no separate Integration Checkpoint is
needed. No other Spec is required; the queue delivers this Spec after 0373,
0374, 0375, 0377 and 0376, and task_01 regenerates the shared Setup Manifest
from whatever the merge base holds.
