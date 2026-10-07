---
schema: roundfix/archive-record/v1
spec: 0196-a-notice-when-profiles-fall-behind
title: A notice when profiles fall behind the recommendation
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0196-a-notice-when-profiles-fall-behind
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0181
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "306"
delivery_commit: 0ce4800c991946aca80771509bf43343e7b48b38
---

# A notice when profiles fall behind the recommendation

Spec 0189 gives Roundfix one dated Recommended Profile per Agent Work Category, and its built-in profiles follow that snapshot with each release. A configured profile does not follow anything:
