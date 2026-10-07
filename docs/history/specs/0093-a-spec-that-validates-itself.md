---
schema: roundfix/archive-record/v1
spec: 0093-a-spec-that-validates-itself
title: A Spec that validates itself
status: archived
created: "2026-08-09"
archived: "2026-08-09"
disposition: qa-override
source: docs/history/specs/0093-a-spec-that-validates-itself
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-09-02.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs: []
sources: []
regeneration: []
promoted: []
---

# A Spec that validates itself

Spec 0090's QA gate returned `fail` on two defects. One of them was a PRD claiming that ADR-0083 makes `make verify` the authoritative gate. ADR-0083 is "Adopted sources move to their owning Spec" and says nothing about verification; the rule exists, but in `docs/agents/specific-repository.md`. The Spec invented a provenance, and the invention travelled into its TechSpec and the queue document.
