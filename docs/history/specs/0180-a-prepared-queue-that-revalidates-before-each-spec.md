---
schema: roundfix/archive-record/v1
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
title: A prepared queue that revalidates before each Spec
status: archived
created: "2026-09-28"
archived: "2026-09-29"
disposition: pass
source: docs/history/specs/0180-a-prepared-queue-that-revalidates-before-each-spec
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_06
qa_report: qa-report-2026-09-29.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-09-28-queue-preparation-inventory.md
  - 2026-09-28-revalidate-later-specs-when-prerequisites-change.md
  - 2026-09-28-queue-wide-limits-and-one-pending-question.md
  - 2026-09-28-implement-spec-skill-delegates-to-roundfix.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "273"
delivery_commit: 658e4073683d3c1ae8a1aa03d8e09343d8ef3a04
---

# A prepared queue that revalidates before each Spec

`roundfix deliver start <slug>...` records whatever slugs it is given and starts working. Nothing tells the operator beforehand which of those Specs is actually approved to run, and nothing re-checks a later Spec against the main it will start from once earlier items have merged. The queue has no limit of its own, and it presents every parked item at once. The owned `implement-spec` skill still tells the Supervisor to run its own Task loop. Retired Spec 0127 carried these gaps as Core Features 1, 4, 5, 6 and 8. They are now four Backlog Entries, and the efficiency waves hit each of them:
