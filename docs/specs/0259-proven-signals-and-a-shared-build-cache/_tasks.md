---
schema: spec-tasks/v1
spec: 0259-proven-signals-and-a-shared-build-cache
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — Proven signals and a shared build cache

| id      | title                                                              | type    | complexity | needs            |
| ------- | ------------------------------------------------------------------ | ------- | ---------- | ---------------- |
| task_01 | Every signal goes to a proven process, and a rule keeps it so      | backend | high       | —                |
| task_02 | No test helper outlives its test binary                            | test    | high       | task_01          |
| task_03 | make shares one Go build cache that the sandbox can write          | infra   | medium     | task_02          |
| task_04 | doctor reports and gc reclaims Worktree Build Caches               | backend | high       | task_02          |
| task_05 | Run the final QA gate                                              | qa      | high       | task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05

task_02 changes the helpers and teardowns task_01 converts and extends task_01's
rule allowlist, and both edit `CONTEXT.md`, so they run in series. task_03 shares
only `docs/agents/specific-repository.md` with task_02, and task_04 shares only
`CONTEXT.md` with it; task_03 and task_04 share no declared file and run side by
side. task_03 is the Makefile change and lands as its own commit. The gate
follows both leaves, and its ten full suites run on the tree every Task built.
