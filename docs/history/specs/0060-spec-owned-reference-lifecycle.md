---
schema: roundfix/archive-record/v1
spec: 0060-spec-owned-reference-lifecycle
title: Spec-owned reference lifecycle
status: archived
created: "2026-07-28"
archived: "2026-07-31"
disposition: pass
source: docs/history/specs/0060-spec-owned-reference-lifecycle
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-31.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0083
sources: []
regeneration: []
promoted: []
pull_request: "57"
delivery_commit: fe36ac6bfd32c91d3b8dd2ebae88f9dd4f6b64e9
---

# Spec-owned reference lifecycle

Promotion from inbox to finding to Spec currently creates links without transferring the source document, so a Spec can be complete and archived while the evidence that shaped it stays behind in a separate lifecycle with links that no longer resolve — archived Specs are linked from findings at their pre-archive paths, and triaged inbox notes outlive the Specs they routed. The 2026-07-28 triage session did all of this reconciliation by hand: removing stale inbox notes, flipping finding statuses, and repointing archived-Spec links, exactly the manual work a lifecycle contract should own. Evidence: [adopted source documents must travel with their owning Spec](../../findings/2026-07-25-spec-owned-reference-lifecycle.md).
