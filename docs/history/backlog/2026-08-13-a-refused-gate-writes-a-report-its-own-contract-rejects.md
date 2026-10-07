---
type: fix # feat | fix | perf | refactor
status: done
created: 2026-08-13
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A refused gate writes a report its own contract rejects

When the QA gate refuses at its authoring precondition, it writes a report whose Results table is empty — correctly, because it stopped before building the matrix. The mechanical stage of every later run then reads that report and refuses:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-13-a-refused-gate-writes-a-report-its-own-contract-rejects.md`.
