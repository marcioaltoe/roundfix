---
type: feat
status: declined
created: 2026-10-03
spec: null
reason: maintainer decided 2026-10-04 to keep the cursor runtime opt-in only; no built-in profile changes
---

# A Grok selection through Cursor as a `docs` and `chore` fallback

The `docs` and `chore` profiles fall back from `codex / gpt-5.6-luna / max` to `claude / sonnet / high`. Spec 0217 made `cursor` an ACP Runtime and measured `cursor / grok-4.7[context=256k,reasoning_effort=high,fast=true] / ""` on two replayed Tasks, one `docs` and one `chore`: both Runs were Clean with one prompt, no Verification repair and no person needed, the same outcome as the default (`docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md`). A Grok fallback would give these categories a third provider, billed to the Cursor plan rather than to the Codex or Claude quota.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-10-03-a-grok-fallback-for-docs-and-chore.md`.
