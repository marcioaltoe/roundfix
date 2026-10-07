---
schema: roundfix/archive-record/v1
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
title: A Jev Router that checks its credit and names its model
status: archived
created: "2026-10-04"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0229-a-jev-router-that-checks-credit-and-names-its-model
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-05.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: row 7b needs historical raw run records outside the authorized scope (the 2026-10-04 measurement addenda record them), and row 10 has no Pull Request yet; the queue opens it. Every behavior row passed against local upstreams. The operator runs one live routed replay after merge to confirm the recorded model.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 932cd2c5987206523d30c7c34b87fe0f7efcc028
adrs:
  - ADR-0234
  - ADR-0114
sources:
  - 2026-10-04-jev-router-cost-controls.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "396"
delivery_commit: 7919b3ffbaef022c14899948ce820cec3dfa0196
---

# A Jev Router that checks its credit and names its model

On 2026-10-04 a routed Task ended at OpenRouter instead of in Roundfix. The Jev Router gate read the month's Jev spend (US$9.69 of the US$50 ceiling) and the key's remaining limit (US$40.31) and sent a Verification repair prompt. The OpenRouter account behind the key held about US$5.11, so OpenRouter refused the request with HTTP 402 ("This request would exceed your available credits given your current in-flight requests").
