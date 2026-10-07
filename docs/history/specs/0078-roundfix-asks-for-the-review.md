---
schema: roundfix/archive-record/v1
spec: 0078-roundfix-asks-for-the-review
title: Roundfix asks for the review
status: archived
created: "2026-08-05"
archived: "2026-08-05"
disposition: pass
source: docs/history/specs/0078-roundfix-asks-for-the-review
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_06
qa_report: qa-report-2026-08-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0036
sources: []
regeneration: []
promoted: []
---

# Roundfix asks for the review

CodeRabbit's rate limit is counted per developer identity over a rolling seven days, and the ladder is progressive: at 60 or more reviews in seven days a Pro account drops from five reviews per hour to one. This repository's Supervisor opens many Pull Requests a day, and every automatic review it never wanted — an archive Pull Request carrying `+4/-1`, a rebase, a changelog fix — spends from the same bucket as the reviews that matter.
