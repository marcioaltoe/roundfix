---
schema: roundfix/archive-record/v1
spec: 0088-a-third-runtime-that-can-run
title: A third ACP Runtime that can actually run
status: archived
created: "2026-08-08"
archived: "2026-08-08"
disposition: pass
source: docs/history/specs/0088-a-third-runtime-that-can-run
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_08
qa_report: qa-report-2026-08-08.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0106
  - ADR-0105
  - ADR-0107
sources:
  - 2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md
regeneration: []
promoted: []
---

# A third ACP Runtime that can actually run

`CONTEXT.md` lists OpenCode as a supported ACP Runtime through `opencode acp`. Measured on 2026-08-08, no Roundfix Run can start on it. Three independent defects stand in the way, and the cost is concrete: the Codex quota was exhausted until 2026-08-12 and the Anthropic weekly limit stood at 66% with two days to reset, leaving a freshly subscribed runtime with full capacity and no route to it.
