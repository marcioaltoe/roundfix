---
schema: roundfix/archive-record/v1
spec: 0044-upgrade-retention-and-formatter-compatibility
title: Context-Driven Baseline upgrade retention and formatter compatibility
status: archived
created: "2026-07-22"
archived: "2026-07-22"
disposition: pass
source: docs/history/specs/0044-upgrade-retention-and-formatter-compatibility
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-22.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0058
  - ADR-0059
  - ADR-0047
  - ADR-0046
sources: []
regeneration: []
promoted: []
---

# Context-Driven Baseline upgrade retention and formatter compatibility

A real 0.9.0 upgrade of a managed repository proved that a Context-Driven Baseline version transition can replace previously managed hard rules wholesale, compress supporting guides below their operational contract, drop normative clauses during prose-to-rule migration, render duplicate skill dispatch, and produce output the target repository's formatter immediately rewrites — all while the current catalog passes its own asset contract. Two 2026-07-22 findings record the evidence: the upgrade preserved bytes outside managed markers but had no boundary for the normative strength of what a previous setup version had placed inside them. This Spec adds that boundary while preserving the modular architecture Spec 0043 delivered.
