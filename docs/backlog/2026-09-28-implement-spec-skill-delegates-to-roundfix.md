---
type: feat
status: open
created: 2026-09-28
spec: null
reason: null
---

# The implement-spec entry point delegates to Roundfix

## Opportunity

Spec 0127 Core Feature 8: the owned `implement-spec` skill still carries its own loop guidance instead of delegating implementation and Verification to `roundfix implement`/`deliver`.

## Value

One implementation loop; the Supervisor never writes code or tests.

## Shape

Rewrite the skill to prepare the queue and hand off to Roundfix, without copying a window script.

Evidence: carried from Spec 0127 when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
