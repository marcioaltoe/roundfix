---
schema: roundfix/archive-record/v1
spec: 0082-the-manifest-already-answered-that
title: Baseline update
status: archived
created: "2026-08-07"
archived: "2026-08-09"
disposition: qa-override
source: docs/history/specs/0082-the-manifest-already-answered-that
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_08
qa_report: qa-report-2026-08-07-02.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0099
  - ADR-0058
sources: []
regeneration: []
promoted: []
---

# Baseline update

A repository that has already adopted the Context-Driven Baseline records every answer it gave in its Setup Manifest — the profile, the ten decisions, the managed artifacts and their digests, the verification gate. Refreshing that repository against a newer Roundfix binary should be a mechanical act. Today it is an interview: the interactive command reads the manifest, announces `update` mode, and then asks all twelve questions anyway, each one offering the stored value as a default that still has to be confirmed. It then spends a supervised ACP turn re-segmenting and re-classifying the entire root instruction corpus and asks for a rule-by-rule review of the result. The non-interactive path is worse: it never opens the manifest at all, refuses with a list of the same decisions the manifest holds, and — once those are supplied — refuses again for an instruction-preservation Decision Document that, in Preservation mode, must bind every Source Baseline Entry.
