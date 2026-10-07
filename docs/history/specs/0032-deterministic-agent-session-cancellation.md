---
schema: roundfix/archive-record/v1
spec: 0032-deterministic-agent-session-cancellation
title: Deterministic Agent Session cancellation
status: archived
created: "2026-07-16"
archived: "2026-07-16"
disposition: pass
source: docs/history/specs/0032-deterministic-agent-session-cancellation
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-17.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "30"
delivery_commit: 7adcd4cf2757875c14cd6ed748916127c586454f
---

# Deterministic Agent Session cancellation

The Agent Session cancellation test failed once because its 20 ms Stop Request grace also bounded the fake cancel subprocess, allowing scheduler delay to kill that process before the helper recorded the invocation. The production default remains 10 seconds and the failure did not reproduce, but the test cannot currently prove cancel-before-close ordering without depending on wall-clock scheduling. The cancellation boundary needs deterministic lifecycle control so the test reports a real ordering regression instead of machine load.
