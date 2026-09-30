---
type: fix
status: open
created: 2026-09-30
spec: null
reason: null
---

# A dispatch trigger names a skill the model cannot invoke

## Symptom

The `rust` Baseline module dispatches the `cut-release` skill with the trigger `trigger.rust.cut-release`, "Preparing or publishing a Rust CLI release." Upstream, `cut-release` is a user-invocable skill only: a person starts it, and the model cannot load it on its own. An Agent that reads the trigger in `docs/agents/skill-dispatch.md` is told to activate a skill it has no way to activate, so the trigger never fires.

## Where

`internal/baseline/assets/modules/rust.json` (`skillDispatch`, `requiredSkills`) and the upstream skill `skills/08-release/cut-release`.

## Expected

A dispatch trigger names only a skill the model can load. For a user-invocable skill, the guide tells the Agent to ask the person to run it, or the module stops dispatching it.

## Evidence

Reported by the stack audit of 2026-09-30. Spec 0195 does not change it: membership and dispatch of upstream skills belong to the Spec that refreshes the upstream skill snapshot, and this case needs a rule for user-invocable skills that no Spec has yet.
