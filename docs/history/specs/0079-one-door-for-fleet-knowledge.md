---
schema: roundfix/archive-record/v1
spec: 0079-one-door-for-fleet-knowledge
title: One door for fleet knowledge
status: archived
created: "2026-08-06"
archived: "2026-08-06"
disposition: partial
source: docs/history/specs/0079-one-door-for-fleet-knowledge
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_07
qa_report: qa-report-2026-08-06-04.md
qa_verdict: partial
unproven:
  - a fleet inventory taken after the window, showing zero pending entries older than fourteen days
  - the pilot's recorded brain commits, their one-path change sets, and the destination checkout's clean status before triage, all read fresh from their commit objects by the gate
  - the triaged entry carrying `resolved_to` in the brain and the typed Backlog Entry citing its inbox provenance in this repository
  - the pending research digest with its cited sources and the recorded advisory qmd result from the pilot
  - 'a supervised session close that commits a pending `capture: auto` draft, followed by an independent later read proving it was not self-triaged'
adrs:
  - ADR-0095
sources:
  - 2026-08-06-findings-accumulate-faster-than-they-become-specs.md
regeneration: []
promoted: []
---

# One door for fleet knowledge

Knowledge the fleet produces dies in three ways today: observations live only in conversations until a commit that Active Runs block, cross-project feedback travels by hand, and findings pile up — 63 in this repository, 111 across eight projects — with no lifecycle stamped, no consolidation, and no archival. This Spec gives every observation one durable door, a triage that lands it in the contract-true home, a rollup that keeps the active set bounded, and a path for acquired research knowledge to reach the Secondbrain's curated layer. The full exploration, scoring, council dissent, and chosen hybrid direction are recorded in [the idea](./_idea.md); the measured evidence is the adopted investigation in [references](./references/2026-08-06-findings-accumulate-faster-than-they-become-specs.md).
