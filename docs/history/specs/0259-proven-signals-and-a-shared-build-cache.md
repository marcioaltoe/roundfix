---
schema: roundfix/archive-record/v1
spec: 0259-proven-signals-and-a-shared-build-cache
title: Proven signals and a shared build cache
status: archived
created: "2026-10-12"
archived: "2026-10-13"
disposition: qa-override
source: docs/specs/0259-proven-signals-and-a-shared-build-cache
source_revision: 8f2370fb01db575a71d5b46453703d025bf673ab
qa_task: task_05
qa_report: qa-report-2026-10-13-01.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: 'Maintainer, AskUserQuestion 2026-10-13 11:15Z: Publicar e investigar depois'
qa_override_reason: 'Every requirement row passed except the ten-suite acceptance of F6: six full suites were clean and suite 7 saw TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout report migrate --check exited -1 (a signal; sender unknown). R09 (Surface Transcript 3 aggregate exit) was blocked by the offline fixture and R17 by the pre-PR Pull Request row. F6 stays open for the next cycle: name the signal first, then trace its sender.'
qa_override_qa_outcome: fail
qa_override_qa_task_status: failed
qa_override_revision: 8f2370fb01db575a71d5b46453703d025bf673ab
adrs:
  - ADR-0266
sources:
  - 2026-10-09-a-test-child-killed-by-an-unknown-signal.md
  - 2026-10-09-each-worktree-keeps-its-own-go-build-cache.md
regeneration:
  - command: make skills-sync
promoted: []
---

# Proven signals and a shared build cache

The report is fail and does not qualify for ordinary QA settlement/archive while R09 remains environment-blocked without equivalent aggregate readiness proof and R10 fails. The no-open-PR row alone would not prevent settlement. No audited-checkout Task status, graph, authorization, commit, branch, push or Pull Request was changed by QA; fixture initialization stayed in disposable repositories. The Daemon owns report integration and terminal Task settlement.
