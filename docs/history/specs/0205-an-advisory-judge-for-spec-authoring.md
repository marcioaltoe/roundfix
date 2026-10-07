---
schema: roundfix/archive-record/v1
spec: 0205-an-advisory-judge-for-spec-authoring
title: An advisory judge for Spec authoring
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0205-an-advisory-judge-for-spec-authoring
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0200
  - ADR-0201
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "329"
delivery_commit: 14e049b3c17b3c050db7748a01ec0bad5ecd8ed0
---

# An advisory judge for Spec authoring

The Spec Consistency Check proves that a cited ADR is listed, that a claim shares words with the record it names, and that every PRD goal has a Coverage Map line. It cannot tell whether the record actually says what the Spec attributes to it, or whether the TechSpec section a Coverage Map line names describes a way to reach the goal. Both defects reach review, where a reviewer settles them only by opening the cited text.
