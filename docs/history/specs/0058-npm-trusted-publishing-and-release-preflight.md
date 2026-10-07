---
schema: roundfix/archive-record/v1
spec: 0058-npm-trusted-publishing-and-release-preflight
title: npm Trusted Publishing and release preflight
status: archived
created: "2026-07-28"
archived: "2026-08-01"
disposition: qa-override
source: docs/history/specs/0058-npm-trusted-publishing-and-release-preflight
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-08-01-04.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0082
  - ADR-0048
  - ADR-0084
sources: []
regeneration: []
promoted: []
pull_request: "65"
delivery_commit: 3ad0a452f797d754105c85430c47887c3285ae1e
---

# npm Trusted Publishing and release preflight

The `v0.0.1` release reset exposed two independent publication risks that remain open. Roundfix's release workflow authenticates six sequential `npm publish` commands with a long-lived repository secret, and it starts publishing platform packages before proving that every package in the release set can accept the target version — during the reset, the launcher name sat inside npm's 24-hour post-unpublish block while the five platform packages were publishable, a combination that without cancellation would have produced a partial release: installable platform packages with no launcher. Evidence: [token authentication and registry state can produce a partial release](../../findings/2026-07-25-npm-trusted-publishing-and-release-preflight.md).
