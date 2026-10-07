---
schema: roundfix/archive-record/v1
spec: 0227-reconcile-releases-the-runs-of-merged-specs
title: Reconcile releases the Runs of merged Specs
status: archived
created: "2026-10-04"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0227-reconcile-releases-the-runs-of-merged-specs
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-05.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: row 07 needs historical Run Database items of the hand-cleaned Specs, which no longer exist (intervention log entries 144-145 record them), and row 10 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 03454c09d9f23c0ec5c69c1d981267fc673a1be7
adrs:
  - ADR-0232
  - ADR-0212
sources:
  - 2026-10-04-reconcile-keeps-runs-of-merged-specs.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "394"
delivery_commit: 96e733771b33a206f95be2e95a3a66512dd820a9
---

# Reconcile releases the Runs of merged Specs

On 2026-10-04 the maintainer found twelve `roundfix/run-…` and item branches in the main checkout. `roundfix reconcile --apply` released only the four it classified `superseded`. It kept six Runs of Specs 0181, 0184 (two), 0200, 0204 and 0207, every one archived and merged on main, as `unintegrated` or `dirty`, and it never reported the item branches of 0205 and 0217 that the Delivery Queue had recreated. The operator removed them by hand (`~/.roundfix-operator/queue-interventions.md`, entry 144).
