---
schema: roundfix/archive-record/v1
spec: 0113-a-gate-report-that-does-not-block-its-successor
title: A gate report that does not block its successor
status: archived
created: "2026-08-14"
archived: "2026-08-15"
disposition: pass
source: docs/history/specs/0113-a-gate-report-that-does-not-block-its-successor
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-15-03.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "165"
delivery_commit: 5d764b9af2443e5979737350d5e3735b39b0fdac
---

# A gate report that does not block its successor

The QA gate refuses before spending an Agent turn when a machine fact says it should, which is right and cheap. What it then writes is a report its own contract calls malformed, and every later run of the same Spec reads that report and refuses on it — with a prescribed fix that is impossible, because the run never built a matrix to materialize rows from. It happened twice in one Spec, and the only exit was deleting evidence, which is the one move this repository's rules single out. A second, independent refusal in the same family names a row's blocked cause as untyped when it is typed correctly: what the parser wants is a literal the diagnostic never mentions, and the mismatch produces a second symptom that sends a reader hunting a counting bug that does not exist.
