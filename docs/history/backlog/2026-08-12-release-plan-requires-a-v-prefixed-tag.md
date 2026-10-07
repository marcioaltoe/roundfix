---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# `release plan` requires a `v`-prefixed tag and cannot see a repository without one

In `fluxus`, all eight release tags carry no prefix — `0.7.2`, `0.7.1`, `0.7.0`, `0.6.1`, `0.6.0`, `0.5.7`, `0.5.6`, `0.5.5`. `roundfix release plan` finds none:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-12-release-plan-requires-a-v-prefixed-tag.md`.
