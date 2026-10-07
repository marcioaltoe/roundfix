---
schema: roundfix/archive-record/v1
spec: 0021-cockpit-visual-fidelity
title: Cockpit Visual Fidelity
status: archived
created: "2026-07-07"
archived: "2026-07-07"
disposition: pass
source: docs/history/specs/0021-cockpit-visual-fidelity
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-07.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
---

# Cockpit Visual Fidelity

Spec 0005 shipped the cockpit's information architecture — two panes, Phase Row, Work Queue, timeline grouping, Detail Modal — and its QA passed against that written contract. The visual layer of the approved design (`docs/specs/_archived/0005-tui-cockpit/design/roundfix-01.png` and `roundfix-02.png`, with the style rules in `design/ui-redesign-plan.md`) never became verifiable acceptance criteria, so it silently degraded: the shipped cockpit is near-monochrome, the timeline renders raw Agent payloads instead of summary rows, selection is a bare marker instead of a highlighted card, and empty states explain nothing. Cockpit Visual Fidelity closes that gap: the cockpit renders to the approved design, and this time the styled render is the acceptance criterion.
