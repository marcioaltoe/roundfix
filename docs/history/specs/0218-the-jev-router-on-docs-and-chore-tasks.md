---
schema: roundfix/archive-record/v1
spec: 0218-the-jev-router-on-docs-and-chore-tasks
title: The Jev Router on docs and chore Tasks
status: archived
created: "2026-10-01"
archived: "2026-10-02"
disposition: qa-override
source: docs/history/specs/0218-the-jev-router-on-docs-and-chore-tasks
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-02-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; blocked rows need a /home/maintainer transcript fixture (2), outbound network the QA sandbox prohibits (6, 7) and an open Pull Request (11)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 8a26592da6b48560aadc84c29bcea8c90b18ccb1
adrs:
  - ADR-0218
sources:
  - 2026-10-01-try-the-jev-router-for-cheap-agent-categories.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "355"
delivery_commit: 05e9c229dcc12b33e275c4bf456814f6469a04e8
---

# The Jev Router on docs and chore Tasks

Roundfix runs every Task of a work category on the same model and reasoning effort, even where most of the work is routine, as in the `docs` and `chore` categories. OpenRouter publishes `typesafe/jev-router`, which uses Jev to pick a model and an effort for each request. If it keeps quality on routine work while it picks cheaper models or lower effort, those categories cost less per Task.
