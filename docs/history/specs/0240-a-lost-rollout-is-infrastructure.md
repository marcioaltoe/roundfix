---
schema: roundfix/archive-record/v1
spec: 0240-a-lost-rollout-is-infrastructure
title: A lost rollout is infrastructure
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: partial
source: docs/history/specs/0240-a-lost-rollout-is-infrastructure
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-06.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0245
sources:
  - 2026-10-06-a-lost-codex-rollout-fails-the-task.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "434"
delivery_commit: c838272f32361debcbfc88ce25a3d2457fb1ea48
---

# A lost rollout is infrastructure

Between 2026-10-05 and 2026-10-06, Runs in Oraculum's Delivery Queue (Roundfix 0.44 to 0.46, `codex-acp` 2.1.1) failed mid-turn with `internal error -32603: no rollout found for thread id <uuid>`. Codex had created each thread and never written its rollout file, so the adapter could not resume it. Roundfix read the failure as the Task's. The Task settled `failed`, a QA gate died before its report, and each attempt spent one of the queue's retries. The same QA, re-run by hand on `claude / opus / high`, reached a verdict ([the adopted Backlog Entry](references/2026-10-06-a-lost-codex-rollout-fails-the-task.md)).
