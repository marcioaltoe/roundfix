---
schema: roundfix/archive-record/v1
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
title: An archive that relinks and an update that lists citations
status: archived
created: "2026-10-12"
archived: "2026-10-12"
disposition: pass
source: docs/specs/0255-an-archive-that-relinks-and-an-update-that-lists-citations
source_revision: 88af2f93bf0fb4978ee76c7bcc9d4bf67fa216f5
qa_task: task_05
qa_report: qa-report-2026-10-12-02.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0260
sources:
  - 2026-10-12-baseline-update-sanitizes-history-without-reporting-the-citations-it-breaks.md
  - 2026-10-12-an-archive-breaks-links-other-specs-make-into-its-references.md
regeneration:
  - command: make skills-sync
promoted: []
---

# An archive that relinks and an update that lists citations

The Archive Record will carry a passing QA verdict: the owner archive preserves resolving Live Document citations through equal-byte promotion and bounded relinks, the queue proves the exact archive, and baseline update lists/binds History Citations, promotes only on request and reports stale targets without changing its outcome. The actual read-only `./bin/roundfix archive 0255-an-archive-that-relinks-and-an-update-that-lists-citations --plan` exited 0 under a private home with no judge credential. It listed the two adopted report candidates with no advice; no candidate was labeled reusable_knowledge and no upstream promotion was indicated. The existing characterization already lives in docs/references. This report makes no claim of completed provider review, open PR or merge acceptance.
