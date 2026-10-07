---
schema: roundfix/archive-record/v1
spec: 0235-one-qa-partial-policy
title: One QA partial policy
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: qa-override
source: docs/history/specs/0235-one-qa-partial-policy
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-06-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after the corrective docs task_06: O3 is an outside-evidence JUnit lookup the Spec''s authorization forbids, and R12 has no Pull Request yet; the queue opens it. Every behavior row passed, F1 and F2 are fixed by task_06.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 6efb825354b51d04a5f1fb49afd17b3dcdb168bd
adrs:
  - ADR-0240
sources:
  - 2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "421"
delivery_commit: 29c31dbe93dcde22e2efafb99e864477ddfa0f30
---

# One QA partial policy

A QA gate runs before any Pull Request exists, and often inside a Run sandbox without network access. Its report then closes `partial` on rows nobody in the Run can reach: the Pull Request row and outside-evidence rows whose source the sandbox denied. The Daemon settles such a QA Task `failed`, `roundfix settle` refuses it, and the Delivery Queue parks the item `qa-environment-partial`. The operator then archives with `--qa-override` and retries. That happened on nearly every delivery from 2026-10-04 to 2026-10-06 (`~/.roundfix-operator/queue-interventions.md`, entries 140 to 186), and two adopters reported the same refusal ([references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md)).
