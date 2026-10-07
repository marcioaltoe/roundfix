---
schema: roundfix/archive-record/v1
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
title: A setup and doctor that recognize what the machine and the repository need
status: archived
created: "2026-10-02"
archived: "2026-10-02"
disposition: qa-override
source: docs/history/specs/0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02-02.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked row is the Pull Request row (11), which needs an open Pull Request
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 97539097ce4062aaa7f71d778e3e088f14887aa9
adrs:
  - ADR-0220
  - ADR-0221
sources:
  - 2026-09-30-an-installed-skill-that-matches-its-lock-can-trail-the-snapshot.md
  - 2026-09-30-declare-this-repositorys-derived-paths.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "348"
delivery_commit: dd0e5de917072b4cf1173697ba0218527654db7c
---

# A setup and doctor that recognize what the machine and the repository need

On 2026-10-01, preparing a second machine, the maintainer asked:
