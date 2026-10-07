---
schema: roundfix/archive-record/v1
spec: 0209-sources-that-share-a-context-share-a-spec
title: Sources that share a context share a Spec
status: archived
created: "2026-10-01"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0209-sources-that-share-a-context-share-a-spec
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0208
  - ADR-0209
sources: []
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text
promoted: []
pull_request: "331"
delivery_commit: d4a0a0daf9e4b79070721227031340206d5a5cec
---

# Sources that share a context share a Spec

Nothing in the Baseline tells an author to look for related work before minting a Spec or a record. One Spec per Finding looks like the safe default, and so does a new Backlog Entry for an intent an open entry already holds. The result is more Specs than the work needs, each paying its own authoring, review and QA, and records that say the same thing twice. In the fluxus repository a Spec was authored without consulting the backlog and left its own Backlog Entry open after it shipped, which a later Spec had to clean up.
