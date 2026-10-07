---
schema: roundfix/archive-record/v1
spec: 0017-run-discovery
title: Run discovery
status: archived
created: "2026-07-06"
archived: "2026-07-07"
disposition: pass
source: docs/history/specs/0017-run-discovery
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-07.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
---

# Run discovery

Detached Runs and Attach solved Run ownership, but discovering what to attach to still depends on the caller having captured the run id from the detach report. A user who lost the id — or an agent session that started fresh — has no command that answers "which Runs exist for this repository, and which are Active?". Run discovery closes that gap with a deterministic listing command and an Attach picker over the same data, so any session can find and follow a Run using only the CLI.
