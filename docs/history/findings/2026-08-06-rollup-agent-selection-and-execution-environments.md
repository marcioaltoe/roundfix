---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
absorbed_by: 0123-runtime-readiness-and-model-capabilities
members:
  - 2026-08-07-a-sixty-four-value-bound-locks-out-the-opencode-runtime.md
  - 2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md
  - 2026-08-14-preflight-starves-when-the-machine-is-busy.md
  - 2026-07-26-claude-adapter-configoptions-migration.md
  - 2026-07-27-claude-adapter-standardization.md
  - 2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md
  - 2026-08-04-a-doctor-next-action-that-does-not-reach-green.md
  - 2026-08-04-the-documented-sandbox-escape-does-not-exist-on-the-codex-adapter.md
  - 2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task.md
---

# Agent selection and execution environments — preflight must prove the Task's reality (2026-08-06)

These findings converge on the same failure mode: configuration can validate while the selected adapter, model controls, or sandbox policy fails when the Agent begins real work. Readiness is useful only when it exercises the same selection and access contract the Task Session will receive.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-rollup-agent-selection-and-execution-environments.md`.
