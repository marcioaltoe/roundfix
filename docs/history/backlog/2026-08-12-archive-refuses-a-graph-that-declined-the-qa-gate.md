---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# `archive` refuses a graph that declined the QA gate

Two parts of the contract disagree. The `write-tasks` skill declares `qa: declined` as one of the **two** valid shapes of a post-contract graph:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-12-archive-refuses-a-graph-that-declined-the-qa-gate.md`.
