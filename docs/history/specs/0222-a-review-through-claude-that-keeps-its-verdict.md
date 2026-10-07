---
schema: roundfix/archive-record/v1
spec: 0222-a-review-through-claude-that-keeps-its-verdict
title: A review through Claude that keeps its verdict
status: archived
created: "2026-10-04"
archived: "2026-10-04"
disposition: qa-override
source: docs/history/specs/0222-a-review-through-claude-that-keeps-its-verdict
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-04.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: rows 6b1-6b3 cite authoring-clone artifacts that no longer exist (ADR-0227 records the measurements); row 6e was cut at a 50-second QA limit, and the operator then ran one real claude/opus review of this candidate on 2026-10-04 (base ff135017): end_turn, outcome findings, finding id F1, estimatedPromptTokens 55912; row 9 has no Pull Request yet, and the queue opens it.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: d65ee0052916d614e1e2264cef8cb20011c50f9c
adrs:
  - ADR-0227
sources:
  - 2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "375"
delivery_commit: 9a053bee70179e5478faab82596dfcd7be9f2bbb
---

# A review through Claude that keeps its verdict

On 2026-10-02 two repositories that select `pre_pr_review.provider: claude` with the `review` profile `claude / opus / high` could not complete the pre-PR review the repository policy requires ([the adopted Backlog Entry](references/2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md)). `roundfix review` blocked as a `review transport anomaly` because acpx exited `5` after the turn had already delivered its verdict: the reviewer had asked to run a terminal command and the read-only session refused it. The block hid a clean `No findings.
