---
schema: roundfix/archive-record/v1
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
title: A Run that reports the tokens and spend it used
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: qa-override
source: docs/history/specs/0204-a-run-that-reports-the-tokens-and-spend-it-used
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: ""
qa_verdict: ""
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing approval of 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'QA partial with environment rows only: rows 04, 06, 07, 09 and 10 need the built binary under a disposable Home with a seeded Run Database, which the QA Agent''s instructions forbid; their in-process equivalents pass exact text, schema, totals and read-only checks (binary-equivalents.log), and the operator ran the unseeded journeys with the built binary (qa/evidence/2026-10-01-operator/binary-journeys.txt). Row 15 is the pre-PR Pull Request row.'
qa_override_qa_outcome: missing
qa_override_qa_task_status: pending
qa_override_revision: 9000a443d784f279493e5946052213f6b53f0f56
adrs:
  - ADR-0198
  - ADR-0199
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "312"
delivery_commit: 58e6d57fd45aa8b46b8bbd3eb57eb1f37790783b
---

# A Run that reports the tokens and spend it used

Roundfix cannot say what a Run consumed. `roundfix deliver status` ends its limits line with `spend not measured`, and a Delivery Queue can be bounded only by a deadline and a retry allowance. The questions a maintainer asks after a wave cannot be answered from the Run Database: how many tokens a Spec took, which Task took most of them, and whether a queue should stop before it drains a quota window.
