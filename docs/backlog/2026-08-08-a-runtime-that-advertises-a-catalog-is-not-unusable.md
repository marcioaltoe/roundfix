---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A runtime that advertises a whole catalog is not unusable

## Opportunity

`internal/agent/selection_capabilities.go` caps an advertised capability at 64
values and rejects the capability outright above it. OpenCode advertises every
model it knows rather than the subset a subscription grants, which is hundreds,
so the cap refuses the capability, no Agent Selection can be proven, and no Run
can start on that runtime. CONTEXT.md lists OpenCode as a supported ACP Runtime
through `opencode acp`, so the product promises a route the cap closes.

## Value

The cap exists to refuse an implausible capability payload, which is right in
intent. It reads a large catalog as malformed input when a large catalog is
exactly what a model aggregator legitimately advertises.

The cost became concrete on 2026-08-08: the Codex quota was exhausted for four
days and the Anthropic weekly limit stood at 66% with two days to reset, leaving
a third runtime as the obvious relief and unreachable. A subscription that
includes `gpt-5.6-luna` — the model this repository's `review` profile already
chose on measured cost — cannot be routed to the review category while the cap
holds.

## Shape

A future design could bound what Roundfix keeps rather than what a runtime may
advertise: accept a large advertised set, retain only what a configured or
requested Agent Selection references, and keep refusing payloads that are
malformed rather than merely large. Exact Agent Selection Proof already gates
correctness downstream, so the cap is not the only guard.

Worth settling in the same work: whether the Model Catalog for such a runtime
should be repository-configured instead of adapter-enumerated, since a
maintainer selecting from hundreds of interactive options is not a usable flow
either. This shape is non-binding.
