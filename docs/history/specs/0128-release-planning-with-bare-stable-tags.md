---
schema: roundfix/archive-record/v1
spec: 0128-release-planning-with-bare-stable-tags
title: Release planning preserves stable tag identity
status: archived
created: "2026-09-08"
archived: "2026-09-20"
disposition: superseded
source: docs/history/specs/0128-release-planning-with-bare-stable-tags
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: ""
qa_report: ""
qa_verdict: ""
unproven: []
superseded_by: 0147-a-planner-that-reads-both-tag-spellings
adrs: []
sources:
  - 2026-09-08-plan-releases-with-bare-stable-tags.md
regeneration: []
promoted: []
delivery_commit: 7806acfa9f4fbe5236be80af5c2388965376bf46
---

# Release planning preserves stable tag identity

Spec 0147 was authored as the whole content of this Spec rather than as a slice of it: both tag spellings parse, selection spans them, an ambiguous highest version refuses in preflight, the proposal keeps the selected spelling, the inventory keeps refs apart, and the public states and exit codes are unchanged.
