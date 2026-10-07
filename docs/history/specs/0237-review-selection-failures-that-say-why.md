---
schema: roundfix/archive-record/v1
spec: 0237-review-selection-failures-that-say-why
title: Review selection failures that say why
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: pass
source: docs/history/specs/0237-review-selection-failures-that-say-why
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0242
sources:
  - 2026-10-06-a-review-agent-selection-failure-says-nothing.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "427"
delivery_commit: bdfaa69aade1c9ef076358a037bc49a9c4530780
---

# Review selection failures that say why

On 2026-10-02 in Pantheon (Roundfix 0.26), `roundfix review --base main` with the default `codex / gpt-5.6-luna / max` reviewer blocked twice with `review runtime failure: Agent Selection failed for runtime "codex": agent/protocol error` and an empty answer file, while the Doctor Command and `codex exec` answered. Three days later the same review ran. The reason named neither the protocol step that failed nor the adapter's message, and the maintainer skipped the review for a Pull Request ([the adopted Backlog Entry](references/2026-10-06-a-review-agent-selection-failure-says-nothing.md)).
