---
schema: roundfix/archive-record/v1
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
title: A Task Context path the Task deletes is not unresolved
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0348-a-task-context-path-the-task-deletes-is-not-unresolved
source_revision: 4ed26d9b100133b0f0734bc3959b5bc6f78be535
qa_task: task_03
qa_report: qa-report-2026-11-04.md
qa_verdict: partial
unproven: []
adrs: []
sources:
  - 2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# A Task Context path the Task deletes is not unresolved

The Spec resolves a missing Task Context path through HEAD or the Delivery Base while continuing to reject paths neither revision tracks, supplies the removal/deletes fix, and documents deletes-only authoring in both skill copies and the guides. Live CLI deletion and typo journeys, Daemon settlement tests, all required consumer suites and all seven isolated sabotage controls confirm that behavior. QA closes as a qualifying partial solely because no Pull Request exists; review and Merge-Ready acceptance remain for delivery.
