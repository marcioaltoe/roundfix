---
schema: roundfix/archive-record/v1
spec: 0048-context-driven-project-decisions-and-spec-constraints
title: Context-Driven project decisions and Spec constraints
status: archived
created: "2026-07-24"
archived: "2026-07-25"
disposition: pass
source: docs/history/specs/0048-context-driven-project-decisions-and-spec-constraints
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-25.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "35"
delivery_commit: 8122ceeab1cc8fa6d46616801489e375e1610e19
---

# Context-Driven project decisions and Spec constraints

The composed Context-Driven guidance still cannot reproduce the accepted project contract unless maintainers manually add an Internal Identifier strategy, a Better Auth route exception, and a broad tooling-authority rule. New Specs can then omit those constraints and authorize incompatible work by silence. This feature collects the project decisions through the public Baseline Command, persists them in the Setup Manifest, renders self-contained guidance, and makes every new Spec account for the applicable Project Constraints. The [greenfield acceptance finding](../../findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md) is the acceptance source.
