---
schema: roundfix/archive-record/v1
spec: 0225-a-jev-ceiling-the-maintainer-sets
title: A Jev ceiling the maintainer sets
status: archived
created: "2026-10-04"
archived: "2026-10-04"
disposition: qa-override
source: docs/history/specs/0225-a-jev-ceiling-the-maintainer-sets
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_03
qa_report: qa-report-2026-10-04.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: row 09 needs the operator''s key-status record and installation; the operator read GET /api/v1/key on 2026-10-04 (limit 50, limit_reset monthly) and the maintainer''s decision is quoted in _authorization.md. Row 12 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 94761ac03fe623575143ccf90703580620e1999c
adrs:
  - ADR-0231
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "382"
delivery_commit: 14c4d8534f70117f3aac101a7dbf3c0f0cb072b2
---

# A Jev ceiling the maintainer sets

Roundfix stops spending on Jev once the month's spend reaches a ceiling that is compiled into the binary at US$5. The judge skips at that point, and a routed Jev Router prompt is refused, together with any OpenRouter key whose monthly credit limit is above US$5. On 2026-10-04 the maintainer raised the limit of Roundfix's OpenRouter key and asked for the work to run without that bound: "A chave do openrouter está com limite de $50 mensal e quero que seja tudo implementado e testado sem limitações. Depois vemos o custo e como reduzir o custo." Today only a release can change the number, and the key's new limit makes every routed prompt fail as an unbounded key.
