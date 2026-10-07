---
schema: roundfix/archive-record/v1
spec: 0005-tui-cockpit
title: TUI Cockpit Redesign
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0005-tui-cockpit
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0008
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# TUI Cockpit Redesign

The Live Run View grew organically: the review path got the full cockpit, the Implement Command received a functional but minimal Task pane, and issue detail occupies layout instead of appearing on demand. A terminal UI redesign was already drafted and mocked before the Implement Command existed (see `design/ui-redesign-plan.md` and the two reference mockups in `design/`); this Spec adopts that direction and generalizes it: one two-surface cockpit for both Run kinds, with Work Items on the left, the session timeline as the dominant surface on the right, a phase row stating where the Run is, and detail as a centered terminal modal.
