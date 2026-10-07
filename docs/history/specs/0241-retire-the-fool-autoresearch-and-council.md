---
schema: roundfix/archive-record/v1
spec: 0241-retire-the-fool-autoresearch-and-council
title: Retire the-fool, autoresearch and council
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: qa-override
source: docs/history/specs/0241-retire-the-fool-autoresearch-and-council
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-06-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after the transcript fix: R4 needs github.com (network denied in the Run sandbox, but its provenance is a requirement row, not an outside-evidence row), R8b is an outside-evidence lookup the Spec''s authorization prohibits, and R11 has no Pull Request yet; the queue opens it. Every behavior row passed, including Surface Transcript 1.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 816411615c1791a1df88fcf155de165d3cc2e0aa
adrs:
  - ADR-0246
sources:
  - 2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
promoted: []
pull_request: "438"
delivery_commit: df324cf68eb1ee77f24e21b2195c778cdb9308ef
---

# Retire the-fool, autoresearch and council

On 2026-10-06 the maintainer asked to retire the skills that no longer have a use, or whose job the Jev judge took over. Asked per skill, the maintainer answered: "remover: the-fool, autoresearch, council // manter: grilling, grill-with-docs, write-idea, business-analyst e handoff". The adopted Backlog Entry, [references/2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md](references/2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md), measured the use: across 220 Runs, `the-fool` and `autoresearch` were never read, and `council` was read in 3 Runs.
