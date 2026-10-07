---
schema: roundfix/archive-record/v1
spec: 0003-dogfood-polish
title: Dogfood Polish
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: no-qa
source: docs/history/specs/0003-dogfood-polish
source_revision: 36ad11279e7cb0cc9ba0162e34621f86b35fd000
qa_task: ""
qa_report: ""
qa_verdict: ""
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# Dogfood Polish

The first two dogfood cycles (Implement Command on spec 0002, watch on a real pull request) shipped working software and a list of small, confirmed irritations: contract inconsistencies, misleading header lines, an environment-sensitive test suite, and opaque infrastructure errors. Each item is small, already diagnosed in the dogfood findings log, and none needs new architecture — this Spec batches them so the debris is cleared before the larger watch and TUI cycles, and so the 0002 QA gate can pass on re-run.
