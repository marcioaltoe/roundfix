---
schema: roundfix/archive-record/v1
spec: 0236-baseline-update-with-a-repository-profile
title: Baseline update with a repository profile
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: pass
source: docs/history/specs/0236-baseline-update-with-a-repository-profile
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_03
qa_report: qa-report-2026-10-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0241
sources:
  - 2026-10-06-baseline-update-refuses-a-repository-profile.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "424"
delivery_commit: 31d19ad858ab968229271bcfc4c46109e91078d4
---

# Baseline update with a repository profile

On Roundfix 0.43.0 `roundfix baseline update` exits 1 in every repository adopted with its own Baseline Profile, after it has already computed a valid plan:
