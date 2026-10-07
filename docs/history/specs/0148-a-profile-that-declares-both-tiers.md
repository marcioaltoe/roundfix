---
schema: roundfix/archive-record/v1
spec: 0148-a-profile-that-declares-both-tiers
title: A profile that declares both tiers
status: archived
created: "2026-09-19"
archived: "2026-09-19"
disposition: pass
source: docs/history/specs/0148-a-profile-that-declares-both-tiers
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-19.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
promoted: []
delivery_commit: 1347d7926a17d7ce6cf56f4cce9dfa3698bedd93
---

# A profile that declares both tiers

Two mandatory clauses in the shipped guidance tell an Agent the same thing: use the active Baseline Profile's **declared incremental verification command** for fast local checks, and keep the complete command for the assembled tree. One of them goes further and says that a missing incremental command leaves the Profile's two-tier contract unmet, and never authorizes skipping the local tier.
