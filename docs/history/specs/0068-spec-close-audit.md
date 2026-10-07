---
schema: roundfix/archive-record/v1
spec: 0068-spec-close-audit
title: Spec close audit
status: archived
created: "2026-08-02"
archived: "2026-08-04"
disposition: pass
source: docs/history/specs/0068-spec-close-audit
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_07
qa_report: qa-report-2026-08-04-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "114"
delivery_commit: 1583e641058166b96b32cdeed3609b6d0a8bfb06
---

# Spec close audit

The per-Spec loop ends at "squash merge and reconcile" and never audits what the cycle created against what survived it. In one session that left four kinds of residue at once: two Supervisor scratch worktrees never removed after their push, a Run Worktree orphaned because the squash merge deleted the target branch `reconcile` resolves by name, a remote backup branch whose purpose ended when its work merged, and two Pull Requests opened and left unreviewed.
