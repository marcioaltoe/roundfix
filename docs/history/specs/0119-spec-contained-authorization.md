---
schema: roundfix/archive-record/v1
spec: 0119-spec-contained-authorization
title: A Spec carries its authority
status: archived
created: "2026-09-08"
archived: "2026-09-13"
disposition: pass
source: docs/history/specs/0119-spec-contained-authorization
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_10
qa_report: qa-report-2026-09-10-04.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0130
sources:
  - 2026-08-26-the-baseline-does-not-name-a-home-for-authorization-records.md
  - 2026-09-08-verification-command-execution-has-an-explicit-trust-boundary.md
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A Spec carries its authority

The maintainer wants every pending Spec to carry the authorizations needed to implement and deliver it. Today grants may live outside the Spec, historical discovery assumes one workflow directory, and executing authored Verification does not have a clearly stated source-trust boundary. A maintainer should be able to inspect one Spec and distinguish approved actions from proposed work.
