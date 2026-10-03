---
type: feat
status: open
created: 2026-10-03
spec: null
reason: null
---

# A Grok selection through Cursor as a `docs` and `chore` fallback

## Opportunity

The `docs` and `chore` profiles fall back from `codex / gpt-5.6-luna / max`
to `claude / sonnet / high`. Spec 0217 made `cursor` an ACP Runtime and
measured `cursor / grok-4.7[context=256k,reasoning_effort=high,fast=true] / ""`
on two replayed Tasks, one `docs` and one `chore`: both Runs were Clean with
one prompt, no Verification repair and no person needed, the same outcome as
the default
(`docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md`).
A Grok fallback would give these categories a third provider, billed to the
Cursor plan rather than to the Codex or Claude quota.

## Value

When the Codex quota runs out, `docs` and `chore` Tasks could settle on the
Cursor plan instead of drawing on the Claude quota that implementation work
also needs. The hypothesis is that bounded documentation and chore Tasks
settle on Grok as often as on the current fallback.

## Shape

Non-binding. Add the Grok selection to the `docs` and `chore` Fallback Chains
in `.roundfixrc.yml`, after or in place of `claude / sonnet / high`, proved
with `roundfix profiles configure`. Before adopting it, settle three points
the measurement left open: the plan charge of a `fast=true` value, which
Cursor does not report to Roundfix as tokens; the wall time, about 1.5 times
the default's on both Tasks; and a larger sample than two low-complexity
Tasks. Reconcile it with Spec 0218's Jev router on `docs` and `chore` Tasks.
