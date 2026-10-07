---
schema: roundfix/archive-record/v1
spec: 0002-acpx-migration
title: ACPX Migration
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0002-acpx-migration
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0017
  - ADR-0018
  - ADR-0011
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# ACPX Migration

Roundfix reaches its Agents through a hand-rolled ACP client layer that spawns a fresh runtime process for every Batch and re-implements session plumbing, permission answering, and filesystem jailing by hand. Every Work Item pays a runtime cold start, a crashed runtime ends the whole Run as Failed, and the client layer is the largest and most fragile subsystem Roundfix maintains. Migrating the agent layer to acpx — a headless ACP client with persistent named sessions, crash resume, and a raw protocol stream — makes multi-item Runs warm and crash-resilient while shrinking Roundfix to process orchestration it can actually own.
