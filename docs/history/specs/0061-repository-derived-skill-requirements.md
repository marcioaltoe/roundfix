---
schema: roundfix/archive-record/v1
spec: 0061-repository-derived-skill-requirements
title: Repository-derived skill requirements
status: archived
created: "2026-07-29"
archived: "2026-07-29"
disposition: qa-override
source: docs/history/specs/0061-repository-derived-skill-requirements
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-29.md
qa_verdict: partial
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

# Repository-derived skill requirements

> **Archived with QA `partial`, by maintainer decision on 2026-07-29.** The > gate recorded **no Spec behavior finding**; the verdict is capped by one row > the managed sandbox cannot execute, because it denies `fork/exec /bin/ps` and > five process-identity tests live inside the repository gate. That row was > executed outside the sandbox and the derivation was exercised against five > real checkouts, in > [`qa/evidence/2026-07-29-supervisor-full-access-gate.md`](qa/evidence/2026-07-29-supervisor-full-access-gate.md). > Spec 0055 removes the `/bin/ps` dependency and Spec 0053 makes this verdict > reachable without an override.
