---
schema: roundfix/archive-record/v1
spec: 0220-tests-and-pins-that-hold-in-every-environment
title: Tests and pins that hold in every environment
status: archived
created: "2026-10-03"
archived: "2026-10-03"
disposition: qa-override
source: docs/history/specs/0220-tests-and-pins-that-hold-in-every-environment
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_03
qa_report: qa-report-2026-10-03.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked row is the Pull Request row (9)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 269a16f911d55b836314632d8806405ce4a944c4
adrs:
  - ADR-0224
sources:
  - 2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md
  - 2026-10-02-restoring-trailing-skills-breaks-the-pinned-skill-digests.md
regeneration: []
promoted: []
pull_request: "367"
delivery_commit: d027c3d90457e046a66ea1777a58d3b6bbd0377d
---

# Tests and pins that hold in every environment

Two checks in `make verify` hold only in the environment they were written in. This Spec adopts both Backlog Entries that name them:
