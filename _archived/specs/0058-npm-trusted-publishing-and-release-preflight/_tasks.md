---
schema: spec-tasks/v1
spec: 0058-npm-trusted-publishing-and-release-preflight
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
      needs: [task_06]
    - id: task_08
      file: task_08.md
      needs: [task_07]
---

# Tasks — npm Trusted Publishing and release preflight

| id      | title                                          | type  | complexity | needs   |
| ------- | ---------------------------------------------- | ----- | ---------- | ------- |
| task_01 | Raise the publishing runtime and grant OIDC     | infra | low        | —       |
| task_02 | Preflight the release set against the registry  | infra | high       | task_01 |
| task_03 | Expose a publish-free preflight rehearsal       | infra | medium     | task_02 |
| task_04 | Publish through OIDC with a bounded fallback    | infra | high       | task_03 |
| task_05 | Document the migration, window, and vocabulary  | docs  | medium     | task_04 |
| task_06 | Attribute a publish failure to its actual cause | infra | medium     | task_05 |
| task_07 | Document the registry-side token shutdown       | docs  | low        | task_06 |
| task_08 | Confine the retained token without interpolating it | infra | medium | task_07 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 ·
6 → task_06 · 7 → task_07 · 8 → task_08

The graph is a chain rather than a fan-out because tasks 01–04, 06 all mutate
the single authorized file `.github/workflows/release.yml`. Parallel waves
would put two Agent Sessions in the same file.

Tasks 06 and 07 are remediation slices added after the QA gate returned
`fail` on 2026-07-31. They close QA-002 (every publish failure was attributed
to identity and token-retried) and QA-003 (the runbook omitted the
registry-side token shutdown). QA-001 was closed by amending the PRD, since
npm offers no read-only way to verify the promise it made.
