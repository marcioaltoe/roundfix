---
status: deferred
created_at: 2026-08-06
updated_at: 2026-08-26
kind: rollup
members:
  - 2026-07-26-claude-adapter-configoptions-migration.md
  - 2026-07-27-claude-adapter-standardization.md
  - 2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md
  - 2026-08-04-a-doctor-next-action-that-does-not-reach-green.md
  - 2026-08-04-the-documented-sandbox-escape-does-not-exist-on-the-codex-adapter.md
  - 2026-08-05-agent-full-access-passes-config-validation-and-fails-every-task.md
---

# Agent selection and execution environments — preflight must prove the Task's reality (2026-08-06)

These findings converge on the same failure mode: configuration can validate
while the selected adapter, model controls, or sandbox policy fails when the
Agent begins real work. Readiness is useful only when it exercises the same
selection and access contract the Task Session will receive.

## Consolidated learning

- Adapter lineage, advertised model controls, and model-name parsing belong to
  one selection proof; package presence alone is not capability evidence.
- Profile updates must merge the named category without deleting unrelated
  selections.
- Access-policy names must map to supported adapter behavior. A documented
  escape hatch or `agent_full_access` value that the runtime cannot honor is a
  preflight defect, not a Task failure.
- Doctor next actions must repair the failing predicate and then re-check the
  same predicate.

## Live edge

Specs 0052 and 0056 absorbed the first adapter and profile defects. The rollup
remains `pending` around end-to-end execution-policy proof: accepted config,
runtime launch, and Task filesystem access must describe the same environment.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
