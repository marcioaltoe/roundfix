---
schema: roundfix/archive-record/v1
spec: 0130-documentation-cleanup-compatibility
title: Verification survives the documentation cleanup
status: archived
created: "2026-09-09"
archived: "2026-09-09"
disposition: pass
source: docs/history/specs/0130-documentation-cleanup-compatibility
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_02
qa_report: qa-report-2026-09-09.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "180"
delivery_commit: 1106d154f149bc7c5a8109d9a470f5d98700a599
---

# Verification survives the documentation cleanup

PR #180 removes obsolete workflow documents as requested, but test fixtures and the regeneration reader still consume their paths. Restore those contracts before integrating the cleanup, preserving real authorization evidence and strict output ownership. This is the five-file repair approved by the maintainer, not the full implementation of Spec 0120.
