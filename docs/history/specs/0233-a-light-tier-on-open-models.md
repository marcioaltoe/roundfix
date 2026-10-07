---
schema: roundfix/archive-record/v1
spec: 0233-a-light-tier-on-open-models
title: A light tier on open models
status: archived
created: "2026-10-05"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0233-a-light-tier-on-open-models
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-05.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: the Run sandbox denied the OpenCode docs, the GitHub issue and the OpenRouter key endpoint (R11-R13), and R16 has no Pull Request yet; the queue opens it. Every behavior row passed. The live light-tier check happens on Spec 0234''s own delivery, whose low Tasks run on the light tier by default; the operator records its cost.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 0a9172b6342f30b2c1ca0d0238a25b648687cea5
adrs:
  - ADR-0238
sources:
  - 2026-10-05-a-judge-assigned-model-tier-per-task.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "412"
delivery_commit: d470b9e75bd5f7933d44b41d5d969a45ba7e2261
---

# A light tier on open models

Every Task in a Run runs on its Agent Work Category's Preferred Selection, so a one-line documentation Task spends the same codex or claude subscription quota as a new subsystem. Since ADR-0235 OpenAI and Anthropic models run only through those subscriptions, which leaves open models on OpenRouter as the cheap place for small work. The adopted Backlog Entry, [a judge-assigned model tier per Task](references/2026-10-05-a-judge-assigned-model-tier-per-task.md), proposed a tier per Task and measured it on 2026-10-05: six archived `complexity: low` Tasks replayed on OpenCode with `deepseek/deepseek-v4.1-flash` all passed Verification on the first attempt, in 0.41 times Codex's agent time, for US$0.
