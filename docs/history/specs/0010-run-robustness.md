---
schema: roundfix/archive-record/v1
spec: 0010-run-robustness
title: Run Robustness
status: archived
created: "2026-07-06"
archived: ""
disposition: no-qa
source: docs/history/specs/0010-run-robustness
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: ""
qa_verdict: ""
unproven: []
adrs:
  - ADR-0027
  - ADR-0028
sources: []
regeneration: []
promoted: []
---

# Run Robustness

The 0009 cycle exposed three process-lifecycle weaknesses in one day. Removing a config key hard-failed every Run on a machine whose config still carried it — the QA "ran overnight" as a corpse. Externally killed Runs leaked forty orphaned adapter processes, because cancelling an Agent Session never reaps its OS process tree. And a Run's lifetime silently depends on its caller's: three Runs died mid-flight when the invoking session reaped its background tasks, leaving `nohup` gymnastics as the only defense. This Spec closes all three: config migrations that never break users, Runs that clean up their processes, and Runs that survive their caller.
