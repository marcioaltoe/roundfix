---
schema: roundfix/archive-record/v1
spec: 0059-run-storage-compaction-and-global-sanitation
title: Run storage compaction and global sanitation
status: archived
created: "2026-07-28"
archived: "2026-08-06"
disposition: pass
source: docs/history/specs/0059-run-storage-compaction-and-global-sanitation
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_06
qa_report: qa-report-2026-08-06-03.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "131"
delivery_commit: 7f6e7d19329a28447a01a9277bf0fc0c78e438d4
---

# Run storage compaction and global sanitation

The GC Command prunes Run Event Journal rows machine-wide but never returns the deleted pages to the filesystem — the Run Database sat at 1.4 GB with its journal already pruned — and its Artifact Directory cleanup resolves only the current repository's Artifact Root, so a machine-wide database coexists with repository-scoped artifact discovery and orphaned roots from other repositories are never reclaimed. Tables added after the retention contract (Run summaries, Agent Selection records) have no defined long-term lifecycle at all. Evidence: [sanitation is repository-scoped and does not compact SQLite](../../findings/2026-07-17-global-run-storage-sanitation-and-compaction.md). The terminal Run Worktree half of that report shipped with Spec 0038; this Spec owns the database and artifact halves.
