---
schema: roundfix/archive-record/v1
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
title: A queue item shows its turn, and the pre-PR review gets its own capacity
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
source_revision: 9fd289d103f33f9430346af9e9447d58a5c303a7
qa_task: task_05
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven:
  - Fiscus's `deliver status` after the release that carries this Spec, reported back through ai-memory
adrs:
  - ADR-0365
sources:
  - 2026-11-05-fiscus-delivery-queue-takes-days.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "761"
type: feat
---

# A queue item shows its turn, and the pre-PR review gets its own capacity

The offline gate supports a qualifying partial for configured review capacity, reviews beside implementation, preserved start order and all six public Turn forms; Fiscus's own queue remains unproven under the PRD declaration, and the absent Pull Request has recorded pre-PR equivalents. The Run-built archive plan exited 0 with the judge skipped and only the adopted Fiscus report as a candidate; it contains external incident history already preserved by the archive cut and adds no reusable knowledge needing promotion. No archive, release or delivery operation was performed, and the Daemon retains settlement ownership.
