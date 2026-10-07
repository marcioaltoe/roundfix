---
schema: roundfix/archive-record/v1
spec: 0155-a-verify-that-runs-what-changed
title: A verify that runs what changed
status: archived
created: "2026-09-24"
archived: "2026-09-24"
disposition: pass
source: docs/history/specs/0155-a-verify-that-runs-what-changed
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-24-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "240"
delivery_commit: d0ab2637f0ca8240df7f9da06716ae59f368e4d8
---

# A verify that runs what changed

Every repository Verification runs every test. A change to the Daemon runs the Baseline's 81 seconds, the Baseline command tests inside `internal/cli` (35 seconds) and the skill tests; a change to a Baseline profile runs the whole core. Measured on 2026-09-24, the Baseline and skills layer is about 26 percent of test time (roughly 126 of 478 seconds), 25 percent of the Go code and 38 percent of the governed paths, and every one of those seconds is paid by every change, including the ones that cannot affect it.
