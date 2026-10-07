---
schema: roundfix/archive-record/v1
spec: 0183-run-storage-that-says-what-it-holds
title: Run storage that says what it holds
status: archived
created: "2026-09-29"
archived: "2026-09-29"
disposition: pass
source: docs/history/specs/0183-run-storage-that-says-what-it-holds
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-29-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0171
  - ADR-0172
sources:
  - 2026-09-25-storage-reclaim-notice.md
  - 2026-09-25-bare-layout-artifact-roots-over-preserved.md
  - 2026-09-29-runs-list-warns-for-runs-whose-checkout-is-gone.md
  - 2026-09-29-gc-reports-already-emptied-runs-as-pruned.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "282"
delivery_commit: f3cbf0c46e8e52384c421b2a7bb5bef073c700c6
---

# Run storage that says what it holds

Roundfix keeps every Run's row, Run Event Journal rows and artifact directory in the machine-wide Run Database and the Artifact Root. Four reports about that storage are wrong or missing, and each misleads the operator who reads them:
