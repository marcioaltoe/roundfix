---
type: feat # feat | fix | perf | refactor
status: deferred
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Speak the Agent Client Protocol directly instead of through acpx

## Opportunity

Roundfix could implement the Agent Client Protocol itself —
`github.com/agentclientprotocol/agent-client-protocol` — rather than driving
agents through the `acpx` CLI. The maintainer raised this as a possible future
direction on 2026-08-07, and this intent came through
`inbox/roundfix/2026-08-07-speak-acp-directly-instead-of-through-acpx.md` in the
Secondbrain.

## Value

Three observations landed the same day while probing adapters for model
identifiers, and all three trace to acpx sitting between Roundfix and the
protocol.

Model identifiers had to be discovered by *provoking an error*. There is no "list
models" surface, so the codex list came from an adapter refusal and the claude and
opencode lists from parsing `configOptions` out of a raw `session/new` result.
Speaking ACP directly would make advertised models, effort values and modes
ordinary reads.

Claude Selections could not be proven, because the adapter answers an unknown
model with a `Custom model` entry rather than a refusal. Owning the client side
would let Roundfix compare a requested identifier against the advertised set
itself instead of waiting for a refusal that never arrives. Spec 0091 has since
answered that specific case by making the refusal that fires first the refusal
(ADR-0119), which lowers the urgency without removing the shape.

`internal/agent/catalog.go` hand-maintained model lists that had already drifted
from what adapters advertise. Separately, the sealed-prompt path was observed
outliving its own five-minute timeout on 2026-08-07, which is a subprocess
lifetime problem a direct protocol client would own rather than inherit.

## Shape

Non-binding, and deliberately unscoped. Replacing a working dependency is a large
change with its own risks: acpx already handles session lifecycle, permission
policies, terminal capabilities, and more than twenty agent runtimes, and
reimplementing that surface would be substantial. The value of this entry is the
list of concrete costs the current indirection imposes, so a future decision
starts from evidence rather than preference. Routing it through `write-idea`
before any PRD is the smaller sufficient route.

Related: `docs/references/model-selection.md` records how each adapter's
identifiers had to be obtained.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
