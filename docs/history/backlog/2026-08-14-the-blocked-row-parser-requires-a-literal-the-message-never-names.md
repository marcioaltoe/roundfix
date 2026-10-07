---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The blocked-row parser requires a literal the message never names

A QA gate cycle executed no journey. It stopped at the mechanical detector, which evaluated the *previous* report and refused it on shape, producing a new report whose entire content is that refusal:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-the-blocked-row-parser-requires-a-literal-the-message-never-names.md`.
