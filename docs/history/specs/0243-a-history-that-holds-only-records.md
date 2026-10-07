---
schema: roundfix/archive-record/v1
spec: 0243-a-history-that-holds-only-records
title: A history that holds only records
status: archived
created: "2026-10-07"
archived: "2026-10-07"
disposition: qa-override
source: docs/specs/0243-a-history-that-holds-only-records
source_revision: 4f885acbfe46d4b20c079cf86d1945f078c451af
qa_task: task_05
qa_report: qa-report-2026-10-07-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after the Transcript 1 fix: row R05 needs a fake judge certificate macOS TLS trust rejects (advice is proven in the injected transport harness), row R08 needs the original authoring-ablation artifact, and row R13 has no Pull Request yet. Every behavior row passed, including a full sanitize in a disposable clone.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 4f885acbfe46d4b20c079cf86d1945f078c451af
adrs:
  - ADR-0248
  - ADR-0247
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
---

# A history that holds only records

QA at audited head 36ad11279e7cb0cc9ba0162e34621f86b35fd000 closes partial with 9 passed and 3 environment-blocked rows. All four Surface Transcripts match; the checkout plan reports 227 units and preserves exact status; tagged fixture batches preserve promotions and Git recovery; all refusals preserve bytes. Full sanitize in a disposable no-local clone converts 227 units (223 Archive Records and 129 reduced entries), reducing history from 60,653,133 to 405,406 bytes; make verify and make verify-docs both exit 0, the no-qa record for 0003-dogfood-polish is captured, the final public plan reports zero pending, and the clone is removed.
