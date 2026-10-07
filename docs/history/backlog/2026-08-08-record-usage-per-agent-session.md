---
type: feat # feat | fix | perf | refactor
status: deferred
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Record what each Agent Session consumed, next to which selection ran it

Roundfix persists the effective Agent Selection for every Task — ACP Runtime, Agent Model, and reasoning effort — and records nothing about what that Session consumed. Run Events cover `task-status`, `agent-selection`, `verification`, and `outcome`; no production path carries a token count. This intent came through `inbox/roundfix/2026-08-08-uso-por-task-com-agent-modelo-e-reasoning.md` in the Secondbrain.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-08-record-usage-per-agent-session.md`.
