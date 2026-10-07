---
schema: roundfix/archive-record/v1
spec: 0174-operator-surfaces-that-tell-the-truth
title: Operator surfaces that tell the truth
status: archived
created: "2026-09-28"
archived: "2026-09-28"
disposition: pass
source: docs/history/specs/0174-operator-surfaces-that-tell-the-truth
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-28-01.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-09-25-read-only-commands-and-schema-versions.md
  - 2026-09-25-help-token-accepted-anywhere.md
  - 2026-09-25-a-qa-report-with-an-empty-front-matter-settles.md
  - 2026-09-25-governed-path-misses-the-skill-ownership-file.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "261"
delivery_commit: cd3945c326bbb3dd35a30bb7c87dc107fd0b813d
---

# Operator surfaces that tell the truth

The commands and checks an operator reads must say what is true. Four defects make them say something else:
