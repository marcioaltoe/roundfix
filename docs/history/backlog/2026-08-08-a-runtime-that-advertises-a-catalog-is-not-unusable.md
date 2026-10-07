---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A runtime that advertises a whole catalog is not unusable

`internal/agent/selection_capabilities.go` caps an advertised capability at 64 values and rejects the capability outright above it. OpenCode advertises every model it knows rather than the subset a subscription grants, which is hundreds, so the cap refuses the capability, no Agent Selection can be proven, and no Run can start on that runtime. CONTEXT.md lists OpenCode as a supported ACP Runtime through `opencode acp`, so the product promises a route the cap closes.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-08-a-runtime-that-advertises-a-catalog-is-not-unusable.md`.
