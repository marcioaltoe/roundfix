---
schema: roundfix/archive-record/v1
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
title: Baseline plans that show what a History Relocation breaks
status: archived
created: "2026-09-29"
archived: "2026-09-29"
disposition: qa-override
source: docs/history/specs/0184-baseline-plans-that-show-what-a-history-move-breaks
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_04
qa_report: qa-report-2026-09-29.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'Maintainer standing authorization of 2026-08-09 (agent project memory ''qa_override por ambiente''): archive with override when the gate closes partial only because of environment, never with rows_blocked_finding above zero. Reconfirmed for this cycle on 2026-09-29 (''Autonomia ampla'').'
qa_override_reason: QA closed partial with rows_blocked_environment 2 and rows_blocked_finding 0. Row 4, the Fluxus outside-evidence replay, cannot run because the exported Fluxus Setup Manifest has no maintained transition to the go-cli-tui profile. Row 9 is the pre-PR Pull Request row, which records equivalent evidence. Every other row passed.
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 2423efe2af34390f6575e649e90dbcd94500091e
adrs:
  - ADR-0173
sources:
  - 2026-09-25-baseline-history-moves-need-separate-approval.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "284"
delivery_commit: 404cf0c8699b0c04d10a2342e6c00601329a3848
---

# Baseline plans that show what a History Relocation breaks

A Baseline Plan that brings an adopted repository to the current layout carries its History Relocations in the same approval as the managed refresh. The plan names each moved file, but not the citations the move breaks. On 2026-09-17, Fiscus's `baseline update` proposed a one-entry managed refresh together with six ADRs moving from the decision directory into the History Root.
