---
schema: roundfix/archive-record/v1
spec: 0019-run-outcome-notifications
title: Run Outcome Notifications
status: archived
created: "2026-07-07"
archived: "2026-07-07"
disposition: pass
source: docs/history/specs/0019-run-outcome-notifications
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

# Run Outcome Notifications

A Detached Run that dies overnight stays dead until someone happens to look: monitors tied to an agent session end with the session, and nothing else wakes the user. Field evidence from dogfooding: a Run failed at the commit boundary after verification had passed and sat dead for ten and a half hours before a human noticed. Run Outcome Notifications make the machine speak when a Run reaches a terminal outcome — a native desktop notification by default, and a user-configured command for any other channel — so an unattended Run's ending is an event someone hears about, not a state someone must poll for.
