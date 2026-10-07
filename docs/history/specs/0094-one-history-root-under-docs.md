---
schema: roundfix/archive-record/v1
spec: 0094-one-history-root-under-docs
title: One history root under `docs/`
status: archived
created: "2026-08-12"
archived: "2026-08-13"
disposition: pass
source: docs/history/specs/0094-one-history-root-under-docs
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_09
qa_report: qa-report-2026-08-13.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0120
  - ADR-0123
sources:
  - 2026-08-12-the-archive-root-sits-beside-docs-instead-of-inside-it.md
regeneration: []
promoted: []
pull_request: "159"
delivery_commit: 5e2d6ba54bcb6df5c5cd6c55c26d1f8cd5abe4ed
---

# One history root under `docs/`

Retired documentation is stored in three places and none of them is right. In this repository it sits at the repository root, beside the Go source trees. In every other repository that adopted Roundfix it sits inside the active directories themselves, one archive folder per documentation tree, so history and live work share a parent. Retired ADRs never moved at all, so an Agent loading decision context reads inactive decisions alongside active ones. And Review Artifacts have no terminal home: a Pull Request whose Spec is known already writes its reviews inside that Spec and travels with it, but the orphan case — fifty folders here — accumulates beside the active Spec Root forever, with nothing deciding when one is finished. This Spec gives every retired family one home under `docs/`, teaches Roundfix to move a repository there on its own, and settles when an orphan review is history rather than work in flight.
