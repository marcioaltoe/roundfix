---
schema: roundfix/archive-record/v1
spec: 0242-an-archive-that-leaves-an-archive-record
title: An archive that leaves an Archive Record
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: qa-override
source: docs/specs/0242-an-archive-that-leaves-an-archive-record
source_revision: 3afb60ad06ae829ed51380a0a42265e04698fffa
qa_task: task_05
qa_report: qa-report-2026-10-06-02.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after corrective tasks 06-07 and the Transcript 5 fix: row 07 needs an injectable fake judge transport the built CLI does not expose, row 12 needs the original authoring-ablation logs, and row 18 has no Pull Request yet. Every behavior row passed, including both folder-free gates.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 3afb60ad06ae829ed51380a0a42265e04698fffa
adrs:
  - ADR-0247
sources:
  - 2026-10-06-history-keeps-only-what-the-secondbrain-needs.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
---

# An archive that leaves an Archive Record

Archives write compact Archive Records, remove committed Spec folders, preserve recoverable Git bytes and copy promoted files with the documented confirmation. Record-only readers, squash-merge evidence and conservative missing-revision behavior pass; the guidance and generated artifacts agree, and both complete gates pass after every archived Spec folder is removed in a disposable clone. QA remains partial because successful public fake-judge advice and the original authoring-ablation evidence remain unverified. The ordinary pre-PR Pull Request absence has its required equivalent controls recorded.
