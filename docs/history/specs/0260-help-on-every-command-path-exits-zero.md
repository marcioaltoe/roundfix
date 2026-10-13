---
schema: roundfix/archive-record/v1
spec: 0260-help-on-every-command-path-exits-zero
title: Help on every command path exits zero
status: archived
created: "2026-10-12"
archived: "2026-10-13"
disposition: pass
source: docs/specs/0260-help-on-every-command-path-exits-zero
source_revision: 07a765e9edf88640147c3556f6199f72152ea3ee
qa_task: task_02
qa_report: qa-report-2026-10-13-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0256
sources:
  - 2026-10-09-qa-report-accept-help-exits-one.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "506"
---

# Help on every command path exits zero

The Archive Record may carry QA verdict pass at audited head `b2c03c24c00af69cb8d10692e9010b4e46917dba`: every built command path answers both help flags with stdout usage and exit 0, the guards detect regressions, and the Behavior Surface Record stays unchanged relative to the delivery base. The gate report records one pre-PR environment block with the specified equivalents. The Daemon owns Task settlement and the report commit.
