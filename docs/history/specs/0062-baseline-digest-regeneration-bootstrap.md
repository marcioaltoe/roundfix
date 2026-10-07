---
schema: roundfix/archive-record/v1
spec: 0062-baseline-digest-regeneration-bootstrap
title: Baseline digest regeneration bootstrap
status: archived
created: "2026-08-01"
archived: "2026-08-01"
disposition: pass
source: docs/history/specs/0062-baseline-digest-regeneration-bootstrap
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-08-01-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0085
sources: []
regeneration: []
promoted: []
pull_request: "63"
delivery_commit: 786b8b382b7012135e88fa08547ae9a7a4453382
---

# Baseline digest regeneration bootstrap

`make baseline-digests` cannot refresh the derived pins it exists to refresh. Every step loads the embedded catalog; catalog validation rejects a formatter `goldenDigest` that disagrees with the goldens on disk; only those steps rewrite that pin. An ordinary Baseline module edit changes the generated guide, which changes the goldens, which invalidates the pin, which refuses the load, which prevents the refresh — and the command's printed remediation is the command itself. A second defect compounds it: adding a Normative Clause reports `catalog.sourceBaseline.required-clause.missing` for every new clause, because the regenerator maintains digests and spans for manifest rows that already exist but never creates one, and the diagnostic never says the tool cannot supply it. Both surface on any module edit that changes a generated guide, which is the ordinary case, so they block the autonomous Spec loop rather than being incidental to one change. Evidence: [`make baseline-digests` cannot bootstrap a Baseline module edit](../../findings/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md).
