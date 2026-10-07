---
schema: roundfix/archive-record/v1
spec: 0084-an-update-that-can-run
title: A Baseline Command update that runs on the repositories that already exist
status: archived
created: "2026-08-08"
archived: "2026-08-09"
disposition: qa-override
source: docs/history/specs/0084-an-update-that-can-run
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_10
qa_report: qa-report-2026-08-08.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0101
  - ADR-0102
  - ADR-0103
  - ADR-0104
sources:
  - 2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md
regeneration: []
promoted: []
---

# A Baseline Command update that runs on the repositories that already exist

Spec 0082 shipped `roundfix baseline update` so a maintainer could refresh a repository's Context-Driven Baseline without re-answering every setup question. Measured against the eight repositories that have actually adopted a Baseline, it refuses to run on six of them. Two of those refusals are this Spec's subject: the command reads a legitimately current managed region as damage and stops before it plans anything. Three more stop on structural clauses the catalog silently stopped emitting. One stops on a Baseline Profile its Setup Manifest names and its checkout no longer contains.
