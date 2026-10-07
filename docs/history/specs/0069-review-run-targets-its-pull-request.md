---
schema: roundfix/archive-record/v1
spec: 0069-review-run-targets-its-pull-request
title: A Review Run targets its Pull Request
status: archived
created: "2026-08-02"
archived: "2026-08-05"
disposition: pass
source: docs/history/specs/0069-review-run-targets-its-pull-request
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_04
qa_report: qa-report-2026-08-05-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "128"
delivery_commit: 8b8bfe5281e5fc6eb8e59aeb55fec16df4720404
---

# A Review Run targets its Pull Request

`roundfix watch --source coderabbit --pr N` names the Pull Request it should process, but resolves the branch it acts on from the main checkout instead. The two are assumed to agree and nothing checks that they do, so a mismatch is silent until its consequences appear.
