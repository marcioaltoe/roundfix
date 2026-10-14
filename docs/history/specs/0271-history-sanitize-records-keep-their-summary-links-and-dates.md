---
schema: roundfix/archive-record/v1
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
title: History sanitize records keep their summary, links and dates
status: archived
created: "2026-10-13"
archived: "2026-10-14"
disposition: partial
source: docs/specs/0271-history-sanitize-records-keep-their-summary-links-and-dates
source_revision: c2dcb593a88b899a20d85feaf4789ddfc041ec8c
qa_task: task_05
qa_report: qa-report-2026-10-14-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0277
  - ADR-0260
sources:
  - 2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "542"
---

# History sanitize records keep their summary, links and dates

Both conversion commands repair the six legacy defect shapes and write identical records and promoted copies. Unchanged shapes match the delivery-base binary, existing Archive Records and Archive Command tests remain unchanged, and all required regressions, glossary and scope checks pass. The report closes as a qualifying partial because original outside publications cannot be retrieved under the Spec’s network prohibition and no Pull Request is open.
