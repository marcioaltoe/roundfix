---
type: feat # feat | fix | perf | refactor
status: deferred
created: 2026-08-06
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Add an atomic helper for durable inbox capture

Fleet sessions could capture a contract-valid Inbox Entry through one observable action instead of manually creating, staging, reviewing, committing, and timing the entry. This intent came through `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md` in the Secondbrain.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-06-atomic-inbox-capture-helper.md`.
