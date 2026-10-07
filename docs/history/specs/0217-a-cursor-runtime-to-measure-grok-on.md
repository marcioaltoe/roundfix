---
schema: roundfix/archive-record/v1
spec: 0217-a-cursor-runtime-to-measure-grok-on
title: A Cursor runtime to measure Grok on
status: archived
created: "2026-10-01"
archived: "2026-10-03"
disposition: pass
source: docs/history/specs/0217-a-cursor-runtime-to-measure-grok-on
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-03.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0217
  - ADR-0208
sources:
  - 2026-09-30-run-grok-through-the-cursor-runtime.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "358"
delivery_commit: 2300080633d2709432e0afb9d66f7dd137a68d81
---

# A Cursor runtime to measure Grok on

Roundfix selects an Agent for each work category from three ACP Runtimes: `codex`, `claude` and `opencode`. The maintainer also has Cursor, whose plan includes Grok models, and its agent CLI, `cursor-agent`, is installed on the maintainer's machine. No Agent Selection can name it, so Grok cannot be tried on a real Task, and a fourth provider that fails independently of the other three cannot join a Fallback Chain.
