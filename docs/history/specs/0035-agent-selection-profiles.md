---
schema: roundfix/archive-record/v1
spec: 0035-agent-selection-profiles
title: Agent Selection Profiles
status: archived
created: "2026-07-16"
archived: "2026-07-17"
disposition: pass
source: docs/history/specs/0035-agent-selection-profiles
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-17.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "32"
delivery_commit: d7fdd8cf3f986bda0cbd3683e239c799dcf6b41a
---

# Agent Selection Profiles

Roundfix currently selects one ACP Runtime, Agent Model, and reasoning effort for an entire Run. Its fallback is discovered dynamically from one runtime's Model Catalog and requires confirmation after failure. That contract cannot guarantee that backend, frontend, QA, and review work use the user's intended models, cannot safely route mixed Task Graphs, and leaves the chosen fallback unknown until the Run has already been delayed. Agent Selection Profiles make the complete preferred-and-fallback policy explicit, proven before a Run, observable per Work Item, and configurable from the Roundfix CLI.
