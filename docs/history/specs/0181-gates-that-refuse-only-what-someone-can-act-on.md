---
schema: roundfix/archive-record/v1
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
title: Gates that refuse only what someone can act on
status: archived
created: "2026-09-29"
archived: "2026-09-29"
disposition: pass
source: docs/history/specs/0181-gates-that-refuse-only-what-someone-can-act-on
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-29-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0166
  - ADR-0167
  - ADR-0088
  - ADR-0161
  - ADR-0168
sources:
  - 2026-09-29-a-file-a-task-creates-fails-the-qa-scope-audit.md
  - 2026-09-29-the-pull-request-row-blocks-every-qualifying-partial.md
  - 2026-09-29-a-new-adr-forces-edits-to-every-active-spec.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "281"
delivery_commit: 122098e7d6d6b332a590c18737091d30176290f3
---

# Gates that refuse only what someone can act on

Waves 2 and 3 shipped on 2026-09-28 and 2026-09-29. Three of their gate refusals had nothing to do with the work being wrong. Each cost reruns or hand edits:
