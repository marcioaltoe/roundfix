---
schema: roundfix/archive-record/v1
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
title: A delivery queue that runs the binary its item builds
status: archived
created: "2026-10-03"
archived: "2026-10-03"
disposition: qa-override
source: docs/history/specs/0221-a-delivery-queue-that-runs-the-binary-its-item-builds
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-03.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; blocked rows need the operator's log outside the repository (09), outbound network the QA sandbox denies (11, 12) and an open Pull Request (15)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: a9a317dcc0ad48ed17956dbd0e4423a99916141c
adrs:
  - ADR-0225
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "365"
delivery_commit: 9c98d516fcb1a443fcb67fae6bbd61a13b75d2d3
---

# A delivery queue that runs the binary its item builds

On 2026-10-02 the Delivery Queue parked Spec 0215 `delivery-error` for a reason the item could not fix: the queue owner, built from `main`, started `roundfix implement` in the item worktree, that child read the item's Project Config, and the item had just added the key `verification.tools` to it. The child refused with `verification.tools is not a supported config key`. The operator built `bin/roundfix` from the item branch and retried with it, and did so again on the next two retries before the item merged. Any Spec that adds a Project Config key to Roundfix's own repository meets the same park at its first retry, archive or review.
