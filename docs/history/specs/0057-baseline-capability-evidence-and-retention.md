---
schema: roundfix/archive-record/v1
spec: 0057-baseline-capability-evidence-and-retention
title: Baseline capability evidence and retention accounting
status: archived
created: "2026-07-28"
archived: "2026-08-02"
disposition: qa-override
source: docs/history/specs/0057-baseline-capability-evidence-and-retention
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-08-02-01.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0058
sources: []
regeneration: []
promoted: []
pull_request: "75"
delivery_commit: ea9c05fe2a06878338ff0a023d02d1dd606e2178
---

# Baseline capability evidence and retention accounting

Two live Baseline sessions — a profile refresh in fluxus and a greenfield alignment in vortex — exposed the same defect family: the Baseline Command's divergences are individually correct but do not carry enough evidence to be acted on, and a changed Profile under an unchanged Baseline identity bypasses retention accounting entirely, letting managed Normative Clauses disappear with an empty Upgrade Retention Contract ledger. Capability probes reject executable symlinks (the norm on Homebrew and Docker Desktop machines), divergence output hides the probe it evaluated, there is no read-only way to re-check capabilities after remediation, and a clean adoption warns about the thirteen files it just wrote. Evidence: [profile refresh applied without semantic retention accounting](../../findings/2026-07-26-baseline-profile-refresh-retention-gap.md) and [capability divergences do not carry enough evidence to be remediated](../../findings/2026-07-26-vortex-baseline-capability-remediation.md). Specs 0049–0051 fixed Preservation idempotency and Doctor hardening; none of this report's items were covered there.
