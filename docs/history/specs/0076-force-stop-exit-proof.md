---
schema: roundfix/archive-record/v1
spec: 0076-force-stop-exit-proof
title: Force Stop exit proof
status: archived
created: "2026-08-04"
archived: "2026-08-04"
disposition: pass
source: docs/history/specs/0076-force-stop-exit-proof
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_03
qa_report: qa-report-2026-08-04-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "106"
delivery_commit: 121de4ed242a281fda728cd1371731c34f50b868
---

# Force Stop exit proof

Force Stop is one of the strongest guarantees in the glossary: it "proves owner identity, cancels registered Agent Sessions, terminates the recorded owning process, and completes the Run as Stopped **only after owner exit is proven**." The test that proves the hardest half of that — that a process ignoring `SIGTERM` is still terminated — has never proved it.
