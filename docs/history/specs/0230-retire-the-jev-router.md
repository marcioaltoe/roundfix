---
schema: roundfix/archive-record/v1
spec: 0230-retire-the-jev-router
title: Retire the Jev Router
status: archived
created: "2026-10-05"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0230-retire-the-jev-router
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-05.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: Q7a needs the operator''s OpenRouter activity export (docs/_inbox/openrouter_activity_2026-10-05.csv in the main checkout; the measurement addendum records its figures), Q7c and Q7d would need OpenRouter requests that this Spec forbids, and Q10 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 73587c4d7ea9ccf5531168868093219ba53b3f82
adrs:
  - ADR-0235
  - ADR-0027
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "402"
delivery_commit: f2022a597923d23c91cf17fbc028f3acc3f4c110
---

# Retire the Jev Router

The Jev Router let a repository hand a Task to OpenRouter's `typesafe/jev-router`, which picks a model for each request. OpenRouter's activity export for Roundfix's key from 2026-10-01 to 2026-10-05 shows that the routed sessions were billed, pay-per-use, for `openai/gpt-6-astra` (US$15.35), `anthropic/claude-opus-5.5` (US$10.12) and `openai/gpt-6.1-sol` (US$4.81), models this repository already reaches through the codex and claude subscriptions, while 456 direct judge calls cost US$0.014. On 2026-10-05 the maintainer set the rule "os modelos da openai e anthropic devem ser utilizado exclusivamente pela assinatura que o codex e claude fornecem", decided "Aposentar o router", and, after weighing a restricted router, confirmed "Vamos seguir com a remoção".
