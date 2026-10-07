---
schema: roundfix/archive-record/v1
spec: 0192-owned-skills-that-describe-the-product-as-it-is
title: Owned skills that describe the product as it is
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: qa-override
source: docs/history/specs/0192-owned-skills-that-describe-the-product-as-it-is
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-30.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing approval of 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'QA partial with environment rows only. Row 10: the QA sandbox proxy refused GitHub (HTTP 403); the operator fetched git/git@e9019fca t/t0450-txt-doc-vs-help.sh (180 lines, sha256 4fdbfd74...) and recorded it under qa/evidence/2026-09-30-operator/. Row 13: no Pull Request exists before the PR. Eleven rows pass; pre-PR review at 46533f6b reported no findings.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 46533f6bff8bf480425006a5198f800245dff677
adrs: []
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "296"
delivery_commit: 93bd6b2a7f9f71c88f178658da8c01e80f2a8cf5
---

# Owned skills that describe the product as it is

An audit on 2026-09-30 compared the fourteen Roundfix-owned skills and the user guide with the v0.21.0 binary and the Baseline clauses. Twelve skills and five guide files state something the product no longer does, or omit something it does:
